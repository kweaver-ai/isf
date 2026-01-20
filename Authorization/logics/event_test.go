package logics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"Authorization/interfaces"
)

func newEvent() *event {
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
	}
	return e
}

func TestUserDeleted(t *testing.T) {
	e := newEvent()
	_ = NewEvent()

	handlerCalled := false
	e.RegisterUserDeleted(func(userID string) error {
		handlerCalled = true
		assert.Equal(t, "testUserID", userID)
		return nil
	})

	err := e.UserDeleted("testUserID")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestDepartmentDeleted(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterDepartmentDeleted(func(depID string) error {
		handlerCalled = true
		assert.Equal(t, "testDepID", depID)
		return nil
	})

	err := e.DepartmentDeleted("testDepID")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestGroupDeleted(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterUserGroupDeleted(func(groupID string) error {
		handlerCalled = true
		assert.Equal(t, "testGroupID", groupID)
		return nil
	})

	err := e.UserGroupDeleted("testGroupID")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestAppDeleted(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterAppDeleted(func(appID string) error {
		handlerCalled = true
		assert.Equal(t, "testAppID", appID)
		return nil
	})

	err := e.AppDeleted("testAppID")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestAppNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterAppNameModified(func(info *interfaces.AppInfo) error {
		handlerCalled = true
		assert.Equal(t, "testAppID", info.ID)
		return nil
	})

	err := e.AppNameModified(&interfaces.AppInfo{ID: "testAppID"})
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestOrgNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterOrgNameModified(func(id, name string) error {
		handlerCalled = true
		assert.Equal(t, "testID", id)
		assert.Equal(t, "testName", name)
		return nil
	})

	err := e.OrgNameModified("testID", "testName", interfaces.AccessorUser)
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestRoleDeleted(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterRoleDeleted(func(roleID string) error {
		handlerCalled = true
		assert.Equal(t, "testRoleID", roleID)
		return nil
	})

	err := e.RoleDeleted("testRoleID")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestRoleNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterRoleNameModified(func(roleID, name string) error {
		handlerCalled = true
		assert.Equal(t, "testRoleID", roleID)
		assert.Equal(t, "testRoleName", name)
		return nil
	})

	err := e.RoleNameModified("testRoleID", "testRoleName")
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestRegisterUserNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterUserNameModified(func(id, name string) error {
		handlerCalled = true
		assert.Equal(t, "testUserID", id)
		assert.Equal(t, "testUserName", name)
		return nil
	})

	// 通过 OrgNameModified 调用
	err := e.OrgNameModified("testUserID", "testUserName", interfaces.AccessorUser)
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestRegisterDepartmentNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterDepartmentNameModified(func(id, name string) error {
		handlerCalled = true
		assert.Equal(t, "testDepID", id)
		assert.Equal(t, "testDepName", name)
		return nil
	})

	// 通过 OrgNameModified 调用
	err := e.OrgNameModified("testDepID", "testDepName", interfaces.AccessorDepartment)
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestRegisterUserGroupNameModified(t *testing.T) {
	e := newEvent()

	handlerCalled := false
	e.RegisterUserGroupNameModified(func(id, name string) error {
		handlerCalled = true
		assert.Equal(t, "testGroupID", id)
		assert.Equal(t, "testGroupName", name)
		return nil
	})

	// 通过 OrgNameModified 调用
	err := e.OrgNameModified("testGroupID", "testGroupName", interfaces.AccessorGroup)
	assert.Nil(t, err)
	assert.True(t, handlerCalled)
}

func TestOrgNameModifiedWithDifferentTypes(t *testing.T) {
	e := newEvent()

	userHandlerCalled := false
	depHandlerCalled := false
	groupHandlerCalled := false
	orgHandlerCalled := false

	e.RegisterUserNameModified(func(id, name string) error {
		userHandlerCalled = true
		return nil
	})

	e.RegisterDepartmentNameModified(func(id, name string) error {
		depHandlerCalled = true
		return nil
	})

	e.RegisterUserGroupNameModified(func(id, name string) error {
		groupHandlerCalled = true
		return nil
	})

	e.RegisterOrgNameModified(func(id, name string) error {
		orgHandlerCalled = true
		return nil
	})

	// 测试用户类型
	err := e.OrgNameModified("testID", "testName", interfaces.AccessorUser)
	assert.Nil(t, err)
	assert.True(t, userHandlerCalled)
	assert.False(t, depHandlerCalled)
	assert.False(t, groupHandlerCalled)
	assert.True(t, orgHandlerCalled)

	// 重置状态
	userHandlerCalled = false
	depHandlerCalled = false
	groupHandlerCalled = false
	orgHandlerCalled = false

	// 测试部门类型
	err = e.OrgNameModified("testID", "testName", interfaces.AccessorDepartment)
	assert.Nil(t, err)
	assert.False(t, userHandlerCalled)
	assert.True(t, depHandlerCalled)
	assert.False(t, groupHandlerCalled)
	assert.True(t, orgHandlerCalled)

	// 重置状态
	userHandlerCalled = false
	depHandlerCalled = false
	groupHandlerCalled = false
	orgHandlerCalled = false

	// 测试组类型
	err = e.OrgNameModified("testID", "testName", interfaces.AccessorGroup)
	assert.Nil(t, err)
	assert.False(t, userHandlerCalled)
	assert.False(t, depHandlerCalled)
	assert.True(t, groupHandlerCalled)
	assert.True(t, orgHandlerCalled)
}

func TestMultipleHandlers(t *testing.T) {
	e := newEvent()

	callCount := 0
	handler1 := func(userID string) error {
		callCount++
		assert.Equal(t, "testUserID", userID)
		return nil
	}
	handler2 := func(userID string) error {
		callCount++
		assert.Equal(t, "testUserID", userID)
		return nil
	}

	e.RegisterUserDeleted(handler1)
	e.RegisterUserDeleted(handler2)

	err := e.UserDeleted("testUserID")
	assert.Nil(t, err)
	assert.Equal(t, 2, callCount)
}

func TestHandlerError(t *testing.T) {
	e := newEvent()

	expectedError := assert.AnError
	e.RegisterUserDeleted(func(userID string) error {
		return expectedError
	})

	err := e.UserDeleted("testUserID")
	assert.Equal(t, expectedError, err)
}
