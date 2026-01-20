package dbaccess

import (
	"context"
	"errors"
	"policy_mgnt/common"
	"policy_mgnt/interfaces"
	"testing"

	mocks "policy_mgnt/interfaces/mock"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
)

func TesGetAuthorizedProducts(t *testing.T) {
	Convey("GetAuthorizedProducts, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		lic := &license{
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
			_, httpErr := lic.GetAuthorizedProducts(ctx, []string{"test"})
			assert.NotEqual(t, httpErr, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_account_id", "f_product"}).AddRow("test", "product1"))
			count, err := lic.GetAuthorizedProducts(ctx, []string{"test"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(count), 1)
			assert.Equal(t, count["test"].Product[0], "product1")
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"f_account_id", "f_product"}).AddRow("test", "product1").AddRow("test2", "product2").AddRow("test", "product3"))
			count, err := lic.GetAuthorizedProducts(ctx, []string{"test"})
			assert.Equal(t, err, nil)
			assert.Equal(t, len(count), 2)

			temp := count["test"]
			if temp.Product[0] == "product1" {
				assert.Equal(t, temp.Product[1], "product3")
			} else {
				assert.Equal(t, temp.Product[0], "product3")
				assert.Equal(t, temp.Product[1], "product1")
			}
			temp = count["test2"]
			assert.Equal(t, temp.Product[0], "product2")
		})
	})
}

func TestDeleteAuthorizedProducts(t *testing.T) {
	Convey("DeleteAuthorizedProducts, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		lic := &license{
			db:    db,
			trace: trace,
			log:   common.NewLogger(),
		}
		ctx := context.Background()

		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New(""))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.DeleteAuthorizedProducts(ctx, []interfaces.ProductInfo{{AccountID: "test", Product: "product1"}}, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.DeleteAuthorizedProducts(ctx, []interfaces.ProductInfo{{AccountID: "test", Product: "product1"}}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestAddAuthorizedProducts(t *testing.T) {
	Convey("AddAuthorizedProducts, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		lic := &license{
			db:    db,
			trace: trace,
			log:   common.NewLogger(),
		}
		ctx := context.Background()

		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New(""))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.AddAuthorizedProducts(ctx, []interfaces.ProductInfo{{AccountID: "test", Product: "product1"}}, tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.AddAuthorizedProducts(ctx, []interfaces.ProductInfo{{AccountID: "test", Product: "product1"}}, tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteUserAuthorizedProducts(t *testing.T) {
	Convey("DeleteUserAuthorizedProducts, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		lic := &license{
			db:    db,
			trace: trace,
			log:   common.NewLogger(),
		}
		ctx := context.Background()

		Convey("execute error", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnError(errors.New(""))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.DeleteUserAuthorizedProducts(ctx, "test", tx)
			assert.NotEqual(t, err, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectBegin()
			mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			tx, _ := db.Begin()

			err := lic.DeleteUserAuthorizedProducts(ctx, "test", tx)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetProductsAuthorizedCount(t *testing.T) {
	Convey("GetProductsAuthorizedCount, db is available", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		trace := mocks.NewMockTraceClient(ctrl)

		db, mock, err := sqlx.New()
		assert.Equal(t, err, nil)

		common.InitARTrace("test")

		lic := &license{
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
			_, httpErr := lic.GetProductsAuthorizedCount(ctx, "test")
			assert.NotEqual(t, httpErr, nil)
		})

		Convey("execute success", func() {
			trace.EXPECT().SetClientSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddClientTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"COUNT(f_account_id)"}).AddRow("1"))
			count, err := lic.GetProductsAuthorizedCount(ctx, "test")
			assert.Equal(t, err, nil)
			assert.Equal(t, count, 1)
		})
	})
}
