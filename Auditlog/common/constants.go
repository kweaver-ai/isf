package common

import (
	"AuditLog/common/constants/logconsts"
	"AuditLog/interfaces"
)

const (
	RedisUniqueID     string = "as:audit_log:unique_id:"
	AuthenticatedUser string = "authenticated_user"
	AnonymousUser     string = "anonymous_user"
	APP               string = "app"
	InternalService   string = "internal_service"
)

// 角色常量的定义
const (
	SuperAdmin string = "super_admin"
	SysAdmin   string = "sys_admin"
	AuditAdmin string = "audit_admin"
	SecAdmin   string = "sec_admin"
	OrgManager string = "org_manager"
	OrgAudit   string = "org_audit"
	NormalUser string = "normal_user"
	App        string = "app"
)

// 三权分立：系统管理员、安全管理员、审计管理员的角色
var MutuallyRoles = []string{SysAdmin, SecAdmin, AuditAdmin}

// 日志类型常量
const (
	Login      string = "login"
	Management string = "management"
	Operation  string = "operation"
	Other      string = "other"
)

// AllLogType 所有的日志类型
var AllLogType = []string{Login, Management, Operation}

// LogTypeMap 日志类型映射
var LogTypeMap = map[string]int{
	Other:      0,
	Login:      10,
	Management: 11,
	Operation:  12,
}

var AllLogTypeInt = []int{10, 11, 12}

// VisitorKey 访问者字段的通用上下文key
const VisitorKey = "visitor"

const (
	KCServiceName = "KnowledgeCenter"
)

// opensearch 索引库
var OpensearchIndexMap = map[string]string{
	"kc_operation":  "as-operation-log-kc_operation",
	"doc_operation": "as-operation-log-doc_operation",
}

var (
	LogPrefix               = "as:personlization:unique_id:"
	DeptPersFeatureLock     = LogPrefix + "dept"
	DeptPersFeatureCronLock = LogPrefix + "cron"
)

// var OpensearchIndexMap = map[string]string{
// 	"kc_operation":  "ar-as_operation",
// 	"doc_operation": "ar-as_operation",
// }

var (
	MapLogType = map[string]interfaces.LogType{
		"login":      interfaces.LogType_Login,
		"management": interfaces.LogType_Management,
		"operation":  interfaces.LogType_Operation,
	}

	MapLevelType = map[string]int{
		"WARN": logconsts.LogLevel.WARN,
		"INFO": logconsts.LogLevel.INFO,
	}
)
