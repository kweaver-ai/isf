package umcmp

import (
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
)

type Um struct {
	umConf *config.UmConf

	arTrace api.Tracer
	logger  api.Logger
}

func NewUmCmp(umConf *config.UmConf,
	logger api.Logger,
	arTrace api.Tracer,
) *Um {
	return &Um{
		umConf:  umConf,
		arTrace: arTrace,
		logger:  logger,
	}
}
