package mbceo

import (
	"AuditLog/common/enums/oprlogenums/mbcenums"
	"AuditLog/common/utils"
)

type Detail struct {
	OpName   string          `json:"op_name"`            // 操作名称（按钮名称）
	OpType   mbcenums.OpType `json:"op_type"`            // 操作类型
	Position *Position       `json:"position,omitempty"` // 操作按钮所在位置
}

func NewDetail() *Detail {
	return &Detail{}
}

func (d *Detail) LoadByInterface(i interface{}) (err error) {
	if i == nil {
		return
	}

	//    通过json来实现
	jsonStr, err := utils.JSON().Marshal(i)
	if err != nil {
		return
	}

	err = utils.JSON().Unmarshal(jsonStr, d)
	if err != nil {
		return
	}

	return
}

type Position struct {
	Type     mbcenums.PositionType `json:"type"`
	PathID   string                `json:"path_id"`
	PathName string                `json:"path_name"`
	AreaType mbcenums.AreaType     `json:"area_type"`
}

func (p *Position) IsDocLibPosition() bool {
	return p.Type == mbcenums.DirPoT
}
