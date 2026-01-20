// Package driveradapters anonymous AnyShare  匿名账户逻辑接口处理层
package driveradapters

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"UserManagement/interfaces"
	"UserManagement/logics"
)

// AnonymousRestHandler user RESTfual API Handler 接口
type AnonymousRestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type anonymousRestHandler struct {
	anonymous interfaces.LogicsAnonymous
}

var (
	anonyOnce    sync.Once
	anonyHandler AnonymousRestHandler
)

// NewAnonymousRestHandler handler 对象
func NewAnonymousRestHandler() AnonymousRestHandler {
	anonyOnce.Do(func() {
		anonyHandler = &anonymousRestHandler{
			anonymous: logics.NewAnonymous(),
		}
	})

	return anonyHandler
}

// RegisterPrivate 注册内部API
func (h *anonymousRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/user-management/v1/anonymity-auth", h.authenticaitonAnonymous)
	engine.GET("/api/user-management/v1/anonymity/:id", h.getAnonymous)
}

// authenticaitonAnonymous 认证匿名账户
func (h *anonymousRestHandler) authenticaitonAnonymous(c *gin.Context) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["account"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	jsonObj := jsonV.(map[string]interface{})
	account := jsonObj["account"].(string)
	password := jsonObj["password"].(string)

	referrer := c.Request.Header.Get("x-referrer")
	err := h.anonymous.Authentication(account, password, referrer)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	resp := map[string]interface{}{
		"result": true,
	}

	rest.ReplyOK(c, http.StatusOK, resp)
}

// getAnonymous 获取匿名账户信息
func (h *anonymousRestHandler) getAnonymous(c *gin.Context) {
	anonymityID := c.Param("id")

	info, err := h.anonymous.GetByID(anonymityID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	resInfo := make(map[string]interface{})
	resInfo["verify_mobile"] = info.VerifyMobile

	rest.ReplyOK(c, http.StatusOK, resInfo)
}
