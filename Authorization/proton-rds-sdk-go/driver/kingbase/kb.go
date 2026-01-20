package kingbase

import (
	"database/sql/driver"

	"github.com/lib/pq"

	"github.com/kweaver-ai/proton-rds-sdk-go/driver/common"
)

func Open(dsn string) (driver.Conn, error) {
	cfg, err := common.ParseMySQLDSN(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := pq.Open(FormatDSN(cfg))
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
	cnct, err := pq.NewConnector(FormatDSN(cfg))
	if err != nil {
		return nil, err
	}
	return &KBCnct{cnct: cnct}, err
}
