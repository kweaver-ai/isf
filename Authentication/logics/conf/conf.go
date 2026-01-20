// Package conf 逻辑层
package conf

import (
	"context"
	"database/sql"
	"errors"

	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	cOnce sync.Once
	c     *conf
)

type conf struct {
	cdb           interfaces.DBConf
	trace         observable.Tracer
	log           common.Logger
	userMgnt      interfaces.DnUserManagement
	ob            interfaces.LogicsOutbox
	pool          *sqlx.DB
	messageBroker interfaces.DrivenMessageBroker
}

// NewConf 创建conf处理对象
func NewConf() *conf {
	cOnce.Do(func() {
		c = &conf{
			cdb:           logics.DBConf,
			trace:         common.SvcARTrace,
			log:           common.NewLogger(),
			userMgnt:      logics.DnUserManagement,
			ob:            logics.NewOutbox(logics.OutboxConfig),
			pool:          logics.DBPool,
			messageBroker: logics.DnMessageBroker,
		}

		c.ob.RegisterHandlers(logics.OutboxAnonymousSmsExpUpdated, c.sendAnonymousSmsExpUpdatedMsg)
	})

	return c
}

func (c *conf) sendAnonymousSmsExpUpdatedMsg(content interface{}) (err error) {
	message := content.(map[string]interface{})
	smsExpiration := int(message["anonymous_sms_expiration"].(float64))
	return c.messageBroker.AnonymousSmsExpUpdated(smsExpiration)
}

// GetConfig 获取认证配置
func (c *conf) GetConfig(ctx context.Context, visitor *interfaces.Visitor, configKeys map[interfaces.ConfigKey]bool) (cfg interfaces.Config, err error) {
	if visitor != nil && visitor.ID != "" {
		err = c.checkAdmin(ctx, visitor)
		if err != nil {
			return
		}
	}

	return c.cdb.GetConfig(configKeys)
}

// GetConfigFromShareMgnt 获取认证配置
func (c *conf) GetConfigFromShareMgnt(ctx context.Context, visitor *interfaces.Visitor, configKeys map[interfaces.ConfigKey]bool) (cfg interfaces.Config, err error) {
	c.trace.SetInternalSpanName("逻辑层-获取认证配置")
	newCtx, span := c.trace.AddInternalTrace(ctx)
	defer func() { c.trace.TelemetrySpanEnd(span, err) }()

	if visitor != nil && visitor.ID != "" {
		err = c.checkAdmin(newCtx, visitor)
		if err != nil {
			return
		}
	}

	return c.cdb.GetConfigFromShareMgnt(newCtx, configKeys)
}

// SetConfig 设置认证配置
func (c *conf) SetConfig(ctx context.Context, visitor *interfaces.Visitor, configKeys map[interfaces.ConfigKey]bool, cfg interfaces.Config) (err error) {
	c.trace.SetInternalSpanName("逻辑层-设置认证配置")
	newCtx, span := c.trace.AddInternalTrace(ctx)
	defer func() { c.trace.TelemetrySpanEnd(span, err) }()

	if visitor != nil && visitor.ID != "" {
		err = c.checkAdmin(newCtx, visitor)
		if err != nil {
			return
		}
	}

	err = c.checkParams(configKeys, cfg)
	if err != nil {
		return
	}

	var tx *sql.Tx
	tx, err = c.pool.Begin()
	if err != nil {
		c.log.Errorf("failed to begin transaction when SetConfig, err: %v", err)
		return err
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			err = tx.Commit()
			if err != nil {
				c.log.Errorf("SetConfig Transaction Commit Error:%v", err)
				return
			}
			// notify outbox推送线程
			c.ob.NotifyPushOutboxThread()
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				c.log.Errorf("SetConfig Transaction Rollback Error:%v", rollbackErr)
			}
		}
	}()

	err = c.cdb.SetConfig(configKeys, cfg)
	if err != nil {
		return
	}

	// 如果未更新SMSExpiration，则直接返回
	if !configKeys[interfaces.SMSExpiration] {
		return nil
	}
	message := map[string]interface{}{
		"anonymous_sms_expiration": cfg.SMSExpiration,
	}
	err = c.ob.AddOutboxInfo(logics.OutboxAnonymousSmsExpUpdated, message, tx)
	if err != nil {
		c.log.Errorf("failed to add outbox info when SetConfig, err: %v", err)
		return err
	}
	return
}

func (c *conf) checkAdmin(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	roleTypes, err := c.userMgnt.GetUserRolesByUserID(ctx, visitor, visitor.ID)
	if err != nil {
		return
	}

	return logics.CheckVisitorType(visitor, roleTypes, []interfaces.VisitorType{interfaces.RealName},
		[]interfaces.RoleType{interfaces.SuperAdmin, interfaces.SecurityAdmin})
}

func (c *conf) checkParams(configKeys map[interfaces.ConfigKey]bool, cfg interfaces.Config) (err error) {
	for k := range configKeys {
		switch k {
		case interfaces.RememberFor:
			if cfg.RememberFor < 0 {
				return rest.NewHTTPError("invalid remember_for", rest.BadRequest, nil)
			}
		case interfaces.RememberVisible:
		case interfaces.SMSExpiration:
			if cfg.SMSExpiration < 1 || cfg.SMSExpiration > 60 {
				return rest.NewHTTPErrorV2(rest.BadRequest, "invalid anonymous_sms_expiration")
			}
		default:
			// 此项不应该被匹配，如果匹配到此项则代表遍历项存在杂项
			return errors.New("this error is unexpected")
		}
	}

	return
}
