package oprlogeo

// Agent 用户代理信息
type Agent struct {
	Type      string `json:"type"`       // 客户端类型，示例："windows"
	OSType    string `json:"os_type"`    // 操作系统类型，示例：windows、linux、android、ios、mac os、unknown
	AppType   string `json:"app_type"`   // 应用类型，示例：同步盘(sync_disk)、富客户端（rich_client）、web、unknown
	IP        string `json:"ip"`         // 操作者IP地址，示例："192.168.50.100"
	UDID      string `json:"udid"`       // 设备硬件码，示例："3C-2F-10-69-AF-E6"
	UserAgent string `json:"user_agent"` // 用户代理，来源于请求头参数中的User-Agent
}
