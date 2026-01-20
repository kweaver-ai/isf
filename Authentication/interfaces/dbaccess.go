// Package interfaces db层定义
package interfaces

import (
	"context"
	"database/sql"
	"time"
)

//go:generate mockgen -package mock -source ../interfaces/dbaccess.go -destination ../interfaces/mock/mock_dbaccess.go

// DBRegisterInfo OAuth2客户端开放注册信息
type DBRegisterInfo struct {
	ClientID               string
	ClientName             string
	ClientSecret           string
	GrantTypes             string
	ResponseTypes          string
	Scope                  string
	RedirectURIs           string
	PostLogoutRedirectURIs string
	Metadata               []byte
}

// DBRegister Context操作接口
type DBRegister interface {
	// CreateClient 创建客户端
	CreateClient(client *DBRegisterInfo) error
}

// Context Context信息
type Context struct {
	Subject   string
	ClientID  string
	SessionID string
	Context   string
	Exp       int64
}

// DBSession Context操作接口
type DBSession interface {
	// Get 获取Context
	Get(sessionID string) (*Context, error)

	// Put 存储Context
	Put(ctx Context) error

	// Delete 删除Context
	Delete(sessionID string) error

	// EcronDelete 定时删除Context
	EcronDelete(exp int64) error
}

// ConfigKey 获取配置类型
type ConfigKey int

// ConfigKey 枚举值
const (
	_ ConfigKey = iota

	RememberFor     // 记住登录状态时间，单位为秒
	RememberVisible // 记住登录状态按钮是否可见

	EnableIDCardLogin  // 是否允许身份证号登录
	EnablePWDLock      // 是否允许密码锁定
	EnableThirdPWDLock // 是否允许第三方密码锁定
	PWDErrCnt          // 密码错误次数上限
	PWDLockTime        // 密码锁定时间，单位为分钟
	VCodeConfig        // 图形验证码配置
	SMSExpiration      // 匿名登录短信验证码过期时间，单位为分钟
	LimitRedirectURI   // 限制客户端的redirect url
	TriSystemStatus    // 三权分立状态
)

// VCodeLoginConfig 图形验证码 配置信息
type VCodeLoginConfig struct {
	Enable    bool // 是否启用图形验证码
	PWDErrCnt int  // 密码错误次数上限
}

// Config 认证配置
type Config struct {
	RememberFor     int
	RememberVisible bool

	EnableIDCardLogin  bool
	EnablePWDLock      bool
	EnableThirdPWDLock bool
	PWDErrCnt          int
	PWDLockTime        int
	VCodeConfig        VCodeLoginConfig
	SMSExpiration      int
	LimitRedirectURI   map[string]bool
	TriSystemStatus    bool
}

// DBConf 认证配置操作接口
type DBConf interface {
	// GetConfig 获取认证配置
	GetConfig(configKeys map[ConfigKey]bool) (cfg Config, err error)

	// GetConfigFromShareMgnt 获取认证配置
	GetConfigFromShareMgnt(ctx context.Context, configKeys map[ConfigKey]bool) (cfg Config, err error)

	// SetConfig 设置认证配置
	SetConfig(configKeys map[ConfigKey]bool, cfg Config) (err error)
}

// DBAccessTokenPerm 访问令牌权限配置接口
type DBAccessTokenPerm interface {
	CheckAppAccessTokenPerm(appID string) (bool, error)
	AddAppAccessTokenPerm(appID string) error
	DeleteAppAccessTokenPerm(appID string) error
	GetAllAppAccessTokenPerm() (permApps []string, err error)
}

// DBLogin 认证接口
type DBLogin interface {
	GetDomainStatus() (enablePrefixMatch bool, err error)
}

// DBOutbox 发件箱数据库接口
type DBOutbox interface {
	// 批量添加outbox消息
	AddOutboxInfos(businessType int, messages []string, tx *sql.Tx) error
	// 获取推送消息
	GetPushMessage(businessType int, tx *sql.Tx) (messageID int64, message string, err error)
	// 根据ID删除outbox消息
	DeleteOutboxInfoByID(messageID int64, tx *sql.Tx) error
}

// DBAnonymousSMS 数据访问层匿名账户短信验证码
type DBAnonymousSMS interface {
	// 创建新的匿名账户短信验证码信息
	Create(ctx context.Context, info *AnonymousSMSInfo) error

	// GetInfoByID 根据ID获取匿名验证码记录
	GetInfoByID(ctx context.Context, id string) (*AnonymousSMSInfo, error)

	// GetExpiredRecords 获取已失效的匿名验证码记录的ID集合
	GetExpiredRecords(ctx context.Context, expiration time.Duration) ([]string, error)

	// DeleteByID 根据ID删除匿名验证码记录
	DeleteByIDs(ctx context.Context, ids []string) error

	// 删除仍在有效期内的匿名验证码
	DeleteRecordWithinValidityPeriod(ctx context.Context, phoneNumber, anonymityID string, expiration time.Duration) error
}

// DBHydra hydra数据库操作接口
type DBHydra interface {
	DeleteExpiredAssertions() (int64, error)
}

// DBTicket 数据访问层Ticket
type DBTicket interface {
	// Create 创建新的ticket
	Create(ctx context.Context, ticketInfo *TicketInfo) error

	// GetTicketByID 根据ID获取ticket信息
	GetTicketByID(ctx context.Context, id string) (*TicketInfo, error)

	// GetExpiredRecords 获取已失效的单点登录凭据的ID集合
	GetExpiredRecords(ctx context.Context, expiration time.Duration) ([]string, error)

	// DeleteByIDs 根据ID批量删除单点登录凭据
	DeleteByIDs(ctx context.Context, ids []string) error
}

// DBFlowClean 数据访问层FlowClean
type DBFlowClean interface {
	// CleanExpiredRefresh 清理过期refresh time 有效期 单位秒
	CleanExpiredRefresh(t int64) (err error)

	// GetAllExpireFlowIDs 获取所有过期的flow_id
	GetAllExpireFlowIDs(t int64) (id []string, err error)

	// CleanFlow 删除过期flow信息
	CleanFlow(ids []string) (err error)
}

// UnorderedOutboxStatus 无序outbox任务状态
type UnorderedOutboxStatus int64

const (
	// OutboxNotStarted 未开始
	OutboxNotStarted UnorderedOutboxStatus = iota
	// OutboxInProgress 处理中
	OutboxInProgress
)

// UnorderedOutbox 无序outbox信息
type UnorderedOutbox struct {
	ID        string
	Message   string
	Status    UnorderedOutboxStatus
	CreatedAt int64
	UpdatedAt int64
}

// DBUnorderedOutbox 审计日志异步任务处理接口
type DBUnorderedOutbox interface {
	// GetUnorderedOutboxInfo 获取无序outbox信息
	GetUnorderedOutboxInfo() (unorderedOutbox UnorderedOutbox, exist bool, err error)

	// DeleteUnorderedOutboxInfoByID 根据ID删除无序outbox信息
	DeleteUnorderedOutboxInfoByID(id string) (err error)

	// UpdateUnorderedOutboxUpdateTimeByID 根据ID更新无序outbox更新时间
	UpdateUnorderedOutboxUpdateTimeByID(id string) (isUpdate bool, err error)

	// AddUnorderedOutboxInfo 添加无序outbox信息
	AddUnorderedOutboxInfo(unorderedOutbox UnorderedOutbox) (err error)

	// RestartUnorderedOutboxInfo 重置无序outbox信息状态
	RestartUnorderedOutboxInfo(updatedTime int64) (err error)
}
