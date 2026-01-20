// Package dbaccess anonymous Anyshare 数据访问层 - 匿名账户数据库操作
package dbaccess

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type anonymous struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	anonyOnce sync.Once
	anonyDB   *anonymous
)

// NewAnonymous 创建数据库操作对象
func NewAnonymous() *anonymous {
	anonyOnce.Do(func() {
		anonyDB = &anonymous{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return anonyDB
}

// Create 创建匿名账户
func (d *anonymous) Create(info *interfaces.AnonymousInfo) error {
	dbName := common.GetDBName("user_management")
	checkSQL := fmt.Sprintf("select f_password from %s.t_anonymity where f_anonymity_id = ?", dbName)
	rows, sqlErr := d.db.Query(checkSQL, info.ID)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				d.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				d.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		d.logger.Errorln(sqlErr, checkSQL)
		return sqlErr
	}

	bExist := false
	for rows.Next() {
		bExist = true
	}

	verifyMobile := 0
	if info.VerifyMobile {
		verifyMobile = 1
	}
	createTimeStamp := common.Now().UnixNano()
	strSQL := fmt.Sprintf("insert into %s.t_anonymity(f_password, f_expires_at, f_limited_times, f_created_at, f_type, f_verify_mobile, f_anonymity_id) values(?, ?, ?, ?, ?, ?, ?)", dbName)
	if bExist {
		strSQL = fmt.Sprintf("update %s.t_anonymity set f_password = ?, f_expires_at = ?, f_limited_times = ?, f_created_at = ?, f_type = ?, f_verify_mobile = ? where f_anonymity_id = ?", dbName)
	}
	_, sqlErr = d.db.Exec(strSQL, info.Password, info.ExpiresAtStamp, info.LimitedTimes, createTimeStamp, info.Type, verifyMobile, info.ID)
	return sqlErr
}

// DeleteByID 删除匿名账户
func (d *anonymous) DeleteByID(id string) error {
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("delete from %s.t_anonymity where f_anonymity_id = ?", dbName)

	_, sqlErr := d.db.Exec(strSQL, id)
	return sqlErr
}

// GetAccount 获取匿名账户
func (d *anonymous) GetAccount(id string) (info interfaces.AnonymousInfo, err error) {
	curTimeStamp := common.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("select f_anonymity_id, f_password, f_expires_at, f_limited_times, f_accessed_times, f_type, f_verify_mobile from %s.t_anonymity"+
		" where f_anonymity_id = ? and (f_expires_at > ? or f_expires_at = 0)", dbName)

	rows, sqlErr := d.db.Query(strSQL, id, curTimeStamp)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				d.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				d.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		d.logger.Errorln(sqlErr, strSQL)
		return info, sqlErr
	}

	verifyMobile := 0
	for rows.Next() {
		if scanErr := rows.Scan(&info.ID, &info.Password, &info.ExpiresAtStamp, &info.LimitedTimes, &info.AccessedTimes, &info.Type, &verifyMobile); scanErr != nil {
			d.logger.Errorln(scanErr, strSQL)
			return info, scanErr
		}
	}

	if verifyMobile == 1 {
		info.VerifyMobile = true
	}
	return info, nil
}

// AddAccessTimes 访问计数+1
func (d *anonymous) AddAccessTimes(id string, tx *sql.Tx) error {
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("update %s.t_anonymity set f_accessed_times = f_accessed_times + 1 where f_anonymity_id = ?", dbName)
	_, sqlErr := tx.Exec(strSQL, id)
	return sqlErr
}

// DeleteByTime 删除过期匿名账户
func (d *anonymous) DeleteByTime(curTime int64) error {
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("delete from %s.t_anonymity where f_expires_at != ? and f_expires_at < ?", dbName)
	if _, err := d.db.Exec(strSQL, 0, curTime); err != nil {
		d.logger.Errorln(err, strSQL, curTime)
		return err
	}
	return nil
}
