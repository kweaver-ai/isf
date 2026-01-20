package logics

import (
	"context"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func newGroup(member interfaces.DBGroupMember, groupdb interfaces.DBGroup, user interfaces.DBUser, department interfaces.DBDepartment,
	eacplog interfaces.DrivenEacpLog, msg interfaces.DrivenMessageBroker, role interfaces.LogicsRole) *group {
	return &group{
		userDB:        user,
		departmentDB:  department,
		groupMemberDB: member,
		groupDB:       groupdb,
		eacpLog:       eacplog,
		messageBroker: msg,
		role:          role,
		i18n: common.NewI18n(common.I18nMap{
			i18nIDObjectsInUnDistributeUserGroup: {
				interfaces.SimplifiedChinese:  "未分配组",
				interfaces.TraditionalChinese: "未分配組",
				interfaces.AmericanEnglish:    "Unassigned Group",
			},
			i18nIDObjectsInUserNotFound: {
				interfaces.SimplifiedChinese:  "用户不存在",
				interfaces.TraditionalChinese: "用戶不存在",
				interfaces.AmericanEnglish:    "This user does not exist",
			},
			i18nIDObjectsInDepartNotFound: {
				interfaces.SimplifiedChinese:  "部门不存在",
				interfaces.TraditionalChinese: "部門不存在",
				interfaces.AmericanEnglish:    "This department does not exist",
			},
			i18nIDObjectsInGroupNotFound: {
				interfaces.SimplifiedChinese:  "用户组不存在",
				interfaces.TraditionalChinese: "用戶組不存在",
				interfaces.AmericanEnglish:    "This group does not exist",
			},
		}),
		logger: common.NewLogger(),
	}
}

//nolint:funlen
func TestUserMatch(t *testing.T) {
	Convey("user match", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		group := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		group.trace = trace
		group.orgPermApp = orgPermApp
		ctx := context.Background()

		vistor := interfaces.Visitor{
			ID:       strID,
			Type:     interfaces.App,
			Language: interfaces.SimplifiedChinese,
		}
		Convey("checkGroupGetAuth2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, testErr)
		})

		appOrgPern := interfaces.AppOrgPerm{
			Value: interfaces.Read,
		}
		perms := map[interfaces.OrgType]interfaces.AppOrgPerm{
			interfaces.Group: appOrgPern,
		}
		var groupInfo interfaces.GroupInfo
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("GetGroupByID2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID2 id is empty", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.URINotExist, "group not exist"))
		})

		groupInfo.ID = strID
		var userInfo interfaces.UserDBInfo
		Convey("GetUserInfoByName error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, testErr)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, testErr)
		})

		Convey("GetUserInfoByName empty", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			result, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, nil)
			assert.Equal(t, result, false)
		})

		userInfo.ID = strID
		Convey("GetGroupMembersByGroupIDs2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, testErr)
		})

		Convey("GetUsersPath2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, testErr)
		})

		D1NameInfo := interfaces.NameInfo{ID: "d1d1", Name: "d1"}
		D2NameInfo := interfaces.NameInfo{ID: "d2d1", Name: "d2"}
		D3NameInfo := interfaces.NameInfo{ID: "d3d1", Name: "d3"}
		D4NameInfo := interfaces.NameInfo{ID: "d4d1", Name: "d4"}
		D5NameInfo := interfaces.NameInfo{ID: "d5d1", Name: "d5"}
		D6NameInfo := interfaces.NameInfo{ID: "d6d1", Name: "d6"}

		member1 := interfaces.GroupMemberInfo{
			ID:         strID2,
			MemberType: 1,
		}
		member2 := interfaces.GroupMemberInfo{
			ID:         D6NameInfo.ID,
			MemberType: 2,
		}

		member3 := interfaces.GroupMemberInfo{
			ID:         userInfo.ID,
			MemberType: 1,
		}

		memberInfos := make(map[string][]interfaces.GroupMemberInfo)
		memberInfos[strID] = []interfaces.GroupMemberInfo{member1, member2, member3}

		tempPath := []string{
			D4NameInfo.ID + "/" + D5NameInfo.ID,
			D6NameInfo.ID + "/" + D3NameInfo.ID + "/" + D1NameInfo.ID,
		}
		userPaths := map[string][]string{
			userInfo.ID: tempPath,
		}

		D6Info := interfaces.DepartmentDBInfo{
			ID:   D6NameInfo.ID,
			Name: D6NameInfo.Name,
			Path: D2NameInfo.ID + "/" + D1NameInfo.ID + "/" + D6NameInfo.ID,
		}

		departInfos := []interfaces.DepartmentDBInfo{D6Info}

		temp1 := interfaces.DepartmentDBInfo{
			ID:   D1NameInfo.ID,
			Name: D1NameInfo.Name,
		}
		temp2 := interfaces.DepartmentDBInfo{
			ID:   D2NameInfo.ID,
			Name: D2NameInfo.Name,
		}
		temp3 := interfaces.DepartmentDBInfo{
			ID:   D3NameInfo.ID,
			Name: D3NameInfo.Name,
		}
		temp4 := interfaces.DepartmentDBInfo{
			ID:   D4NameInfo.ID,
			Name: D4NameInfo.Name,
		}
		temp5 := interfaces.DepartmentDBInfo{
			ID:   D5NameInfo.ID,
			Name: D5NameInfo.Name,
		}
		temp6 := interfaces.DepartmentDBInfo{
			ID:   D6NameInfo.ID,
			Name: D6NameInfo.Name,
		}
		allDepartInfos := []interfaces.DepartmentDBInfo{temp1, temp2, temp3, temp4, temp5, temp6}

		// 有两个members存在  member3和member2 因为paths包含 handlePathParams的userids为userid为userInfo.ID， 部门id是departmentid为D6
		// alldepartInfo 包括
		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(memberInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(departInfos, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(allDepartInfos, nil)
			result, uInfo, mInfos, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, nil)
			assert.Equal(t, result, true)

			assert.Equal(t, uInfo.ID, userInfo.ID)
			assert.Equal(t, len(uInfo.ParentDeps), 2)
			assert.Equal(t, len(uInfo.ParentDeps[0]), 2)
			assert.Equal(t, len(uInfo.ParentDeps[1]), 3)

			assert.Equal(t, uInfo.ParentDeps[0][0], D4NameInfo)
			assert.Equal(t, uInfo.ParentDeps[0][1], D5NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][0], D6NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][1], D3NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][2], D1NameInfo)

			assert.Equal(t, len(mInfos), 2)
			assert.Equal(t, mInfos[0].ID, D6NameInfo.ID)
			assert.Equal(t, mInfos[0].Name, D6NameInfo.Name)
			assert.Equal(t, len(mInfos[0].ParentDeps), 1)
			assert.Equal(t, len(mInfos[0].ParentDeps[0]), 2)
			assert.Equal(t, mInfos[0].ParentDeps[0][0], D2NameInfo)
			assert.Equal(t, mInfos[0].ParentDeps[0][1], D1NameInfo)

			assert.Equal(t, mInfos[1].ID, userInfo.ID)
			assert.Equal(t, mInfos[1].Name, strID2)
			assert.Equal(t, len(mInfos[1].ParentDeps), 2)
			assert.Equal(t, len(mInfos[1].ParentDeps[0]), 2)
			assert.Equal(t, len(mInfos[1].ParentDeps[1]), 3)
			assert.Equal(t, mInfos[1].ParentDeps[0][0], D4NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[0][1], D5NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][0], D6NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][1], D3NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][2], D1NameInfo)
		})

		member11 := interfaces.GroupMemberInfo{
			ID:         strID,
			MemberType: 1,
		}
		memberInfos1 := make(map[string][]interfaces.GroupMemberInfo)
		memberInfos1[strID] = []interfaces.GroupMemberInfo{member11}

		tempPath1 := []string{"-1"}
		userPaths1 := map[string][]string{
			userInfo.ID: tempPath1,
		}
		temp11 := interfaces.DepartmentDBInfo{
			ID:   "-1",
			Name: "",
		}
		allDepartInfos1 := []interfaces.DepartmentDBInfo{temp11}
		Convey("success undistrbute group", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(memberInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths1, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths1, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(allDepartInfos1, nil)
			result, uInfo, mInfos, err := group.UserMatch(ctx, &vistor, strID, strID2)

			assert.Equal(t, err, nil)
			assert.Equal(t, result, true)

			assert.Equal(t, uInfo.ID, userInfo.ID)
			assert.Equal(t, uInfo.Name, userInfo.Name)
			assert.Equal(t, len(uInfo.ParentDeps), 1)
			assert.Equal(t, len(uInfo.ParentDeps[0]), 1)

			tempName1 := group.i18n.Load(i18nIDObjectsInUnDistributeUserGroup, vistor.Language)
			assert.Equal(t, uInfo.ParentDeps[0][0].ID, "-1")
			assert.Equal(t, uInfo.ParentDeps[0][0].Name, tempName1)

			assert.Equal(t, len(mInfos), 1)
			assert.Equal(t, mInfos[0].ID, member11.ID)
			assert.Equal(t, mInfos[0].Name, strID2)
			assert.Equal(t, len(mInfos[0].ParentDeps), 1)
			assert.Equal(t, mInfos[0].ParentDeps[0][0].ID, "-1")
			assert.Equal(t, mInfos[0].ParentDeps[0][0].Name, tempName1)
		})
	})
}

