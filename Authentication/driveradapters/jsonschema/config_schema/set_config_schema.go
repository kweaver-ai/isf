// Package configschema jsonschema定义层
package configschema

import (
	_ "embed" // 标准用法
)

var (
	// SetConfigSchemaStr 更新配置schema str
	//go:embed set_config.json
	SetConfigSchemaStr string
)
