// Package telemetryschema jsonschema定义层
package telemetryschema

import (
	_ "embed" // 标准用法
)

var (
	// TelemetryLogSchema 发送可观测性日志schema
	//go:embed telemetry_log.json
	TelemetryLogSchema string
)
