// Package dbaccess 数据访问层 -outbox发件箱
package dbaccess

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
)

type outbox struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	outboxOnce sync.Once
	ob         *outbox
)

// NewOutbox 创建outbox 数据库对象
func NewOutbox() *outbox {
	outboxOnce.Do(func() {
		ob = &outbox{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})
	return ob
}

// 批量添加outbox消息
func (o *outbox) AddOutboxInfos(businessType int, messages []string, tx *sql.Tx) error {
	var valuesStr []string
	var inserts []interface{}
	curTime := time.Now().UnixNano() / 1000 // 获取当前时间
	for i := range messages {
		valuesStr = append(valuesStr, "(?,?,?)")
		inserts = append(inserts, businessType, messages[i], curTime)
	}
	valueStr := strings.Join(valuesStr, ",")
	dbName := common.GetDBName("user_management")
	strSQL := "insert into %s.t_outbox (f_business_type, f_message, f_create_time) values " + valueStr
	strSQL = fmt.Sprintf(strSQL, dbName)
	_, err := tx.Exec(strSQL, inserts...)
	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return nil
}

// 获取一条待推送的outbox消息
func (o *outbox) GetPushMessage(businessType int, tx *sql.Tx) (messageID int64, message string, err error) {
	dbName := common.GetDBName("user_management")
	// 使用 select ... for update 来加锁，防止事务处理时其他进程读取数据进行处理，间接实现了分布式锁
	getLockSQL := "select f_business_type from %s.t_outbox_lock where f_business_type = ? for update"
	getLockSQL = fmt.Sprintf(getLockSQL, dbName)
	rows, err := tx.Query(getLockSQL, businessType)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	if closeErr := rows.Close(); closeErr != nil {
		o.logger.Errorln(closeErr)
	}

	strSQL := "select f_id, f_message from %s.t_outbox where f_business_type = ? order by f_create_time asc limit 1"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, err = tx.Query(strSQL, businessType)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()

	for rows.Next() {
		if err = rows.Scan(&messageID, &message); err != nil {
			o.logger.Errorln(err, strSQL)
			return
		}
	}
	return
}

// 删除outbox消息
func (o *outbox) DeleteOutboxInfoByID(messageID int64, tx *sql.Tx) error {
	dbName := common.GetDBName("user_management")
	strSQL := "DELETE FROM %s.t_outbox WHERE f_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	_, err := tx.Exec(strSQL, messageID)

	if err != nil {
		o.logger.Errorln(err)
		return err
	}
	return nil
}
