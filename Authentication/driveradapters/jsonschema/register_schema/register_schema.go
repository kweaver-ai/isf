// Package registerschema jsonschema定义层
package registerschema

import (
	_ "embed" // 标准用法
)

var (
	// RegisterSchema 注册客户端schema
	//go:embed register.json
	RegisterSchema string
)
