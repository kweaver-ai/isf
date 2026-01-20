// Package logics AnyShare
package logics

import (
	"policy_mgnt/interfaces"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

var dbLicense interfaces.DBLicense

// SetDBLicense 设置实例
func SetDBLicense(i interfaces.DBLicense) {
	dbLicense = i
}

var dbTracePool *sqlx.DB

// SetDBTracePool 设置带有trace数据库连接
func SetDBTracePool(i *sqlx.DB) {
	dbTracePool = i
}

var dbConfig interfaces.DBConfig

// SetDBConfig 设置实例
func SetDBConfig(i interfaces.DBConfig) {
	dbConfig = i
}

var dbOutbox interfaces.DBOutbox

// SetDBOutbox 设置实例
func SetDBOutbox(i interfaces.DBOutbox) {
	dbOutbox = i
}
