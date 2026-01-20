// Package error
package error

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	prefix = "Public."
)

const (
	// BadRequest 客户端请求错误
	BadRequest = "BadRequest"
	// Unauthorized 未授权或者授权已过期
	Unauthorized = "Unauthorized"
	// Forbidden 禁止访问
	Forbidden = "Forbidden"
	// NotFound 请求URI资源不存在
	NotFound = "NotFound"
	// Conflict 资源冲突
	Conflict = "Conflict"
)

const (
	// InternalServerError 服务端内部错误
	InternalServerError = "InternalServerError"
	// ServiceUnavailable 服务器暂时不可用
	ServiceUnavailable = "ServiceUnavailable"
)

const (
	// PublicBadRequest 通用错误码，客户端请求错误
	PublicBadRequest = prefix + BadRequest
	// PublicUnauthorized 通用错误码，未授权或者授权已过期
	PublicUnauthorized = prefix + Unauthorized
	// PublicForbidden 通用错误码，禁止访问
	PublicForbidden = prefix + Forbidden
	// PublicNotFound 通用错误码，请求URI资源不存在
	PublicNotFound = prefix + NotFound
	// PublicConflict 通用错误码，资源冲突
	PublicConflict = prefix + Conflict
)

const (
	// PublicInternalServerError 通用错误码，服务端内部错误
	PublicInternalServerError = prefix + InternalServerError
	// PublicServiceUnavailable 通用错误码，服务器暂时不可用
	PublicServiceUnavailable = prefix + ServiceUnavailable
)

var (
	codesSet = map[string]struct{}{
		BadRequest:          struct{}{},
		Unauthorized:        struct{}{},
		Forbidden:           struct{}{},
		NotFound:            struct{}{},
		Conflict:            struct{}{},
		InternalServerError: struct{}{},
		ServiceUnavailable:  struct{}{},
	}
)

// Error 错误信息
type Error struct {
	Code        string
	Description string
	Solution    string
	Detail      map[string]interface{}
	Link        string
}

// SetErrAttribute 设置参数
type SetErrAttribute func(*Error)

// SetSolution 设置solution参数
func SetSolution(solution string) SetErrAttribute {
	return func(err *Error) {
		err.Solution = solution
	}
}

// SetDetail 设置detail参数
func SetDetail(detail map[string]interface{}) SetErrAttribute {
	return func(e *Error) {
		e.Detail = detail
	}
}

// SetLink 设置Link参数
func SetLink(link string) SetErrAttribute {
	return func(e *Error) {
		e.Link = link
	}
}

// NewError 新建一个Error
func NewError(code string, description string, setters ...SetErrAttribute) *Error {
	CheckCodeValid(code)
	e := &Error{
		Code:        code,
		Description: description,
	}

	// 设置可选属性
	for _, setter := range setters {
		setter(e)
	}

	return e
}

func (e *Error) Error() string {
	data := map[string]interface{}{
		"code":        e.Code,
		"description": e.Description,
	}
	if e.Solution != "" {
		data["solution"] = e.Solution
	}
	if len(e.Detail) != 0 {
		data["detail"] = e.Detail
	}
	if e.Link != "" {
		data["link"] = e.Link
	}
	errstr, err := jsoniter.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(errstr)
}

func CheckCodeValid(code string) {
	strs := strings.Split(code, ".")
	if len(strs) != 2 && len(strs) != 3 {
		panic("the parameter 'code' should be in 'Public.xxx' or 'xxx.xxx.xxx' format")
	}
	strCode := strs[1]
	if _, isExist := codesSet[strCode]; !isExist {
		panic("invalid parameter 'code'")
	}
	// Public.<错误标识>
	if len(strs) == 2 && prefix != strs[0]+"." {
		panic("the parameter 'code' should be 'Public.xxx' format")
	}
	// <服务名>.<错误标识>.<错误说明>
	if len(strs) == 3 {
		if len(strs[0]) == 0 || len(strs[0]) > 16 {
			panic("the length of the service name should be greater than 0 and not more than 16")
		}
		if len(strs[2]) == 0 || len(strs[2]) > 36 {
			panic("the length of the error description should be greater than 0 and not more than 36")
		}
	}
}
