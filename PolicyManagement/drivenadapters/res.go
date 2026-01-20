// Package drivenadapters 翻译文件
package drivenadapters

import (
	"sync"

	"policy_mgnt/common"
	"policy_mgnt/common/config"
)

var (
	resOnce sync.Once
	resMap  = make(map[string]string)
)

// loadString 获取字符串型语言项
func loadString(key string) string {
	resOnce.Do(initResMap)

	return resMap[key]
}

//nolint:funlen
func initResMap() {
	switch config.Config.Language {
	case "zh_CN":
		resMap["IDS_ADD_AUTHORIZED_PRODUCTS"] = "给用户“%s” 授权产品“%s” 成功"
		resMap["IDS_DELETE_AUTHORIZED_PRODUCTS"] = "取消用户“%s” 的产品“%s” 授权成功"
		resMap["IDS_UPDATE_AUTHORIZED_PRODUCTS"] = "将用户“%s” 授权的产品“%s” 更新为“%s”成功"

	case "en_US":
		resMap["IDS_ADD_AUTHORIZED_PRODUCTS"] = "The authorization of \"%s\" for user \"%s\" has been completed successfully."
		resMap["IDS_DELETE_AUTHORIZED_PRODUCTS"] = "The authorization for the products \"%s\" has been successfully revoked for user \"%s\"."
		resMap["IDS_UPDATE_AUTHORIZED_PRODUCTS"] = "The authorized products for user \"%s\" have been successfully updated from \"%s\" to \"%s\"."

	case "zh_TW":
		resMap["IDS_ADD_AUTHORIZED_PRODUCTS"] = "給用戶“%s” 授權產品“A%s” 成功。"
		resMap["IDS_DELETE_AUTHORIZED_PRODUCTS"] = "取消用戶“%s” 的產品“%s” 授權成功"
		resMap["IDS_UPDATE_AUTHORIZED_PRODUCTS"] = "將用戶“%s” 授權的產品“%s” 更新為“%s”成功。"

	default:
		common.NewLogger().Fatalln("service language not set")
	}
}
