package driveradapters

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authorization/common"
	"Authorization/interfaces"
	"Authorization/interfaces/mock"
)

// 创建测试用的initData实例
func newTestInitData(
	resourceType interfaces.LogicsResourceType,
	role interfaces.LogicsRole,
	policy interfaces.LogicsPolicy,
	logger common.Logger,
) *initData {
	return &initData{
		log:          logger,
		resourceType: resourceType,
		role:         role,
		policy:       policy,
		memberStringTypes: map[string]interfaces.AccessorType{
			"user":       interfaces.AccessorUser,
			"department": interfaces.AccessorDepartment,
			"group":      interfaces.AccessorGroup,
			"app":        interfaces.AccessorApp,
		},
		accessorStrToType: map[string]interfaces.AccessorType{
			"user":       interfaces.AccessorUser,
			"department": interfaces.AccessorDepartment,
			"group":      interfaces.AccessorGroup,
			"role":       interfaces.AccessorRole,
			"app":        interfaces.AccessorApp,
		},
	}
}

func TestInitData_InitResourceType(t *testing.T) {
	Convey("测试InitResourceType方法", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResourceType := mock.NewMockLogicsResourceType(ctrl)
		mockRole := mock.NewMockLogicsRole(ctrl)
		mockPolicy := mock.NewMockLogicsPolicy(ctrl)
		logger := common.NewLogger()

		initData := newTestInitData(mockResourceType, mockRole, mockPolicy, logger)

		Convey("当JSON解析失败时", func() {
			// 这里我们无法直接测试JSON解析失败，因为数据是嵌入的
			// 但我们可以测试InitResourceTypes调用失败的情况
			mockResourceType.EXPECT().InitResourceTypes(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

			// 由于嵌入的JSON数据是有效的，这个方法会正常执行
			// 但会调用InitResourceTypes，我们可以验证这个调用
			initData.InitResourceType()
		})

		Convey("当InitResourceTypes成功时", func() {
			mockResourceType.EXPECT().InitResourceTypes(gomock.Any(), gomock.Any()).Return(nil)

			initData.InitResourceType()
		})
	})
}

func TestInitData_InitRole(t *testing.T) {
	Convey("测试InitRole方法", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResourceType := mock.NewMockLogicsResourceType(ctrl)
		mockRole := mock.NewMockLogicsRole(ctrl)
		mockPolicy := mock.NewMockLogicsPolicy(ctrl)
		logger := common.NewLogger()

		initData := newTestInitData(mockResourceType, mockRole, mockPolicy, logger)

		Convey("当JSON解析失败时", func() {
			// 这里我们无法直接测试JSON解析失败，因为数据是嵌入的
			// 但我们可以测试InitResourceTypes调用失败的情况
			mockRole.EXPECT().InitRoles(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

			// 由于嵌入的JSON数据是有效的，这个方法会正常执行
			// 但会调用InitResourceTypes，我们可以验证这个调用
			initData.InitRole()
		})

		Convey("当InitResourceTypes成功时", func() {
			mockRole.EXPECT().InitRoles(gomock.Any(), gomock.Any()).Return(nil)

			initData.InitRole()
		})
	})
}

func TestInitData_InitRoleMembers(t *testing.T) {
	Convey("测试InitRoleMembers方法", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResourceType := mock.NewMockLogicsResourceType(ctrl)
		mockRole := mock.NewMockLogicsRole(ctrl)
		mockPolicy := mock.NewMockLogicsPolicy(ctrl)
		logger := common.NewLogger()

		initData := newTestInitData(mockResourceType, mockRole, mockPolicy, logger)

		Convey("当JSON解析失败时", func() {
			// 这里我们无法直接测试JSON解析失败，因为数据是嵌入的
			// 但我们可以测试InitResourceTypes调用失败的情况
			mockRole.EXPECT().InitRoleMemebers(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

			// 由于嵌入的JSON数据是有效的，这个方法会正常执行
			// 但会调用InitResourceTypes，我们可以验证这个调用
			initData.InitRoleMembers()
		})

		Convey("当InitResourceTypes成功时", func() {
			mockRole.EXPECT().InitRoleMemebers(gomock.Any(), gomock.Any()).Return(nil)

			initData.InitRoleMembers()
		})
	})
}

func TestInitData_InitPolicy(t *testing.T) {
	Convey("测试InitPolicy方法", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResourceType := mock.NewMockLogicsResourceType(ctrl)
		mockRole := mock.NewMockLogicsRole(ctrl)
		mockPolicy := mock.NewMockLogicsPolicy(ctrl)
		logger := common.NewLogger()

		initData := newTestInitData(mockResourceType, mockRole, mockPolicy, logger)

		Convey("当JSON解析失败时", func() {
			// 这里我们无法直接测试JSON解析失败，因为数据是嵌入的
			// 但我们可以测试InitResourceTypes调用失败的情况
			mockPolicy.EXPECT().InitPolicy(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

			// 由于嵌入的JSON数据是有效的，这个方法会正常执行
			// 但会调用InitResourceTypes，我们可以验证这个调用
			initData.InitPolicy()
		})

		Convey("当InitResourceTypes成功时", func() {
			mockPolicy.EXPECT().InitPolicy(gomock.Any(), gomock.Any()).Return(nil)

			initData.InitPolicy()
		})
	})
}
