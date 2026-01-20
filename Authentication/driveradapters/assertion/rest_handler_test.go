// Package assertion 协议层
package assertion

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

const (
	getAssertionURL = "/api/authentication/v1/jwt?user_id=ebb09008-fb85-11ed-9b9f-761dc86fda9f"
	tokanHookURL    = "/api/authentication/v1/token-hook"
	appID           = "b550af01-06d0-446d-be5b-b44cfcd97906"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func newAssertionHandler(assertion interfaces.Assertion, hydra interfaces.Hydra) *restHandler {
	return &restHandler{
		assertion: assertion,
		hydra:     hydra,
	}
}

func TestGetAssertionByUserID(t *testing.T) {
	Convey("TestGetAssertionByUserID", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		assertion := mock.NewMockAssertion(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		asHandler := newAssertionHandler(assertion, hydra)
		asHandler.RegisterPrivate(engine)
		asHandler.RegisterPublic(engine)

		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: appID}
		testErr := errors.New("test")

		Convey("token校验失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, testErr)
			req := httptest.NewRequest("GET", getAssertionURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("token过期", func() {
			introspectInfo.Active = false
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", getAssertionURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("未传user_id", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", "/api/authentication/v1/jwt", http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAssertionByUserID 逻辑层失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			assertion.EXPECT().GetAssertionByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return("", testErr)
			req := httptest.NewRequest("GET", getAssertionURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAssertionByUserID 成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			assertion.EXPECT().GetAssertionByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return("xxx", nil)
			req := httptest.NewRequest("GET", getAssertionURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			message, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(message, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody["assertion"].(string), "xxx")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestTokenHook(t *testing.T) {
	Convey("TestTokenHook", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		assertion := mock.NewMockAssertion(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		asHandler := newAssertionHandler(assertion, hydra)
		asHandler.RegisterPrivate(engine)
		asHandler.RegisterPublic(engine)

		reqParam := `
		{
			"request": {"client_id":"a9112f5f-281d-4fe9-8e8a-6b9fd3cd73d7","payload":{"assertion":["xxx"]}}
		}
		`
		testErr := errors.New("test")

		Convey("request.payload不包含assertion，响应204", func() {
			reqParam = `
			{
				"request": {"payload":{}}
			}
			`
			req := httptest.NewRequest("POST", tokanHookURL, bytes.NewReader([]byte(reqParam)))
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request.payload包含assertion，不包含client_id 抛错400", func() {
			reqParam = `
			{
				"request": {"payload":{"assertion":["xxx"]}}
			}
			`
			req := httptest.NewRequest("POST", tokanHookURL, bytes.NewReader([]byte(reqParam)))
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("TokenHook 逻辑层失败", func() {
			assertion.EXPECT().TokenHook(gomock.Any(), gomock.Any()).Return(map[string]interface{}{}, testErr)
			req := httptest.NewRequest("POST", tokanHookURL, bytes.NewReader([]byte(reqParam)))
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("TokenHook 解析成功", func() {
			res := map[string]interface{}{
				"visitor_type": "realname",
				"login_ip":     "",
				"account_type": "other",
				"udid":         "",
				"client_type":  "app",
			}
			assertion.EXPECT().TokenHook(gomock.Any(), gomock.Any()).Return(res, nil)
			req := httptest.NewRequest("POST", tokanHookURL, bytes.NewReader([]byte(reqParam)))
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			message, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(message, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody["session"].(map[string]interface{})["access_token"].(map[string]interface{}), res)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
