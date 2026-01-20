package config

import (
	"strconv"
	"sync"

	"AuditLog/common"
)

var (
	ldOnce sync.Once
	ld     LogDumpConfigRepo
)

//go:generate mockgen -package mock -source ../../infra/config/log_dump_config.go -destination ./mock/log_dump_config.go

type LogDumpConfigRepo interface {
	GetDumpThresholdByType(logType string) (threshold int64)
	GetDumpLogNum() int64
	GetDumpIntervalTime() int64
}

type LogDumpConfig struct {
	DumpThresholdLogin      string `yaml:"log_dump_threshold_login"`
	DumpThresholdManagement string `yaml:"log_dump_threshold_management"`
	DumpThresholdOperation  string `yaml:"log_dump_threshold_operation"`
	DumpLogNum              string `yaml:"log_dump_log_num"`
	DumpIntervalTime        string `yaml:"log_dump_interval_time"`
}

func loadLogDumpConfig() LogDumpConfigRepo {
	ldOnce.Do(func() {
		ld = &LogDumpConfig{}

		err := common.Configure(ld, "log_dump_config.yaml")
		if err != nil {
			panic(err)
		}
	})

	return ld
}

// GetDumpThresholdByType 获取日志类型对应的阈值
func (ld *LogDumpConfig) GetDumpThresholdByType(logType string) (threshold int64) {
	var err error

	switch logType {
	case common.Login:
		threshold, err = strconv.ParseInt(ld.DumpThresholdLogin, 10, 64)
	case common.Management:
		threshold, err = strconv.ParseInt(ld.DumpThresholdManagement, 10, 64)
	case common.Operation:
		threshold, err = strconv.ParseInt(ld.DumpThresholdOperation, 10, 64)
	}

	if err != nil {
		panic(err)
	}

	if threshold == 0 {
		switch logType {
		case common.Login:
			threshold = 1000000
		case common.Management:
			threshold = 1000000
		case common.Operation:
			threshold = 8000000
		}
	}

	return
}

// GetDumpLogNum 获取每次sql执行删除的日志数量
func (ld *LogDumpConfig) GetDumpLogNum() int64 {
	dumpLogNum, err := strconv.ParseInt(ld.DumpLogNum, 10, 64)
	if err != nil {
		panic(err)
	}

	if dumpLogNum == 0 {
		return 50000
	}

	return dumpLogNum
}

// GetDumpIntervalTime 获取sql执行删除时间间隔
func (ld *LogDumpConfig) GetDumpIntervalTime() int64 {
	dumpIntervalTime, err := strconv.ParseInt(ld.DumpIntervalTime, 10, 64)
	if err != nil {
		panic(err)
	}

	return dumpIntervalTime
}
