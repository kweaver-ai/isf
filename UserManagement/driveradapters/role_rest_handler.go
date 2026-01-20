// Package driveradapters AnyShare 公共接口处理层
package driveradapters

import (
	"net/http"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// RoleRestHandler driveradapters 角色 RESTfual API Handler 接口
type RoleRestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type roleRestHandler struct {
	roleNameIDMap map[string]interfaces.Role
	roleIDNameMap map[interfaces.Role]string
	role          interfaces.LogicsRole
}

var (
	ronce    sync.Once
	rhandler RoleRestHandler
)

// NewRoleRestHandler 角色 restful api handler 对象
func NewRoleRestHandler() RoleRestHandler {
	ronce.Do(func() {
		rhandler = &roleRestHandler{
			roleNameIDMap: map[string]interfaces.Role{
				EnumSuperAdmin: interfaces.SystemRoleSuperAdmin,
				EnumSysAdmin:   interfaces.SystemRoleSysAdmin,
				EnumAuditAdmin: interfaces.SystemRoleAuditAdmin,
				EnumSecAdmin:   interfaces.SystemRoleSecAdmin,
				EnumOrgManager: interfaces.SystemRoleOrgManager,
				EnumOrgAudit:   interfaces.SystemRoleOrgAudit,
			},
			roleIDNameMap: map[interfaces.Role]string{
				interfaces.SystemRoleSuperAdmin: EnumSuperAdmin,
				interfaces.SystemRoleSysAdmin:   EnumSysAdmin,
				interfaces.SystemRoleAuditAdmin: EnumAuditAdmin,
				interfaces.SystemRoleSecAdmin:   EnumSecAdmin,
				interfaces.SystemRoleOrgManager: EnumOrgManager,
				interfaces.SystemRoleOrgAudit:   EnumOrgAudit,
			},
			role: logics.NewRole(),
		}
	})

	return rhandler
}

// RegisterPrivate 注册内部API
func (ro *roleRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/user-management/v1/role-members/:roles", observable.MiddlewareTrace(common.SvcARTrace), ro.getUserIDsByRoleIDs)
}

// getUserIDsByRoleIDs 根据角色ID获取角色成员
func (ro *roleRestHandler) getUserIDsByRoleIDs(ctx *gin.Context) {
	// 检查角色字符串
	roleStrs := strings.Split(ctx.Param("roles"), ",")
	roles := make([]interfaces.Role, 0, len(roleStrs))
	for _, data := range roleStrs {
		if _, ok := ro.roleNameIDMap[data]; !ok {
			err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid role")
			rest.ReplyError(ctx, err)
			return
		}
		roles = append(roles, ro.roleNameIDMap[data])
	}

	// 获取角色成员
	roleData, err := ro.role.GetUserIDsByRoleIDs(ctx, roles)
	if err != nil {
		rest.ReplyError(ctx, err)
		return
	}

	// 配置信息
	out := make([]interface{}, 0, len(roleData))
	tempData := make(map[interfaces.Role]bool)
	for _, v := range roles {
		if _, ok := tempData[v]; ok {
			continue
		}
		tempData[v] = true

		roleInfo := make(map[string]interface{})
		roleInfo["role"] = ro.roleIDNameMap[v]

		members := make([]interface{}, 0, len(roleData[v]))
		for _, userID := range roleData[v] {
			members = append(members, map[string]interface{}{
				"id":   userID,
				"type": strUser,
			})
		}
		roleInfo["members"] = members
		out = append(out, roleInfo)
	}

	rest.ReplyOK(ctx, http.StatusOK, out)
}
