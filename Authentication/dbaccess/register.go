package dbaccess

import (
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type register struct {
	db  *sqlx.DB
	log common.Logger
}

var (
	rOnce sync.Once
	r     *register
)

// NewRegister 创建Register操作对象
func NewRegister() *register {
	rOnce.Do(func() {
		r = &register{
			db:  dbPool,
			log: common.NewLogger(),
		}
	})

	return r
}

func (r *register) CreateClient(client *interfaces.DBRegisterInfo) error {
	dbName := common.GetDBName("authentication")
	sqlStr := "insert into %s.t_client_public " +
		"(`id`, `client_name`, `client_secret`, `redirect_uris`, `grant_types`, `response_types`, `scope`, `post_logout_redirect_uris`, `metadata`)" +
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?) "
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := r.db.Exec(sqlStr,
		client.ClientID,
		client.ClientName,
		client.ClientSecret,
		client.RedirectURIs,
		client.GrantTypes,
		client.ResponseTypes,
		client.Scope,
		client.PostLogoutRedirectURIs,
		client.Metadata,
	); err != nil {
		r.log.Errorln(err, sqlStr, client)
		return err
	}

	return nil
}
