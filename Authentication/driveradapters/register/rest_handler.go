// Package register 协议层
package register

import (
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	registerschema "Authentication/driveradapters/jsonschema/register_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/register"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)
}

type restHandler struct {
	register               interfaces.Register
	scopeMember            map[string]bool
	errInvalidParameter    *RFC6749Error
	errInternalServerError *RFC6749Error
	registerSchema         *gojsonschema.Schema
}

var (
	once sync.Once
	r    RESTHandler
)

// NewRESTHandler 创建register handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		registerSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(registerschema.RegisterSchema))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		r = &restHandler{
			register: register.NewRegister(),
			scopeMember: map[string]bool{
				"offline": true,
				"openid":  true,
				"all":     true,
			},
			errInvalidParameter: &RFC6749Error{
				Name:        "invalid_request",
				Description: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
				Code:        http.StatusBadRequest,
			},
			errInternalServerError: &RFC6749Error{
				Name:        "internal_server_error",
				Description: "Internal Server Error",
				Code:        http.StatusInternalServerError,
			},
			registerSchema: registerSchema,
		}
	})

	return r
}

// RegisterPublic 注册开放API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/oauth2/clients", r.publicRegister)
}

func (r *restHandler) publicRegister(c *gin.Context) {
	// 检查请求参数与文档是否匹配
	var jsonV map[string]interface{}
	if err := util.ValidateAndBindGin(c, r.registerSchema, &jsonV); err != nil {
		r.ReplyHydraError(c, err)
		return
	}

	// 检查参数是否合法
	clientName := jsonV["client_name"].(string)
	if clientName == "" {
		r.ReplyHydraError(c, r.errInvalidParameter.withHint("Invalid client_name"))
		return
	}

	grantTypes := []string{}
	for _, grantType := range jsonV["grant_types"].([]interface{}) {
		grantTypes = append(grantTypes, grantType.(string))
	}

	responseTypes := []string{}
	for _, responseType := range jsonV["response_types"].([]interface{}) {
		responseTypes = append(responseTypes, responseType.(string))
	}

	scopeStr := jsonV["scope"].(string)
	scopesMap := make(map[string]bool)
	scopes := strings.Split(scopeStr, " ")
	for _, scope := range scopes {
		scopesMap[scope] = true
	}
	if !reflect.DeepEqual(scopesMap, r.scopeMember) || len(scopes) != len(r.scopeMember) {
		r.ReplyHydraError(c, r.errInvalidParameter.withHint("Invalid scope"))
		return
	}

	metadata := jsonV["metadata"].(map[string]interface{})
	device := metadata["device"].(map[string]interface{})

	nDevice := map[string]interface{}{
		"client_type": device["client_type"],
	}
	if name, ok := device["name"]; ok {
		nDevice["name"] = name
	}
	if description, ok := device["description"]; ok {
		nDevice["description"] = description
	}
	metadata["device"] = nDevice

	redirectURIs := []string{}
	for _, redirectURI := range jsonV["redirect_uris"].([]interface{}) {
		redirectURIs = append(redirectURIs, redirectURI.(string))
	}

	logoutUris := []string{}
	for _, logoutURI := range jsonV["post_logout_redirect_uris"].([]interface{}) {
		logoutUris = append(logoutUris, logoutURI.(string))
	}
	if len(logoutUris) == 0 {
		r.ReplyHydraError(c, r.errInvalidParameter.withHint("Invalid post_logout_redirect_uris"))
		return
	}

	registerInfo := &interfaces.RegisterInfo{
		ClientName:             clientName,
		GrantTypes:             grantTypes,
		ResponseTypes:          responseTypes,
		Scope:                  scopeStr,
		RedirectURIs:           redirectURIs,
		PostLogoutRedirectURIs: logoutUris,
		Metadata:               metadata,
	}

	info, err := r.register.PublicRegister(registerInfo)
	if err != nil {
		r.ReplyHydraError(c, err)
		return
	}

	// 响应201
	rest.ReplyOK(c, http.StatusCreated, info)
}

// ReplyHydraError 响应 hydra 错误
func (r *restHandler) ReplyHydraError(c *gin.Context, err error) {
	var statusCode int
	var body string
	switch e := err.(type) {
	case *RFC6749Error:
		statusCode = e.Code
		body = e.Error()
	case interfaces.HydraError:
		statusCode = e.Status
		body = e.Error()
	case *rest.HTTPError:
		err := r.errInvalidParameter
		statusCode = err.Code
		body = err.withHint(e.Cause).Error()
	default:
		err := r.errInternalServerError
		statusCode = err.Code
		body = err.withHint(e.Error()).Error()
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(statusCode, body)
}

// RFC6749Error rfc6749错误
type RFC6749Error struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Hint        string `json:"error_hint,omitempty"`
	Code        int    `json:"status_code,omitempty"`
}

func (e *RFC6749Error) Error() string {
	errstr, _ := jsoniter.Marshal(e)
	return string(errstr)
}

// WithHint 报错中添加hint
func (e *RFC6749Error) withHint(hint string) *RFC6749Error {
	err := *e
	err.Hint = hint

	return &err
}
