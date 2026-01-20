package kingbase

import (
	"database/sql/driver"

	"github.com/kweaver-ai/proton-rds-sdk-go/driver/common"
	"github.com/kweaver-ai/proton-rds-sdk-go/driver/kingbase/gokb"
)

func Open(dsn string) (driver.Conn, error) {
	cfg, err := common.ParseMySQLDSN(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := gokb.Open(FormatDSN(cfg))
	if err != nil {
		return nil, err
	}
	return KBConn{conn: conn}, err
}

func OpenConnector(dsn string) (driver.Connector, error) {
	cfg, err := common.ParseMySQLDSN(dsn)
	if err != nil {
		return nil, err
	}
	cnct, err := gokb.NewConnector(FormatDSN(cfg))
	if err != nil {
		return nil, err
	}
	return &KBCnct{cnct: cnct}, err
}
