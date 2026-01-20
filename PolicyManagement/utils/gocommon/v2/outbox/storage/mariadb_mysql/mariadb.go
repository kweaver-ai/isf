// Package mariadb 支持 mariadb和 mysql
package mariadb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	outboxer "policy_mgnt/utils/gocommon/v2/outbox"
	cutils "policy_mgnt/utils/gocommon/v2/utils"
)

// TODO 恢复日志

const (
	// DefaultEventStoreTable is the default table name.
	DefaultEventStoreTable = "t_event_store"
)

var (
	ErrNoDatabaseName = errors.New("no database name")
)

// Mariadb is the implementation of the data store.
type Mariadb struct {
	db              *gorm.DB
	DatabaseName    string
	EventStoreTable string
}

// WithInstance creates a mariadb data store with an existing db connection.
func WithInstance(db *gorm.DB) (*Mariadb, error) {
	m := Mariadb{
		db:              db,
		EventStoreTable: DefaultEventStoreTable,
	}

	// var databaseName string
	// if tx := db.Raw(`SELECT DATABASE()`).Scan(&databaseName); tx.Error != nil {
	// 	return nil, tx.Error
	// }

	// if len(databaseName) == 0 {
	// 	return nil, ErrNoDatabaseName
	// }

	// m.DatabaseName = databaseName

	//if len(m.EventStoreTable) == 0 {
	//	m.EventStoreTable = DefaultEventStoreTable
	//}

	return &m, nil
}

func (m *Mariadb) GetTx() *gorm.DB { // TODO 返回包装函数,从而实现*gorm.DB的解耦
	return m.db
}

func (m *Mariadb) Add(evt *outboxer.OutboxMessage) error {
	var err error
	evt.ID, err = cutils.GetSonyflakeID()
	if err != nil {
		return fmt.Errorf("GetSonyflakeID: %w", err)
	}
	if tx := m.db.Create(evt); tx.Error != nil {
		return fmt.Errorf("failed to insert event into the data store: %w", tx.Error)
	}
	return nil
}

func (m *Mariadb) GetEvent(tx *gorm.DB) (evt *outboxer.OutboxMessage, hasEvent bool, err error) {
	evt = &outboxer.OutboxMessage{}
	tmpTx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("f_dispatched = ?", false).Order("f_id asc").Limit(1).Find(evt)
	if tmpTx.Error != nil {
		return evt, false, fmt.Errorf("failed to get events from store: %w", tmpTx.Error)
	}

	return evt, tmpTx.RowsAffected > 0, nil
}

func (m *Mariadb) SetAsDispatched(tx *gorm.DB, id uint64) error {
	err := tx.Model(&outboxer.OutboxMessage{ID: id}).Updates(outboxer.OutboxMessage{
		Dispatched:   true,
		DispatchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}).Error

	if err != nil {
		return fmt.Errorf("failed to set event as dispatched: %w", err)
	}

	return nil
}

// Remove removes old messages from the data store.
func (m *Mariadb) Remove(retentionTime time.Duration, batchSize int32) error {
	tx := m.db.Where("f_dispatched = ?", true).Where("f_dispatched_at < ?", time.Now().Add(-retentionTime)).Limit(int(batchSize)).Delete(&outboxer.OutboxMessage{})
	if tx.Error != nil {
		return fmt.Errorf("failed to remove messages from the data store: %w", tx.Error)
	}
	return nil
}
