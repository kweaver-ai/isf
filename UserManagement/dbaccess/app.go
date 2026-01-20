package dbaccess

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type app struct {
	db  *sqlx.DB
	log common.Logger
}

var (
	aOnce sync.Once
	a     *app
)

// NewApp 创建应用账户操作对象
func NewApp() *app {
	aOnce.Do(func() {
		a = &app{
			db:  dbPool,
			log: common.NewLogger(),
		}
	})

	return a
}

// Register 注册应用账户
func (a *app) RegisterApp(info *interfaces.AppCompleteInfo, tx *sql.Tx) (err error) {
	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("insert into %s.t_app "+
		"(f_id, f_name, f_password, f_type, f_created_time, f_credential_type)"+
		"values (?, ?, ?, ?, ?, ?)", dbName)

	if _, err = tx.Exec(sqlStr,
		info.AppInfo.ID,
		info.AppInfo.Name,
		info.Password,
		info.Type,
		currentTime,
		info.CredentialType,
	); err != nil {
		a.log.Errorln(err, sqlStr, info)
		return
	}

	return
}

// Delete 删除应用账户
func (a *app) DeleteApp(id string, tx *sql.Tx) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("delete from %s.t_app where f_id = ?", dbName)
	if _, err = tx.Exec(sqlStr, id); err != nil {
		a.log.Errorln(err, sqlStr)
		return
	}

	return
}

// Update 更新应用账户
func (a *app) UpdateApp(id string, bName bool, name string, bPwd bool, pwd string, tx *sql.Tx) (err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("update %s.t_app set ", dbName)

	if bName {
		args = append(args, name)
		sqlStr += "f_name = ? "
		if pwd != "" {
			sqlStr += ","
		}
	}
	if bPwd {
		args = append(args, pwd)
		sqlStr += "f_password = ? "
	}

	sqlStr += "where f_id = ? "
	args = append(args, id)
	if _, err = tx.Exec(sqlStr, args...); err != nil {
		a.log.Errorln(err, sqlStr)
		return
	}

	return
}

// List 应用账户列表
func (a *app) AppList(searchInfo *interfaces.SearchInfo) (info *[]interfaces.AppInfo, err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("select f_id, f_name, f_credential_type from %s.t_app where f_type = 1 ", dbName)

	var args []interface{}
	if searchInfo.Keyword != "" {
		args = append(args, "%"+searchInfo.Keyword+"%")
		sqlStr += "and f_name like ? "
	}

	if searchInfo.Sort == interfaces.DateCreated {
		sqlStr += "order by f_created_time "
	} else {
		sqlStr += "order by upper(f_name) "
	}

	if searchInfo.Direction == interfaces.Desc {
		sqlStr += strDesc
	} else {
		sqlStr += strAsc
	}

	sqlStr += ", f_id asc LIMIT ?, ?"
	args = append(args, searchInfo.Offset, searchInfo.Limit)
	rows, err := a.db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.log.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				a.log.Errorln(closeErr)
			}
		}
	}()

	tmpInfo := interfaces.AppInfo{}
	infoList := make([]interfaces.AppInfo, 0)
	for rows.Next() {
		if err := rows.Scan(
			&tmpInfo.ID,
			&tmpInfo.Name,
			&tmpInfo.CredentialType,
		); err != nil {
			return nil, err
		}

		infoList = append(infoList, tmpInfo)
	}

	return &infoList, nil
}

// ListNum 应用账户总数
func (a *app) AppListCount(searchInfo *interfaces.SearchInfo) (num int, err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("select count(*) from %s.t_app where f_type = 1", dbName)

	var args []interface{}
	if searchInfo.Keyword != "" {
		args = append(args, "%"+searchInfo.Keyword+"%")
		sqlStr += " and f_name like ? "
	}

	rows, err := a.db.Query(sqlStr, args...)
	if err != nil {
		return
	}
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.log.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				a.log.Errorln(closeErr)
			}
		}
	}()

	for rows.Next() {
		if err = rows.Scan(&num); err != nil {
			a.log.Errorln(err, sqlStr)
			return
		}
	}

	return
}

// GetAppByName 获取应用账户
func (a *app) GetAppByName(name string) (appinfo *interfaces.AppInfo, err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("select f_id, f_name from %s.t_app where f_name = ?", dbName)

	row := a.db.QueryRow(sqlStr, name)
	info := interfaces.AppInfo{}
	err = row.Scan(
		&info.ID,
		&info.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return
	}

	return &info, nil
}

// GetAppByID 获取应用账户
func (a *app) GetAppByID(id string) (appinfo *interfaces.AppInfo, err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := fmt.Sprintf("select f_id, f_name, f_credential_type from %s.t_app where f_id = ?", dbName)

	row := a.db.QueryRow(sqlStr, id)

	info := interfaces.AppInfo{}
	err = row.Scan(
		&info.ID,
		&info.Name,
		&info.CredentialType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return
	}

	return &info, nil
}

// GetAppName 获取应用账户名
func (a *app) GetAppName(ids []string) (nameInfo []interfaces.NameInfo, exsitIDs []string, err error) {
	if len(ids) == 0 {
		return
	}

	// 分割ID 防止sql过长
	splitedIDs := SplitArray(ids)

	nameInfo = make([]interfaces.NameInfo, 0)
	exsitIDs = make([]string, 0)
	for _, ids := range splitedIDs {
		infoTmp, idTmp, tmpErr := a.getAppNameSingle(ids)
		if tmpErr != nil {
			return nil, nil, tmpErr
		}

		nameInfo = append(nameInfo, infoTmp...)
		exsitIDs = append(exsitIDs, idTmp...)
	}
	return nameInfo, exsitIDs, err
}

func (a *app) getAppNameSingle(ids []string) (nameInfo []interfaces.NameInfo, exsitIDs []string, err error) {
	set, argIDs := GetFindInSetSQL(ids)
	dbName := common.GetDBName("user_management")
	strSQL := "select f_id, f_name from %s.t_app where f_id in ("
	strSQL += set
	strSQL += ")"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := a.db.Query(strSQL, argIDs...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				a.log.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				a.log.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		a.log.Errorln(sqlErr, strSQL, argIDs)
		return nil, nil, sqlErr
	}

	nameInfo = make([]interfaces.NameInfo, 0)
	exsitIDs = make([]string, 0)
	var tmpInfo interfaces.NameInfo
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name); scanErr != nil {
			a.log.Errorln(scanErr, strSQL)
			return nil, nil, scanErr
		}

		nameInfo = append(nameInfo, tmpInfo)
		exsitIDs = append(exsitIDs, tmpInfo.ID)
	}

	return nameInfo, exsitIDs, nil
}
