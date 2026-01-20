package dbaccess

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	mocks "UserManagement/interfaces/mock"
)

func newOrgPermAppDB(ptrDB *sqlx.DB) *orgPermApp {
	return &orgPermApp{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestNewOrgPermApp(t *testing.T) {
	Convey("NewOrgPermApp, db is available", t, func() {
		data := NewOrgPermApp()
		assert.NotEqual(t, data, nil)
	})
}

func TestGetAppPermByID2(t *testing.T) {
	Convey("GetAppPermByID2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)
		perm.dbTrace = db
		perm.trace = trace

		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnError(errors.New("unknown"))
			_, err := perm.GetAppPermByID2(ctx, strID)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			fields := []string{
				"f_app_name",
				"f_org_type",
				"f_perm_value",
				"f_end_time",
			}

			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("name1", 2, 3, -1))
			out, err := perm.GetAppPermByID2(ctx, strID)
			assert.Equal(t, err, nil)

			data, ok := out[interfaces.Department]
			assert.Equal(t, ok, true)
			assert.Equal(t, data.Subject, strID)
			assert.Equal(t, data.Object, interfaces.Department)
			assert.Equal(t, data.Name, "name1")
			assert.Equal(t, data.EndTime, int64(-1))
			assert.Equal(t, int32(data.Value), int32(3))
		})
	})
}

func TestGetAppPermByID(t *testing.T) {
	Convey("GetAppPermByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)

		Convey("execute error", func() {
			id := "bbb-bbb-bbb-bbb"
			mock.ExpectQuery("").WillReturnError(errors.New("unknown"))
			_, err := perm.GetAppPermByID(id)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			id := "xxx-xxx-xxx-xxx-xx1"
			fields := []string{
				"f_app_name",
				"f_org_type",
				"f_perm_value",
				"f_end_time",
			}

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("name1", 2, 3, -1))
			out, err := perm.GetAppPermByID(id)
			assert.Equal(t, err, nil)

			data, ok := out[interfaces.Department]
			assert.Equal(t, ok, true)
			assert.Equal(t, data.Subject, id)
			assert.Equal(t, data.Object, interfaces.Department)
			assert.Equal(t, data.Name, "name1")
			assert.Equal(t, data.EndTime, int64(-1))
			assert.Equal(t, int32(data.Value), int32(3))
		})
	})
}

func TestUpdateAppName(t *testing.T) {
	Convey("UpdateAppName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)

		info := interfaces.AppInfo{}
		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := perm.UpdateAppName(&info)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := perm.UpdateAppName(&info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateAppOrgPermInfo(t *testing.T) {
	Convey("UpdateAppOrgPerm, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)

		var inputInfo interfaces.AppOrgPerm
		inputInfo.Value = 3
		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.UpdateAppOrgPerm(inputInfo, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.UpdateAppOrgPerm(inputInfo, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddAppOrgPermInfo(t *testing.T) {
	Convey("AddAppOrgPerm, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)

		var inputInfo interfaces.AppOrgPerm
		inputInfo.Value = 1
		assert.Equal(t, err, nil)
		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.AddAppOrgPerm(inputInfo, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.AddAppOrgPerm(inputInfo, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAppOrgPermInfo(t *testing.T) {
	Convey("DeleteAppOrgPerm, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		perm := newOrgPermAppDB(db)
		strUserID := "strUserID1"
		types := make([]interfaces.OrgType, 0)

		Convey("types are empty", func() {
			err = perm.DeleteAppOrgPerm(strUserID, types)
			assert.Equal(t, err, nil)
		})

		types = append(types, interfaces.User)
		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			assert.Equal(t, err, nil)
			err = perm.DeleteAppOrgPerm(strUserID, types)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err = perm.DeleteAppOrgPerm(strUserID, types)
			assert.Equal(t, err, nil)
		})
	})
}
