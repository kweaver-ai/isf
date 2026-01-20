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

func newaSMS(ptrDB *sqlx.DB) *anonymousSMS {
	aSMS := &anonymousSMS{
		dbTrace:     ptrDB,
		batchNumber: 10,
		desKey:      "Ea8ek&ah",
		logger:      common.NewLogger(),
	}

	return aSMS
}

func TestCreate(t *testing.T) {
	Convey("create asms code", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		aSMS := newaSMS(db)
		aSMS.trace = trace
		ctx := context.Background()

		Convey("faild, db unavailable", func() {
			info := &interfaces.AnonymousSMSInfo{
				ID:          "id",
				AnonymityID: "a-id",
				PhoneNumber: "tel_number",
				Content:     "123456",
				CreateTime:  time.Now().Unix(),
			}
			tmpErr := fmt.Errorf("unknown error")
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnError(tmpErr)

			err := aSMS.Create(ctx, info)

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			info := &interfaces.AnonymousSMSInfo{
				ID:          "id",
				AnonymityID: "a-id",
				PhoneNumber: "tel_number",
				Content:     "123456",
				CreateTime:  time.Now().Unix(),
			}
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			err := aSMS.Create(ctx, info)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetInfoByID(t *testing.T) {
	Convey("get asms info by id", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		aSMS := newaSMS(db)
		aSMS.trace = trace
		ctx := context.Background()

		fields := []string{
			"f_phone_number",
			"f_anonymity_id",
			"f_content",
			"f_create_time",
		}
		info := &interfaces.AnonymousSMSInfo{}
		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectQuery("").WillReturnError(tmpErr)

			_, err = aSMS.GetInfoByID(ctx, "id")

			assert.Equal(t, err, tmpErr)
		})
		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(fields).AddRow("f2ca5307f4cb540953f0fe6526730482", "a-id", "123456", time.Unix(time.Now().Unix(), 0)))

			info, err = aSMS.GetInfoByID(ctx, "id")

			assert.Equal(t, err, nil)
			assert.Equal(t, info.PhoneNumber, "13100005678")
			assert.Equal(t, info.AnonymityID, "a-id")
			assert.Equal(t, info.Content, "123456")
			assert.Equal(t, info.CreateTime, time.Now().Unix())
		})
	})
}

func TestDeleteRecordWithinValidityPeriod(t *testing.T) {
	Convey("DeleteRecordWithinValidityPeriod", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		aSMS := newaSMS(db)
		aSMS.trace = trace
		ctx := context.Background()

		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnError(tmpErr)

			err := aSMS.DeleteRecordWithinValidityPeriod(ctx, "tel_number", "a-id", time.Minute*2)

			assert.Equal(t, err, tmpErr)
		})
		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err := aSMS.DeleteRecordWithinValidityPeriod(ctx, "tel_number", "a-id", time.Minute*2)

			assert.Equal(t, err, nil)
		})
	})
}

func TestGetExpiredRecords(t *testing.T) {
	Convey("GetExpiredRecords", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		aSMS := newaSMS(db)
		aSMS.trace = trace
		ctx := context.Background()

		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectQuery("").WillReturnError(tmpErr)

			_, err = aSMS.GetExpiredRecords(ctx, time.Minute*2)

			assert.Equal(t, err, tmpErr)
		})
		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_id"}).AddRow("id-1"))

			ids, err := aSMS.GetExpiredRecords(ctx, time.Minute*2)

			assert.Equal(t, err, nil)
			assert.Equal(t, ids[0], "id-1")
		})
	})
}

func TestDeleteByIDs(t *testing.T) {
	Convey("DeleteByIDs", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)
		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		aSMS := newaSMS(db)
		aSMS.trace = trace
		ctx := context.Background()

		Convey("failed, db unavailable", func() {
			tmpErr := fmt.Errorf("unknown error")
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnError(tmpErr)

			err = aSMS.DeleteByIDs(ctx, []string{"id-1"})

			assert.Equal(t, err, tmpErr)
		})
		Convey("success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			err = aSMS.DeleteByIDs(ctx, []string{"id-1"})

			assert.Equal(t, err, nil)
		})
	})
}
