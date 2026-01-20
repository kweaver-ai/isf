package umhttpaccess

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/infra/cmp/umcmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

type umHttpAcc struct {
	um      *umcmp.Um
	logger  api.Logger
	arTrace api.Tracer
}

var _ ihttpaccess.UmHttpAcc = &umHttpAcc{}

func NewUmHttpAcc(arTrace api.Tracer,
	logger api.Logger, umConf *config.UmConf,
) ihttpaccess.UmHttpAcc {
	_um := umcmp.NewUmCmp(
		umConf,
		logger,
		arTrace,
	)

	umImpl := &umHttpAcc{
		um:      _um,
		logger:  logcmp.GetLogger(),
		arTrace: arTrace,
	}

	return umImpl
}
