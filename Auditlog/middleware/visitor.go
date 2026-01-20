package middleware

import (
	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/models"
)

// VisitorParser 访问者信息解析中间件
func VisitorParser(c *gin.Context) {
	// Email和Mac字段没有在前面的中间件上设置，所以拿不到
	c.Set(common.VisitorKey, &models.Visitor{
		ID:        c.GetString("userId"),
		Name:      c.GetString("name"),
		CsfLevel:  c.GetFloat64("csfLevel"),
		IP:        c.GetString("ip"),
		Udid:      c.GetString("udid"),
		AgentType: c.GetString("clientType"),
		Roles:     c.GetStringSlice("userRoles"),
		Token:     c.GetString("token"),
		Type:      c.GetString("visitorType"),
	})
}
