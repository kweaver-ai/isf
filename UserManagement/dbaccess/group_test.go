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

func newGroupDB(ptrDB *sqlx.DB) *group {
	return &group{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestNewGroup(t *testing.T) {
	Convey("NewGroup", t, func() {
		data := NewGroup()
		assert.NotEqual(t, data, nil)
	})
}

func TestGetGroupIDByName(t *testing.T) {
	Convey("GetGroupIDByName, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		userName := "user_name"
		Convey("no userName", func() {
			fields := []string{"f_group_name"}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			ret, httpErr := group.GetGroupIDByName(userName)
			assert.Equal(t, ret, "")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			fields := []string{"f_group_name"}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("group_name"))
			ret, httpErr := group.GetGroupIDByName(userName)
			assert.Equal(t, ret, "group_name")
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestGetGroupIDByName2(t *testing.T) {
	Convey("GetGroupIDByName2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)
		group.dbTrace = db
		group.trace = trace

		userName := "user_name"
		Convey("no userName", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			fields := []string{"f_group_name"}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			ret, httpErr := group.GetGroupIDByName2(ctx, userName)
			assert.Equal(t, ret, "")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			fields := []string{"f_group_name"}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("group_name"))
			ret, httpErr := group.GetGroupIDByName2(ctx, userName)
			assert.Equal(t, ret, "group_name")
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestAddGroup(t *testing.T) {
	Convey("AddGroup, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)
		group.dbTrace = db
		group.trace = trace

		Convey("execute error", func() {
			mock.ExpectBegin()
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			tx, _ := db.Begin()
			err := group.AddGroup(ctx, "id", "name", "notes", tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			tx, _ := db.Begin()
			err := group.AddGroup(ctx, "id", "name", "notes", tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteGroup(t *testing.T) {
	Convey("DeleteGroup, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := group.DeleteGroup("id")
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := group.DeleteGroup("id")
			assert.Equal(t, err, nil)
		})
	})
}

func TestModifyGroup(t *testing.T) {
	Convey("DeleteGroup, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := group.ModifyGroup("id", "xxx", true, "xxx", true)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := group.ModifyGroup("id", "xxx", false, "xxx", false)
			assert.Equal(t, err, nil)
		})

		Convey("success1", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := group.ModifyGroup("id", "xxx", true, "xxx", false)
			assert.Equal(t, err, nil)
		})

		Convey("success2", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := group.ModifyGroup("id", "xxx", false, "xxx", true)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetGroups(t *testing.T) {
	Convey("GetGroups, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		info := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    10,
			Limit:     20,
			Keyword:   "xxx",
		}

		fields := []string{
			"f_group_id",
			"f_group_name",
			"f_notes",
		}

		Convey("no group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groups, httpErr := group.GetGroups(info)
			assert.Equal(t, len(groups), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			tempInfo := interfaces.GroupInfo{
				ID:    "xxx",
				Name:  "yyyy",
				Notes: "zzz",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.ID, tempInfo.Name, tempInfo.Notes))
			outInfos, httpErr := group.GetGroups(info)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, outInfos[0], tempInfo)
		})
	})
}

func TestGetGroupsNum(t *testing.T) {
	Convey("GetGroupsNum, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		info := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    10,
			Limit:     20,
			Keyword:   "xxx",
		}

		fields := []string{
			"count",
		}

		Convey("no group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			num, httpErr := group.GetGroupsNum(info)
			assert.Equal(t, num, 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(1))
			num, httpErr := group.GetGroupsNum(info)
			assert.Equal(t, num, 1)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success1", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(2))
			num, httpErr := group.GetGroupsNum(info)
			assert.Equal(t, num, 2)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestGetGroupByID(t *testing.T) {
	Convey("GetGroupByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		fields := []string{
			"f_group_id",
			"f_group_name",
			"f_notes",
		}

		Convey("no group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			_, httpErr := group.GetGroupByID("")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			tempInfo := interfaces.GroupInfo{
				ID:    "xxx",
				Name:  "yyyy",
				Notes: "zzz",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.ID, tempInfo.Name, tempInfo.Notes))
			outInfos, httpErr := group.GetGroupByID("")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, outInfos, tempInfo)
		})
	})
}

func TestGetGroupByID2(t *testing.T) {
	Convey("GetGroupByID2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)
		group.dbTrace = db
		group.trace = trace

		fields := []string{
			"f_group_id",
			"f_group_name",
			"f_notes",
		}

		Convey("no group", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			_, httpErr := group.GetGroupByID2(ctx, "")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			tempInfo := interfaces.GroupInfo{
				ID:    "xxx",
				Name:  "yyyy",
				Notes: "zzz",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.ID, tempInfo.Name, tempInfo.Notes))
			outInfos, httpErr := group.GetGroupByID2(ctx, "")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, outInfos, tempInfo)
		})
	})
}

