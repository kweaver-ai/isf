package api

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

// AssertString 断言是否字符串，不是返回空字符串
func AssertString(s interface{}) (r string) {
	if str, ok := s.(string); ok {
		r = str
	}
	return
}

// ErrorResponse 返回错误
func ErrorResponse(c *gin.Context, err error) {
	apiErr, ok := err.(*Error)
	if !ok {
		apiErr = ErrInternalServerErrorPublic(&ErrorInfo{Cause: err.Error()})
	}

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

// ValidJson 校验json数据，返回错误字段、错误原因
func ValidJson(shema, doc string) (invalide_params []string, cause string) {
	jsonLoader := gojsonschema.NewStringLoader(shema)
	schema, err := gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		return []string{"input json shema"}, err.Error()
	}

	documentLoader := gojsonschema.NewStringLoader(doc)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return []string{"input json data"}, err.Error()
	}

	if result.Valid() {
		return
	}
	for _, desc := range result.Errors() {
		invalide_params = append(invalide_params, desc.Field())
		if cause == "" {
			cause = desc.String()
			continue
		}
		cause = strings.Join([]string{cause, desc.String()}, "; ")
	}
	return
}

// ParseBody 解析body中的json数据，如果数据类型不符，报错
// 错误信息中包含解析失败的field
func ParseBody(c *gin.Context, params interface{}) (err error) {
	err = c.ShouldBindJSON(&params)
	if err != nil {
		var invalideParams []string
		var field string
		if err1, ok := err.(*json.UnmarshalTypeError); ok {
			field = err1.Field
		}
		invalideParams = append(invalideParams, field)
		err = ErrBadRequestPublic(&ErrorInfo{Detail: map[string]interface{}{"invalid_params": invalideParams}})
		return
	}
	return
}

// SliceFind 判断字符串是否存在数组和切片中
func SliceFind(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}

// GetEnv 封装os.Getenv(),可以指定默认值
func GetEnv(key, defaultV string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultV
	}
	return v
}
