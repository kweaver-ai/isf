package umhttpaccess

import (
	"context"

	"github.com/pkg/errors"

	"AuditLog/common/helpers"
	"AuditLog/infra/cmp/umcmp"
	"AuditLog/infra/cmp/umcmp/dto/umarg"
	"AuditLog/infra/cmp/umcmp/umtypes"
)

func (u *umHttpAcc) GetUserDeptIDs(ctx context.Context, userID string) (deptIDs []string, err error) {
	ctx, span := u.arTrace.AddInternalTrace(ctx)
	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()

	deptIDs, err = u.um.GetUserDeptIDs(ctx, userID)
	if err != nil {
		helpers.RecordErrLogWithPos(u.logger, err, "umHttpAcc.GetUserDeptIDs")

		return nil, errors.Wrap(err, "[GetUserDeptIDs]: 获取用户部门ids失败")
	}

	return
}

func (u *umHttpAcc) GetDeptInfoMapByIDs(ctx context.Context, deptIDs []string) (deptInfoMap map[string]*umtypes.DepartmentInfo, err error) {
	ctx, span := u.arTrace.AddInternalTrace(ctx)
	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()

	deptInfoMap = make(map[string]*umtypes.DepartmentInfo)

	args := &umarg.GetDeptInfoArgDto{
		DeptIds: deptIDs,
		Fields: umarg.DeptFields{
			umarg.DeptFieldName,
			umarg.DeptFieldParentDeps,
		},
	}

	deptInfoMap, err = u.um.GetDeptInfoMap(ctx, args)
	if err != nil {
		helpers.RecordErrLogWithPos(u.logger, err, "umHttpAcc.GetDeptInfoMapByIDs")
		return nil, errors.Wrap(err, "[GetDeptInfoMapByIDs]: 获取部门信息失败")
	}

	return
}

//func (u *umHttpAcc) GetDeptIDNameMap(ctx context.Context, deptIDs []string) (idNameMap map[string]string, err error) {
//	ctx, span := u.arTrace.AddInternalTrace(ctx)
//	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()
//
//	deptInfoMap, err := u.GetDeptInfoMapByIDs(ctx, deptIDs)
//	if err != nil {
//		return
//	}
//
//	idNameMap = make(map[string]string, len(deptInfoMap))
//	for id, info := range deptInfoMap {
//		idNameMap[id] = info.Name
//	}
//
//	return
//}

// GetUserDep 获取用户部门信息（直属部门和上级部门）
func (u *umHttpAcc) GetUserDep(ctx context.Context, userID string) (depts [][]umcmp.ObjectBaseInfo, err error) {
	ctx, span := u.arTrace.AddInternalTrace(ctx)
	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()

	depts, err = u.um.GetUserDept(ctx, userID)
	if err != nil {
		helpers.RecordErrLogWithPos(u.logger, err, "umHttpAcc.GetUserDep")
		return nil, errors.Wrap(err, "[GetUserDep]: 获取用户部门信息失败")
	}

	return
}
