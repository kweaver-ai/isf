// Package driveradapters 定时清理线程
package driveradapters

import (
	"sync"
	"time"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// Timer 定时任务接口
type Timer interface {
	// StartCleanThread 启动清理过期匿名账户线程
	StartCleanThread()
}

var (
	timerOnce sync.Once
	t         Timer
)

type timer struct {
	once      sync.Once
	log       common.Logger
	anonymity interfaces.LogicsAnonymous
}

// NewTimer 创建定时任务对象
func NewTimer() Timer {
	timerOnce.Do(func() {
		t = &timer{
			log:       common.NewLogger(),
			anonymity: logics.NewAnonymous(),
		}
	})

	return t
}

// StartCleanThread 启动清理过期匿名账户线程
func (tt *timer) StartCleanThread() {
	tt.once.Do(func() {
		const t time.Duration = 24
		go func() {
			for {
				curTime := common.Now().UnixNano()
				err := tt.anonymity.DeleteByTime(curTime)
				if err != nil {
					tt.log.Errorln("delete expired anonymity failed:", err)
				}

				time.Sleep(t * time.Hour)
			}
		}()
	})
}
