package apierr

import "AuditLog/errors"

// 通用错误码
const (
	// BadRequestErr 通用错误码，客户端请求错误
	BadRequestErr = errors.BadRequestErr
	// UnauthorizedErr 通用错误码，未授权或者授权已过期
	UnauthorizedErr = errors.UnauthorizedErr
	// ForbiddenErr 通用错误码，禁止访问
	ForbiddenErr = errors.ForbiddenErr
	// ResourceNotFoundErr 通用错误码，请求资源不存在
	ResourceNotFoundErr = errors.ResourceNotFoundErr
	// MethodNotAllowedErr 通用错误码，目标资源不支持该方法
	MethodNotAllowedErr = errors.MethodNotAllowedErr
	// ConflictErr 通用错误码，资源冲突
	ConflictErr = errors.ConflictErr
	// TooManyRequestsErr 通用错误码，请求过于频繁
	TooManyRequestsErr = errors.TooManyRequestsErr
	// InternalErr 通用错误码，服务端内部错误
	InternalErr = errors.InternalErr
)
