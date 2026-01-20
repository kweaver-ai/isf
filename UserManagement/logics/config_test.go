package logics

import (
	"testing"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestNewConfig(t *testing.T) {
	sqlDB, _, err := sqlx.New()
	assert.Equal(t, err, nil)
	dbPool = sqlDB

	data := NewConfig()
	assert.NotEqual(t, data, nil)
}

//nolint:dupl,funlen
func TestUpdateConfig(t *testing.T) {
	Convey("更新配置信息", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		configDB := mock.NewMockDBConfig(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		con := &config{
			db:      configDB,
			pool:    dPool,
			logger:  common.NewLogger(),
			role:    role,
			eacpLog: eacplog,
			ob:      ob,
		}

		rg := make(map[interfaces.ConfigKey]bool)
		var config interfaces.Config
		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		rg[interfaces.UserDefaultPWD] = true
		config.UserDefaultPWD = "aaaaa"
		visitor := interfaces.Visitor{}
		visitor.ID = strID

		Convey("visitor没有权限，则报错", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := con.UpdateConfig(&visitor, rg, &config)
			assert.Equal(t, err, testErr)
		})

		Convey("密码未加密，报错", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			err := con.UpdateConfig(&visitor, rg, &config)
			assert.NotEqual(t, err, nil)
		})

		config.UserDefaultPWD = "zhMytQF/dmSfreO1Qgkdr8wBEtzi/2QcFwroQV8y+AnFjqhS6aAkVExtgk1VpjwtBk6DlmtSedTFngRLbc61aDQKKJhXTtGYucnRGgOOqD2uu" +
			"I+MaxxAk5t7Vys29XzEyHXB5OAvETjfxkNV/5jAxmQ8k29NDraxpz/yhZ/SsnviskBaGE+l/n+7EvhL2VIVhf9Yp3FB96tOxrjfApf+7a0iIN5NgM+5YjazKnN8nHAJ5Em" +
			"SINBbb+nK+7ciC+IkEBLXRms5Hv5KWpUdPP23iw55Nl3ffjXvqtUyCfVBOqItWDAd32DA7U8Qg8Ver7Tn3wScuGokNijmis0dlEbRVw=="
		Convey("密码加密，密码长度小于6报错", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.UpdateConfig(&visitor, rg, &config)
			assert.Equal(t, err, rest.NewHTTPError("default password formt error", rest.BadRequest, nil))
		})

		config.UserDefaultPWD = ""
		Convey("默认密码长度问题,小于6", func() {
			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, rest.NewHTTPError("default password formt error", rest.BadRequest, nil))
		})

		config.UserDefaultPWD = "xxxxxxxxxxaaaaaaaaaaxxxxxxxxxxaaaaaaaaaaxxxxxxxxxxaaaaaaaaaaxxxxxxxxxxaaaaaaaaaaxxxxxxxxxxaaaaaaaaaa1"
		Convey("默认密码长度问题,大于100", func() {
			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, rest.NewHTTPError("default password formt error", rest.BadRequest, nil))
		})

		config.UserDefaultPWD = "xxxxxx"
		Convey("pool begin 失败", func() {
			txMock.ExpectBegin().WillReturnError(testErr)

			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, testErr)
		})

		Convey("更新密码配置，", func() {
			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, testErr)
		})

		Convey("更新密码配置，UpdateConfig失败", func() {
			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, testErr)
		})

		Convey("成功", func() {
			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()

			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, nil)
		})

		config.UserDefaultPWD = "xxxxxxxx"
		Convey("成功1", func() {
			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()

			err := con.UpdateConfig(nil, rg, &config)
			assert.Equal(t, err, nil)
		})

		config.UserDefaultPWD = "OnPfN16wF+kw6dck2HyR/CLzftN80LLext5ao/PPkQy3qLlOFc6WAzuDOYT+wrp/IMNfJ8QVR8G2yfnT/MpBKgpcYrLSzez2Sybf7Gul820zv2h4w4zbPt7" +
			"zwdEzAZofw7ZBcys1TwO3TzZV/5e00Wkp3w0rPfL9Mq34Lozz7yGjIQtUtdBtMsQPMkEGKhnw907SEfvFxU6HSKV3GAHPS3iEulMPs0A/XG7S4/2i69W30GU17dV9OaRvAxYmWttv1EuW/O1" +
			"AjboahH1UJP/XQyCFnwZH0Aepoz4kXm7UpxjiGOZZgcWuuJy+k3Rri7sl1w4GKdF0nj7PjxnoG7y3nQ=="
		Convey("发送消息失败", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := con.UpdateConfig(&visitor, rg, &config)
			assert.Equal(t, err, testErr)
		})

		Convey("成功2 发送消息", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			txMock.ExpectBegin()
			configDB.EXPECT().SetConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()

			err := con.UpdateConfig(&visitor, rg, &config)
			assert.Equal(t, err, nil)
		})

		Convey("设置密级枚举，无权限", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}
			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, gerrors.NewError(gerrors.PublicForbidden, "this user do not has the authority"))
		})

		Convey("设置密级枚举，获取当前密级报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}
			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{}, testErr)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, testErr)
		})

		Convey("设置密级枚举，已有密级1，报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}
			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			currentConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(currentConfig, nil)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, gerrors.NewError(gerrors.PublicConflict, "csf level enum already set"))
		})

		Convey("设置密级枚举，已有密级2，报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}
			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			currentConfig := interfaces.Config{
				CSFLevel2Enum: map[string]int{
					"test": 1,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSysAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(currentConfig, nil)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, gerrors.NewError(gerrors.PublicConflict, "csf level2 enum already set"))
		})

		Convey("设置密级枚举，密级1为空，报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}
			testConfig := interfaces.Config{
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleAuditAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, gerrors.NewError(gerrors.PublicBadRequest, "csf level enum is empty"))
		})

		Convey("设置密级枚举，密级2为空，报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}

			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSecAdmin] = true
			roleInfos[visitor.ID] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, gerrors.NewError(gerrors.PublicBadRequest, "csf level2 enum is empty"))
		})

		Convey("设置密级枚举，记录日志1 报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}

			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSecAdmin] = true
			roleInfos[visitor.ID] = userRoles

			currentConfig := interfaces.Config{}

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(currentConfig, nil)
			txMock.ExpectBegin()
			configDB.EXPECT().SetShareMgntConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(testErr)
			txMock.ExpectCommit()

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, testErr)
		})

		Convey("设置密级枚举，记录日志2 报错", func() {
			testRg := map[interfaces.ConfigKey]bool{
				interfaces.CSFLevelEnum:  true,
				interfaces.CSFLevel2Enum: true,
			}

			testConfig := interfaces.Config{
				CSFLevelEnum: map[string]int{
					"test": 1,
				},
				CSFLevel2Enum: map[string]int{
					"test2": 2,
				},
			}

			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSecAdmin] = true
			roleInfos[visitor.ID] = userRoles

			currentConfig := interfaces.Config{}

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(currentConfig, nil)
			txMock.ExpectBegin()
			configDB.EXPECT().SetShareMgntConfig(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(testErr)
			txMock.ExpectCommit()

			err := con.UpdateConfig(&visitor, testRg, &testConfig)
			assert.Equal(t, err, testErr)
		})
	})
}

