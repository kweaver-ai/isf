package driveradapters

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	healthOnce      sync.Once
	healthSingleton *healthHandler
)

type healthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() RestHandler {
	healthOnce.Do(func() {
		healthSingleton = &healthHandler{}
	})

	return healthSingleton
}

// RegisterPrivate 注册内部API
func (h *healthHandler) RegisterPrivate(engine *gin.Engine) {
	// 就绪检查
	engine.GET("/health/ready", h.ready)
	// 存活检查
	engine.GET("/health/live", h.live)
}

// RegisterPublic 注册外部API
func (h *healthHandler) RegisterPublic(engine *gin.Engine) {
	// 就绪检查
	engine.GET("/health/ready", h.ready)
	// 存活检查
	engine.GET("/health/live", h.live)
}

func (h *healthHandler) ready(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *healthHandler) live(c *gin.Context) {
	c.Status(http.StatusOK)
}
