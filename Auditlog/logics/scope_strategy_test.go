package logics

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/models/lsmodels"
	"AuditLog/test/mock_log"
	"AuditLog/test/mock_trace"
)

func newScopeDependencies(t *testing.T) (
	*mock_log.MockLogger,
	*mock_trace.MockTracer,
	*mock.MockLogScopeStrategyRepo,
	*mock.MockLogMgnt,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl),
		mock_trace.NewMockTracer(ctrl),
		mock.NewMockLogScopeStrategyRepo(ctrl),
		mock.NewMockLogMgnt(ctrl)
}

func newScopeStrategy(
	logger *mock_log.MockLogger,
	tracer *mock_trace.MockTracer,
	ssRepo *mock.MockLogScopeStrategyRepo,
	logMgnt *mock.MockLogMgnt,
) *ScopeStrategy {
	return &ScopeStrategy{
		logger:  logger,
		tracer:  tracer,
		ssRepo:  ssRepo,
		logMgnt: logMgnt,
	}
}

func TestGetStrategy(t *testing.T) {
	Convey("GetStrategy", t, func() {
		logger, tracer, ssRepo, logMgnt := newScopeDependencies(t)
		scopeStrategy := newScopeStrategy(logger, tracer, ssRepo, logMgnt)

		Convey("获取策略列表成功", func() {
			ctx := context.Background()
			req := &lsmodels.GetScopeStrategyReq{
				Category: 1,
				Type:     1,
				Role:     "admin",
			}

			// Mock策略数据返回
			strategies := []*lsmodels.ScopeStrategyPO{
				{
					ID:          1,
					LogType:     1,
					LogCategory: 1,
					Role:        "admin",
					Scope:       "scope1,scope2",
				},
			}
			ssRepo.EXPECT().GetStrategiesByCondition(gomock.Any(), gomock.Any()).Return(strategies, nil)

			res, err := scopeStrategy.GetStrategy(ctx, req)
			assert.NoError(t, err)
			assert.Equal(t, int64(1), res.TotalCount)
			assert.Equal(t, 1, len(res.Entries))
			assert.Equal(t, 1, res.Entries[0].ID)
		})
	})
}

func TestNewStrategy(t *testing.T) {
	Convey("NewStrategy", t, func() {
		logger, tracer, ssRepo, logMgnt := newScopeDependencies(t)
		scopeStrategy := newScopeStrategy(logger, tracer, ssRepo, logMgnt)

		Convey("创建新策略成功", func() {
			ctx := context.WithValue(context.Background(), common.VisitorKey, &models.Visitor{
				ID:   "test_user",
				Name: "Test User",
			})
			req := &lsmodels.ScopeStrategyVO{
				LogType:     1,
				LogCategory: 1,
				Role:        "admin",
				Scope:       []string{"scope1", "scope2"},
			}

			// Mock检查已存在策略
			ssRepo.EXPECT().GetStrategiesByCondition(gomock.Any(), gomock.Any()).Return([]*lsmodels.ScopeStrategyPO{}, nil)
			// Mock创建新策略
			ssRepo.EXPECT().NewStrategy(gomock.Any()).Return(nil)
			// Mock日志记录
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil)

			id, err := scopeStrategy.NewStrategy(ctx, req)
			assert.NoError(t, err)
			assert.NotZero(t, id)
		})
	})
}

func TestUpdateStrategy(t *testing.T) {
	Convey("UpdateStrategy", t, func() {
		logger, tracer, ssRepo, logMgnt := newScopeDependencies(t)
		scopeStrategy := newScopeStrategy(logger, tracer, ssRepo, logMgnt)

		Convey("更新策略成功", func() {
			ctx := context.WithValue(context.Background(), common.VisitorKey, &models.Visitor{
				ID:   "test_user",
				Name: "Test User",
			})
			req := &lsmodels.ScopeStrategyVO{
				LogType:     1,
				LogCategory: 1,
				Role:        "admin",
				Scope:       []string{"scope1", "scope2"},
			}

			// Mock获取策略
			strategy := &lsmodels.ScopeStrategyPO{
				ID:          1,
				LogType:     1,
				LogCategory: 1,
				Role:        "admin",
				Scope:       "scope1,scope2",
			}
			ssRepo.EXPECT().GetStrategyByID(gomock.Any()).Return(strategy, nil)
			// Mock检查已存在策略
			ssRepo.EXPECT().GetStrategiesByCondition(gomock.Any(), gomock.Any()).Return([]*lsmodels.ScopeStrategyPO{}, nil)
			// Mock更新策略
			ssRepo.EXPECT().UpdateStrategy(gomock.Any()).Return(nil)
			// Mock日志记录
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil)

			err := scopeStrategy.UpdateStrategy(ctx, 1, req)
			assert.NoError(t, err)
		})
	})
}

func TestDeleteStrategy(t *testing.T) {
	Convey("DeleteStrategy", t, func() {
		logger, tracer, ssRepo, logMgnt := newScopeDependencies(t)
		scopeStrategy := newScopeStrategy(logger, tracer, ssRepo, logMgnt)

		Convey("删除策略成功", func() {
			ctx := context.WithValue(context.Background(), common.VisitorKey, &models.Visitor{
				ID:   "test_user",
				Name: "Test User",
			})

			// Mock获取策略
			strategy := &lsmodels.ScopeStrategyPO{
				ID:          1,
				LogType:     1,
				LogCategory: 1,
				Role:        "admin",
				Scope:       "scope1,scope2",
			}
			ssRepo.EXPECT().GetStrategyByID(gomock.Any()).Return(strategy, nil)
			// Mock删除策略
			ssRepo.EXPECT().DeleteStrategy(gomock.Any()).Return(nil)
			// Mock日志记录
			logMgnt.EXPECT().SendLog(gomock.Any()).Return(nil)

			err := scopeStrategy.DeleteStrategy(ctx, 1)
			assert.NoError(t, err)
		})
	})
}
