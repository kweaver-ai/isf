// Package util 协议层
package util

import (
	"fmt"
	"io"
	"strings"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/interfaces"
)

var (
	errCodeTypeMap = map[string]interfaces.ErrorCodeType{
		"string": interfaces.Str,
		"number": interfaces.Number,
	}
)

// IsEmpty 检查map中所有字符串是否为空
func IsEmpty(strs map[string]string) error {
	for key, str := range strs {
		if str == "" {
			cause := fmt.Sprintf("The value of parameter %v is invalid", key)
			return rest.NewHTTPError(cause, rest.BadRequest, nil)
		}
	}
	return nil
}

func GetErrorCodeType(c *gin.Context) interfaces.ErrorCodeType {
	errCodeType := strings.ToLower(c.GetHeader("x-error-code"))
	if val, isExist := errCodeTypeMap[errCodeType]; isExist {
		return val
	}
	return interfaces.Number
}

// Verify 内省令牌
func Verify(c *gin.Context, hydra interfaces.Hydra) (visitor interfaces.Visitor, err error) {
	tokenID := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenID, "Bearer ")
	info, err := hydra.Introspect(token)
	if err != nil {
		return
	}

	if !info.Active {
		err = rest.NewHTTPError("token expired", rest.Unauthorized, nil)
		return
	}

	visitor = interfaces.Visitor{
		ID:            info.VisitorID,
		TokenID:       tokenID,
		IP:            c.ClientIP(),
		Mac:           c.GetHeader("X-Request-MAC"),
		UserAgent:     c.GetHeader("User-Agent"),
		Type:          info.VisitorTyp,
		ErrorCodeType: GetErrorCodeType(c),
	}

	return
}

// ValidateAndBindGin 校验json数据
func ValidateAndBindGin(c *gin.Context, schema *gojsonschema.Schema, bind interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return ValidateAndBind(body, schema, bind)
}

// ValidateAndBind 校验json数据
func ValidateAndBind(body []byte, schema *gojsonschema.Schema, bind interface{}) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return rest.NewHTTPError(err.Error(), rest.BadRequest, nil)
	}
	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return rest.NewHTTPError(strings.Join(msgList, "; "), rest.BadRequest, nil)
	}

	if err := jsoniter.Unmarshal(body, bind); err != nil {
		return err
	}

	return nil
}
