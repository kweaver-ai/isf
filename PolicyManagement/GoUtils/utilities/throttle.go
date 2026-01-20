package utilities

import "time"

/*
Throttle 是一个限流函数，一定时间内的多次调用会被合并为一次调用。

sayHi := Throttle(func() { fmt.Println("Hi") }, time.Second)

for range make([]struct{}, 3) {
	sayHi()
}

// 只会打印一次“Hi”
*/
func Throttle(f func(), period time.Duration) func() {
	var (
		timer    *time.Timer
		lastCall time.Time
	)

	return func() {
		if timer == nil {
			timer = time.AfterFunc(period, f)
		} else {
			// 未到达period时间的调用，会先将前一次定时器注销掉并重新计时
			if time.Since(lastCall) < period {
				timer.Stop()
				timer.Reset(period)
			} else {
				timer = time.AfterFunc(period, f)
			}
		}

		lastCall = time.Now()
	}
}
