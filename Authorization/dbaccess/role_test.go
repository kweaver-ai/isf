//nolint:goconst
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

func TestDBAddRoles(t *testing.T) {
	Convey("TestDBAddRoles", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}
		roles := []interfaces.RoleInfo{
			{
				ID:          "1",
				Name:        "test",
				Description: "test",
				ResourceTypeScopeInfo: interfaces.ResourceTypeScopeInfo{
					Unlimited: false,
					Types: []interfaces.ResourceTypeScope{
						{
							ResourceTypeID: "test",
						},
					},
				},
			},
		}

		Convey("入参 列表为空 ", func() {
			err := b.AddRoles(ctx, []interfaces.RoleInfo{})
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("insert error", func() {
			mock.ExpectExec("^insert").WillReturnError(mockErr)
			err := b.AddRoles(ctx, roles)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^insert").WillReturnResult(sqlmock.NewResult(1, 1))
			err := b.AddRoles(ctx, roles)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBDeleteRole(t *testing.T) {
	Convey("TestDBDeleteRole", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := b.DeleteRole(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.DeleteRole(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBModifyRole(t *testing.T) {
	Convey("TestDBModifyRole", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		resourceTypeScopes := interfaces.ResourceTypeScopeInfo{
			Unlimited: false,
			Types: []interfaces.ResourceTypeScope{
				{
					ResourceTypeID: "test",
				},
			},
		}

		Convey("update error", func() {
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := b.ModifyRole(ctx, "test-id", "new-name", true, "new-desc", true, resourceTypeScopes, true)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - all fields changed", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.ModifyRole(ctx, "test-id", "new-name", true, "new-desc", true, resourceTypeScopes, true)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - only name changed", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.ModifyRole(ctx, "test-id", "new-name", true, "", false, interfaces.ResourceTypeScopeInfo{}, false)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - only description changed", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.ModifyRole(ctx, "test-id", "", false, "new-desc", true, interfaces.ResourceTypeScopeInfo{}, false)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - only resourceTypeScopes changed", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.ModifyRole(ctx, "test-id", "", false, "", false, resourceTypeScopes, true)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleByID(t *testing.T) {
	Convey("TestDBGetRoleByID", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			info, err := b.GetRoleByID(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, info.ID, "")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_source", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", interfaces.RoleSourceSystem, "invalid-json", "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			info, err := b.GetRoleByID(ctx, "test-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_source", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", interfaces.RoleSourceSystem, resourceTypeScopes, "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			info, err := b.GetRoleByID(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, info.Name, "test-name")
			assert.Equal(t, info.Description, "test-desc")
			assert.Equal(t, info.RoleSource, interfaces.RoleSourceSystem)
			assert.Equal(t, len(info.Types), 1)
			assert.Equal(t, info.Types[0].ResourceTypeID, "test")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleByName(t *testing.T) {
	Convey("TestDBGetRoleByName", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			info, err := b.GetRoleByName(ctx, "test-name")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, info.ID, "")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", "invalid-json", "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			info, err := b.GetRoleByName(ctx, "test-name")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", resourceTypeScopes, "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			info, err := b.GetRoleByName(ctx, "test-name")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, info.Name, "test-name")
			assert.Equal(t, info.Description, "test-desc")
			assert.Equal(t, len(info.Types), 1)
			assert.Equal(t, info.Types[0].ResourceTypeID, "test")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoleByIDs(t *testing.T) {
	Convey("TestDBGetRoleByIDs", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		Convey("empty ids", func() {
			infoMap, err := b.GetRoleByIDs(ctx, []string{})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infoMap), 0)
		})

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			infoMap, err := b.GetRoleByIDs(ctx, []string{"test-id-1", "test-id-2"})
			assert.Equal(t, err, mockErr)
			assert.Equal(t, infoMap, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", "invalid-json", "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			_, err := b.GetRoleByIDs(ctx, []string{"test-id"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_source", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-1", "test-name-1", interfaces.RoleSourceSystem, resourceTypeScopes, "test-desc-1", 1234567890, 1234567890).
				AddRow("test-id-2", "test-name-2", interfaces.RoleSourceSystem, resourceTypeScopes, "test-desc-2", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			infoMap, err := b.GetRoleByIDs(ctx, []string{"test-id-1", "test-id-2"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infoMap), 2)
			assert.Equal(t, infoMap["test-id-1"].Name, "test-name-1")
			assert.Equal(t, infoMap["test-id-2"].Name, "test-name-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRoles(t *testing.T) {
	Convey("TestDBGetRoles", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		searchInfo := interfaces.RoleSearchInfo{
			Offset:      0,
			Limit:       10,
			RoleSources: []interfaces.RoleSource{interfaces.RoleSourceSystem},
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			roles, err := b.GetRoles(ctx, searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, roles, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", "invalid-json", "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			roles, err := b.GetRoles(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, roles, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_source", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-1", "test-name-1", interfaces.RoleSourceSystem, resourceTypeScopes, "test-desc-1", 1234567890, 1234567890).
				AddRow("test-id-2", "test-name-2", interfaces.RoleSourceSystem, resourceTypeScopes, "test-desc-2", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			roles, err := b.GetRoles(ctx, searchInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 2)
			assert.Equal(t, roles[0].ID, "test-id-1")
			assert.Equal(t, roles[0].Name, "test-name-1")
			assert.Equal(t, roles[0].RoleSource, interfaces.RoleSourceSystem)
			assert.Equal(t, roles[1].ID, "test-id-2")
			assert.Equal(t, roles[1].Name, "test-name-2")
			assert.Equal(t, roles[1].RoleSource, interfaces.RoleSourceSystem)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetRolesSum(t *testing.T) {
	Convey("TestDBGetRolesSum", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		searchInfo := interfaces.RoleSearchInfo{
			Offset:      0,
			Limit:       10,
			RoleSources: []interfaces.RoleSource{interfaces.RoleSourceSystem},
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			num, err := b.GetRolesSum(ctx, searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, num, 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"count(*)"}).AddRow("invalid")
			mock.ExpectQuery("^select").WillReturnRows(rows)
			num, err := b.GetRolesSum(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, num, 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(5)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			num, err := b.GetRolesSum(ctx, searchInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, num, 5)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBSetRoleByID(t *testing.T) {
	Convey("TestDBSetRoleByID", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		roleInfo := &interfaces.RoleInfo{
			ID:          "test-id",
			Name:        "test-name",
			Description: "test-description",
			ResourceTypeScopeInfo: interfaces.ResourceTypeScopeInfo{
				Unlimited: false,
				Types: []interfaces.ResourceTypeScope{
					{
						ResourceTypeID: "test-type",
					},
				},
			},
		}

		Convey("update error", func() {
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := b.SetRoleByID(ctx, "test-id", roleInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.SetRoleByID(ctx, "test-id", roleInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success with unlimited resource scope", func() {
			roleInfoUnlimited := &interfaces.RoleInfo{
				ID:          "test-id",
				Name:        "test-name",
				Description: "test-description",
				ResourceTypeScopeInfo: interfaces.ResourceTypeScopeInfo{
					Unlimited: true,
					Types:     []interfaces.ResourceTypeScope{},
				},
			}
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.SetRoleByID(ctx, "test-id", roleInfoUnlimited)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success with empty resource types", func() {
			roleInfoEmpty := &interfaces.RoleInfo{
				ID:          "test-id",
				Name:        "test-name",
				Description: "test-description",
				ResourceTypeScopeInfo: interfaces.ResourceTypeScopeInfo{
					Unlimited: false,
					Types:     []interfaces.ResourceTypeScope{},
				},
			}
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := b.SetRoleByID(ctx, "test-id", roleInfoEmpty)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestDBGetAllRolesInternal(t *testing.T) {
	Convey("TestDBGetAllRolesInternal", t, func() {
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

		b := &role{
			db:      db,
			logger:  common.NewLogger(),
			dbTrace: db,
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select").WillReturnError(mockErr)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, roles, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id", "test-name", "invalid-json", "test-desc", 1234567890, 1234567890)
			mock.ExpectQuery("^select").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, roles, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - empty keyword", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-1", "test-name-1", resourceTypeScopes, "test-desc-1", 1234567890, 1234567890).
				AddRow("test-id-2", "test-name-2", resourceTypeScopes, "test-desc-2", 1234567890, 1234567890)
			mock.ExpectQuery("^select.*where").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 2)
			assert.Equal(t, roles[0].ID, "test-id-1")
			assert.Equal(t, roles[0].Name, "test-name-1")
			assert.Equal(t, roles[0].Description, "test-desc-1")
			assert.Equal(t, roles[0].Unlimited, false)
			assert.Equal(t, len(roles[0].Types), 1)
			assert.Equal(t, roles[0].Types[0].ResourceTypeID, "test")
			assert.Equal(t, roles[1].ID, "test-id-2")
			assert.Equal(t, roles[1].Name, "test-name-2")
			assert.Equal(t, roles[1].Description, "test-desc-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - with keyword", func() {
			resourceTypeScopes := `{"unlimited":true,"types":[]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-3", "admin-role", resourceTypeScopes, "admin description", 1234567890, 1234567890)
			mock.ExpectQuery("^select.*where").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "admin")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 1)
			assert.Equal(t, roles[0].ID, "test-id-3")
			assert.Equal(t, roles[0].Name, "admin-role")
			assert.Equal(t, roles[0].Description, "admin description")
			assert.Equal(t, roles[0].Unlimited, true)
			assert.Equal(t, len(roles[0].Types), 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - empty result", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"})
			mock.ExpectQuery("^select.*where").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - complex resource type scope", func() {
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"type1","name":"type1"},{"id":"type2","name":"type2"}]}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-4", "complex-role", resourceTypeScopes, "complex description", 1234567890, 1234567890)
			mock.ExpectQuery("^select.*where").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 1)
			assert.Equal(t, roles[0].ID, "test-id-4")
			assert.Equal(t, roles[0].Name, "complex-role")
			assert.Equal(t, roles[0].Unlimited, false)
			assert.Equal(t, len(roles[0].Types), 2)
			assert.Equal(t, roles[0].Types[0].ResourceTypeID, "type1")
			assert.Equal(t, roles[0].Types[1].ResourceTypeID, "type2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success - malformed JSON but valid structure", func() {
			// Test with JSON that has extra fields but valid structure
			resourceTypeScopes := `{"unlimited":false,"types":[{"id":"test","name":"test","extra":"field"}],"extra_field":"value"}`
			rows := sqlmock.NewRows([]string{"f_id", "f_name", "f_resource_scope", "f_description", "f_created_time", "f_modify_time"}).
				AddRow("test-id-5", "extra-role", resourceTypeScopes, "extra description", 1234567890, 1234567890)
			mock.ExpectQuery("^select.*where").WillReturnRows(rows)
			roles, err := b.GetAllUserRolesInternal(ctx, "")
			assert.Equal(t, err, nil)
			assert.Equal(t, len(roles), 1)
			assert.Equal(t, roles[0].ID, "test-id-5")
			assert.Equal(t, roles[0].Name, "extra-role")
			assert.Equal(t, roles[0].Unlimited, false)
			assert.Equal(t, len(roles[0].Types), 1)
			assert.Equal(t, roles[0].Types[0].ResourceTypeID, "test")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}
