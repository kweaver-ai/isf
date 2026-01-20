package driveradapters

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/interfaces"
)

var (
	healthOnce sync.Once
	h          interfaces.PublicRESTHandler
)

type healthHandler struct{}

// NewHealthHandler 创建document handler对象
func NewHealthHandler() interfaces.PublicRESTHandler {
	healthOnce.Do(func() {
		h = &healthHandler{}
	})

	return h
}

// RegisterPublic 注册外部API
func (h *healthHandler) RegisterPublic(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/health/ready", h.healthReady)
	routerGroup.GET("/health/alive", h.healthAlive)
}

func (h *healthHandler) healthReady(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "ready")
}

func (h *healthHandler) healthAlive(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "alive")
}
