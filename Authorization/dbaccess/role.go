// Package dbaccess group Anyshare 数据访问层 - 用户组数据库操作
package dbaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authorization/common"
	"Authorization/interfaces"
)

type role struct {
	db      *sqlx.DB
	logger  common.Logger
	dbTrace *sqlx.DB
}

var (
	roleOnce sync.Once
	roleDB   *role
)

// NewRole 创建数据库角色
func NewRole() *role {
	roleOnce.Do(func() {
		roleDB = &role{
			db:      dbPool,
			logger:  common.NewLogger(),
			dbTrace: dbTracePool,
		}
	})

	return roleDB
}

// resourceTypeScopeInfoToString 将 ResourceTypeScopeInfo 转为 str
func (r *role) resourceTypeScopeInfoToString(info interfaces.ResourceTypeScopeInfo) (resp string, err error) {
	result := make(map[string]any)
	result["unlimited"] = info.Unlimited
	types := make([]map[string]any, 0, len(info.Types))
	for _, t := range info.Types {
		typeMap := map[string]any{
			"id": t.ResourceTypeID,
		}
		types = append(types, typeMap)
	}
	result["types"] = types
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		r.logger.Errorf("json.Marshal: %v", err)
		return "", err
	}
	resp = string(jsonBytes)
	return
}

// resourceTypeScopeInfoFromString 将字符串转为 ResourceTypeScopeInfo
func (r *role) resourceTypeScopeInfoFromString(resourceTypeScopeStr string) (info interfaces.ResourceTypeScopeInfo, err error) {
	var result map[string]any
	err = json.Unmarshal([]byte(resourceTypeScopeStr), &result)
	if err != nil {
		r.logger.Errorf("json.Unmarshal: %v", err)
		return info, err
	}

	// 解析 unlimited 字段
	if unlimitedVal, ok := result["unlimited"]; ok {
		info.Unlimited = unlimitedVal.(bool)
	}

	// 解析 types 字段
	if typesVal, ok := result["types"]; ok {
		if typesArray, ok := typesVal.([]any); ok {
			info.Types = make([]interfaces.ResourceTypeScope, 0, len(typesArray))
			for _, typeVal := range typesArray {
				if typeMap, ok := typeVal.(map[string]any); ok {
					scope := interfaces.ResourceTypeScope{}
					scope.ResourceTypeID = typeMap["id"].(string)
					info.Types = append(info.Types, scope)
				}
			}
		}
	}

	return info, nil
}

// AddRole 添加角色
func (r *role) AddRoles(ctx context.Context, roles []interfaces.RoleInfo) (err error) {
	if len(roles) == 0 {
		return
	}
	currentTime := common.GetCurrentMicrosecondTimestamp()
	type tmpRoleInfo struct {
		ID                 string
		Name               string
		Description        string
		RoleSource         interfaces.RoleSource
		ResourceTypeScopes string
	}

	tmpRoles := []tmpRoleInfo{}
	for _, role := range roles {
		var resourceTypeScopeStr string
		resourceTypeScopeStr, err = r.resourceTypeScopeInfoToString(role.ResourceTypeScopeInfo)
		if err != nil {
			return err
		}
		tmpRoles = append(tmpRoles, tmpRoleInfo{
			ID:                 role.ID,
			Name:               role.Name,
			Description:        role.Description,
			RoleSource:         role.RoleSource,
			ResourceTypeScopes: resourceTypeScopeStr,
		})
	}

	var valuesStr []string
	var inserts []any
	// 批量插入
	// 数据库中 f_visibility 字段已经弃用， 直接设置成0值
	visibility := 0
	for i := range tmpRoles {
		valuesStr = append(valuesStr, "(?, ?, ?, ?, ?, ?, ?,?)")
		inserts = append(inserts, tmpRoles[i].ID, tmpRoles[i].Name, tmpRoles[i].RoleSource,
			tmpRoles[i].ResourceTypeScopes, tmpRoles[i].Description, visibility, currentTime, currentTime)
	}
	valueStr := strings.Join(valuesStr, ",")

	strSQL := "insert into " + common.GetDBName(databaseName) +
		".t_role(f_id, f_name, f_source, f_resource_scope, f_description, f_visibility, f_created_time, f_modify_time) values " + valueStr

	_, err = r.db.Exec(strSQL, inserts...)
	if err != nil {
		r.logger.Errorf("sql: %s, err: %v", strSQL, err)
		return err
	}
	return nil
}

