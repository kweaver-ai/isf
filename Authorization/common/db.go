// Package common dbPool
package common

import (
	"fmt"
	"os"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	yaml "gopkg.in/yaml.v2"

	// _ 注册proton-rds驱动
	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
)

var (
	dbOnce sync.Once
	dbPool *sqlx.DB = nil

	// dbonfig 服务配置信息
	dbConfig sqlx.DBConfig
)

// NewDBPool 获取数据库连接池
func NewDBPool() *sqlx.DB {
	dbOnce.Do(func() {
		dbLog := NewLogger()
		file, err := os.ReadFile("/sysvol/conf/service_conf/authorization.yaml")
		if err != nil {
			dbLog.Fatalf("load /sysvol/conf/service_conf/authorization.yaml failed: %v\n", err)
		}

		secretFile, err := os.ReadFile("/sysvol/conf/secret_conf/secret.yaml")
		if err != nil {
			dbLog.Fatalf("load /sysvol/conf/secret_conf/secret.yaml failed: %v\n", err)
		}

		if err = yaml.Unmarshal(file, &dbConfig); err != nil {
			dbLog.Fatalf("unmarshal yaml file failed: %v\n", err)
		}

		if err = yaml.Unmarshal(secretFile, &dbConfig); err != nil {
			dbLog.Fatalf("unmarshal yaml secretFile failed: %v\n", err)
		}

		dbConfig.Database = GetDBName(dbConfig.Database)

		dbPool, err = sqlx.NewDB(&dbConfig)
		if err != nil {
			dbLog.Fatalf("new db operator failed: %v\n", err)
		}
	})

	return dbPool
}

// GetDBName 获取数据库名称
func GetDBName(dbName string) string {
	return fmt.Sprintf("%s%s", SvcConfig.SystemID, dbName)
}
