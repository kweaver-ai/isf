package logics

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"AuditLog/common"
	"AuditLog/common/constants/logconsts"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/common/utils"
	"AuditLog/common/utils/dumplogutils"
	"AuditLog/common/utils/rclogutils"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/models/rcvo"
)

var (
	h     *HistoryLog
	hOnce sync.Once
)

type HistoryLog struct {
	logger               api.Logger
	tracer               api.Tracer
	cache                redis.Cmdable
	historyLogRepo       interfaces.HistoryRepo
	userMgntRepo         interfaces.UserMgntRepo
	logStrategyRepo      interfaces.LogStrategyRepo
	logScopeStrategyRepo interfaces.LogScopeStrategyRepo
	ossGateway           interfaces.OssGatewayRepo
	logMgnt              interfaces.LogMgnt
}

func NewHistoryLog() interfaces.HistoryLog {
	hOnce.Do(func() {
		h = &HistoryLog{
			logger:               logger,
			tracer:               tracer,
			cache:                redisClient,
			historyLogRepo:       historyRepo,
			userMgntRepo:         userMgntRepo,
			logStrategyRepo:      logStrategyRepo,
			logScopeStrategyRepo: logScopeStrategyRepo,
			ossGateway:           ossGateway,
			logMgnt:              NewLogMgnt(),
		}
	})

	return h
}

// GetHistoryMetadata 获取历史审计日志元数据
func (h *HistoryLog) GetHistoryMetadata() (meta *rcvo.ReportMetadataRes, err error) {
	ctx := context.TODO()
	return rclogutils.GetHistoryMetadata(ctx)
}

// GetHistoryDataList 获取历史审计日志
func (h *HistoryLog) GetHistoryDataList(ctx context.Context, category string, req *rcvo.ReportGetDataListReq, userID string) (res *rcvo.HistoryReportListRes, err error) {
	var tErr error
	_, span := h.tracer.AddInternalTrace(ctx)
	defer func() { h.tracer.TelemetrySpanEnd(span, tErr) }()

	userInfos, _, err := h.userMgntRepo.GetUserInfoByID([]string{userID})
	if err != nil {
		h.logger.Errorf("[GetHistoryDataList]: get user info error: %v", err)
		return
	}

	if len(userInfos) != 0 {
		hasPermission := false
		userRoles := userInfos[0].Roles

		// 超级管理员 直接访问
		if common.InArray(common.SuperAdmin, userRoles) {
			hasPermission = true
		} else {
			// 三权分立角色 系统管理员、安全管理员、审计管理员根据策略访问
			in := utils.Intersection(common.MutuallyRoles, userRoles)
			for _, r := range in {
				count, err := h.logScopeStrategyRepo.GetHistoryScopeCountBy(common.LogTypeMap[category], r)
				if err != nil {
					h.logger.Errorf("[GetHistoryDataList]: get scope count by role error: %v", err)
					return nil, err
				}
				if count > 0 {
					hasPermission = true
					break
				}
			}
		}

		if hasPermission {
			sqlStr, err := rclogutils.BuildHistoryCondition(category, req.Condition, req.OrderBy, req.IDs)
			if err != nil {
				return nil, err
			}

			logs, err := h.historyLogRepo.FindByCondition(req.Offset, req.Limit, sqlStr, req.IDs)
			if err != nil {
				return nil, err
			}

			entries := make(rcvo.HistoryLogReports, 0, len(logs))

			for _, log := range logs {
				entry := rcvo.HistoryLogReport{
					ID:       log.ID,
					DumpDate: log.DumpDate / 1000,
					FileName: log.Name,
					Size:     common.FormatFileSize(int(log.Size), 2),
				}
				entries = append(entries, entry)
			}

			// 获取日志总数
			nTotalCount := len(req.IDs)
			if nTotalCount == 0 {
				nTotalCount, err = h.historyLogRepo.FindCountByCondition(sqlStr)
				if err != nil {
					return nil, err
				}
			}

			res = &rcvo.HistoryReportListRes{
				Entries:    entries,
				TotalCount: nTotalCount,
			}
		} else {
			return nil, errors.NewCtx(ctx, errors.ForbiddenErr, "No permission", nil)
		}
	}

	return
}

