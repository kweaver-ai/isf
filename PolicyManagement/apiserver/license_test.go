package apiserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/ory/gojsonschema"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"policy_mgnt/common"
	"policy_mgnt/interfaces"
	"policy_mgnt/interfaces/mock"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func newTestLicenseHandler() *licenseHandler {
	getAuthorizedProductsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getAuthorizedProductsSchemaString))
	if err != nil {
		panic(err)
	}

	updateAuthorizedProductsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateAuthorizedProductsSchemaString))
	if err != nil {
		panic(err)
	}
	return &licenseHandler{
		getAuthorizedProductsSchema: getAuthorizedProductsSchema,
		mapObjectTypeToString: map[interfaces.ObjectType]string{
			interfaces.ObjectTypeUser: "user",
		},
		updateAuthorizedProductsSchema: updateAuthorizedProductsSchema,
		mapStringToObjectType: map[string]interfaces.ObjectType{
			"user": interfaces.ObjectTypeUser,
		},
	}
}

func TestGetLicenses(t *testing.T) {
	Convey("getLicenses", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lic := mock.NewMockLogicsLicense(ctrl)
		h := mock.NewMockHydra(ctrl)

		group := r.Group("/api/license/v1")
		licenseHandler := newTestLicenseHandler()
		licenseHandler.license = lic
		licenseHandler.hydra = h
		licenseHandler.AddRouters(group)

		target := "/api/license/v1/console/licenses"

		Convey("token失效，报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: false}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		// GetLicenses 报错
		Convey("GetLicenses 报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().GetLicenses(gomock.Any(), gomock.Any()).Return(nil, gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable"))

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		// GetLicenses 成功
		Convey("GetLicenses 成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().GetLicenses(gomock.Any(), gomock.Any()).Return(map[string]interfaces.LicenseInfo{
				"product1": {Product: "product1", TotalUserQuota: 100, AuthorizedUserCount: 10},
				"product2": {Product: "product2", TotalUserQuota: 200, AuthorizedUserCount: 20},
			}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var out []interface{}
			err := jsoniter.Unmarshal(respBody, &out)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 2)

			if out[0].(map[string]interface{})["product"] == "product1" {
				assert.Equal(t, out[0].(map[string]interface{})["product"], "product1")
				assert.Equal(t, out[0].(map[string]interface{})["total_user_quota"].(float64), float64(100))
				assert.Equal(t, out[0].(map[string]interface{})["authorized_user_count"].(float64), float64(10))
			}
			if out[1].(map[string]interface{})["product"] == "product2" {
				assert.Equal(t, out[1].(map[string]interface{})["product"], "product2")
				assert.Equal(t, out[1].(map[string]interface{})["total_user_quota"].(float64), float64(200))
				assert.Equal(t, out[1].(map[string]interface{})["authorized_user_count"].(float64), float64(20))
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetAuthorizedProducts(t *testing.T) {
	Convey("getAuthorizedProducts", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lic := mock.NewMockLogicsLicense(ctrl)
		h := mock.NewMockHydra(ctrl)

		group := r.Group("/api/license/v1")
		licenseHandler := newTestLicenseHandler()
		licenseHandler.license = lic
		licenseHandler.hydra = h
		licenseHandler.AddRouters(group)

		target := "/api/license/v1/console/query-authorized-products"

		Convey("token失效，报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: false}, nil)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts 参数错误，没有method", func() {
			jsonReq := map[string]interface{}{
				"user_ids": []string{"1234567890"},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "method: method is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts 参数错误，method不为GET", func() {
			jsonReq := map[string]interface{}{
				"method":   "POST",
				"user_ids": []string{"1234567890"},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "method: method must be one of the following: \"GET\"")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts 参数错误，method不为string", func() {
			jsonReq := map[string]interface{}{
				"method":   123,
				"user_ids": []string{"1234567890"},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "method: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts 参数错误，user_ids不为[]string", func() {
			jsonReq := map[string]interface{}{
				"method":   "GET",
				"user_ids": 123,
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "user_ids: Invalid type. Expected: array, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts 参数错误，user_ids不存在", func() {
			jsonReq := map[string]interface{}{
				"method": "GET",
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "user_ids: user_ids is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAuthorizedProducts GetAuthorizedProducts 报错", func() {
			jsonReq := map[string]interface{}{
				"method":   "GET",
				"user_ids": []string{"1234567890"},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().GetAuthorizedProducts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable"))

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			jsonReq := map[string]interface{}{
				"method":   "GET",
				"user_ids": []string{"1", "2"},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().GetAuthorizedProducts(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string]interfaces.AuthorizedProduct{
				"1": {ID: "1", Type: interfaces.ObjectTypeUser, Product: []string{"product1", "product2"}},
				"2": {ID: "2", Type: interfaces.ObjectTypeUser, Product: []string{"product2"}},
			}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var out []interface{}
			err = jsoniter.Unmarshal(respBody, &out)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 2)

			for _, v := range out {
				temp := v.(map[string]interface{})

				if temp["id"].(string) == "1" {
					assert.Equal(t, temp["id"].(string), "1")
					assert.Equal(t, temp["type"].(string), "user")
					assert.Equal(t, temp["products"].([]interface{}), []interface{}{"product1", "product2"})
				} else if temp["id"].(string) == "2" {
					assert.Equal(t, temp["id"].(string), "2")
					assert.Equal(t, temp["type"].(string), "user")
					assert.Equal(t, temp["products"].([]interface{}), []interface{}{"product2"})
				}
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

	})
}

func TestCheckProductAuthorized(t *testing.T) {
	Convey("checkProductAuthorized", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lic := mock.NewMockLogicsLicense(ctrl)
		h := mock.NewMockHydra(ctrl)

		group := r.Group("/api/license/v1")
		licenseHandler := newTestLicenseHandler()
		licenseHandler.license = lic
		licenseHandler.hydra = h
		licenseHandler.AddRouters(group)

		target := "/api/license/v1/check-product-authorized"

		Convey("token失效，报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: false}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("CheckProductAuthorized 参数错误，没有product", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "product is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("CheckProductAuthorized报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().CheckProductAuthorized(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "", gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable"))

			req := httptest.NewRequest("GET", target+"?product=product1", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().CheckProductAuthorized(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, "xxxx", nil)

			req := httptest.NewRequest("GET", target+"?product=product1", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := map[string]interface{}{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam["authorized"].(bool), false)
			assert.Equal(t, respParam["unauthorized_reason"].(string), "xxxx")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUpdateAuthorizedProducts(t *testing.T) {
	Convey("updateAuthorizedProducts", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("policy-mgnt")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lic := mock.NewMockLogicsLicense(ctrl)
		h := mock.NewMockHydra(ctrl)

		group := r.Group("/api/license/v1")
		licenseHandler := newTestLicenseHandler()
		licenseHandler.license = lic
		licenseHandler.hydra = h
		licenseHandler.AddRouters(group)

		target := "/api/license/v1/console/authorized-products"

		Convey("token失效，报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: false}, nil)

			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，没有id", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"type":     "user",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.id: id is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，id不不为string", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       1111,
					"type":     "user",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.id: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，type不不为string", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"type":     123,
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.type: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，type不为user", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"type":     "user1",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.type: 0.type must be one of the following: \"user\"")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，type不存在", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.type: type is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，products不存在", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":   "1",
					"type": "user",
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.products: products is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts 参数错误，products不为[]string", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"type":     "user",
					"products": 123,
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "0.products: Invalid type. Expected: array, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateAuthorizedProducts报错", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"type":     "user",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().UpdateAuthorizedProducts(gomock.Any(), gomock.Any(), gomock.Any()).Return(gerrors.NewError(gerrors.PublicForbidden, "service unavailable"))

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := gerrors.NewError(gerrors.PublicServiceUnavailable, "service unavailable")
			err = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			jsonReq := []interface{}{
				map[string]interface{}{
					"id":       "1",
					"type":     "user",
					"products": []string{"product1"},
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			lic.EXPECT().UpdateAuthorizedProducts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			body, err := json.Marshal(jsonReq)
			assert.Equal(t, err, nil)
			req := httptest.NewRequest("PUT", target, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

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
