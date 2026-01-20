package umhttpaccess

import (
	"context"

	"github.com/pkg/errors"

	"AuditLog/common/helpers"
	"AuditLog/infra/cmp/umcmp"
	"AuditLog/infra/cmp/umcmp/dto/umarg"
)

// GetUserUserGroupIDs 获取用户的用户组ID列表
func (u *umHttpAcc) GetUserUserGroupIDs(ctx context.Context, userID string) (userGroupIDs []string, err error) {
	ctx, span := u.arTrace.AddInternalTrace(ctx)
	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()

	dto := &umarg.GetUserInfoArgDto{
		UserIds: []string{userID},
		Fields: umarg.Fields{
			umarg.FieldGroups,
		},
	}

	var uim umcmp.UserInfoMap

	uim, err = u.um.GetUserInfo(ctx, dto)
	if err != nil {
		helpers.RecordErrLogWithPos(u.logger, err, "umHttpAcc.GetUserUserGroupIDs")
		return nil, errors.Wrap(err, "[GetUserUserGroupIDs]:获取用户的用户组信息失败")
	}

	userGroupIDs = make([]string, 0)

	for _, ui := range uim {
		for _, group := range ui.Groups {
			userGroupIDs = append(userGroupIDs, group.ID)
		}
	}

	return
}
