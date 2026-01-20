// Package driveradapters AnyShare 入站适配器
package driveradapters

import (
	"context"
	_ "embed" // 标准用法
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"

	"Authorization/interfaces"
	"Authorization/logics"
)

//go:embed jsonschema/obligation_type/set.json
var setObligationTypeSchemaStr string

var (
	obligationTypeOnce    sync.Once
	obligationTypeHandler RestHandler
)

type obligationTypeRestHandler struct {
	obligationType interfaces.ObligationType
	hydra          interfaces.Hydra
	setSchemaStr   *gojsonschema.Schema
}

// NewObligationTemplateRestHandler 权限适配器接口
func NewObligationTemplateRestHandler() RestHandler {
	obligationTypeOnce.Do(func() {
		obligationTypeHandler = &obligationTypeRestHandler{
			obligationType: logics.NewObligationType(),
			hydra:          newHydra(),
			setSchemaStr:   newJSONSchema(setObligationTypeSchemaStr),
		}
	})
	return obligationTypeHandler
}

// RegisterPrivate 注册内部API
func (o *obligationTypeRestHandler) RegisterPrivate(_ *gin.Engine) {
}

// RegisterPublic 注册外部API
func (o *obligationTypeRestHandler) RegisterPublic(engine *gin.Engine) {
	// 管理接口
	engine.PUT("/api/authorization/v1/obligation-types/:id", o.set)
	engine.GET("/api/authorization/v1/obligation-types/:id", o.getByID)
	engine.DELETE("/api/authorization/v1/obligation-types/:id", o.delete)
	engine.GET("/api/authorization/v1/obligation-types", o.get)

	// 查询接口
	// 义务类型查询
	engine.GET("/api/authorization/v1/query-obligation-types", o.queryObligationTypes)
}

func (o *obligationTypeRestHandler) set(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	ID := c.Param("id")
	if ID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	var obligationTypeDr map[string]any
	if err = validateAndBindGin(c, o.setSchemaStr, &obligationTypeDr); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationType := interfaces.ObligationTypeInfo{
		ID:     ID,
		Name:   obligationTypeDr["name"].(string),
		Schema: obligationTypeDr["schema"],
	}

	defJson, ok := obligationTypeDr["default_value"]
	if ok {
		obligationType.DefaultValue = defJson
	}

	descriptionJson, ok := obligationTypeDr["description"]
	if ok {
		obligationType.Description = descriptionJson.(string)
	}

	uiSchemaJson, ok := obligationTypeDr["ui_schema"]
	if ok {
		obligationType.UiSchema = uiSchemaJson
	}

	resourceTypeScopesJson := obligationTypeDr["applicable_resource_types"].(map[string]any)
	obligationType.ResourceTypeScope.Unlimited = resourceTypeScopesJson["unlimited"].(bool)

	// 如果资源类型有范围限制，则需要设置资源类型范围
	if !obligationType.ResourceTypeScope.Unlimited {
		_, resourceTypesExist := resourceTypeScopesJson["resource_types"]
		if !resourceTypesExist {
			rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "resource_types is required"))
			return
		}
		resourceTypesJson := resourceTypeScopesJson["resource_types"].([]any)
		for _, resourceTypeJson := range resourceTypesJson {
			// 遍历资源类型，每个资源类型信息 放入 resourceTypeScope
			var resourceTypeScope interfaces.ObligationResourceTypeScope
			resourceTypeJsonMap := resourceTypeJson.(map[string]any)
			resourceTypeID := resourceTypeJsonMap["id"].(string)
			operationsScopeJson := resourceTypeJsonMap["applicable_operations"].(map[string]any)
			var operationsScopeInfo interfaces.ObligationOperationsScopeInfo
			operationsScopeInfo.Unlimited = operationsScopeJson["unlimited"].(bool)
			if !operationsScopeInfo.Unlimited {
				_, operationsExist := operationsScopeJson["operations"]
				if !operationsExist {
					rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "operations is required"))
					return
				}
				// 获取资源类型上的操作
				operationsJson := operationsScopeJson["operations"].([]any)
				for _, operationJson := range operationsJson {
					operationJsonMap := operationJson.(map[string]any)
					operationID := operationJsonMap["id"].(string)
					var operation interfaces.ObligationOperation
					operation.ID = operationID
					operationsScopeInfo.Operations = append(operationsScopeInfo.Operations, operation)
				}
			}
			resourceTypeScope.ResourceTypeID = resourceTypeID
			resourceTypeScope.OperationsScope = operationsScopeInfo
			// 将资源类型信息放入资源类型范围
			obligationType.ResourceTypeScope.Types = append(obligationType.ResourceTypeScope.Types, resourceTypeScope)
		}
	}

	err = o.obligationType.Set(context.Background(), &visitor, &obligationType)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *obligationTypeRestHandler) delete(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationTypeID := c.Param("id")
	if obligationTypeID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	err = o.obligationType.Delete(context.Background(), &visitor, obligationTypeID)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (o *obligationTypeRestHandler) get(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取列举用户组信息
	queryInfo, err := getListQueryParam(c)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationType := interfaces.ObligationTypeSearchInfo{
		Offset: queryInfo.offset,
		Limit:  queryInfo.limit,
	}

	count, infos, err := o.obligationType.Get(context.Background(), &visitor, &obligationType)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 按照这个返回结果 拼一下 返回结果
	resultTmp := make([]map[string]any, 0, len(infos))
	for i := range infos {
		resourceTypeScope := o.resourceTypeScopeInfoToString(infos[i].ResourceTypeScope)

		var def any
		if infos[i].DefaultValue != nil {
			def = infos[i].DefaultValue
		} else {
			def = map[string]any{}
		}

		// 解析UiSchema
		var uiSchema any
		if infos[i].UiSchema != nil {
			uiSchema = infos[i].UiSchema
		} else {
			uiSchema = map[string]any{}
		}
		resultTmp = append(resultTmp, map[string]any{
			"id":                        infos[i].ID,
			"name":                      infos[i].Name,
			"description":               infos[i].Description,
			"schema":                    infos[i].Schema,
			"default_value":             def,
			"ui_schema":                 uiSchema,
			"applicable_resource_types": resourceTypeScope,
		})
	}

	rest.ReplyOK(c, http.StatusOK, map[string]any{
		"total_count": count,
		"entries":     resultTmp,
	})
}

func (o *obligationTypeRestHandler) getByID(c *gin.Context) {
	visitor, err := verify(c, o.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	obligationTypeID := c.Param("id")
	if obligationTypeID == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "id is required"))
		return
	}

	info, err := o.obligationType.GetByID(context.Background(), &visitor, obligationTypeID)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	resourceTypeScope := o.resourceTypeScopeInfoToString(info.ResourceTypeScope)

	var def any
	if info.DefaultValue != nil {
		def = info.DefaultValue
	} else {
		def = map[string]any{}
	}

	// 解析UiSchema
	var uiSchema any
	if info.UiSchema != nil {
		uiSchema = info.UiSchema
	} else {
		uiSchema = map[string]any{}
	}
	result := map[string]any{
		"id":                        info.ID,
		"name":                      info.Name,
		"description":               info.Description,
		"schema":                    info.Schema,
		"default_value":             def,
		"ui_schema":                 uiSchema,
		"applicable_resource_types": resourceTypeScope,
	}
	rest.ReplyOK(c, http.StatusOK, result)
}

