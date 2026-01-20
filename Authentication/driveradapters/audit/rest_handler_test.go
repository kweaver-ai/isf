package audit

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"

	auditschema "Authentication/driveradapters/jsonschema/audit_schema"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func newRESTHandler(tmp interfaces.LogicsAudit, auditLogAsyncTask interfaces.LogicsAuditLogAsyncTask) RESTHandler {
	logSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditschema.AuditLogSchema))
	unorderedLogSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditschema.UnorderedAuditLogSchema))
	return &restHandler{
		audit:              tmp,
		logSchema:          logSchema,
		auditLogAsyncTask:  auditLogAsyncTask,
		unorderedLogSchema: unorderedLogSchema,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestLog(t *testing.T) {
	Convey("private Log", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		audit := mock.NewMockLogicsAudit(ctrl)
		testRestHandler := newRESTHandler(audit, nil)
		testRestHandler.RegisterPrivate(r)

		target := "/api/authentication/v1/audit-log"

		msg := map[string]interface{}{
			"xx": 1,
		}
		Convey("topic is required", func() {
			reqInfo := map[string]interface{}{
				"message": msg,
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam.Cause, "(root): topic is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("topic is string", func() {
			reqInfo := map[string]interface{}{
				"message": msg,
				"topic":   1,
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam.Cause, "topic: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("message is required", func() {
			reqInfo := map[string]interface{}{
				"topic": "xxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam.Cause, "(root): message is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("message is object", func() {
			reqInfo := map[string]interface{}{
				"message": "",
				"topic":   "x",
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			respParam := rest.NewHTTPErrorV2(rest.InternalServerError, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam.Cause, "message: Invalid type. Expected: object, given: string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			reqInfo := map[string]interface{}{
				"message": msg,
				"topic":   "x",
			}
			reqParamByte, _ := jsoniter.Marshal(reqInfo)
			audit.EXPECT().Log(gomock.Any(), gomock.Any())

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
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

func TestAuditLog(t *testing.T) {
	Convey("private unordered log", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		audit := mock.NewMockLogicsAudit(ctrl)
		auditLog := mock.NewMockLogicsAuditLogAsyncTask(ctrl)
		testRestHandler := newRESTHandler(audit, auditLog)
		testRestHandler.RegisterPrivate(r)

		target := "/api/authentication/v2/audit-log"

		Convey("message is required", func() {
			reqParamByte, _ := jsoniter.Marshal(map[string]interface{}{
				"topic": "",
			})
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("invalid topic", func() {
			reqParamByte, _ := jsoniter.Marshal(map[string]interface{}{
				"topic":   "",
				"message": map[string]interface{}{},
			})
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("Log error", func() {
			reqParamByte, _ := jsoniter.Marshal(map[string]interface{}{
				"topic":   "as.audit_log.log_operation",
				"message": map[string]interface{}{},
			})
			auditLog.EXPECT().Log(gomock.Any(), gomock.Any()).Return(errors.New("Log error"))
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("success", func() {
			reqParamByte, _ := jsoniter.Marshal(map[string]interface{}{
				"topic":   "as.audit_log.log_operation",
				"message": map[string]interface{}{},
			})
			auditLog.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest("POST", target, bytes.NewBuffer(reqParamByte))
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
