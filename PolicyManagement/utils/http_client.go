package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	oauthcli "golang.org/x/oauth2"
	credentmgnt "golang.org/x/oauth2/clientcredentials"
)

// Management httpclient管理接口
type Management interface {
	Do(*http.Request) (resp *http.Response, err error)
}

type mgnt struct {
	client *http.Client
}

// NewManagement 初始化管理实例
func NewManagement() (Management, error) {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost:   100,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr}
	ctx := context.WithValue(context.Background(), oauthcli.HTTPClient, client)
	conf := &credentmgnt.Config{
		ClientID:     os.Getenv("SERVICE_CLIENT_ID"),
		ClientSecret: os.Getenv("SERVICE_CLIENT_SECRET"),
		Scopes:       []string{},
		TokenURL:     getTokenEndpoint(),
	}
	oauthHTTPClient := conf.Client(ctx)
	return &mgnt{client: oauthHTTPClient}, nil
}

func (m *mgnt) Do(request *http.Request) (resp *http.Response, err error) {
	resp, err = m.client.Do(request)
	return
}

func getHydraPublicURL() url.URL {
	schema := os.Getenv("HYDRA_PUBLIC_PROTOCOL")
	host := os.Getenv("HYDRA_PUBLIC_HOST")
	port := os.Getenv("HYDRA_PUBLIC_PORT")
	url := url.URL{
		Scheme: schema,
		Host:   fmt.Sprintf("%s:%s", host, port),
	}
	return url
}

func getTokenEndpoint() string {
	url := getHydraPublicURL()
	url.Path = "/oauth2/token"
	return url.String()
}
