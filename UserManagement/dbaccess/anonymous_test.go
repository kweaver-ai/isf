package dbaccess

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"UserManagement/common"
	"UserManagement/interfaces"
)

func newAnonymous(ptrDB *sqlx.DB) *anonymous {
	return &anonymous{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestNewAnonymous(t *testing.T) {
	Convey("NewAnonymous, db is available", t, func() {
		data := NewAnonymous()
		assert.NotEqual(t, data, nil)
	})
}

func TestCreate(t *testing.T) {
	Convey("Create, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		anonymous := newAnonymous(db)

		var info interfaces.AnonymousInfo
		Convey("query error", func() {
			mock.ExpectQuery("").WillReturnError(errors.New("unknown"))
			err := anonymous.Create(&info)
			assert.NotEqual(t, err, nil)
		})

		Convey("query none", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_password"}))
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := anonymous.Create(&info)
			assert.Equal(t, err, nil)
		})

		Convey("query one success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_password"}).AddRow("11"))
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := anonymous.Create(&info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteByID(t *testing.T) {
	Convey("DeleteByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		anonymous := newAnonymous(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := anonymous.DeleteByID("")
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := anonymous.DeleteByID("")
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAccount(t *testing.T) {
	Convey("GetAccount, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		anonymous := newAnonymous(db)

		fields := []string{
			"f_anonymity_id",
			"f_password",
			"f_expires_at",
			"f_limited_times",
			"f_accessed_times",
			"f_type",
			"f_verify_mobile",
		}

		Convey("contacor not exist", func() {
			mock.ExpectQuery("").WillReturnError(errors.New("unknown"))
			_, httpErr := anonymous.GetAccount("")
			assert.NotEqual(t, httpErr, nil)
		})

		Convey("success", func() {
			info := interfaces.AnonymousInfo{
				ID:             "zzz",
				Password:       "kkk",
				ExpiresAtStamp: 1,
				LimitedTimes:   2,
				AccessedTimes:  3,
				Type:           "document",
				VerifyMobile:   false,
			}

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(info.ID, info.Password, info.ExpiresAtStamp, info.LimitedTimes, info.AccessedTimes, info.Type, 0))
			out, httpErr := anonymous.GetAccount("")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, out, info)
		})
	})
}

func TestAddAccessTimes(t *testing.T) {
	Convey("AddAccessTimes, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		anonymous := newAnonymous(db)

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = anonymous.AddAccessTimes("", tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			err = anonymous.AddAccessTimes("", tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAnonymoysOutOfDate(t *testing.T) {
	Convey("DeleteAnonymoysOutOfDate, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		anonymous := newAnonymous(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := anonymous.DeleteByTime(0)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := anonymous.DeleteByTime(0)
			assert.Equal(t, err, nil)
		})
	})
}
