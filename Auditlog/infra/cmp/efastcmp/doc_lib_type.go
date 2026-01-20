package efastcmp

import (
	"context"

	"AuditLog/common/enums"
	"AuditLog/infra/cmp/efastcmp/dto/efastarg"
	"AuditLog/infra/cmp/efastcmp/dto/efastret"
)

// GetFileDocLibType 获取文件所属文档库类型
func (e *EFast) GetFileDocLibType(ctx context.Context, ids []string) (m map[string]enums.DocLibType, err error) {
	m = make(map[string]enums.DocLibType)

	dto := &efastarg.GetFsMetadataArgDto{
		IDs: ids,
		Fields: []efastarg.IbField{
			efastarg.IbFieldDocLibTypes,
		},
	}

	var ret efastret.GetFsMetadataRetDto

	ret, err = e.GetFsMetadata(ctx, dto)
	if err != nil {
		return
	}

	for _, v := range ret {
		m[v.ID] = v.DocLibType
	}

	return
}
