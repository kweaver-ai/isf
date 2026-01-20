// Package interfaces 逻辑层定义
package interfaces

import (
	"context"
	"database/sql"
	"time"
)

//go:generate mockgen -package mock -source ../interfaces/logics.go -destination ../interfaces/mock/mock_logics.go

// AnonymousCredential 匿名认证凭据
type AnonymousCredential struct {
	Account  string
	Password string
}

// AnonymousLoginInfo 请求信息
type AnonymousLoginInfo struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	Scope        string
	Credential   AnonymousCredential
}

// AnonymousLoginInfo2 请求信息
type AnonymousLoginInfo2 struct {
	ClientID     string // OAuth2客户端唯一标识
	ClientSecret string // OAuth2客户端密码
	Credential   AnonymousCredential
	VCode        AnonymousSMSInfo
	VisitorName  string
}

// SSOCredential 第三方认证凭据
type SSOCredential struct {
	ID     string
	Params interface{}
}

// SSOLoginInfo 请求信息
type SSOLoginInfo struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	Scope        string
	Udids        []string
	IP           string
	Credential   SSOCredential
}

// NCTVcodeType 验证码类型
type NCTVcodeType int

const (
	_ NCTVcodeType = iota

	// ImageVCode 图片验证码
	ImageVCode

	// NumVCode 数字验证码
	NumVCode

	// DualAuthSMS 双因子认证--短信验证码
	DualAuthSMS

	// DualAuthOTP 双因子认证--OTP
	DualAuthOTP
)

// ClientLoginReq 用户登录请求信息
type ClientLoginReq struct {
	Method   string            `json:"method"`
	Account  string            `json:"account"`
	Password string            `json:"password"`
	Option   ClientLoginOption `json:"option"`
}

// ClientLoginOption 用户登录附带选项信息
type ClientLoginOption struct {
	UUID      string       `json:"uuid"`
	VCode     string       `json:"vcode"`
	VCodeType NCTVcodeType `json:"vcodeType"`
}

// AccessTokenReq 获取访问令牌时的请求信息结构体
type AccessTokenReq struct {
	Account      string // 账户名（f_login_name）
	Password     string // 账户密码
	ClientID     string // OAuth2客户端ID
	ClientSecret string // OAuth2客户端密码
}

// Login 登录接口
type Login interface {
	// 客户端登录
	ClientAccountAuth(ctx context.Context, visitor *Visitor, req *ClientLoginReq) (string, error)

	// Anonymous 匿名登录
	Anonymous(visitor *Visitor, req *AnonymousLoginInfo) (*TokenInfo, error)

	// Anonymous2 匿名登录
	Anonymous2(ctx context.Context, visitor *Visitor, req *AnonymousLoginInfo2, referrer string) (*TokenInfo, error)

	// SingleSignOn 单点登录
	SingleSignOn(ctx context.Context, visitor *Visitor, req *SSOLoginInfo) (*TokenInfo, error)

	// PwdAuth 账户密码校验
	PwdAuth(ctx context.Context, visitor *Visitor, req *AccessTokenReq) (*TokenInfo, error)

	// GetAccessToken 获取用户访问令牌
	GetAccessToken(ctx context.Context, visitor *Visitor, reqInfo *AccessTokenReq) (*TokenInfo, error)
}

// Register 注册接口
type Register interface {
	// RegisterClient 公开注册
	PublicRegister(client *RegisterInfo) (ClientInfo, error)
}

// Session 操作接口
type Session interface {
	// Get 获取Context
	Get(sessionID string) (Context, error)

	// Put 存储Context
	Put(ctx Context) error

	// Delete 删除Context
	Delete(sessionID string) error

	// EcronDelete 定时删除Context
	EcronDelete(exp int64) error
}

// Conf 操作接口
type Conf interface {
	// GetConfig 获取认证配置
	GetConfig(ctx context.Context, visitor *Visitor, configKeys map[ConfigKey]bool) (cfg Config, err error)
	// GetConfigFromShareMgnt 获取认证配置
	GetConfigFromShareMgnt(ctx context.Context, visitor *Visitor, configKeys map[ConfigKey]bool) (cfg Config, err error)
	// SetConfig 设置认证配置
	SetConfig(ctx context.Context, visitor *Visitor, configKeys map[ConfigKey]bool, cfg Config) (err error)
}

// ErrorCodeType 错误码类型
type ErrorCodeType int