//nolint:funlen
func TestSearchInAllGroupOrg(t *testing.T) {
	Convey("SearchInAllGroupOrg", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		group := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		group.trace = trace
		group.orgPermApp = orgPermApp
		ctx := context.Background()

		vistor := interfaces.Visitor{
			ID:       strID,
			Type:     interfaces.App,
			Language: interfaces.SimplifiedChinese,
		}
		Convey("checkGroupGetAuth2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 0)

			assert.Equal(t, err, testErr)
		})

		appOrgPern := interfaces.AppOrgPerm{
			Value: interfaces.Read,
		}
		perms := map[interfaces.OrgType]interfaces.AppOrgPerm{
			interfaces.Group: appOrgPern,
		}
		var groupInfo interfaces.GroupInfo
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("GetGroupByID2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 0)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID2 id is empty", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 0)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.URINotExist, "group not exist"))
		})

		groupInfo.ID = strID
		userInfo := make([]interfaces.UserDBInfo, 0)
		Convey("SearchUserInfoByName error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, testErr)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 0)

			assert.Equal(t, err, testErr)
		})

		Convey("SearchUserInfoByName empty", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

			num, out, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 0)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 0)
			assert.Equal(t, num, 0)
		})

		tempUser := interfaces.UserDBInfo{
			ID:   strID,
			Name: strID2,
		}
		userInfo = append(userInfo, tempUser)
		Convey("GetGroupMembersByGroupIDs2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey("GetUsersPath2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, _, _, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 1, 1)

			assert.Equal(t, err, testErr)
		})

		D1NameInfo := interfaces.NameInfo{ID: "d1d1", Name: "d1"}
		D2NameInfo := interfaces.NameInfo{ID: "d2d1", Name: "d2"}
		D3NameInfo := interfaces.NameInfo{ID: "d3d1", Name: "d3"}
		D4NameInfo := interfaces.NameInfo{ID: "d4d1", Name: "d4"}
		D5NameInfo := interfaces.NameInfo{ID: "d5d1", Name: "d5"}
		D6NameInfo := interfaces.NameInfo{ID: "d6d1", Name: "d6"}

		member1 := interfaces.GroupMemberInfo{
			ID:         strID2,
			MemberType: 1,
		}
		member2 := interfaces.GroupMemberInfo{
			ID:         D6NameInfo.ID,
			MemberType: 2,
		}

		member3 := interfaces.GroupMemberInfo{
			ID:         strID,
			MemberType: 1,
		}

		memberInfos := make(map[string][]interfaces.GroupMemberInfo)
		memberInfos[strID] = []interfaces.GroupMemberInfo{member1, member2, member3}

		tempPath := []string{
			D4NameInfo.ID + "/" + D5NameInfo.ID,
			D6NameInfo.ID + "/" + D3NameInfo.ID + "/" + D1NameInfo.ID,
		}
		userPaths := map[string][]string{
			strID: tempPath,
		}

		D6Info := interfaces.DepartmentDBInfo{
			ID:   D6NameInfo.ID,
			Name: D6NameInfo.Name,
			Path: D2NameInfo.ID + "/" + D1NameInfo.ID + "/" + D6NameInfo.ID,
		}

		departInfos := []interfaces.DepartmentDBInfo{D6Info}

		temp1 := interfaces.DepartmentDBInfo{
			ID:   D1NameInfo.ID,
			Name: D1NameInfo.Name,
		}
		temp2 := interfaces.DepartmentDBInfo{
			ID:   D2NameInfo.ID,
			Name: D2NameInfo.Name,
		}
		temp3 := interfaces.DepartmentDBInfo{
			ID:   D3NameInfo.ID,
			Name: D3NameInfo.Name,
		}
		temp4 := interfaces.DepartmentDBInfo{
			ID:   D4NameInfo.ID,
			Name: D4NameInfo.Name,
		}
		temp5 := interfaces.DepartmentDBInfo{
			ID:   D5NameInfo.ID,
			Name: D5NameInfo.Name,
		}
		temp6 := interfaces.DepartmentDBInfo{
			ID:   D6NameInfo.ID,
			Name: D6NameInfo.Name,
		}
		allDepartInfos := []interfaces.DepartmentDBInfo{temp1, temp2, temp3, temp4, temp5, temp6}

		// 只有一个user，有两个members存在  member3和member2 因为paths包含 handlePathParams的userids为userid为userInfo.ID， 部门id是departmentid为D6
		// alldepartInfo 包括
		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(memberInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(departInfos, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(allDepartInfos, nil)
			num, userIDs, uInfox, mInfosx, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 4)

			uInfo := uInfox[strID]
			mInfos := mInfosx[strID]
			assert.Equal(t, err, nil)
			assert.Equal(t, num, 1)

			assert.Equal(t, len(userIDs), 1)
			assert.Equal(t, userIDs[0], strID)

			assert.Equal(t, uInfo.ID, strID)
			assert.Equal(t, len(uInfo.ParentDeps), 2)
			assert.Equal(t, len(uInfo.ParentDeps[0]), 2)
			assert.Equal(t, len(uInfo.ParentDeps[1]), 3)

			assert.Equal(t, uInfo.ParentDeps[0][0], D4NameInfo)
			assert.Equal(t, uInfo.ParentDeps[0][1], D5NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][0], D6NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][1], D3NameInfo)
			assert.Equal(t, uInfo.ParentDeps[1][2], D1NameInfo)

			assert.Equal(t, len(mInfos), 2)
			assert.Equal(t, mInfos[0].ID, D6NameInfo.ID)
			assert.Equal(t, mInfos[0].Name, D6NameInfo.Name)
			assert.Equal(t, len(mInfos[0].ParentDeps), 1)
			assert.Equal(t, len(mInfos[0].ParentDeps[0]), 2)
			assert.Equal(t, mInfos[0].ParentDeps[0][0], D2NameInfo)
			assert.Equal(t, mInfos[0].ParentDeps[0][1], D1NameInfo)

			assert.Equal(t, mInfos[1].ID, strID)
			assert.Equal(t, mInfos[1].Name, strID2)
			assert.Equal(t, len(mInfos[1].ParentDeps), 2)
			assert.Equal(t, len(mInfos[1].ParentDeps[0]), 2)
			assert.Equal(t, len(mInfos[1].ParentDeps[1]), 3)
			assert.Equal(t, mInfos[1].ParentDeps[0][0], D4NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[0][1], D5NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][0], D6NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][1], D3NameInfo)
			assert.Equal(t, mInfos[1].ParentDeps[1][2], D1NameInfo)
		})

		member11 := interfaces.GroupMemberInfo{
			ID:         strID,
			MemberType: 1,
		}
		memberInfos1 := make(map[string][]interfaces.GroupMemberInfo)
		memberInfos1[strID] = []interfaces.GroupMemberInfo{member11}

		tempPath1 := []string{"-1"}
		userPaths1 := map[string][]string{
			strID: tempPath1,
		}
		temp11 := interfaces.DepartmentDBInfo{
			ID:   "-1",
			Name: "",
		}
		allDepartInfos1 := []interfaces.DepartmentDBInfo{temp11}
		Convey("success undistrbute group", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(memberInfos1, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths1, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(userPaths1, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(departInfos, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(allDepartInfos1, nil)
			num, userIDs, uInfox, mInfosx, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 4)

			assert.Equal(t, err, nil)
			assert.Equal(t, num, 1)

			assert.Equal(t, len(userIDs), 1)
			assert.Equal(t, userIDs[0], strID)

			uInfo := uInfox[strID]
			mInfos := mInfosx[strID]
			assert.Equal(t, uInfo.ID, strID)
			assert.Equal(t, uInfo.Name, userInfo[0].Name)
			assert.Equal(t, len(uInfo.ParentDeps), 1)
			assert.Equal(t, len(uInfo.ParentDeps[0]), 1)

			tempName1 := group.i18n.Load(i18nIDObjectsInUnDistributeUserGroup, vistor.Language)
			assert.Equal(t, uInfo.ParentDeps[0][0].ID, "-1")
			assert.Equal(t, uInfo.ParentDeps[0][0].Name, tempName1)

			assert.Equal(t, len(mInfos), 1)
			assert.Equal(t, mInfos[0].ID, member11.ID)
			assert.Equal(t, mInfos[0].Name, strID2)
			assert.Equal(t, len(mInfos[0].ParentDeps), 1)
			assert.Equal(t, mInfos[0].ParentDeps[0][0].ID, "-1")
			assert.Equal(t, mInfos[0].ParentDeps[0][0].Name, tempName1)
		})

		tempUser1 := interfaces.UserDBInfo{
			ID:   strID,
			Name: strName,
		}
		tempUser2 := interfaces.UserDBInfo{
			ID:   strID1,
			Name: strName1,
		}
		userInfo2 := []interfaces.UserDBInfo{tempUser1, tempUser2}

		member21 := interfaces.GroupMemberInfo{
			ID:         strID,
			Name:       strName,
			MemberType: 1,
		}
		member22 := interfaces.GroupMemberInfo{
			ID:         strID1,
			Name:       strName1,
			MemberType: 1,
		}
		memberInfos2 := make(map[string][]interfaces.GroupMemberInfo)
		memberInfos2[strID] = []interfaces.GroupMemberInfo{member21, member22}

		Convey("success 2 user out", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(perms, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().SearchUserInfoByName(gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo2, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(memberInfos2, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(nil, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			num, userIDs, uInfox, mInfosx, err := group.SearchInAllGroupOrg(ctx, &vistor, strID, strID2, 0, 4)

			assert.Equal(t, err, nil)
			assert.Equal(t, num, 2)

			assert.Equal(t, len(userIDs), 2)
			assert.Equal(t, userIDs[0], strID)
			assert.Equal(t, userIDs[1], strID1)

			uInfo := uInfox[strID]
			mInfos := mInfosx[strID]
			assert.Equal(t, uInfo.ID, strID)
			assert.Equal(t, uInfo.Name, userInfo2[0].Name)
			assert.Equal(t, len(uInfo.ParentDeps), 0)

			assert.Equal(t, len(mInfos), 1)
			assert.Equal(t, mInfos[0].ID, member21.ID)
			assert.Equal(t, mInfos[0].Name, member21.Name)
			assert.Equal(t, len(mInfos[0].ParentDeps), 0)

			uInfo = uInfox[strID1]
			mInfos = mInfosx[strID1]
			assert.Equal(t, uInfo.ID, strID1)
			assert.Equal(t, uInfo.Name, userInfo2[1].Name)
			assert.Equal(t, len(uInfo.ParentDeps), 0)

			assert.Equal(t, len(mInfos), 1)
			assert.Equal(t, mInfos[0].ID, member22.ID)
			assert.Equal(t, mInfos[0].Name, member22.Name)
			assert.Equal(t, len(mInfos[0].ParentDeps), 0)
		})
	})
}

//nolint:funlen
func TestAddGroup(t *testing.T) {
	Convey("AddGroup, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := sqlDB.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.trace = trace
		groupMember.pool = sqlDB
		groupMember.ob = ob
		ctx := context.Background()

		common.SvcConfig.Lang = "zh_CN"

		v := &interfaces.Visitor{}
		v.Type = interfaces.RealName
		Convey("GetRolesByUserID err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := groupMember.AddGroup(ctx, v, "", "", nil)

			assert.Equal(t, err, testErr)
		})

		roleInfos := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfos[v.ID] = userRole
		Convey("user dont have the role to manager", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr1 := rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority")
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			_, err := groupMember.AddGroup(ctx, v, "", "", nil)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleSysAdmin] = true
		Convey("checkName err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			_, err := groupMember.AddGroup(ctx, v, "", "", nil)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "param name is illegal"))
		})

		notes := `xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx1111111111xxxxxxxxxx1111111111
		xxxxxxxxxx11111111111`
		Convey("note too long", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			_, err := groupMember.AddGroup(ctx, v, "zzz", notes, nil)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "param notes is illegal"))
		})

		notes = strID1
		name := strID1
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey("GetGroupIDByName err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			_, err := groupMember.AddGroup(ctx, v, name, notes, nil)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupIDByName failed", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return(strID1, nil)
			_, err := groupMember.AddGroup(ctx, v, name, notes, nil)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.Conflict, "this group name is existing",
				rest.SetCodeStr(errors.StrConflictGroup)))
		})

		initalGroupIDs := []string{strID1, strID1}
		Convey("group id  duplicate", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "initial group ids duplicate"))
		})
		initalGroupIDs = []string{strID1}
		Convey("GetGroupName2  error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("group id not exist  error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil, nil)
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "initial group ids not exist"))
		})

		Convey("GetGroupMembersByGroupIDs2 error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("dbpool begin error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			txMock.ExpectBegin().WillReturnError(testErr)
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("AddGroup  error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			txMock.ExpectBegin()
			groupDB.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("AddGroupMembers  error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			txMock.ExpectBegin()
			groupDB.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupMemberDB.EXPECT().AddGroupMembers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("outbox error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			txMock.ExpectBegin()
			groupDB.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupMemberDB.EXPECT().AddGroupMembers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfos, nil)
			groupDB.EXPECT().GetGroupIDByName2(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().GetGroupName2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, initalGroupIDs, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			txMock.ExpectBegin()
			groupDB.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			groupMemberDB.EXPECT().AddGroupMembers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()
			_, err := groupMember.AddGroup(ctx, v, name, notes, initalGroupIDs)

			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteGroup(t *testing.T) {
	Convey("DeleteGroup, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.ob = ob
		groupMember.pool = sqlDB

		v := &interfaces.Visitor{}
		v.Type = interfaces.RealName

		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo[v.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleSysAdmin] = true
		Convey("GetGroupByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, testErr)
		})

		Convey("DeleteGroupAllMembersByGroupID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().DeleteGroupMemberByID(gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, testErr)
		})

		Convey("DeleteGroup err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().DeleteGroupMemberByID(gomock.Any()).AnyTimes().Return(nil)
			groupDB.EXPECT().DeleteGroup(gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, testErr)
		})

		Convey("DeleteGroup Success", func() {
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().DeleteGroupMemberByID(gomock.Any()).AnyTimes().Return(nil)
			groupDB.EXPECT().DeleteGroup(gomock.Any()).AnyTimes().Return(nil)
			msg.EXPECT().Publish(gomock.Any(), gomock.Any()).AnyTimes()
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()
			err := groupMember.DeleteGroup(v, "")

			assert.Equal(t, err, nil)
		})
	})
}

