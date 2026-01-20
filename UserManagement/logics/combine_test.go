package logics

import (
	"context"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestNewCombine(t *testing.T) {
	Convey("NewCombine", t, func() {
		sqlDB, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		dbPool = sqlDB

		data := NewCombine()
		assert.NotEqual(t, data, nil)
	})
}

func TestSearchInOrgTree(t *testing.T) {
	Convey("组织架构树用户和部门搜索-逻辑层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.NewMockLogicsUser(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		depart := mock.NewMockLogicsDepartment(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		comb := combine{
			user:       user,
			role:       role,
			department: depart,
			trace:      trace,
		}

		var visitor interfaces.Visitor
		visitor.ID = "zzzew"
		mapRoles := make(map[interfaces.Role]bool)
		roleInfo := make(map[string]map[interfaces.Role]bool)
		roleInfo[visitor.ID] = mapRoles
		var info interfaces.OrgShowPageInfo

		testErr := rest.NewHTTPError("error", 503000000, nil)
		var ctx context.Context
		Convey("获取用户角色错误-报错 ", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Times(1).Return(roleInfo, testErr)

			_, _, _, _, err := comb.SearchInOrgTree(ctx, &visitor, info)
			assert.Equal(t, err, testErr)
		})

		mapRoles[interfaces.SystemRoleAuditAdmin] = true
		Convey("用户没有指定角色-报错", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			_, _, _, _, err := comb.SearchInOrgTree(ctx, &visitor, info)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has this role", rest.BadRequest, nil))
		})

		mapRoles[interfaces.SystemRoleAuditAdmin] = true
		info.Role = interfaces.SystemRoleAuditAdmin
		info.BShowUsers = true
		Convey("用户信息获取报错", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			user.EXPECT().SearchUsersByKeyInScope(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, testErr)

			_, _, _, _, err := comb.SearchInOrgTree(ctx, &visitor, info)
			assert.Equal(t, err, testErr)
		})

		info.BShowDeparts = true
		Convey("部门信息获取报错", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			user.EXPECT().SearchUsersByKeyInScope(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, nil)
			depart.EXPECT().SearchDepartsByKey(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, testErr)

			_, _, _, _, err := comb.SearchInOrgTree(ctx, &visitor, info)
			assert.Equal(t, err, testErr)
		})

		Convey("成功", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			user.EXPECT().SearchUsersByKeyInScope(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 2, nil)
			depart.EXPECT().SearchDepartsByKey(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 3, nil)

			users, departments, num1, num2, err := comb.SearchInOrgTree(ctx, &visitor, info)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(users), 0)
			assert.Equal(t, len(departments), 0)
			assert.Equal(t, num1, 2)
			assert.Equal(t, num2, 3)
		})
	})
}

