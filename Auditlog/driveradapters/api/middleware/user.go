package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/common/types"
	"AuditLog/common/utils"
)

// SetUserIDToCtx 设置用户ID到上下文
// 有SetVisitorUserToCtx后，这个可能就不需要了（后面考虑删除）
func SetUserIDToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		if helpers.IsLocalDev() {
			// 本地开发环境，模拟用户信息
			c.Set("userId", "c40542b2-cf3a-11ef-ae53-52116f620329")
		}

		uid := c.GetString("userId")
		if uid == "" {
			return
		}

		ctxKey := enums.VisitUserIDCtxKey.String()
		c.Set(ctxKey, uid)

		// 设置到c.Request.Context()中
		//nolint:staticcheck
		newCtx := context.WithValue(c.Request.Context(), ctxKey, uid)
		utils.UpdateGinReqCtx(c, newCtx)
	}
}

// SetUserTokenToCtx 设置用户token到上下文
// 有SetVisitorUserToCtx后，这个可能就不需要了（后面考虑删除）
func SetUserTokenToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	1. 从header中获取token
		tokenID := c.GetHeader("Authorization")

		if tokenID == "" {
			if helpers.IsAaronLocalDev() {
				// 本地开发环境，模拟用户信息
				tokenID = "Bearer fa6ab300-abdb-11ef-81f6-ce93e9c3810d_token"
			}
		}

		token := strings.TrimPrefix(tokenID, "Bearer ")
		if token == "" {
			return
		}

		//	2. 设置到c.Request.Context()中
		ctxKey := enums.VisitUserTokenCtxKey.String()
		c.Set(ctxKey, token)

		//nolint:staticcheck
		newCtx := context.WithValue(c.Request.Context(), ctxKey, token)
		utils.UpdateGinReqCtx(c, newCtx)
	}
}

// SetVisitorUserToCtx 设置访问用户到上下文
func SetVisitorUserToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		info := types.NewVisitUserInfo()
		info.LoadFromGinCtx(c)

		// 设置到c.Request.Context()中
		ctxKey := enums.VisitUserInfoCtxKey.String()

		c.Set(ctxKey, info)

		//nolint:staticcheck
		newCtx := context.WithValue(c.Request.Context(), ctxKey, info)
		utils.UpdateGinReqCtx(c, newCtx)
	}
}
