package session

import (
	"errors"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"Authentication/common"
	"Authentication/interfaces/mock"
)

func TestDeleteHydraSession(t *testing.T) {
	Convey("Delete login and consent session", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		loOutbox := mock.NewMockLogicsOutbox(ctrl)

		var err error
		db, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := db.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		hSession := &hydraSession{
			hydraAdmin: hydraAdmin,
			ob:         loOutbox,
			pool:       db,
			logger:     common.NewLogger(),
		}

		Convey("AddOutboxInfo error", func() {
			tErr := errors.New("err")
			txMock.ExpectBegin()
			loOutbox.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tErr)
			txMock.ExpectRollback()

			err := hSession.Delete("xx", "")

			assert.Equal(t, err, tErr)
		})

		Convey("sql commit error", func() {
			tErr := errors.New("err")
			txMock.ExpectBegin()
			loOutbox.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit().WillReturnError(tErr)

			err := hSession.Delete("xx", "")

			assert.Equal(t, err, nil)
		})

		Convey("delete session success", func() {
			txMock.ExpectBegin()
			loOutbox.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			loOutbox.EXPECT().NotifyPushOutboxThread()

			err := hSession.Delete("xx", "")

			assert.Equal(t, err, nil)
		})
	})
}
