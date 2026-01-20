package driveradapters

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/utils/paramutils"
	"AuditLog/errors"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/middleware"
	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/models/rcvo"
)

var (
	hlOnce             sync.Once
	hl                 interfaces.PublicRESTHandler
	avaliableHisParams *paramutils.AvailableParam
)

type historyLogHandler struct {
	hlSvc interfaces.HistoryLog
}

func NewHistoryLogHandler() interfaces.PublicRESTHandler {
	hlOnce.Do(func() {
		hlc := &historyLogHandler{
			hlSvc: logics.NewHistoryLog(),
		}
		avaliableHisParams = paramutils.GetAvaliableParams(hlc.hlSvc.GetHistoryMetadata)
		hl = hlc
	})

	return hl
}

func (h *historyLogHandler) RegisterPublic(routerGroup *gin.RouterGroup) {
	setRoler := middleware.PermissionMiddleware([]string{common.SuperAdmin, common.SecAdmin})
	routerGroup.GET(
		"/history-log/download/pwdstatus",
		setRoler,
		h.getDownloadPwdStatus,
	)
	routerGroup.PUT(
		"/history-log/download/pwdstatus",
		setRoler,
		middleware.ValidateMiddleware(common.PutHistoryPwdStatus),
		middleware.VisitorParser,
		h.setDownloadPwdStatus,
	)

	downloadRoler := middleware.PermissionMiddleware([]string{
		common.SuperAdmin, common.SecAdmin, common.AuditAdmin, common.SysAdmin,
	})
	routerGroup.POST(
		"/history-log/download/task",
		downloadRoler,
		middleware.ValidateMiddleware(common.PostHistoryTask),
		middleware.VisitorParser,
		h.createDownloadTask,
	)
	routerGroup.GET(
		"/history-log/download/:taskId/progress",
		downloadRoler,
		middleware.VisitorParser,
		h.getDownloadProgress,
	)
	routerGroup.GET(
		"/history-log/download/:taskId/result",
		downloadRoler,
		middleware.VisitorParser,
		h.getDownloadResult,
	)

	getRoler := middleware.PermissionMiddleware([]string{common.SuperAdmin, common.SecAdmin, common.AuditAdmin, common.SysAdmin, common.OrgAudit})
	routerGroup.POST(
		"/report-center/history/:category/list",
		getRoler,
		middleware.VisitorParser,
		h.getDataList,
	)
}

func (h *historyLogHandler) getDownloadPwdStatus(c *gin.Context) {
	res, err := h.hlSvc.GetHistoryDownloadPwdStatus(c)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *historyLogHandler) setDownloadPwdStatus(c *gin.Context) {
	req := &lsmodels.HistoryLogDownloadPwdStatus{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrResponse(c, err)
		return
	}

	err := h.hlSvc.SetHistoryDownloadPwdStatus(c, req)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *historyLogHandler) createDownloadTask(c *gin.Context) {
	req := &lsmodels.HistoryLogDownloadReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrResponse(c, err)
		return
	}

	res, err := h.hlSvc.CreateDownloadTask(c, req)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *historyLogHandler) getDownloadProgress(c *gin.Context) {
	taskId := c.Param("taskId")
	if taskId == "" {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "taskId is required", nil))
		return
	}

	res, err := h.hlSvc.GetDownloadProgress(c, taskId)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *historyLogHandler) getDownloadResult(c *gin.Context) {
	taskId := c.Param("taskId")
	if taskId == "" {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "taskId is required", nil))
		return
	}

	res, err := h.hlSvc.GetHistoryDownloadResult(c, taskId)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// getDataList 获取历史日志报表数据列表
func (h *historyLogHandler) getDataList(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")

	// 验证 logType 是否合法
	if err := paramutils.CategoryCheck(c, logType); err != nil {
		common.ErrResponse(c, err)
		return
	}

	req := &rcvo.ReportGetDataListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrResponse(c, err)
		return
	}

	if req.Limit == 0 {
		req.Limit = 5000
	}

	// 校验搜索参数是否合法
	var searchParams []string
	for k := range req.Condition {
		searchParams = append(searchParams, k)
	}

	if err := paramutils.ParamsCheck(c, searchParams, avaliableHisParams.SearchFields, "condition"); err != nil {
		common.ErrResponse(c, err)
		return
	}

	// 校验排序参数是否合法
	var orderParams []string
	for _, v := range req.OrderBy {
		orderParams = append(orderParams, v.Field)
	}

	if err := paramutils.ParamsCheck(c, orderParams, avaliableHisParams.OrderFields, "order"); err != nil {
		common.ErrResponse(c, err)
		return
	}

	visitor := c.Value(common.VisitorKey).(*models.Visitor)
	userID := visitor.ID

	metaData, err := h.hlSvc.GetHistoryDataList(c, logType, req, userID)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, metaData)
}
