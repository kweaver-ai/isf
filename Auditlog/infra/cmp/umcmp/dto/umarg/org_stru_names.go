package umarg

import (
	"net/http"

	"AuditLog/common/utils"
)

// GetOsnArgDto 获取组织架构对象的names的参数
// osn: org structure names
type GetOsnArgDto struct {
	UserIDs       []string `json:"user_ids,omitempty"`
	DepartmentIDs []string `json:"department_ids,omitempty"`
	GroupIDs      []string `json:"group_ids,omitempty"`
	AppIDs        []string `json:"app_ids,omitempty"`
}

// DeDupl 去重
func (d *GetOsnArgDto) DeDupl() {
	d.UserIDs = utils.DeduplGeneric(d.UserIDs)
	d.DepartmentIDs = utils.DeduplGeneric(d.DepartmentIDs)
	d.GroupIDs = utils.DeduplGeneric(d.GroupIDs)
}

type GetOsnUMArgDto struct {
	*GetOsnArgDto
	Method string `json:"method"`
}

func NewGetOsnUMArgDto(getOsnArgDto *GetOsnArgDto) *GetOsnUMArgDto {
	// 去重
	getOsnArgDto.DeDupl()

	return &GetOsnUMArgDto{
		GetOsnArgDto: getOsnArgDto,
		Method:       http.MethodGet,
	}
}
