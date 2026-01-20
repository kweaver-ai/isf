package dbaccess

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"gotest.tools/assert"

	"Authentication/common"
	"Authentication/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	mocks "Authentication/interfaces/mock"
)

func newDBTicket(ptrDB *sqlx.DB) *ticket {
	t := &ticket{
		dbTrace:     ptrDB,
		batchNumber: 100,
		logger:      common.NewLogger(),
	}

	return t
}

func TestCreateTicket(t *testing.T) {
	Convey("Create ticket", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBTicket(db)
		tic.trace = trace

		info := &interfaces.TicketInfo{
			ID:         "id",
			UserID:     "dfc9b098-dac4-11ee-b50a-028586548cf7",
			ClientID:   "0c8839f4-894c-452c-ae96-911df5e04c64",
			CreateTime: time.Now().Unix(),
		}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("faild, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			mock.ExpectExec("").WillReturnError(tmpErr)

			err := tic.Create(ctx, info)

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err := tic.Create(ctx, info)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetTicketByID(t *testing.T) {
	Convey("get ticket info by id", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBTicket(db)
		tic.trace = trace

		fields := []string{
			"f_user_id",
			"f_client_id",
			"f_create_time",
		}
		info := &interfaces.TicketInfo{}

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")

			mock.ExpectQuery("").WillReturnError(tmpErr)

			_, err = tic.GetTicketByID(ctx, "id")

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("dfc9b098-dac4-11ee-b50a-028586548cf7", "0c8839f4-894c-452c-ae96-911df5e04c64", time.Now().Unix()))

			info, err = tic.GetTicketByID(ctx, "id")

			assert.Equal(t, err, nil)
			assert.Equal(t, info.ID, "id")
			assert.Equal(t, info.UserID, "dfc9b098-dac4-11ee-b50a-028586548cf7")
			assert.Equal(t, info.ClientID, "0c8839f4-894c-452c-ae96-911df5e04c64")
			assert.Equal(t, info.CreateTime, time.Now().Unix())
		})
	})
}

func TestDeleteTicketByIDs(t *testing.T) {
	Convey("DeleteTicketByIDs", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		tic := newDBTicket(db)
		tic.trace = trace

		ctx := context.Background()
		trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			mock.ExpectExec("").WillReturnError(tmpErr)

			err = tic.DeleteByIDs(ctx, []string{"id-1"})

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err := tic.DeleteByIDs(ctx, []string{"id-1"})

			assert.Equal(t, err, nil)
		})
	})
}
