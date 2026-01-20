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

const (
	operationJSONStr = `[{"id":"read","name":[{"language":"zh","value":"读取"}],"description":"读取文档","scope":["type","instance"]}]`
)

func TestResourceType_GetPagination(t *testing.T) {
	Convey("TestResourceType_GetPagination", t, func() {
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

		d := &resourceType{
			db:     db,
			logger: common.NewLogger(),
		}

		params := interfaces.ResourceTypePagination{
			Offset: 0,
			Limit:  10,
		}

		// 准备测试数据
		operationJSON := operationJSONStr
		expectedResource := interfaces.ResourceType{
			ID:          "doc",
			Name:        "文档",
			Description: "文档资源类型",
			InstanceURL: "https://example.com/doc",
			DataStruct:  "{}",
			CreateTime:  1234567890,
			ModifyTime:  1234567890,
		}

		Convey("count query error", func() {
			mock.ExpectQuery("^select count\\(1\\)").WillReturnError(mockErr)
			count, resources, err := d.GetPagination(ctx, params)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, count, 0)
			assert.Equal(t, resources, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("count scan error", func() {
			rows := sqlmock.NewRows([]string{"count"}).AddRow("invalid")
			mock.ExpectQuery("^select count\\(1\\)").WillReturnRows(rows)
			count, resources, err := d.GetPagination(ctx, params)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, count, 0)
			assert.Equal(t, resources, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("data query error", func() {
			// 先设置count查询成功
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery("^select").WillReturnRows(countRows)

			// 设置数据查询失败
			mock.ExpectQuery("^select.*limit.*offset").WillReturnError(mockErr)
			count, resources, err := d.GetPagination(ctx, params)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, count, 0)
			assert.Equal(t, resources, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("data scan error", func() {
			// 先设置count查询成功
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery("^select count\\(1\\)").WillReturnRows(countRows)

			// 设置数据查询返回无效数据
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation", "f_create_time", "f_modify_time"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", operationJSON, "invalid_time", "invalid_time")
			mock.ExpectQuery("^select.*limit.*offset").WillReturnRows(dataRows)

			count, resources, err := d.GetPagination(ctx, params)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, count, 1)
			assert.Equal(t, resources, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("json unmarshal error", func() {
			// 先设置count查询成功
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery("^select count\\(1\\)").WillReturnRows(countRows)

			// 设置数据查询返回无效JSON
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation", "f_create_time", "f_modify_time"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", "invalid_json", 1234567890, 1234567890)
			mock.ExpectQuery("^select.*limit.*offset").WillReturnRows(dataRows)

			count, resources, err := d.GetPagination(ctx, params)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, count, 1)
			assert.Equal(t, resources, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			// 设置count查询成功
			countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
			mock.ExpectQuery("^select count\\(1\\)").WillReturnRows(countRows)

			// 设置数据查询成功
			dataRows := sqlmock.NewRows([]string{
				"f_id", "f_name", "f_description",
				"f_instance_url", "f_data_struct", "f_operation", "f_create_time", "f_modify_time",
			}).
				AddRow(expectedResource.ID, expectedResource.Name, expectedResource.Description,
					expectedResource.InstanceURL, expectedResource.DataStruct, operationJSON, expectedResource.CreateTime, expectedResource.ModifyTime)
			mock.ExpectQuery("^select.*limit.*offset").WillReturnRows(dataRows)

			count, resources, err := d.GetPagination(ctx, params)
			assert.Equal(t, err, nil)
			assert.Equal(t, count, 1)
			assert.Equal(t, len(resources), 1)
			assert.Equal(t, resources[0].ID, expectedResource.ID)
			assert.Equal(t, resources[0].Name, expectedResource.Name)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestResourceType_Set(t *testing.T) {
	Convey("TestResourceType_Set", t, func() {
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

		d := &resourceType{
			db:     db,
			logger: common.NewLogger(),
		}

		// 准备测试数据
		resource := &interfaces.ResourceType{
			ID:          "doc",
			Name:        "文档",
			Description: "文档资源类型",
			InstanceURL: "https://example.com/doc",
			DataStruct:  "{}",
			Operation: []interfaces.ResourceTypeOperation{
				{
					ID:          "read",
					Name:        []interfaces.OperationName{{Language: "zh", Value: "读取"}},
					Description: "读取文档",
					Scope:       []interfaces.OperationScopeType{interfaces.ScopeType, interfaces.ScopeInstance},
				},
			},
		}

		Convey("check existence query error", func() {
			mock.ExpectQuery("^select f_id").WillReturnError(mockErr)
			err := d.Set(ctx, resource)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("update existing resource error", func() {
			// 设置查询返回存在的数据
			checkRows := sqlmock.NewRows([]string{"f_id"}).AddRow("doc")
			mock.ExpectQuery("^select f_id").WillReturnRows(checkRows)

			// 设置更新失败
			mock.ExpectExec("^update").WillReturnError(mockErr)
			err := d.Set(ctx, resource)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("update existing resource success", func() {
			// 设置查询返回存在的数据
			checkRows := sqlmock.NewRows([]string{"f_id"}).AddRow("doc")
			mock.ExpectQuery("^select f_id").WillReturnRows(checkRows)

			// 设置更新成功
			mock.ExpectExec("^update").WillReturnResult(sqlmock.NewResult(0, 1))
			err := d.Set(ctx, resource)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("insert new resource error", func() {
			// 设置查询返回不存在的数据
			checkRows := sqlmock.NewRows([]string{"f_id"})
			mock.ExpectQuery("^select f_id").WillReturnRows(checkRows)

			// 设置插入失败
			mock.ExpectExec("^insert").WillReturnError(mockErr)
			err := d.Set(ctx, resource)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("insert new resource success", func() {
			// 设置查询返回不存在的数据
			checkRows := sqlmock.NewRows([]string{"f_id"})
			mock.ExpectQuery("^select f_id").WillReturnRows(checkRows)

			// 设置插入成功
			mock.ExpectExec("^insert").WillReturnResult(sqlmock.NewResult(1, 1))
			err := d.Set(ctx, resource)
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

//nolint:dupl
func TestResourceType_Delete(t *testing.T) {
	Convey("TestResourceType_Delete", t, func() {
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

		d := &resourceType{
			db:     db,
			logger: common.NewLogger(),
		}

		Convey("delete error", func() {
			mock.ExpectExec("^delete").WillReturnError(mockErr)
			err := d.Delete(ctx, "test-id")
			assert.Equal(t, err, mockErr)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			mock.ExpectExec("^delete").WillReturnResult(sqlmock.NewResult(0, 1))
			err := d.Delete(ctx, "test-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestResourceType_GetByIDs(t *testing.T) {
	Convey("TestResourceType_GetByIDs", t, func() {
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

		d := &resourceType{
			db:     db,
			logger: common.NewLogger(),
		}

		// 准备测试数据
		operationJSON := operationJSONStr
		expectedResource := interfaces.ResourceType{
			ID:          "doc",
			Name:        "文档",
			Description: "文档资源类型",
			InstanceURL: "https://example.com/doc",
			DataStruct:  "{}",
		}

		Convey("empty resourceTypeIDs", func() {
			resourceMap, err := d.GetByIDs(ctx, []string{})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceMap), 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("query error", func() {
			mock.ExpectQuery("^select.*in").WillReturnError(mockErr)
			resourceMap, err := d.GetByIDs(ctx, []string{"doc", "folder"})
			assert.Equal(t, err, mockErr)
			assert.Equal(t, resourceMap, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			// 设置查询返回无效数据
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", "invalid_json")
			mock.ExpectQuery("^select.*in").WillReturnRows(dataRows)

			resourceMap, err := d.GetByIDs(ctx, []string{"doc"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, resourceMap, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("json unmarshal error", func() {
			// 设置查询返回无效JSON
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", "invalid_json")
			mock.ExpectQuery("^select.*in").WillReturnRows(dataRows)

			resourceMap, err := d.GetByIDs(ctx, []string{"doc"})
			assert.NotEqual(t, err, nil)
			assert.Equal(t, resourceMap, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Success", func() {
			// 设置查询成功
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow(expectedResource.ID, expectedResource.Name, expectedResource.Description, expectedResource.InstanceURL, expectedResource.DataStruct, operationJSON)
			mock.ExpectQuery("^select.*in").WillReturnRows(dataRows)

			resourceMap, err := d.GetByIDs(ctx, []string{"doc"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceMap), 1)
			assert.Equal(t, resourceMap["doc"].ID, expectedResource.ID)
			assert.Equal(t, resourceMap["doc"].Name, expectedResource.Name)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("Multiple resources success", func() {
			// 设置查询成功，返回多个资源
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", operationJSON).
				AddRow("folder", "文件夹", "文件夹资源类型", "https://example.com/folder", "{}", operationJSON)
			mock.ExpectQuery("^select.*in").WillReturnRows(dataRows)

			resourceMap, err := d.GetByIDs(ctx, []string{"doc", "folder"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceMap), 2)
			assert.Equal(t, resourceMap["doc"].ID, "doc")
			assert.Equal(t, resourceMap["folder"].ID, "folder")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}

func TestResourceType_NewResource(t *testing.T) {
	Convey("TestResourceType_NewResource", t, func() {
		Convey("Singleton pattern", func() {
			// 测试单例模式
			instance1 := NewResource()
			instance2 := NewResource()

			assert.Equal(t, instance1, instance2)
			assert.NotEqual(t, instance1, nil)
			assert.NotEqual(t, instance2, nil)
		})
	})
}

func TestResourceType_GetAllInternal(t *testing.T) {
	Convey("TestResourceType_GetAllInternal", t, func() {
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

		d := &resourceType{
			db:     db,
			logger: common.NewLogger(),
		}

		// 准备测试数据
		operationJSON := operationJSONStr
		expectedResource := interfaces.ResourceType{
			ID:          "doc",
			Name:        "文档",
			Description: "文档资源类型",
			InstanceURL: "https://example.com/doc",
			DataStruct:  "{}",
		}

		Convey("query error", func() {
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnError(mockErr)
			resourceTypes, err := d.GetAllInternal(ctx)
			assert.Equal(t, err, mockErr)
			assert.Equal(t, resourceTypes, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("scan error", func() {
			// 设置查询返回无效数据
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", "invalid_json")
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnRows(dataRows)

			resourceTypes, err := d.GetAllInternal(ctx)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, resourceTypes, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("json unmarshal error", func() {
			// 设置查询返回无效JSON
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", "invalid_json")
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnRows(dataRows)

			resourceTypes, err := d.GetAllInternal(ctx)
			assert.NotEqual(t, err, nil)
			assert.Equal(t, resourceTypes, nil)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("empty result success", func() {
			// 设置查询返回空结果
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"})
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnRows(dataRows)

			resourceTypes, err := d.GetAllInternal(ctx)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceTypes), 0)
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("single resource success", func() {
			// 设置查询成功，返回单个资源
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow(expectedResource.ID, expectedResource.Name, expectedResource.Description, expectedResource.InstanceURL, expectedResource.DataStruct, operationJSON)
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnRows(dataRows)

			resourceTypes, err := d.GetAllInternal(ctx)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceTypes), 1)
			assert.Equal(t, resourceTypes[0].ID, expectedResource.ID)
			assert.Equal(t, resourceTypes[0].Name, expectedResource.Name)
			assert.Equal(t, resourceTypes[0].Description, expectedResource.Description)
			assert.Equal(t, resourceTypes[0].InstanceURL, expectedResource.InstanceURL)
			assert.Equal(t, resourceTypes[0].DataStruct, expectedResource.DataStruct)
			assert.Equal(t, len(resourceTypes[0].Operation), 1)
			assert.Equal(t, resourceTypes[0].Operation[0].ID, "read")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})

		Convey("multiple resources success", func() {
			// 设置查询成功，返回多个资源
			dataRows := sqlmock.NewRows([]string{"f_id", "f_name", "f_description", "f_instance_url", "f_data_struct", "f_operation"}).
				AddRow("doc", "文档", "文档资源类型", "https://example.com/doc", "{}", operationJSON).
				AddRow("folder", "文件夹", "文件夹资源类型", "https://example.com/folder", "{}", operationJSON).
				AddRow("file", "文件", "文件资源类型", "https://example.com/file", "{}", operationJSON)
			mock.ExpectQuery("^select.*f_id.*f_name.*f_description.*f_instance_url.*f_data_struct.*f_operation").WillReturnRows(dataRows)

			resourceTypes, err := d.GetAllInternal(ctx)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(resourceTypes), 3)
			assert.Equal(t, resourceTypes[0].ID, "doc")
			assert.Equal(t, resourceTypes[1].ID, "folder")
			assert.Equal(t, resourceTypes[2].ID, "file")
			assert.Equal(t, resourceTypes[0].Name, "文档")
			assert.Equal(t, resourceTypes[1].Name, "文件夹")
			assert.Equal(t, resourceTypes[2].Name, "文件")
			assert.Equal(t, mock.ExpectationsWereMet(), nil)
		})
	})
}
