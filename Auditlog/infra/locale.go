package infra

import (
	"fmt"
	"strings"
	"sync"

	"AuditLog/common"
)

var (
	localeOnce sync.Once
	textMap    = make(map[string]string)
)

func TextDomain(format string, a ...interface{}) (s string) {
	localeOnce.Do(initTextMap)

	s = fmt.Sprintf(textMap[format], a...)

	return
}

func initTextMap() {
	lang := strings.ToLower(common.SvcConfig.Languaue)

	switch lang {
	case "zh-cn", "zh_cn":
		textMap["IDS_UNDISTRIBUTED_GROUP"] = "未分配组"

	case "zh-tw", "zh_tw":
		textMap["IDS_UNDISTRIBUTED_GROUP"] = "未分配組"

	case "en-us", "en_us":
		textMap["IDS_UNDISTRIBUTED_GROUP"] = "Unassigned Group"
	default:
		fmt.Println("invalid language")
	}
}
