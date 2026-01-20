package logics

import (
	"fmt"
	"testing"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/errors"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

var timeNow time.Time

func newAnonymous(cdb interfaces.DBAnonymous, mb interfaces.DrivenMessageBroker, outbox interfaces.LogicsOutbox, dbPool *sqlx.DB) *anonymous {
	return &anonymous{
		db:            cdb,
		messageBroker: mb,
		ob:            outbox,
		logger:        common.NewLogger(),
		pool:          dbPool,
	}
}

func TestNewAnonymous(t *testing.T) {
	Convey("NewAnonymous", t, func() {
		sqlDB, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		dbPool = sqlDB

		data := NewAnonymous()
		assert.NotEqual(t, data, nil)
	})
}

func TestCreate(t *testing.T) {
	Convey("Create, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymous(ctrl)
		mb := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		anony := newAnonymous(db, mb, ob, nil)

		Convey("success", func() {
			var info interfaces.AnonymousInfo
			db.EXPECT().Create(gomock.Any()).AnyTimes().Return(nil)
			err := anony.Create(info)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteByID(t *testing.T) {
	Convey("DeleteByID, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymous(ctrl)
		mb := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)
		anony := newAnonymous(db, mb, ob, nil)

		Convey("success", func() {
			db.EXPECT().DeleteByID(gomock.Any()).AnyTimes().Return(nil)
			err := anony.DeleteByID("")
			assert.Equal(t, err, nil)
		})
	})
}

func TestAuthentication(t *testing.T) {
	Convey("Authentication, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymous(ctrl)
		mb := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		sqlDB, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := sqlDB.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		anony := newAnonymous(db, mb, ob, sqlDB)

		var info interfaces.AnonymousInfo
		Convey("GetAccount error", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, testErr)
			err := anony.Authentication("", "", "")
			assert.Equal(t, err, testErr)
		})

		Convey("no anonymous ", func() {
			tmpErr := rest.NewHTTPError("record not exist", errors.AnonymityNotFound, nil)
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, tmpErr)
			err := anony.Authentication("xxx", "", "")
			assert.Equal(t, err, tmpErr)
		})

		anonymousID := "xxx1"
		anonymousePass := "zzz"
		Convey("password empty ", func() {
			info.ID = anonymousID
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, rest.NewHTTPError("wrong password", errors.AnonymityWrongPassword, nil))
		})

		Convey("access time out of range  ", func() {
			info.ID = anonymousID
			info.Password = anonymousePass
			info.AccessedTimes = 10
			info.LimitedTimes = 9
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, rest.NewHTTPError("the visits has reached the limit", errors.AnonymityReachLimitTimes, nil))
		})

		Convey("pool open error   ", func() {
			timeNow = time.Now()
			info.ID = anonymousID
			info.Password = anonymousePass
			info.LimitedTimes = -1
			info.ExpiresAtStamp = timeNow.UnixNano() + 10
			testErr := rest.NewHTTPError("error", 503000000, nil)

			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			mb.EXPECT().Publish(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectBegin().WillReturnError(testErr)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, testErr)
		})

		Convey("AddAccessTimes error   ", func() {
			testErr := rest.NewHTTPError("error", 503000000, nil)
			timeNow = time.Now()
			info.ID = anonymousID
			info.Password = anonymousePass
			info.LimitedTimes = -1
			info.ExpiresAtStamp = timeNow.UnixNano() + 10

			txMock.ExpectBegin().WillReturnError(testErr)
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			db.EXPECT().AddAccessTimes(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, testErr)
		})

		Convey("AddOutboxInfo error   ", func() {
			timeNow = time.Now()
			info.ID = anonymousID
			info.Password = anonymousePass
			info.LimitedTimes = -1
			info.ExpiresAtStamp = timeNow.UnixNano() + 10
			info.Type = "document"
			testErr := rest.NewHTTPError("error", 503000000, nil)

			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			db.EXPECT().AddAccessTimes(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			mb.EXPECT().AnonymityAuth(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, testErr)
		})

		Convey("sql commit error   ", func() {
			timeNow = time.Now()
			info.ID = anonymousID
			info.Password = anonymousePass
			info.LimitedTimes = -1
			info.ExpiresAtStamp = timeNow.UnixNano() + 10
			info.Type = "document"
			testErr := rest.NewHTTPError("error", 503000000, nil)

			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			db.EXPECT().AddAccessTimes(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			mb.EXPECT().Publish(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit().WillReturnError(testErr)
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, testErr)
		})

		Convey("Success ", func() {
			timeNow = time.Now()
			info.ID = anonymousID
			info.Password = anonymousePass
			info.LimitedTimes = -1
			info.ExpiresAtStamp = timeNow.UnixNano() + 10

			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(info, nil)
			db.EXPECT().AddAccessTimes(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			mb.EXPECT().Publish(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectBegin()
			ob.EXPECT().AddOutboxInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			txMock.ExpectCommit()
			ob.EXPECT().NotifyPushOutboxThread()
			err := anony.Authentication(anonymousID, anonymousePass, "")
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteAnonymoysOutOfDate(t *testing.T) {
	Convey("DeleteAnonymoysOutOfDate, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymous(ctrl)
		mb := mock.NewMockDrivenMessageBroker(ctrl)
		ob := mock.NewMockLogicsOutbox(ctrl)

		anony := newAnonymous(db, mb, ob, nil)

		Convey("DeleteAnonymoysOutOfDate Success", func() {
			db.EXPECT().DeleteByTime(gomock.Any()).AnyTimes().Return(nil)
			err := anony.DeleteByTime(0)
			assert.Equal(t, err, nil)
		})
	})
}

func TestGetByID(t *testing.T) {
	Convey("GetByID, db is available", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymous(ctrl)
		anony := newAnonymous(db, nil, nil, nil)

		Convey("GetByID failed", func() {
			tmpErr := fmt.Errorf("error")
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(interfaces.AnonymousInfo{}, tmpErr)
			_, err := anony.GetByID("anonymityID")
			assert.Equal(t, err, tmpErr)
		})
		Convey("GetByID success", func() {
			db.EXPECT().GetAccount(gomock.Any()).AnyTimes().Return(interfaces.AnonymousInfo{ID: "xx", VerifyMobile: true}, nil)
			info, err := anony.GetByID("anonymityID")
			assert.Equal(t, err, nil)
			assert.Equal(t, info.VerifyMobile, true)
		})
	})
}
