// Package audit 协议层
package audit

import (
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	auditschema "Authentication/driveradapters/jsonschema/audit_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/audit"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type restHandler struct {
	logSchema          *gojsonschema.Schema
	unorderedLogSchema *gojsonschema.Schema
	audit              interfaces.LogicsAudit
	auditLogAsyncTask  interfaces.LogicsAuditLogAsyncTask
}

var (
	once sync.Once
	r    RESTHandler
)

// NewRESTHandler 创建 handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		logSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditschema.AuditLogSchema))
		unorderedLogSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditschema.UnorderedAuditLogSchema))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		r = &restHandler{
			logSchema:          logSchema,
			unorderedLogSchema: unorderedLogSchema,
			audit:              audit.NewAudit(),
			auditLogAsyncTask:  audit.NewAuditLogAsyncTask(),
		}
	})

	return r
}

// RegisterPublic 注册开放API
func (r *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/authentication/v1/audit-log", r.log)

	engine.POST("/api/authentication/v2/audit-log", r.unorderedLog)
}

// log 记录审计日志
func (r *restHandler) log(c *gin.Context) {
	// 检查请求参数与文档是否匹配
	var jsonV map[string]interface{}
	if err := util.ValidateAndBindGin(c, r.logSchema, &jsonV); err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 检查参数是否合法
	topic := jsonV["topic"].(string)
	if topic == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid topic"))
		return
	}

	msg := jsonV["message"]
	err := r.audit.Log(topic, msg)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// unorderedLog 记录审计日志
func (r *restHandler) unorderedLog(c *gin.Context) {
	// 检查请求参数与文档是否匹配
	var jsonV map[string]interface{}
	if err := util.ValidateAndBindGin(c, r.unorderedLogSchema, &jsonV); err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 检查参数是否合法
	topic := jsonV["topic"].(string)
	if topic == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid topic"))
		return
	}

	msg := jsonV["message"]
	err := r.auditLogAsyncTask.Log(topic, msg)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}
