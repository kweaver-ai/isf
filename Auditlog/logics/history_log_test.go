package logics

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/errors"
	"AuditLog/infra/cmp/langcmp"
	"AuditLog/interfaces"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/models/rcvo"
	"AuditLog/test/mock_log"
	"AuditLog/test/mock_trace"
)

func TestMain(m *testing.M) {
	langcmp.NewLangCmp().SetSysDefLang(string(langcmp.ZhCN))
	os.Exit(m.Run())
}

func newHistoryDependencies(t *testing.T) (
	*mock_log.MockLogger,
	*mock_trace.MockTracer,
	*mock.MockHistoryRepo,
	*mock.MockUserMgntRepo,
	*mock.MockLogStrategyRepo,
	*mock.MockLogScopeStrategyRepo,
	*mock.MockOssGatewayRepo,
	*mock.MockLogMgnt,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl),
		mock_trace.NewMockTracer(ctrl),
		mock.NewMockHistoryRepo(ctrl),
		mock.NewMockUserMgntRepo(ctrl),
		mock.NewMockLogStrategyRepo(ctrl),
		mock.NewMockLogScopeStrategyRepo(ctrl),
		mock.NewMockOssGatewayRepo(ctrl),
		mock.NewMockLogMgnt(ctrl)
}

func newHistoryLog(
	logger *mock_log.MockLogger,
	tracer *mock_trace.MockTracer,
	cache *redis.Client,
	historyRepo *mock.MockHistoryRepo,
	userMgnt *mock.MockUserMgntRepo,
	logStrategyRepo *mock.MockLogStrategyRepo,
	logScopeStrategyRepo *mock.MockLogScopeStrategyRepo,
	ossGateway *mock.MockOssGatewayRepo,
	logMgnt *mock.MockLogMgnt,
) interfaces.HistoryLog {
	return &HistoryLog{
		logger:               logger,
		tracer:               tracer,
		cache:                cache,
		historyLogRepo:       historyRepo,
		userMgntRepo:         userMgnt,
		logStrategyRepo:      logStrategyRepo,
		logScopeStrategyRepo: logScopeStrategyRepo,
		ossGateway:           ossGateway,
		logMgnt:              logMgnt,
	}
}

