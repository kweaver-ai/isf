// Package logics group AnyShare 用户组业务逻辑层
package logics

import (
	"context"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/uuid"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
)

type group struct {
	groupMemberDB interfaces.DBGroupMember
	userDB        interfaces.DBUser
	departmentDB  interfaces.DBDepartment
	groupDB       interfaces.DBGroup
	eacpLog       interfaces.DrivenEacpLog
	messageBroker interfaces.DrivenMessageBroker
	role          interfaces.LogicsRole
	orgPermApp    interfaces.DBOrgPermApp
	event         interfaces.LogicsEvent
	trace         observable.Tracer
	pool          *sqlx.DB
	logger        common.Logger
	i18n          *common.I18n
	ob            interfaces.LogicsOutbox
}

var (
	gOnce   sync.Once
	gLogics *group
)

// NewGroup 创建新的group对象
func NewGroup() *group {
	gOnce.Do(func() {
		gLogics = &group{
			groupMemberDB: dbGroupMember,
			groupDB:       dbGroup,
			userDB:        dbUser,
			departmentDB:  dbDepartment,
			eacpLog:       dnEacpLog,
			messageBroker: dnMessageBroker,
			role:          NewRole(),
			orgPermApp:    dbOrgPermApp,
			event:         NewEvent(),
			trace:         common.SvcARTrace,
			pool:          dbPool,
			logger:        common.NewLogger(),
			i18n: common.NewI18n(common.I18nMap{
				i18nIDObjectsInUnDistributeUserGroup: {
					interfaces.SimplifiedChinese:  "未分配组",
					interfaces.TraditionalChinese: "未分配組",
					interfaces.AmericanEnglish:    "Unassigned Group",
				},
				i18nIDObjectsInUserNotFound: {
					interfaces.SimplifiedChinese:  "用户不存在",
					interfaces.TraditionalChinese: "用戶不存在",
					interfaces.AmericanEnglish:    "This user does not exist",
				},
				i18nIDObjectsInDepartNotFound: {
					interfaces.SimplifiedChinese:  "部门不存在",
					interfaces.TraditionalChinese: "部門不存在",
					interfaces.AmericanEnglish:    "This department does not exist",
				},
				i18nIDObjectsInGroupNotFound: {
					interfaces.SimplifiedChinese:  "用户组不存在",
					interfaces.TraditionalChinese: "用戶組不存在",
					interfaces.AmericanEnglish:    "This group does not exist",
				},
			}),
			ob: NewOutbox(OutboxBusinessGroup),
		}

		gLogics.event.RegisterDeptDeleted(gLogics.DeleteGroupMemberByMemberID)

		gLogics.ob.RegisterHandlers(outboxGroupAddedLog, gLogics.sendAddGroupAuditLog)
		gLogics.ob.RegisterHandlers(outboxGroupDeletedLog, gLogics.sendDeleteGroupAuditLog)
		gLogics.ob.RegisterHandlers(outboxGroupModifiedLog, gLogics.sendModifyGroupAuditLog)
		gLogics.ob.RegisterHandlers(outboxGroupMembersAddedLog, gLogics.sendAddGroupMembersAuditLog)
		gLogics.ob.RegisterHandlers(outboxGroupMembersDeletedLog, gLogics.sendDeleteGroupMembersAuditLog)
	})

	return gLogics
}

// UserMatch 组内用户匹配
func (g *group) UserMatch(ctx context.Context, visitor *interfaces.Visitor, groupID, userName string) (
	exist bool, uInfo interfaces.GroupMemberInfo, mInfo []interfaces.GroupMemberInfo, err error) {
	// trace
	g.trace.SetInternalSpanName("业务逻辑-组内用户匹配")
	newCtx, span := g.trace.AddInternalTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	// 判断用户是否有用户组的列举权限（实名用户和应用账户）
	err = g.checkGroupGetAuth2(newCtx, visitor)
	if err != nil {
		return
	}

	// 判断组织是否存在
	var groupInfo interfaces.GroupInfo
	groupInfo, err = g.groupDB.GetGroupByID2(newCtx, groupID)
	if err != nil {
		return
	}
	if groupInfo.ID == "" {
		err = rest.NewHTTPErrorV2(rest.URINotExist, "group not exist")
		return
	}

	// 判断用户是否存在，如果不存在，则返回
	userInfo, err := g.userDB.GetUserInfoByName(newCtx, userName)
	if err != nil {
		return
	}

	if userInfo.ID == "" {
		return false, uInfo, nil, nil
	}

	// 获取用户组成员信息
	memberInfos, err := g.groupMemberDB.GetGroupMembersByGroupIDs2(newCtx, []string{groupID})
	if err != nil {
		return
	}

	// 获取用户所有的部门
	paths, err := g.userDB.GetUsersPath2(newCtx, []string{userInfo.ID})
	if err != nil {
		return
	}
	userParentDeps := make(map[string]bool)
	for _, v := range paths[userInfo.ID] {
		tempPaths := strings.Split(v, "/")
		for _, v1 := range tempPaths {
			userParentDeps[v1] = true
		}
	}

	// 判断用户在那些成员内
	existGroupMember := make([]interfaces.GroupMemberInfo, 0)
	userIDs := []string{userInfo.ID}
	departIDs := []string{}
	for _, v := range memberInfos[groupID] {
		if v.MemberType-1 == 0 && v.ID == userInfo.ID {
			v.Name = userName
			existGroupMember = append(existGroupMember, v)
		} else if v.MemberType-2 == 0 {
			if _, ok := userParentDeps[v.ID]; ok {
				existGroupMember = append(existGroupMember, v)
				departIDs = append(departIDs, v.ID)
			}
		}
	}

	// 获取成员和用户的path
	memberPaths, departNames, err := g.handlePathParams(newCtx, visitor, userIDs, departIDs)
	if err != nil {
		return
	}

	// 数据数据
	uInfo.ID = userInfo.ID
	uInfo.ParentDeps = memberPaths[userInfo.ID]
	mInfo = make([]interfaces.GroupMemberInfo, 0, len(existGroupMember))
	for _, v := range existGroupMember {
		v.ParentDeps = memberPaths[v.ID]
		if v.MemberType-2 == 0 {
			v.Name = departNames[v.ID]
		}
		mInfo = append(mInfo, v)
	}

	exist = true
	return
}

