package logics

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"AuditLog/common"
	"AuditLog/interfaces"
	"AuditLog/interfaces/mock"
	"AuditLog/models"
	"AuditLog/models/rcvo"
	"AuditLog/test/mock_log"
	"AuditLog/test/mock_trace"
)

func newActiveDependencies(t *testing.T) (
	*mock_log.MockLogger,
	*mock_trace.MockTracer,
	*mock.MockLogRepo,
	*mock.MockUserMgntRepo,
	*mock.MockLogScopeStrategyRepo,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return mock_log.NewMockLogger(ctrl),
		mock_trace.NewMockTracer(ctrl),
		mock.NewMockLogRepo(ctrl),
		mock.NewMockUserMgntRepo(ctrl),
		mock.NewMockLogScopeStrategyRepo(ctrl)
}

func newActiveLog(
	logger *mock_log.MockLogger,
	tracer *mock_trace.MockTracer,
	logRepo *mock.MockLogRepo,
	userMgnt *mock.MockUserMgntRepo,
	logScopeStrategy *mock.MockLogScopeStrategyRepo,
) interfaces.ActiveLog {
	return &ActiveLog{
		logger:           logger,
		tracer:           tracer,
		loginLogRepo:     logRepo,
		mgntLogRepo:      logRepo,
		operLogRepo:      logRepo,
		userMgntRepo:     userMgnt,
		logScopeStrategy: logScopeStrategy,
	}
}