func TestModifyGroup(t *testing.T) {
	Convey("ModifyGroup, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.pool = sqlDB
		groupMember.ob = ob

		v := &interfaces.Visitor{
			Language: interfaces.SimplifiedChinese,
		}
		v.Type = interfaces.RealName
		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo[v.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleSuperAdmin] = true
		Convey("GetGroupByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID failed", func() {
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound)))
		})

		Convey("GetGroupIDByName err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupDB.EXPECT().GetGroupIDByName(gomock.Any()).AnyTimes().Return("", testErr)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupIDByName failed", func() {
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupDB.EXPECT().GetGroupIDByName(gomock.Any()).AnyTimes().Return("xxxx", nil)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.Conflict, "the name of this group is existing",
				rest.SetCodeStr(errors.StrConflictGroup)))
		})

		Convey("ModifyGroup failed", func() {
			testErr := rest.NewHTTPErrorV2(errors.Conflict, "the name of this group is existing",
				rest.SetCodeStr(errors.StrConflictGroup))
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupDB.EXPECT().GetGroupIDByName(gomock.Any()).AnyTimes().Return("xxxx", nil)
			groupDB.EXPECT().ModifyGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, testErr)
		})

		Convey("ModifyGroup Success", func() {
			groupInfo := interfaces.GroupInfo{ID: "xxxxyy"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupDB.EXPECT().GetGroupIDByName(gomock.Any()).AnyTimes().Return("", nil)
			groupDB.EXPECT().ModifyGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			msg.EXPECT().Publish(gomock.Any(), gomock.Any()).AnyTimes()
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()

			err := groupMember.ModifyGroup(v, "", "", false, "", false)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetGroup(t *testing.T) {
	Convey("GetGroup, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)

		searchInfo := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    0,
			Limit:     50,
			Keyword:   "xxx",
		}

		var visitor interfaces.Visitor
		visitor.Type = interfaces.RealName
		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroup(&visitor, searchInfo)

			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo[visitor.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			_, _, err := groupMember.GetGroup(&visitor, searchInfo)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleAuditAdmin] = true
		Convey("GetGroupsNum err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(0, testErr)
			_, _, err := groupMember.GetGroup(&visitor, searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupsNum 0", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(0, nil)
			num, outInfo, err := groupMember.GetGroup(&visitor, searchInfo)

			assert.Equal(t, num, 0)
			assert.Equal(t, err, nil)
			assert.Equal(t, outInfo, nil)
		})

		Convey("GetGroups err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(1, nil)
			groupDB.EXPECT().GetGroups(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroup(&visitor, searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroups Success", func() {
			groups := make([]interfaces.GroupInfo, 0)
			group1 := interfaces.GroupInfo{}
			group2 := interfaces.GroupInfo{}
			groups = append(groups, group1, group2)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(10, nil)
			groupDB.EXPECT().GetGroups(gomock.Any()).AnyTimes().Return(groups, nil)
			num, outData, err := groupMember.GetGroup(&visitor, searchInfo)

			var oList []interfaces.GroupInfo
			oList = append(oList, group1, group2)
			assert.Equal(t, err, nil)
			assert.Equal(t, outData, oList)
			assert.Equal(t, 10, num)
		})
	})
}

func TestGetGroupByID(t *testing.T) {
	Convey("GetGroupByID, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)

		var visitor interfaces.Visitor
		visitor.Type = interfaces.RealName
		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := groupMember.GetGroupByID(&visitor, "")

			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo[visitor.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			_, err := groupMember.GetGroupByID(&visitor, "")

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleAuditAdmin] = true
		Convey("GetGroupByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			_, err := groupMember.GetGroupByID(&visitor, "")

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID Success", func() {
			groupInfo := interfaces.GroupInfo{
				ID:   "xxxx",
				Name: "xxxtt",
			}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			out, err := groupMember.GetGroupByID(&visitor, "")

			assert.Equal(t, err, nil)
			assert.Equal(t, out, groupInfo)
		})
	})
}

