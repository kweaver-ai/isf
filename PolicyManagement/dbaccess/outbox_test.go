package dbaccess

import (
	"errors"
	"policy_mgnt/common"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func newOutboxDB(ptrDB *sqlx.DB) *outbox {
	return &outbox{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestAddOutboxInfo(t *testing.T) {
	Convey("AddOutboxInfo", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		outboxDB := newOutboxDB(db)

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = outboxDB.AddOutboxInfos(1, []string{""}, tx)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, 1, 1)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = outboxDB.AddOutboxInfos(1, []string{""}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetPushMessage(t *testing.T) {
	Convey("GetPushMessage", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		outboxDB := newOutboxDB(db)

		lockFields := []string{
			"f_business_type",
		}

		fields := []string{
			"f_id",
			"f_message",
		}

		Convey("get lock error", func() {
			mock.ExpectBegin()
			mock.ExpectQuery("").WillReturnError(errors.New("get lock unknown error"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			_, _, err = outboxDB.GetPushMessage(1, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("query error", func() {
			mock.ExpectBegin()
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(lockFields).AddRow(1))
			mock.ExpectQuery("").WillReturnError(errors.New("query unknown error"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			_, _, err = outboxDB.GetPushMessage(1, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("query not exist", func() {
			mock.ExpectBegin()
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(lockFields).AddRow(1))
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			messageID, message, err := outboxDB.GetPushMessage(1, tx)
			assert.Equal(t, err, nil)
			assert.Equal(t, messageID, int64(0))
			assert.Equal(t, message, "")
		})

		Convey("query exist", func() {
			mock.ExpectBegin()
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(lockFields).AddRow(1))
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(1, "xxx"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			messageID, message, err := outboxDB.GetPushMessage(1, tx)
			assert.Equal(t, err, nil)
			assert.Equal(t, messageID, int64(1))
			assert.Equal(t, message, "xxx")
		})
	})
}

func TestDeleteOutboxInfoByID(t *testing.T) {
	Convey("DeleteOutboxInfoByID", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		outboxDB := newOutboxDB(db)

		Convey("error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = outboxDB.DeleteOutboxInfoByID(1, tx)
			assert.NotEqual(t, err, nil)
		})

		var outboxID int64 = 1
		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = outboxDB.DeleteOutboxInfoByID(outboxID, tx)
			assert.Equal(t, err, nil)
		})
	})
}
