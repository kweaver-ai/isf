package utils

import (
	"github.com/spf13/viper"
)

// Visitor ACL配置
type Visitor struct {
	Name     string `mapstructure:"name"`
	ClientID string `mapstructure:"client_id"`
}

// InitConfig 初始化配置
func InitConfig() (err error) {
	SetDefaultConfig()
	return nil
}

// SetDefaultConfig 默认配置
func SetDefaultConfig() {
	viper.SetDefault("oauth_on", true)
}

type TokenInfo struct {
	UserID string
	IP     string
}
