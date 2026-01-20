package logcmp

import (
	"sync"

	"AuditLog/common"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
)

var (
	logOnce       sync.Once
	l             api.Logger
	isInitialized bool
)

// InitLogger 初始化日志
func InitLogger(logConf *config.LogConf) {
	logOnce.Do(func() {
		initLogger(logConf)

		isInitialized = true
	})
}

func initLogger(logConf *config.LogConf) {
	l = common.SvcConfig.Logger
}

// GetLogger 获取日志句柄，如果没有初始化，会panic
func GetLogger() api.Logger {
	if !isInitialized {
		panic("logger not initialized")
	}

	return l
}
