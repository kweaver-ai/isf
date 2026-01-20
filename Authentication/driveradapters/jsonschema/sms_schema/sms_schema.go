// Package smsschema jsonschema定义层
package smsschema

import (
	_ "embed" // 标准用法
)

var (
	// SMSSchemaStr 发送匿名账户短信验证码 schema str
	//go:embed sms_schema.json
	SMSSchemaStr string
)