func TestGetActiveDataList(t *testing.T) {
	Convey("GetActiveDataList", t, func() {
		logger, tracer, logRepo, userMgnt, logScopeStrategy := newActiveDependencies(t)
		activeLog := newActiveLog(logger, tracer, logRepo, userMgnt, logScopeStrategy)

		Convey("超级管理员获取活跃日志列表成功", func() {
			ctx := context.Background()
			category := common.Login
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

			// Mock活跃日志数据
			logs := []*models.LogPO{
				{
					LogID:    "1211",
					Date:     1717670099640000,
					UserName: "admin",
					Level:    1,
					OpType:   1,
					IP:       "127.0.0.1",
					MAC:      "00:00:00:00:00:00",
					Msg:      "test message",
				},
			}
			logRepo.EXPECT().FindByCondition(req.Offset, req.Limit, gomock.Any(), req.IDs).Return(logs, nil)
			logRepo.EXPECT().FindCountByCondition(gomock.Any()).Return(10, nil)

			// Mock tracer
			tracer.EXPECT().AddInternalTrace(gomock.Any()).MaxTimes(1)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			res, err := activeLog.GetActiveDataList(ctx, category, req, userID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(res.Entries))
			assert.Equal(t, "1211", res.Entries[0].ID)
			assert.Equal(t, 10, res.TotalCount)
		})

		Convey("无权限用户获取活跃日志列表失败", func() {
			ctx := context.Background()
			category := common.Login
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
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			res, err := activeLog.GetActiveDataList(ctx, category, req, userID)
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

func TestGetActiveMetadata(t *testing.T) {
	Convey("GetActiveMetadata", t, func() {
		logger, tracer, logRepo, userMgnt, logScopeStrategy := newActiveDependencies(t)
		activeLog := newActiveLog(logger, tracer, logRepo, userMgnt, logScopeStrategy)

		Convey("获取活跃日志元数据成功", func() {
			meta, err := activeLog.GetActiveMetadata()
			assert.NoError(t, err)
			assert.NotNil(t, meta)
		})
	})
}

func TestGetActiveFieldValues(t *testing.T) {
	Convey("GetActiveFieldValues", t, func() {
		logger, tracer, logRepo, userMgnt, logScopeStrategy := newActiveDependencies(t)
		activeLog := newActiveLog(logger, tracer, logRepo, userMgnt, logScopeStrategy)

		Convey("获取活跃日志报表字段值成功", func() {
			// Mock tracer
			tracer.EXPECT().AddInternalTrace(gomock.Any()).MaxTimes(1)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			req := &rcvo.ReportGetFieldValuesReq{
				Field: "level",
				ReportGetFieldValuesReqBody: rcvo.ReportGetFieldValuesReqBody{
					Limit:     10,
					Offset:    0,
					Condition: map[string]any{},
					KeyWord:   "",
				},
			}
			res, err := activeLog.GetActiveFieldValues(context.Background(), common.Login, req)
			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
	})
}

func TestGetUserIDsByRoleName(t *testing.T) {
	Convey("GetUserIDsByRoleName", t, func() {
		logger, tracer, logRepo, userMgnt, logScopeStrategy := newActiveDependencies(t)
		activeLog := newActiveLog(logger, tracer, logRepo, userMgnt, logScopeStrategy)

		Convey("成功获取角色成员ID", func() {
			roleName := "sec_admin"
			expectedUserIDs := []string{"user1", "user2"}

			// Mock 角色成员信息返回
			roleMemberInfos := []*models.RoleMemberInfo{
				{
					Role: roleName,
					Members: []models.MemberInfo{
						{ID: "user1"},
						{ID: "user2"},
					},
				},
			}
			userMgnt.EXPECT().GetUserIDsByRoleNames([]string{roleName}).Return(roleMemberInfos, 200, nil)

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			userIDs, err := concreteActiveLog.getUserIDsByRoleName(roleName)
			assert.NoError(t, err)
			assert.Equal(t, expectedUserIDs, userIDs)
		})

		Convey("获取角色成员ID失败", func() {
			roleName := "sec_admin"
			userMgnt.EXPECT().GetUserIDsByRoleNames([]string{roleName}).Return(nil, 500, errors.New("mock error"))

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			userIDs, err := concreteActiveLog.getUserIDsByRoleName(roleName)
			assert.Error(t, err)
			assert.Nil(t, userIDs)
		})
	})
}

func TestGetUserIds(t *testing.T) {
	Convey("GetUserIds", t, func() {
		logger, tracer, logRepo, userMgnt, logScopeStrategy := newActiveDependencies(t)
		activeLog := newActiveLog(logger, tracer, logRepo, userMgnt, logScopeStrategy)

		ctx := context.Background()
		logType := common.Login
		userID := "test_user"

		Convey("超级管理员访问", func() {
			// Mock 用户信息返回
			userInfos := []models.User{
				{
					ID:    userID,
					Roles: []string{common.SuperAdmin},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(userInfos, 200, nil)

			// Mock tracer
			tracer.EXPECT().AddInternalTrace(gomock.Any()).MaxTimes(1)
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			tracer.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).MaxTimes(1)

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			includeIDs, excludeIDs, err := concreteActiveLog.getUserIds(ctx, logType, userID)
			assert.NoError(t, err)
			assert.Nil(t, includeIDs)
			assert.Nil(t, excludeIDs)
		})

		Convey("三权分立角色访问-系统管理员", func() {
			// Mock 用户信息返回
			userInfos := []models.User{
				{
					ID:    userID,
					Roles: []string{common.SysAdmin},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(userInfos, 200, nil)

			// Mock 获取角色可访问范围
			scope := []string{common.NormalUser}
			logScopeStrategy.EXPECT().GetActiveScopeBy(common.LogTypeMap[logType], common.SysAdmin).Return(scope, nil)

			// Mock 获取需要排除的用户ID
			for _, role := range common.MutuallyRoles {
				if !common.InArray(role, scope) {
					excludeUserIDs := []*models.RoleMemberInfo{
						{
							Role: role,
							Members: []models.MemberInfo{
								{ID: "exclude_user1"},
								{ID: "exclude_user2"},
							},
						},
					}
					userMgnt.EXPECT().GetUserIDsByRoleNames([]string{role}).Return(excludeUserIDs, 200, nil)
				}
			}

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			includeIDs, excludeIDs, err := concreteActiveLog.getUserIds(ctx, logType, userID)
			assert.NoError(t, err)
			assert.Nil(t, includeIDs)
			assert.NotNil(t, excludeIDs)
		})

		Convey("无权限用户访问", func() {
			// Mock 用户信息返回
			userInfos := []models.User{
				{
					ID:    userID,
					Roles: []string{"normal_user"},
				},
			}
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(userInfos, 200, nil)

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			includeIDs, excludeIDs, err := concreteActiveLog.getUserIds(ctx, logType, userID)
			assert.Error(t, err)
			assert.Nil(t, includeIDs)
			assert.Nil(t, excludeIDs)
		})

		Convey("获取用户信息失败", func() {
			userMgnt.EXPECT().GetUserInfoByID([]string{userID}).Return(nil, 500, errors.New("mock error"))
			logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

			// 将 activeLog 转换为具体类型 *ActiveLog
			concreteActiveLog := activeLog.(*ActiveLog)
			includeIDs, excludeIDs, err := concreteActiveLog.getUserIds(ctx, logType, userID)
			assert.Error(t, err)
			assert.Nil(t, includeIDs)
			assert.Nil(t, excludeIDs)
		})
	})
}