// searchInAllGroupOrg 组内用户搜索
func (g *group) SearchInAllGroupOrg(ctx context.Context, visitor *interfaces.Visitor, groupID, key string, offset, limit int) (
	num int, userIDs []string, uInfos map[string]interfaces.GroupMemberInfo, mInfos map[string][]interfaces.GroupMemberInfo, err error) {
	// trace
	g.trace.SetInternalSpanName("业务逻辑-组内用户搜索")
	newCtx, span := g.trace.AddInternalTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	// 判断用户是否有用户组的列举权限（实名用户和应用账户）
	err = g.checkGroupGetAuth2(newCtx, visitor)
	if err != nil {
		return
	}

	// 判断组织是否存在
	var groupInfo interfaces.GroupInfo
	groupInfo, err = g.groupDB.GetGroupByID2(newCtx, groupID)
	if err != nil {
		return
	}
	if groupInfo.ID == "" {
		err = rest.NewHTTPErrorV2(rest.URINotExist, "group not exist")
		return
	}

	// 获取所有的名称类似用户
	userInfos, err := g.userDB.SearchUserInfoByName(newCtx, key)
	if err != nil {
		return
	}

	if len(userInfos) == 0 {
		return 0, nil, nil, nil, nil
	}

	// 获取用户组成员信息
	memberInfos, err := g.groupMemberDB.GetGroupMembersByGroupIDs2(newCtx, []string{groupID})
	if err != nil {
		return
	}

	// 获取用户所在的所有部门，包括父部门
	allUserIDs := make([]string, 0)
	allUserNames := make(map[string]string)
	for k := range userInfos {
		allUserIDs = append(allUserIDs, userInfos[k].ID)
		allUserNames[userInfos[k].ID] = userInfos[k].Name
	}
	paths, err := g.userDB.GetUsersPath2(newCtx, allUserIDs)
	if err != nil {
		return
	}
	userParentDeps := make(map[string]map[string]bool)
	for k, v2 := range paths {
		tempDatas := make(map[string]bool)
		for _, v := range v2 {
			tempPaths := strings.Split(v, "/")
			for _, v1 := range tempPaths {
				tempDatas[v1] = true
			}
		}
		userParentDeps[k] = tempDatas
	}

	// 返回在特定组的用户和其所在的组成员，以及按照limit offset排序结果
	var departIDs []string
	var existGroupMember map[string][]interfaces.GroupMemberInfo
	num, userIDs, departIDs, existGroupMember = g.getExistUsersGroupMembers(allUserIDs, userParentDeps, memberInfos[groupID], offset, limit)

	// 获取用户和部门的路径
	memberPaths, departNames, err := g.handlePathParams(newCtx, visitor, userIDs, departIDs)
	if err != nil {
		return
	}

	// 数据数据
	uInfos = make(map[string]interfaces.GroupMemberInfo)
	mInfos = make(map[string][]interfaces.GroupMemberInfo)
	for _, v := range userIDs {
		tempGroupMember := existGroupMember[v]
		if len(tempGroupMember) == 0 {
			continue
		}

		uInfo := interfaces.GroupMemberInfo{
			ID:         v,
			Name:       allUserNames[v],
			ParentDeps: memberPaths[v],
		}
		uInfos[v] = uInfo

		mInfo := make([]interfaces.GroupMemberInfo, 0, len(tempGroupMember))
		for _, v1 := range tempGroupMember {
			v1.ParentDeps = memberPaths[v1.ID]
			if v1.MemberType-2 == 0 {
				v1.Name = departNames[v1.ID]
			} else {
				v1.Name = allUserNames[v]
			}
			mInfo = append(mInfo, v1)
		}
		mInfos[v] = mInfo
	}
	return
}

// 获取用户所在的组成员
func (g *group) getExistUsersGroupMembers(userIDs []string, userParentDeps map[string]map[string]bool, memberInfos []interfaces.GroupMemberInfo, offset, limit int) (
	num int, existUserIDs []string, existDepartIDs []string, existGroupMember map[string][]interfaces.GroupMemberInfo) {
	existGroupMember = make(map[string][]interfaces.GroupMemberInfo)
	existUserIDs = make([]string, 0)
	existDepartIDs = make([]string, 0)

	nValidIndex := 0
	for _, v11 := range userIDs {
		// 检查此用户是否是这个组成员，或者此用户在此组成员下
		tempMembers := make([]interfaces.GroupMemberInfo, 0)
		for _, v := range memberInfos {
			if v.MemberType-1 == 0 && v.ID == v11 {
				tempMembers = append(tempMembers, v)
			} else if v.MemberType-2 == 0 {
				if tempPath, ok := userParentDeps[v11]; ok {
					if _, ok1 := tempPath[v.ID]; ok1 {
						tempMembers = append(tempMembers, v)
					}
				}
			}
		}

		// 如果在此组，并且在limit offset范围内，则记录需要返回的用户
		if len(tempMembers) > 0 {
			if nValidIndex >= offset && len(existUserIDs) < limit {
				existUserIDs = append(existUserIDs, v11)
				existGroupMember[v11] = tempMembers
			}
			nValidIndex++
		}
	}

	// 有效用户总数
	num = nValidIndex

	// 获取limit offset范围内用户所在组成员中type为department的部门id
	for _, v := range existGroupMember {
		for _, v1 := range v {
			if v1.MemberType-2 == 0 {
				existDepartIDs = append(existDepartIDs, v1.ID)
			}
		}
	}
	RemoveDuplicatStrs(&existDepartIDs)

	return
}

