package types

import "github.com/gin-gonic/gin"

// VisitUserInfo 访问用户信息
type VisitUserInfo struct {
	// 用户ID
	UserID string `json:"userId"`
	// 用户Token
	UserToken string `json:"userToken"`
	// IP地址
	IP string `json:"ip"`
	// 客户端类型
	ClientType string `json:"clientType"`
	// 设备唯一标识
	UDID string `json:"udid"`
	// 访问者类型
	VisitorType string `json:"visitorType"`
	// 用户角色
	UserRoles interface{} `json:"userRoles"`
	// CSF等级
	CSFLevel interface{} `json:"csfLevel"`
	// 用户名称
	Name string `json:"name"`
	// 账户类型
	AccountType string `json:"accountType"`
	// 部门路径
	DeptPaths interface{} `json:"deptPaths"`
}

func NewVisitUserInfo() *VisitUserInfo {
	return &VisitUserInfo{}
}

// LoadFromGinCtx 从gin上下文中加载访问用户信息
// 和 goCommon Oauth2Middleware方法中的设置对应
func (v *VisitUserInfo) LoadFromGinCtx(ctx *gin.Context) {
	v.UserID = ctx.GetString("userId")
	v.UserToken = ctx.GetString("userToken")
	v.IP = ctx.GetString("ip")
	v.ClientType = ctx.GetString("clientType")
	v.UDID = ctx.GetString("udid")
	v.VisitorType = ctx.GetString("visitorType")
	v.UserRoles, _ = ctx.Get("userRoles")
	v.CSFLevel, _ = ctx.Get("csfLevel")
	v.Name = ctx.GetString("name")
	v.AccountType = ctx.GetString("accountType")
	v.DeptPaths, _ = ctx.Get("dept_paths")
}
