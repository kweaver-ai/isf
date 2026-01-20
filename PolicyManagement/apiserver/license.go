package apiserver

import (
	_ "embed" // 标准用法
	"net/http"
	"policy_mgnt/common"
	"policy_mgnt/interfaces"
	"policy_mgnt/logics"
	"sync"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/ory/gojsonschema"
	"golang.org/x/text/language"
)

var (
	licenseSchemaOnce sync.Once
	lics              *licenseHandler

	langMatcher = language.NewMatcher([]language.Tag{
		language.SimplifiedChinese,
		language.TraditionalChinese,
		language.AmericanEnglish,
	})
	langMap = map[language.Tag]interfaces.Language{
		language.SimplifiedChinese:  interfaces.SimplifiedChinese,
		language.TraditionalChinese: interfaces.TraditionalChinese,
		language.AmericanEnglish:    interfaces.AmericanEnglish,
	}

	//go:embed jsonschema/license/get_authorized_product.json
	getAuthorizedProductsSchemaString string

	//go:embed jsonschema/license/update_authorized_product.json
	updateAuthorizedProductsSchemaString string
)

type licenseHandler struct {
	license                        interfaces.LogicsLicense
	hydra                          interfaces.Hydra
	getAuthorizedProductsSchema    *gojsonschema.Schema
	updateAuthorizedProductsSchema *gojsonschema.Schema
	mapObjectTypeToString          map[interfaces.ObjectType]string
	mapStringToObjectType          map[string]interfaces.ObjectType
}

func newLicenseHandler() (*licenseHandler, error) {
	licenseSchemaOnce.Do(func() {
		getAuthorizedProductsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getAuthorizedProductsSchemaString))
		if err != nil {
			panic(err)
		}

		updateAuthorizedProductsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(updateAuthorizedProductsSchemaString))
		if err != nil {
			panic(err)
		}

		lics = &licenseHandler{
			license:                     logics.NewLicense(),
			hydra:                       newHydra(),
			getAuthorizedProductsSchema: getAuthorizedProductsSchema,
			mapObjectTypeToString: map[interfaces.ObjectType]string{
				interfaces.ObjectTypeUser: "user",
			},
			updateAuthorizedProductsSchema: updateAuthorizedProductsSchema,
			mapStringToObjectType: map[string]interfaces.ObjectType{
				"user": interfaces.ObjectTypeUser,
			},
		}
	})
	return lics, nil
}

func (h *licenseHandler) AddRouters(r *gin.RouterGroup) {
	r.GET("/console/licenses", observable.MiddlewareTrace(common.SvcARTrace), h.getLicenses)
	r.POST("/console/query-authorized-products", observable.MiddlewareTrace(common.SvcARTrace), h.getAuthorizedProducts)
	r.GET("/check-product-authorized", observable.MiddlewareTrace(common.SvcARTrace), h.checkProductAuthorized)
	r.PUT("/console/authorized-products", observable.MiddlewareTrace(common.SvcARTrace), h.updateAuthorizedProducts)
}

func (h *licenseHandler) updateAuthorizedProducts(c *gin.Context) {
	// 访问者信息校验
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyErrorV2(c, vErr)
		return
	}

	// jsonschema校验
	var jsonReq []interface{}
	if err := validateAndBindGin(c, h.updateAuthorizedProductsSchema, &jsonReq); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 数据整理
	productsMap := make(map[string]interfaces.AuthorizedProduct)
	for _, item := range jsonReq {
		temp := item.(map[string]interface{})
		authorizedProduct := interfaces.AuthorizedProduct{
			ID:      temp["id"].(string),
			Type:    h.mapStringToObjectType[temp["type"].(string)],
			Product: make([]string, 0),
		}

		for _, product := range temp["products"].([]interface{}) {
			authorizedProduct.Product = append(authorizedProduct.Product, product.(string))
		}
		productsMap[authorizedProduct.ID] = authorizedProduct
	}

	// 更新已授权产品
	err := h.license.UpdateAuthorizedProducts(c, &visitor, productsMap)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (h *licenseHandler) checkProductAuthorized(c *gin.Context) {
	// 访问者信息校验
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyErrorV2(c, vErr)
		return
	}

	product := c.Query("product")
	if product == "" {
		rest.ReplyErrorV2(c, gerrors.NewError(gerrors.PublicBadRequest, "product is required"))
		return
	}

	authorized, unauthorized_reason, err := h.license.CheckProductAuthorized(c, &visitor, product)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回结果
	out := map[string]interface{}{
		"authorized": authorized,
	}
	if !authorized {
		out["unauthorized_reason"] = unauthorized_reason
	}
	rest.ReplyOK(c, http.StatusOK, out)

}

func (h *licenseHandler) getAuthorizedProducts(c *gin.Context) {
	// 访问者信息校验
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyErrorV2(c, vErr)
		return
	}

	// jsonschema校验
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, h.getAuthorizedProductsSchema, &jsonReq); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	userIds := make([]string, 0)
	tempUserIds := jsonReq["user_ids"].([]interface{})
	for _, userId := range tempUserIds {
		userIds = append(userIds, userId.(string))
	}

	products, err := h.license.GetAuthorizedProducts(c, &visitor, userIds)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	out := make([]interface{}, 0, len(products))
	for _, v := range products {
		out = append(out, map[string]interface{}{
			"id":       v.ID,
			"type":     h.mapObjectTypeToString[v.Type],
			"products": v.Product,
		})
	}

	// 整合数据
	rest.ReplyOK(c, http.StatusOK, out)
}

func (h *licenseHandler) getLicenses(c *gin.Context) {
	// 访问者信息校验
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyErrorV2(c, vErr)
		return
	}

	// 获取信息
	infos, err := h.license.GetLicenses(c, &visitor)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 整合数据
	out := make([]interface{}, 0, len(infos))
	for i := range infos {
		temp := map[string]interface{}{
			"product":               infos[i].Product,
			"total_user_quota":      infos[i].TotalUserQuota,
			"authorized_user_count": infos[i].AuthorizedUserCount,
		}
		out = append(out, temp)
	}
	rest.ReplyOK(c, http.StatusOK, out)
}
