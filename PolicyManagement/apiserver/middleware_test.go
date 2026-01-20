package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"policy_mgnt/test"
	"policy_mgnt/test/mock_api"
	"policy_mgnt/utils"
	"policy_mgnt/utils/errors"

	"policy_mgnt/utils/gocommon/api"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setMiddlewareRouter(o api.OAuth2) *gin.Engine {
	router := gin.Default()
	router.Use(oauth2Middleware(o))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "ok",
		})
	})
	return router
}

func Test_oauth2Middleware(t *testing.T) {
	viper.Set("oauth_on", true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	o := mock_api.NewMockOAuth2(ctrl)
	o.EXPECT().Introspection(gomock.Any(), "xxx", nil).Return(api.IntrospectionResult{Active: true, ClientID: "test"}, nil)

	teardown := test.SetUpGin(t)
	defer teardown(t)

	viper.SetDefault("visitors", []utils.Visitor{
		{
			Name:     "test",
			ClientID: "test",
		},
	})

	router := setMiddlewareRouter(o)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, string(`{"result":"ok"}`), string(resp.Body.Bytes()))

	o.EXPECT().Introspection(gomock.Any(), "xxx", nil).Return(api.IntrospectionResult{Active: true, ClientID: "test1", Subject: "test2"}, nil)

	router = setMiddlewareRouter(o)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.JSONEq(t, string(`{"result":"ok"}`), string(resp.Body.Bytes()))

	o.EXPECT().Introspection(gomock.Any(), "xxx", nil).Return(api.IntrospectionResult{Active: false}, nil)

	router = setMiddlewareRouter(o)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	router.ServeHTTP(resp, req)

	assert.Equal(t, 401, resp.Code)

	o.EXPECT().Introspection(gomock.Any(), "", nil).AnyTimes().Return(api.IntrospectionResult{Active: false}, nil)

	router = setMiddlewareRouter(o)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "")
	router.ServeHTTP(resp, req)

	assert.Equal(t, 401, resp.Code)

}

func Test_extraBearerToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	actualToken, _ := extraBearerToken(req)
	assert.Equal(t, "xxx", actualToken)

	req.Header.Set("Authorization", "xxx")
	actualToken, err := extraBearerToken(req)
	assert.Equal(t, "invalid", actualToken)
	assert.Equal(t, err, errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token invalid"}))

	req.Header.Set("Authorization", "Other xxx")
	actualToken, err = extraBearerToken(req)
	assert.Equal(t, "invalid", actualToken)
	assert.Equal(t, err, errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token invalid"}))
}
