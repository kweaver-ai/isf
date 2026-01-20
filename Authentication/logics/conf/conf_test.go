package conf

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
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

func newConf(db interfaces.DBConf, userMgnt interfaces.DnUserManagement) *conf {
	return &conf{
		cdb:      db,
		userMgnt: userMgnt,
		log:      common.NewLogger(),
	}
}

//nolint:lll
func TestGetConfig(t *testing.T) {
	Convey("get config", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBConf(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		conf := newConf(db, userMgnt)

		testErr := errors.New("some error")

		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		cofigKeys := map[interfaces.ConfigKey]bool{
			interfaces.RememberFor: true,
		}
		ctx := context.Background()

		Convey("GetUserRolesByUserID error", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, testErr)
			_, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, err, testErr)
		})

		Convey("组织管理员、组织审计员、普通用户获取配置", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.OrganizationAdmin, interfaces.OrganizationAudit, interfaces.NormalUser}, nil)
			_, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})

		Convey("admin、aduit管理员获取认证配置", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SystemAdmin, interfaces.AuditAdmin}, nil)
			_, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})

		Convey("应用账户获取配置", func() {
			visitor.Type = interfaces.Business
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			_, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})

		Convey("GetConfig failed", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			db.EXPECT().GetConfig(gomock.Any()).Return(interfaces.Config{}, testErr)
			_, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, err, testErr)
		})

		Convey("admin管理员获取认证配置", func() {
			config := interfaces.Config{
				RememberFor: 600,
			}
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			db.EXPECT().GetConfig(gomock.Any()).Return(interfaces.Config{RememberFor: 600}, nil)
			cfg, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, cfg, config)
			assert.Equal(t, err, nil)
		})

		Convey("security管理员获取认证配置", func() {
			config := interfaces.Config{
				RememberFor: 600,
			}
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SecurityAdmin}, nil)
			db.EXPECT().GetConfig(gomock.Any()).Return(interfaces.Config{RememberFor: 600}, nil)
			cfg, err := conf.GetConfig(ctx, visitor, cofigKeys)
			assert.Equal(t, cfg, config)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetConfigFromShareMgnt(t *testing.T) {
	Convey("get config from sharemgnt", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)
		db := mock.NewMockDBConf(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		conf := newConf(db, userMgnt)
		conf.trace = trace

		ctx := context.Background()
		testErr := errors.New("some error")
		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		cofigKeys := map[interfaces.ConfigKey]bool{
			interfaces.RememberFor: true,
		}

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get config error", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			db.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any()).Return(interfaces.Config{}, testErr)

			_, err := conf.GetConfigFromShareMgnt(ctx, visitor, cofigKeys)

			assert.Equal(t, err, testErr)
		})
		Convey("get config success", func() {
			config := interfaces.Config{
				EnablePWDLock:      true,
				EnableThirdPWDLock: true,
			}
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			db.EXPECT().GetConfigFromShareMgnt(gomock.Any(), gomock.Any()).Return(interfaces.Config{EnablePWDLock: true, EnableThirdPWDLock: true}, nil)

			cfg, err := conf.GetConfigFromShareMgnt(ctx, visitor, cofigKeys)

			assert.Equal(t, err, nil)
			assert.Equal(t, cfg, config)
		})
	})
}

//nolint:lll
func TestSetConfig(t *testing.T) {
	Convey("set config", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)
		db := mock.NewMockDBConf(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		sqldb, sqlmock, _ := sqlx.New()
		ob := mock.NewMockLogicsOutbox(ctrl)
		conf := newConf(db, userMgnt)
		conf.pool = sqldb
		conf.ob = ob
		conf.trace = trace

		testErr := errors.New("some error")
		visitor := &interfaces.Visitor{
			Type: interfaces.RealName,
			ID:   "some-id",
		}
		cofigKeys := map[interfaces.ConfigKey]bool{
			interfaces.RememberFor: true,
		}
		cfg := interfaces.Config{
			RememberFor: 900,
		}
		ctx := context.Background()

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("GetUserRolesByUserID error", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, testErr)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err, testErr)
		})
		Convey("组织管理员、组织审计员、普通用户设置配置", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.OrganizationAdmin, interfaces.OrganizationAudit, interfaces.NormalUser}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})
		Convey("admin、aduit管理员设置认证配置", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SystemAdmin, interfaces.AuditAdmin}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})
		Convey("应用账户设置配置", func() {
			visitor.Type = interfaces.Business
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.OrganizationAdmin, interfaces.OrganizationAudit, interfaces.NormalUser}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Unauthorized)
		})
		Convey("remember_for小于0", func() {
			cfg.RememberFor = -100
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
		})
		Convey("SetConfig error", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			sqlmock.ExpectBegin()
			sqlmock.ExpectRollback()
			db.EXPECT().SetConfig(gomock.Any(), gomock.Any()).Return(testErr)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err, testErr)
		})
		Convey("admin管理员设置认证配置", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SuperAdmin}, nil)
			sqlmock.ExpectBegin()
			sqlmock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().Return()
			db.EXPECT().SetConfig(gomock.Any(), gomock.Any()).Return(nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err, nil)
		})
		Convey("security管理员设置认证配置/remember_for等于0", func() {
			cfg.RememberFor = 0
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SecurityAdmin}, nil)
			sqlmock.ExpectBegin()
			sqlmock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread().Return()
			db.EXPECT().SetConfig(gomock.Any(), gomock.Any()).Return(nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err, nil)
		})
		Convey("smsExpiration小于1", func() {
			cfg.SMSExpiration = 0
			cofigKeys[interfaces.SMSExpiration] = true
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SecurityAdmin}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
		})
		Convey("smsExpiration大于60", func() {
			cfg.SMSExpiration = 61
			cofigKeys[interfaces.SMSExpiration] = true
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{interfaces.SecurityAdmin}, nil)
			err := conf.SetConfig(ctx, visitor, cofigKeys, cfg)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
		})
	})
}
