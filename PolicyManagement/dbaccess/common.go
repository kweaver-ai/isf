package dbaccess

import "strings"

// GetFindInSetSQL 获取集合查询部分的SQL
func GetFindInSetSQL(value []string) (sql string, args []interface{}) {
	set := make([]string, 0)
	for _, v := range value {
		set = append(set, "?")
		args = append(args, v)
	}
	sql = strings.Join(set, ",")

	return
}
