package httpclientcmp

import (
	"net/http"

	"AuditLog/common/enums"
	"AuditLog/common/utils"
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/icmp"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/pkg/errors"
)

type Option func(c *httpClient)

type httpClient struct {
	token   string
	client  *gclient.Client
	arTrace api.Tracer
}

var _ icmp.IHttpClient = &httpClient{}

func NewHTTPClient(arTrace api.Tracer, opts ...Option) icmp.IHttpClient {
	gClient := GetNewGClientWithDefaultStdClient()

	c := &httpClient{
		token:   "",
		client:  gClient,
		arTrace: arTrace,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithToken(token string) Option {
	return func(c *httpClient) {
		if token == "" {
			return
		}

		c.token = token
		c.client.SetHeader("Authorization", "Bearer "+token)
	}
}

func WithHeader(k, v string) Option {
	return func(c *httpClient) {
		c.client.SetHeader(k, v)
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *httpClient) {
		for k, v := range headers {
			c.client.SetHeader(k, v)
		}
	}
}

func WithClient(client *gclient.Client) Option {
	return func(c *httpClient) {
		c.client = client
	}
}

func (c *httpClient) errExpect2xx(r *gclient.Response) (err error) {
	//nolint:nestif
	if utils.IsHttpErr(r.Response) {
		resp := &CommonRespError{}

		body := r.ReadAll()
		err = utils.JSON().Unmarshal(body, &resp)

		if err != nil {
			prefixLen := 30
			_url := r.Request.URL

			if len(body) > prefixLen {
				err = errors.Errorf("httpClient failed(not 2xx), url: [%s], http code: [%v], body[0:%d] is [%q]",
					_url, r.StatusCode, prefixLen, string(body)[:prefixLen])
			} else {
				err = errors.Errorf("httpClient failed(not 2xx), url: [%s], http code: [%v], body is [%q]", _url, r.StatusCode, string(body))
			}
		} else {
			err = errors.Wrap(resp, "httpClient failed(not 2xx), response")
		}

		debugResLog(debugResLogger{
			Err:      err,
			RespBody: body,
		})
	}

	return
}

func (c *httpClient) errExpect2xxStd(r *http.Response) (err error) {
	if utils.IsHttpErr(r) {
		err = errors.Errorf("httpClient failed(not 2xx), url: [%s], http code: [%v]", r.Request.URL, r.StatusCode)
	}

	return
}

func (c *httpClient) setContentType(contentType string) {
	c.client.SetHeader(enums.HTTPHct, contentType)
}

func (c *httpClient) GetClient() *gclient.Client {
	return c.client
}
