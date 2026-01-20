// Package common 配置服务模块
package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

// SvcConfig 服务配置信息
var SvcConfig Config

// Config 配置文件
type Config struct {
	Lang                  string `yaml:"lang"`
	SvcHost               string `yaml:"svc_host"`
	SvcPublicPort         int    `yaml:"svc_public_port"`
	SvcPrivatePort        int    `yaml:"svc_private_port"`
	OAuthAdminHost        string `yaml:"oauth_admin_host"`
	OAuthAdminPort        int    `yaml:"oauth_admin_port"`
	OAuthPublicHost       string `yaml:"oauth_public_host"`
	OAuthPublicPort       int    `yaml:"oauth_public_port"`
	BusinessTimeOffset    int64  `yaml:"business_time_offset"`
	OSSGateWayPrivateHost string `yaml:"oss_gateway_private_host"`
	OSSGateWayPrivatePort int    `yaml:"oss_gateway_private_port"`
	CleanAvatarOffsetTime int    `yaml:"clean_avatar_offset_time"`
	LogLevel              int    `yaml:"log_level"`
	SystemID              string `yaml:"system_id"`
}

// InitConfig 读取服务配置
func InitConfig() {
	configLog := NewLogger()
	file, err := os.ReadFile("/sysvol/conf/user-management.yaml")
	if err != nil {
		configLog.Fatalln("load /sysvol/conf/user-management.yaml failed:", err)
	}

	err = yaml.Unmarshal(file, &SvcConfig)
	if err != nil {
		configLog.Fatalln("unmarshal yaml file failed:", err)
	}

	configLog.Infoln(SvcConfig)
}
