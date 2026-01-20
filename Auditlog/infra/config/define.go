package config

import (
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/common"
)

var (
	config     *Config
	configOnce sync.Once
)

type Config struct {
	dbConfig         *sqlx.DBConfig
	um               *UmConf
	logConf          *LogConf
	EFast            *EFastConf `yaml:"efast"`
	OpenSearch       *OpenSearchConf
	Rec              *RecConf
	ShareMgnt        *ShareMgntConf  `yaml:"sharemgnt"`
	DocCenter        *DocCenterConf  `yaml:"doc-center"`
	OssGateway       *OssGatewayConf `yaml:"ossgateway"`
	logDumpConfig    LogDumpConfigRepo
	Kc               *KcConf             `yaml:"kc"`
	selfServerConfig *SelfServerConfig   `yaml:"self_server_config"`
	PersonalConfig   *PersonalConfigConf `yaml:"personal-config"`

	redisConfig *RedisConfig
}

func (c *Config) loadDevSvcConf() {
	err := common.Configure(c, "dep_svc.yaml")
	if err != nil {
		panic(err)
	}
}
