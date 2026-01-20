package driveradapters

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/utils/paramutils"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/models/rcvo"
)

var (
	historyOnce        sync.Once
	his                interfaces.PrivateRESTHandler
	avaliableHisParams *paramutils.AvailableParam
)

type historyHandler struct {
	historyLog interfaces.HistoryLog
}

// NewHistoryHandler 创建document handler对象
func NewHistoryHandler() interfaces.PrivateRESTHandler {
	historyOnce.Do(func() {
		hh := &historyHandler{
			historyLog: logics.NewHistoryLog(),
		}
		avaliableHisParams = paramutils.GetAvaliableParams(hh.historyLog.GetHistoryMetadata)
		his = hh
	})

	return his
}

// RegisterPrivate 注册内部API
func (h *historyHandler) RegisterPrivate(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/report-center/history/:category/metadata", h.getMetaData)
	routerGroup.POST("/report-center/history/:category/list", h.getDataList)
}

// getMetaData 获取历史日志报表元数据
func (h *historyHandler) getMetaData(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")

	// 验证 logType 是否合法
	if err := paramutils.CategoryCheck(c, logType); err != nil {
		common.ErrResponse(c, err)
		return
	}

	metaData, err := h.historyLog.GetHistoryMetadata()
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, metaData)
}

// getDataList 获取历史日志报表数据列表
func (h *historyHandler) getDataList(c *gin.Context) {
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

	userID := c.GetHeader("X-RC-VISIT-USER-ID")

	metaData, err := h.historyLog.GetHistoryDataList(c, logType, req, userID)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, metaData)
}