// GetHistoryDownloadPwdStatus 获取历史审计日志下载密码状态
func (h *HistoryLog) GetHistoryDownloadPwdStatus(ctx context.Context) (res *lsmodels.HistoryLogDownloadPwdStatus, err error) {
	status, err := h.logStrategyRepo.GetHistoryIsDownloadWithPwd()
	if err != nil {
		return
	}

	res = &lsmodels.HistoryLogDownloadPwdStatus{
		Status: status,
	}

	return
}

// SetHistoryDownloadPwdStatus 设置历史审计日志下载密码状态
func (h *HistoryLog) SetHistoryDownloadPwdStatus(ctx context.Context, req *lsmodels.HistoryLogDownloadPwdStatus) (err error) {
	if err = h.logStrategyRepo.SetHistoryIsDownloadWithPwd(req.Status); err != nil {
		h.logger.Errorf("[SetHistoryDownloadPwdStatus] set history log download pwd status error: %v", err)
		return
	}

	// 记录设置下载密码状态日志
	go func() {
		visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
		msg := locale.GetI18nCtx(ctx, locale.SetHistoryEncrypted)
		if !req.Status {
			msg = locale.GetI18nCtx(ctx, locale.CancelHistoryEncrypted)
		}
		err := h.logMgnt.SendLog(&models.SendLogVo{
			LogType:  common.Management,
			Language: "",
			LogContent: &models.AuditLog{
				UserID:    visitor.ID,
				UserName:  visitor.Name,
				UserType:  common.AuthenticatedUser,
				Level:     logconsts.LogLevel.INFO,
				OpType:    logconsts.OpType.ManagementType.SET,
				Date:      time.Now().UnixMicro(),
				IP:        visitor.IP,
				Mac:       visitor.Mac,
				Msg:       msg,
				Exmsg:     "",
				UserAgent: visitor.AgentType,
				OutBizID:  uuid.NewString(),
			},
		})
		if err != nil {
			h.logger.Warnf("[SetHistoryDownloadPwdStatus] send log error: %v", err)
		}
	}()
	return
}

// CreateDownloadTask 创建历史审计日志下载任务
func (h *HistoryLog) CreateDownloadTask(ctx context.Context, req *lsmodels.HistoryLogDownloadReq) (res *lsmodels.HistoryLogDownloadRes, err error) {
	// 获取历史审计日志信息
	logInfo, err := h.historyLogRepo.GetHistoryLogByID(req.ObjId)
	if err != nil {
		h.logger.Errorf("[CreateDownloadTask] get history log info error: %v", err)
		return
	}
	if logInfo == nil {
		return nil, errors.NewCtx(ctx, errors.BadRequestErr, "history log %s not found", []string{req.ObjId})
	}
	// 检查是否需要密码
	pwdStatus, err := h.logStrategyRepo.GetHistoryIsDownloadWithPwd()
	if err != nil {
		h.logger.Errorf("[CreateDownloadTask] get history log info error: %v", err)
		return
	}
	if pwdStatus && req.Password == "" {
		return nil, errors.NewCtx(ctx, errors.PasswordRequiredErr, "password is required", nil)
	}
	if req.Password != "" {
		decryptedPwd, err := utils.Rsa2048Decrypt(req.Password)
		if err != nil {
			h.logger.Errorf("[CreateDownloadTask] decrypt password error: %v", err)
			return nil, err
		}
		if !utils.CheckPassword(decryptedPwd) {
			return nil, errors.NewCtx(ctx, errors.PasswordInvalidErr, "password is invalid", nil)
		}
		req.Password = decryptedPwd
	}

	downloadTaskId := uuid.New().String()
	err = utils.SetCache(
		ctx,
		h.cache,
		lsconsts.TaskCacheKey+downloadTaskId,
		lsmodels.HistoryLogDownloadTaskInfo{
			Status: lsconsts.PendingStatus,
		},
		lsconsts.TaskInfoExpire,
	)
	if err != nil {
		h.logger.Errorf("[CreateDownloadTask] set redis cache failed, err is: %v", err)
		return
	}
	res = &lsmodels.HistoryLogDownloadRes{
		TaskId: downloadTaskId,
	}

	go h.compressLog(ctx, downloadTaskId, req, logInfo)
	return
}

