package driveradapters

import (
	"strconv"

	"github.com/gin-gonic/gin"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

// RestHandler 通用适配器接口
type RestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}

type listQueryInfo struct {
	offset int
	limit  int
}

// listResultInfo 分页列举结构
type listResultInfo struct {
	Entries    []any `json:"entries" `
	TotalCount int   `json:"total_count"`
}

func getListQueryParam(c *gin.Context) (info listQueryInfo, err error) {
	// 获取数据下标
	info.offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		err = gerrors.NewError(gerrors.PublicBadRequest, "offset is illeagal")
		return
	}

	if info.offset < 0 {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid offset(>=0)")
		return
	}

	// 获取分页数据
	info.limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		err = gerrors.NewError(gerrors.PublicBadRequest, "limit is illeagal")
		return
	}
	if info.limit < 1 || info.limit > 1000 {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid limit([ 1 .. 1000 ])")
		return
	}

	return
}
