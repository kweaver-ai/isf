package dbaccess

import (
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"Authentication/interfaces"
)

func newSessionDB(ptrDB *sqlx.DB) *session {
	return &session{
		db: ptrDB,
	}
}

func TestGet(t *testing.T) {
	Convey("Get, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		session := newSessionDB(db)
		Convey("success", func() {
			sessionID := "52535da8-5528-4b28-9afd-4b1abb30ec2e"
			fields := []string{
				"f_subject",
				"f_client_id",
				"f_login_session_id",
				"f_exp",
				"f_session_access_token",
			}
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields))
			sessionInfo, _ := session.Get(sessionID)
			assert.Equal(t, sessionInfo, nil)
		})
	})
}

func TestPut(t *testing.T) {
	Convey("Put, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		session := newSessionDB(db)
		Convey("success", func() {
			ctx := interfaces.Context{
				Subject:   "test",
				ClientID:  "test",
				SessionID: "52535da8-5528-4b28-9afd-4b1abb30ec2e",
				Exp:       111111,
				Context:   "{test}",
			}
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := session.Put(ctx)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestDelete(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		session := newSessionDB(db)
		Convey("success", func() {
			sessionID := "52535da8-5528-4b28-9afd-4b1abb30ec2e"
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := session.Delete(sessionID)
			assert.Equal(t, httpErr, nil)
		})
	})
}

func TestEcronDelete(t *testing.T) {
	Convey("Delete, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		session := newSessionDB(db)
		Convey("success", func() {
			var exp int64 = 111111
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := session.EcronDelete(exp)
			assert.Equal(t, httpErr, nil)
		})
	})
}
