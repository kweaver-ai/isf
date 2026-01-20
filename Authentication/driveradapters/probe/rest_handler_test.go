package probe

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newRESTHandler() RESTHandler {
	return &restHandler{}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestPublicProbe(t *testing.T) {
	Convey("getHealth and getAlive", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler := newRESTHandler()
		handler.RegisterPublic(engine)

		Convey("getHealth", func() {
			req := httptest.NewRequest("GET", "/health/ready", http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getAlive", func() {
			req := httptest.NewRequest("GET", "/health/alive", http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestPrivateProbe(t *testing.T) {
	Convey("getHealth and getAlive", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := newRESTHandler()
		private.RegisterPrivate(engine)
	})
}
