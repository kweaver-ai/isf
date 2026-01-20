// Package login 协议层
package login

import (
	_ "embed" // 标准用法
	"net/http"
	"reflect"
	"strings"
	"sync"

	"Authentication/common"
	"Authentication/driveradapters"
	authSchema "Authentication/driveradapters/jsonschema/auth_schema"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	accessTokenSchema "Authentication/driveradapters/jsonschema/access_token_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/login"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)
	RegisterPrivate(engine *gin.Engine)
}

type restHandler struct {
	login                interfaces.Login
	clientLoginSchema    *gojsonschema.Schema
	anonyousLogin2Schema *gojsonschema.Schema
	pwdAuthSchemaStr     *gojsonschema.Schema
	accessTokenSchema    *gojsonschema.Schema
}

var (
	once sync.Once
	r    RESTHandler
)

type jsonValueDesc = rest.JSONValueDesc

// NewRESTHandler 创建login RESTHandler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		clientLoginSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.ClientLoginSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		anonyousLogin2Schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.AnonymousLogin2))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		pwdAuthSchemaStr, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.PwdAuthSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		accessTokenSchema1, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(accessTokenSchema.AccessTokenSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		r = &restHandler{
			login:                login.NewLogin(),
			clientLoginSchema:    clientLoginSchema,
			anonyousLogin2Schema: anonyousLogin2Schema,
			pwdAuthSchemaStr:     pwdAuthSchemaStr,
			accessTokenSchema:    accessTokenSchema1,
		}
	})

	return r
}

// RegisterPublic 注册开放API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/sso", observable.MiddlewareTrace(common.SvcARTrace), r.singleSignOn)
	engine.POST("/api/authentication/v1/pwd-auth", observable.MiddlewareTrace(common.SvcARTrace), r.pwdAuth)
	engine.POST("/api/authentication/v1/access_token", observable.MiddlewareTrace(common.SvcARTrace), r.getAccessToken)
	engine.POST("/api/authentication/v1/anonymous", r.anonymous)
	engine.POST("/api/authentication/v2/anonymous", observable.MiddlewareTrace(common.SvcARTrace), r.anonymous2)
}

func (r *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/client-account-auth", r.clientAccountAuth)
}

