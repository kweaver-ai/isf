package logics

import (
	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/interfaces"
)

const (
	// SuperAdminID 超级管理员
	SuperAdminID = "266c6a42-6131-4d62-8f39-853e7093701c"
	// AdminAdminID 系统管理员
	AdminAdminID = "266c6a42-6131-4d62-8f39-853e7093701c"
	// SecurityAdmin 安全管理员
	SecurityAdminID = "4bb41612-a040-11e6-887d-005056920bea"
	// AuditAdmin 审计管理员
	AuditAdminID = "94752844-BDD0-4B9E-8927-1CA8D427E699"
	// SystemAdmin 此账户已废弃
	SystemAdminID = "234562BE-88FF-4440-9BFF-447F139871A2"
)

const (
	// 使用iota，增加新的枚举时，必须按照顺序添加，不得在中间插入
	// 否则会改变枚举的值，造成outbox handler与db中记录的type对应不上
	_ = iota
	// OutboxDeleteHydraSession 清除用户login、consent会话
	OutboxDeleteHydraSession

	// OutboxSetAppAccessTokenPermLog 发送设置应用账户权限审计日志
	OutboxSetAppAccessTokenPermLog

	// OutboxDeleteAppAccessTokenPermLog 发送删除应用账户权限审计日志
	OutboxDeleteAppAccessTokenPermLog

	// OutboxSendAuditLog 发送审计日志
	OutboxSendAuditLog

	// OutboxAnonymousSmsExpUpdated 匿名登录短信过期时间更新
	OutboxAnonymousSmsExpUpdated
)

// outbox业务类型
const (
	// 默认业务类型
	_ = iota
	// OutboxBusinessHydraSession 业务类型
	OutboxBusinessHydraSession

	// OutboxBusinessAccessTokenPerm 权限相关
	OutboxBusinessAccessTokenPerm

	// OutboxAuditLog 审计日志
	OutboxAuditLog

	// OutboxConfig 配置
	OutboxConfig

	// 若新增新的业务类型，需在initdb中对anyshare.t_outbox_lock表进行数据初始化
)

// CheckVisitorType 检测访问者类型
func CheckVisitorType(visitor *interfaces.Visitor, roleTypes []interfaces.RoleType, acceptVisitorTypes []interfaces.VisitorType, acceptRoleTypes []interfaces.RoleType) (err error) {
	isVisitorTypePass := false
	isRoleTypePass := false
	// 访问者类型是否在检测范围内
	for _, val := range acceptVisitorTypes {
		if visitor.Type == val {
			isVisitorTypePass = true
		}
	}
	if !isVisitorTypePass {
		err = rest.NewHTTPError("Unsupported user type", rest.Unauthorized, nil)
		return
	}

	if visitor.Type == interfaces.RealName {
		// 访问者角色是否在检测范围内
		for _, roleType := range roleTypes {
			for _, val := range acceptRoleTypes {
				if roleType == val {
					isRoleTypePass = true
				}
			}
		}
		if !isRoleTypePass {
			err = rest.NewHTTPError("Unsupported user role type", rest.Unauthorized, nil)
			return
		}
	}
	return nil
}
