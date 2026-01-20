// Package logics contactor AnyShare 部门业务逻辑层
package logics

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/satori/uuid"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
)

type contactor struct {
	db            interfaces.DBContactor
	pool          *sqlx.DB
	ob            interfaces.LogicsOutbox
	logger        common.Logger
	messageBroker interfaces.DrivenMessageBroker
	event         interfaces.LogicsEvent
	trace         observable.Tracer
}

var (
	contactorOnce sync.Once
	cLogics       *contactor
)

// NewContactor 创建联系人组对象
func NewContactor() *contactor {
	contactorOnce.Do(func() {
		cLogics = &contactor{
			db:            dbContactor,
			pool:          dbPool,
			ob:            NewOutbox(OutboxBusinessContactor),
			logger:        common.NewLogger(),
			messageBroker: dnMessageBroker,
			event:         NewEvent(),
			trace:         common.SvcARTrace,
		}

		// 添加联系人组被删除事件消息
		cLogics.ob.RegisterHandlers(outboxContactorDeleted, cLogics.sendContactorDeletedMsg)

		// 注册用户删除事件
		cLogics.event.RegisterUserDeleted(cLogics.onUserDeleted)
	})

	return cLogics
}

// ConvertContactorName 根据Contactorid批量获取联系人组名称
func (d *contactor) ConvertContactorName(contactorIDs []string, bStrict bool) ([]interfaces.NameInfo, error) {
	infoArray := make([]interfaces.NameInfo, 0)
	if len(contactorIDs) == 0 {
		return infoArray, nil
	}

	// 去掉重复id
	RemoveDuplicatStrs(&contactorIDs)

	tempInfoMap, existIDs, err := d.db.GetContactorName(contactorIDs)
	if err != nil {
		return nil, err
	}

	// 如果严格模式， 且有联系人组不存在，则返回错误
	if bStrict && len(contactorIDs) != len(existIDs) {
		// 获取不存在的联系人组id
		notExistIDs := Difference(contactorIDs, existIDs)
		err := rest.NewHTTPError("record not exist", errors.ContactorNotFound,
			map[string]interface{}{"ids": notExistIDs})
		return nil, err
	}

	infoArray = append(infoArray, tempInfoMap...)
	return infoArray, nil
}

// DeleteContactor 删除联系人组
func (d *contactor) DeleteContactor(visitor *interfaces.Visitor, contactorID string) (err error) {
	// 检测contactorID是否合法
	_, err = uuid.FromString(contactorID)
	if err != nil {
		err = rest.NewHTTPError("param groupid is illegal", rest.BadRequest, nil)
		return
	}

	// 获取联系人组信息
	result, info, err := d.db.GetContactorInfo(contactorID)
	if err != nil {
		return
	}

	// 判断联系人组是否存在，并且是访问用户的
	if !result || info.UserID != visitor.ID {
		err = rest.NewHTTPError("group is not exist", rest.BadRequest, nil)
		return
	}

	// 获取事务处理器
	tx, err := d.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				d.logger.Errorf("DeleteContactor Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			d.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				d.logger.Errorf("DeleteContactor Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 删除联系人组成员
	err = d.db.DeleteContactorMembers([]string{contactorID}, tx)
	if err != nil {
		return
	}

	// 删除联系人组
	err = d.db.DeleteContactors([]string{contactorID}, tx)
	if err != nil {
		return
	}

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["ids"] = []string{contactorID}

	err = d.ob.AddOutboxInfo(outboxContactorDeleted, contentJSON, tx)
	if err != nil {
		d.logger.Errorf("Add Outbox Info err:%v", err)
		return
	}
	return
}

// onUserDeleted 用户删除事件-联系人组相关
func (d *contactor) onUserDeleted(userID string) (err error) {
	// 获取用户所有的联系人组
	contactorInfos, err := d.db.GetAllContactorInfos(userID)
	if err != nil {
		return err
	}

	// 获取所有的联系人组ID
	contactorIDs := []string{}
	for _, v := range contactorInfos {
		contactorIDs = append(contactorIDs, v.ContactorID)
	}

	// 获取事务处理器
	tx, err := d.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				d.logger.Errorf("onUserDeleted Transaction Commit Error:%v", err)
				return
			}

			// 更新联系人组的联系人数量信息
			err = d.db.UpdateContactorCount()
			if err != nil {
				d.logger.Errorf("UpdateContactorCount:%v", err)
				return
			}

			// notify outbox推送线程
			if len(contactorIDs) != 0 {
				d.ob.NotifyPushOutboxThread()
			}
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				d.logger.Errorf("onUserDeleted Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 删除联系人组成员
	err = d.db.DeleteContactorMembers(contactorIDs, tx)
	if err != nil {
		d.logger.Errorf("DeleteContactorMembers:%v", err)
		return
	}

	// 删除联系人组
	err = d.db.DeleteContactors(contactorIDs, tx)
	if err != nil {
		d.logger.Errorf("DeleteContactors:%v", err)
		return
	}

	// 在联系人组内删除此用户
	err = d.db.DeleteUserInContactors(userID, tx)
	if err != nil {
		d.logger.Errorf("DeleteUserInContactors:%v", err)
		return
	}

	// 如果存在被删除的联系人组，则发送联系人组删除事件消息
	if len(contactorIDs) != 0 {
		// 插入outbox 信息
		contentJSON := make(map[string]interface{})
		contentJSON["ids"] = contactorIDs

		err = d.ob.AddOutboxInfo(outboxContactorDeleted, contentJSON, tx)
		if err != nil {
			d.logger.Errorf("Add Outbox Info err:%v", err)
			return
		}
	}
	return
}

// sendContactorDeletedMsg 发送联系人组被删除消息
func (d *contactor) sendContactorDeletedMsg(content interface{}) error {
	info := content.(map[string]interface{})
	contactorIDs := []string{}
	for _, v := range info["ids"].([]interface{}) {
		contactorIDs = append(contactorIDs, v.(string))
	}

	if err := d.messageBroker.ContactorDeleted(contactorIDs); err != nil {
		return err
	}
	return nil
}

// GetContactorMembers 获取联系人组成员
func (d *contactor) GetContactorMembers(ctx context.Context, contactorIDs []string) (contactorMembers []interfaces.ContactorMemberInfo, err error) {
	//
	d.trace.SetInternalSpanName("业务逻辑-获取联系人组成员")
	newCtx, span := d.trace.AddInternalTrace(ctx)
	defer func() { d.trace.TelemetrySpanEnd(span, err) }()

	// id去重
	RemoveDuplicatStrs(&contactorIDs)

	// 获取存在的联系人组
	_, existIDs, err := d.db.GetContactorName(contactorIDs)
	if err != nil {
		return nil, err
	}

	// 获取联系人组成员
	contactorMemberIDs, err := d.db.GetContactorMemberIDs(newCtx, contactorIDs)
	if err != nil {
		return nil, err
	}

	// 转换联系人组成员信息
	contactorMembers = make([]interfaces.ContactorMemberInfo, 0)
	for _, contactorID := range existIDs {
		memberIDs, ok := contactorMemberIDs[contactorID]
		if !ok {
			memberIDs = []string{}
		}
		contactorMembers = append(contactorMembers, interfaces.ContactorMemberInfo{ContactorID: contactorID, MemberIDs: memberIDs})
	}

	return contactorMembers, nil
}
