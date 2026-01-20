package locale

import (
	"AuditLog/common"
	"AuditLog/common/constants"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/infra/cmp/langcmp"
)

// i18ns 国际化信息映射表
var i18ns = map[string]map[langcmp.Lang]string{
	constants.SystemID: {
		langcmp.ZhCN: "系统",
		langcmp.ZhTW: "系統",
		langcmp.En:   "System",
	},
	RCLogID: {
		langcmp.ZhCN: "日志ID",
		langcmp.ZhTW: "日誌ID",
		langcmp.En:   "Log ID",
	},
	RCLogFileName: {
		langcmp.ZhCN: "文件名称",
		langcmp.ZhTW: "文件名稱",
		langcmp.En:   "File Name",
	},
	RCLogDumpDate: {
		langcmp.ZhCN: "转存时间",
		langcmp.ZhTW: "轉存時間",
		langcmp.En:   "Dump Time",
	},
	RCLogSize: {
		langcmp.ZhCN: "大小",
		langcmp.ZhTW: "大小",
		langcmp.En:   "Size",
	},
	RCLogOperation: {
		langcmp.ZhCN: "操作",
		langcmp.ZhTW: "操作",
		langcmp.En:   "Operation",
	},
	RCLogLevel: {
		langcmp.ZhCN: "级别",
		langcmp.ZhTW: "級別",
		langcmp.En:   "Level",
	},
	RCLogDate: {
		langcmp.ZhCN: "时间",
		langcmp.ZhTW: "時間",
		langcmp.En:   "Time",
	},
	RCLogMac: {
		langcmp.ZhCN: "设备地址",
		langcmp.ZhTW: "裝置位址",
		langcmp.En:   "Device Address",
	},
	RCLogIP: {
		langcmp.ZhCN: "IP地址",
		langcmp.ZhTW: "IP位址",
		langcmp.En:   "IP",
	},
	RCLogUser: {
		langcmp.ZhCN: "用户",
		langcmp.ZhTW: "使用者",
		langcmp.En:   "User",
	},
	RCLogUserPaths: {
		langcmp.ZhCN: "部门",
		langcmp.ZhTW: "部門",
		langcmp.En:   "Department",
	},
	RCLogOpType: {
		langcmp.ZhCN: "操作",
		langcmp.ZhTW: "操作",
		langcmp.En:   "Operation",
	},
	RCLogMsg: {
		langcmp.ZhCN: "日志描述",
		langcmp.ZhTW: "日誌描述",
		langcmp.En:   "Details",
	},
	RCLogExMsg: {
		langcmp.ZhCN: "附加信息",
		langcmp.ZhTW: "其他資訊",
		langcmp.En:   "Additional Info",
	},
	RCLogObjName: {
		langcmp.ZhCN: "对象名称",
		langcmp.ZhTW: "對象名稱",
		langcmp.En:   "Object Name",
	},
	RCLogObjType: {
		langcmp.ZhCN: "操作对象",
		langcmp.ZhTW: "操作對象",
		langcmp.En:   "Object Type",
	},
	RCLogUserAgent: {
		langcmp.ZhCN: "用户代理",
		langcmp.ZhTW: "使用者代理",
		langcmp.En:   "User Agent",
	},
	RCLogAdditionalInfo: {
		langcmp.ZhCN: "备注",
		langcmp.ZhTW: "備註",
		langcmp.En:   "Remark",
	},
	RCLogDataSourceGroup: {
		langcmp.ZhCN: "审计日志",
		langcmp.ZhTW: "審計日誌",
		langcmp.En:   "Audit Logs",
	},
	RCLogReportGroup: {
		langcmp.ZhCN: "审计日志",
		langcmp.ZhTW: "審計日誌",
		langcmp.En:   "Audit Logs",
	},
	RCLogReportLogin: {
		langcmp.ZhCN: "访问日志",
		langcmp.ZhTW: "存取日誌",
		langcmp.En:   "Access Logs",
	},
	RCLogReportMgnt: {
		langcmp.ZhCN: "管理日志",
		langcmp.ZhTW: "管理日誌",
		langcmp.En:   "Management Logs",
	},
	RCLogReportOp: {
		langcmp.ZhCN: "操作日志",
		langcmp.ZhTW: "操作日誌",
		langcmp.En:   "Operation Logs",
	},
	RCLogReportHistoryLogin: {
		langcmp.ZhCN: "历史访问日志",
		langcmp.ZhTW: "歷史存取日誌",
		langcmp.En:   "Historical Access Logs",
	},
	RCLogReportHistoryMgnt: {
		langcmp.ZhCN: "历史管理日志",
		langcmp.ZhTW: "歷史管理日誌",
		langcmp.En:   "Historical Management Logs",
	},
	RCLogReportHistoryOp: {
		langcmp.ZhCN: "历史操作日志",
		langcmp.ZhTW: "歷史操作日誌",
		langcmp.En:   "Historical Operation Logs",
	},
	LogDumpExMsg: {
		langcmp.ZhCN: "系统自动执行了一次日志转存，已产生归档文件<%s>",
		langcmp.ZhTW: "系統自動執行了一次日誌轉存，已產生歸檔檔案<%s>",
		langcmp.En:   "The system automatically performed a log dump and generated an archive file <%s>.",
	},
	LogDumpMsg: {
		langcmp.ZhCN: "转存 历史%s 成功",
		langcmp.ZhTW: "轉存 歷史%s 成功",
		langcmp.En:   "Successfully dumped Historical %s",
	},
	NewLogScopeStrategy: {
		langcmp.ZhCN: "新建 %s策略 成功",
		langcmp.ZhTW: "新建 %s策略 成功",
		langcmp.En:   "Successfully created %s policy",
	},
	EditLogScopeStrategy: {
		langcmp.ZhCN: "编辑 %s策略 成功",
		langcmp.ZhTW: "編輯 %s策略 成功",
		langcmp.En:   "Successfully edited %s policy",
	},
	DeleteLogScopeStrategy: {
		langcmp.ZhCN: "删除 %s策略 成功",
		langcmp.ZhTW: "刪除 %s策略 成功",
		langcmp.En:   "Successfully deleted %s policy",
	},
	SetLogDumpStrategy: {
		langcmp.ZhCN: "设置 日志转存策略 成功",
		langcmp.ZhTW: "設定 日誌轉存策略 成功",
		langcmp.En:   "Successfully set log dump policy",
	},
	SetHistoryEncrypted: {
		langcmp.ZhCN: "设置 历史日志下载加密 成功",
		langcmp.ZhTW: "設定 歷史日誌下載加密 成功",
		langcmp.En:   "Successfully enabled encryption for historical log downloads",
	},
	CancelHistoryEncrypted: {
		langcmp.ZhCN: "取消 设置历史日志下载加密 成功",
		langcmp.ZhTW: "取消 設定歷史日誌下載加密 成功",
		langcmp.En:   "Successfully canceled encryption for historical log downloads",
	},
	ExportLogSuccess: {
		langcmp.ZhCN: "导出历史日志文件 <%s> 成功",
		langcmp.ZhTW: "匯出歷史日誌檔案 <%s> 成功",
		langcmp.En:   "Successfully export historical log <%s> successfully",
	},
	LogDumpPeriod: {
		langcmp.ZhCN: "转存周期",
		langcmp.ZhTW: "轉存週期",
		langcmp.En:   "Dump Cycle",
	},
	LogDumpFormat: {
		langcmp.ZhCN: "转存格式",
		langcmp.ZhTW: "轉存格式",
		langcmp.En:   "Dump Format",
	},
	LogDumpTime: {
		langcmp.ZhCN: "转存时间",
		langcmp.ZhTW: "轉存時間",
		langcmp.En:   "Dump Time",
	},
	lsconsts.Day: {
		langcmp.ZhCN: "天",
		langcmp.ZhTW: "天",
		langcmp.En:   "Day",
	},
	lsconsts.Week: {
		langcmp.ZhCN: "周",
		langcmp.ZhTW: "週",
		langcmp.En:   "Week",
	},
	lsconsts.Month: {
		langcmp.ZhCN: "月",
		langcmp.ZhTW: "月",
		langcmp.En:   "Month",
	},
	lsconsts.Year: {
		langcmp.ZhCN: "年",
		langcmp.ZhTW: "年",
		langcmp.En:   "Year",
	},
	LogType: {
		langcmp.ZhCN: "日志类型",
		langcmp.ZhTW: "日誌類型",
		langcmp.En:   "Log Type",
	},
	LogCategory: {
		langcmp.ZhCN: "日志分类",
		langcmp.ZhTW: "日誌分類",
		langcmp.En:   "Log Category",
	},
	LogRole: {
		langcmp.ZhCN: "查看者",
		langcmp.ZhTW: "查看者",
		langcmp.En:   "Viewer",
	},
	LogScope: {
		langcmp.ZhCN: "可见范围",
		langcmp.ZhTW: "可見範圍",
		langcmp.En:   "Visible Range",
	},
	LogTypeLogin: {
		langcmp.ZhCN: "登录日志",
		langcmp.ZhTW: "登入日誌",
		langcmp.En:   "Login Log",
	},
	LogTypeMgnt: {
		langcmp.ZhCN: "管理日志",
		langcmp.ZhTW: "管理日誌",
		langcmp.En:   "Management Log",
	},
	LogTypeOp: {
		langcmp.ZhCN: "操作日志",
		langcmp.ZhTW: "操作日誌",
		langcmp.En:   "Operation Log",
	},
	LogTypeOther: {
		langcmp.ZhCN: "其他",
		langcmp.ZhTW: "其他",
		langcmp.En:   "Other",
	},
	LogCategoryActive: {
		langcmp.ZhCN: "活跃日志",
		langcmp.ZhTW: "活躍日誌",
		langcmp.En:   "Active Log",
	},
	LogCategoryHistory: {
		langcmp.ZhCN: "历史日志",
		langcmp.ZhTW: "歷史日誌",
		langcmp.En:   "Historical Log",
	},
	common.SuperAdmin: {
		langcmp.ZhCN: "超级管理员",
		langcmp.ZhTW: "超級管理員",
		langcmp.En:   "Super Admin",
	},
	common.SysAdmin: {
		langcmp.ZhCN: "系统管理员",
		langcmp.ZhTW: "系統管理員",
		langcmp.En:   "System Admin",
	},
	common.AuditAdmin: {
		langcmp.ZhCN: "审计管理员",
		langcmp.ZhTW: "審計管理員",
		langcmp.En:   "Audit Admin",
	},
	common.SecAdmin: {
		langcmp.ZhCN: "安全管理员",
		langcmp.ZhTW: "安全管理員",
		langcmp.En:   "Security Admin",
	},
	common.OrgManager: {
		langcmp.ZhCN: "组织管理员",
		langcmp.ZhTW: "組織管理員",
		langcmp.En:   "Org Manager",
	},
	common.OrgAudit: {
		langcmp.ZhCN: "组织审计员",
		langcmp.ZhTW: "組織審計員",
		langcmp.En:   "Org Audit",
	},
	common.NormalUser: {
		langcmp.ZhCN: "普通用户",
		langcmp.ZhTW: "一般使用者",
		langcmp.En:   "Regular User",
	},
	common.AnonymousUser: {
		langcmp.ZhCN: "匿名用户",
		langcmp.ZhTW: "匿名使用者",
		langcmp.En:   "Anonymous User",
	},
}

// rcLevelI18n 报表中心日志等级国际化信息映射表
var rcLevelI18n = map[string]map[langcmp.Lang]string{
	RCLogLevelInfo: {
		langcmp.ZhCN: "信息",
		langcmp.ZhTW: "資訊",
		langcmp.En:   "Information",
	},
	RCLogLevelWarn: {
		langcmp.ZhCN: "警告",
		langcmp.ZhTW: "警告",
		langcmp.En:   "Warning",
	},
}
