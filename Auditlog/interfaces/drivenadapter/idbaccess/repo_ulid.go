package idbaccess

import (
	"context"
	"database/sql"

	uniqidenums "AuditLog/common/enums/uniqueid"
)

//go:generate mockgen -source=./repo_ulid.go -destination ./dbmock/repo_ulid.go -package dbmock
type UlidRepo interface {
	GenDBID(ctx context.Context, tx *sql.Tx) (id string, err error)
	BatchGenDBID(ctx context.Context, tx *sql.Tx, num int) (ids []string, err error)

	GenUniqID(ctx context.Context, flag uniqidenums.UniqueIDFlag) (id string, err error)
	DelUniqID(ctx context.Context, flag uniqidenums.UniqueIDFlag, id string) (err error)
}
