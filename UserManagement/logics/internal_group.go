// Package logics group AnyShare 内部组业务逻辑层
package logics

import (
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/satori/uuid"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type internalGroup struct {
	groupDB       interfaces.DBInternalGroup
	groupMemberDB interfaces.DBInternalGroupMember
	userDB        interfaces.DBUser
	pool          *sqlx.DB
	logger        common.Logger
	ob            interfaces.LogicsOutbox
	messageBroker interfaces.DrivenMessageBroker
}

var (
	igOnce   sync.Once
	igLogics *internalGroup
)

// NewInternalGroup 创建新的internal group对象
func NewInternalGroup() *internalGroup {
	igOnce.Do(func() {
		igLogics = &internalGroup{
			groupMemberDB: dbInternalGroupMember,
			groupDB:       dbInternalGroup,
			userDB:        dbUser,
			pool:          dbPool,
			logger:        common.NewLogger(),
			ob:            NewOutbox(OutboxBusinessInternalGroup),
			messageBroker: dnMessageBroker,
		}

		igLogics.ob.RegisterHandlers(outboxInternalGroupDeleted, igLogics.sendInternalGroupDeletedInfo)
	})

	return igLogics
}

// AddGroup 增加内部组
func (g *internalGroup) AddGroup() (id string, err error) {
	id = uuid.NewV4().String()
	err = g.groupDB.Add(id)
	return id, err
}

// DeleteGroup 删除内部组
func (g *internalGroup) DeleteGroup(ids []string) (err error) {
	// 获取存在的部门
	outInfos, err := g.groupDB.Get(ids)
	if err != nil {
		return
	}

	existIDs := make([]string, 0, len(outInfos))
	for k := range outInfos {
		existIDs = append(existIDs, k)
	}

	if len(existIDs) == 0 {
		return nil
	}

	// 获取事务处理器
	tx, err := g.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				g.logger.Errorf("DeleteGroup Transaction Commit Error:%v", err)
				return
			}

			// notify outbox推送线程
			g.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				g.logger.Errorf("DeleteGroup Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 删除内部组成员
	err = g.groupMemberDB.DeleteAll(existIDs, tx)
	if err != nil {
		return
	}

	// 删除内部组
	err = g.groupDB.Delete(existIDs, tx)
	if err != nil {
		return
	}

	// 插入outbox 信息
	contentJSON := make(map[string]interface{})
	contentJSON["ids"] = existIDs

	err = g.ob.AddOutboxInfo(outboxInternalGroupDeleted, contentJSON, tx)
	if err != nil {
		g.logger.Errorf("Add Outbox Info err:%v", err)
	}
	return
}

// GetGroupMemberByID 根据内部组ID获取成员ID
func (g *internalGroup) GetGroupMemberByID(id string) (outInfos []interfaces.InternalGroupMember, err error) {
	// 检查内部组是否存在
	out, err := g.groupDB.Get([]string{id})
	if err != nil {
		return
	}

	if _, ok := out[id]; !ok {
		err = rest.NewHTTPError("internal group do not exist", rest.URINotExist, nil)
		return
	}

	// 获取内部组成员
	return g.groupMemberDB.Get(id)
}

// UpdateMembers 更新成员
func (g *internalGroup) UpdateMembers(id string, infos []interfaces.InternalGroupMember) (err error) {
	// 检查内部组是否存在
	out, err := g.groupDB.Get([]string{id})
	if err != nil {
		return
	}

	if _, ok := out[id]; !ok {
		err = rest.NewHTTPError("internal group do not exist", rest.URINotExist, nil)
		return
	}

	// 检查成员是否重复
	userIDs := make([]string, 0, len(infos))
	for _, v := range infos {
		userIDs = append(userIDs, v.ID)
	}
	RemoveDuplicatStrs(&userIDs)
	if len(userIDs) != len(infos) {
		err = rest.NewHTTPError("internal group member do not unique", rest.BadRequest, nil)
		return
	}

	// 检查成员是否存在
	_, existIDs, err := g.userDB.GetUserName(userIDs)
	if err != nil {
		return err
	}

	if len(existIDs) != len(userIDs) {
		notExists := Difference(userIDs, existIDs)
		err = rest.NewHTTPError("some members are not existing", rest.BadRequest, map[string]interface{}{"ids": notExists})
		return
	}

	// 获取事务处理器
	tx, err := g.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				g.logger.Errorf("internal group AddMembers Transaction Commit Error:%v", err)
				return
			}
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				g.logger.Errorf("internal group AddMembers Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 删除组所有成员
	err = g.groupMemberDB.DeleteAll([]string{id}, tx)
	if err != nil {
		return
	}

	// 添加成员
	return g.groupMemberDB.Add(id, infos, tx)
}

// sendInternalGroupDeletedInfo 发送内部组被删除消息
func (g *internalGroup) sendInternalGroupDeletedInfo(content interface{}) error {
	info := content.(map[string]interface{})
	ids := make([]string, 0)
	for _, v := range info["ids"].([]interface{}) {
		ids = append(ids, v.(string))
	}

	if err := g.messageBroker.InternalGroupDeleted(ids); err != nil {
		return err
	}
	return nil
}

// GetBelongGroupByID 获取用户成员所属内部组
func (g *internalGroup) GetBelongGroups(info interfaces.InternalGroupMember) (ids []string, err error) {
	return g.groupMemberDB.GetBelongGroups(info)
}
