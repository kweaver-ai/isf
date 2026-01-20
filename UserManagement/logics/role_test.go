// Package logics
package logics

import (
	"context"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestGetRolesByUserIDs(t *testing.T) {
	Convey("批量获取用户角色", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		roleDB := mock.NewMockDBRole(ctrl)

		u := &role{
			db: roleDB,
		}

		Convey("获取用户角色数据失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			roleDB.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetRolesByUserIDs([]string{"xxx"})
			assert.Equal(t, err, testErr)
		})

		Convey("获取用户角色数据成功", func() {
			roleInfo := make(map[interfaces.Role]bool)
			roleInfo["role1"] = true
			userRoleInfo := make(map[string]map[interfaces.Role]bool)
			userRoleInfo["xxx"] = roleInfo
			roleDB.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(userRoleInfo, nil)

			out, err := u.GetRolesByUserIDs([]string{"xxx", "yyy", interfaces.SystemSysAdmin})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 3)
			temp1, ok := out["xxx"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp1["role1"], true)
			assert.Equal(t, temp1[interfaces.SystemRoleNormalUser], true)
			temp2, ok := out["yyy"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp2[interfaces.SystemRoleNormalUser], true)
			temp3, ok := out[interfaces.SystemSysAdmin]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp3), 0)
		})

		Convey("获取用户角色数据成功1", func() {
			roleDB.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, nil)

			out, err := u.GetRolesByUserIDs([]string{"xxx", "yyy", interfaces.SystemSysAdmin})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 3)
			temp1, ok := out["xxx"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp1[interfaces.SystemRoleNormalUser], true)
			assert.Equal(t, len(temp1), 1)
			temp2, ok := out["yyy"]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp2), 1)
			assert.Equal(t, temp2[interfaces.SystemRoleNormalUser], true)
			temp3, ok := out[interfaces.SystemSysAdmin]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp3), 0)
		})
	})
}

func TestGetRolesByUserIDs2(t *testing.T) {
	Convey("批量获取用户角色2", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		roleDB := mock.NewMockDBRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		var ctx context.Context

		u := &role{
			db:    roleDB,
			trace: trace,
		}

		Convey("获取用户角色数据失败", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			roleDB.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetRolesByUserIDs2(ctx, []string{"xxx"})
			assert.Equal(t, err, testErr)
		})

		Convey("获取用户角色数据成功", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			roleInfo := make(map[interfaces.Role]bool)
			roleInfo["role1"] = true
			userRoleInfo := make(map[string]map[interfaces.Role]bool)
			userRoleInfo["xxx"] = roleInfo
			roleDB.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(userRoleInfo, nil)

			out, err := u.GetRolesByUserIDs2(ctx, []string{"xxx", "yyy", interfaces.SystemSysAdmin})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 3)
			temp1, ok := out["xxx"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp1["role1"], true)
			assert.Equal(t, temp1[interfaces.SystemRoleNormalUser], true)
			temp2, ok := out["yyy"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp2[interfaces.SystemRoleNormalUser], true)
			temp3, ok := out[interfaces.SystemSysAdmin]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp3), 0)
		})

		Convey("获取用户角色数据成功1", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			roleDB.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)

			out, err := u.GetRolesByUserIDs2(ctx, []string{"xxx", "yyy", interfaces.SystemSysAdmin})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 3)
			temp1, ok := out["xxx"]
			assert.Equal(t, ok, true)
			assert.Equal(t, temp1[interfaces.SystemRoleNormalUser], true)
			assert.Equal(t, len(temp1), 1)
			temp2, ok := out["yyy"]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp2), 1)
			assert.Equal(t, temp2[interfaces.SystemRoleNormalUser], true)
			temp3, ok := out[interfaces.SystemSysAdmin]
			assert.Equal(t, ok, true)
			assert.Equal(t, len(temp3), 0)
		})
	})
}

