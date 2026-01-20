// Package sms 协议层
package sms

import (
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	"Authentication/driveradapters"
	smsSchema "Authentication/driveradapters/jsonschema/sms_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/sms"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)
}

type restHandler struct {
	aSMS       interfaces.LogicsAnonymousSMS
	aSMSSchema *gojsonschema.Schema
	redis      interfaces.RedisConn
	logger     common.Logger
}

var (
	rOnce sync.Once
	r     RESTHandler
)

// NewRESTHandler 创建sms RESTHandler对象
func NewRESTHandler() RESTHandler {
	rOnce.Do(func() {
		aSMSSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(smsSchema.SMSSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		handler := &restHandler{
			aSMS:       sms.NewAnonymousSMS(),
			aSMSSchema: aSMSSchema,
			redis:      common.NewRedisConn(),
			logger:     common.NewLogger(),
		}
		handler.redisSub()
		r = handler
	})

	return r
}

// RegisterPublic 注册开放API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/anonymous-sms-vcode", observable.MiddlewareTrace(common.SvcARTrace), r.createAndSendAnonymousSMSCode)
}

func (r *restHandler) createAndSendAnonymousSMSCode(c *gin.Context) {
	visitor := interfaces.Visitor{
		Language:      driveradapters.GetXLang(c),
		ErrorCodeType: util.GetErrorCodeType(c),
	}

	var err error
	var reqJSON map[string]interface{}
	if err = util.ValidateAndBindGin(c, r.aSMSSchema, &reqJSON); err != nil {
		rest.ReplyError(c, err)
		return
	}
	phoneNumber := reqJSON["phone_number"].(string)
	anonymityID := reqJSON["account"].(string)

	vcodeID, err := r.aSMS.CreateAndSendVCode(c, &visitor, phoneNumber, anonymityID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := make(map[string]interface{})
	resInfo["vcode_id"] = vcodeID
	// 响应200
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) redisSub() {
	r.redis.Subscribe("authentication.config.anonymous_sms_expiration.updated", r.anonymousSmsExpUpdated)
}

func (r *restHandler) anonymousSmsExpUpdated(message []byte) {
	var payload map[string]interface{}
	if err := jsoniter.Unmarshal(message, &payload); err != nil {
		r.logger.Errorln("failed to unmarshal anonymous sms expiration, err:", err)
		return
	}
	r.aSMS.UpdateSMSExpiration(int(payload["anonymous_sms_expiration"].(float64)))
}
