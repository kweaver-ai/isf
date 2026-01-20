package common

import (
	"time"
)

// Now 业务系统时间
func Now() time.Time {
	return time.Now().Add(time.Second * time.Duration(SvcConfig.BusinessTimeOffset))
}