func TestGetOrgManagersInfo(t *testing.T) {
	Convey("获取组织管理员信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		roleDB := mock.NewMockDBRole(ctrl)
		departDB := mock.NewMockDBDepartment(ctrl)
		userDB := mock.NewMockDBUser(ctrl)

		u := &role{
			db:       roleDB,
			departDB: departDB,
			userDB:   userDB,
		}

		Convey("获取组织管理员管理部门失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			userDB.EXPECT().GetOrgManagersDepartInfo(gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetOrgManagersInfo([]string{"xxx"}, interfaces.OrgManagerInfoRange{ShowSubUserIDs: true})
			assert.Equal(t, err, testErr)
		})

		Convey("获取部门信息失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			orgManagerDepInfos := make(map[string][]string)
			orgManagerDepInfos["xxx"] = []string{"dep1", "dep2"}
			userDB.EXPECT().GetOrgManagersDepartInfo(gomock.Any()).AnyTimes().Return(orgManagerDepInfos, nil)
			departDB.EXPECT().GetDepartmentInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetOrgManagersInfo([]string{"xxx"}, interfaces.OrgManagerInfoRange{ShowSubUserIDs: true})
			assert.Equal(t, err, testErr)
		})

		Convey("获取部门子成员失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			orgManagerDepInfos := make(map[string][]string)
			orgManagerDepInfos["xxx"] = []string{"dep1", "dep2"}
			departInfo := []interfaces.DepartmentDBInfo{
				{
					ID:   "dep1",
					Path: "path1",
				},
				{
					ID:   "dep2",
					Path: "path2",
				},
			}
			userDB.EXPECT().GetOrgManagersDepartInfo(gomock.Any()).AnyTimes().Return(orgManagerDepInfos, nil)
			departDB.EXPECT().GetDepartmentInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(departInfo, nil)
			departDB.EXPECT().GetAllSubUserIDsByDepartPath(gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetOrgManagersInfo([]string{"xxx"}, interfaces.OrgManagerInfoRange{ShowSubUserIDs: true})
			assert.Equal(t, err, testErr)
		})

		Convey("获取组织管理员信息成功", func() {
			orgManagerDepInfos := make(map[string][]string)
			orgManagerDepInfos["xxx"] = []string{"dep1", "dep2"}
			departInfo := []interfaces.DepartmentDBInfo{
				{
					ID:   "dep1",
					Path: "path1",
				},
				{
					ID:   "dep2",
					Path: "path2",
				},
			}
			depSubUsersIDs := make(map[string][]string)
			depSubUsersIDs["dep1"] = []string{"user1", "user2"}
			depSubUsersIDs["dep2"] = []string{"user3", "user4"}
			userDB.EXPECT().GetOrgManagersDepartInfo(gomock.Any()).AnyTimes().Return(orgManagerDepInfos, nil)
			departDB.EXPECT().GetDepartmentInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(departInfo, nil)
			departDB.EXPECT().GetAllSubUserIDsByDepartPath(gomock.Any()).AnyTimes().Return([]string{"user1", "user2"}, nil)
			departDB.EXPECT().GetAllSubUserIDsByDepartPath(gomock.Any()).AnyTimes().Return([]string{"user3", "user4"}, nil)

			_, err := u.GetOrgManagersInfo([]string{"xxx"}, interfaces.OrgManagerInfoRange{ShowSubUserIDs: true})
			assert.Equal(t, err, nil)
		})
	})
}

// 根据角色ID获取角色成员
func TestGetUserIDsByRoleIDs(t *testing.T) {
	Convey("根据角色ID获取角色成员", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		roleDB := mock.NewMockDBRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		var ctx context.Context

		u := &role{
			db:    roleDB,
			trace: trace,
		}

		Convey("获取用户角色数据失败", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			roleDB.EXPECT().GetUserIDsByRoleIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			_, err := u.GetUserIDsByRoleIDs(ctx, []interfaces.Role{interfaces.SystemRoleSuperAdmin})
			assert.Equal(t, err, testErr)
		})

		Convey("获取用户角色数据成功， 去重检查", func() {
			roles := []interfaces.Role{interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSuperAdmin}
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			roleDB.EXPECT().GetUserIDsByRoleIDs(gomock.Any(), []interfaces.Role{interfaces.SystemRoleSuperAdmin}).AnyTimes().Return(nil, nil)

			_, err := u.GetUserIDsByRoleIDs(ctx, roles)
			assert.Equal(t, err, nil)
		})
	})
}
