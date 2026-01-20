// Package driveradapters config AnyShare  用户管理配置接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/language"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

// ConfigRestHandler RESTful api Handler接口
type ConfigRestHandler interface {
	// RegisterPrivate 注册开放API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}

type configRestHandler struct {
	config            interfaces.LogicsConfig
	strDefaultUserPWd string
	strCSFLevelEnum   string
	strCSFLevel2Enum  string
	strShowCSFLevel2  string
	hydra             interfaces.Hydra
	mapLangs          map[string]interfaces.LangType
	setConfigSchema   *gojsonschema.Schema
}

var (
	configonce    sync.Once
	configHandler ConfigRestHandler

	//go:embed jsonschema/config/set_config.json
	setConfigSchemaStr string
)

// NewConfigRESTHandler 创建配置操作对象
func NewConfigRESTHandler() ConfigRestHandler {
	configonce.Do(func() {
		setConfigSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(setConfigSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		configHandler = &configRestHandler{
			config:            logics.NewConfig(),
			strDefaultUserPWd: "default_user_pwd",
			strCSFLevelEnum:   "csf_level_enum",
			strCSFLevel2Enum:  "csf_level2_enum",
			strShowCSFLevel2:  "show_csf_level2",
			hydra:             newHydra(),
			mapLangs: map[string]interfaces.LangType{
				"zh_CN": interfaces.LTZHCN,
				"zh_TW": interfaces.LTZHTW,
				"en_US": interfaces.LTENUS,
			},
			setConfigSchema: setConfigSchema,
		}
	})

	return configHandler
}

// RegisterPrivate 注册内部API
func (con *configRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.PUT("/api/user-management/v1/configs/:fields", con.updateConfig)
	engine.GET("/api/user-management/v1/configs/:fields", con.getConfigsPrivate)
}

// RegisterPublic 注册外部API
func (con *configRestHandler) RegisterPublic(engine *gin.Engine) {
	engine.PUT("/api/user-management/v1/management/configs/:fields", con.updateManageConfig)
	engine.GET("/api/user-management/v1/management/default-pwd-valid", con.checkPWDValid)
	engine.GET("/api/user-management/v1/configs/:fields", con.getConfigs)
}

func (con *configRestHandler) getConfigsPrivate(c *gin.Context) {
	fields := strings.Split(c.Param("fields"), ",")
	rg := make(map[interfaces.ConfigKey]bool)
	for _, v := range fields {
		if v == con.strCSFLevelEnum {
			rg[interfaces.CSFLevelEnum] = true
		} else if v == con.strCSFLevel2Enum {
			rg[interfaces.CSFLevel2Enum] = true
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "invalid fields type")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	config, err := con.config.GetConfig(rg)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回值处理，按照value的值从小到大排序
	outData := make(map[string]interface{})
	if rg[interfaces.CSFLevelEnum] {
		temp := make([]map[string]interface{}, 0)
		tempValues := make([]int, 0)
		tempNames := make(map[int]string, 0)
		for k, v := range config.CSFLevelEnum {
			tempValues = append(tempValues, v)
			tempNames[v] = k
		}
		sort.Ints(tempValues)
		for _, v := range tempValues {
			temp = append(temp, map[string]interface{}{
				"name":  tempNames[v],
				"value": v,
			})
		}
		outData["csf_level_enum"] = temp
	}
	if rg[interfaces.CSFLevel2Enum] {
		temp := make([]map[string]interface{}, 0)
		tempValues := make([]int, 0)
		tempNames := make(map[int]string, 0)
		for k, v := range config.CSFLevel2Enum {
			tempValues = append(tempValues, v)
			tempNames[v] = k
		}
		sort.Ints(tempValues)
		for _, v := range tempValues {
			temp = append(temp, map[string]interface{}{
				"name":  tempNames[v],
				"value": v,
			})
		}
		outData["csf_level2_enum"] = temp
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

func (con *configRestHandler) getConfigs(c *gin.Context) {
	// token验证
	_, vErr := verify(c, con.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	fields := strings.Split(c.Param("fields"), ",")
	rg := make(map[interfaces.ConfigKey]bool)
	for _, v := range fields {
		if v == con.strCSFLevelEnum {
			rg[interfaces.CSFLevelEnum] = true
		} else if v == con.strCSFLevel2Enum {
			rg[interfaces.CSFLevel2Enum] = true
		} else if v == con.strShowCSFLevel2 {
			rg[interfaces.ShowCSFLevel2] = true
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "invalid fields type")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	config, err := con.config.GetConfig(rg)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回值处理，按照value的值从小到大排序
	outData := make(map[string]interface{})
	if rg[interfaces.CSFLevelEnum] {
		temp := make([]map[string]interface{}, 0)
		tempValues := make([]int, 0)
		tempNames := make(map[int]string, 0)
		for k, v := range config.CSFLevelEnum {
			tempValues = append(tempValues, v)
			tempNames[v] = k
		}
		sort.Ints(tempValues)
		for _, v := range tempValues {
			temp = append(temp, map[string]interface{}{
				"name":  tempNames[v],
				"value": v,
			})
		}
		outData["csf_level_enum"] = temp
	}
	if rg[interfaces.CSFLevel2Enum] {
		temp := make([]map[string]interface{}, 0)
		tempValues := make([]int, 0)
		tempNames := make(map[int]string, 0)
		for k, v := range config.CSFLevel2Enum {
			tempValues = append(tempValues, v)
			tempNames[v] = k
		}
		sort.Ints(tempValues)
		for _, v := range tempValues {
			temp = append(temp, map[string]interface{}{
				"name":  tempNames[v],
				"value": v,
			})
		}
		outData["csf_level2_enum"] = temp
	}

	if rg[interfaces.ShowCSFLevel2] {
		outData["show_csf_level2"] = config.ShowCSFLevel2
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// checkPWDValid 检查密码可用性
func (con *configRestHandler) checkPWDValid(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, con.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 检查语言
	ok := false
	lang := con.GetLanguage(c.GetHeader("x-language"))
	visitor.LangType, ok = con.mapLangs[lang]
	if !ok {
		err := rest.NewHTTPError("invalid x-language", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 检查参数
	password, ok := c.GetQuery("password")
	if !ok {
		err := rest.NewHTTPError("invalid password", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 检查密码
	result, msg, err := con.config.CheckDefaultPWD(&visitor, password)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 返回值处理
	outData := make(map[string]interface{})
	outData["result"] = result
	if !result {
		outData["err_msg"] = msg
	}
	rest.ReplyOK(c, http.StatusOK, outData)
}

// updateConfig 更新涉密配置
func (con *configRestHandler) updateManageConfig(c *gin.Context) {
	// token验证
	visitor, vErr := verifyNewError(c, con.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	con.updateConfigCommon(&visitor, c)
}

// updateConfig 更新涉密配置
func (con *configRestHandler) updateConfig(c *gin.Context) {
	con.updateConfigCommon(nil, c)
}

// updateConfig 更新涉密配置
func (con *configRestHandler) updateConfigCommon(visitor *interfaces.Visitor, c *gin.Context) {
	fields := strings.Split(c.Param("fields"), ",")

	// 检查url参数是否合法, 内外部接口支持设置默认密码，只有外部接口才支持密级枚举设置
	rg := make(map[interfaces.ConfigKey]bool)
	for _, v := range fields {
		if v == con.strDefaultUserPWd {
			rg[interfaces.UserDefaultPWD] = true
		} else if v == con.strCSFLevelEnum && visitor != nil {
			rg[interfaces.CSFLevelEnum] = true
		} else if v == con.strCSFLevel2Enum && visitor != nil {
			rg[interfaces.CSFLevel2Enum] = true
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "invalid fields type")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	// jsonschema校验
	var jsonReq map[string]interface{}
	var err error
	if err = validateAndBindGinNewError(c, con.setConfigSchema, &jsonReq); err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取参数，并且检查fields的有效性
	configInfo := interfaces.Config{}
	if _, ok := rg[interfaces.UserDefaultPWD]; ok {
		if value, ok := jsonReq[con.strDefaultUserPWd]; ok {
			configInfo.UserDefaultPWD = value.(string)
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "body.default_user_pwd is required")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	if _, ok := rg[interfaces.CSFLevelEnum]; ok {
		if value, ok := jsonReq[con.strCSFLevelEnum]; ok {
			configInfo.CSFLevelEnum, err = con.checkCSFLevelEnum(value.([]interface{}))
			if err != nil {
				rest.ReplyErrorV2(c, err)
				return
			}
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "body.csf_level_enum is required")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	if _, ok := rg[interfaces.CSFLevel2Enum]; ok {
		if value, ok := jsonReq[con.strCSFLevel2Enum]; ok {
			configInfo.CSFLevel2Enum, err = con.checkCSFLevelEnum(value.([]interface{}))
			if err != nil {
				rest.ReplyErrorV2(c, err)
				return
			}
		} else {
			err := gerrors.NewError(gerrors.PublicBadRequest, "body.csf_level2_enum is required")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	// 获取信息
	err = con.config.UpdateConfig(visitor, rg, &configInfo)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 返回信息
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (con *configRestHandler) GetLanguage(xlang string) string {
	matcher := language.NewMatcher([]language.Tag{
		language.MustParse(common.SvcConfig.Lang),
		language.MustParse("zh_CN"), // 默认语言放第一位，均未匹配时返回此项
		language.MustParse("zh_TW"),
		language.MustParse("en_US"),
	})
	tag, _ := language.MatchStrings(matcher, xlang)
	base, _ := tag.Base()
	reg, _ := tag.Region()
	return fmt.Sprintf("%s_%s", base, reg)
}

// checkCSFLevelEnum 检查密级枚举是否合法， value范围限制+name value不允许重复
/*
"csf_level_enum": [
{
	"name": "公开",
	"value": 5
},
{
	"name": "秘密",
	"value": 6
},
{
	"name": "机密",
	"value": 7
},
{
	"name": "绝密",
	"value": 8
}
],*/
func (con *configRestHandler) checkCSFLevelEnum(csfLevelEnum []interface{}) (csfLevelEnumMap map[string]int, err error) {
	csfLevelEnumMap = make(map[string]int)
	tempValues := make(map[int]bool)
	tempNames := make(map[string]bool)
	for _, v := range csfLevelEnum {
		tempName := v.(map[string]interface{})["name"].(string)
		tempValue := int(v.(map[string]interface{})["value"].(float64))

		// 不允许重复设置密级枚举
		if _, ok := tempNames[tempName]; ok {
			return nil, gerrors.NewError(gerrors.PublicBadRequest, "csf level enum name is not unique")
		}
		if _, ok := tempValues[tempValue]; ok {
			return nil, gerrors.NewError(gerrors.PublicBadRequest, "csf level enum value is not unique")
		}
		tempNames[tempName] = true
		tempValues[tempValue] = true
		csfLevelEnumMap[tempName] = tempValue
	}
	return csfLevelEnumMap, nil
}
