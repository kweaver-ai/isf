// Package sms 逻辑层
package sms

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/oklog/ulid/v2"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
	"Authentication/logics/conf"
)

var (
	aSMSOnce sync.Once
	aSMS     *anonymousSMS
)

// 定义错误码唯一标识
const (
	_ int = iota
	i18nSMSVCodeExpired
	i18nSMSVerificationFailed
	i18nSMSInvalidPhoneNumber
)

type anonymousSMS struct {
	dbAnonymousSMS interfaces.DBAnonymousSMS
	dbTracePool    *sqlx.DB
	sharemgnt      interfaces.DnShareMgnt
	conf           interfaces.Conf
	numerics       [10]byte
	smsWidth       int           // 匿名账户短信验证码位数
	smsExpiration  time.Duration // 匿名账户短信验证码过期时间
	re             *regexp.Regexp
	privateKey     *rsa.PrivateKey
	logger         common.Logger
	trace          observable.Tracer
	i18n           *common.I18n
	smsExpRWMutex  sync.RWMutex // 匿名账户短信验证码过期时间的读写锁
}

// NewAnonymousSMS 创建anonymousSMS对象
func NewAnonymousSMS() *anonymousSMS {
	aSMSOnce.Do(func() {
		// pem 解码
		blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
		// X509解码
		privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		// 使用正则表达式进行匹配
		pattern := `^1[3456789]\d{9}$`
		re := regexp.MustCompile(pattern)

		aSMS = &anonymousSMS{
			dbAnonymousSMS: logics.DBAnonymousSMS,
			dbTracePool:    logics.DBTracePool,
			sharemgnt:      logics.DnShareMgnt,
			conf:           conf.NewConf(),
			numerics:       [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			smsWidth:       6,
			re:             re,
			privateKey:     privateKey,
			logger:         common.NewLogger(),
			trace:          common.SvcARTrace,
			i18n: common.NewI18n(common.I18nMap{
				i18nSMSVCodeExpired: {
					interfaces.SimplifiedChinese:  "验证码已过期",
					interfaces.TraditionalChinese: "驗證碼已過期",
					interfaces.AmericanEnglish:    "The verification code has expired",
				},
				i18nSMSVerificationFailed: {
					interfaces.SimplifiedChinese:  "验证码校验失败",
					interfaces.TraditionalChinese: "驗證碼校驗失敗",
					interfaces.AmericanEnglish:    "The CAPTCHA verification failed. Please try again",
				},
				i18nSMSInvalidPhoneNumber: {
					interfaces.SimplifiedChinese:  "手机号不合法",
					interfaces.TraditionalChinese: "手機號不合法",
					interfaces.AmericanEnglish:    "Invalid tel number",
				},
			}),
		}

		// 每天凌晨清理一次
		go aSMS.cronDelete()

		// 初始化短信验证码过期时间
		err = aSMS.initSMSExpiration()
		if err != nil {
			aSMS.logger.Fatalln(err)
		}
	})

	return aSMS
}

func (aSMS *anonymousSMS) create(ctx context.Context, phoneNumber, anonymityID string) (*interfaces.AnonymousSMSInfo, error) {
	aSMS.smsExpRWMutex.RLock()
	smsExpiration := aSMS.smsExpiration
	aSMS.smsExpRWMutex.RUnlock()

	// 在创建新的验证码之前，先清除之前申请的、依旧有效的验证码。
	err := aSMS.dbAnonymousSMS.DeleteRecordWithinValidityPeriod(ctx, phoneNumber, anonymityID, smsExpiration)
	if err != nil {
		return nil, err
	}

	aSMSInfo := &interfaces.AnonymousSMSInfo{
		ID:          ulid.Make().String(),
		PhoneNumber: phoneNumber,
		AnonymityID: anonymityID,
		Content:     aSMS.genValidateCode(aSMS.smsWidth),
		CreateTime:  time.Now().Unix(),
	}
	err = aSMS.dbAnonymousSMS.Create(ctx, aSMSInfo)
	if err != nil {
		return nil, err
	}

	return aSMSInfo, nil
}

func (aSMS *anonymousSMS) Validate(ctx context.Context, visitor *interfaces.Visitor, vcodeID, vcode, anonymityID string) (phoneNumber string, err error) {
	aSMS.trace.SetInternalSpanName("逻辑层-匿名验证码校验")
	newCtx, span := aSMS.trace.AddInternalTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	aSMSInfo, err := aSMS.dbAnonymousSMS.GetInfoByID(newCtx, vcodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSVerificationFailed, visitor.Language))
		}
		return "", err
	}

	aSMS.smsExpRWMutex.RLock()
	defer aSMS.smsExpRWMutex.RUnlock()

	// 判断记录是否失效
	if time.Now().After(time.Unix(aSMSInfo.CreateTime, 0).Add(aSMS.smsExpiration)) {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSVCodeExpired, visitor.Language))
	}
	if !strings.EqualFold(vcode, aSMSInfo.Content) {
		return "", rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSVerificationFailed, visitor.Language))
	}

	// 删除验证码的逻辑移到 login.Anonymous2 中
	return aSMSInfo.PhoneNumber, nil
}

