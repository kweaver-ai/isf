// Package authschema jsonschema定义层
package authschema

import (
	_ "embed" // 标准用法
)

var (
	// PwdAuthSchemaStr 用户账号密码登录schema str
	//go:embed pwd_auth_schema.json
	PwdAuthSchemaStr string
)
