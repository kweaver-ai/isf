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

func TestObligation_Add(t *testing.T) {
	Convey("TestObligation_Add", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		info := &interfaces.ObligationInfo{
			ID:          "test-id",
			TypeID:      "test-type-id",
			Name:        "test-name",
			Description: "test-desc",
			Value: map[string]any{
				"key1": "value1",
				"key2": 123,
			},
		}

		Convey("Add success", func() {
			mock.ExpectExec("^insert into").WillReturnResult(sqlmock.NewResult(1, 1))
			err := o.Add(ctx, info)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Add error", func() {
			mock.ExpectExec("^insert into").WillReturnError(mockErr)
			err := o.Add(ctx, info)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Add with nil Value", func() {
			infoNilValue := &interfaces.ObligationInfo{
				ID:          "test-id",
				TypeID:      "test-type-id",
				Name:        "test-name",
				Description: "test-desc",
				Value:       nil,
			}
			mock.ExpectExec("^insert into").WillReturnResult(sqlmock.NewResult(1, 1))
			err := o.Add(ctx, infoNilValue)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligation_Update(t *testing.T) {
	Convey("TestObligation_Update", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Update name only", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Update(ctx, "test-id", "new-name", true, "", false, nil, false)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update description only", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Update(ctx, "test-id", "", false, "new-desc", true, nil, false)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update value only", func() {
			value := map[string]any{"key": "value"}
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Update(ctx, "test-id", "", false, "", false, value, true)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update all fields", func() {
			value := map[string]any{"key": "value"}
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Update(ctx, "test-id", "new-name", true, "new-desc", true, value, true)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update error", func() {
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := o.Update(ctx, "test-id", "new-name", true, "", false, nil, false)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Update no fields changed", func() {
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := o.Update(ctx, "test-id", "", false, "", false, nil, false)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

//nolint:dupl
func TestObligation_Delete(t *testing.T) {
	Convey("TestObligation_Delete", t, func() {
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

		o := &obligation{
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

func TestObligation_GetByID(t *testing.T) {
	Convey("TestObligation_GetByID", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Query error", func() {
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnError(mockErr)
			_, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "type-id", "test-name", "test-desc",
				`{"key":"value"}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			info, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "test-id")
			assert.Equal(t, info.TypeID, "type-id")
			assert.Equal(t, info.Name, "test-name")
			assert.Equal(t, info.Description, "test-desc")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with empty result", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			})
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			info, err := o.GetByID(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id",
			}).AddRow("test-id", "type-id")
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByID(ctx, "test-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid JSON value", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "type-id", "test-name", "test-desc",
				`{invalid json}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByID(ctx, "test-id")
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligation_Get(t *testing.T) {
	Convey("TestObligation_Get", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		searchInfo := &interfaces.ObligationSearchInfo{
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
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnError(mockErr)
			_, _, err := o.Get(ctx, searchInfo)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success", func() {
			countRows := sqlmock.NewRows([]string{"count(1)"}).AddRow(2)
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			dataRows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id-1", "type-id-1", "test-name-1", "test-desc-1",
				`{"key":"value1"}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"test-id-2", "type-id-2", "test-name-2", "test-desc-2",
				`{"key":"value2"}`,
				int64(1234567891), int64(1234567891),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(dataRows)
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
			dataRows := sqlmock.NewRows([]string{"f_id", "f_type_id"}).AddRow("test-id", "type-id")
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(dataRows)
			_, _, err := o.Get(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid JSON in data", func() {
			countRows := sqlmock.NewRows([]string{"count(1)"}).AddRow(1)
			mock.ExpectQuery("^select count").WillReturnRows(countRows)
			dataRows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"test-id", "type-id", "test-name", "test-desc",
				`{invalid json}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(dataRows)
			_, _, err := o.Get(ctx, searchInfo)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligation_GetByObligationTypeIDs(t *testing.T) {
	Convey("TestObligation_GetByObligationTypeIDs", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Empty obligationTypeIDMap", func() {
			resultInfos, err := o.GetByObligationTypeIDs(ctx, map[string]bool{})
			assert.Equal(t, err, nil)
			assert.Equal(t, resultInfos, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query error", func() {
			typeIDMap := map[string]bool{
				"type-id-1": true,
				"type-id-2": true,
			}
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnError(mockErr)
			_, err := o.GetByObligationTypeIDs(ctx, typeIDMap)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with multiple type IDs", func() {
			typeIDMap := map[string]bool{
				"type-id-1": true,
				"type-id-2": true,
			}
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"obl-id-1", "type-id-1", "name-1", "desc-1",
				`{"key":"value1"}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"obl-id-2", "type-id-1", "name-2", "desc-2",
				`{"key":"value2"}`,
				int64(1234567891), int64(1234567891),
			).AddRow(
				"obl-id-3", "type-id-2", "name-3", "desc-3",
				`{"key":"value3"}`,
				int64(1234567892), int64(1234567892),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			resultInfos, err := o.GetByObligationTypeIDs(ctx, typeIDMap)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resultInfos), 2)
			assert.Equal(t, len(resultInfos["type-id-1"]), 2)
			assert.Equal(t, len(resultInfos["type-id-2"]), 1)
			assert.Equal(t, resultInfos["type-id-1"][0].ID, "obl-id-1")
			assert.Equal(t, resultInfos["type-id-1"][1].ID, "obl-id-2")
			assert.Equal(t, resultInfos["type-id-2"][0].ID, "obl-id-3")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			typeIDMap := map[string]bool{
				"type-id-1": true,
			}
			rows := sqlmock.NewRows([]string{"f_id", "f_type_id"}).AddRow("obl-id-1", "type-id-1")
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByObligationTypeIDs(ctx, typeIDMap)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid JSON value", func() {
			typeIDMap := map[string]bool{
				"type-id-1": true,
			}
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"obl-id-1", "type-id-1", "name-1", "desc-1",
				`{invalid json}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByObligationTypeIDs(ctx, typeIDMap)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestObligation_GetByIDs(t *testing.T) {
	Convey("TestObligation_GetByIDs", t, func() {
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

		o := &obligation{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("Empty IDs", func() {
			resultInfos, err := o.GetByIDs(ctx, []string{})
			assert.Equal(t, err, nil)
			assert.Equal(t, resultInfos, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query error", func() {
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnError(mockErr)
			_, err := o.GetByIDs(ctx, []string{"id-1", "id-2"})
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with multiple IDs", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"id-1", "type-id-1", "name-1", "desc-1",
				`{"key":"value1"}`,
				int64(1234567890), int64(1234567890),
			).AddRow(
				"id-2", "type-id-2", "name-2", "desc-2",
				`{"key":"value2"}`,
				int64(1234567891), int64(1234567891),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			resultInfos, err := o.GetByIDs(ctx, []string{"id-1", "id-2"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resultInfos), 2)
			assert.Equal(t, resultInfos[0].ID, "id-1")
			assert.Equal(t, resultInfos[1].ID, "id-2")
			assert.Equal(t, resultInfos[0].TypeID, "type-id-1")
			assert.Equal(t, resultInfos[1].TypeID, "type-id-2")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Query success with single ID", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"id-1", "type-id-1", "name-1", "desc-1",
				`{"key":"value1"}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			resultInfos, err := o.GetByIDs(ctx, []string{"id-1"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resultInfos), 1)
			assert.Equal(t, resultInfos[0].ID, "id-1")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Scan error", func() {
			rows := sqlmock.NewRows([]string{"f_id", "f_type_id"}).AddRow("id-1", "type-id-1")
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByIDs(ctx, []string{"id-1"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Invalid JSON value", func() {
			rows := sqlmock.NewRows([]string{
				"f_id", "f_type_id", "f_name", "f_description", "f_value", "f_created_at", "f_modified_at",
			}).AddRow(
				"id-1", "type-id-1", "name-1", "desc-1",
				`{invalid json}`,
				int64(1234567890), int64(1234567890),
			)
			mock.ExpectQuery("^select f_id, f_type_id").WillReturnRows(rows)
			_, err := o.GetByIDs(ctx, []string{"id-1"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestNewObligation(t *testing.T) {
	Convey("TestNewObligation", t, func() {
		Convey("Create new obligation service", func() {
			// Reset the singleton for testing
			obligationOnce = sync.Once{}
			obligationService = nil

			service1 := NewObligation()
			assert.NotEqual(t, service1, nil)

			// Test singleton pattern
			service2 := NewObligation()
			assert.Equal(t, service1, service2)
		})
	})
}
