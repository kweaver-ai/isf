// Package logics combine Anyshare 业务逻辑层 -组合
package logics

import (
	"sync"

	"Authorization/interfaces"
)

var (
	eOnce sync.Once
	e     *event
)

type event struct {
	userDeletedHandlers       []func(string) error
	depDeletedHandlers        []func(string) error
	userGroupDeletedHandlers  []func(string) error
	orgNameModifiedHandlers   []func(string, string) error
	userNameModifiedHandlers  []func(string, string) error
	depNameModifiedHandlers   []func(string, string) error
	groupNameModifiedHandlers []func(string, string) error
	appDeletedHandlers        []func(string) error
	appNameModifiedHandlers   []func(*interfaces.AppInfo) error
	roleDeletedHandlers       []func(string) error
	roleNameModifiedHandlers  []func(string, string) error
}

// NewEvent 创建新的event对象
func NewEvent() *event {
	eOnce.Do(func() {
		e = &event{
			userDeletedHandlers:       make([]func(string) error, 0),
			depDeletedHandlers:        make([]func(string) error, 0),
			userGroupDeletedHandlers:  make([]func(string) error, 0),
			orgNameModifiedHandlers:   make([]func(string, string) error, 0),
			userNameModifiedHandlers:  make([]func(string, string) error, 0),
			depNameModifiedHandlers:   make([]func(string, string) error, 0),
			groupNameModifiedHandlers: make([]func(string, string) error, 0),
			appDeletedHandlers:        make([]func(string) error, 0),
			appNameModifiedHandlers:   make([]func(*interfaces.AppInfo) error, 0),
			roleDeletedHandlers:       make([]func(string) error, 0),
			roleNameModifiedHandlers:  make([]func(string, string) error, 0),
		}
	})
	return e
}

// UserDeleted 用户被删除
func (e *event) UserDeleted(userID string) (err error) {
	for _, f := range e.userDeletedHandlers {
		err = f(userID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterUserDeleted 注册用户被删除
func (e *event) RegisterUserDeleted(f func(string) error) {
	e.userDeletedHandlers = append(e.userDeletedHandlers, f)
}

// DepartmentDeleted 部门被删除
func (e *event) DepartmentDeleted(depID string) (err error) {
	for _, f := range e.depDeletedHandlers {
		err = f(depID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterDepartmentDeleted 注册部门被删除
func (e *event) RegisterDepartmentDeleted(f func(string) error) {
	e.depDeletedHandlers = append(e.depDeletedHandlers, f)
}

// UserGroupDeleted 用户组被删除
func (e *event) UserGroupDeleted(groupID string) (err error) {
	for _, f := range e.userGroupDeletedHandlers {
		err = f(groupID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterUserGroupDeleted 注册用户组被删除
func (e *event) RegisterUserGroupDeleted(f func(string) error) {
	e.userGroupDeletedHandlers = append(e.userGroupDeletedHandlers, f)
}

// AppDeleted 删除应用账户权限(文档权限和文档库权限)和所有者
func (e *event) AppDeleted(appID string) (err error) {
	for _, f := range e.appDeletedHandlers {
		err = f(appID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterAppDeleted 注册删除应用账户权限(文档权限和文档库权限)和所有者
func (e *event) RegisterAppDeleted(f func(string) error) {
	e.appDeletedHandlers = append(e.appDeletedHandlers, f)
}

// AppNameModified 更新应用账户名称(文档和文档库权限表、所有者表)
func (e *event) AppNameModified(info *interfaces.AppInfo) (err error) {
	for _, f := range e.appNameModifiedHandlers {
		err = f(info)
		if err != nil {
			return
		}
	}
	return
}

// RegisterAppNameModified 注册更新应用账户名称
func (e *event) RegisterAppNameModified(f func(*interfaces.AppInfo) error) {
	e.appNameModifiedHandlers = append(e.appNameModifiedHandlers, f)
}

// OrgNameModified 更新组织架构显示名, 不区分类型使用 RegisterOrgNameModified。区分类型使用具体类型如 RegisterUserNameModified
func (e *event) OrgNameModified(id, name string, orgType interfaces.AccessorType) (err error) {
	for _, f := range map[interfaces.AccessorType][]func(string, string) error{
		interfaces.AccessorUser:       e.userNameModifiedHandlers,
		interfaces.AccessorDepartment: e.depNameModifiedHandlers,
		interfaces.AccessorGroup:      e.groupNameModifiedHandlers,
	}[orgType] {
		err = f(id, name)
		if err != nil {
			return
		}
	}

	for _, f := range e.orgNameModifiedHandlers {
		err = f(id, name)
		if err != nil {
			return
		}
	}
	return
}

// RegisterOrgNameModified 注册组织架构显示名变更
func (e *event) RegisterOrgNameModified(f func(string, string) error) {
	e.orgNameModifiedHandlers = append(e.orgNameModifiedHandlers, f)
}

// RegisterUserNameModified 注册用户名称变更
func (e *event) RegisterUserNameModified(f func(string, string) error) {
	e.userNameModifiedHandlers = append(e.userNameModifiedHandlers, f)
}

// RegisterDepartmentNameModified 注册部门名称变更
func (e *event) RegisterDepartmentNameModified(f func(string, string) error) {
	e.depNameModifiedHandlers = append(e.depNameModifiedHandlers, f)
}

// RegisterGroupNameModified 注册组名称变更
func (e *event) RegisterUserGroupNameModified(f func(string, string) error) {
	e.groupNameModifiedHandlers = append(e.groupNameModifiedHandlers, f)
}

/*
	Role
*/
// RoleDeleted 角色被删除
func (e *event) RoleDeleted(roleID string) (err error) {
	for _, f := range e.roleDeletedHandlers {
		err = f(roleID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterRoleDeleted 注册角色被删除
func (e *event) RegisterRoleDeleted(f func(string) error) {
	e.roleDeletedHandlers = append(e.roleDeletedHandlers, f)
}

/*
	Role
*/
// RoleNameModified 更新角色名称
func (e *event) RoleNameModified(roleID, name string) (err error) {
	for _, f := range e.roleNameModifiedHandlers {
		err = f(roleID, name)
		if err != nil {
			return
		}
	}
	return
}

// RegisterRoleNameModified 注册角色名称更新
func (e *event) RegisterRoleNameModified(f func(string, string) error) {
	e.roleNameModifiedHandlers = append(e.roleNameModifiedHandlers, f)
}
