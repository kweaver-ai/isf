package logics

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"AuditLog/common"
	"AuditLog/common/constants"
	"AuditLog/common/constants/logconsts"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/common/constants/rclogconsts"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/common/utils/dumplogutils"
	"AuditLog/common/utils/rclogutils"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
	"AuditLog/models/rcvo"
)

var (
	dlOnce sync.Once
	dl     *DumpLog
)

type DumpLog struct {
	accountID       string
	logger          api.Logger
	tracer          api.Tracer
	dlmLock         interfaces.DLM
	logDumpConfig   config.LogDumpConfigRepo
	loginLogRepo    interfaces.LogRepo     // 数据库对象
	mgntLogRepo     interfaces.LogRepo     // 数据库对象
	operLogRepo     interfaces.LogRepo     // 数据库对象
	historyLogRepo  interfaces.HistoryRepo // 数据库对象
	ossGateway      interfaces.OssGatewayRepo
	logStrategyRepo interfaces.LogStrategyRepo
	logMgnt         interfaces.LogMgnt
}

func NewDumpLog() interfaces.DumpLog {
	dlOnce.Do(func() {
		dl = &DumpLog{
			accountID:       "",
			logger:          logger,
			tracer:          tracer,
			dlmLock:         dlmLock,
			logDumpConfig:   config.GetLogDumpConfig(),
			loginLogRepo:    loginLogRepo,
			mgntLogRepo:     mgntLogRepo,
			operLogRepo:     operLogRepo,
			historyLogRepo:  historyRepo,
			ossGateway:      ossGateway,
			logStrategyRepo: logStrategyRepo,
			logMgnt:         NewLogMgnt(),
		}
	})

	return dl
}

func (d *DumpLog) InitDumpLog(ctx context.Context) {
	if helpers.IsAaronLocalDev() {
		return
	}

	// 初始化accountID
	if err := d.initAccountID(); err != nil {
		d.logger.Warnf("[InitDumpLog] init account id error: %v", err)
		return
	}

	for {
		dumpTime, err := d.logStrategyRepo.GetDumpTime()
		if err != nil {
			d.logger.Warnf("[InitDumpLog] get dump time error: %v", err)
			panic(err)
		}
		d.logger.Infof("[InitDumpLog] get dump time success, dump time %s", dumpTime)

		waitDuration := d.getWaitTimeTillNextDumpAlt(dumpTime)
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitDuration):
			{
				// 获取锁
				locked, err := d.dlmLock.TryLock(constants.DumpLogLockKey)
				if err != nil {
					d.logger.Warnf("dlm lock error: %v", err)
					break
				}
				if !locked {
					continue
				}

				d.checkDumpExpiredLog(ctx)
				d.checkDumpOutOfRangeLog(ctx)

				if err := d.dlmLock.UnLock(constants.DumpLogLockKey); err != nil {
					d.logger.Warnf("[InitDumpLog] dlm unlock error: %v", err)
				}
			}
		}
	}
}

// 初始化accountID
func (d *DumpLog) initAccountID() (err error) {
	if d.accountID == "" {
		historyLogs, err := d.historyLogRepo.GetHistoryLogsByType(int8(common.LogTypeMap[common.Other]))
		if err != nil {
			d.logger.Warnf("[InitDumpLog] get history log error: %v", err)
			return err
		}

		if len(historyLogs) == 0 {
			id := common.GenerateID()

			err = d.historyLogRepo.New(&models.HistoryPO{
				ID:       id,
				Name:     "",
				Size:     -1,
				Type:     int8(common.LogTypeMap[common.Other]),
				Date:     0,
				DumpDate: 0,
				OssID:    "",
			})
			if err != nil {
				d.logger.Warnf("[InitDumpLog] create history log error: %v", err)
				return err
			}

			d.accountID = id
		} else {
			d.accountID = historyLogs[0].ID
		}
	}

	return nil
}

