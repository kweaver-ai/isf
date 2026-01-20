// Package interfaces 杂项接口
package interfaces

import (
	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/observable"
	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"
)

//go:generate mockgen -package mock -source ./misc.go -destination ./mock/mock_misc.go

// MsgBrokerClient 接口嵌入，为driveradapter/mq_client提供mock对象
type MsgBrokerClient interface {
	msqclient.ProtonMQClient
}

// DepHTTPSvc 接口嵌入，为drivenadapter下http接口提供mock对象
type DepHTTPSvc interface {
	// HandleRequest 处理函数
	HandleRequest(method, url string, reqBody interface{}) (code int, resBody []byte)
}

type DnHTTPClient interface {
	httpclient.HTTPClient
}

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

// TraceClient 接口嵌入，提供mock对象
type TraceClient interface {
	observable.Tracer
}
