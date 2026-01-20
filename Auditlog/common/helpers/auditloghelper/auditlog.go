package auditloghelper

type AuditLog struct {
	UserID   string         `json:"user_id"`
	UserName string         `json:"user_name"`
	UserType NcTLogUserType `json:"user_type"`

	Level NcTLogLevel `json:"level"`

	Date int64  `json:"date"`
	IP   string `json:"ip"`
	Mac  string `json:"mac"`

	Msg   string `json:"msg"`
	ExMsg string `json:"ex_msg"`

	UserAgent string `json:"user_agent"`

	OpType int `json:"op_type"`

	OutBizID string `json:"out_biz_id"`

	DeptPaths string `json:"dept_paths"`
}
