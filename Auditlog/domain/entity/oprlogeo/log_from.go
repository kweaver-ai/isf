package oprlogeo

// LogFrom 日志来源信息
type LogFrom struct {
	Package string       `json:"package"` // 大包名，示例：package项目application.json文件中的name字段的值，如as主模块为："AnyShareMainModule"
	Service *ServiceInfo `json:"service"` // 服务信息，示例：服务的信息
}
