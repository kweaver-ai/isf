package dbaccess

import (
	"context"
	"errors"
	"policy_mgnt/common"
	"testing"

	mocks "policy_mgnt/interfaces/mock"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

func TestGetConfig(t *testing.T) {
	Convey("GetConfig, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		cof := &config{
			db:    db,
			trace: trace,
			log:   common.NewLogger(),
		}
		ctx := context.Background()

		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnError(errors.New(""))
			_, httpErr := cof.GetConfig(ctx, "test")
			assert.NotEqual(t, httpErr, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow("1"))
			count, err := cof.GetConfig(ctx, "test")
			assert.Equal(t, err, nil)
			assert.Equal(t, count, "1")
		})
	})
}
