package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver" // 注册数据库驱动
	dm "github.com/kweaver-ai/proton_dm_dialect_go"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var sqlDB *sql.DB

func getDBConnURL() string {
	var host, port, user, pwd, name, dsnSuffix, dsn, systemID string
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user = os.Getenv("DB_USER")
	pwd = os.Getenv("DB_PASSWORD")
	name = os.Getenv("DB_NAME")
	if name == "" {
		name = "policy_mgnt"
	}
	systemID = os.Getenv("DB_SYSTEM_ID")
	name = fmt.Sprintf("%s%s", systemID, name)

	dsnSuffix = "charset=utf8&parseTime=True&loc=Local"
	if isIPv6 := strings.Contains(host, ":"); isIPv6 {
		host = fmt.Sprintf("[%s]", host)
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pwd, host, port, name, dsnSuffix)
	return dsn
}

// ConnectDB return orm.DB
func ConnectDB() (*gorm.DB, error) {
	var err error
	dbtype := os.Getenv("DB_TYPE")
	if db == nil {
		var driver, url string
		driver = os.Getenv("DB_DRIVER")
		if driver == "" {
			driver = "proton-rds"
			os.Setenv("DB_DRIVER", driver)
		}
		if driver == "proton-rds" {
			url = getDBConnURL()
		} else {
			url = os.Getenv("DB_URL")
		}
		operation, err := sql.Open(driver, url)
		if err != nil {
			return nil, err
		}
		var dialector gorm.Dialector
		if driver == "sqlite3" {
			dialector = sqlite.Open(url)
		} else { // driver=proton-rds
			if dbtype == "DM8" {
				dialector = dm.New(dm.Config{Conn: operation})
			} else if strings.HasPrefix(dbtype, "KDB") {
				dialector = postgres.New(postgres.Config{Conn: operation})
			} else { // mysql mariadb tidb
				dialector = mysql.New(mysql.Config{Conn: operation})
			}
		}
		// gorm logger 配置
		loggerDefault := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Warn,            // Log level
				Colorful:                  false,                  // 彩色打印
				IgnoreRecordNotFoundError: true,                   // 关闭 not found错误
			},
		)
		db, err = gorm.Open(dialector, &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 loggerDefault, // gorm的log设置
		})
		if err != nil {
			return nil, err
		}
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic("连接数据库失败：" + err.Error())
		}

		// TODO: 通过文件配置修改
		maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
		if maxIdleConn < 2 {
			maxIdleConn = 2
		}
		maxOpenConn, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
		if maxOpenConn < 0 {
			maxOpenConn = 0
		}
		sqlDB, err = db.DB()
		if err != nil {
			return db, err
		}
		sqlDB.SetMaxIdleConns(maxIdleConn)
		sqlDB.SetMaxOpenConns(maxOpenConn)
	}

	return db, err
}

// DisconnectDB ...
func DisconnectDB() (err error) {
	if db != nil {
		sqlDB, _ = db.DB()
		err = sqlDB.Close()
		db = nil
		return err
	}
	return nil
}