func TestConvertIDToName(t *testing.T) {
	Convey("根据ID获取NAME-逻辑层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.NewMockLogicsUser(ctrl)
		depart := mock.NewMockLogicsDepartment(ctrl)
		contactor := mock.NewMockLogicsContactor(ctrl)
		group := mock.NewMockLogicsGroup(ctrl)
		app := mock.NewMockLogicsApp(ctrl)

		c := combine{
			user:       user,
			department: depart,
			contactor:  contactor,
			group:      group,
			app:        app,
		}

		nameInfo := []interfaces.NameInfo{}
		info := interfaces.OrgIDInfo{}
		visitor := interfaces.Visitor{}
		Convey("获取用户名称错误-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, testErr)

			_, err := c.ConvertIDToName(&visitor, &info, false, true)
			assert.Equal(t, err, testErr)
		})

		Convey("获取部门名称错误-报错", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			depart.EXPECT().ConvertDepartmentName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, testErr)
			_, err := c.ConvertIDToName(&visitor, &info, false, true)
			assert.Equal(t, err, testErr)
		})

		Convey("获取联系人组名称错误-报错", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), false).AnyTimes().Return(nameInfo, nil)
			depart.EXPECT().ConvertDepartmentName(gomock.Any(), gomock.Any(), false).AnyTimes().Return(nameInfo, nil)
			contactor.EXPECT().ConvertContactorName(gomock.Any(), false).AnyTimes().Return(nameInfo, testErr)
			_, err := c.ConvertIDToName(&visitor, &info, false, false)
			assert.Equal(t, err, testErr)
		})

		Convey("获取用户组名称错误-报错", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), true).AnyTimes().Return(nameInfo, nil)
			depart.EXPECT().ConvertDepartmentName(gomock.Any(), gomock.Any(), true).AnyTimes().Return(nameInfo, nil)
			contactor.EXPECT().ConvertContactorName(gomock.Any(), true).AnyTimes().Return(nameInfo, nil)
			group.EXPECT().ConvertGroupName(gomock.Any(), gomock.Any(), true).AnyTimes().Return(nameInfo, testErr)
			_, err := c.ConvertIDToName(&visitor, &info, false, true)
			assert.Equal(t, err, testErr)
		})

		Convey("获取应用账户名称错误-报错", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			depart.EXPECT().ConvertDepartmentName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			contactor.EXPECT().ConvertContactorName(gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			group.EXPECT().ConvertGroupName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			app.EXPECT().ConvertAppName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := c.ConvertIDToName(&visitor, &info, false, true)
			assert.Equal(t, err, testErr)
		})

		Convey("成功-报错", func() {
			user.EXPECT().ConvertUserName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			depart.EXPECT().ConvertDepartmentName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			contactor.EXPECT().ConvertContactorName(gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			group.EXPECT().ConvertGroupName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			app.EXPECT().ConvertAppName(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)
			_, err := c.ConvertIDToName(&visitor, &info, false, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetEmails(t *testing.T) {
	Convey("获取部门和用户Email-逻辑层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.NewMockLogicsUser(ctrl)
		depart := mock.NewMockLogicsDepartment(ctrl)
		contactor := mock.NewMockLogicsContactor(ctrl)
		group := mock.NewMockLogicsGroup(ctrl)

		c := combine{
			user:       user,
			department: depart,
			contactor:  contactor,
			group:      group,
		}

		eInfo := []interfaces.EmailInfo{}
		info := interfaces.OrgIDInfo{}
		visitor := interfaces.Visitor{}

		Convey("获取用户Email信息-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().GetUserEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(eInfo, testErr)

			_, err := c.GetEmails(&visitor, &info)
			assert.Equal(t, err, testErr)
		})

		Convey("获取部门Email信息-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			user.EXPECT().GetUserEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(eInfo, nil)
			depart.EXPECT().GetDepartEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(eInfo, testErr)

			_, err := c.GetEmails(&visitor, &info)
			assert.Equal(t, err, testErr)
		})

		Convey("成功", func() {
			user.EXPECT().GetUserEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(eInfo, nil)
			depart.EXPECT().GetDepartEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(eInfo, nil)

			_, err := c.GetEmails(&visitor, &info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetUserAndDepartmentInScope(t *testing.T) {
	Convey("获取范围内的用户和部门-逻辑层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.NewMockLogicsUser(ctrl)
		depart := mock.NewMockLogicsDepartment(ctrl)
		contactor := mock.NewMockLogicsContactor(ctrl)
		group := mock.NewMockLogicsGroup(ctrl)

		c := combine{
			user:       user,
			department: depart,
			contactor:  contactor,
			group:      group,
		}

		depIds := []string{}
		Convey("获取范围内所有子部门信息-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			depart.EXPECT().GetAllChildDeparmentIDs(gomock.Any()).AnyTimes().Return(depIds, testErr)

			_, _, err := c.GetUserAndDepartmentInScope(nil, nil, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("获取存在的用户-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			depart.EXPECT().GetAllChildDeparmentIDs(gomock.Any()).AnyTimes().Return(depIds, nil)
			user.EXPECT().GetUsersInDepartments(gomock.Any(), gomock.Any()).AnyTimes().Return(depIds, testErr)

			_, _, err := c.GetUserAndDepartmentInScope(nil, nil, nil)
			assert.Equal(t, err, testErr)
		})

		InDepartIds := []string{"D1"}
		scopeDepIds := []string{"D1", "D2", "D3", "D1"}
		userIDs := []string{"U1"}
		Convey("成功", func() {
			depart.EXPECT().GetAllChildDeparmentIDs(gomock.Any()).AnyTimes().Return(scopeDepIds, nil)
			user.EXPECT().GetUsersInDepartments(gomock.Any(), gomock.Any()).AnyTimes().Return(userIDs, nil)

			out1, out2, err := c.GetUserAndDepartmentInScope(nil, InDepartIds, nil)
			assert.Equal(t, err, nil)
			assert.Equal(t, out1, userIDs)
			assert.Equal(t, out2, InDepartIds)
		})
	})
}

func TestSearchGroupAndMemberInfoByKey(t *testing.T) {
	Convey("客户端搜索组和组成员-逻辑层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mock.NewMockLogicsUser(ctrl)
		depart := mock.NewMockLogicsDepartment(ctrl)
		contactor := mock.NewMockLogicsContactor(ctrl)
		group := mock.NewMockLogicsGroup(ctrl)

		c := combine{
			user:       user,
			department: depart,
			contactor:  contactor,
			group:      group,
		}

		info := interfaces.SearchClientInfo{
			BShowMember: true,
		}
		Convey("符合搜索条件的组数量获取失败-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			group.EXPECT().SearchMemberNumByKeyword(gomock.Any()).AnyTimes().Return(1, testErr)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, testErr)
		})

		memebers := []interfaces.MemberInfo{}
		Convey("符合搜索条件的组信息获取失败-报错 ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			group.EXPECT().SearchMemberNumByKeyword(gomock.Any()).AnyTimes().Return(1, nil)
			group.EXPECT().SearchMembersByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(memebers, testErr)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, testErr)
		})

		Convey("成功1", func() {
			group.EXPECT().SearchMemberNumByKeyword(gomock.Any()).AnyTimes().Return(1, nil)
			group.EXPECT().SearchMembersByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(memebers, nil)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, nil)
		})

		info.BShowGroup = true
		info.BShowMember = false
		Convey("获取符合条件的组数量失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			group.EXPECT().SearchGroupNumByKeyword(gomock.Any()).AnyTimes().Return(1, testErr)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, testErr)
		})

		nameInfo := []interfaces.NameInfo{}
		Convey("获取符合条件的组信息失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			group.EXPECT().SearchGroupNumByKeyword(gomock.Any()).AnyTimes().Return(1, nil)
			group.EXPECT().SearchGroupByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, testErr)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, testErr)
		})

		Convey("成功2", func() {
			group.EXPECT().SearchGroupNumByKeyword(gomock.Any()).AnyTimes().Return(1, nil)
			group.EXPECT().SearchGroupByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nameInfo, nil)

			_, err := c.SearchGroupAndMemberInfoByKey(&info)
			assert.Equal(t, err, nil)
		})
	})
}
