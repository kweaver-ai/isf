package driveradapters

import (
	"github.com/gin-gonic/gin"
)

// getXLang 解析获取 Header x-language
func getXLang(c *gin.Context) string {
	return c.GetHeader("x-language")
}
