package config

// TODO pkg common

import (
	jsoniter "github.com/json-iterator/go"
)

var Config config

type config struct {
	selfConfig
	serviceAccessConfig
}

func (c config) String() string {
	pretty, err := jsoniter.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(pretty)
}

func InitConfig() {
	Config = config{
		selfConfig:          initSelfConfig(),
		serviceAccessConfig: initServiceAccess(),
	}
}
