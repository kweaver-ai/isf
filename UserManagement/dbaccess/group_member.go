// Package dbaccess group member Anyshare 数据访问层 - 用户组成员数据库操作
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

type groupMember struct {
	db      *sqlx.DB
	logger  common.Logger
	trace   observable.Tracer
	dbTrace *sqlx.DB
}

var (
	gmOnce sync.Once
	gmDB   *groupMember
)

// NewGroupMember 创建数据库操作对象--和用户组成员相关
func NewGroupMember() *groupMember {
	gmOnce.Do(func() {
		gmDB = &groupMember{
			db:      dbPool,
			dbTrace: dbTracePool,
			logger:  common.NewLogger(),
			trace:   common.SvcARTrace,
		}
	})

	return gmDB
}

// DeleteGroupMemberByID 按照用户组id删除用户组成员
func (g *groupMember) DeleteGroupMemberByID(id string) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_group_member where f_group_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := g.db.Exec(sqlStr, id); err != nil {
		return err
	}
	return nil
}

// AddGroupMember 添加用户组成员
func (g *groupMember) AddGroupMember(id string, info *interfaces.GroupMemberInfo) (err error) {
	currentTime := time.Now().UnixNano()
	dbName := common.GetDBName("user_management")
	sqlStr := "insert into %s.t_group_member(f_group_id, f_member_id, f_member_type, f_added_time)values(?, ?, ?, ?)"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := g.db.Exec(sqlStr, id, info.ID, info.MemberType, currentTime); err != nil {
		return err
	}
	return nil
}

// AddGroupMembers 批量添加用户组成员
func (g *groupMember) AddGroupMembers(ctx context.Context, id string, infos []interfaces.GroupMemberInfo, tx *sql.Tx) (err error) {
	// trace
	g.trace.SetClientSpanName("数据库操作-批量添加用户组成员")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	if len(infos) == 0 {
		return
	}

	// 参数检查
	currentTime := time.Now().UnixNano()
	args := []interface{}{}
	sqlParamsStr := []string{}
	for _, v := range infos {
		args = append(args, id, v.ID, v.MemberType, currentTime)
		sqlParamsStr = append(sqlParamsStr, "(?, ?, ?, ?)")
	}

	dbName := common.GetDBName("user_management")
	sqlStr := "insert into %s.t_group_member(f_group_id, f_member_id, f_member_type, f_added_time) values "
	sqlStr = fmt.Sprintf(sqlStr, dbName)
	sqlStr += strings.Join(sqlParamsStr, ",")
	if _, err := g.db.ExecContext(newCtx, sqlStr, args...); err != nil {
		return err
	}
	return nil
}

// DeleteGroupMember 删除用户组成员
func (g *groupMember) DeleteGroupMember(id string, info *interfaces.GroupMemberInfo) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_group_member where f_group_id = ? and f_member_id = ? and f_member_type = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	if _, err := g.db.Exec(sqlStr, id, info.ID, info.MemberType); err != nil {
		return err
	}
	return nil
}

// DeleteGroupMemberByMemberID 根据成员ID删除用户组成员关系
func (g *groupMember) DeleteGroupMemberByMemberID(id string) (err error) {
	dbName := common.GetDBName("user_management")
	sqlStr := "delete from %s.t_group_member where  f_member_id = ?"
	sqlStr = fmt.Sprintf(sqlStr, dbName)

	_, err = g.db.Exec(sqlStr, id)
	if err != nil {
		g.logger.Errorln(err, sqlStr, id)
	}
	return err
}

// GetGroupMembersByGroupIDs 批量获取用户组成员的id和类型
func (g *groupMember) GetGroupMembersByGroupIDs(groupIDs []string) (outInfos []interfaces.GroupMemberInfo, err error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}

	set, argIDs := GetFindInSetSQL(groupIDs)
	dbName := common.GetDBName("user_management")
	strSQL := "select distinct(f_member_id), f_member_type from %s.t_group_member where f_group_id in ("
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

	var tmpInfo interfaces.GroupMemberInfo
	outInfos = make([]interfaces.GroupMemberInfo, 0)
	for rows.Next() {
		if scanErr := rows.Scan(&tmpInfo.ID, &tmpInfo.MemberType); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		outInfos = append(outInfos, tmpInfo)
	}

	return outInfos, nil
}

