package db

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/common/constants/lsconsts"
	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/models/lsmodels"
)

var (
	ssOnce sync.Once
	ss     *scopeStrategy
)

type scopeStrategy struct {
	db     *sqlx.DB
	logger api.Logger
}

func NewScopeStrategy() interfaces.LogScopeStrategyRepo {
	ssOnce.Do(func() {
		ss = &scopeStrategy{
			db:     drivenadapters.DBPool,
			logger: drivenadapters.Logger,
		}
	})
	return ss
}

// GetStrategiesByCondition 根据条件查询日志范围策略
func (s *scopeStrategy) GetStrategiesByCondition(condition string, params []interface{}) (res []*lsmodels.ScopeStrategyPO, err error) {
	sqlStr := "SELECT f_id, f_log_type, f_log_category, f_role, f_scope FROM " + infra.GetDBName() + ".t_log_scope_strategy " + condition
	rows, err := s.db.Query(sqlStr, params...)
	if err != nil {
		s.logger.Errorf("db query scope strategy error: %v", err)
		return
	}
	defer rows.Close()

	res = make([]*lsmodels.ScopeStrategyPO, 0)
	for rows.Next() {
		scopeStrategy := lsmodels.ScopeStrategyPO{}
		err = rows.Scan(
			&scopeStrategy.ID,
			&scopeStrategy.LogType,
			&scopeStrategy.LogCategory,
			&scopeStrategy.Role,
			&scopeStrategy.Scope,
		)
		if err != nil {
			s.logger.Errorf("db scan scope strategy error: %v", err)
			return
		}
		res = append(res, &scopeStrategy)
	}

	return
}

// NewStrategy 新增日志范围策略
func (s *scopeStrategy) NewStrategy(req *lsmodels.ScopeStrategyPO) (err error) {
	sqlStr := "INSERT INTO " + infra.GetDBName() +
		`.t_log_scope_strategy (
		f_id,
		f_log_type,
		f_log_category,
		f_role,
		f_scope,
		f_created_at,
		f_created_by
		) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = s.db.Exec(sqlStr, req.ID, req.LogType, req.LogCategory, req.Role, req.Scope, req.CreatedAt, req.CreatedBy)
	if err != nil {
		s.logger.Errorf("db insert scope strategy error: %v", err)
		return
	}
	return
}

// UpdateStrategy 更新日志范围策略
func (s *scopeStrategy) UpdateStrategy(req *lsmodels.ScopeStrategyPO) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() +
		`.t_log_scope_strategy SET 
		f_log_type = ?,
		f_log_category = ?,
		f_role = ?,
		f_scope = ?,
		f_updated_at = ?,
		f_updated_by = ?
		WHERE f_id = ?`
	_, err = s.db.Exec(sqlStr, req.LogType, req.LogCategory, req.Role, req.Scope, req.UpdatedAt, req.UpdatedBy, req.ID)
	if err != nil {
		s.logger.Errorf("db update scope strategy error: %v", err)
		return
	}
	return
}

// DeleteStrategy 删除日志范围策略
func (s *scopeStrategy) DeleteStrategy(id int64) (err error) {
	sqlStr := "DELETE FROM " + infra.GetDBName() + ".t_log_scope_strategy WHERE f_id = ?"
	_, err = s.db.Exec(sqlStr, id)
	if err != nil {
		s.logger.Errorf("db delete scope strategy error: %v", err)
		return
	}
	return
}

// GetActiveScopeBy 根据条件获取日志查看范围
func (s *scopeStrategy) GetActiveScopeBy(logType int, role string) (scope []string, err error) {
	sqlStr := "SELECT f_scope FROM " + infra.GetDBName() + ".t_log_scope_strategy WHERE f_log_category = ? AND f_log_type = ? AND f_role = ?"
	var nullableValue sql.NullString
	err = s.db.QueryRow(sqlStr, lsconsts.ActiveLog, logType, role).Scan(&nullableValue)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查询不到数据时返回空切片
			return []string{}, nil
		}
		s.logger.Errorf("db query scope strategy error: %v", err)
		return
	}
	if nullableValue.Valid {
		scope = strings.Split(nullableValue.String, ",")
	}
	return
}

// GetHistoryScopeCountBy 根据条件获取历史日志查看范围
func (s *scopeStrategy) GetHistoryScopeCountBy(logType int, role string) (count int, err error) {
	sqlStr := "SELECT COUNT(f_id) FROM " + infra.GetDBName() + ".t_log_scope_strategy WHERE f_log_category = ? AND f_log_type = ? AND f_role = ?"
	var nullableValue sql.NullInt64
	err = s.db.QueryRow(sqlStr, lsconsts.HistoryLog, logType, role).Scan(&nullableValue)
	if err != nil {
		s.logger.Errorf("db query scope strategy error: %v", err)
		return
	}
	if nullableValue.Valid {
		count = int(nullableValue.Int64)
	}
	return
}

// GetStrategyByID 根据ID获取日志范围策略
func (s *scopeStrategy) GetStrategyByID(id int64) (res *lsmodels.ScopeStrategyPO, err error) {
	sqlStr := "SELECT f_id, f_log_type, f_log_category, f_role, f_scope FROM " + infra.GetDBName() + ".t_log_scope_strategy WHERE f_id = ?"
	res = &lsmodels.ScopeStrategyPO{}
	err = s.db.QueryRow(sqlStr, id).Scan(&res.ID, &res.LogType, &res.LogCategory, &res.Role, &res.Scope)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		s.logger.Errorf("db query scope strategy error: %v", err)
		return
	}
	return
}
