package conf

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
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	configSchema "Authentication/driveradapters/jsonschema/config_schema"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func newRESTHandler(conf interfaces.Conf, hydra interfaces.Hydra) RESTHandler {
	return &restHandler{
		conf:  conf,
		hydra: hydra,
		keyConfMap: map[string]interfaces.ConfigKey{
			"remember_for":             interfaces.RememberFor,
			"remember_visible":         interfaces.RememberVisible,
			"anonymous_sms_expiration": interfaces.SMSExpiration,
		},
		resKeyConfMap: map[interfaces.ConfigKey]string{
			interfaces.RememberFor:     "remember_for",
			interfaces.RememberVisible: "remember_visible",
			interfaces.SMSExpiration:   "anonymous_sms_expiration",
		},
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestPrivateGetConfig(t *testing.T) {
	Convey("get config", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("authentication")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		re := mock.NewMockConf(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		testRestHandler := newRESTHandler(re, hydra)
		testRestHandler.RegisterPrivate(r)

		testErr := errors.New("some error")
		target := "/api/authentication/v1/config/remember_for"

		Convey("fields传除remember_for之外的参数", func() {
			target += ",remember"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetConfig error", func() {
			re.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.Config{}, testErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields包含remember_for", func() {
			re.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.Config{RememberFor: 600}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBody, &resBody)
			assert.Equal(t, resBody["remember_for"].(float64), float64(600))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestPublicGetConfig(t *testing.T) {
	Convey("get config", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("authentication")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		re := mock.NewMockConf(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		testRestHandler := newRESTHandler(re, hydra)
		testRestHandler.RegisterPublic(r)

		testErr := errors.New("some error")
		target := "/api/authentication/v1/config/remember_for,remember_visible,anonymous_sms_expiration"
		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: "user_id"}

		Convey("introspect error", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, testErr)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("introspect failed", func() {
			introspectInfo.Active = false
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields传除remember_for、remember_visible、anonymous_sms_expiration之外的参数", func() {
			target += ",remember"
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetConfig error", func() {
			re.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.Config{}, testErr)
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields包含remember_for、remember_visible、anonymous_sms_expiration", func() {
			re.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(interfaces.Config{RememberFor: 600, RememberVisible: true, SMSExpiration: 2}, nil)
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBody, &resBody)
			assert.Equal(t, resBody["remember_for"].(float64), float64(600))
			assert.Equal(t, resBody["remember_visible"].(bool), true)
			assert.Equal(t, resBody["anonymous_sms_expiration"].(float64), float64(2))
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestPublicSetConfig(t *testing.T) {
	Convey("set config", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		re := mock.NewMockConf(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		setConfigSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(configSchema.SetConfigSchemaStr))
		assert.Equal(t, err, nil)

		handler := &restHandler{
			conf:  re,
			hydra: hydra,
			keyConfMap: map[string]interfaces.ConfigKey{
				"remember_for":             interfaces.RememberFor,
				"remember_visible":         interfaces.RememberVisible,
				"anonymous_sms_expiration": interfaces.SMSExpiration,
			},
			resKeyConfMap: map[interfaces.ConfigKey]string{
				interfaces.RememberFor:     "remember_for",
				interfaces.RememberVisible: "remember_visible",
				interfaces.SMSExpiration:   "anonymous_sms_expiration",
			},
			setConfigSchema: setConfigSchema,
		}
		var testRestHandler RESTHandler = handler
		testRestHandler.RegisterPublic(r)

		testErr := errors.New("some error")
		target := "/api/authentication/v1/config/remember_for,remember_visible,anonymous_sms_expiration"
		introspectInfo := interfaces.TokenIntrospectInfo{Active: true, VisitorID: "user_id"}

		Convey("introspect failed", func() {
			introspectInfo.Active = false
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields传除remember_for、remember_visible、anonymous_sms_expiration之外的参数", func() {
			target += ",remember"
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("body does not have remember_for、remember_visible", func() {
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		reqParam := map[string]interface{}{
			"remember_for":             600,
			"remember_visible":         false,
			"anonymous_sms_expiration": 2.0,
		}

		cfg := interfaces.Config{
			RememberFor:     600,
			RememberVisible: false,
			SMSExpiration:   2,
		}

		keys := map[interfaces.ConfigKey]bool{
			interfaces.RememberFor:     true,
			interfaces.RememberVisible: true,
			interfaces.SMSExpiration:   true,
		}

		Convey("字段类型不正确", func() {
			Convey("remember_for类型错误", func() {
				url := "/api/authentication/v1/config/remember_for"
				reqParam["remember_for"] = "600"
				reqParamByte, _ := jsoniter.Marshal(reqParam)
				hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
				req := httptest.NewRequest("PUT", url, bytes.NewReader(reqParamByte))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				result := w.Result()

				assert.Equal(t, result.StatusCode, http.StatusBadRequest)

				if err := result.Body.Close(); err != nil {
					assert.Equal(t, err, nil)
				}
			})

			Convey("remember_visible类型错误", func() {
				url := "/api/authentication/v1/config/remember_visible"
				reqParam["remember_visible"] = 100
				reqParamByte, _ := jsoniter.Marshal(reqParam)
				hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
				req := httptest.NewRequest("PUT", url, bytes.NewReader(reqParamByte))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				result := w.Result()

				assert.Equal(t, result.StatusCode, http.StatusBadRequest)

				if err := result.Body.Close(); err != nil {
					assert.Equal(t, err, nil)
				}
			})

			Convey("anonymous_sms_expiration类型错误", func() {
				url := "/api/authentication/v1/config/anonymous_sms_expiration"
				reqParam["anonymous_sms_expiration"] = 1.1
				reqParamByte, _ := jsoniter.Marshal(reqParam)
				hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
				req := httptest.NewRequest("PUT", url, bytes.NewReader(reqParamByte))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				result := w.Result()

				assert.Equal(t, result.StatusCode, http.StatusBadRequest)

				if err := result.Body.Close(); err != nil {
					assert.Equal(t, err, nil)
				}
			})
		})

		Convey("SetConfig error", func() {
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			re.EXPECT().SetConfig(gomock.Any(), gomock.Any(), keys, cfg).Return(testErr)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("更新配置成功", func() {
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			re.EXPECT().SetConfig(gomock.Any(), gomock.Any(), keys, cfg).Return(nil)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("remember_for大于0,为9223372036854775295", func() {
			reqParam["remember_for"] = 9223372036854775295
			cfg.RememberFor = 9223372036854774784
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydra.EXPECT().Introspect(gomock.Any()).Return(introspectInfo, nil)
			re.EXPECT().SetConfig(gomock.Any(), gomock.Any(), keys, cfg).Return(nil)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
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
