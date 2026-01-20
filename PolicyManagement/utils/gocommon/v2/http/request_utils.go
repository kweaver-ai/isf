package http

import (
	stderrors "errors"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"policy_mgnt/utils/gocommon/v2/errors"
)

// ParseBody 解析body中的json数据, 并绑定到结构体上面, 出错会返回400 TODO: 不应该是Public应该是对应服务码
func ParseBody(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(&obj); err != nil {
		return errors.ErrBadRequestPublic(&errors.ErrorInfo{Cause: err.Error(), Detail: err})
	}
	return nil
}

// ErrorResponse 处理error并生成 依次生成Response
func ErrorResponse(c *gin.Context, err error) {
	originalErr := err
	for stderrors.Unwrap(err) != nil {
		err = stderrors.Unwrap(err)
	}
	apiErr, ok := err.(*errors.Error)
	if !ok {
		err = originalErr
		apiErr = errors.ErrInternalServerErrorPublic(&errors.ErrorInfo{Cause: err.Error(), Detail: err}) // TODO: err 应该放在哪里
	}
	// 其实逻辑层Handler调用c.JSON(code, jsonObj)即可，
	// 因为最后一个Handler里面调用c.Abort()是没有意义的

	var detail map[string]any
	detail, ok = apiErr.Detail.(map[string]any)
	if !ok {
		detail = map[string]any{
			"detail": apiErr.Detail,
		}
	}

	restErr := &rest.HTTPError{
		Code:    apiErr.Code,
		Message: apiErr.Message,
		Detail:  detail,
		Cause:   apiErr.Cause,
	}

	rest.ReplyError(c, restErr)

	c.Abort()
}
