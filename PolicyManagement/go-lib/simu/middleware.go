// Package simu simulate方法包
// @File middleware.go
// @Description  简化mock 测试
package simu

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// HandlerFunc 执行函数
type HandlerFunc func()

// MiddlewareFunc 中间件函数
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type middlewareIterator struct {
	index int
	mws   []MiddlewareFunc
}

func (m *middlewareIterator) hasNext() bool {
	return m.index > 0
}
func (m *middlewareIterator) Next() MiddlewareFunc {
	if m.hasNext() {
		defer func() {
			m.index--
		}()
		return m.mws[m.index-1]
	}
	return nil
}

// ConveyTest 使用ConveyTest
func ConveyTest(name string, t *testing.T, do HandlerFunc, mws ...MiddlewareFunc) {
	mi := &middlewareIterator{index: len(mws), mws: mws}
	convey.Convey(name, t, func() {
		for mw := mi.Next(); mw != nil; mw = mi.Next() {
			do = mw(do)
		}
		do()
	})
}
