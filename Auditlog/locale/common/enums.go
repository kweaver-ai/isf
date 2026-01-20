package localecommon

import "strings"

type SubKey string

const (
	DlSKSuccess SubKey = "success"
	DlSKFail    SubKey = "fail"
	DlSKExt     SubKey = "ext"
)

type LangType string

func (l *LangType) FromLang(lang string) {
	lang = strings.ToLower(lang)
	switch lang {
	case "zh-cn", "zh_cn":
		*l = LangTypeCn
	case "zh-tw", "zh_tw":
		*l = LangTypeTw
	case "en-us", "en_us":
		*l = LangTypeEn
	default:
		panic("LangType: invalid language")
	}
}

// FromLang 从lang转换为LangType，兼容大小写和下划线
func FromLang(lang string) (l LangType) {
	lang = strings.ToLower(lang)
	lang = strings.ReplaceAll(lang, "_", "-")

	switch lang {
	case "zh-cn":
		l = LangTypeCn
	case "zh-tw":
		l = LangTypeTw
	case "en-us":
		l = LangTypeEn
	default:
		panic("LangType: invalid language")
	}

	return
}

const (
	LangTypeCn LangType = "cn"
	LangTypeTw LangType = "tw"
	LangTypeEn LangType = "en"
)

type LangLower = string

const (
	LangLowerCn LangLower = "zh-cn"
	LangLowerTw LangLower = "zh-tw"
	LangLowerEn LangLower = "en-us"
)
