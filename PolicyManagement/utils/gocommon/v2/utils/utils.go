package utils

import (
	"os"
	"strings"
)

// GetEnv 获取环境变量的重新实现, 可以指定默认值
func GetEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

// ContactStr 连接字符串
func ContactStr(strSlice ...string) string {
	var result strings.Builder
	for _, v := range strSlice {
		result.WriteString(v)
	}
	return result.String()
}
