package dbaccess

import (
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
	"Authentication/interfaces"
)

type session struct {
	db *sqlx.DB
}

var (
	aOnce sync.Once
	s     *session
)

// NewSession 创建Session操作对象
func NewSession() *session {
	aOnce.Do(func() {
		s = &session{
			db: dbPool,
		}
	})

	return s
}

func (s *session) Get(sessionID string) (*interfaces.Context, error) {
	dbName := common.GetDBName("authentication")
	sqlStr := "select `f_subject`,`f_client_id`,`f_login_session_id`, " +
		"`f_exp`,`f_session_access_token` from %s.t_session where `f_login_session_id` = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	row := s.db.QueryRow(sqlStr, sessionID)

	info := interfaces.Context{}
	if err := row.Scan(
		&info.Subject,
		&info.ClientID,
		&info.SessionID,
		&info.Exp,
		&info.Context,
	); err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *session) Put(ctx interfaces.Context) error {
	dbName := common.GetDBName("authentication")
	sqlStr := "insert into %s.t_session " +
		"(`f_subject`, `f_client_id`, `f_login_session_id`, `f_exp`, `f_session_access_token`) " +
		"values (?, ?, ?, ?, ?) "
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := s.db.Exec(sqlStr,
		ctx.Subject,
		ctx.ClientID,
		ctx.SessionID,
		ctx.Exp,
		ctx.Context,
	); err != nil {
		return err
	}

	return nil
}

func (s *session) Delete(sessionID string) error {
	dbName := common.GetDBName("authentication")
	sqlStr := "delete from %s.t_session where `f_login_session_id` = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := s.db.Exec(sqlStr, sessionID); err != nil {
		return err
	}

	return nil
}

func (s *session) EcronDelete(exp int64) error {
	dbName := common.GetDBName("authentication")
	sqlStr := "delete from %s.t_session where f_exp < ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := s.db.Exec(sqlStr, exp); err != nil {
		return err
	}

	return nil
}
