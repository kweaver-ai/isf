package locale

// 日志报表中心国际化信息 key
const (
	// 日志等级
	RCLogLevelInfo string = "rc_log_level_info" // 信息
	RCLogLevelWarn string = "rc_log_level_warn" // 警告
)

// 报表中心
const (
	// 日志列表表头
	RCLogID             string = "rc_log_id"              // 日志ID
	RCLogFileName       string = "rc_log_file_name"       // 文件名称
	RCLogDumpDate       string = "rc_log_dump_date"       // 转存时间
	RCLogSize           string = "rc_log_size"            // 文件大小
	RCLogOperation      string = "rc_log_operation"       // 操作
	RCLogLevel          string = "rc_log_level"           // 级别
	RCLogDate           string = "rc_log_date"            // 时间
	RCLogMac            string = "rc_log_mac"             // 设备地址
	RCLogIP             string = "rc_log_ip"              // IP地址
	RCLogUser           string = "rc_log_user"            // 用户
	RCLogUserPaths      string = "rc_log_user_paths"      // 部门
	RCLogOpType         string = "rc_log_op_type"         // 操作类型
	RCLogMsg            string = "rc_log_msg"             // 日志描述
	RCLogExMsg          string = "rc_log_ex_msg"          // 附加信息
	RCLogUserAgent      string = "rc_log_user_agent"      // 用户代理
	RCLogAdditionalInfo string = "rc_log_additional_info" // 备注
	RCLogObjName        string = "rc_log_obj_name"        // 对象名称
	RCLogObjType        string = "rc_log_obj_type"        // 对象类型

	// 报表中心
	RCLogDataSourceGroup    string = "rc_log_data_source_group"    // 数据源组
	RCLogReportLogin        string = "rc_log_report_login"         // 访问日志
	RCLogReportMgnt         string = "rc_log_report_mgnt"          // 管理日志
	RCLogReportOp           string = "rc_log_report_op"            // 操作日志
	RCLogReportHistoryLogin string = "rc_log_report_history_login" // 历史登录日志
	RCLogReportHistoryMgnt  string = "rc_log_report_history_mgnt"  // 历史管理日志
	RCLogReportHistoryOp    string = "rc_log_report_history_op"    // 历史操作日志
	RCLogReportGroup        string = "rc_log_report_group"         // 报表组
)

// 日志策略
const (
	LogDumpMsg             string = "log_dump_msg"             // 转存管理日志
	LogDumpExMsg           string = "log_dump_ex_msg"          // 转存管理日志附加信息
	SetLogDumpStrategy     string = "set_log_dump_strategy"    // 设置转存策略日志
	LogDumpPeriod          string = "log_dump_period"          // 转存周期
	LogDumpFormat          string = "log_dump_format"          // 转存格式
	LogDumpTime            string = "log_dump_time"            // 转存时间
	SetHistoryEncrypted    string = "set_history_encrypted"    // 设置历史日志加密
	CancelHistoryEncrypted string = "cancel_history_encrypted" // 取消历史日志加密

	// 策略管控
	NewLogScopeStrategy    string = "new_log_scope_strategy"    // 新建策略管控日志
	EditLogScopeStrategy   string = "edit_log_scope_strategy"   // 编辑策略管控日志
	DeleteLogScopeStrategy string = "delete_log_scope_strategy" // 删除策略管控日志
	LogType                string = "log_type"                  // 日志类型
	LogCategory            string = "log_category"              // 日志分类
	LogRole                string = "log_role"                  // 查看者
	LogScope               string = "log_scope"                 // 可见范围
	LogTypeLogin           string = "log_type_login"            // 登录日志
	LogTypeMgnt            string = "log_type_mgnt"             // 管理日志
	LogTypeOp              string = "log_type_op"               // 操作日志
	LogTypeOther           string = "log_type_other"            // 其他
	LogCategoryActive      string = "log_category_active"       // 活跃日志
	LogCategoryHistory     string = "log_category_history"      // 历史日志
	ExportLogSuccess       string = "export_log_success"        // 导出日志成功
)

var LogTypeMap = map[int]string{
	0:  LogTypeOther,
	10: LogTypeLogin,
	11: LogTypeMgnt,
	12: LogTypeOp,
}

var LogCategoryMap = map[int]string{
	1: LogCategoryActive,
	2: LogCategoryHistory,
}

// LogLevelMap 日志级别映射
var LogLevelMap = map[int]string{
	1: RCLogLevelInfo,
	2: RCLogLevelWarn,
}
