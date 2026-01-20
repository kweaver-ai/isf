package cmphelper

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/httpclientcmp"
	"AuditLog/infra/cmp/icmp"
	"AuditLog/infra/config"
)

func GetClientWithTimeout(timeout time.Duration, arTrace api.Tracer, opts ...httpclientcmp.Option) (c icmp.IHttpClient) {
	tran := httpclientcmp.GetDefaultTp()

	config.SetTpTlsInsecureSkipVerify(tran)

	client := &http.Client{
		Transport: tran,
		Timeout:   timeout,
	}

	gClient := g.Client()
	gClient.Client = *client
	opt := httpclientcmp.WithClient(gClient)

	// 注意这个顺序，先设置client，再设置其他option
	opts = append([]httpclientcmp.Option{opt}, opts...)

	c = httpclientcmp.NewHTTPClient(arTrace, opts...)

	return
}

func GetClient(arTrace api.Tracer, opts ...httpclientcmp.Option) (c icmp.IHttpClient) {
	c = GetClientWithTimeout(httpclientcmp.DefaultTimeout, arTrace, opts...)

	return
}
