// Package assertion 逻辑层
package assertion

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/satori/uuid"
	"go.step.sm/crypto/jose"
	"gopkg.in/square/go-jose.v2/jwt"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
	accesstokenperm "Authentication/logics/access_token_perm"
)

var (
	aOnce sync.Once
	a     *assertion
)

type assertion struct {
	userMgnt        interfaces.DnUserManagement
	hydraAdmin      interfaces.DnHydraAdmin
	hydraPublic     interfaces.DnHydraPublic
	accessTokenPerm interfaces.AccessTokenPerm
	log             common.Logger
	issuer          string
	kid             string
	trace           observable.Tracer
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
}

// NewAssertion 创建assertion操作对象
func NewAssertion() *assertion {
	aOnce.Do(func() {
		// pem 解码
		blockPub, _ := pem.Decode([]byte(pubKeyStr))
		// X509解码
		pubKey, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		publicKey := pubKey.(*rsa.PublicKey)

		// pem 解码
		blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
		// X509解码
		privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		a = &assertion{
			userMgnt:        logics.DnUserManagement,
			hydraAdmin:      logics.DnHydraAdmin,
			hydraPublic:     logics.DnHydraPublic,
			accessTokenPerm: accesstokenperm.NewAccessTokenPerm(),
			log:             common.NewLogger(),
			issuer:          "authentication",
			kid:             "401e6962-ab50-4c90-bf07-f9d23d4def3b",
			trace:           common.SvcARTrace,
			privateKey:      privateKey,
			publicKey:       publicKey,
		}
	})

	return a
}

// GetAssertionByUserID 获取指定用户的访问令牌断言
func (a *assertion) GetAssertionByUserID(ctx context.Context, visitor *interfaces.Visitor, userID string) (assertion string, err error) {
	a.trace.SetInternalSpanName("逻辑层-获取指定用户的访问令牌断言")
	newCtx, span := a.trace.AddInternalTrace(ctx)
	defer func() { a.trace.TelemetrySpanEnd(span, err) }()

	// (1) 权限校验
	// 访问者类型检查
	if visitor.Type != interfaces.Business {
		err = rest.NewHTTPError("Visitor type must be app.", rest.BadRequest, nil)
		return
	}

	// 检查userID是否存在
	_, err = a.userMgnt.GetUserRolesByUserID(newCtx, visitor, userID)
	if err != nil {
		return
	}

	// 检查应用账户获取任意用户访问令牌权限
	hasPerm, err := a.accessTokenPerm.CheckAppAccessTokenPerm(visitor.ID)
	if err != nil {
		return
	}
	if !hasPerm {
		err = rest.NewHTTPError("This app doesn't have permission to access token.", rest.Forbidden, nil)
		return
	}

	// (2) 生成断言
	privateClaims := map[string]interface{}{
		"ext": map[string]interface{}{
			"visitor_type": "realname",
			"login_ip":     "",
			"account_type": "other",
			"udid":         "",
			"client_type":  "app",
			"client_id":    visitor.ID, // 限制生成断言与使用断言的app必须一致
		},
	}
	return a.CreateJWK(context.Background(), userID, time.Hour, privateClaims)
}

func (a *assertion) CreateJWK(ctx context.Context, subject string, ttl time.Duration, privateClaims map[string]interface{}) (assertion string, err error) {
	a.trace.SetInternalSpanName("逻辑层-生成断言")
	_, span := a.trace.AddInternalTrace(ctx)
	defer func() { a.trace.TelemetrySpanEnd(span, err) }()

	var trustedPair map[string]bool
	trustedPair, err = a.hydraAdmin.GetKidTrustedPairByIssuer(a.issuer)
	if err != nil {
		return
	}
	if !trustedPair[a.kid] {
		err = a.hydraAdmin.CreateTrustedPair(a.publicKey, a.issuer, a.kid)
		if err != nil {
			return
		}
	}

	signingKey := jose.SigningKey{Algorithm: jose.SignatureAlgorithm("RS256"), Key: a.privateKey}
	signer, _ := jose.NewSigner(signingKey, (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", a.kid))

	audience, err := a.hydraPublic.GetTokenEndpoint()
	if err != nil {
		return
	}

	jti := uuid.NewV4()
	claims := jwt.Claims{
		Issuer:    a.issuer,
		Subject:   subject,
		Audience:  []string{audience},
		Expiry:    jwt.NewNumericDate(time.Now().Add(ttl)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        jti.String(),
	}

	assertion, err = jwt.Signed(signer).Claims(claims).Claims(privateClaims).CompactSerialize()
	if err != nil {
		a.log.Errorln(err)
	}
	return
}

func (a *assertion) TokenHook(assertion, clientID string) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	segs := strings.SplitN(assertion, ".", -1)
	claims, err := a.decodeSegment(segs[1])
	if err != nil {
		a.log.Errorln(err)
		return
	}
	result = claims["ext"].(map[string]interface{})
	// 限制生成断言与使用断言的app必须一致
	if !(result["client_id"].(string) == clientID) {
		// 根据hydra webhook要求，返回403标识拒绝此次请求
		err = rest.NewHTTPError("Invalid client_id.", rest.Forbidden, nil)
		return
	}
	delete(result, "client_id")
	return
}

func (a *assertion) decodeSegment(seg string) (map[string]any, error) {
	const remainder = 4
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", remainder-l)
	}
	bi, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	err = json.Unmarshal(bi, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

const pvtKeyStr = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const pubKeyStr = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArxTpPffUO5x5c0Wx/xS8
ocYcKgCAwfD0Uo/3W3oFjwbO42S/UGzP1YNzIKZP3cXRV+8Ps74+2NSNqljN85gr
yfmZkz/gOILpd3mNemYbYpdnStgin0pc4FA62IN+lV/d3Sa/xds7qk9yY8ut85Kw
RxZnmr+gnppAI0E1sZKKOGfkClz1N5J0o65Qdu90JdZ0lU/B1Uu29MLCi5YRHQs8
Pzm8sU5W7+g9EHxKbM995zWpSC65gUKoxNSHo+ED9leDTzPDEF2SEUIpSJwx7Fnw
kB3RuD9UYvT+FOsn2edIijj7t3VrImGeqpfQnlrVewje17SNf77HSau6y3zJfwvZ
3g0PoxwPAbjXFITrMFxOd73r1mFNwNYxQXSQyGNi+Z+pfNa6C19pL7ehsjRw7DeC
qxqsJ9rgxjgGWNnKkBMsUTxck22kYuwcECLYkBikx2WqwTwb648zjDbzGKBQmyuM
AVDbTbyRDXONV5OGRJ3vrC7xCzoksaTc+Tnms+OKRQzJ/n1Njv+7sulHf4qrhtt5
DYc5u8arfxwJOoEwt++nSK4IiyhWz1w2pYxae67kRfr2C5njaiPah9UWlv/e1flb
/YnZPdkqSwdj/v8kmIPMzFAjeOO2FsljNfLJ4rWeWmu0k4jkqyR3C/jLEd/DyIlk
sQfkDZGF3IFOKPCINxzY7EcCAwEAAQ==
-----END PUBLIC KEY-----`
