package localecommon

import "fmt"

// LangS 语言结构体，用于存储消息不同语言的文本
type LangS struct {
	Cn string `yaml:"cn"`
	Tw string `yaml:"tw"`
	En string `yaml:"en"`
}

// ToMap 将LangS转换为map，方便后续使用
func (l LangS) ToMap() map[LangType]string {
	m := make(map[LangType]string)
	m[LangTypeCn] = l.Cn
	m[LangTypeTw] = l.Tw
	m[LangTypeEn] = l.En
	return m
}

// LogMsg 日志消息结构体，用于存储日志消息
type LogMsg struct {
	Success *LangS `yaml:"success"`
	Fail    *LangS `yaml:"fail"`
	Ext     *LangS `yaml:"ext"`
}

// ToMap 将LogMsg转换为map，方便后续使用
func (l LogMsg) ToMap() map[SubKey]*LangS {
	m := make(map[SubKey]*LangS)
	m[DlSKSuccess] = l.Success
	m[DlSKFail] = l.Fail
	m[DlSKExt] = l.Ext
	return m
}

// BizKey 业务key，用于标识不同的业务
type BizKey string

// I18nMap 用于存储不同业务的日志消息
type I18nMap = map[BizKey]*LogMsg

// Msg 根据业务key和子key获取对应语言的消息
// 1. 如果没有对应的消息，则返回空字符串
// 2. 如果有对应的消息，但是没有对应的语言，则返回空字符串
// 3. 支持根据args替换消息中的占位符
func Msg(m I18nMap, key BizKey, subKey SubKey, lang string, args ...interface{}) (msg string) {
	logMsg := m[key]
	if logMsg == nil {
		return
	}

	msgMap := logMsg.ToMap()

	langMsg := msgMap[subKey]

	if langMsg == nil {
		return
	}

	var l LangType
	l.FromLang(lang)
	msg = langMsg.ToMap()[l]

	msg = fmt.Sprintf(msg, args...)
	return
}
