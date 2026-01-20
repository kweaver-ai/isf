// Package usermanagementchema jsonschema定义层
package usermanagementchema

import (
	_ "embed" // 标准用法
)

var (
	// UserDeleteSchemaStr 用户删除schema str
	//go:embed user_delete_schema.json
	UserDeleteSchemaStr string
	// UserPasswordModifySchemaStr 用户密码更改schema str
	//go:embed user_password_modify_schema.json
	UserPasswordModifySchemaStr string
	// UserStatusChangeSchemaStr 用户状态变更schema str
	//go:embed user_status_change_schema.json
	UserStatusChangeSchemaStr string
)
