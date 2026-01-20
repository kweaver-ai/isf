package login

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
	"Authentication/logics"
	"Authentication/tapi/ethriftexception"
)

// nolintlint
func newLogin(ha interfaces.DnHydraAdmin, hp interfaces.DnHydraPublic, um interfaces.DnUserManagement,
	eacp interfaces.DnEacp, shareMgnt interfaces.DnShareMgnt, loginDB interfaces.DBLogin,
	config interfaces.Conf, assertion interfaces.Assertion, aSMS interfaces.LogicsAnonymousSMS) *login {
	authFailedMap := make(map[string]int)
	authFailedMap["invalid_password"] = common.InvalidAccountORPassword
	authFailedMap["initial_password"] = common.PasswordISInitial
	authFailedMap["password_not_safe"] = common.PasswordNotSafe
	authFailedMap["under_control_password_expire"] = common.ControledPasswordExpire
	authFailedMap["password_expire"] = common.PasswordExpire

	// pem 解码
	blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
	// X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
	if err != nil {
		common.NewLogger().Fatalln(err)
	}

	return &login{
		hydraAdmin:     ha,
		hydraPublic:    hp,
		userManagement: um,
		eacp:           eacp,
		sharemgnt:      shareMgnt,
		loginDB:        loginDB,
		config:         config,
		authFailedMap:  authFailedMap,
		assertion:      assertion,
		aSMS:           aSMS,
		privateKey:     privateKey,
		logger:         common.NewLogger(),
		i18n: common.NewI18n(common.I18nMap{
			i18nSMSVCodeInfoRequired: {
				interfaces.SimplifiedChinese:  "请输入验证码",
				interfaces.TraditionalChinese: "請輸入驗證碼訊息",
				interfaces.AmericanEnglish:    "Please enter verification code information",
			},
		}),
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

//nolint:lll
func TestSingleSignOn(t *testing.T) {
	Convey("singlesignon", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		um := mock.NewMockDnUserManagement(ctrl)
		eacp := mock.NewMockDnEacp(ctrl)
		loTicket := mock.NewMockLogicsTicket(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		n := newLogin(hydraAdmin, hydraPublic, um, eacp, nil, nil, nil, nil, nil)
		n.trace = trace
		n.ticket = loTicket

		loginInfo := &interfaces.SSOLoginInfo{
			ClientID:     "00002da3-b64f-4d61-9269-48fc77966ec8",
			RedirectURI:  "https://127.0.0.1:9010/callback",
			ResponseType: "token",
			Scope:        "offline",
			Udids:        []string{"127.0.0.1"},
			Credential: interfaces.SSOCredential{
				ID:     "test",
				Params: "XXXXXXXX",
			},
		}
		ctx := context.Background()
		visitor := &interfaces.Visitor{Language: interfaces.SimplifiedChinese}

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("authorize request failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("", nil, testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("get login request information failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(nil, testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("third party authentication failed", func() {
			device := &interfaces.DeviceInfo{}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("third party authentication failed2", func() {
			loginInfo.Credential.ID = "aishu"
			loginInfo.Credential.Params = map[string]interface{}{
				"ticket": "Lpw51Y/KYx2UefjxdNtezRHtIEStHRAbooPexyexXDAz6tYj/5+gHL2/jqup91+kxHTW+GsG0vFxkkJyPVM5jBUOqDiFzXn3OUEhYXoeS6sp2g4g9df5lYrZtyO+SPiBHYXMECQoAoCVfiXqFM+4LU3OPUqHi3UMteAojpoJlL5ep1u+7Vk4Sbq7h0ZamZUPBkiqWiIWfC1aDBa94s/lXOy5xFstgbC9JToFIxA8QUP25puYYjM7OdOC/4aNAL9wnoW7qPb7DybN7R8GlOPRqEyCxWzCDF3fb1lelqowkv7SOEZ1veojMJ9oqxUvXJ+crhNBw2Nu3U4nW9WXnmYRkA==",
			}
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.UserBaseInfo{
				Account: "account",
			}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)

			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			loTicket.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("user_id", nil)
			um.EXPECT().GetUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)

			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("accept login failed", func() {
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.LoginInfo{}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("verifier login failed", func() {
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.LoginInfo{}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil, testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("accept consent request failed", func() {
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.LoginInfo{}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("verifier consent failed", func() {
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.LoginInfo{}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierConsent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("singlesignon success", func() {
			device := &interfaces.DeviceInfo{}
			userInfo := &interfaces.LoginInfo{}
			token := &interfaces.TokenInfo{}
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			hydraAdmin.EXPECT().GetLoginRequestInformation(gomock.Any()).AnyTimes().Return(device, nil)
			eacp.EXPECT().ThirdPartyAuthentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierConsent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(token, nil)
			tokenInfo, err := n.SingleSignOn(ctx, visitor, loginInfo)
			assert.NotEqual(t, tokenInfo, nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAnonymous(t *testing.T) {
	Convey("anonumous", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		um := mock.NewMockDnUserManagement(ctrl)
		eacp := mock.NewMockDnEacp(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		n := newLogin(hydraAdmin, hydraPublic, um, eacp, nil, nil, nil, nil, nil)
		n.trace = trace

		loginInfo := &interfaces.AnonymousLoginInfo{
			ClientID:     "00002da3-b64f-4d61-9269-48fc77966ec8",
			RedirectURI:  "https://127.0.0.1:9010/callback",
			ResponseType: "token",
			Scope:        "offline",
			Credential: interfaces.AnonymousCredential{
				Account:  "test",
				Password: "XXXXXXXX",
			},
		}

		visitor := &interfaces.Visitor{
			Language:      interfaces.SimplifiedChinese,
			ErrorCodeType: interfaces.Number,
		}
		Convey("authorize request failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("", nil, testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("Anonymous authentication failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("accept login failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("verifier login failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil, testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("accept consent request failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("verifier consent failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierConsent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.Equal(t, tokenInfo, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("anonymous success", func() {
			token := &interfaces.TokenInfo{}
			hydraPublic.EXPECT().AuthorizeRequest(gomock.Any()).AnyTimes().Return("login_challenge", nil, nil)
			um.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			hydraAdmin.EXPECT().AcceptLoginRequest(gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierLogin(gomock.Any(), gomock.Any()).AnyTimes().Return("consent_challenge", nil, nil)
			hydraAdmin.EXPECT().AcceptConsentRequest(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("redirURL", nil)
			hydraPublic.EXPECT().VerifierConsent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(token, nil)
			tokenInfo, err := n.Anonymous(visitor, loginInfo)
			assert.NotEqual(t, tokenInfo, nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAnonymous2(t *testing.T) {
	Convey("anonumous2", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mock.NewMockTraceClient(ctrl)

		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		assertion := mock.NewMockAssertion(ctrl)
		aSMS := mock.NewMockLogicsAnonymousSMS(ctrl)
		dbAnonymousSMS := mock.NewMockDBAnonymousSMS(ctrl)
		n := newLogin(nil, hydraPublic, userMgnt, nil, nil, nil, nil, assertion, aSMS)
		n.trace = trace
		n.dbAnonymousSMS = dbAnonymousSMS

		ctx := context.Background()
		loginInfo := &interfaces.AnonymousLoginInfo2{
			ClientID:     "clientID",
			ClientSecret: "clientSecret",
			Credential: interfaces.AnonymousCredential{
				Account:  "account",
				Password: "password",
			},
		}
		visitor := &interfaces.Visitor{Language: interfaces.SimplifiedChinese}

		Convey("GetAnonymityInfoByID error", func() {
			tmpErr := errors.New("unknow error")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, tmpErr)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, tmpErr)
		})

		Convey("needVerifyMobile is true, vcode_id is null, content is not null", func() {
			loginInfo.VCode.Content = "xx"
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, n.i18n.Load(i18nSMSVCodeNotSend, interfaces.SimplifiedChinese)))
		})

		Convey("needVerifyMobile is true, vcode_id null, content is null", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, n.i18n.Load(i18nSMSVCodeInfoRequired, interfaces.SimplifiedChinese)))
		})

		Convey("visitor_name is required", func() {
			loginInfo.VCode.ID = "id"
			loginInfo.VCode.Content = "content"
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "visitor_name required"))
		})

		Convey("needVerifyMobile is true, smsValidate failed", func() {
			loginInfo.VCode.ID = "id"
			loginInfo.VCode.Content = "content"
			loginInfo.VisitorName = "xx"
			tmpErr := fmt.Errorf("unknown")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			aSMS.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("", tmpErr)
			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, tmpErr)
		})

		Convey("访问密码不正确", func() {
			testErr := rest.NewHTTPError("Wrong password", 403007001, nil)
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			userMgnt.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.NotEqual(t, err, nil)
		})

		Convey("数据不存在", func() {
			testErr := rest.NewHTTPError("Not Found", 400007001, nil)
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			userMgnt.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.NotEqual(t, err, nil)
		})

		Convey("访问次数已达上限", func() {
			testErr := rest.NewHTTPError("Not Found", 403007002, nil)
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			userMgnt.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)

			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.NotEqual(t, err, nil)
		})

		Convey("验证码删除失败", func() {
			resInfo := &interfaces.TokenInfo{
				AccessToken: "AccessTokenStr",
				TokenType:   "bearer",
				Scope:       "all",
				ExpirsesIn:  3600,
			}
			loginInfo.VCode.ID = "id"
			loginInfo.VisitorName = "xx"
			loginInfo.VCode.Content = "content"
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			aSMS.EXPECT().Validate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("phoneNumber", nil)
			userMgnt.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			assertion.EXPECT().CreateJWK(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("assertionStr", nil)
			hydraPublic.EXPECT().AssertionForToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(resInfo, nil)
			dbAnonymousSMS.EXPECT().DeleteByIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("error"))
			_, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, errors.New("error"))
		})

		Convey("断言换取令牌成功", func() {
			resInfo := &interfaces.TokenInfo{
				AccessToken: "AccessTokenStr",
				TokenType:   "bearer",
				Scope:       "all",
				ExpirsesIn:  3600,
			}
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			userMgnt.EXPECT().GetAnonymityInfoByID(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			userMgnt.EXPECT().AnonymousAuthentication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)
			assertion.EXPECT().CreateJWK(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("assertionStr", nil)
			hydraPublic.EXPECT().AssertionForToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(resInfo, nil)
			dbAnonymousSMS.EXPECT().DeleteByIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			res, err := n.Anonymous2(ctx, visitor, loginInfo, "")

			assert.Equal(t, err, nil)
			assert.Equal(t, res.AccessToken, "AccessTokenStr")
			assert.Equal(t, res.TokenType, "bearer")
			assert.Equal(t, res.Scope, "all")
			assert.Equal(t, res.ExpirsesIn, int64(3600))
		})
	})
}

//nolint:funlen
func TestClientAccountAuth(t *testing.T) {
	Convey("ClientAccountAuth", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userMgnt := mock.NewMockDnUserManagement(ctrl)
		shareMgnt := mock.NewMockDnShareMgnt(ctrl)
		loginDB := mock.NewMockDBLogin(ctrl)
		config := mock.NewMockConf(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		login := newLogin(nil, nil, userMgnt, nil, shareMgnt, loginDB, config, nil, nil)
		login.trace = trace

		reqInfo := &interfaces.ClientLoginReq{
			Method:   "GET",
			Account:  "account1",
			Password: "111111",
			Option: interfaces.ClientLoginOption{
				VCodeType: interfaces.ImageVCode,
				UUID:      "uuid1",
				VCode:     "code1",
			},
		}
		cfg := interfaces.Config{
			EnableIDCardLogin:  true,
			EnablePWDLock:      true,
			EnableThirdPWDLock: true,
			PWDErrCnt:          3,
			PWDLockTime:        10,
			VCodeConfig: interfaces.VCodeLoginConfig{
				Enable:    true,
				PWDErrCnt: 3,
			},
		}

		detail := make(map[string]interface{})
		detail["isShowStatus"] = false
		ctx := context.Background()
		visitor := interfaces.Visitor{
			Language: interfaces.AmericanEnglish,
		}

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("用户不存在，且需要短信验证码，报错用户或者密码错误", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthSMS
			detail["isShowStatus"] = true
			userInfo := interfaces.UserBaseInfo{}
			tmpErr := rest.NewHTTPError("", common.ImageVCodeISWrong, detail)
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, userInfo, nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, tmpErr)
		})

		Convey("用户不存在，且需要动态密码，报错用户或者密码错误", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			detail["isShowStatus"] = true
			userInfo := interfaces.UserBaseInfo{}
			tmpErr := rest.NewHTTPError("", common.OTPWrong, detail)
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, userInfo, nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, tmpErr)
		})

		Convey("三元分立下，用户为system，报错用户或者密码错误", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			detail["isShowStatus"] = true
			userInfo := interfaces.UserBaseInfo{
				ID: logics.SystemAdminID,
			}
			tempCfg := cfg
			tempCfg.TriSystemStatus = true
			tmpErr := rest.NewHTTPError("", common.OTPWrong, detail)
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempCfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, tmpErr)
		})

		Convey("validate image vcode, timeout", func() {
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      1,
				PwdErrLastTime: time.Now().Unix(),
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			tmpErr := &ethriftexception.NcTException{
				ErrID: int32(common.ImageVCodeTimeout % 1e6),
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmIMAGECodeValidate(gomock.Any(), gomock.Any()).Return(tmpErr)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError("", int(rest.Unauthorized+int32(common.ImageVCodeTimeout%1e6)), detail))
		})
		Convey("validate sms vcode, failed", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthSMS
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      1,
				PwdErrLastTime: time.Now().Unix(),
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			tmpErr := &ethriftexception.NcTException{
				ErrID: int32(common.ImageVCodeISWrong % 1e6),
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmSMSValidate(gomock.Any(), gomock.Any()).Return(tmpErr)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError("", int(rest.Unauthorized+int32(common.ImageVCodeISWrong%1e6)), detail))
		})
		Convey("local auth failed", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = false
			cfg.EnableThirdPWDLock = false
			cfg.VCodeConfig.Enable = false
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      1,
				PwdErrLastTime: time.Now().Unix(),
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			tmpErr := rest.NewHTTPError("", login.authFailedMap["invalid_password"], detail)
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, "invalid_password", nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, tmpErr)
		})
		Convey("domain auth failed", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = false
			cfg.EnableThirdPWDLock = false
			cfg.VCodeConfig.Enable = false
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Domain,
				PwdErrCnt:      1,
				PwdErrLastTime: time.Now().Unix(),
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			tmpErr := &ethriftexception.NcTException{
				ErrID: int32(common.DomainNotExist % 1e6),
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			shareMgnt.EXPECT().UsrmDomainAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, tmpErr)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError("", int(rest.Unauthorized+int32(common.DomainNotExist%1e6)), detail))
		})
		Convey("third auth failed", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = false
			cfg.EnableThirdPWDLock = false
			cfg.VCodeConfig.Enable = false
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Third,
				PwdErrCnt:      1,
				PwdErrLastTime: time.Now().Unix(),
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			tmpErr := errors.New("Unknown error")
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			shareMgnt.EXPECT().UsrmThirdAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, tmpErr)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError(tmpErr.Error(), rest.InternalServerError, detail))
		})
		Convey("auth failed, account locked", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = true
			cfg.EnableThirdPWDLock = false
			cfg.VCodeConfig.Enable = true
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      2,
				PwdErrLastTime: time.Now().Unix() - 1,
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			detail["isShowStatus"] = true
			tmpErr := rest.NewHTTPError("", common.InvalidAccountORPassword, detail)
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(false, "invalid_password", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, tmpErr)
		})
		Convey("auth success", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = true
			cfg.EnableThirdPWDLock = true
			cfg.VCodeConfig.Enable = true
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      2,
				PwdErrLastTime: time.Now().Unix() - 1,
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, "", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "id1")
			assert.Equal(t, err, nil)
		})

		Convey("user disaled", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = true
			cfg.EnableThirdPWDLock = true
			cfg.VCodeConfig.Enable = true
			detail["isShowStatus"] = false
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      2,
				PwdErrLastTime: time.Now().Unix() - 1,
				DisableStatus:  true,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, "", nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError("", common.UserDisabled, detail))
		})

		Convey("user locked", func() {
			reqInfo.Option.VCodeType = interfaces.DualAuthOTP
			cfg.EnablePWDLock = true
			cfg.EnableThirdPWDLock = true
			cfg.VCodeConfig.Enable = false
			userInfo := interfaces.UserBaseInfo{
				ID:             "id1",
				AuthType:       interfaces.Local,
				PwdErrCnt:      3,
				PwdErrLastTime: time.Now().Unix() - 1,
				DisableStatus:  false,
				LDAPType:       interfaces.LDAPServerType(0),
				DomainPath:     "",
			}
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			loginDB.EXPECT().GetDomainStatus().AnyTimes().Return(true, nil)
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, userInfo, nil)
			shareMgnt.EXPECT().UsrmOTPValidate(gomock.Any(), gomock.Any()).Return(nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(true, "", nil)

			userID, err := login.ClientAccountAuth(ctx, &visitor, reqInfo)

			detail["isShowStatus"] = false
			detail["remainlockTime"] = int64(10)
			assert.Equal(t, userID, "")
			assert.Equal(t, err, rest.NewHTTPError("", common.AccountLocked, detail))
		})
	})
}

//nolint:nolintlint, lll
func TestPwdAuth(t *testing.T) {
	Convey("ClientAccountAuth", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userMgnt := mock.NewMockDnUserManagement(ctrl)
		shareMgnt := mock.NewMockDnShareMgnt(ctrl)
		loginDB := mock.NewMockDBLogin(ctrl)
		config := mock.NewMockConf(ctrl)
		assertion := mock.NewMockAssertion(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		login := newLogin(nil, nil, userMgnt, nil, shareMgnt, loginDB, config, nil, nil)
		login.assertion = assertion
		login.hydraPublic = hydraPublic
		login.trace = trace
		reqInfo := &interfaces.AccessTokenReq{
			Account:      "account1",
			Password:     "ZQ5YxG74U2ts3NGva9nkCtvJacGGSRUbWl4uMm9bz/Zj7vynQDmu4TriXCXzIwGbj7wUoagCkU7JfmvSX9gdmULRWZ2pfax0SHm5YPifjlBZJ2mtlqtZBIfX3qZY5WQU29PWwLQV+WdNfGAhGLfT+k2iTj9P5TRkn/Ist3JQydSHFFYu1RBRsew0KowGJj7C6JrDZCNRKEM1Aap+GiqU4tFOkaOK2TxshH819n5hmwcIisc/8+sahxcDBi0ELdTC18Aia8cyFNdMZVYcbITv/RNyQm0YxVVGFl3++qLpsj0IArkOr+MAMffOJZfTV7wm2LYgxMeLaVDaDPrvtMniqQ==",
			ClientID:     "xx",
			ClientSecret: "xx",
		}
		cfg := interfaces.Config{
			EnablePWDLock:      true,
			EnableThirdPWDLock: true,
			PWDErrCnt:          3,
			PWDLockTime:        10,
		}
		userInfo := interfaces.UserBaseInfo{
			ID:             "id1",
			AuthType:       interfaces.Local,
			PwdErrCnt:      2,
			PwdErrLastTime: time.Now().Unix() - 1,
			DisableStatus:  false,
			LDAPType:       interfaces.LDAPServerType(0),
			DomainPath:     "",
		}

		ctx := context.Background()
		tmpErr := errors.New("Unknown error")
		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

		visitor := interfaces.Visitor{
			Language: interfaces.AmericanEnglish,
		}
		Convey("base64 decoded error", func() {
			reqInfo.Password = "1"
			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "illegal base64 data at input byte 0"))
		})
		Convey("rsa decrypted error", func() {
			reqInfo.Password = "MTExCg=="
			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "crypto/rsa: decryption error"))
		})
		Convey("GetConfigFromShareMgnt error", func() {
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, tmpErr)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("user disabled", func() {
			userInfo.DisableStatus = true
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(common.UserDisabled, "账户被禁用"))
		})
		Convey("user locked", func() {
			userInfo.PwdErrCnt = 3
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.Unauthorized, "账户被锁定, 请联系管理员"))
		})
		Convey("UserAuth failed, pwd lock disable", func() {
			userInfo.PwdErrCnt = 0
			cfg.EnablePWDLock = false
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "invalid_password", nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPError("", login.authFailedMap["invalid_password"], nil))
		})
		Convey("UserAuth failed, UpdatePWDErrInfo error", func() {
			userInfo.PwdErrLastTime = time.Now().Add(-10 * time.Minute).Unix()
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "invalid_password", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tmpErr)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("UserAuth failed, UpdatePWDErrInfo success", func() {
			userInfo.PwdErrCnt = 0
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "invalid_password", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPError("", login.authFailedMap["invalid_password"], nil))
		})
		Convey("UserAuth failed, UpdatePWDErrInfo success, account locked", func() {
			userInfo.PwdErrCnt = 2
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "invalid_password", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(common.PWDThirdFailed, "密码错误次数已达上限，账户将被锁定."))
		})
		Convey("UserAuth success, UpdatePWDErrInfo success", func() {
			userInfo.PwdErrCnt = 0
			config.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(cfg, nil)
			userMgnt.EXPECT().UserAuth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, "", nil)
			userMgnt.EXPECT().UpdatePWDErrInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := login.pwdAuth(ctx, &visitor, reqInfo, &userInfo)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAccessToken(t *testing.T) {
	Convey("GetAccesstoken", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userMgnt := mock.NewMockDnUserManagement(ctrl)
		assertion := mock.NewMockAssertion(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		accesstokenPerm := mock.NewMockAccessTokenPerm(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		login := newLogin(nil, nil, userMgnt, nil, nil, nil, nil, nil, nil)
		login.assertion = assertion
		login.accessTokenPerm = accesstokenPerm
		login.hydraPublic = hydraPublic
		login.trace = trace

		reqInfo := &interfaces.AccessTokenReq{
			Account:      "account1",
			ClientID:     "xx",
			ClientSecret: "xx",
		}
		userInfo := interfaces.UserBaseInfo{}
		tokenInfo := &interfaces.TokenInfo{}
		ctx := context.Background()
		tmpErr := errors.New("unknown error")
		visitor := interfaces.Visitor{
			Language: interfaces.AmericanEnglish,
		}
		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("account match failed", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, userInfo, tmpErr)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("account match failed, account doesn't exists", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, userInfo, nil)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(common.InvalidAccountORPassword, "账户密码错误"))
		})
		Convey("perm check failed, unknown error", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, userInfo, nil)
			accesstokenPerm.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, tmpErr)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("perm check failed, no perm", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, userInfo, nil)
			accesstokenPerm.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.Forbidden, "This app doesn't have permission to access token."))
		})
		Convey("createJWK failed, unknown error", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, userInfo, nil)
			accesstokenPerm.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			assertion.EXPECT().CreateJWK(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", tmpErr)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("assertionfortoken failed, unknown error", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, userInfo, nil)
			accesstokenPerm.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			assertion.EXPECT().CreateJWK(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("assertionStr", nil)
			hydraPublic.EXPECT().AssertionForToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, tmpErr)

			_, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, tmpErr)
		})
		Convey("get access_token success", func() {
			userMgnt.EXPECT().AccountMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, userInfo, nil)
			accesstokenPerm.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			assertion.EXPECT().CreateJWK(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("assertionStr", nil)
			hydraPublic.EXPECT().AssertionForToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(tokenInfo, nil)

			token, err := login.GetAccessToken(ctx, &visitor, reqInfo)

			assert.Equal(t, err, nil)
			assert.Equal(t, token.AccessToken, tokenInfo.AccessToken)
		})
	})
}
