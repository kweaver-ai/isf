package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// mgnt httpclient结构体
type mgnt struct {
	client        *http.Client
	retryCount    int           //重试次数
	retryInterval time.Duration //等待间隔单位毫秒
	trace         *arTrace
}

const (
	KeepAliveTime = 60 * time.Second //连接时间
)

func NewHttpClient() Client {
	// TODO: InsecureSkipVerify 需要可以通过参数指定
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   500 * time.Millisecond, // 连接超时时间
				KeepAlive: KeepAliveTime,          // 保持长连接的时间
			}).DialContext, // 设置连接的参数
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
			MaxIdleConns:          100,              // 最大空闲连接
			IdleConnTimeout:       90 * time.Second, // 空闲连接的超时时间
			TLSHandshakeTimeout:   10 * time.Second, //限制 TLS握手的时间
			ExpectContinueTimeout: 1 * time.Second,  //限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待。
		},
		Timeout: KeepAliveTime,
	}
	trace := NewARTrace()
	trace.SetClientSpanName("http接口调用")
	return &mgnt{client: client, retryCount: 4, retryInterval: 200 * time.Millisecond, trace: trace}
}

func (m *mgnt) Do(ctx context.Context, request *http.Request) (resp *http.Response, err error) {
	_, span := m.trace.AddClientTrace(ctx)
	defer func() {
		span.SetAttributes(attribute.String("request.method", request.Method))
		span.SetAttributes(attribute.String("request.url", request.URL.String()))
		// var respBody map[string]interface{}
		if resp != nil {
			span.SetAttributes(attribute.String("response.statusCode", fmt.Sprintf("%v", resp.StatusCode)))
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 0 {
				// _ = json.Unmarshal(body, &respBody)
				span.SetAttributes(attribute.String("response.body", string(body)))
			}
		}
		m.trace.TelemetrySpanEnd(span, err)
	}()

	for i := 0; i < m.retryCount; i++ {
		resp, err = m.client.Do(request)
		if err == nil {
			break
		}
		netError, _ := err.(net.Error)
		if !netError.Timeout() || i+1 == m.retryCount {
			break
		}
		fmt.Printf("url: %s, %s request %d retry %d times\n", request.URL, time.Now().String(), m.retryCount, i+1)
		time.Sleep(m.retryInterval)
	}
	return resp, err
}