func TestDeleteGroupMembers(t *testing.T) {
	Convey("DeleteGroupMembers, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.ob = ob
		groupMember.pool = sqlDB

		v := &interfaces.Visitor{
			ID:       "xxx",
			Language: interfaces.SimplifiedChinese,
		}
		v.Type = interfaces.RealName
		infos := make(map[string]interfaces.GroupMemberInfo)
		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		roleInfo := make(map[string]map[interfaces.Role]bool)
		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo[v.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleSuperAdmin] = true
		Convey("GetGroupByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID true", func() {
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound)))
		})

		Convey("DeleteGroupMember err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			data := interfaces.GroupMemberInfo{}
			infos["xxxx"] = data
			groupInfo := interfaces.GroupInfo{ID: "xxxs"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().DeleteGroupMember(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("DeleteGroupMember success", func() {
			groupInfo := interfaces.GroupInfo{ID: "xxxs"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().DeleteGroupMember(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()

			err := groupMember.DeleteGroupMembers(v, "", infos)

			assert.Equal(t, err, nil)
		})
	})
}

func TestAddGroupMembers(t *testing.T) {
	Convey("AddGroupMembers, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.ob = ob

		infos := make(map[string]interfaces.GroupMemberInfo)
		oneUser := interfaces.GroupMemberInfo{
			ID:              "xxxxid",
			MemberType:      1,
			DepartmentNames: []string{"xxxdepart"},
		}
		oneDepartment := interfaces.GroupMemberInfo{
			ID:              "yyyyid",
			MemberType:      2,
			DepartmentNames: []string{"yyydepart"},
		}
		infos["xxxxid"] = oneUser
		infos["yyyyid"] = oneDepartment

		v := &interfaces.Visitor{}
		v.Type = interfaces.RealName
		testID := "xxxxid"

		Convey("GetRolesByUserID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo := make(map[string]map[interfaces.Role]bool)
		roleInfo[v.ID] = userRole
		Convey("user dont have the role to manager", func() {
			testErr1 := rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleSuperAdmin] = true
		Convey("GetGroupByID err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID fail", func() {
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound)))
		})

		Convey("CheckUsersExist err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("CheckUsersExist fail", func() {
			outUserIDs := make([]interfaces.UserDBInfo, 0)
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(outUserIDs, nil, nil)
			err := groupMember.AddGroupMembers(v, "", infos)

			outInfo := make(map[string]interface{})
			outInfo["ids"] = []string{testID}

			test1Err := rest.NewHTTPErrorV2(errors.UserNotFound,
				groupMember.i18n.Load(i18nIDObjectsInUserNotFound, v.Language),
				rest.SetDetail(map[string]interface{}{"ids": []string{testID}}),
				rest.SetCodeStr(errors.StrBadRequestUserNotFound),
			)
			assert.Equal(t, err, test1Err)
		})

		Convey("CheckDepartmentsExist err", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			outInfo := make([]interfaces.UserDBInfo, 0)
			nameInfo := interfaces.UserDBInfo{
				ID:   "xxxxid",
				Name: "sdad",
			}
			outInfo = append(outInfo, nameInfo)
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(outInfo, nil, nil)
			departmentDB.EXPECT().GetDepartmentName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("CheckDepartmentsExist fail", func() {
			outInfo := make([]interfaces.UserDBInfo, 0)
			nameInfo := interfaces.UserDBInfo{
				ID:   "xxxxid",
				Name: "sdad",
			}
			outInfo = append(outInfo, nameInfo)
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return(outInfo, nil, nil)
			departmentDB.EXPECT().GetDepartmentName(gomock.Any()).AnyTimes().Return(nil, nil, nil)
			err := groupMember.AddGroupMembers(v, "", infos)

			outInfo2 := make(map[string]interface{})
			outInfo2["ids"] = []string{"yyyyid"}

			test2Err := rest.NewHTTPErrorV2(errors.DepartmentNotFound,
				groupMember.i18n.Load(i18nIDObjectsInDepartNotFound, v.Language),
				rest.SetDetail(outInfo2),
				rest.SetCodeStr(errors.StrBadRequestDepartmentNotFound))
			assert.Equal(t, err, test2Err)
		})

		Convey("CheckGroupMembersExist fail", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			nameInfo1 := interfaces.UserDBInfo{
				ID:   "xxxxid",
				Name: "sdad",
			}
			nameInfo2 := interfaces.NameInfo{
				ID:   "yyyyid",
				Name: "sdad",
			}
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{nameInfo1}, nil, nil)
			departmentDB.EXPECT().GetDepartmentName(gomock.Any()).AnyTimes().Return([]interfaces.NameInfo{nameInfo2}, nil, nil)
			groupMemberDB.EXPECT().CheckGroupMembersExist(gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)
			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})

		Convey("AddGroupMember fail", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			nameInfo1 := interfaces.UserDBInfo{
				ID:   "xxxxid",
				Name: "sdad",
			}
			nameInfo2 := interfaces.NameInfo{
				ID:   "yyyyid",
				Name: "sdad",
			}
			groupInfo := interfaces.GroupInfo{ID: "xxx"}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{nameInfo1}, nil, nil)
			departmentDB.EXPECT().GetDepartmentName(gomock.Any()).AnyTimes().Return([]interfaces.NameInfo{nameInfo2}, nil, nil)
			groupMemberDB.EXPECT().CheckGroupMembersExist(gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			groupMemberDB.EXPECT().AddGroupMember(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, testErr)
		})
	})
}