func TestGetConfig(t *testing.T) {
	Convey("获取配置信息", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configDB := mock.NewMockDBConfig(ctrl)
		con := &config{
			db:     configDB,
			logger: common.NewLogger(),
		}

		rg := make(map[interfaces.ConfigKey]bool)
		var config interfaces.Config
		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		rg[interfaces.UserDefaultPWD] = true
		config.UserDefaultPWD = "aaaaa"

		Convey("失败", func() {
			configDB.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(config, testErr)
			_, err := con.GetConfig(rg)
			assert.Equal(t, err, testErr)
		})
	})
}

func TestCheckDefaultPWDAuthority(t *testing.T) {
	Convey("检查修改初始密码权限", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configDB := mock.NewMockDBConfig(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		con := &config{
			db:     configDB,
			logger: common.NewLogger(),
			role:   role,
		}

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)

		Convey("GetRolesByUserIDs 报错", func() {
			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(nil, testErr)

			err := con.checkDefaultPWDAuthority(strID1)
			assert.Equal(t, err, testErr)
		})

		Convey("无权限", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.checkDefaultPWDAuthority(strID1)
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil))
		})

		Convey("成功", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			err := con.checkDefaultPWDAuthority(strID1)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCheckIsValidDefaultPWD(t *testing.T) {
	Convey("检查修改初始密码格式", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configDB := mock.NewMockDBConfig(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		con := &config{
			db:     configDB,
			logger: common.NewLogger(),
			role:   role,
		}

		Convey("长度等于5", func() {
			pwd := "11111"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, false)
		})

		Convey("长度等于6", func() {
			pwd := "111111"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, true)
		})

		Convey("长度等于100", func() {
			pwd := "1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, true)
		})

		Convey("长度等于101", func() {
			pwd := "1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111" +
				"1111111111a"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, false)
		})

		Convey("包含其他特殊字符", func() {
			pwd := "1111111111\t"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, false)
		})

		Convey("包含允许的特殊字符", func() {
			pwd := " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
			out := con.checkIsValidDefaultPWD(pwd)
			assert.Equal(t, out, true)
		})
	})
}