// DeleteRole 删除角色
func (r *role) DeleteRole(ctx context.Context, id string) (err error) {
	dbName := common.GetDBName(databaseName)
	sqlStr := "delete from %s.t_role where f_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := r.db.Exec(sqlStr, id); err != nil {
		return err
	}
	return nil
}

// ModifyRole 修改角色
func (r *role) ModifyRole(ctx context.Context, id, name string, nameChanged bool, description string,
	descriptionChanged bool, resourceTypeScopes interfaces.ResourceTypeScopeInfo, resourceTypeScopesChanged bool,
) (err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	sqlStr := "update %s.t_role set "
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if nameChanged {
		args = append(args, name)
		sqlStr += "f_name = ?, "
	}

	if descriptionChanged {
		args = append(args, description)
		sqlStr += "f_description = ?, "
	}

	if resourceTypeScopesChanged {
		resourceTypeScopeStr, err := r.resourceTypeScopeInfoToString(resourceTypeScopes)
		if err != nil {
			return err
		}
		args = append(args, resourceTypeScopeStr)
		sqlStr += "f_resource_scope = ?, "
	}
	currentTime := common.GetCurrentMicrosecondTimestamp()
	sqlStr += "f_modify_time = ? where f_id = ? "
	args = append(args, currentTime, id)
	if _, err := r.db.Exec(sqlStr, args...); err != nil {
		r.logger.Errorf("err: %v, sqlStr: %s, args: %v", err, sqlStr, args)
		return err
	}
	return nil
}

// GetRoleByID 获取指定的角色
func (r *role) GetRoleByID(ctx context.Context, id string) (info interfaces.RoleInfo, err error) {
	dbName := common.GetDBName(databaseName)
	strSQL := "select f_id, f_name,  f_source, f_resource_scope, f_description, f_created_time, f_modify_time from %s.t_role where f_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.db.Query(strSQL, id)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return info, sqlErr
	}

	for rows.Next() {
		var resourceTypeScopeStr string
		if scanErr := rows.Scan(&info.ID, &info.Name, &info.RoleSource, &resourceTypeScopeStr, &info.Description, &info.CreateTime, &info.ModifyTime); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return info, scanErr
		}
		// Parse resourceTypeScopeStr to ResourceTypeScopeInfo
		info.ResourceTypeScopeInfo, err = r.resourceTypeScopeInfoFromString(resourceTypeScopeStr)
		if err != nil {
			r.logger.Errorln(err, strSQL)
			return info, err
		}
	}

	return info, nil
}

// GetRoleByName 根据名称获取指定的角色
func (r *role) GetRoleByName(ctx context.Context, name string) (info interfaces.RoleInfo, err error) {
	dbName := common.GetDBName(databaseName)
	strSQL := "select f_id, f_name, f_resource_scope, f_description, f_created_time, f_modify_time from %s.t_role where f_name = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.db.Query(strSQL, name)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return info, sqlErr
	}

	for rows.Next() {
		var resourceTypeScopeStr string
		if scanErr := rows.Scan(&info.ID, &info.Name, &resourceTypeScopeStr, &info.Description, &info.CreateTime, &info.ModifyTime); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return info, scanErr
		}
		info.ResourceTypeScopeInfo, err = r.resourceTypeScopeInfoFromString(resourceTypeScopeStr)
		if err != nil {
			r.logger.Errorln(err, strSQL)
			return info, err
		}
	}

	return info, nil
}

