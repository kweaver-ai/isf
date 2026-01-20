package interfaces

import (
	"context"
	"database/sql"
	"net/http"

	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/models/rcvo"
	"AuditLog/tapi/sharemgnt"
)

//go:generate mockgen -package mock -source ../interfaces/drivenadapters.go -destination ../interfaces/mock/mock_drivenadapters.go
type UserMgntRepo interface {
	// 获取用户信息 api/user-management/v1/users/{user_ids}/{fields}
	GetUserInfoByID(userIDs []string) (userinfos []models.User, statusCode int, err error)
	GetAppInfoByID(id string) (appInfos models.App, statusCode int, err error)
	GetDeptAllUserIDs(deptID string) (userIDs map[string][]string, statusCode int, err error)
	GetUserIDsByRoleNames(roleNames []string) (roleMemberInfos []*models.RoleMemberInfo, statusCode int, err error)
	GetDeptInfoByIDs(deptIDs []string) (deptInfo []*models.DeptInfo, statusCode int, err error)
	GetDepsByLevel(level int) (deptInfos []*models.DepInfo, statusCode int, err error)
}

type ShareMgntRepo interface {
	GetRoleMemberInfos(roleID string) (res []*sharemgnt.NcTRoleMemberInfo, err error)
	GetTriSystemStatus() (status bool, err error)
}

type DocCenterRepo interface {
	NewDataSourceGroup(body *rcvo.DCNewDataSourceGroupBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error)
	NewDataSource(body *rcvo.DCNewDataSourceBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error)
	GetDataSourceFields(dataSourceID int) (res *rcvo.DCDataSourceFieldsRes, statusCode int, err error)
	NewReportGroup(body *rcvo.DCNewBizGroupBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error)
	NewReport(body *rcvo.DCNewReportBody) (res *rcvo.DCResponse, statusCode int, err error)
}

type OssGatewayRepo interface {
	GetLocalOSSInfo() (res []*models.OSSInfo, statusCode int, err error)
	GetAvailableOSSID() (res string, err error)
	GetUploadInfo(ossID, objName string) (res *models.OSSUploadInfo, statusCode int, err error)
	GetUploadPartRequestInfo(ossID, objName, uploadID string, partNumber int) (res *models.OSSRequestInfo, statusCode int, err error)
	UploadPartByURL(url string, method string, body string, headers map[string]string) (res *models.OSSUploadPartInfo, statusCode int, err error)
	GetCompleteUploadRequestInfo(ossID, objName, uploadID string, multiPartInfo map[int]models.OSSUploadPartInfo) (resp *models.OSSRequestInfo, statusCode int, err error)
	CompleteUploadByURL(url string, method string, body string, headers map[string]string) (resp *http.Response, statusCode int, err error)
	GetDownLoadInfo(ossID, objName, fileName string, isInternal bool) (res *models.OSSRequestInfo, statusCode int, err error)
	DownloadBlockByURL(url string, method string, body string, headers map[string]string, start, end int64) (resp *http.Response, statusCode int, err error)
	GetDeleteRequestInfo(ossID, objName string) (res *models.OSSRequestInfo, statusCode int, err error)
	DeleteObjectByURL(url string, method string, body string, headers map[string]string) (resp *http.Response, statusCode int, err error)
}

type AuditLogPub interface {
	WriteAuditLoginLog(entity *models.AuditLog) (err error)
	WriteAuditManagementLog(entity *models.AuditLog) (err error)
	WriteAuditOperationLog(entity *models.AuditLog) (err error)
}

type LogRepo interface {
	New(log *models.AuditLog) (logID string, err error)
	FindByCondition(offset, limit int, condition string, ids []string) (logs []*models.LogPO, err error)
	FindCountByCondition(condition string) (count int, err error)
	GetFirstLogTime() (timeMicro int64, err error)
	ClearOutdatedLog(logID, date, batchSize, sleepTime int64) (err error)
	GetLogCount() (count int64, err error)
}

type HistoryRepo interface {
	New(log *models.HistoryPO) (err error)
	FindByCondition(offset, limit int, condition string, ids []string) (logs []*models.HistoryPO, err error)
	FindCountByCondition(condition string) (count int, err error)
	GetHistoryLogsByType(logType int8) (logs []*models.HistoryPO, err error)
	GetHistoryLogByID(id string) (log *models.HistoryPO, err error)
}

type DumpStrategyTx interface {
	SetRetentionPeriod(period int) (err error)
	SetRetentionPeriodUnit(unit string) (err error)
	SetDumpFormat(format string) (err error)
	SetDumpTime(time string) (err error)
	Commit() (err error)
	Rollback() (err error)
}

type LogStrategyRepo interface {
	GetRetentionPeriod() (period int, err error)
	GetRetentionPeriodUnit() (unit string, err error)
	GetDumpFormat() (format string, err error)
	GetDumpTime() (time string, err error)
	BeginDumpTx() (tx DumpStrategyTx, err error)
	GetHistoryIsDownloadWithPwd() (isDownload bool, err error)
	SetHistoryIsDownloadWithPwd(isDownload bool) (err error)
	GetLogPrefix() (prefix string, err error)
	SetLogPrefix(prefix string) (err error)
}

type LogScopeStrategyRepo interface {
	GetStrategiesByCondition(condition string, params []interface{}) (res []*lsmodels.ScopeStrategyPO, err error)
	NewStrategy(req *lsmodels.ScopeStrategyPO) (err error)
	UpdateStrategy(req *lsmodels.ScopeStrategyPO) (err error)
	DeleteStrategy(id int64) (err error)
	GetActiveScopeBy(logType int, role string) (scope []string, err error)
	GetHistoryScopeCountBy(logType int, role string) (count int, err error)
	GetStrategyByID(id int64) (res *lsmodels.ScopeStrategyPO, err error)
}

type DBOutbox interface {
	AddOutboxInfos(businessType string, messages []string, tx *sql.Tx) error
	GetPushMessage(businessType string, tx *sql.Tx) (messageID int64, message string, err error)
	DeleteOutboxInfoByID(messageID int64, tx *sql.Tx) error
}

type DLM interface {
	// TryLock 获取锁，只尝试一次，失败不会阻塞
	TryLock(key string) (success bool, err error)
	// UnLock 解锁
	UnLock(key string) (err error)
}

type OpenSearch interface {
	Query(ctx context.Context, dslQuery string, index string) (*models.OSResp, error)
}

type KcRepo interface {
	GetUserInfoByIDS(userIDs []string) (res []*models.KcUserInfo, err error)
}

type PersonalConfigRepo interface {
	GetModuleInfoByName(moduleName string) (res *models.ServiceModuleInfo, statusCode int, err error)
}
