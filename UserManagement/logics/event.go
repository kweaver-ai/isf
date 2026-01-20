// Package logics event Anyshare 业务逻辑层 -事件
package logics

import (
	"sync"
)

var (
	eOnce sync.Once
	e     *event
)

type event struct {
	deptDeletedHandlers            []func(string) error
	userDeletedHandlers            []func(string) error
	userNameChangedHandlers        []func(string, string) error
	departResponserChangedHandlers []func([]string) error
}

// NewEvent 创建新的event对象
func NewEvent() *event {
	eOnce.Do(func() {
		e = &event{
			deptDeletedHandlers:            make([]func(string) error, 0),
			userDeletedHandlers:            make([]func(string) error, 0),
			departResponserChangedHandlers: make([]func([]string) error, 0),
			userNameChangedHandlers:        make([]func(string, string) error, 0),
		}
	})
	return e
}

// 部门删除
func (e *event) DeptDeleted(deptID string) (err error) {
	for _, f := range e.deptDeletedHandlers {
		err = f(deptID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterDeptDeleted
func (e *event) RegisterDeptDeleted(f func(string) error) {
	e.deptDeletedHandlers = append(e.deptDeletedHandlers, f)
}

// 部门管理员变更事件
func (e *event) OrgManagerChanged(userIDs []string) (err error) {
	for _, f := range e.departResponserChangedHandlers {
		err = f(userIDs)
		if err != nil {
			return
		}
	}
	return
}

// RegisterDepartResponserChanged
func (e *event) RegisterDepartResponserChanged(f func([]string) error) {
	e.departResponserChangedHandlers = append(e.departResponserChangedHandlers, f)
}

func (e *event) UserDeleted(userID string) (err error) {
	for _, f := range e.userDeletedHandlers {
		err = f(userID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterUserDeleted
func (e *event) RegisterUserDeleted(f func(string) error) {
	e.userDeletedHandlers = append(e.userDeletedHandlers, f)
}

// UserNameChanged 用户名变更
func (e *event) UserNameChanged(userID, newName string) (err error) {
	for _, f := range e.userNameChangedHandlers {
		err = f(userID, newName)
		if err != nil {
			return
		}
	}
	return
}

// RegisterUserNameChanged
func (e *event) RegisterUserNameChanged(f func(string, string) error) {
	e.userNameChangedHandlers = append(e.userNameChangedHandlers, f)
}
