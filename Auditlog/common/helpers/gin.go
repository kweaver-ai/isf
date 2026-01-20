package helpers

import (
	"github.com/gin-gonic/gin"

	"AuditLog/common/enums"
)

func GetXLanguage(c *gin.Context) string {
	return c.GetHeader(enums.XLanguage)
}
