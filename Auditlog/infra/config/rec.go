package config

import (
	"AuditLog/common"
)

// RecConf 日志
type RecConf struct {
	SaveDays                       int `yaml:"save_days"`
	RemoveOldLogTaskIntervalSecond int `yaml:"remove_old_log_task_interval_second"`
}

func (l *RecConf) loadConf() {
	err := common.Configure(l, "rec_config.yaml")
	if err != nil {
		panic(err)
	}
}
