package efasthttp

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/efastcmp"
	"AuditLog/infra/cmp/icmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

type eFastHttpAcc struct {
	eFast   icmp.IEFast
	logger  api.Logger
	arTrace api.Tracer
}

var _ ihttpaccess.EFastHttpAcc = &eFastHttpAcc{}

func NewEFastHttpAcc(arTrace api.Tracer,
	logger api.Logger, efConf *config.EFastConf,
	// httpClient icmp.IHttpClient,
) ihttpaccess.EFastHttpAcc {
	conf := &efastcmp.EFastConf{
		PrivateScheme: efConf.Private.Protocol,
		PrivateHost:   efConf.Private.Host,
		PrivatePort:   efConf.Private.Port,

		PublicScheme: efConf.Public.Protocol,
		PublicHost:   efConf.Public.Host,
		PublicPort:   efConf.Public.Port,

		ArTrace: arTrace,
		Logger:  logger,
	}
	_eFast := efastcmp.NewEFast(conf)

	eFastImpl := &eFastHttpAcc{
		eFast:   _eFast,
		logger:  logger,
		arTrace: arTrace,
	}

	return eFastImpl
}