// GetRoleByIDs 获取指定的角色
func (r *role) GetRoleByIDs(ctx context.Context, ids []string) (infoMap map[string]interfaces.RoleInfo, err error) {
	infoMap = make(map[string]interfaces.RoleInfo)
	if len(ids) == 0 {
		return infoMap, nil
	}
	dbName := common.GetDBName(databaseName)
	idsSet, idsGroup := getFindInSetSQL(ids)
	strSQL := "select f_id, f_name, f_source, f_resource_scope, f_description, f_created_time, f_modify_time from %s.t_role where f_id in (" + idsSet + ")"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.db.Query(strSQL, idsGroup...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	for rows.Next() {
		var resourceTypeScopeStr string
		var info interfaces.RoleInfo
		if scanErr := rows.Scan(&info.ID, &info.Name, &info.RoleSource, &resourceTypeScopeStr, &info.Description, &info.CreateTime, &info.ModifyTime); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		// Parse resourceTypeScopeStr to ResourceTypeScopeInfo
		info.ResourceTypeScopeInfo, err = r.resourceTypeScopeInfoFromString(resourceTypeScopeStr)
		if err != nil {
			r.logger.Errorln(err, strSQL)
			return infoMap, err
		}
		infoMap[info.ID] = info
	}
	return infoMap, nil
}

const (
	dataAdminRoleID        = "00990824-4bf7-11f0-8fa7-865d5643e61f"
	aiAdminRoleID          = "3fb94948-5169-11f0-b662-3a7bdba2913f"
	applicationAdminRoleID = "1572fb82-526f-11f0-bde6-e674ec8dde71"
)

/*
GetRoles 列举符合条件的角色
排序逻辑：
1. 首先按 f_source 升序排序（系统角色 -> 业务角色 -> 用户角色）
2. 当 f_source = 2（业务内置角色）时，按照预定义的内置角色顺序排序：
  - dataAdminRoleID（数据管理员）
  - aiAdminRoleID（AI管理员）
  - applicationAdminRoleID（应用管理员）

3. 最后按修改时间倒序、主键倒序排序
*/
func (r *role) GetRoles(ctx context.Context, info interfaces.RoleSearchInfo) (roles []interfaces.RoleInfo, err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	// 如果角色来源为空，则返回空
	if len(info.RoleSources) == 0 {
		return
	}

	set := make([]string, 0)
	for _, v := range info.RoleSources {
		set = append(set, "?")
		args = append(args, v)
	}
	setStr := strings.Join(set, ",")
	whereClause := ""
	if info.Keyword != "" {
		whereClause = " and f_name like ?"
		args = append(args, "%"+info.Keyword+"%")
	}

	var pinnedRoleIDs []string
	// 构建业务内置角色的排序CASE语句
	var businessRoleOrder string
	for _, roleSource := range info.RoleSources {
		if roleSource == interfaces.RoleSourceBusiness {
			pinnedRoleIDs = []string{
				dataAdminRoleID,
				aiAdminRoleID,
				applicationAdminRoleID,
			}
			businessRoleOrder = "case f_id "
			for i, id := range pinnedRoleIDs {
				businessRoleOrder += fmt.Sprintf("when '%s' then %d ", id, i)
			}
			businessRoleOrder += fmt.Sprintf("else %d end", len(pinnedRoleIDs))
		}
	}

	args = append(args, info.Offset, info.Limit)

	// 构建排序逻辑：先按f_source升序，当f_source=2时按内置角色顺序排序，然后按修改时间倒序
	var orderByClause string
	if businessRoleOrder != "" {
		orderByClause = fmt.Sprintf("f_source asc, case when f_source = 2 then %s else %d end, f_modify_time desc, f_primary_id desc", businessRoleOrder, len(pinnedRoleIDs))
	} else {
		orderByClause = "f_source asc, f_modify_time desc, f_primary_id desc"
	}

	strSQL := "select f_id, f_name,  f_source, f_resource_scope, f_description, f_created_time, f_modify_time from %s.t_role where f_source in (" + setStr + ") " + whereClause +
		" order by " + orderByClause + " limit ?, ? "
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := r.db.Query(strSQL, args...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	roles = make([]interfaces.RoleInfo, 0)
	for rows.Next() {
		var resourceTypeScopeStr string
		var tmpInfo interfaces.RoleInfo
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name, &tmpInfo.RoleSource,
			&resourceTypeScopeStr, &tmpInfo.Description, &tmpInfo.CreateTime, &tmpInfo.ModifyTime); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		// Parse resourceTypeScopeStr to ResourceTypeScopeInfo
		tmpInfo.ResourceTypeScopeInfo, err = r.resourceTypeScopeInfoFromString(resourceTypeScopeStr)
		if err != nil {
			r.logger.Errorln(err, strSQL)
			return nil, err
		}

		roles = append(roles, tmpInfo)
	}

	return roles, nil
}

// GetRolesSum 列举符合条件的角色数量
func (r *role) GetRolesSum(ctx context.Context, info interfaces.RoleSearchInfo) (num int, err error) {
	if len(info.RoleSources) == 0 {
		return 0, nil
	}
	dbName := common.GetDBName(databaseName)
	whereClause := ""
	args := make([]any, 0)

	set := make([]string, 0)
	for _, v := range info.RoleSources {
		set = append(set, "?")
		args = append(args, v)
	}
	setStr := strings.Join(set, ",")
	if info.Keyword != "" {
		whereClause = " and f_name like ?"
		args = append(args, "%"+info.Keyword+"%")
	}

	strSQL := "select count(*) from %s.t_role where f_source in (" + setStr + ") " + whereClause
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := r.db.Query(strSQL, args...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return num, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&num); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return num, scanErr
		}
	}

	return num, nil
}