// 获取用户组内所有成员
func (g *group) getGroupMembers(ctx context.Context, initalGroupIDs []string) (info []interfaces.GroupMemberInfo, err error) {
	// 获取初始组成员
	mapMemberInfos, err := g.groupMemberDB.GetGroupMembersByGroupIDs2(ctx, initalGroupIDs)
	if err != nil {
		return nil, err
	}

	// 获取所有的成员并且去重
	tempMemberInfos := make(map[string]interfaces.GroupMemberInfo)
	for _, v := range mapMemberInfos {
		for _, v1 := range v {
			tempMemberInfos[v1.ID] = v1
		}
	}

	info = make([]interfaces.GroupMemberInfo, 0, len(tempMemberInfos))
	for _, v := range tempMemberInfos {
		info = append(info, v)
	}
	return
}

// AddGroup 用户组创建
func (g *group) AddGroup(ctx context.Context, visitor *interfaces.Visitor, name, notes string, initalGroupIDs []string) (id string, err error) {
	// trace
	g.trace.SetInternalSpanName("业务逻辑-用户组创建")
	newCtx, span := g.trace.AddInternalTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	// 判断用户是否有用户组的管理权限（实名用户和应用账户）
	err = g.checkGroupManageAuth2(newCtx, visitor)
	if err != nil {
		return
	}
	// 参数检查
	err = g.checkName(name)
	if err != nil {
		return
	}
	// 中英文混合
	maxLen := 300
	if utf8.RuneCountInString(notes) > maxLen {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "param notes is illegal")
		return
	}

	// 检测用户组名是否重复
	var groupID string
	groupID, err = g.groupDB.GetGroupIDByName2(newCtx, name)
	if err != nil {
		return id, err
	}

	if groupID != "" {
		err = rest.NewHTTPErrorV2(errors.Conflict, "this group name is existing",
			rest.SetCodeStr(errors.StrConflictGroup))
		return id, err
	}

	// 创建userID
	id = uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		return id, err
	}

	// 检查用户组id是否重复
	initalGroupNum := len(initalGroupIDs)
	RemoveDuplicatStrs(&initalGroupIDs)
	if initalGroupNum != len(initalGroupIDs) && initalGroupNum != 0 {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "initial group ids duplicate")
		return id, err
	}

	// 检查用户组id是否存在
	_, existIDs, err := g.groupDB.GetGroupName2(newCtx, initalGroupIDs)
	if err != nil {
		return id, err
	}

	if len(existIDs) != len(initalGroupIDs) {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "initial group ids not exist")
		return id, err
	}

	// 获取初始组成员
	memberInfos, err := g.getGroupMembers(newCtx, initalGroupIDs)
	if err != nil {
		return id, err
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
				g.logger.Errorf("AddGroup Transaction Commit Error:%v", err)
				return
			}

			g.ob.NotifyPushOutboxThread()

		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				g.logger.Errorf("AddGroup Rollback err:%v", rollbackErr)
			}
		}
	}()

	// 创建用户组
	err = g.groupDB.AddGroup(newCtx, id, name, notes, tx)
	if err != nil {
		return id, err
	}

	// 添加用户组成员
	err = g.groupMemberDB.AddGroupMembers(newCtx, id, memberInfos, tx)
	if err != nil {
		return id, err
	}

	// 记录审计日志
	groupInfo := interfaces.GroupInfo{
		Name:  name,
		Notes: notes,
	}

	content := make(map[string]interface{})
	content["visitor"] = *visitor
	content["group_info"] = groupInfo
	err = g.ob.AddOutboxInfo(outboxGroupAddedLog, content, tx)
	return id, err
}

// sendAddGroupAuditLog 添加用户组发送审计消息
func (g *group) sendAddGroupAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	groupInfo := interfaces.GroupInfo{}
	err = mapstructure.Decode(info["group_info"], &groupInfo)
	if err != nil {
		g.logger.Errorf("sendAddGroupAuditLog group_info mapstructure.Decode err:%v", err)
		return
	}

	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		g.logger.Errorf("sendAddGroupAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = g.eacpLog.EacpLog(&v, interfaces.OpAddGroup, groupInfo)
	if err != nil {
		g.logger.Errorf("sendAddGroupAuditLog err:%v", err)
	}
	return err
}

// DeleteGroup 用户组删除
func (g *group) DeleteGroup(visitor *interfaces.Visitor, groupID string) (err error) {
	// 判断用户是否有用户组的管理权限（实名用户和应用账户）
	err = g.checkGroupManageAuth(visitor)
	if err != nil {
		return
	}

	// 判断用户组是否存在,如果存在 获取用户组信息
	var info interfaces.GroupInfo
	info, err = g.groupDB.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// 如果用户组不存在， 则不处理
	if info.ID != "" {
		// 删除用户组内成员
		err = g.groupMemberDB.DeleteGroupMemberByID(groupID)
		if err != nil {
			return err
		}

		// 删除用户组
		err = g.groupDB.DeleteGroup(groupID)
		if err != nil {
			return err
		}
	}

	// 消息发送
	go func() {
		if info.ID == "" {
			return
		}
		_ = g.messageBroker.Publish(interfaces.DeleteGroup, info.ID)
	}()

	// 记录日志
	go func() {
		// 如果用户组不存在 则不记录删除日志
		if info.ID == "" {
			return
		}

		// 获取事务处理器
		var err1 error
		tx, err1 := g.pool.Begin()
		if err1 != nil {
			g.logger.Errorf("DeleteGroup send audit log  pool begin error:%v", err1)
			return
		}

		// 异常时Rollback
		defer func() {
			switch err1 {
			case nil:
				// 提交事务
				if err1 = tx.Commit(); err1 != nil {
					g.logger.Errorf("DeleteGroup Transaction Commit Error:%v", err1)
					return
				}

				g.ob.NotifyPushOutboxThread()

			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					g.logger.Errorf("DeleteGroup Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 记录审计日志
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["group_info"] = info
		err1 = g.ob.AddOutboxInfo(outboxGroupDeletedLog, content, tx)
		if err1 != nil {
			g.logger.Errorf("DeleteGroup AddOutboxInfo Error:%v", err1)
		}
	}()
	return err
}

// sendDeleteGroupAuditLog 删除用户组发送审计消息
func (g *group) sendDeleteGroupAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	groupInfo := interfaces.GroupInfo{}
	err = mapstructure.Decode(info["group_info"], &groupInfo)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupAuditLog group_info mapstructure.Decode err:%v", err)
		return
	}

	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = g.eacpLog.EacpLog(&v, interfaces.OpDeleteGroup, groupInfo)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupAuditLog err:%v", err)
	}
	return err
}

