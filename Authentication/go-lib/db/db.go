// Package db 数据库连接
package db

import (
	"database/sql"
	"fmt"

	"github.com/kweaver-ai/go-lib/util"
)

// DBConfig 数据库配置信息
type DBConfig struct {
	DriverName   string `yaml:"driver_name"`
	User         string `yaml:"user_name"`
	Password     string `yaml:"user_pwd"`
	Host         string `yaml:"db_host"`
	Port         int    `yaml:"db_port"`
	Database     string `yaml:"db_name"`
	Charset      string `yaml:"db_charset"`
	Timeout      int    `yaml:"timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// NewDB 新建数据库连接
func NewDB(dbConfig *DBConfig) (*sql.DB, error) {
	dbConfig.Host = util.ParseHost(dbConfig.Host)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=%ds&readTimeout=%ds&writeTimeout=%ds",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.Charset,
		dbConfig.Timeout,
		dbConfig.ReadTimeout,
		dbConfig.WriteTimeout)
	// Open may just validate its arguments without creating a connection to the database.
	// To verify that the data source name is valid, call Ping.
	// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
	// Thus, the Open function should be called just once.
	// It is rarely necessary to close a DB.
	db, err := sql.Open(dbConfig.DriverName, dsn)
	if err != nil {
		return nil, err
	}
	// If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
