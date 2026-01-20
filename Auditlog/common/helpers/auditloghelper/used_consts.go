package auditloghelper

// LogLevel 日志级别（别名）
type LogLevel = NcTLogLevel

const (
	All  LogLevel = NcTLogLevel_NCT_LL_ALL  // 所有日志级别
	Info LogLevel = NcTLogLevel_NCT_LL_INFO // 信息
	Warn LogLevel = NcTLogLevel_NCT_LL_WARN // 警告 删除日志调用此
)

// MgtOpLogType 管理日志操作类型（别名）
type MgtOpLogType = NcTManagementType

const (
	Create MgtOpLogType = NcTManagementType_NCT_MNT_CREATE // 创建
	Update MgtOpLogType = NcTManagementType_NCT_MNT_SET    // 修改
	Delete MgtOpLogType = NcTManagementType_NCT_MNT_DELETE // 删除
)