// ModifyGroup 用户组修改
func (g *group) ModifyGroup(visitor *interfaces.Visitor, groupID, name string, nameChanged bool, notes string, notesChanged bool) (err error) {
	// 判断用户是否有用户组的管理权限（实名用户和应用账户）
	err = g.checkGroupManageAuth(visitor)
	if err != nil {
		return
	}

	// 判断用户组是否存在
	var info interfaces.GroupInfo
	info, err = g.groupDB.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if info.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound))
		return err
	}

	// 检测用户组名是否重复
	var id string
	id, err = g.groupDB.GetGroupIDByName(name)
	if err != nil {
		return err
	}
	if id != "" && id != groupID {
		err = rest.NewHTTPErrorV2(errors.Conflict, "the name of this group is existing",
			rest.SetCodeStr(errors.StrConflictGroup))
		return err
	}

	// 修改用户组信息
	err = g.groupDB.ModifyGroup(groupID, name, nameChanged, notes, notesChanged)
	if err != nil {
		return err
	}

	go func() {
		if !nameChanged {
			return
		}
		// 组织架构显示名变更消息信息
		obj := interfaces.NameChangeMsg{
			ID:      groupID,
			NewName: name,
			OType:   "group",
		}
		_ = g.messageBroker.Publish(interfaces.OrgNameChange, obj)
	}()

	// 记录日志
	go func() {
		nowName := info.Name
		if nameChanged {
			nowName = name
		}
		nowNotes := info.Notes
		if notesChanged {
			nowNotes = notes
		}

		groupInfo := interfaces.GroupInfo{
			Name:  nowName,
			Notes: nowNotes,
		}

		// 获取事务处理器
		var err1 error
		tx, err1 := g.pool.Begin()
		if err1 != nil {
			g.logger.Errorf("Modify group send audit log  pool begin error:%v", err1)
			return
		}

		// 异常时Rollback
		defer func() {
			switch err1 {
			case nil:
				// 提交事务
				if err1 = tx.Commit(); err1 != nil {
					g.logger.Errorf("ModifyGroup Transaction Commit Error:%v", err1)
					return
				}

				g.ob.NotifyPushOutboxThread()

			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					g.logger.Errorf("ModifyGroup Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 记录审计日志
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["group_info"] = groupInfo
		err1 = g.ob.AddOutboxInfo(outboxGroupModifiedLog, content, tx)
		if err1 != nil {
			g.logger.Errorf("ModifyGroup AddOutboxInfo err:%v", err1)
		}
	}()

	return err
}

// sendModifyGroupAuditLog 修改用户组发送审计消息
func (g *group) sendModifyGroupAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	groupInfo := interfaces.GroupInfo{}
	err = mapstructure.Decode(info["group_info"], &groupInfo)
	if err != nil {
		g.logger.Errorf("sendModifyGroupAuditLog groupInfo mapstructure.Decode err:%v", err)
		return
	}

	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		g.logger.Errorf("sendModifyGroupAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = g.eacpLog.EacpLog(&v, interfaces.OpModifyGroup, groupInfo)
	if err != nil {
		g.logger.Errorf("sendModifyGroupAuditLog err:%v", err)
	}
	return err
}

// GetGroupByID 获取指定的用户组
func (g *group) GetGroupByID(visitor *interfaces.Visitor, groupID string) (info interfaces.GroupInfo, err error) {
	// 判断用户是否有用户组的列举权限（实名用户和应用账户）
	err = g.checkGroupGetAuth(visitor)
	if err != nil {
		return
	}

	// 获取符合条件的用户组
	info, err = g.groupDB.GetGroupByID(groupID)
	if err != nil {
		return info, err
	}
	if info.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound))
		return info, err
	}

	return info, err
}

// GetGroup 用户组列举
func (g *group) GetGroup(visitor *interfaces.Visitor, info interfaces.SearchInfo) (num int, outInfo []interfaces.GroupInfo, err error) {
	// 判断用户是否有用户组的列举权限（实名用户和应用账户）
	err = g.checkGroupGetAuth(visitor)
	if err != nil {
		return
	}

	// 获取所有符合条件的用户组数量
	num, err = g.groupDB.GetGroupsNum(info)
	if err != nil {
		return num, outInfo, err
	}

	var infos []interfaces.GroupInfo
	if num > 0 {
		// 获取符合条件并且分页的用户组
		infos, err = g.groupDB.GetGroups(info)
		if err != nil {
			return num, outInfo, err
		}
	}

	outInfo = append(outInfo, infos...)
	return num, outInfo, err
}

// AddOrDeleteGroupMemebers 批量删除或者添加用户组成员
func (g *group) AddOrDeleteGroupMemebers(visitor *interfaces.Visitor, method, groupID string, infos map[string]interfaces.GroupMemberInfo) (err error) {
	if method == "POST" {
		err = g.AddGroupMembers(visitor, groupID, infos)
	} else if method == "DELETE" {
		err = g.DeleteGroupMembers(visitor, groupID, infos)
	}

	return err
}

