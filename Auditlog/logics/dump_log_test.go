package logics

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	configMock "AuditLog/infra/config/mock"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/test/mock_log"
	"AuditLog/test/mock_trace"
)

func newDumpLogDependencies(t *testing.T) (
	*mock_log.MockLogger,
	*mock_trace.MockTracer,
	*mock.MockDLM,
	*configMock.MockLogDumpConfigRepo,
	*mock.MockLogRepo,
	*mock.MockLogRepo,
	*mock.MockLogRepo,
	*mock.MockHistoryRepo,
	*mock.MockOssGatewayRepo,
	*mock.MockLogStrategyRepo,
	*mock.MockLogMgnt,
) {
	ctrl := gomock.NewController(t)
	return mock_log.NewMockLogger(ctrl),
		mock_trace.NewMockTracer(ctrl),
		mock.NewMockDLM(ctrl),
		configMock.NewMockLogDumpConfigRepo(ctrl),
		mock.NewMockLogRepo(ctrl),
		mock.NewMockLogRepo(ctrl),
		mock.NewMockLogRepo(ctrl),
		mock.NewMockHistoryRepo(ctrl),
		mock.NewMockOssGatewayRepo(ctrl),
		mock.NewMockLogStrategyRepo(ctrl),
		mock.NewMockLogMgnt(ctrl)
}

func newDumpLog(
	logger *mock_log.MockLogger,
	tracer *mock_trace.MockTracer,
	dlmLock *mock.MockDLM,
	logDumpConfig *configMock.MockLogDumpConfigRepo,
	loginLogRepo *mock.MockLogRepo,
	mgntLogRepo *mock.MockLogRepo,
	operLogRepo *mock.MockLogRepo,
	historyRepo *mock.MockHistoryRepo,
	ossGateway *mock.MockOssGatewayRepo,
	logStrategyRepo *mock.MockLogStrategyRepo,
	logMgnt *mock.MockLogMgnt,
) *DumpLog {
	return &DumpLog{
		logger:          logger,
		tracer:          tracer,
		dlmLock:         dlmLock,
		logDumpConfig:   logDumpConfig,
		loginLogRepo:    loginLogRepo,
		mgntLogRepo:     mgntLogRepo,
		operLogRepo:     operLogRepo,
		historyLogRepo:  historyRepo,
		ossGateway:      ossGateway,
		logStrategyRepo: logStrategyRepo,
		logMgnt:         logMgnt,
	}
}

func TestGetWaitTimeTillNextDumpAlt(t *testing.T) {
	Convey("GetWaitTimeTillNextDumpAlt", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)

		Convey("计算等待时间成功", func() {
			now := time.Now()
			dumpTime := "00:00:00"

			duration := dumpLog.getWaitTimeTillNextDumpAlt(dumpTime)

			// 如果当前时间已过今天的转储时间点,应该等到明天
			if now.Hour() > 0 || (now.Hour() == 0 && (now.Minute() > 0 || now.Second() > 0)) {
				assert.True(t, duration > 0)
				assert.True(t, duration <= 24*time.Hour)
			} else {
				// 否则应该等到今天的转储时间点
				assert.True(t, duration >= 0)
				assert.True(t, duration < 24*time.Hour)
			}
		})
	})
}

