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
	userCreatedHandlers       []func(string) error
	userStatusChangedHandlers []func(string, bool) error
	userDeletedHandlers       []func(string) error
}

// NewEvent 创建新的event对象
func NewEvent() *event {
	eOnce.Do(func() {
		e = &event{
			userCreatedHandlers:       make([]func(string) error, 0),
			userStatusChangedHandlers: make([]func(string, bool) error, 0),
			userDeletedHandlers:       make([]func(string) error, 0),
		}
	})
	return e
}

// 用户删除
func (e *event) UserCreated(userID string) (err error) {
	for _, f := range e.userCreatedHandlers {
		err = f(userID)
		if err != nil {
			return
		}
	}
	return
}

// RegisterDeptDeleted
func (e *event) RegisterUserCreated(f func(string) error) {
	e.userCreatedHandlers = append(e.userCreatedHandlers, f)
}

// 用户状态改变
func (e *event) UserStatusChanged(userID string, status bool) (err error) {
	for _, f := range e.userStatusChangedHandlers {
		err = f(userID, status)
		if err != nil {
			return
		}
	}
	return
}

// RegisterUserStatusChanged
func (e *event) RegisterUserStatusChanged(f func(string, bool) error) {
	e.userStatusChangedHandlers = append(e.userStatusChangedHandlers, f)
}

// 用户删除
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
