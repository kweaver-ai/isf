package common

import (
	"github.com/go-sql-driver/mysql"
)

func ParseMySQLDSN(dsn string) (mysql.Config, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return mysql.Config{}, err
	}
	return *cfg, err
}
