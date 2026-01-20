// Package driveradapters 应用账户组织管理权限设置测试
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
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func newOrgPermRESTHandler(perm interfaces.LogicsOrgPerm) *orgPermHandler {
	updateSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateSchemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	return &orgPermHandler{
		orgPerm: perm,
		permOrgStrType: map[string]interfaces.OrgType{
			"user":       interfaces.User,
			"department": interfaces.Department,
			"group":      interfaces.Group,
		},
		orgPermStrType: map[string]interfaces.OrgPermValue{
			"read": interfaces.OPRead,
		},
		subTypeStrType: map[string]interfaces.VisitorType{
			"user": interfaces.RealName,
		},
		updateSchema: updateSchema,
	}
}

//nolint:funlen,dupl
func TestUpdateOrgPerm(t *testing.T) {
	Convey("更新账户权限信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockLogicsOrgPerm(ctrl)
		org := newOrgPermRESTHandler(o)

		common.InitARTrace("user-management")

		org.RegisterPrivate(r)

		const target = "/api/user-management/v1/org-perm"
		Convey("sub type 错误", func() {
			tempTarget := target + "/xxx/xxx/xxxx"

			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url subject_type error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url object 错误", func() {
			tempTarget := target + "/user/xxx/xxxx"

			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url object 不唯一", func() {
			tempTarget := target + "/user/xxx/user,user"

			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects are not uniqued")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body object error", func() {
			tempTarget := target + "/user/xxx/user"

			data1 := gin.H{
				"object": "xxx",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "0.object: 0.object must be one of the following: \"user\", \"department\", \"group\"")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body perms error", func() {
			tempTarget := target + "/user/xxx/user"

			data1 := gin.H{
				"object": "user",
				"perms": []string{
					"xxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "0.perms.0: 0.perms.0 must be one of the following: \"read\"")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body object 和url object 不一致", func() {
			tempTarget := target + "/user/xxx/user"

			data1 := gin.H{
				"object": "department",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url object not equal object object")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
		Convey("SetOrgPerm error", func() {
			tempTarget := target + "/user/xxx/user"

			data1 := gin.H{
				"object": "user",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})
			o.EXPECT().SetOrgPerm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			tempTarget := target + "/user/xxx/user"

			data1 := gin.H{
				"object": "user",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})
			o.EXPECT().SetOrgPerm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
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

func TestDeleteOrgPerm(t *testing.T) {
	Convey("删除账户权限信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockLogicsOrgPerm(ctrl)
		org := newOrgPermRESTHandler(o)

		common.InitARTrace("user-management")

		org.RegisterPrivate(r)

		const target = "/api/user-management/v1/org-perm"
		Convey("sub type 错误", func() {
			tempTarget := target + "/xxx/xxx/xxxx"

			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url subject_type error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url object 错误", func() {
			tempTarget := target + "/user/xxx/xxxx"

			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
		Convey("DeleteOrgPerm 报错", func() {
			tempTarget := target + "/user/xxx/user"

			o.EXPECT().DeleteOrgPerm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			tempTarget := target + "/user/xxx/user"

			o.EXPECT().DeleteOrgPerm(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
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
