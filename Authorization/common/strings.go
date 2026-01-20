package common

import (
	"fmt"
	"strings"
)

// Joinable 接口用于 Join 函数，表示可连接的类型
type Joinable interface {
	~string
}

// Join 函数对给定的 items 切片，使用 extract 函数从每个元素中提取值，
// 将这些值转换为字符串后，用指定的 separator 分隔符连接成一个字符串，并返回结果。
//
// 类型参数：
//   - T：items 切片中元素的类型。
//   - V：提取值的类型，必须实现 Joinable 接口。
//
// 参数：
//   - items：类型为 []T，要处理的元素切片。
//   - extract：func(T) V，函数类型，从每个元素中提取值。
//   - separator：string，分隔符，用于连接提取的字符串。
//
// 返回值：
//   - string：提取值转换为字符串并用分隔符连接后的结果字符串。
func Join[T any, V Joinable](items []T, extract func(T) V, separator string) string {
	var fields []string
	for i := range items {
		fields = append(fields, fmt.Sprint(extract(items[i])))
	}
	return strings.Join(fields, separator)
}
