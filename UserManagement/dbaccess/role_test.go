// Package dbaccess user Anyshare 数据访问层 - 用户数据库操作
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

func TestGetRolesByUserIDs(t *testing.T) {
	Convey("GetRolesByUserIDs, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)

		userID := "test1"

		role := &role{
			db:     db,
			logger: common.NewLogger(),
		}

		fields := []string{
			"f_user_id",
			"f_role_id",
		}
		Convey("roles id not exist", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			roleIDs, httpErr := role.GetRolesByUserIDs([]string{})
			assert.Equal(t, len(roleIDs), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("Success1", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", "role_id"))
			roleIDs, httpErr := role.GetRolesByUserIDs([]string{userID})
			assert.Equal(t, len(roleIDs), 1)
			assert.Equal(t, httpErr, nil)
			roleInfo, ok := roleIDs["user_id"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo["role_id"], true)
		})

		Convey("Success2", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", "role_id").AddRow("user_id1", "role_id1").AddRow("user_id", "role_id1"))
			roleIDs, httpErr := role.GetRolesByUserIDs([]string{userID})
			assert.Equal(t, len(roleIDs), 2)
			assert.Equal(t, httpErr, nil)
			roleInfo, ok := roleIDs["user_id"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo["role_id"], true)
			assert.Equal(t, roleInfo["role_id1"], true)
			roleInfo1, ok := roleIDs["user_id1"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo1["role_id1"], true)
		})
	})
}

func TestGetRolesByUserIDs2(t *testing.T) {
	Convey("GetRolesByUserIDs2, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)

		userID := "test1"

		role := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
			trace:   trace,
		}

		fields := []string{
			"f_user_id",
			"f_role_id",
		}
		Convey("roles id not exist", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			roleIDs, httpErr := role.GetRolesByUserIDs2(ctx, []string{})
			assert.Equal(t, len(roleIDs), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("Success1", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", "role_id"))
			roleIDs, httpErr := role.GetRolesByUserIDs2(ctx, []string{userID})
			assert.Equal(t, len(roleIDs), 1)
			assert.Equal(t, httpErr, nil)
			roleInfo, ok := roleIDs["user_id"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo["role_id"], true)
		})

		Convey("Success2", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("user_id", "role_id").AddRow("user_id1", "role_id1").AddRow("user_id", "role_id1"))
			roleIDs, httpErr := role.GetRolesByUserIDs2(ctx, []string{userID})
			assert.Equal(t, len(roleIDs), 2)
			assert.Equal(t, httpErr, nil)
			roleInfo, ok := roleIDs["user_id"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo["role_id"], true)
			assert.Equal(t, roleInfo["role_id1"], true)
			roleInfo1, ok := roleIDs["user_id1"]
			assert.Equal(t, ok, true)
			assert.Equal(t, roleInfo1["role_id1"], true)
		})
	})
}

// 测试GetUserIDsByRoleIDs
func TestGetUserIDsByRoleIDs(t *testing.T) {
	Convey("GetUserIDsByRoleIDs, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		ctx := context.Background()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		assert.NotEqual(t, db, nil)

		role := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
			trace:   trace,
		}

		Convey("roles为空", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			roleIDs, httpErr := role.GetUserIDsByRoleIDs(ctx, []interfaces.Role{})
			assert.Equal(t, len(roleIDs), 0)
			assert.Equal(t, httpErr, nil)
		})

		Convey("query error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnError(errors.New("query error"))
			_, httpErr := role.GetUserIDsByRoleIDs(ctx, []interfaces.Role{interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin})
			assert.NotEqual(t, httpErr, nil)
		})

		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_user_id", "f_role_id"}).AddRow("user_id", interfaces.SystemRoleSuperAdmin).
				AddRow("user_id1", interfaces.SystemRoleSysAdmin).AddRow("user_id", interfaces.SystemRoleSysAdmin))
			roleIDs, httpErr := role.GetUserIDsByRoleIDs(ctx, []interfaces.Role{interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin})
			assert.Equal(t, len(roleIDs), 2)
			assert.Equal(t, httpErr, nil)
			assert.Equal(t, roleIDs[interfaces.SystemRoleSuperAdmin], []string{"user_id"})
			data := map[string]bool{
				"user_id":  true,
				"user_id1": true,
			}
			assert.Equal(t, data[roleIDs[interfaces.SystemRoleSysAdmin][0]], true)
			assert.Equal(t, data[roleIDs[interfaces.SystemRoleSysAdmin][1]], true)
		})
	})
}