// 获取下一次转储时间
func (d *DumpLog) getWaitTimeTillNextDumpAlt(dumpTime string) time.Duration {
	DumpTimeHour, DumpTimeMin, DumpTimeSec, err := utils.ParseTime(dumpTime)
	if err != nil {
		d.logger.Warnf("[getWaitTimeTillNextDump] parse dump time error: %v", err)
		return 0
	}

	now := time.Now()

	// 创建今天的转储时间点
	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		DumpTimeHour, DumpTimeMin, DumpTimeSec,
		0, now.Location(),
	)

	// 如果当前时间已经过了今天的转储时间，则等待到明天的转储时间
	if now.After(today) {
		today = today.Add(24 * time.Hour)
	}

	return today.Sub(now)
}

// 检查日志是否过期
func (d *DumpLog) checkDumpExpiredLog(ctx context.Context) {
	retentionPeriod, err := d.logStrategyRepo.GetRetentionPeriod()
	if err != nil {
		d.logger.Warnf("[checkDumpExpiredLog] get retention period error: %v", err)
		return
	}

	retentionPeriodUnit, err := d.logStrategyRepo.GetRetentionPeriodUnit()
	if err != nil {
		d.logger.Warnf("[checkDumpExpiredLog] get retention period unit error: %v", err)
		return
	}

	for _, logType := range common.AllLogType {
		firstLogTime, err := d.getFirstLogTime(logType)
		if err != nil {
			d.logger.Warnf("[checkDumpExpiredLog] get first log time error: %v", err)
			continue
		}

		if firstLogTime.IsZero() {
			continue
		}

		firstLogTimeLocal := firstLogTime.Local()
		dayStart := time.Date(firstLogTimeLocal.Year(), firstLogTimeLocal.Month(), firstLogTimeLocal.Day(), 0, 0, 0, 0, firstLogTimeLocal.Location())
		var periodStartTime time.Time
		var periodEndTime time.Time
		switch retentionPeriodUnit {
		case lsconsts.Day:
			periodStartTime = dayStart
			periodEndTime = dayStart.AddDate(0, 0, retentionPeriod)
		case lsconsts.Week:
			weekStart := dayStart.AddDate(0, 0, -int(dayStart.Weekday()))
			periodStartTime = weekStart
			periodEndTime = weekStart.AddDate(0, 0, retentionPeriod*7)
		case lsconsts.Month:
			monthStart := time.Date(dayStart.Year(), dayStart.Month(), 1, 0, 0, 0, 0, dayStart.Location())
			periodStartTime = monthStart
			periodEndTime = monthStart.AddDate(0, retentionPeriod, 0)
		case lsconsts.Year:
			yearStart := time.Date(dayStart.Year(), 1, 1, 0, 0, 0, 0, dayStart.Location())
			periodStartTime = yearStart
			periodEndTime = yearStart.AddDate(retentionPeriod, 0, 0)
		}

		if time.Now().After(periodEndTime) {
			// 日志已过期，执行相应操作
			d.logger.Infof("[checkDumpExpiredLog] log type %s has expired", logType)
			if err := d.dumpLog(ctx, logType, periodStartTime.UTC(), periodEndTime.UTC()); err != nil {
				d.logger.Warnf("[checkDumpExpiredLog] dump log error: %v", err)
			}
		}
	}
}

