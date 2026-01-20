package efastcmp

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/icmp"
)

type EFast struct {
	privateScheme string
	privateHost   string
	privatePort   int

	publicScheme string
	publicHost   string
	publicPort   int

	arTrace api.Tracer
	logger  api.Logger
}

type EFastConf struct {
	PrivateScheme string
	PrivateHost   string
	PrivatePort   int

	PublicScheme string
	PublicHost   string
	PublicPort   int

	ArTrace api.Tracer
	Logger  api.Logger

	// HttpClient icmp.IHttpClient
}

var _ icmp.IEFast = &EFast{}

func NewEFast(conf *EFastConf) icmp.IEFast {
	arTrace := conf.ArTrace

	return &EFast{
		privateScheme: conf.PrivateScheme,
		privateHost:   conf.PrivateHost,
		privatePort:   conf.PrivatePort,

		publicScheme: conf.PublicScheme,
		publicHost:   conf.PublicHost,
		publicPort:   conf.PublicPort,

		arTrace: arTrace,
		logger:  conf.Logger,
		// httpClient: conf.HttpClient,
	}
}
