package driveradapters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestPublicReady(t *testing.T) {
	Convey("ready", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testCRestHandler := NewHealthHandler()
		testCRestHandler.RegisterPublic(r)

		req := httptest.NewRequest("GET", "/health/ready", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}

func TestPrivateReady(t *testing.T) {
	Convey("ready", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testCRestHandler := NewHealthHandler()
		testCRestHandler.RegisterPrivate(r)

		req := httptest.NewRequest("GET", "/health/ready", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}

func TestPublicLive(t *testing.T) {
	Convey("live", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testCRestHandler := NewHealthHandler()
		testCRestHandler.RegisterPublic(r)

		req := httptest.NewRequest("GET", "/health/live", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}

func TestPrivateLive(t *testing.T) {
	Convey("live", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testCRestHandler := NewHealthHandler()
		testCRestHandler.RegisterPrivate(r)

		req := httptest.NewRequest("GET", "/health/live", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}
