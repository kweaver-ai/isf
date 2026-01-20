package logics

import (
	"time"

	"Authentication/common"
	"Authentication/dbaccess"
)

// StartCronThread 开启定时线程任务
func StartCronThread() {
	go func() {
		hydra := dbaccess.NewDBHydra()
		const t time.Duration = 24

		for {
			// 计算下一个零点
			next := common.Now().Add(time.Hour * t)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Local().Sub(common.Now()))
			<-t.C

			var totalCount int64 = 0
			for {
				affectedRows, err := hydra.DeleteExpiredAssertions()
				totalCount += affectedRows
				if err != nil {
					common.NewLogger().Errorf("DeleteExpiredAssertions failed:%v", err)
					break
				} else if affectedRows == 0 {
					common.NewLogger().Infof("DeleteExpiredAssertions succeeded, %d rows deleted", totalCount)
					break
				} else {
					continue
				}
			}
		}
	}()
}
