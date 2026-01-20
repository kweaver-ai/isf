package efastret

import (
	"AuditLog/common/enums"
)

type FsMetadata struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	DocLibType enums.DocLibType `json:"doc_lib_type"`
	Path       string           `json:"path"`
	Size       int64            `json:"size"`
}

// GetFsMetadataRetDto 响应dto
type GetFsMetadataRetDto []*FsMetadata
