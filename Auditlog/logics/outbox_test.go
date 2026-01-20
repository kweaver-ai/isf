// Package logics outbox Anyshare 业务逻辑层 -outbox发件箱UT
package logics

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"AuditLog/interfaces/mock"
	"AuditLog/test"
	"AuditLog/test/mock_log"
)

func newOutboxDepend(t *testing.T) (*mock_log.MockLogger, *mock.MockDBOutbox) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl), mock.NewMockDBOutbox(ctrl)
}

func newOutbox(logger *mock_log.MockLogger, dbOutbox *mock.MockDBOutbox, dbpool *sqlx.DB) *outbox {
	return &outbox{
		db:           dbOutbox,
		pushChan:     make(chan struct{}, 1),
		logger:       logger,
		pool:         dbpool,
		businessType: "test",
	}
}

func TestAddOutboxInfo(t *testing.T) {
	Convey("AddOutboxInfo", t, func() {
		teardown := test.SetUpDB(t)
		defer teardown(t)

		dPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		logger, dbOutbox := newOutboxDepend(t)
		ob := newOutbox(logger, dbOutbox, dPool)

		outboxHandlers = make(map[string]func(interface{}) error)
		outboxHandlers["test"] = func(interface{}) error { return nil }

		Convey("db error", func() {
			dbOutbox.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("errors"))
			err := ob.AddOutboxInfo("test", "xxx", nil)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			dbOutbox.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			err := ob.AddOutboxInfo("test", "xxx", nil)
			assert.Equal(t, err, nil)
		})
	})
}

// func TestAddOutboxInfos(t *testing.T) {
// 	Convey("AddOutboxInfos", t, func() {
// 		teardown := test.SetUpDB(t)
// 		defer teardown(t)
// 		engine := gin.New()
// 		engine.Use(gin.Recovery())

// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		db := mock.NewMockDBOutbox(ctrl)
// 		dPool, _, err := sqlx.New()
// 		assert.Equal(t, err, nil)
// 		defer func() {
// 			if closeErr := dPool.Close(); closeErr != nil {
// 				assert.Equal(t, 1, 1)
// 			}
// 		}()

// 		ob := &outbox{
// 			db:           db,
// 			pushChan:     make(chan struct{}, 1),
// 			logger:       common.NewLogger(),
// 			pool:         dPool,
// 			businessType: 1,
// 		}

// 		outboxHandlers = make(map[int]func(interface{}) error)
// 		outboxHandlers[1] = func(interface{}) error { return nil }

// 		msg := interfaces.OutboxMsg{
// 			Type:    1,
// 			Content: "xxx",
// 		}

// 		Convey("type error", func() {
// 			msg.Type = 999
// 			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
// 			assert.Equal(t, err.(*rest.HTTPError).Code, rest.InternalError)
// 		})

// 		Convey("db error", func() {
// 			dbOutbox.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("errors"))
// 			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
// 			assert.NotEqual(t, err, nil)
// 		})

// 		Convey("success", func() {
// 			dbOutbox.EXPECT().AddOutboxInfos(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
// 			err := ob.AddOutboxInfos([]interfaces.OutboxMsg{msg}, nil)
// 			assert.Equal(t, err, nil)
// 		})
// 	})
// }

func TestNotifyPushOutboxThread(t *testing.T) {
	Convey("NotifyPushOutboxThread", t, func() {
		teardown := test.SetUpDB(t)
		defer teardown(t)

		dPool, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		logger, dbOutbox := newOutboxDepend(t)
		ob := newOutbox(logger, dbOutbox, dPool)

		Convey("success", func() {
			ob.NotifyPushOutboxThread()
			assert.Equal(t, 1, 1)
		})
	})
}

func TestPush(t *testing.T) {
	Convey("Push", t, func() {
		teardown := test.SetUpDB(t)
		defer teardown(t)

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		logger, dbOutbox := newOutboxDepend(t)
		ob := newOutbox(logger, dbOutbox, dPool)
		logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()

		testErr := errors.New("Error")
		outboxHandlers["test"] = func(content interface{}) error {
			if content.(string) == "error" {
				return testErr
			}
			return nil
		}
		outboxJSON1 := make(map[string]interface{})
		outboxJSON1["type"] = "test1"
		outboxJSON1["content"] = ""
		errTypeMsg, _ := jsoniter.MarshalToString(outboxJSON1)

		outboxJSON2 := make(map[string]interface{})
		outboxJSON2["type"] = "test"
		outboxJSON2["content"] = ""
		permMsg, _ := jsoniter.MarshalToString(outboxJSON2)

		outboxJSON3 := make(map[string]interface{})
		outboxJSON3["type"] = "test"
		outboxJSON3["content"] = "error"
		errPermMsg, _ := jsoniter.MarshalToString(outboxJSON3)
		logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		Convey("Get TX Failed", func() {
			txMock.ExpectBegin().WillReturnError(testErr)
			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Get OutBox Info Failed", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(0), "", testErr)
			txMock.ExpectRollback()
			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Push All Outbox Message, CommitError", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "", nil)
			txMock.ExpectCommit().WillReturnError(testErr)

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Push All Outbox Message", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "", nil)
			txMock.ExpectCommit()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, false)
			assert.Equal(t, isAllFinished, true)
		})

		Convey("jsoniter UnmarshalFromString Failed", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), "xxx", nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("ERROR message TYPE", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), errTypeMsg, nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("message handler function Error", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), errPermMsg, nil)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("DeleteOutboxInfoByID Error", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			dbOutbox.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Commit Error", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			dbOutbox.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit().WillReturnError(testErr)

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, true)
			assert.Equal(t, isAllFinished, false)
		})

		Convey("Successful", func() {
			txMock.ExpectBegin()
			dbOutbox.EXPECT().GetPushMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(int64(1), permMsg, nil)
			dbOutbox.EXPECT().DeleteOutboxInfoByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()

			needTimer, isAllFinished := ob.push()
			assert.Equal(t, needTimer, false)
			assert.Equal(t, isAllFinished, false)
		})
	})
}
