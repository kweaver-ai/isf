package driveradapters

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/constants/rclogconsts"
	"AuditLog/common/utils/paramutils"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/models/rcvo"
)

var (
	activeOnce      sync.Once
	ac              interfaces.PrivateRESTHandler
	avaliableParams *paramutils.AvailableParam
)

type activeHandler struct {
	activeLog interfaces.ActiveLog
}

// NewActiveHandler 创建document handler对象
func NewActiveHandler() interfaces.PrivateRESTHandler {
	activeOnce.Do(func() {
		ach := &activeHandler{
			activeLog: logics.NewActiveLog(),
		}
		avaliableParams = paramutils.GetAvaliableParams(ach.activeLog.GetActiveMetadata)
		ac = ach
	})

	return ac
}

// RegisterPrivate 注册内部API
func (ac *activeHandler) RegisterPrivate(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/report-center/active/:category/metadata", ac.getMetaData)
	routerGroup.POST("/report-center/active/:category/list", ac.getDataList)
	routerGroup.POST("/report-center/active/:category/field/:field/values", ac.getFieldValues)
}

// getMetaData 获取活跃日志报表元数据
func (ac *activeHandler) getMetaData(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")

	// 验证 logType 是否合法
	if err := paramutils.CategoryCheck(c, logType); err != nil {
		common.ErrResponse(c, err)
		return
	}

	metaData, err := ac.activeLog.GetActiveMetadata()
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, metaData)
}

// getDataList 获取活跃日志报表数据列表
func (ac *activeHandler) getDataList(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")

	// 验证 logType 是否合法
	if err := paramutils.CategoryCheck(c, logType); err != nil {
		common.ErrResponse(c, err)
		return
	}

	req := &rcvo.ReportGetActiveDataListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrResponse(c, err)
		return
	}

	body := &rcvo.ReportGetDataListReq{}
	body.Limit = req.Limit
	body.Offset = req.Offset
	body.Condition = make(map[string]interface{})
	body.OrderBy = req.OrderBy

	if len(req.IDs) > 0 {
		strIDs := make([]string, len(req.IDs))
		for i, id := range req.IDs {
			strIDs[i] = strconv.FormatUint(id, 10)
		}

		body.IDs = strIDs
	}

	for k, v := range req.Condition {
		if k == rclogconsts.Date {
			body.Condition[k] = []interface{}{
				float64(v.([]interface{})[0].(float64)) * 1000,
				float64(v.([]interface{})[1].(float64)) * 1000,
			}
		} else {
			body.Condition[k] = v
		}
	}

	// 当请求没传limit时，设置默认值5000
	if body.Limit == 0 {
		body.Limit = 5000
	}

	// 校验搜索参数是否合法
	var searchParams []string
	for k := range req.Condition {
		searchParams = append(searchParams, k)
	}

	if err := paramutils.ParamsCheck(c, searchParams, avaliableParams.SearchFields, "condition"); err != nil {
		common.ErrResponse(c, err)
		return
	}

	// 校验排序参数是否合法
	var orderParams []string
	for _, v := range req.OrderBy {
		orderParams = append(orderParams, v.Field)
	}

	if err := paramutils.ParamsCheck(c, orderParams, avaliableParams.OrderFields, "order"); err != nil {
		common.ErrResponse(c, err)
		return
	}

	userID := c.GetHeader("X-RC-VISIT-USER-ID")

	list, err := ac.activeLog.GetActiveDataList(c, logType, body, userID)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// getFieldValues 获取活跃日志报表字段值
func (ac *activeHandler) getFieldValues(c *gin.Context) {
	// 获取path中参数
	logType := c.Param("category")
	field := c.Param("field")

	// 验证 logType 是否合法
	if err := paramutils.CategoryCheck(c, logType); err != nil {
		common.ErrResponse(c, err)
		return
	}

	// 解析请求体
	var body rcvo.ReportGetFieldValuesReq
	if err := c.ShouldBindJSON(&body); err != nil {
		common.ErrResponse(c, err)
		return
	}
	// 当请求没传limit时，设置默认值5000
	if body.Limit == 0 {
		body.Limit = 5000
	}

	// 校验搜索参数是否合法
	var searchParams []string
	for k := range body.Condition {
		searchParams = append(searchParams, k)
	}

	if err := paramutils.ParamsCheck(c, searchParams, avaliableParams.SearchFields, "condition"); err != nil {
		common.ErrResponse(c, err)
		return
	}
	// 校验关键字参数是否合法
	if body.KeyWord != "" {
		if err := paramutils.ParamsCheck(c, []string{field}, avaliableParams.KeyWordFields, "keyword"); err != nil {
			common.ErrResponse(c, err)
			return
		}
	}

	// 创建请求值对象
	req := &rcvo.ReportGetFieldValuesReq{
		Field: field,
	}
	req.Limit = body.Limit
	req.Offset = body.Offset
	req.Condition = body.Condition
	req.KeyWord = body.KeyWord

	values, err := ac.activeLog.GetActiveFieldValues(c, logType, req)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, values)
}
