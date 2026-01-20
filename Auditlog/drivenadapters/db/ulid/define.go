package dbaulid

import (
	"sync"

	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces/drivenadapter/idbaccess"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

var (
	ulidRepoOnce sync.Once
	ulidRepoImpl idbaccess.UlidRepo
)

type ulidRepo struct {
	db     *sqlx.DB
	logger api.Logger
}

var _ idbaccess.UlidRepo = &ulidRepo{}

func NewUlidRepo(logger api.Logger) idbaccess.UlidRepo {
	ulidRepoOnce.Do(func() {
		ulidRepoImpl = &ulidRepo{
			db:     infra.NewDBPool(),
			logger: logger,
		}
	})

	return ulidRepoImpl
}