func (aSMS *anonymousSMS) CreateAndSendVCode(ctx context.Context, visitor *interfaces.Visitor, phoneNumber, anonymityID string) (vcodeID string, err error) {
	aSMS.trace.SetInternalSpanName("逻辑层-创建并发送匿名验证码")
	newCtx, span := aSMS.trace.AddInternalTrace(ctx)
	defer func() { aSMS.trace.TelemetrySpanEnd(span, err) }()

	// 接口向后兼容：如果为11位手机号明文，代表旧的客户端发送的请求，不需要解密
	if !aSMS.re.MatchString(phoneNumber) {
		phoneNumber, err = logics.DecodeAndDecrypt(phoneNumber, aSMS.privateKey)
		if err != nil {
			return "", err
		}

		// 需要校验解密后的手机号，是否合法
		if !aSMS.re.MatchString(phoneNumber) {
			return "", rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSInvalidPhoneNumber, visitor.Language))
		}
	}

	aSMSInfo, err := aSMS.create(newCtx, phoneNumber, anonymityID)
	if err != nil {
		return "", err
	}

	err = aSMS.sharemgnt.UsrmSendAnonymousSMSVCode(newCtx, aSMSInfo.PhoneNumber, aSMSInfo.Content)
	if err != nil {
		return "", err
	}

	return aSMSInfo.ID, nil
}

func (aSMS *anonymousSMS) cronDelete() {
	const t time.Duration = 24
	for {
		// 计算下一个零点
		next := common.Now().Add(time.Hour * t)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Local().Sub(common.Now()))
		<-t.C

		aSMS.smsExpRWMutex.RLock()
		smsExpiration := aSMS.smsExpiration
		aSMS.smsExpRWMutex.RUnlock()

		for {
			ctx := context.Background()
			ids, err := aSMS.dbAnonymousSMS.GetExpiredRecords(ctx, smsExpiration)
			if err != nil {
				break
			}

			if len(ids) > 0 {
				err = aSMS.dbAnonymousSMS.DeleteByIDs(ctx, ids)
				if err != nil {
					break
				}
			} else {
				break
			}
		}
	}
}

func (aSMS *anonymousSMS) genValidateCode(width int) string {
	r := len(aSMS.numerics)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", aSMS.numerics[rand.Intn(r)])
	}

	return sb.String()
}

// UpdateSMSExpiration 更新短信验证码过期时间
func (aSMS *anonymousSMS) UpdateSMSExpiration(expiration int) {
	aSMS.smsExpRWMutex.Lock()
	defer aSMS.smsExpRWMutex.Unlock()

	aSMS.smsExpiration = time.Minute * time.Duration(expiration)
	aSMS.logger.Infof("Update anonymous login SMS verification code expiration time to %d minutes", expiration)
}

func (aSMS *anonymousSMS) initSMSExpiration() error {
	cfgKeys := map[interfaces.ConfigKey]bool{
		interfaces.SMSExpiration: true,
	}

	visitor := interfaces.Visitor{
		ErrorCodeType: interfaces.Number,
	}
	ctx := context.Background()
	cfg, err := aSMS.conf.GetConfig(ctx, &visitor, cfgKeys)
	if err != nil {
		return err
	}
	aSMS.smsExpiration = time.Minute * time.Duration(cfg.SMSExpiration)
	return nil
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
