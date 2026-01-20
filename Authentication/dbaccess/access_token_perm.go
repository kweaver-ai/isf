// Package dbaccess 数据访问层
package dbaccess

import (
	"fmt"
	"sync"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authentication/common"
)

type accessTokenPerm struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	acOnce sync.Once
	a      *accessTokenPerm
)

// NewAccessTokenPerm 创建数据库对象
func NewAccessTokenPerm() *accessTokenPerm {
	acOnce.Do(func() {
		a = &accessTokenPerm{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})
	return a
}

func (a *accessTokenPerm) CheckAppAccessTokenPerm(appID string) (bool, error) {
	dbName := common.GetDBName("authentication")
	strSQL := "select f_app_id from %s.t_access_token_perm where f_app_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, err := a.db.Query(strSQL, appID)
	if err != nil {
		a.logger.Errorln(err, strSQL)
		return false, err
	}

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

	for rows.Next() {
		return true, nil
	}
	return false, nil
}

func (a *accessTokenPerm) AddAppAccessTokenPerm(appID string) (err error) {
	curTime := time.Now().UnixNano() / 1000
	strSQL := "insert into t_access_token_perm(f_app_id, f_create_time) values(?,?)"
	_, err = a.db.Exec(strSQL, appID, curTime)
	if err != nil {
		a.logger.Errorln(err)
		return err
	}
	return
}

func (a *accessTokenPerm) DeleteAppAccessTokenPerm(appID string) (err error) {
	strSQL := "delete from t_access_token_perm where f_app_id = ?"
	_, err = a.db.Exec(strSQL, appID)
	if err != nil {
		a.logger.Errorln(err)
		return err
	}
	return
}

func (a *accessTokenPerm) GetAllAppAccessTokenPerm() (permApps []string, err error) {
	permApps = make([]string, 0)
	dbName := common.GetDBName("authentication")
	strSQL := fmt.Sprintf("select f_app_id from %s.t_access_token_perm", dbName)
	rows, err := a.db.Query(strSQL)
	if err != nil {
		a.logger.Errorln(err, strSQL)
		return
	}

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

	for rows.Next() {
		var app string
		err = rows.Scan(&app)
		if err != nil {
			return
		}
		permApps = append(permApps, app)
	}
	return
}
