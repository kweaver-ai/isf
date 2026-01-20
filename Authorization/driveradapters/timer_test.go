package driveradapters

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"Authorization/common"
	"Authorization/interfaces"
	"Authorization/interfaces/mock"
)

func newTimer(a interfaces.LogicsPolicy) Timer {
	return &timer{
		log:    common.NewLogger(),
		policy: a,
	}
}

func TestStartCleanThread(t *testing.T) {
	Convey("StartCleanThread", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		a := mock.NewMockLogicsPolicy(ctrl)
		timer := newTimer(a)
		assert.NotEqual(t, timer, nil)

		Convey("DeleteByTime error", func() {
			a.EXPECT().DeleteByEndTime(gomock.Any()).AnyTimes().Return(errors.New("delete error"))
			timer.StartCleanThread()
		})

		Convey("success", func() {
			a.EXPECT().DeleteByEndTime(gomock.Any()).AnyTimes().Return(nil)
			timer.StartCleanThread()
		})
	})
}