func TestGetExistGroupIDs(t *testing.T) {
	Convey("GetExistGroupIDs, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		fields := []string{
			"f_group_id",
		}

		groupIDs := make([]string, 0)
		Convey("groupIDs is empty", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			existGroupIDs, err := group.GetExistGroupIDs(groupIDs)
			assert.Equal(t, existGroupIDs, nil)
			assert.Equal(t, err, nil)
		})

		groupIDs = append(groupIDs, "group_id")
		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("group_id"))
			existGroupIDs, err := group.GetExistGroupIDs(groupIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, existGroupIDs, []string{"group_id"})
		})
	})
}

func TestSearchGroupByKey(t *testing.T) {
	Convey("SearchGroupByKeyword, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		fields := []string{
			"f_group_id",
			"f_group_name",
		}

		Convey("no group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			_, httpErr := group.SearchGroupByKeyword("", 1, 1)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			tempInfo := interfaces.NameInfo{
				ID:   "xxx",
				Name: "yyyy",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(tempInfo.ID, tempInfo.Name))
			outInfos, httpErr := group.SearchGroupByKeyword("", 1, 1)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, outInfos[0], tempInfo)
		})
	})
}

func TestSearchGroupNumByKeyword(t *testing.T) {
	Convey("SearchGroupNumByKeyword, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)

		fields := []string{
			"f_group_id",
		}

		Convey("no group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			_, httpErr := group.SearchGroupNumByKeyword("")
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(2))
			outInfos, httpErr := group.SearchGroupNumByKeyword("")
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, outInfos, 2)
		})
	})
}

func TestConvertGroupName(t *testing.T) {
	Convey("GetGroupName, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)
		group.dbTrace = db
		group.trace = trace

		fields := []string{
			"f_group_id",
			"f_group_name",
		}

		Convey("no groupID", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			outInfo1, outInfo2, httpErr := group.GetGroupName(make([]string, 0))
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("no group", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			outInfo1, outInfo2, httpErr := group.GetGroupName([]string{"xxxx"})
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("zzz", "xxxxx").AddRow("kkk", "yyyyy"))
			outInfos, outInfo2, httpErr := group.GetGroupName([]string{"xxxx"})
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

func TestGetGroupName2(t *testing.T) {
	Convey("GetGroupName2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		group := newGroupDB(db)
		group.dbTrace = db
		group.trace = trace

		fields := []string{
			"f_group_id",
			"f_group_name",
		}

		Convey("no groupID", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			outInfo1, outInfo2, httpErr := group.GetGroupName2(ctx, make([]string, 0))
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("no group", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			outInfo1, outInfo2, httpErr := group.GetGroupName2(ctx, []string{"xxxx"})
			assert.Equal(t, len(outInfo1), 0)
			assert.Equal(t, len(outInfo2), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("zzz", "xxxxx").AddRow("kkk", "yyyyy"))
			outInfos, outInfo2, httpErr := group.GetGroupName2(ctx, []string{"xxxx"})
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