// GetGroupMembersByGroupIDs2 批量获取用户组成员的id和类型，支持trace
func (g *groupMember) GetGroupMembersByGroupIDs2(ctx context.Context, groupIDs []string) (outInfos map[string][]interfaces.GroupMemberInfo, err error) {
	// trace
	g.trace.SetClientSpanName("数据库操作-根据组ID获取组成员ID")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	if len(groupIDs) == 0 {
		return nil, nil
	}

	set, argIDs := GetFindInSetSQL(groupIDs)
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id, f_member_id, f_member_type from %s.t_group_member where f_group_id in ("
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
		g.logger.Errorln(sqlErr, strSQL)
		return nil, sqlErr
	}

	outInfos = make(map[string][]interfaces.GroupMemberInfo)
	for rows.Next() {
		var tmpInfo interfaces.GroupMemberInfo
		groupID := ""
		if scanErr := rows.Scan(&groupID, &tmpInfo.ID, &tmpInfo.MemberType); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		if _, ok := outInfos[groupID]; !ok {
			outInfos[groupID] = make([]interfaces.GroupMemberInfo, 0)
		}

		tempInfo := outInfos[groupID]
		tempInfo = append(tempInfo, tmpInfo)
		outInfos[groupID] = tempInfo
	}

	return outInfos, nil
}

// GetGroupMembers 列举用户组成员
func (g *groupMember) GetGroupMembers(ctx context.Context, id string, info interfaces.SearchInfo) (outInfo []interfaces.GroupMemberInfo, err error) {
	// trace
	g.trace.SetClientSpanName("数据库操作-列举用户组成员")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	var args []interface{}

	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL :=
		`select c.f_name , b.f_display_name, a.f_member_id, a.f_member_type, a.f_added_time from( (
		select  f_member_id ,f_member_type, f_added_time from %s.t_group_member where f_group_id = ?) a
		left join
		%s.t_user as b on a.f_member_id = b.f_user_id
		left join
		%s.t_department as c on a.f_member_id = c.f_department_id ) `
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB)

	args = append(args, id)
	if info.HasKeyWord {
		args = append(args, "%"+info.Keyword+"%", "%"+info.Keyword+"%")
		strSQL += "where c.f_name like ? or b.f_display_name like ? "
	}

	if info.Sort == interfaces.Name {
		strSQL += "order by UPPER(IFNULL(c.f_name , b.f_display_name))"
	} else {
		strSQL += "order by a.f_added_time "
	}

	if info.Direction == interfaces.Desc {
		strSQL += strDesc
	} else {
		strSQL += strAsc
	}

	// 添加第二排序规则 防止重复
	strSQL += ",  a.f_member_id asc LIMIT ?, ?"

	args = append(args, info.Offset, info.Limit)

	rows, sqlErr := g.dbTrace.QueryContext(newCtx, strSQL, args...)
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

	var tmpInfo interfaces.GroupMemberInfo
	outInfo = make([]interfaces.GroupMemberInfo, 0)
	var nAddTime int64
	for rows.Next() {
		var strDepartName, strUserName sql.NullString
		if scanErr := rows.Scan(&strDepartName, &strUserName, &tmpInfo.ID, &tmpInfo.MemberType, &nAddTime); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}
		tmpInfo.Name = strDepartName.String + strUserName.String
		outInfo = append(outInfo, tmpInfo)
	}

	return outInfo, nil
}

