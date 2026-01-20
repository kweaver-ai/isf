// Package ticket 逻辑层
package ticket

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/oklog/ulid/v2"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

type ticket struct {
	hydraAdmin       interfaces.DnHydraAdmin
	userMgnt         interfaces.DnUserManagement
	dbTicket         interfaces.DBTicket
	ticketExpiration time.Duration
	privateKey       *rsa.PrivateKey
	logger           common.Logger
	trace            observable.Tracer
}

var (
	tOnce sync.Once
	t     *ticket
)

// NewTicket 创建新的ticket对象
func NewTicket() *ticket {
	tOnce.Do(func() {
		// pem 解码
		blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
		// X509解码
		privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		t = &ticket{
			hydraAdmin:       logics.DnHydraAdmin,
			userMgnt:         logics.DnUserManagement,
			dbTicket:         logics.DBTicket,
			ticketExpiration: time.Minute * 5,
			privateKey:       privateKey,
			logger:           common.NewLogger(),
			trace:            common.SvcARTrace,
		}

		go t.cronDelete()
	})
	return t
}

func (tic *ticket) CreateTicket(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.TicketReq) (ticketID string, err error) {
	tic.trace.SetInternalSpanName("逻辑层-生成单点登录凭据")
	newCtx, span := tic.trace.AddInternalTrace(ctx)
	defer func() { tic.trace.TelemetrySpanEnd(span, err) }()

	reqInfo.RefreshToken, err = logics.DecodeAndDecrypt(reqInfo.RefreshToken, tic.privateKey)
	if err != nil {
		return "", err
	}

	// 刷新令牌内省，获取userID
	refreshTokeninfo, err := tic.hydraAdmin.IntrospectRefreshToken(reqInfo.RefreshToken)
	if err != nil {
		tic.logger.Errorf("IntrospectRefreshToken failed, err: %v", err)
		return "", err
	}
	// 不能自己给自己申请单点登录凭据。这是没有意义的
	if refreshTokeninfo.ClientID == reqInfo.ClientID {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, "The client cannot apply for single sign-on credentials for itself")
	}

	// 判断用户是否存在
	_, err = tic.userMgnt.GetUserInfo(newCtx, visitor, refreshTokeninfo.Sub)
	if err != nil {
		tic.logger.Errorf("GetUserInfo failed, err: %v", err)
		return "", err
	}

	// 判断reqInfo.ClientID代表的客户端是否存在
	_, err = tic.hydraAdmin.GetClientInfo(reqInfo.ClientID)
	if err != nil {
		tic.logger.Errorf("GetClientInfo failed, err: %v", err)
		return "", err
	}

	ticketInfo := &interfaces.TicketInfo{
		ID:         ulid.Make().String(),
		UserID:     refreshTokeninfo.Sub,
		ClientID:   reqInfo.ClientID,
		CreateTime: time.Now().Unix(),
	}

	if err = tic.dbTicket.Create(newCtx, ticketInfo); err != nil {
		tic.logger.Errorf("Create ticket failed, err: %v", err)
		return "", err
	}

	return ticketInfo.ID, nil
}

func (tic *ticket) Validate(ctx context.Context, ticketID, clientID string) (userID string, err error) {
	tic.trace.SetInternalSpanName("逻辑层-验证单点登录凭据")
	newCtx, span := tic.trace.AddInternalTrace(ctx)
	defer func() { tic.trace.TelemetrySpanEnd(span, err) }()

	info, err := tic.dbTicket.GetTicketByID(newCtx, ticketID)
	if err != nil {
		tic.logger.Errorf("GetTicketByID failed, err: %v", err)
		return "", err
	}
	if info.ClientID != clientID {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, "Invalid clientID")
	}
	if time.Since(time.Unix(info.CreateTime, 0)) > tic.ticketExpiration {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, "Ticket expired")
	}

	err = tic.dbTicket.DeleteByIDs(newCtx, []string{ticketID})

	return info.UserID, err
}

