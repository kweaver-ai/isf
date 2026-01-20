package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
)

func ValidateMiddleware(schema string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解决读取两次body失败
		var bodyBytes []byte
		if c.Request.Body == nil {
			common.ErrorResponse(c, errors.New(c.GetHeader("x-language"), errors.BadRequestErr, "Invalid Params", map[string]interface{}{"invalid_params": []string{"request body"}}))

			return
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			common.ErrorResponse(c, errors.New(c.GetHeader("x-language"), errors.BadRequestErr, "Invalid Params", map[string]interface{}{"invalid_params": []string{"request body"}}))
			return
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			common.ErrorResponse(c, errors.New(c.GetHeader("x-language"), errors.BadRequestErr, "Request body is needed.", map[string]interface{}{"invalid_params": []string{"request body"}}))
			return
		}
		invalideParams, cause := api.ValidJson(schema, string(bodyBytes))
		// 无错误字段，返回
		if len(invalideParams) == 0 {
			return
		}
		common.ErrorResponse(c, errors.New(c.GetHeader("x-language"), errors.BadRequestErr, cause, map[string]interface{}{"invalid_params": invalideParams}))
	}
}
