package logics

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"AuditLog/common"
	"AuditLog/common/conf"
	"AuditLog/common/utils/rclogutils"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
	"AuditLog/models/rcvo"
	"AuditLog/tapi/sharemgnt"
)

var (
	a     *ActiveLog
	aOnce sync.Once
)

type ActiveLog struct {
	logger           api.Logger
	tracer           api.Tracer
	loginLogRepo     interfaces.LogRepo              // 数据库对象
	mgntLogRepo      interfaces.LogRepo              // 数据库对象
	operLogRepo      interfaces.LogRepo              // 数据库对象
	logScopeStrategy interfaces.LogScopeStrategyRepo // 数据库对象
	userMgntRepo     interfaces.UserMgntRepo
	shareMgntRepo    interfaces.ShareMgntRepo
	docCenterRepo    interfaces.DocCenterRepo
}

func NewActiveLog() interfaces.ActiveLog {
	aOnce.Do(func() {
		a = &ActiveLog{
			logger:           logger,
			tracer:           tracer,
			loginLogRepo:     loginLogRepo,
			operLogRepo:      operLogRepo,
			mgntLogRepo:      mgntLogRepo,
			logScopeStrategy: logScopeStrategyRepo,
			userMgntRepo:     userMgntRepo,
			shareMgntRepo:    shareMgntRepo,
			docCenterRepo:    docCenterRepo,
		}
	})

	return a
}

// 获取角色成员ID
func (al *ActiveLog) getUserIDsByRoleName(roleName string) (userIDs []string, err error) {
	roleMemberInfos, _, err := al.userMgntRepo.GetUserIDsByRoleNames([]string{roleName})
	if err != nil {
		return nil, err
	}
	for _, roleMemberInfo := range roleMemberInfos {
		if roleMemberInfo.Role == roleName {
			for _, member := range roleMemberInfo.Members {
				userIDs = append(userIDs, member.ID)
			}
		}
	}
	return
}

// 获取可查看范围内的用户ID
func (al *ActiveLog) getUserIds(ctx context.Context, logType string, userID string) (includeIDs []string, excludeIDs []string, err error) {
	userInfos, _, err := al.userMgntRepo.GetUserInfoByID([]string{userID})
	if err != nil {
		al.logger.Errorf("[getUserIds]: get user info error: %v", err)
		return
	}
	if len(userInfos) > 0 {
		userRoles := userInfos[0].Roles

		// 超级管理员访问
		if common.InArray(common.SuperAdmin, userRoles) {
			return nil, nil, nil
		}

		// 三权分立角色 系统管理员、安全管理员、审计管理员访问
		for _, r := range common.MutuallyRoles {
			if common.InArray(r, userRoles) {
				scope, scopeErr := al.logScopeStrategy.GetActiveScopeBy(common.LogTypeMap[logType], r)
				if scopeErr != nil {
					al.logger.Errorf("[getUserIds]: get scope by role error: %v", scopeErr)
					return nil, nil, scopeErr
				}
				if len(scope) > 0 {
					// 如果包含普通用户，则使用排除法
					if common.InArray(common.NormalUser, scope) {
						for _, mr := range common.MutuallyRoles {
							if !common.InArray(mr, scope) {
								userIDs, err := al.getUserIDsByRoleName(mr)
								if err != nil {
									al.logger.Errorf("[getUserIds]: get user id by role error: %v", err)
									return nil, nil, err
								}
								excludeIDs = append(excludeIDs, userIDs...)
							}
						}
					} else {
						for _, s := range scope {
							userIDs, err := al.getUserIDsByRoleName(s)
							if err != nil {
								al.logger.Errorf("[getUserIds]: get user id by role error: %v", err)
								return nil, nil, err
							}
							includeIDs = append(includeIDs, userIDs...)
						}
					}

					return
				}
			}
		}

		// 组织审计员角色只能查看自己组织下用户日志，但不能查看超级管理员、系统管理员、安全管理员、审计管理员的日志
		if common.InArray(common.OrgAudit, userRoles) {
			// 获取组织审计员角色成员信息以及管辖部门信息
			memberInfos, err := al.shareMgntRepo.GetRoleMemberInfos(sharemgnt.NCT_SYSTEM_ROLE_ORG_AUDIT)
			if err != nil {
				al.logger.Errorf("[getUserIds]: get role member infos error: %v", err)
				return nil, nil, err
			}
			// 获取搜索管辖用户ID
			for _, member := range memberInfos {
				if member.UserId == userID {
					for _, deptID := range member.ManageDeptInfo.DepartmentIds {
						subUserIDs, _, err := al.userMgntRepo.GetDeptAllUserIDs(deptID)
						if err != nil {
							al.logger.Errorf("[getUserIds]: get dept all user ids error: %v", err)
							return nil, nil, err
						}
						includeIDs = append(includeIDs, subUserIDs["all_user_ids"]...)
					}
					break
				}
			}

			// 排除超级管理员、系统管理员、安全管理员、审计管理员
			for _, mr := range append(common.MutuallyRoles, common.SuperAdmin) {
				userIDs, err := al.getUserIDsByRoleName(mr)
				if err != nil {
					al.logger.Errorf("[getUserIds]: get user id by role error: %v", err)
					return nil, nil, err
				}
				excludeIDs = append(excludeIDs, userIDs...)
			}

			return includeIDs, excludeIDs, nil
		}

		return nil, nil, errors.NewCtx(ctx, errors.ForbiddenErr, "No permission", nil)
	}

	return
}

