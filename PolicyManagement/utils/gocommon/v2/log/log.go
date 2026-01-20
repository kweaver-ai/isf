package log

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO: https://pkg.go.dev/github.com/rs/zerolog/log, https://pkg.go.dev/go.uber.org/zap

type Logger = *logrus.Logger

var loggerSingleton Logger
var loggerOnce sync.Once

type Config struct {
}

// InitLogger 初始化日志配置
func InitLogger(config *Config) {
	loggerOnce.Do(func() {
		// TODO: 仔细研究一下 https://github.com/sirupsen/logrus
		// TODO: 日志持久化 https://github.com/sirupsen/logrus/issues/784
		// TODO：set level
		// TODO：color
		// TODO: logger withFields
		// TODO: 使用entry自己实现

		loggerSingleton = logrus.New() // TODO: init by struct
		loggerSingleton.SetReportCaller(true)

		loggerSingleton.Out = os.Stdout
		loggerSingleton.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			PadLevelText:              true,
			TimestampFormat:           time.RFC3339Nano,
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
		})
	})
}

// NewLogger 获取Logger实例
func NewLogger() Logger {
	return loggerSingleton
}
