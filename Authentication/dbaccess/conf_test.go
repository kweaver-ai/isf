package dbaccess

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces"
	mocks "Authentication/interfaces/mock"
)

func newConf(ptrDB *sqlx.DB) *conf {
	c := &conf{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
	c.tableInfoMapInit()
	return c
}

func TestGetConfig(t *testing.T) {
	Convey("Get, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		conf := newConf(db)
		conf.trace = trace
		conf.dbTrace = db

		configKeys := map[interfaces.ConfigKey]bool{interfaces.RememberFor: true, interfaces.RememberVisible: true}
		Convey("configTypes is empty", func() {
			configTypesTmp := make(map[interfaces.ConfigKey]bool)
			cfg, err := conf.GetConfig(configTypesTmp)
			assert.Equal(t, cfg, interfaces.Config{})
			assert.Equal(t, err, nil)
		})

		testErr := errors.New("some error")
		Convey("get config err", func() {
			mock.ExpectQuery("").WillReturnError(testErr)
			configMap, err := conf.GetConfig(configKeys)
			assert.Equal(t, configMap, interfaces.Config{})
			assert.Equal(t, err, testErr)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		Convey("get authentication.t_conf config success", func() {
			fields := []string{
				"f_key",
				"f_value",
			}
			cfg := interfaces.Config{
				RememberFor:     600,
				RememberVisible: false,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("remember_for", "600").AddRow("remember_visible", "false"))
			configMap, err := conf.GetConfig(configKeys)
			assert.Equal(t, configMap, cfg)
			assert.Equal(t, err, nil)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		configKeys = map[interfaces.ConfigKey]bool{
			interfaces.EnableIDCardLogin:  true,
			interfaces.EnablePWDLock:      true,
			interfaces.EnableThirdPWDLock: true,
			interfaces.PWDErrCnt:          true,
			interfaces.PWDLockTime:        true,
			interfaces.VCodeConfig:        true,
			interfaces.LimitRedirectURI:   true,
		}
		Convey("get sharemgnt_db.t_sharemgnt_config config success", func() {
			fields := []string{
				"f_key",
				"f_value",
			}
			vcodeConfig := interfaces.VCodeLoginConfig{
				Enable:    false,
				PWDErrCnt: 0,
			}
			cfg := interfaces.Config{
				EnableIDCardLogin:  false,
				EnablePWDLock:      false,
				EnableThirdPWDLock: false,
				PWDErrCnt:          0,
				PWDLockTime:        0,
				VCodeConfig:        vcodeConfig,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).
				AddRow("id_card_login_status", "false").
				AddRow("enable_pwd_lock", "false").
				AddRow("enable_third_pwd_lock", "false").
				AddRow("pwd_err_cnt", 0).
				AddRow("pwd_lock_time", 0).
				AddRow("vcode_login_config", `{"isEnable":false,"passwdErrCnt":0}`).
				AddRow("oauth2_redirect_uris", `["https://www.baidu.com", "https://www.google.com"]`))
			configMap, err := conf.GetConfig(configKeys)
			assert.Equal(t, err, nil)
			assert.Equal(t, configMap.EnableIDCardLogin, cfg.EnableIDCardLogin)
			assert.Equal(t, configMap.EnablePWDLock, cfg.EnablePWDLock)
			assert.Equal(t, configMap.EnableThirdPWDLock, cfg.EnableThirdPWDLock)
			assert.Equal(t, configMap.PWDErrCnt, cfg.PWDErrCnt)
			assert.Equal(t, configMap.PWDLockTime, cfg.PWDLockTime)
			assert.Equal(t, configMap.VCodeConfig, cfg.VCodeConfig)
			assert.Equal(t, configMap.LimitRedirectURI["https://www.baidu.com"], true)
			assert.Equal(t, configMap.LimitRedirectURI["https://www.google.com"], true)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		Convey("get sharemgnt_db.t_sharemgnt_config config success1", func() {
			fields := []string{
				"f_key",
				"f_value",
			}

			tempConfigKeys := map[interfaces.ConfigKey]bool{
				interfaces.LimitRedirectURI: true,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			configMap, err := conf.GetConfig(tempConfigKeys)
			assert.Equal(t, err, nil)
			assert.Equal(t, configMap, interfaces.Config{
				LimitRedirectURI: make(map[string]bool),
			})

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		Convey("get sharemgnt_db.t_sharemgnt_config config success2", func() {
			fields := []string{
				"f_key",
				"f_value",
			}
			tempConfigKeys := map[interfaces.ConfigKey]bool{
				interfaces.TriSystemStatus: true,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("enable_tri_system_status", "1"))
			configMap, err := conf.GetConfig(tempConfigKeys)
			assert.Equal(t, err, nil)
			assert.Equal(t, configMap.TriSystemStatus, true)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})
}

func TestGetConfigFromShareMgnt(t *testing.T) {
	Convey("Get, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		conf := newConf(db)
		conf.trace = trace
		conf.dbTrace = db

		ctx := context.Background()
		testErr := errors.New("some error")
		configKeys := map[interfaces.ConfigKey]bool{
			interfaces.EnablePWDLock:      true,
			interfaces.EnableThirdPWDLock: true,
			interfaces.PWDErrCnt:          true,
			interfaces.PWDLockTime:        true,
		}

		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("get config from sharemgnt err", func() {
			mock.ExpectQuery("").WillReturnError(testErr)
			configMap, err := conf.GetConfigFromShareMgnt(ctx, configKeys)
			assert.Equal(t, configMap, interfaces.Config{})
			assert.Equal(t, err, testErr)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		Convey("get config from sharemgnt success", func() {
			fields := []string{
				"f_key",
				"f_value",
			}
			cfg := interfaces.Config{
				EnablePWDLock:      false,
				EnableThirdPWDLock: false,
				PWDErrCnt:          0,
				PWDLockTime:        0,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).
				AddRow("enable_pwd_lock", "false").
				AddRow("enable_third_pwd_lock", "false").
				AddRow("pwd_err_cnt", 0).
				AddRow("pwd_lock_time", 0))
			configMap, err := conf.GetConfigFromShareMgnt(ctx, configKeys)
			assert.Equal(t, configMap, cfg)
			assert.Equal(t, err, nil)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})
}

func TestSetConfig(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		mock.ExpectBegin()

		conf := newConf(db)
		configKeys := map[interfaces.ConfigKey]bool{interfaces.RememberFor: true}
		cfg := interfaces.Config{
			RememberFor: 600,
		}

		testErr := errors.New("some error")
		Convey("set config err", func() {
			mock.ExpectExec("").WithArgs("600", "remember_for").WillReturnError(testErr)
			mock.ExpectRollback()
			err := conf.SetConfig(configKeys, cfg)
			assert.Equal(t, err, testErr)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		mock.MatchExpectationsInOrder(false)
		Convey("success", func() {
			cfg := interfaces.Config{
				RememberFor:     600,
				RememberVisible: false,
			}
			configKeys[interfaces.RememberVisible] = true
			mock.ExpectExec("").WithArgs("600", "remember_for").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("").WithArgs("false", "remember_visible").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			err := conf.SetConfig(configKeys, cfg)
			assert.Equal(t, err, nil)

			// 判断是否所有期望都被达到
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})
}
