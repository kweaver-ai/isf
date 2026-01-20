package dbaccess

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type unorderedOutbox struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	unOnce sync.Once
	un     *unorderedOutbox
)

// NewUnorderedOutbox 创建无序outbox数据库对象
func NewUnorderedOutbox() *unorderedOutbox {
	unOnce.Do(func() {
		un = &unorderedOutbox{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})
	return un
}

// GetUnorderedOutboxInfo 获取无序outbox信息
func (un *unorderedOutbox) GetUnorderedOutboxInfo() (auditLogAsyncInfo interfaces.UnorderedOutbox, exist bool, err error) {
	dbName := common.GetDBName("authentication")
	// 乐观锁
	sqlStr := "select id, f_message, f_status from %s.t_outbox_unordered where f_status = 0 limit 1"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	for {
		row := un.db.QueryRow(sqlStr)
		if err := row.Scan(&auditLogAsyncInfo.ID, &auditLogAsyncInfo.Message, &auditLogAsyncInfo.Status); err != nil {
			if err == sql.ErrNoRows {
				return auditLogAsyncInfo, false, nil
			}
			return auditLogAsyncInfo, false, err
		}

		updateSQL := "update %s.t_outbox_unordered set f_status = 1, f_updated_at = ? where id = ? and f_status = 0"
		updateSQL = fmt.Sprintf(updateSQL, dbName)
		result, err := un.db.Exec(updateSQL, common.Now().UnixNano()/(1e3), auditLogAsyncInfo.ID)
		// 判断影响行数，影响行数为1，则代表获取成功，否则获取失败，需要重新获取
		if err != nil {
			return auditLogAsyncInfo, false, err
		}
		rowAffected, err := result.RowsAffected()
		if err != nil {
			return auditLogAsyncInfo, false, err
		}

		if rowAffected == 1 {
			auditLogAsyncInfo.Status = interfaces.OutboxInProgress
			return auditLogAsyncInfo, true, nil
		}
	}
}

// DeleteUnorderedOutboxInfoByID 根据ID删除无序outbox信息
func (un *unorderedOutbox) DeleteUnorderedOutboxInfoByID(id string) (err error) {
	dbName := common.GetDBName("authentication")
	deleteSQL := "delete from %s.t_outbox_unordered where id = ? and f_status = 1"
	deleteSQL = fmt.Sprintf(deleteSQL, dbName)
	_, err = un.db.Exec(deleteSQL, id)
	if err != nil {
		return
	}
	return
}

// UpdateUnorderedOutboxUpdateTimeByID 根据ID更新无序outbox更新时间
func (un *unorderedOutbox) UpdateUnorderedOutboxUpdateTimeByID(id string) (isUpdate bool, err error) {
	dbName := common.GetDBName("authentication")
	sqlStr := "update %s.t_outbox_unordered set f_updated_at = ? where id = ? and f_status =1"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	result, err := un.db.Exec(sqlStr, common.Now().UnixNano()/(1e3), id)
	if err != nil {
		return false, err
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected != 1 {
		return false, nil
	}
	return true, nil
}

// AddUnorderedOutboxInfo 添加无序outbox信息
func (un *unorderedOutbox) AddUnorderedOutboxInfo(auditLogAsyncInfo interfaces.UnorderedOutbox) (err error) {
	dbName := common.GetDBName("authentication")
	sqlStr := "INSERT INTO %s.t_outbox_unordered (`f_message`, `f_status`, `f_created_at`, `f_updated_at`) VALUES (?,?,?,?)"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	_, err = un.db.Exec(sqlStr, auditLogAsyncInfo.Message, auditLogAsyncInfo.Status, auditLogAsyncInfo.CreatedAt, auditLogAsyncInfo.UpdatedAt)
	return
}

// RestartUnorderedOutboxInfo 重置无序outbox信息状态
func (un *unorderedOutbox) RestartUnorderedOutboxInfo(updatedTime int64) (err error) {
	dbName := common.GetDBName("authentication")
	sqlStr := "update %s.t_outbox_unordered set f_status = 0 where f_status = 1 and f_updated_at < ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	_, err = un.db.Exec(sqlStr, updatedTime)
	return
}