// GetGroupMembersNum 列举用户组成员数量
func (g *groupMember) GetGroupMembersNum(id string, info interfaces.SearchInfo) (num int, err error) {
	var args []interface{}
	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL :=
		`select count(a.f_member_id) from( (
		select  f_member_id, f_member_type from %s.t_group_member where f_group_id = ?) a
		left join
		%s.t_user as b on a.f_member_id = b.f_user_id
		left join
		%s.t_department as c on a.f_member_id = c.f_department_id )
		`
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB)
	args = append(args, id)

	if info.HasKeyWord || info.NotShowDisabledUser {
		strSQL += "where "
	}

	if info.HasKeyWord {
		args = append(args, "%"+info.Keyword+"%", "%"+info.Keyword+"%")
		strSQL += "c.f_name like ? or b.f_display_name like ? "
	}

	if info.HasKeyWord && info.NotShowDisabledUser {
		strSQL += "and "
	}

	if info.NotShowDisabledUser {
		strSQL += "((b.f_status = 0 and b.f_auto_disable_status = 0) or (a.f_member_type <> 1))"
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

// GetGroupMembersNum2 列举用户组成员数量,支持trace
func (g *groupMember) GetGroupMembersNum2(ctx context.Context, id string, info interfaces.SearchInfo) (num int, err error) {
	// trace
	g.trace.SetClientSpanName("数据库操作-列举用户组成员数量")
	newCtx, span := g.trace.AddClientTrace(ctx)
	defer func() { g.trace.TelemetrySpanEnd(span, err) }()

	var args []interface{}
	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL :=
		`select count(a.f_member_id) from( (
		select  f_member_id, f_member_type from %s.t_group_member where f_group_id = ?) a
		left join
		%s.t_user as b on a.f_member_id = b.f_user_id
		left join
		%s.t_department as c on a.f_member_id = c.f_department_id )
		`
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB)

	args = append(args, id)

	if info.HasKeyWord || info.NotShowDisabledUser {
		strSQL += "where "
	}

	if info.HasKeyWord {
		args = append(args, "%"+info.Keyword+"%", "%"+info.Keyword+"%")
		strSQL += "c.f_name like ? or b.f_display_name like ? "
	}

	if info.HasKeyWord && info.NotShowDisabledUser {
		strSQL += "and "
	}

	if info.NotShowDisabledUser {
		strSQL += "((b.f_status = 0 and b.f_auto_disable_status = 0) or (a.f_member_type <> 1))"
	}

	rows, sqlErr := g.dbTrace.QueryContext(newCtx, strSQL, args...)
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

// CheckGroupMembersExist 用户组成员是否存在
func (g *groupMember) CheckGroupMembersExist(id string, info *interfaces.GroupMemberInfo) (ret bool, err error) {
	dbName := common.GetDBName("user_management")
	strSQL := "select f_group_id from %s.t_group_member where f_group_id = ? and f_member_id = ? and f_member_type = ?"
	strSQL = fmt.Sprintf(strSQL, dbName)
	rows, sqlErr := g.db.Query(strSQL, id, info.ID, info.MemberType)
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
		return false, sqlErr
	}

	for rows.Next() {
		ret = true
	}

	return ret, nil
}

// GetMembersBelongGroupIDs 获取包含成员的用户组ID或用户组信息
func (g *groupMember) GetMembersBelongGroupIDs(memberIDs []string) (groupIDs []string, groups []interfaces.GroupInfo, err error) {
	groupIDs = make([]string, 0)
	groups = make([]interfaces.GroupInfo, 0)
	if len(memberIDs) == 0 {
		return groupIDs, groups, err
	}

	set, argIDs := GetFindInSetSQL(memberIDs)
	dbName := common.GetDBName("user_management")
	strSQL := `select distinct a.f_group_id, b.f_group_name
			from %s.t_group_member as a
			left join %s.t_group as b on a.f_group_id = b.f_group_id
			where a.f_member_id in (`
	strSQL += set
	strSQL += ")"
	strSQL = fmt.Sprintf(strSQL, dbName, dbName)

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
		return groupIDs, groups, sqlErr
	}

	var groupInfo interfaces.GroupInfo
	for rows.Next() {
		if scanErr := rows.Scan(&groupInfo.ID, &groupInfo.Name); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, nil, scanErr
		}

		groups = append(groups, groupInfo)
		groupIDs = append(groupIDs, groupInfo.ID)
	}

	return groupIDs, groups, nil
}

