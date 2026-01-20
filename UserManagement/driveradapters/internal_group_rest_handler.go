// Package driveradapters group AnyShare  内部组逻辑接口处理层
package driveradapters

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"UserManagement/interfaces"
	"UserManagement/logics"
)

// InternalGroupRestHandler driveradapters 内部组 RESTfual API Handler 接口
type InternalGroupRestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type internalGroupRestHandler struct {
	group           interfaces.LogicsInternalGroup
	orgTypeToString map[interfaces.OrgType]string
	stringToOrgType map[string]interfaces.OrgType
}

var (
	igOnce    sync.Once
	ighandler *internalGroupRestHandler
)

// NewInternalGroupRESTHandler 内部组 restful api handler 对象
func NewInternalGroupRESTHandler() InternalGroupRestHandler {
	igOnce.Do(func() {
		ighandler = &internalGroupRestHandler{
			group: logics.NewInternalGroup(),
			orgTypeToString: map[interfaces.OrgType]string{
				interfaces.User: "user",
			},
			stringToOrgType: map[string]interfaces.OrgType{
				"user": interfaces.User,
			},
		}
	})

	return ighandler
}

// RegisterPrivate 注册内部API
func (g *internalGroupRestHandler) RegisterPrivate(engine *gin.Engine) {
	// 组管理
	engine.POST("/api/user-management/v1/internal-groups", g.createGroup)
	engine.DELETE("/api/user-management/v1/internal-groups/:ids", g.deleteGroup)

	// 组成员
	engine.GET("/api/user-management/v1/internal-group-members/:id", g.getMembersByID)
	engine.PUT("/api/user-management/v1/internal-group-members/:id", g.updateMembers)
}

// createGroup 创建内部组
func (g *internalGroupRestHandler) createGroup(c *gin.Context) {
	strID, err := g.group.AddGroup()
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/user-management/v1/internal-groups/%s", strID))
	rest.ReplyOK(c, http.StatusCreated, gin.H{"id": strID})
}

// deleteGroup 删除内部组
func (g *internalGroupRestHandler) deleteGroup(c *gin.Context) {
	ids := strings.Split(c.Param("ids"), ",")

	err := g.group.DeleteGroup(ids)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// getMembersByID 获取组成员
func (g *internalGroupRestHandler) getMembersByID(c *gin.Context) {
	id := c.Param("id")

	outInfos, err := g.group.GetGroupMemberByID(id)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	out := make([]interface{}, 0, len(outInfos))
	for _, v := range outInfos {
		temp := make(map[string]interface{})
		temp["id"] = v.ID
		temp["type"] = g.orgTypeToString[v.Type]
		out = append(out, temp)
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// updateMembers 更新成员
func (g *internalGroupRestHandler) updateMembers(c *gin.Context) {
	id := c.Param("id")

	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	memberDesc := make(map[string]*jsonValueDesc)
	memberDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	memberDesc["type"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	membersDesc := make(map[string]*jsonValueDesc)
	membersDesc["element"] = &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: memberDesc}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: membersDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	jsonObj := jsonV.([]interface{})
	memberInfos := make([]interfaces.InternalGroupMember, 0)
	for _, v := range jsonObj {
		member := interfaces.InternalGroupMember{}
		temp := v.(map[string]interface{})
		member.ID = temp["id"].(string)
		strType := temp["type"].(string)
		var ok bool
		member.Type, ok = g.stringToOrgType[strType]
		if !ok {
			err := rest.NewHTTPError("param member type is illegal", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		if member.ID == "" {
			err := rest.NewHTTPError("param member id is illegal", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}

		memberInfos = append(memberInfos, member)
	}

	// 更新成员
	err := g.group.UpdateMembers(id, memberInfos)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}
