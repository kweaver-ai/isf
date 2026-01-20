package kingbase

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type KBConn struct {
	conn driver.Conn
}

func (KC KBConn) Prepare(query string) (driver.Stmt, error) {
	return KC.conn.Prepare(Rebind(query))
}
func (KC KBConn) Begin() (driver.Tx, error) {
	return KC.conn.Begin()
}

func (KC KBConn) Close() error {
	return KC.conn.Close()
}

func Rebind(query string) string {
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)
		rqb = append(rqb, '$')

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)
		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
