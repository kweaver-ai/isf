package dbaccess

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
)

type login struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	lOnce sync.Once
	l     *login
)

// NewLogin 创建Login操作对象
func NewLogin() *login {
	lOnce.Do(func() {
		l = &login{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return l
}

// GetDomainStatus 检查是否有开启的域
func (l *login) GetDomainStatus() (enablePrefixMatch bool, err error) {
	var status int
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "select f_status from %s.t_domain where f_status = 1"
	strSQL = fmt.Sprintf(strSQL, dbName)
	err = l.db.QueryRow(strSQL).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return
	}

	return true, nil
}
