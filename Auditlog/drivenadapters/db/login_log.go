package db

import (
	"database/sql"
	"strconv"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
	"AuditLog/models"
)

var (
	lOnce sync.Once
	ll    *loginLog
)

type loginLog struct {
	db     *sqlx.DB
	logger api.Logger
}

// NewLoginLog 创建登录审计日志数据库对象
func NewLoginLog() interfaces.LogRepo {
	lOnce.Do(func() {
		ll = &loginLog{
			db:     drivenadapters.DBPool,
			logger: drivenadapters.Logger,
		}
	})
	return ll
}

func (repo *loginLog) New(log *models.AuditLog) (logID string, err error) {
	// 数据库操作
	uid, err := infra.GetUniqueID()
	if err != nil {
		repo.logger.Errorf("new sonyflake id error: %v", err)
		return
	}
	sqlStr := "INSERT INTO " + infra.GetDBName() + ".t_log_login" +
		" (f_log_id, f_user_id,f_user_name,f_user_type,f_obj_id,f_level,f_op_type,f_date,f_ip,f_mac,f_msg,f_exmsg,f_user_agent,f_additional_info,f_user_paths,f_obj_name,f_obj_type) " +
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		// 需要考虑事务一致性的问题
	_, err = repo.db.Exec(sqlStr, uid, log.UserID, log.UserName, log.UserType, log.ObjID, log.Level, log.OpType, log.Date, log.IP, log.Mac, log.Msg, log.Exmsg, log.UserAgent, log.AdditionalInfo, log.DeptPaths, log.ObjName, log.ObjType)
	if err != nil {
		repo.logger.Errorf("insert log error: %v, business key: %v", err, log.OutBizID)
		return
	}

	return strconv.FormatUint(uid, 10), nil
}

// FindCountByCondition 根据条件查询登录审计日志数量
func (repo *loginLog) FindCountByCondition(condition string) (count int, err error) {
	sqlStr := "SELECT COUNT(f_log_id) FROM " + infra.GetDBName() + ".t_log_login " + condition
	err = repo.db.QueryRow(sqlStr).Scan(&count)
	if err != nil {
		repo.logger.Errorf("db login log [FindCountByCondition] error: %v", err)
		return
	}
	return
}

// FindByCondition 根据条件查询登录审计日志
func (repo *loginLog) FindByCondition(offset, limit int, condition string, ids []string) (logs []*models.LogPO, err error) {
	sqlStr := `SELECT
		f_log_id,
		f_user_id,
		f_user_name,
		f_obj_id,
		f_level,
		f_op_type,
		f_date,
		f_ip,
		f_mac,
		f_msg,
		f_exmsg,
		f_user_agent,
		f_additional_info,
		f_user_paths,
		f_obj_name,
		f_obj_type
		FROM ` + infra.GetDBName() + `.t_log_login ` + condition + ` LIMIT ? OFFSET ?
	`

	rows, err := repo.db.Query(sqlStr, limit, offset)
	if err != nil {
		repo.logger.Errorf("db query login log error: %v", err)
		return
	}
	defer rows.Close()

	if len(ids) > 0 {
		resMap := make(map[string]*models.LogPO)
		for rows.Next() {
			log := models.LogPO{}
			err = rows.Scan(
				&log.LogID,
				&log.UserID,
				&log.UserName,
				&log.ObjID,
				&log.Level,
				&log.OpType,
				&log.Date,
				&log.IP,
				&log.MAC,
				&log.Msg,
				&log.ExMsg,
				&log.UserAgent,
				&log.AdditionalInfo,
				&log.UserPaths,
				&log.ObjName,
				&log.ObjType,
			)
			if err != nil {
				repo.logger.Errorf("db scan login log error: %v", err)
				return
			}
			resMap[log.LogID] = &log
		}

		for _, id := range ids {
			log, ok := resMap[id]
			if ok {
				logs = append(logs, log)
			}
		}
	} else {
		for rows.Next() {
			log := models.LogPO{}
			err = rows.Scan(
				&log.LogID,
				&log.UserID,
				&log.UserName,
				&log.ObjID,
				&log.Level,
				&log.OpType,
				&log.Date,
				&log.IP,
				&log.MAC,
				&log.Msg,
				&log.ExMsg,
				&log.UserAgent,
				&log.AdditionalInfo,
				&log.UserPaths,
				&log.ObjName,
				&log.ObjType,
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

// GetFirstLogTime 获取第一个登录审计日志的时间 单位：微秒
func (repo *loginLog) GetFirstLogTime() (timeMicro int64, err error) {
	var nullableTime sql.NullInt64
	sqlStr := "SELECT MIN(f_date) FROM " + infra.GetDBName() + ".t_log_login"
	err = repo.db.QueryRow(sqlStr).Scan(&nullableTime)
	if err != nil {
		repo.logger.Errorf("db login log [GetFirstLogTime] error: %v", err)
		return
	}

	if nullableTime.Valid {
		timeMicro = nullableTime.Int64
	} else {
		timeMicro = -1
	}
	return
}

// ClearOutdatedLog 清除过期日志
func (repo *loginLog) ClearOutdatedLog(logID, date, batchSize, sleepTime int64) (err error) {
	if batchSize <= 0 {
		batchSize = 50000
	}

	sqlStr := "DELETE FROM " + infra.GetDBName() + ".t_log_login WHERE f_log_id <= ? AND f_date < ? LIMIT ?"
	for {
		stat, err := repo.db.Exec(sqlStr, logID, date, batchSize)
		if err != nil {
			repo.logger.Errorf("db login log [ClearOutdatedLog] error: %v", err)
			return err
		}

		affected, err := stat.RowsAffected()
		if err != nil {
			repo.logger.Errorf("db login log [ClearOutdatedLog] error: %v", err)
			return err
		}

		if affected < batchSize {
			break
		}

		if sleepTime > 0 {
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	return
}

// GetLogCount 获取登录审计日志数量
func (repo *loginLog) GetLogCount() (count int64, err error) {
	sqlStr := "SELECT COUNT(f_log_id) FROM " + infra.GetDBName() + ".t_log_login"
	err = repo.db.QueryRow(sqlStr).Scan(&count)
	if err != nil {
		repo.logger.Errorf("db login log [GetLogCount] error: %v", err)
		return
	}
	return
}
