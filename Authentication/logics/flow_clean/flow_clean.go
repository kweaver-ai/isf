// Package flowclean 逻辑层
package flowclean

import (
	"sync"
	"time"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/logics"
)

var (
	fcOnce sync.Once
	fc     *flowClean
)

type flowClean struct {
	logger common.Logger
	db     interfaces.DBFlowClean
}

// NewFlowClean 创建flowClean操作对象
func NewFlowClean() *flowClean {
	fcOnce.Do(func() {
		fc = &flowClean{
			logger: common.NewLogger(),
			db:     logics.DBFlowClean,
		}

		// 清理线程
		fc.StartInValidFlowThread()
	})

	return fc
}

// CleanFlow 清理flow
func (fc *flowClean) CleanFlow() (err error) {
	fc.logger.Debugf("logics CleanFlow begin")

	// 清理过期的refresh token信息(默认三个月)
	expiredTime := common.SvcConfig.FlowExpiredTime * 24 * 3600
	err = fc.db.CleanExpiredRefresh(expiredTime)
	if err != nil {
		fc.logger.Errorf("CleanExpiredRefresh failed, err: %v", err)
		return
	}

	// 获取所有过期且无refresh_token的flow记录
	ids, err := fc.db.GetAllExpireFlowIDs(expiredTime)
	if err != nil {
		fc.logger.Errorf("GetAllExpireFlowIDs failed, err: %v", err)
		return
	}

	// 删除flow
	err = fc.db.CleanFlow(ids)
	if err != nil {
		fc.logger.Errorf("CleanFlow failed, err: %v", err)
	}

	fc.logger.Debugf("logics CleanFlow end")
	return
}

// StartCleanThread 启动清理线程
func (fc *flowClean) StartInValidFlowThread() {
	go func() {
		// 获取清理时间，设置下次清理时间
		tCleanTime, err := time.Parse("15:04:05", common.SvcConfig.FlowCleanTime)
		if err != nil {
			// 必须退出，否则会导致定时任务一直执行
			fc.logger.Errorln("flow clean thread stop !!!!!, parse flow clean time failed:", err)
			return
		}

		for {
			err := fc.CleanFlow()
			if err != nil {
				fc.logger.Errorln("delete expired flow failed:", err)
			}

			now := common.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), tCleanTime.Hour(), tCleanTime.Minute(), tCleanTime.Second(), 0, now.Location())
			if now.After(next) {
				// 如果当前时间已经过了指定时间，则计算明天的指定时间
				next = next.AddDate(0, 0, 1)
			}

			t := time.NewTimer(next.Sub(now))

			<-t.C
		}
	}()
}
