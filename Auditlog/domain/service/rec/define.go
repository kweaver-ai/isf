package recsvc

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/icmp"
	persrecrepoi "AuditLog/interfaces/drivenadapter/idbaccess/persrec"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
	recdriveri "AuditLog/interfaces/driveradapter/rec"
)

type recSvc struct {
	logger     api.Logger
	opsHttpAcc ihttpaccess.OpsHttpAcc

	svcConfigRepo persrecrepoi.IPersSvcConfigRepo

	dmlCmp icmp.RedisDlmCmp
}

func NewRecSvc(
	logger api.Logger,
	oprHttpAcc ihttpaccess.OpsHttpAcc,
	svcConfigRepo persrecrepoi.IPersSvcConfigRepo,
	dmlCmp icmp.RedisDlmCmp,
) recdriveri.IRecSvc {
	svc := &recSvc{
		logger:        logger,
		opsHttpAcc:    oprHttpAcc,
		svcConfigRepo: svcConfigRepo,
		dmlCmp:        dmlCmp,
	}

	return svc
}
