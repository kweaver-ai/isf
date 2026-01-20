package langcmp

import "strings"

type Lang string

const (
	ZhCN Lang = "zh_cn" // zh_CN
	ZhTW Lang = "zh_tw" // zh_TW
	En   Lang = "en_us" // en_US
)

// Languages 支持的语言
var Languages = [3]Lang{ZhCN, ZhTW, En}

// Check 检查语言是否支持
func (l *Lang) Check() (ok bool) {
	lowerL := Lang(strings.ToLower(string(*l)))
	for i := range Languages {
		if Languages[i] == lowerL {
			ok = true
			return
		}
	}

	return
}

func (l *Lang) ToLower() {
	*l = Lang(strings.ToLower(string(*l)))
}

func (l *Lang) ToString() string {
	return string(*l)
}

func NewFromStr(lang string) (l Lang) {
	if lang == "" {
		return
	}

	lang = strings.ToLower(lang)
	switch lang {
	case "zh-cn", "zh_cn":
		l = ZhCN
	case "zh-tw", "zh_tw":
		l = ZhTW
	case "en-us", "en_us":
		l = En
	}

	return
}
