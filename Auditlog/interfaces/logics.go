package interfaces

import (
	"context"
	"database/sql"

	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/models/rcvo"
)

//go:generate mockgen -package mock -source ../interfaces/logics.go -destination ../interfaces/mock/mock_logics.go
type LogMgnt interface {
	SendLog(info *models.SendLogVo) (err error)
	ReceiveLog(info *models.ReceiveLogVo) (err error)
	ReceiveAuditLog(info *models.ReceiveLogVo) (err error)
	WriteAuditLoginLog(entity interface{}) (err error)
	WriteAuditManagementLog(entity interface{}) (err error)
	WriteAuditOperationLog(entity interface{}) (err error)
	AddAuditLog(visitor Visitor, logType LogType, info *models.ReceiveLogVo) (logID string, err error)
}

type Outbox interface {
	AddOutboxInfo(opType string, content interface{}, tx *sql.Tx) error
	NotifyPushOutboxThread()
	RegisterHandlers(opType string, op func(interface{}) error)
}

type ActiveLog interface {
	GetActiveMetadata() (meta *rcvo.ReportMetadataRes, err error)
	GetActiveDataList(ctx context.Context, logType string, req *rcvo.ReportGetDataListReq, userID string) (res *rcvo.ActiveReportListRes, err error)
	GetActiveFieldValues(ctx context.Context, logType string, req *rcvo.ReportGetFieldValuesReq) (res *rcvo.ReportFieldValuesRes, err error)
}

type HistoryLog interface {
	GetHistoryMetadata() (meta *rcvo.ReportMetadataRes, err error)
	GetHistoryDataList(ctx context.Context, logType string, req *rcvo.ReportGetDataListReq, userID string) (res *rcvo.HistoryReportListRes, err error)
	GetHistoryDownloadPwdStatus(ctx context.Context) (res *lsmodels.HistoryLogDownloadPwdStatus, err error)
	SetHistoryDownloadPwdStatus(ctx context.Context, req *lsmodels.HistoryLogDownloadPwdStatus) (err error)
	CreateDownloadTask(ctx context.Context, req *lsmodels.HistoryLogDownloadReq) (res *lsmodels.HistoryLogDownloadRes, err error)
	GetDownloadProgress(ctx context.Context, taskId string) (res *lsmodels.HistoryLogDownloadProgress, err error)
	GetHistoryDownloadResult(ctx context.Context, taskId string) (res *models.OSSRequestInfo, err error)
}

type DumpLog interface {
	InitDumpLog(ctx context.Context)
}

type LogStrategy interface {
	GetDumpStrategy(ctx context.Context, fields []string) (res map[string]interface{}, err error)
	SetDumpStrategy(ctx context.Context, req map[string]interface{}) (err error)
}

type LogScopeStrategy interface {
	GetStrategy(ctx context.Context, req *lsmodels.GetScopeStrategyReq) (res *lsmodels.GetScopeStrategyRes, err error)
	NewStrategy(ctx context.Context, req *lsmodels.ScopeStrategyVO) (id int64, err error)
	UpdateStrategy(ctx context.Context, id int64, req *lsmodels.ScopeStrategyVO) (err error)
	DeleteStrategy(ctx context.Context, id int64) (err error)
}
