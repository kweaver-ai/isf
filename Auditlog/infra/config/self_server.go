package config

import (
	"crypto/tls"
	"net/http"

	"AuditLog/common"
)

// SelfServerConfig 本服务配置
type SelfServerConfig struct {
	HttpClientInsecureSkipVerify bool `yaml:"httpClientInsecureSkipVerify"`
}

func (l *SelfServerConfig) loadConf() {
	err := common.Configure(l, "self_server_config.yaml")
	if err != nil {
		panic(err)
	}
}

func SetTpTlsInsecureSkipVerify(tp *http.Transport) {
	if GetSelfServerConfig().HttpClientInsecureSkipVerify {
		if tp.TLSClientConfig == nil {
			tp.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			tp.TLSClientConfig.InsecureSkipVerify = true
		}
	}
}
