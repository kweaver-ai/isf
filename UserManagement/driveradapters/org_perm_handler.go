// Package driveradapters 应用账户组织管理权限设置
package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
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

// OrgPermHandler 账户组织管理权限接口配置
type OrgPermHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type orgPermHandler struct {
	orgPerm        interfaces.LogicsOrgPerm
	permOrgStrType map[string]interfaces.OrgType
	orgPermStrType map[string]interfaces.OrgPermValue
	subTypeStrType map[string]interfaces.VisitorType

	updateSchema *gojsonschema.Schema
}

var (
	oponce    sync.Once
	opHandler OrgPermHandler

	//go:embed jsonschema/org_perm/update.json
	updateSchemaStr string
)

// NewOrgPermApHandler 新建应用账户组织架构权限对象
func NewOrgPermApHandler() OrgPermHandler {
	updateSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateSchemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}

	oponce.Do(func() {
		opHandler = &orgPermHandler{
			permOrgStrType: map[string]interfaces.OrgType{
				"user":       interfaces.User,
				"department": interfaces.Department,
				"group":      interfaces.Group,
			},
			orgPermStrType: map[string]interfaces.OrgPermValue{
				"read": interfaces.OPRead,
			},
			subTypeStrType: map[string]interfaces.VisitorType{
				"user": interfaces.RealName,
			},
			updateSchema: updateSchema,
			orgPerm:      logics.NewOrgPerm(),
		}
	})

	return opHandler
}

// RegisterPublic 注册开放API
func (o *orgPermHandler) RegisterPrivate(engine *gin.Engine) {
	engine.DELETE("/api/user-management/v1/org-perm/:subject_type/:subject_id/:objects", observable.MiddlewareTrace(common.SvcARTrace), o.deleteOrgPerm)
	engine.PUT("/api/user-management/v1/org-perm/:subject_type/:subject_id/:objects", observable.MiddlewareTrace(common.SvcARTrace), o.updateOrgPerm)
}

func (o *orgPermHandler) updateOrgPerm(c *gin.Context) {
	// 获取id数组
	subjectTypeStr := c.Param("subject_type")
	subject := c.Param("subject_id")
	objects := strings.Split(c.Param("objects"), ",")

	// subject type 检查
	var subjectType interfaces.VisitorType
	var ok bool
	if subjectType, ok = o.subTypeStrType[subjectTypeStr]; !ok {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "url subject_type error")
		rest.ReplyError(c, err)
		return
	}

	// object 检查
	urlObjects := make(map[interfaces.OrgType]bool)
	for _, v := range objects {
		temp, ok := o.permOrgStrType[v]
		if !ok {
			err := rest.NewHTTPErrorV2(rest.BadRequest, "url objects error")
			rest.ReplyError(c, err)
			return
		}
		urlObjects[temp] = true
	}

	if len(urlObjects) != len(objects) {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "url objects are not uniqued")
		rest.ReplyError(c, err)
		return
	}

	// jsonschema校验
	var jsonReq []interface{}
	if err := validateAndBindGin(c, o.updateSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	infos := make([]interfaces.OrgPerm, 0, len(jsonReq))
	jsonObjects := make(map[interfaces.OrgType]bool)
	var err error
	for _, v := range jsonReq {
		var tempInfo interfaces.OrgPerm
		temp := v.(map[string]interface{})
		tempInfo.Object = o.permOrgStrType[temp["object"].(string)]

		tempPerms := temp["perms"].([]interface{})
		tempInfo.Value, err = o.permArrayToInt(tempPerms)
		if err != nil {
			rest.ReplyError(c, err)
			return
		}

		tempInfo.EndTime = -1
		tempInfo.SubjectID = subject
		tempInfo.SubjectType = subjectType
		infos = append(infos, tempInfo)

		jsonObjects[tempInfo.Object] = true
	}

	// 检查url object和json object是否匹配
	result := o.areMapsEqual(urlObjects, jsonObjects)
	if !result {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "url object not equal object object"))
		return
	}

	// 设置 权限信息
	err = o.orgPerm.SetOrgPerm(c, subject, subjectType, infos)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *orgPermHandler) deleteOrgPerm(c *gin.Context) {
	// 获取 id数组
	subjectTypeStr := c.Param("subject_type")
	subject := c.Param("subject_id")
	objects := strings.Split(c.Param("objects"), ",")

	// subject type 检查
	var subjectType interfaces.VisitorType
	var ok bool
	if subjectType, ok = o.subTypeStrType[subjectTypeStr]; !ok {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "url subject_type error")
		rest.ReplyError(c, err)
		return
	}

	types := make([]interfaces.OrgType, 0, len(objects))
	for _, v := range objects {
		if typ, ok := o.permOrgStrType[v]; ok {
			types = append(types, typ)
		} else {
			err := rest.NewHTTPErrorV2(rest.BadRequest, "url objects error")
			rest.ReplyError(c, err)
			return
		}
	}

	// 删除 权限信息
	err := o.orgPerm.DeleteOrgPerm(c, subject, subjectType, types)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *orgPermHandler) permArrayToInt(perms []interface{}) (value interfaces.OrgPermValue, err error) {
	for _, v := range perms {
		tempValue, ok := o.orgPermStrType[v.(string)]
		if !ok {
			err = rest.NewHTTPErrorV2(rest.BadRequest, "request body perms error")
			return
		}
		value |= tempValue
	}
	return
}

func (o *orgPermHandler) areMapsEqual(map1, map2 map[interfaces.OrgType]bool) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value := range map1 {
		if map2Value, ok := map2[key]; !ok || map2Value != value {
			return false
		}
	}

	return true
}
