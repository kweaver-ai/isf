//nolint:govet
package driveradapters

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authorization/common"
)

func TestNewServiceConfigDriver(t *testing.T) {
	Convey("Test NewServiceConfigDriver", t, func() {
		Convey("Should create singleton instance", func() {
			driver1 := NewServiceConfigDriver()
			driver2 := NewServiceConfigDriver()

			So(driver1, ShouldNotBeNil)
			So(driver2, ShouldNotBeNil)
			So(driver1, ShouldEqual, driver2) // 单例模式测试
		})
	})
}

func TestServiceConfigDriver_RegisterPrivate(t *testing.T) {
	Convey("Test RegisterPrivate", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		driver := NewServiceConfigDriver()
		driver.RegisterPrivate(r)

		Convey("Should register GET /api/authorization/v1/config/log/level", func() {
			req := httptest.NewRequest("GET", "/api/authorization/v1/config/log/level", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should register PUT /api/authorization/v1/config/log/level", func() {
			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// PUT 请求需要 body，所以会返回 400 或 500，这是正常的
			So(result.StatusCode, ShouldBeGreaterThanOrEqualTo, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestServiceConfigDriver_RegisterPublic(t *testing.T) {
	Convey("Test RegisterPublic", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		driver := NewServiceConfigDriver()
		driver.RegisterPublic(r)

		Convey("Should not register any public endpoints", func() {
			req := httptest.NewRequest("GET", "/api/authorization/v1/config/log/level", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// 应该返回 404，因为没有注册任何路由
			assert.Equal(t, result.StatusCode, http.StatusNotFound)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestServiceConfigDriver_getLogLevel(t *testing.T) {
	Convey("Test getLogLevel", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		driver := NewServiceConfigDriver()
		driver.RegisterPrivate(r)

		Convey("Should return current log level", func() {
			req := httptest.NewRequest("GET", "/api/authorization/v1/config/log/level", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			body, err := io.ReadAll(result.Body)
			So(err, ShouldBeNil)

			var response map[string]any
			err = json.Unmarshal(body, &response)
			So(err, ShouldBeNil)

			So(response, ShouldContainKey, "level")
			So(response["level"], ShouldNotBeEmpty)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestServiceConfigDriver_setLogLevel(t *testing.T) {
	Convey("Test setLogLevel", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		driver := NewServiceConfigDriver()
		driver.RegisterPrivate(r)

		Convey("Should set log level successfully with valid level", func() {
			// 保存原始日志级别
			originalLevel := common.GetLogLevel()

			// 测试设置 debug 级别
			requestBody := map[string]string{"level": "debug"}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			// 验证日志级别是否被设置
			currentLevel := common.GetLogLevel()
			So(currentLevel, ShouldEqual, "debug")

			// 恢复原始日志级别
			_ = common.SetLogLevel(originalLevel)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should return error for invalid log level", func() {
			requestBody := map[string]string{"level": "invalid_level"}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// 应该返回 400 错误
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should return error for missing level key", func() {
			requestBody := map[string]string{"other_key": "debug"}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// 应该返回 400 错误
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should return error for invalid JSON body", func() {
			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBufferString("invalid json"))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// 应该返回 400 错误
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should return error for empty body", func() {
			req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", http.NoBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			// 应该返回 400 或 500 错误
			So(result.StatusCode, ShouldBeGreaterThanOrEqualTo, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Should handle different valid log levels", func() {
			// 保存原始日志级别
			originalLevel := common.GetLogLevel()

			validLevels := []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}

			for _, level := range validLevels {
				requestBody := map[string]string{"level": level}
				jsonBody, _ := json.Marshal(requestBody)

				req := httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				result := w.Result()

				assert.Equal(t, result.StatusCode, http.StatusNoContent)

				// 验证日志级别是否被设置
				currentLevel := common.GetLogLevel()
				So(currentLevel, ShouldEqual, level)

				if err := result.Body.Close(); err != nil {
					assert.Equal(t, err, nil)
				}
			}

			// 恢复原始日志级别
			_ = common.SetLogLevel(originalLevel)
		})
	})
}

func TestServiceConfigDriver_Integration(t *testing.T) {
	Convey("Integration test for service config driver", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		driver := NewServiceConfigDriver()
		driver.RegisterPrivate(r)

		Convey("Should be able to get and set log level", func() {
			// 保存原始日志级别
			originalLevel := common.GetLogLevel()

			// 1. 获取当前日志级别
			req := httptest.NewRequest("GET", "/api/authorization/v1/config/log/level", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			body, err := io.ReadAll(result.Body)
			So(err, ShouldBeNil)

			var getResponse map[string]any
			err = json.Unmarshal(body, &getResponse)
			So(err, ShouldBeNil)
			So(getResponse, ShouldContainKey, "level")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}

			// 2. 设置新的日志级别
			requestBody := map[string]string{"level": "info"}
			jsonBody, _ := json.Marshal(requestBody)

			req = httptest.NewRequest("PUT", "/api/authorization/v1/config/log/level", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result = w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}

			// 3. 验证日志级别已被设置
			req = httptest.NewRequest("GET", "/api/authorization/v1/config/log/level", http.NoBody)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result = w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			body, err = io.ReadAll(result.Body)
			So(err, ShouldBeNil)

			var setResponse map[string]any
			err = json.Unmarshal(body, &setResponse)
			So(err, ShouldBeNil)
			So(setResponse["level"], ShouldEqual, "info")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}

			// 恢复原始日志级别
			_ = common.SetLogLevel(originalLevel)
		})
	})
}
