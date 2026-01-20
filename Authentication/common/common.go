package common

import (
	"Authentication/interfaces"
)

// ClientTypeToStringMap 客户端类型转换Map
var ClientTypeToStringMap = map[interfaces.ClientType]string{
	interfaces.Unknown:      "unknown",
	interfaces.IOS:          "ios",
	interfaces.Android:      "android",
	interfaces.WindowsPhone: "windows_phone",
	interfaces.Windows:      "windows",
	interfaces.MacOS:        "mac_os",
	interfaces.Web:          "web",
	interfaces.MobileWeb:    "mobile_web",
	interfaces.ConsoleWeb:   "console_web",
	interfaces.DeployWeb:    "deploy_web",
	interfaces.Linux:        "linux",
	interfaces.APP:          "app",
}

// ReverseClientTypeToStringMap 获取客户端类型string转int Map
func ReverseClientTypeToStringMap() (clientTypeStringToIntMap map[string]interfaces.ClientType) {
	clientTypeStringToIntMap = make(map[string]interfaces.ClientType)
	for k, v := range ClientTypeToStringMap {
		clientTypeStringToIntMap[v] = k
	}
	return
}
