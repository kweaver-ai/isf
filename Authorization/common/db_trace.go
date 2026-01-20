// Package common dbTracePool
package common

import (
	"database/sql"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	protonRDS "github.com/kweaver-ai/proton-rds-sdk-go/driver"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	sqlhooks "github.com/qustavo/sqlhooks/v2"
)

var (
	dbTraceOnce sync.Once
	dbTracePool *sqlx.DB = nil
)

const (
	traceDriverName = "rds-trace"
)

func initTraceHook() {
	hook := &observable.RDSHook{System: "rds"}
	sql.Register(traceDriverName, sqlhooks.Wrap(new(protonRDS.RDSDriver), hook))
}

// NewDBTracePool 获取带有trace数据库连接池
func NewDBTracePool() *sqlx.DB {
	dbTraceOnce.Do(func() {
		initTraceHook()
		dbLog := NewLogger()
		connInfo := dbConfig
		connInfo.CustomDriver = traceDriverName
		var err error
		dbTracePool, err = sqlx.NewDB(&connInfo)
		if err != nil {
			// 判断err里的错误
			if err.Error() == "driver must implement driver.ConnBeginTx" {
				connInfo.CustomDriver = "proton-rds"
				dbTracePool, err = sqlx.NewDB(&connInfo)
			}
			if err != nil {
				dbLog.Errorf("new db operator failed; error:%s, connInfo:%+v, configLoader.DB:%+v",
					err.Error(), connInfo, dbConfig.Database)
				panic(err)
			}
		}
	})

	return dbTracePool
}
