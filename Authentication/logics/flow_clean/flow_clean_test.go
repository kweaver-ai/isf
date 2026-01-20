package flowclean

import (
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces/mock"
)

func TestCleanFlow(t *testing.T) {
	Convey("CleanFlow", t, func() {
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fc := mock.NewMockDBFlowClean(ctrl)
		n := flowClean{
			logger: common.NewLogger(),
			db:     fc,
		}

		Convey("CleanExpiredRefresh error", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			fc.EXPECT().CleanExpiredRefresh(gomock.Any()).AnyTimes().Return(testErr)
			err := n.CleanFlow()
			assert.Equal(t, err, testErr)
		})

		Convey("GetAllExpireFlowIDs error", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			fc.EXPECT().CleanExpiredRefresh(gomock.Any()).AnyTimes().Return(nil)
			fc.EXPECT().GetAllExpireFlowIDs(gomock.Any()).AnyTimes().Return(nil, testErr)
			err := n.CleanFlow()
			assert.Equal(t, err, testErr)
		})

		Convey("CleanFlow error", func() {
			testErr := rest.NewHTTPError("error", rest.InternalServerError, nil)
			fc.EXPECT().CleanExpiredRefresh(gomock.Any()).AnyTimes().Return(nil)
			fc.EXPECT().GetAllExpireFlowIDs(gomock.Any()).AnyTimes().Return(nil, nil)
			fc.EXPECT().CleanFlow(gomock.Any()).AnyTimes().Return(testErr)
			err := n.CleanFlow()
			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			fc.EXPECT().CleanExpiredRefresh(gomock.Any()).AnyTimes().Return(nil)
			fc.EXPECT().GetAllExpireFlowIDs(gomock.Any()).AnyTimes().Return(nil, nil)
			fc.EXPECT().CleanFlow(gomock.Any()).AnyTimes().Return(nil)
			err := n.CleanFlow()
			assert.Equal(t, err, nil)
		})
	})
}
