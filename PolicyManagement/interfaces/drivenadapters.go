package interfaces

import "context"

//go:generate mockgen -package mock -source ../interfaces/drivenadapters.go -destination ../interfaces/mock/mock_drivenadapters.go

// License 许可证
type License struct {
	Product        string
	TotalUserQuota int
}

// DrivenLicense 许可证驱动接口
type DrivenLicense interface {
	GetLicenses(ctx context.Context) (infos map[string]License, err error)
}

type Role int

const (
	_ Role = iota
	SystemRoleSuperAdmin
	SystemRoleSysAdmin
	SystemRoleSecAdmin
	SystemRoleAuditAdmin
	SystemRoleOrgManager
	SystemRoleOrgAudit
	SystemRoleNormalUser
)

type UserInfo struct {
	ID    string
	Name  string
	Roles map[Role]bool
}

// DrivenUserManagement 用户管理驱动接口
type DrivenUserManagement interface {
	// GetUserInfos 获取用户信息
	GetUserInfos(ctx context.Context, ids []string) (infos map[string]UserInfo, err error)
}

// DrivenEacpLog 日志处理接口
type DrivenEacpLog interface {
	// OpAddAuthorizedProducts 新增产品授权
	OpAddAuthorizedProducts(visitor *Visitor, name string, products []string) (err error)

	// OpDeleteAuthorizedProducts 删除产品授权
	OpDeleteAuthorizedProducts(visitor *Visitor, name string, products []string) (err error)

	// OpUpdateAuthorizedProducts 更新产品授权
	OpUpdateAuthorizedProducts(visitor *Visitor, name string, currentProducts []string, futureProducts []string) (err error)
}