func TestAddGroupMembersTrue(t *testing.T) {
	Convey("AddGroupMembers, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.ob = ob
		groupMember.pool = sqlDB

		infos := make(map[string]interfaces.GroupMemberInfo)
		oneUser := interfaces.GroupMemberInfo{
			ID:              "xxxxid",
			MemberType:      1,
			DepartmentNames: []string{"xxxdepart"},
		}
		oneDepartment := interfaces.GroupMemberInfo{
			ID:              "yyyyid",
			MemberType:      2,
			DepartmentNames: []string{"yyydepart"},
		}
		infos["xxxxid"] = oneUser
		infos["yyyyid"] = oneDepartment

		v := &interfaces.Visitor{}
		v.Type = interfaces.RealName

		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleSysAdmin] = true
		roleInfo := make(map[string]map[interfaces.Role]bool)
		roleInfo[v.ID] = userRole

		Convey("DeleteGroupMember success", func() {
			nameInfo1 := interfaces.UserDBInfo{
				ID:   "xxxxid",
				Name: "sdad",
			}
			nameInfo2 := interfaces.NameInfo{
				ID:   "yyyyid",
				Name: "sdad",
			}
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			userDB.EXPECT().GetUserName(gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{nameInfo1}, nil, nil)
			departmentDB.EXPECT().GetDepartmentName(gomock.Any()).AnyTimes().Return([]interfaces.NameInfo{nameInfo2}, nil, nil)
			groupMemberDB.EXPECT().CheckGroupMembersExist(gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
			groupMemberDB.EXPECT().AddGroupMember(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().AnyTimes()

			err := groupMember.AddGroupMembers(v, "", infos)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetGroupMembersID(t *testing.T) {
	Convey("GetGroupMembersID, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)

		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)
		groupIDs := make([]string, 0)
		visitor := interfaces.Visitor{}

		Convey("groupIDs is empty", func() {
			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, true)
			assert.Equal(t, len(userIDs), 0)
			assert.Equal(t, len(departmentIDs), 0)
			assert.Equal(t, err, nil)
		})

		groupIDs = append(groupIDs, "group_id")
		Convey("DB GetExistGroupIDs error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(nil, testErr)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, true)
			assert.Equal(t, userIDs, nil)
			assert.Equal(t, departmentIDs, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("group id is not exist", func() {
			testErr := rest.NewHTTPErrorV2(errors.GroupNotFound,
				groupMember.i18n.Load(i18nIDObjectsInGroupNotFound, visitor.Language),
				rest.SetDetail(map[string]interface{}{"ids": []string{"group_id"}}),
				rest.SetCodeStr(errors.StrBadRequestGroupNotFound))
			testExistGroupIDs := make([]string, 0)
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, true)
			assert.Equal(t, userIDs, nil)
			assert.Equal(t, departmentIDs, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("DB GetGroupMembersByGroupIDs error", func() {
			testExistGroupIDs := []string{"group_id"}
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs(gomock.Any()).AnyTimes().Return(nil, testErr)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, true)
			assert.Equal(t, userIDs, nil)
			assert.Equal(t, departmentIDs, nil)
			assert.Equal(t, err, testErr)
		})

		Convey("show unbale user success", func() {
			testExistGroupIDs := []string{"group_id"}
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			testGroupMemberInfos := make([]interfaces.GroupMemberInfo, 0)
			groupMemberInfo1 := interfaces.GroupMemberInfo{
				ID:              "user_member_id",
				MemberType:      1,
				Name:            "user_name",
				DepartmentNames: []string{},
			}
			groupMemberInfo2 := interfaces.GroupMemberInfo{
				ID:              "group_member_id",
				MemberType:      2,
				Name:            "group_name",
				DepartmentNames: []string{},
			}
			testGroupMemberInfos = append(testGroupMemberInfos, groupMemberInfo1, groupMemberInfo2)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs(gomock.Any()).AnyTimes().Return(testGroupMemberInfos, nil)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, true)
			assert.Equal(t, len(userIDs), 1)
			assert.Equal(t, len(departmentIDs), 1)
			assert.Equal(t, err, nil)
		})

		Convey("GetUserDBInfo error", func() {
			testExistGroupIDs := []string{"group_id"}
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			testGroupMemberInfos := make([]interfaces.GroupMemberInfo, 0)
			groupMemberInfo1 := interfaces.GroupMemberInfo{
				ID:              "user_member_id",
				MemberType:      1,
				Name:            "user_name",
				DepartmentNames: []string{},
			}
			groupMemberInfo2 := interfaces.GroupMemberInfo{
				ID:              "group_member_id",
				MemberType:      2,
				Name:            "group_name",
				DepartmentNames: []string{},
			}
			testErr := rest.NewHTTPError("error", 503000000, nil)
			testGroupMemberInfos = append(testGroupMemberInfos, groupMemberInfo1, groupMemberInfo2)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs(gomock.Any()).AnyTimes().Return(testGroupMemberInfos, nil)
			userDB.EXPECT().GetUserDBInfo(gomock.Any()).AnyTimes().Return(nil, testErr)

			_, _, err := groupMember.GetGroupMembersID(&visitor, groupIDs, false)
			assert.Equal(t, err, testErr)
		})

		groupMemberInfo1 := interfaces.GroupMemberInfo{
			ID:              "user_member_id",
			MemberType:      1,
			Name:            "user_name",
			DepartmentNames: []string{},
		}
		groupMemberInfo2 := interfaces.GroupMemberInfo{
			ID:              "group_member_id",
			MemberType:      2,
			Name:            "group_name",
			DepartmentNames: []string{},
		}

		user1 := interfaces.UserDBInfo{
			ID:                groupMemberInfo1.ID,
			DisableStatus:     interfaces.Enabled,
			AutoDisableStatus: interfaces.AEnabled,
		}

		Convey("no show unbale user success", func() {
			testExistGroupIDs := []string{"group_id"}
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			testGroupMemberInfos := make([]interfaces.GroupMemberInfo, 0)
			testGroupMemberInfos = append(testGroupMemberInfos, groupMemberInfo1, groupMemberInfo2)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs(gomock.Any()).AnyTimes().Return(testGroupMemberInfos, nil)
			userDB.EXPECT().GetUserDBInfo(gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{user1}, nil)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, false)
			assert.Equal(t, len(userIDs), 1)
			assert.Equal(t, len(departmentIDs), 1)
			assert.Equal(t, err, nil)
		})

		Convey("no show unbale user success, no user", func() {
			testExistGroupIDs := []string{"group_id"}
			groupDB.EXPECT().GetExistGroupIDs(gomock.Any()).AnyTimes().Return(testExistGroupIDs, nil)

			testGroupMemberInfos := make([]interfaces.GroupMemberInfo, 0)
			user2 := interfaces.UserDBInfo{
				ID:                groupMemberInfo1.ID,
				DisableStatus:     interfaces.Disabled,
				AutoDisableStatus: interfaces.ADisabled,
			}
			testGroupMemberInfos = append(testGroupMemberInfos, groupMemberInfo1, groupMemberInfo2)
			groupMemberDB.EXPECT().GetGroupMembersByGroupIDs(gomock.Any()).AnyTimes().Return(testGroupMemberInfos, nil)
			userDB.EXPECT().GetUserDBInfo(gomock.Any()).AnyTimes().Return([]interfaces.UserDBInfo{user2}, nil)

			userIDs, departmentIDs, err := groupMember.GetGroupMembersID(&visitor, groupIDs, false)
			assert.Equal(t, len(userIDs), 0)
			assert.Equal(t, len(departmentIDs), 1)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetGroupMembers(t *testing.T) {
	Convey("GetGroupMembers, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.trace = trace

		searchInfo := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    0,
			Limit:     50,
			Keyword:   "xxx",
		}

		var visitor interfaces.Visitor
		visitor.Type = interfaces.RealName
		visitor.Language = interfaces.SimplifiedChinese
		var ctx context.Context
		Convey("GetRolesByUserID err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		userRole := make(map[interfaces.Role]bool)
		userRole[interfaces.SystemRoleNormalUser] = true
		roleInfo := make(map[string]map[interfaces.Role]bool)
		roleInfo[visitor.ID] = userRole
		Convey("user dont have the role to manager", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr1 := rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority")
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr1)
		})

		userRole[interfaces.SystemRoleAuditAdmin] = true
		Convey("GetGroupByID err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupByID fail", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			groupInfo := interfaces.GroupInfo{}
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.NotFound, "this group is not existing"))
		})

		Convey("GetGroupMembersNum err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(0, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetGroupMembers err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(10, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		memberInfo := interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      1,
			DepartmentNames: []string{"xxxx"},
		}
		memberInfo1 := interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      2,
			DepartmentNames: []string{"xxxx"},
		}
		outInfos := make([]interfaces.GroupMemberInfo, 0)
		outInfos = append(outInfos, memberInfo, memberInfo1)
		Convey("GetUsersPath2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(10, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetDepartmentInfo2 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(10, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})

		Convey("GetDepartmentInfo22 err", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupInfo := interfaces.GroupInfo{ID: "xxxx"}

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(10, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, testErr)
			_, _, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			assert.Equal(t, err, testErr)
		})
	})
}

func TestGetGroupMembersSuccess(t *testing.T) {
	Convey("GetGroupMembers, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, role)
		groupMember.trace = trace

		searchInfo := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    0,
			Limit:     50,
			Keyword:   "xxx",
		}

		var visitor interfaces.Visitor
		visitor.Type = interfaces.RealName
		visitor.ID = strID
		Convey("Success", func() {
			memberInfo := interfaces.GroupMemberInfo{
				ID:              "xxxx1",
				MemberType:      1,
				DepartmentNames: []string{"xxxx"},
			}
			memberInfo1 := interfaces.GroupMemberInfo{
				ID:              "xxxx2",
				MemberType:      2,
				DepartmentNames: []string{"xxxx"},
			}
			var ctx context.Context

			userRole := make(map[interfaces.Role]bool)
			userRole[interfaces.SystemRoleSuperAdmin] = true
			roleInfo := make(map[string]map[interfaces.Role]bool)
			roleInfo[visitor.ID] = userRole

			outInfos := make([]interfaces.GroupMemberInfo, 0)
			outInfos = append(outInfos, memberInfo, memberInfo1)

			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			groupInfo := interfaces.GroupInfo{ID: "xxxx"}

			userPaths := make(map[string][]string)
			userPaths[memberInfo.ID] = []string{"D1ID/D2ID/D3ID", "D4ID"}

			depInfo := interfaces.DepartmentDBInfo{
				ID:   memberInfo1.ID,
				Path: "D1ID/D4ID/D5ID/" + memberInfo1.ID,
			}

			nameInfo1 := interfaces.DepartmentDBInfo{ID: "D1ID", Name: "D1"}
			nameInfo2 := interfaces.DepartmentDBInfo{ID: "D2ID", Name: "D2"}
			nameInfo3 := interfaces.DepartmentDBInfo{ID: "D3ID", Name: "D3"}
			nameInfo4 := interfaces.DepartmentDBInfo{ID: "D4ID", Name: "D4"}
			nameInfo5 := interfaces.DepartmentDBInfo{ID: "D5ID", Name: "D5"}
			nameInfo6 := interfaces.DepartmentDBInfo{ID: memberInfo1.ID, Name: memberInfo1.ID}

			depNameInfos := []interfaces.DepartmentDBInfo{nameInfo1, nameInfo2, nameInfo3, nameInfo4, nameInfo5, nameInfo6}

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(16, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(userPaths, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return([]interfaces.DepartmentDBInfo{depInfo}, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(depNameInfos, nil)
			num, memberInfos, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			name1 := interfaces.NameInfo{ID: "D1ID", Name: "D1"}
			name2 := interfaces.NameInfo{ID: "D2ID", Name: "D2"}
			name3 := interfaces.NameInfo{ID: "D3ID", Name: "D3"}
			name4 := interfaces.NameInfo{ID: "D4ID", Name: "D4"}
			name5 := interfaces.NameInfo{ID: "D5ID", Name: "D5"}
			path1 := []interfaces.NameInfo{name1, name2, name3}
			path2 := []interfaces.NameInfo{name4}
			path3 := []interfaces.NameInfo{name1, name4, name5}

			var oList []interfaces.GroupMemberInfo
			memberInfo.DepartmentNames = []string{"D3", "D4"}
			memberInfo.ParentDeps = [][]interfaces.NameInfo{path1, path2}
			memberInfo1.DepartmentNames = []string{"D5"}
			memberInfo1.ParentDeps = [][]interfaces.NameInfo{path3}
			oList = append(oList, memberInfo, memberInfo1)

			assert.Equal(t, err, nil)
			assert.Equal(t, memberInfos, oList)
			assert.Equal(t, num, 16)
		})

		Convey("Success undistribute group", func() {
			memberInfo := interfaces.GroupMemberInfo{
				ID:         "xxxx1",
				MemberType: 1,
			}
			var ctx context.Context

			userRole := make(map[interfaces.Role]bool)
			userRole[interfaces.SystemRoleSuperAdmin] = true
			roleInfo := make(map[string]map[interfaces.Role]bool)
			roleInfo[visitor.ID] = userRole

			outInfos := make([]interfaces.GroupMemberInfo, 0)
			outInfos = append(outInfos, memberInfo)

			groupInfo := interfaces.GroupInfo{ID: "xxxx"}

			userPaths := make(map[string][]string)
			userPaths[memberInfo.ID] = []string{"-1"}

			nameInfo1 := interfaces.DepartmentDBInfo{ID: "-1", Name: ""}

			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(roleInfo, nil)
			groupDB.EXPECT().GetGroupByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum2(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(16, nil)
			groupMemberDB.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(outInfos, nil)
			userDB.EXPECT().GetUsersPath2(gomock.Any(), gomock.Any()).AnyTimes().Return(userPaths, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return(nil, nil)
			departmentDB.EXPECT().GetDepartmentInfoByIDs(gomock.Any(), gomock.Any()).Return([]interfaces.DepartmentDBInfo{nameInfo1}, nil)
			num, memberInfos, err := groupMember.GetGroupMembers(ctx, &visitor, "", searchInfo)

			var oList []interfaces.GroupMemberInfo
			tempData := groupMember.i18n.Load(i18nIDObjectsInUnDistributeUserGroup, visitor.Language)
			tempNames1 := interfaces.NameInfo{ID: "-1", Name: tempData}
			path1 := []interfaces.NameInfo{tempNames1}
			memberInfo.DepartmentNames = []string{tempData}
			memberInfo.ParentDeps = [][]interfaces.NameInfo{path1}
			oList = append(oList, memberInfo)

			assert.Equal(t, err, nil)
			assert.Equal(t, memberInfos, oList)
			assert.Equal(t, num, 16)
		})
	})
}

func TestSearchGroupByKey(t *testing.T) {
	Convey("SearchGroupByKeyword, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)

		Convey(" SearchGroupByKeyword Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().SearchGroupByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := groupMember.SearchGroupByKeyword("", 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			info1 := interfaces.NameInfo{
				ID:   "xxxx",
				Name: "xxxx1",
			}
			info2 := interfaces.NameInfo{
				ID:   "yyyy",
				Name: "yyyy1",
			}
			tmp := []interfaces.NameInfo{
				info1,
				info2,
			}
			groupDB.EXPECT().SearchGroupByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tmp, nil)
			outInfos, err := groupMember.SearchGroupByKeyword("", 1, 1)

			assert.Equal(t, err, nil)
			assert.Equal(t, outInfos[0], info1)
			assert.Equal(t, outInfos[1], info2)
			assert.Equal(t, len(outInfos), 2)
		})
	})
}

