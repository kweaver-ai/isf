// Package sessionsschema jsonschema定义层
package sessionsschema

import (
	_ "embed" // 标准用法
)

var (
	// HydraSessionsSchemaStr 注册客户端schema str
	//go:embed sessions_schema.json
	HydraSessionsSchemaStr string
)
