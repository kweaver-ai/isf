package dbaccess

import (
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"Authentication/common"
)

func newLogin(ptrDB *sqlx.DB) *login {
	l := &login{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
	return l
}

func TestGetDomainStatus(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		login := newLogin(db)
		fields := []string{
			"f_status",
		}
		var enablePrefixMatch bool
		Convey("f_status = 1", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(0))
			enablePrefixMatch, err = login.GetDomainStatus()
			assert.Equal(t, enablePrefixMatch, true)
			assert.Equal(t, err, nil)
		})
		Convey("f_status = 0", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			enablePrefixMatch, err = login.GetDomainStatus()
			assert.Equal(t, enablePrefixMatch, false)
			assert.Equal(t, err, nil)
		})

		// 判断是否所有期望都被达到
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
