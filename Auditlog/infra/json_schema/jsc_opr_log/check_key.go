package jsc_opr_log

// OprJSONSchemaCheckKey 运营日志模式键的类型定义
type OprJSONSchemaCheckKey string

// 定义不同类型的运营日志模式键
const (
	// OprJSCommon 通用基础模式
	OprJSCommon OprJSONSchemaCheckKey = "common"

	// OprJSDirVisit 目录访问模式
	OprJSDirVisit OprJSONSchemaCheckKey = "dir_visit"

	// OprJSMenuButtonClick 菜单按钮点击模式
	OprJSMenuButtonClick OprJSONSchemaCheckKey = "menu_button_click"

	// OprJSCommonClient 客户端通用模式
	OprJSCommonClient OprJSONSchemaCheckKey = "common_client"

	// OprJSCommonServer 服务端通用模式
	OprJSCommonServer OprJSONSchemaCheckKey = "common_server"
)
