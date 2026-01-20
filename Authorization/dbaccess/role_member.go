// Package dbaccess group member Anyshare 数据访问层 - 角色成员数据库操作
package dbaccess

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"Authorization/common"
	"Authorization/interfaces"
)

type roleMember struct {
	db     *sqlx.DB
	logger common.Logger
}

var (
	rmOnce sync.Once
	rmDB   *roleMember
)

// NewRoleMember 创建数据库操作对象--和角色成员相关
func NewRoleMember() *roleMember {
	rmOnce.Do(func() {
		rmDB = &roleMember{
			db:     dbPool,
			logger: common.NewLogger(),
		}
	})

	return rmDB
}

// DeleteRoleMemberByID 按照角色id删除角色成员
func (r *roleMember) DeleteRoleMemberByID(ctx context.Context, id string) (err error) {
	dbName := common.GetDBName(databaseName)
	sqlStr := "delete from %s.t_role_member where f_role_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := r.db.Exec(sqlStr, id); err != nil {
		return err
	}
	return nil
}

// AddRoleMembers 批量添加角色成员
func (r *roleMember) AddRoleMembers(ctx context.Context, id string, infos []interfaces.RoleMemberInfo) (err error) {
	if len(infos) == 0 {
		return nil
	}
	currentTime := common.GetCurrentMicrosecondTimestamp()
	var valuesStr []string
	var inserts []any
	// 批量插入
	for i := range infos {
		valuesStr = append(valuesStr, "(?, ?, ?, ?, ?, ?)")
		inserts = append(inserts, id, infos[i].ID, infos[i].MemberType, infos[i].Name, currentTime, currentTime)
	}
	valueStr := strings.Join(valuesStr, ",")

	strSQL := "insert into " + common.GetDBName(databaseName) +
		".t_role_member(f_role_id, f_member_id, f_member_type, f_member_name, f_created_time, f_modify_time) values " + valueStr

	_, err = r.db.Exec(strSQL, inserts...)
	if err != nil {
		r.logger.Errorf("sql: %s, err: %v", strSQL, err)
		return err
	}
	return nil
}

// DeleteRoleMembers 批量删除角色成员
func (r *roleMember) DeleteRoleMembers(ctx context.Context, id string, membersIDs []string) (err error) {
	if len(membersIDs) == 0 {
		return nil
	}
	membersSet, membersGroup := getFindInSetSQL(membersIDs)
	dbName := common.GetDBName(databaseName)
	sqlStr := "delete from %s.t_role_member where f_role_id = ? and f_member_id in (" + membersSet + ")"
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	var args []any
	args = append(args, id)
	args = append(args, membersGroup...)

	if _, err := r.db.Exec(sqlStr, args...); err != nil {
		r.logger.Errorf("DeleteRoleMembers sql: %s, err: %v", sqlStr, err)
		return err
	}
	return nil
}

// GetRoleMembersNum 列举用户组成员数量
func (r *roleMember) GetRoleMembersNum(ctx context.Context, id string, info interfaces.RoleMemberSearchInfo) (num int, err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	args = append(args, id)
	fliterStr := ""
	if len(info.MemberTypes) > 0 {
		memberTypes := make([]string, 0, len(info.MemberTypes))
		for _, v := range info.MemberTypes {
			memberTypes = append(memberTypes, "?")
			args = append(args, v)
		}
		memberTypesStr := strings.Join(memberTypes, ",")
		fliterStr = "and f_member_type in (" + memberTypesStr + ")"
	}
	if info.Keyword != "" {
		fliterStr += "and f_member_name like ?"
		args = append(args, "%"+info.Keyword+"%")
	}
	strSQL := `select count(f_member_id) from %s.t_role_member where f_role_id = ? ` + fliterStr
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := r.db.Query(strSQL, args...)
	defer func() {
		if rows != nil {
			if rowsErr := rows.Err(); rowsErr != nil {
				r.logger.Errorln(rowsErr)
				return
			}

			// 1、判断是否为空再关闭，2、如果不关闭而数据行并没有被scan的话，连接一直会被占用直到超时断开
			if closeErr := rows.Close(); closeErr != nil {
				r.logger.Errorln(closeErr)
				return
			}
		}
	}()

	if sqlErr != nil {
		r.logger.Errorln(sqlErr, strSQL)
		return num, sqlErr
	}

	for rows.Next() {
		if scanErr := rows.Scan(&num); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return num, scanErr
		}
	}

	return num, nil
}

