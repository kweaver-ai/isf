// Package driveradapters AnyShare 入站适配器
package driveradapters

import (
	"context"
	_ "embed" // 标准用法
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"

	"Authorization/interfaces"
	"Authorization/logics"
)

//go:embed jsonschema/obligation/add.json
var addObligationSchemaStr string

//go:embed jsonschema/obligation/update.json
var updateObligationSchemaStr string

var (
	obligationOnce    sync.Once
	obligationHandler RestHandler
)

type obligationRestHandler struct {
	obligation      interfaces.LogicsObligation
	hydra           interfaces.Hydra
	addSchemaStr    *gojsonschema.Schema
	updateSchemaStr *gojsonschema.Schema
}

// NewObligationTemplateRestHandler 权限适配器接口
func NewObligationRestHandler() RestHandler {
	obligationOnce.Do(func() {
		obligationHandler = &obligationRestHandler{
			obligation:      logics.NewObligation(),
			hydra:           newHydra(),
			addSchemaStr:    newJSONSchema(addObligationSchemaStr),
			updateSchemaStr: newJSONSchema(updateObligationSchemaStr),
		}
	})
	return obligationHandler
}

// RegisterPrivate 注册内部API
func (o *obligationRestHandler) RegisterPrivate(_ *gin.Engine) {
}

// RegisterPublic 注册外部API
func (o *obligationRestHandler) RegisterPublic(engine *gin.Engine) {
	// 义务管理接口
	engine.POST("/api/authorization/v1/obligations", o.add)
	engine.GET("/api/authorization/v1/obligations", o.get)
	engine.DELETE("/api/authorization/v1/obligations/:id", o.delete)
	engine.GET("/api/authorization/v1/obligations/:id", o.getByID)
	engine.PUT("/api/authorization/v1/obligations/:id/:fields", o.update)

	// 查询接口
	// 义务查询
	engine.GET("/api/authorization/v1/query-obligations", o.queryObligation)
}

func (o *obligationRestHandler) add(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	var obligationDr map[string]any
	if err = validateAndBindGin(c, o.addSchemaStr, &obligationDr); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligation := interfaces.ObligationInfo{
		TypeID: obligationDr["type_id"].(string),
		Name:   obligationDr["name"].(string),
		Value:  obligationDr["value"],
	}

	descriptionJson, ok := obligationDr["description"]
	if ok {
		obligation.Description = descriptionJson.(string)
	}

	id, err := o.obligation.Add(context.Background(), &visitor, &obligation)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/authorization/v1/obligations/%s", id))

	rest.ReplyOK(c, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (o *obligationRestHandler) update(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationID := c.Param("id")
	if obligationID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	fields := strings.Split(c.Param("fields"), ",")
	var nameExist, descriptionExist, valueExist bool
	for _, v := range fields {
		if v == "name" {
			nameExist = true
		}
		if v == "description" {
			descriptionExist = true
		}
		if v == "value" {
			valueExist = true
		}
	}

	var obligationDr map[string]any
	if err = validateAndBindGin(c, o.updateSchemaStr, &obligationDr); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	name := ""
	if nameExist {
		nameJson, ok := obligationDr["name"]
		if ok {
			name = nameJson.(string)
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "param name is illegal")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	description := ""
	if descriptionExist {
		descriptionJson, ok := obligationDr["description"]
		if ok {
			description = descriptionJson.(string)
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "param description is illegal")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	var value any
	if valueExist {
		valueJson, ok := obligationDr["value"]
		if ok {
			value = valueJson
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "param value is illegal")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	err = o.obligation.Update(context.Background(), &visitor, obligationID, name, nameExist,
		description, descriptionExist, value, valueExist)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *obligationRestHandler) delete(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationID := c.Param("id")
	if obligationID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	err = o.obligation.Delete(context.Background(), &visitor, obligationID)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *obligationRestHandler) getByID(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationID := c.Param("id")
	if obligationID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	obligation, err := o.obligation.GetByID(context.Background(), &visitor, obligationID)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	result := map[string]any{
		"id":          obligation.ID,
		"type_id":     obligation.TypeID,
		"name":        obligation.Name,
		"description": obligation.Description,
		"value":       obligation.Value,
	}
	rest.ReplyOK(c, http.StatusOK, result)
}

func (o *obligationRestHandler) get(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取列举义务信息
	queryInfo, err := getListQueryParam(c)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}
	obligationSearchInfo := interfaces.ObligationSearchInfo{
		Offset: queryInfo.offset,
		Limit:  queryInfo.limit,
	}

	count, obligations, err := o.obligation.Get(context.Background(), &visitor, &obligationSearchInfo)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	resultTmp := make([]map[string]any, 0, len(obligations))
	for _, obligation := range obligations {
		resultTmp = append(resultTmp, map[string]any{
			"id":          obligation.ID,
			"type_id":     obligation.TypeID,
			"name":        obligation.Name,
			"description": obligation.Description,
			"value":       obligation.Value,
		})
	}

	rest.ReplyOK(c, http.StatusOK, map[string]any{
		"total_count": count,
		"entries":     resultTmp,
	})
}

func (o *obligationRestHandler) queryObligation(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	resourceTypeID := c.DefaultQuery("resource_type_id", "")
	if resourceTypeID == "" {
		err := gerrors.NewError(gerrors.PublicBadRequest, "param resource_type_id is required")
		rest.ReplyErrorV2(c, err)
		return
	}

	operationIDs := c.QueryArray("operation_ids")
	obligationTypeIDs := c.QueryArray("obligation_type_ids")

	queryInfo := interfaces.QueryObligationInfo{
		ResourceType:      resourceTypeID,
		Operation:         operationIDs,
		ObligationTypeIDs: obligationTypeIDs,
	}

	mapInfos, err := o.obligation.Query(context.Background(), &visitor, &queryInfo)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回结果
	results := make([]map[string]any, 0, len(mapInfos))
	for operationID, infos := range mapInfos {
		resultTmp := make([]map[string]any, 0, len(infos))
		for _, obligation := range infos {
			resultTmp = append(resultTmp, map[string]any{
				"id":          obligation.ID,
				"type_id":     obligation.TypeID,
				"name":        obligation.Name,
				"description": obligation.Description,
				"value":       obligation.Value,
			})
		}
		results = append(results, map[string]any{
			"operation_id": operationID,
			"obligations":  resultTmp,
		})
	}
	rest.ReplyOK(c, http.StatusOK, results)
}
