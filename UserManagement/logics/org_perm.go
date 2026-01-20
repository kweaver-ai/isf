// Package logics orgPerm AnyShare 应用账户组织架构权限管理业务逻辑层
package logics

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type orgPerm struct {
	db     interfaces.DBOrgPerm
	user   interfaces.DBUser
	pool   *sqlx.DB
	logger common.Logger
	event  interfaces.LogicsEvent
	trace  observable.Tracer
}

var (
	opOnce   sync.Once
	opLogics *orgPerm
)

// NewOrgPerm 创建新的org perm对象
func NewOrgPerm() *orgPerm {
	opOnce.Do(func() {
		opLogics = &orgPerm{
			db:     dbOrgPerm,
			pool:   dbPool,
			logger: common.NewLogger(),
			trace:  common.SvcARTrace,
			event:  NewEvent(),
			user:   dbUser,
		}

		opLogics.event.RegisterUserDeleted(opLogics.onUserDeleted)
		opLogics.event.RegisterUserNameChanged(opLogics.UpdateName)
	})

	return opLogics
}

// UpdateName 更新账户组织架构权限表内账户名称
func (o *orgPerm) UpdateName(id, newName string) (err error) {
	return o.db.UpdateName(id, newName)
}

// DeleteOrgPerm 删除账户组织架构管理权限
func (o *orgPerm) DeleteOrgPerm(ctx context.Context, subjectID string, subjectType interfaces.VisitorType, objects []interfaces.OrgType) (err error) {
	// trace
	o.trace.SetInternalSpanName("业务逻辑-删除账户组织架构管理权限")
	newCtx, span := o.trace.AddInternalTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	if subjectType != interfaces.RealName {
		return rest.NewHTTPErrorV2(rest.BadRequest, "subject type is not supported")
	}

	// 检查objects是否重复
	orgTypes := make(map[interfaces.OrgType]bool)
	for _, v := range objects {
		orgTypes[v] = true
	}

	if len(orgTypes) != len(objects) {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "org type is not uniqued")
		return
	}

	// 删除权限信息
	err = o.db.DeleteOrgPerm(newCtx, subjectID, objects)
	return
}

// SetOrgPerm 设置账户组织架构管理权限
func (o *orgPerm) SetOrgPerm(ctx context.Context, subjectID string, subjectType interfaces.VisitorType, infos []interfaces.OrgPerm) (err error) {
	// trace
	o.trace.SetInternalSpanName("业务逻辑-设置账户组织架构管理权限")
	newCtx, span := o.trace.AddInternalTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	// 检查id是否存在
	strName := ""
	if subjectType == interfaces.RealName {
		var userInfo []interfaces.UserDBInfo
		userInfo, err = o.user.GetUserDBInfo2(newCtx, []string{subjectID})
		if err != nil {
			return err
		}

		if len(userInfo) != 1 {
			return rest.NewHTTPErrorV2(rest.URINotExist, "user not found")
		}
		strName = userInfo[0].Name
	} else {
		return rest.NewHTTPErrorV2(rest.BadRequest, "subject type is not supported")
	}

	// 获取已有权限应用账户id
	outPerms, err := o.db.GetPermByID(newCtx, subjectID)
	if err != nil {
		return err
	}

	// 处理权限
	insertData, updateData, err := o.handlePermInfo(subjectID, strName, outPerms, infos)
	if err != nil {
		return
	}

	// 获取事务处理器
	tx, err := o.pool.Begin()
	if err != nil {
		return
	}

	// 异常时Rollback
	defer func() {
		switch err {
		case nil:
			// 提交事务
			if err = tx.Commit(); err != nil {
				o.logger.Errorf("SetOrgPerm Transaction Commit Error:%v", err)
				return
			}

		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				o.logger.Errorf("SetOrgPerm Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 增加信息
	for _, v := range insertData {
		err = o.db.AddOrgPerm(newCtx, v, tx)
		if err != nil {
			return err
		}
	}

	// 更新信息
	for _, v := range updateData {
		err = o.db.UpdateOrgPerm(newCtx, v, tx)
		if err != nil {
			return err
		}
	}

	return
}

// handlePermInfo 处理权限信息
func (o *orgPerm) handlePermInfo(id, name string, currentPerms map[interfaces.OrgType]interfaces.OrgPerm, insertPerms []interfaces.OrgPerm) (
	insertData []interfaces.OrgPerm, updateData []interfaces.OrgPerm, err error) {
	updateData = make([]interfaces.OrgPerm, 0)
	insertData = make([]interfaces.OrgPerm, 0)

	// 检查配置的权限
	for _, v := range insertPerms {
		temp := v
		temp.EndTime = -1
		temp.Name = name

		// 检查id和subject是否匹配
		if v.SubjectID != id {
			err = rest.NewHTTPErrorV2(rest.BadRequest, "subject is not same as subject id")
			return
		}

		// 如果存在，则更新，否则插入
		if _, ok := currentPerms[v.Object]; ok {
			updateData = append(updateData, temp)
		} else {
			insertData = append(insertData, temp)
		}
	}
	return
}

// CheckPerms 账户组织架构管理权限检查
func (o *orgPerm) CheckPerms(ctx context.Context, subjectID string, orgTyp interfaces.OrgType, checkPerm interfaces.OrgPermValue) (result bool, err error) {
	// trace
	o.trace.SetInternalSpanName("业务逻辑-账户组织架构管理权限检查")
	newCtx, span := o.trace.AddInternalTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	perms, err := o.db.GetPermByID(newCtx, subjectID)
	if err != nil {
		return false, err
	}

	if v, ok := perms[orgTyp]; ok {
		if v.Value&checkPerm != 0 {
			return true, nil
		}
		return false, nil
	}

	return false, nil
}

// onUserDeleted 用户被删除
func (o *orgPerm) onUserDeleted(userID string) (err error) {
	return o.db.DeleteOrgPermByID(userID)
}
