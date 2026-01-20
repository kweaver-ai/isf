package ihttpaccess

import (
	"context"

	"AuditLog/infra/cmp/umcmp"
	"AuditLog/infra/cmp/umcmp/dto/umarg"
	"AuditLog/infra/cmp/umcmp/umtypes"
)

//go:generate mockgen -source=./um.go -destination ./httpaccmock/um_mock.go -package httpaccmock
type UmHttpAcc interface {
	GetUserInfo(ctx context.Context, dto *umarg.GetUserInfoArgDto) (uim umcmp.UserInfoMap, err error)
	GetAppIDNameKv(ctx context.Context, appIDs []string) (idNameKvMap map[string]string, err error)

	GetUserDeptIDs(ctx context.Context, userID string) (deptIDs []string, err error)

	GetUserUserGroupIDs(ctx context.Context, userID string) (userGroupIDs []string, err error)

	GetDeptInfoMapByIDs(ctx context.Context, deptIDs []string) (deptInfoMap map[string]*umtypes.DepartmentInfo, err error)

	GetUserDep(ctx context.Context, userID string) (depts [][]umcmp.ObjectBaseInfo, err error)

	// GetDeptIDNameMap(ctx context.Context, deptIDs []string) (idNameMap map[string]string, err error)

	GetOsnNames(ctx context.Context, dto *umarg.GetOsnArgDto) (ret *umtypes.OsnInfoMapS, err error)
}
