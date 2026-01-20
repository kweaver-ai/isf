package middleware

import (
	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/utils"
	"AuditLog/errors"
)

func PermissionMiddleware(allow []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roles := ctx.GetStringSlice("userRoles")
		if len(utils.Intersection(roles, allow)) == 0 {
			err := errors.NewCtx(ctx, errors.ForbiddenErr, "invalid user role", nil)
			common.ErrResponse(ctx, err)
			return
		}
	}
}
