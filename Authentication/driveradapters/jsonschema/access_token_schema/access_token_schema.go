// Package accesstokenschema jsonschema定义层
package accesstokenschema

import (
	_ "embed" // 标准用法
)

var (
	// AccessTokenSchemaStr 申请访问令牌schema str
	//go:embed access_token_schema.json
	AccessTokenSchemaStr string
)
