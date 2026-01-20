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

func newOrgPermDB(ptrDB *sqlx.DB) *orgPerm {
	return &orgPerm{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestNewOrgPerm(t *testing.T) {
	Convey("NewOrgPerm, db is available", t, func() {
		data := NewOrgPerm()
		assert.NotEqual(t, data, nil)
	})
}

func TestUpdateName(t *testing.T) {
	Convey("UpdateAppName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		perm := newOrgPermDB(db)
		perm.trace = trace

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := perm.UpdateName("", "")
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := perm.UpdateName("", "")
			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateOrgPermInfo(t *testing.T) {
	Convey("UpdateOrgPermInfo, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		perm := newOrgPermDB(db)
		perm.trace = trace

		ctx := context.Background()

		var inputInfo interfaces.OrgPerm
		inputInfo.Value = 3
		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.UpdateOrgPerm(ctx, inputInfo, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.UpdateOrgPerm(ctx, inputInfo, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddOrgPermInfo(t *testing.T) {
	Convey("AddOrgPermInfo, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		perm := newOrgPermDB(db)
		perm.trace = trace

		ctx := context.Background()

		var inputInfo interfaces.OrgPerm
		inputInfo.Value = 1
		assert.Equal(t, err, nil)
		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.AddOrgPerm(ctx, inputInfo, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = perm.AddOrgPerm(ctx, inputInfo, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteOrgPermInfo(t *testing.T) {
	Convey("DeleteOrgPermInfo, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		perm := newOrgPermDB(db)
		perm.trace = trace

		ctx := context.Background()
		strUserID := "strUserID1"
		types := make([]interfaces.OrgType, 0)

		Convey("types are empty", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			err = perm.DeleteOrgPerm(ctx, strUserID, types)
			assert.Equal(t, err, nil)
		})

		types = append(types, interfaces.User)
		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			assert.Equal(t, err, nil)
			err = perm.DeleteOrgPerm(ctx, strUserID, types)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err = perm.DeleteOrgPerm(ctx, strUserID, types)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetPermByID(t *testing.T) {
	Convey("GetPermByID, db is available", t, func() {
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
			id := "bbb-bbb-bbb-bbb"
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnError(errors.New("unknown"))
			_, err := perm.GetAppPermByID2(ctx, id)
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

			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("name1", 2, 3, -1))
			out, err := perm.GetAppPermByID2(ctx, id)
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

func TestDeleteOrgPermByID(t *testing.T) {
	Convey("DeleteOrgPermByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		perm := newOrgPermDB(db)
		perm.trace = trace

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := perm.DeleteOrgPermByID(strID)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := perm.DeleteOrgPermByID(strID)
			assert.Equal(t, err, nil)
		})
	})
}
