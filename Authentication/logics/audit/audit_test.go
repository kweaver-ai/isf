package audit

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func TestLog(t *testing.T) {
	Convey("Log", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var err error
		db, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := db.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ob := mock.NewMockLogicsOutbox(ctrl)
		eacplog := mock.NewMockDnEacpLog(ctrl)
		auditlog := &audit{
			logger:  common.NewLogger(),
			pool:    db,
			ob:      ob,
			eacpLog: eacplog,
		}

		testErr := errors.New("some error")

		strTopic := "topic"
		message := make(map[string]interface{})
		Convey("pool beigin error", func() {
			txMock.ExpectBegin().WillReturnError(testErr)
			err := auditlog.Log(strTopic, message)
			assert.Equal(t, err, testErr)
		})

		Convey("AddOutboxInfo error", func() {
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			txMock.ExpectRollback()
			err := auditlog.Log(strTopic, message)
			assert.Equal(t, err, testErr)
		})

		Convey("success error", func() {
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()
			err := auditlog.Log(strTopic, message)
			assert.Equal(t, err, nil)
		})
	})
}