func TestSearchGroupNumByKeyword(t *testing.T) {
	Convey("SearchGroupNumByKeyword, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		gDB := mock.NewMockDBGroup(ctrl)
		gmDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(gmDB, gDB, userDB, departmentDB, eacplog, msg, nil)

		Convey(" SearchGroupNumByKeyword Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			gDB.EXPECT().SearchGroupNumByKeyword(gomock.Any()).AnyTimes().Return(0, testErr)
			_, err := groupMember.SearchGroupNumByKeyword("")

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			gDB.EXPECT().SearchGroupNumByKeyword(gomock.Any()).AnyTimes().Return(2, nil)
			outInfos, err := groupMember.SearchGroupNumByKeyword("")

			assert.Equal(t, err, nil)
			assert.Equal(t, outInfos, 2)
		})

		Convey(" for golangci", func() {
			assert.Equal(t, 1, 1)
		})
	})
}

func TestSearchMembersByKeyword(t *testing.T) {
	Convey("SearchMembersByKeyword, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)

		Convey(" SearchGroupNumByKeyword Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupMemberDB.EXPECT().SearchMembersByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := groupMember.SearchMembersByKeyword("", 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			test1 := interfaces.MemberInfo{
				ID:         "zzz",
				Name:       "kkkk",
				NType:      1,
				GroupNames: []string{"xxxx"},
			}
			groupMemberDB.EXPECT().SearchMembersByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.MemberInfo{test1}, nil)
			outInfos, err := groupMember.SearchMembersByKeyword("", 1, 1)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, outInfos[0], test1)
		})
	})
}

