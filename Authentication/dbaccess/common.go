package dbaccess

import "strings"

// GetFindInSetSQL 获取集合查询部分的SQL
func GetFindInSetSQL(ids []string) (sql string, argIDs []interface{}) {
	set := make([]string, 0)
	for _, id := range ids {
		set = append(set, "?")
		argIDs = append(argIDs, id)
	}
	sql = strings.Join(set, ",")

	return
}

// SplitArray 拆分数组，保证in值列表限制在500以内
func SplitArray(arr []string) [][]string {
	length := 500
	total := len(arr)
	count := total / length
	if total%length != 0 {
		count++
	}

	resArray := make([][]string, 0)
	start := 0
	end := 0
	for i := 0; i < count; i++ {
		end = (i + 1) * length
		if i != (count - 1) {
			resArray = append(resArray, arr[start:end])
		} else {
			resArray = append(resArray, arr[start:])
		}
		start = end
	}

	return resArray
}
