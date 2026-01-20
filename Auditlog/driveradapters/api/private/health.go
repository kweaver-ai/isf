package driveradapters

import (
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/interfaces"
)

var (
	healthOnce sync.Once
	h          interfaces.PrivateRESTHandler
)

type healthHandler struct{}

// NewHealthHandler 创建document handler对象
func NewHealthHandler() interfaces.PrivateRESTHandler {
	healthOnce.Do(func() {
		h = &healthHandler{}
	})

	return h
}

// RegisterPrivate 注册私有API
func (h *healthHandler) RegisterPrivate(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/health/ready", h.healthTest)
	routerGroup.GET("/health/alive", h.healthTest)
}

func (h *healthHandler) healthTest(c *gin.Context) {
}