// DeleteGroupMembers 用户组成员删除
func (g *group) DeleteGroupMembers(visitor *interfaces.Visitor, groupID string, infos map[string]interfaces.GroupMemberInfo) (err error) {
	// 判断用户是否有用户组的管理权限（实名用户和应用账户）
	err = g.checkGroupManageAuth(visitor)
	if err != nil {
		return
	}

	// 判断用户组是否存在
	var info interfaces.GroupInfo
	info, err = g.groupDB.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if info.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound))
		return err
	}

	// 判断用户和部门是否存在, 如果不存在 则删除存在的用户
	_, memberInfos, _ := g.checkUserAndDepartmentExist(visitor, infos)

	// 删除用户
	for k := range infos {
		temp := infos[k]
		err = g.groupMemberDB.DeleteGroupMember(groupID, &temp)
		if err != nil {
			return err
		}
	}

	// 记录日志
	go func() {
		j := 0
		strNames := make([]string, len(memberInfos))
		for _, v := range memberInfos {
			strNames[j] = v.Name
			j++
		}

		groupMembers := interfaces.GroupMemberNames{
			GroupName:   info.Name,
			MemberNames: strNames,
		}

		// 获取事务处理器
		var err1 error
		tx, err1 := g.pool.Begin()
		if err1 != nil {
			g.logger.Errorf("DeleteGroupMembers send audit log  pool begin error:%v", err1)
			return
		}

		// 异常时Rollback
		defer func() {
			switch err {
			case nil:
				// 提交事务
				if err1 = tx.Commit(); err1 != nil {
					g.logger.Errorf("DeleteGroupMembers Transaction Commit Error:%v", err1)
					return
				}

				g.ob.NotifyPushOutboxThread()

			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					g.logger.Errorf("DeleteGroupMembers Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 记录审计日志
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["group_members"] = groupMembers
		err1 = g.ob.AddOutboxInfo(outboxGroupMembersDeletedLog, content, tx)
		if err1 != nil {
			g.logger.Errorf("DeleteGroupMembers AddOutboxInfo err:%v", err1)
		}
	}()

	return err
}

// sendDeleteGroupMembersAuditLog 删除用户组成员发送审计消息
func (g *group) sendDeleteGroupMembersAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	groupMembers := interfaces.GroupMemberNames{}
	err = mapstructure.Decode(info["group_members"], &groupMembers)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupMembersAuditLog group_members mapstructure.Decode err:%v", err)
		return
	}

	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupMembersAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = g.eacpLog.EacpLog(&v, interfaces.OpDeleteGroupMembers, groupMembers)
	if err != nil {
		g.logger.Errorf("sendDeleteGroupMembersAuditLog err:%v", err)
	}
	return err
}

// AddGroupMembers 用户组成员添加
func (g *group) AddGroupMembers(visitor *interfaces.Visitor, groupID string, infos map[string]interfaces.GroupMemberInfo) (err error) {
	// 判断用户是否有用户组的管理权限（实名用户和应用账户）
	err = g.checkGroupManageAuth(visitor)
	if err != nil {
		return
	}

	// 判断用户组是否存在
	var info interfaces.GroupInfo
	info, err = g.groupDB.GetGroupByID(groupID)
	if err != nil {
		return err
	}
	if info.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "group does not exist", rest.SetCodeStr(errors.StrNotFoundGroupNotFound))
		return err
	}

	// 判断用户和部门是否存在
	var memberInfos []interfaces.NameInfo
	var ret bool
	ret, memberInfos, err = g.checkUserAndDepartmentExist(visitor, infos)
	if err != nil && !ret {
		return err
	}

	// 添加用户组成员
	for k := range infos {
		temp := infos[k]

		// 判断用户组成员是否存在
		var ret bool
		ret, err = g.groupMemberDB.CheckGroupMembersExist(groupID, &temp)
		if err != nil {
			return err
		}
		if ret {
			continue
		}

		// 添加用户组
		err = g.groupMemberDB.AddGroupMember(groupID, &temp)
		if err != nil {
			return err
		}
	}

	// 记录日志
	go func() {
		var err1 error
		j := 0
		strNames := make([]string, len(memberInfos))
		for _, v := range memberInfos {
			strNames[j] = v.Name
			j++
		}

		groupMembers := interfaces.GroupMemberNames{
			GroupName:   info.Name,
			MemberNames: strNames,
		}

		// 获取事务处理器
		tx, err1 := g.pool.Begin()
		if err != nil {
			g.logger.Errorf("AddGroupMembers send audit log  pool begin error:%v", err1)
			return
		}

		// 异常时Rollback
		defer func() {
			switch err1 {
			case nil:
				// 提交事务
				if err1 = tx.Commit(); err1 != nil {
					g.logger.Errorf("AddGroupMembers Transaction Commit Error:%v", err1)
					return
				}

				g.ob.NotifyPushOutboxThread()

			default:
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					g.logger.Errorf("AddGroupMembers Rollback err:%v", rollbackErr)
				}
			}
		}()

		// 记录审计日志
		content := make(map[string]interface{})
		content["visitor"] = *visitor
		content["group_members"] = groupMembers
		err1 = g.ob.AddOutboxInfo(outboxGroupMembersAddedLog, content, tx)
		if err != nil {
			g.logger.Errorf("AddGroupMembers AddOutboxInfo err:%v", err1)
		}
	}()

	return err
}

// sendAddGroupMembersAuditLog 添加用户组成员发送审计消息
func (g *group) sendAddGroupMembersAuditLog(content interface{}) (err error) {
	info := content.(map[string]interface{})
	groupMembers := interfaces.GroupMemberNames{}
	err = mapstructure.Decode(info["group_members"], &groupMembers)
	if err != nil {
		g.logger.Errorf("sendAddGroupMembersAuditLog group_members mapstructure.Decode err:%v", err)
		return
	}

	v := interfaces.Visitor{}
	err = mapstructure.Decode(info["visitor"], &v)
	if err != nil {
		g.logger.Errorf("sendAddGroupMembersAuditLog mapstructure.Decode err:%v", err)
		return
	}

	err = g.eacpLog.EacpLog(&v, interfaces.OpAddGroupMembers, groupMembers)
	if err != nil {
		g.logger.Errorf("sendAddGroupMembersAuditLog err:%v", err)
	}
	return err
}

