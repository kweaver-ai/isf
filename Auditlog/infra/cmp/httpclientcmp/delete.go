package httpclientcmp

import (
	"context"
	"log"
	"time"

	"github.com/gogf/gf/v2/net/gclient"
)

func (c *httpClient) Delete(ctx context.Context, url string) (resp *gclient.Response, err error) {
	return c.client.Retry(3, time.Second*RetryInterval).Delete(ctx, url)
}

func (c *httpClient) DeleteExpect2xx(ctx context.Context, url string) (resp string, err error) {
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

	r, requestErr = c.client.Retry(3, time.Second*RetryInterval).Delete(ctx, url)

	if requestErr != nil {
		// todo logger 从外面传进来
		log.Printf("[DeleteExpect2xx] request error: %v\n", requestErr)
		return
	}

	defer func(r *gclient.Response) {
		_ = r.Close()
	}(r)

	err = c.errExpect2xx(r)
	if err != nil {
		return
	}

	resp = r.ReadAllString()

	return
}
