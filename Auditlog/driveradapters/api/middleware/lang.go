package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
)

func SetLangFromHEADER() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := helpers.GetLangFromHeader(c)

		ctxKey := enums.VisitLangCtxKey.String()
		c.Set(ctxKey, lang)

		// 将lang设置到c.Request.Context()中
		//nolint:staticcheck
		newCtx := context.WithValue(c.Request.Context(), ctxKey, lang)
		utils.UpdateGinReqCtx(c, newCtx)
	}
}
