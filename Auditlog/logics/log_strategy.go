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
	"AuditLog/common/constants/lsconsts"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/locale"
	"AuditLog/models"
)

var (
	lsOnce sync.Once
	ls     *LogStrategy
)

type LogStrategy struct {
	logger          api.Logger
	tracer          api.Tracer
	logStrategyRepo interfaces.LogStrategyRepo
	logMgnt         interfaces.LogMgnt
}

func NewLogStrategy() interfaces.LogStrategy {
	lsOnce.Do(func() {
		ls = &LogStrategy{
			logger:          logger,
			tracer:          tracer,
			logStrategyRepo: logStrategyRepo,
			logMgnt:         NewLogMgnt(),
		}
	})
	return ls
}

// GetDumpStrategy 获取日志转存策略
func (l *LogStrategy) GetDumpStrategy(ctx context.Context, fields []string) (strategy map[string]interface{}, err error) {
	// 创建结果对象
	strategy = make(map[string]interface{})
	allFields := fields
	if len(allFields) == 0 {
		allFields = lsconsts.AllDumpFields
	}

	for _, field := range allFields {
		switch field {
		case lsconsts.RetentionPeriod:
			period, err := l.logStrategyRepo.GetRetentionPeriod()
			if err != nil {
				l.logger.Errorf("[GetDumpStrategy] failed to get retention_period: %v", err)
				return nil, err
			}
			strategy[lsconsts.RetentionPeriod] = period

		case lsconsts.RetentionPeriodUnit:
			unit, err := l.logStrategyRepo.GetRetentionPeriodUnit()
			if err != nil {
				l.logger.Errorf("[GetDumpStrategy] failed to get retention_period_unit: %v", err)
				return nil, err
			}
			strategy[lsconsts.RetentionPeriodUnit] = unit

		case lsconsts.DumpFormat:
			format, err := l.logStrategyRepo.GetDumpFormat()
			if err != nil {
				l.logger.Errorf("[GetDumpStrategy] failed to get dump_format: %v", err)
				return nil, err
			}
			strategy[lsconsts.DumpFormat] = format

		case lsconsts.DumpTime:
			time, err := l.logStrategyRepo.GetDumpTime()
			if err != nil {
				l.logger.Errorf("[GetDumpStrategy] failed to get dump_time: %v", err)
				return nil, err
			}
			strategy[lsconsts.DumpTime] = time
		}
	}

	return strategy, nil
}

// SetDumpStrategy 设置日志转存策略
func (l *LogStrategy) SetDumpStrategy(ctx context.Context, req map[string]interface{}) (err error) {
	tx, err := l.logStrategyRepo.BeginDumpTx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				l.logger.Warnf("[SetDumpStrategy] rollback error: %v", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				l.logger.Warnf("[SetDumpStrategy] commit error: %v", err)
			}
		}
	}()

	exmsg := []string{}
	for field, value := range req {
		switch field {
		case lsconsts.RetentionPeriod:
			if err = tx.SetRetentionPeriod(int(value.(float64))); err != nil {
				return err
			}
			exmsg = append(exmsg, fmt.Sprintf(locale.GetI18nCtx(ctx, locale.LogDumpPeriod)+": %d%s", int(value.(float64)), locale.GetI18nCtx(ctx, req[lsconsts.RetentionPeriodUnit].(string))))
		case lsconsts.RetentionPeriodUnit:
			if err = tx.SetRetentionPeriodUnit(value.(string)); err != nil {
				return err
			}
		case lsconsts.DumpFormat:
			if err = tx.SetDumpFormat(value.(string)); err != nil {
				return err
			}
			exmsg = append(exmsg, fmt.Sprintf(locale.GetI18nCtx(ctx, locale.LogDumpFormat)+": %s", value))
		case lsconsts.DumpTime:
			if err = tx.SetDumpTime(value.(string)); err != nil {
				return err
			}
			exmsg = append(exmsg, fmt.Sprintf(locale.GetI18nCtx(ctx, locale.LogDumpTime)+": %s", value))
		}
	}

	go func() {
		visitor := ctx.Value(common.VisitorKey).(*models.Visitor)
		if visitor == nil {
			return
		}
		err := l.logMgnt.SendLog(&models.SendLogVo{
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
				Msg:       locale.GetI18nCtx(ctx, locale.SetLogDumpStrategy),
				Exmsg:     strings.Join(exmsg, ", "),
				UserAgent: visitor.AgentType,
				OutBizID:  uuid.NewString(),
			},
		})
		if err != nil {
			l.logger.Warnf("[SetDumpStrategy] send log error: %v", err)
		}
	}()

	return nil
}
