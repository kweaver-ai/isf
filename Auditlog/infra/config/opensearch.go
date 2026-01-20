package config

import "AuditLog/common"

type OpenSearchConf struct {
	PrivateHost     string `yaml:"private_host"`
	PrivatePort     string `yaml:"private_port"`
	PrivateProtocol string `yaml:"private_protocol"`

	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (o *OpenSearchConf) loadConf() {
	err := common.Configure(o, "opensearch.yaml")
	if err != nil {
		panic(err)
	}
}

func (o *OpenSearchConf) GetAddress() string {
	return o.PrivateProtocol + "://" + o.PrivateHost + ":" + o.PrivatePort
}
