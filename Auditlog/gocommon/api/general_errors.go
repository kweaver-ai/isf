package api

import (
	"net/http"
)

var (
	publicErrorBuilder           TypedErrorBuilder
	ErrBadRequestPublic          func(*ErrorInfo) *Error
	ErrInternalServerErrorPublic func(*ErrorInfo) *Error
	ErrUnauthorization           func(*ErrorInfo) *Error
)

var (
	// Public 所有服务公用
	Public = 0
)

func init() {
	// 初始化public错误码构建函数，错误码中间三位为000，目前有[400000000,500000000]
	publicErrorBuilder = NewErrorBuilder(Public)

	// 500000000
	ErrInternalServerErrorPublic = publicErrorBuilder.OfType(&ErrorType{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  0,
		ErrorInfo: ErrorInfo{
			Message: "Server unavailable.",
		},
	})
	// 400000000
	ErrBadRequestPublic = publicErrorBuilder.OfType(&ErrorType{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  0,
		ErrorInfo: ErrorInfo{
			Message: "Invalid request.",
		},
	})
	// 401000000
	ErrUnauthorization = publicErrorBuilder.OfType(&ErrorType{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  0,
		ErrorInfo: ErrorInfo{
			Message: "Not authorized.",
		},
	})
}
