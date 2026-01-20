// Package driveradapters contactor AnyShare  部门逻辑接口处理层
package driveradapters

import (
	"bytes"
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

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestDeleteContactor(t *testing.T) {
	Convey("deleteContactor", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contactor := mock.NewMockLogicsContactor(ctrl)
		hydra := mock.NewMockHydra(ctrl)
		testCRestHandler := &contactorRestHandler{
			contactor: contactor,
			hydra:     hydra,
		}
		testCRestHandler.RegisterPublic(r)
		target := "/api/eacp/v1/contactor/deletegroup"

		tokenInfo := interfaces.TokenIntrospectInfo{
			Active: false,
		}
		Convey("token 检测报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, tempErr)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("token 检测失败", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tokenInfo.Active = true
		reqBody := []byte{'x', 'z'}
		Convey("request body 非json", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tmpBody := map[string]interface{}{
			"name": "client_1",
		}
		reqBody, _ = jsoniter.Marshal(tmpBody)
		Convey("缺少参数groupid", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tmpBody["groupid"] = "e1c0eed3-8ba7-4eec-9a1f-747591cc3661"
		reqBody, _ = jsoniter.Marshal(tmpBody)
		Convey("DeleteContactor 失败", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			contactor.EXPECT().DeleteContactor(gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			hydra.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(tokenInfo, nil)
			contactor.EXPECT().DeleteContactor(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
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

// 测试GetContactorMembers
func TestGetContactorMembers(t *testing.T) {
	Convey("GetContactorMembers", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		contactor := mock.NewMockLogicsContactor(ctrl)
		hydra := mock.NewMockHydra(ctrl)

		getContactorMembersPostSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getContactorMembersPostSchemaStr))
		assert.Equal(t, err, nil)
		testCRestHandler := &contactorRestHandler{
			contactor:                     contactor,
			hydra:                         hydra,
			getContactorMembersPostSchema: getContactorMembersPostSchema,
		}
		testCRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/contactor-members"

		Convey("参数检查，没有method字段", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["contactor_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("参数检查，没有contactor_ids字段", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("参数检查，method字段不为字符串", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = 11
			tmpBody["contactor_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("参数检查，contactor_ids字段不为数组", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["contactor_ids"] = "test"
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("method字段不为GET", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = strXXX
			tmpBody["contactor_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetContactorMembers 失败", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			contactor.EXPECT().GetContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, tempErr)

			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["contactor_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			outData := []interfaces.ContactorMemberInfo{
				{
					ContactorID: strXXX,
					MemberIDs:   []string{strXXX},
				},
			}
			contactor.EXPECT().GetContactorMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(outData, nil)

			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["contactor_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var tmpData []interface{}
			err := jsoniter.Unmarshal(w.Body.Bytes(), &tmpData)
			assert.Equal(t, err, nil)

			assert.Equal(t, len(tmpData), 1)

			assert.Equal(t, tmpData[0].(map[string]interface{})["contactor_id"], outData[0].ContactorID)

			data := tmpData[0].(map[string]interface{})["members"].([]interface{})
			assert.Equal(t, data[0].(map[string]interface{})["id"], outData[0].MemberIDs[0])
			assert.Equal(t, data[0].(map[string]interface{})["type"], "user")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
