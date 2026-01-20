package oprlogeo

import "AuditLog/common/enums"

// Operator 操作员信息
type Operator struct {
	ID             string            `json:"id"`                        // 操作员的唯一标识符，示例："8b085b72-567c-11ed-aecc-063c8a32c7bf"
	Name           string            `json:"name"`                      // 操作员的姓名，示例："李宇（Aaron）"
	Type           enums.UserType    `json:"type"`                      // 操作员的类型，示例："authenticated_user"
	DepartmentPath []*DepartmentPath `json:"department_path,omitempty"` // 用户所属部门路径，仅当operator.type为"authenticated_user"时需要
	Agent          *Agent            `json:"agent,omitempty"`           // 用户代理，仅当operator.type为"authenticated_user"或"anonymous_user"时需要

	IsSystemOp bool `json:"is_system_op"` // 是否为系统操作（即非用户直接触发，如：系统的一些周期性任务等），默认：false
}
