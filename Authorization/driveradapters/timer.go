// Package driveradapters 定时清理线程
package driveradapters

import (
	"sync"
	"time"

	"Authorization/common"
	"Authorization/interfaces"
	"Authorization/logics"
)

// Timer 定时任务接口
type Timer interface {
	// StartCleanThread 启动清理过期匿名SharedLink线程
	StartCleanThread()
}

var (
	timerOnce sync.Once
	t         Timer
)

type timer struct {
	once   sync.Once
	log    common.Logger
	policy interfaces.LogicsPolicy
}

// NewTimer 创建定时任务对象
func NewTimer() Timer {
	timerOnce.Do(func() {
		t = &timer{
			log:    common.NewLogger(),
			policy: logics.NewPolicy(),
		}
	})

	return t
}

// StartCleanThread 启动清理过期匿名SharedLink线程
func (tt *timer) StartCleanThread() {
	tt.once.Do(func() {
		const t time.Duration = 24
		go func() {
			for {
				curTime := common.GetCurrentMicrosecondTimestamp()
				err := tt.policy.DeleteByEndTime(curTime)
				if err != nil {
					tt.log.Errorln("policy DeleteByEndTime failed:", err)
				}
				time.Sleep(t * time.Hour)
			}
		}()
	})
}
