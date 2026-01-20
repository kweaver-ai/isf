package utilities

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestThrottle(t *testing.T) {
	Convey("测试Throttle限流函数", t, func() {
		assert := assert.New(t)

		Convey("多次连续调用，应当只执行一次。", func() {
			count := 0

			fn := func() {
				count++
			}

			lazyAdder := Throttle(fn, time.Millisecond)

			// 连续多次调用
			for range make([]struct{}, 5) {
				lazyAdder()
			}

			// 适当增加一些冗余时长
			time.Sleep(time.Millisecond * 2)

			assert.Equal(1, count)
		})

		Convey("Throttle返回的函数可以复用", func() {

			count := 0

			fn := func() {
				count++
			}

			lazyAdder := Throttle(fn, time.Millisecond)

			lazyAdder()
			time.Sleep(time.Millisecond * 2)
			assert.Equal(1, count)

			lazyAdder()
			time.Sleep(time.Millisecond * 2)
			assert.Equal(2, count)

			// 再次触发连续多次调用
			for range make([]struct{}, 5) {
				lazyAdder()
			}
			time.Sleep(time.Millisecond * 2)
			assert.Equal(3, count)

		})
	})

}
