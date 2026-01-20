package dbaccess

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
)

func newAccessTokenPerm(ptrDB *sqlx.DB) *accessTokenPerm {
	return &accessTokenPerm{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

const appID = "b550af01-06d0-446d-be5b-b44cfcd97906"

func TestCheckAppAccessTokenPerm(t *testing.T) {
	Convey("TestCheckAppAccessTokenPerm", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		access := newAccessTokenPerm(db)
		dbErr := errors.New("db test error")

		Convey("execute error", func() {
			mock.ExpectQuery("").WithArgs(appID).WillReturnError(dbErr)
			_, err := access.CheckAppAccessTokenPerm(appID)
			assert.Equal(t, err, dbErr)
		})

		Convey("has no rows", func() {
			mock.ExpectQuery("").WithArgs(appID).WillReturnRows(sqlmock.NewRows([]string{"f_app_id"}))
			res, err := access.CheckAppAccessTokenPerm(appID)
			assert.Equal(t, err, nil)
			assert.Equal(t, res, false)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WithArgs(appID).WillReturnRows(sqlmock.NewRows([]string{"f_app_id"}).AddRow(appID))
			res, err := access.CheckAppAccessTokenPerm(appID)
			assert.Equal(t, err, nil)
			assert.Equal(t, res, true)
		})
	})
}

func TestAddAppAccessTokenPerm(t *testing.T) {
	Convey("TestAddAppAccessTokenPerm", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		access := newAccessTokenPerm(db)
		dbErr := errors.New("db test error")

		Convey("execute error", func() {
			mock.ExpectExec("").WithArgs(appID, sqlmock.AnyArg()).WillReturnError(dbErr)
			err := access.AddAppAccessTokenPerm(appID)
			assert.Equal(t, err, dbErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WithArgs(appID, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err := access.AddAppAccessTokenPerm(appID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAppAccessTokenPerm(t *testing.T) {
	Convey("TestDeleteAppAccessTokenPerm", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		access := newAccessTokenPerm(db)
		dbErr := errors.New("db test error")

		Convey("execute error", func() {
			mock.ExpectExec("").WithArgs(appID).WillReturnError(dbErr)
			err := access.DeleteAppAccessTokenPerm(appID)
			assert.Equal(t, err, dbErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WithArgs(appID).WillReturnResult(sqlmock.NewResult(1, 1))
			err := access.DeleteAppAccessTokenPerm(appID)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAllAppAccessTokenPerm(t *testing.T) {
	Convey("TestGetAllAppAccessTokenPerm", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		access := newAccessTokenPerm(db)
		dbErr := errors.New("db test error")

		Convey("execute error", func() {
			mock.ExpectQuery("").WillReturnError(dbErr)
			_, err := access.GetAllAppAccessTokenPerm()
			assert.Equal(t, err, dbErr)
		})

		Convey("has no rows", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_app_id"}))
			res, err := access.GetAllAppAccessTokenPerm()
			assert.Equal(t, err, nil)
			assert.Equal(t, len(res), 0)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_app_id"}).AddRow(appID))
			res, err := access.GetAllAppAccessTokenPerm()
			assert.Equal(t, err, nil)
			assert.Equal(t, len(res), 1)
			assert.Equal(t, res[0], appID)
		})
	})
}
