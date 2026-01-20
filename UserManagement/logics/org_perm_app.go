// Package logics orgPermApp AnyShare 应用账户组织架构权限管理业务逻辑层
package logics

import (
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/mitchellh/mapstructure"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type orgPermApp struct {
	db      interfaces.DBOrgPermApp
	app     interfaces.LogicsApp
	pool    *sqlx.DB
	logger  common.Logger
	role    interfaces.LogicsRole
	eacpLog interfaces.DrivenEacpLog
	ob      interfaces.LogicsOutbox
}

var (
	opaOnce   sync.Once
	opaLogics *orgPermApp
)

// NewOrgPermApp 创建新的app perm对象
func NewOrgPermApp() *orgPermApp {
	opaOnce.Do(func() {
		opaLogics = &orgPermApp{
			db:      dbOrgPermApp,
			app:     NewApp(),
			pool:    dbPool,
			logger:  common.NewLogger(),
			role:    NewRole(),
			eacpLog: dnEacpLog,
			ob:      NewOutbox(OutboxBusinessOrgPermApp),
		}

		opaLogics.ob.RegisterHandlers(outboxOrgPermAppAddedLog, opaLogics.sendAddOrgPermAppAuditLog)
		opaLogics.ob.RegisterHandlers(outboxOrgPermAppUpdatedLog, opaLogics.sendUpdateOrgPermAppAuditLog)
		opaLogics.ob.RegisterHandlers(outboxOrgPermAppDeletedLog, opaLogics.sendDeleteOrgPermAppAuditLog)
	})

	return opaLogics
}

// UpdateAppName 更新应用账户组织架构权限表内应用账户名称
func (o *orgPermApp) UpdateAppName(info *interfaces.AppInfo) error {
	return o.db.UpdateAppName(info)
}

// GetAppOrgPerm 获取应用账户组织架构管理权限
func (o *orgPermApp) GetAppOrgPerm(visitor *interfaces.Visitor, id string, objects []interfaces.OrgType) (infos []interfaces.AppOrgPerm, err error) {
	// 检查objects是否重复
	orgTypes := make(map[interfaces.OrgType]bool)
	for _, v := range objects {
		orgTypes[v] = true
	}

	if len(orgTypes) != len(objects) {
		err = rest.NewHTTPError("org type is not uniqued", rest.BadRequest, nil)
		return
	}

	// 检查调用者权限
	err = o.checkManageAuthority(visitor)
	if err != nil {
		return
	}

	// 检查id是否存在
	_, err = o.app.GetApp(id)
	if err != nil {
		return nil, err
	}

	// 获取数据
	outInfos, err := o.db.GetAppPermByID(id)
	if err != nil {
		return nil, err
	}

	// 整理数据
	infos = make([]interfaces.AppOrgPerm, 0, len(outInfos))
	for _, v := range objects {
		data, ok := outInfos[v]
		if ok {
			infos = append(infos, data)
		}
	}

	return
}

// DeleteAppOrgPerm 删除应用账户组织架构管理权限
func (o *orgPermApp) DeleteAppOrgPerm(visitor *interfaces.Visitor, id string, objects []interfaces.OrgType) (err error) {
	// 检查objects是否重复
	orgTypes := make(map[interfaces.OrgType]bool)
	for _, v := range objects {
		orgTypes[v] = true
	}

	if len(orgTypes) != len(objects) {
		err = rest.NewHTTPError("org type is not uniqued", rest.BadRequest, nil)
		return
	}

	// 检查调用者权限
	err = o.checkManageAuthority(visitor)
	if err != nil {
		return
	}

	// 检查id是否存在
	_, err = o.app.GetApp(id)
	if err != nil {
		return err
	}

	// 获取数据
	outInfos, err := o.db.GetAppPermByID(id)
	if err != nil {
		return err
	}

	// 删除权限信息
	err = o.db.DeleteAppOrgPerm(id, objects)

	// 发送日志
	go func() {
		// 获取事务处理器
		var err1 error
		tx, err1 := o.pool.Begin()
		if err1 != nil {
			o.logger.Errorf("DeleteAppOrgPerm pool Begin Error:%v", err1)
			return
		}

		// 异常时Rollback
		defer func() {
			switch err1 {
			case nil:
				// 提交事务
				if err1 = tx.Commit(); err1 != nil {
					o.logger.Errorf("DeleteAppOrgPerm Transaction Commit Error:%v", err1)
					return
				}

				o.ob.NotifyPushOutboxThread()
			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					o.logger.Errorf("DeleteAppOrgPerm Rollback err:%v", rollbackErr)
				}
			}
		}()

		for _, v := range objects {
			if _, ok := outInfos[v]; !ok {
				continue
			}

			contentJSON := make(map[string]interface{})
			contentJSON["visitor"] = *visitor
			contentJSON["perm"] = outInfos[v]
			err1 = o.ob.AddOutboxInfo(outboxOrgPermAppDeletedLog, contentJSON, tx)
			if err1 != nil {
				o.logger.Errorf("DeleteAppOrgPerm AddOutboxInfo err:%v", err)
			}
		}
	}()
	return
}

