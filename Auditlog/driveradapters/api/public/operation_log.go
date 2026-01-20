package driveradapters

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/utils"
	oprinject "AuditLog/domain/service/inject/operation_log"
	"AuditLog/driveradapters/api/middleware"
	"AuditLog/errors"
	"AuditLog/interfaces"
	oprlogdriveri "AuditLog/interfaces/driveradapter/operation_log"
)

type logHandler struct {
	oprLogSvc oprlogdriveri.IOprLogSvc
}

// NewOperationLogHandler 创建运营日志 handler对象
func NewOperationLogHandler() interfaces.PublicRESTHandler {
	return &logHandler{
		oprLogSvc: oprinject.NewOprLogSvc(),
	}
}

// RegisterPublic 注册外部API
func (l *logHandler) RegisterPublic(routerGroup *gin.RouterGroup) {
	group := routerGroup.Group("/operation-log")

	allowRoles := []string{common.NormalUser}

	opt := middleware.RolePermissionOption{
		IsFromOprLogAPI: true,
	}
	group.Use(middleware.RolePermission(allowRoles, opt))

	// 单个日志上报
	group.POST("/:biz_type", l.report)

	// 批量
	group.POST("/batch/:biz_type", l.report)
}

// report 上报运营日志（客户端调用）
func (l *logHandler) report(c *gin.Context) {
	//  获取path中参数
	bizTypeStr := c.Param("biz_type")

	bizType := oprlogenums.BizType(bizTypeStr)

	// 1 检查业务类型
	if !bizType.Check() {
		err := errors.NewCtx(c, errors.BadRequestErr, "invalid biz_type", "")

		common.ErrResponse(c, err)

		return
	}

	// 2. 获取body中的数据 byte
	bodyData, err := c.GetRawData()
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	if len(bodyData) == 0 {
		err = errors.NewCtx(c, errors.BadRequestErr, "request body is empty", "")
		common.ErrResponse(c, err)

		return
	}

	// 3. 如果不是批量，将数据转换为数组
	// 判断是否是批量
	isBatch := strings.Contains(c.Request.URL.Path, "operation-log/batch")
	if !isBatch {
		bodyData = utils.JSONObjectToArray(bodyData)
	}

	// 4. 添加业务类型到日志中
	bodyData, err = utils.AddKeyToJSONArrayBys(bodyData, "biz_type", bizTypeStr)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	// 5. 校验日志
	err = l.checkLogs(c, bodyData, bizType)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	// 6. 发送到mq
	// 每次最多发送20条
	splits := utils.SplitJSONArrayBys(bodyData, 20)
	for _, split := range splits {
		err = l.oprLogSvc.WriteOperationLogToMQ(c, bizType.ToTopic(), split)
		if err != nil {
			common.ErrResponse(c, err)
			return
		}
	}

	// 6. 返回结果
	c.Status(http.StatusCreated)
}
