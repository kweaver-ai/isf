package complete_info

import (
	"context"

	"AuditLog/common/enums"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/infra/cmp/umcmp"
)

func (l *CompleteInfo) completeDepInfo(c context.Context, maps []map[string]interface{}, eos []*oprlogeo.LogEntry) (err error) {
	if len(eos) == 0 {
		return
	}

	if eos[0].Operator == nil {
		return
	}

	uid := eos[0].Operator.ID
	if uid == "" {
		return
	}

	userType := eos[0].Operator.Type

	var depts [][]umcmp.ObjectBaseInfo

	// 1. 获取用户部门信息
	if userType == enums.AuthUser {
		depts, err = l.umHttpAcc.GetUserDep(c, uid)
		if err != nil {
			return
		}
	}

	if len(depts) == 0 {
		return
	}

	// 2. 填充用户部门信息 for eos
	for i := range eos {
		if len(depts) > 0 {
			eos[i].Operator.DepartmentPath = oprlogeo.ParseDepartments(depts)
		}
	}

	//	3. 填充用户部门信息 for maps
	err = l.mergeEosToMap(maps, eos)

	return
}
