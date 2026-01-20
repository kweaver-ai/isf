// Package driveradapters 协议层
package driveradapters

import (
	"sync"
	"time"

	"Authentication/common"
	"Authentication/logics/session"
)

var (
	timerOnce sync.Once
	log       = common.NewLogger()
)

// StartCleanThread 启动清理过期context
func StartCleanThread() {
	sessn := session.NewSession()
	timerOnce.Do(func() {
		const t time.Duration = 24
		go func() {
			for {
				now := time.Now()
				// 计算下一个零点
				next := now.Add(time.Hour * t)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				t := time.NewTimer(next.Sub(now))
				<-t.C
				exp := now.UnixNano()

				err := sessn.EcronDelete(exp)
				if err != nil {
					log.Errorln("delete context failed:", err)
				}
			}
		}()
	})
}
