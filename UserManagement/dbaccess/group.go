// Package dbaccess group Anyshare 数据访问层 - 用户组数据库操作
package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type group struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	gOnce sync.Once
	gDB   *group
)

const (
	strDesc = "desc"
	strAsc  = "asc"
)

// NewGroup 创建数据库操作对象--和用户组相关
func NewGroup() *group {
	gOnce.Do(func() {
		gDB = &group{
			db:      dbPool,
			logger:  common.NewLogger(),
			dbTrace: dbTracePool,
			trace:   common.SvcARTrace,
		}
	})

	return gDB
}

// 根据用户名获取用户ID
func (g *group) GetGroupIDByName(name string) (id string, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id from %s.t_group where f_group_name = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := g.db.Query(strSQL, name)
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
		return id, sqlErr
	}

	if rows.Next() {
		if scanErr := rows.Scan(&id); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return id, scanErr
		}
	}

	return id, nil
}

// 根据用户名获取用户ID2 支持trace
func (g *group) GetGroupIDByName2(ctx context.Context, name string) (id string, err error) {
	g.trace.SetClientSpanName("数据库操作-根据组名获取组ID")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id from %s.t_group where f_group_name = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := g.dbTrace.QueryContext(newCtx, strSQL, name)
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
		return id, sqlErr
	}

	if rows.Next() {
		if scanErr := rows.Scan(&id); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return id, scanErr
		}
	}

	return id, nil
}

// AddGroup 添加用户组
func (g *group) AddGroup(ctx context.Context, id, name, notes string, tx *sql.Tx) (err error) {
	g.trace.SetClientSpanName("数据库操作-添加用户组")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	sqlStr := "insert into %s.t_group(f_group_id, f_group_name, f_notes, f_created_time)values(?, ?, ?, ?)"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := tx.ExecContext(newCtx, sqlStr, id, name, notes, currentTime); err != nil {
		return err
	}
	return nil
}

// DeleteGroup 删除用户组
func (g *group) DeleteGroup(id string) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_group where f_group_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if _, err := g.db.Exec(sqlStr, id); err != nil {
		return err
	}
	return nil
}

// ModifyGroup 修改用户组
func (g *group) ModifyGroup(id, name string, nameChanged bool, notes string, notesChanged bool) (err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	sqlStr := "update %s.t_group set "
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	if nameChanged {
		args = append(args, name)
		sqlStr += "f_group_name = ? "
		if notesChanged {
			sqlStr += ","
		}
	}
	if notesChanged {
		args = append(args, notes)
		sqlStr += "f_notes = ? "
	}
	sqlStr += "where f_group_id = ? "
	args = append(args, id)
	if _, err := g.db.Exec(sqlStr, args...); err != nil {
		return err
	}
	return nil
}

// GetGroupByID2 获取指定的用户组， 支持trace
func (g *group) GetGroupByID2(ctx context.Context, id string) (info interfaces.GroupInfo, err error) {
	g.trace.SetClientSpanName("数据库操作-获取指定的用户组信息")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id, f_group_name, f_notes from %s.t_group where f_group_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := g.dbTrace.QueryContext(newCtx, strSQL, id)
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
		return info, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&info.ID, &info.Name, &info.Notes); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return info, scanErr
		}
	}

	return info, nil
}

// GetGroupByID 获取指定的用户组
func (g *group) GetGroupByID(id string) (info interfaces.GroupInfo, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id, f_group_name, f_notes from %s.t_group where f_group_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := g.db.Query(strSQL, id)
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
		return info, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&info.ID, &info.Name, &info.Notes); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return info, scanErr
		}
	}

	return info, nil
}

// GetExistGroupIDs 获取存在的用户组id
func (g *group) GetExistGroupIDs(groupIDs []string) (existGroupIDs []string, err error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}

	set, argIDs := GetFindInSetSQL(groupIDs)
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id from %s.t_group where f_group_id in ("
	strSQL += set
	strSQL += ")"
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

	var existGroupID string
	existGroupIDs = make([]string, 0)
	for rows.Next() {
		if scanErr := rows.Scan(&existGroupID); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		existGroupIDs = append(existGroupIDs, existGroupID)
	}

	return existGroupIDs, nil
}

// GetGroups 列举符合条件的用户组
func (g *group) GetGroups(info interfaces.SearchInfo) (groups []interfaces.GroupInfo, err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id, f_group_name, f_notes from %s.t_group "
	strSQL = fmt.Sprintf(strSQL, dbName)
	if info.HasKeyWord {
		args = append(args, "%"+info.Keyword+"%")
		strSQL += "where f_group_name like ? "
	}
	if info.Sort == interfaces.Name {
		strSQL += "order by upper(f_group_name) "
	} else {
		strSQL += "order by f_created_time "
	}

	if info.Direction == interfaces.Desc {
		strSQL += strDesc
	} else {
		strSQL += strAsc
	}

	// 添加第二排序规则 防止重复
	strSQL += ",  f_group_id asc limit ?, ? "

	args = append(args, info.Offset, info.Limit)

	rows, sqlErr := g.db.Query(strSQL, args...)
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

	var tmpInfo interfaces.GroupInfo
	groups = make([]interfaces.GroupInfo, 0)
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name, &tmpInfo.Notes); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		groups = append(groups, tmpInfo)
	}

	return groups, nil
}

