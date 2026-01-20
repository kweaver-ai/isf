// Package common dbTracePool
package common

import (
	"database/sql"
	"os"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	protonRDS "github.com/kweaver-ai/proton-rds-sdk-go/driver"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/qustavo/sqlhooks/v2"
	"gopkg.in/yaml.v2"
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

		// 读取sql配置文件
		file, err := os.ReadFile("/sysvol/conf/user-management.yaml")
		if err != nil {
			dbLog.Fatalf("load /sysvol/conf/user-management.yaml failed: %v\n", err)
		}

		secretFile, err := os.ReadFile("/sysvol/secret_conf/secret.yaml")
		if err != nil {
			dbLog.Fatalf("load /sysvol/secret_conf/secret.yaml failed: %v\n", err)
		}

		connInfo := sqlx.DBConfig{}
		if err = yaml.Unmarshal(file, &connInfo); err != nil {
			dbLog.Fatalf("unmarshal yaml file failed: %v\n", err)
		}

		if err = yaml.Unmarshal(secretFile, &connInfo); err != nil {
			dbLog.Fatalf("unmarshal yaml secretFile failed: %v\n", err)
		}
		connInfo.CustomDriver = traceDriverName
		connInfo.Database = GetDBName(connInfo.Database)

		dbTracePool, err = sqlx.NewDB(&connInfo)
		if err != nil {
			// 判断err里
			if err.Error() == "driver must implement driver.ConnBeginTx" {
				connInfo.CustomDriver = "proton-rds"
				dbTracePool, err = sqlx.NewDB(&connInfo)
			}
			if err != nil {
				dbLog.Errorf("new db operator failed; error:%s, connInfo:%+v, configLoader.DB:%+v",
					err.Error(), connInfo, connInfo.Database)
				panic(err)
			}
		}
	})

	return dbTracePool
}
