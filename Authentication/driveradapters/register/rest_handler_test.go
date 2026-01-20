package register

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	registerschema "Authentication/driveradapters/jsonschema/register_schema"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

//nolint:unparam
func getAllJSONKeys(v any) []string {
	var keys []string

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return keys
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			jsonKey := strings.Split(jsonTag, ",")[0]
			if jsonKey != "-" {
				keys = append(keys, jsonKey)
			}
		}
	}

	return keys
}

func newRESTHandler(tmpRegister interfaces.Register) RESTHandler {
	registerSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(registerschema.RegisterSchema))
	return &restHandler{
		register: tmpRegister,
		scopeMember: map[string]bool{
			"offline": true,
			"openid":  true,
			"all":     true,
		},
		errInvalidParameter: &RFC6749Error{
			Name:        "invalid_request",
			Description: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
			Code:        http.StatusBadRequest,
		},
		errInternalServerError: &RFC6749Error{
			Name:        "internal_server_error",
			Description: "Internal Server Error",
			Code:        http.StatusInternalServerError,
		},
		registerSchema: registerSchema,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

//nolint:funlen
func TestPublicRegisterCheckRequired(t *testing.T) {
	Convey("public register", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		registerLogics := mock.NewMockRegister(ctrl)
		testRestHandler := newRESTHandler(registerLogics)
		testRestHandler.RegisterPublic(r)

		target := "/oauth2/clients"

		reqInfo := map[string]interface{}{
			"client_name":               "test",
			"grant_types":               []string{"authorization_code", "implicit", "refresh_token"},
			"redirect_uris":             []string{"https://10.2.176.204:9010/callback"},
			"response_types":            []string{"token", "token id_token", "code"},
			"scope":                     "offline openid all",
			"post_logout_redirect_uris": []string{"https://10.2.176.204:9010/successful-logout"},
			"metadata": map[string]interface{}{
				"device": map[string]interface{}{
					"client_type": "ios",
				},
			},
		}

		Convey("post_logout_redirect_uris is required", func() {
			delete(reqInfo, "post_logout_redirect_uris")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("client_name is required", func() {
			delete(reqInfo, "client_name")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("grant_types is required", func() {
			delete(reqInfo, "grant_types")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("redirect_uris is required", func() {
			delete(reqInfo, "redirect_uris")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("scope is required", func() {
			delete(reqInfo, "scope")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("response_types is required", func() {
			delete(reqInfo, "response_types")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("metadata is required", func() {
			delete(reqInfo, "metadata")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("client_type is required", func() {
			device := reqInfo["metadata"].(map[string]interface{})["device"].(map[string]interface{})
			delete(device, "client_type")
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			var actualBody map[string]any
			err := jsoniter.Unmarshal(respBody, &actualBody)
			assert.Equal(t, err, nil)
			for _, key := range getAllJSONKeys(RFC6749Error{}) {
				_, ok := actualBody[key]
				assert.Equal(t, ok, true)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestPublicRegisterCheckEmpty(t *testing.T) {
	Convey("public register", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		regis := mock.NewMockRegister(ctrl)
		testRestHandler := newRESTHandler(regis)
		testRestHandler.RegisterPublic(r)

		target := "/oauth2/clients"

		reqInfo := map[string]interface{}{
			"client_name":               "test",
			"grant_types":               []string{"authorization_code", "implicit", "refresh_token"},
			"redirect_uris":             []string{"https://10.2.176.204:9010/callback"},
			"response_types":            []string{"token", "token id_token", "code"},
			"scope":                     "offline openid all",
			"post_logout_redirect_uris": []string{"https://10.2.176.204:9010/successful-logout"},
			"metadata": map[string]interface{}{
				"device": map[string]interface{}{
					"client_type": "ios",
				},
			},
		}

		clientInfo := interfaces.ClientInfo{
			ClientID:     "",
			ClientSecret: "",
		}

		tempErr := rest.NewHTTPError("xxx is not empty", 400000000, nil)
		gomock.InOrder(
			regis.EXPECT().PublicRegister(gomock.Any()).AnyTimes().Return(clientInfo, tempErr),
		)

		Convey("post_logout_redirect_uris is not empty", func() {
			reqInfo["post_logout_redirect_uris"] = []string{""}

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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("client_name is not empty", func() {
			reqInfo["client_name"] = ""
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("grant_types is not empty", func() {
			reqInfo["grant_types"] = []string{""}
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("redirect_uris is not empty", func() {
			reqInfo["redirect_uris"] = []string{""}
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("scope is not empty", func() {
			reqInfo["scope"] = ""
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("response_types is not empty", func() {
			reqInfo["response_types"] = ""
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("metadata is not empty", func() {
			reqInfo["metadata"] = map[string]interface{}{}
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("client_type is not empty", func() {
			device := reqInfo["metadata"].(map[string]interface{})["device"].(map[string]interface{})
			delete(device, "client_type")
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestPublicRegisterCheckInvalid(t *testing.T) {
	Convey("public register", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		registerLogics := mock.NewMockRegister(ctrl)
		testRestHandler := newRESTHandler(registerLogics)
		testRestHandler.RegisterPublic(r)

		target := "/oauth2/clients"

		reqInfo := map[string]interface{}{
			"client_name":               "test",
			"grant_types":               []string{"authorization_code", "implicit", "refresh_token"},
			"redirect_uris":             []string{"https://10.2.176.204:9010/callback"},
			"response_types":            []string{"token", "token id_token", "code"},
			"scope":                     "offline openid all",
			"post_logout_redirect_uris": []string{"https://10.2.176.204:9010/successful-logout"},
			"metadata": map[string]interface{}{
				"device": map[string]interface{}{
					"client_type": "ios",
				},
			},
		}

		Convey("grant_types is invalid", func() {
			reqInfo["grant_types"] = []string{"auth"}
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("scope is invalid", func() {
			reqInfo["scope"] = "all"
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("response_types is invalid", func() {
			reqInfo["response_types"] = "token"
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("client_type is invalid", func() {
			device := reqInfo["metadata"].(map[string]interface{})["device"].(map[string]interface{})
			device["client_type"] = "test"
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

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
