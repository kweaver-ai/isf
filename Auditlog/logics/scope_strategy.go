package logics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"AuditLog/common"
	"AuditLog/common/constants/logconsts"
	"AuditLog/errors"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
	"AuditLog/models/lsmodels"
)

var (
	ssOnce sync.Once
	ss     *ScopeStrategy
)

type ScopeStrategy struct {
	logger  api.Logger
	tracer  api.Tracer
	ssRepo  interfaces.LogScopeStrategyRepo
	logMgnt interfaces.LogMgnt
}

func NewScopeStrategy() interfaces.LogScopeStrategy {
	ssOnce.Do(func() {
		ss = &ScopeStrategy{
			logger:  logger,
			tracer:  tracer,
			ssRepo:  logScopeStrategyRepo,
			logMgnt: NewLogMgnt(),
		}
	})
	return ss
}

func (s *ScopeStrategy) GetStrategy(ctx context.Context, req *lsmodels.GetScopeStrategyReq) (res *lsmodels.GetScopeStrategyRes, err error) {
	res = &lsmodels.GetScopeStrategyRes{}

	var conditions []string
	params := []interface{}{}
	if req.Category != 0 {
		conditions = append(conditions, "f_log_category=?")
		params = append(params, req.Category)
	}
	if req.Type != 0 {
		conditions = append(conditions, "f_log_type=?")
		params = append(params, req.Type)
	}
	if req.Role != "" {
		conditions = append(conditions, "f_role=?")
		params = append(params, req.Role)
	}

	var condition string
	if len(conditions) > 0 {
		condition = "WHERE " + strings.Join(conditions, " AND ")
	}

	if req.Limit > 0 {
		condition += " LIMIT ? OFFSET ?"
		params = append(params, req.Limit, req.Offset)
	}

	strategies, err := s.ssRepo.GetStrategiesByCondition(condition, params)
	if err != nil {
		return nil, fmt.Errorf("[GetScopeStrategy] get strategies failed: %w", err)
	}

	entries := make([]*lsmodels.ScopeStrategyVO, 0)
	for _, strategy := range strategies {
		scope := []string{}
		if strategy.Scope != "" {
			scope = strings.Split(strategy.Scope, ",")
		}
		entries = append(entries, &lsmodels.ScopeStrategyVO{
			ID:          int(strategy.ID),
			LogType:     strategy.LogType,
			LogCategory: strategy.LogCategory,
			Role:        strategy.Role,
			Scope:       scope,
		})
	}
	res.Entries = entries
	res.TotalCount = int64(len(strategies))
	return
}

func (s *ScopeStrategy) NewStrategy(ctx context.Context, req *lsmodels.ScopeStrategyVO) (id int64, err error) {
	visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
	existing, err := s.ssRepo.GetStrategiesByCondition(
		"WHERE f_log_type=? AND f_log_category=? AND f_role=?",
		[]interface{}{req.LogType, req.LogCategory, req.Role},
	)
	if err != nil {
		return 0, err
	}
	if len(existing) > 0 {
		return 0, errors.NewCtx(ctx, errors.ScopeStrategyConflictErr, "Policy already exists", nil)
	}

	uid, err := infra.GetUniqueID()
	if err != nil {
		s.logger.Errorf("new sonyflake id error: %v", err)
		return 0, err
	}
	strategy := &lsmodels.ScopeStrategyPO{
		ID:          int64(uid),
		LogType:     req.LogType,
		LogCategory: req.LogCategory,
		Role:        req.Role,
		Scope:       strings.Join(req.Scope, ","),
		CreatedBy:   visitor.ID,
		CreatedAt:   time.Now().UnixMicro(),
	}
	err = s.ssRepo.NewStrategy(strategy)
	if err != nil {
		return 0, err
	}

	go s.autilog(
		ctx,
		req,
		logconsts.LogLevel.INFO,
		logconsts.OpType.ManagementType.CREATE,
		locale.NewLogScopeStrategy,
	)

	return int64(uid), nil
}

