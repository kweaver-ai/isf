// Package conf 数据库&依赖服务相关的配置
package common

import (
	"os"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/log"

	"AuditLog/gocommon/api"
)

var SvcConfig Config

// GetEnv 封装os.Getenv(),可以指定默认值
func GetEnv(key, defaultV string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultV
	}
	return v
}

// LogConfig 日志配置
type LogConfig struct {
	Tags string
}

type Config struct {
	ServiceName             string
	PodIP                   string
	CommitID                string
	SystemID                string
	ServicePublicPort       string
	ServicePrivatePort      string
	Languaue                string
	DBType                  string
	UserMgntPrivateProtocol string
	UserMgntPrivateHost     string
	UserMgntPrivatePort     string
	UserMgntPublicProtocol  string
	UserMgntPublicHost      string
	UserMgntPublicPort      string
	DocumentPrivateProtocol string
	DocumentPrivateHost     string
	DocumentPrivatePort     string
	OAuthAdminHost          string
	OAuthAdminPort          string
	LogConfig               LogConfig
	Logger                  api.Logger
}

func init() {
	SvcConfig.ServiceName = GetEnv("SERVICE_NAME", "audit-log")
	SvcConfig.PodIP = GetEnv("POD_IP", "")
	SvcConfig.SystemID = GetEnv("DB_SYSTEM_ID", "")
	SvcConfig.Languaue = GetEnv("LANGUAUE", "zh-cn")
	SvcConfig.DBType = GetEnv("DB_TYPE", "")
	SvcConfig.ServicePublicPort = GetEnv("SERVICE_PUBLIC_PORT", "30569")
	SvcConfig.ServicePrivatePort = GetEnv("SERVICE_PRIVATE_PORT", "30570")
	SvcConfig.UserMgntPrivateProtocol = GetEnv("USER_MANAGEMENT_PRIVATE_PROTOCOL", "http")
	SvcConfig.UserMgntPrivateHost = GetEnv("USER_MANAGEMENT_PRIVATE_HOST", "user-management-private.anyshare")
	SvcConfig.UserMgntPrivatePort = GetEnv("USER_MANAGEMENT_PRIVATE_PORT", "30980")
	SvcConfig.UserMgntPublicProtocol = GetEnv("USER_MANAGEMENT_PUBLIC_PROTOCOL", "http")
	SvcConfig.UserMgntPublicHost = GetEnv("USER_MANAGEMENT_PUBLIC_HOST", "user-management-public.anyshare")
	SvcConfig.UserMgntPublicPort = GetEnv("USER_MANAGEMENT_PUBLIC_PORT", "30980")
	SvcConfig.DocumentPrivateProtocol = GetEnv("DOCUMENT_PRIVATE_PROTOCOL", "http")
	SvcConfig.DocumentPrivateHost = GetEnv("DOCUMENT_PRIVATE_HOST", "document-private.anyshare")
	SvcConfig.DocumentPrivatePort = GetEnv("DOCUMENT_PRIVATE_PORT", "30920")
	SvcConfig.OAuthAdminHost = GetEnv("HYDRA_ADMIN_HOST", "hydra-admin.anyshare")
	SvcConfig.OAuthAdminPort = GetEnv("HYDRA_ADMIN_PORT", "4445")
	l := api.NewTelemetryLogger(os.Stdout, log.InfoLevel, &api.LogOptionServiceInfo{
		Name:     SvcConfig.ServiceName,
		Version:  SvcConfig.CommitID,
		Instance: SvcConfig.PodIP,
	})
	SvcConfig.LogConfig.Tags = "audit_log"
	tagMarker := func(level field.Field, _ string, message field.Field) (tags []string) {
		if string(level.(field.StringField)) == "Error" {
			return []string{SvcConfig.LogConfig.Tags}
		}
		return []string{}
	}

	l.AddTagMarker(tagMarker)
	SvcConfig.Logger = l
}
