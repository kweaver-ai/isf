package apiserver

import (
	"io"
	"policy_mgnt/interfaces"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/ory/gojsonschema"
	"golang.org/x/text/language"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

// verify token有效性检查
func verify(c *gin.Context, hydra interfaces.Hydra) (visitor interfaces.Visitor, err error) {
	tokenID := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenID, "Bearer ")
	info, err := hydra.Introspect(token)
	if err != nil {
		return
	}

	if !info.Active {
		err = gerrors.NewError(gerrors.PublicUnauthorized, "token expired")
		return
	}

	visitor = interfaces.Visitor{
		ID:         info.VisitorID,
		TokenID:    tokenID,
		IP:         c.ClientIP(),
		Mac:        c.GetHeader("X-Request-MAC"),
		UserAgent:  c.GetHeader("User-Agent"),
		Type:       info.VisitorTyp,
		Language:   getXLang(c),
		ClientType: info.ClientTyp,
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

// ValidateAndBindGin 校验json数据
func validateAndBindGin(c *gin.Context, schema *gojsonschema.Schema, bind interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return validateAndBindNewError(body, schema, bind)
}

// validateAndBind 校验json数据
func validateAndBindNewError(body []byte, schema *gojsonschema.Schema, bind interface{}) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return gerrors.NewError(gerrors.PublicBadRequest, err.Error())
	}
	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return gerrors.NewError(gerrors.PublicBadRequest, strings.Join(msgList, "; "))
	}

	if err := jsoniter.Unmarshal(body, bind); err != nil {
		return err
	}

	return nil
}
