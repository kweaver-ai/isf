```go
// Agent 用户代理信息
type Agent struct {
	Type      string `json:"type"` // 客户端类型，示例："windows"
	OSType    string `json:"os_type"`                 // 操作系统类型，示例：windows、linux、android、ios、mac os、unknown
	AppType   string `json:"app_type"`                // 应用类型，示例：同步盘(sync_disk)、富客户端（rich_client）、web、unknown
	IP        string `json:"ip"`   // 操作者IP地址，示例："192.168.50.100"
	UDID      string `json:"udid"`                    // 设备硬件码，示例："3C-2F-10-69-AF-E6"
	UserAgent string `json:"user_agent"`              // 用户代理，来源于请求头参数中的User-Agent
}

// DepartmentPath 部门路径信息
type DepartmentPath struct {
	IDPath   string `json:"id_path"`   // 用户所属部门id全路径，示例："4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0"
	NamePath string `json:"name_path"` // 用户所属部门名称全路径，示例："爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
}

// DocLib 文档库信息
type DocLib struct {
	ID   string `json:"id"`   // 文档库的唯一标识符，示例："gns://D42F2729C56E489A948985D4E75C5813"
	Type string `json:"type"` // 文档库的类型，示例："department_doc_lib"
	Name string `json:"name"`                    // 文档库的名称，示例："部门文档库1"
}

// LogEntry 主日志结构
type LogEntry struct {
	Recorder     string         `json:"recorder"` // 日志记录者的身份，示例："AnyShare"
	//BizType    BizType         `json:"biz_type"`     // 业务类型
	Operation    string         `json:"operation"`   // 执行的操作，示例："cd"
	Description  string         `json:"description"` // 操作的描述，示例："用户"张三"从""进入到"部门文档库1/a"。"
	IsSystemOp   bool           `json:"is_system_op"`                   // 是否为系统操作（即非用户直接触发，如：系统的一些周期性任务等），默认：false
	Operator     *Operator      `json:"operator"`    // 操作员信息
	Object       *ObjectInfo    `json:"object"`      // 操作对象信息
	TargetObject *ObjectInfo    `json:"target_object"`                  // 目标对象信息
	LogFrom      *LogFrom       `json:"log_from"`    // 日志来源
	Rec          *RecommendInfo `json:"rec"`                            // 推荐相关对象
	Detail       interface{}    `json:"detail"`                         // 业务模块扩展的其他字段
}

// LogFrom 日志来源信息
type LogFrom struct {
	Package string       `json:"package"` // 大包名，示例：package项目application.json文件中的name字段的值，如as主模块为："AnyShareMainModule"
	Service *ServiceInfo `json:"service"` // 服务信息，示例：服务的信息
}

// ObjectInfo 对象信息
type ObjectInfo struct {
	ID     string   `json:"id"`      // 对象的唯一标识符，示例："gns://D42F2729C56E489A948985D4E75C5813"
	Path   string   `json:"path"`    // 对象的路径，示例："部门文档库1/a/b"
	Name   string   `json:"name"`                       // 对象的名称，Path的最后一部分，示例："b"
	Type   string   `json:"type"`    // 对象的类型，示例："folder"
	Size   uint64   `json:"size"`                       // 对象的大小，示例：0
	Tags   []string `json:"tags"`                       // 对象的标签名，示例：["tag1", "tag2"]。目前可能只针对“object", 不针对"target_object"
	DocLib *DocLib  `json:"doc_lib"` // 文档库信息
}

// Operator 操作员信息
type Operator struct {
	ID             string            `json:"id"`   // 操作员的唯一标识符，示例："8b085b72-567c-11ed-aecc-063c8a32c7bf"
	Name           string            `json:"name"` // 操作员的姓名，示例："李宇（Aaron）"
	Type           string            `json:"type"` // 操作员的类型，示例："authenticated_user"
	DepartmentPath []*DepartmentPath `json:"department_path"`         // 用户所属部门路径，仅当operator.type为"authenticated_user"时需要
	Agent          *Agent            `json:"agent"`                   // 用户代理，仅当operator.type为"authenticated_user"或"anonymous_user"时需要
}

// RecommendInfo 推荐相关信息
type RecommendInfo struct {
	NotUseForRec bool   `json:"not_use_for_rec"` // 是否不用于推荐，默认：false
	ExtInfoJSON  string `json:"ext_info_json"`   // 推荐用扩展信息，示例：{"k1":{"k2::"v2"}}
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Instance struct {
		ID string `json:"id"` // 示例：实例的ID
	} `json:"instance"`
	Name    string `json:"name"` // 示例：服务名称
	Version string `json:"version"` // 示例：服务版本
}
```