package config

type LogStrategy struct {
	Items []*LogStrategyItem
}

/**
 * 日志配置
 */
type LogStrategyItem struct {
	Type       string      `yaml:"type"`
	ActiveLog  []RoleScope `yaml:"active_log"`
	HistoryLog []string    `yaml:"history_log"`
}

/**
 * 角色以及查看范围
 */
type RoleScope struct {
	Role  string   `yaml:"role"`
	Scope []string `yaml:"scope"`
}

// GetActiveScopeByRole 获取活动日志配置
func (ls *LogStrategy) GetActiveScopeByRole(logType string, role string) []string {
	for _, v := range ls.Items {
		if v.Type == logType {
			for _, ac := range v.ActiveLog {
				if ac.Role == role {
					return ac.Scope
				}
			}
		}
	}

	return nil
}

// GetHistoryScope 获取历史日志配置
func (ls *LogStrategy) GetHistoryScope(logType string) []string {
	for _, v := range ls.Items {
		if v.Type == logType {
			return v.HistoryLog
		}
	}

	return nil
}
