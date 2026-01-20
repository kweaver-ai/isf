package db

import (
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/models"
)

var (
	hOnce sync.Once
	hl    *historyLog
)

type historyLog struct {
	db     *sqlx.DB
	logger api.Logger
}

// NewHisotryLog 创建历史日志数据库对象
func NewHisotryLog() interfaces.HistoryRepo {
	hOnce.Do(func() {
		hl = &historyLog{
			db:     drivenadapters.DBPool,
			logger: drivenadapters.Logger,
		}
	})
	return hl
}

// New 创建历史日志
func (repo *historyLog) New(log *models.HistoryPO) (err error) {
	sqlStr := "INSERT INTO " + infra.GetDBName() + ".t_history_log_info" +
		" (f_id, f_name, f_size, f_type, f_date, f_dump_date, f_oss_id) " +
		" VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = repo.db.Exec(sqlStr, log.ID, log.Name, log.Size, log.Type, log.Date, log.DumpDate, log.OssID)
	if err != nil {
		repo.logger.Errorf("insert history log error: %v", err)
		return
	}
	return
}

// FindCountByCondition 查询历史日志数量
func (repo *historyLog) FindCountByCondition(condition string) (count int, err error) {
	sqlStr := "SELECT COUNT(f_id) FROM " + infra.GetDBName() + ".t_history_log_info " + condition
	err = repo.db.QueryRow(sqlStr).Scan(&count)
	if err != nil {
		repo.logger.Errorf("db query history log error: %v", err)
		return
	}
	return
}

// FindByCondition 查询历史日志
func (repo *historyLog) FindByCondition(offset, limit int, condition string, ids []string) (logs []*models.HistoryPO, err error) {
	sqlStr := `SELECT
		f_id,
		f_name,
		f_size,
		f_type,
		f_date,
		f_dump_date,
		f_oss_id
		FROM ` + infra.GetDBName() + `.t_history_log_info ` + condition + ` LIMIT ? OFFSET ?
	`

	rows, err := repo.db.Query(sqlStr, limit, offset)
	if err != nil {
		repo.logger.Errorf("db query login log error: %v", err)
		return
	}
	defer rows.Close()

	if len(ids) > 0 {
		resMap := make(map[string]*models.HistoryPO)
		for rows.Next() {
			log := models.HistoryPO{}
			err = rows.Scan(
				&log.ID,
				&log.Name,
				&log.Size,
				&log.Type,
				&log.Date,
				&log.DumpDate,
				&log.OssID,
			)
			if err != nil {
				repo.logger.Errorf("db scan login log error: %v", err)
				return
			}
			resMap[log.ID] = &log
		}

		for _, id := range ids {
			log, ok := resMap[id]
			if ok {
				logs = append(logs, log)
			}
		}
	} else {
		for rows.Next() {
			log := models.HistoryPO{}
			err = rows.Scan(
				&log.ID,
				&log.Name,
				&log.Size,
				&log.Type,
				&log.Date,
				&log.DumpDate,
				&log.OssID,
			)
			if err != nil {
				repo.logger.Errorf("db scan login log error: %v", err)
				return
			}
			logs = append(logs, &log)
		}
	}
	return
}

// GetHistoryLogsByType 获取历史日志
func (repo *historyLog) GetHistoryLogsByType(logType int8) (logs []*models.HistoryPO, err error) {
	sqlStr := `SELECT
		f_id,
		f_name,
		f_size,
		f_type,
		f_date,
		f_dump_date,
		f_oss_id
		FROM ` + infra.GetDBName() + `.t_history_log_info WHERE f_type = ?`

	rows, err := repo.db.Query(sqlStr, logType)
	if err != nil {
		repo.logger.Errorf("db query history log error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		log := models.HistoryPO{}
		err = rows.Scan(
			&log.ID,
			&log.Name,
			&log.Size,
			&log.Type,
			&log.Date,
			&log.DumpDate,
			&log.OssID,
		)
		if err != nil {
			repo.logger.Errorf("db scan history log error: %v", err)
			return
		}
		logs = append(logs, &log)
	}
	return
}

// GetHistoryLogByID 获取历史日志
func (repo *historyLog) GetHistoryLogByID(id string) (log *models.HistoryPO, err error) {
	sqlStr := `SELECT
		f_id,
		f_name,
		f_size,
		f_type,
		f_date,
		f_dump_date,
		f_oss_id
		FROM ` + infra.GetDBName() + `.t_history_log_info WHERE f_id = ?`

	var logInfo models.HistoryPO
	err = repo.db.QueryRow(sqlStr, id).Scan(
		&logInfo.ID,
		&logInfo.Name,
		&logInfo.Size,
		&logInfo.Type,
		&logInfo.Date,
		&logInfo.DumpDate,
		&logInfo.OssID,
	)
	if err != nil {
		repo.logger.Errorf("db get history log error: %v", err)
		return
	}
	log = &logInfo
	return
}
