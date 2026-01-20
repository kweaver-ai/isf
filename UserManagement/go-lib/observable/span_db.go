package observable

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ctxSpanKey struct{}

// RDSHook rds hook
type RDSHook struct {
	System string
}

func generateSQL(sqlstr string, args ...interface{}) string {
	var str string
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			str = fmt.Sprint(v)
		case string:
			str = fmt.Sprintf("'%s'", v)
		case time.Time:
			str = v.Format("2006-01-02 15:04:05")
		default:
			f, err := strconv.ParseFloat(fmt.Sprint(v), 64)
			if err == nil {
				str = fmt.Sprint(f)
			} else {
				str = fmt.Sprintf("'%s'", fmt.Sprint(v))
			}
		}
		sqlstr = strings.Replace(sqlstr, "?", str, 1)
	}
	return sqlstr
}

// Before 开始
func (h *RDSHook) Before(ctx context.Context, sqltmp string, args ...interface{}) (context.Context, error) {
	sqlstr := generateSQL(sqltmp, args...)
	if ar_trace.Tracer != nil {
		newCtx, span := ar_trace.Tracer.Start(ctx, "sql", trace.WithSpanKind(trace.SpanKindInternal))
		span.SetAttributes(attribute.Key("db.system").String(h.System))
		span.SetAttributes(attribute.Key("db.statement").String(sqlstr))
		newCtx = context.WithValue(newCtx, (*ctxSpanKey)(nil), span)
		return newCtx, nil
	}
	return ctx, nil
}

// After 结束
func (h *RDSHook) After(ctx context.Context, sqltmp string, args ...interface{}) (context.Context, error) {
	if _, ok := ctx.Value((*ctxSpanKey)(nil)).(trace.Span); ok {
		ar_trace.EndSpan(ctx, nil)
	}
	return ctx, nil
}

// OnError 错误
func (h *RDSHook) OnError(ctx context.Context, err error, sqltmp string, args ...interface{}) error {
	if err != nil && !errors.Is(err, driver.ErrSkip) {
		if span, ok := ctx.Value((*ctxSpanKey)(nil)).(trace.Span); ok {
			span.SetAttributes(attribute.Key("db.error").String(err.Error()))
			ar_trace.EndSpan(ctx, err)
		}
	}
	return nil
}
