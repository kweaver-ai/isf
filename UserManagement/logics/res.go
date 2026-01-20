// Package logics 翻译文件
package logics

import (
	"sync"

	"UserManagement/common"
	"UserManagement/interfaces"
)

const (
	strZHCN = "zh_CN"
	strZHTW = "zh_TW"
	strENUS = "en_US"
)

var (
	resOnce sync.Once
	resMap  = make(map[interfaces.LangType]map[string]string)
)

// loadString 获取字符串型语言项
func loadString(langType interfaces.LangType, key string) string {
	resOnce.Do(initResMap)

	// 如果语言存在问题，则返回空
	if _, ok := resMap[langType]; !ok {
		common.NewLogger().Errorf("api language set error :%s", langType)
		return ""
	}

	return resMap[langType][key]
}

func initResMap() {
	zhCHMap := make(map[string]string)
	zhCHMap["IDS_DEFAULT_PWD_INVALID"] = "密码只能包含英文或数字或~!%#$@-_.字符，长度范围6~100个字符，请重新输入。"

	zhTWMap := make(map[string]string)
	zhTWMap["IDS_DEFAULT_PWD_INVALID"] = "密碼只能包含英文或數字或~!%#$@-_.字元，長度範圍6~100個字元，請重新輸入。"

	enUSMap := make(map[string]string)
	enUSMap["IDS_DEFAULT_PWD_INVALID"] = "The password should be letters, numbers or ~!%#$@-_. within 6 ~ 100 characters, please re-enter."

	resMap[interfaces.LTZHCN] = zhCHMap
	resMap[interfaces.LTZHTW] = zhTWMap
	resMap[interfaces.LTENUS] = enUSMap

	switch common.SvcConfig.Lang {
	case strENUS, strZHCN, strZHTW:
	default:
		// 如果后端服务国际化设置错误，则报错
		common.NewLogger().Fatalln("service language not set")
	}
}
