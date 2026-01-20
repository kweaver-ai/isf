package login

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	accessTokenSchema "Authentication/driveradapters/jsonschema/access_token_schema"
	authSchema "Authentication/driveradapters/jsonschema/auth_schema"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func newRESTHandler(tmpLogin interfaces.Login) RESTHandler {
	clientLoginSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.ClientLoginSchemaStr))
	anonyousLogin2Schema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.AnonymousLogin2))
	pwdAuthSchemaStr, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(authSchema.PwdAuthSchemaStr))
	accessTokenSchema1, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(accessTokenSchema.AccessTokenSchemaStr))

	return &restHandler{
		login:                tmpLogin,
		clientLoginSchema:    clientLoginSchema,
		anonyousLogin2Schema: anonyousLogin2Schema,
		pwdAuthSchemaStr:     pwdAuthSchemaStr,
		accessTokenSchema:    accessTokenSchema1,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestSingleSignOn(t *testing.T) {
	Convey("single sign on", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("authentication")

		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPublic(r)

		target := "/api/authentication/v1/sso"

		Convey("client not exist", func() {
			reqInfo := map[string]interface{}{
				"redirect_uri":  "https://10.2.176.204:9010/callback",
				"response_type": "token id_token",
				"scope":         "offline openid",
				"credential": map[string]interface{}{
					"id": "test",
					"params": map[string]interface{}{
						"ticket": "ST-238-gAJhDR7DJekt6eoi5s0x-cas01.example.org",
					},
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
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

		Convey("singlesignon", func() {
			reqInfo := map[string]interface{}{
				"client_id":     "test",
				"redirect_uri":  "https://10.2.176.204:9010/callback",
				"response_type": "token id_token",
				"scope":         "offline openid",
				"credential": map[string]interface{}{
					"id": "test",
					"params": map[string]interface{}{
						"ticket": "ST-238-gAJhDR7DJekt6eoi5s0x-cas01.example.org",
					},
				},
			}
			tokenInfo := &interfaces.TokenInfo{}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			lo.EXPECT().SingleSignOn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Set("X-Forwarded-For", "127.0.0.1")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAnonymous(t *testing.T) {
	Convey("anonymous", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPublic(r)

		target := "/api/authentication/v1/anonymous"

		Convey("client not exist", func() {
			reqInfo := map[string]interface{}{
				"redirect_uri":  "https://10.2.176.204:9010/callback",
				"response_type": "token id_token",
				"scope":         "offline openid",
				"credential": map[string]interface{}{
					"account":  "test",
					"password": "test",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
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

		Convey("anonymous", func() {
			reqInfo := map[string]interface{}{
				"client_id":     "test",
				"redirect_uri":  "https://10.2.176.204:9010/callback",
				"response_type": "token",
				"scope":         "all",
				"credential": map[string]interface{}{
					"account":  "test",
					"password": "test",
				},
			}
			tokenInfo := &interfaces.TokenInfo{}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			lo.EXPECT().Anonymous(gomock.Any(), gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAnonymous2(t *testing.T) {
	Convey("anonymous2", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")
		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPublic(r)

		target := "/api/authentication/v2/anonymous"

		Convey("http basic auth info not exist", func() {
			reqInfo := map[string]interface{}{
				"credential": map[string]interface{}{
					"account":  "test",
					"password": "test",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("anonymou2 failed", func() {
			reqInfo := map[string]interface{}{
				"credential": map[string]interface{}{
					"account":  "test",
					"password": "test",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)
			resInfo := &interfaces.TokenInfo{}
			tmpErr := errors.New("xx")
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().Anonymous2(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(resInfo, tmpErr)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("anonymou2 success", func() {
			reqInfo := map[string]interface{}{
				"credential": map[string]interface{}{
					"account":  "test",
					"password": "test",
				},
				"vcode": map[string]interface{}{
					"id":      "01HWH7HH48SZRT35WNX8ZYX20R",
					"content": "123456",
				},
				"visitor_name": "xx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)
			resInfo := &interfaces.TokenInfo{
				AccessToken: "AccessTokenStr",
				TokenType:   "bearer",
				Scope:       "all",
				ExpirsesIn:  3600,
			}
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().Anonymous2(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(resInfo, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resBodyByte, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(resBodyByte, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody["access_token"].(string), "AccessTokenStr")
			assert.Equal(t, resBody["token_type"].(string), "bearer")
			assert.Equal(t, resBody["scope"].(string), "all")
			assert.Equal(t, int64(resBody["expires_in"].(float64)), int64(3600))
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestClientAccountAuth(t *testing.T) {
	Convey("ClientAccountAuth", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPrivate(r)

		target := "/api/authentication/v1/client-account-auth"
		userID := "user_id"
		Convey("invalid method", func() {
			reqParam := map[string]interface{}{
				"method":   "POST",
				"account":  "account1",
				"password": "password1",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid account", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "",
				"password": "password1",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid password", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid clientLoginOption1", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
				"option": map[string]interface{}{
					"vcodeType": 1,
					"vcode":     "",
					"uuid":      "uuid1",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid clientLoginOption2", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
				"option": map[string]interface{}{
					"vcodeType": 2,
					"vcode":     "vcodexxx",
					"uuid":      "uuid1",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid clientLoginOption3", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
				"option": map[string]interface{}{
					"vcodeType": 3,
					"vcode":     "",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid clientLoginOption4", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
				"option": map[string]interface{}{
					"vcodeType": 4,
					"vcode":     "",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("success1", func() {
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			respBodyByte, _ := io.ReadAll(result.Body)
			respBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBodyByte, &respBody)
			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, userID, respBody["user_id"].(string))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("success2", func() {
			userID := "user_id"
			reqParam := map[string]interface{}{
				"method":   "GET",
				"account":  "account1",
				"password": "password1",
				"option": map[string]interface{}{
					"vcodeType": 1,
					"vcode":     "vcodexxx",
					"uuid":      "uuid1",
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			lo.EXPECT().ClientAccountAuth(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userID, nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			respBodyByte, _ := io.ReadAll(result.Body)
			respBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBodyByte, &respBody)
			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, userID, respBody["user_id"].(string))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:nolintlint, lll
func TestPwdAuth(t *testing.T) {
	Convey("User Auth", t, func() {
		common.InitARTrace("Authentication-test")
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPublic(r)

		common.InitARTrace("test")
		target := "/api/authentication/v1/pwd-auth"
		Convey("account can not be empty", func() {
			reqInfo := map[string]interface{}{
				"account":  "",
				"password": "xx",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPError("account: String length must be greater than or equal to 1", rest.BadRequest, nil))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("password can not be empty", func() {
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPError("password: String length must be greater than or equal to 1", rest.BadRequest, nil))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("http basic auth param missing", func() {
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "UVtjEgp6dCM8UAss2M3P+fe3zy6qzwkcjZh5v7E5GlCq6Yd7a2CHENSzQ7p7UhZdpcFzy/XiD9MFv0t/n+LjnsewO5S++7y6osRACnzmGrDxrDy4Ypco0VkWsQT9QUpNlG2XW9YQUcZn9MZ4OVb24H6LhpZ3XxXbYYK+S8PkH0ANQql58QQvg1D0zPpcqa7AfhQhn1qpliaXV5w1EmUxd4Sc++mTt+zxgmCwPCI8PsAT6x1hchsBmoVhFRVBh6P76YLvjluRyUh0IMMzDXroVuljg9+M8ATmWsW79Ir2nBku40QR2HnXug6ycG9CuhOz+9XUG46zgP8+uvWj0LhY1g==",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("http basic auth param can not be empty", func() {
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "UVtjEgp6dCM8UAss2M3P+fe3zy6qzwkcjZh5v7E5GlCq6Yd7a2CHENSzQ7p7UhZdpcFzy/XiD9MFv0t/n+LjnsewO5S++7y6osRACnzmGrDxrDy4Ypco0VkWsQT9QUpNlG2XW9YQUcZn9MZ4OVb24H6LhpZ3XxXbYYK+S8PkH0ANQql58QQvg1D0zPpcqa7AfhQhn1qpliaXV5w1EmUxd4Sc++mTt+zxgmCwPCI8PsAT6x1hchsBmoVhFRVBh6P76YLvjluRyUh0IMMzDXroVuljg9+M8ATmWsW79Ir2nBku40QR2HnXug6ycG9CuhOz+9XUG46zgP8+uvWj0LhY1g==",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":")))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("user auth failed", func() {
			tmpErr := errors.New("unknow error")
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "UVtjEgp6dCM8UAss2M3P+fe3zy6qzwkcjZh5v7E5GlCq6Yd7a2CHENSzQ7p7UhZdpcFzy/XiD9MFv0t/n+LjnsewO5S++7y6osRACnzmGrDxrDy4Ypco0VkWsQT9QUpNlG2XW9YQUcZn9MZ4OVb24H6LhpZ3XxXbYYK+S8PkH0ANQql58QQvg1D0zPpcqa7AfhQhn1qpliaXV5w1EmUxd4Sc++mTt+zxgmCwPCI8PsAT6x1hchsBmoVhFRVBh6P76YLvjluRyUh0IMMzDXroVuljg9+M8ATmWsW79Ir2nBku40QR2HnXug6ycG9CuhOz+9XUG46zgP8+uvWj0LhY1g==",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().PwdAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, tmpErr)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("user auth success", func() {
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "UVtjEgp6dCM8UAss2M3P+fe3zy6qzwkcjZh5v7E5GlCq6Yd7a2CHENSzQ7p7UhZdpcFzy/XiD9MFv0t/n+LjnsewO5S++7y6osRACnzmGrDxrDy4Ypco0VkWsQT9QUpNlG2XW9YQUcZn9MZ4OVb24H6LhpZ3XxXbYYK+S8PkH0ANQql58QQvg1D0zPpcqa7AfhQhn1qpliaXV5w1EmUxd4Sc++mTt+zxgmCwPCI8PsAT6x1hchsBmoVhFRVBh6P76YLvjluRyUh0IMMzDXroVuljg9+M8ATmWsW79Ir2nBku40QR2HnXug6ycG9CuhOz+9XUG46zgP8+uvWj0LhY1g==",
			}
			tokenInfo := &interfaces.TokenInfo{
				AccessToken: "ory_at_xx",
				ExpirsesIn:  3600,
				Scope:       "all",
				TokenType:   "bearer",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().PwdAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(tokenInfo, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resBodyByte, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(resBodyByte, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody["access_token"].(string), tokenInfo.AccessToken)
			assert.Equal(t, int64(resBody["expires_in"].(float64)), tokenInfo.ExpirsesIn)
			assert.Equal(t, resBody["scope"].(string), tokenInfo.Scope)
			assert.Equal(t, resBody["token_type"].(string), tokenInfo.TokenType)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:nolintlint, lll
func TestGetAccessToken(t *testing.T) {
	Convey("Get AccessToken", t, func() {
		common.InitARTrace("Authentication-test")
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lo := mock.NewMockLogin(ctrl)
		testRestHandler := newRESTHandler(lo)
		testRestHandler.RegisterPublic(r)

		common.InitARTrace("test")
		target := "/api/authentication/v1/access_token"
		Convey("account can not be empty", func() {
			reqInfo := map[string]interface{}{
				"account": "",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", rest.InternalServerError, nil)
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPError("account: String length must be greater than or equal to 1", rest.BadRequest, nil))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("http basic auth param missing", func() {
			reqInfo := map[string]interface{}{
				"account":  "xx",
				"password": "UVtjEgp6dCM8UAss2M3P+fe3zy6qzwkcjZh5v7E5GlCq6Yd7a2CHENSzQ7p7UhZdpcFzy/XiD9MFv0t/n+LjnsewO5S++7y6osRACnzmGrDxrDy4Ypco0VkWsQT9QUpNlG2XW9YQUcZn9MZ4OVb24H6LhpZ3XxXbYYK+S8PkH0ANQql58QQvg1D0zPpcqa7AfhQhn1qpliaXV5w1EmUxd4Sc++mTt+zxgmCwPCI8PsAT6x1hchsBmoVhFRVBh6P76YLvjluRyUh0IMMzDXroVuljg9+M8ATmWsW79Ir2nBku40QR2HnXug6ycG9CuhOz+9XUG46zgP8+uvWj0LhY1g==",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("http basic auth param can not be empty", func() {
			reqInfo := map[string]interface{}{
				"account": "xx",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":")))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam, rest.NewHTTPErrorV2(rest.BadRequest, "http basic auth param missing or can not be empty"))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("get access_token failed", func() {
			tmpErr := errors.New("unknow error")
			reqInfo := map[string]interface{}{
				"account": "xx",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().GetAccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, tmpErr)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("get access_token success", func() {
			reqInfo := map[string]interface{}{
				"account": "xx",
			}
			tokenInfo := &interfaces.TokenInfo{
				AccessToken: "ory_at_xx",
				ExpirsesIn:  3600,
				Scope:       "all",
				TokenType:   "bearer",
			}
			reqParamByte, _ := json.Marshal(reqInfo)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("clientID"+":"+"clientSecret")))
			lo.EXPECT().GetAccessToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(tokenInfo, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resBodyByte, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(resBodyByte, &resBody)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, resBody["access_token"].(string), tokenInfo.AccessToken)
			assert.Equal(t, int64(resBody["expires_in"].(float64)), tokenInfo.ExpirsesIn)
			assert.Equal(t, resBody["scope"].(string), tokenInfo.Scope)
			assert.Equal(t, resBody["token_type"].(string), tokenInfo.TokenType)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