func (o *orgPermApp) sendDeleteOrgPermAppAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		o.logger.Errorf("sendDeleteOrgPermAppAuditLog mapstructure.Decode err:%v", err)
		return
	}

	i := interfaces.AppOrgPerm{}
	err = mapstructure.Decode(info["perm"], &i)
	if err != nil {
		o.logger.Errorf("sendDeleteOrgPermAppAuditLog log_info mapstructure.Decode err:%v", err)
		return
	}

	err = o.eacpLog.OpDeleteOrgPermAppLog(&v, &i)
	if err != nil {
		o.logger.Errorf("sendDeleteOrgPermAppAuditLog err:%v", err)
	}
	return err
}

// SetAppOrgPerm 设置应用账户组织架构管理权限
func (o *orgPermApp) SetAppOrgPerm(visitor *interfaces.Visitor, id string, infos []interfaces.AppOrgPerm) (err error) {
	// 检查调用者权限
	err = o.checkManageAuthority(visitor)
	if err != nil {
		return
	}

	// 检查id是否存在
	appInfo, err := o.app.GetApp(id)
	if err != nil {
		return err
	}

	// 获取已有权限应用账户id
	outPerms, err := o.db.GetAppPermByID(id)
	if err != nil {
		return err
	}

	// 处理权限
	insertData, updateData, err := o.handlePermInfo(id, appInfo.Name, outPerms, infos)
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
				o.logger.Errorf("SetAppOrgPerm Transaction Commit Error:%v", err)
				return
			}

			o.ob.NotifyPushOutboxThread()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				o.logger.Errorf("SetAppOrgPerm Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 增加信息
	for _, v := range insertData {
		err = o.db.AddAppOrgPerm(v, tx)
		if err != nil {
			return err
		}

		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["perm"] = v
		err = o.ob.AddOutboxInfo(outboxOrgPermAppAddedLog, contentJSON, tx)
		if err != nil {
			return err
		}
	}

	// 更新信息
	for _, v := range updateData {
		err = o.db.UpdateAppOrgPerm(v, tx)
		if err != nil {
			return err
		}

		contentJSON := make(map[string]interface{})
		contentJSON["visitor"] = *visitor
		contentJSON["perm"] = v
		err = o.ob.AddOutboxInfo(outboxOrgPermAppUpdatedLog, contentJSON, tx)
		if err != nil {
			return err
		}
	}

	return
}

func (o *orgPermApp) sendAddOrgPermAppAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		o.logger.Errorf("sendAddOrgPermAppAuditLog mapstructure.Decode err:%v", err)
		return
	}

	i := interfaces.AppOrgPerm{}
	err = mapstructure.Decode(info["perm"], &i)
	if err != nil {
		o.logger.Errorf("sendAddOrgPermAppAuditLog log_info mapstructure.Decode err:%v", err)
		return
	}

	err = o.eacpLog.OpAddOrgPermAppLog(&v, &i)
	if err != nil {
		o.logger.Errorf("sendAddOrgPermAppAuditLog err:%v", err)
	}
	return err
}

func (o *orgPermApp) sendUpdateOrgPermAppAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		o.logger.Errorf("sendUpdateOrgPermAppAuditLog mapstructure.Decode err:%v", err)
		return
	}

	i := interfaces.AppOrgPerm{}
	err = mapstructure.Decode(info["perm"], &i)
	if err != nil {
		o.logger.Errorf("sendUpdateOrgPermAppAuditLog log_info mapstructure.Decode err:%v", err)
		return
	}

	err = o.eacpLog.OpUpdateOrgPermAppLog(&v, &i)
	if err != nil {
		o.logger.Errorf("sendUpdateOrgPermAppAuditLog err:%v", err)
	}
	return err
}

// checkManageAuthority 检查应用账户权限管理权限
func (o *orgPermApp) checkManageAuthority(visitor *interfaces.Visitor) (err error) {
	if visitor.Type == interfaces.RealName {
		return checkUserRole(o.role, visitor.ID, []interfaces.Role{interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin})
	}
	return rest.NewHTTPError("this user do not has the authority", rest.Forbidden, nil)
}

// handlePermInfo 处理权限信息
func (o *orgPermApp) handlePermInfo(id, name string, currentPerms map[interfaces.OrgType]interfaces.AppOrgPerm, insertPerms []interfaces.AppOrgPerm) (
	insertData []interfaces.AppOrgPerm, updateData []interfaces.AppOrgPerm, err error) {
	updateData = make([]interfaces.AppOrgPerm, 0)
	insertData = make([]interfaces.AppOrgPerm, 0)

	// 检查配置的权限
	for _, v := range insertPerms {
		temp := v
		temp.EndTime = -1
		temp.Name = name

		// 检查id和subject是否匹配
		if v.Subject != id {
			err = rest.NewHTTPError("subject is not same as app id", rest.BadRequest, nil)
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
