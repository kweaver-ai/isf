// Package dbaccess user Anyshare 数据访问层 - 用户数据库操作
package dbaccess

import (
	"context"
	"fmt"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type role struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	rOnce sync.Once
	rDB   *role
)

// NewRole 创建数据库操作对象--和用户相关
func NewRole() *role {
	rOnce.Do(func() {
		rDB = &role{
			db:      dbPool,
			dbTrace: dbTracePool,
			logger:  common.NewLogger(),
			trace:   common.SvcARTrace,
		}
	})

	return rDB
}

// GetRolesByUserIDs 根据用户ID数组批量获取用户角色
func (r *role) GetRolesByUserIDs(userIDs []string) (out map[string]map[interfaces.Role]bool, err error) {
	if len(userIDs) == 0 {
		return
	}

	userSet, userArgIDs := GetFindInSetSQL(userIDs)
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "select f_user_id, f_role_id from %s.t_user_role_relation where f_user_id in ( "
	strSQL += userSet
	strSQL += " )"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.db.Query(strSQL, userArgIDs...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	var roleID string
	var userID string
	out = make(map[string]map[interfaces.Role]bool)
	for rows.Next() {
		if err := rows.Scan(&userID, &roleID); err != nil {
			r.logger.Errorln(err, strSQL)
			return nil, err
		}

		// 数据汇总
		roleInfos, ok := out[userID]
		if !ok {
			temp := make(map[interfaces.Role]bool)
			temp[interfaces.Role(roleID)] = true
			out[userID] = temp
		} else {
			roleInfos[interfaces.Role(roleID)] = true
		}
	}

	return out, nil
}

// GetRolesByUserIDs2 根据用户ID数组批量获取用户角色，支持trace
func (r *role) GetRolesByUserIDs2(ctx context.Context, userIDs []string) (out map[string]map[interfaces.Role]bool, err error) {
	// trace
	r.trace.SetClientSpanName("数据库操作-根据用户ID数组批量获取用户角色")
	newCtx, span := r.trace.AddClientTrace(ctx)
	defer func() { r.trace.TelemetrySpanEnd(span, err) }()

	if len(userIDs) == 0 {
		return
	}

	userSet, userArgIDs := GetFindInSetSQL(userIDs)
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "select f_user_id, f_role_id from %s.t_user_role_relation where f_user_id in ( "
	strSQL += userSet
	strSQL += " )"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.dbTrace.QueryContext(newCtx, strSQL, userArgIDs...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	var roleID string
	var userID string
	out = make(map[string]map[interfaces.Role]bool)
	for rows.Next() {
		if err := rows.Scan(&userID, &roleID); err != nil {
			r.logger.Errorln(err, strSQL)
			return nil, err
		}

		// 数据汇总
		roleInfos, ok := out[userID]
		if !ok {
			temp := make(map[interfaces.Role]bool)
			temp[interfaces.Role(roleID)] = true
			out[userID] = temp
		} else {
			roleInfos[interfaces.Role(roleID)] = true
		}
	}

	return out, nil
}

// GetUserIDsByRoleIDs 根据roleid获取用户信息
func (r *role) GetUserIDsByRoleIDs(ctx context.Context, roles []interfaces.Role) (out map[interfaces.Role][]string, err error) {
	// trace
	r.trace.SetClientSpanName("数据库操作-根据roleid获取用户信息")
	newCtx, span := r.trace.AddClientTrace(ctx)
	defer func() { r.trace.TelemetrySpanEnd(span, err) }()

	if len(roles) == 0 {
		return
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, string(role))
	}

	roleSet, roleArgIDs := GetFindInSetSQL(roleIDs)
	dbName := common.GetDBName("sharemgnt_db")
	strSQL := "select f_user_id, f_role_id from %s.t_user_role_relation where f_role_id in ( "
	strSQL += roleSet
	strSQL += " )"
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := r.dbTrace.QueryContext(newCtx, strSQL, roleArgIDs...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	out = make(map[interfaces.Role][]string)
	for rows.Next() {
		var userID, roleID string
		if err := rows.Scan(&userID, &roleID); err != nil {
			r.logger.Errorln(err, strSQL)
			return nil, err
		}

		// 数据汇总
		if _, ok := out[interfaces.Role(roleID)]; !ok {
			out[interfaces.Role(roleID)] = make([]string, 0)
		}

		out[interfaces.Role(roleID)] = append(out[interfaces.Role(roleID)], userID)
	}

	return out, nil
}
