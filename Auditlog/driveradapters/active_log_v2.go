package driveradapters

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	gerror "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/language"

	"AuditLog/common"
	"AuditLog/common/conf"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/models"
)

var (
	activeLogV2Once sync.Once
	activeLogV2     interfaces.PublicRESTHandler

	//go:embed jsonschema/audit_log_schema.json
	auditLogSchemaStr string

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
)

type activeLogV2Handler struct {
	auditLogSchema *gojsonschema.Schema
	hydra          interfaces.Hydra
	logMgnt        interfaces.LogMgnt
}

// NewActiveLogV2Handler 创建document handler对象
func NewActiveLogV2Handler() interfaces.PublicRESTHandler {
	auditLogSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditLogSchemaStr))
	if err != nil {
		panic(err)
	}

	activeLogV2Once.Do(func() {
		activeLogV2 = &activeLogV2Handler{
			auditLogSchema: auditLogSchema,
			hydra:          newHydra(),
			logMgnt:        logics.NewLogMgnt(),
		}
	})

	return activeLogV2
}

// RegisterPublic 注册外部API
func (ac *activeLogV2Handler) RegisterPublic(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/logs", ac.addLog)
}

// addLog 添加活跃日志
func (ac *activeLogV2Handler) addLog(c *gin.Context) {
	visitor, err := verify(c, ac.hydra)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}
	// 检查json格式
	var jsonReq map[string]interface{}
	err = ac.validateAndBindGin(c, ac.auditLogSchema, &jsonReq)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 检查operation
	logInfo := models.AuditLog{}
	strLogType := jsonReq["type"].(string)
	logType := common.MapLogType[strLogType]
	bExistOpType := false
	switch logType {
	case interfaces.LogType_Login:
		if data, ok := conf.MapLoginOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			logInfo.OpType = data
			bExistOpType = true
		}
	case interfaces.LogType_Management:
		if data, ok := conf.MapManageOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			logInfo.OpType = data
			bExistOpType = true
		}
	case interfaces.LogType_Operation:
		if data, ok := conf.MapOperOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			logInfo.OpType = data
			bExistOpType = true
		}
	}

	if !bExistOpType {
		err := gerror.NewError(gerror.PublicBadRequest, "operation string is not enums")
		rest.ReplyErrorV2(c, err)
		return
	}

	logInfo.Msg = jsonReq["description"].(string)
	logInfo.Date = int64(jsonReq["op_time"].(float64) / 1000)

	// 获取object
	if object, ok := jsonReq["object"]; ok {
		objData := object.(map[string]interface{})
		if logInfo.ObjType, ok = conf.MapObjectTypeStrToint[objData["type"].(string)]; !ok {
			err := gerror.NewError(gerror.PublicBadRequest, "object type is not enums")
			rest.ReplyErrorV2(c, err)
			return
		}

		if _, ok := objData["id"]; ok {
			logInfo.ObjID = objData["id"].(string)
		}
		if _, ok := objData["name"]; ok {
			logInfo.ObjName = objData["name"].(string)
		}
	}

	// 获取日志其他信息
	if exMsg, ok := jsonReq["ex_msg"]; ok {
		logInfo.Exmsg = exMsg.(string)
	}
	logInfo.Level = common.MapLevelType[jsonReq["level"].(string)]
	logInfo.OutBizID = jsonReq["out_biz_id"].(string)

	if detail, ok := jsonReq["detail"]; ok {
		detailBytes, err := jsoniter.Marshal(detail)
		if err != nil {
			rest.ReplyErrorV2(c, err)
			return
		}
		logInfo.AdditionalInfo = string(detailBytes)
	}

	info := &models.ReceiveLogVo{
		Language:   common.SvcConfig.Languaue,
		LogType:    strLogType,
		LogContent: &logInfo,
	}

	// 记录日志信息
	logID, err := ac.logMgnt.AddAuditLog(visitor, logType, info)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/audit_log/v1/logs/%s", logID))
	rest.ReplyOK(c, http.StatusCreated, gin.H{"id": logID})
}

// ValidateAndBindGin 校验json数据
func (ac *activeLogV2Handler) validateAndBindGin(c *gin.Context, schema *gojsonschema.Schema, bind interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return ac.validateAndBind(body, schema, bind)
}

// validateAndBind 校验json数据
func (ac *activeLogV2Handler) validateAndBind(body []byte, schema *gojsonschema.Schema, bind interface{}) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return gerror.NewError(gerror.PublicBadRequest, err.Error())
	}
	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return gerror.NewError(gerror.PublicBadRequest, strings.Join(msgList, "; "))
	}

	if err := jsoniter.Unmarshal(body, bind); err != nil {
		return err
	}

	return nil
}

// verify token有效性检查
func verify(c *gin.Context, hydra interfaces.Hydra) (visitor interfaces.Visitor, err error) {
	tokenID := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenID, "Bearer ")
	info, err := hydra.Introspect(token)
	if err != nil {
		return
	}

	if !info.Active {
		err = gerror.NewError(gerror.PublicUnauthorized, "token expired")
		return
	}

	visitor = interfaces.Visitor{
		ID:        info.VisitorID,
		TokenID:   tokenID,
		IP:        c.ClientIP(),
		Mac:       c.GetHeader("X-Request-MAC"),
		UserAgent: c.GetHeader("User-Agent"),
		Type:      info.VisitorTyp,
		Language:  getXLang(c),
	}

	return
}

// getXLang 解析获取 Header x-language
func getXLang(c *gin.Context) interfaces.Language {
	tag, _ := language.MatchStrings(langMatcher, getBCP47(c.GetHeader("x-language")))
	return langMap[tag]
}

// getBCP47 将约定的语言标签转换为符合BCP47标准的语言标签
// 默认值为 zh-Hans, 中国大陆简体中文
// https://www.rfc-editor.org/info/bcp47
func getBCP47(s string) string {
	switch strings.ToLower(s) {
	case "zh_cn", "zh-cn":
		return "zh-Hans"
	case "zh_tw", "zh-tw":
		return "zh-Hant"
	case "en_us", "en-us":
		return "en-US"
	default:
		return "zh-Hans"
	}
}
