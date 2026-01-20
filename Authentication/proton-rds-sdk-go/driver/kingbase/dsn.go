package kingbase

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func FormatDSN(cfg mysql.Config) string {
	dsn := ""
	if cfg.User != "" {
		dsn += fmt.Sprintf("user=%s ", cfg.User)
	}
	if cfg.Passwd != "" {
		dsn += fmt.Sprintf("password=%s ", cfg.Passwd)
	}
	if cfg.Addr != "" {
		s := strings.Split(cfg.Addr, ":")
		port := s[len(s)-1]
		host := cfg.Addr[:len(cfg.Addr)-len(port)-1]
		dsn += fmt.Sprintf("host=%s port=%s ", host, port)
	}
	if cfg.DBName != "" {
		dsn += fmt.Sprintf("search_path=%s ", cfg.DBName)
	}
	dsn += fmt.Sprintf("connect_timeout=%d ", cfg.Timeout/(1000*1000*1000))
	dsn += "sslmode=disable dbname=proton"
	return dsn
}
