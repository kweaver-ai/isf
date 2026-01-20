// Package dbaccess contactor Anyshare 数据访问层 - 联系人组数据库操作
package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type contactor struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	contactorOnce sync.Once
	cDB           *contactor
)

// NewContactor 创建数据库操作对象
func NewContactor() *contactor {
	contactorOnce.Do(func() {
		cDB = &contactor{
			db:      dbPool,
			logger:  common.NewLogger(),
			trace:   common.SvcARTrace,
			dbTrace: dbTracePool,
		}
	})

	return cDB
}

// GetContactorName 批量获取联系人组名
func (d *contactor) GetContactorName(contactorIDs []string) (info []interfaces.NameInfo, existIDs []string, err error) {
	// 分割ID 防止sql过长
	splitedIDs := SplitArray(contactorIDs)

	info = make([]interfaces.NameInfo, 0)
	existIDs = make([]string, 0)
	for _, ids := range splitedIDs {
		infoTmp, idTmp, tmpErr := d.getContactorNameSingle(ids)
		if tmpErr != nil {
			return nil, nil, tmpErr
		}

		info = append(info, infoTmp...)
		existIDs = append(existIDs, idTmp...)
	}
	return info, existIDs, err
}

func (d *contactor) getContactorNameSingle(contactorIDs []string) (info []interfaces.NameInfo, existIDs []string, err error) {
	existIDs = make([]string, 0)
	info = make([]interfaces.NameInfo, 0)
	if len(contactorIDs) == 0 {
		return nil, nil, nil
	}

	set, argIDs := GetFindInSetSQL(contactorIDs)
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "select f_group_id, f_group_name from %s.t_person_group where f_group_id in ("
	strSQL += set
	strSQL += ")"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := d.db.Query(strSQL, argIDs...)
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
		d.logger.Errorln(sqlErr, strSQL, argIDs)
		return nil, nil, sqlErr
	}

	var tmpInfo interfaces.NameInfo
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name); scanErr != nil {
			d.logger.Errorln(scanErr, strSQL)
			return nil, nil, scanErr
		}

		existIDs = append(existIDs, tmpInfo.ID)
		info = append(info, tmpInfo)
	}

	return info, existIDs, err
}

// GetUserAllBelongContactorIDs 获取包含用户的所有联系人组
func (d *contactor) GetUserAllBelongContactorIDs(userID string) (contactorIDs []string, err error) {
	contactorIDs = make([]string, 0)

	dbName := common.GetDBName("sharemgnt_db")
	strSQL := fmt.Sprintf("select f_group_id from %s.t_contact_person where f_user_id = ?", dbName)

	rows, sqlErr := d.db.Query(strSQL, userID)
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
		return nil, sqlErr
	}

	var contactID string
	for rows.Next() {
		if scanErr := rows.Scan(&contactID); scanErr != nil {
			d.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		contactorIDs = append(contactorIDs, contactID)
	}

	return contactorIDs, err
}

// GetContactorInfo  获取联系人组信息
func (d *contactor) GetContactorInfo(contactorID string) (result bool, info interfaces.ContactorInfo, err error) {
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := fmt.Sprintf("select f_user_id from %s.t_person_group where f_group_id = ?", dbName)

	rows, sqlErr := d.db.Query(strSQL, contactorID)
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
		return false, info, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&info.UserID); scanErr != nil {
			d.logger.Errorln(scanErr, strSQL)
			return false, info, scanErr
		}
		result = true
	}
	return result, info, nil
}

// DeleteContactor 删除联系人组
func (d *contactor) DeleteContactors(contactorIDs []string, tx *sql.Tx) (err error) {
	if len(contactorIDs) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(contactorIDs)
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "delete from %s.t_person_group where f_group_id in ("
	sqlStr += set
	sqlStr += ")"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.Exec(sqlStr, argIDs...); err != nil {
		return err
	}
	return nil
}

