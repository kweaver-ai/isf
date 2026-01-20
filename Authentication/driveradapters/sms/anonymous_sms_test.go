package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"

	"Authentication/common"
	smsSchema "Authentication/driveradapters/jsonschema/sms_schema"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"

	. "github.com/smartystreets/goconvey/convey"
)

func newRESTHandler(aSMS interfaces.LogicsAnonymousSMS) RESTHandler {
	aSMSSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(smsSchema.SMSSchemaStr))
	return &restHandler{
		aSMS:       aSMS,
		aSMSSchema: aSMSSchema,
	}
}

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestCreateAndSendAnonymousSMSCode(t *testing.T) {
	Convey("CreateAndSendAnonymousSMSCode", t, func() {
		test := setGinMode()
		defer test()

		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		aSMS := mock.NewMockLogicsAnonymousSMS(ctrl)

		common.InitARTrace("test")
		testRestHandler := newRESTHandler(aSMS)
		testRestHandler.RegisterPublic(r)
		target := "/api/authentication/v1/anonymous-sms-vcode"

		Convey("failed, phone_number is empty string", func() {
			reqBody := map[string]interface{}{
				"phone_number": "",
				"account":      "a-id",
			}
			reqBodyByte, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("failed, unknown error", func() {
			tmpErr := fmt.Errorf("unknown error")
			reqBody := map[string]interface{}{
				"phone_number": "tel_number",
				"account":      "a-id",
			}
			reqBodyByte, _ := json.Marshal(reqBody)
			aSMS.EXPECT().CreateAndSendVCode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", tmpErr)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			_, _ = io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("success", func() {
			reqBody := map[string]interface{}{
				"phone_number": "tel_number",
				"account":      "a-id",
			}
			reqBodyByte, _ := json.Marshal(reqBody)
			aSMS.EXPECT().CreateAndSendVCode(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("vcode-id", nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBodyByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			resParamByte, _ := io.ReadAll(result.Body)

			var jsonV map[string]interface{}
			err := jsoniter.Unmarshal(resParamByte, &jsonV)
			assert.Equal(t, err, nil)
			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, jsonV["vcode_id"].(string), "vcode-id")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