// 检查日志是否超出范围
func (d *DumpLog) checkDumpOutOfRangeLog(ctx context.Context) {
	for _, logType := range common.AllLogType {
		count, err := d.getLogCountByType(logType)
		if err != nil {
			d.logger.Warnf("[checkDumpOutOfRangeLog] get log count error: %v", err)
			continue
		}

		threshold := d.logDumpConfig.GetDumpThresholdByType(logType)
		if count < threshold {
			continue
		}

		// 循环处理直到日志数量低于阈值的一半
		for count >= threshold/2 {
			firstLogTime, err := d.getFirstLogTime(logType)
			if err != nil {
				d.logger.Warnf("[checkDumpOutOfRangeLog] get first log time error: %v", err)
				continue
			}

			if firstLogTime.IsZero() {
				break
			}

			// 计算结束时间
			localBegin := firstLogTime.Local()
			var localEnd time.Time
			if localBegin.Month() == time.December {
				localEnd = time.Date(localBegin.Year()+1, time.January, 1, 0, 0, 0, 0, localBegin.Location())
			} else {
				localEnd = time.Date(localBegin.Year(), localBegin.Month()+1, 1, 0, 0, 0, 0, localBegin.Location())
			}

			if err = d.dumpLog(ctx, logType, firstLogTime, localEnd.UTC()); err != nil {
				d.logger.Warnf("[checkDumpOutOfRangeLog] log type: %s, dump error: %v", logType, err)
				break
			}

			count, err = d.getLogCountByType(logType)
			if err != nil {
				d.logger.Warnf("[checkDumpOutOfRangeLog] get log count error: %v", err)
				break
			}
		}
	}
}