// SearchMembersByKeyword 用户组成员关键字搜索
func (g *groupMember) SearchMembersByKeyword(keyword string, start, limit int) (out []interfaces.MemberInfo, err error) {
	var args []interface{}
	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL := `select a.f_member_id, a.f_member_type , d.f_group_name, c.f_name , b.f_display_name from
			(select  f_member_id ,f_member_type, f_group_id from %s.t_group_member ) a
			left join
			%s.t_user as b on a.f_member_id = b.f_user_id
			left join
			%s.t_department as c on a.f_member_id = c.f_department_id
			left join
			%s.t_group as d on a.f_group_id = d.f_group_id
			where c.f_name like ? or  b.f_display_name like ? and b.f_status = 0 and b.f_auto_disable_status = 0
			order by case when  c.f_name = ? or b.f_display_name = ? then 0 when c.f_name like ? or b.f_display_name like ? then 1 else 2 end,
			UPPER(IFNULL(c.f_name , b.f_display_name))
			limit ?, ? `
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB, userMgntDB)

	args = append(args, "%"+keyword+"%", "%"+keyword+"%", keyword, keyword, "%"+keyword+"%", "%"+keyword+"%", start, limit)
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

	tmpData := make(map[string]interfaces.MemberInfo)
	tmpExsitData := make(map[string]bool)
	nameList := make([]string, 0)
	for rows.Next() {
		var temp interfaces.MemberInfo
		var groupName string
		var strUserName, strDepartName sql.NullString
		if scanErr := rows.Scan(&temp.ID, &temp.NType, &groupName, &strDepartName, &strUserName); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		temp.Name = strUserName.String + strDepartName.String
		if tmpExsitData[temp.ID] {
			tmp := tmpData[temp.ID]
			tmp.GroupNames = append(tmp.GroupNames, groupName)
			tmpData[temp.ID] = tmp
		} else {
			temp.GroupNames = make([]string, 0)
			temp.GroupNames = append(temp.GroupNames, groupName)
			tmpExsitData[temp.ID] = true
			tmpData[temp.ID] = temp

			nameList = append(nameList, temp.ID)
		}
	}

	out = make([]interfaces.MemberInfo, 0)
	for _, v := range nameList {
		out = append(out, tmpData[v])
	}

	return out, nil
}

// SearchMemberNumByKeyword 用户组成员关键字搜索符合条件的用户组总数目
func (g *groupMember) SearchMemberNumByKeyword(keyword string) (num int, err error) {
	var args []interface{}
	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL := `select count(distinct(a.f_member_id)) from
			(select  f_member_id ,f_member_type, f_group_id from %s.t_group_member ) a
			left join
			%s.t_user as b on a.f_member_id = b.f_user_id
			left join
			%s.t_department as c on a.f_member_id = c.f_department_id
			left join
			%s.t_group as d on a.f_group_id = d.f_group_id
			where c.f_name like ? or  b.f_display_name like ? and b.f_status = 0 and b.f_auto_disable_status = 0 `
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB, userMgntDB)

	args = append(args, "%"+keyword+"%", "%"+keyword+"%")
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

// GetMemberOnClient 客户端列举组成员
func (g *groupMember) GetMemberOnClient(id string, offset, limit int) (info []interfaces.MemberSimpleInfo, err error) {
	var args []interface{}
	userMgntDB := common.GetDBName("user_management")
	shareMgntDB := common.GetDBName("sharemgnt_db")
	strSQL := `select a.f_member_id, a.f_member_type , c.f_name , b.f_display_name, IFNULL(b.f_priority ,c.f_priority ) as f_new_priority from( (
			select  f_member_id ,f_member_type from %s.t_group_member where f_group_id = ?) a
			left join
			%s.t_user as b on a.f_member_id = b.f_user_id
			left join
			%s.t_department as c on a.f_member_id = c.f_department_id )
			where a.f_member_type != 1 or b.f_status = 0 and b.f_auto_disable_status = 0
			order by f_new_priority asc ,UPPER(IFNULL(c.f_name , b.f_display_name)) asc
			limit ?, ? `
	strSQL = fmt.Sprintf(strSQL, userMgntDB, shareMgntDB, shareMgntDB)

	args = append(args, id, offset, limit)
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

	info = make([]interfaces.MemberSimpleInfo, 0)
	for rows.Next() {
		var temp interfaces.MemberSimpleInfo
		var nPriority int
		var strDepartName, strUserName sql.NullString
		if scanErr := rows.Scan(&temp.ID, &temp.NType, &strDepartName, &strUserName, &nPriority); scanErr != nil {
			g.logger.Errorln(scanErr, strSQL)
			return nil, scanErr
		}

		temp.Name = strDepartName.String + strUserName.String
		info = append(info, temp)
	}
	return
}
