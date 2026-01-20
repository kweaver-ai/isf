package dbaccess

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type reservedName struct {
	logger common.Logger
	db     *sqlx.DB
}

var (
	rnOnce sync.Once
	rn     *reservedName
)

// NewReservedName 创建数据库操作对象
func NewReservedName() *reservedName {
	rnOnce.Do(func() {
		rn = &reservedName{
			logger: common.NewLogger(),
			db:     dbPool,
		}
	})
	return rn
}

// AddReservedName 添加保留名称
func (r *reservedName) AddReservedName(name interfaces.ReservedNameInfo, tx *sql.Tx) error {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "insert into %s.t_reserved_name(f_id, f_name, f_create_time, f_update_time) values(?, ?, ?, ?)"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	_, err := tx.Exec(sqlStr, name.ID, name.Name, name.CreateTime, name.UpdateTime)
	if err != nil {
		r.logger.Errorln("failed to add reserved name, err:", err)
	}
	return err
}

// UpdateReservedName 更新保留名称
func (r *reservedName) UpdateReservedName(name interfaces.ReservedNameInfo, tx *sql.Tx) error {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "update %s.t_reserved_name set f_name = ?, f_update_time = ? where f_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	_, err := tx.Exec(sqlStr, name.Name, name.UpdateTime, name.ID)
	if err != nil {
		r.logger.Errorln("failed to update reserved name, err:", err)
	}
	return err
}

// GetReservedNameByID 根据ID获取保留名称
func (r *reservedName) GetReservedNameByID(id string, tx *sql.Tx) (interfaces.ReservedNameInfo, bool, error) {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "select f_id, f_name, f_create_time, f_update_time from %s.t_reserved_name where f_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	var res interfaces.ReservedNameInfo
	err := tx.QueryRow(sqlStr, id).Scan(&res.ID, &res.Name, &res.CreateTime, &res.UpdateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, false, nil
		}
		r.logger.Errorln("failed to get record when GetReservedNameByID, err:", err)
		return res, false, err
	}
	return res, true, nil
}

// GetReservedNameByName 根据name获取保留名称
func (r *reservedName) GetReservedNameByName(name string, tx *sql.Tx) (interfaces.ReservedNameInfo, bool, error) {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "select f_id, f_name, f_create_time, f_update_time from %s.t_reserved_name where f_name = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	var res interfaces.ReservedNameInfo
	var err error
	if tx != nil {
		err = tx.QueryRow(sqlStr, name).Scan(&res.ID, &res.Name, &res.CreateTime, &res.UpdateTime)
	} else {
		err = r.db.QueryRow(sqlStr, name).Scan(&res.ID, &res.Name, &res.CreateTime, &res.UpdateTime)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, false, nil
		}
		r.logger.Errorln("failed to get record when GetReservedNameByName, err:", err)
		return res, false, err
	}
	return res, true, nil
}

// DeleteReservedName 删除保留名称
func (r *reservedName) DeleteReservedName(id string) error {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "delete from %s.t_reserved_name where f_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	_, err := r.db.Exec(sqlStr, id)
	if err != nil {
		r.logger.Errorln("failed to delete reserved name, err:", err)
	}
	return err
}

func (r *reservedName) GetLock(tx *sql.Tx) error {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "select f_value from %s.t_sharemgnt_config where f_key = 'reserved_name_lock' for update"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	_, err := tx.Exec(sqlStr)
	if err != nil {
		r.logger.Errorln("failed to get lock, err:", err)
	}
	return err
}
