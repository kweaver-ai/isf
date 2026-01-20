// Package drivenadapters 翻译文件
package drivenadapters

import (
	"fmt"
	"sync"

	"Authentication/common"
)

var (
	resOnce sync.Once
	resMap  = make(map[string]string)
)

// loadString 获取字符串型语言项
// 如果传入了args, args仅为 format string的value
// 与resMap[key]返回的字符串是对应的
func loadString(key string, args ...any) string {
	resOnce.Do(initResMap)
	if len(args) > 0 {
		return fmt.Sprintf(resMap[key], args...)
	}
	return resMap[key]
}

const (
	simpleChinese  = "zh_CN"
	complexChinese = "zh_TW"
	english        = "en_US"
)

func initResMap() {
	initAccessTokenPermResMap()
}

func initAccessTokenPermResMap() {
	switch common.SvcConfig.Lang {
	case simpleChinese:
		resMap["IDS_SET_ACCESSTOKENPERM"] = "设置应用账户“%s”用户访问令牌权限 成功"
		resMap["IDS_DELETE_ACCESSTOKENPERM"] = "删除应用账户“%s”用户访问令牌权限 成功"
		resMap["IDS_ACCESSTOKENPERM_EXMSG"] = "权限类型：获取用户访问令牌权限；权限：用户访问令牌"
	case complexChinese:
		resMap["IDS_SET_ACCESSTOKENPERM"] = "設定應用帳戶“%s”使用者存取令牌權限 成功"
		resMap["IDS_DELETE_ACCESSTOKENPERM"] = "刪除應用帳戶“%s”使用者存取令牌權限 成功"
		resMap["IDS_ACCESSTOKENPERM_EXMSG"] = "權限類型：獲取使用者存取令牌權限；權限：使用者存取令牌"
	case english:
		resMap["IDS_SET_ACCESSTOKENPERM"] = "Set the permission for user access token of the app account \"%s\" successfully"
		resMap["IDS_DELETE_ACCESSTOKENPERM"] = "Delete the permission for user access token of the app account \"%s\" successfully"
		resMap["IDS_ACCESSTOKENPERM_EXMSG"] = "Get Permission for User Access Token; Permission: User Access Token"
	}
}
