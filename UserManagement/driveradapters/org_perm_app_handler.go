// Package driveradapters 应用账户组织管理权限设置
package driveradapters

import (
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"UserManagement/interfaces"
	"UserManagement/logics"
)

// OrgPermAppHandler 应用账户组织管理权限接口配置
type OrgPermAppHandler interface {
	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}

type orgPermAppHandler struct {
	hydra             interfaces.Hydra
	orgPermApp        interfaces.LogicsOrgPermApp
	appPermOrgTypeStr map[interfaces.OrgType]string
	appOrgPermTypeStr map[interfaces.AppOrgPermValue]string
	appPermOrgStrType map[string]interfaces.OrgType
	appOrgPermStrType map[string]interfaces.AppOrgPermValue
}

var (
	oaonce    sync.Once
	oaHandler OrgPermAppHandler
)

// NewOrgPermAppHandler 新建应用账户组织架构权限对象
func NewOrgPermAppHandler() OrgPermAppHandler {
	oaonce.Do(func() {
		oaHandler = &orgPermAppHandler{
			hydra:      newHydra(),
			orgPermApp: logics.NewOrgPermApp(),
			appPermOrgTypeStr: map[interfaces.OrgType]string{
				interfaces.User:       "user",
				interfaces.Department: "department",
				interfaces.Group:      "group",
			},
			appOrgPermTypeStr: map[interfaces.AppOrgPermValue]string{
				interfaces.Modify: "modify",
				interfaces.Read:   "read",
			},
			appPermOrgStrType: map[string]interfaces.OrgType{
				"user":       interfaces.User,
				"department": interfaces.Department,
				"group":      interfaces.Group,
			},
			appOrgPermStrType: map[string]interfaces.AppOrgPermValue{
				"modify": interfaces.Modify,
				"read":   interfaces.Read,
			},
		}
	})

	return oaHandler
}

// RegisterPublic 注册开放API
func (o *orgPermAppHandler) RegisterPublic(engine *gin.Engine) {
	engine.GET("/api/user-management/v1/app-perms/:subject/:objects", o.getAppOrgPerm)
	engine.PUT("/api/user-management/v1/app-perms/:subject/:objects", o.updateAppOrgPerm)
	engine.DELETE("/api/user-management/v1/app-perms/:subject/:objects", o.deleteAppOrgPerm)
}

func (o *orgPermAppHandler) getAppOrgPerm(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, o.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取 subject和objects
	strSubject := c.Param("subject")
	strObjects := strings.Split(c.Param("objects"), ",")

	// 检查object类型
	objects := make([]interfaces.OrgType, 0, len(strObjects))
	for _, v := range strObjects {
		if object, ok := o.appPermOrgStrType[v]; ok {
			objects = append(objects, object)
		} else {
			err := rest.NewHTTPError("objects error", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取app 权限信息
	infos, err := o.orgPermApp.GetAppOrgPerm(&visitor, strSubject, objects)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回信息
	outData := make([]interface{}, len(infos))
	for k, v := range infos {
		temp := make(map[string]interface{})
		temp["subject"] = v.Subject
		temp["object"] = o.appPermOrgTypeStr[v.Object]
		temp["perms"] = o.permIntToArray(v.Value)

		outData[k] = temp
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

func (o *orgPermAppHandler) updateAppOrgPerm(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, o.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取app id数组
	subject := c.Param("subject")
	objects := strings.Split(c.Param("objects"), ",")

	urlObjects := make(map[interfaces.OrgType]bool)
	for _, v := range objects {
		temp, ok := o.appPermOrgStrType[v]
		if !ok {
			err := rest.NewHTTPError("url objects error", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		urlObjects[temp] = true
	}

	if len(urlObjects) != len(objects) {
		err := rest.NewHTTPError("url objects are not uniqued", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 获取请求参数
	var jsonV []interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	objDesc := make(map[string]*jsonValueDesc)
	permObjDesc := make(map[string]*jsonValueDesc)
	sliceArr := make(map[string]*jsonValueDesc)
	sliceArr["element"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	permObjDesc["subject"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	permObjDesc["object"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	permObjDesc["perms"] = &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: sliceArr}
	objDesc["element"] = &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: permObjDesc}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: objDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查权限参数有效性
	mapTypes := make(map[interfaces.OrgType]bool)
	infos := make([]interfaces.AppOrgPerm, 0, len(jsonV))
	var ok bool
	var err error
	for _, v := range jsonV {
		tempInter := v.(map[string]interface{})
		var temp interfaces.AppOrgPerm

		temp.Subject = tempInter["subject"].(string)
		if temp.Subject != subject {
			err = rest.NewHTTPError("request subject error", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		temp.Object, ok = o.appPermOrgStrType[tempInter["object"].(string)]
		if !ok {
			err = rest.NewHTTPError("request body object error", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		mapTypes[temp.Object] = true
		tempPerms := tempInter["perms"].([]interface{})
		if len(tempPerms) == 0 {
			err = rest.NewHTTPError("request body perms are empty", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}

		temp.Value, err = o.permArrayToInt(tempPerms)
		if err != nil {
			rest.ReplyError(c, err)
			return
		}
		infos = append(infos, temp)
	}

	// 检查url内subject和object与body内subject和object是否匹配
	if len(urlObjects) != len(jsonV) {
		err := rest.NewHTTPError("request body objects error", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	for k := range urlObjects {
		if _, ok := mapTypes[k]; !ok {
			err := rest.NewHTTPError("request body objects are not same with url objects", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 设置app 权限信息
	err = o.orgPermApp.SetAppOrgPerm(&visitor, subject, infos)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *orgPermAppHandler) deleteAppOrgPerm(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, o.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取app id数组
	subject := c.Param("subject")
	objects := strings.Split(c.Param("objects"), ",")

	types := make([]interfaces.OrgType, 0, len(objects))
	for _, v := range objects {
		if typ, ok := o.appPermOrgStrType[v]; ok {
			types = append(types, typ)
		} else {
			err := rest.NewHTTPError("objects error", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 删除app 权限信息
	err := o.orgPermApp.DeleteAppOrgPerm(&visitor, subject, types)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *orgPermAppHandler) permIntToArray(value interfaces.AppOrgPermValue) (out []string) {
	out = make([]string, 0)
	for k, v := range o.appOrgPermStrType {
		if v&value != 0 {
			out = append(out, k)
		}
	}
	return
}

func (o *orgPermAppHandler) permArrayToInt(perms []interface{}) (value interfaces.AppOrgPermValue, err error) {
	for _, v := range perms {
		tempValue, ok := o.appOrgPermStrType[v.(string)]
		if !ok {
			err = rest.NewHTTPError("request body perms error", rest.BadRequest, nil)
			return
		}
		value |= tempValue
	}
	return
}
