package common

import (
	"fmt"
	"strings"
)

// I18nMap map[int]map[string]string 别名
type I18nMap = map[int]map[string]string

// I18n 国际化
type I18n struct {
	m I18nMap
}

// NewI18n 实例化I18n
func NewI18n(m I18nMap) *I18n {
	return &I18n{m: m}
}

// Load 国际化
func (i *I18n) Load(id int, language string, a ...any) string {
	language = strings.ToLower(language)
	if len(a) > 0 {
		return fmt.Sprintf(i.m[id][language], a...)
	}
	return i.m[id][language]
}
