package http

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"policy_mgnt/utils/gocommon/v2/errors"
	"policy_mgnt/utils/gocommon/v2/utils"
)

// JsonSchemaValidationMiddleware schema验证中间件
func JsonSchemaValidationMiddleware(jsonSchema string) gin.HandlerFunc {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(jsonSchema))
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		// 解决读取两次body失败
		var bodyBytes []byte
		if c.Request.Body == nil {
			ErrorResponse(c, errors.ErrBadRequestPublic(&errors.ErrorInfo{Detail: map[string]string{"invalid_params": "request body"}}))
			return
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			ErrorResponse(c, errors.ErrBadRequestPublic(&errors.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"request body"}}}))
			return
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if len(bodyBytes) == 0 {
			ErrorResponse(c, errors.ErrBadRequestPublic(&errors.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"request body"}}, Cause: "Request body is needed."}))
			return
		}
		err = utils.ValidateJsonByJsonSchema(schema, bodyBytes)
		// 无错误字段，返回
		if err != nil {
			ErrorResponse(c, errors.ErrBadRequestPublic(&errors.ErrorInfo{Cause: err.Error()}))
		}
	}
}
