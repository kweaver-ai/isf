// Package logics anonymous Anyshare 业务逻辑层 -匿名共享
package logics

import (
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
)

type anonymous struct {
	db            interfaces.DBAnonymous
	messageBroker interfaces.DrivenMessageBroker
	ob            interfaces.LogicsOutbox
	logger        common.Logger
	pool          *sqlx.DB
}

var (
	anonyOnce   sync.Once
	anonyLogics *anonymous
)

// NewAnonymous 创建新的anonymous对象
func NewAnonymous() *anonymous {
	anonyOnce.Do(func() {
		anonyLogics = &anonymous{
			db:            dbAnonymous,
			messageBroker: dnMessageBroker,
			ob:            NewOutbox(OutboxBusinessAnonymous),
			logger:        common.NewLogger(),
			pool:          dbPool,
		}

		anonyLogics.ob.RegisterHandlers(outboxAnonymityAuth, func(content interface{}) error {
			contentJSON := content.(map[string]interface{})
			info := map[string]interface{}{
				"id":             contentJSON["id"].(string),
				"accessed_times": int32(contentJSON["accessed_times"].(float64)),
				"referrer":       contentJSON["referrer"].(string),
			}
			msgType := contentJSON["type"].(string)

			err := anonyLogics.messageBroker.AnonymityAuth(msgType, info)
			if err != nil {
				return err
			}
			return nil
		})
	})

	return anonyLogics
}

// Create 创建匿名账户
func (a *anonymous) Create(info interfaces.AnonymousInfo) error {
	return a.db.Create(&info)
}

// DeleteByID 根据ID删除匿名账户
func (a *anonymous) DeleteByID(id string) error {
	return a.db.DeleteByID(id)
}

// Authentication 认证匿名账户
func (a *anonymous) Authentication(id, password, referrer string) error {
	info, err := a.db.GetAccount(id)
	if err != nil {
		return err
	}

	if info.ID != id {
		return rest.NewHTTPError("record not exist", errors.AnonymityNotFound, nil)
	}

	// 检查访问密码
	if password != info.Password {
		return rest.NewHTTPError("wrong password", errors.AnonymityWrongPassword, nil)
	}

	// 检查访问次数
	if info.LimitedTimes != -1 && info.AccessedTimes >= info.LimitedTimes {
		return rest.NewHTTPError("the visits has reached the limit", errors.AnonymityReachLimitTimes, nil)
	}

	// 验证成功，获取事务处理器
	tx, err := a.pool.Begin()
	if err != nil {
		return err
	}
	// 异常时Rollback
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				a.logger.Errorf("Anonymity Authentication Rollback Error:%v", rollbackErr)
			}
		}
	}()

	// 访问计数+1
	err = a.db.AddAccessTimes(id, tx)
	if err != nil {
		return err
	}

	if info.Type != "" {
		contentJSON := map[string]any{
			"id":             id,
			"type":           info.Type,
			"accessed_times": info.AccessedTimes + 1,
			"referrer":       referrer,
		}

		err = a.ob.AddOutboxInfo(outboxAnonymityAuth, contentJSON, tx)
		if err != nil {
			a.logger.Errorf("Add Outbox Info err:%v", err)
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		a.logger.Errorf("Anonymity Authentication Commit Error:%v", err)
		return err
	}

	// notify outbox推送线程
	a.ob.NotifyPushOutboxThread()

	return nil
}

// DeleteByTime 删除过期匿名账户
func (a *anonymous) DeleteByTime(curTime int64) error {
	return a.db.DeleteByTime(curTime)
}

// GetByID 获取匿名账户信息
func (a *anonymous) GetByID(id string) (*interfaces.AnonymousInfo, error) {
	info, err := a.db.GetAccount(id)
	if err != nil {
		return nil, err
	}
	if info.ID == "" {
		return nil, rest.NewHTTPError("record not exist", errors.AnonymityNotFound, nil)
	}

	return &info, nil
}
