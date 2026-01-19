package kingbase

import (
	"context"
	"database/sql/driver"
)

type KBCnct struct {
	cnct driver.Connector
}

func (KCT *KBCnct) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := KCT.cnct.Connect(ctx)
	return KBConn{conn: conn}, err
}

func (KCT *KBCnct) Driver() driver.Driver {
	return Driver{}
}