// 转储日志
func (d *DumpLog) dumpLog(ctx context.Context, logType string, beginLogTime time.Time, endLogTime time.Time) (err error) {
	var suffix string = "csv"
	if s, err := d.logStrategyRepo.GetDumpFormat(); err == nil {
		suffix = s
	} else {
		d.logger.Warnf("[dumpLog] get retention format error: %v, use default suffix %s", err, suffix)
	}

	fileName, localBegin, localEnd := d.getDumpFileName(ctx, logType, suffix, beginLogTime, endLogTime.Add(-time.Nanosecond))

	logs, err := d.getPeriodsOfLogLimit(logType, -1, beginLogTime, endLogTime, lsconsts.HistoryMaxBatchSize)
	if err != nil {
		d.logger.Warnf("[dumpLog] get periods of log limit error: %v", err)
		return
	}

	if len(logs) == 0 {
		return
	}

	endLogID, err := strconv.Atoi(logs[0].LogID)
	if err != nil {
		d.logger.Warnf("[dumpLog] convert end log id to int error: %v", err)
		return
	}

	d.logger.Infof(
		"[dumpLog] Begin to dump log, file_name: %s, records: %d, end_log_id: %d, begin_time: %s, end_time: %s",
		fileName, len(logs), endLogID, localBegin, localEnd,
	)

	// 获取可用的OSS ID
	ossID, err := d.ossGateway.GetAvailableOSSID()
	if err != nil {
		d.logger.Warnf("[dumpLog] get available oss id error: %v", err)
		return
	}

	// 生成对象ID和前缀
	objID := common.GenerateID()
	prefix, err := d.getLogPrefix()
	if err != nil {
		d.logger.Warnf("[dumpLog] get log prefix error: %v", err)
		return
	}
	// 构建对象名称
	var objectName string
	if prefix != "" {
		objectName = fmt.Sprintf("%s/%s/%s", prefix, d.accountID, objID)
	} else {
		objectName = fmt.Sprintf("%s/%s", d.accountID, objID)
	}

	// 初始化上传
	uploadInfo, _, err := d.ossGateway.GetUploadInfo(ossID, objectName)
	if err != nil {
		d.logger.Errorf("[dumpLog] init upload error: %v", err)
		return
	}

	// 分片信息
	partInfos := make(map[int]*models.OSSUploadPartInfo)
	size, logCount := int64(0), int64(len(logs))
	sn := 1
	content := strings.Builder{}
	content.Grow(int(uploadInfo.PartSize))

	if suffix == lsconsts.XMLSuffix {
		content.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<log>\n")
	}

	for {
		for i, log := range logs {
			var str string
			if suffix == lsconsts.XMLSuffix {
				str, err = dumplogutils.LogInfo2XMLString(log, logType)
				if err != nil {
					d.logger.Warnf("[dumpLog] log info to xml string error: %v", err)
					return
				}
			} else {
				str, err = dumplogutils.LogInfo2CSVString(log, logType)
				if err != nil {
					d.logger.Warnf("[dumpLog] log info to string error: %v", err)
					return
				}
			}
			str += "\n"

			// 添加XML尾部
			if suffix == lsconsts.XMLSuffix && len(logs) != lsconsts.HistoryMaxBatchSize && i == len(logs)-1 {
				str += "</log>\n"
			}

			// 检查是否需要上传当前块
			if content.Len()+len(str) > int(uploadInfo.PartSize) {
				appendingSize := int(uploadInfo.PartSize) - content.Len()
				content.WriteString(str[:appendingSize])
				str = str[appendingSize:]

				// 上传数据块
				partInfo, err := d.uploadBlock(
					objectName,
					ossID,
					uploadInfo.UploadID,
					sn,
					content.String(),
				)
				if err != nil {
					d.logger.Warnf("[dumpLog] upload block failed: %v", err)
					return err
				}
				d.logger.Infof("[dumpLog] Uploaded part, account %s, objID %s, sn %d", d.accountID, objID, sn)

				partInfos[sn] = partInfo
				sn++
				size += int64(content.Len())

				content.Reset()
			}

			content.WriteString(str)
		}

		if len(logs) != lsconsts.HistoryMaxBatchSize {
			break
		}

		// 获取下一批日志
		lastLog := logs[len(logs)-1]
		lastLogID, err := strconv.Atoi(lastLog.LogID)
		if err != nil {
			d.logger.Warnf("[dumpLog] convert last log id to int error: %v", err)
			return err
		}

		logs, err = d.getPeriodsOfLogLimit(
			logType,
			lastLogID-1,
			beginLogTime,
			time.UnixMicro(lastLog.Date),
			lsconsts.HistoryMaxBatchSize,
		)
		if err != nil {
			d.logger.Warnf("[dumpLog] get next batch logs failed: %v", err)
			return err
		}
		logCount += int64(len(logs))

		d.logger.Infof("[dumpLog] Got next batch records, count %d", len(logs))
	}

	// 上传最后的数据块
	if content.Len() > 0 {
		partInfo, err := d.uploadBlock(
			objectName,
			ossID,
			uploadInfo.UploadID,
			sn,
			content.String(),
		)
		if err != nil {
			d.logger.Errorf("[dumpLog] upload final block failed: %v", err)
			return err
		}

		partInfos[sn] = partInfo
		size += int64(content.Len())
	}

	// 完成上传
	if err := d.completeUpload(objectName, ossID, uploadInfo.UploadID, partInfos); err != nil {
		d.logger.Errorf("[dumpLog] complete upload failed: %v", err)
		return err
	}

	d.logger.Infof(
		"[dumpLog] dump log success, file_name: %s, records: %d, end_log_id: %d, begin_time: %s, end_time: %s",
		fileName, logCount, endLogID, localBegin, localEnd,
	)

	// 记录历史日志
	historyInfo := &models.HistoryPO{
		ID:       path.Join(d.accountID, objID),
		Name:     fileName,
		Size:     int64(size),
		Type:     int8(common.LogTypeMap[logType]),
		Date:     beginLogTime.UnixMicro(),
		DumpDate: time.Now().UnixMicro(),
		OssID:    ossID,
	}

	if err := d.historyLogRepo.New(historyInfo); err != nil {
		d.logger.Errorf("[dumpLog] add history log failed: %v", err)
		return err
	}

	batchSize := d.logDumpConfig.GetDumpLogNum()
	sleepTime := d.logDumpConfig.GetDumpIntervalTime()

	// 清理过期日志
	switch logType {
	case common.Login:
		err = d.loginLogRepo.ClearOutdatedLog(int64(endLogID), endLogTime.UnixMicro(), batchSize, sleepTime)
	case common.Management:
		err = d.mgntLogRepo.ClearOutdatedLog(int64(endLogID), endLogTime.UnixMicro(), batchSize, sleepTime)
	case common.Operation:
		err = d.operLogRepo.ClearOutdatedLog(int64(endLogID), endLogTime.UnixMicro(), batchSize, sleepTime)
	}

	if err != nil {
		d.logger.Errorf("[dumpLog] clear outdated log failed: %v", err)
		return err
	}

	// 记录管理日志
	go func() {
		userInfo := dumplogutils.GetSystemAccount(ctx)
		if err := d.logMgnt.SendLog(&models.SendLogVo{
			LogType:  common.Management,
			Language: "",
			LogContent: &models.AuditLog{
				UserID:         userInfo.ID,
				UserName:       userInfo.DisplayName,
				UserType:       common.InternalService,
				Level:          logconsts.LogLevel.INFO,
				OpType:         logconsts.OpType.ManagementType.SET,
				Date:           time.Now().UnixMicro(),
				IP:             "127.0.0.1",
				Mac:            "",
				Msg:            fmt.Sprintf(locale.GetI18nCtx(ctx, locale.LogDumpMsg), d.getLogTypeName(ctx, logType)),
				Exmsg:          fmt.Sprintf(locale.GetI18nCtx(ctx, locale.LogDumpExMsg), fileName),
				UserAgent:      "",
				ObjID:          "",
				AdditionalInfo: fmt.Sprintf("{\"user_account\": \"%s\"}", userInfo.Account),
				OutBizID:       uuid.NewString(),
				DeptPaths:      "",
			},
		}); err != nil {
			d.logger.Warnf("[dumpLog] send log error: %v", err)
		}
	}()

	return
}

