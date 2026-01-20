package dbaccess

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	mocks "UserManagement/interfaces/mock"
)

func newContactorDB(ptrDB *sqlx.DB) *contactor {
	return &contactor{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestNewContactor(t *testing.T) {
	Convey("NewContactor", t, func() {
		data := NewContactor()
		assert.NotEqual(t, data, nil)
	})
}

func TestGetContactorName(t *testing.T) {
	Convey("GetContactorName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)
		contactorIDs := make([]string, 0)

		Convey("contacor not exist", func() {
			fields := []string{
				"f_group_id",
				"f_group_name",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			infoMap, name, httpErr := contactor.GetContactorName(contactorIDs)
			assert.Equal(t, len(name), 0)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(infoMap), 0)
		})

		Convey("success", func() {
			contactorIDs = append(contactorIDs, "test")

			fields := []string{
				"f_group_id",
				"f_group_name",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("test", "f_name"))
			infoMap, name, httpErr := contactor.GetContactorName(contactorIDs)
			assert.Equal(t, httpErr, nil)
			tempNameInfo := interfaces.NameInfo{
				ID:   "test",
				Name: "f_name",
			}
			assert.Equal(t, infoMap, []interfaces.NameInfo{tempNameInfo})
			assert.Equal(t, name, []string{"test"})
			assert.Equal(t, 0, 0)
		})
	})
}

func TestGetUserAllBelongGroupIDs(t *testing.T) {
	Convey("GetUserAllBelongContactorIDs, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)

		Convey("contactor not exist", func() {
			fields := []string{
				"f_group_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			infoMap, httpErr := contactor.GetUserAllBelongContactorIDs("xxxx")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(infoMap), 0)
		})

		Convey("success", func() {
			fields := []string{
				"f_group_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("test"))
			infoMap, httpErr := contactor.GetUserAllBelongContactorIDs("xxxxx")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, infoMap, []string{"test"})
		})
	})
}

func TestGetContactorInfo(t *testing.T) {
	Convey("GetContactorInfo, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)

		Convey("contactor not exist", func() {
			fields := []string{
				"f_user_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			ret, _, httpErr := contactor.GetContactorInfo("xxxx")
			assert.Equal(t, ret, false)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			fields := []string{
				"f_user_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("test"))
			ret, infoMap, httpErr := contactor.GetContactorInfo("xxxxx")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, ret, true)
			assert.Equal(t, infoMap.UserID, "test")
		})
	})
}

func TestDeleteContactors(t *testing.T) {
	Convey("DeleteContactorMembers, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)
		Convey("empty contactors", func() {
			mock.ExpectBegin()
			tx, _ := db.Begin()
			err := contactor.DeleteContactorMembers(nil, tx)
			assert.Equal(t, err, nil)
		})

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteContactorMembers([]string{"xxx"}, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(nil)
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteContactorMembers([]string{"xxx"}, tx)
			assert.NotEqual(t, err, nil)
		})
	})
}

func TestDeleteContactorMembers(t *testing.T) {
	Convey("DeleteContactors, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)
		Convey("empty contactors", func() {
			mock.ExpectBegin()
			tx, _ := db.Begin()
			err := contactor.DeleteContactors(nil, tx)
			assert.Equal(t, err, nil)
			assert.Equal(t, 1, 1)
		})

		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteContactors([]string{"xxx"}, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(nil)
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteContactors([]string{"xxx"}, tx)
			assert.NotEqual(t, err, nil)
		})
	})
}

func TestGetAllContactorInfos(t *testing.T) {
	Convey("GetAllContactorInfos, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)

		Convey("Query error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			mock.ExpectQuery("").WillReturnError(testErr)
			_, httpErr := contactor.GetAllContactorInfos("xxxx")
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			fields := []string{
				"f_group_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("test"))
			info, httpErr := contactor.GetAllContactorInfos("xxxxx")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(info), 1)
			assert.Equal(t, info[0].ContactorID, "test")
		})
	})
}

func TestDeleteUserInContactors(t *testing.T) {
	Convey("DeleteUserInContactors, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)
		Convey("execute error", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteUserInContactors("", tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(nil)
			mock.ExpectCommit()
			tx, _ := db.Begin()
			err := contactor.DeleteUserInContactors("", tx)
			assert.NotEqual(t, err, nil)
		})
	})
}

func TestUpdateContactorCount(t *testing.T) {
	Convey("UpdateContactorCount, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		contactor := newContactorDB(db)
		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := contactor.UpdateContactorCount()
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnError(nil)
			err := contactor.UpdateContactorCount()
			assert.NotEqual(t, err, nil)
		})
	})
}

// TestGetContactorMemberIDs 测试批量获取联系人组成员ID
func TestGetContactorMemberIDs(t *testing.T) {
	Convey("GetContactorMemberIDs, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		contactor := newContactorDB(db)
		contactor.trace = trace
		contactor.dbTrace = db

		Convey("Query error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			testErr := rest.NewHTTPError("error", 503000000, nil)
			mock.ExpectQuery("").WillReturnError(testErr)
			_, httpErr := contactor.GetContactorMemberIDs(ctx, []string{"xxxx"})
			assert.Equal(t, httpErr, testErr)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			fields := []string{
				"f_group_id",
				"f_user_id",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("test", "test1").AddRow("test", "test2"))
			info, httpErr := contactor.GetContactorMemberIDs(ctx, []string{"test"})
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(info), 1)
			assert.Equal(t, info["test"], []string{"test1", "test2"})
		})
	})
}