// GetGroupMemberIds 批量获取用户组成员id
func (g *group) GetGroupMembersID(visitor *interfaces.Visitor, groupIDs []string, bShowAllUser bool) (userIDs, departmentIDs []string, err error) {
	userIDs = make([]string, 0)
	departmentIDs = make([]string, 0)
	// 判断用户组是否存在
	RemoveDuplicatStrs(&groupIDs)

	splitedIDs := SplitArray(groupIDs)
	for _, ids := range splitedIDs {
		existIDs, err := g.groupDB.GetExistGroupIDs(ids)
		if err != nil {
			return nil, nil, err
		}

		// 有用户组不存在
		if len(ids) != len(existIDs) {
			// 获取不存在的部门id
			notExistIDs := Difference(ids, existIDs)
			err = rest.NewHTTPErrorV2(errors.GroupNotFound,
				g.i18n.Load(i18nIDObjectsInGroupNotFound, visitor.Language),
				rest.SetDetail(map[string]interface{}{"ids": notExistIDs}),
				rest.SetCodeStr(errors.StrBadRequestGroupNotFound),
			)
			return nil, nil, err
		}

		// 批量获取用户组成员信息
		var outInfos []interfaces.GroupMemberInfo
		outInfos, err = g.groupMemberDB.GetGroupMembersByGroupIDs(ids)
		if err != nil {
			return nil, nil, err
		}
		// 根据成员类型封装成员id
		for _, v := range outInfos {
			if v.MemberType-1 == 0 {
				userIDs = append(userIDs, v.ID)
			} else if v.MemberType-2 == 0 {
				departmentIDs = append(departmentIDs, v.ID)
			}
		}
	}

	// 如果不显示被禁用用户，则筛选
	if !bShowAllUser {
		// 获取用户信息，判断用户是否被禁用
		userInfos, err := g.userDB.GetUserDBInfo(userIDs)
		if err != nil {
			return nil, nil, err
		}

		info := make([]string, 0)
		for i := 0; i < len(userInfos); i++ {
			if userInfos[i].DisableStatus == interfaces.Enabled && userInfos[i].AutoDisableStatus == interfaces.AEnabled {
				info = append(info, userInfos[i].ID)
			}
		}

		userIDs = info
	}

	return userIDs, departmentIDs, nil
}

// GetGroupMembers 用户组成员列举
func (g *group) GetGroupMembers(ctx context.Context, visitor *interfaces.Visitor, groupID string, info interfaces.SearchInfo) (num int, outInfo []interfaces.GroupMemberInfo, err error) {
	// trace
	g.trace.SetInternalSpanName("业务逻辑-用户组成员列举")
	newCtx, span := g.trace.AddInternalTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	// 判断用户是否有用户组的列举权限（实名用户和应用账户）
	err = g.checkGroupGetAuth2(newCtx, visitor)
	if err != nil {
		return
	}

	// 判断用户组是否存在
	var groupInfo interfaces.GroupInfo
	groupInfo, err = g.groupDB.GetGroupByID2(newCtx, groupID)
	if err != nil {
		return num, outInfo, err
	}
	if groupInfo.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "this group is not existing")
		return num, outInfo, err
	}

	// 获取符合条件的用户组成员数量
	num, err = g.groupMemberDB.GetGroupMembersNum2(newCtx, groupID, info)
	if err != nil {
		return num, outInfo, err
	}

	// 获取符合条件的分页用户组成员信息
	var outInfos []interfaces.GroupMemberInfo
	if num > 0 {
		outInfos, err = g.groupMemberDB.GetGroupMembers(newCtx, groupID, info)
		if err != nil {
			return num, outInfo, err
		}
	}

	// 获取父部门路径和直属部门id
	outInfo = make([]interfaces.GroupMemberInfo, 0)
	departIDs := make([]string, 0)
	userIDs := make([]string, 0)

	// 区分成员类型
	for _, v := range outInfos {
		if v.MemberType-1 == 0 {
			userIDs = append(userIDs, v.ID)
		} else if v.MemberType-2 == 0 {
			departIDs = append(departIDs, v.ID)
		}
	}

	// 获取路径信息
	pathNameInfos, _, err := g.handlePathParams(newCtx, visitor, userIDs, departIDs)
	if err != nil {
		return num, outInfo, err
	}

	// 组织用户的直属部门信息
	for _, v := range outInfos {
		v.ParentDeps = pathNameInfos[v.ID]
		v.DepartmentNames = make([]string, 0, len(v.ParentDeps))
		for _, v1 := range v.ParentDeps {
			if len(v1) > 0 {
				v.DepartmentNames = append(v.DepartmentNames, v1[len(v1)-1].Name)
			}
		}
		outInfo = append(outInfo, v)
	}
	return num, outInfo, err
}

func (g *group) handlePathParams(ctx context.Context, visitor *interfaces.Visitor, userIDs, departIDs []string) (data map[string][][]interfaces.NameInfo, departNames map[string]string, err error) {
	// 获取成员的路径信息(用户路径为直属部门路径，部门路径需要去除最后一个自身id)
	parentDepPaths, err := g.userDB.GetUsersPath2(ctx, userIDs)
	if err != nil {
		return nil, nil, err
	}

	departInfos, err := g.departmentDB.GetDepartmentInfoByIDs(ctx, departIDs)
	if err != nil {
		return nil, nil, err
	}

	paths := make(map[string][][]string)
	departNames = make(map[string]string)
	for k := range departInfos {
		tempPath := strings.Split(departInfos[k].Path, "/")
		paths[departInfos[k].ID] = [][]string{tempPath[:len(tempPath)-1]}
		departNames[departInfos[k].ID] = departInfos[k].Name
	}

	for k := range parentDepPaths {
		temp := make([][]string, 0, len(parentDepPaths[k]))
		for _, v1 := range parentDepPaths[k] {
			temp = append(temp, strings.Split(v1, "/"))
		}
		paths[k] = temp
	}

	// 获取所有成员路径上部门的id和name
	allDepartIDs := make([]string, 0)
	for _, v := range paths {
		for _, v1 := range v {
			allDepartIDs = append(allDepartIDs, v1...)
		}
	}
	RemoveDuplicatStrs(&allDepartIDs)

	allDepartInfos, err := g.departmentDB.GetDepartmentInfoByIDs(ctx, allDepartIDs)
	if err != nil {
		return nil, nil, err
	}

	depNameMap := make(map[string]string)
	for k := range allDepartInfos {
		depNameMap[allDepartInfos[k].ID] = allDepartInfos[k].Name
	}
	// 处理未分配组的信息
	depNameMap["-1"] = g.i18n.Load(i18nIDObjectsInUnDistributeUserGroup, visitor.Language)

	// 组织用户的直属部门信息
	data = make(map[string][][]interfaces.NameInfo)
	for k, v := range paths {
		// 如果没有直属部门，用户为未分配组，部门为顶级部门
		tempPaths := make([][]interfaces.NameInfo, 0, len(v))
		for _, v1 := range v {
			tempPath := make([]interfaces.NameInfo, len(v1))
			for i, v2 := range v1 {
				tempPath[i] = interfaces.NameInfo{
					Name: depNameMap[v2],
					ID:   v2,
				}
			}
			tempPaths = append(tempPaths, tempPath)
		}

		data[k] = tempPaths
	}
	return
}

