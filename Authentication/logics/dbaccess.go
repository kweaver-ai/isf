// Package logics AnyShare
package logics

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/interfaces"
)

var (
	// DBSession 接口实例
	DBSession interfaces.DBSession
	// DBRegister 接口实例
	DBRegister interfaces.DBRegister
	// DBConf 接口实例
	DBConf interfaces.DBConf
	// DBAccessTokenPerm 接口实例
	DBAccessTokenPerm interfaces.DBAccessTokenPerm
	// DBLogin 接口实例
	DBLogin interfaces.DBLogin
	// DBPool 实例
	DBPool *sqlx.DB
	// DBOutbox 实例
	DBOutbox interfaces.DBOutbox
	// DBAnonymousSMS 实例
	DBAnonymousSMS interfaces.DBAnonymousSMS
	// DBTracePool 实例
	DBTracePool *sqlx.DB
	// DBTicket 实例
	DBTicket interfaces.DBTicket
	// DBFlowClean 实例
	DBFlowClean interfaces.DBFlowClean
	// DBUnorderedOutbox 实例
	DBUnorderedOutbox interfaces.DBUnorderedOutbox
)

// SetDBSession 设置实例
func SetDBSession(i interfaces.DBSession) {
	DBSession = i
}

// SetDBRegister 设置实例
func SetDBRegister(i interfaces.DBRegister) {
	DBRegister = i
}

// SetDBConf 设置实例
func SetDBConf(i interfaces.DBConf) {
	DBConf = i
}

// SetDBAccessTokenPerm 设置实例
func SetDBAccessTokenPerm(i interfaces.DBAccessTokenPerm) {
	DBAccessTokenPerm = i
}

// SetDBLogin 设置实例
func SetDBLogin(i interfaces.DBLogin) {
	DBLogin = i
}

// SetDBPool  设置实例
func SetDBPool(i *sqlx.DB) {
	DBPool = i
}

// SetDBOutbox 设置实例
func SetDBOutbox(i interfaces.DBOutbox) {
	DBOutbox = i
}

// SetDBAnonymousSMS 设置实例
func SetDBAnonymousSMS(i interfaces.DBAnonymousSMS) {
	DBAnonymousSMS = i
}

// SetDBTracePool 设置实例
func SetDBTracePool(i *sqlx.DB) {
	DBTracePool = i
}

// SetDBTicket 设置实例
func SetDBTicket(i interfaces.DBTicket) {
	DBTicket = i
}

// SetDBFlowClean 设置实例
func SetDBFlowClean(i interfaces.DBFlowClean) {
	DBFlowClean = i
}

func SetDBUnorderedOutbox(i interfaces.DBUnorderedOutbox) {
	DBUnorderedOutbox = i
}
