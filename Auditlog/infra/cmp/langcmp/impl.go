package langcmp

import (
	"strings"
	"sync"
)

var (
	langOnce sync.Once
	langImpl LangInterface
)

// langCmp 语言组件
type langCmp struct {
	isSetSysDefLang bool
	sysDefLang      Lang // 系统默认语言
}

var _ LangInterface = &langCmp{}

func NewLangCmp() LangInterface {
	langOnce.Do(func() {
		langImpl = &langCmp{}
	})

	return langImpl
}

func (p *langCmp) SetSysDefLang(lang string) {
	lowerL := Lang(strings.ToLower(lang))

	for i := range Languages {
		if Languages[i] == lowerL {
			p.sysDefLang = lowerL
			p.isSetSysDefLang = true

			return
		}
	}

	panic("[langCmp][SetSysDefLang]invalid lang")
}

// GetSysDefaultLang 获取系统默认语言
func (p *langCmp) GetSysDefaultLang() Lang {
	if !p.isSetSysDefLang || p.sysDefLang == "" {
		panic("[langCmp][GetSysDefaultLang]: language not set")
	}

	return p.sysDefLang
}
