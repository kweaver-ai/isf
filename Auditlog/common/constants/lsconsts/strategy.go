package lsconsts

import "time"

// 日志转存周期单位
const (
	Day   string = "day"
	Week  string = "week"
	Month string = "month"
	Year  string = "year"
)

// HistoryMaxBatchSize 历史日志最大批量大小
const HistoryMaxBatchSize = 100000

// 日志转存文件后缀
const (
	XMLSuffix string = "xml"
	CSVSuffix string = "csv"
)

// 转存策略字段
const (
	RetentionPeriod     string = "retention_period"
	RetentionPeriodUnit string = "retention_period_unit"
	DumpTime            string = "dump_time"
	DumpFormat          string = "dump_format"
)

// 所有转存策略字段
var AllDumpFields = []string{RetentionPeriod, RetentionPeriodUnit, DumpTime, DumpFormat}

// 日志查看角色
const (
	ActiveLog  int = 1
	HistoryLog int = 2
)

var AllLogCategory = []int{ActiveLog, HistoryLog}

// 历史审计日志下载任务状态
const (
	ErrorStatus    int = -1
	PendingStatus  int = 0
	FinishedStatus int = 1
)

// 历史审计日志下载任务缓存键
const (
	TaskCacheKey   string        = "as:audit_log:history_log_download_task:"
	TaskInfoExpire time.Duration = 1800 * time.Second
)
