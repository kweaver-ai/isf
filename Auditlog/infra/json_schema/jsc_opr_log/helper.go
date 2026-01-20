package jsc_opr_log

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/utils"
)

// ValidateOprLogJSONSchema 验证运营日志的json schema
func ValidateOprLogJSONSchema(logBys []byte, bizType oprlogenums.BizType) (operation string, invalidFields []string, err error) {
	key := OprJSONSchemaCheckKey(bizType)

	// 获取对应类型的校验路径
	schemas, exists := OprJsonSchemaCheckPathMap[key]

	operation = gjson.GetBytes(logBys, "0.operation").String()

	// 如果不存在对应的校验路径，则使用通用的校验路径
	if !exists {
		if bizType.IsServerBizType() {
			schemas = OprJsonSchemaCheckPathMap[OprJSCommonServer]

			// 特殊处理：bizType为ContentAutomation且operation为use_create时，使用客户端通用模式验证（为前端上报）
			if IsSpecialClientType(bizType, operation) {
				schemas = OprJsonSchemaCheckPathMap[OprJSCommonClient]
			}
		} else if bizType.IsClientBizType() {
			schemas = OprJsonSchemaCheckPathMap[OprJSCommonClient]
		} else {
			panic("[ValidateOprLogJSONSchema]: 未知的业务类型")
		}
	}

	// 依次使用路径中的每个模式验证日志数据
	for i, schema := range schemas {
		if bizType != oprlogenums.ClientOperation && i == 0 {
			arrStr := `["id","type"]`

			schema, err = utils.AddJSONArrayToJSON(schema, "properties.object.required", arrStr)
			if err != nil {
				err = errors.Wrap(err, "[ValidateOprLogJSONSchema]:AddToJSON")
				return
			}
		}

		schema = ToBatchSchema(schema)

		var _invalidFields []string

		_invalidFields, err = utils.ValidJsonSchema(schema, string(logBys))
		if err != nil {
			err = errors.Wrap(err, "[ValidateOprLogJSONSchema]:ValidJsonSchema")
			return
		}

		if len(_invalidFields) > 0 {
			invalidFields = append(invalidFields, _invalidFields...)
		}
	}

	invalidFields = utils.DeduplGeneric(invalidFields)

	return
}

type FieldInvalidS struct {
	InvalidField string `json:"invalid_field"`
	Cause        string `json:"cause"`
}

func ToBatchSchema(schema string) (newSchema string) {
	newSchema = fmt.Sprintf(`{
  "type": "array",
  "items": %s,
  "minItems": 1
}
`, schema)

	return
}

func IsSpecialClientType(bizType oprlogenums.BizType, operation string) bool {
	b1 := bizType == oprlogenums.ContentAutomation && operation == "use_create"

	return b1
}
