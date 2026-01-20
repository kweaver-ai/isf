// Package dbaccess org_perm Anyshare 数据访问层 - 组织架构管理权限管理
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

type orgPerm struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	opOnce sync.Once
	op     *orgPerm
)

// NewOrgPerm 创建数据库操作对象
func NewOrgPerm() *orgPerm {
	opOnce.Do(func() {
		op = &orgPerm{
			db:      dbPool,
			dbTrace: dbTracePool,
			logger:  common.NewLogger(),
			trace:   common.SvcARTrace,
		}
	})

	return op
}

// UpdateName 更新账户组织架构权限表 名称
func (o *orgPerm) UpdateName(id, newName string) (err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "update %s.t_org_perm set f_name = ? where f_id = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	_, err = o.db.Exec(strSQL, newName, id)
	if err != nil {
		o.logger.Errorln(err)
		return
	}
	return
}

// UpdateOrgPerm 更新账户权限信息
func (o *orgPerm) UpdateOrgPerm(ctx context.Context, info interfaces.OrgPerm, tx *sql.Tx) (err error) {
	// trace
	o.trace.SetClientSpanName("数据库操作-更新账户权限信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	// 处理数据
	currentTime := time.Now().UnixNano()

	// 整理sql
	dbName := common.GetDBName("user_management")
	strSQL := "update %s.t_org_perm set f_perm_value = ?, f_end_time = ?, f_modify_time = ? where f_id = ? and f_org_type = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	_, err = tx.ExecContext(newCtx, strSQL, info.Value, info.EndTime, currentTime, info.SubjectID, info.Object)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}

// AddOrgPerm 增加账户权限信息
func (o *orgPerm) AddOrgPerm(ctx context.Context, info interfaces.OrgPerm, tx *sql.Tx) (err error) {
	// trace
	o.trace.SetClientSpanName("数据库操作-增加账户权限信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	// 处理数据
	currentTime := time.Now().UnixNano()

	// 整理sql
	dbName := common.GetDBName("user_management")
	strSQL := "insert into %s.t_org_perm (f_id, f_name, f_type, f_org_type, f_perm_value, f_end_time, f_modify_time, f_create_time) values (?, ?, ?, ?, ?, ?, ?, ?) "
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	_, err = tx.ExecContext(newCtx, strSQL, info.SubjectID, info.Name, info.SubjectType, info.Object, info.Value, info.EndTime, currentTime, currentTime)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}

// DeleteOrgPermInfo 删除账户权限信息
func (o *orgPerm) DeleteOrgPerm(ctx context.Context, id string, types []interfaces.OrgType) (err error) {
	// trace
	o.trace.SetClientSpanName("数据库操作-删除账户权限信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

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
	strSQL := "delete from %s.t_org_perm where f_org_type in ( "
	strSQL += sql1
	strSQL += " ) and f_id = ? "
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	args = append(args, id)
	_, err = o.db.ExecContext(newCtx, strSQL, args...)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}

// GetPermByID 根据账户ID获取账户对组织架构的权限信息，支持trace
func (o *orgPerm) GetPermByID(ctx context.Context, id string) (out map[interfaces.OrgType]interfaces.OrgPerm, err error) {
	// trace
	o.trace.SetClientSpanName("数据库操作-根据账户ID获取账户对组织架构的权限信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	strSQL := "select f_name, f_type, f_org_type, f_perm_value, f_end_time from %s.t_org_perm where f_id = ? and (f_end_time = -1 or f_end_time > ?) "
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

	out = make(map[interfaces.OrgType]interfaces.OrgPerm)
	for rows.Next() {
		var tmpInfo interfaces.OrgPerm
		if scanErr := rows.Scan(&tmpInfo.Name, &tmpInfo.SubjectType, &tmpInfo.Object, &tmpInfo.Value, &tmpInfo.EndTime); scanErr != nil {
			o.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		tmpInfo.SubjectID = id
		out[tmpInfo.Object] = tmpInfo
	}

	return
}

// DeleteOrgPermByID 删除账户权限信息
func (o *orgPerm) DeleteOrgPermByID(id string) (err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "delete from %s.t_org_perm where f_id = ? "
	strSQL = fmt.Sprintf(strSQL, dbName)

	// 执行sql
	_, err = o.db.Exec(strSQL, id)
	if err != nil {
		o.logger.Errorln(err)
	}
	return
}
