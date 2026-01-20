// Package accesstokenperm 逻辑层
package accesstokenperm

import (
	"context"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

const appID = "b550af01-06d0-446d-be5b-b44cfcd97906"

func newAccessTokenPerm(db interfaces.DBAccessTokenPerm, userMgnt interfaces.DnUserManagement, hydraAdmin interfaces.DnHydraAdmin, eacpLog interfaces.DnEacpLog) *accessTokenPerm {
	return &accessTokenPerm{
		db:         db,
		userMgnt:   userMgnt,
		hydraAdmin: hydraAdmin,
		eacpLog:    eacpLog,
	}
}

func TestSetAppAccessTokenPerm(t *testing.T) {
	Convey("TestSetAppAccessTokenPerm", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		db1, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := db1.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)
		ac.ob = ob
		ac.pool = db1

		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		roleType := []interfaces.RoleType{interfaces.SuperAdmin}
		testErr := errors.New("some error")
		ctx := context.Background()

		Convey("GetUserRolesByUserID 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, testErr)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("访问者类型不是实名用户", func() {
			visitor.Type = interfaces.Business
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user type")
		})

		Convey("访问者角色不是超级/系统管理员", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.AuditAdmin}, nil)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user role type")
		})

		Convey("GetAppInfo 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, testErr)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("CheckAppAccessTokenPerm 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, testErr)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("已添加过权限", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, nil)
		})

		Convey("未添加过权限，更新客户端client_type失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)
			hydraAdmin.EXPECT().SetAppAsUserAgent(gomock.Any()).Return(testErr)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("未添加过权限，添加失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)
			hydraAdmin.EXPECT().SetAppAsUserAgent(gomock.Any()).Return(nil)
			db.EXPECT().AddAppAccessTokenPerm(gomock.Any()).Return(testErr)
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("未添加过权限，添加成功", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)
			hydraAdmin.EXPECT().SetAppAsUserAgent(gomock.Any()).Return(nil)
			db.EXPECT().AddAppAccessTokenPerm(gomock.Any()).Return(nil)

			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()
			err := ac.SetAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAppAccessTokenPerm(t *testing.T) {
	Convey("TestDeleteAppAccessTokenPerm", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		db1, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := db1.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)
		ac.ob = ob
		ac.pool = db1

		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		roleType := []interfaces.RoleType{interfaces.SuperAdmin}
		testErr := errors.New("some error")
		ctx := context.Background()

		Convey("GetUserRolesByUserID 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, testErr)
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("访问者类型不是实名用户", func() {
			visitor.Type = interfaces.Business
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user type")
		})

		Convey("访问者角色不是超级/系统管理员", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.AuditAdmin}, nil)
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user role type")
		})

		Convey("GetAppInfo 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, testErr)
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteAppAccessTokenPerm 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().DeleteAppAccessTokenPerm(gomock.Any()).Return(testErr)
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteAppAccessTokenPerm 成功", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			userMgnt.EXPECT().GetAppInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.AppInfo{}, nil)
			db.EXPECT().DeleteAppAccessTokenPerm(gomock.Any()).Return(nil)

			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()
			err := ac.DeleteAppAccessTokenPerm(ctx, visitor, appID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAppDeleted(t *testing.T) {
	Convey("TestAppDeleted", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)

		Convey("AppDeleted", func() {
			db.EXPECT().DeleteAppAccessTokenPerm(gomock.Any()).Return(nil)
			err := ac.AppDeleted(appID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckAppAccessTokenPerm(t *testing.T) {
	Convey("TestCheckAppAccessTokenPerm", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)

		Convey("CheckAppAccessTokenPerm", func() {
			db.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)
			res, err := ac.CheckAppAccessTokenPerm(appID)
			assert.Equal(t, err, nil)
			assert.Equal(t, res, false)
		})
	})
}

func TestGetAllAppAccessTokenPerm(t *testing.T) {
	Convey("TestGetAllAppAccessTokenPerm", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)

		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		roleType := []interfaces.RoleType{interfaces.SuperAdmin}
		testErr := errors.New("some error")
		ctx := context.Background()
		Convey("GetUserRolesByUserID 失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, testErr)
			_, err := ac.GetAllAppAccessTokenPerm(ctx, visitor)
			assert.Equal(t, err, testErr)
		})

		Convey("访问者类型不是实名用户", func() {
			visitor.Type = interfaces.Business
			_, err := ac.GetAllAppAccessTokenPerm(ctx, visitor)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user type")
		})

		Convey("访问者角色不是超级/系统管理员", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.AuditAdmin}, nil)
			_, err := ac.GetAllAppAccessTokenPerm(ctx, visitor)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Unsupported user role type")
		})

		Convey("GetAllAppAccessTokenPerm", func() {
			permApps := []string{"d8521454-c8ff-402f-9ccb-e7f2c0a0723c"}
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleType, nil)
			db.EXPECT().GetAllAppAccessTokenPerm().Return(permApps, nil)
			res, err := ac.GetAllAppAccessTokenPerm(ctx, visitor)
			assert.Equal(t, len(res), 1)
			assert.Equal(t, res[0], "d8521454-c8ff-402f-9ccb-e7f2c0a0723c")
			assert.Equal(t, err, nil)
		})
	})
}

