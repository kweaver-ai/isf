package interfaces

//go:generate mockgen -package mock -source ../interfaces/dbaccess.go -destination ../interfaces/mock/mock_dbaccess.go

import (
	"context"
	"database/sql"
)

// ProductInfo 产品信息
type ProductInfo struct {
	AccountID string
	Product   string
}

// 日志信息，记录变更的产品信息
type LogProductInfo struct {
	CurrentProducts []string
	FutureProducts  []string
}

// DBLicense 许可证数据库访问接口
type DBLicense interface {
	// GetAuthorizedProducts 获取已授权产品
	GetAuthorizedProducts(ctx context.Context, userIDs []string) (products map[string]AuthorizedProduct, err error)

	// DeleteAuthorizedProducts 删除已授权产品
	DeleteAuthorizedProducts(ctx context.Context, products []ProductInfo, tx *sql.Tx) (err error)

	// AddAuthorizedProducts 新增已授权产品
	AddAuthorizedProducts(ctx context.Context, products []ProductInfo, tx *sql.Tx) (err error)

	// DeleteUserAuthorizedProducts 删除用户已授权产品
	DeleteUserAuthorizedProducts(ctx context.Context, userID string, tx *sql.Tx) (err error)

	// GetProductsAuthorizedCount 获取产品已授权用户数量
	GetProductsAuthorizedCount(ctx context.Context, product string) (count int, err error)
}

type DBConfig interface {
	// GetConfig 获取配置
	GetConfig(ctx context.Context, key string) (value string, err error)
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
