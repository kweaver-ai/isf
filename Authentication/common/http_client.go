// Package common HTTP客户端服务接口
package common

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

var (
	rawOnce2   sync.Once
	rawClient2 *http.Client
)

// NewRawHTTPClient2 创建没有超时限制的原生HTTP客户端对象
func NewRawHTTPClient2() *http.Client {
	rawOnce2.Do(func() {
		rawClient2 = &http.Client{
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
		}
	})

	return rawClient2
}
