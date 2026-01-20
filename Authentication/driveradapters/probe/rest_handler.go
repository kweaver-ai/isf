// Package probe 协议层
package probe

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// RESTHandler  健康检查接口
type RESTHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册开放API
	RegisterPublic(engine *gin.Engine)
}

type restHandler struct{}

var (
	once sync.Once
	r    RESTHandler
)

// NewRESTHandler 创建健康检查 handler对象
func NewRESTHandler() RESTHandler {
	once.Do(func() {
		r = &restHandler{}
	})

	return r
}

// RegisterPublic 注册内部API
func (r *restHandler) RegisterPrivate(engine *gin.Engine) {
	// 注册探针
	engine.GET("/health/ready", r.getHealth)
	engine.GET("/health/alive", r.getAlive)
}

// RegisterPublic 注册公开API
func (r *restHandler) RegisterPublic(engine *gin.Engine) {
	// 注册探针
	engine.GET("/health/ready", r.getHealth)
	engine.GET("/health/alive", r.getAlive)
}

func (r *restHandler) getHealth(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "OK")
}

func (r *restHandler) getAlive(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "OK")
}
