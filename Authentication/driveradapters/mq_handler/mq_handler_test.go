package mq

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/common"
	sessionSchema "Authentication/driveradapters/jsonschema/sessions_schema"
	userManagementSchema "Authentication/driveradapters/jsonschema/user_management_schema"
	"Authentication/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestSubscribe(t *testing.T) {
	Convey("订阅", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		h := mock.NewMockMsgBrokerClient(ctrl)
		mq := &msgQueue{
			log:    common.NewLogger(),
			client: h,
		}
		handler := mqHandler{
			mqClient:        mq,
			accessTokenPerm: accesstokenperm,
			log:             common.NewLogger(),
		}

		Convey("Subscribe ok", func() {
			h.EXPECT().Sub(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			handler.Subscribe()
		})
	})
}

func TestAppDeleted(t *testing.T) {
	Convey("TestAppDeleted", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		access := mock.NewMockAccessTokenPerm(ctrl)
		handler := mqHandler{
			accessTokenPerm: access,
			log:             common.NewLogger(),
		}

		msg := `{"id":"app-id"}`
		testErr := errors.New("test")

		Convey("msg 格式错误", func() {
			errMsg, _ := jsoniter.Marshal("[")
			err := handler.appDeleted(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 不存在", func() {
			errMsg, _ := jsoniter.Marshal("{}")
			err := handler.appDeleted(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 为空", func() {
			errMsg := `{"id":""}`
			err := handler.appDeleted([]byte(errMsg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteAppAccessTokenPerm 失败", func() {
			access.EXPECT().AppDeleted(gomock.Any()).Return(testErr)
			err := handler.appDeleted([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteAppAccessTokenPerm 成功", func() {
			access.EXPECT().AppDeleted(gomock.Any()).Return(nil)
			err := handler.appDeleted([]byte(msg))
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteHydraSession(t *testing.T) {
	Convey("TestDeleteHydraSession", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(sessionSchema.HydraSessionsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		session := mock.NewMockHydraSession(ctrl)
		handler := mqHandler{
			hydraSession:   session,
			log:            common.NewLogger(),
			sessionsSchema: sessionsSchema,
		}

		msg := `{"user_id":"userID"}`
		msg2 := `{"user_id":"userID", "client_id":"client_id"}`

		Convey("msg 格式错误", func() {
			errMsg, _ := jsoniter.Marshal("[")
			err := handler.sessionDelete(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg user_id 不存在", func() {
			errMsg, _ := jsoniter.Marshal("{}")
			err := handler.sessionDelete(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg user_id 为空", func() {
			errMsg := `{"user_id":""}`
			err := handler.sessionDelete([]byte(errMsg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 失败1", func() {
			testErr := errors.New("test")
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.sessionDelete([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败4", func() {
			testErr := &rest.HTTPError{
				Code: rest.InternalServerError,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.sessionDelete([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败5", func() {
			testErr := &rest.HTTPError{
				Code: rest.URINotExist,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.sessionDelete([]byte(msg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 成功1", func() {
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			err := handler.sessionDelete([]byte(msg))
			assert.Equal(t, err, nil)
		})
		Convey("DeleteHydraSession 成功2", func() {
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			err := handler.sessionDelete([]byte(msg2))
			assert.Equal(t, err, nil)
		})
	})
}

//nolint:dupl
func TestUserDelete(t *testing.T) {
	Convey("TestUserDelete", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(sessionSchema.HydraSessionsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userDeleteSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserDeleteSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		session := mock.NewMockHydraSession(ctrl)
		handler := mqHandler{
			hydraSession:     session,
			log:              common.NewLogger(),
			sessionsSchema:   sessionsSchema,
			userDeleteSchema: userDeleteSchema,
		}

		msg := `{"id":"userID"}`

		Convey("msg 格式错误", func() {
			errMsg, _ := jsoniter.Marshal("[")
			err := handler.userDelete(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 不存在", func() {
			errMsg, _ := jsoniter.Marshal("{}")
			err := handler.userDelete(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 为空", func() {
			errMsg := `{"id":""}`
			err := handler.userDelete([]byte(errMsg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 失败1", func() {
			testErr := errors.New("test")
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userDelete([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败4", func() {
			testErr := &rest.HTTPError{
				Code: rest.InternalServerError,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userDelete([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败5", func() {
			testErr := &rest.HTTPError{
				Code: rest.URINotExist,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userDelete([]byte(msg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 成功", func() {
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			err := handler.userDelete([]byte(msg))
			assert.Equal(t, err, nil)
		})
	})
}

//nolint:dupl
func TestUserPasswordModify(t *testing.T) {
	Convey("TestUserPasswordModify", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(sessionSchema.HydraSessionsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userPasswordModifySchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserPasswordModifySchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		session := mock.NewMockHydraSession(ctrl)
		handler := mqHandler{
			hydraSession:             session,
			log:                      common.NewLogger(),
			sessionsSchema:           sessionsSchema,
			userPasswordModifySchema: userPasswordModifySchema,
		}

		msg := `{"user_id":"userID"}`

		Convey("msg 格式错误", func() {
			errMsg, _ := jsoniter.Marshal("[")
			err := handler.userPasswordModify(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 不存在", func() {
			errMsg, _ := jsoniter.Marshal("{}")
			err := handler.userPasswordModify(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg id 为空", func() {
			errMsg := `{"user_id":""}`
			err := handler.userPasswordModify([]byte(errMsg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 失败1", func() {
			testErr := errors.New("test")
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userPasswordModify([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败4", func() {
			testErr := &rest.HTTPError{
				Code: rest.InternalServerError,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userPasswordModify([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败5", func() {
			testErr := &rest.HTTPError{
				Code: rest.URINotExist,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userPasswordModify([]byte(msg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 成功", func() {
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			err := handler.userPasswordModify([]byte(msg))
			assert.Equal(t, err, nil)
		})
	})
}

//nolint:dupl
func TestUserStatusChange(t *testing.T) {
	Convey("TestUserStatusChange", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(sessionSchema.HydraSessionsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userStatusChanageSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserStatusChangeSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		session := mock.NewMockHydraSession(ctrl)
		handler := mqHandler{
			hydraSession:            session,
			log:                     common.NewLogger(),
			sessionsSchema:          sessionsSchema,
			userStatusChanageSchema: userStatusChanageSchema,
		}

		msg := `{"user_id":"userID", "status": false}`

		Convey("msg 格式错误", func() {
			errMsg, _ := jsoniter.Marshal("[")
			err := handler.userStatusChange(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg user_id 不存在", func() {
			errMsg, _ := jsoniter.Marshal("{}")
			err := handler.userStatusChange(errMsg)
			assert.Equal(t, err, nil)
		})

		Convey("msg user_id 为空", func() {
			errMsg := `{"user_id":""}`
			err := handler.userStatusChange([]byte(errMsg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 失败1", func() {
			testErr := errors.New("test")
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userStatusChange([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败4", func() {
			testErr := &rest.HTTPError{
				Code: rest.InternalServerError,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userStatusChange([]byte(msg))
			assert.Equal(t, err, testErr)
		})

		Convey("DeleteHydraSession 失败5", func() {
			testErr := &rest.HTTPError{
				Code: rest.URINotExist,
			}
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(testErr)
			err := handler.userStatusChange([]byte(msg))
			assert.Equal(t, err, nil)
		})

		Convey("DeleteHydraSession 成功", func() {
			session.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			err := handler.userStatusChange([]byte(msg))
			assert.Equal(t, err, nil)
		})
	})
}
