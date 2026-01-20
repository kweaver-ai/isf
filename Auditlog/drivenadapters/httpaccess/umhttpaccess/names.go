package umhttpaccess

import (
	"context"

	"github.com/pkg/errors"

	"AuditLog/common/helpers"
	"AuditLog/infra/cmp/umcmp/dto/umarg"
	"AuditLog/infra/cmp/umcmp/umtypes"
)

func (u *umHttpAcc) GetOsnNames(ctx context.Context, dto *umarg.GetOsnArgDto) (ret *umtypes.OsnInfoMapS, err error) {
	ctx, span := u.arTrace.AddInternalTrace(ctx)
	defer func() { u.arTrace.TelemetrySpanEnd(span, err) }()

	if helpers.IsLocalDev() && len(dto.AppIDs) > 0 {
		ret = &umtypes.OsnInfoMapS{
			AppNameMap: map[string]string{
				"app_id_1": "app_id_1_name",
				"app_id_3": "app_id_3_name",
			},
		}

		return
	}

	ret, err = u.um.GetOsnNames(ctx, dto)
	if err != nil {
		return nil, errors.Wrap(err, "获取组织架构names失败")
	}

	return
}
