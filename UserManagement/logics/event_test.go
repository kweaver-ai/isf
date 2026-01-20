// Package logics event Anyshare 业务逻辑层 -事件
package logics

import (
	"errors"
	"testing"

	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestDeptDeleted(t *testing.T) {
	Convey("部门被删除事件处理-逻辑层", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := event{}

		testErr := errors.New("test1")
		errorFunc := func(string) error {
			return testErr
		}
		Convey("异常报错", func() {
			c.deptDeletedHandlers = []func(string) error{errorFunc}

			err := c.DeptDeleted("")
			assert.Equal(t, err, testErr)
		})

		Convey("正常返回", func() {
			c.deptDeletedHandlers = []func(string) error{}

			err := c.DeptDeleted("")
			assert.Equal(t, err, nil)
		})
	})
}
