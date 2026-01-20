// Package logics outbox Anyshare 业务逻辑层 -outbox发件箱UT
package logics

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
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

func TestAddOutboxInfo(t *testing.T) {
	Convey("AddOutboxInfo", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBOutbox(ctrl)
		dPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ob := &outbox{
			db:           db,
			pushChan:     make(chan struct{}, 1),
			logger:       common.NewLogger(),
			pool:         dPool,
			businessType: 1,
		}

		outboxHandlers = make(map[int]func(interface{}) error)
		outboxHandlers[1] = func(interface{}) error { return nil }

		Convey("type error", func() {
			err := ob.AddOutboxInfo(999, "xxx", nil)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.InternalServerError)
		})

		Convey("db error", func() {
			db.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("errors"))
			err := ob.AddOutboxInfo(1, "xxx", nil)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			db.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			err := ob.AddOutboxInfo(1, "xxx", nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddOutboxInfos(t *testing.T) {
	Convey("AddOutboxInfos", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBOutbox(ctrl)
		dPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ob := &outbox{
			db:           db,
			pushChan:     make(chan struct{}, 1),
			logger:       common.NewLogger(),
			pool:         dPool,
			businessType: 1,
		}

		outboxHandlers = make(map[int]func(interface{}) error)
		outboxHandlers[1] = func(interface{}) error { return nil }

		msg := interfaces.OutboxMsg{
			Type:    1,
			Content: "xxx",
		}

		Convey("type error", func() {
			msg.Type = 999
			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.InternalServerError)
		})

		Convey("db error", func() {
			db.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("errors"))
			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			db.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
			assert.Equal(t, err, nil)
		})
	})
}

func TestNotifyPushOutboxThread(t *testing.T) {
	Convey("NotifyPushOutboxThread", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBOutbox(ctrl)
		dPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ob := &outbox{
			db:           db,
			pushChan:     make(chan struct{}, 1),
			logger:       common.NewLogger(),
			pool:         dPool,
			businessType: 1,
		}

		Convey("success", func() {
			ob.NotifyPushOutboxThread()
			assert.Equal(t, 1, 1)
		})
	})
}

func TestPush(t *testing.T) {
	Convey("Push", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBOutbox(ctrl)
		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ob := &outbox{
			db:           db,
			pushChan:     make(chan struct{}, 1),
			logger:       common.NewLogger(),
			pool:         dPool,
			businessType: 1,
		}

		testErr := errors.New("Error")
		outboxHandlers[1] = func(content interface{}) error {
			if content.(string) == "error" {
				return testErr
			}
			return nil
		}
		outboxJSON1 := make(map[string]interface{})
		outboxJSON1["type"] = 999
		outboxJSON1["content"] = ""
		errTypeMsg, _ := jsoniter.MarshalToString(outboxJSON1)

		outboxJSON2 := make(map[string]interface{})
		outboxJSON2["type"] = 1
		outboxJSON2["content"] = ""
		permMsg, _ := jsoniter.MarshalToString(outboxJSON2)

		outboxJSON3 := make(map[string]interface{})
		outboxJSON3["type"] = 1
		outboxJSON3["content"] = "error"
		errPermMsg, _ := jsoniter.MarshalToString(outboxJSON3)

		Convey("Get TX Failed", func() {
			txMock.ExpectBegin().WillReturnError(testErr)
			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Get OutBox Info Failed", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(0), "", testErr)
			txMock.ExpectRollback()
			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Push All Outbox Message, CommitError", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "", nil)
			txMock.ExpectCommit().WillReturnError(testErr)

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Push All Outbox Message", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "", nil)
			txMock.ExpectCommit()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, false)
			assert.Equal(t, isAllFinished, true)
		})

		Convey("jsoniter UnmarshalFromString Failed", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "xxx", nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("ERROR message TYPE", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), errTypeMsg, nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("message handler function Error", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), errPermMsg, nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("DeleteOutboxInfoByID Error", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			db.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Commit Error", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			db.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit().WillReturnError(testErr)

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Successful", func() {
			txMock.ExpectBegin()
			db.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			db.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, false)
			assert.Equal(t, isAllFinished, false)
		})
	})
}
