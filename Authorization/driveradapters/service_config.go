// Package driveradapters service config 服务配置处理
package driveradapters

import (
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"

	"Authorization/common"
)

var (
	serviceConfigOnce      sync.Once
	serviceConfigSingleton *serviceConfigDriver
)

type serviceConfigDriver struct {
	logger common.Logger
}

// NewServiceConfigDriver 创建serviceConfigDriver
func NewServiceConfigDriver() RestHandler {
	serviceConfigOnce.Do(func() {
		serviceConfigSingleton = &serviceConfigDriver{
			logger: common.NewLogger(),
		}
	})
	return serviceConfigSingleton
}

// RegisterPrivate 注册内部API
func (s *serviceConfigDriver) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/authorization/v1/config/log/level", s.getLogLevel)
	engine.PUT("/api/authorization/v1/config/log/level", s.setLogLevel) // 可在维持服务正常的运行的情况下动态修改日志等级
}

// RegisterPublic 注册外部API
func (s *serviceConfigDriver) RegisterPublic(_ *gin.Engine) {
}

func (s *serviceConfigDriver) setLogLevel(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.logger.Errorln(err)
		rest.ReplyErrorV2(c, err)
		return
	}

	var kv map[string]string
	if err = jsoniter.Unmarshal(body, &kv); err != nil {
		s.logger.Errorln(err)
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid request body: "+err.Error())
		rest.ReplyErrorV2(c, err)
		return
	}

	levelStr, ok := kv["level"]
	if !ok {
		err = gerrors.NewError(gerrors.PublicBadRequest, "key 'level' was needed")
		rest.ReplyErrorV2(c, err)
		return
	}

	if err = common.SetLogLevel(levelStr); err != nil {
		s.logger.Errorln(err)
		err = gerrors.NewError(gerrors.PublicBadRequest, err.Error())
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (s *serviceConfigDriver) getLogLevel(c *gin.Context) {
	rest.ReplyOK(c, http.StatusOK, map[string]string{"level": common.GetLogLevel()})
}
