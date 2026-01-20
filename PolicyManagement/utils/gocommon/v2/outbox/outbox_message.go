package outboxer

import (
	"database/sql"
)

/*
CREATE TABLE IF NOT EXISTS %s (
	f_id BIGINT AUTO_INCREMENT not null primary key,
	f_dispatched BOOL not null default false,
	f_dispatched_at DATETIME,
	f_payload LONGBLOB not null,
	f_options json,
	f_headers json
);
*/

// OutboxMessage represents a message that will be sent.
type OutboxMessage struct {
	ID           uint64        `gorm:"column:f_id;type:bigint(20);primary_key;not null;autoIncrement:false"`
	Dispatched   bool          `gorm:"column:f_dispatched;type:bool;not null;default:false"`
	DispatchedAt sql.NullTime  `gorm:"column:f_dispatched_at;type:datetime"`
	Payload      []byte        `gorm:"column:f_payload;type:longblob;not null"`
	Options      DynamicValues `gorm:"column:f_options;type:json;serializer:json"`
	Headers      DynamicValues `gorm:"column:f_headers;type:json;serializer:json"`
}

// TableName 表名
func (OutboxMessage) TableName() string {
	return "t_event_store"
}

type DynamicValues map[string]interface{}
