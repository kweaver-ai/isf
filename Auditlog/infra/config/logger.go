package config

import (
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/sirupsen/logrus"

	"AuditLog/common"
)

// LogConf 日志
type LogConf struct {
	LogDir   string `yaml:"log_dir"`
	LogFile  string `yaml:"log_file"`
	LogLevel string `yaml:"log_level"` // 日志级别 需要能通过logrus.ParseLevel转化成logrus.Level
}

func (l *LogConf) loadConf() {
	err := common.Configure(l, "log_config.yaml")
	if err != nil {
		panic(err)
	}
}

func (l *LogConf) MkLogDir() {
	err := gfile.Mkdir(l.LogDir)
	if err != nil {
		panic("create log dir error:" + err.Error())
	}
}

func (l *LogConf) GetLogrusLevel() logrus.Level {
	level, err := logrus.ParseLevel(l.LogLevel)
	if err != nil {
		panic("[GetLogrusLevel]: logrus.ParseLevel error: " + err.Error())
	}

	return level
}