func TestGetHistoryDataList(t *testing.T) {
	Convey("GetHistoryDataList", t, func() {
		redisClient, _ := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		Convey("超级管理员获取历史数据列表成功", func() {
			ctx := context.Background()
			category := "login"
			userID := "111"
			req := &rcvo.ReportGetDataListReq{
				Offset:    0,
				Limit:     10,
				Condition: map[string]any{},
				OrderBy:   rcvo.OrderFields{},
				IDs:       []string{},
			}

			// Mock用户信息返回
			userInfo := []models.User{
				{
					ID:    "111",
					Roles: []string{common.SuperAdmin},
					Name:  "admin",
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(userInfo, 200, nil)

			// Mock历史日志数据
			historyLogs := []*models.HistoryPO{
				{
					ID:       "1211",
					DumpDate: 1717670099640000,
					Name:     "test.log",
					Size:     1024,
				},
			}
			historyRepo.EXPECT().FindByCondition(req.Offset, req.Limit, gomock.Any(), req.IDs).Return(historyLogs, nil)
			historyRepo.EXPECT().FindCountByCondition(gomock.Any()).Return(10, nil)

			// Mock tracer
			tracer.EXPECT().AddInternalTrace(gomock.Any()).MaxTimes(1)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infoln(gomock.Any()).AnyTimes()
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			res, err := historyLog.GetHistoryDataList(ctx, category, req, userID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(res.Entries))
			assert.Equal(t, "1211", res.Entries[0].ID)
			assert.Equal(t, 10, res.TotalCount)
		})

		Convey("无权限用户获取历史数据列表失败", func() {
			ctx := context.Background()
			category := "login"
			userID := "222"
			req := &rcvo.ReportGetDataListReq{
				Offset:    0,
				Limit:     10,
				Condition: map[string]any{},
				OrderBy:   rcvo.OrderFields{},
				IDs:       []string{},
			}

			// Mock用户信息返回
			userInfo := []models.User{
				{
					ID:    "222",
					Roles: []string{"normal_user"},
					Name:  "user",
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(userInfo, 200, nil)

			// Mock tracer
			tracer.EXPECT().AddInternalTrace(gomock.Any()).MaxTimes(1)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infoln(gomock.Any()).AnyTimes()
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			res, err := historyLog.GetHistoryDataList(ctx, category, req, userID)
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

func TestGetHistoryMetadata(t *testing.T) {
	Convey("GetHistoryMetadata", t, func() {
		redisClient, _ := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		Convey("获取历史元数据成功", func() {
			meta, err := historyLog.GetHistoryMetadata()
			assert.NoError(t, err)
			assert.NotNil(t, meta)
		})
	})
}

func TestGetHistoryDownloadPwdStatus(t *testing.T) {
	Convey("GetHistoryDownloadPwdStatus", t, func() {
		redisClient, _ := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		Convey("获取下载密码状态成功", func() {
			// Mock 获取密码状态
			logStrategyRepo.EXPECT().GetHistoryIsDownloadWithPwd().Return(true, nil)

			res, err := historyLog.GetHistoryDownloadPwdStatus(context.Background())
			assert.NoError(t, err)
			assert.True(t, res.Status)
		})

		Convey("获取下载密码状态失败", func() {
			// Mock 获取密码状态失败
			logStrategyRepo.EXPECT().GetHistoryIsDownloadWithPwd().Return(false, errors.NewCtx(context.Background(), errors.InternalErr, "mock error", nil))

			res, err := historyLog.GetHistoryDownloadPwdStatus(context.Background())
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

func TestSetHistoryDownloadPwdStatus(t *testing.T) {
	Convey("SetHistoryDownloadPwdStatus", t, func() {
		redisClient, _ := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		ctx := context.WithValue(context.Background(), common.VisitorKey, &models.Visitor{
			ID:        "test_user",
			Name:      "Test User",
			IP:        "127.0.0.1",
			Mac:       "00:00:00:00:00:00",
			AgentType: "test_agent",
		})

		Convey("设置下载密码状态成功", func() {
			req := &lsmodels.HistoryLogDownloadPwdStatus{
				Status: true,
			}

			// Mock 设置密码状态
			logStrategyRepo.EXPECT().SetHistoryIsDownloadWithPwd(true).Return(nil).AnyTimes()

			// Mock 发送日志
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil).AnyTimes()

			err := historyLog.SetHistoryDownloadPwdStatus(ctx, req)
			assert.NoError(t, err)
		})

		Convey("设置下载密码状态失败", func() {
			req := &lsmodels.HistoryLogDownloadPwdStatus{
				Status: true,
			}

			// Mock 设置密码状态失败
			logStrategyRepo.EXPECT().SetHistoryIsDownloadWithPwd(true).Return(errors.NewCtx(context.Background(), errors.InternalErr, "mock error", nil)).AnyTimes()

			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			err := historyLog.SetHistoryDownloadPwdStatus(ctx, req)
			assert.Error(t, err)
		})
	})
}

func TestGetDownloadProgress(t *testing.T) {
	Convey("GetDownloadProgress", t, func() {
		redisClient, redisMock := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		taskId := "test_task_id"

		Convey("获取下载进度-任务进行中", func() {
			// Mock Redis缓存返回pending状态
			redisMock.ExpectGet(lsconsts.TaskCacheKey + taskId).SetVal(`{"status": 0}`)

			res, err := historyLog.GetDownloadProgress(context.Background(), taskId)
			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.False(t, res.Status)
		})

		Convey("获取下载进度-任务不存在", func() {
			// Mock Redis缓存返回nil
			redisMock.ExpectGet(lsconsts.TaskCacheKey + taskId).RedisNil()

			res, err := historyLog.GetDownloadProgress(context.Background(), taskId)
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

func TestCompressLog(t *testing.T) {
	Convey("CompressLog", t, func() {
		redisClient, redisMock := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
		ctx := context.Background()
		downloadTaskId := "test_task_id"
		logInfo := &models.HistoryPO{
			ID:       "test_id",
			Name:     "test.csv",
			Size:     1024,
			OssID:    "test_oss_id",
			DumpDate: time.Now().UnixMicro(),
		}
		req := &lsmodels.HistoryLogDownloadReq{
			ObjId: "test_id",
		}

		Convey("压缩日志成功", func() {
			// Mock获取日志前缀
			logStrategyRepo.EXPECT().GetLogPrefix().Return("test_prefix", nil)

			// Mock获取下载信息
			ossGateway.EXPECT().GetDownLoadInfo(logInfo.OssID, gomock.Any(), logInfo.Name, true).
				Return(&models.OSSRequestInfo{
					URL:     "test_url",
					Method:  "GET",
					Headers: map[string]string{},
				}, 200, nil)

			// Mock下载块
			ossGateway.EXPECT().DownloadBlockByURL(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte("test content"))),
			}, 200, nil).AnyTimes()

			// Mock获取OSS ID
			ossGateway.EXPECT().GetAvailableOSSID().Return("test_oss_id", nil)

			// Mock获取上传信息
			ossGateway.EXPECT().GetUploadInfo("test_oss_id", downloadTaskId).
				Return(&models.OSSUploadInfo{
					UploadID: "test_upload_id",
					PartSize: 1024 * 1024,
				}, 200, nil)

			// Mock分块上传
			ossGateway.EXPECT().GetUploadPartRequestInfo(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&models.OSSRequestInfo{
				URL:     "test_url",
				Method:  "PUT",
				Headers: map[string]string{},
			}, 200, nil).AnyTimes()

			ossGateway.EXPECT().UploadPartByURL(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&models.OSSUploadPartInfo{
				Etag: "test_etag",
				Size: 100,
			}, 200, nil).AnyTimes()

			// Mock完成上传
			ossGateway.EXPECT().GetCompleteUploadRequestInfo(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&models.OSSRequestInfo{
				URL:     "test_url",
				Method:  "POST",
				Headers: map[string]string{},
			}, 200, nil)

			ossGateway.EXPECT().CompleteUploadByURL(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&http.Response{StatusCode: 200}, 200, nil)

			// Mock Redis缓存操作
			redisMock.ExpectSet(lsconsts.TaskCacheKey+"xxx", gomock.Any(), lsconsts.TaskInfoExpire).SetVal("OK")

			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			// 将 historyLog 转换为具体类型 *HistoryLog
			concreteHistoryLog := historyLog.(*HistoryLog)
			concreteHistoryLog.compressLog(ctx, downloadTaskId, req, logInfo)
		})

		Convey("压缩日志失败 - 获取文件内容失败", func() {
			// Mock获取日志前缀失败
			logStrategyRepo.EXPECT().GetLogPrefix().Return("", errors.NewCtx(ctx, errors.InternalErr, "mock error", nil))

			// Mock Redis缓存操作
			redisMock.ExpectSet(lsconsts.TaskCacheKey+"xxx", gomock.Any(), lsconsts.TaskInfoExpire).SetVal("OK")

			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

			// 将 historyLog 转换为具体类型 *HistoryLog
			concreteHistoryLog := historyLog.(*HistoryLog)
			concreteHistoryLog.compressLog(ctx, downloadTaskId, req, logInfo)
		})
	})
}

func TestDeleteLogFile(t *testing.T) {
	Convey("DeleteLogFile", t, func() {
		redisClient, _ := redismock.NewClientMock()
		logger, tracer, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt := newHistoryDependencies(t)
		historyLog := newHistoryLog(logger, tracer, redisClient, historyRepo, userMgnt, logStrategyRepo, logScopeStrategyRepo, ossGateway, logMgnt)

		ossID := "test_oss_id"
		fileId := "test_file_id"
		fileName := "test.csv"

		Convey("删除文件成功", func() {
			// Mock获取删除信息
			ossGateway.EXPECT().GetDeleteRequestInfo(ossID, fileId).
				Return(&models.OSSRequestInfo{
					URL:     "test_url",
					Method:  "DELETE",
					Headers: map[string]string{},
				}, 200, nil)

			// Mock删除文件
			ossGateway.EXPECT().DeleteObjectByURL(
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(&http.Response{StatusCode: 200}, 200, nil)

			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			// 将 historyLog 转换为具体类型 *HistoryLog
			concreteHistoryLog := historyLog.(*HistoryLog)
			concreteHistoryLog.deleteLogFile(ossID, fileId, fileName)
		})

		Convey("删除文件失败", func() {
			// Mock获取删除信息失败
			ossGateway.EXPECT().GetDeleteRequestInfo(ossID, fileId).
				Return(nil, 0, errors.NewCtx(context.Background(), errors.InternalErr, "mock error", nil))

			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

			// 将 historyLog 转换为具体类型 *HistoryLog
			concreteHistoryLog := historyLog.(*HistoryLog)
			concreteHistoryLog.deleteLogFile(ossID, fileId, fileName)
		})
	})
}
