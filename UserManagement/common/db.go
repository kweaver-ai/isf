// Package common dbPool
package common

import (
	"fmt"
	"os"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"gopkg.in/yaml.v2"

	// _ 注册proton-rds驱动
	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
)

var (
	dbOnce sync.Once
	dbPool *sqlx.DB = nil
)

// NewDBPool 获取数据库连接池
func NewDBPool() *sqlx.DB {
	dbOnce.Do(func() {
		dbLog := NewLogger()
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

		connInfo.Database = GetDBName(connInfo.Database)

		dbPool, err = sqlx.NewDB(&connInfo)
		if err != nil {
			dbLog.Fatalf("new db operator failed: %v\n", err)
		}
	})

	return dbPool
}

// GetDBName 获取数据库名
func GetDBName(dbName string) string {
	return fmt.Sprintf("%s%s", SvcConfig.SystemID, dbName)
}
