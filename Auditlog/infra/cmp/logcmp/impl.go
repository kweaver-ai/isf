package logcmp

import (
	"github.com/sirupsen/logrus"
)

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

	l.logger.Warnf(format, args...)
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