// 获取活跃日志报表元数据
func (al *ActiveLog) GetActiveMetadata() (meta *rcvo.ReportMetadataRes, err error) {
	return rclogutils.GetActiveMetadata()
}

// 获取活跃日志报表数据列表
func (al *ActiveLog) GetActiveDataList(ctx context.Context, logType string, req *rcvo.ReportGetDataListReq, userID string) (res *rcvo.ActiveReportListRes, err error) {
	var tErr error
	_, span := al.tracer.AddInternalTrace(ctx)
	defer func() { al.tracer.TelemetrySpanEnd(span, tErr) }()

	var (
		inUserIDs []string
		exUserIDs []string
	)

	if len(req.IDs) == 0 {
		inUserIDs, exUserIDs, err = al.getUserIds(ctx, logType, userID)
		if err != nil {
			return nil, err
		}
	}

	sqlStr, err := rclogutils.BuildActiveCondition(req.Condition, req.OrderBy, req.IDs, inUserIDs, exUserIDs)
	if err != nil {
		return nil, err
	}

	// 获取日志信息
	var logs []*models.LogPO

	switch logType {
	case common.Login:
		logs, err = al.loginLogRepo.FindByCondition(req.Offset, req.Limit, sqlStr, req.IDs)
	case common.Management:
		logs, err = al.mgntLogRepo.FindByCondition(req.Offset, req.Limit, sqlStr, req.IDs)
	case common.Operation:
		logs, err = al.operLogRepo.FindByCondition(req.Offset, req.Limit, sqlStr, req.IDs)
	default:
		return nil, fmt.Errorf("[GetActiveDataList]: invalid logType: %s", logType)
	}

	entries := make(rcvo.ActiveLogReports, 0, len(logs))

	if len(logs) == 0 {
		res = &rcvo.ActiveReportListRes{
			Entries: entries,
		}
		return
	}

	for _, log := range logs {
		var opType string
		switch logType {
		case common.Management:
			opType = locale.GetRCLogMgntI18n(ctx, log.OpType)
		case common.Login:
			opType = locale.GetRCLogLoginI18n(ctx, log.OpType)
		case common.Operation:
			opType = locale.GetRCLogOpI18n(ctx, log.OpType)
		}

		entry := rcvo.ActiveLogReport{
			ID:          log.LogID,
			UserName:    log.UserName,
			CreatedTime: log.Date / 1000,
			IP:          log.IP,
			Mac:         log.MAC,
			Msg:         log.Msg,
			ExMsg:       log.ExMsg,
			OpType:      opType,
			UserPaths:   log.UserPaths,
			Level:       locale.GetRCLogLevelI18n(ctx, locale.LogLevelMap[log.Level]),
			ObjName:     log.ObjName,
			ObjType:     locale.GetRCLogObjTypeI18n(ctx, log.ObjType),
		}
		entries = append(entries, entry)
	}

	// 获取日志总数
	nTotalCount := len(req.IDs)
	if nTotalCount == 0 {
		switch logType {
		case common.Login:
			nTotalCount, err = al.loginLogRepo.FindCountByCondition(sqlStr)
		case common.Management:
			nTotalCount, err = al.mgntLogRepo.FindCountByCondition(sqlStr)
		case common.Operation:
			nTotalCount, err = al.operLogRepo.FindCountByCondition(sqlStr)
		}

		if err != nil {
			return nil, err
		}
	}

	res = &rcvo.ActiveReportListRes{
		Entries:    entries,
		TotalCount: nTotalCount,
	}
	return
}

