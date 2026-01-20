package hlartrace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_trace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"AuditLog/common"
	"AuditLog/gocommon/api"
)

var traceOnce sync.Once

// NewARTrace 内部函数调用、数据库调用和依赖服务调用需要链路数据埋点
func NewARTrace() Tracer {
	// 初始化一个实例
	traceOnce.Do(func() {
		serviceName := os.Getenv("SERVICE_NAME")
		ar_trace.InitTracer("cm", "anyshare-telemetry-sdk", serviceName)
	})

	return &arTrace{
		logger: common.SvcConfig.Logger,
	}
}

type arTrace struct {
	logger           api.Logger
	clientSpanName   string
	internalSpanName string
}

func (at *arTrace) SetClientSpanName(spanName string) {
	at.clientSpanName = spanName
}

func (at *arTrace) SetInternalSpanName(spanName string) {
	at.internalSpanName = spanName
}

// AddInternalTrace 内部方法调用链路埋点
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

//AddClientTrace 依赖服务调用链路埋点
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

// AddServerTrace 接口层调用时使用
func (at *arTrace) AddServerTrace(c *gin.Context) (ctx context.Context, span trace.Span) {
	if strings.Contains(c.FullPath(), "health/ready") {
		return
	}

	newCtx := c.Request.Context()
	for k, v := range c.Keys {
		//nolint:staticcheck
		newCtx = context.WithValue(newCtx, k, v)
	}

	ctx = otel.GetTextMapPropagator().Extract(newCtx, propagation.HeaderCarrier(c.Request.Header))
	ctx, span = ar_trace.Tracer.Start(ctx, c.FullPath(), trace.WithSpanKind(trace.SpanKindServer))

	span.SetAttributes(attribute.String("http.method", c.Request.Method))
	span.SetAttributes(attribute.String("http.route", c.FullPath()))
	span.SetAttributes(attribute.String("http.client_ip", c.ClientIP()))

	return ctx, span
}

// AddConsumerTrace 消费者消费消息时记录使用
/*
@description 消费者消费消息时记录使用
@param     ctx           context.Context         链路上下文
@param     topic         string                  消费的topic
@return    newCtx        context.Context         新的链路上下文
@return    span          trace.Span              a unit of work or operation
*/
func (at *arTrace) AddConsumerTrace(ctx context.Context, topic string) (newCtx context.Context, span trace.Span) {
	_, file, linkNo, ok := runtime.Caller(1)
	if !ok {
		at.logger.Warnf("start consumer span error")

		newCtx, span = ar_trace.Tracer.Start(ctx, "unknow", trace.WithSpanKind(trace.SpanKindConsumer))

		return
	}

	spanName := "mqConsumer: " + topic
	newCtx, span = ar_trace.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))
	span.SetAttributes(attribute.String("func.path", fmt.Sprintf("%s:%v", file, linkNo)))

	return
}

// AddProducerTrace 生产者生产消息时记录使用
func (at *arTrace) AddProducerTrace(ctx context.Context) (newCtx context.Context, span trace.Span) {
	pc, file, linkNo, ok := runtime.Caller(1)
	if !ok {
		at.logger.Warnf("start producer span error")

		newCtx, span := ar_trace.Tracer.Start(ctx, "unknow", trace.WithSpanKind(trace.SpanKindProducer))

		return newCtx, span
	}

	funcPaths := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	spanName := funcPaths[len(funcPaths)-1]
	newCtx, span = ar_trace.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindProducer))
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

// TelemetrySpanEndIgnoreDBNotFound gorm span end ignore gorm.ErrRecordNotFound
// 注意：使用这个要注意，要根据需要来。如果不清楚是否需要，那么就不要使用这个。
// 使用这个需要注意使用的位置
func (at *arTrace) TelemetrySpanEndIgnoreDBNotFound(span trace.Span, err error) {
	if span == nil {
		return
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	at.TelemetrySpanEnd(span, err)
}
