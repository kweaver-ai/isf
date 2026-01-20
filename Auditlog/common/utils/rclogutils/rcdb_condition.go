package rclogutils

import (
	"fmt"
	"strconv"
	"strings"

	"AuditLog/common"
	"AuditLog/common/constants/rclogconsts"
	"AuditLog/models/rcvo"
)

// SetDefaultSortField 设置默认排序字段
func SetDefaultSortField(orderBy *rcvo.OrderFields, ids []string) {
	// 当没有传入排序字段和ID列表时，按照创建时间倒序排序
	if len(*orderBy) == 0 && len(ids) == 0 {
		*orderBy = rcvo.OrderFields{
			{
				Field:     rclogconsts.Date,
				Direction: "desc",
			},
		}
	}
}

// AddWhereToSql 添加查询条件
func AddWhereToSql(sqlStr *string, where map[string]any) (err error) {
	// 没有查询条件时
	if len(where) == 0 {
		*sqlStr = strings.Replace(*sqlStr, "[where]", "", 1)
		return
	}
	// 存在查询条件时
	var conditions []string

	for k, v := range where {
		if k == rclogconsts.Date {
			left := v.([]interface{})[0].(float64)
			right := v.([]interface{})[1].(float64)
			conditions = append(conditions, "(f_"+k+" BETWEEN "+fmt.Sprintf("%d", int(left))+" AND "+fmt.Sprintf("%d", int(right))+")")

			continue
		}

		if k == rclogconsts.LogID {
			if v.(int) > 0 {
				conditions = append(conditions, "f_"+k+" <= "+fmt.Sprintf("%d", v))
			}

			continue
		}

		// 精确匹配
		if common.InArray(k, []string{rclogconsts.LogLevel, rclogconsts.OpType, rclogconsts.LogID}) {
			var value int
			switch vt := v.(type) {
			case int:
				value = vt
			case string:
				value, err = strconv.Atoi(vt)
			}

			if err != nil {
				return
			}

			conditions = append(conditions, "f_"+k+" = "+fmt.Sprintf("%d", value))

			continue
		} else if common.InArray(k, []string{rclogconsts.UserName, rclogconsts.Mac, rclogconsts.IP, rclogconsts.UserPaths}) {
			conditions = append(conditions, "f_"+k+" = "+fmt.Sprintf("'%v'", v))
			continue
		} else {
			conditions = append(conditions, "f_"+k+" LIKE "+fmt.Sprintf("'%%%v%%'", v))
			continue
		}
	}
	// 拼接查询条件
	match := strings.Join(conditions, " AND ")
	if strings.Contains(*sqlStr, "WHERE") {
		// 存在其它查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[where]", "AND "+match, 1)
	} else {
		// 不存在其它查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[where]", "WHERE "+match, 1)
	}

	return
}

// AddWhereToSql2 添加查询条件,支持左闭右开
func AddWhereToSql2(sqlStr *string, where map[string]any) (err error) {
	// 没有查询条件时
	if len(where) == 0 {
		*sqlStr = strings.Replace(*sqlStr, "[where]", "", 1)
		return
	}
	// 存在查询条件时
	var conditions []string

	for k, v := range where {
		if k == rclogconsts.Date {
			left := v.([]interface{})[0].(float64)
			right := v.([]interface{})[1].(float64)
			conditions = append(conditions, "(f_"+k+" >= "+fmt.Sprintf("%d", int(left))+" AND "+" f_"+k+" < "+fmt.Sprintf("%d", int(right))+")")

			continue
		}

		if k == rclogconsts.LogID {
			if v.(int) > 0 {
				conditions = append(conditions, "f_"+k+" <= "+fmt.Sprintf("%d", v))
			}

			continue
		}

		// 精确匹配
		if common.InArray(k, []string{rclogconsts.LogLevel, rclogconsts.OpType, rclogconsts.LogID}) {
			var value int
			switch vt := v.(type) {
			case int:
				value = vt
			case string:
				value, err = strconv.Atoi(vt)
			}

			if err != nil {
				return
			}

			conditions = append(conditions, "f_"+k+" = "+fmt.Sprintf("%d", value))

			continue
		} else if common.InArray(k, []string{rclogconsts.UserName, rclogconsts.Mac, rclogconsts.IP, rclogconsts.UserPaths}) {
			conditions = append(conditions, "f_"+k+" = "+fmt.Sprintf("'%v'", v))
			continue
		} else {
			conditions = append(conditions, "f_"+k+" LIKE "+fmt.Sprintf("'%%%v%%'", v))
			continue
		}
	}
	// 拼接查询条件
	match := strings.Join(conditions, " AND ")
	if strings.Contains(*sqlStr, "WHERE") {
		// 存在其它查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[where]", "AND "+match, 1)
	} else {
		// 不存在其它查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[where]", "WHERE "+match, 1)
	}

	return
}

