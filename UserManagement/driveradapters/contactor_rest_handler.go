// Package driveradapters contactor AnyShare  部门逻辑接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
	"reflect"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// ContactorRestHandler RESTful api Handler接口
type ContactorRestHandler interface {
	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)

	// RegisterPrivate 注册私有API
	RegisterPrivate(engine *gin.Engine)
}

type contactorRestHandler struct {
	contactor                     interfaces.LogicsContactor
	hydra                         interfaces.Hydra
	getContactorMembersPostSchema *gojsonschema.Schema
}

var (
	contactoronce sync.Once
	contactor     ContactorRestHandler
)

var (
	//go:embed jsonschema/contactor/get_contactor_members_post.json
	getContactorMembersPostSchemaStr string
)

// NewContactorRESTHandler 创建联系人组操作对象
func NewContactorRESTHandler() ContactorRestHandler {
	contactoronce.Do(func() {
		getContactorMembersPostSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getContactorMembersPostSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		contactor = &contactorRestHandler{
			contactor:                     logics.NewContactor(),
			hydra:                         newHydra(),
			getContactorMembersPostSchema: getContactorMembersPostSchema,
		}
	})

	return contactor
}

// RegisterPublic 注册开放API
func (con *contactorRestHandler) RegisterPublic(engine *gin.Engine) {
	engine.POST("/api/eacp/v1/contactor/deletegroup", con.deleteContactor)
}

// RegisterPrivate 注册私有API
func (con *contactorRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.POST("/api/user-management/v1/contactor-members", observable.MiddlewareTrace(common.SvcARTrace), con.getContactorMembers)
}

// getContactorMembers 获取联系人组成员
func (con *contactorRestHandler) getContactorMembers(c *gin.Context) {
	// jsonschema校验
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, con.getContactorMembersPostSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 检查method
	if jsonReq["method"].(string) != strHTTPGET {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid method"))
		return
	}

	// 获取所有的联系人组id
	contactorIDs := make([]string, 0)
	for _, v := range jsonReq["contactor_ids"].([]interface{}) {
		contactorIDs = append(contactorIDs, v.(string))
	}

	// 获取联系人组成员
	contactorMembers, err := con.contactor.GetContactorMembers(c, contactorIDs)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回信息
	outInfo := make([]interface{}, 0, len(contactorMembers))
	for _, v := range contactorMembers {
		out := make(map[string]interface{})
		out["contactor_id"] = v.ContactorID
		tmpMembers := make([]interface{}, 0, len(v.MemberIDs))
		for _, memberID := range v.MemberIDs {
			member := make(map[string]interface{})
			member["id"] = memberID
			member["type"] = strUser
			tmpMembers = append(tmpMembers, member)
		}
		out["members"] = tmpMembers
		outInfo = append(outInfo, out)
	}

	// 响应
	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// deleteContactor 删除联系人组
func (con *contactorRestHandler) deleteContactor(c *gin.Context) {
	// token 检查
	visitor, vErr := verify(c, con.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["groupid"] = &jsonValueDesc{Kind: reflect.String, Required: true}

	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 获取具体请求参数
	jsonObj := jsonV.(map[string]interface{})
	groupID := jsonObj["groupid"].(string)

	// 删除联系人组
	err := con.contactor.DeleteContactor(&visitor, groupID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 沿用EACP逻辑 响应200
	temp := make(map[string]interface{})
	rest.ReplyOK(c, http.StatusOK, temp)
}
