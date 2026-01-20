package dbaccess

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authorization/common"
	"Authorization/interfaces"
)

//nolint:dupl
func TestDBDeleteRoleMemberByID(t *testing.T) {
	Convey("TestDBDeleteRoleMemberByID", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := b.DeleteRoleMemberByID(ctx, "test-role-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.DeleteRoleMemberByID(ctx, "test-role-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBAddRoleMembers(t *testing.T) {
	Convey("TestDBAddRoleMembers", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		roleMembers := []interfaces.RoleMemberInfo{
			{
				ID:         "test-member-1",
				MemberType: interfaces.AccessorUser,
				Name:       "test-member-1-name",
			},
			{
				ID:         "test-member-2",
				MemberType: interfaces.AccessorDepartment,
				Name:       "test-member-2-name",
			},
		}

		Convey("empty members", func() {
			err := b.AddRoleMembers(ctx, "test-role-id", []interfaces.RoleMemberInfo{})
			assert.Equal(t, err, nil)
		})

		Convey("insert error", func() {
			mock.ExpectExec("^insert").WillReturnError(mockErr)
			err := b.AddRoleMembers(ctx, "test-role-id", roleMembers)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^insert").WillReturnResult(sqlmock.NewResult(2, 2))
			err := b.AddRoleMembers(ctx, "test-role-id", roleMembers)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBDeleteRoleMembers(t *testing.T) {
	Convey("TestDBDeleteRoleMembers", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		memberIDs := []string{"test-member-1", "test-member-2"}

		Convey("empty member ids", func() {
			err := b.DeleteRoleMembers(ctx, "test-role-id", []string{})
			assert.Equal(t, err, nil)
		})

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := b.DeleteRoleMembers(ctx, "test-role-id", memberIDs)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 2))
			err := b.DeleteRoleMembers(ctx, "test-role-id", memberIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleMembersNum(t *testing.T) {
	Convey("TestDBGetRoleMembersNum", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		searchInfo := interfaces.RoleMemberSearchInfo{
			Offset: 0,
			Limit:  10,
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			num, err := b.GetRoleMembersNum(ctx, "test-role-id", searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, num, 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"count(f_member_id)"}).AddRow("invalid")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			num, err := b.GetRoleMembersNum(ctx, "test-role-id", searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, num, 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"count(f_member_id)"}).AddRow(5)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			num, err := b.GetRoleMembersNum(ctx, "test-role-id", searchInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, num, 5)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetPaginationByRoleID(t *testing.T) {
	Convey("TestDBGetPaginationByRoleID", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		searchInfo := interfaces.RoleMemberSearchInfo{
			Offset: 0,
			Limit:  10,
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			members, err := b.GetPaginationByRoleID(ctx, "test-role-id", searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, members, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_member_id", "f_member_type", "f_member_name"}).
				AddRow("test-member-id", "invalid-type", "test-member-name")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			members, err := b.GetPaginationByRoleID(ctx, "test-role-id", searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, members, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"f_member_id", "f_member_type", "f_member_name"}).
				AddRow("test-member-1", interfaces.AccessorUser, "test-member-1-name").
				AddRow("test-member-2", interfaces.AccessorDepartment, "test-member-2-name")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			members, err := b.GetPaginationByRoleID(ctx, "test-role-id", searchInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(members), 2)
			assert.Equal(t, members[0].ID, "test-member-1")
			assert.Equal(t, members[0].MemberType, interfaces.AccessorUser)
			assert.Equal(t, members[0].Name, "test-member-1-name")
			assert.Equal(t, members[1].ID, "test-member-2")
			assert.Equal(t, members[1].MemberType, interfaces.AccessorDepartment)
			assert.Equal(t, members[1].Name, "test-member-2-name")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleMembersByRoleID(t *testing.T) {
	Convey("TestDBGetRoleMembersByRoleID", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			members, err := b.GetRoleMembersByRoleID(ctx, "test-role-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, members, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_member_id", "f_member_type", "f_member_name"}).
				AddRow("test-member-id", "invalid-type", "test-member-name")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			members, err := b.GetRoleMembersByRoleID(ctx, "test-role-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, members, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"f_member_id", "f_member_type", "f_member_name"}).
				AddRow("test-member-1", interfaces.AccessorUser, "test-member-1-name").
				AddRow("test-member-2", interfaces.AccessorDepartment, "test-member-2-name")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			members, err := b.GetRoleMembersByRoleID(ctx, "test-role-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(members), 2)
			assert.Equal(t, members[0].ID, "test-member-1")
			assert.Equal(t, members[0].MemberType, interfaces.AccessorUser)
			assert.Equal(t, members[0].Name, "test-member-1-name")
			assert.Equal(t, members[1].ID, "test-member-2")
			assert.Equal(t, members[1].MemberType, interfaces.AccessorDepartment)
			assert.Equal(t, members[1].Name, "test-member-2-name")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleByMembers(t *testing.T) {
	Convey("TestDBGetRoleByMembers", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		memberIDs := []string{"test-member-1", "test-member-2"}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			roles, err := b.GetRoleByMembers(ctx, memberIDs)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, roles, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan Success has one role", func() {
			rows := sqlmock.NewRows([]string{"f_role_id"}).AddRow("test-role-id")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			roles, err := b.GetRoleByMembers(ctx, memberIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 1)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"f_role_id"}).
				AddRow("test-role-1").
				AddRow("test-role-2")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			roles, err := b.GetRoleByMembers(ctx, memberIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 2)
			assert.Equal(t, roles[0].ID, "test-role-1")
			assert.Equal(t, roles[1].ID, "test-role-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

//nolint:dupl
func TestDBDeleteByMemberIDs(t *testing.T) {
	Convey("TestDBDeleteByMemberIDs", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		memberIDs := []string{"test-member-1", "test-member-2"}

		Convey("empty member ids", func() {
			err := b.DeleteByMemberIDs([]string{})
			assert.Equal(t, err, nil)
		})

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := b.DeleteByMemberIDs(memberIDs)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 2))
			err := b.DeleteByMemberIDs(memberIDs)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

//nolint:dupl
func TestDBDeleteByRoleID(t *testing.T) {
	Convey("TestDBDeleteByRoleID", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")
		ctx := context.Background()

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := b.DeleteByRoleID(ctx, "test-role-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.DeleteByRoleID(ctx, "test-role-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBUpdateMemberName(t *testing.T) {
	Convey("TestDBUpdateMemberName", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				return
			}
		}(db)

		mockErr := errors.New("test error")

		b := &roleMember{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("update error", func() {
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := b.UpdateMemberName("test-member-id", "new-member-name")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.UpdateMemberName("test-member-id", "new-member-name")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}
