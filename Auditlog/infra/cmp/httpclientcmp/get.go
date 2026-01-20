package httpclientcmp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
)

func (c *httpClient) Get(ctx context.Context, url string, queryData ...interface{}) (resp *gclient.Response, err error) {
	ctx, span := c.arTrace.AddClientTrace(ctx)
	defer func() {
		c.arTraceRecord(span, resp, nil, err)
	}()

	debugReqLog(debugReqLogger{
		URL:    url,
		Data:   queryData,
		Method: http.MethodGet,
	})

	return c.client.Retry(3, time.Second*RetryInterval).Get(ctx, url, queryData...)
}

func (c *httpClient) GetExpect2xx(ctx context.Context, url string, queryData ...interface{}) (resp string, err error) {
	resByte, err := c.GetExpect2xxByte(ctx, url, queryData...)
	if err != nil {
		return
	}

	resp = string(resByte)

	return
}

func (c *httpClient) GetExpect2xxByte(ctx context.Context, url string, queryData ...interface{}) (resp []byte, err error) {
	var (
		r          *gclient.Response
		requestErr error
	)

	ctx, span := c.arTrace.AddClientTrace(ctx)
	defer func() {
		c.arTraceRecord(span, r, err, requestErr)

		if requestErr != nil {
			err = requestErr
		}
	}()

	debugReqLog(debugReqLogger{
		URL:    url,
		Data:   queryData,
		Method: http.MethodGet,
	})

	r, requestErr = c.client.Retry(3, time.Second*RetryInterval).Get(ctx, url, queryData...)

	if requestErr != nil {
		// todo logger 从外面传进来
		log.Printf("[GetExpect2xx] request error: %v\n", requestErr)
		return
	}

	defer func(r *gclient.Response) {
		_ = r.Close()
	}(r)

	err = c.errExpect2xx(r)
	if err != nil {
		return
	}

	resp = r.ReadAll()

	return
}
