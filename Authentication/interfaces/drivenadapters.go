// Package interfaces 适配器层定义
package interfaces

import (
	"context"
	"crypto/rsa"
)

//go:generate mockgen -package mock -source ../interfaces/drivenadapters.go -destination ../interfaces/mock/mock_drivenadapters.go

// DeviceInfo 设备信息
type DeviceInfo struct {
	Name        string
	ClientType  string
	Description string
}

// LoginInfo 认证需求参数
type LoginInfo struct {
	Subject string
	Context interface{}
}

// ThirdPartyAuthInfo 第三方认证请求信息
type ThirdPartyAuthInfo struct {
	Udids      []string
	IP         string
	Credential ThirdPartyCredential
	Device     DeviceInfo
}

// ThirdPartyCredential 第三方认证凭据
type ThirdPartyCredential struct {
	ID     string
	Params interface{}
}

// DnEacp eacp接口
type DnEacp interface {
	// ThirdPartyAuthentication 第三方认证
	ThirdPartyAuthentication(ctx context.Context, visitor *Visitor, req *ThirdPartyAuthInfo) (*LoginInfo, error)
}

// RegisterInfo 客户端注册信息
type RegisterInfo struct {
	ClientName             string                 `json:"client_name"`
	GrantTypes             []string               `json:"grant_types"`
	ResponseTypes          []string               `json:"response_types"`
	Scope                  string                 `json:"scope"`
	RedirectURIs           []string               `json:"redirect_uris"`
	PostLogoutRedirectURIs []string               `json:"post_logout_redirect_uris"`
	Metadata               map[string]interface{} `json:"metadata"`
}

// AuthorizeInfo 请求信息
type AuthorizeInfo struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	Scope        string
}

// TokenInfo 单点登录响应参数
type TokenInfo struct {
	Code         string
	AccessToken  string
	IDToken      string
	TokenType    string
	Scope        string
	ExpirsesIn   int64
	ResponseType string
}

// ClientInfo 客户端信息
type ClientInfo struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// RefreshTokenIntrospectInfo 刷新令牌内省结果。其它字段暂不列出，只列出当前业务需要的字段
type RefreshTokenIntrospectInfo struct {
	ClientID string
	Sub      string // 令牌主体，也就是用户ID
}

// HydraError 储存 hydra 服务返回的错误
// 用于本服务转发的 hydra 接口，将 hydra 返回的错误原样返回给接口调用方
// 目前只有 DnHydraAdmin interface 内定义的 PublicRegister 方法中会使用
type HydraError struct {
	Status int
	Body   []byte
}

// Error HydraError 实现 Error interface
func (err HydraError) Error() string {
	return string(err.Body)
}

// DnHydraAdmin hydra admin接口
type DnHydraAdmin interface {
	// GetLoginRequestInformation 获取登录请求信息
	GetLoginRequestInformation(loginChallenge string) (*DeviceInfo, error)

	// AcceptLoginRequest 接受登录请求
	AcceptLoginRequest(subject, loginChallenge string) (string, error)

	// AcceptConsentRequest 接受授权请求
	AcceptConsentRequest(scope, consentChallenge string, context interface{}) (string, error)

	// RegisterClient 公开注册
	PublicRegister(client *RegisterInfo) (*ClientInfo, error)

	CreateTrustedPair(publicKey *rsa.PublicKey, issuer, kid string) (err error)
	GetKidTrustedPairByIssuer(issuer string) (trustedPair map[string]bool, err error)
	SetAppAsUserAgent(clientID string) (err error)

	// DeleteSession 清除用户login、consent会话
	DeleteSession(userID, clientID string) (err error)

	// IntrospectRefreshToken 对刷新令牌做内省
	IntrospectRefreshToken(refreshToken string) (info *RefreshTokenIntrospectInfo, err error)

	// GetClientInfo 获取客户端信息
	GetClientInfo(clientID string) (*ClientInfo, error)
}

// DnHydraPublic hydra public接口
type DnHydraPublic interface {
	// AuthorizeRequest 授权请求
	AuthorizeRequest(reqInfo *AuthorizeInfo) (challenge string, context interface{}, err error)

	// VerifierLogin 验证认证请求
	VerifierLogin(redirURL string, context interface{}) (challenge string, newContext interface{}, err error)

	// VerifierConsent 验证授权请求
	VerifierConsent(redirURL, responseType string, context interface{}) (*TokenInfo, error)

	GetTokenEndpoint() (tokenEndpoint string, err error)

	// AssertionForToken 通过断言换取令牌
	AssertionForToken(string, string, string) (*TokenInfo, error)
}

