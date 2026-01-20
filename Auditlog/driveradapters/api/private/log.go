package driveradapters

import (
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/middleware"
	"AuditLog/models"
)

var (
	logOnce sync.Once
	l       interfaces.PrivateRESTHandler
)

type logHandler struct {
	logMgnt interfaces.LogMgnt
}

// NewLogHandler 创建document handler对象
func NewLogHandler() interfaces.PrivateRESTHandler {
	logOnce.Do(func() {
		l = &logHandler{
			logMgnt: logics.NewLogMgnt(),
		}
	})

	return l
}

// RegisterPrivate 注册私有API
func (l *logHandler) RegisterPrivate(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/log/:category", middleware.ValidateMiddleware(common.PostAuditLog), l.add)
}

func (l *logHandler) add(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")

	var params models.AuditLog
	if err := common.ParseBody(c, &params); err != nil {
		common.ErrorResponse(c, err)
		return
	}
	// 逻辑层
	info := &models.SendLogVo{
		Language:   common.SvcConfig.Languaue,
		LogType:    logType,
		LogContent: &params,
	}

	err := l.logMgnt.SendLog(info)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	c.Status(204)
}
