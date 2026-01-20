package config

import (
	"os"

	clog "policy_mgnt/utils/gocommon/v2/log"
	cutils "policy_mgnt/utils/gocommon/v2/utils"
)

// selfConfig 本服务配置
type selfConfig struct {
	Host      string
	PodIP     string
	LogConfig clog.Config
}

func getDBPassword() string {
	passwd, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("The environment variable 'DB_PASSWORD' must be set.")
	}
	return passwd
}

func getPodIP() string {
	podIP, ok := os.LookupEnv("POD_IP")
	if !ok {
		panic("The environment variable 'POD_IP' must be set.")
	}
	return podIP
}

func initSelfConfig() selfConfig {
	return selfConfig{
		Host:  cutils.GetEnv("HOST", "0.0.0.0"),
		PodIP: getPodIP(),
	}
}