// RoleType 用户角色类型
type RoleType int32

// 用户角色类型定义
const (
	SuperAdmin        RoleType = iota // 超级管理员
	SystemAdmin                       // 系统管理员
	AuditAdmin                        // 审计管理员
	SecurityAdmin                     // 安全管理员
	OrganizationAdmin                 // 组织管理员
	OrganizationAudit                 // 组织审计员
	NormalUser                        // 普通用户
)

// AuthType 用户认证类型
type AuthType int

const (
	_ AuthType = iota

	// Local 本地认证
	Local

	// Domain 域认证
	Domain

	// Third 第三方认证
	Third
)

// LDAPServerType LDAP Server类型
type LDAPServerType int

const (
	_ LDAPServerType = iota

	// WindowAD Window AD
	WindowAD

	// OtherLDAP Other LDAP
	OtherLDAP
)

// UserBaseInfo 用户认证基本信息
type UserBaseInfo struct {
	ID             string
	Account        string
	AuthType       AuthType
	PwdErrCnt      int
	PwdErrLastTime int64
	DisableStatus  bool
	LDAPType       LDAPServerType
	DomainPath     string
}

// DnUserManagement UserManagement接口
type DnUserManagement interface {
	// AnonymousAuthentication 匿名认证
	AnonymousAuthentication(ctx context.Context, visitor *Visitor, account, password, referrer string) (bool, error)
	// GetUserNameByUserID 通过用户id获取用户名
	GetUserRolesByUserID(ctx context.Context, visitor *Visitor, userID string) (roleTypes []RoleType, err error)
	// GetAppInfo 检查应用账户是否存在
	GetAppInfo(ctx context.Context, visitor *Visitor, appID string) (appInfo AppInfo, err error)
	// AccountMatch 账户匹配
	AccountMatch(ctx context.Context, visitor *Visitor, account string, idCardLogin, prefixMatch bool) (bool, UserBaseInfo, error)
	// UserAuth 账户认证
	UserAuth(ctx context.Context, visitor *Visitor, userID, password string) (result bool, reason string, err error)
	// UpdatePWDErrInfo 更新账户密码错误信息
	UpdatePWDErrInfo(ctx context.Context, visitor *Visitor, userID string, pwdErrCnt int) (err error)
	// GetAnonymityInfoByID 根据ID获取匿名账户相关信息
	GetAnonymityInfoByID(ctx context.Context, visitor *Visitor, anonymityID string) (verifyMobile bool, err error)
	// GetUserInfo 获取实名账户信息
	GetUserInfo(ctx context.Context, visitor *Visitor, userID string) (info *UserBaseInfo, err error)
}

// AppInfo 应用账户信息
type AppInfo struct {
	ID   string //  应用账户ID
	Name string //  应用账户名称
}

// DnShareMgnt ShareMgnt接口
type DnShareMgnt interface {
	// 校验短信验证码
	UsrmSMSValidate(userID, vcode string) error

	// 校验动态密码
	UsrmOTPValidate(userID, otp string) error

	// 校验图形验证码
	UsrmIMAGECodeValidate(uuid, vcode string) error

	// 域认证
	UsrmDomainAuth(ctx context.Context, loginName, domainPath, password string, ldapType int32) (bool, error)

	// 第三方认证
	UsrmThirdAuth(ctx context.Context, loginName, password string) (bool, error)

	// 发送匿名账户短信验证码
	UsrmSendAnonymousSMSVCode(ctx context.Context, phoneNumber, vcode string) error
}

// DnEacpLog 日志处理接口
type DnEacpLog interface {
	// Publish 发送消息
	Publish(topic string, msg interface{}) error
	// 记录设置应用账户获取任意用户访问令牌权限日志
	OpSetAppAccessTokenPerm(visitor *Visitor, appName string) error
	// 记录删除应用账户获取任意用户访问令牌权限日志
	OpDeleteAppAccessTokenPerm(visitor *Visitor, appName string) error
}

// DnTelemetry ar 可观测性对象
type DnTelemetry interface {
	Log(typ string, message interface{})
}

// DrivenMessageBroker 消息代理接口
type DrivenMessageBroker interface {
	// AnonymousSmsExpUpdated 更新匿名登录短信验证码过期时间
	AnonymousSmsExpUpdated(smsExpiration int) error
}