// 获取活跃日志报表字段值
func (al *ActiveLog) GetActiveFieldValues(ctx context.Context, logType string, req *rcvo.ReportGetFieldValuesReq) (res *rcvo.ReportFieldValuesRes, err error) {
	var tErr error
	_, span := al.tracer.AddInternalTrace(ctx)
	defer func() { al.tracer.TelemetrySpanEnd(span, tErr) }()

	entries := make(rcvo.ReportFieldValues, 0)

	if common.InArray(req.Field, []string{"level", "op_type", "obj_type"}) {
		switch req.Field {
		case "level":
			for k, v := range locale.LogLevelMap {
				i18nValue := locale.GetRCLogLevelI18n(ctx, v)
				if req.KeyWord == "" || strings.Contains(i18nValue, req.KeyWord) {
					entries = append(entries, rcvo.ReportFieldValue{
						ValueCode: strconv.Itoa(k),
						ValueName: i18nValue,
					})
				}
			}
		case "op_type":
			switch logType {
			case common.Management:
				for k := range conf.MapManageOperTypeLang {
					i18nValue := locale.GetRCLogMgntI18n(ctx, k)
					if req.KeyWord == "" || strings.Contains(i18nValue, req.KeyWord) {
						entries = append(entries, rcvo.ReportFieldValue{
							ValueCode: strconv.Itoa(k),
							ValueName: i18nValue,
						})
					}
				}
			case common.Operation:
				for k := range conf.MapOperOperTypeLang {
					i18nValue := locale.GetRCLogOpI18n(ctx, k)
					if req.KeyWord == "" || strings.Contains(i18nValue, req.KeyWord) {
						entries = append(entries, rcvo.ReportFieldValue{
							ValueCode: strconv.Itoa(k),
							ValueName: i18nValue,
						})
					}
				}
			case common.Login:
				for k := range conf.MapLoginOperTypeLang {
					i18nValue := locale.GetRCLogLoginI18n(ctx, k)
					if req.KeyWord == "" || strings.Contains(i18nValue, req.KeyWord) {
						entries = append(entries, rcvo.ReportFieldValue{
							ValueCode: strconv.Itoa(k),
							ValueName: i18nValue,
						})
					}
				}
			}
		case "obj_type":
			for k := range conf.MapObjectTypeLang {
				i18nValue := locale.GetRCLogObjTypeI18n(ctx, k)
				if req.KeyWord == "" || strings.Contains(i18nValue, req.KeyWord) {
					entries = append(entries, rcvo.ReportFieldValue{
						ValueCode: strconv.Itoa(k),
						ValueName: i18nValue,
					})
				}
			}
		}
	}

	res = &rcvo.ReportFieldValuesRes{
		TotalCount: len(entries),
		Entries:    entries,
	}

	return
}
