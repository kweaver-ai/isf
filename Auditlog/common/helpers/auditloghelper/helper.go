package auditloghelper

import (
	"context"
	"time"

	"AuditLog/common/helpers"
)

func NewManagementLog(ctx context.Context, logLevel LogLevel, opType MgtOpLogType) *AuditLog {
	ui := helpers.GetVisitUserInfoFromCtx(ctx)

	aLog := &AuditLog{
		UserID: ui.UserID,

		UserType: NcTLogUserType_NCT_LUT_AUTHUSER,

		Level: logLevel,

		Date: time.Now().UnixMicro(),
		IP:   ui.IP,
		Mac:  "",

		UserAgent: ui.ClientType,

		OpType: int(opType),
	}

	return aLog
}

func NewMgtCreateLog(ctx context.Context) *AuditLog {
	return NewManagementLog(ctx, Info, Create)
}

func NewMgtUpdateLog(ctx context.Context) *AuditLog {
	return NewManagementLog(ctx, Info, Update)
}

func NewMgtDeleteLog(ctx context.Context) *AuditLog {
	return NewManagementLog(ctx, Warn, Delete)
}
