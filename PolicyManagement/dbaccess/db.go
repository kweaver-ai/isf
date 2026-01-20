// Package dbaccess 数据访问层
package dbaccess

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

var dbTracePool *sqlx.DB = nil

// SetDBTracePool 设置带有trace数据库连接
func SetDBTracePool(i *sqlx.DB) {
	dbTracePool = i
}
