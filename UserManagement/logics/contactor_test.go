package logics

import (
	"context"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func newContactor(cdb interfaces.DBContactor) *contactor {
	return &contactor{
		db: cdb,
	}
}

func TestConvertContactorName(t *testing.T) {
	Convey("ConvertContactorName, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contactorDB := mock.NewMockDBContactor(ctrl)
		contactorLogics := newContactor(contactorDB)
		contactorIDs := make([]string, 0)

		Convey("contactor id is empty", func() {
			outContactorIDs, err := contactorLogics.ConvertContactorName(contactorIDs, true)
			assert.Equal(t, len(outContactorIDs), 0)
			assert.Equal(t, err, nil)
		})

		contactorIDs = append(contactorIDs, "user_id")
		Convey("DB GetDepartmentName error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			outContactorIDs, err := contactorLogics.ConvertContactorName(contactorIDs, true)
			assert.Equal(t, len(outContactorIDs), 0)
			assert.Equal(t, err, testErr)
		})

		Convey("contactor is not exist", func() {
			testErr := rest.NewHTTPError("record not exist", errors.ContactorNotFound, map[string]interface{}{"ids": []string{"user_id"}})
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, nil, nil)
			outContactorIDs, err := contactorLogics.ConvertContactorName(contactorIDs, true)
			assert.Equal(t, len(outContactorIDs), 0)
			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			testLogics := interfaces.NameInfo{ID: "user_id", Name: "user_name"}
			testInfoMap := make([]interfaces.NameInfo, 0)
			testInfoMap = append(testInfoMap, testLogics)
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(testInfoMap, []string{testLogics.Name}, nil)
			outContactorIDs, err := contactorLogics.ConvertContactorName(contactorIDs, true)
			assert.Equal(t, len(outContactorIDs), 1)
			assert.Equal(t, outContactorIDs[0], testLogics)
			assert.Equal(t, err, nil)
		})

		Convey("success1", func() {
			testLogics := interfaces.NameInfo{ID: "user_id", Name: "user_name"}
			testInfoMap := make([]interfaces.NameInfo, 0)
			testInfoMap = append(testInfoMap, testLogics)
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(testInfoMap, []string{testLogics.Name}, nil)
			outContactorIDs, err := contactorLogics.ConvertContactorName([]string{strID, strID2}, false)
			assert.Equal(t, len(outContactorIDs), 1)
			assert.Equal(t, outContactorIDs[0], testLogics)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteContactor(t *testing.T) {
	Convey("ConvertContactorName, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		contactorDB := mock.NewMockDBContactor(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		msgBroker := mock.NewMockDrivenMessageBroker(ctrl)
		contactorLogics := &contactor{
			db:            contactorDB,
			pool:          dPool,
			ob:            ob,
			logger:        common.NewLogger(),
			messageBroker: msgBroker,
		}

		contactorID := ""
		info := interfaces.ContactorInfo{}
		visitor := &interfaces.Visitor{}
		Convey("contactorID 非UUID", func() {
			testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		contactorID = "e9348980-6868-11ec-8b2c-fa9a2cadb275"
		Convey("获取联系人组信息失败", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(false, info, testErr)
			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		Convey("联系人组不存在", func() {
			testErr := rest.NewHTTPError("group is not exist", rest.BadRequest, nil)
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(false, info, nil)
			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		info.UserID = "aa"
		Convey("联系人组不是指定用户的", func() {
			testErr := rest.NewHTTPError("group is not exist", rest.BadRequest, nil)
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		info.UserID = ""
		Convey("db pool begin 失败", func() {
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			testErr := rest.NewHTTPError("error", 503000000, nil)
			txMock.ExpectBegin().WillReturnError(testErr)

			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		Convey("删除联系人组成员失败", func() {
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			testErr := rest.NewHTTPError("error", 503000000, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		Convey("删除联系人组失败", func() {
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			testErr := rest.NewHTTPError("error", 503000000, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		Convey("插入outbox失败", func() {
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			testErr := rest.NewHTTPError("error", 503000000, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()

			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, testErr)
		})

		Convey("删除成功", func() {
			contactorDB.EXPECT().GetContactorInfo(gomock.Any()).AnyTimes().Return(true, info, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().NotifyPushOutboxThread()

			txMock.ExpectCommit()

			err := contactorLogics.DeleteContactor(visitor, contactorID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestContactorOnUserDeleted(t *testing.T) {
	Convey("用户被删除事件处理", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contactorDB := mock.NewMockDBContactor(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		u := &contactor{
			db:     contactorDB,
			pool:   dPool,
			ob:     ob,
			logger: common.NewLogger(),
		}

		Convey("获取用户所有联系人组报错", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(nil, tempErr1)
			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		info := interfaces.ContactorInfo{
			ContactorID: "zzz",
		}
		contactorInfos := []interfaces.ContactorInfo{info}
		Convey("dbpool Begin 失败", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin().WillReturnError(tempErr1)

			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		Convey("DeleteContactorMembers error", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr1)
			txMock.ExpectRollback()

			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		Convey("DeleteContactors error", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr1)
			txMock.ExpectRollback()

			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		Convey("DeleteUserInContactors error", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteUserInContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr1)
			txMock.ExpectRollback()

			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		Convey("AddOutboxInfo error", func() {
			tempErr1 := rest.NewHTTPError("Unsupported user type", rest.Forbidden, nil)
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteUserInContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr1)
			txMock.ExpectRollback()

			out := u.onUserDeleted("")
			assert.Equal(t, out, tempErr1)
		})

		Convey("success", func() {
			contactorDB.EXPECT().GetAllContactorInfos(gomock.Any()).AnyTimes().Return(contactorInfos, nil)
			txMock.ExpectBegin()
			contactorDB.EXPECT().DeleteContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().DeleteUserInContactors(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			contactorDB.EXPECT().UpdateContactorCount().AnyTimes().Return(nil)
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()

			out := u.onUserDeleted("")
			assert.Equal(t, out, nil)
		})
	})
}

func TestSendContactorDeletedMsg(t *testing.T) {
	Convey("sendContactorDeletedMsg, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contactorDB := mock.NewMockDBContactor(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		msgBroker := mock.NewMockDrivenMessageBroker(ctrl)
		contactorLogics := &contactor{
			db:            contactorDB,
			ob:            ob,
			logger:        common.NewLogger(),
			messageBroker: msgBroker,
		}

		msg := make(map[string]interface{})
		msg["ids"] = []string{"xxx", "zzz"}

		outboxJSON := make(map[string]interface{})
		outboxJSON["type"] = outboxContactorDeleted
		outboxJSON["content"] = msg
		outboxMsg, _ := jsoniter.MarshalToString(outboxJSON)

		var messageJSON interface{}
		err := jsoniter.UnmarshalFromString(outboxMsg, &messageJSON)
		assert.Equal(t, err, nil)

		content := messageJSON.(map[string]interface{})["content"]

		Convey("发送消息失败", func() {
			testErr := rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
			msgBroker.EXPECT().ContactorDeleted(gomock.Any()).AnyTimes().Return(testErr)
			err := contactorLogics.sendContactorDeletedMsg(content)
			assert.Equal(t, err, testErr)
		})

		Convey("发送消息成功", func() {
			msgBroker.EXPECT().ContactorDeleted(gomock.Any()).AnyTimes().Return(nil)
			err := contactorLogics.sendContactorDeletedMsg(content)
			assert.Equal(t, err, nil)
		})
	})
}

// TestGetContactorMemberIDs 测试批量获取联系人组成员ID
func TestGetContactorMemberIDs(t *testing.T) {
	Convey("GetContactorMemberIDs, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contactorDB := mock.NewMockDBContactor(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		msgBroker := mock.NewMockDrivenMessageBroker(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		ctx := context.Background()

		contactorLogics := &contactor{
			db:            contactorDB,
			ob:            ob,
			logger:        common.NewLogger(),
			messageBroker: msgBroker,
			trace:         trace,
		}

		Convey("GetContactorName error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, nil, testErr)
			_, httpErr := contactorLogics.GetContactorMembers(ctx, []string{"test", "test1"})
			assert.Equal(t, httpErr, testErr)
		})

		Convey("GetContactorMemberIDs error", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, nil, nil)
			contactorDB.EXPECT().GetContactorMemberIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, httpErr := contactorLogics.GetContactorMembers(ctx, []string{"test", "test1"})
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, []string{"test"}, nil)
			contactorDB.EXPECT().GetContactorMemberIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(map[string][]string{"test": {"test1", "test2"}}, nil)

			infos, httpErr := contactorLogics.GetContactorMembers(ctx, []string{"test"})
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(infos), 1)
			assert.Equal(t, infos[0].ContactorID, "test")
			assert.Equal(t, infos[0].MemberIDs, []string{"test1", "test2"})
		})

		Convey("success1", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			contactorDB.EXPECT().GetContactorName(gomock.Any()).AnyTimes().Return(nil, []string{"test"}, nil)
			contactorDB.EXPECT().GetContactorMemberIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(map[string][]string{}, nil)

			infos, httpErr := contactorLogics.GetContactorMembers(ctx, []string{"test"})
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(infos), 1)
			assert.Equal(t, infos[0].ContactorID, "test")
			assert.Equal(t, len(infos[0].MemberIDs), 0)
		})
	})
}
