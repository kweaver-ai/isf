package logics

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/test/mock_log"
	"AuditLog/test/mock_trace"
)

func newLogDependencies(t *testing.T) (
	*mock_log.MockLogger,
	*mock_trace.MockTracer,
	*mock.MockLogStrategyRepo,
	*mock.MockLogMgnt,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl),
		mock_trace.NewMockTracer(ctrl),
		mock.NewMockLogStrategyRepo(ctrl),
		mock.NewMockLogMgnt(ctrl)
}

func newLogStrategy(
	logger *mock_log.MockLogger,
	tracer *mock_trace.MockTracer,
	lsRepo *mock.MockLogStrategyRepo,
	logMgnt *mock.MockLogMgnt,
) *LogStrategy {
	return &LogStrategy{
		logger:          logger,
		tracer:          tracer,
		logStrategyRepo: lsRepo,
		logMgnt:         logMgnt,
	}
}

func TestGetDumpStrategy(t *testing.T) {
	Convey("GetDumpStrategy", t, func() {
		logger, tracer, lsRepo, logMgnt := newLogDependencies(t)
		logStrategy := newLogStrategy(logger, tracer, lsRepo, logMgnt)

		Convey("获取所有转存策略成功", func() {
			ctx := context.Background()

			// Mock返回值
			lsRepo.EXPECT().GetRetentionPeriod().Return(30, nil)
			lsRepo.EXPECT().GetRetentionPeriodUnit().Return("day", nil)
			lsRepo.EXPECT().GetDumpFormat().Return("csv", nil)
			lsRepo.EXPECT().GetDumpTime().Return("00:00", nil)

			strategy, err := logStrategy.GetDumpStrategy(ctx, []string{})
			assert.NoError(t, err)
			assert.Equal(t, 30, strategy[lsconsts.RetentionPeriod])
			assert.Equal(t, "day", strategy[lsconsts.RetentionPeriodUnit])
			assert.Equal(t, "csv", strategy[lsconsts.DumpFormat])
			assert.Equal(t, "00:00", strategy[lsconsts.DumpTime])
		})

		Convey("获取指定字段策略成功", func() {
			ctx := context.Background()
			fields := []string{lsconsts.RetentionPeriod, lsconsts.DumpFormat}

			lsRepo.EXPECT().GetRetentionPeriod().Return(30, nil)
			lsRepo.EXPECT().GetDumpFormat().Return("csv", nil)

			strategy, err := logStrategy.GetDumpStrategy(ctx, fields)
			assert.NoError(t, err)
			assert.Equal(t, 30, strategy[lsconsts.RetentionPeriod])
			assert.Equal(t, "csv", strategy[lsconsts.DumpFormat])
		})
	})
}

func TestSetDumpStrategy(t *testing.T) {
	Convey("SetDumpStrategy", t, func() {
		logger, tracer, lsRepo, logMgnt := newLogDependencies(t)
		logStrategy := newLogStrategy(logger, tracer, lsRepo, logMgnt)

		Convey("设置转存策略成功", func() {
			ctx := context.WithValue(context.Background(), common.VisitorKey, &models.Visitor{
				ID:   "test_user",
				Name: "Test User",
			})

			req := map[string]interface{}{
				lsconsts.RetentionPeriod:     float64(30),
				lsconsts.RetentionPeriodUnit: "day",
				lsconsts.DumpFormat:          "csv",
				lsconsts.DumpTime:            "00:00",
			}

			// Mock事务操作
			mockTx := mock.NewMockDumpStrategyTx(gomock.NewController(t))
			lsRepo.EXPECT().BeginDumpTx().Return(mockTx, nil)
			mockTx.EXPECT().SetRetentionPeriod(30).Return(nil)
			mockTx.EXPECT().SetRetentionPeriodUnit("day").Return(nil)
			mockTx.EXPECT().SetDumpFormat("csv").Return(nil)
			mockTx.EXPECT().SetDumpTime("00:00").Return(nil)
			mockTx.EXPECT().Commit().Return(nil)

			// Mock日志记录
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil)

			err := logStrategy.SetDumpStrategy(ctx, req)
			assert.NoError(t, err)
		})
	})
}
