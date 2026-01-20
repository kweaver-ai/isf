package logics

import (
	"context"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/errors"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestRemoveDuplicatStrs(t *testing.T) {
	Convey("RemoveDuplicatStrs delelte the same string", t, func() {
		str1 := []string{
			"aaaa",
			"bbbb",
			"ccccc",
			"cccdd",
			"ddddd",
		}

		str2 := []string{
			"ccccc",
			"bbbb",
			"ddddd",
			"aaaa",
			"bbbb",
			"ccccc",
			"aaaa",
			"cccdd",
			"ddddd",
			"bbbb",
		}

		RemoveDuplicatStrs(&str2)
		assert.Equal(t, str2, str1)
	})
}

func TestSplitArray(t *testing.T) {
	Convey("SplitArray ", t, func() {
		strTemp := "xxxxxx"
		strMax := make([]string, 2650)
		for i := 0; i < 2650; {
			strMax[i] = strTemp
			i++
		}

		strArry := SplitArray(strMax)
		count := len(strArry)
		assert.Equal(t, count, 6)
		for k, v := range strArry {
			if k < count-1 {
				assert.Equal(t, len(v), 500)
			}
			if k == count-1 {
				assert.Equal(t, len(v), 2650%500)
			}

			for _, v1 := range v {
				assert.Equal(t, v1, strTemp)
			}
		}
	})
}

func TestDifference(t *testing.T) {
	Convey("Difference ", t, func() {
		str1 := []string{
			"aaaa",
			"bbbb",
			"cccc",
		}

		str2 := []string{
			"ddd",
			"eee",
			"cccc",
		}

		str := Difference(str1, str2)

		assert.Equal(t, str[0], "aaaa")
		assert.Equal(t, str[1], "bbbb")
	})
}

func TestIntersection(t *testing.T) {
	Convey("Intersection ", t, func() {
		str1 := []string{
			"aaaa",
			"bbbb",
			"cccc",
		}

		str2 := []string{
			"ddd",
			"eee",
			"cccc",
			"dddd",
		}

		str := Intersection(str1, str2)

		assert.Equal(t, len(str), 1)
		assert.Equal(t, str[0], "cccc")
	})
}

func TestGetRolesByUserID(t *testing.T) {
	Convey("根据用户ID获取用户角色 ", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		role := mock.NewMockLogicsRole(ctrl)
		Convey("获取用户角色失败", func() {
			tempErr1 := rest.NewHTTPError("xxx", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, tempErr1)

			_, err := getRolesByUserID(role, "")
			assert.Equal(t, err, tempErr1)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo["xxx"] = userRole
		Convey("获取用户角色成功", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			out, err := getRolesByUserID(role, "xxx")
			assert.Equal(t, err, nil)
			assert.Equal(t, out[interfaces.SystemRoleNormalUser], true)
			assert.Equal(t, len(out), 1)
		})
	})
}

func TestGetRolesByUserID2(t *testing.T) {
	Convey("根据用户ID获取用户角色2 ", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var ctx context.Context

		role := mock.NewMockLogicsRole(ctrl)
		Convey("获取用户角色失败", func() {
			tempErr1 := rest.NewHTTPError("xxx", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, tempErr1)

			_, err := getRolesByUserID2(ctx, role, "")
			assert.Equal(t, err, tempErr1)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo["xxx"] = userRole
		Convey("获取用户角色成功", func() {
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			out, err := getRolesByUserID2(ctx, role, "xxx")
			assert.Equal(t, err, nil)
			assert.Equal(t, out[interfaces.SystemRoleNormalUser], true)
			assert.Equal(t, len(out), 1)
		})
	})
}

func TestCheckManageAuthority(t *testing.T) {
	Convey("检查管理权限", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		role := mock.NewMockLogicsRole(ctrl)
		userID := "xxxxasdas"
		Convey("获取用户角色失败", func() {
			tempErr1 := rest.NewHTTPError("xxx", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, tempErr1)

			err := checkManageAuthority(role, userID)
			assert.Equal(t, err, tempErr1)
		})

		Convey("超级管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkManageAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("系统管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkManageAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("安全管理员、审计管理员、普通用户无权限", func() {
			tempErr := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleAuditAdmin] = true
			userRoles[interfaces.SystemRoleSecAdmin] = true
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkManageAuthority(role, userID)
			assert.Equal(t, err, tempErr)
		})
	})
}

func TestCheckGetInfoAuthority(t *testing.T) {
	Convey("检查列举权限", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		role := mock.NewMockLogicsRole(ctrl)
		userID := "xxxx111"
		Convey("获取用户角色失败", func() {
			tempErr1 := rest.NewHTTPError("xxx", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, tempErr1)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, tempErr1)
		})

		Convey("超级管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("系统管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("安全管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSecAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("审计管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleAuditAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("普通用户无权限权限", func() {
			tempErr := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority(role, userID)
			assert.Equal(t, err, tempErr)
		})
	})
}