// GetPaginationByRoleID 列举角色成员
func (r *roleMember) GetPaginationByRoleID(ctx context.Context, id string, info interfaces.RoleMemberSearchInfo) (outInfo []interfaces.RoleMemberInfo, err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	args = append(args, id)
	fliterStr := ""
	if len(info.MemberTypes) > 0 {
		memberTypes := make([]string, 0, len(info.MemberTypes))
		for _, v := range info.MemberTypes {
			memberTypes = append(memberTypes, "?")
			args = append(args, v)
		}
		memberTypesStr := strings.Join(memberTypes, ",")
		fliterStr = "and f_member_type in (" + memberTypesStr + ")"
	}
	if info.Keyword != "" {
		fliterStr += "and f_member_name like ?"
		args = append(args, "%"+info.Keyword+"%")
	}
	strSQL := `select f_member_id, f_member_type, f_member_name from %s.t_role_member where f_role_id = ? ` + fliterStr + ` order by f_modify_time desc, f_primary_id desc LIMIT ?, ?`
	strSQL = fmt.Sprintf(strSQL, dbName)
	args = append(args, info.Offset, info.Limit)

	rows, sqlErr := r.db.Query(strSQL, args...)
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

	outInfo = make([]interfaces.RoleMemberInfo, 0)
	for rows.Next() {
		var tmpInfo interfaces.RoleMemberInfo
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.MemberType, &tmpInfo.Name); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		outInfo = append(outInfo, tmpInfo)
	}

	return outInfo, nil
}

// GetRoleMembersByRoleID 根据角色ID获取角色成员
func (r *roleMember) GetRoleMembersByRoleID(ctx context.Context, id string) (outInfo []interfaces.RoleMemberInfo, err error) {
	var args []any
	dbName := common.GetDBName(databaseName)
	strSQL := `select f_member_id, f_member_type, f_member_name from %s.t_role_member where f_role_id = ? `
	strSQL = fmt.Sprintf(strSQL, dbName)
	args = append(args, id)

	rows, sqlErr := r.db.Query(strSQL, args...)
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

	outInfo = make([]interfaces.RoleMemberInfo, 0)
	for rows.Next() {
		var tmpInfo interfaces.RoleMemberInfo
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.MemberType, &tmpInfo.Name); scanErr != nil {
			r.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		outInfo = append(outInfo, tmpInfo)
	}

	return outInfo, nil
}

// GetRoleByMembers 通过成员获取角色
func (r *roleMember) GetRoleByMembers(ctx context.Context, memberIDs []string) (outInfo []interfaces.RoleInfo, err error) {
	memberSet, memberIDGroup := getFindInSetSQL(memberIDs)
	var paramList []any
	paramList = append(paramList, memberIDGroup...)
	dbName := common.GetDBName(databaseName)
	strSQL := "select f_role_id from %s.t_role_member where f_member_id in (" + memberSet + ")"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, err := r.db.Query(strSQL, paramList...)
	if err != nil {
		r.logger.Errorf("GetRoleByMembers err:%v, strSQL:%s", err, strSQL)
		return nil, err
	}
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

	for rows.Next() {
		var role interfaces.RoleInfo
		err = rows.Scan(&role.ID)
		if err != nil {
			r.logger.Errorf("GetRoleByMembers err:%v, strSQL:%s", err, strSQL)
			return nil, err
		}
		outInfo = append(outInfo, role)
	}
	return
}

// DeleteByMemberIDs 删除成员 根据成员id
func (r *roleMember) DeleteByMemberIDs(memberIDs []string) error {
	if len(memberIDs) == 0 {
		return nil
	}
	IDsSet, IDsGroup := getFindInSetSQL(memberIDs)
	strSQL := "delete from " + common.GetDBName(databaseName) + ".t_role_member where f_member_id in (" + IDsSet + ")"
	_, err := r.db.Exec(strSQL, IDsGroup...)
	if err != nil {
		r.logger.Errorf("DeleteByMemberIDs sql: %s, err: %v", strSQL, err)
		return err
	}
	return nil
}

// DeleteByRoleID 根据角色ID删除成员
func (r *roleMember) DeleteByRoleID(ctx context.Context, roleID string) error {
	strSQL := "delete from " + common.GetDBName(databaseName) + ".t_role_member where f_role_id = ?"
	_, err := r.db.Exec(strSQL, roleID)
	if err != nil {
		r.logger.Errorf("DeleteByRoleID sql: %s, err: %v", strSQL, err)
		return err
	}
	return nil
}

func (r *roleMember) UpdateMemberName(memberID, name string) error {
	strSQL := "update " + common.GetDBName(databaseName) + ".t_role_member set f_member_name = ? where f_member_id = ?"
	_, err := r.db.Exec(strSQL, name, memberID)
	if err != nil {
		r.logger.Errorf("UpdateMemberName sql: %s, err: %v", strSQL, err)
		return err
	}
	return nil
}
