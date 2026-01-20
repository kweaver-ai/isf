package docaccret

import "AuditLog/common/enums"

// DocLibItem http://{host}:{port}/api/document/v1/batch-doc-libs/{fields} 响应体
type DocLibItem struct {
	ID   string           `json:"id"`
	Name string           `json:"name"`
	Type enums.DocLibType `json:"doc_type"`
}