func TestCheckGetInfoAuthority2(t *testing.T) {
	Convey("检查列举权限2", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var ctx context.Context

		role := mock.NewMockLogicsRole(ctrl)
		userID := "xxxx111"
		Convey("获取用户角色失败", func() {
			tempErr1 := rest.NewHTTPError("xxx", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, tempErr1)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, tempErr1)
		})

		Convey("超级管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("系统管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("安全管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSecAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("审计管理员具有权限", func() {
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleAuditAdmin] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, nil)
		})

		Convey("普通用户无权限权限", func() {
			tempErr := rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority")
			roleInfo := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfo[userID] = userRoles
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkGetInfoAuthority2(ctx, role, userID)
			assert.Equal(t, err, tempErr)
		})
	})
}

func TestCheckAppPerm(t *testing.T) {
	Convey("检查应用账户权限", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)

		allPerm := make(map[interfaces.OrgType]interfaces.AppOrgPerm)
		Convey("获取权限信息报错", func() {
			testErr := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			orgPermApp.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(allPerm, testErr)

			err := checkAppPerm(orgPermApp, "", interfaces.User, interfaces.Modify)
			assert.Equal(t, err, testErr)
		})

		var info interfaces.AppOrgPerm
		info.Value = 1
		allPerm[interfaces.Department] = info
		Convey("获取权限信息成功，但是无权限", func() {
			orgPermApp.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(allPerm, nil)

			testErr1 := rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil)
			err := checkAppPerm(orgPermApp, "", interfaces.User, interfaces.Read)
			assert.Equal(t, err, testErr1)
		})

		Convey("获取权限信息成功，有权限", func() {
			orgPermApp.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(allPerm, nil)

			err := checkAppPerm(orgPermApp, "", interfaces.Department, interfaces.Modify)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckUserRole(t *testing.T) {
	Convey("检查用户是否具有权限", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		role := mock.NewMockLogicsRole(ctrl)
		userID := "xxxxdddd"
		allowRoles := []interfaces.Role{interfaces.SystemRoleSecAdmin}
		Convey("获取用户角色信息失败", func() {
			testErr := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)

			err := checkUserRole(role, userID, allowRoles)
			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRoles := make(map[interfaces.Role]bool)
		userRoles[interfaces.SystemRoleAuditAdmin] = true
		roleInfo[userID] = userRoles
		Convey("获取用户角色信息，无允许的角色", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			testErr1 := rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil)
			err := checkUserRole(role, userID, allowRoles)
			assert.Equal(t, err, testErr1)
		})

		userRoles[interfaces.SystemRoleSecAdmin] = true
		Convey("获取用户角色信息，有允许的角色", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkUserRole(role, userID, allowRoles)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckUserRole2(t *testing.T) {
	Convey("检查用户是否具有权限2", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var ctx context.Context

		role := mock.NewMockLogicsRole(ctrl)
		userID := "xxxxdddd1"
		allowRoles := []interfaces.Role{interfaces.SystemRoleSecAdmin}
		Convey("获取用户角色信息失败", func() {
			testErr := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			err := checkUserRole2(ctx, role, userID, allowRoles)
			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRoles := make(map[interfaces.Role]bool)
		userRoles[interfaces.SystemRoleAuditAdmin] = true
		roleInfo[userID] = userRoles
		Convey("获取用户角色信息，无允许的角色", func() {
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			testErr1 := rest.NewHTTPErrorV2(rest.Forbidden, "this user do not has the authority")
			err := checkUserRole2(ctx, role, userID, allowRoles)
			assert.Equal(t, err, testErr1)
		})

		userRoles[interfaces.SystemRoleSecAdmin] = true
		Convey("获取用户角色信息，有允许的角色", func() {
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)

			err := checkUserRole2(ctx, role, userID, allowRoles)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckAppPerm2(t *testing.T) {
	Convey("检查应用账户是否具有权限2", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var ctx context.Context

		orgPerm := mock.NewMockDBOrgPermApp(ctrl)
		userID := "xxxxdddd"
		Convey("获取应用账户权限失败", func() {
			testErr := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			orgPerm.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			err := checkAppPerm2(ctx, orgPerm, userID, interfaces.Group, interfaces.Modify)
			assert.Equal(t, err, testErr)
		})

		Convey("权限检查成失败", func() {
			orgPerm.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)

			err := checkAppPerm2(ctx, orgPerm, userID, interfaces.Department, interfaces.Read)
			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.Forbidden, "this user do not has the authority"))
		})

		out := make(map[interfaces.OrgType]interfaces.AppOrgPerm)
		out[interfaces.Group] = interfaces.AppOrgPerm{
			Subject: userID,
			Value:   interfaces.Read,
		}
		Convey("权限检查成功", func() {
			orgPerm.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(out, nil)

			err := checkAppPerm2(ctx, orgPerm, userID, interfaces.Group, interfaces.Read)
			assert.Equal(t, err, nil)
		})
	})
}
