package log

import (
	clog "policy_mgnt/utils/gocommon/v2/log"

	"policy_mgnt/common/config"
)

// InitLogger 调用GoCommon.git/v2/log中的InitLogger，以初始化日志配置
func InitLogger() {
	clog.InitLogger(&config.Config.LogConfig)
}
