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

func newGroupMemberDB(ptrDB *sqlx.DB) *groupMember {
	return &groupMember{
		db:     ptrDB,
		logger: common.NewLogger(),
	}
}

func TestDeleteGroupMemberByID(t *testing.T) {
	Convey("DeleteGroupMemberByID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := groupMember.DeleteGroupMemberByID("id")
			assert.NotEqual(t, err, nil)
		})

		Convey("execute error1", func() {
			mock.ExpectExec("").WillReturnError(errors.New("xxxx"))
			err := groupMember.DeleteGroupMemberByID("id")
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := groupMember.DeleteGroupMemberByID("id")
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddGroupMember(t *testing.T) {
	Convey("AddGroupMember, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		info := &interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      1,
			Name:            "xxxxx",
			DepartmentNames: []string{"xxxx"},
		}

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := groupMember.AddGroupMember("id", info)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := groupMember.AddGroupMember("id", info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddGroupMembers(t *testing.T) {
	Convey("AddGroupMembers, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)
		groupMember.dbTrace = db
		groupMember.trace = trace

		info := interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      1,
			Name:            "xxxxx",
			DepartmentNames: []string{"xxxx"},
		}

		Convey("execute error", func() {
			mock.ExpectBegin()
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = groupMember.AddGroupMembers(ctx, "id", []interfaces.GroupMemberInfo{info}, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectBegin()
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			tx, err := db.Begin()
			assert.Equal(t, err, nil)
			err = groupMember.AddGroupMembers(ctx, "id", []interfaces.GroupMemberInfo{info}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteGroupMember(t *testing.T) {
	Convey("DeleteGroupMember, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		info := &interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      1,
			Name:            "xxxxx",
			DepartmentNames: []string{"xxxx"},
		}

		Convey("execute error", func() {
			mock.ExpectExec("").WillReturnError(errors.New("unknown"))
			err := groupMember.DeleteGroupMember("id", info)
			assert.NotEqual(t, err, nil)
		})

		Convey("execute error1", func() {
			mock.ExpectExec("").WillReturnError(errors.New("xxxxx"))
			err := groupMember.DeleteGroupMember("id", info)
			assert.NotEqual(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := groupMember.DeleteGroupMember("id", info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetGroupMembers(t *testing.T) {
	Convey("GetGroupMembers, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)
		groupMember.dbTrace = db
		groupMember.trace = trace

		info := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    50,
			Limit:     50,
			Keyword:   "xxx",
		}
		groupID := strID

		fields := []string{
			"c.f_name",
			"b.f_display_name",
			"b.f_member_id",
			"b.f_MemberType",
			"a.f_added_time",
		}

		Convey("no group member", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			members, httpErr := groupMember.GetGroupMembers(ctx, groupID, info)
			assert.Equal(t, len(members), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			tempInfo := interfaces.GroupMemberInfo{
				ID:              "xxxx",
				MemberType:      1,
				Name:            "kkkk",
				DepartmentNames: nil,
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("", tempInfo.Name, tempInfo.ID, 1, 1))
			outInfos, httpErr := groupMember.GetGroupMembers(ctx, groupID, info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, len(outInfos), 1)
			assert.Equal(t, outInfos[0], tempInfo)
		})
	})
}

func TestGetGroupMembersNum2(t *testing.T) {
	Convey("GetGroupMembersNum2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)
		groupMember.dbTrace = db
		groupMember.trace = trace

		info := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    50,
			Limit:     50,
			Keyword:   "xxx",
		}
		groupID := "xxxxx"

		fields := []string{
			"count",
		}

		Convey("no group member", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			num, httpErr := groupMember.GetGroupMembersNum2(ctx, groupID, info)
			assert.Equal(t, num, 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(1))
			num, httpErr := groupMember.GetGroupMembersNum2(ctx, groupID, info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, num, 1)
		})

		Convey("success1", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(2))
			num, httpErr := groupMember.GetGroupMembersNum2(ctx, groupID, info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, num, 2)
		})
	})
}

func TestGetGroupMembersNum(t *testing.T) {
	Convey("GetGroupMembersNum, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		info := interfaces.SearchInfo{
			Direction: interfaces.Asc,
			Sort:      interfaces.DateCreated,
			Offset:    50,
			Limit:     50,
			Keyword:   "xxx",
		}
		groupID := "xxxxx"

		fields := []string{
			"count",
		}

		Convey("no group member", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			num, httpErr := groupMember.GetGroupMembersNum(groupID, info)
			assert.Equal(t, num, 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(1))
			num, httpErr := groupMember.GetGroupMembersNum(groupID, info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, num, 1)
		})

		Convey("success1", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(2))
			num, httpErr := groupMember.GetGroupMembersNum(groupID, info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, num, 2)
		})
	})
}

func TestCheckGroupMembersExist(t *testing.T) {
	Convey("CheckGroupMembersExist, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"f_group_id",
		}

		info := interfaces.GroupMemberInfo{
			ID:              "xxxx",
			MemberType:      1,
			Name:            "xxxxx",
			DepartmentNames: []string{"xxxx"},
		}

		Convey("no group member", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			ret, httpErr := groupMember.CheckGroupMembersExist("xxxxx", &info)
			assert.Equal(t, ret, false)
			assert.Equal(t, httpErr, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(info.ID))
			ret, httpErr := groupMember.CheckGroupMembersExist("xxxx", &info)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, ret, true)
		})
	})
}

