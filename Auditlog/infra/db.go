// Package infra dbPool
package infra

import (
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/common"
	"AuditLog/infra/config"

	// _ 注册proton-rds驱动
	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
)

var (
	dbOnce sync.Once
	dbPool *sqlx.DB = nil
)

// GetDBName 获取拼接数据库名称 system_id + dbname
func GetDBName() string {
	return config.GetDBConfig().Database
}

// NewDBPool 获取数据库连接池
func NewDBPool() *sqlx.DB {
	//dbOnce.Do(func() {
	//	dbLog := common.SvcConfig.Logger
	//	file, err := os.ReadFile("/sysvol/conf/dbrw.yaml")
	//	if err != nil {
	//		dbLog.Fatalf("load /sysvol/conf/dbrw.yaml failed: %v\n", err)
	//	}
	//
	//	if err = yaml.Unmarshal(file, &dbConfig); err != nil {
	//		dbLog.Fatalf("unmarshal yaml file failed: %v\n", err)
	//	}
	//
	//	dbPool, err = sqlx.NewDB(&dbConfig)
	//	if err != nil {
	//		dbLog.Fatalf("new db operator failed: %v\n", err)
	//	}
	//})
	dbOnce.Do(func() {
		dbLog := common.SvcConfig.Logger

		var err error

		dbConfig := config.GetDBConfig()

		dbPool, err = sqlx.NewDB(dbConfig)
		if err != nil {
			dbLog.Errorf("[NewDBPool]: sqlx.NewDB failed: %v\n", err)
		}
	})

	return dbPool
}