// SetRoleByID 根据ID修改角色
func (r *role) SetRoleByID(ctx context.Context, id string, role *interfaces.RoleInfo) (err error) {
	dbName := common.GetDBName(databaseName)
	currentTime := common.GetCurrentMicrosecondTimestamp()
	sqlStr := "update %s.t_role set f_name = ?, f_description = ?, f_resource_scope = ?, f_modify_time = ? where f_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	resourceTypeScopeStr, err := r.resourceTypeScopeInfoToString(role.ResourceTypeScopeInfo)
	if err != nil {
		return err
	}
	args := []any{role.Name, role.Description, resourceTypeScopeStr, currentTime, id}
	_, err = r.db.Exec(sqlStr, args...)
	if err != nil {
		r.logger.Errorln(err, sqlStr)
		return err
	}
	return err
}

// GetAllUserRolesInternal 列举所有用户创建的角色
func (r *role) GetAllUserRolesInternal(ctx context.Context, keyword string) (roles []interfaces.RoleInfo, err error) {
	dbName := common.GetDBName(databaseName)
	args := make([]any, 0)
	args = append(args, interfaces.RoleSourceUser)
	filterStr := ""
	if keyword != "" {
		filterStr = " and f_name like ?"
		args = append(args, "%"+keyword+"%")
	}
	strSQL := "select f_id, f_name, f_resource_scope, f_description, f_created_time, f_modify_time from %s.t_role where f_source = ?" + filterStr
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := r.db.Query(strSQL, args...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	roles = make([]interfaces.RoleInfo, 0)
	for rows.Next() {
		var resourceTypeScopeStr string
		var tmpInfo interfaces.RoleInfo
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name, &resourceTypeScopeStr, &tmpInfo.Description, &tmpInfo.CreateTime, &tmpInfo.ModifyTime); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		// Parse resourceTypeScopeStr to ResourceTypeScopeInfo
		tmpInfo.ResourceTypeScopeInfo, err = r.resourceTypeScopeInfoFromString(resourceTypeScopeStr)
		if err != nil {
			r.logger.Errorln(err, strSQL)
			return nil, err
		}

		roles = append(roles, tmpInfo)
	}
	return roles, nil
}