func TestDumpLog(t *testing.T) {
	Convey("DumpLog", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)

		Convey("转储日志成功", func() {
			ctx := context.Background()
			beginTime := int64(1734571503937827)

			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Return().AnyTimes()

			// Mock获取转储格式
			logStrategyRepo.EXPECT().GetDumpFormat().Return("csv", nil)

			// Mock获取日志数据
			logs := []*models.LogPO{
				{
					LogID:    "1",
					Date:     beginTime,
					UserName: "test_user",
					Level:    1,
					OpType:   2,
					IP:       "127.0.0.1",
					MAC:      "00:00:00:00:00:00",
					Msg:      "test_msg",
				},
			}
			loginLogRepo.EXPECT().FindByCondition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(logs, nil)

			// Mock OSS相关操作
			ossGateway.EXPECT().GetAvailableOSSID().Return("test_oss_id", nil)
			logStrategyRepo.EXPECT().GetLogPrefix().Return("test_prefix", nil)
			ossGateway.EXPECT().GetUploadInfo(gomock.Any(), gomock.Any()).Return(&models.OSSUploadInfo{
				UploadID: "test_upload_id",
				PartSize: 1024 * 1024,
			}, 200, nil)

			// Mock上传分块
			ossGateway.EXPECT().GetUploadPartRequestInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.OSSRequestInfo{
				URL:     "test_url",
				Method:  "PUT",
				Headers: map[string]string{},
			}, 200, nil).AnyTimes()
			ossGateway.EXPECT().UploadPartByURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.OSSUploadPartInfo{
				Etag: "test_etag",
				Size: 100,
			}, 200, nil).AnyTimes()

			// Mock完成上传
			ossGateway.EXPECT().GetCompleteUploadRequestInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.OSSRequestInfo{
				URL:         "test_url",
				Method:      "POST",
				RequestBody: "",
				Headers:     map[string]string{},
			}, 200, nil)
			ossGateway.EXPECT().CompleteUploadByURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&http.Response{
				StatusCode: 200,
			}, 200, nil)

			logDumpConfig.EXPECT().GetDumpLogNum().Return(int64(100))
			logDumpConfig.EXPECT().GetDumpIntervalTime().Return(int64(100))

			// Mock保存历史记录
			historyRepo.EXPECT().New(gomock.Any()).Return(nil)

			// Mock清理日志
			loginLogRepo.EXPECT().ClearOutdatedLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			// Mock发送管理日志
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil).AnyTimes()

			err := dumpLog.dumpLog(ctx, common.Login, time.UnixMicro(beginTime), time.Time{})
			assert.NoError(t, err)
		})
	})
}

func TestGetLogCountByType(t *testing.T) {
	Convey("GetLogCountByType", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)

		Convey("获取登录日志数量", func() {
			loginLogRepo.EXPECT().GetLogCount().Return(int64(100), nil)
			count, err := dumpLog.getLogCountByType(common.Login)
			assert.NoError(t, err)
			assert.Equal(t, int64(100), count)
		})

		Convey("获取管理日志数量", func() {
			mgntLogRepo.EXPECT().GetLogCount().Return(int64(200), nil)
			count, err := dumpLog.getLogCountByType(common.Management)
			assert.NoError(t, err)
			assert.Equal(t, int64(200), count)
		})

		Convey("获取操作日志数量", func() {
			operLogRepo.EXPECT().GetLogCount().Return(int64(300), nil)
			count, err := dumpLog.getLogCountByType(common.Operation)
			assert.NoError(t, err)
			assert.Equal(t, int64(300), count)
		})
	})
}

func TestInitDumpLog(t *testing.T) {
	Convey("InitDumpLog", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)

		Convey("初始化 accountID 成功", func() {
			historyLogs := []*models.HistoryPO{
				{
					ID:   "test_account_id",
					Type: int8(common.LogTypeMap[common.Other]),
				},
			}
			historyRepo.EXPECT().GetHistoryLogsByType(int8(common.LogTypeMap[common.Other])).Return(historyLogs, nil)

			err := dumpLog.initAccountID()
			So(err, ShouldBeNil)
			So(dumpLog.accountID, ShouldEqual, "test_account_id")
		})

		Convey("初始化 accountID - 无历史记录", func() {
			historyRepo.EXPECT().GetHistoryLogsByType(int8(common.LogTypeMap[common.Other])).Return([]*models.HistoryPO{}, nil)
			historyRepo.EXPECT().New(gomock.Any()).Return(nil)

			err := dumpLog.initAccountID()
			So(err, ShouldBeNil)
			So(dumpLog.accountID, ShouldNotBeEmpty)
		})
	})
}