func TestSearchMemberNumByKeyword(t *testing.T) {
	Convey("SearchMemberNumByKeyword, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)

		Convey(" SearchGroupNumByKeyword Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupMemberDB.EXPECT().SearchMemberNumByKeyword(gomock.Any()).AnyTimes().Return(0, testErr)
			_, err := groupMember.SearchMemberNumByKeyword("")

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			groupMemberDB.EXPECT().SearchMemberNumByKeyword(gomock.Any()).AnyTimes().Return(2, nil)
			outInfos, err := groupMember.SearchMemberNumByKeyword("")

			assert.Equal(t, err, nil)
			assert.Equal(t, outInfos, 2)
		})
	})
}

func TestGetGroupOnClient(t *testing.T) {
	Convey("GetGroupOnClient, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)

		Convey(" GetGroupsNum Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(0, testErr)
			_, _, err := groupMember.GetGroupOnClient(1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey(" GetGroups Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(2, nil)
			groupDB.EXPECT().GetGroups(gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetGroupOnClient(1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			test1 := interfaces.GroupInfo{
				ID:    "xxxx",
				Name:  "zzzz",
				Notes: "zzzz",
			}
			groupDB.EXPECT().GetGroupsNum(gomock.Any()).AnyTimes().Return(2, nil)
			groupDB.EXPECT().GetGroups(gomock.Any()).AnyTimes().Return([]interfaces.GroupInfo{test1}, nil)
			outInfos, num, err := groupMember.GetGroupOnClient(1, 1)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, num, 2)
			assert.Equal(t, outInfos[0].ID, test1.ID)
			assert.Equal(t, outInfos[0].Name, test1.Name)
		})
	})
}

func TestGetMemberOnClient(t *testing.T) {
	Convey("GetMemberOnClient, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)

		var groupInfo interfaces.GroupInfo
		id := strID1
		Convey(" GetGroupByID Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, testErr)
			_, _, err := groupMember.GetMemberOnClient(id, 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey(" group not exsit", func() {
			testErr := rest.NewHTTPErrorV2(errors.NotFound, "group does not exist")
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			_, _, err := groupMember.GetMemberOnClient(id, 1, 1)

			assert.Equal(t, err, testErr)
		})

		groupInfo.ID = id
		groupInfo.Name = "zzzz"
		Convey(" GetGroupMembersNum Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum(gomock.Any(), gomock.Any()).AnyTimes().Return(0, testErr)
			_, _, err := groupMember.GetMemberOnClient(id, 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey(" GetGroupOnClient Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum(gomock.Any(), gomock.Any()).AnyTimes().Return(2, nil)
			groupMemberDB.EXPECT().GetMemberOnClient(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, _, err := groupMember.GetMemberOnClient(id, 1, 1)

			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			test1 := interfaces.MemberSimpleInfo{
				ID:    id,
				Name:  "zzzz",
				NType: 1,
			}
			groupDB.EXPECT().GetGroupByID(gomock.Any()).AnyTimes().Return(groupInfo, nil)
			groupMemberDB.EXPECT().GetGroupMembersNum(gomock.Any(), gomock.Any()).AnyTimes().Return(2, nil)
			groupMemberDB.EXPECT().GetMemberOnClient(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.MemberSimpleInfo{test1}, nil)
			outInfos, num, err := groupMember.GetMemberOnClient(id, 1, 1)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, num, 2)
			assert.Equal(t, outInfos[0].ID, test1.ID)
			assert.Equal(t, outInfos[0].Name, test1.Name)
			assert.Equal(t, outInfos[0].NType, test1.NType)
		})
	})
}

func TestConverGroupName(t *testing.T) {
	Convey("ConvertGroupName, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userDB := mock.NewMockDBUser(ctrl)
		departmentDB := mock.NewMockDBDepartment(ctrl)
		groupDB := mock.NewMockDBGroup(ctrl)
		groupMemberDB := mock.NewMockDBGroupMember(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		msg := mock.NewMockDrivenMessageBroker(ctrl)
		groupMember := newGroup(groupMemberDB, groupDB, userDB, departmentDB, eacplog, msg, nil)
		visitor := interfaces.Visitor{}

		Convey(" ids empty", func() {
			outInfo, err := groupMember.ConvertGroupName(&visitor, make([]string, 0), true)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfo), 0)
		})

		Convey(" ConvertGroupName Error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			groupDB.EXPECT().GetGroupName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			_, err := groupMember.ConvertGroupName(&visitor, []string{"xxx"}, true)

			assert.Equal(t, err, testErr)
		})

		Convey(" ids not exsit", func() {
			testInfo := make([]interfaces.NameInfo, 0)
			exsitIDs := []string{"xxx"}
			groupDB.EXPECT().GetGroupName(gomock.Any()).AnyTimes().Return(testInfo, exsitIDs, nil)
			_, err := groupMember.ConvertGroupName(&visitor, []string{"xxx", "yyy"}, true)

			tmpIds := []string{"yyy"}
			testErr := rest.NewHTTPErrorV2(errors.GroupNotFound,
				groupMember.i18n.Load(i18nIDObjectsInGroupNotFound, visitor.Language),
				rest.SetDetail(map[string]interface{}{"ids": tmpIds}),
				rest.SetCodeStr(errors.StrBadRequestGroupNotFound),
			)
			assert.Equal(t, err, testErr)
		})

		Convey("Success", func() {
			tmp1 := interfaces.NameInfo{
				ID:   "xxx",
				Name: "zzzzz",
			}
			tmp2 := interfaces.NameInfo{
				ID:   "yyy",
				Name: "zzzzz1",
			}
			testInfo := []interfaces.NameInfo{
				tmp1,
				tmp2,
			}
			exsitIDs := []string{"xxx", "yyy"}
			groupDB.EXPECT().GetGroupName(gomock.Any()).AnyTimes().Return(testInfo, exsitIDs, nil)
			outInfo, err := groupMember.ConvertGroupName(&visitor, []string{"xxx", "yyy", "yyy"}, true)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfo), 2)
		})

		Convey("Success1", func() {
			tmp1 := interfaces.NameInfo{
				ID:   "xxx",
				Name: "zzzzz",
			}
			tmp2 := interfaces.NameInfo{
				ID:   "yyy",
				Name: "zzzzz1",
			}
			testInfo := []interfaces.NameInfo{
				tmp1,
				tmp2,
			}
			exsitIDs := []string{"xxx", "yyy"}
			groupDB.EXPECT().GetGroupName(gomock.Any()).AnyTimes().Return(testInfo, exsitIDs, nil)
			outInfo, err := groupMember.ConvertGroupName(&visitor, []string{"xxx", "yyy", "zzz"}, false)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(outInfo), 2)
		})
	})
}

func TestCheckGroupManageAuth(t *testing.T) {
	Convey("checkGroupManageAuth, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		role := mock.NewMockLogicsRole(ctrl)

		groupMember := group{}
		groupMember.orgPermApp = orgPermApp
		groupMember.role = role

		visitor := interfaces.Visitor{}
		visitor.Type = interfaces.Anonymous
		Convey(" visitor type为Anonymous", func() {
			err := groupMember.checkGroupManageAuth(&visitor)

			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		var tempRoles map[string]map[interfaces.Role]bool
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey(" visitor type为RealName", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(tempRoles, testErr)
			err := groupMember.checkGroupManageAuth(&visitor)

			assert.Equal(t, err, testErr)
		})

		visitor.Type = interfaces.App
		Convey(" visitor type为App", func() {
			orgPermApp.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.checkGroupManageAuth(&visitor)

			assert.Equal(t, err, testErr)
			assert.Equal(t, 1, 1)
		})
	})
}

func TestCheckGroupManageAuth2(t *testing.T) {
	Convey("checkGroupManageAuth2, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		role := mock.NewMockLogicsRole(ctrl)

		groupMember := group{}
		groupMember.orgPermApp = orgPermApp
		groupMember.role = role
		ctx := context.Background()

		visitor := interfaces.Visitor{}
		visitor.Type = interfaces.Anonymous
		Convey(" visitor type为Anonymous", func() {
			err := groupMember.checkGroupManageAuth2(ctx, &visitor)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority"))
		})

		visitor.Type = interfaces.RealName
		var tempRoles map[string]map[interfaces.Role]bool
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey(" visitor type为RealName", func() {
			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(tempRoles, testErr)
			err := groupMember.checkGroupManageAuth2(ctx, &visitor)

			assert.Equal(t, err, testErr)
		})

		visitor.Type = interfaces.App
		Convey(" visitor type为App", func() {
			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.checkGroupManageAuth2(ctx, &visitor)

			assert.Equal(t, err, testErr)
			assert.Equal(t, 1, 1)
		})
	})
}

