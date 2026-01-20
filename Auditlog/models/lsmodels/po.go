package lsmodels

// 日志查看范围策略
type ScopeStrategyPO struct {
	ID          int64  `gorm:"column:f_id;primaryKey"`               // 主键ID
	LogType     int8   `gorm:"column:f_log_type;index"`              // 日志类型：10-访问日志，11-管理日志，12-操作日志
	LogCategory int8   `gorm:"column:f_log_category;index"`          // 日志分类：1-活跃日志，2-历史日志
	Role        string `gorm:"column:f_role;type:char(128);index"`   // 查看者角色名
	Scope       string `gorm:"column:f_scope;type:varchar(1024)"`    // 查看范围
	CreatedAt   int64  `gorm:"column:f_created_at"`                  // 创建时间
	CreatedBy   string `gorm:"column:f_created_by;type:varchar(64)"` // 创建者ID
	UpdatedAt   int64  `gorm:"column:f_updated_at"`                  // 更新时间
	UpdatedBy   string `gorm:"column:f_updated_by;type:varchar(64)"` // 更新者ID
}
