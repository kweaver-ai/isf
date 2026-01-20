package driveradapters

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"UserManagement/interfaces/mock"
)

func TestUpdateReservedName(t *testing.T) {
	Convey("添加/更新保留名称", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reservedName := mock.NewMockLogicsReservedName(ctrl)
		updateReservedNameSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateReservedNameSchemaStr))
		assert.Equal(t, err, nil)
		handler := &reservedNameHandler{
			reservedName:             reservedName,
			updateReservedNameSchema: updateReservedNameSchema,
		}

		handler.RegisterPrivate(r)
		target := "/api/user-management/v1/reserved-names/12341234123412341234123412341234"
		payload := map[string]interface{}{
			"name": "test",
		}

		Convey("name长度非法", func() {
			payload["name"] = "123412341234123412341234123412341234123412341234123412341234123412341234123412341234123412341234123412341234123412341234123412341"
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)
		})

		Convey("name开头包含空格", func() {
			payload["name"] = " test"
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)
		})

		Convey("name结尾包含空格", func() {
			payload["name"] = "test "
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)
		})

		Convey("name包含特殊字符", func() {
			payload["name"] = "test|123"
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test\123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test/123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test*123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test?123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test"123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test<123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)

			payload["name"] = `test>123`
			reqParamByte, _ = jsoniter.Marshal(payload)
			req = httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusBadRequest)
		})

		Convey("添加失败", func() {
			reservedName.EXPECT().UpdateReservedName(gomock.Any()).Return(errors.New("test"))
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusInternalServerError)
		})

		Convey("添加成功", func() {
			reservedName.EXPECT().UpdateReservedName(gomock.Any()).Return(nil)
			reqParamByte, _ := jsoniter.Marshal(payload)
			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusNoContent)
		})
	})
}

func TestDeleteReservedName(t *testing.T) {
	Convey("删除保留名称", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reservedName := mock.NewMockLogicsReservedName(ctrl)
		handler := &reservedNameHandler{
			reservedName: reservedName,
		}

		handler.RegisterPrivate(r)
		target := "/api/user-management/v1/reserved-names/12341234123412341234123412341234"

		Convey("删除失败", func() {
			reservedName.EXPECT().DeleteReservedName(gomock.Any()).Return(errors.New("test"))
			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusInternalServerError)
		})

		Convey("删除成功", func() {
			reservedName.EXPECT().DeleteReservedName(gomock.Any()).Return(nil)
			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, http.StatusNoContent)
		})
	})
}
