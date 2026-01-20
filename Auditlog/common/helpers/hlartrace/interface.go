package hlartrace

import (
	"go.opentelemetry.io/otel/trace"

	"AuditLog/gocommon/api"
)

//go:generate mockgen -source=./interface.go -destination ./hlartracemock/artrace_mock.go -package hlartracemock
type Tracer interface {
	api.Tracer

	TelemetrySpanEndIgnoreDBNotFound(span trace.Span, err error)
}
