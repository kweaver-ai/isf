// Package common authentication服务配置
package common

import (
	"os"
	"time"

	"github.com/kweaver-ai/go-lib/util"
	"gopkg.in/yaml.v2"
)

// SvcConfig 服务配置信息
var SvcConfig Config

// Config 配置文件
type Config struct {
	Lang                      string      `yaml:"lang"`
	SvcHost                   string      `yaml:"svc_host"`
	SvcPublicPort             int         `yaml:"svc_public_port"`
	SvcPrivatePort            int         `yaml:"svc_private_port"`
	DBClientPingSleep         int         `yaml:"db_client_ping_sleep"`
	MaxOpenConns              int         `yaml:"max_open_conns"`
	OAuthPublicHost           string      `yaml:"oauth_public_host"`
	OAuthPublicPort           int         `yaml:"oauth_public_port"`
	OAuthAdminHost            string      `yaml:"oauth_admin_host"`
	OAuthAdminPort            int         `yaml:"oauth_admin_port"`
	UserManagementPrivateHost string      `yaml:"user_management_private_host"`
	UserManagementPrivatePort int         `yaml:"user_management_private_port"`
	EacpPrivateHost           string      `yaml:"eacp_private_host"`
	EacpPrivatePort           int         `yaml:"eacp_private_port"`
	ShareMgntHost             string      `yaml:"sharemgnt_host"`
	ShareMgntPort             int         `yaml:"sharemgnt_port"`
	BusinessTimeOffset        int64       `yaml:"business_time_offset"`
	FlowExpiredTime           int64       `yaml:"flow_expired_time"`
	FlowCleanTime             string      `yaml:"flow_clean_time"`
	LogLevel                  int         `yaml:"log_level"`
	SystemID                  string      `yaml:"system_id"`
	Redis                     RedisConfig `yaml:"redis"`
}

// RedisConfig 配置信息
type RedisConfig struct {
	ConnectType string           `yaml:"connectType"` // sentinel/standalone/master-slave 对应哨兵、单机、主从三种连接方式
	ConnectInfo RedisConnectInfo `yaml:"connectInfo"`
	EnableSSL   bool             `yaml:"enableSSL"`
	SecretName  string           `yaml:"secretName"` // 当 enableSSL 为 true 时需要
	CaName      string           `yaml:"caName"`     // 当 enableSSL 为 true 时需要，表示secret里 ca 证书的名字
	CertName    string           `yaml:"certName"`   // 当 enableSSL 为 true 时需要，表示secret里 cert 证书的名字
	KeyName     string           `yaml:"keyName"`    // 当 enableSSL 为 true 时需要，表示secret里 key 密钥的名字
}

// RedisConnectInfo 配置信息
type RedisConnectInfo struct {
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	MasterHost       string `yaml:"masterHost"`
	MasterPort       int    `yaml:"masterPort"`
	SlaveHost        string `yaml:"slaveHost"`
	SlavePort        int    `yaml:"slavePort"`
	SentinelHost     string `yaml:"sentinelHost"`
	SentinelPort     int    `yaml:"sentinelPort"`
	SentinelUsername string `yaml:"sentinelUsername"`
	SentinelPassword string `yaml:"sentinelPassword"`
	MasterGroupName  string `yaml:"masterGroupName"`
}

// InitConfig 读取服务配置
func InitConfig() {
	configLog := NewLogger()
	file, err := os.ReadFile("/sysvol/conf/service_conf/authentication.yaml")
	if err != nil {
		configLog.Fatalln("load /sysvol/conf/service_conf/authentication.yaml failed:", err)
	}

	err = yaml.Unmarshal(file, &SvcConfig)
	if err != nil {
		configLog.Fatalln("unmarshal yaml file failed:", err)
	}
	secretFile, err := os.ReadFile("/sysvol/conf/secret_conf/secret.yaml")
	if err != nil {
		configLog.Fatalf("load /sysvol/conf/secret_conf/secret.yaml failed: %v\n", err)
	}

	if err = yaml.Unmarshal(secretFile, &SvcConfig); err != nil {
		configLog.Fatalf("unmarshal yaml secretFile failed: %v\n", err)
	}

	SvcConfig.Redis.ConnectInfo.Host = util.ParseHost(SvcConfig.Redis.ConnectInfo.Host)
	SvcConfig.Redis.ConnectInfo.MasterHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.MasterHost)
	SvcConfig.Redis.ConnectInfo.SlaveHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.SlaveHost)
	SvcConfig.Redis.ConnectInfo.SentinelHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.SentinelHost)

	svcConfig := SvcConfig
	svcConfig.Redis.ConnectInfo.Password = ""
	svcConfig.Redis.ConnectInfo.SentinelPassword = ""

	configLog.Infoln(svcConfig)

	// 配置检测主要处理两种情况
	// 1. 升级时如果只升级镜像，未设置flow_expired_time和flow_clean_time，会导致清理最新的flow，登录失败
	// 2. 回退时如果只回退了chart，镜像没变会导致登陆失败，且定时任务一直执行，因为没有flow_clean_time和flow_expired_time

	// 检查过期时间有效性
	if SvcConfig.FlowExpiredTime <= 0 {
		configLog.Fatalln("flow_expired_time must be greater than 0")
	}

	// 检查清理时间有效性
	_, err = time.Parse("15:04:05", SvcConfig.FlowCleanTime)
	if err != nil {
		configLog.Fatalln("flow_clean_time must be set, error:", err)
	}
}