// GetDownloadProgress 获取历史审计日志下载进度
func (h *HistoryLog) GetDownloadProgress(ctx context.Context, taskId string) (res *lsmodels.HistoryLogDownloadProgress, err error) {
	taskInfo, err := utils.GetCache[lsmodels.HistoryLogDownloadTaskInfo](ctx, h.cache, lsconsts.TaskCacheKey+taskId)
	if err != nil {
		h.logger.Errorf("[GetDownloadProgress] get redis cache failed, err is: %v", err)
		return
	}
	if taskInfo == nil {
		return nil, errors.NewCtx(ctx, errors.BadRequestErr, "task not found", nil)
	}

	if taskInfo.Status == lsconsts.PendingStatus {
		res = &lsmodels.HistoryLogDownloadProgress{
			Status: false,
		}
	} else {
		res = &lsmodels.HistoryLogDownloadProgress{
			Status: true,
		}
	}

	return
}

// GetHistoryDownloadResult 获取历史审计日志下载结果
func (h *HistoryLog) GetHistoryDownloadResult(ctx context.Context, taskId string) (res *models.OSSRequestInfo, err error) {
	taskInfo, err := utils.GetCache[lsmodels.HistoryLogDownloadTaskInfo](ctx, h.cache, lsconsts.TaskCacheKey+taskId)
	if err != nil {
		h.logger.Errorf("[GetHistoryDownloadResult] get redis cache failed, err is: %v", err)
		return
	}
	if taskInfo == nil {
		return nil, errors.NewCtx(ctx, errors.BadRequestErr, "task not found", nil)
	}

	if taskInfo.Status == lsconsts.FinishedStatus {
		res, _, err = h.ossGateway.GetDownLoadInfo(taskInfo.OssId, taskId, "log.zip", false)
		if err != nil {
			h.logger.Errorf("[GetHistoryDownloadResult] get pub download info error: %v", err)
			return nil, err
		}
		// 记录下载日志
		go func() {
			visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
			err := h.logMgnt.SendLog(&models.SendLogVo{
				LogType:  common.Management,
				Language: "",
				LogContent: &models.AuditLog{
					UserID:    visitor.ID,
					UserName:  visitor.Name,
					UserType:  common.AuthenticatedUser,
					Level:     logconsts.LogLevel.INFO,
					OpType:    logconsts.OpType.ManagementType.EXPORT,
					Date:      time.Now().UnixMicro(),
					IP:        visitor.IP,
					Mac:       visitor.Mac,
					Msg:       fmt.Sprintf(locale.GetI18nCtx(ctx, locale.ExportLogSuccess), taskInfo.FileName),
					Exmsg:     "",
					UserAgent: visitor.AgentType,
					OutBizID:  uuid.NewString(),
				},
			})
			if err != nil {
				h.logger.Warnf("[GetHistoryDownloadResult] send log error: %v", err)
			}
		}()
	} else if taskInfo.Status == lsconsts.ErrorStatus {
		return nil, errors.NewCtx(ctx, errors.InternalErr, "download task error", nil)
	} else {
		return nil, errors.NewCtx(ctx, errors.InternalErr, "download task not finished", nil)
	}

	return
}

// compressLog 压缩历史审计日志
func (h *HistoryLog) compressLog(ctx context.Context, downloadTaskId string, req *lsmodels.HistoryLogDownloadReq, logInfo *models.HistoryPO) {
	ossID, err := h.doCompressLog(ctx, downloadTaskId, req, logInfo)
	if err != nil {
		h.logger.Warnf("[compressLog] do compress log error: %v", err)
		if err = utils.SetCache(
			ctx,
			h.cache,
			lsconsts.TaskCacheKey+downloadTaskId,
			lsmodels.HistoryLogDownloadTaskInfo{
				Status: lsconsts.ErrorStatus,
			},
			lsconsts.TaskInfoExpire,
		); err != nil {
			h.logger.Warnf("[compressLog] set redis cache failed, err is: %v", err)
		}
		return
	}

	if err = utils.SetCache(
		ctx,
		h.cache,
		lsconsts.TaskCacheKey+downloadTaskId,
		lsmodels.HistoryLogDownloadTaskInfo{
			Status:   lsconsts.FinishedStatus,
			OssId:    ossID,
			FileName: logInfo.Name,
		},
		lsconsts.TaskInfoExpire,
	); err != nil {
		h.logger.Warnf("[compressLog] set redis cache failed, err is: %v", err)
		return
	}
}

