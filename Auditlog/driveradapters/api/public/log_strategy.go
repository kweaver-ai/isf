package driveradapters

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"AuditLog/common"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/common/utils"
	"AuditLog/common/utils/paramutils"
	"AuditLog/errors"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/middleware"
	"AuditLog/models/lsmodels"
)

var (
	lsOnce sync.Once
	ls     interfaces.PublicRESTHandler
)

type logStrategyHandler struct {
	lsSvc interfaces.LogStrategy
	ssSvc interfaces.LogScopeStrategy
}

func NewLogStrategyHandler() interfaces.PublicRESTHandler {
	lsOnce.Do(func() {
		ls = &logStrategyHandler{
			lsSvc: logics.NewLogStrategy(),
			ssSvc: logics.NewScopeStrategy(),
		}
	})

	return ls
}

func (l *logStrategyHandler) RegisterPublic(routerGroup *gin.RouterGroup) {
	roler := middleware.PermissionMiddleware([]string{common.SuperAdmin, common.SecAdmin})
	routerGroup.GET(
		"/log-strategy/dump",
		roler,
		l.getDumpStrategy,
	)
	routerGroup.PUT(
		"/log-strategy/dump",
		roler,
		middleware.ValidateMiddleware(common.PutDumpStrategy),
		middleware.VisitorParser,
		l.updateDumpStrategy,
	)

	secRoler := middleware.PermissionMiddleware([]string{common.SecAdmin})
	routerGroup.GET(
		"/log-strategy/scope",
		secRoler,
		l.getScopeStrategies,
	)
	routerGroup.POST(
		"/log-strategy/scope",
		secRoler,
		middleware.ValidateMiddleware(common.ScopeStrategy),
		middleware.VisitorParser,
		l.newScopeStrategy,
	)
	routerGroup.PUT(
		"/log-strategy/scope/:id",
		secRoler,
		middleware.ValidateMiddleware(common.ScopeStrategy),
		middleware.VisitorParser,
		l.updateScopeStrategy,
	)
	routerGroup.DELETE(
		"/log-strategy/scope/:id",
		secRoler,
		middleware.VisitorParser,
		l.deleteScopeStrategy,
	)
}

// 获取日志转存策略
func (l *logStrategyHandler) getDumpStrategy(c *gin.Context) {
	fields := c.QueryArray("field")
	if err := paramutils.DumpFieldsCheck(c, fields); err != nil {
		common.ErrResponse(c, err)
		return
	}

	dumpStrategy, err := l.lsSvc.GetDumpStrategy(c, fields)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, dumpStrategy)
}

// 更新日志转存策略
func (l *logStrategyHandler) updateDumpStrategy(c *gin.Context) {
	fields := c.QueryArray("field")
	if len(fields) == 0 {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "field is required", nil))
		return
	}

	if err := paramutils.DumpFieldsCheck(c, fields); err != nil {
		common.ErrResponse(c, err)
		return
	}

	req := make(map[string]interface{})
	if err := common.ParseBody(c, &req); err != nil {
		common.ErrResponse(c, err)
		return
	}

	strategy := make(map[string]interface{})
	// 检查req中的字段是否在fields中
	for _, field := range fields {
		fieldValue := req[field]
		if fieldValue == nil {
			cause := fmt.Sprintf("invalid field: %s", field)
			common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, cause, nil))

			return
		}

		strategy[field] = fieldValue
	}

	err := l.lsSvc.SetDumpStrategy(c, strategy)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// 获取日志查看范围配置
func (l *logStrategyHandler) getScopeStrategies(c *gin.Context) {
	req := &lsmodels.GetScopeStrategyReq{
		Limit:  200,
		Offset: 0,
	}

	// 统一错误处理函数
	handleError := func(msg string) {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, msg, nil))
	}

	// 解析并验证整数类型参数
	parseIntParam := func(value string, min int, max int, dest *int, field string) bool {
		if value == "" {
			return true
		}

		val, err := strconv.Atoi(value)
		if err != nil || val < min || val > max {
			handleError("invalid " + field)
			return false
		}

		*dest = val

		return true
	}

	// 解析基本参数
	if !parseIntParam(c.Query("limit"), 1, 1000, &req.Limit, "limit") ||
		!parseIntParam(c.Query("offset"), 0, math.MaxInt, &req.Offset, "offset") {
		return
	}

	// 解析并验证 category
	if category := c.Query("category"); category != "" {
		categoryInt, err := strconv.Atoi(category)
		if err != nil || !utils.ExistsGeneric(lsconsts.AllLogCategory, categoryInt) {
			handleError("invalid category")
			return
		}

		req.Category = int8(categoryInt)
	}

	// 解析并验证 type
	if logType := c.Query("type"); logType != "" {
		logTypeInt, err := strconv.Atoi(logType)
		if err != nil || !utils.ExistsGeneric(common.AllLogTypeInt, logTypeInt) {
			handleError("invalid type")
			return
		}

		req.Type = int8(logTypeInt)
	}

	// 验证 role
	if role := c.Query("role"); role != "" {
		if !common.InArray(role, common.MutuallyRoles) {
			handleError("invalid role")
			return
		}

		req.Role = role
	}

	scopeStrategies, err := l.ssSvc.GetStrategy(c, req)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, scopeStrategies)
}

// 新增日志查看范围配置
func (l *logStrategyHandler) newScopeStrategy(c *gin.Context) {
	reqBody := &lsmodels.ScopeStrategyVO{}
	if err := common.ParseBody(c, reqBody); err != nil {
		common.ErrResponse(c, err)
		return
	}

	id, err := l.ssSvc.NewStrategy(c, reqBody)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (l *logStrategyHandler) updateScopeStrategy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "id is required", nil))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "invalid id", nil))
		return
	}

	reqBody := &lsmodels.ScopeStrategyVO{}
	if err := common.ParseBody(c, reqBody); err != nil {
		common.ErrResponse(c, err)
		return
	}

	err = l.ssSvc.UpdateStrategy(c, idInt, reqBody)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (l *logStrategyHandler) deleteScopeStrategy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "id is required", nil))
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		common.ErrResponse(c, errors.NewCtx(c, errors.BadRequestErr, "invalid id", nil))
		return
	}

	err = l.ssSvc.DeleteStrategy(c, idInt)
	if err != nil {
		common.ErrResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
