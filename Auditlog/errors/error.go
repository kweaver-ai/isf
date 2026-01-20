/* 错误码构造函数 */
package errors

import (
	"context"
	"fmt"
	"strings"

	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/infra/cmp/langcmp"
)

// New 新建错误码返回体
func New(language string, code int, cause string, detail interface{}) *ErrorResp {
	var lang langcmp.Lang
	// 未设置语言，使用系统默认语言
	language = strings.ToLower(language)
	if language == "" {
		lang = langcmp.NewLangCmp().GetSysDefaultLang()
	} else {
		lang = langcmp.NewFromStr(language)
	}
	return &ErrorResp{
		code:        code,
		cause:       cause,
		detail:      detail,
		description: i18ns[code].Description[lang],
		solution:    i18ns[code].Solution[lang],
	}
}

func NewCtx(ctx context.Context, code int, cause string, detail interface{}) *ErrorResp {
	lang := helpers.GetLangFromCtx(ctx)
	return New(string(lang), code, cause, detail)
}

// 内部使用的错误响应结构
type ErrorResp struct {
	code        int         // 错误码 前三位等于状态码，中间三位为服务标识，后三位为错误标识
	cause       string      // 错误原因，产生错误的具体原因
	detail      interface{} // 错误码拓展信息，补充说明错误信息
	description string      // 错误描述，客户端采用此字段做错误提示（需要符合国际化要求）
	solution    string      // 操作提示，针对当前错误的操作提示（需要符合国际化要求）
}

// MarshalJSON .
func (e *ErrorResp) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["code"] = e.code
	result["cause"] = e.cause
	result["detail"] = e.detail
	result["description"] = e.description
	result["message"] = e.description
	result["solution"] = e.solution
	return utils.JSON().Marshal(&result)
}

// Error .
func (e *ErrorResp) Error() string {
	errInfo := []string{}

	if e.code != 0 {
		errInfo = append(errInfo, fmt.Sprintf("code: %d", e.code))
	}

	if e.description != "" {
		errInfo = append(errInfo, fmt.Sprintf("Description: %s", e.description))
	}

	if e.cause != "" {
		errInfo = append(errInfo, fmt.Sprintf("Cause: %s", e.cause))
	}

	if e.solution != "" {
		errInfo = append(errInfo, fmt.Sprintf("Solution: %s", e.solution))
	}

	return strings.Join(errInfo, ", ")
}

func (e *ErrorResp) Code() int {
	return e.code
}

func (e *ErrorResp) HTTPCode() int {
	// 前3位为http状态码
	return e.code / 1000000
}

func (e *ErrorResp) Cause() string {
	return e.cause
}

func (e *ErrorResp) Description() string {
	return e.description
}

func (e *ErrorResp) Detail() interface{} {
	return e.detail
}

func (e *ErrorResp) Solution() string {
	return e.solution
}
