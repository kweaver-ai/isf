package errors

import (
	"net/http"

	"policy_mgnt/utils/gocommon/api"
)

var (
	policyErrBuilder             api.TypedErrorBuilder
	publicErrorBuilder           api.TypedErrorBuilder
	ErrBadRequestPublic          func(*api.ErrorInfo) *api.Error
	ErrInternalServerErrorPublic func(*api.ErrorInfo) *api.Error
	ErrPolicyLocked              func(*api.ErrorInfo) *api.Error
	ErrNotFound                  func(*api.ErrorInfo) *api.Error
	ErrConflict                  func(*api.ErrorInfo) *api.Error
	ErrTooManyRequests           func(*api.ErrorInfo) *api.Error
	ErrUnauthorization           func(*api.ErrorInfo) *api.Error
	ErrNoPermission              func(*api.ErrorInfo) *api.Error
	ErrInvalideName              func(*api.ErrorInfo) *api.Error
)

func init() {
	// 初始化policy错误码构建函数，错误码中间三位为013
	policyErrBuilder = api.NewErrorBuilder(api.PolicyManagement)
	// 初始化public错误码构建函数，错误码中间三位为000，目前有[400000000,500000000]
	publicErrorBuilder = api.NewErrorBuilder(api.Public)

	// 500000000
	ErrInternalServerErrorPublic = publicErrorBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Server unavailable.",
		},
	})
	// 400000000
	ErrBadRequestPublic = publicErrorBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Invalid request.",
		},
	})
	// 401000000
	ErrUnauthorization = publicErrorBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Not authorized.",
		},
	})
	// 400013001 名称错误
	ErrInvalideName = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  1,
		ErrorInfo: api.ErrorInfo{
			Message: "Invalid name.",
		},
	})
	// 400013100
	ErrPolicyLocked = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  100,
		ErrorInfo: api.ErrorInfo{
			Message: "Policy locked.",
		},
	})
	// 403013000
	ErrNoPermission = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusForbidden,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "No permission to do this service.",
		},
	})
	// 404013000
	ErrNotFound = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusNotFound,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Resource not found.",
		},
	})
	// 409013000
	ErrConflict = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusConflict,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Resource conflicts.",
		},
	})
	// 429013000
	ErrTooManyRequests = policyErrBuilder.OfType(&api.ErrorType{
		StatusCode: http.StatusTooManyRequests,
		ErrorCode:  0,
		ErrorInfo: api.ErrorInfo{
			Message: "Too many requests.",
		},
	})
}
