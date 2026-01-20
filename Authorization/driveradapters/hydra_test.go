// Package driveradapters
package driveradapters

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"Authorization/common"
	"Authorization/interfaces"
	dmock "Authorization/interfaces/mock"
)

var hydraHandler interfaces.DepHTTPSvc

func hydraHandle(t *testing.T, w http.ResponseWriter, r *http.Request) {
	urll := r.URL.EscapedPath()
	method := r.Method
	body, _ := io.ReadAll(r.Body)
	resCode, resBody := hydraHandler.HandleRequest(method, urll, string(body))
	w.WriteHeader(resCode)
	_, err := w.Write(resBody)
	assert.Equal(t, err, nil)
}

func TestHydraHTTPInterface(t *testing.T) {
	Convey("Introspect", t, func() {
		tsHydra := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hydraHandle(t, w, r)
		}))
		u, err := url.Parse(tsHydra.URL)
		if err != nil {
			t.Fatalf("%v", err)
		}

		port, _ := strconv.Atoi(u.Port())
		common.SvcConfig.OAuthAdminHost = u.Hostname()
		common.SvcConfig.OAuthAdminPort = port

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hydraMock := dmock.NewMockDepHTTPSvc(ctrl)
		hydraHandler = hydraMock

		hydra := newHydra()

		Convey("token内省成功", func() {
			tokeninfo, _ := jsoniter.Marshal(gin.H{
				"active":    true,
				"sub":       "user_1",
				"scope":     "some-scope",
				"client_id": "some-client-id",
				"ext": map[string]any{
					"visitor_type": "realname",
					"login_ip":     "1.1.1.1",
					"udid":         "aa-bb-cc-dd",
					"account_type": "other",
					"client_type":  "unknown",
				},
			})

			tmpInfo := interfaces.TokenIntrospectInfo{
				Active:     true,
				VisitorID:  "user_1",
				Scope:      "some-scope",
				ClientID:   "some-client-id",
				VisitorTyp: interfaces.RealName,
				LoginIP:    "1.1.1.1",
				Udid:       "aa-bb-cc-dd",
				AccountTyp: interfaces.Other,
				ClientTyp:  interfaces.Unknown,
			}

			hydraMock.EXPECT().HandleRequest("POST", "/admin/oauth2/introspect", "token=some-token-id").Return(http.StatusOK, tokeninfo)
			info, err := hydra.Introspect("some-token-id")
			assert.Equal(t, err, nil)
			assert.Equal(t, info, tmpInfo)
		})
	})
}
