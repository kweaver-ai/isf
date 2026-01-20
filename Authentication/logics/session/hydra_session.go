package session

import (
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	hydraOnce sync.Once
	hSession  *hydraSession
)

type hydraSession struct {
	hydraAdmin interfaces.DnHydraAdmin
	ob         interfaces.LogicsOutbox
	logger     common.Logger
	pool       *sqlx.DB
}

// NewHydraSession 创建hydraSession处理对象
func NewHydraSession() *hydraSession {
	hydraOnce.Do(func() {
		hSession = &hydraSession{
			hydraAdmin: logics.DnHydraAdmin,
			ob:         logics.NewOutbox(logics.OutboxBusinessHydraSession),
			logger:     common.NewLogger(),
			pool:       logics.DBPool,
		}
	})

	hSession.ob.RegisterHandlers(logics.OutboxDeleteHydraSession, func(content interface{}) error {
		info := content.(map[string]interface{})
		if err := hSession.hydraAdmin.DeleteSession(info["user_id"].(string), info["client_id"].(string)); err != nil {
			return err
		}
		return nil
	})

	return hSession
}

// Delete 删除login、consent会话
func (s *hydraSession) Delete(userID, clientID string) error {
	// 获取事务处理器
	tx, err := s.pool.Begin()
	if err != nil {
		return err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				s.logger.Errorf("DeleteHydraSession Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			s.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.logger.Errorf("DeleteHydraSession Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["user_id"] = userID
	contentJSON["client_id"] = clientID

	err = s.ob.AddOutboxInfo(logics.OutboxDeleteHydraSession, contentJSON, tx)
	if err != nil {
		s.logger.Errorf("Add Outbox Info err:%v", err)
		return err
	}

	return nil
}