// SearchGroupByKeyword 用户组关键字搜索
func (g *group) SearchGroupByKeyword(keyword string, start, limit int) (out []interfaces.NameInfo, err error) {
	out = make([]interfaces.NameInfo, 0)
	tmpData, err := g.groupDB.SearchGroupByKeyword(keyword, start, limit)
	if err != nil {
		return nil, err
	}

	out = append(out, tmpData...)
	return
}

// SearchGroupNumByKeyword 用户组关键字搜索符合条件的用户组总数目
func (g *group) SearchGroupNumByKeyword(keyword string) (num int, err error) {
	return g.groupDB.SearchGroupNumByKeyword(keyword)
}

// SearchMembersByKeyword 用户组成员关键字搜索
func (g *group) SearchMembersByKeyword(keyword string, start, limit int) (out []interfaces.MemberInfo, err error) {
	out = make([]interfaces.MemberInfo, 0)
	tmpData, err := g.groupMemberDB.SearchMembersByKeyword(keyword, start, limit)
	if err != nil {
		return nil, err
	}

	out = append(out, tmpData...)

	return
}

// SearchMemberNumByKeyword 用户组成员关键字搜索符合条件的用户组总数目
func (g *group) SearchMemberNumByKeyword(keyword string) (num int, err error) {
	return g.groupMemberDB.SearchMemberNumByKeyword(keyword)
}

// ConvertGroupName 根据用户组ID获取用户组名
func (g *group) ConvertGroupName(visitor *interfaces.Visitor, ids []string, bStrict bool) (nameInfo []interfaces.NameInfo, err error) {
	nameInfo = make([]interfaces.NameInfo, 0)
	if len(ids) == 0 {
		return nameInfo, nil
	}

	// 去重
	copyIDs := make([]string, len(ids))
	copy(copyIDs, ids)
	RemoveDuplicatStrs(&copyIDs)

	outInfo, exsitIDs, err := g.groupDB.GetGroupName(copyIDs)
	if err != nil {
		return nameInfo, err
	}

	// 如果严格模式， 且有用户组不存在，则返回错误
	if bStrict && len(exsitIDs) != len(copyIDs) {
		// 获取不存在的部门id
		notExistIDs := Difference(copyIDs, exsitIDs)
		err = rest.NewHTTPErrorV2(errors.GroupNotFound,
			g.i18n.Load(i18nIDObjectsInGroupNotFound, visitor.Language),
			rest.SetDetail(map[string]interface{}{"ids": notExistIDs}),
			rest.SetCodeStr(errors.StrBadRequestGroupNotFound),
		)
		return nil, err
	}

	nameInfo = append(nameInfo, outInfo...)

	return
}

// GetGroupOnClient 客户端列举组
func (g *group) GetGroupOnClient(offset, limit int) (info []interfaces.NameInfo, num int, err error) {
	// 设置搜索条件
	var searchInfo interfaces.SearchInfo
	searchInfo.HasKeyWord = false
	searchInfo.Offset = offset
	searchInfo.Limit = limit
	searchInfo.Sort = interfaces.Name
	searchInfo.Direction = interfaces.Asc

	// 获取所有符合条件的用户组数量
	num, err = g.groupDB.GetGroupsNum(searchInfo)
	if err != nil {
		return nil, num, err
	}

	// 获取信息
	tmpInfo, err := g.groupDB.GetGroups(searchInfo)
	if err != nil {
		return nil, num, err
	}

	info = make([]interfaces.NameInfo, 0)
	var tmp interfaces.NameInfo
	for _, v := range tmpInfo {
		tmp.ID = v.ID
		tmp.Name = v.Name

		info = append(info, tmp)
	}
	return
}

// GetMemberOnClient 客户端列举组成员
func (g *group) GetMemberOnClient(id string, offset, limit int) (info []interfaces.MemberSimpleInfo, num int, err error) {
	// 判断用户组是否存在
	var groupInfo interfaces.GroupInfo
	groupInfo, err = g.groupDB.GetGroupByID(id)
	if err != nil {
		return nil, num, err
	}
	if groupInfo.ID == "" {
		err = rest.NewHTTPErrorV2(errors.NotFound, "group does not exist")
		return nil, num, err
	}

	// 设置搜索条件
	var searchInfo interfaces.SearchInfo
	searchInfo.HasKeyWord = false
	searchInfo.Offset = offset
	searchInfo.Limit = limit
	searchInfo.Sort = interfaces.Name
	searchInfo.Direction = interfaces.Asc
	searchInfo.NotShowDisabledUser = true

	// 获取符合条件的用户组成员数量
	num, err = g.groupMemberDB.GetGroupMembersNum(id, searchInfo)
	if err != nil {
		return nil, num, err
	}

	// 获取组成员信息
	tmp, err := g.groupMemberDB.GetMemberOnClient(id, offset, limit)
	if err != nil {
		return nil, num, err
	}
	info = append(info, tmp...)
	return
}