func TestGetMembersBelongGroupIDs(t *testing.T) {
	Convey("GetMembersBelongGroupIDs, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"f_group_id",
			"f_group_name",
		}

		info := interfaces.GroupInfo{
			ID:   "xxxx",
			Name: "xxxxx",
		}

		Convey("no members belong group", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groupIds, groups, Err := groupMember.GetMembersBelongGroupIDs([]string{"xxxxx"})
			assert.Equal(t, len(groupIds), 0)
			assert.Equal(t, len(groups), 0)
			assert.Equal(t, Err, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(info.ID, info.Name))
			groupIds, groups, Err := groupMember.GetMembersBelongGroupIDs([]string{"xxxxx"})
			assert.Equal(t, len(groupIds), 1)
			assert.Equal(t, groups, []interfaces.GroupInfo{info})
			assert.Equal(t, Err, nil)
		})
	})
}

func TestGetGroupMembersByGroupIDs(t *testing.T) {
	Convey("GetGroupMembersByGroupIDs, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"f_member_id",
			"f_member_type",
		}

		groupIDs := make([]string, 0)
		Convey("groupIDs is empty", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groupMemberInfos, err := groupMember.GetGroupMembersByGroupIDs(groupIDs)
			assert.Equal(t, groupMemberInfos, nil)
			assert.Equal(t, err, nil)
		})

		groupIDs = append(groupIDs, "group_id")
		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", 1))
			groupMemberInfos, err := groupMember.GetGroupMembersByGroupIDs(groupIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(groupMemberInfos), 1)
			assert.Equal(t, groupMemberInfos[0].ID, "user_id")
			assert.Equal(t, groupMemberInfos[0].MemberType, 1)
		})
	})
}

func TestGetGroupMembersByGroupIDs2(t *testing.T) {
	Convey("GetGroupMembersByGroupIDs2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)
		groupMember.dbTrace = db
		groupMember.trace = trace

		fields := []string{
			"f_group_id",
			"f_member_id",
			"f_member_type",
		}

		groupIDs := make([]string, 0)
		Convey("groupIDs is empty", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groupMemberInfos, err := groupMember.GetGroupMembersByGroupIDs2(ctx, groupIDs)
			assert.Equal(t, groupMemberInfos, nil)
			assert.Equal(t, err, nil)
		})

		groupIDs = append(groupIDs, "group_id")
		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("group_id", "user_id", 1))
			groupMemberInfos, err := groupMember.GetGroupMembersByGroupIDs2(ctx, groupIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(groupMemberInfos), 1)
			assert.Equal(t, len(groupMemberInfos["group_id"]), 1)
			assert.Equal(t, groupMemberInfos["group_id"][0].ID, "user_id")
			assert.Equal(t, groupMemberInfos["group_id"][0].MemberType, 1)
		})
	})
}

