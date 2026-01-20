package config

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/common"
)

func InitConfig() *Config {
	configOnce.Do(func() {
		// 1. 加载数据库配置
		dbConfig := &sqlx.DBConfig{}

		err := common.Configure(dbConfig, "dbrw.yaml")
		if err != nil {
			panic(err)
		}

		dbConfig.Database = common.SvcConfig.SystemID + dbConfig.Database

		// 2. 加载用户管理配置
		um := &UmConf{}
		um.loadConf()

		// 3. 初始化log配置
		logConf := &LogConf{}
		logConf.loadConf()

		// 4. 初始化openSearch配置
		openSearch := &OpenSearchConf{}
		openSearch.loadConf()

		// 5. 加载推荐配置
		recConf := &RecConf{}
		recConf.loadConf()

		// 6. 初始化日志转存配置
		logDumpConfig := loadLogDumpConfig()

		// 8. 初始化redis配置
		redisConfig := &RedisConfig{}

		err = common.Configure(redisConfig, "redis.yaml")
		if err != nil {
			panic(err)
		}

		// 10. 初始化本服务配置
		selfServerConfig := &SelfServerConfig{}
		selfServerConfig.loadConf()

		// 12. 初始化配置
		config = &Config{
			dbConfig:         dbConfig,
			um:               um,
			logConf:          logConf,
			OpenSearch:       openSearch,
			Rec:              recConf,
			logDumpConfig:    logDumpConfig,
			redisConfig:      redisConfig,
			selfServerConfig: selfServerConfig,
		}

		// 13. 加载依赖服务配置
		config.loadDevSvcConf()
	})

	return config
}

// GetConfig 获取配置
func GetConfig() *Config {
	if config == nil {
		panic("[GetConfig]: config is nil")
	}

	return config
}

// GetDBConfig 获取数据库配置
func GetDBConfig() *sqlx.DBConfig {
	if config == nil {
		panic("[GetDBConfig]: dbConfig is nil")
	}

	return config.dbConfig
}

// GetUserManagementConfig 获取用户管理配置
func GetUserManagementConfig() *UmConf {
	if config == nil {
		panic("[GetUserManagementConfig]: config is nil")
	}

	return config.um
}

// GetEFastConfig 获取EFast配置
func GetEFastConfig() *EFastConf {
	if config == nil {
		panic("[GetEFastConfig]: config is nil")
	}

	return config.EFast
}

func GetLogConf() *LogConf {
	if config == nil {
		panic("[GetLogConf]: config is nil")
	}

	return config.logConf
}

// GetLogDumpConfig 获取日志策略配置
func GetLogDumpConfig() LogDumpConfigRepo {
	if config == nil {
		panic("[GetLogDumpConfig]: config is nil")
	}

	return config.logDumpConfig
}

// GetShareMgntConf 获取共享管理配置
func GetShareMgntConf() *ShareMgntConf {
	if config == nil {
		panic("[GetShareMgntConf]: config is nil")
	}

	return config.ShareMgnt
}

// GetDocCenterConf 获取文档中心配置
func GetDocCenterConf() *DocCenterConf {
	if config == nil {
		panic("[GetDocCenterConf]: config is nil")
	}

	return config.DocCenter
}

// GetOssGatewayConf 获取对象存储配置
func GetOssGatewayConf() *OssGatewayConf {
	if config == nil {
		panic("[GetOssGatewayConf]: config is nil")
	}

	return config.OssGateway
}

func GetRedisConfig() *RedisConfig {
	if config == nil {
		panic("[GetRedisConfig]: config is nil")
	}

	return config.redisConfig
}

func GetKcConf() *KcConf {
	if config == nil {
		panic("[GetKcConf]: config is nil")
	}

	return config.Kc
}

func GetSelfServerConfig() *SelfServerConfig {
	if config == nil {
		panic("[GetSelfServerConfig]: config is nil")
	}

	return config.selfServerConfig
}

func GetPersonalConfigConf() *PersonalConfigConf {
	if config == nil {
		panic("[GetPersonalConfigConf]: config is nil")
	}

	return config.PersonalConfig
}
