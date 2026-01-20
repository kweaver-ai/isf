// Package interfaces 接口
package interfaces

import "context"

//go:generate mockgen -package mock -source ../interfaces/drivenadapters.go -destination ../interfaces/mock/mock_drivenadapters.go

// AppInfo 文档信息
type AppInfo struct {
	ID   string //  应用账户ID
	Name string //  应用账户名称
}

// SystemRoleType 用户角色类型，这里全部是系统角色
type SystemRoleType int32

// 用户角色类型定义
const (
	SuperAdmin        SystemRoleType = iota // 超级管理员
	SystemAdmin                             // 系统管理员
	AuditAdmin                              // 审计管理员
	SecurityAdmin                           // 安全管理员
	OrganizationAdmin                       // 组织管理员
	OrganizationAudit                       // 组织审计员
	NormalUser                              // 普通用户
)

// UserInfo 用户基本信息
type UserInfo struct {
	ID         string                  // 用户id
	Account    string                  // 用户名称
	VisionName string                  // 显示名
	CsfLevel   int                     // 密级
	Frozen     bool                    // 冻结状态
	Roles      map[SystemRoleType]bool // 角色
	Email      string                  // 邮箱地址
	Telephone  string                  // 电话号码
	ThirdAttr  string                  // 第三方应用属性
	ThirdID    string                  // 第三方应用id
	UserType   AccessorType            // 用户类型
	Groups     []Group                 // 用户及其所属部门所在的用户组
	ParentDeps [][]Department          // 父部门信息
}

// Department 组织结构部门
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Group 组基本信息(可包含用户和部门)
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// DrivenUserMgnt 服务处理接口
type DrivenUserMgnt interface {
	GetUserRolesByUserID(ctx context.Context, userID string) (roleTypes []SystemRoleType, err error)
	// GetNameByAccessorIDs 获取用户,部门,组织, 应用账户的名称
	GetNameByAccessorIDs(ctx context.Context, accessorIDs map[string]AccessorType) (accessorNames map[string]string, err error)
	// BatchGetUserInfoByID 批量获取用户的基础信息
	BatchGetUserInfoByID(ctx context.Context, userIDs []string) (userInfoMap map[string]UserInfo, err error)
	// GetAccessorIDsByUserID 获取指定用户的访问令牌
	GetAccessorIDsByUserID(ctx context.Context, userID string) (accessorIDs []string, err error)
	// GetParentDepartmentsByDepartmentID 根据部门ID获取父部门信息
	GetParentDepartmentsByDepartmentID(ctx context.Context, departmentID string) (parentDeps []Department, err error)
}
