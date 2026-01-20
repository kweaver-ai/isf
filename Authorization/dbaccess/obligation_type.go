// Package dbaccess 数据访问层
package dbaccess

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authorization/common"
	"Authorization/interfaces"
)

type obligationType struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	obligationTypeOnce    sync.Once
	obligationTypeService *obligationType
)

// NewObligationType 创建数据库对象
func NewObligationType() *obligationType {
	obligationTypeOnce.Do(func() {
		obligationTypeService = &obligationType{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})
	return obligationTypeService
}

func (o *obligationType) Set(ctx context.Context, info *interfaces.ObligationTypeInfo) (err error) {
	curTime := common.GetCurrentMicrosecondTimestamp()
	// info.Config 转成string
	configByte, err := json.Marshal(info.Schema)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}

	var defaultValueByte []byte
	if info.DefaultValue != nil {
		defaultValueByte, err = json.Marshal(info.DefaultValue)
		if err != nil {
			o.logger.Errorln(err)
			return err
		}
	}

	var uiSchemaByte []byte
	if info.UiSchema != nil {
		uiSchemaByte, err = json.Marshal(info.UiSchema)
		if err != nil {
			o.logger.Errorln(err)
			return err
		}
	}

	// info.ResourceTypeScope 转成string
	resourceTypeScopeStr, err := o.resourceTypeScopeInfoToString(info.ResourceTypeScope)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}

	// 先判断是否存在
	strSQL := "select f_id from " + common.GetDBName(databaseName) + ".t_obligation_type where f_id = ?"
	rows, err := o.db.Query(strSQL, info.ID)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}
		}
	}()
	if rows.Next() {
		strSQL = "update " + common.GetDBName(databaseName) +
			".t_obligation_type set f_name = ?, f_description = ?, f_schema = ?, f_default_value = ?, f_ui_schema = ?, f_applicable_resource_types = ?, f_modified_at = ? where f_id = ?"
		_, err = o.db.Exec(strSQL, info.Name, info.Description, string(configByte), string(defaultValueByte), string(uiSchemaByte), resourceTypeScopeStr, curTime, info.ID)
	} else {
		strSQL = "insert into " + common.GetDBName(databaseName) +
			".t_obligation_type(f_id, f_name, f_description, f_schema, f_default_value , f_ui_schema, f_applicable_resource_types, f_created_at, f_modified_at) values(?,?,?,?,?,?,?,?,?)"
		_, err = o.db.Exec(strSQL, info.ID, info.Name, info.Description, string(configByte), string(defaultValueByte), string(uiSchemaByte), resourceTypeScopeStr, curTime, curTime)
	}
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return
}

// resourceTypeScopeInfoToString 将 ResourceTypeScopeInfo 转为 str
func (o *obligationType) resourceTypeScopeInfoToString(info interfaces.ObligationResourceTypeScopeInfo) (resp string, err error) {
	result := make(map[string]any)
	// 如果不限制
	if info.Unlimited {
		result["unlimited"] = true
		result["resource_types"] = make([]any, 0)
	} else {
		result["unlimited"] = false
		resourceTypes := make([]any, 0, len(info.Types))
		for _, rt := range info.Types {
			resourceType := make(map[string]any)
			resourceType["id"] = rt.ResourceTypeID
			OperationsScope := make(map[string]any)
			if rt.OperationsScope.Unlimited {
				OperationsScope["unlimited"] = true
				OperationsScope["operations"] = make([]any, 0)
			} else {
				OperationsScope["unlimited"] = false
				var operations []any
				for _, op := range rt.OperationsScope.Operations {
					operations = append(operations, map[string]any{
						"id": op.ID,
					})
				}
				OperationsScope["operations"] = operations
			}
			resourceType["applicable_operations"] = OperationsScope
			resourceTypes = append(resourceTypes, resourceType)
		}
		result["resource_types"] = resourceTypes
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		o.logger.Errorf("json.Marshal: %v", err)
		return "", err
	}
	resp = string(jsonBytes)
	return
}

