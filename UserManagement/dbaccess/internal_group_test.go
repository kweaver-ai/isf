// dbaccess 数据库模块ut
package dbaccess

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"UserManagement/common"
)

const (
	strID  = "sdadadassa"
	strID1 = "zzsdasdada"
)

func TestNewInternalGroup(t *testing.T) {
	Convey("NewInternalGroup", t, func() {
		data := NewInternalGroup()
		assert.NotEqual(t, data, nil)
	})
}

func TestInternalGroupAdd(t *testing.T) {
	Convey("Add, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroup := &internalGroup{
			db:     db,
			logger: common.NewLogger(),
		}

		testErr := errors.New("xxx")
		Convey("exec fail", func() {
			mock.ExpectExec("").WillReturnError(testErr)
			httpErr := internalGroup.Add(strID)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := internalGroup.Add(strID)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestInternalGroupDelete(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroup := &internalGroup{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("ids len = 0", func() {
			httpErr := internalGroup.Delete(nil, nil)
			assert.Equal(t, httpErr, nil)
		})

		testErr := errors.New("xxx")
		Convey("exec fail", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(testErr)

			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := internalGroup.Delete([]string{strID}, tx)
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)

			httpErr := internalGroup.Delete([]string{strID}, tx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestInternalGroupGet(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		internalGroup := &internalGroup{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("ids len = 0", func() {
			_, httpErr := internalGroup.Get(nil)
			assert.Equal(t, httpErr, nil)
		})

		testErr := errors.New("xxx")
		Convey("query fail", func() {
			mock.ExpectQuery("").WillReturnError(testErr)

			_, httpErr := internalGroup.Get([]string{strID})
			assert.Equal(t, httpErr, testErr)
		})

		Convey("query success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_id"}).AddRow(strID))

			out, httpErr := internalGroup.Get([]string{strID})
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(out), 1)
			assert.Equal(t, out[strID].ID, strID)
		})
	})
}