// GetGroups 列举符合条件的用户组
func (g *group) GetGroupsNum(info interfaces.SearchInfo) (num int, err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	strSQL := "select count(*) from %s.t_group "
	strSQL = fmt.Sprintf(strSQL, dbName)
	if info.HasKeyWord {
		args = append(args, "%"+info.Keyword+"%")
		strSQL += "where f_group_name like ? "
	}

	rows, sqlErr := g.db.Query(strSQL, args...)
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
		return num, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&num); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return num, scanErr
		}
	}

	return num, nil
}

// SearchGroupByKeyword 用户组关键字搜索
func (g *group) SearchGroupByKeyword(keyword string, start, limit int) (out []interfaces.NameInfo, err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	strSQL := `select f_group_id, f_group_name
			from %s.t_group
			where f_group_name like ?
			order by case when f_group_name = ? then 0 when f_group_name like ? then 1 else 2 end,upper(f_group_name)
			limit ?, ? `
	strSQL = fmt.Sprintf(strSQL, dbName)
	args = append(args, "%"+keyword+"%", keyword, "%"+keyword+"%", start, limit)
	rows, sqlErr := g.db.Query(strSQL, args...)
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

	out = make([]interfaces.NameInfo, 0)
	for rows.Next() {
		var temp interfaces.NameInfo
		if scanErr := rows.Scan(&temp.ID, &temp.Name); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		out = append(out, temp)
	}

	return out, nil
}

// SearchGroupNumByKeyword 用户组关键字搜索符合条件的用户组数量
func (g *group) SearchGroupNumByKeyword(keyword string) (num int, err error) {
	var args []interface{}
	dbName := common.GetDBName("user_management")
	strSQL := "select count(*) from %s.t_group where f_group_name like ? "
	strSQL = fmt.Sprintf(strSQL, dbName)

	args = append(args, "%"+keyword+"%")
	rows, sqlErr := g.db.Query(strSQL, args...)
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
		return num, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&num); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return num, scanErr
		}
	}

	return num, nil
}

// GetGroupName 根据用户组ID获取用户组名
func (g *group) GetGroupName(ids []string) (nameInfo []interfaces.NameInfo, exsitIDs []string, err error) {
	// 分割ID 防止sql过长
	splitedIDs := SplitArray(ids)
	ctx := context.Background()

	nameInfo = make([]interfaces.NameInfo, 0)
	exsitIDs = make([]string, 0)
	for _, ids := range splitedIDs {
		infoTmp, idTmp, tmpErr := g.getGroupNameSingle(ctx, ids)
		if tmpErr != nil {
			return nil, nil, tmpErr
		}

		nameInfo = append(nameInfo, infoTmp...)
		exsitIDs = append(exsitIDs, idTmp...)
	}
	return nameInfo, exsitIDs, err
}

// GetGroupName2 根据用户组ID获取用户组名，支持trace
func (g *group) GetGroupName2(ctx context.Context, ids []string) (nameInfo []interfaces.NameInfo, exsitIDs []string, err error) {
	// 分割ID 防止sql过长
	splitedIDs := SplitArray(ids)

	nameInfo = make([]interfaces.NameInfo, 0)
	exsitIDs = make([]string, 0)
	for _, ids := range splitedIDs {
		infoTmp, idTmp, tmpErr := g.getGroupNameSingle(ctx, ids)
		if tmpErr != nil {
			return nil, nil, tmpErr
		}

		nameInfo = append(nameInfo, infoTmp...)
		exsitIDs = append(exsitIDs, idTmp...)
	}
	return nameInfo, exsitIDs, err
}

func (g *group) getGroupNameSingle(ctx context.Context, ids []string) (nameInfo []interfaces.NameInfo, exsitIDs []string, err error) {
	g.trace.SetClientSpanName("数据库操作-获取组名")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	if len(ids) == 0 {
		return nil, nil, nil
	}

	set, argIDs := GetFindInSetSQL(ids)
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id, f_group_name from %s.t_group where f_group_id in ("
	strSQL += set
	strSQL += ")"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := g.dbTrace.QueryContext(newCtx, strSQL, argIDs...)
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
		g.logger.Errorln(sqlErr, strSQL, argIDs)
		return nil, nil, sqlErr
	}

	nameInfo = make([]interfaces.NameInfo, 0)
	exsitIDs = make([]string, 0)
	var tmpInfo interfaces.NameInfo
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.Name); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, nil, scanErr
		}

		nameInfo = append(nameInfo, tmpInfo)
		exsitIDs = append(exsitIDs, tmpInfo.ID)
	}

	return nameInfo, exsitIDs, nil
}
