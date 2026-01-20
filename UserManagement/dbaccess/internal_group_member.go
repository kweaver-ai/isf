package dbaccess

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type internalGroupMember struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	igmOnce sync.Once
	igmDB   *internalGroupMember
)

// NewInternalGroupMember 创建数据库操作对象--和内部组成员相关
func NewInternalGroupMember() *internalGroupMember {
	igmOnce.Do(func() {
		igmDB = &internalGroupMember{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return igmDB
}

// Add 增加内部组成员
func (gm *internalGroupMember) Add(groupID string, infos []interfaces.InternalGroupMember, tx *sql.Tx) (err error) {
	if len(infos) == 0 {
		return nil
	}

	currentTime := time.Now().UnixNano()
	strArgs := make([]string, 0, len(infos))
	strValues := make([]interface{}, 0)
	for _, v := range infos {
		strArgs = append(strArgs, "(?, ?, ?, ?)")
		strValues = append(strValues, groupID, v.ID, v.Type, currentTime)
	}
	dbName := common.GetDBName("user_management")
	sqlStr := "insert into %s.t_internal_group_member(f_internal_group_id, f_member_id, f_member_type, f_added_time) values "
	sqlStr += strings.Join(strArgs, ",")
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.Exec(sqlStr, strValues...); err != nil {
		return err
	}
	return nil
}

// DeleteAll 删除内部组内所有成员
func (gm *internalGroupMember) DeleteAll(ids []string, tx *sql.Tx) (err error) {
	if len(ids) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_internal_group_member where f_internal_group_id in ( "
	sqlStr += set
	sqlStr += " )"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.Exec(sqlStr, argIDs...); err != nil {
		return err
	}
	return nil
}

// Get 获取内部组成员
func (gm *internalGroupMember) Get(groupID string) (outInfo []interfaces.InternalGroupMember, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := `select f_member_id, f_member_type from %s.t_internal_group_member where f_internal_group_id = ?`
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := gm.db.Query(strSQL, groupID)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				gm.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				gm.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		gm.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	var tmpInfo interfaces.InternalGroupMember
	outInfo = make([]interfaces.InternalGroupMember, 0)
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Type); scanErr != nil {
			gm.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		outInfo = append(outInfo, tmpInfo)
	}

	return outInfo, nil
}

// GetBelongGroups 获取用户成员所属内部组
func (gm *internalGroupMember) GetBelongGroups(info interfaces.InternalGroupMember) (ids []string, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "select f_internal_group_id from %s.t_internal_group_member where f_member_id = ? and f_member_type = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := gm.db.Query(strSQL, info.ID, info.Type)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				gm.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				gm.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		gm.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	ids = make([]string, 0)
	for rows.Next() {
		var strID string
		if scanErr := rows.Scan(&strID); scanErr != nil {
			gm.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		ids = append(ids, strID)
	}

	return ids, nil
}
