// Package common 日志服务模块
package common

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 日志服务，可适配其他日志组件
type Logger interface {
	Infof(format string, args ...any)
	Infoln(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)
	Tracef(format string, args ...any)
	Traceln(args ...any)
	Panicf(format string, args ...any)
	Panicln(args ...any)
	Fatalf(format string, args ...any)
	Fatalln(args ...any)
}

var (
	logHandle      *logrus.Logger
	logHandlerOnce sync.Once
)

// NewLogger 获取日志对象
func NewLogger() *logrus.Logger {
	logHandlerOnce.Do(func() {
		logHandle = logrus.New()
		logHandle.SetReportCaller(true)
		logHandle.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			PadLevelText:              true,
			TimestampFormat:           time.RFC3339Nano,
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return strings.TrimPrefix(frame.Function, "devops.aishu.cn/AISHUDevOps/AnyShareFamily/_git"),
					fmt.Sprintf("%s:%d", frame.File, frame.Line)
			},
		})
		logout := os.Getenv("LOGOUT")
		if logout != "" {
			logHandle.SetOutput(os.Stdout)
		} else {
			logDir := "/var/log/authorization/"
			logName := "authorization.log"
			var filePerm os.FileMode = 0o750
			err := os.MkdirAll(logDir, filePerm)
			if err != nil {
				fmt.Println("mkdir err", err)
			}
			logFileName := path.Join(logDir, logName)
			logFile, err := os.OpenFile(path.Clean(logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
			if err != nil {
				fmt.Println("open log file err", err)
			}
			logHandle.SetOutput(logFile)
		}
	})
	return logHandle
}

// SetLogLevel 设置日志等级
func SetLogLevel(levelStr string) error {
	// 日志等级枚举 panic, fatal, error, warn/warning, info, debug, trace
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return err
	}
	logHandle.SetLevel(level)
	return nil
}

// GetLogLevel 获取日志等级
func GetLogLevel() string {
	return logHandle.GetLevel().String()
}