// stringToResourceTypeScopeInfo 将JSON字符串转为 ResourceTypeScopeInfo
func (o *obligationType) stringToResourceTypeScopeInfo(jsonStr string) (info interfaces.ObligationResourceTypeScopeInfo, err error) {
	var data map[string]any
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		o.logger.Errorf("json.Unmarshal: %v", err)
		return info, err
	}

	info.Unlimited = data["unlimited"].(bool)

	// 如果资源类型有范围限制，则需要设置资源类型范围
	if !info.Unlimited {
		resourceTypesJson := data["resource_types"].([]any)
		for _, resourceTypeJson := range resourceTypesJson {
			// 遍历资源类型，每个资源类型信息 放入 resourceTypeScope
			var resourceTypeScope interfaces.ObligationResourceTypeScope
			resourceTypeJsonMap := resourceTypeJson.(map[string]any)
			resourceTypeID := resourceTypeJsonMap["id"].(string)
			operationsScopeJson := resourceTypeJsonMap["applicable_operations"].(map[string]any)
			var operationsScopeInfo interfaces.ObligationOperationsScopeInfo
			operationsScopeInfo.Unlimited = operationsScopeJson["unlimited"].(bool)
			if !operationsScopeInfo.Unlimited {
				// 获取资源类型上的操作
				operationsJson := operationsScopeJson["operations"].([]any)
				for _, operationJson := range operationsJson {
					operationJsonMap := operationJson.(map[string]any)
					operationID := operationJsonMap["id"].(string)
					var operation interfaces.ObligationOperation
					operation.ID = operationID
					operationsScopeInfo.Operations = append(operationsScopeInfo.Operations, operation)
				}
			}
			resourceTypeScope.ResourceTypeID = resourceTypeID
			resourceTypeScope.OperationsScope = operationsScopeInfo
			// 将资源类型信息放入资源类型范围
			info.Types = append(info.Types, resourceTypeScope)
		}
	}

	return info, nil
}

func (o *obligationType) Delete(ctx context.Context, obligationTypeID string) (err error) {
	strSQL := "delete from " + common.GetDBName(databaseName) + ".t_obligation_type where f_id = ?"
	_, err = o.db.Exec(strSQL, obligationTypeID)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return
}

func (o *obligationType) GetByID(ctx context.Context, obligationTypeID string) (info interfaces.ObligationTypeInfo, err error) {
	strSQL := "select f_id, f_name, f_description, f_schema, f_default_value, f_ui_schema, f_applicable_resource_types, f_created_at, f_modified_at from " +
		common.GetDBName(databaseName) + ".t_obligation_type where f_id = ?"
	// 帮我写一个逻辑 生成 sql 语句
	var args []any
	args = append(args, obligationTypeID)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var schemaStr string
		var defaultValueStr string
		var uiSchemaStr string
		var resourceTypeScopeStr string
		err = rows.Scan(&info.ID, &info.Name, &info.Description, &schemaStr, &defaultValueStr, &uiSchemaStr, &resourceTypeScopeStr, &info.CreatedTime, &info.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return
		}
		// 解析Config
		err = json.Unmarshal([]byte(schemaStr), &info.Schema)
		if err != nil {
			o.logger.Errorln(err)
			return
		}

		// 解析DefaultValue
		if defaultValueStr != "" {
			err = json.Unmarshal([]byte(defaultValueStr), &info.DefaultValue)
			if err != nil {
				o.logger.Errorln(err)
				return
			}
		}

		// 解析UiSchema
		if uiSchemaStr != "" {
			err = json.Unmarshal([]byte(uiSchemaStr), &info.UiSchema)
			if err != nil {
				o.logger.Errorln(err)
				return
			}
		}

		// 解析ResourceTypeScope
		var resourceTypeScope interfaces.ObligationResourceTypeScopeInfo
		resourceTypeScope, err = o.stringToResourceTypeScopeInfo(resourceTypeScopeStr)
		if err != nil {
			o.logger.Errorln(err)
			return
		}
		info.ResourceTypeScope = resourceTypeScope
	}
	return
}

//nolint:gocyclo
func (o *obligationType) Get(ctx context.Context, info *interfaces.ObligationTypeSearchInfo) (count int, resultInfos []interfaces.ObligationTypeInfo, err error) {
	var countRows *sql.Rows
	countRows, err = o.db.Query("select count(1) from " + common.GetDBName(databaseName) + ".t_obligation_type")
	if err != nil {
		o.logger.Errorln(err)
		return 0, nil, err
	}

	for countRows.Next() {
		err = countRows.Scan(&count)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
	}

	if countRows != nil {
		if countRowsErr := countRows.Err(); countRowsErr != nil {
			o.logger.Errorln(countRowsErr)
		}
		if closeErr := countRows.Close(); closeErr != nil {
			o.logger.Errorln(closeErr)
		}
	}

	strSQL := "select f_id, f_name, f_description, f_schema, f_default_value, f_ui_schema, f_applicable_resource_types, f_created_at, f_modified_at from " + common.GetDBName(databaseName) +
		".t_obligation_type limit ? offset ?"
	// 帮我写一个逻辑 生成 sql 语句
	var args []any
	args = append(args, info.Limit, info.Offset)
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return 0, nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationTypeInfo
		var schemaStr string
		var defaultValueStr string
		var uiSchemaStr string
		var resourceTypeScopeStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.Name, &resultInfo.Description, &schemaStr, &defaultValueStr, &uiSchemaStr,
			&resourceTypeScopeStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
		// 解析Config
		err = json.Unmarshal([]byte(schemaStr), &resultInfo.Schema)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}

		// 解析DefaultValue
		if defaultValueStr != "" {
			err = json.Unmarshal([]byte(defaultValueStr), &resultInfo.DefaultValue)
			if err != nil {
				o.logger.Errorln(err)
				return 0, nil, err
			}
		}
		// 解析UiSchema
		if uiSchemaStr != "" {
			err = json.Unmarshal([]byte(uiSchemaStr), &resultInfo.UiSchema)
			if err != nil {
				o.logger.Errorln(err)
				return 0, nil, err
			}
		}
		// 解析ResourceTypeScope
		var resourceTypeScope interfaces.ObligationResourceTypeScopeInfo
		resourceTypeScope, err = o.stringToResourceTypeScopeInfo(resourceTypeScopeStr)
		if err != nil {
			o.logger.Errorln(err)
			return 0, nil, err
		}
		resultInfo.ResourceTypeScope = resourceTypeScope
		resultInfos = append(resultInfos, resultInfo)
	}
	return
}

