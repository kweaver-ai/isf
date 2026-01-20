// Package accesstokenperm 逻辑层
package accesstokenperm

import (
	"context"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/mitchellh/mapstructure"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	aOnce sync.Once
	a     *accessTokenPerm
)

type accessTokenPerm struct {
	db         interfaces.DBAccessTokenPerm
	userMgnt   interfaces.DnUserManagement
	hydraAdmin interfaces.DnHydraAdmin
	eacpLog    interfaces.DnEacpLog
	ob         interfaces.LogicsOutbox
	pool       *sqlx.DB
	logger     common.Logger
}

// NewAccessTokenPerm 创建accessTokenPerm处理对象
func NewAccessTokenPerm() *accessTokenPerm {
	aOnce.Do(func() {
		a = &accessTokenPerm{
			db:         logics.DBAccessTokenPerm,
			userMgnt:   logics.DnUserManagement,
			hydraAdmin: logics.DnHydraAdmin,
			eacpLog:    logics.DnEacpLog,
			ob:         logics.NewOutbox(logics.OutboxBusinessAccessTokenPerm),
			pool:       logics.DBPool,
			logger:     common.NewLogger(),
		}

		a.ob.RegisterHandlers(logics.OutboxSetAppAccessTokenPermLog, a.setAppPermAuditLog)
		a.ob.RegisterHandlers(logics.OutboxDeleteAppAccessTokenPermLog, a.deleteAppPermAuditLog)
	})

	return a
}

func (a *accessTokenPerm) checkAdmin(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	var roleTypes []interfaces.RoleType
	// 实名用户获取对应角色信息
	if visitor.Type == interfaces.RealName {
		roleTypes, err = a.userMgnt.GetUserRolesByUserID(ctx, visitor, visitor.ID)
		if err != nil {
			return
		}
	}

	return logics.CheckVisitorType(visitor, roleTypes, []interfaces.VisitorType{interfaces.RealName},
		[]interfaces.RoleType{interfaces.SuperAdmin, interfaces.SystemAdmin})
}

// SetAccessTokenPerm 设置应用账户获取任意用户访问令牌权限
func (a *accessTokenPerm) SetAppAccessTokenPerm(ctx context.Context, visitor *interfaces.Visitor, appID string) (err error) {
	if visitor != nil && visitor.ID != "" {
		err = a.checkAdmin(ctx, visitor)
		if err != nil {
			return
		}
	}

	// 检查应用账户是否存在
	appInfo, err := a.userMgnt.GetAppInfo(ctx, visitor, appID)
	if err != nil {
		return
	}

	hasPerm, err := a.CheckAppAccessTokenPerm(appID)
	if err != nil {
		return err
	}

	if hasPerm {
		return
	}

	// 更新客户端client_type，使所有作为用户代理的客户端，metadata保持一致，便于后续扩展
	// 更新客户端grant_types，使其能够使用断言获取token
	err = a.hydraAdmin.SetAppAsUserAgent(appID)
	if err != nil {
		return
	}

	err = a.db.AddAppAccessTokenPerm(appID)
	if err != nil {
		return
	}

	// 记录日志
	if visitor != nil && visitor.ID != "" {
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
					a.logger.Errorf("SetAppAccessTokenPerm Transaction Commit Error:%v", err)
					return
				}

				// notify outbox推送线程
				a.ob.NotifyPushOutboxThread()
			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					a.logger.Errorf("SetAppAccessTokenPerm Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 插入outbox 信息
		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["app_name"] = appInfo.Name

		err = a.ob.AddOutboxInfo(logics.OutboxSetAppAccessTokenPermLog, contentJSON, tx)
		if err != nil {
			a.logger.Errorf("SetAppAccessTokenPerm Add Outbox Info err:%v", err)
			return err
		}
	}

	return
}

// setAppPermAuditLog 设置权限审计日志
func (a *accessTokenPerm) setAppPermAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("setAppPermAuditLog mapstructure.Decode err:%v", err)
		return
	}
	return a.eacpLog.OpSetAppAccessTokenPerm(&v, info["app_name"].(string))
}

// DeleteAppAccessTokenPerm 删除应用账户获取任意用户访问令牌权限，visitor 不为nil时需检验访问者身份并记录管理日志
func (a *accessTokenPerm) DeleteAppAccessTokenPerm(ctx context.Context, visitor *interfaces.Visitor, appID string) (err error) {
	if visitor != nil && visitor.ID != "" {
		err = a.checkAdmin(ctx, visitor)
		if err != nil {
			return
		}
	}

	// 检查应用账户是否存在
	appInfo, err := a.userMgnt.GetAppInfo(ctx, visitor, appID)
	if err != nil {
		return
	}

	err = a.db.DeleteAppAccessTokenPerm(appID)
	if err != nil {
		return
	}

	// 记录日志
	if visitor != nil && visitor.ID != "" {
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
					a.logger.Errorf("DeleteAppAccessTokenPerm Transaction Commit Error:%v", err)
					return
				}

				// notify outbox推送线程
				a.ob.NotifyPushOutboxThread()
			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					a.logger.Errorf("DeleteAppAccessTokenPerm Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 插入outbox 信息
		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["app_name"] = appInfo.Name

		err = a.ob.AddOutboxInfo(logics.OutboxDeleteAppAccessTokenPermLog, contentJSON, tx)
		if err != nil {
			a.logger.Errorf("DeleteAppAccessTokenPerm Add Outbox Info err:%v", err)
			return err
		}
	}

	return
}

// deleteAppPermAuditLog 删除权限审计日志
func (a *accessTokenPerm) deleteAppPermAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		a.logger.Errorf("deleteAppPermAuditLog mapstructure.Decode err:%v", err)
		return
	}
	return a.eacpLog.OpDeleteAppAccessTokenPerm(&v, info["app_name"].(string))
}

// AppDeleted 删除应用账户获取任意用户访问令牌权限
func (a *accessTokenPerm) AppDeleted(appID string) error {
	return a.db.DeleteAppAccessTokenPerm(appID)
}

func (a *accessTokenPerm) CheckAppAccessTokenPerm(appID string) (bool, error) {
	return a.db.CheckAppAccessTokenPerm(appID)
}

func (a *accessTokenPerm) GetAllAppAccessTokenPerm(ctx context.Context, visitor *interfaces.Visitor) (permApps []string, err error) {
	if visitor != nil && visitor.ID != "" {
		err = a.checkAdmin(ctx, visitor)
		if err != nil {
			return
		}
	}

	return a.db.GetAllAppAccessTokenPerm()
}
