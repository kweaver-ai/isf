// Package driveradapters AnyShare 公共接口处理层
package driveradapters

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

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

type convertNamesParam struct {
	Method       *string  `json:"method" binding:"required"`
	UserIDs      []string `json:"user_ids" binding:"required"`
	DeptIDs      []string `json:"department_ids" binding:"required"`
	GroupIDs     []string `json:"group_ids" binding:"required"`
	AppIDs       []string `json:"app_ids" binding:"required"`
	ContactorIDs []string `json:"contactor_ids" binding:"required"`
}

type GetUserAndDepartmentInRangeParam struct {
	Method  *string  `json:"method" binding:"required"`
	UserIDs []string `json:"user_ids" binding:"required"`
	DeptIDs []string `json:"department_ids" binding:"required"`
	Range   []string `json:"scope" binding:"required"`
}

const (
	mstrGET                     = "GET"
	mstrGetUserAndDepInRangeURL = "/api/user-management/v1/search-org"
)

func newRESTHandler(combine interfaces.LogicsCombine, h interfaces.Hydra) RestHandler {
	return &restHandler{
		combine: combine,
		hydra:   h,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestGetHealth(t *testing.T) {
	Convey("getHealth", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPublic(r)

		req := httptest.NewRequest(mstrGET, "/health/ready", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}

func TestGetAlive(t *testing.T) {
	Convey("getAlive", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)

		req := httptest.NewRequest(mstrGET, "/health/alive", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		result := w.Result()

		assert.Equal(t, result.StatusCode, http.StatusOK)

		if err := result.Body.Close(); err != nil {
			assert.Equal(t, err, nil)
		}
	})
}

func TestConvertIDToNameFail(t *testing.T) {
	Convey("convertIDToName", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/names"

		method := "error"
		Convey("unsupported method", func() {
			reqParam := convertNamesParam{
				Method:  &method,
				UserIDs: []string{"user_id"},
				DeptIDs: []string{},
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

		method = mstrGET
		Convey("DeptIDs param no unique", func() {
			reqParam := convertNamesParam{
				Method:   &method,
				DeptIDs:  []string{"xxxxxx", "xxxxxx"},
				UserIDs:  []string{"xxxxxx"},
				GroupIDs: []string{"xxxxxx"},
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

		Convey("UserIDs param no unique", func() {
			reqParam := convertNamesParam{
				Method:   &method,
				DeptIDs:  []string{"xxxxxx"},
				UserIDs:  []string{"xxxxxx", "xxxxxx"},
				GroupIDs: []string{"xxxxxx"},
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

		Convey("Contacotrs param no unique", func() {
			reqParam := convertNamesParam{
				Method:       &method,
				DeptIDs:      []string{"xxxxxx"},
				ContactorIDs: []string{"xxxxxx", "xxxxxx"},
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

		Convey("appIDs param no unique", func() {
			reqParam := convertNamesParam{
				Method:   &method,
				DeptIDs:  []string{"xxxxxx"},
				UserIDs:  []string{"xxxxxx"},
				GroupIDs: []string{"xxxxxx"},
				AppIDs:   []string{"xxxxxx", "xxxxxx"},
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

		Convey("strict 参数不为bool", func() {
			reqParam := map[string]interface{}{
				"method":  &method,
				"app_ids": []string{"xxxxxx", "xxxxxx"},
				"strict":  "xxxxxx",
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
			assert.Equal(t, respParam.Cause, "type of body.strict should be bool")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestConvertIDToNameSuccess(t *testing.T) {
	Convey("convertIDToName", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/names"

		method := mstrGET

		Convey("GroupIDs param no unique", func() {
			reqParam := convertNamesParam{
				Method:   &method,
				DeptIDs:  []string{"xxxxxx"},
				GroupIDs: []string{"xxxxxx", "xxxxxx"},
				UserIDs:  []string{"xxxxxx"},
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
			assert.Equal(t, 1, 1)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			testName := interfaces.NameInfo{ID: "user_id", Name: "user_name"}
			testDepts := interfaces.NameInfo{ID: "department_id", Name: "department_name"}
			testGroups := interfaces.NameInfo{ID: "group_id", Name: "group_name"}
			testContactors := interfaces.NameInfo{ID: "contactor_id", Name: "contactor_name"}
			testApp := interfaces.NameInfo{ID: "app_id", Name: "app_name"}
			testNameInfo := []interfaces.NameInfo{testName}
			testDeptsInfo := []interfaces.NameInfo{testDepts}
			testGroupInfo := []interfaces.NameInfo{testGroups}
			testAppInfo := []interfaces.NameInfo{testApp}
			testContactorInfo := []interfaces.NameInfo{testContactors}
			info := interfaces.OrgNameInfo{
				UserNames:      testNameInfo,
				DepartNames:    testDeptsInfo,
				ContactorNames: testContactorInfo,
				GroupNames:     testGroupInfo,
				AppNames:       testAppInfo,
			}
			combineLogics.EXPECT().ConvertIDToName(gomock.Any(), gomock.Any(), gomock.Any(), true).Return(info, nil)

			reqParam := convertNamesParam{
				Method:       &method,
				DeptIDs:      []string{"xxxx"},
				UserIDs:      []string{"xxxxxx"},
				GroupIDs:     []string{"xxxxxxx"},
				AppIDs:       []string{"xxxxxxxx"},
				ContactorIDs: []string{"yyyy"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var respParam interface{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success1", func() {
			testName := interfaces.NameInfo{ID: "user_id", Name: "user_name"}
			testDepts := interfaces.NameInfo{ID: "department_id", Name: "department_name"}
			testGroups := interfaces.NameInfo{ID: "group_id", Name: "group_name"}
			testContactors := interfaces.NameInfo{ID: "contactor_id", Name: "contactor_name"}
			testApp := interfaces.NameInfo{ID: "app_id", Name: "app_name"}
			testNameInfo := []interfaces.NameInfo{testName}
			testDeptsInfo := []interfaces.NameInfo{testDepts}
			testGroupInfo := []interfaces.NameInfo{testGroups}
			testAppInfo := []interfaces.NameInfo{testApp}
			testContactorInfo := []interfaces.NameInfo{testContactors}
			info := interfaces.OrgNameInfo{
				UserNames:      testNameInfo,
				DepartNames:    testDeptsInfo,
				ContactorNames: testContactorInfo,
				GroupNames:     testGroupInfo,
				AppNames:       testAppInfo,
			}
			combineLogics.EXPECT().ConvertIDToName(gomock.Any(), gomock.Any(), gomock.Any(), false).Return(info, nil)

			reqParam := map[string]interface{}{
				"method":         &method,
				"department_ids": []string{"xxxx", "xxxx"},
				"user_ids":       []string{"xxxxxx", "xxxxxx"},
				"group_ids":      []string{"xxxxxxx", "xxxxxxx"},
				"app_ids":        []string{"xxxxxxxx", "xxxxxxxx"},
				"contactor_ids":  []string{"yyyy", "yyyy"},
				"strict":         false,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var respParam interface{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetUserAndDepartmentInRangeParams1(t *testing.T) {
	Convey("getUserAndDepartmentInRange", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := mstrGetUserAndDepInRangeURL

		method := ""
		Convey("unsupported method", func() {
			reqParam := GetUserAndDepartmentInRangeParam{
				Method:  &method,
				UserIDs: []string{"user_id"},
				DeptIDs: []string{},
				Range:   []string{},
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

		method = "xxxxx"
		Convey("unsupported method 1", func() {
			reqParam := GetUserAndDepartmentInRangeParam{
				Method:  &method,
				UserIDs: []string{"user_id"},
				DeptIDs: []string{},
				Range:   []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, 1, 1)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("unsupported method 2", func() {
			type GetUserAndDepartmentInRangeParamMethodInt struct {
				Method  *int     `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			tmp := 1
			reqParam := GetUserAndDepartmentInRangeParamMethodInt{
				Method:  &tmp,
				UserIDs: []string{"user_id"},
				DeptIDs: []string{},
				Range:   []string{},
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

		Convey("unsupported method 3", func() {
			type GetUserAndDepartmentInRangeParamMethodInt struct {
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}
			reqParam := GetUserAndDepartmentInRangeParamMethodInt{
				UserIDs: []string{"user_id"},
				DeptIDs: []string{},
				Range:   []string{},
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

		method = mstrGET
		Convey("unsupported UserIDs", func() {
			type GetUserAndDepartmentInRangeParamUserIDsInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs int      `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamUserIDsInt{
				Method:  &method,
				UserIDs: 1,
				DeptIDs: []string{},
				Range:   []string{},
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

func TestGetUserAndDepartmentInRangeParams2(t *testing.T) {
	Convey("getUserAndDepartmentInRange", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := mstrGetUserAndDepInRangeURL

		method := mstrGET

		Convey("unsupported UserIDs 1", func() {
			type GetUserAndDepartmentInRangeParamUserIDsArryOfInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []int    `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamUserIDsArryOfInt{
				Method:  &method,
				UserIDs: []int{1},
				DeptIDs: []string{},
				Range:   []string{},
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

		Convey("unsupported UserIDs 2", func() {
			type GetUserAndDepartmentInRangeParamUserIDsString struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs string   `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamUserIDsString{
				Method:  &method,
				UserIDs: method,
				DeptIDs: []string{},
				Range:   []string{},
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

		Convey("unsupported DeptIDs", func() {
			type GetUserAndDepartmentInRangeParamDeptIDsInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs int      `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamDeptIDsInt{
				Method:  &method,
				UserIDs: []string{},
				DeptIDs: 1,
				Range:   []string{},
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

		Convey("unsupported DeptIDs 1", func() {
			type GetUserAndDepartmentInRangeParamDeptIDsArryOfInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []int    `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamDeptIDsArryOfInt{
				Method:  &method,
				UserIDs: []string{},
				DeptIDs: []int{1},
				Range:   []string{},
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

		Convey("unsupported DeptIDs 2", func() {
			type GetUserAndDepartmentInRangeParamDeptIDsString struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs string   `json:"department_ids" `
				Range   []string `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamDeptIDsString{
				Method:  &method,
				DeptIDs: method,
				UserIDs: []string{},
				Range:   []string{},
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

func TestGetUserAndDepartmentInRangeParams3(t *testing.T) {
	Convey("getUserAndDepartmentInRange", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := mstrGetUserAndDepInRangeURL

		method := mstrGET

		Convey("unsupported range", func() {
			type GetUserAndDepartmentInRangeParamRangeInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   int      `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamRangeInt{
				Method:  &method,
				UserIDs: []string{},
				Range:   1,
				DeptIDs: []string{},
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

		Convey("unsupported range 1", func() {
			type GetUserAndDepartmentInRangeParamRangeArryOfInt struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   []int    `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamRangeArryOfInt{
				Method:  &method,
				UserIDs: []string{},
				Range:   []int{1},
				DeptIDs: []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, 1, 1)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("unsupported Range 2", func() {
			type GetUserAndDepartmentInRangeParamRangeString struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
				Range   string   `json:"scope" binding:"required"`
			}

			reqParam := GetUserAndDepartmentInRangeParamRangeString{
				Method:  &method,
				Range:   method,
				UserIDs: []string{},
				DeptIDs: []string{},
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

		Convey("unsupported Range 3", func() {
			type GetUserAndDepartmentInRangeParamRangeString struct {
				Method  *string  `json:"method" binding:"required"`
				UserIDs []string `json:"user_ids" `
				DeptIDs []string `json:"department_ids" `
			}

			reqParam := GetUserAndDepartmentInRangeParamRangeString{
				Method:  &method,
				UserIDs: []string{},
				DeptIDs: []string{},
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

func TestGetUserAndDepartmentInRange(t *testing.T) {
	Convey("getUserAndDepartmentInRange", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := mstrGetUserAndDepInRangeURL

		method := mstrGET

		Convey("GetUserAndDepartmentInScope fail", func() {
			testErr := rest.NewHTTPError("error", 401019006, nil)
			reqParam := GetUserAndDepartmentInRangeParam{
				Method:  &method,
				UserIDs: []string{"xxx"},
				Range:   []string{"xxx"},
				DeptIDs: []string{"xxx"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			combineLogics.EXPECT().GetUserAndDepartmentInScope(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, testErr)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 401019006)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			reqParam := GetUserAndDepartmentInRangeParam{
				Method:  &method,
				UserIDs: []string{},
				Range:   []string{},
				DeptIDs: []string{},
			}
			outDeps := []string{
				0: "xxxxx",
				1: "zzzz",
			}
			outUsers := []string{
				0: "yyy",
				1: "kkk",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			combineLogics.EXPECT().GetUserAndDepartmentInScope(gomock.Any(), gomock.Any(), gomock.Any()).Return(outUsers, outDeps, nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			type outData struct {
				UserIDs       []string `json:"user_ids" binding:"required"`
				DepartmentIDs []string `json:"department_ids" binding:"required"`
			}
			testData := outData{
				UserIDs:       outUsers,
				DepartmentIDs: outDeps,
			}
			var respParam outData
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, testData, respParam)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSearchInTree(t *testing.T) {
	Convey("搜索用户和部门信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("xxxxx")
		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/search-in-org-tree"

		Convey("token失效-报错", func() {
			tempTarget := target + "?type=user"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "department_id",
		}

		Convey("type参数错误-报错", func() {
			tempTarget := target + "?type=xxx&role=sys_admin&keyword=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset为-1-报错", func() {
			tempTarget := target + "?type=user&offset=-1&role=sys_admin&keyword=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("limit为0-报错", func() {
			tempTarget := target + "?type=user&limit=0&role=sys_admin&keyword=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("limit为1001-报错", func() {
			tempTarget := target + "?type=user&limit=1001&role=sys_admin&keyword=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("未传keyword报错", func() {
			tempTarget := target + "?type=user&limit=101&role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("用户不传role-报错", func() {
			tempTarget := target + "?type=user&limit=101&keyword=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("role不符合规则-报错", func() {
			tempTarget := target + "?type=user&limit=101&keyword=sys_admin&role=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type不符合规则-报错", func() {
			tempTarget := target + "?type=user11&limit=101&keyword=sys_admin&role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type不传-报错", func() {
			tempTarget := target + "?limit=101&keyword=sys_admin&role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("内部错误-报错", func() {
			tempTarget := target + "?type=user&role=sys_admin&keyword=xxx"
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			combineLogics.EXPECT().SearchInOrgTree(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil, 0, 0, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSearchInTreeSuccess(t *testing.T) {
	Convey("搜索用户和部门信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("xxxxx")
		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/search-in-org-tree"

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "department_id",
		}
		temp1 := interfaces.SearchUserInfo{
			ID:   "ID1",
			Name: "Name1",
			Type: "department",
		}

		temp2 := interfaces.SearchDepartInfo{
			ID:   "ID2",
			Name: "Name2",
			Type: "user",
		}

		userInfos := []interfaces.SearchUserInfo{temp1}
		departInfos := []interfaces.SearchDepartInfo{temp2}

		Convey("接口调用成功", func() {
			tempTarget := target + "?type=user&type=department&role=sys_admin&keyword=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			combineLogics.EXPECT().SearchInOrgTree(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfos, departInfos, 3, 6, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			respParam := make(map[string]ListInfo)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["users"].TotalCount, 3)
			outTemp := respParam["users"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp["id"], temp1.ID)
			assert.Equal(t, respParam["departments"].TotalCount, 6)
			outTemp1 := respParam["departments"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp1["id"], temp2.ID)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetEmails(t *testing.T) {
	Convey("获取邮箱信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/emails"

		temp := interfaces.OrgEmailInfo{}

		Convey("参数错误", func() {
			tempTarget := target

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("接口调用失败", func() {
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			tempTarget := target + "?user_id=xxxx&department_id=xxxx"
			combineLogics.EXPECT().GetEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(temp, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("接口调用成功", func() {
			tempTarget := target + "?user_id=xxxx&department_id=xxxx"
			combineLogics.EXPECT().GetEmails(gomock.Any(), gomock.Any()).AnyTimes().Return(temp, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
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

func TestConsoleSearchInTree(t *testing.T) {
	Convey("管理控制台-搜索用户和部门信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("xxxxx")
		combineLogics := mock.NewMockLogicsCombine(ctrl)
		h := mock.NewMockHydra(ctrl)
		testCRestHandler := newRESTHandler(combineLogics, h)
		testCRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/console/search-in-org-tree"

		Convey("token失效-报错", func() {
			tempTarget := target + "?type=user"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "department_id",
		}

		Convey("user_enabled 参数错误，报错", func() {
			tempTarget := target + "?type=user&role=sys_admin&user_enabled=xxxx&keyword=xx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBody, &resBody)
			assert.Equal(t, resBody["code"].(float64), float64(rest.BadRequest))
			assert.Equal(t, resBody["cause"], "invalid user_enabled")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("user_assigned 参数错误，报错", func() {
			tempTarget := target + "?type=user&role=sys_admin&user_assigned=xxxx&keyword=xx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			resBody := make(map[string]interface{})
			_ = jsoniter.Unmarshal(respBody, &resBody)
			assert.Equal(t, resBody["code"].(float64), float64(rest.BadRequest))
			assert.Equal(t, resBody["cause"], "invalid user_assigned")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("内部错误-报错", func() {
			tempTarget := target + "?type=user&role=sys_admin&keyword=xx"
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			combineLogics.EXPECT().SearchInOrgTree(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil, 0, 0, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		temp1 := interfaces.SearchUserInfo{
			ID:   "ID1",
			Name: "Name1",
			Type: "department",
		}

		temp2 := interfaces.SearchDepartInfo{
			ID:   "ID2",
			Name: "Name2",
			Type: "user",
		}

		userInfos := []interfaces.SearchUserInfo{temp1}
		departInfos := []interfaces.SearchDepartInfo{temp2}

		Convey("接口调用成功", func() {
			tempTarget := target + "?type=user&type=department&role=sys_admin&keyword=xxx&user_enabled=true&user_assigned=true"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			combineLogics.EXPECT().SearchInOrgTree(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfos, departInfos, 1, 1, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			respParam := make(map[string]ListInfo)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["departments"].TotalCount, 1)
			outTemp1 := respParam["departments"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp1["id"], temp2.ID)
			assert.Equal(t, respParam["users"].TotalCount, 1)
			outTemp := respParam["users"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp["id"], temp1.ID)
			assert.Equal(t, 1, 1)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
