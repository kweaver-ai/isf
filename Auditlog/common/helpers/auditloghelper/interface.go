package auditloghelper

import (
	"context"
)

type IBuildAuditLog interface {
	AuditLogCreate(ctx context.Context) *AuditLog
	AuditLogUpdate(ctx context.Context) *AuditLog
	AuditLogDelete(ctx context.Context) *AuditLog
}
