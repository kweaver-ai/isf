package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"AuditLog/common/enums"
	"AuditLog/common/utils"
)

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 设置traceID
		traceID := utils.UlidMake()

		ctxKey := enums.TraceIDCtxKey.String()
		c.Set(ctxKey, traceID)

		// 2. 设置到c.Request.Context()中
		//nolint:staticcheck
		newCtx := context.WithValue(c.Request.Context(), ctxKey, traceID)
		utils.UpdateGinReqCtx(c, newCtx)
	}
}
