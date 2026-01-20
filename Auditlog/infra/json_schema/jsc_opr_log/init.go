package jsc_opr_log

import (
	"AuditLog/common"
)

func Init() {
	// 1. 初始化模式映射
	m := make(map[OprJSONSchemaCheckKey]string)

	// 2. 定义需要加载的JSON模式文件列表
	files := []string{
		"common.json",            // 基础模式
		"dir_visit.json",         // 目录访问模式
		"menu_button_click.json", // 菜单按钮点击模式
		"common_client.json",     // 客户端通用模式
		"common_server.json",     // 服务端通用模式
	}

	// 3. 加载所有模式文件
	for _, file := range files {
		// 3.1. 从配置文件中读取JSON模式内容
		bys, _, err := common.GetConfigFileContent("operation_log_schema/" + file)
		if err != nil {
			panic(err)
		}

		schema := string(bys)

		// 3.2. 根据文件名将内容存入相应的模式映射
		switch file {
		case "common.json":
			m[OprJSCommon] = schema

		case "dir_visit.json":
			m[OprJSDirVisit] = schema

		case "menu_button_click.json":
			m[OprJSMenuButtonClick] = schema

		case "common_client.json":
			m[OprJSCommonClient] = schema

		case "common_server.json":
			m[OprJSCommonServer] = schema
		}
	}

	// 4. 初始化校验路径映射
	checkPathMap := make(oprJsonSchemaCheckPathMap, len(m))

	// 5. 构建校验路径
	// 5.1. 示例：当验证目录访问日志时，需要同时验证通用模式和目录访问特定模式
	checkPathMap[OprJSCommon] = []string{m[OprJSCommon]}

	// 5.2. 客户端通用路径 = 基础模式 + 客户端通用模式
	checkPathMap[OprJSCommonClient] = append(checkPathMap[OprJSCommon], m[OprJSCommonClient])

	// 5.3. 服务端通用路径 = 基础模式 + 服务端通用模式
	checkPathMap[OprJSCommonServer] = append(checkPathMap[OprJSCommon], m[OprJSCommonServer])

	// 5.4. 菜单按钮点击路径 = 客户端通用路径 + 按钮点击特定模式
	checkPathMap[OprJSMenuButtonClick] = append(checkPathMap[OprJSCommonClient], m[OprJSMenuButtonClick])

	// 5.5. 目录访问路径 = 客户端通用路径 + 目录访问特定模式
	checkPathMap[OprJSDirVisit] = append(checkPathMap[OprJSCommonClient], m[OprJSDirVisit])

	// 6. 将构建好的校验路径映射赋值给全局变量
	OprJsonSchemaCheckPathMap = checkPathMap
}