func (tic *ticket) cronDelete() {
	const duration time.Duration = 24
	for {
		// 计算下一个零点
		next := common.Now().Add(time.Hour * duration)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Local().Sub(common.Now()))
		<-t.C

		for {
			ctx := context.Background()
			ids, err := tic.dbTicket.GetExpiredRecords(ctx, tic.ticketExpiration)
			if err != nil {
				break
			}

			if len(ids) > 0 {
				err = tic.dbTicket.DeleteByIDs(ctx, ids)
				if err != nil {
					break
				}
			} else {
				break
			}
		}
	}
}

const pvtKeyStr = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsyOstgbYuubBi2PUqeVjGKlkwVUY6w1Y8d4k116dI2SkZI8f
xcjHALv77kItO4jYLVplk9gO4HAtsisnNE2owlYIqdmyEPMwupaeFFFcg751oiTX
JiYbtX7ABzU5KQYPjRSEjMq6i5qu/mL67XTkhvKwrC83zme66qaKApmKupDODPb0
RRkutK/zHfd1zL7sciBQ6psnNadh8pE24w8O2XVy1v2bgSNkGHABgncR7seyIg81
JQ3c/Axxd6GsTztjLnlvGAlmT1TphE84mi99fUaGD2A1u1qdIuNc+XuisFeNcUW6
fct0+x97eS2eEGRr/7qxWmO/P20sFVzXc2bF1QIDAQABAoIBAACDungGYoJ87bLl
DUQUqtl0CRxODoWEUwxUz0XIGYrzu84nJBf5GOs9Xv6i9YbNgJN2xkJrtTU7VUJF
AfaSP4kZXqqAO9T1Id9zVc5oomuldSiLUwviwaMek1Yh9sFRqWNGGxBdd7Y1ckm8
Roy+kHZ7xXqlIxOmdCC+7DgQMVgSV64wzQY8p7L9kTLIkeDodEolkUkGsreF9I9S
kzlLjGU9flPt13319G0KSaQUWEpxF/UBr2gKJvQPQHSRzzl5HlRwznZkU4Hs6RID
ue6E68ZJNMRn3FUAvLMCRw9C4PQQR/x/50WH4BXJ9veVIOIpTVCJedI0QZjbVuBk
RPKHTMkCgYEA2XjGIw9Vp0qu/nCeo5Nk15xt/SJCn0jIhyRpckHtCidotkiZmFdU
vUK7IwbAUPqEJcgmS/zwREV8Gff8S324C2RoDN4FxFtBMZgQjqV1zYqGLQSbTJUh
GlpTe7jKVskuSPSf00OqqAIlYNtzZK3mWj8MadFD99Wo9gktXRAFdf0CgYEA0uBe
wuE007XLqb8ANS+4U0CkexeVDkDzI9yXN2CB+L5wmJ/WsNF8iD53xHxpwZWRiizX
ArBdhWL9yv4YkbryyD15NRSQhLanRcs0MqGh1GJJ9vpGzBjfJJ3Bw0hBfkwnf/C6
nNzGjNWNTeNKwlcFaVhBADyGYZt9Len9YYFNKrkCgYEAmsn7BYNprOxciCAy2i0U
Lt9Z7j3Pe757dK13HGtOQ9bvEie0o5ktaJSxzGmGw1y8aIQAtj9v6Lgob/dxrW3r
bLhn0xjItA1b5ufciRu+MLFzdWF9BFJ1QGOgXkSWSJVji2wKwn28X18/qaQpizS3
6+5KcJsRrLp4S78WedHogSUCgYEAomb5k8wtCv7vIoNefZeKtVMLWWEIAjozBmNU
cel5L0A7Js+yX+p1pde2FTRbniK6O1fdHs0EuT1Lh5G5CkKXx27QcfisdAjXOgEM
6hFguFgZ7oNBEt30vBZiqypyhfnQUc/rZ/L/VmcAtANgB9tM55x4Mt5p/7Hn7fxO
j1EtRMECgYEAp2sI035BcCR2kFW1vC9eXLAPZ0anyy1/T1dEgFJ/ELqmGEMEWZKA
9H1KH6YIkDdXabwfaSTRebaEescCxRtgmo5WEdZxw4Nz66SSomc24aD0iem7+VSl
x2qRWdif0jHG8fOdMey3NrY7NF4xQTzuO9jDnLpBTwFg3o7QlywIBlM=
-----END RSA PRIVATE KEY-----`
