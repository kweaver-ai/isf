// Package driveradapters app AnyShare  应用账户逻辑接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// AppRestHandler RESTful api Handler接口
type AppRestHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)

	// RegisterPrivate 注册开放API
	RegisterPrivate(engine *gin.Engine)
}

type appRestHandler struct {
	hydra                  interfaces.Hydra
	app                    interfaces.LogicsApp
	user                   interfaces.LogicsUser
	appType                map[string]interfaces.AppType
	credentialType         map[interfaces.CredentialType]string
	appTokenGenerateSchema *gojsonschema.Schema
}

var (
	aonce    sync.Once
	aHandler AppRestHandler

	//go:embed jsonschema/app/app-token-generate.json
	appTokenGenerateSchemaStr string
)

// NewAppRESTHandler 创建应用账户操作对象
func NewAppRESTHandler() AppRestHandler {
	aonce.Do(func() {
		appTokenGenerateSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(appTokenGenerateSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		aHandler = &appRestHandler{
			hydra: newHydra(),
			app:   logics.NewApp(),
			user:  logics.NewUser(),
			appType: map[string]interfaces.AppType{
				"specified": interfaces.Specified,
				"internal":  interfaces.Internal,
			},
			credentialType: map[interfaces.CredentialType]string{
				interfaces.CredentialTypePassword: "password",
				interfaces.CredentialTypeToken:    "token",
			},
			appTokenGenerateSchema: appTokenGenerateSchema,
		}
	})

	return aHandler
}

// RegisterPublic 注册开放API
func (a *appRestHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/api/user-management/v1/apps", a.generalAppRegister)
	engine.DELETE("/api/user-management/v1/apps/:id", a.deleteApp)
	engine.PUT("/api/user-management/v1/apps/:id/:fields", a.updateApp)
	engine.GET("/api/user-management/v1/apps", a.getAppList)
	engine.POST("/api/user-management/v1/console/app-tokens", observable.MiddlewareTrace(common.SvcARTrace), a.generateAppToken)
}

// RegisterPrivate 注册内部API
func (a *appRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/user-management/v1/apps", a.specifiedAppRegister)
	engine.GET("/api/user-management/v1/apps/:id", a.getApp)
}

func (a *appRestHandler) generateAppToken(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, a.hydra)
	if vErr != nil {
		rest.ReplyErrorV2(c, vErr)
		return
	}

	// 获取请求参数
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, a.appTokenGenerateSchema, &jsonReq); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取应用账户ID
	appID := jsonReq["id"].(string)

	token, err := a.app.GenerateAppToken(c, &visitor, appID)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 响应200
	rest.ReplyOK(c, http.StatusOK, gin.H{"token": token})
}

func (a *appRestHandler) generalAppRegister(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, a.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	var jsonV map[string]interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["name"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	name := jsonV["name"].(string)
	pwd := jsonV["password"].(string)

	id, err := a.app.RegisterApp(&visitor, name, pwd, interfaces.General)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应201
	c.Writer.Header().Set("Location", "/api/user-management/v1/apps/"+id)
	rest.ReplyOK(c, http.StatusCreated, gin.H{"id": id})
}

func (a *appRestHandler) deleteApp(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, a.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取客户端id
	id := c.Param("id")

	err := a.app.DeleteApp(&visitor, id)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (a *appRestHandler) updateApp(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, a.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取修改参数
	id := c.Param("id")
	fields := strings.Split(c.Param("fields"), ",")
	var nameExist, pwdExist bool
	for _, v := range fields {
		if v == strName {
			nameExist = true
		} else if v == "password" {
			pwdExist = true
		} else {
			err := rest.NewHTTPError("Invalid fields", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取请求参数bool
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数
	objDesc := make(map[string]*jsonValueDesc)
	if nameExist {
		objDesc["name"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	}

	if pwdExist {
		objDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数合法
	var err error
	jsonObj := jsonV.(map[string]interface{})
	var name, pwd string
	if nameExist {
		name = jsonObj["name"].(string)
	}
	if pwdExist {
		pwd = jsonObj["password"].(string)
	}

	err = a.app.UpdateApp(&visitor, id, nameExist, name, pwdExist, pwd)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (a *appRestHandler) getAppList(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, a.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	queryInfo, err := getListQueryParam(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	searchInfo := interfaces.SearchInfo{
		Direction: queryInfo.direction,
		Sort:      queryInfo.sort,
		Keyword:   queryInfo.keyword,
		Offset:    queryInfo.offset,
		Limit:     queryInfo.limit,
	}

	info, num, err := a.app.AppList(&visitor, &searchInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	var outInfo ListInfo
	outInfo.TotalCount = num
	outInfo.Entries = make([]interface{}, 0)

	for _, v := range *info {
		temp := map[string]interface{}{
			"id":              v.ID,
			"name":            v.Name,
			"credential_type": a.credentialType[v.CredentialType],
		}
		outInfo.Entries = append(outInfo.Entries, temp)
	}

	// 响应200
	rest.ReplyOK(c, http.StatusOK, outInfo)
}

func (a *appRestHandler) specifiedAppRegister(c *gin.Context) {
	// 获取请求参数
	var jsonV map[string]interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["name"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["type"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	name := jsonV["name"].(string)
	if name == "" {
		err := rest.NewHTTPError("Invalid name", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}
	pwd := jsonV["password"].(string)

	appType := jsonV["type"].(string)
	if _, ok := a.appType[appType]; !ok {
		err := rest.NewHTTPError("Invalid type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	id, err := a.app.RegisterApp(nil, name, pwd, a.appType[appType])
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应201
	c.Writer.Header().Set("Location", "/api/user-management/v1/apps/"+id)
	rest.ReplyOK(c, http.StatusCreated, gin.H{"id": id})
}

func (a *appRestHandler) getApp(c *gin.Context) {
	// 获取客户端id
	id := c.Param("id")

	info, err := a.app.GetApp(id)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	rest.ReplyOK(c, http.StatusOK, info)
}