func TestCheckDumpExpiredLog(t *testing.T) {
	Convey("CheckDumpExpiredLog", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)
		ctx := context.Background()

		Convey("检查过期日志 - 日期单位", func() {
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			logStrategyRepo.EXPECT().GetRetentionPeriod().Return(7, nil)
			logStrategyRepo.EXPECT().GetRetentionPeriodUnit().Return("day", nil)

			// Mock 所有类型日志的首条时间
			firstLogTime := time.Now().AddDate(0, 0, -10) // 10天前
			loginLogRepo.EXPECT().GetFirstLogTime().Return(firstLogTime.UnixMicro(), nil).AnyTimes()
			mgntLogRepo.EXPECT().GetFirstLogTime().Return(firstLogTime.UnixMicro(), nil).AnyTimes()
			operLogRepo.EXPECT().GetFirstLogTime().Return(firstLogTime.UnixMicro(), nil).AnyTimes()

			// Mock 转储相关操作
			logStrategyRepo.EXPECT().GetDumpFormat().Return("csv", nil).AnyTimes()
			ossGateway.EXPECT().GetAvailableOSSID().Return("test_oss_id", nil).AnyTimes()
			logStrategyRepo.EXPECT().GetLogPrefix().Return("test_prefix", nil).AnyTimes()

			// Mock 上传相关操作
			ossGateway.EXPECT().GetUploadInfo(gomock.Any(), gomock.Any()).Return(&models.OSSUploadInfo{
				UploadID: "test_upload_id",
				PartSize: 1024 * 1024,
			}, 200, nil).AnyTimes()

			// Mock 查找日志条件
			loginLogRepo.EXPECT().FindByCondition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*models.LogPO{}, nil).AnyTimes()
			mgntLogRepo.EXPECT().FindByCondition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*models.LogPO{}, nil).AnyTimes()
			operLogRepo.EXPECT().FindByCondition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*models.LogPO{}, nil).AnyTimes()

			// Mock 上传部分和完成上传
			ossGateway.EXPECT().GetUploadPartRequestInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&models.OSSRequestInfo{
					URL:     "test_url",
					Method:  "PUT",
					Headers: map[string]string{},
				}, 200, nil).AnyTimes()

			ossGateway.EXPECT().UploadPartByURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&models.OSSUploadPartInfo{
					Etag: "test_etag",
					Size: 100,
				}, 200, nil).AnyTimes()

			ossGateway.EXPECT().GetCompleteUploadRequestInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&models.OSSRequestInfo{
					URL:     "test_url",
					Method:  "POST",
					Headers: map[string]string{},
				}, 200, nil).AnyTimes()

			ossGateway.EXPECT().CompleteUploadByURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&http.Response{StatusCode: 200}, 200, nil).AnyTimes()

			// Mock 历史记录
			historyRepo.EXPECT().New(gomock.Any()).Return(nil).AnyTimes()

			// Mock 清理日志
			loginLogRepo.EXPECT().ClearOutdatedLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			mgntLogRepo.EXPECT().ClearOutdatedLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			operLogRepo.EXPECT().ClearOutdatedLog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()

			// Mock 发送管理日志
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil).AnyTimes()

			dumpLog.checkDumpExpiredLog(ctx)
		})
	})
}

func TestGetFirstLogTime(t *testing.T) {
	Convey("GetFirstLogTime", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)

		Convey("获取登录日志首条时间", func() {
			expectedTime := time.Now().UnixMicro()
			loginLogRepo.EXPECT().GetFirstLogTime().Return(expectedTime, nil)

			firstTime, err := dumpLog.getFirstLogTime(common.Login)
			So(err, ShouldBeNil)
			So(firstTime.UnixMicro(), ShouldEqual, expectedTime)
		})

		Convey("获取管理日志首条时间", func() {
			expectedTime := time.Now().UnixMicro()
			mgntLogRepo.EXPECT().GetFirstLogTime().Return(expectedTime, nil)

			firstTime, err := dumpLog.getFirstLogTime(common.Management)
			So(err, ShouldBeNil)
			So(firstTime.UnixMicro(), ShouldEqual, expectedTime)
		})

		Convey("无日志记录", func() {
			operLogRepo.EXPECT().GetFirstLogTime().Return(int64(-1), nil)

			firstTime, err := dumpLog.getFirstLogTime(common.Operation)
			So(err, ShouldBeNil)
			So(firstTime.IsZero(), ShouldBeTrue)
		})
	})
}

func TestGetDumpFileName(t *testing.T) {
	Convey("GetDumpFileName", t, func() {
		logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt := newDumpLogDependencies(t)
		dumpLog := newDumpLog(logger, tracer, dlmLock, logDumpConfig, loginLogRepo, mgntLogRepo, operLogRepo, historyRepo, ossGateway, logStrategyRepo, logMgnt)
		ctx := context.Background()

		Convey("生成单日文件名", func() {
			beginTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.Local)
			endTime := time.Date(2024, 3, 15, 23, 59, 59, 0, time.Local)

			name, _, _ := dumpLog.getDumpFileName(ctx, common.Login, "csv", beginTime, endTime)
			So(name, ShouldContainSubstring, "2024-3-15")
		})

		Convey("生成日期范围文件名", func() {
			beginTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.Local)
			endTime := time.Date(2024, 3, 20, 23, 59, 59, 0, time.Local)

			name, _, _ := dumpLog.getDumpFileName(ctx, common.Login, "csv", beginTime, endTime)
			So(name, ShouldContainSubstring, "2024-3-15~2024-3-20")
		})
	})
}