// 获取所有义务类型
//
//nolint:dupl
func (o *obligationType) GetAll(ctx context.Context) (resultInfos []interfaces.ObligationTypeInfo, err error) {
	strSQL := "select f_id, f_name, f_description, f_schema, f_default_value, f_ui_schema, f_applicable_resource_types, f_created_at, f_modified_at from " + common.GetDBName(databaseName) +
		".t_obligation_type order by f_modified_at desc"
	// 帮我写一个逻辑 生成 sql 语句
	var args []any
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationTypeInfo
		var schemaStr string
		var defaultValueStr string
		var uiSchemaStr string
		var resourceTypeScopeStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.Name, &resultInfo.Description, &schemaStr, &defaultValueStr, &uiSchemaStr,
			&resourceTypeScopeStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		// 解析Config
		err = json.Unmarshal([]byte(schemaStr), &resultInfo.Schema)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}

		// 解析DefaultValue
		if defaultValueStr != "" {
			err = json.Unmarshal([]byte(defaultValueStr), &resultInfo.DefaultValue)
			if err != nil {
				o.logger.Errorln(err)
				return
			}
		}

		// 解析UiSchema
		if uiSchemaStr != "" {
			err = json.Unmarshal([]byte(uiSchemaStr), &resultInfo.UiSchema)
			if err != nil {
				o.logger.Errorln(err)
				return nil, err
			}
		}
		// 解析ResourceTypeScope
		var resourceTypeScope interfaces.ObligationResourceTypeScopeInfo
		resourceTypeScope, err = o.stringToResourceTypeScopeInfo(resourceTypeScopeStr)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		resultInfo.ResourceTypeScope = resourceTypeScope
		resultInfos = append(resultInfos, resultInfo)
	}
	return
}

// 通过ID批量获取义务类型
//
//nolint:dupl
func (o *obligationType) GetByIDs(ctx context.Context, ids []string) (infos []interfaces.ObligationTypeInfo, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	IDSet, IDGroup := getFindInSetSQL(ids)
	var args []any
	args = append(args, IDGroup...)
	strSQL := "select f_id, f_name, f_description, f_schema, f_default_value, f_ui_schema, f_applicable_resource_types, f_created_at, f_modified_at from " + common.GetDBName(databaseName) +
		".t_obligation_type  where f_id in (" + IDSet + ")"
	rows, err := o.db.Query(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
		return nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()
	for rows.Next() {
		var resultInfo interfaces.ObligationTypeInfo
		var schemaStr string
		var defaultValueStr string
		var uiSchemaStr string
		var resourceTypeScopeStr string
		err = rows.Scan(&resultInfo.ID, &resultInfo.Name, &resultInfo.Description, &schemaStr, &defaultValueStr, &uiSchemaStr,
			&resourceTypeScopeStr, &resultInfo.CreatedTime, &resultInfo.ModifyTime)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		// 解析Config
		err = json.Unmarshal([]byte(schemaStr), &resultInfo.Schema)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}

		// 解析DefaultValue
		if defaultValueStr != "" {
			err = json.Unmarshal([]byte(defaultValueStr), &resultInfo.DefaultValue)
			if err != nil {
				o.logger.Errorln(err)
				return
			}
		}
		// 解析UiSchema
		if uiSchemaStr != "" {
			err = json.Unmarshal([]byte(uiSchemaStr), &resultInfo.UiSchema)
			if err != nil {
				o.logger.Errorln(err)
				return nil, err
			}
		}
		// 解析ResourceTypeScope
		var resourceTypeScope interfaces.ObligationResourceTypeScopeInfo
		resourceTypeScope, err = o.stringToResourceTypeScopeInfo(resourceTypeScopeStr)
		if err != nil {
			o.logger.Errorln(err)
			return nil, err
		}
		resultInfo.ResourceTypeScope = resourceTypeScope
		infos = append(infos, resultInfo)
	}
	return
}
