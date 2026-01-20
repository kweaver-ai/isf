package receo

import (
	"AuditLog/common/utils"
	"AuditLog/domain/entity/oprlogeo"
)

type RecLogEntry struct {
	*oprlogeo.LogEntry
	CreatedTime    int64  `json:"__rec_log_created_time__"`
	CreatedTimeTxt string `json:"__rec_log_created_time_txt__"`
	ULid           string `json:"__rec_log_ulid__"` // 可能不能保证肯定唯一，仅用于方便查询定位等
}

func NewRecLogEntry(logEntry *oprlogeo.LogEntry, createdTime int64) *RecLogEntry {
	return &RecLogEntry{
		LogEntry:       logEntry,
		CreatedTime:    createdTime,
		CreatedTimeTxt: utils.FormatTimeUnix(createdTime),
		ULid:           utils.UlidMake(),
	}
}
