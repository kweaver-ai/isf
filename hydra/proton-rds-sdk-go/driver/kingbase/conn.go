package kingbase

import (
	"context"
	"database/sql/driver"
	"strconv"
	"strings"
)

type KBConn struct {
	conn driver.Conn
}

func (KC KBConn) ExecContext(ctx context.Context, sql string, args []driver.NamedValue) (driver.Result, error) {
	//sql = Rebind(sql)
	return KC.conn.(driver.ExecerContext).ExecContext(ctx, sql, args)
}

func (KC KBConn) QueryContext(ctx context.Context, sql string, args []driver.NamedValue) (driver.Rows, error) {
	//sql = Rebind(sql)
	return KC.conn.(driver.QueryerContext).QueryContext(ctx, sql, args)
}

func (KC KBConn) PrepareContext(ctx context.Context, sql string) (driver.Stmt, error) {
	//sql = Rebind(sql)
	return KC.conn.Prepare(sql)
}

func (KC KBConn) Prepare(sql string) (driver.Stmt, error) {
	//sql = Rebind(sql)
	return KC.conn.Prepare(sql)
}
func (KC KBConn) Begin() (driver.Tx, error) {
	return KC.conn.Begin()
}

func (KC KBConn) Close() error {
	return KC.conn.Close()
}

func RebindOld(sql string) string {
	rqb := make([]byte, 0, len(sql)+10)

	var i, j int

	for i = strings.Index(sql, "?"); i != -1; i = strings.Index(sql, "?") {
		rqb = append(rqb, sql[:i]...)
		rqb = append(rqb, '$')

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)
		sql = sql[i+1:]
	}

	return string(append(rqb, sql...))
}

func Rebind(sql string) string {
	bytes := []byte(sql)
	for i, c := range bytes {
		if c == '?' {
			bytes[i] = '$'
		}
	}
	newSql := string(bytes)
	return newSql
}
