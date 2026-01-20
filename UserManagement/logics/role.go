// Package logics 角色管理
package logics

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type role struct {
	db       interfaces.DBRole
	trace    observable.Tracer
	departDB interfaces.DBDepartment
	userDB   interfaces.DBUser
}

var (
	rOnce   sync.Once
	rLogics *role
)

// NewRole 创建新的role对象
func NewRole() *role {
	rOnce.Do(func() {
		rLogics = &role{
			db:       dbRole,
			trace:    common.SvcARTrace,
			departDB: dbDepartment,
			userDB:   dbUser,
		}
	})

	return rLogics
}

// GetRolesByUserIDs 根据用户id数组获取用户角色id,包含普通用户
func (r *role) GetRolesByUserIDs(userIDs []string) (out map[string]map[interfaces.Role]bool, err error) {
	// 获取用户角色信息
	roleInfos, err := r.db.GetRolesByUserIDs(userIDs)
	if err != nil {
		return
	}

	if roleInfos == nil {
		roleInfos = make(map[string]map[interfaces.Role]bool)
	}

	// 如果用户不是管理员则增加普通用户角色
	for _, v := range userIDs {
		temp, ok := roleInfos[v]
		if !ok {
			temp = make(map[interfaces.Role]bool)
			roleInfos[v] = temp
		}

		if !AdminIDMap[v] {
			temp[interfaces.SystemRoleNormalUser] = true
		}
	}
	return roleInfos, err
}

// GetRolesByUserIDs2 根据用户id数组获取用户角色id,包含普通用户，支持trace
func (r *role) GetRolesByUserIDs2(ctx context.Context, userIDs []string) (out map[string]map[interfaces.Role]bool, err error) {
	// trace
	r.trace.SetInternalSpanName("业务逻辑-根据用户id数组获取用户角色id,包含普通用户")
	newCtx, span := r.trace.AddInternalTrace(ctx)
	defer func() { r.trace.TelemetrySpanEnd(span, err) }()

	// 获取用户角色信息
	roleInfos, err := r.db.GetRolesByUserIDs2(newCtx, userIDs)
	if err != nil {
		return
	}

	if roleInfos == nil {
		roleInfos = make(map[string]map[interfaces.Role]bool)
	}

	// 如果用户不是管理员则增加普通用户角色
	for _, v := range userIDs {
		temp, ok := roleInfos[v]
		if !ok {
			temp = make(map[interfaces.Role]bool)
			roleInfos[v] = temp
		}

		if !AdminIDMap[v] {
			temp[interfaces.SystemRoleNormalUser] = true
		}
	}
	return roleInfos, err
}

// GetOrgManagersInfo 获取组织管理员信息
func (r *role) GetOrgManagersInfo(orgIDs []string, rangeInfo interfaces.OrgManagerInfoRange) (out []interfaces.OrgManagerInfo, err error) {
	// 获取组织管理员下所有部门id
	if rangeInfo.ShowSubUserIDs {
		orgManagerDepInfos, err := r.userDB.GetOrgManagersDepartInfo(orgIDs)
		if err != nil {
			return out, err
		}

		depIDs := make([]string, 0)
		for _, v := range orgManagerDepInfos {
			depIDs = append(depIDs, v...)
		}

		// 获取所有部门的路径
		depInfos, err := r.departDB.GetDepartmentInfo(depIDs, false, 0, -1)
		if err != nil {
			return out, err
		}

		// 获取各个部门管理的所有子用户
		depSubUsersIDs := make(map[string][]string, 0)
		for k := range depInfos {
			var tempIDs []string
			tempIDs, err = r.departDB.GetAllSubUserIDsByDepartPath(depInfos[k].Path)
			if err != nil {
				return out, err
			}

			depSubUsersIDs[depInfos[k].ID] = tempIDs
		}

		// 获取组织管理员下所有用户
		for _, v := range orgIDs {
			var orgManagerInfo interfaces.OrgManagerInfo
			tempDepIDs := orgManagerDepInfos[v]
			for _, v1 := range tempDepIDs {
				orgManagerInfo.SubUserIDs = append(orgManagerInfo.SubUserIDs, depSubUsersIDs[v1]...)
			}
			// 去重
			orgManagerInfo.ID = v
			RemoveDuplicatStrs(&orgManagerInfo.SubUserIDs)
			out = append(out, orgManagerInfo)
		}
		return out, err
	}
	return
}

// 根据角色获取角色成员
func (r *role) GetUserIDsByRoleIDs(ctx context.Context, roles []interfaces.Role) (out map[interfaces.Role][]string, err error) {
	// trace
	r.trace.SetInternalSpanName("业务逻辑-根据角色获取角色成员")
	newCtx, span := r.trace.AddInternalTrace(ctx)
	defer func() { r.trace.TelemetrySpanEnd(span, err) }()

	// 去重
	mapRoles := make(map[interfaces.Role]bool)
	for _, v := range roles {
		mapRoles[v] = true
	}
	tempRoles := make([]interfaces.Role, 0)
	for k := range mapRoles {
		tempRoles = append(tempRoles, k)
	}

	return r.db.GetUserIDsByRoleIDs(newCtx, tempRoles)
}