func (o *obligationTypeRestHandler) resourceTypeScopeInfoToString(info interfaces.ObligationResourceTypeScopeInfo) (result map[string]any) {
	result = make(map[string]any)
	// 如果不限制
	if info.Unlimited {
		result["unlimited"] = true
		result["resource_types"] = make([]any, 0)
	} else {
		result["unlimited"] = false
		resourceTypes := make([]any, 0, len(info.Types))
		for _, rt := range info.Types {
			resourceType := make(map[string]any)
			resourceType["id"] = rt.ResourceTypeID
			resourceType["name"] = rt.ResourceTypeName
			OperationsScope := make(map[string]any)
			if rt.OperationsScope.Unlimited {
				OperationsScope["unlimited"] = true
				OperationsScope["operations"] = make([]any, 0)
			} else {
				OperationsScope["unlimited"] = false
				var operations []any
				for _, op := range rt.OperationsScope.Operations {
					operations = append(operations, map[string]any{
						"id":   op.ID,
						"name": op.Name,
					})
				}
				OperationsScope["operations"] = operations
			}
			resourceType["applicable_operations"] = OperationsScope
			resourceTypes = append(resourceTypes, resourceType)
		}
		result["resource_types"] = resourceTypes
	}
	return result
}

func (o *obligationTypeRestHandler) queryObligationTypes(c *gin.Context) {
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
	queryInfo := interfaces.QueryObligationTypeInfo{
		ResourceType: resourceTypeID,
		Operation:    operationIDs,
	}

	mapInfos, err := o.obligationType.Query(context.Background(), &visitor, &queryInfo)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回结果
	results := make([]map[string]any, 0, len(mapInfos))
	for operationID, infos := range mapInfos {
		resultTmps := make([]map[string]any, 0, len(infos))
		for i := range infos {
			var def any
			if infos[i].DefaultValue != nil {
				def = infos[i].DefaultValue
			} else {
				def = map[string]any{}
			}
			// 解析UiSchema
			var uiSchema any
			if infos[i].UiSchema != nil {
				uiSchema = infos[i].UiSchema
			} else {
				uiSchema = map[string]any{}
			}
			resultTmps = append(resultTmps, map[string]any{
				"id":            infos[i].ID,
				"name":          infos[i].Name,
				"description":   infos[i].Description,
				"schema":        infos[i].Schema,
				"default_value": def,
				"ui_schema":     uiSchema,
			})
		}

		results = append(results, map[string]any{
			"operation_id":     operationID,
			"obligation_types": resultTmps,
		})
	}
	rest.ReplyOK(c, http.StatusOK, results)
}
