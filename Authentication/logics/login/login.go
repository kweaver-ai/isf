// Package login 逻辑层
package login

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	accesstokenperm "Authentication/logics/access_token_perm"
	"Authentication/tapi/ethriftexception"

	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/interfaces"
	"Authentication/logics"
	Assertion "Authentication/logics/assertion"
	"Authentication/logics/conf"
	"Authentication/logics/sms"
	tic "Authentication/logics/ticket"
)

var (
	lOnce sync.Once
	l     *login
)

// 定义错误码唯一标识
const (
	_ int = iota
	i18nSMSVCodeNotSend
	i18nSMSVCodeInfoRequired
)

type login struct {
	hydraAdmin        interfaces.DnHydraAdmin
	hydraPublic       interfaces.DnHydraPublic
	userManagement    interfaces.DnUserManagement
	eacp              interfaces.DnEacp
	sharemgnt         interfaces.DnShareMgnt
	loginDB           interfaces.DBLogin
	dbTracePool       *sqlx.DB
	config            interfaces.Conf
	assertion         interfaces.Assertion
	aSMS              interfaces.LogicsAnonymousSMS
	accessTokenPerm   interfaces.AccessTokenPerm
	ticket            interfaces.LogicsTicket
	privateKey        *rsa.PrivateKey
	trace             observable.Tracer
	i18n              *common.I18n
	logger            common.Logger
	authFailedMap     map[string]int
	ldapServerTypeMap map[interfaces.LDAPServerType]int32
	dbAnonymousSMS    interfaces.DBAnonymousSMS
	tracer            observable.Tracer
}

// NewLogin 新建Login接口操作对象
func NewLogin() *login {
	lOnce.Do(func() {
		authFailedMap := make(map[string]int)
		authFailedMap["invalid_password"] = common.InvalidAccountORPassword
		authFailedMap["initial_password"] = common.PasswordISInitial
		authFailedMap["password_not_safe"] = common.PasswordNotSafe
		authFailedMap["under_control_password_expire"] = common.ControledPasswordExpire
		authFailedMap["password_expire"] = common.PasswordExpire

		ldapServerTypeMap := make(map[interfaces.LDAPServerType]int32)
		ldapServerTypeMap[interfaces.WindowAD] = 1
		ldapServerTypeMap[interfaces.OtherLDAP] = 2

		// pem 解码
		blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
		// X509解码
		privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		l = &login{
			hydraAdmin:      logics.DnHydraAdmin,
			hydraPublic:     logics.DnHydraPublic,
			userManagement:  logics.DnUserManagement,
			eacp:            logics.DnEacp,
			sharemgnt:       logics.DnShareMgnt,
			loginDB:         logics.DBLogin,
			dbTracePool:     logics.DBTracePool,
			config:          conf.NewConf(),
			assertion:       Assertion.NewAssertion(),
			aSMS:            sms.NewAnonymousSMS(),
			accessTokenPerm: accesstokenperm.NewAccessTokenPerm(),
			ticket:          tic.NewTicket(),
			privateKey:      privateKey,
			trace:           common.SvcARTrace,
			i18n: common.NewI18n(common.I18nMap{
				i18nSMSVCodeNotSend: {
					interfaces.SimplifiedChinese:  "无效的输入，请先获取验证码",
					interfaces.TraditionalChinese: "無效的輸入，請先取得驗證碼",
					interfaces.AmericanEnglish:    "Invalid input, please obtain verification code first",
				},
				i18nSMSVCodeInfoRequired: {
					interfaces.SimplifiedChinese:  "请输入验证码",
					interfaces.TraditionalChinese: "請輸入驗證碼訊息",
					interfaces.AmericanEnglish:    "Please enter verification code information",
				},
			}),
			logger:            common.NewLogger(),
			authFailedMap:     authFailedMap,
			ldapServerTypeMap: ldapServerTypeMap,
			dbAnonymousSMS:    logics.DBAnonymousSMS,
			tracer:            common.SvcARTrace,
		}
	})

	return l
}