func TestCheckDefaultPWD(t *testing.T) {
	Convey("检查初始密码格式", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configDB := mock.NewMockDBConfig(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		con := &config{
			db:     configDB,
			logger: common.NewLogger(),
			role:   role,
		}

		visitor := interfaces.Visitor{
			ID: strID1,
		}

		Convey("无权限，报错403", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleNormalUser] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			_, _, err := con.CheckDefaultPWD(&visitor, "xxxx")
			assert.Equal(t, err, rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil))
		})

		Convey("密码解压报错", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			visitor.LangType = interfaces.LTENUS
			_, _, err := con.CheckDefaultPWD(&visitor, "xxxx")
			assert.NotEqual(t, err, nil)
		})

		Convey("成功", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			visitor.LangType = interfaces.LTENUS
			pwd := "OnPfN16wF+kw6dck2HyR/CLzftN80LLext5ao/PPkQy3qLlOFc6WAzuDOYT+wrp/IMNfJ8QVR8G2yfnT/MpBKgpcYrLSzez2Sybf7Gul820zv2h4w4zbPt7" +
				"zwdEzAZofw7ZBcys1TwO3TzZV/5e00Wkp3w0rPfL9Mq34Lozz7yGjIQtUtdBtMsQPMkEGKhnw907SEfvFxU6HSKV3GAHPS3iEulMPs0A/XG7S4/2i69W30GU17dV9OaRvAxYmWttv1EuW/O1" +
				"AjboahH1UJP/XQyCFnwZH0Aepoz4kXm7UpxjiGOZZgcWuuJy+k3Rri7sl1w4GKdF0nj7PjxnoG7y3nQ=="
			result, msg, err := con.CheckDefaultPWD(&visitor, pwd)
			assert.Equal(t, err, nil)
			assert.Equal(t, result, true)
			assert.Equal(t, msg, "")
		})

		Convey("失败1", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			visitor.LangType = interfaces.LTZHCN
			pwd := "zhMytQF/dmSfreO1Qgkdr8wBEtzi/2QcFwroQV8y+AnFjqhS6aAkVExtgk1VpjwtBk6DlmtSedTFngRLbc61aDQKKJhXTtGYucnRGgOOqD2uu" +
				"I+MaxxAk5t7Vys29XzEyHXB5OAvETjfxkNV/5jAxmQ8k29NDraxpz/yhZ/SsnviskBaGE+l/n+7EvhL2VIVhf9Yp3FB96tOxrjfApf+7a0iIN5NgM+5YjazKnN8nHAJ5Em" +
				"SINBbb+nK+7ciC+IkEBLXRms5Hv5KWpUdPP23iw55Nl3ffjXvqtUyCfVBOqItWDAd32DA7U8Qg8Ver7Tn3wScuGokNijmis0dlEbRVw=="
			result, msg, err := con.CheckDefaultPWD(&visitor, pwd)
			assert.Equal(t, err, nil)
			assert.Equal(t, result, false)
			assert.Equal(t, msg, "密码只能包含英文或数字或~!%#$@-_.字符，长度范围6~100个字符，请重新输入。")
		})

		Convey("失败2", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			visitor.LangType = interfaces.LTENUS
			pwd := "zhMytQF/dmSfreO1Qgkdr8wBEtzi/2QcFwroQV8y+AnFjqhS6aAkVExtgk1VpjwtBk6DlmtSedTFngRLbc61aDQKKJhXTtGYucnRGgOOqD2uu" +
				"I+MaxxAk5t7Vys29XzEyHXB5OAvETjfxkNV/5jAxmQ8k29NDraxpz/yhZ/SsnviskBaGE+l/n+7EvhL2VIVhf9Yp3FB96tOxrjfApf+7a0iIN5NgM+5YjazKnN8nHAJ5Em" +
				"SINBbb+nK+7ciC+IkEBLXRms5Hv5KWpUdPP23iw55Nl3ffjXvqtUyCfVBOqItWDAd32DA7U8Qg8Ver7Tn3wScuGokNijmis0dlEbRVw=="
			result, msg, err := con.CheckDefaultPWD(&visitor, pwd)
			assert.Equal(t, err, nil)
			assert.Equal(t, result, false)
			assert.Equal(t, msg, "The password should be letters, numbers or ~!%#$@-_. within 6 ~ 100 characters, please re-enter.")
		})

		Convey("失败3", func() {
			roleInfos := make(map[string]map[interfaces.Role]bool)
			userRoles := make(map[interfaces.Role]bool)
			userRoles[interfaces.SystemRoleSuperAdmin] = true
			roleInfos[strID1] = userRoles

			role.EXPECT().GetRolesByUserIDs(gomock.Any()).AnyTimes().Return(roleInfos, nil)

			visitor.LangType = interfaces.LTZHTW
			pwd := "zhMytQF/dmSfreO1Qgkdr8wBEtzi/2QcFwroQV8y+AnFjqhS6aAkVExtgk1VpjwtBk6DlmtSedTFngRLbc61aDQKKJhXTtGYucnRGgOOqD2uu" +
				"I+MaxxAk5t7Vys29XzEyHXB5OAvETjfxkNV/5jAxmQ8k29NDraxpz/yhZ/SsnviskBaGE+l/n+7EvhL2VIVhf9Yp3FB96tOxrjfApf+7a0iIN5NgM+5YjazKnN8nHAJ5Em" +
				"SINBbb+nK+7ciC+IkEBLXRms5Hv5KWpUdPP23iw55Nl3ffjXvqtUyCfVBOqItWDAd32DA7U8Qg8Ver7Tn3wScuGokNijmis0dlEbRVw=="
			result, msg, err := con.CheckDefaultPWD(&visitor, pwd)
			assert.Equal(t, err, nil)
			assert.Equal(t, result, false)
			assert.Equal(t, msg, "密碼只能包含英文或數字或~!%#$@-_.字元，長度範圍6~100個字元，請重新輸入。")
		})
	})
}

func TestSendDefaultPWDModifiedAuditLog(t *testing.T) {
	Convey("发送默认密码被修改事件", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configDB := mock.NewMockDBConfig(ctrl)
		role := mock.NewMockLogicsRole(ctrl)
		eacplog := mock.NewMockDrivenEacpLog(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		con := &config{
			db:      configDB,
			logger:  common.NewLogger(),
			role:    role,
			eacpLog: eacplog,
			ob:      ob,
		}

		testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		data := map[string]interface{}{
			"visitor": make(map[string]interface{}),
		}
		Convey("visitor没有权限，则报错", func() {
			eacplog.EXPECT().OpSetDefaultPWDLog(gomock.Any()).AnyTimes().Return(testErr)
			err := con.sendDefaultPWDModifiedAuditLog(data)
			assert.Equal(t, err, testErr)
		})
	})
}
