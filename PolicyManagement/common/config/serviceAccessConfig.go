package config

import (
	"os"
	"strconv"

	cdb "policy_mgnt/utils/gocommon/v2/db"
	cservice "policy_mgnt/utils/gocommon/v2/service"
	cutils "policy_mgnt/utils/gocommon/v2/utils"
)

// serviceAccessConfig 服务访问
type serviceAccessConfig struct {
	ProtonPolicyEngine      cservice.Access
	UserMgmtPvt             cservice.Access
	DB                      cdb.Config
	Redis                   RedisConfig
	ProtonApplicationConfig cservice.Access
	HydraAdminConfig        cservice.Access
	Language                string
}

// RedisConfig 配置信息
type RedisConfig struct {
	ConnectType string // sentinel/standalone/master-slave 对应哨兵、单机、主从三种连接方式
	ConnectInfo RedisConnectInfo
	EnableSSL   bool
	SecretName  string // 当 enableSSL 为 true 时需要
	CaName      string // 当 enableSSL 为 true 时需要，表示secret里 ca 证书的名字
	CertName    string // 当 enableSSL 为 true 时需要，表示secret里 cert 证书的名字
	KeyName     string // 当 enableSSL 为 true 时需要，表示secret里 key 密钥的名字
}

// RedisConnectInfo 配置信息
type RedisConnectInfo struct {
	Username         string
	Password         string
	Host             string
	Port             int
	MasterHost       string
	MasterPort       int
	SlaveHost        string
	SlavePort        int
	SentinelHost     string
	SentinelPort     int
	SentinelUsername string
	SentinelPassword string
	MasterGroupName  string
}

func initServiceAccess() serviceAccessConfig {
	return serviceAccessConfig{
		ProtonPolicyEngine: cservice.Access{
			Protocol: cutils.GetEnv("PROTON_POLICY_ENGINE_PROTOCOL", "http"),
			Host:     cutils.GetEnv("PROTON_POLICY_ENGINE_HOST", "proton-policy-engine-proton-policy-engine-cluster.resource.svc.cluster.local"),
			Port:     cutils.GetEnv("PROTON_POLICY_ENGINE_PORT", "9800"),
		},
		UserMgmtPvt: cservice.Access{
			Protocol: cutils.GetEnv("USER_MANAGEMENT_PROTOCOL", "http"),
			Host:     cutils.GetEnv("USER_MANAGEMENT_HOST", "user-management-private.anyshare.svc.cluster.local"),
			Port:     cutils.GetEnv("USER_MANAGEMENT_PORT", "30980"),
		},
		DB: cdb.Config{
			DBType:          cutils.GetEnv("DB_TYPE", "mysql"),
			Host:            cutils.GetEnv("DB_HOST", "mariadb-mariadb-cluster.resource.svc.cluster.local"),
			Port:            cutils.GetEnv("DB_PORT", "3330"),
			Driver:          cutils.GetEnv("DB_DRIVER", "mysql"),
			DataBaseName:    cutils.GetEnv("DB_NAME", "policy_mgnt"),
			User:            cutils.GetEnv("DB_USER", ""),
			Password:        getDBPassword(),
			Timezone:        cutils.GetEnv("TIMEZONE", "Asia/Shanghai"),
			MaxIdleConns:    cutils.GetEnv("DB_MAX_IDLE_CONNS", "2"),
			MaxOpenConns:    cutils.GetEnv("DB_MAX_OPEN_CONNS", "0"),
			ConnMaxLifetime: cutils.GetEnv("DB_CONN_MAX_LIFE_TIME", "0m"), // time.Duration
		},
		Redis: RedisConfig{
			ConnectType: cutils.GetEnv("REDIS_CONNECT_TYPE", "sentinel"),
			ConnectInfo: RedisConnectInfo{
				Username:         cutils.GetEnv("REDIS_USERNAME", "root"),
				Password:         cutils.GetEnv("REDIS_PASSWORD", ""),
				Host:             cutils.GetEnv("REDIS_HOST", ""),
				Port:             mustGetIntEnv("REDIS_PORT", 0),
				MasterHost:       cutils.GetEnv("REDIS_MASTER_HOST", ""),
				MasterPort:       mustGetIntEnv("REDIS_MASTER_PORT", 0),
				SlaveHost:        cutils.GetEnv("REDIS_SLAVE_HOST", ""),
				SlavePort:        mustGetIntEnv("REDIS_SLAVE_PORT", 0),
				SentinelHost:     cutils.GetEnv("REDIS_SENTINEL_HOST", "proton-redis-proton-redis-sentinel.resource"),
				SentinelPort:     mustGetIntEnv("REDIS_SENTINEL_PORT", 26379),
				SentinelUsername: cutils.GetEnv("REDIS_SENTINEL_USERNAME", "root"),
				SentinelPassword: cutils.GetEnv("REDIS_SENTINEL_PASSWORD", ""),
				MasterGroupName:  cutils.GetEnv("REDIS_MASTER_GROUP_NAME", "mymaster"),
			},
			EnableSSL:  getRedisEnableSSL(),
			SecretName: cutils.GetEnv("REDIS_SECRET_NAME", ""),
			CaName:     cutils.GetEnv("REDIS_CA_NAME", ""),
			CertName:   cutils.GetEnv("REDIS_CERT_NAME", ""),
			KeyName:    cutils.GetEnv("REDIS_KEY_NAME", ""),
		},
		ProtonApplicationConfig: cservice.Access{
			Host:     cutils.GetEnv("PROTON_APPLICATION_HOST", "deploy-web-service"),
			Port:     cutils.GetEnv("PROTON_APPLICATION_PORT", "18880"),
			Protocol: cutils.GetEnv("PROTON_APPLICATION_PROTOCOL", "http"),
		},
		HydraAdminConfig: cservice.Access{
			Protocol: cutils.GetEnv("HYDRA_ADMIN_PROTOCOL", "http"),
			Host:     cutils.GetEnv("HYDRA_ADMIN_HOST", "hydra-hydra-cluster.resource.svc.cluster.local"),
			Port:     cutils.GetEnv("HYDRA_ADMIN_PORT", "4444"),
		},
		Language: cutils.GetEnv("LANGUAGE", "zh_CN"),
	}
}

func getRedisEnableSSL() bool {
	return cutils.GetEnv("REDIS_ENABLE_SSL", "false") == "true"
}

func mustGetIntEnv(envName string, defaultV int) int {
	v, ok := os.LookupEnv(envName)
	if !ok {
		return defaultV
	}
	intV, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return intV
}