func (h *HistoryLog) doCompressLog(ctx context.Context, downloadTaskId string, req *lsmodels.HistoryLogDownloadReq, logInfo *models.HistoryPO) (ossID string, err error) {
	// 下载历史日志文件
	fileContent, err := h.genDownloadFile(ctx, logInfo)
	if err != nil {
		h.logger.Errorf("[doCompressLog] gen download file error: %v", err)
		return
	}

	// 生成zip文件
	zipContent, err := dumplogutils.GenZipFile(fileContent, logInfo.Name, req.Password)
	if err != nil {
		h.logger.Errorf("[doCompressLog] gen zip file error: %v", err)
		return
	}

	ossID, err = h.ossGateway.GetAvailableOSSID()
	if err != nil {
		h.logger.Errorf("[doCompressLog] get oss id error: %v", err)
		return
	}

	uploadInfo, _, err := h.ossGateway.GetUploadInfo(ossID, downloadTaskId)
	if err != nil {
		h.logger.Errorf("[doCompressLog] get upload info error: %v", err)
		return
	}

	// 对zip文件进行分块上传
	uploadID := uploadInfo.UploadID
	partInfos := make(map[int]models.OSSUploadPartInfo)

	parts, err := dumplogutils.SplitFile(zipContent, int64(uploadInfo.PartSize))
	if err != nil {
		h.logger.Errorf("[doCompressLog] split file error: %v", err)
		return
	}

	partNumber := 1
	for _, partData := range parts {
		uploadPartRequestInfo, _, err := h.ossGateway.GetUploadPartRequestInfo(ossID, downloadTaskId, uploadID, partNumber)
		if err != nil || uploadPartRequestInfo == nil {
			h.logger.Errorf("[doCompressLog] get upload part request info error: %v", err)
			return "", err
		}

		partInfo, _, err := h.ossGateway.UploadPartByURL(
			uploadPartRequestInfo.URL,
			uploadPartRequestInfo.Method,
			string(partData),
			uploadPartRequestInfo.Headers,
		)
		if err != nil || partInfo == nil {
			h.logger.Errorf("[doCompressLog] upload part error: %v", err)
			return "", err
		}

		partInfos[partNumber] = models.OSSUploadPartInfo{
			Etag: partInfo.Etag,
			Size: partInfo.Size,
		}
		partNumber++
	}

	completeUploadRequestInfo, _, err := h.ossGateway.GetCompleteUploadRequestInfo(ossID, downloadTaskId, uploadID, partInfos)
	if err != nil || completeUploadRequestInfo == nil {
		h.logger.Errorf("[doCompressLog] get complete upload request info error: %v", err)
		return "", err
	}

	completeUploadResponse, _, err := h.ossGateway.CompleteUploadByURL(
		completeUploadRequestInfo.URL,
		completeUploadRequestInfo.Method,
		completeUploadRequestInfo.RequestBody,
		completeUploadRequestInfo.Headers,
	)
	if err != nil || completeUploadResponse == nil {
		h.logger.Errorf("[doCompressLog] complete upload error: %v", err)
		return "", err
	}
	h.logger.Infof("[doCompressLog] complete upload zip file of %s to oss, ossID: %s, fileId: %s", logInfo.Name, ossID, downloadTaskId)

	// 30分钟后删除文件
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Minute):
			h.deleteLogFile(ossID, downloadTaskId, logInfo.Name)
		}
	}()

	return
}

