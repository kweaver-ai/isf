// Package ticketschema jsonschema定义层
package ticketschema

import (
	_ "embed" // 标准用法
)

var (
	// TicketSchemaStr 生成单点凭据schema str
	//go:embed ticket_schema.json
	TicketSchemaStr string
)
