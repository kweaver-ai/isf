package common

import (
	"time"
)

// Now 业务系统时间
func Now() time.Time {
	return time.Now().Add(time.Second * time.Duration(SvcConfig.BusinessTimeOffset))
}

// GetCurrentMicrosecondTimestamp 返回当前时间的微秒级别的时间戳
/*
	`time.Now()`：这个函数调用返回当前的本地时间
	`.UnixNano()`：这是 `time.Time` 对象的一个方法，它返回自 Unix 纪元（1970年1月1日 00:00:00 UTC）以来的纳秒数
	`/ 1000`：将纳秒数除以1000，将其转换为微秒数
*/
func GetCurrentMicrosecondTimestamp() int64 {
	return time.Now().UnixNano() / 1000 //nolint:mnd
}
