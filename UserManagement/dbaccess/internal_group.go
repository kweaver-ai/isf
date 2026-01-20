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

type internalGroup struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	igOnce sync.Once
	igDB   *internalGroup
)

// NewInternalGroup 创建数据库操作对象--和内部组相关
func NewInternalGroup() *internalGroup {
	igOnce.Do(func() {
		igDB = &internalGroup{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return igDB
}

// Add 新增内部组
func (g *internalGroup) Add(id string) (err error) {
	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	sqlStr := "insert into %s.t_internal_group (f_id, f_created_time)values(?, ?) "
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := g.db.Exec(sqlStr, id, currentTime); err != nil {
		return err
	}
	return nil
}

// Delete 删除内部组
func (g *internalGroup) Delete(ids []string, tx *sql.Tx) (err error) {
	if len(ids) == 0 {
		return nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_internal_group where f_id in ( "
	sqlStr += set
	sqlStr += " )"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := tx.Exec(sqlStr, argIDs...); err != nil {
		return err
	}
	return nil
}

// Get 获取内部组信息
func (g *internalGroup) Get(ids []string) (infos map[string]interfaces.InternelGroup, err error) {
	infos = make(map[string]interfaces.InternelGroup)
	if len(ids) == 0 {
		return nil, nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	dbName := common.GetDBName("user_management")
	strSQL := "select f_id from %s.t_internal_group where f_id in ( "
	strSQL += set
	strSQL += " )"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := g.db.Query(strSQL, argIDs...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				g.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				g.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		g.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	var tmpInfo interfaces.InternelGroup
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		infos[tmpInfo.ID] = tmpInfo
	}

	return infos, nil
}
