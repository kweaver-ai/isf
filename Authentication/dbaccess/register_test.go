package dbaccess

import (
	"encoding/json"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"

	"Authentication/interfaces"
)

func newRegisterDB(ptrDB *sqlx.DB) *register {
	return &register{
		db: ptrDB,
	}
}

func TestPublicRegister(t *testing.T) {
	Convey("public register, db is available", t, func() {
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		register := newRegisterDB(db)
		Convey("success", func() {
			str := map[string]interface{}{
				"device": map[string]interface{}{
					"client_type": "ios",
				},
			}
			metadata, _ := json.Marshal(str)
			registerInfo := &interfaces.DBRegisterInfo{
				ClientName:             "test",
				GrantTypes:             "authorization_code | implicit |refresh_token",
				ResponseTypes:          "token | id_token | code | token",
				Scope:                  "offline all",
				RedirectURIs:           "https://10.2.176.204:9010/callback",
				PostLogoutRedirectURIs: "https://10.2.176.204:9010/successful-logout",
				Metadata:               metadata,
			}
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			httpErr := register.CreateClient(registerInfo)
			assert.Equal(t, httpErr, nil)
		})
	})
}
