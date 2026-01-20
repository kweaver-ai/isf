// Package errors 服务错误
package errors

import (
	gerrors "github.com/kweaver-ai/go-lib/error"
)

// 服务错误码
const (
	strPrefix = "Authorization."
	// 角色名称冲突 Authorization.Conflict.RoleNameConflict
	RoleNameConflict = strPrefix + gerrors.Conflict + ".RoleNameConflict"
	// 角色不存在 Authorization.NotFound.RoleNotFound
	RoleNotFound = strPrefix + gerrors.NotFound + ".RoleNotFound"
)
