package register

import (
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func newRegister(ha interfaces.DnHydraAdmin, db interfaces.DBRegister) *register {
	return &register{
		hydraAdmin: ha,
		db:         db,
		cost:       12,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestPublicRegister(t *testing.T) {
	Convey("public register", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		db := mock.NewMockDBRegister(ctrl)
		n := newRegister(hydraAdmin, db)
		conf := mock.NewMockDBConf(ctrl)
		n.conf = conf

		registerInfo := &interfaces.RegisterInfo{
			ClientName:             "test",
			GrantTypes:             []string{"authorization_code", "implicit", "refresh_token"},
			ResponseTypes:          []string{"token id_token", "code", "token"},
			Scope:                  "offline all",
			RedirectURIs:           []string{"https://10.2.176.204:9010/callback/xxx"},
			PostLogoutRedirectURIs: []string{"https://10.2.176.204:9010/successful-logout"},
			Metadata: map[string]interface{}{
				"test": "test",
			},
		}

		Convey("get config failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			conf.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{}, testErr)
			clientInfo, err := n.PublicRegister(registerInfo)
			assert.Equal(t, clientInfo, interfaces.ClientInfo{})
			assert.Equal(t, err, testErr)
		})

		Convey("redirect url is not allowed", func() {
			conf.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{
				LimitRedirectURI: map[string]bool{
					"https://10.2.176.204:9010/callback1": true,
				},
			}, nil)
			clientInfo, err := n.PublicRegister(registerInfo)
			assert.Equal(t, clientInfo, interfaces.ClientInfo{})
			assert.Equal(t, err, rest.NewHTTPError("redirect url is not allowed: https://10.2.176.204:9010/callback/xxx", rest.BadRequest, nil))
		})

		Convey("register failed", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			conf.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{
				LimitRedirectURI: map[string]bool{
					"https://10.2.176.204:9010/callback": true,
				},
			}, nil)
			hydraAdmin.EXPECT().PublicRegister(gomock.Any()).AnyTimes().Return(nil, testErr)
			clientInfo, err := n.PublicRegister(registerInfo)
			assert.Equal(t, clientInfo, interfaces.ClientInfo{})
			assert.Equal(t, err, testErr)
		})

		Convey("dbRegister failed", func() {
			clientInfo := &interfaces.ClientInfo{
				ClientID:     "test",
				ClientSecret: "some-secret",
			}
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			conf.EXPECT().GetConfig(gomock.Any()).AnyTimes().Return(interfaces.Config{
				LimitRedirectURI: map[string]bool{},
			}, nil)
			hydraAdmin.EXPECT().PublicRegister(gomock.Any()).AnyTimes().Return(clientInfo, nil)
			db.EXPECT().CreateClient(gomock.Any()).AnyTimes().Return(testErr)
			client, err := n.PublicRegister(registerInfo)
			assert.Equal(t, client, interfaces.ClientInfo{})
			assert.Equal(t, err, testErr)
		})
	})
}
