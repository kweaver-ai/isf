// Package common 配置服务模块
package common

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"

	"github.com/kweaver-ai/go-lib/util"
)

// SvcConfig 服务配置信息
var SvcConfig Config

// Config 配置文件
type Config struct {
	Lang                string `yaml:"lang"`
	GinMode             string `yaml:"gin_mode"`
	SvcHost             string `yaml:"svc_host"`
	SvcPublicPort       int    `yaml:"svc_public_port"`
	SvcPrivatePort      int    `yaml:"svc_private_port"`
	SystemID            string `yaml:"system_id"`
	UserMgntPrivateHost string `yaml:"user_management_private_host"`
	UserMgntPrivatePort int    `yaml:"user_management_private_port"`
	OAuthAdminHost      string `yaml:"oauth_admin_host"`
	OAuthAdminPort      int    `yaml:"oauth_admin_port"`
	BusinessTimeOffset  int64  `yaml:"business_time_offset"`
	LogLevel            string `yaml:"log_level"`

	Redis RedisConfig `yaml:"redis"`

	// 以下从环境变量中获取

	MQConfigFilePath string
	PodIP            string
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
	file, err := os.ReadFile("/sysvol/conf/service_conf/authorization.yaml")
	if err != nil {
		configLog.Fatalln("load /sysvol/conf/service_conf/authorization failed:", err)
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

	SvcConfig.GinMode = getGinMode()

	SvcConfig.Redis.ConnectInfo.Host = util.ParseHost(SvcConfig.Redis.ConnectInfo.Host)
	SvcConfig.Redis.ConnectInfo.MasterHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.MasterHost)
	SvcConfig.Redis.ConnectInfo.SlaveHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.SlaveHost)
	SvcConfig.Redis.ConnectInfo.SentinelHost = util.ParseHost(SvcConfig.Redis.ConnectInfo.SentinelHost)

	SvcConfig.MQConfigFilePath = mustGetEnv("MQ_CONFIG_FILE_PATH")
	SvcConfig.PodIP = mustGetEnv("POD_IP")

	svcConfig := SvcConfig
	svcConfig.Redis.ConnectInfo.Password = ""
	svcConfig.Redis.ConnectInfo.SentinelPassword = ""

	configLog.Infoln(svcConfig)
}

// 限制不能使用test模式 https://github.com/gin-gonic/gin/issues/2984
func getGinMode() string {
	if SvcConfig.GinMode != gin.ReleaseMode && SvcConfig.GinMode != gin.DebugMode {
		panic("invalid gin mode: " + SvcConfig.GinMode)
	}
	return SvcConfig.GinMode
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("The environment variable '%s' must be set.", key))
	}
	return v
}
