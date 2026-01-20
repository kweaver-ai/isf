package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// ReservedNameHandler 保留名称接口
type ReservedNameHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}

type reservedNameHandler struct {
	updateReservedNameSchema *gojsonschema.Schema
	logger                   common.Logger
	reservedName             interfaces.LogicsReservedName
}

var (
	rnOnce    sync.Once
	rnHandler ReservedNameHandler
)

//go:embed jsonschema/reserved_name/update_reserved_name.json
var updateReservedNameSchemaStr string

// NewReservedNameHander 新建对象
func NewReservedNameHander() ReservedNameHandler {
	rnOnce.Do(func() {
		updateReservedNameSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateReservedNameSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		rnHandler = &reservedNameHandler{
			updateReservedNameSchema: updateReservedNameSchema,
			logger:                   common.NewLogger(),
			reservedName:             logics.NewReservedName(),
		}
	})
	return rnHandler
}

func (r *reservedNameHandler) RegisterPrivate(engine *gin.Engine) {
	engine.PUT("/api/user-management/v1/reserved-names/:id", r.updateReservedName)
	engine.DELETE("/api/user-management/v1/reserved-names/:id", r.deleteReservedName)
}

func (r *reservedNameHandler) updateReservedName(c *gin.Context) {
	var payload map[string]interface{}
	err := validateAndBindGin(c, r.updateReservedNameSchema, &payload)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	name := interfaces.ReservedNameInfo{
		ID:         c.Param("id"),
		Name:       payload["name"].(string),
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}
	err = r.reservedName.UpdateReservedName(name)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (r *reservedNameHandler) deleteReservedName(c *gin.Context) {
	err := r.reservedName.DeleteReservedName(c.Param("id"))
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}