func (s *ScopeStrategy) UpdateStrategy(ctx context.Context, id int64, req *lsmodels.ScopeStrategyVO) (err error) {
	checked, err := s.ssRepo.GetStrategyByID(id)
	if err != nil {
		return err
	}
	if checked == nil {
		return errors.NewCtx(
			ctx,
			errors.ScopeStrategyNotFoundErr,
			"Policy not found",
			map[string]interface{}{
				"id": []int64{id},
			},
		)
	}

	existing, err := s.ssRepo.GetStrategiesByCondition(
		"WHERE f_log_type=? AND f_log_category=? AND f_role=?",
		[]interface{}{req.LogType, req.LogCategory, req.Role},
	)
	if err != nil {
		return err
	}
	if len(existing) > 0 && existing[0].ID != id {
		return errors.NewCtx(ctx, errors.ScopeStrategyConflictErr, "Policy already exists", nil)
	}

	visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
	strategy := &lsmodels.ScopeStrategyPO{
		ID:          id,
		LogType:     req.LogType,
		LogCategory: req.LogCategory,
		Role:        req.Role,
		Scope:       strings.Join(req.Scope, ","),
		UpdatedBy:   visitor.ID,
		UpdatedAt:   time.Now().UnixMicro(),
	}
	err = s.ssRepo.UpdateStrategy(strategy)
	if err != nil {
		return err
	}

	go s.autilog(
		ctx,
		req,
		logconsts.LogLevel.INFO,
		logconsts.OpType.ManagementType.EDIT,
		locale.EditLogScopeStrategy,
	)

	return
}

func (s *ScopeStrategy) DeleteStrategy(ctx context.Context, id int64) (err error) {
	strategy, err := s.ssRepo.GetStrategyByID(id)
	if err != nil {
		return err
	}
	if strategy == nil {
		return
	}
	err = s.ssRepo.DeleteStrategy(id)
	if err != nil {
		return err
	}

	strategyVO := &lsmodels.ScopeStrategyVO{
		ID:          int(strategy.ID),
		LogType:     strategy.LogType,
		LogCategory: strategy.LogCategory,
		Role:        strategy.Role,
		Scope:       strings.Split(strategy.Scope, ","),
	}

	go s.autilog(
		ctx,
		strategyVO,
		logconsts.LogLevel.WARN,
		logconsts.OpType.ManagementType.DELETE,
		locale.DeleteLogScopeStrategy,
	)

	return
}

// 记录审计日志
func (s *ScopeStrategy) autilog(ctx context.Context, strategy *lsmodels.ScopeStrategyVO, level int, opType int, opKey string) {
	visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
	scope := []string{}
	for _, s := range strategy.Scope {
		scope = append(scope, locale.GetI18nCtx(ctx, s))
	}
	exmsg := ""
	if len(scope) > 0 {
		exmsg = "; " + fmt.Sprintf(
			locale.GetI18nCtx(ctx, locale.LogScope)+": %s",
			strings.Join(scope, ","),
		)
	}
	err := s.logMgnt.SendLog(&models.SendLogVo{
		LogType:  common.Management,
		Language: "",
		LogContent: &models.AuditLog{
			UserID:   visitor.ID,
			UserName: visitor.Name,
			UserType: common.AuthenticatedUser,
			Level:    level,
			OpType:   opType,
			Date:     time.Now().UnixMicro(),
			IP:       visitor.IP,
			Mac:      visitor.Mac,
			Msg: fmt.Sprintf(
				locale.GetI18nCtx(ctx, opKey),
				locale.GetI18nCtx(ctx, locale.LogCategoryMap[int(strategy.LogCategory)]),
			),
			Exmsg: fmt.Sprintf(
				locale.GetI18nCtx(ctx, locale.LogType)+": %s; "+
					locale.GetI18nCtx(ctx, locale.LogRole)+": %s",
				locale.GetI18nCtx(ctx, locale.LogTypeMap[int(strategy.LogType)]),
				locale.GetI18nCtx(ctx, strategy.Role),
			) + exmsg,
			UserAgent: visitor.AgentType,
			OutBizID:  uuid.NewString(),
		},
	})
	if err != nil {
		s.logger.Warnf("[LogScopeStrategy] send log error: %v", err)
	}
}
