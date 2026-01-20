package driveradapters

import (
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"Authentication/interfaces"
)

var (
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

// GetXLang 解析获取 Header x-language
func GetXLang(c *gin.Context) interfaces.Language {
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