// DeleteGroupMemberByMemberID 删除组成员
func (g *group) DeleteGroupMemberByMemberID(id string) (err error) {
	return g.groupMemberDB.DeleteGroupMemberByMemberID(id)
}

// 检测是否所有id都存在
func (g *group) checkNotExistIDs(allIDs []string, existIDs map[string]bool) (ret bool, outInfo map[string]interface{}) {
	var outIDs []string
	outInfo = make(map[string]interface{})
	if len(allIDs) != len(existIDs) {
		for _, v := range allIDs {
			if !existIDs[v] {
				outIDs = append(outIDs, v)
			}
		}
		outInfo["ids"] = outIDs
		ret = false
	} else {
		ret = true
	}
	return ret, outInfo
}

// 检测部门或者用户是否存在
func (g *group) checkUserAndDepartmentExist(visitor *interfaces.Visitor, memberInfos map[string]interfaces.GroupMemberInfo) (ret bool, infos []interfaces.NameInfo, err error) {
	// 数据处理
	userIDs := make([]string, 0)
	departmentIDs := make([]string, 0)
	for _, v := range memberInfos {
		if v.MemberType-1 == 0 {
			userIDs = append(userIDs, v.ID)
		} else if v.MemberType-2 == 0 {
			departmentIDs = append(departmentIDs, v.ID)
		}
	}

	// 判断用户是否存在
	if len(userIDs) > 0 {
		outUserIDs := make(map[string]bool)
		var userInfos []interfaces.UserDBInfo
		userInfos, _, err = g.userDB.GetUserName(userIDs)
		if err != nil {
			return false, infos, err
		}
		outUserInfo := make([]interfaces.NameInfo, 0)
		for k := range userInfos {
			tmp := interfaces.NameInfo{
				Name: userInfos[k].Name,
				ID:   userInfos[k].ID,
			}
			outUserIDs[userInfos[k].ID] = true
			outUserInfo = append(outUserInfo, tmp)
		}

		retUser, mapUser := g.checkNotExistIDs(userIDs, outUserIDs)
		if !retUser {
			return false, infos, rest.NewHTTPErrorV2(errors.UserNotFound,
				g.i18n.Load(i18nIDObjectsInUserNotFound, visitor.Language),
				rest.SetDetail(mapUser),
				rest.SetCodeStr(errors.StrBadRequestUserNotFound),
			)
		}
		infos = append(infos, outUserInfo...)
	}

	// 判断部门是否存在
	if len(departmentIDs) > 0 {
		outDepartmentIDs := make(map[string]bool)
		var departmentInfos []interfaces.NameInfo
		departmentInfos, _, err = g.departmentDB.GetDepartmentName(departmentIDs)
		if err != nil {
			return false, infos, err
		}
		outDepartmentInfo := make([]interfaces.NameInfo, 0)
		for _, v := range departmentInfos {
			outDepartmentIDs[v.ID] = true
			outDepartmentInfo = append(outDepartmentInfo, v)
		}

		retDepart, mapDepartment := g.checkNotExistIDs(departmentIDs, outDepartmentIDs)
		if !retDepart {
			err = rest.NewHTTPErrorV2(errors.DepartmentNotFound,
				g.i18n.Load(i18nIDObjectsInDepartNotFound, visitor.Language),
				rest.SetDetail(mapDepartment),
				rest.SetCodeStr(errors.StrBadRequestDepartmentNotFound))
			return false, infos, err
		}

		infos = append(infos, outDepartmentInfo...)
	}

	return true, infos, nil
}

// checkGroupManageAuth 支持实名用户和应用账户权限检查
//
//nolint:exhaustive
func (g *group) checkGroupManageAuth(visitor *interfaces.Visitor) (err error) {
	switch visitor.Type {
	case interfaces.RealName:
		err = checkManageAuthority(g.role, visitor.ID)
	case interfaces.App:
		err = checkAppPerm(g.orgPermApp, visitor.ID, interfaces.Group, interfaces.Modify)
	default:
		err = rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
	}
	return err
}

// checkGroupManageAuth2 支持实名用户和应用账户权限检查，支持trace
//
//nolint:exhaustive
func (g *group) checkGroupManageAuth2(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	switch visitor.Type {
	case interfaces.RealName:
		err = checkManageAuthority2(ctx, g.role, visitor.ID)
	case interfaces.App:
		err = checkAppPerm2(ctx, g.orgPermApp, visitor.ID, interfaces.Group, interfaces.Modify)
	default:
		err = rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority")
	}
	return err
}

// checkGroupGetAuth 支持实名用户和应用账户权限检查
//
//nolint:exhaustive
func (g *group) checkGroupGetAuth(visitor *interfaces.Visitor) (err error) {
	switch visitor.Type {
	case interfaces.RealName:
		err = checkGetInfoAuthority(g.role, visitor.ID)
	case interfaces.App:
		err = checkAppPerm(g.orgPermApp, visitor.ID, interfaces.Group, interfaces.Read)
	default:
		err = rest.NewHTTPError("this user do not has the authority", errors.Forbidden, nil)
	}
	return err
}

// checkGroupGetAuth 支持实名用户和应用账户权限检查，支持trace
//
//nolint:exhaustive
func (g *group) checkGroupGetAuth2(ctx context.Context, visitor *interfaces.Visitor) (err error) {
	switch visitor.Type {
	case interfaces.RealName:
		err = checkGetInfoAuthority2(ctx, g.role, visitor.ID)
	case interfaces.App:
		err = checkAppPerm2(ctx, g.orgPermApp, visitor.ID, interfaces.Group, interfaces.Read)
	default:
		err = rest.NewHTTPErrorV2(errors.Forbidden, "this user do not has the authority")
	}
	return err
}

// 检查有效性
func (g *group) checkName(name string) (err error) {
	illegalChars := " |\\/:*?\"<>"
	length := utf8.RuneCountInString(name)

	if strings.ContainsAny(name, illegalChars) || length < 1 || length > 128 {
		return rest.NewHTTPErrorV2(rest.BadRequest, "param name is illegal")
	}

	return nil
}
