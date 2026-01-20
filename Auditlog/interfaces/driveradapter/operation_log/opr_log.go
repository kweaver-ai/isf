package oprlogdriveri

import (
	"context"

	"AuditLog/common/enums/oprlogenums"
)

//go:generate mockgen -package mock -source opr_log.go -destination ../mock/opr_log.go
type IOprLogSvc interface {
	HandleMsg(ctx context.Context, msg []byte, bizType oprlogenums.BizType) (err error)
	WriteOperationLogToMQ(ctx context.Context, topic oprlogenums.OperationLogTopic, msgByte []byte) (err error)
	MsgToMap(msg []byte) (logMaps []map[string]interface{}, err error)
}
