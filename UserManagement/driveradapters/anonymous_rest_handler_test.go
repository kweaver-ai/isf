package driveradapters

import (
	"bytes"
	"errors"
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

	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func newAnonymousHandler(dp interfaces.LogicsAnonymous) *anonymousRestHandler {
	return &anonymousRestHandler{
		anonymous: dp,
	}
}

func TestAuthenticaitonAnonymousNoParams(t *testing.T) {
	Convey("AuthenticaitonAnonymous", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsAnonymous(ctrl)
		h := newAnonymousHandler(dp)
		h.RegisterPrivate(r)

		target := "/api/user-management/v1/anonymity-auth"

		Convey("no password", func() {
			type AuthenticationParams struct {
				Account string `json:"account" binding:"required"`
			}

			reqParam := AuthenticationParams{
				Account: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("no account", func() {
			type AuthenticationParams struct {
				Password string `json:"password" binding:"required"`
			}

			reqParam := AuthenticationParams{
				Password: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAuthenticaitonAnonymousError(t *testing.T) {
	Convey("AuthenticaitonAnonymous", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsAnonymous(ctrl)
		h := newAnonymousHandler(dp)
		h.RegisterPrivate(r)

		target := "/api/user-management/v1/anonymity-auth"

		Convey("anonymous Authentication error", func() {
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			dp.EXPECT().Authentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			type AuthenticationParams struct {
				Account  string `json:"account" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			reqParam := AuthenticationParams{
				Account:  "AAxxxx",
				Password: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, testErr.Code)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			dp.EXPECT().Authentication(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			type AuthenticationParams struct {
				Account  string `json:"account" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			reqParam := AuthenticationParams{
				Account:  "BAxxxx",
				Password: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetAnonymous(t *testing.T) {
	Convey("GetAnonymous", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsAnonymous(ctrl)
		h := newAnonymousHandler(dp)
		h.RegisterPrivate(r)

		target := "/api/user-management/v1/anonymity/user_id"
		Convey("get anonymous error", func() {
			testErr := errors.New("unknown")
			dp.EXPECT().GetByID(gomock.Any()).AnyTimes().Return(&interfaces.AnonymousInfo{}, testErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("get anonymous success1", func() {
			dp.EXPECT().GetByID(gomock.Any()).AnyTimes().Return(&interfaces.AnonymousInfo{VerifyMobile: true}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resParam, _ := io.ReadAll(result.Body)
			var jsonV map[string]interface{}
			_ = jsoniter.Unmarshal(resParam, &jsonV)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, true, jsonV["verify_mobile"].(bool))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("get anonymous success2", func() {
			dp.EXPECT().GetByID(gomock.Any()).AnyTimes().Return(&interfaces.AnonymousInfo{VerifyMobile: false}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resParam, _ := io.ReadAll(result.Body)
			var jsonV map[string]interface{}
			_ = jsoniter.Unmarshal(resParam, &jsonV)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, false, jsonV["verify_mobile"].(bool))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
