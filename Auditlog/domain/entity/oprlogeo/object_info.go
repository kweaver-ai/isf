package oprlogeo

// ObjectInfo 对象信息
type ObjectInfo struct {
	ID     string   `json:"id"`                // 对象的唯一标识符，示例："gns://D42F2729C56E489A948985D4E75C5813"
	Path   string   `json:"path"`              // 对象的路径，示例："部门文档库1/a/b"
	Name   string   `json:"name"`              // 对象的名称，Path的最后一部分，示例："b"
	Type   string   `json:"type"`              // 对象的类型，示例："folder"
	Size   int      `json:"size"`              // 对象的大小，示例：0
	Tags   []string `json:"tags,omitempty"`    // 对象的标签名，示例：["tag1", "tag2"]。目前可能只针对“object", 不针对"target_object"
	DocLib *DocLib  `json:"doc_lib,omitempty"` // 文档库信息
}