func TestCheckGroupGetAuth(t *testing.T) {
	Convey("checkGroupGetAuth, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		role := mock.NewMockLogicsRole(ctrl)

		groupMember := group{}
		groupMember.orgPermApp = orgPermApp
		groupMember.role = role

		visitor := interfaces.Visitor{}
		visitor.Type = interfaces.Anonymous
		Convey(" visitor type为Anonymous", func() {
			err := groupMember.checkGroupGetAuth(&visitor)

			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil))
		})

		visitor.Type = interfaces.RealName
		var tempRoles map[string]map[interfaces.Role]bool
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey(" visitor type为RealName", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(tempRoles, testErr)
			err := groupMember.checkGroupGetAuth(&visitor)

			assert.Equal(t, err, testErr)
		})

		visitor.Type = interfaces.App
		Convey(" visitor type为App", func() {
			orgPermApp.EXPECT().GetAppPermByID(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.checkGroupGetAuth(&visitor)

			assert.Equal(t, err, testErr)
		})
	})
}

func TestCheckGroupGetAuth2(t *testing.T) {
	Convey("checkGroupGetAuth2, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orgPermApp := mock.NewMockDBOrgPermApp(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		groupMember := group{}
		groupMember.orgPermApp = orgPermApp
		groupMember.role = role
		groupMember.trace = trace

		visitor := interfaces.Visitor{}
		visitor.Type = interfaces.Anonymous
		var ctx context.Context
		Convey(" visitor type为Anonymous", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			err := groupMember.checkGroupGetAuth2(ctx, &visitor)

			assert.Equal(t, err, rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority"))
		})

		visitor.Type = interfaces.RealName
		var tempRoles map[string]map[interfaces.Role]bool
		testErr := rest.NewHTTPError("error", 503000000, nil)
		Convey(" visitor type为RealName", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			role.EXPECT().GetRolesByUserIDs2(gomock.Any(), gomock.Any()).AnyTimes().Return(tempRoles, testErr)
			err := groupMember.checkGroupGetAuth2(ctx, &visitor)

			assert.Equal(t, err, testErr)
		})

		visitor.Type = interfaces.App
		Convey(" visitor type为App", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			orgPermApp.EXPECT().GetAppPermByID2(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			err := groupMember.checkGroupGetAuth2(ctx, &visitor)

			assert.Equal(t, err, testErr)
		})
	})
}

func TestDeleteGroupMemberByMemberID(t *testing.T) {
	Convey("DeleteGroupMemberByMemberID, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gm := mock.NewMockDBGroupMember(ctrl)

		groupMember := &group{
			groupMemberDB: gm,
		}

		Convey("success", func() {
			gm.EXPECT().DeleteGroupMemberByMemberID(gomock.Any()).AnyTimes().Return(nil)
			err := groupMember.DeleteGroupMemberByMemberID("")

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetExistUsersGroupMembers(t *testing.T) {
	Convey("getExistUsersGroupMembers, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := &group{}

		userIDs := []string{strID, strID1, strID2, strName}

		userpath1 := map[string]bool{strName1: true, strName2: true}
		userpath2 := map[string]bool{strName2: true}
		userpath3 := map[string]bool{strRSA2048PrivateKey: true}
		userpath4 := map[string]bool{}

		tempUserParentDeps := make(map[string]map[string]bool)
		tempUserParentDeps[strID] = userpath1
		tempUserParentDeps[strID1] = userpath2
		tempUserParentDeps[strID2] = userpath3
		tempUserParentDeps[strName] = userpath4

		member1 := interfaces.GroupMemberInfo{ID: strName1, MemberType: 2}
		member2 := interfaces.GroupMemberInfo{ID: strName, MemberType: 1}
		member3 := interfaces.GroupMemberInfo{ID: strName2, MemberType: 2}
		tempMemberInfos := []interfaces.GroupMemberInfo{member1, member2, member3}

		Convey("success", func() {
			num, existUserIDs, existDepartIDs, existGroupMember := g.getExistUsersGroupMembers(userIDs, tempUserParentDeps, tempMemberInfos, 0, 4)

			assert.Equal(t, num, 3)

			assert.Equal(t, len(existUserIDs), 3)
			assert.Equal(t, existUserIDs[0], strID)
			assert.Equal(t, existUserIDs[1], strID1)
			assert.Equal(t, existUserIDs[2], strName)

			assert.Equal(t, len(existDepartIDs), 2)
			assert.Equal(t, existDepartIDs[0], strName1)
			assert.Equal(t, existDepartIDs[1], strName2)

			assert.Equal(t, len(existGroupMember), 3)
			assert.Equal(t, len(existGroupMember[strID]), 2)
			assert.Equal(t, existGroupMember[strID][0], member1)
			assert.Equal(t, existGroupMember[strID][1], member3)
			assert.Equal(t, len(existGroupMember[strID1]), 1)
			assert.Equal(t, existGroupMember[strID1][0], member3)
			assert.Equal(t, len(existGroupMember[strName]), 1)
			assert.Equal(t, existGroupMember[strName][0], member2)
		})

		Convey("success 1", func() {
			num, existUserIDs, existDepartIDs, existGroupMember := g.getExistUsersGroupMembers(userIDs, tempUserParentDeps, tempMemberInfos, 1, 1)

			assert.Equal(t, num, 3)

			assert.Equal(t, len(existUserIDs), 1)
			assert.Equal(t, existUserIDs[0], strID1)

			assert.Equal(t, len(existDepartIDs), 1)
			assert.Equal(t, existDepartIDs[0], strName2)

			assert.Equal(t, len(existGroupMember), 1)
			assert.Equal(t, len(existGroupMember[strID1]), 1)
			assert.Equal(t, existGroupMember[strID1][0], member3)
		})
	})
}

//nolint:dupl
func TestSendAddGroupAuditLog(t *testing.T) {
	Convey("sendAddGroupAuditLog, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		groupMember := newGroup(nil, nil, nil, nil, eacplog, nil, nil)

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor":      make(map[string]interface{}),
			"nalog_infome": make(map[string]interface{}),
		}
		Convey("error", func() {
			eacplog.EXPECT().EacpLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.sendAddGroupAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

//nolint:dupl
func TestSendDeleteGroupAuditLog(t *testing.T) {
	Convey("sendDeleteGroupAuditLog, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		groupMember := newGroup(nil, nil, nil, nil, eacplog, nil, nil)

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor":      make(map[string]interface{}),
			"nalog_infome": make(map[string]interface{}),
		}
		Convey("error", func() {
			eacplog.EXPECT().EacpLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.sendDeleteGroupAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

//nolint:dupl
func TestSendModifyGroupAuditLog(t *testing.T) {
	Convey("sendModifyGroupAuditLog, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		groupMember := newGroup(nil, nil, nil, nil, eacplog, nil, nil)

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor":      make(map[string]interface{}),
			"nalog_infome": make(map[string]interface{}),
		}
		Convey("error", func() {
			eacplog.EXPECT().EacpLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.sendModifyGroupAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

//nolint:dupl
func TestSendDeleteGroupMembersAuditLog(t *testing.T) {
	Convey("sendDeleteGroupMembersAuditLog, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		groupMember := newGroup(nil, nil, nil, nil, eacplog, nil, nil)

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor":      make(map[string]interface{}),
			"nalog_infome": make(map[string]interface{}),
		}
		Convey("error", func() {
			eacplog.EXPECT().EacpLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.sendDeleteGroupMembersAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}

//nolint:dupl
func TestSendAddGroupMembersAuditLog(t *testing.T) {
	Convey("sendAddGroupMembersAuditLog, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		groupMember := newGroup(nil, nil, nil, nil, eacplog, nil, nil)

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor":      make(map[string]interface{}),
			"nalog_infome": make(map[string]interface{}),
		}
		Convey("error", func() {
			eacplog.EXPECT().EacpLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := groupMember.sendAddGroupMembersAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}