func TestSearchMembersByKeyword(t *testing.T) {
	Convey("SearchMembersByKeyword, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"f_member_id",
			"f_member_type",
			"f_group_name",
			"f_name1",
			"f_name2",
		}

		Convey("groupIDs is empty", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groupMemberInfos, err := groupMember.SearchMembersByKeyword("", 0, 0)
			assert.Equal(t, len(groupMemberInfos), 0)
			assert.Equal(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", 1, "zzzz", "ss", nil).AddRow("user_id", 1, "kkk", "ss", "").AddRow("xxx", 2, "zzz", nil, "sasd"))
			groupMemberInfos, err := groupMember.SearchMembersByKeyword("", 0, 0)

			test1Info := interfaces.MemberInfo{
				ID:    "user_id",
				Name:  "ss",
				NType: 1,
				GroupNames: []string{
					"zzzz",
					"kkk",
				},
			}

			test2Info := interfaces.MemberInfo{
				ID:         "xxx",
				Name:       "sasd",
				NType:      2,
				GroupNames: []string{"zzz"},
			}

			assert.Equal(t, err, nil)
			assert.Equal(t, len(groupMemberInfos), 2)
			assert.Equal(t, groupMemberInfos[0], test1Info)
			assert.Equal(t, groupMemberInfos[1], test2Info)
		})
	})
}

func TestSearchMemberNumByKeyword(t *testing.T) {
	Convey("SearchMemberNumByKeyword, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"num",
		}

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow(2))
			groupMemberInfos, err := groupMember.SearchMemberNumByKeyword("")

			assert.Equal(t, err, nil)
			assert.Equal(t, groupMemberInfos, 2)
		})
	})
}

func TestGetMemberOnClient(t *testing.T) {
	Convey("GetMemberOnClient, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		fields := []string{
			"f_member_id",
			"f_member_type",
			"f_name",
			"f_name1",
			"f_priority",
		}

		Convey("groupIDs is empty", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			groupMemberInfos, err := groupMember.GetMemberOnClient("", 0, 0)
			assert.Equal(t, len(groupMemberInfos), 0)
			assert.Equal(t, err, nil)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", 1, "zzzz", nil, 22).AddRow("xxx", 2, "zzz", nil, 33))
			groupMemberInfos, err := groupMember.GetMemberOnClient("", 0, 0)

			test1Info := interfaces.MemberSimpleInfo{
				ID:    "user_id",
				Name:  "zzzz",
				NType: 1,
			}

			test2Info := interfaces.MemberSimpleInfo{
				ID:    "xxx",
				Name:  "zzz",
				NType: 2,
			}

			assert.Equal(t, err, nil)
			assert.Equal(t, len(groupMemberInfos), 2)
			assert.Equal(t, groupMemberInfos[0], test1Info)
			assert.Equal(t, groupMemberInfos[1], test2Info)
		})
	})
}

func TestDeleteGroupMemberByMemberID(t *testing.T) {
	Convey("DeleteGroupMemberByMemberID, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)
		groupMember := newGroupMemberDB(db)

		Convey("DeleteGroupMemberByMemberID error", func() {
			mock.ExpectExec("").WillReturnError(errors.New(""))

			err = groupMember.DeleteGroupMemberByMemberID(userID)
			assert.NotEqual(t, err, nil)
		})

		Convey("DeleteGroupMemberByMemberID successfully", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err = groupMember.DeleteGroupMemberByMemberID(userID)
			assert.Equal(t, err, nil)
		})
	})
}
