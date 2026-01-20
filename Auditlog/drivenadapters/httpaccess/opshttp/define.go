package opshttp

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/icmp"
	"AuditLog/infra/cmp/opensearchcmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

type opsHttpAcc struct {
	opsCmp  icmp.IOpsCmp
	logger  api.Logger
	arTrace api.Tracer
}

var _ ihttpaccess.OpsHttpAcc = &opsHttpAcc{}

// NewOpsHttpAcc 实例化open search adapter
func NewOpsHttpAcc(arTrace api.Tracer,
	logger api.Logger, opsConf *config.OpenSearchConf,
) (impl ihttpaccess.OpsHttpAcc, err error) {
	conf := &opensearchcmp.OpsCmpConf{
		Address:  opsConf.GetAddress(),
		Username: opsConf.User,
		Password: opsConf.Password,

		ArTrace: arTrace,
		Logger:  logger,
	}

	cmp, err := opensearchcmp.NewOpsCmp(conf)
	if err != nil {
		return
	}

	impl = &opsHttpAcc{
		opsCmp:  cmp,
		logger:  logger,
		arTrace: arTrace,
	}

	return
}
