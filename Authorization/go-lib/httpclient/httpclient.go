// Package httpclient HTTP客户端服务接口
package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	errorv2 "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/log"
	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"go.opentelemetry.io/otel/attribute"
)

// HTTPClient HTTP客户端服务接口
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (respParam interface{}, err error)
	Post(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respParam interface{}, err error)
	// respHeaders带有返回报头
	PostEx(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respHeaders map[string]string, respParam interface{}, err error)
	Put(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respParam interface{}, err error)
	Delete(ctx context.Context, url string, headers map[string]string) (respParam interface{}, err error)
}

var (
	rawOnce   sync.Once
	rawClient *http.Client
	httpOnce  sync.Once
	client    HTTPClient
)

// httpClient HTTP客户端结构
type httpClient struct {
	client *http.Client
	trace  observable.Tracer
}

// NewRawHTTPClient 创建原生HTTP客户端对象
func NewRawHTTPClient() *http.Client {
	rawOnce.Do(func() {
		rawClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConnsPerHost:   100,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 10 * time.Second,
		}
	})

	return rawClient
}

// NewHTTPClient 创建HTTP客户端对象
func NewHTTPClient(trace observable.Tracer) HTTPClient {
	httpOnce.Do(func() {
		client = &httpClient{
			client: NewRawHTTPClient(),
			trace:  trace,
		}
	})

	return client
}

// NewHTTPClientEx 创建HTTP客户端对象, 自定义超时时间
func NewHTTPClientEx(timeout time.Duration, trace observable.Tracer) HTTPClient {
	return &httpClient{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConnsPerHost:   100,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: timeout * time.Second,
		},
		trace: trace,
	}
}

// Get http client get
func (c *httpClient) Get(ctx context.Context, url string, headers map[string]string) (respParam interface{}, err error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return
	}

	_, respParam, err = c.httpDo(ctx, req, headers)
	return
}

// Post http client post
func (c *httpClient) Post(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respParam interface{}, err error) {
	var reqBody []byte
	if v, ok := reqParam.([]byte); ok {
		reqBody = v
	} else {
		reqBody, err = jsoniter.Marshal(reqParam)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	respCode, respParam, err = c.httpDo(ctx, req, headers)
	return
}

// Put http client put
func (c *httpClient) Put(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respParam interface{}, err error) {
	reqBody, err := jsoniter.Marshal(reqParam)
	if err != nil {
		return
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	respCode, respParam, err = c.httpDo(ctx, req, headers)
	return
}

// Delete http client delete
func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) (respParam interface{}, err error) {
	req, err := http.NewRequest("DELETE", url, http.NoBody)
	if err != nil {
		return
	}

	_, respParam, err = c.httpDo(ctx, req, headers)
	return
}

func (c *httpClient) httpDoCommon(ctx context.Context, req *http.Request, headers map[string]string) (respCode int, respParam interface{}, resp *http.Response, err error) {
	if c.client == nil {
		return 0, nil, nil, errors.New("http client is unavailable")
	}
	c.trace.SetClientSpanName("http " + req.Method + " 接口调用")
	ctx, span := c.trace.AddClientTrace(ctx)
	defer func() {
		span.SetAttributes(attribute.String("request.method", req.Method))
		span.SetAttributes(attribute.String("request.url", req.URL.String()))
		// var respBody map[string]interface{}
		if resp != nil {
			span.SetAttributes(attribute.String("response.statusCode", fmt.Sprintf("%v", resp.StatusCode)))
			if respParam != nil {
				if str, ok := respParam.(string); ok {
					span.SetAttributes(attribute.String("response.body", str))
				} else {
					span.SetAttributes(attribute.String("response.body", fmt.Sprintf("%v", respParam)))
				}
			}
		}
		c.trace.TelemetrySpanEnd(span, err)
	}()

	c.addHeaders(req, headers)

	resp, err = c.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.NewLogger().Errorln(closeErr)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respCode = resp.StatusCode
	if (respCode < http.StatusOK) || (respCode >= http.StatusMultipleChoices) {
		data := make(map[string]interface{})
		err = jsoniter.Unmarshal(body, &data)
		if err != nil {
			// Unmarshal失败时转成内部错误, body为空Unmarshal失败
			err = fmt.Errorf("code:%v,header:%v,body:%v", respCode, resp.Header, string(body))
			return
		}
		err = c.parseError(data)
		return
	}

	if len(body) != 0 {
		err = jsoniter.Unmarshal(body, &respParam)
	}

	return
}

func (c *httpClient) httpDo(ctx context.Context, req *http.Request, headers map[string]string) (respCode int, respParam interface{}, err error) {
	respCode, respParam, _, err = c.httpDoCommon(ctx, req, headers)
	return
}

func (c *httpClient) addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		if len(v) > 0 {
			req.Header.Add(k, v)
		}
	}
}

// 返回 header
func (c *httpClient) httpDoV2(ctx context.Context, req *http.Request, headers map[string]string) (respCode int, respHeaders map[string]string, respParam interface{}, err error) {
	respCode, respParam, resp, err := c.httpDoCommon(ctx, req, headers)
	if err != nil {
		return
	}
	respHeaders = map[string]string{
		"x-request-id": resp.Header.Get("x-request-id"),
	}
	return
}

func (c *httpClient) parseError(data map[string]interface{}) (err error) {
	var description, solution string
	var detail map[string]interface{}

	if val, ok := data["description"].(string); ok {
		description = val
	}

	if val, ok := data["solution"].(string); ok {
		solution = val
	}

	if val, ok := data["detail"].(map[string]interface{}); ok {
		detail = val
	}

	switch code := data["code"].(type) {
	case string:
		var link string
		if val, ok := data["link"].(string); ok {
			link = val
		}
		err = &errorv2.Error{
			Code:        code,
			Description: description,
			Solution:    solution,
			Detail:      detail,
			Link:        link,
		}
	case float64:
		var cause, message string
		if val, ok := data["cause"].(string); ok {
			cause = val
		}
		if val, ok := data["message"].(string); ok {
			message = val
		}
		err = &rest.HTTPError{
			Cause:       cause,
			Code:        int(code),
			Message:     message,
			Detail:      detail,
			Description: description,
			Solution:    solution,
		}
	}
	return err
}

// Post http client post
func (c *httpClient) PostEx(ctx context.Context, url string, headers map[string]string, reqParam interface{}) (respCode int, respHeaders map[string]string, respParam interface{}, err error) {
	var reqBody []byte
	switch v := reqParam.(type) {
	case []byte:
		reqBody = v
	case string:
		reqBody = []byte(v)
	default:
		reqBody, err = jsoniter.Marshal(reqParam)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	respCode, respHeaders, respParam, err = c.httpDoV2(ctx, req, headers)
	return
}
