// Package authschema jsonschema定义层
package authschema

import (
	_ "embed" // 标准用法
)

var (
	// ClientLoginSchemaStr 注册客户端schema str
	//go:embed client_auth_schema.json
	ClientLoginSchemaStr string
)
