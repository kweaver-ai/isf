// Package session 协议层
package session

import (
	"encoding/json"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	"Authentication/logics/session"
)

// RESTHandler RESTful api Handler接口
type RESTHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type restHandler struct {
	session interfaces.Session
}

var (
	once sync.Once
	r    RESTHandler
)

type jsonValueDesc = rest.JSONValueDesc

// NewRESTHandler 创建context handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		r = &restHandler{
			session: session.NewSession(),
		}
	})

	return r
}

// RegisterPublic 注册内部API
func (r *restHandler) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/authentication/v1/session/:session_id", r.getSession)
	engine.PUT("/api/authentication/v1/session/:session_id", r.putSession)
	engine.DELETE("/api/authentication/v1/session/:session_id", r.deleteSession)
}

func (r *restHandler) getSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	context, err := r.session.Get(sessionID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	ctxInfo := make(map[string]interface{})
	if err := json.Unmarshal([]byte(context.Context), &ctxInfo); err != nil {
		rest.ReplyError(c, err)
		return
	}

	resInfo := map[string]interface{}{
		"subject":    context.Subject,
		"client_id":  context.ClientID,
		"session_id": context.SessionID,
		"context":    ctxInfo,
	}

	// 响应200
	rest.ReplyOK(c, http.StatusOK, resInfo)
}

func (r *restHandler) putSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["subject"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["client_id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["remember_for"] = &jsonValueDesc{Kind: reflect.Float64, Required: true}
	objDesc["context"] = &jsonValueDesc{Kind: reflect.Map, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	info := interfaces.Context{}
	jsonObj := jsonV.(map[string]interface{})
	jsonStr := make(map[string]string)
	jsonStr["subject"] = jsonObj["subject"].(string)
	jsonStr["client_id"] = jsonObj["client_id"].(string)
	rememberFor := int64(jsonObj["remember_for"].(float64))
	info.Exp = time.Now().UnixNano() + rememberFor*1e9
	info.SessionID = sessionID
	context, err := json.Marshal(jsonObj["context"])
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	info.Context = string(context)

	err = util.IsEmpty(jsonStr)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	info.Subject = jsonStr["subject"]
	info.ClientID = jsonStr["client_id"]

	err = r.session.Put(info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	c.Writer.Header().Set("Location", "api/authentication/v1/session/"+sessionID)
	rest.ReplyOK(c, http.StatusCreated, nil)
}

func (r *restHandler) deleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	err := r.session.Delete(sessionID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}