// AddInToSql 添加枚举条件
func AddInToSql(sqlStr *string, enums map[string][]string) {
	// 没有枚举条件时
	if len(enums) == 0 {
		*sqlStr = strings.Replace(*sqlStr, "[in]", "", 1)
		return
	}
	// 存在枚举条件时
	matches := make([]string, 0)

	for field, values := range enums {
		if len(values) > 0 {
			match := "f_" + field + " IN ('" + strings.Join(values, "','") + "')"
			matches = append(matches, match)
		}
	}

	condition := strings.Join(matches, " AND ")
	if strings.Contains(*sqlStr, "WHERE") {
		// 存在其他查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[in]", "AND "+condition, 1)
	} else {
		// 不存在其他查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[in]", "WHERE "+condition, 1)
	}
}

// AddNotInToSql 添加排除枚举条件
func AddNotInToSql(sqlStr *string, enums map[string][]string) {
	// 没有枚举条件时
	if len(enums) == 0 {
		*sqlStr = strings.Replace(*sqlStr, "[not_in]", "", 1)
		return
	}
	// 存在枚举条件时
	matches := make([]string, 0)

	for field, values := range enums {
		if len(values) > 0 {
			match := "f_" + field + " NOT IN ('" + strings.Join(values, "','") + "')"
			matches = append(matches, match)
		}
	}

	condition := strings.Join(matches, " AND ")
	if strings.Contains(*sqlStr, "WHERE") {
		// 存在其他查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[not_in]", "AND "+condition, 1)
	} else {
		// 不存在其他查询条件时
		*sqlStr = strings.Replace(*sqlStr, "[not_in]", "WHERE "+condition, 1)
	}
}

// AddOrderToSql 添加排序条件
func AddOrderToSql(sqlStr *string, orderBy rcvo.OrderFields) {
	// 没有排序字段时
	if len(orderBy) == 0 {
		*sqlStr = strings.Replace(*sqlStr, "[order]", "", 1)
		return
	}
	// 存在排序字段时
	var orders []string
	for _, v := range orderBy {
		orders = append(orders, "f_"+v.Field+" "+strings.ToUpper(v.Direction))
	}

	*sqlStr = strings.Replace(*sqlStr, "[order]", "ORDER BY "+strings.Join(orders, ","), 1)
}

// BuildActiveCondition 构建查询活跃日志条件
func BuildActiveCondition(condition map[string]any, orderBy rcvo.OrderFields, ids []string, inUserIDs []string, exUserIDs []string) (sqlStr string, err error) {
	sqlStr = `[where] [in] [not_in] [order]`

	SetDefaultSortField(&orderBy, ids)

	if err = AddWhereToSql(&sqlStr, condition); err != nil {
		return
	}

	inenums := make(map[string][]string)
	if len(ids) > 0 {
		inenums["log_id"] = ids
	}

	if len(inUserIDs) > 0 {
		inenums["user_id"] = inUserIDs
	}

	AddInToSql(&sqlStr, inenums)

	exenums := make(map[string][]string)
	if len(exUserIDs) > 0 {
		exenums["user_id"] = exUserIDs
	}

	AddNotInToSql(&sqlStr, exenums)

	AddOrderToSql(&sqlStr, orderBy)

	return
}

// BuildActiveCondition2 构建查询活跃日志条件,支持左闭右开
func BuildActiveCondition2(condition map[string]any, orderBy rcvo.OrderFields, ids []string, inUserIDs []string, exUserIDs []string) (sqlStr string, err error) {
	sqlStr = `[where] [in] [not_in] [order]`

	SetDefaultSortField(&orderBy, ids)

	// 支持左闭右开
	if err = AddWhereToSql2(&sqlStr, condition); err != nil {
		return
	}

	inenums := make(map[string][]string)
	if len(ids) > 0 {
		inenums["log_id"] = ids
	}

	if len(inUserIDs) > 0 {
		inenums["user_id"] = inUserIDs
	}

	AddInToSql(&sqlStr, inenums)

	exenums := make(map[string][]string)
	if len(exUserIDs) > 0 {
		exenums["user_id"] = exUserIDs
	}

	AddNotInToSql(&sqlStr, exenums)

	AddOrderToSql(&sqlStr, orderBy)

	return
}

// BuildHistoryCondition 构建查询历史日志条件
func BuildHistoryCondition(category string, condition map[string]any, orderBy rcvo.OrderFields, ids []string) (sqlStr string, err error) {
	sqlStr = `[where] [in] [order]`

	SetDefaultSortField(&orderBy, ids)

	if condition == nil {
		condition = make(map[string]any)
	}

	condition["type"] = common.LogTypeMap[category]

	if err = AddWhereToSql(&sqlStr, condition); err != nil {
		return
	}

	inenums := make(map[string][]string)
	if len(ids) > 0 {
		inenums["id"] = ids
	}

	AddInToSql(&sqlStr, inenums)

	AddOrderToSql(&sqlStr, orderBy)

	return
}
