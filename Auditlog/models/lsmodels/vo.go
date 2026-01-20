package lsmodels

// 转存策略
type DumpStrategy struct {
	RetentionPeriod     int    `json:"retention_period" validate:"required"`
	RetentionPeriodUnit string `json:"retention_period_unit" validate:"required"`
	DumpTime            string `json:"dump_time" validate:"required"`
	DumpFormat          string `json:"dump_format" validate:"required"`
}

type ScopeStrategyVO struct {
	ID          int      `json:"id"`
	LogType     int8     `json:"type" validate:"required"`
	LogCategory int8     `json:"category" validate:"required"`
	Role        string   `json:"role" validate:"required"`
	Scope       []string `json:"scope" validate:"required"`
}

type GetScopeStrategyReq struct {
	Category int8   `json:"category"`
	Role     string `json:"role"`
	Type     int8   `json:"type"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

type GetScopeStrategyRes struct {
	Entries    []*ScopeStrategyVO `json:"entries"`
	TotalCount int64              `json:"total_count"`
}

// 日志查看范围策略
type LogScopeStrategy struct {
	ActiveLog  ActiveLogScope  `json:"active_log" validate:"required"`
	HistoryLog HistoryLogScope `json:"history_log" validate:"required"`
}

// 活跃日志查看范围
type ActiveLogScope []RoleScope

// 历史日志查看范围
type HistoryLogScope []string

// 角色查看范围
type RoleScope struct {
	Role  string   `json:"role" validate:"required"`
	Scope []string `json:"scope" validate:"required"`
}

// 历史日志下载加密状态
type HistoryLogDownloadPwdStatus struct {
	Status bool `json:"status"`
}

// 历史日志下载请求
type HistoryLogDownloadReq struct {
	ObjId    string `json:"obj_id"`
	Password string `json:"pwd"`
}

// 历史日志下载响应
type HistoryLogDownloadRes struct {
	TaskId string `json:"task_id"`
}

// 历史日志下载进度
type HistoryLogDownloadProgress struct {
	Status bool `json:"status"`
}

// 历史日志下载任务信息
type HistoryLogDownloadTaskInfo struct {
	Status   int    `json:"status"`
	OssId    string `json:"oss_id"`
	FileName string `json:"file_name"`
}