//nolint:dupl
func TestSetAppPermAuditLog(t *testing.T) {
	Convey("TestSetAppPermAuditLog", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)

		var messageJSON interface{}
		data := `{"type":3,"content":{"visitor":{"ID":"266c6a42-6131-4d62-8f39-853e7093701c","TokenID":"token_id1",
		"IP":"10.4.36.153","Mac":"mac1","UserAgent":"user_agent1","Type":1,"Language":1},"app_name":"xx"}}`
		err := jsoniter.UnmarshalFromString(data, &messageJSON)
		assert.Equal(t, err, nil)
		content := messageJSON.(map[string]interface{})["content"]

		test1 := interfaces.Visitor{
			ID:        "266c6a42-6131-4d62-8f39-853e7093701c",
			TokenID:   "token_id1",
			IP:        "10.4.36.153",
			Mac:       "mac1",
			UserAgent: "user_agent1",
			Type:      interfaces.RealName,
			Language:  interfaces.SimplifiedChinese,
		}
		Convey("success", func() {
			eacpLog.EXPECT().OpSetAppAccessTokenPerm(&test1, gomock.Any()).Return(nil)
			err := ac.setAppPermAuditLog(content)
			assert.Equal(t, err, nil)
		})
	})
}

//nolint:dupl
func TestDeleteAppPermAuditLog(t *testing.T) {
	Convey("TestDeleteAppPermAuditLog", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAccessTokenPerm(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		eacpLog := mock.NewMockDnEacpLog(ctrl)
		ac := newAccessTokenPerm(db, userMgnt, hydraAdmin, eacpLog)

		var messageJSON interface{}
		data := `{"type":3,"content":{"visitor":{"ID":"266c6a42-6131-4d62-8f39-853e7093701c","TokenID":"token_id1",
		"IP":"10.4.36.153","Mac":"mac1","UserAgent":"user_agent1","Type":1,"Language":1},"app_name":"xx"}}`
		err := jsoniter.UnmarshalFromString(data, &messageJSON)
		assert.Equal(t, err, nil)
		content := messageJSON.(map[string]interface{})["content"]

		test1 := interfaces.Visitor{
			ID:        "266c6a42-6131-4d62-8f39-853e7093701c",
			TokenID:   "token_id1",
			IP:        "10.4.36.153",
			Mac:       "mac1",
			UserAgent: "user_agent1",
			Type:      interfaces.RealName,
			Language:  interfaces.SimplifiedChinese,
		}
		Convey("success", func() {
			eacpLog.EXPECT().OpDeleteAppAccessTokenPerm(&test1, gomock.Any()).Return(nil)
			err := ac.deleteAppPermAuditLog(content)
			assert.Equal(t, err, nil)
		})
	})
}
