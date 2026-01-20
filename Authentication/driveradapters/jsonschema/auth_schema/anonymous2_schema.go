// Package authschema jsonschema定义层
package authschema

import (
	_ "embed" // 标准用法
)

var (
	// AnonymousLogin2 匿名登录schema str
	//go:embed anonymous2_schema.json
	AnonymousLogin2 string
)
