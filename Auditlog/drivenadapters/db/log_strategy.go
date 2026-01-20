package db

import (
	"database/sql"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	"AuditLog/interfaces"
)

var (
	lsOnce sync.Once
	ls     *logStrategy
)

type logStrategy struct {
	db     *sqlx.DB
	logger api.Logger
}

func NewLogStrategy() interfaces.LogStrategyRepo {
	lsOnce.Do(func() {
		ls = &logStrategy{
			db:     drivenadapters.DBPool,
			logger: drivenadapters.Logger,
		}
	})
	return ls
}

// GetRetentionPeriod 获取日志保留周期
func (l *logStrategy) GetRetentionPeriod() (period int, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'retention_period'"
	var nullableValue sql.NullInt64
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query retention_period error: %v", err)
		return
	}
	if nullableValue.Valid {
		period = int(nullableValue.Int64)
	}
	return
}

// GetRetentionPeriodUnit 获取日志保留周期单位
func (l *logStrategy) GetRetentionPeriodUnit() (unit string, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'retention_period_unit'"
	var nullableValue sql.NullString
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query retention_period_unit error: %v", err)
		return
	}
	if nullableValue.Valid {
		unit = nullableValue.String
	}
	return
}

// GetDumpFormat 获取日志保留格式
func (l *logStrategy) GetDumpFormat() (format string, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'dump_format'"
	var nullableValue sql.NullString
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query dump_format error: %v", err)
		return
	}
	if nullableValue.Valid {
		format = nullableValue.String
	}
	return
}

// GetDumpTime 获取日志转储时间
func (l *logStrategy) GetDumpTime() (time string, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'dump_time'"
	var nullableValue sql.NullString
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query dump_time error: %v", err)
		return
	}
	if nullableValue.Valid {
		time = nullableValue.String
	}
	return
}

// GetHistoryIsDownloadWithPwd 获取历史日志下载是否需要密码
func (l *logStrategy) GetHistoryIsDownloadWithPwd() (isDownload bool, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'history_log_export_with_pwd'"
	var nullableValue sql.NullBool
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query history_log_export_with_pwd error: %v", err)
		return
	}
	if nullableValue.Valid {
		isDownload = nullableValue.Bool
	}
	return
}

// SetHistoryIsDownloadWithPwd 设置历史日志下载是否需要密码
func (l *logStrategy) SetHistoryIsDownloadWithPwd(isDownload bool) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() + ".t_log_config SET f_value = ? WHERE f_key = 'history_log_export_with_pwd'"
	_, err = l.db.Exec(sqlStr, isDownload)
	if err != nil {
		l.logger.Errorf("db update history_log_export_with_pwd error: %v", err)
		return
	}
	return
}

// GetLogPrefix 获取日志前缀
func (l *logStrategy) GetLogPrefix() (prefix string, err error) {
	sqlStr := "SELECT f_value FROM " + infra.GetDBName() + ".t_log_config WHERE f_key = 'oss_storage_prefix'"
	var nullableValue sql.NullString
	err = l.db.QueryRow(sqlStr).Scan(&nullableValue)
	if err != nil {
		l.logger.Errorf("db query oss_storage_prefix error: %v", err)
		return
	}
	if nullableValue.Valid {
		prefix = nullableValue.String
	}
	return
}

// SetLogPrefix 设置日志前缀 prefix guid
func (l *logStrategy) SetLogPrefix(prefix string) (err error) {
	sqlStr := "INSERT INTO " + infra.GetDBName() + ".t_log_config (f_key, f_value) VALUES ('oss_storage_prefix', ?)"
	_, err = l.db.Exec(sqlStr, prefix)
	if err != nil {
		l.logger.Errorf("db update oss_storage_prefix error: %v", err)
		return
	}
	return
}

// 转存策略配置事务
type dumpStrategyTx struct {
	tx     *sql.Tx
	logger api.Logger
}

// BeginDumpTx 开始日志转存事务
func (l *logStrategy) BeginDumpTx() (tx interfaces.DumpStrategyTx, err error) {
	btx, err := l.db.Begin()
	if err != nil {
		l.logger.Errorf("db begin dump tx error: %v", err)
		return
	}
	tx = &dumpStrategyTx{tx: btx, logger: l.logger}
	return
}

func (d *dumpStrategyTx) Commit() (err error) {
	return d.tx.Commit()
}

func (d *dumpStrategyTx) Rollback() (err error) {
	return d.tx.Rollback()
}

func (d *dumpStrategyTx) SetRetentionPeriod(period int) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() + ".t_log_config SET f_value = ? WHERE f_key = 'retention_period'"
	_, err = d.tx.Exec(sqlStr, period)
	if err != nil {
		d.logger.Errorf("db update retention_period error: %v", err)
		return
	}
	return
}

func (d *dumpStrategyTx) SetRetentionPeriodUnit(unit string) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() + ".t_log_config SET f_value = ? WHERE f_key = 'retention_period_unit'"
	_, err = d.tx.Exec(sqlStr, unit)
	if err != nil {
		d.logger.Errorf("db update retention_period_unit error: %v", err)
		return
	}
	return
}

func (d *dumpStrategyTx) SetDumpFormat(format string) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() + ".t_log_config SET f_value = ? WHERE f_key = 'dump_format'"
	_, err = d.tx.Exec(sqlStr, format)
	if err != nil {
		d.logger.Errorf("db update dump_format error: %v", err)
		return
	}
	return
}

func (d *dumpStrategyTx) SetDumpTime(time string) (err error) {
	sqlStr := "UPDATE " + infra.GetDBName() + ".t_log_config SET f_value = ? WHERE f_key = 'dump_time'"
	_, err = d.tx.Exec(sqlStr, time)
	if err != nil {
		d.logger.Errorf("db update dump_time error: %v", err)
		return
	}
	return
}