// 获取过期日志的周期
func (d *DumpLog) getPeriodsOfLogLimit(logType string, maxLogID int, beginLogTime time.Time, endLogTime time.Time, maxBatchSize int) (logs []*models.LogPO, err error) {
	sqlStr, err := rclogutils.BuildActiveCondition2(
		map[string]any{
			rclogconsts.LogID: int(maxLogID),
			rclogconsts.Date: []interface{}{
				float64(beginLogTime.UnixMicro()),
				float64(endLogTime.UnixMicro()),
			},
		},
		rcvo.OrderFields{
			{
				Field:     rclogconsts.LogID,
				Direction: "desc",
			},
		},
		[]string{}, []string{}, []string{},
	)
	if err != nil {
		return nil, err
	}

	switch logType {
	case common.Operation:
		logs, err = d.operLogRepo.FindByCondition(0, maxBatchSize, sqlStr, []string{})
	case common.Management:
		logs, err = d.mgntLogRepo.FindByCondition(0, maxBatchSize, sqlStr, []string{})
	case common.Login:
		logs, err = d.loginLogRepo.FindByCondition(0, maxBatchSize, sqlStr, []string{})
	}
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// 获取转储文件名
func (d *DumpLog) getDumpFileName(ctx context.Context, logType, suffix string, beginLogTime, endLogTime time.Time) (name string, localBegin, localEnd time.Time) {
	// 处理时间
	localBegin = beginLogTime.Local()
	localEnd = endLogTime.Local()

	if localBegin.Year() == localEnd.Year() &&
		localBegin.Month() == localEnd.Month() &&
		localBegin.Day() == localEnd.Day() {
		name = fmt.Sprintf("%s.%d-%d-%d.%s",
			d.getLogTypeName(ctx, logType),
			localBegin.Year(), localBegin.Month(), localBegin.Day(),
			suffix,
		)
	} else {
		name = fmt.Sprintf("%s.%d-%d-%d~%d-%d-%d.%s",
			d.getLogTypeName(ctx, logType),
			localBegin.Year(), localBegin.Month(), localBegin.Day(),
			localEnd.Year(), localEnd.Month(), localEnd.Day(),
			suffix,
		)
	}

	return
}

// 获取日志前缀
func (d *DumpLog) getLogPrefix() (prefix string, err error) {
	prefix, _ = d.logStrategyRepo.GetLogPrefix()

	if prefix == "" {
		prefix = uuid.New().String()
		err = d.logStrategyRepo.SetLogPrefix(prefix)
		if err != nil {
			d.logger.Errorf("[getLogPrefix] set log prefix error: %v", err)
			return
		}
		d.logger.Infof("[getLogPrefix] set log prefix success, prefix %s", prefix)
	}

	return prefix, nil
}

// uploadBlock 上传分块
func (d *DumpLog) uploadBlock(objectName, ossID, uploadID string, partNumber int, data string) (partInfo *models.OSSUploadPartInfo, err error) {
	uploadPartRequestInfo, _, err := d.ossGateway.GetUploadPartRequestInfo(ossID, objectName, uploadID, partNumber)
	if err != nil || uploadPartRequestInfo == nil {
		return nil, fmt.Errorf("get upload part request info failed: %w", err)
	}

	partInfo, _, err = d.ossGateway.UploadPartByURL(
		uploadPartRequestInfo.URL,
		uploadPartRequestInfo.Method,
		data,
		uploadPartRequestInfo.Headers,
	)
	if err != nil || partInfo == nil {
		return nil, fmt.Errorf("upload part failed: %w", err)
	}

	return partInfo, nil
}

// completeUpload 完成上传
func (d *DumpLog) completeUpload(objectName, ossID, uploadID string, partInfos map[int]*models.OSSUploadPartInfo) (err error) {
	// 转换分片信息
	ossPartInfos := make(map[int]models.OSSUploadPartInfo)
	for partNum, info := range partInfos {
		ossPartInfos[partNum] = models.OSSUploadPartInfo{
			Etag: info.Etag,
			Size: info.Size,
		}
	}

	completeUploadRequestInfo, _, err := d.ossGateway.GetCompleteUploadRequestInfo(ossID, objectName, uploadID, ossPartInfos)
	if err != nil || completeUploadRequestInfo == nil {
		return fmt.Errorf("get complete upload request info failed: %w", err)
	}

	comRes, _, err := d.ossGateway.CompleteUploadByURL(
		completeUploadRequestInfo.URL,
		completeUploadRequestInfo.Method,
		completeUploadRequestInfo.RequestBody,
		completeUploadRequestInfo.Headers,
	)
	if err != nil || comRes == nil {
		return fmt.Errorf("complete upload failed: %w", err)
	}

	return nil
}

// 获取第一个日志时间
func (d *DumpLog) getFirstLogTime(logType string) (firstLogTime time.Time, err error) {
	var firstLogTimeMicro int64
	switch logType {
	case common.Login:
		firstLogTimeMicro, err = d.loginLogRepo.GetFirstLogTime()
	case common.Management:
		firstLogTimeMicro, err = d.mgntLogRepo.GetFirstLogTime()
	case common.Operation:
		firstLogTimeMicro, err = d.operLogRepo.GetFirstLogTime()
	}
	if err != nil {
		d.logger.Errorf("[getFirstLogTime] get first log time error: %v", err)
		return time.Time{}, err
	}

	if firstLogTimeMicro == -1 {
		return time.Time{}, nil
	}

	return time.UnixMicro(firstLogTimeMicro), nil
}

// 获取日志类型名称
func (d *DumpLog) getLogTypeName(ctx context.Context, logType string) string {
	switch logType {
	case common.Login:
		return locale.GetI18nCtx(ctx, locale.RCLogReportLogin)
	case common.Management:
		return locale.GetI18nCtx(ctx, locale.RCLogReportMgnt)
	case common.Operation:
		return locale.GetI18nCtx(ctx, locale.RCLogReportOp)
	}

	return ""
}

// 获取日志数量
func (d *DumpLog) getLogCountByType(logType string) (count int64, err error) {
	switch logType {
	case common.Login:
		count, err = d.loginLogRepo.GetLogCount()
	case common.Management:
		count, err = d.mgntLogRepo.GetLogCount()
	case common.Operation:
		count, err = d.operLogRepo.GetLogCount()
	}

	return count, err
}