// DeleteContactorMembers 删除联系人组中的成员
func (d *contactor) DeleteContactorMembers(contactorIDs []string, tx *sql.Tx) (err error) {
	if len(contactorIDs) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(contactorIDs)
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "delete from %s.t_contact_person where f_group_id in ("
	sqlStr += set
	sqlStr += ")"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.Exec(sqlStr, argIDs...); err != nil {
		return err
	}
	return nil
}

// GetAllContactorInfos 获取用户所有的联系人组信息
func (d *contactor) GetAllContactorInfos(userID string) (infos []interfaces.ContactorInfo, err error) {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "select f_group_id from %s.t_person_group where f_user_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	rows, sqlErr := d.db.Query(sqlStr, userID)
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
		d.logger.Errorln(sqlErr, sqlStr)
		return nil, sqlErr
	}

	infos = make([]interfaces.ContactorInfo, 0)
	for rows.Next() {
		var info interfaces.ContactorInfo
		if scanErr := rows.Scan(&info.ContactorID); scanErr != nil {
			d.logger.Errorln(scanErr, sqlStr)
			return nil, scanErr
		}

		infos = append(infos, info)
	}
	return infos, nil
}

// DeleteUserInContactors 在联系人组内删除用户
func (d *contactor) DeleteUserInContactors(userID string, tx *sql.Tx) (err error) {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "delete from %s.t_contact_person where f_user_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.Exec(sqlStr, userID); err != nil {
		return err
	}
	return nil
}

// UpdateContactorCount 更新联系人组的联系人数量信息
func (d *contactor) UpdateContactorCount() (err error) {
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "update %s.t_person_group as g set f_person_count = (select count(*) from %s.t_contact_person where f_group_id = g.f_group_id)"
	sqlStr = fmt.Sprintf(sqlStr, dbName, dbName)
	if _, err := d.db.Exec(sqlStr); err != nil {
		return err
	}
	return nil
}

// GetContactorMemberIDs 批量获取联系人组成员ID
func (d *contactor) GetContactorMemberIDs(ctx context.Context, contactorIDs []string) (infos map[string][]string, err error) {
	// 分割ID 防止sql过长
	splitedIDs := SplitArray(contactorIDs)

	infos = make(map[string][]string)
	for _, ids := range splitedIDs {
		infoTmp, tmpErr := d.getContactorMemberIDsSingle(ctx, ids)
		if tmpErr != nil {
			return nil, tmpErr
		}

		for contactorID, userIDs := range infoTmp {
			infos[contactorID] = userIDs
		}
	}
	return infos, err
}

func (d *contactor) getContactorMemberIDsSingle(ctx context.Context, contactorIDs []string) (infos map[string][]string, err error) {
	d.trace.SetClientSpanName("数据库操作-批量获取联系人组成员ID")
	newCtx, span := d.trace.AddClientTrace(ctx)
	defer func() { d.trace.TelemetrySpanEnd(span, err) }()

	infos = make(map[string][]string)

	set, argIDs := GetFindInSetSQL(contactorIDs)
	dbName := common.GetDBName("sharemgnt_db")
	sqlStr := "select f_group_id, f_user_id from %s.t_contact_person where f_group_id in ("
	sqlStr += set
	sqlStr += ")"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	rows, sqlErr := d.dbTrace.QueryContext(newCtx, sqlStr, argIDs...)
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
		d.logger.Errorln(sqlErr, sqlStr)
		return nil, sqlErr
	}

	for rows.Next() {
		var contactorID, userID string
		if scanErr := rows.Scan(&contactorID, &userID); scanErr != nil {
			d.logger.Errorln(scanErr, sqlStr)
			return nil, scanErr
		}

		if _, ok := infos[contactorID]; !ok {
			infos[contactorID] = make([]string, 0)
		}

		infos[contactorID] = append(infos[contactorID], userID)
	}
	return infos, nil
}
