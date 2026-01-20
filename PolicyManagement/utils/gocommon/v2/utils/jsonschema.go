package utils

import (
	"strings"

	"errors"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJsonByJsonSchema 校验json数据
func ValidateJsonByJsonSchema(schema *gojsonschema.Schema, jsonData []byte) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(jsonData))
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	} else {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.Description())
		}
		return errors.New(strings.Join(msgList, "; "))
	}
}
