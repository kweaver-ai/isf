package rest

import (
	"fmt"

	errorv2 "github.com/kweaver-ai/go-lib/error"
	jsoniter "github.com/json-iterator/go"
)

var (
	// Languages 支持的语言
	Languages = [3]string{"zh_CN", "zh_TW", "en_US"}
)

var (
	i18n         = make(map[int]map[string]string)
	code2Message = make(map[int]string)
	intCodeToStr = map[int]string{
		400: errorv2.PublicBadRequest,
		401: errorv2.PublicUnauthorized,
		403: errorv2.PublicForbidden,
		404: errorv2.PublicNotFound,
		409: errorv2.PublicConflict,
		500: errorv2.PublicInternalServerError,
		503: errorv2.PublicServiceUnavailable,
	}
)

// SetLang 设置语言
func SetLang(lang string) {
	valid := false
	for _, l := range Languages {
		if l == lang {
			valid = true
		}
	}
	if !valid {
		panic("invalid lang")
	}
	for code := range i18n {
		code2Message[code] = i18n[code][lang]
	}
}

// Register 注册code对应message
func Register(langRes map[int]map[string]string) {
	for code, message := range langRes {
		if _, ok := i18n[code]; ok {
			panic(fmt.Sprintf("duplicate code: %v", code))
		}
		i18n[code] = make(map[string]string)
		for _, lang := range Languages {
			if m, ok := message[lang]; ok {
				i18n[code][lang] = m
			} else {
				panic(fmt.Sprintf("language %v not exists", lang))
			}
		}
	}
}

// HTTPError 服务错误结构体
type HTTPError struct {
	Cause       string                 `json:"cause"`
	Code        int                    `json:"code"`
	Message     string                 `json:"message"`
	Detail      map[string]interface{} `json:"detail,omitempty"`
	Description string                 `json:"description,omitempty"`
	Solution    string                 `json:"solution,omitempty"`
	CodeStr     string
	useCodeStr  bool
}

// SetErrAttribute 设置参数
type SetErrAttribute func(*HTTPError)

// SetDescription 设置description参数
// Deprecated: 如需设置 Description, 请使用 NewHTTPErrorV2
func SetDescription(description string) SetErrAttribute {
	return func(err *HTTPError) {
		err.Description = description
	}
}

// SetSolution 设置solution参数
func SetSolution(solution string) SetErrAttribute {
	return func(err *HTTPError) {
		err.Solution = solution
	}
}

// SetDetail 设置detail参数
func SetDetail(detail map[string]interface{}) SetErrAttribute {
	return func(e *HTTPError) {
		e.Detail = detail
	}
}

// SetCodeStr 设置codeStr
func SetCodeStr(codeStr string) SetErrAttribute {
	errorv2.CheckCodeValid(codeStr)
	return func(e *HTTPError) {
		e.CodeStr = codeStr
	}
}

// NewHTTPErrorV2 新建一个HTTPError
func NewHTTPErrorV2(code int, description string, setters ...SetErrAttribute) *HTTPError {
	checkCodeValid(code)
	e := &HTTPError{
		Code:        code,
		Message:     code2Message[code],
		Cause:       description,
		Description: description,
	}

	// 设置可选属性
	for _, setter := range setters {
		setter(e)
	}

	return e
}

// NewHTTPError 新建一个HTTPError
func NewHTTPError(cause string, code int, detail map[string]interface{}, params ...SetErrAttribute) *HTTPError {
	checkCodeValid(code)
	httpError := &HTTPError{
		Cause:   cause,
		Code:    code,
		Message: code2Message[code],
		Detail:  detail,
	}
	for _, p := range params {
		p(httpError)
	}

	return httpError
}

func (e *HTTPError) Error() string {
	data := map[string]interface{}{
		"cause": e.Cause,
	}
	if len(e.Detail) != 0 {
		data["detail"] = e.Detail
	}
	if e.Description != "" {
		data["description"] = e.Description
	}
	if e.Solution != "" {
		data["solution"] = e.Solution
	}

	if e.useCodeStr {
		// 使用NewXXX方法创建HTTPError实例时已经检查了code合法性
		if len(e.CodeStr) == 0 {
			data["code"] = intCodeToStr[e.Code/1e6]
		} else {
			data["code"] = e.CodeStr
		}
	} else {
		data["code"] = e.Code
		data["message"] = e.Message
	}
	errStr, err := jsoniter.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(errStr)
}

// ExHTTPError 其他服务响应的错误结构体
type ExHTTPError struct {
	Status int
	Body   []byte
}

func (err ExHTTPError) Error() string {
	return string(err.Body)
}

// checkCodeValid 检查错误码是否合法
func checkCodeValid(code int) {
	if code < BadRequest || code >= 600000000 {
		panic("the parameter 'code' length should be 9")
	}
	digit := code / 1e6
	if _, isExist := intCodeToStr[digit]; !isExist {
		panic("the first three digits of 'code' should be one of 400/401/403/404/409/500/503")
	}
}
