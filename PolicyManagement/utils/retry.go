package utils

import (
	"time"

	"policy_mgnt/utils/gocommon/api"
)

// 失败重试机制可设置重试次数和间隔时长
func Retry(sleep time.Duration, path, datapath string, callback func(string, string) error) (err error) {
	l := api.NewLogger()
	count := 1
	for {
		err = callback(path, datapath)
		if err != nil {
			l.Errorf("Push package failed retry count:%v error info:%s.", count, err)
			time.Sleep(sleep)
			count++
		} else {
			break
		}
	}
	return
}
