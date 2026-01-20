package httpclientcmp

import (
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/net/gclient"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *httpClient) arTraceRecord(span trace.Span, r *gclient.Response, err error, requestErr error) {
	if requestErr != nil {
		c.arTrace.TelemetrySpanEnd(span, requestErr)
		return
	}

	if r != nil {
		req := r.Request
		resp := r.Response

		span.SetAttributes(attribute.String("request.method", req.Method))
		span.SetAttributes(attribute.String("request.url", req.URL.String()))

		if resp != nil {
			span.SetAttributes(attribute.String("response.statusCode", fmt.Sprintf("%v", resp.StatusCode)))
		}
	}

	c.arTrace.TelemetrySpanEnd(span, err)
}

func (c *httpClient) arTraceRecord2(span trace.Span, req *http.Request, resp *http.Response, err error, requestErr error) {
	if requestErr != nil {
		c.arTrace.TelemetrySpanEnd(span, requestErr)
		return
	}

	span.SetAttributes(attribute.String("request.method", req.Method))
	span.SetAttributes(attribute.String("request.url", req.URL.String()))

	if resp != nil {
		span.SetAttributes(attribute.String("response.statusCode", fmt.Sprintf("%v", resp.StatusCode)))
	}

	c.arTrace.TelemetrySpanEnd(span, err)
}
