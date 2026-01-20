package audit

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"

	"github.com/go-playground/assert"
)

func TestAddUnorderedLog(t *testing.T) {
	Convey("添加信息", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbOb := mock.NewMockDBUnorderedOutbox(ctrl)
		dnEacpLog := mock.NewMockDnEacpLog(ctrl)

		auditLog := &auditLogAsyncTask{
			unOb:     dbOb,
			pushChan: make(chan struct{}, 1),
			eacpLog:  dnEacpLog,
			logger:   common.NewLogger(),
		}

		Convey("AddUnorderedOutboxInfo error", func() {
			dbOb.EXPECT().AddUnorderedOutboxInfo(gomock.Any()).Return(errors.New("AddUnorderedOutboxInfo error"))
			err := auditLog.Log("as.audit_log.log_operation", "")
			assert.Equal(t, err, errors.New("AddUnorderedOutboxInfo error"))
		})
		Convey("success", func() {
			dbOb.EXPECT().AddUnorderedOutboxInfo(gomock.Any()).Return(nil)
			err := auditLog.Log("as.audit_log.log_operation", "")
			assert.Equal(t, err, nil)
		})
	})
}

func TestAuditLogPush(t *testing.T) {
	Convey("发送审计日志", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbOb := mock.NewMockDBUnorderedOutbox(ctrl)
		dnEacpLog := mock.NewMockDnEacpLog(ctrl)

		auditLog := &auditLogAsyncTask{
			unOb:     dbOb,
			pushChan: make(chan struct{}, 1),
			eacpLog:  dnEacpLog,
			logger:   common.NewLogger(),
		}
		Convey("RestartUnorderedOutboxInfo error", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(errors.New("RestartUnorderedOutboxInfo error"))
			isAllFinished := auditLog.push()
			assert.Equal(t, isAllFinished, false)
		})
		Convey("GetUnorderedOutboxInfo error", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(nil)
			dbOb.EXPECT().GetUnorderedOutboxInfo().Return(interfaces.UnorderedOutbox{}, false, errors.New("GetUnorderedOutboxInfo error"))
			isAllFinished := auditLog.push()
			assert.Equal(t, isAllFinished, false)
		})
		Convey("finished all UnorderedOutboxInfo  ", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(nil)
			dbOb.EXPECT().GetUnorderedOutboxInfo().Return(interfaces.UnorderedOutbox{}, false, nil)
			isAllFinished := auditLog.push()
			assert.Equal(t, isAllFinished, true)
		})
		Convey("success", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(nil)
			dbOb.EXPECT().GetUnorderedOutboxInfo().Return(interfaces.UnorderedOutbox{ID: "",
				Message: "{}"}, true, nil)
			dbOb.EXPECT().UpdateUnorderedOutboxUpdateTimeByID(gomock.Any()).AnyTimes().Return(true, nil)
			dnEacpLog.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			dbOb.EXPECT().DeleteUnorderedOutboxInfoByID(gomock.Any()).Return(nil)
			isAllFinished := auditLog.push()
			assert.Equal(t, isAllFinished, false)
		})
	})
}

func TestAuditLogAutoRenewal(t *testing.T) {
	Convey("自动续期", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbOb := mock.NewMockDBUnorderedOutbox(ctrl)
		dnEacpLog := mock.NewMockDnEacpLog(ctrl)

		auditLog := &auditLogAsyncTask{
			unOb:     dbOb,
			pushChan: make(chan struct{}, 1),
			eacpLog:  dnEacpLog,
			logger:   common.NewLogger(),
		}
		id := "xxx"
		Convey("UpdateUnorderedOutboxUpdateTimeByID error", func() {
			dbOb.EXPECT().UpdateUnorderedOutboxUpdateTimeByID(gomock.Any()).Return(false, errors.New("UpdateUnorderedOutboxUpdateTimeByID error"))
			runAble, err := auditLog.autoRenewal(id)
			assert.Equal(t, runAble, false)
			assert.Equal(t, err, errors.New("UpdateUnorderedOutboxUpdateTimeByID error"))
		})
		Convey("success, affected rows is 0", func() {
			dbOb.EXPECT().UpdateUnorderedOutboxUpdateTimeByID(gomock.Any()).Return(false, nil)
			runAble, err := auditLog.autoRenewal(id)
			assert.Equal(t, runAble, false)
			assert.Equal(t, err, nil)
		})
		Convey("success, affected rows is not 0", func() {
			dbOb.EXPECT().UpdateUnorderedOutboxUpdateTimeByID(gomock.Any()).Return(true, nil)
			runAble, err := auditLog.autoRenewal(id)
			assert.Equal(t, runAble, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAuditLogRestartUnAsyncTask(t *testing.T) {
	Convey("重置异常的任务", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbOb := mock.NewMockDBUnorderedOutbox(ctrl)
		dnEacpLog := mock.NewMockDnEacpLog(ctrl)

		auditLog := &auditLogAsyncTask{
			unOb:     dbOb,
			pushChan: make(chan struct{}, 1),
			eacpLog:  dnEacpLog,
			logger:   common.NewLogger(),
		}
		Convey("RestartUnorderedOutboxInfo error", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(errors.New("RestartUnorderedOutboxInfo error"))
			auditLog.restartUnAsyncTask()
		})
		Convey("success", func() {
			dbOb.EXPECT().RestartUnorderedOutboxInfo(gomock.Any()).Return(nil)
			auditLog.restartUnAsyncTask()
		})
	})
}
