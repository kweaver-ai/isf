// Package accesstokenperm 协议层
package accesstokenperm

import (
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
	acURL    = "/api/authentication/v1/access-token-perm/app/b550af01-06d0-446d-be5b-b44cfcd97906"
	acGetURL = "/api/authentication/v1/access-token-perm/app"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func newAccessHandler(accesstokenperm interfaces.AccessTokenPerm, hydra interfaces.Hydra) *restHandler {
	return &restHandler{
		accessTokenPerm: accesstokenperm,
		hydra:           hydra,
	}
}

//nolint:dupl
func TestSetAppAccessTokenPermPvt(t *testing.T) {
	Convey("TestSetAppAccessTokenPermPvt", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPrivate(engine)

		Convey("SetAppAccessTokenPerm 逻辑层失败", func() {
			lErr := errors.New("test")
			accesstokenperm.EXPECT().SetAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(lErr)
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SetAppAccessTokenPerm 成功", func() {
			accesstokenperm.EXPECT().SetAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSetAppAccessTokenPermPub(t *testing.T) {
	Convey("TestSetAppAccessTokenPermPub", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPublic(engine)

		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: "99e68254-77cd-4280-b147-f987ec53bc7e"}
		testErr := errors.New("test")

		Convey("token校验失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, testErr)
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
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
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SetAppAccessTokenPerm 逻辑层失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().SetAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SetAppAccessTokenPerm 成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().SetAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest("PUT", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:dupl
func TestDeleteAppAccessTokenPermPvt(t *testing.T) {
	Convey("TestDeleteAppAccessTokenPermPvt", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPrivate(engine)

		Convey("DeleteAppAccessTokenPerm 逻辑层失败", func() {
			lErr := errors.New("test")
			accesstokenperm.EXPECT().DeleteAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(lErr)
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("DeleteAppAccessTokenPerm 成功", func() {
			accesstokenperm.EXPECT().DeleteAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteAppAccessTokenPermPub(t *testing.T) {
	Convey("TestDeleteAppAccessTokenPermPub", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPublic(engine)

		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: "99e68254-77cd-4280-b147-f987ec53bc7e"}
		testErr := errors.New("test")

		Convey("token校验失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, testErr)
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
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
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("DeleteAppAccessTokenPerm 逻辑层失败", func() {
			lErr := errors.New("test")
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().DeleteAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(lErr)
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("DeleteAppAccessTokenPerm 成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().DeleteAppAccessTokenPerm(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest("DELETE", acURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetAllAppAccessTokenPermPvt(t *testing.T) {
	Convey("TestGetAllAppAccessTokenPermPvt", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPrivate(engine)

		Convey("GetAllAppAccessTokenPerm 逻辑层失败", func() {
			lErr := errors.New("test")
			accesstokenperm.EXPECT().GetAllAppAccessTokenPerm(gomock.Any(), gomock.Any()).Return([]string{}, lErr)
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAllAppAccessTokenPerm 成功", func() {
			accesstokenperm.EXPECT().GetAllAppAccessTokenPerm(gomock.Any(), gomock.Any()).Return([]string{"d8521454-c8ff-402f-9ccb-e7f2c0a0723c"}, nil)
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			message, _ := io.ReadAll(result.Body)
			resBody := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(message, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody[0].(string), "d8521454-c8ff-402f-9ccb-e7f2c0a0723c")
			assert.Equal(t, len(resBody), 1)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetAllAppAccessTokenPermPub(t *testing.T) {
	Convey("TestGetAllAppAccessTokenPermPub", t, func() {
		test := setGinMode()
		defer test()
		engine := gin.New()
		engine.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		accesstokenperm := mock.NewMockAccessTokenPerm(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		acHandler := newAccessHandler(accesstokenperm, hydra)
		acHandler.RegisterPublic(engine)

		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: "99e68254-77cd-4280-b147-f987ec53bc7e"}
		testErr := errors.New("test")

		Convey("token校验失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, testErr)
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
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
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAllAppAccessTokenPerm 逻辑层失败", func() {
			lErr := errors.New("test")
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().GetAllAppAccessTokenPerm(gomock.Any(), gomock.Any()).Return([]string{}, lErr)
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAllAppAccessTokenPerm 成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			accesstokenperm.EXPECT().GetAllAppAccessTokenPerm(gomock.Any(), gomock.Any()).Return([]string{"d8521454-c8ff-402f-9ccb-e7f2c0a0723c"}, nil)
			req := httptest.NewRequest("GET", acGetURL, http.NoBody)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			result := w.Result()
			message, _ := io.ReadAll(result.Body)
			resBody := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(message, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody[0].(string), "d8521454-c8ff-402f-9ccb-e7f2c0a0723c")
			assert.Equal(t, len(resBody), 1)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
