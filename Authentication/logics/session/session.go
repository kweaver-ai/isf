// Package session 逻辑层
package session

import (
	"database/sql"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	aOnce sync.Once
	s     *session
)

type session struct {
	db interfaces.DBSession
}

// NewSession 创建session处理对象
func NewSession() *session {
	aOnce.Do(func() {
		s = &session{
			db: logics.DBSession,
		}
	})

	return s
}

func (s *session) Get(sessionID string) (interfaces.Context, error) {
	ctx, err := s.db.Get(sessionID)
	if err == sql.ErrNoRows {
		return interfaces.Context{}, rest.NewHTTPError("session_id not exist", rest.URINotExist, nil)
	} else if err != nil {
		return interfaces.Context{}, err
	}

	return *ctx, err
}

func (s *session) Put(ctx interfaces.Context) error {
	if ctx.SessionID == "" {
		return rest.NewHTTPError("session_id is invalid", rest.URINotExist, nil)
	}

	return s.db.Put(ctx)
}

func (s *session) Delete(sessionID string) error {
	return s.db.Delete(sessionID)
}

func (s *session) EcronDelete(exp int64) error {
	return s.db.EcronDelete(exp)
}
