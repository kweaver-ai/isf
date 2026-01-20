package oprlogmq

import (
	"sync"

	"AuditLog/common"
	oprinject "AuditLog/domain/service/inject/operation_log"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	oprlogdriveri "AuditLog/interfaces/driveradapter/operation_log"
)

type oprLogMqHandler struct {
	logger    api.Logger
	client    api.MQClient
	oprLogSvc oprlogdriveri.IOprLogSvc
}

var (
	olMqOnce sync.Once
	olMq     interfaces.MQHandler
)

func NewOprLogMqHandler() interfaces.MQHandler {
	olMqOnce.Do(func() {
		olMq = &oprLogMqHandler{
			logger:    common.SvcConfig.Logger,
			client:    api.NewMQClient(),
			oprLogSvc: oprinject.NewOprLogSvc(),
		}
	})

	return olMq
}
