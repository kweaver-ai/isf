package dbaccess

import (
	"fmt"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"gotest.tools/assert"

	"Authentication/common"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func newDBFlowClean(ptrDB *sqlx.DB) *flowClean {
	t := &flowClean{
		db:     ptrDB,
		logger: common.NewLogger(),
	}

	return t
}

func TestCleanExpiredRefresh(t *testing.T) {
	Convey("CleanExpiredRefresh", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBFlowClean(db)

		Convey("faild, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			mock.ExpectExec("").WillReturnError(tmpErr)

			err := tic.CleanExpiredRefresh(0)

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 0))

			err := tic.CleanExpiredRefresh(0)

			assert.Equal(t, err, nil)
		})
	})
}

func TestCleanFlow(t *testing.T) {
	Convey("CleanFlow", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBFlowClean(db)

		Convey("faild, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			mock.ExpectExec("").WillReturnError(tmpErr)

			err := tic.CleanFlow([]string{appID})

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err := tic.CleanFlow([]string{appID})

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetAllExpireFlowIDs(t *testing.T) {
	Convey("GetAllExpireFlowIDs", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBFlowClean(db)

		Convey("faild, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			mock.ExpectQuery("").WillReturnError(tmpErr)

			_, err := tic.GetAllExpireFlowIDs(10)

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"login_challenge", "challenge_id", "requested_at"}).AddRow(appID, "xxx", appID))

			out, err := tic.GetAllExpireFlowIDs(10)

			assert.Equal(t, err, nil)
			assert.Equal(t, len(out), 1)
			assert.Equal(t, out[0], appID)
		})
	})
}
