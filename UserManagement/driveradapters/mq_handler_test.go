package driveradapters

import (
	"errors"
	"fmt"
	"testing"
	"time"

	errorv2 "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func newMQHandler(an interfaces.LogicsAnonymous, c interfaces.LogicsCombine) *mqHandler {
	return &mqHandler{
		logger:    common.NewLogger(),
		combine:   c,
		anonymous: an,
	}
}

func TestCreateAnonymousNoParams(t *testing.T) {
	Convey("CreateAnonymous", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		an := mock.NewMockLogicsAnonymous(ctrl)
		handler := newMQHandler(an, nil)

		Convey("no ID", func() {
			type CreateAnonymousParams struct {
				LimitTimes float64 `json:"limited_times" binding:"required"`
				ExpiresAt  string  `json:"expires_at" binding:"required"`
				Password   string  `json:"password" binding:"required"`
			}

			reqParam := CreateAnonymousParams{
				LimitTimes: 1,
				ExpiresAt:  "xxxx",
				Password:   "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("no password", func() {
			type CreateAnonymousParams struct {
				LimitTimes float64 `json:"limited_times" binding:"required"`
				ExpiresAt  string  `json:"expires_at" binding:"required"`
				ID         string  `json:"id" binding:"required"`
			}

			reqParam := CreateAnonymousParams{
				LimitTimes: 1,
				ExpiresAt:  "xxxx",
				ID:         "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("no limited times", func() {
			type CreateAnonymousParams struct {
				Password  string `json:"password" binding:"required"`
				ExpiresAt string `json:"expires_at" binding:"required"`
				ID        string `json:"id" binding:"required"`
			}

			reqParam := CreateAnonymousParams{
				Password:  "xxxx",
				ExpiresAt: "xxxx",
				ID:        "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("no expires at", func() {
			type CreateAnonymousParams struct {
				Password   string  `json:"password" binding:"required"`
				LimitTimes float64 `json:"limited_times" binding:"required"`
				ID         string  `json:"id" binding:"required"`
			}

			reqParam := CreateAnonymousParams{
				Password:   "xxxx",
				LimitTimes: 1,
				ID:         "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCreateAnonymousParamsError(t *testing.T) {
	Convey("CreateAnonymous", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		an := mock.NewMockLogicsAnonymous(ctrl)
		handler := newMQHandler(an, nil)

		type CreateAnonymousParams struct {
			ExpiresAt  string  `json:"expires_at" binding:"required"`
			Password   string  `json:"password" binding:"required"`
			LimitTimes float64 `json:"limited_times" binding:"required"`
			ID         string  `json:"id" binding:"required"`
			Type       string  `json:"type" binding:"required"`
		}

		Convey("password err", func() {
			reqParam := CreateAnonymousParams{
				LimitTimes: 10,
				ExpiresAt:  "2002-10-02T15:00:00Z",
				Password:   "z",
				ID:         "zzzz",
				Type:       "document",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("expires at error", func() {
			reqParam := CreateAnonymousParams{
				Password:   "xxxxzzzzz",
				ExpiresAt:  "123131",
				LimitTimes: 10,
				ID:         "zzzz",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("limited times error", func() {
			common.SvcConfig.BusinessTimeOffset = 0

			reqParam := CreateAnonymousParams{
				Password:   "xxxxzzzzz",
				ExpiresAt:  "2002-10-02T15:00:00Z",
				LimitTimes: -2,
				ID:         "zzzz",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestCreateAnonymous(t *testing.T) {
	Convey("CreateAnonymous", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		an := mock.NewMockLogicsAnonymous(ctrl)
		handler := newMQHandler(an, nil)

		common.SvcConfig.BusinessTimeOffset = 0

		type CreateGroupParams struct {
			LimitTimes int32  `json:"limited_times" binding:"required"`
			ExpiresAt  string `json:"expires_at" binding:"required"`
			Password   string `json:"password" binding:"required"`
			ID         string `json:"id" binding:"required"`
			Type       string `json:"type" binding:"required"`
		}
		t1 := fmt.Sprint(time.Now().Year() + 1)

		Convey("Create error", func() {
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			reqParam := CreateGroupParams{
				Password:   "xxxxzzzzz",
				ExpiresAt:  t1 + "-10-02T15:00:00Z",
				LimitTimes: 10,
				ID:         "zzzz",
				Type:       "document",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			an.EXPECT().Create(gomock.Any()).AnyTimes().Return(testErr)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("Create Success", func() {
			reqParam := CreateGroupParams{
				Password:   "xxxxzzzzz",
				ExpiresAt:  t1 + "-10-02T15:00:00Z",
				LimitTimes: 10,
				ID:         "zzzz",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			an.EXPECT().Create(gomock.Any()).AnyTimes().Return(nil)

			err := handler.onCreateAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAnonymous(t *testing.T) {
	Convey("deleteAnonymous", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		an := mock.NewMockLogicsAnonymous(ctrl)
		handler := newMQHandler(an, nil)

		common.SvcConfig.BusinessTimeOffset = 0

		type DeleteParamsErr struct {
			ID int `json:"ids" binding:"required"`
		}

		Convey("params error", func() {
			reqParam := DeleteParamsErr{
				ID: 111,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.onDeleteAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		type DeleteParams struct {
			ID []string `json:"ids" binding:"required"`
		}

		Convey("DeleteByID error", func() {
			reqParam := DeleteParams{
				ID: []string{"zzz", "zzz"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			an.EXPECT().DeleteByID(gomock.Any()).AnyTimes().Return(testErr)

			err := handler.onDeleteAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("Success", func() {
			reqParam := DeleteParams{
				ID: []string{"zzz", "zzz"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			an.EXPECT().DeleteByID(gomock.Any()).AnyTimes().Return(nil)

			err := handler.onDeleteAnonymous(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestOnDeleteUser(t *testing.T) {
	Convey("OnUserDeleted", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ev := mock.NewMockLogicsEvent(ctrl)
		handler := newMQHandler(nil, nil)
		handler.event = ev

		common.SvcConfig.BusinessTimeOffset = 0

		type DeleteParamsErr struct {
			ID string `json:"ids" binding:"required"`
		}

		Convey("params error", func() {
			reqParam := DeleteParamsErr{
				ID: "xxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.OnUserDeleted(reqParamByte)
			assert.Equal(t, err, nil)
		})

		type DeleteParamsRight struct {
			ID string `json:"id" binding:"required"`
		}
		Convey("OnUserDeleted error", func() {
			reqParam := DeleteParamsRight{
				ID: "xxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			testErr := errors.New("xxxx")
			ev.EXPECT().UserDeleted(gomock.Any()).AnyTimes().Return(testErr)

			err := handler.OnUserDeleted(reqParamByte)
			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			reqParam := DeleteParamsRight{
				ID: "xxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			ev.EXPECT().UserDeleted(gomock.Any()).AnyTimes().Return(nil)

			err := handler.OnUserDeleted(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestMQDeleteManageDepart(t *testing.T) {
	Convey("OnDepartDeleted", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combin := mock.NewMockLogicsCombine(ctrl)
		event := mock.NewMockLogicsEvent(ctrl)
		handler := &mqHandler{
			logger:  common.NewLogger(),
			combine: combin,
			event:   event,
		}

		common.SvcConfig.BusinessTimeOffset = 0

		type DeleteParamsErr struct {
			ID string `json:"ids" binding:"required"`
		}

		Convey("params error", func() {
			reqParam := DeleteParamsErr{
				ID: "xxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.OnDepartDeleted(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("success", func() {
			reqParam := map[string]interface{}{
				"ids": []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			event.EXPECT().DeptDeleted(gomock.Any()).AnyTimes().Return(nil)
			err := handler.OnDepartDeleted(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestMQDepartResponserChanged(t *testing.T) {
	Convey("OrgManagerChanged", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		idsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(idsSchemaStr))
		assert.Equal(t, err, nil)

		combin := mock.NewMockLogicsCombine(ctrl)
		event := mock.NewMockLogicsEvent(ctrl)
		handler := &mqHandler{
			logger:    common.NewLogger(),
			combine:   combin,
			event:     event,
			idsSchema: idsSchema,
		}

		common.SvcConfig.BusinessTimeOffset = 0

		Convey("params error", func() {
			reqParam := map[string]interface{}{
				"id": []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			err := handler.OrgManagerChanged(reqParamByte)
			assert.Equal(t, err, nil)
		})

		Convey("success", func() {
			reqParam := map[string]interface{}{
				"ids": []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			event.EXPECT().OrgManagerChanged(gomock.Any()).AnyTimes().Return(nil)
			err := handler.OrgManagerChanged(reqParamByte)
			assert.Equal(t, err, nil)
		})
	})
}

func TestNeedReTry(t *testing.T) {
	Convey("NeedReTry", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combine := mock.NewMockLogicsCombine(ctrl)
		handler := newMQHandler(nil, combine)

		Convey("rest err  nil", func() {
			tmp := handler.needReTry(nil)
			assert.Equal(t, tmp, false)
		})

		Convey("rest err", func() {
			testErr := rest.NewHTTPError("error", rest.BadRequest, nil)
			tmp := handler.needReTry(testErr)
			assert.Equal(t, tmp, false)
		})

		Convey("status code 500", func() {
			testErr := &errorv2.Error{
				Code: errorv2.PublicInternalServerError,
			}
			tmp := handler.needReTry(testErr)
			assert.Equal(t, tmp, true)
		})

		Convey("status code 450", func() {
			testErr := &errorv2.Error{
				Code: "Public.BlockedByWindowsParentalControls",
			}
			tmp := handler.needReTry(testErr)
			assert.Equal(t, tmp, false)
		})

		Convey("common err", func() {
			testErr := errors.New("common err")
			tmp := handler.needReTry(testErr)
			assert.Equal(t, tmp, true)
		})
	})
}
