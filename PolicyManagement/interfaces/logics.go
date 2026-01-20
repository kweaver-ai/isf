package interfaces

//go:generate mockgen -package mock -source ../interfaces/logics.go -destination ../interfaces/mock/mock_logics.go

import (
	"context"
	"database/sql"
)

// outbox业务类型
const (
	// 默认业务类型
	_ = iota
	// OutboxProductAuthorizedUpdated 产品授权更新
	OutboxProductAuthorizedUpdated
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
	LangType  LangType
	Language
	ClientType ClientType
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

// LangType 语言类型
type LangType int

const (
	// 使用iota，增加新的枚举时，必须按照顺序添加，不得在中间插入
	_ LangType = iota

	// LTZHCN 中文简体
	LTZHCN

	// LTZHTW 中文繁体
	LTZHTW

	// LTENUS 英文
	LTENUS
)

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

// LicenseInfo 许可证信息
type LicenseInfo struct {
	Product             string
	TotalUserQuota      int
	AuthorizedUserCount int
}

// ObjectType 对象类型
type ObjectType int

const (
	_ ObjectType = iota
	ObjectTypeUser
)

// AuthorizedProduct 已授权产品
type AuthorizedProduct struct {
	ID      string
	Type    ObjectType
	Product []string
}

// LogicsLicense 许可证业务逻辑接口
type LogicsLicense interface {
	// GetLicenses 获取许可证信息
	GetLicenses(ctx context.Context, visitor *Visitor) (infos map[string]LicenseInfo, err error)

	// GetAuthorizedProducts 获取已授权产品
	GetAuthorizedProducts(ctx context.Context, visitor *Visitor, userIDs []string) (products map[string]AuthorizedProduct, err error)

	// CheckProductAuthorized 检查产品是否已授权
	CheckProductAuthorized(ctx context.Context, visitor *Visitor, product string) (authorized bool, unauthorizedReason string, err error)

	// UpdateAuthorizedProducts 更新已授权产品
	UpdateAuthorizedProducts(ctx context.Context, visitor *Visitor, products map[string]AuthorizedProduct) (err error)
}

// LogicsEvent  处理来自其他服务的事件
type LogicsEvent interface {
	// 用户创建
	UserCreated(userID string) (err error)

	// 注册用户创建事件
	RegisterUserCreated(f func(string) error)

	// 用户状态改变
	UserStatusChanged(userID string, status bool) (err error)

	// 注册用户状态改变事件
	RegisterUserStatusChanged(f func(string, bool) error)

	// 用户删除
	UserDeleted(userID string) (err error)

	// 注册用户删除事件
	RegisterUserDeleted(f func(string) error)
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
