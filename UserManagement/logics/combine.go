// Package logics combine Anyshare 业务逻辑层 -组合
package logics

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type combine struct {
	user       interfaces.LogicsUser
	department interfaces.LogicsDepartment
	group      interfaces.LogicsGroup
	contactor  interfaces.LogicsContactor
	role       interfaces.LogicsRole
	app        interfaces.LogicsApp
	trace      observable.Tracer
}

var (
	cOnce     sync.Once
	comLogics *combine
)

// NewCombine 创建新的combine
func NewCombine() *combine {
	cOnce.Do(func() {
		comLogics = &combine{
			user:       NewUser(),
			department: NewDepartment(),
			group:      NewGroup(),
			contactor:  NewContactor(),
			role:       NewRole(),
			app:        NewApp(),
			trace:      common.SvcARTrace,
		}
	})

	return comLogics
}

// ConvertIDToName  根据用户/部门/联系人组/用户组/应用账户id 获取用户/部门/联系人组/用户组/应用账户显示名
func (l *combine) ConvertIDToName(visitor *interfaces.Visitor, info *interfaces.OrgIDInfo, bV2, bStrict bool) (nameInfo interfaces.OrgNameInfo, err error) {
	nameInfo.UserNames, err = l.user.ConvertUserName(visitor, info.UserIDs, bStrict)
	if err != nil {
		return
	}
	nameInfo.DepartNames, err = l.department.ConvertDepartmentName(visitor, info.DepartIDs, bStrict)
	if err != nil {
		return
	}
	nameInfo.ContactorNames, err = l.contactor.ConvertContactorName(info.ContactorIDs, bStrict)
	if err != nil {
		return
	}
	nameInfo.GroupNames, err = l.group.ConvertGroupName(visitor, info.GroupIDs, bStrict)
	if err != nil {
		return
	}
	nameInfo.AppNames, err = l.app.ConvertAppName(info.AppIDs, bV2, bStrict)
	if err != nil {
		return
	}
	return
}

// GetEmails  根据用户/部门id 获取用户/部门emails
func (l *combine) GetEmails(visitor *interfaces.Visitor, info *interfaces.OrgIDInfo) (emailInfo interfaces.OrgEmailInfo, err error) {
	emailInfo.UserEmails, err = l.user.GetUserEmails(visitor, info.UserIDs)
	if err != nil {
		return
	}
	emailInfo.DepartEmails, err = l.department.GetDepartEmails(visitor, info.DepartIDs)
	if err != nil {
		return
	}
	return
}

// GetUserAndDepartmentInScope 获取在范围内的用户和部门
func (l *combine) GetUserAndDepartmentInScope(userIDs, deptIDs, rangeIDs []string) (outUserIDs, outDepIDs []string, err error) {
	// 获取范围部门的所有子部门
	emptyString := make([]string, 0)
	var scopeAllDepIDs []string
	scopeAllDepIDs, err = l.department.GetAllChildDeparmentIDs(rangeIDs)
	if err != nil {
		return emptyString, emptyString, err
	}
	scopeAllDepIDs = append(scopeAllDepIDs, rangeIDs...)
	RemoveDuplicatStrs(&scopeAllDepIDs)

	// 检查存在的部门ID
	outDepIDs = Intersection(deptIDs, scopeAllDepIDs)
	RemoveDuplicatStrs(&outDepIDs)

	// 检查存在的用户ID
	outUserIDs, err = l.user.GetUsersInDepartments(userIDs, scopeAllDepIDs)
	if err != nil {
		return emptyString, emptyString, err
	}
	return
}

// SearchGroupAndMemberInfoByKey 客户端搜索组和组成员信息
func (l *combine) SearchGroupAndMemberInfoByKey(info *interfaces.SearchClientInfo) (out interfaces.GMSearchOutInfo, err error) {
	out.GroupInfos = make([]interfaces.NameInfo, 0)
	out.MemberInfos = make([]interfaces.MemberInfo, 0)
	// 符合条件的组成员获取
	if info.BShowMember {
		out.MemberNum, err = l.group.SearchMemberNumByKeyword(info.Keyword)

		if err != nil {
			return out, err
		}

		if out.MemberNum > 0 {
			out.MemberInfos, err = l.group.SearchMembersByKeyword(info.Keyword, info.Offset, info.Limit)
			if err != nil {
				return out, err
			}
		}
	}

	if !info.BShowGroup {
		return out, err
	}

	// 添加策略逻辑

	// 获取符合条件的组
	out.GroupNum, err = l.group.SearchGroupNumByKeyword(info.Keyword)
	if err != nil {
		return out, err
	}

	if out.GroupNum > 0 {
		out.GroupInfos, err = l.group.SearchGroupByKeyword(info.Keyword, info.Offset, info.Limit)
	}

	return out, err
}

// SearchInOrgTree 在组织架构内搜索用户和信息
func (l *combine) SearchInOrgTree(ctx context.Context, visitor *interfaces.Visitor, info interfaces.OrgShowPageInfo) (
	users []interfaces.SearchUserInfo, departs []interfaces.SearchDepartInfo, userNum, departNum int, err error) {
	// trace
	l.trace.SetInternalSpanName("业务逻辑-在组织架构内搜索用户和信息")
	newCtx, span := l.trace.AddInternalTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	// 检测调用者是否拥有指定角色
	var mapRoles map[interfaces.Role]bool
	mapRoles, err = getRolesByUserID(l.role, visitor.ID)
	if err != nil {
		return
	}
	if !mapRoles[info.Role] {
		err = rest.NewHTTPError("this user do not has this role", rest.BadRequest, nil)
		return
	}

	// 数据获取
	users = make([]interfaces.SearchUserInfo, 0)
	if info.BShowUsers {
		users, userNum, err = l.user.SearchUsersByKeyInScope(newCtx, visitor, info)
		if err != nil {
			return
		}
	}

	departs = make([]interfaces.SearchDepartInfo, 0)
	if info.BShowDeparts {
		departs, departNum, err = l.department.SearchDepartsByKey(newCtx, visitor, info)
		if err != nil {
			return
		}
	}

	return
}