// genDownloadFile 生成历史审计日志下载文件
func (h *HistoryLog) genDownloadFile(ctx context.Context, logInfo *models.HistoryPO) (fileContent []byte, err error) {
	prefix, err := h.logStrategyRepo.GetLogPrefix()
	if err != nil {
		h.logger.Errorf("[genDownloadFile] get log prefix error: %v", err)
		return
	}

	objectName := logInfo.ID
	if prefix != "" {
		objectName = fmt.Sprintf("%s/%s", prefix, objectName)
	}

	dlInfo, _, err := h.ossGateway.GetDownLoadInfo(logInfo.OssID, objectName, logInfo.Name, true)
	if err != nil {
		h.logger.Errorf("[genDownloadFile] get download info error: %v", err)
		return
	}

	readCount := int64(0)
	readBlockSize := int64(5 * 1024 * 1024) // 5MB

	// 写入Excel BOM头
	fileContent = append(fileContent, []byte{0xEF, 0xBB, 0xBF}...)

	// 如果是CSV文件，写入表头
	if strings.HasSuffix(logInfo.Name, ".csv") {
		headers := []string{
			locale.GetI18nCtx(ctx, locale.RCLogDate),
			locale.GetI18nCtx(ctx, locale.RCLogUser),
			locale.GetI18nCtx(ctx, locale.RCLogUserPaths),
			locale.GetI18nCtx(ctx, locale.RCLogLevel),
			locale.GetI18nCtx(ctx, locale.RCLogOperation),
			locale.GetI18nCtx(ctx, locale.RCLogIP),
			locale.GetI18nCtx(ctx, locale.RCLogMac),
			locale.GetI18nCtx(ctx, locale.RCLogMsg),
			locale.GetI18nCtx(ctx, locale.RCLogExMsg),
			locale.GetI18nCtx(ctx, locale.RCLogUserAgent),
			locale.GetI18nCtx(ctx, locale.RCLogAdditionalInfo),
			locale.GetI18nCtx(ctx, locale.RCLogObjName),
			locale.GetI18nCtx(ctx, locale.RCLogObjType),
		}
		fileContent = append(fileContent, []byte(strings.Join(headers, ","))...)
		fileContent = append(fileContent, []byte("\n")...)
	}

	// 分块下载并写入文件
	for blockNum := int64(0); blockNum < (logInfo.Size+readBlockSize-1)/readBlockSize; blockNum++ {
		start := blockNum * readBlockSize
		end := common.Min(start+readBlockSize, logInfo.Size)

		rsp, _, err := h.ossGateway.DownloadBlockByURL(
			dlInfo.URL,
			dlInfo.Method,
			dlInfo.RequestBody,
			dlInfo.Headers,
			start,
			end-1,
		)
		if err != nil {
			h.logger.Errorf("[genDownloadFile] download block by url error: %v", err)
			return nil, err
		}

		body, err := io.ReadAll(rsp.Body)
		if err != nil {
			h.logger.Errorf("[genDownloadFile] read block body error: %v", err)
			return nil, err
		}

		fileContent = append(fileContent, body...)
		readCount += int64(len(body))
	}

	// 验证下载的数据大小
	if logInfo.Size != readCount {
		cause := fmt.Sprintf("download history log failed, expected size: %d, actual size: %d", logInfo.Size, readCount)
		return nil, errors.NewCtx(ctx, errors.InternalErr, cause, nil)
	}

	return fileContent, nil
}

// deleteLogFile 删除oss上压缩后的历史审计日志文件
func (h *HistoryLog) deleteLogFile(ossID string, fileId string, fileName string) {
	delInfo, _, err := h.ossGateway.GetDeleteRequestInfo(ossID, fileId)
	if err != nil {
		h.logger.Errorf("[deleteLogFile] get delete info error: %v", err)
		return
	}

	_, _, err = h.ossGateway.DeleteObjectByURL(delInfo.URL, delInfo.Method, delInfo.RequestBody, delInfo.Headers)
	if err != nil {
		h.logger.Errorf("[deleteLogFile] delete zip file of %s from oss error: %v", fileName, err)
		return
	}

	h.logger.Infof("[deleteLogFile] delete zip file of %s from oss success, ossID: %s, fileId: %s", fileName, ossID, fileId)
}
