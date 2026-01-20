// Package audit 逻辑层
package audit

import (
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/drivenadapters"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	aOnce sync.Once
	a     *audit
)

type audit struct {
	logger  common.Logger
	pool    *sqlx.DB
	ob      interfaces.LogicsOutbox
	eacpLog interfaces.DnEacpLog
}

// NewAudit 创建audit log处理对象
func NewAudit() *audit {
	aOnce.Do(func() {
		a = &audit{
			logger:  common.NewLogger(),
			pool:    logics.DBPool,
			ob:      logics.NewOutbox(logics.OutboxAuditLog),
			eacpLog: drivenadapters.NewEacpLog(),
		}

		a.ob.RegisterHandlers(logics.OutboxSendAuditLog, a.sendAuditLog)
	})

	return a
}

// Log 记录日志
func (a *audit) Log(topic string, message interface{}) (err error) {
	// 获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				a.logger.Errorf("Log Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			a.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				a.logger.Errorf("Log Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["topic"] = topic
	contentJSON["message"] = message

	err = a.ob.AddOutboxInfo(logics.OutboxSendAuditLog, contentJSON, tx)
	if err != nil {
		a.logger.Errorf("Log Add Outbox Info err:%v", err)
		return err
	}

	return nil
}

// sendAuditLog 发送审计日志
func (a *audit) sendAuditLog(content interface{}) error {
	info := content.(map[string]interface{})
	topic := info["topic"].(string)
	msg := info["message"]
	return a.eacpLog.Publish(topic, msg)
}
