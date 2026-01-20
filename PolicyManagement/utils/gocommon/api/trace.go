package api

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_trace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	traceOnce sync.Once
	at        *arTrace
)

type Tracer interface {
	AddInternalTrace(ctx context.Context) (newCtx context.Context, span trace.Span)
	SetInternalSpanName(spanName string)
	TelemetrySpanEnd(span trace.Span, err error)
}

// NewARTrace 内部函数调用、数据库调用和依赖服务调用需要链路数据埋点
func NewARTrace() *arTrace {
	// 初始化一个实例
	traceOnce.Do(func() {
		serviceName := os.Getenv("SERVICE_NAME")
		ar_trace.InitTracer("cm", "anyshare-telemetry-sdk", serviceName)
	})
	at = &arTrace{
		logger: NewTelemetryLogger(),
	}

	return at
}

type arTrace struct {
	logger           Logger
	clientSpanName   string
	internalSpanName string
}

func (at *arTrace) SetClientSpanName(spanName string) {
	at.clientSpanName = spanName
}

func (at *arTrace) SetInternalSpanName(spanName string) {
	at.internalSpanName = spanName
}

/*
@description 内部方法调用链路埋点
@param     ctx           context.Context         链路上下文
@return    newCtx        context.Context         新的链路上下文
@return    span          trace.Span              a unit of work or operation
*/
func (at *arTrace) AddInternalTrace(ctx context.Context) (newCtx context.Context, span trace.Span) {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	pc, file, linkNo, ok := runtime.Caller(1)
	if !ok {
		at.logger.Warnf("start internal span error")
		newCtx, span = ar_trace.Tracer.Start(ctx, "unknow", trace.WithSpanKind(trace.SpanKindInternal))
		return
	}

	var spanName string
	if at.internalSpanName == "" {
		funcPaths := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		spanName = funcPaths[len(funcPaths)-1]
	} else {
		spanName = at.internalSpanName
	}

	newCtx, span = ar_trace.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal))

	span.SetAttributes(attribute.String("func.path", fmt.Sprintf("%s:%v", file, linkNo)))

	return
}

/*
@description 依赖服务调用链路埋点
@param     ctx           context.Context         链路上下文
@return    newCtx        context.Context         新的链路上下文
@return    span          trace.Span              a unit of work or operation
*/
func (at *arTrace) AddClientTrace(ctx context.Context) (newCtx context.Context, span trace.Span) {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	pc, file, linkNo, ok := runtime.Caller(1)
	if !ok {
		at.logger.Warnf("start client span error")
		newCtx, span = ar_trace.Tracer.Start(ctx, "unknow", trace.WithSpanKind(trace.SpanKindClient))
		return
	}

	var spanName string
	if at.clientSpanName == "" {
		funcPaths := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		spanName = funcPaths[len(funcPaths)-1]
	} else {
		spanName = at.clientSpanName
	}
	newCtx, span = ar_trace.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))

	span.SetAttributes(attribute.String("func.path", fmt.Sprintf("%s:%v", file, linkNo)))

	return
}

// TelemetrySpanEnd 关闭span
func (at *arTrace) TelemetrySpanEnd(span trace.Span, err error) {
	if span == nil {
		return
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}
	span.End()
}