func (l *login) SingleSignOn(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.SSOLoginInfo) (info *interfaces.TokenInfo, err error) {
	l.trace.SetInternalSpanName("逻辑层-单点登录")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	oauthReqInfo := &interfaces.AuthorizeInfo{
		ClientID:     reqInfo.ClientID,
		RedirectURI:  reqInfo.RedirectURI,
		ResponseType: reqInfo.ResponseType,
		Scope:        reqInfo.Scope,
	}

	loginChallenge, loginSession, err := l.hydraPublic.AuthorizeRequest(oauthReqInfo)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := l.hydraAdmin.GetLoginRequestInformation(loginChallenge)
	if err != nil {
		return nil, err
	}

	authReqInfo := &interfaces.ThirdPartyAuthInfo{
		Udids: reqInfo.Udids,
		IP:    reqInfo.IP,
		Credential: interfaces.ThirdPartyCredential{
			ID:     reqInfo.Credential.ID,
			Params: reqInfo.Credential.Params,
		},
		Device: interfaces.DeviceInfo{
			Name:        deviceInfo.Name,
			ClientType:  deviceInfo.ClientType,
			Description: deviceInfo.Description,
		},
	}
	if reqInfo.Credential.ID == "aishu" {
		var account string
		// 校验ticket。检验通过后，用登录名覆盖reqInfo.Credential.Params，从而兼容eacp getbythirdparty逻辑
		account, err = l.validateTicket(newCtx, visitor, reqInfo.Credential.Params.(map[string]interface{})["ticket"].(string), reqInfo.ClientID)
		if err != nil {
			return nil, err
		}
		authReqInfo.Credential.Params = map[string]interface{}{"account": account}
	}

	loginInfo, err := l.eacp.ThirdPartyAuthentication(newCtx, visitor, authReqInfo)
	if err != nil {
		return nil, err
	}

	redirURL, err := l.hydraAdmin.AcceptLoginRequest(loginInfo.Subject, loginChallenge)
	if err != nil {
		return nil, err
	}

	consentChallenge, loginSession, err := l.hydraPublic.VerifierLogin(redirURL, loginSession)
	if err != nil {
		return nil, err
	}

	redirURL, err = l.hydraAdmin.AcceptConsentRequest(reqInfo.Scope, consentChallenge, loginInfo.Context)
	if err != nil {
		return nil, err
	}

	info, err = l.hydraPublic.VerifierConsent(redirURL, reqInfo.ResponseType, loginSession)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (l *login) Anonymous(visitor *interfaces.Visitor, reqInfo *interfaces.AnonymousLoginInfo) (*interfaces.TokenInfo, error) {
	oauthReqInfo := &interfaces.AuthorizeInfo{
		ClientID:     reqInfo.ClientID,
		RedirectURI:  reqInfo.RedirectURI,
		ResponseType: reqInfo.ResponseType,
		Scope:        reqInfo.Scope,
	}

	loginChallenge, loginSession, err := l.hydraPublic.AuthorizeRequest(oauthReqInfo)
	if err != nil {
		return nil, err
	}

	_, err = l.userManagement.AnonymousAuthentication(context.Background(), visitor, reqInfo.Credential.Account, reqInfo.Credential.Password, "")
	if err != nil {
		return nil, err
	}

	redirURL, err := l.hydraAdmin.AcceptLoginRequest(reqInfo.Credential.Account, loginChallenge)
	if err != nil {
		return nil, err
	}

	consentChallenge, loginSession, err := l.hydraPublic.VerifierLogin(redirURL, loginSession)
	if err != nil {
		return nil, err
	}

	redirURL, err = l.hydraAdmin.AcceptConsentRequest(reqInfo.Scope, consentChallenge, map[string]interface{}{"visitor_type": "anonymous"})
	if err != nil {
		return nil, err
	}

	info, err := l.hydraPublic.VerifierConsent(redirURL, reqInfo.ResponseType, loginSession)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (l *login) Anonymous2(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AnonymousLoginInfo2, referrer string) (token *interfaces.TokenInfo, err error) {
	l.trace.SetInternalSpanName("逻辑层-匿名登录v2")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	// 判断该匿名帐户是否需要手机验证码校验
	needVerifyMobile, err := l.userManagement.GetAnonymityInfoByID(newCtx, visitor, reqInfo.Credential.Account)
	if err != nil {
		return nil, err
	}
	phoneNumber := ""
	if needVerifyMobile {
		if reqInfo.VCode.ID == "" && reqInfo.VCode.Content != "" {
			return nil, rest.NewHTTPErrorV2(rest.BadRequest, l.i18n.Load(i18nSMSVCodeNotSend, visitor.Language))
		}
		if reqInfo.VCode.ID == "" || reqInfo.VCode.Content == "" {
			return nil, rest.NewHTTPErrorV2(rest.BadRequest, l.i18n.Load(i18nSMSVCodeInfoRequired, visitor.Language))
		}
		if reqInfo.VisitorName == "" {
			return nil, rest.NewHTTPErrorV2(rest.BadRequest, "visitor_name required")
		}
		phoneNumber, err = l.aSMS.Validate(newCtx, visitor, reqInfo.VCode.ID, reqInfo.VCode.Content, reqInfo.Credential.Account)
		if err != nil {
			return nil, err
		}
		// 手机号脱敏
		phoneNumber = phoneNumber[:3] + "****" + phoneNumber[7:]
	}

	// 权限校验
	_, err = l.userManagement.AnonymousAuthentication(newCtx, visitor, reqInfo.Credential.Account, reqInfo.Credential.Password, referrer)
	if err != nil {
		return nil, err
	}

	// 申请断言
	privateClaims := map[string]interface{}{
		"ext": map[string]interface{}{
			"client_id":    reqInfo.ClientID,
			"visitor_type": "anonymous",
			"phone_number": phoneNumber,
			"visitor_name": reqInfo.VisitorName,
		},
	}
	assertionStr, err := l.assertion.CreateJWK(context.Background(), reqInfo.Credential.Account, time.Hour, privateClaims)
	if err != nil {
		return nil, err
	}

	// 断言换取令牌
	token, err = l.hydraPublic.AssertionForToken(reqInfo.ClientID, reqInfo.ClientSecret, assertionStr)
	if err != nil {
		return nil, err
	}

	/*
		删除验证码的逻辑原本在 LogicsAnonymousSMS.Validate 中，验证码校验成功后，立即删除数据库中的验证码。
		客户要求先验证提取码，再验证验证码。
		如果提取码不正确，就不验证验证码，避免重复发送验证码，浪费短信配额。
		因此删除验证码的这一步，在上述一系列的逻辑完成之后，再执行。
	*/
	if needVerifyMobile {
		err = l.dbAnonymousSMS.DeleteByIDs(newCtx, []string{reqInfo.VCode.ID})
		if err != nil {
			return nil, err
		}
	}
	return token, nil
}

// checkAccountMatch 检查账户匹配
func (l *login) checkAccountMatch(ctx context.Context, visitor *interfaces.Visitor, account string, enablePrefixMatch bool,
	config *interfaces.Config) (exist bool, userInfo interfaces.UserBaseInfo, err error) {
	// 账户匹配
	exist, userInfo, err = l.userManagement.AccountMatch(ctx, visitor, account, config.EnableIDCardLogin, enablePrefixMatch)
	if err != nil {
		return
	}

	// 如果账户存在
	// 如果三权分立未开启，屏蔽security、audit、system
	// 如果三权分立开启，则屏蔽system
	if exist {
		if !config.TriSystemStatus && (userInfo.ID == logics.SecurityAdminID || userInfo.ID == logics.AuditAdminID || userInfo.ID == logics.SystemAdminID) {
			exist = false
			userInfo = interfaces.UserBaseInfo{}
		}
		if config.TriSystemStatus && userInfo.ID == logics.SystemAdminID {
			exist = false
			userInfo = interfaces.UserBaseInfo{}
		}
	}

	return exist, userInfo, nil
}

func (l *login) ClientAccountAuth(ctx context.Context, visitor *interfaces.Visitor, req *interfaces.ClientLoginReq) (userID string, err error) {
	l.trace.SetInternalSpanName("逻辑层-客户端账户认证")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	// 获取配置信息
	rg := make(map[interfaces.ConfigKey]bool)
	rg[interfaces.EnableIDCardLogin] = true
	rg[interfaces.EnablePWDLock] = true
	rg[interfaces.EnableThirdPWDLock] = true
	rg[interfaces.PWDErrCnt] = true
	rg[interfaces.PWDLockTime] = true
	rg[interfaces.VCodeConfig] = true
	rg[interfaces.TriSystemStatus] = true
	config, err := l.config.GetConfigFromShareMgnt(newCtx, visitor, rg)
	if err != nil {
		return
	}

	// 获取开启的域状态信息
	enablePrefixMatch, err := l.loginDB.GetDomainStatus()
	if err != nil {
		return
	}

	// 认证失败时，detail作为error的一部分信息返回。通常有两个字段：
	// isShowStatus: 告知客户端下次登录，是否需要输入图形验证码，用于锁定场景和图形验证码场景
	// remainlockTime：账户剩余待锁定时间
	detail := make(map[string]interface{})
	detail["isShowStatus"] = false

	// 账户匹配
	exist, userInfo, err := l.checkAccountMatch(newCtx, visitor, req.Account, enablePrefixMatch, &config)
	if err != nil {
		return
	}

	// 如果需要验证短信或者动态密码，但是账户不存在，则返回验证码错误，如果是图形验证码，则继续验证
	// config.VCodeConfig.Enable 只有在开启图形验证码功能时，才为true，其他验证码或者没有验证码时为false
	if !exist {
		detail["isShowStatus"] = config.VCodeConfig.Enable
		if req.Option.VCodeType == interfaces.DualAuthSMS {
			err = rest.NewHTTPError("", common.ImageVCodeISWrong, detail)
			return
		} else if req.Option.VCodeType == interfaces.DualAuthOTP {
			err = rest.NewHTTPError("", common.OTPWrong, detail)
			return
		}
	}

	// 当图形验证码开启，且密码错误次数达到上限时，需要显示图形验证码
	if config.VCodeConfig.Enable && userInfo.PwdErrCnt >= config.VCodeConfig.PWDErrCnt {
		detail["isShowStatus"] = true
	}

	// 校验验证码
	if err = l.validateVCode(&userInfo, &req.Option, &config, detail); err != nil {
		return
	}

	// 检查账户是否被锁定
	bLocked := false
	nRemaingLockTime := int64(0)
	bEnablePWDLock := config.EnablePWDLock && (userInfo.AuthType == interfaces.Local || config.EnableThirdPWDLock)
	if userInfo.ID != "" && bEnablePWDLock {
		if userInfo.PwdErrCnt >= config.PWDErrCnt {
			nRemaingLockTime = int64(config.PWDLockTime) - (common.Now().Unix()-userInfo.PwdErrLastTime)/60
			bLocked = true
		}
	}

	// 账户密码验证
	pwdErrCnt := 0
	if err = l.validatePassword(newCtx, visitor, &userInfo, req.Password, detail); err != nil {
		// 如果账户不存在或者用户被禁用或者用户被锁定，则跳过错误次数处理
		if userInfo.ID == "" || userInfo.DisableStatus || bLocked {
			return
		}

		// 既没开启密码错误锁定，也没开启图形验证码功能，此时不关心用户密码错误次数，直接返回。
		if !bEnablePWDLock && !config.VCodeConfig.Enable {
			return
		}

		if _, ok := err.(*rest.HTTPError); !ok {
			return
		}
		if err.(*rest.HTTPError).Code != common.InvalidAccountORPassword {
			return
		}

		// 根据本次密码错误和上次密码错误的时间差。计算出新的密码错误次数
		var interval int64 = 5
		if (common.Now().Unix()-userInfo.PwdErrLastTime)/60 < interval {
			pwdErrCnt = userInfo.PwdErrCnt + 1
		} else {
			pwdErrCnt = 1
		}

		// 更新账户密码错误信息
		if tmpErr := l.userManagement.UpdatePWDErrInfo(newCtx, visitor, userInfo.ID, pwdErrCnt); tmpErr != nil {
			return "", tmpErr
		}

		// 计算该账户下次登录，是否需要显示校验码
		if config.VCodeConfig.Enable && pwdErrCnt >= config.VCodeConfig.PWDErrCnt {
			detail["isShowStatus"] = true
			err = rest.NewHTTPError("", common.InvalidAccountORPassword, detail)
		}

		return
	}

	// 检查账户是否被禁用/自动禁用
	if userInfo.DisableStatus {
		err = rest.NewHTTPError("", common.UserDisabled, detail)
		return
	}

	// 检查账户是否被锁定
	if bLocked {
		// 计算该账户剩余锁定时间
		detail["remainlockTime"] = nRemaingLockTime
		err = rest.NewHTTPError("", common.AccountLocked, detail)
		return
	}

	err = l.userManagement.UpdatePWDErrInfo(newCtx, visitor, userInfo.ID, pwdErrCnt)
	return userInfo.ID, err
}

//nolint:lll
func (l *login) getAccessToken(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AccessTokenReq, permCheck func(context.Context, *interfaces.Visitor, *interfaces.AccessTokenReq, *interfaces.UserBaseInfo) error) (token *interfaces.TokenInfo, err error) {
	l.trace.SetInternalSpanName("逻辑层-申请用户令牌")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	// 账户匹配
	exist, userInfo, err := l.userManagement.AccountMatch(newCtx, visitor, reqInfo.Account, false, false)
	if err != nil {
		return nil, err
	}
	if !exist {
		err = rest.NewHTTPErrorV2(common.InvalidAccountORPassword, "账户密码错误")
		return nil, err
	}

	// 权限校验
	err = permCheck(newCtx, visitor, reqInfo, &userInfo)
	if err != nil {
		return nil, err
	}

	// 生成断言 && 断言换取令牌
	privateClaims := map[string]interface{}{
		"ext": map[string]interface{}{
			"visitor_type": "realname",
			"login_ip":     "",
			"account_type": "other",
			"udid":         "",
			"client_type":  "app",
			"client_id":    reqInfo.ClientID, // 限制生成断言与使用断言的app必须一致
		},
	}
	assertion, err := l.assertion.CreateJWK(newCtx, userInfo.ID, time.Hour, privateClaims)
	if err != nil {
		return nil, err
	}
	token, err = l.hydraPublic.AssertionForToken(reqInfo.ClientID, reqInfo.ClientSecret, assertion)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (l *login) GetAccessToken(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AccessTokenReq) (*interfaces.TokenInfo, error) {
	return l.getAccessToken(ctx, visitor, reqInfo, l.checkAppAccessTokenPerm)
}

func (l *login) PwdAuth(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AccessTokenReq) (*interfaces.TokenInfo, error) {
	return l.getAccessToken(ctx, visitor, reqInfo, l.pwdAuth)
}

func (l *login) pwdAuth(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AccessTokenReq, userInfo *interfaces.UserBaseInfo) (err error) {
	l.trace.SetInternalSpanName("逻辑层-账号密码校验")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	reqInfo.Password, err = logics.DecodeAndDecrypt(reqInfo.Password, l.privateKey)
	if err != nil {
		return err
	}

	// 获取配置信息
	rg := make(map[interfaces.ConfigKey]bool)
	rg[interfaces.EnablePWDLock] = true
	rg[interfaces.EnableThirdPWDLock] = true
	rg[interfaces.PWDErrCnt] = true
	rg[interfaces.PWDLockTime] = true
	config, err := l.config.GetConfigFromShareMgnt(newCtx, nil, rg)
	if err != nil {
		return err
	}

	// 检查账户状态：是否被禁用? 是否被锁定?
	if userInfo.DisableStatus {
		return rest.NewHTTPErrorV2(common.UserDisabled, "账户被禁用")
	}
	if config.EnablePWDLock && (userInfo.AuthType == interfaces.Local || config.EnableThirdPWDLock) {
		if userInfo.PwdErrCnt >= config.PWDErrCnt {
			return rest.NewHTTPErrorV2(rest.Unauthorized, "账户被锁定, 请联系管理员")
		}
	}

	// 校验密码
	pwdErrCnt := 0
	if err = l.validatePassword(newCtx, visitor, userInfo, reqInfo.Password, nil); err != nil {
		// 没有开启密码错误锁定，直接返回。
		enableAccountLock := config.EnablePWDLock && (userInfo.AuthType == interfaces.Local || config.EnableThirdPWDLock)
		if !enableAccountLock {
			l.logger.Errorln(err)
			return err
		}

		// 非密码不正确导致的错误，直接返回
		if err.(*rest.HTTPError).Code != common.InvalidAccountORPassword {
			l.logger.Errorln(err)
			return err
		}

		// 根据本次密码错误和上次密码错误的时间差。计算出新的密码错误次数。
		// 计算规则：5分钟内密码连续错误，密码错误次数加1。密码连续错误，但间隔在5分钟以上，则记为1次。
		var interval int64 = 5
		if (common.Now().Unix()-userInfo.PwdErrLastTime)/60 < interval {
			pwdErrCnt = userInfo.PwdErrCnt + 1
		} else {
			pwdErrCnt = 1
		}

		// 更新账户密码错误信息
		if tmpErr := l.userManagement.UpdatePWDErrInfo(newCtx, visitor, userInfo.ID, pwdErrCnt); tmpErr != nil {
			l.logger.Errorln(err)
			return tmpErr
		}

		// 根据新的密码错误信息，计算当前账户是否被锁定
		if enableAccountLock && pwdErrCnt >= config.PWDErrCnt {
			err = rest.NewHTTPErrorV2(common.PWDThirdFailed, "密码错误次数已达上限，账户将被锁定.")
			l.logger.Errorln(err)
			return err
		}
		return err
	}

	// 更新密码错误次数
	if err = l.userManagement.UpdatePWDErrInfo(newCtx, visitor, userInfo.ID, pwdErrCnt); err != nil {
		l.logger.Errorln(err)
		return err
	}
	return nil
}

func (l *login) checkAppAccessTokenPerm(ctx context.Context, visitor *interfaces.Visitor, reqInfo *interfaces.AccessTokenReq, userInfo *interfaces.UserBaseInfo) (err error) {
	l.trace.SetInternalSpanName("逻辑层-获取任意用户令牌权限校验")
	_, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	hasPerm, err := l.accessTokenPerm.CheckAppAccessTokenPerm(reqInfo.ClientID)
	if err != nil {
		return err
	}
	if !hasPerm {
		err = rest.NewHTTPErrorV2(rest.Forbidden, "This app doesn't have permission to access token.")
		return err
	}

	return nil
}

func (l *login) validatePassword(ctx context.Context, visitor *interfaces.Visitor, userInfo *interfaces.UserBaseInfo, password string, detail map[string]interface{}) (err error) {
	if userInfo.ID == "" {
		err = rest.NewHTTPError("", common.InvalidAccountORPassword, detail)
		return
	}

	var result bool
	var reason string
	switch userInfo.AuthType {
	case interfaces.Local:
		result, reason, err = l.userManagement.UserAuth(ctx, visitor, userInfo.ID, password)
		if err != nil {
			return
		}
		if !result {
			err = rest.NewHTTPError("", l.authFailedMap[reason], detail)
		}
	case interfaces.Domain:
		result, err = l.sharemgnt.UsrmDomainAuth(ctx, userInfo.Account, userInfo.DomainPath, password, l.ldapServerTypeMap[userInfo.LDAPType])
		if err != nil {
			switch v := err.(type) {
			case *ethriftexception.NcTException:
				err = rest.NewHTTPError("", int(rest.Unauthorized+v.GetErrID()), detail)
			default:
				err = rest.NewHTTPError(err.Error(), rest.InternalServerError, detail)
			}

			return
		}
		if !result {
			err = rest.NewHTTPError("", common.InvalidAccountORPassword, detail)
		}
	case interfaces.Third:
		result, err = l.sharemgnt.UsrmThirdAuth(ctx, userInfo.Account, password)
		if err != nil {
			switch v := err.(type) {
			case *ethriftexception.NcTException:
				err = rest.NewHTTPError("", int(rest.Unauthorized+v.GetErrID()), detail)
			default:
				err = rest.NewHTTPError(err.Error(), rest.InternalServerError, detail)
			}

			return
		}
		if !result {
			err = rest.NewHTTPError("", common.InvalidAccountORPassword, detail)
		}
	}

	return
}

func (l *login) validateVCode(userInfo *interfaces.UserBaseInfo, option *interfaces.ClientLoginOption, config *interfaces.Config, detail map[string]interface{}) (err error) {
	// 有点奇怪，正常只有图形验证码 config.VCodeConfig.Enable才为true
	// 没看懂
	if config.VCodeConfig.Enable && userInfo.PwdErrCnt >= config.VCodeConfig.PWDErrCnt {
		option.VCodeType = interfaces.ImageVCode
	}

	switch option.VCodeType {
	case interfaces.ImageVCode:
		err = l.sharemgnt.UsrmIMAGECodeValidate(option.UUID, option.VCode)
	case interfaces.DualAuthSMS:
		err = l.sharemgnt.UsrmSMSValidate(userInfo.ID, option.VCode)
	case interfaces.DualAuthOTP:
		err = l.sharemgnt.UsrmOTPValidate(userInfo.ID, option.VCode)
	}

	if err != nil {
		switch v := err.(type) {
		case *ethriftexception.NcTException:
			err = rest.NewHTTPError("", int(rest.Unauthorized+v.GetErrID()), detail)
		default:
			err = rest.NewHTTPError(err.Error(), rest.InternalServerError, detail)
		}

		return
	}

	return
}

func (l *login) validateTicket(ctx context.Context, visitor *interfaces.Visitor, ticket, clientID string) (account string, err error) {
	ticketID, err := logics.DecodeAndDecrypt(ticket, l.privateKey)
	if err != nil {
		l.logger.Errorf("解密ticket失败: %v", err)
		return "", err
	}

	userID, err := l.ticket.Validate(ctx, ticketID, clientID)
	if err != nil {
		l.logger.Errorf("Validate ticket failed, err: %v", err)
		return "", err
	}

	userBaseInfo, err := l.userManagement.GetUserInfo(ctx, visitor, userID)

	return userBaseInfo.Account, err
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

const _ = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsyOstgbYuubBi2PUqeVj
GKlkwVUY6w1Y8d4k116dI2SkZI8fxcjHALv77kItO4jYLVplk9gO4HAtsisnNE2o
wlYIqdmyEPMwupaeFFFcg751oiTXJiYbtX7ABzU5KQYPjRSEjMq6i5qu/mL67XTk
hvKwrC83zme66qaKApmKupDODPb0RRkutK/zHfd1zL7sciBQ6psnNadh8pE24w8O
2XVy1v2bgSNkGHABgncR7seyIg81JQ3c/Axxd6GsTztjLnlvGAlmT1TphE84mi99
fUaGD2A1u1qdIuNc+XuisFeNcUW6fct0+x97eS2eEGRr/7qxWmO/P20sFVzXc2bF
1QIDAQAB
-----END PUBLIC KEY-----`
