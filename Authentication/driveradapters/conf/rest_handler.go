// Package conf 协议层
package conf

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	configSchema "Authentication/driveradapters/jsonschema/config_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/conf"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type restHandler struct {
	conf            interfaces.Conf
	hydra           interfaces.Hydra
	keyConfMap      map[string]interfaces.ConfigKey
	resKeyConfMap   map[interfaces.ConfigKey]string
	setConfigSchema *gojsonschema.Schema
}

var (
	ronce sync.Once
	r     RESTHandler
)

// NewRESTHandler 创建context handler对象
func NewRESTHandler() RESTHandler {
	ronce.Do(func() {
		setConfigSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(configSchema.SetConfigSchemaStr))
		if err != nil {
			common.NewLogger().Errorf("failed to create setConfigSchema, err: %v", err)
			return
		}
		r = &restHandler{
			conf:  conf.NewConf(),
			hydra: util.NewHydra(),
			keyConfMap: map[string]interfaces.ConfigKey{
				"remember_for":             interfaces.RememberFor,
				"remember_visible":         interfaces.RememberVisible,
				"anonymous_sms_expiration": interfaces.SMSExpiration,
			},
			resKeyConfMap: map[interfaces.ConfigKey]string{
				interfaces.RememberFor:     "remember_for",
				interfaces.RememberVisible: "remember_visible",
				interfaces.SMSExpiration:   "anonymous_sms_expiration",
			},
			setConfigSchema: setConfigSchema,
		}
	})

	return r
}

// RegisterPublic 注册内部API
func (r *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/authentication/v1/config/:fields", observable.MiddlewareTrace(common.SvcARTrace), r.privateGetConfig)
	engine.PUT("/api/authentication/v1/config/:fields", observable.MiddlewareTrace(common.SvcARTrace), r.privateSetConfig)
}

// RegisterPublic 注册外部API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.GET("/api/authentication/v1/config/:fields", observable.MiddlewareTrace(common.SvcARTrace), r.publicGetConfig)
	engine.PUT("/api/authentication/v1/config/:fields", observable.MiddlewareTrace(common.SvcARTrace), r.publicSetConfig)
}

func (r *restHandler) privateGetConfig(c *gin.Context) {
	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	r.getConfig(c, &visitor)
}

func (r *restHandler) publicGetConfig(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, r.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	r.getConfig(c, &visitor)
}

func (r *restHandler) getConfig(c *gin.Context, visitor *interfaces.Visitor) {
	keyMap, err := r.convert(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	cfg, err := r.conf.GetConfig(c, visitor, keyMap)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := make(map[string]interface{})
	for key := range keyMap {
		switch key {
		case interfaces.RememberFor:
			resInfo[r.resKeyConfMap[key]] = cfg.RememberFor
		case interfaces.RememberVisible:
			resInfo[r.resKeyConfMap[key]] = cfg.RememberVisible
		case interfaces.SMSExpiration:
			resInfo[r.resKeyConfMap[key]] = cfg.SMSExpiration
		}
	}
	// 响应200
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) privateSetConfig(c *gin.Context) {
	visitor := interfaces.Visitor{
		ErrorCodeType: util.GetErrorCodeType(c),
	}
	r.setConfig(c, &visitor)
}

func (r *restHandler) publicSetConfig(c *gin.Context) {
	// token内省
	visitor, err := util.Verify(c, r.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	r.setConfig(c, &visitor)
}

func (r *restHandler) setConfig(c *gin.Context, visitor *interfaces.Visitor) {
	keyMap, err := r.convert(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	var payload map[string]interface{}
	err = util.ValidateAndBindGin(c, r.setConfigSchema, &payload)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	var cfg interfaces.Config
	for k := range keyMap {
		fieldName := r.resKeyConfMap[k]
		// 检查请求体参数是否包含path中定义的field
		if _, ok := payload[fieldName]; !ok {
			rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, fmt.Sprintf("missing required field '%s'", fieldName)))
			return
		}
		switch k {
		case interfaces.RememberFor:
			cfg.RememberFor = int(payload[fieldName].(float64))
		case interfaces.RememberVisible:
			cfg.RememberVisible = payload[fieldName].(bool)
		case interfaces.SMSExpiration:
			cfg.SMSExpiration = int(payload[fieldName].(float64))
		default:
		}
	}

	err = r.conf.SetConfig(c, visitor, keyMap, cfg)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (r *restHandler) convert(c *gin.Context) (keyMap map[interfaces.ConfigKey]bool, err error) {
	fields := strings.Split(c.Param("fields"), ",")

	keyMap = make(map[interfaces.ConfigKey]bool)
	for _, v := range fields {
		key, exist := r.keyConfMap[v]
		if !exist {
			return keyMap, rest.NewHTTPError("Invalid fields", rest.BadRequest, nil)
		}
		keyMap[key] = true
	}

	return
}
