package oprlogeo

import "AuditLog/common/enums/oprlogenums"

// LogEntry 主日志结构
type LogEntry struct {
	Recorder    string              `json:"recorder"`    // 日志记录者的身份，示例："AnyShare"
	BizType     oprlogenums.BizType `json:"biz_type"`    // 业务类型
	Operation   string              `json:"operation"`   // 执行的操作，示例："cd"
	Description string              `json:"description"` // 操作的描述，示例："用户"张三"从""进入到"部门文档库1/a"。"

	Operator     *Operator      `json:"operator,omitempty"`      // 操作员信息
	Object       *ObjectInfo    `json:"object,omitempty"`        // 操作对象信息
	TargetObject *ObjectInfo    `json:"target_object,omitempty"` // 目标对象信息
	LogFrom      *LogFrom       `json:"log_from,omitempty"`      // 日志来源
	Rec          *RecommendInfo `json:"rec,omitempty"`           // 推荐相关对象
	Detail       interface{}    `json:"detail,omitempty"`        // 业务模块扩展的其他字段
	Referer      *Referer       `json:"referer,omitempty"`       // 来源
}
