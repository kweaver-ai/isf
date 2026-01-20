package http

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

var (
	clientOnce      sync.Once
	clientSingleton *client
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type client struct {
	client *http.Client
}

// NewClient 初始化管理实例 // TODO in env
func NewClient() Client {
	clientOnce.Do(func() {
		clientSingleton = &client{client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConnsPerHost:   100,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 10 * time.Second, // TODO in env
		}}
	})
	return clientSingleton
}

// Get Just Get
func (c *client) Get(url string) (resp *http.Response, err error) {
	return c.client.Get(url)
}
