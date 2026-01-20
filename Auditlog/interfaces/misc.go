// Package interfaces 杂项接口
package interfaces

//go:generate mockgen -package mock -source ./misc.go -destination ./mock/mock_misc.go

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

// TokenIntrospectInfo 令牌内省结果
type TokenIntrospectInfo struct {
	Active     bool        // 令牌状态
	VisitorID  string      // 访问者ID
	Scope      string      // 权限范围
	ClientID   string      // 客户端ID
	VisitorTyp VisitorType // 访问者类型
	// 以下字段只在visitorType=1，即实名用户时才存在
	LoginIP    string      // 登陆IP
	Udid       string      // 设备码
	AccountTyp AccountType // 账户类型
	ClientTyp  ClientType  // 设备类型
}

// VisitorType 访问者类型
type VisitorType int32

// 访问者类型定义
const (
	RealName  VisitorType = 1 // 实名用户
	Anonymous VisitorType = 4 // 匿名用户
	App       VisitorType = 6 // 应用账户
)

// AccountType 登录账号类型
type AccountType int32

// 登录账号类型定义
const (
	Other  AccountType = 0
	IDCard AccountType = 1
)

// ClientType 设备类型
type ClientType int32

// 设备类型定义
const (
	Unknown ClientType = iota
	IOS
	Android
	WindowsPhone
	Windows
	MacOS
	Web
	MobileWeb
	Nas
	ConsoleWeb
	DeployWeb
	Linux
	APP
)

// Hydra 授权服务接口
type Hydra interface {
	// Introspect token内省
	Introspect(token string) (info TokenIntrospectInfo, err error)
}
