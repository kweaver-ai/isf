// Package auditschema jsonschema定义层
package auditschema

import (
	_ "embed" // 标准用法
)

var (
	// AuditLogSchema 发送审计日志schema
	//go:embed audit_log.json
	AuditLogSchema string

	//go:embed unordered_audit_log.json
	UnorderedAuditLogSchema string
)
