package constants

const (
	SvcNameEnvKey = "SERVICE_NAME"

	SvcErrFlagNo = 62 // audit-log 服务错误码标识号（见：https://confluence.aishu.cn/pages/viewpage.action?pageId=73108671）

	// SvcErrFlagBizBasePers 个性化推荐 后三位错误码 基础偏移
	SvcErrFlagBizBasePers = 200

	RedisKeyPrefix = "audit-log"
)

const (
	ReportInitLockKey = "as:audit_log:report_init_lock"
)

// DumpLogLockKey 日志转存锁
const (
	DumpLogLockKey string = "as:audit_log:dump_log_lock"
)

// SystemID 系统账户ID
const (
	SystemID string = "da5bfdc4-cb4b-4b28-90c2-9eca46c3e500"
)
