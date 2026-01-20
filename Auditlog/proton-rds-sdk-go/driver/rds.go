package driver

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"strings"

	"github.com/kweaver-ai/proton-rds-sdk-go/driver/dmdb"
	"github.com/kweaver-ai/proton-rds-sdk-go/driver/goldendb"
	"github.com/kweaver-ai/proton-rds-sdk-go/driver/kingbase"
	"github.com/kweaver-ai/proton-rds-sdk-go/driver/mysql"
	"github.com/kweaver-ai/proton-rds-sdk-go/driver/tidb"
)

type RDSDriver struct {
}

var supportedOpen = map[string]func(string) (driver.Conn, error){
	"MYSQL":    mysql.Open,
	"MARIADB":  mysql.Open,
	"GOLDENDB": goldendb.Open,
	"DM8":      dmdb.Open,
	"TIDB":     tidb.Open,
	"DEFAULT":  mysql.Open,
	"KDB9":     kingbase.Open,
}

var supportedOpenConnector = map[string]func(string) (driver.Connector, error){
	"MYSQL":    mysql.OpenConnector,
	"MARIADB":  mysql.OpenConnector,
	"GOLDENDB": goldendb.OpenConnector,
	"DM8":      dmdb.OpenConnector,
	"TIDB":     tidb.OpenConnector,
	"DEFAULT":  mysql.OpenConnector,
	"KDB9":     kingbase.OpenConnector,
}

func (d RDSDriver) Open(dsn string) (driver.Conn, error) {
	dbType := os.Getenv("DB_TYPE")
	dbType = strings.ToUpper(dbType)
	if v, ok := supportedOpen[dbType]; ok {
		return v(dsn)
	}
	return supportedOpen["DEFAULT"](dsn)
}

func (d RDSDriver) OpenConnector(dsn string) (driver.Connector, error) {
	dbType := os.Getenv("DB_TYPE")
	dbType = strings.ToUpper(dbType)
	if v, ok := supportedOpenConnector[dbType]; ok {
		return v(dsn)
	}
	return supportedOpenConnector["DEFAULT"](dsn)
}

func init() {
	sql.Register("proton-rds", &RDSDriver{})
}
