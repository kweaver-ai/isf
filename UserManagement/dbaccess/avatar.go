package dbaccess

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type avatar struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	avOnce sync.Once
	av     *avatar
)

// NewAvatar 创建头像操作对象
func NewAvatar() *avatar {
	avOnce.Do(func() {
		av = &avatar{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return av
}

// Get 获取用户当前头像信息
func (a *avatar) Get(userID string) (info interfaces.AvatarOSSInfo, err error) {
	if userID == "" {
		return
	}
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("select f_oss_id, f_key, f_type, f_time from %s.t_avatar where f_user_id = ? and f_status = 1", dbName)

	rows, sqlErr := a.db.Query(strSQL, userID)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				a.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		a.logger.Errorln(sqlErr, strSQL)
		return info, sqlErr
	}

	for rows.Next() {
		if err = rows.Scan(&info.OSSID, &info.Key, &info.Type, &info.Time); err != nil {
			a.logger.Errorln(err, strSQL)
			return
		}

		info.BUseful = true
		info.UserID = userID
	}

	return
}

// Add 新增用户头像信息
func (a *avatar) Add(info *interfaces.AvatarOSSInfo) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("insert into %s.t_avatar ( f_user_id, f_oss_id, f_key, f_type, f_status, f_time) values ( ?, ?, ?, ?, ?, ?)", dbName)
	_, err = a.db.Exec(sqlStr, info.UserID, info.OSSID, info.Key, info.Type, info.BUseful, info.Time)
	return
}

// UpdateStatusByKey 根据key更新用户头像状态
func (a *avatar) UpdateStatusByKey(key string, status bool, tx *sql.Tx) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("update %s.t_avatar set f_status = ? where f_key = ?", dbName)
	_, err = tx.Exec(sqlStr, status, key)
	return
}

// SetAvatarUnableByID 根据用户ID禁用用户头像
func (a *avatar) SetAvatarUnableByID(userID string, tx *sql.Tx) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("update %s.t_avatar set f_status = 0 where f_user_id = ? and f_status = 1", dbName)
	_, err = tx.Exec(sqlStr, userID)
	return
}

// Delete 删除用户头像信息
func (a *avatar) Delete(key string) (err error) {
	if key == "" {
		return
	}
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("delete from %s.t_avatar where f_key = ?", dbName)
	_, err = a.db.Exec(sqlStr, key)
	return
}

// GetUselessAvatar 获取超时未用到的头像信息
func (a *avatar) GetUselessAvatar(time int64) (data []interfaces.AvatarOSSInfo, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := fmt.Sprintf("select f_user_id, f_oss_id, f_key, f_type, f_status, f_time from %s.t_avatar where f_time < ? and f_status = 0", dbName)

	rows, sqlErr := a.db.Query(strSQL, time)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				a.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		a.logger.Errorln(sqlErr, strSQL)
		return data, sqlErr
	}

	data = make([]interfaces.AvatarOSSInfo, 0)
	var info interfaces.AvatarOSSInfo
	for rows.Next() {
		if err = rows.Scan(&info.UserID, &info.OSSID, &info.Key, &info.Type, &info.BUseful, &info.Time); err != nil {
			a.logger.Errorln(err, strSQL)
			return
		}

		data = append(data, info)
	}

	return
}
