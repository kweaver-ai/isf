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

func newAppDB(ptrDB *sqlx.DB) *app {
	return &app{
		db:  ptrDB,
		log: common.NewLogger(),
	}
}

func TestNewApp(t *testing.T) {
	Convey("NewApp, db is available", t, func() {
		data := NewApp()
		assert.NotEqual(t, data, nil)
	})
}

func TestRegisterApp(t *testing.T) {
	Convey("Register, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		info := &interfaces.AppCompleteInfo{
			AppInfo: interfaces.AppInfo{
				ID:             "xxx-xxx-xxx-xxx",
				Name:           "test",
				CredentialType: interfaces.CredentialTypePassword,
			},
			Password: "some-secret",
			Type:     interfaces.General,
		}

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.RegisterApp(info, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.RegisterApp(info, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteApp(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		id := "xxxxx-xxxxx-xxxxxx-xxxxxx-xxxxxx"

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.DeleteApp(id, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.DeleteApp(id, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateApp(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		id := "xxxxx-xxxxx-xxxxxx-xxxxxx-xxxxxx"
		name := "test"
		pwd := "some-secret"

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.UpdateApp(id, true, name, true, pwd, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := app.UpdateApp(id, true, name, true, pwd, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAppList(t *testing.T) {
	Convey("List, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		searchInfo := &interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    0,
			Limit:     20,
		}

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			_, err := app.AppList(searchInfo)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			fields := []string{
				"f_name",
				"f_id",
				"f_credential_type",
			}

			info := interfaces.AppInfo{
				ID:             "xxx-xxx-xxx-xxx-xx1",
				Name:           "test1",
				CredentialType: interfaces.CredentialTypePassword,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("xxx-xxx-xxx-xxx-xx1", "test1", 1))
			tmpInfo, err := app.AppList(searchInfo)
			assert.Equal(t, tmpInfo, &[]interfaces.AppInfo{info})
			assert.Equal(t, err, nil)
		})
	})
}

func TestAppListCount(t *testing.T) {
	Convey("ListCount, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		info := &interfaces.SearchInfo{}
		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			_, err := app.AppListCount(info)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			fields := []string{
				"count(*)",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("27"))
			tmp, err := app.AppListCount(info)
			assert.Equal(t, tmp, 27)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAppByName(t *testing.T) {
	Convey("GetAppByName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		name := "test"

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			_, err := app.GetAppByName(name)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			fields := []string{
				"f_id",
				"f_name",
			}

			info := interfaces.AppInfo{
				ID:   "aaa-aaa-aaa-aaa",
				Name: "test1",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("aaa-aaa-aaa-aaa", "test1"))
			tmpInfo, err := app.GetAppByName(name)
			assert.Equal(t, tmpInfo, &info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAppByID(t *testing.T) {
	Convey("GetAppByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		app := newAppDB(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			_, err := app.GetAppByID(strID)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			id := strID
			fields := []string{
				"f_id",
				"f_name",
				"f_credential_type",
			}

			info := interfaces.AppInfo{
				ID:             strID,
				Name:           "test1",
				CredentialType: interfaces.CredentialTypePassword,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(strID, "test1", 1))
			tmpInfo, err := app.GetAppByID(id)
			assert.Equal(t, tmpInfo, &info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestConvertAppName(t *testing.T) {
	Convey("ConvertAppName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		app := newAppDB(db)

		Convey("no appID", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_id", "f_name"}))
			outInfo1, outInfo2, httpErr := app.GetAppName(make([]string, 0))
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("no app", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_id", "f_name"}))
			outInfo1, outInfo2, httpErr := app.GetAppName([]string{"xxxx"})
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_id", "f_name"}).AddRow("zzz", "xxxxx").AddRow("kkk", "yyyyy"))
			outInfos, outInfo2, httpErr := app.GetAppName([]string{"xxxx"})
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(outInfos), 2)
			assert.Equal(t, len(outInfos), 2)
			assert.Equal(t, outInfos[0], interfaces.NameInfo{ID: "zzz", Name: "xxxxx"})
			assert.Equal(t, outInfos[1], interfaces.NameInfo{ID: "kkk", Name: "yyyyy"})
			assert.Equal(t, outInfo2[0], "zzz")
			assert.Equal(t, outInfo2[1], "kkk")
		})
	})
}