func (r *restHandler) getAccessToken(c *gin.Context) {
	var err error
	var jsonV map[string]interface{}
	if err = util.ValidateAndBindGin(c, r.accessTokenSchema, &jsonV); err != nil {
		rest.ReplyError(c, err)
		return
	}
	var ok bool
	reqInfo := &interfaces.AccessTokenReq{}
	reqInfo.Account = jsonV["account"].(string)
	reqInfo.ClientID, reqInfo.ClientSecret, ok = c.Request.BasicAuth()
	if !ok || reqInfo.ClientID == "" || reqInfo.ClientSecret == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))
		return
	}

	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	res, err := r.login.GetAccessToken(c, &visitor, reqInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"access_token": res.AccessToken,
		"expires_in":   res.ExpirsesIn,
		"scope":        res.Scope,
		"token_type":   res.TokenType,
	}
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) pwdAuth(c *gin.Context) {
	var reqJSON map[string]interface{}
	if err := util.ValidateAndBindGin(c, r.pwdAuthSchemaStr, &reqJSON); err != nil {
		rest.ReplyError(c, err)
		return
	}

	var ok bool
	reqInfo := &interfaces.AccessTokenReq{}
	reqInfo.Account = reqJSON["account"].(string)
	reqInfo.Password = reqJSON["password"].(string)
	reqInfo.ClientID, reqInfo.ClientSecret, ok = c.Request.BasicAuth()
	if !ok || reqInfo.ClientID == "" || reqInfo.ClientSecret == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))
		return
	}

	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	tokenInfo, err := r.login.PwdAuth(c, &visitor, reqInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"access_token": tokenInfo.AccessToken,
		"expires_in":   tokenInfo.ExpirsesIn,
		"scope":        tokenInfo.Scope,
		"token_type":   tokenInfo.TokenType,
	}
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) clientAccountAuth(c *gin.Context) {
	var clientLoginReq interfaces.ClientLoginReq
	var err error
	if err = util.ValidateAndBindGin(c, r.clientLoginSchema, &clientLoginReq); err != nil {
		rest.ReplyError(c, err)
		return
	}
	if err = r.checkLoginOption(&clientLoginReq.Option); err != nil {
		rest.ReplyError(c, err)
		return
	}

	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	userID, err := r.login.ClientAccountAuth(c, &visitor, &clientLoginReq)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := make(map[string]interface{})
	resInfo["user_id"] = userID
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

// SingleSignOn 单点登录
func (r *restHandler) singleSignOn(c *gin.Context) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	strDesc := make(map[string]*jsonValueDesc)
	credentialObjDesc := make(map[string]*jsonValueDesc)
	credentialObjDesc["params"] = &jsonValueDesc{Kind: reflect.Map, Required: true}
	credentialObjDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["credential"] = &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: credentialObjDesc}
	objDesc["client_id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["redirect_uri"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["response_type"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["scope"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	strDesc["element"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["udids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	req := interfaces.SSOLoginInfo{
		Udids: []string{},
	}
	jsonObj := jsonV.(map[string]interface{})
	jsonStr := make(map[string]string)
	jsonStr["client_id"] = jsonObj["client_id"].(string)
	jsonStr["redirect_uri"] = jsonObj["redirect_uri"].(string)
	jsonStr["response_type"] = jsonObj["response_type"].(string)
	jsonStr["scope"] = jsonObj["scope"].(string)
	if objDesc["udids"].Exist {
		udids := jsonObj["udids"].([]interface{})
		for _, udid := range udids {
			req.Udids = append(req.Udids, udid.(string))
		}
	}
	req.IP = strings.Split(c.Request.Header["X-Forwarded-For"][0], ",")[0]
	req.Credential.ID = jsonObj["credential"].(map[string]interface{})["id"].(string)
	req.Credential.Params = jsonObj["credential"].(map[string]interface{})["params"]

	err := r.checkSSOResponseType(jsonStr["response_type"])
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	err = util.IsEmpty(jsonStr)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	req.ClientID = jsonStr["client_id"]
	req.RedirectURI = jsonStr["redirect_uri"]
	req.ResponseType = jsonStr["response_type"]
	req.Scope = jsonStr["scope"]

	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	tokenInfo, err := r.login.SingleSignOn(c, &visitor, &req)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := make(map[string]interface{})
	switch tokenInfo.ResponseType {
	case "code":
		resInfo = map[string]interface{}{
			"code":  tokenInfo.Code,
			"scope": tokenInfo.Scope,
		}
	case "token id_token":
		resInfo = map[string]interface{}{
			"access_token": tokenInfo.AccessToken,
			"expirses_in":  tokenInfo.ExpirsesIn,
			"id_token":     tokenInfo.IDToken,
			"scope":        tokenInfo.Scope,
			"token_type":   tokenInfo.TokenType,
		}
	}

	rest.ReplyOK(c, http.StatusOK, resInfo)
}

