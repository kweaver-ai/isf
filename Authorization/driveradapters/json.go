/*
json 通用方法
- json schema 创建
- json schema 校验
*/
package driveradapters

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/xeipuuv/gojsonschema"

	"Authorization/common"
)

// ValidateAndBindGin 校验json数据
func validateAndBindGin(c *gin.Context, schema *gojsonschema.Schema, bind any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return validateAndBind(body, schema, bind)
}

// validateAndBind 校验json数据
func validateAndBind(body []byte, schema *gojsonschema.Schema, bind any) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return gerrors.NewError(gerrors.PublicBadRequest, err.Error())
	}
	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return gerrors.NewError(gerrors.PublicBadRequest, strings.Join(msgList, "; "))
	}

	if err := jsoniter.Unmarshal(body, bind); err != nil {
		return err
	}

	return nil
}

// newJSONSchema 创建一个 *gojsonschema.Schema, 失败程序退出
func newJSONSchema(schemaStr string) (schema *gojsonschema.Schema) {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaStr))
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	return
}