const (
	_ ErrorCodeType = iota
	// Number 整数类型错误码
	Number
	// Str 字符串类型错误码
	// NOTE: 编目属性值类型已存在枚举值String，所以这里使用Str
	Str
)

// Visitor 请求信息
type Visitor struct {
	ID string

	// TokenID 在 JSON 序列化和反序列化时会被忽略，用于防止令牌在持久化过程中泄露
	// 如需在反序列化时获取 TokenID，请通过代码手动处理
	TokenID   string `json:"-"`
	IP        string
	Mac       string
	UserAgent string
	Type      VisitorType
	Language
	ErrorCodeType
}

// AccessTokenPerm 访问令牌权限配置接口
type AccessTokenPerm interface {
	SetAppAccessTokenPerm(ctx context.Context, visitor *Visitor, appID string) error
	DeleteAppAccessTokenPerm(ctx context.Context, visitor *Visitor, appID string) error
	AppDeleted(appID string) error
	CheckAppAccessTokenPerm(appID string) (bool, error)
	GetAllAppAccessTokenPerm(ctx context.Context, visitor *Visitor) (permApps []string, err error)
}

// Assertion 断言操作
type Assertion interface {
	// 根据应用账户获取断言
	GetAssertionByUserID(ctx context.Context, visitor *Visitor, userID string) (assertion string, err error)

	// 获取断言
	CreateJWK(ctx context.Context, subject string, ttl time.Duration, privateClaims map[string]interface{}) (assertion string, err error)

	TokenHook(assertion, clientID string) (result map[string]interface{}, err error)
}

// HydraSession 操作接口
type HydraSession interface {
	// Delete 删除login、consent会话
	Delete(userID, clientID string) error
}

// OutboxMsg outbox消息结构体
type OutboxMsg struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

// LogicsOutbox 逻辑层 发件箱处理接口
type LogicsOutbox interface {
	// 添加outbox消息
	AddOutboxInfo(opType int, content interface{}, tx *sql.Tx) error

	// 批量添加outbox消息
	AddOutboxInfos(msgs []OutboxMsg, tx *sql.Tx) error

	// RegisterHandlers 注册异步处理函数
	RegisterHandlers(opType int, op func(interface{}) error)

	// notify推送线程
	NotifyPushOutboxThread()
}

// AnonymousSMSInfo 匿名账户短信验证码 信息结构体
type AnonymousSMSInfo struct {
	ID          string // 验证码唯一标识
	PhoneNumber string // 手机号码
	AnonymityID string // 匿名账户ID
	Content     string // 验证码内容
	CreateTime  int64  // 记录创建时间
}

// LogicsAnonymousSMS 逻辑层匿名账户短信验证码
type LogicsAnonymousSMS interface {
	// CreateAndSendVCode 创建并发送匿名账户短信验证码
	CreateAndSendVCode(ctx context.Context, visitor *Visitor, phoneNumber string, account string) (vcodeID string, err error)
	// 校验
	Validate(ctx context.Context, visitor *Visitor, vcodeID, vcode, anonymityID string) (phoneNumber string, err error)
	// UpdateSMSExpiration 更新短信验证码过期时间
	UpdateSMSExpiration(expiration int)
}

// Language 语言类型
type Language int

// 语言类型
const (
	_                  Language = iota
	SimplifiedChinese           // 简体中文
	TraditionalChinese          // 繁体中文
	AmericanEnglish             // 美国英语
)

// TicketReq 生成单点登录凭据请求信息结构体
type TicketReq struct {
	ClientID     string // OAuth2客户端ID
	RefreshToken string // OAuth2刷新令牌
}

// TicketInfo 单点登录凭据信息结构体
type TicketInfo struct {
	ID         string // ticket唯一标识
	UserID     string // 用户唯一标识
	ClientID   string // OAuth2客户端ID
	CreateTime int64  // ticket创建时间
}

// LogicsTicket 逻辑层单点登录凭据
type LogicsTicket interface {
	// CreateTicket 生成单点登录凭据
	CreateTicket(ctx context.Context, visitor *Visitor, reqInfo *TicketReq) (ticketID string, err error)

	// Validate 验证单点登录凭据
	Validate(ctx context.Context, ticketID, clientID string) (userID string, err error)
}

// LogicsAudit 审计日志
type LogicsAudit interface {
	Log(topic string, message interface{}) (err error)
}

// LogicsAuditLogAsyncTask 审计日志
type LogicsAuditLogAsyncTask interface {
	Log(topic string, message interface{}) (err error)
}
