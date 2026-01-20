package session

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func newRESTHandler(tmpSession interfaces.Session) RESTHandler {
	return &restHandler{
		session: tmpSession,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestGetSession(t *testing.T) {
	Convey("get session", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionLogics := mock.NewMockSession(ctrl)

		testRestHandler := newRESTHandler(sessionLogics)
		testRestHandler.RegisterPrivate(r)
		target := "/api/authentication/v1/session/c2aa9e98-98d4-41a2-bf2e-94a13059aa09"

		Convey("get err", func() {
			getErr := rest.NewHTTPError("session_id not exist", rest.URINotExist, nil)
			sessionInfo := interfaces.Context{}

			sessionLogics.EXPECT().Get(gomock.Any()).AnyTimes().Return(sessionInfo, getErr)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusNotFound)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			umErr := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, umErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestPutSession(t *testing.T) {
	Convey("get session", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionLogics := mock.NewMockSession(ctrl)

		testRestHandler := newRESTHandler(sessionLogics)
		testRestHandler.RegisterPrivate(r)
		target := "/api/authentication/v1/session/c2aa9e98-98d4-41a2-bf2e-94a13059aa09"

		Convey("subject not exist", func() {
			reqInfo := map[string]interface{}{
				"client_id":    "client1",
				"remember_for": 123456,
				"context": map[string]interface{}{
					"account_type": "other",
					"client_type":  "ios",
					"login_ip":     "10.2.176.204",
					"udid":         "0a-23-fd-dd-aa-dd-xc",
					"visitor_type": "realname",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("subject invalid", func() {
			reqInfo := map[string]interface{}{
				"subject":      "test",
				"client_id":    "client1",
				"remember_for": 123456,
				"context": map[string]interface{}{
					"account_type": "other",
					"client_type":  "ios",
					"login_ip":     "10.2.176.204",
					"udid":         "0a-23-fd-dd-aa-dd-xc",
					"visitor_type": "realname",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)
			sessionLogics.EXPECT().Put(gomock.Any()).AnyTimes().Return(nil)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusCreated)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteSession(t *testing.T) {
	Convey("delete session", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionLogics := mock.NewMockSession(ctrl)

		testRestHandler := newRESTHandler(sessionLogics)
		testRestHandler.RegisterPrivate(r)
		target := "/api/authentication/v1/session/c2aa9e98-98d4-41a2-bf2e-94a13059aa09"

		Convey("delete session success", func() {
			sessionLogics.EXPECT().Delete(gomock.Any()).AnyTimes().Return(nil)
			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
