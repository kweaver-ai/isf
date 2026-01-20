package dbaccess

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"Authentication/common"
	"Authentication/interfaces"
)

func TestGetAuditLogAsyncTaskInfo(t *testing.T) {
	Convey("TestGetAuditLogAsyncTaskInfo", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		un = &unorderedOutbox{
			db:     db,
			logger: common.NewLogger(),
		}
		Convey("sql no rows", func() {
			mock.ExpectQuery("^select").WillReturnError(sql.ErrNoRows)
			_, exist, err := un.GetUnorderedOutboxInfo()
			assert.Equal(t, exist, false)
			assert.Equal(t, err, nil)
		})
		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(errors.New("query error"))
			_, exist, err := un.GetUnorderedOutboxInfo()
			assert.Equal(t, exist, false)
			assert.Equal(t, err, errors.New("query error"))
		})

		Convey("update error", func() {
			mock.ExpectQuery("^select").WillReturnRows(sqlmock.NewRows([]string{"id", "f_message", "f_status"}).AddRow("1", "", 0))
			mock.ExpectExec("^update").WillReturnError(errors.New("update error"))
			_, exist, err := un.GetUnorderedOutboxInfo()
			assert.Equal(t, exist, false)
			assert.Equal(t, err, errors.New("update error"))
		})
		Convey("update success", func() {
			mock.ExpectQuery("^select").WillReturnRows(sqlmock.NewRows([]string{"id", "f_message", "f_status"}).AddRow("1", "", 0))
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			auditLogAsyncInfo, exist, err := un.GetUnorderedOutboxInfo()
			assert.Equal(t, auditLogAsyncInfo.ID, "1")
			assert.Equal(t, auditLogAsyncInfo.Message, "")
			assert.Equal(t, auditLogAsyncInfo.Status, interfaces.OutboxInProgress)
			assert.Equal(t, exist, true)
			assert.Equal(t, err, nil)
		})
		Convey("get info cycle 1", func() {
			mock.ExpectQuery("^select").WillReturnRows(sqlmock.NewRows([]string{"id", "f_message", "f_status"}).AddRow("1", "", 0))
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectQuery("^select").WillReturnRows(sqlmock.NewRows([]string{"id", "f_message", "f_status"}).AddRow("2", "", 0))
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			auditLogAsyncInfo, exist, err := un.GetUnorderedOutboxInfo()
			assert.Equal(t, auditLogAsyncInfo.ID, "2")
			assert.Equal(t, auditLogAsyncInfo.Message, "")
			assert.Equal(t, auditLogAsyncInfo.Status, interfaces.OutboxInProgress)
			assert.Equal(t, exist, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAuditLogAsyncTaskByID(t *testing.T) {
	Convey("DeleteUnorderedOutboxInfoByID", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		un = &unorderedOutbox{
			db:     db,
			logger: common.NewLogger(),
		}
		Convey("delete error", func() {
			mock.ExpectExec("^delete from authentication.t_outbox_unordered").WithArgs(sqlmock.AnyArg()).WillReturnError(errors.New("delete error"))
			err := un.DeleteUnorderedOutboxInfoByID("1")
			assert.Equal(t, err, errors.New("delete error"))
		})
		Convey("update success", func() {
			mock.ExpectExec("^delete from authentication.t_outbox_unordered").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := un.DeleteUnorderedOutboxInfoByID("1")
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateAuditLogAsyncTaskUpdateTimeByID(t *testing.T) {
	Convey("UpdateUnorderedOutboxUpdateTimeByID", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		un = &unorderedOutbox{
			db:     db,
			logger: common.NewLogger(),
		}
		Convey("update error", func() {
			mock.ExpectExec("^update").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("update error"))
			isUpdate, err := un.UpdateUnorderedOutboxUpdateTimeByID("1")
			assert.Equal(t, isUpdate, false)
			assert.Equal(t, err, errors.New("update error"))
		})
		Convey("affect rows 0", func() {
			mock.ExpectExec("^update").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0))
			isUpdate, err := un.UpdateUnorderedOutboxUpdateTimeByID("1")
			assert.Equal(t, isUpdate, false)
			assert.Equal(t, err, nil)
		})
		Convey("affect rows 1", func() {
			mock.ExpectExec("^update").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
			isUpdate, err := un.UpdateUnorderedOutboxUpdateTimeByID("1")
			assert.Equal(t, isUpdate, true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddAuditLogAsyncTask(t *testing.T) {
	Convey("AddUnorderedOutboxInfo", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		un = &unorderedOutbox{
			db:     db,
			logger: common.NewLogger(),
		}
		auditLogAsyncInfo := interfaces.UnorderedOutbox{
			ID:      "1",
			Message: "",
			Status:  interfaces.OutboxNotStarted,
		}
		Convey("insert error", func() {
			mock.ExpectExec("^INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("insert error"))
			err := un.AddUnorderedOutboxInfo(auditLogAsyncInfo)
			assert.Equal(t, err, errors.New("insert error"))
		})
		Convey("success", func() {
			mock.ExpectExec("^INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := un.AddUnorderedOutboxInfo(auditLogAsyncInfo)
			assert.Equal(t, err, nil)
		})
	})
}

func TestRestartUnorderedOutboxInfo(t *testing.T) {
	Convey("RestartUnorderedOutboxInfo", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		un = &unorderedOutbox{
			db:     db,
			logger: common.NewLogger(),
		}
		Convey("update error", func() {
			mock.ExpectExec("^update").WithArgs(sqlmock.AnyArg()).WillReturnError(errors.New("update error"))
			err := un.RestartUnorderedOutboxInfo(123)
			assert.Equal(t, err, errors.New("update error"))
		})
		Convey("success", func() {
			mock.ExpectExec("^update").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := un.RestartUnorderedOutboxInfo(123)
			assert.Equal(t, err, nil)
		})
	})
}