// Anonymous 匿名登录
func (r *restHandler) anonymous(c *gin.Context) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	credentialObjDesc := make(map[string]*jsonValueDesc)
	credentialObjDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	credentialObjDesc["account"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["credential"] = &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: credentialObjDesc}
	objDesc["client_id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["redirect_uri"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["response_type"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["scope"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	req := &interfaces.AnonymousLoginInfo{}
	jsonObj := jsonV.(map[string]interface{})
	jsonStr := make(map[string]string)
	jsonStr["client_id"] = jsonObj["client_id"].(string)
	jsonStr["redirect_uri"] = jsonObj["redirect_uri"].(string)
	jsonStr["response_type"] = jsonObj["response_type"].(string)
	jsonStr["scope"] = jsonObj["scope"].(string)
	jsonStr["account"] = jsonObj["credential"].(map[string]interface{})["account"].(string)
	req.Credential.Password = jsonObj["credential"].(map[string]interface{})["password"].(string)

	if jsonStr["response_type"] != "token" {
		rest.ReplyError(c, rest.NewHTTPError("response_type is invalid", rest.BadRequest, nil))
		return
	}

	err := util.IsEmpty(jsonStr)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	req.ClientID = jsonStr["client_id"]
	req.RedirectURI = jsonStr["redirect_uri"]
	req.ResponseType = jsonStr["response_type"]
	req.Scope = jsonStr["scope"]
	req.Credential.Account = jsonStr["account"]

	visitor := interfaces.Visitor{
		Language:      driveradapters.GetXLang(c),
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	res, err := r.login.Anonymous(&visitor, req)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"access_token": res.AccessToken,
		"expirses_in":  res.ExpirsesIn,
		"scope":        res.Scope,
		"token_type":   res.TokenType,
	}

	rest.ReplyOK(c, http.StatusOK, resInfo)
}

// Anonymous 匿名登录
func (r *restHandler) anonymous2(c *gin.Context) {
	visitor := interfaces.Visitor{
		Language:      driveradapters.GetXLang(c),
		ErrorCodeType: util.GetErrorCodeType(c),
	}

	var err error
	var reqJSON map[string]interface{}
	if err = util.ValidateAndBindGin(c, r.anonyousLogin2Schema, &reqJSON); err != nil {
		rest.ReplyError(c, err)
		return
	}
	clientID, clientSecret, ok := c.Request.BasicAuth()
	if !ok {
		err = rest.NewHTTPError("http basic auth param missing", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	var reqInfo interfaces.AnonymousLoginInfo2
	reqInfo.ClientID = clientID
	reqInfo.ClientSecret = clientSecret
	reqInfo.Credential.Account = reqJSON["credential"].(map[string]interface{})["account"].(string)
	reqInfo.Credential.Password = reqJSON["credential"].(map[string]interface{})["password"].(string)
	if v, ok := reqJSON["vcode"].(map[string]interface{}); ok {
		reqInfo.VCode.ID = v["id"].(string)
		reqInfo.VCode.Content = v["content"].(string)

		if _, ok := reqJSON["visitor_name"].(string); ok {
			reqInfo.VisitorName = reqJSON["visitor_name"].(string)
		}
	}

	referrer := c.Request.Header.Get("x-referrer")
	res, err := r.login.Anonymous2(c, &visitor, &reqInfo, referrer)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"access_token": res.AccessToken,
		"expires_in":   res.ExpirsesIn,
		"scope":        res.Scope,
		"token_type":   res.TokenType,
	}
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) checkSSOResponseType(responsetype string) error {
	if (responsetype != "code") && (responsetype != "token id_token") {
		return rest.NewHTTPError("response_type is invalid", rest.BadRequest, nil)
	}

	return nil
}

func (r *restHandler) checkLoginOption(option *interfaces.ClientLoginOption) (err error) {
	// 判断是否是合法的验证码类型
	switch option.VCodeType {
	case interfaces.ImageVCode:
		if option.UUID == "" {
			err = rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "option.uuid"})
			return
		}
		if option.VCode == "" {
			err = rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "option.vcode"})
			return
		}
	case interfaces.NumVCode:
		// 认证接口，暂不支持数字验证码
		err = rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "option.vcodeType"})
		return
	case interfaces.DualAuthSMS:
		if option.VCode == "" {
			err = rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "option.vcode"})
			return
		}
	case interfaces.DualAuthOTP:
		if option.VCode == "" {
			err = rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "option.vcode"})
			return
		}
	}

	return
}
