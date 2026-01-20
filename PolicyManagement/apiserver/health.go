package apiserver

import "github.com/gin-gonic/gin"

type healthHandler struct {
}

func newHealthHandler() (*healthHandler, error) {
	return newHealthHandlerWithMgnt(), nil
}

func newHealthHandlerWithMgnt() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) AddRouters(r *gin.RouterGroup) {
	r.GET("/health/ready", h.healthTest)
	r.GET("/health/alive", h.healthTest)
}

func (h *healthHandler) healthTest(c *gin.Context) {
}
