package dbaccess

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authorization/common"
	"Authorization/interfaces"
)

func TestObligationType_Set(t *testing.T) {
	Convey("TestObligationType_Set", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		info := &interfaces.ObligationTypeInfo{
			ID:          "test-id",
			Name:        "test-name",
			Description: "test-desc",
			Schema: map[string]any{
				"type": "object",
			},
			DefaultValue: map[string]any{
				"value": "default",
			},
			UiSchema: map[string]any{
				"ui:widget": "text",
			},
			ResourceTypeScope: interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: true,
				Types:     []interfaces.ObligationResourceTypeScope{},
			},
		}

		Convey("Query error when checking existence", func() {
			mock.ExpectQuery("^select f_id from").WillReturnError(mockErr)
			err := o.Set(ctx, info)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Insert new record - success", func() {
			rows := sqlmock.NewRows([]string{"f_id"})
			mock.ExpectQuery("^select f_id from").WillReturnRows(rows)
			mock.ExpectExec("^insert into").WillReturnResult(sqlmock.NewResult(1, 1))
			err := o.Set(ctx, info)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Insert new record - error", func() {
			rows := sqlmock.NewRows([]string{"f_id"})
			mock.ExpectQuery("^select f_id from").WillReturnRows(rows)
			mock.ExpectExec("^insert into").WillReturnError(mockErr)
			err := o.Set(ctx, info)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update existing record - success", func() {
			rows := sqlmock.NewRows([]string{"f_id"}).AddRow("test-id")
			mock.ExpectQuery("^select f_id from").WillReturnRows(rows)
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Set(ctx, info)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update existing record - error", func() {
			rows := sqlmock.NewRows([]string{"f_id"}).AddRow("test-id")
			mock.ExpectQuery("^select f_id from").WillReturnRows(rows)
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := o.Set(ctx, info)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Set with nil DefaultValue and UiSchema", func() {
			infoNoOptional := &interfaces.ObligationTypeInfo{
				ID:          "test-id",
				Name:        "test-name",
				Description: "test-desc",
				Schema: map[string]any{
					"type": "object",
				},
				ResourceTypeScope: interfaces.ObligationResourceTypeScopeInfo{
					Unlimited: false,
					Types: []interfaces.ObligationResourceTypeScope{
						{
							ResourceTypeID: "resource-type-1",
							OperationsScope: interfaces.ObligationOperationsScopeInfo{
								Unlimited: false,
								Operations: []interfaces.ObligationOperation{
									{ID: "op-1"},
									{ID: "op-2"},
								},
							},
						},
					},
				},
			}
			rows := sqlmock.NewRows([]string{"f_id"})
			mock.ExpectQuery("^select f_id from").WillReturnRows(rows)
			mock.ExpectExec("^insert into").WillReturnResult(sqlmock.NewResult(1, 1))
			err := o.Set(ctx, infoNoOptional)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

//nolint:dupl
func TestObligationType_Delete(t *testing.T) {
	Convey("TestObligationType_Delete", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Delete error", func() {
			mock.ExpectExec("^delete from").WillReturnError(mockErr)
			err := o.Delete(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Delete success", func() {
			mock.ExpectExec("^delete from").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Delete(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligationType_GetByID(t *testing.T) {
	Convey("TestObligationType_GetByID", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Query error", func() {
			mock.ExpectQuery("^select f_id, f_name").WillReturnError(mockErr)
			_, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with all fields", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "test-name", "test-desc",
				`{"type":"object"}`, `{"value":"default"}`, `{"ui:widget":"text"}`,
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			info, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, info.Name, "test-name")
			assert.Equal(t, info.Description, "test-desc")
			assert.Equal(t, info.ResourceTypeScope.Unlimited, true)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with empty optional fields", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "test-name", "test-desc",
				`{"type":"object"}`, "", "",
				`{"unlimited":false,"resource_types":[{"id":"rt-1","applicable_operations":{"unlimited":false,"operations":[{"id":"op-1"}]}}]}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			info, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, info.ResourceTypeScope.Unlimited, false)
			assert.Equal(t, len(info.ResourceTypeScope.Types), 1)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description",
			}).AddRow("test-id", "test-name", "test-desc")
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetByID(ctx, "test-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid JSON schema", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "test-name", "test-desc",
				`{invalid json}`, "", "",
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetByID(ctx, "test-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligationType_Get(t *testing.T) {
	Convey("TestObligationType_Get", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		searchInfo := &interfaces.ObligationTypeSearchInfo{
			Offset: 0,
			Limit:  10,
		}

		Convey("Count query error", func() {
			mock.ExpectQuery("^select count").WillReturnError(mockErr)
			_, _, err := o.Get(ctx, searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Count scan error", func() {
			countRows := sqlmock.NewRows([]string{"count"})
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			_, _, err := o.Get(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Data query error", func() {
			countRows := sqlmock.NewRows([]string{"count(1)"}).AddRow(2)
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			mock.ExpectQuery("^select f_id, f_name").WillReturnError(mockErr)
			_, _, err := o.Get(ctx, searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success", func() {
			countRows := sqlmock.NewRows([]string{"count(1)"}).AddRow(2)
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			dataRows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id-1", "test-name-1", "test-desc-1",
				`{"type":"object"}`, `{"value":"default1"}`, `{"ui:widget":"text"}`,
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"test-id-2", "test-name-2", "test-desc-2",
				`{"type":"string"}`, "", "",
				`{"unlimited":false,"resource_types":[]}`,
				int64(1234567891), int64(1234567891),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(dataRows)
			count, infos, err := o.Get(ctx, searchInfo)
			assert.Equal(t, err, nil)
			assert.Equal(t, count, 2)
			assert.Equal(t, len(infos), 2)
			assert.Equal(t, infos[0].ID, "test-id-1")
			assert.Equal(t, infos[1].ID, "test-id-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Data scan error", func() {
			countRows := sqlmock.NewRows([]string{"count(1)"}).AddRow(1)
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name"}).AddRow("test-id", "test-name")
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(dataRows)
			_, _, err := o.Get(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligationType_GetAll(t *testing.T) {
	Convey("TestObligationType_GetAll", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Query error", func() {
			mock.ExpectQuery("^select f_id, f_name").WillReturnError(mockErr)
			_, err := o.GetAll(ctx)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with multiple records", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id-1", "test-name-1", "test-desc-1",
				`{"type":"object"}`, `{"value":"default1"}`, `{"ui:widget":"text"}`,
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"test-id-2", "test-name-2", "test-desc-2",
				`{"type":"string"}`, "", "",
				`{"unlimited":false,"resource_types":[{"id":"rt-1","applicable_operations":{"unlimited":true,"operations":[]}}]}`,
				int64(1234567891), int64(1234567891),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			infos, err := o.GetAll(ctx)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infos), 2)
			assert.Equal(t, infos[0].ID, "test-id-1")
			assert.Equal(t, infos[1].ID, "test-id-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with empty result", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			})
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			infos, err := o.GetAll(ctx)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infos), 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name"}).AddRow("test-id", "test-name")
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetAll(ctx)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid resource type scope JSON", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "test-name", "test-desc",
				`{"type":"object"}`, "", "",
				`{invalid json}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetAll(ctx)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligationType_GetByIDs(t *testing.T) {
	Convey("TestObligationType_GetByIDs", t, func() {
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

		o := &obligationType{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Empty IDs", func() {
			infos, err := o.GetByIDs(ctx, []string{})
			assert.Equal(t, err, nil)
			assert.Equal(t, infos, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query error", func() {
			mock.ExpectQuery("^select f_id, f_name").WillReturnError(mockErr)
			_, err := o.GetByIDs(ctx, []string{"id-1", "id-2"})
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with multiple IDs", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"id-1", "name-1", "desc-1",
				`{"type":"object"}`, `{"value":"default1"}`, `{"ui:widget":"text"}`,
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"id-2", "name-2", "desc-2",
				`{"type":"string"}`, "", "",
				`{"unlimited":false,"resource_types":[{"id":"rt-1","applicable_operations":{"unlimited":false,"operations":[{"id":"op-1"}]}}]}`,
				int64(1234567891), int64(1234567891),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			infos, err := o.GetByIDs(ctx, []string{"id-1", "id-2"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(infos), 2)
			assert.Equal(t, infos[0].ID, "id-1")
			assert.Equal(t, infos[1].ID, "id-2")
			assert.Equal(t, infos[1].ResourceTypeScope.Unlimited, false)
			assert.Equal(t, len(infos[1].ResourceTypeScope.Types), 1)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_name"}).AddRow("id-1", "name-1")
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetByIDs(ctx, []string{"id-1"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid default value JSON", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description", "f_schema", "f_default_value",
				"f_ui_schema", "f_applicable_resource_types", "f_created_at", "f_modified_at",
			}).AddRow(
				"id-1", "name-1", "desc-1",
				`{"type":"object"}`, `{invalid json}`, "",
				`{"unlimited":true,"resource_types":[]}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_name").WillReturnRows(rows)
			_, err := o.GetByIDs(ctx, []string{"id-1"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligationType_ResourceTypeScopeInfoToString(t *testing.T) {
	Convey("TestObligationType_ResourceTypeScopeInfoToString", t, func() {
		o := &obligationType{
			logger: common.NewLogger(),
		}

		Convey("Unlimited resource type scope", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: true,
				Types:     []interfaces.ObligationResourceTypeScope{},
			}
			result, err := o.resourceTypeScopeInfoToString(info)
			assert.Equal(t, err, nil)
			assert.NotEqual(t, result, "")
			assert.NotEqual(t, len(result), 0)
		})

		Convey("Limited resource type scope with operations", func() {
			info := interfaces.ObligationResourceTypeScopeInfo{
				Unlimited: false,
				Types: []interfaces.ObligationResourceTypeScope{
					{
						ResourceTypeID: "resource-type-1",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited: false,
							Operations: []interfaces.ObligationOperation{
								{ID: "op-1"},
								{ID: "op-2"},
							},
						},
					},
					{
						ResourceTypeID: "resource-type-2",
						OperationsScope: interfaces.ObligationOperationsScopeInfo{
							Unlimited:  true,
							Operations: []interfaces.ObligationOperation{},
						},
					},
				},
			}
			result, err := o.resourceTypeScopeInfoToString(info)
			assert.Equal(t, err, nil)
			assert.NotEqual(t, result, "")
		})
	})
}

func TestObligationType_StringToResourceTypeScopeInfo(t *testing.T) {
	Convey("TestObligationType_StringToResourceTypeScopeInfo", t, func() {
		o := &obligationType{
			logger: common.NewLogger(),
		}

		Convey("Unlimited resource type scope", func() {
			jsonStr := `{"unlimited":true,"resource_types":[]}`
			info, err := o.stringToResourceTypeScopeInfo(jsonStr)
			assert.Equal(t, err, nil)
			assert.Equal(t, info.Unlimited, true)
			assert.Equal(t, len(info.Types), 0)
		})

		Convey("Limited resource type scope with operations", func() {
			jsonStr := `{
				"unlimited": false,
				"resource_types": [
					{
						"id": "resource-type-1",
						"applicable_operations": {
							"unlimited": false,
							"operations": [
								{"id": "op-1"},
								{"id": "op-2"}
							]
						}
					},
					{
						"id": "resource-type-2",
						"applicable_operations": {
							"unlimited": true,
							"operations": []
						}
					}
				]
			}`
			info, err := o.stringToResourceTypeScopeInfo(jsonStr)
			assert.Equal(t, err, nil)
			assert.Equal(t, info.Unlimited, false)
			assert.Equal(t, len(info.Types), 2)
			assert.Equal(t, info.Types[0].ResourceTypeID, "resource-type-1")
			assert.Equal(t, info.Types[0].OperationsScope.Unlimited, false)
			assert.Equal(t, len(info.Types[0].OperationsScope.Operations), 2)
			assert.Equal(t, info.Types[0].OperationsScope.Operations[0].ID, "op-1")
			assert.Equal(t, info.Types[1].ResourceTypeID, "resource-type-2")
			assert.Equal(t, info.Types[1].OperationsScope.Unlimited, true)
		})

		Convey("Invalid JSON", func() {
			jsonStr := `{invalid json}`
			_, err := o.stringToResourceTypeScopeInfo(jsonStr)
			assert.NotEqual(t, err, nil)
		})

		Convey("Empty resource types", func() {
			jsonStr := `{"unlimited":false,"resource_types":[]}`
			info, err := o.stringToResourceTypeScopeInfo(jsonStr)
			assert.Equal(t, err, nil)
			assert.Equal(t, info.Unlimited, false)
			assert.Equal(t, len(info.Types), 0)
		})
	})
}

func TestNewObligationType(t *testing.T) {
	Convey("TestNewObligationType", t, func() {
		Convey("Create new obligation type service", func() {
			// Reset the singleton for testing
			obligationTypeOnce = sync.Once{}
			obligationTypeService = nil

			service1 := NewObligationType()
			assert.NotEqual(t, service1, nil)

			// Test singleton pattern
			service2 := NewObligationType()
			assert.Equal(t, service1, service2)
		})
	})
}
