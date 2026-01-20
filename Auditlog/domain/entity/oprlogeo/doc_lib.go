package oprlogeo

import "AuditLog/common/enums"

// DocLib 文档库信息
type DocLib struct {
	ID                  string           `json:"id"`                      // 文档库的唯一标识符，示例："gns://D42F2729C56E489A948985D4E75C5813"
	Type                enums.DocLibType `json:"type"`                    // 文档库的类型，示例："department_doc_lib"
	Name                string           `json:"name"`                    // 文档库的名称，示例："部门文档库1"
	CustomDocLibSubType string           `json:"custom_doc_lib_sub_type"` // doc_lib.type="custom_doc_lib"时的库分类
}
