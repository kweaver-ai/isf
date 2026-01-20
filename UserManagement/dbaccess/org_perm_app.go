// Package dbaccess org_perm_app Anyshare 数据访问层 - 应用账户组织架构管理权限管理
package dbaccess

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type orgPermApp struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	opaOnce sync.Once
	opa     *orgPermApp
)

// NewOrgPermApp 创建数据库操作对象
func NewOrgPermApp() *orgPermApp {
	opaOnce.Do(func() {
		opa = &orgPermApp{
			db:      dbPool,
			dbTrace: dbTracePool,
			logger:  common.NewLogger(),
			trace:   common.SvcARTrace,
		}
	})

	return opa
}

// GetAppPermByID2 根据应用账户ID获取应用账户对组织架构的权限信息，支持trace
func (o *orgPermApp) GetAppPermByID2(ctx context.Context, id string) (out map[interfaces.OrgType]interfaces.AppOrgPerm, err error) {
	// trace
	o.trace.SetClientSpanName("数据库操作-根据应用账户ID获取应用账户对组织架构的权限信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	strSQL := "select f_app_name, f_org_type, f_perm_value, f_end_time from %s.t_org_perm_app where f_app_id = ? and (f_end_time = -1 or f_end_time > ?) "
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := o.dbTrace.QueryContext(newCtx, strSQL, id, currentTime)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		o.logger.Errorln(sqlErr, strSQL, id, currentTime)
		return nil, sqlErr
	}

	out = make(map[interfaces.OrgType]interfaces.AppOrgPerm)
	for rows.Next() {
		var tmpInfo interfaces.AppOrgPerm
		if scanErr := rows.Scan(&tmpInfo.Name, &tmpInfo.Object, &tmpInfo.Value, &tmpInfo.EndTime); scanErr != nil {
			o.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		tmpInfo.Subject = id
		out[tmpInfo.Object] = tmpInfo
	}

	return
}

// GetAppPermByID 根据应用账户ID获取应用账户对组织架构的权限信息
func (o *orgPermApp) GetAppPermByID(id string) (out map[interfaces.OrgType]interfaces.AppOrgPerm, err error) {
	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	strSQL := "select f_app_name, f_org_type, f_perm_value, f_end_time from %s.t_org_perm_app where f_app_id = ? and (f_end_time = -1 or f_end_time > ?) "
	strSQL = fmt.Sprintf(strSQL, dbName)

	rows, sqlErr := o.db.Query(strSQL, id, currentTime)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				o.logger.Errorln(rowsErr)
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				o.logger.Errorln(closeErr)
			}
		}
	}()

	if sqlErr != nil {
		o.logger.Errorln(sqlErr, strSQL, id, currentTime)
		return nil, sqlErr
	}

	out = make(map[interfaces.OrgType]interfaces.AppOrgPerm)
	for rows.Next() {
		var tmpInfo interfaces.AppOrgPerm
		if scanErr := rows.Scan(&tmpInfo.Name, &tmpInfo.Object, &tmpInfo.Value, &tmpInfo.EndTime); scanErr != nil {
			o.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		tmpInfo.Subject = id
		out[tmpInfo.Object] = tmpInfo
	}

	return
}

// UpdateAppName 更新应用账户组织架构权限表 名称
func (o *orgPermApp) UpdateAppName(info *interfaces.AppInfo) (err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "update %s.t_org_perm_app set f_app_name = ? where f_app_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	_, err = o.db.Exec(strSQL, info.Name, info.ID)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	return
}

// UpdateAppOrgPermInfo 更新应用账户权限信息
func (o *orgPermApp) UpdateAppOrgPerm(info interfaces.AppOrgPerm, tx *sql.Tx) (err error) {
	// 处理数据
	currentTime := time.Now().UnixNano()

	// 整理sql
	dbName := common.GetDBName("user_management")
	strSQL := "update %s.t_org_perm_app set f_perm_value = ?, f_end_time = ?, f_modify_time = ? where f_app_id = ? and f_org_type = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	_, err = tx.Exec(strSQL, info.Value, info.EndTime, currentTime, info.Subject, info.Object)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}

// AddAppOrgPermInfo 增加应用账户权限信息
func (o *orgPermApp) AddAppOrgPerm(info interfaces.AppOrgPerm, tx *sql.Tx) (err error) {
	// 处理数据
	currentTime := time.Now().UnixNano()

	// 整理sql
	dbName := common.GetDBName("user_management")
	strSQL := "insert into %s.t_org_perm_app (f_app_id, f_app_name, f_org_type, f_perm_value, f_end_time, f_modify_time, f_create_time) values (?, ?, ?, ?, ?, ?, ?) "
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	_, err = tx.Exec(strSQL, info.Subject, info.Name, info.Object, info.Value, info.EndTime, currentTime, currentTime)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}

// DeleteAppOrgPermInfo 删除应用账户权限信息
func (o *orgPermApp) DeleteAppOrgPerm(id string, types []interfaces.OrgType) (err error) {
	if len(types) == 0 {
		return nil
	}

	// 整理sql
	set := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range types {
		set = append(set, "?")
		args = append(args, v)
	}
	sql1 := strings.Join(set, ",")

	dbName := common.GetDBName("user_management")
	strSQL := "delete from %s.t_org_perm_app where f_org_type in ( "
	strSQL += sql1
	strSQL += " ) and f_app_id = ? "
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	args = append(args, id)
	_, err = o.db.Exec(strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}
