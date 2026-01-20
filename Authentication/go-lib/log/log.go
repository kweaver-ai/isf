package log

import (
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger 提供基本日志接口
type Logger interface {
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Tracef(format string, args ...interface{})
	Traceln(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

var (
	logOnce sync.Once
	l       Logger
)

// NewLogger 获取日志句柄
func NewLogger() Logger {
	// 初始化一个实例
	logOnce.Do(initLogger)

	return l
}

type serverLog struct {
	logger *logrus.Logger
}

// Infof 普通信息
func (l *serverLog) Infof(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Infof(format, args...)
}

// Infoln 普通信息
func (l *serverLog) Infoln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Infoln(args...)
}

// Warnf 警告信息
func (l *serverLog) Warnf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Warnf(format, args...)
}

// Warnln 警告信息
func (l *serverLog) Warnln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Warnln(args...)
}

// Errorf 错误信息
func (l *serverLog) Errorf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Errorf(format, args...)
}

// Errorln 错误信息
func (l *serverLog) Errorln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Errorln(args...)
}

// Debugf 调试信息
func (l *serverLog) Debugf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Debugf(format, args...)
}

// Debugln 调试信息
func (l *serverLog) Debugln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Debugln(args...)
}

// Tracef 跟踪信息
func (l *serverLog) Tracef(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Tracef(format, args...)
}

// Traceln 跟踪信息
func (l *serverLog) Traceln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Traceln(args...)
}

// Fatalf 致命错误
func (l *serverLog) Fatalf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Fatalf(format, args...)
}

// Fatalln 致命错误
func (l *serverLog) Fatalln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Fatalln(args...)
}

// Panicf 恐慌错误
func (l *serverLog) Panicf(format string, args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Panicf(format, args...)
}

// Panicln 恐慌错误
func (l *serverLog) Panicln(args ...interface{}) {
	if l.logger == nil {
		return
	}
	l.logger.Panicln(args...)
}

// 初始化
func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	logHandle := &serverLog{}
	logHandle.logger = logrus.New()
	logHandle.logger.SetFormatter(&logrus.JSONFormatter{})
	initLogOut(logHandle.logger)
	l = logHandle
}

// 设置log输出位置
func initLogOut(logger *logrus.Logger) {
	logout := os.Getenv("LOGOUT")
	// 不设置、设置为NOTSTDOUT，默认屏幕输出
	if logout != "NOTSTDOUT" {
		logger.SetOutput(os.Stdout)
	} else {
		logDir := os.Getenv("LOGDIR")
		if len(logDir) == 0 {
			panic("invalid log directory!")
		}

		logName := os.Getenv("LOGNAME")
		if len(logName) == 0 {
			panic("invalid log file name!")
		}

		logFileName := path.Join(logDir, logName)
		err := os.MkdirAll(logDir, 0750)
		if err != nil {
			panic(err)
		}
		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
		if err != nil {
			panic(err)
		}
		logger.SetOutput(logFile)
	}
}

// gin 中间件，记录日志
func Ginrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger := logrus.New()
		logger.Level = logrus.ErrorLevel
		logger.SetFormatter(&logrus.JSONFormatter{})
		initLogOut(logger)
		entry := logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			// "time":       end.Format(time.RFC3339),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
