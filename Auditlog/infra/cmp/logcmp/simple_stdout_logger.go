package logcmp

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"AuditLog/infra/cmp/icmp"
)

var simpleStdoutLogger *serverLog

// GetStdoutLogger 不依赖配置文件，直接输出到stdout
func GetStdoutLogger() (logger icmp.Logger) {
	if simpleStdoutLogger != nil {
		return simpleStdoutLogger
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	simpleStdoutLogger = &serverLog{}
	simpleStdoutLogger.logger = logrus.New()
	simpleStdoutLogger.logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	simpleStdoutLogger.logger.SetOutput(os.Stdout)

	return simpleStdoutLogger
}
