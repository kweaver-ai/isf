// Package logics avatar AnyShare 用户头像业务逻辑层
package logics

import (
	"context"
	"errors"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func TestNewAvatar(t *testing.T) {
	Convey("NewAvatar", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		avatarDB := mock.NewMockDBAvatar(ctrl)

		sqlDB, _, err := sqlx.New()
		assert.Equal(t, err, nil)
		dbPool = sqlDB

		SetDBAvatar(avatarDB)

		outInfo := make([]interfaces.AvatarOSSInfo, 0)
		avatarDB.EXPECT().GetUselessAvatar(gomock.Any()).AnyTimes().Return(outInfo, nil)
		data := NewAvatar()
		assert.NotEqual(t, data, nil)
	})
}

func TestGet(t *testing.T) {
	Convey("Get, db is available", t, func() {
		test := setGinMode()
		defer test()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		avatarDB := mock.NewMockDBAvatar(ctrl)
		oss := mock.NewMockDnOSSGateWay(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		ava := &avatar{
			avatarDB:   avatarDB,
			ossGateway: oss,
			logger:     common.NewLogger(),
			trace:      trace,
		}

		userID := "xxxssss"
		testErr := errors.New("sdad")
		info := interfaces.AvatarOSSInfo{}
		visitor := &interfaces.Visitor{
			ID:   interfaces.SystemAuditAdmin,
			Type: interfaces.App,
		}
		ctx := context.Background()

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

		Convey("avatarDB Get error", func() {
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(info, testErr)
			_, err := ava.Get(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		info.Key = "key"
		req := ""
		Convey("ossGateway GetDownloadURL error", func() {
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(info, nil)
			oss.EXPECT().GetDownloadURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(req, testErr)
			_, err := ava.Get(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		req = "url1"
		Convey("success", func() {
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(info, nil)
			oss.EXPECT().GetDownloadURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(req, nil)
			out, err := ava.Get(ctx, visitor, userID)
			assert.Equal(t, err, nil)
			assert.Equal(t, out, req)
		})
	})
}

func TestUpdate(t *testing.T) {
	Convey("Update, db is available", t, func() {
		test := setGinMode()
		defer test()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		avatarDB := mock.NewMockDBAvatar(ctrl)
		oss := mock.NewMockDnOSSGateWay(ctrl)
		trace := mock.NewMockTraceClient(ctrl)

		dPool, txMock, err := sqlx.New()
		assert.Equal(t, err, nil)
		defer func() {
			if closeErr := dPool.Close(); closeErr != nil {
				assert.Equal(t, 1, 1)
			}
		}()

		ava := &avatar{
			avatarDB:          avatarDB,
			ossGateway:        oss,
			pool:              dPool,
			logger:            common.NewLogger(),
			maxAvatarSize:     48 * 1024,
			maxContentTypeLen: 50,
			trace:             trace,
		}

		visitor := &interfaces.Visitor{
			ID:   interfaces.SystemAuditAdmin,
			Type: interfaces.App,
		}
		ctx := context.Background()

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes().Return(ctx, nil)
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

		Convey("visitor type is error", func() {
			err := ava.Update(ctx, visitor, "png", nil)
			assert.Equal(t, err, rest.NewHTTPError("only support normal user", rest.BadRequest, nil))
		})

		visitor.Type = interfaces.RealName
		testErr := errors.New("xxx1")

		Convey("user is not normal", func() {
			err := ava.Update(ctx, visitor, "png", nil)
			assert.Equal(t, err, rest.NewHTTPError("only support normal user", rest.BadRequest, nil))
		})

		visitor.ID = "111xxx"
		tempBuff := make([]byte, 50*1024)
		Convey("file type too long", func() {
			err := ava.Update(ctx, visitor, "pngpngpng1pngpngpng1pngpngpng1pngpngpng1pngpngpng11", tempBuff)
			assert.Equal(t, err, rest.NewHTTPError("invalid params, file type error", rest.BadRequest, nil))
		})

		Convey("file is too big", func() {
			err := ava.Update(ctx, visitor, "png", tempBuff)
			assert.Equal(t, err, rest.NewHTTPError("invalid params, file is too big", rest.BadRequest, nil))
		})

		buff := make([]byte, 10)
		ossInfos := make([]interfaces.OSSInfo, 0)
		Convey("evfs GetSiteDefaultOSS error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		Convey("evfs get no oss", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, rest.NewHTTPError("no available oss", rest.InternalServerError, nil))
		})

		data := interfaces.OSSInfo{ID: "xxxxx"}
		ossInfos = append(ossInfos, data)
		Convey("avatarDB Add error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		Convey("ossGateway UploadFile error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			oss.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		tmp1Info := interfaces.AvatarOSSInfo{}
		Convey("pool begin error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			oss.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(tmp1Info, nil)
			txMock.ExpectBegin().WillReturnError(testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		Convey("avatarDB UpdateStatusByKey 1 error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			oss.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(tmp1Info, nil)
			txMock.ExpectBegin()
			avatarDB.EXPECT().SetAvatarUnableByID(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		Convey("avatarDB UpdateStatusByKey 2 error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			oss.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(tmp1Info, nil)
			txMock.ExpectBegin()
			avatarDB.EXPECT().SetAvatarUnableByID(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			avatarDB.EXPECT().UpdateStatusByKey(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(testErr)
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, testErr)
		})

		Convey("success", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(ossInfos, nil)
			avatarDB.EXPECT().Add(gomock.Any()).AnyTimes().Return(nil)
			oss.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			avatarDB.EXPECT().Get(gomock.Any()).AnyTimes().Return(tmp1Info, nil)
			txMock.ExpectBegin()
			avatarDB.EXPECT().SetAvatarUnableByID(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			avatarDB.EXPECT().UpdateStatusByKey(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			txMock.ExpectCommit()
			err := ava.Update(ctx, visitor, "png", buff)
			assert.Equal(t, err, nil)
		})
	})
}

func TestDeleteUselessAvatar(t *testing.T) {
	Convey("deleteUselessAvatar, db is available", t, func() {
		test := setGinMode()
		defer test()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		avatarDB := mock.NewMockDBAvatar(ctrl)
		oss := mock.NewMockDnOSSGateWay(ctrl)

		ava := &avatar{
			avatarDB:   avatarDB,
			ossGateway: oss,
			logger:     common.NewLogger(),
		}

		testErr := errors.New("xxx")
		data := make([]interfaces.AvatarOSSInfo, 1)
		Convey("avatarDB GetUselessAvatar error", func() {
			avatarDB.EXPECT().GetUselessAvatar(gomock.Any()).AnyTimes().Return(data, testErr)
			ava.deleteUselessAvatar()
		})

		data[0] = interfaces.AvatarOSSInfo{}
		Convey("ossGateway DeleteFile error", func() {
			avatarDB.EXPECT().GetUselessAvatar(gomock.Any()).AnyTimes().Return(data, nil)
			oss.EXPECT().DeleteFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			ava.deleteUselessAvatar()
		})

		Convey("avatarDB Delete error", func() {
			avatarDB.EXPECT().GetUselessAvatar(gomock.Any()).AnyTimes().Return(data, nil)
			oss.EXPECT().DeleteFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			avatarDB.EXPECT().Delete(gomock.Any()).AnyTimes().Return(testErr)
			ava.deleteUselessAvatar()
		})
	})
}

func TestGetAvaliableOSS(t *testing.T) {
	Convey("getAvaliableOSS, db is available", t, func() {
		test := setGinMode()
		defer test()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		avatarDB := mock.NewMockDBAvatar(ctrl)
		oss := mock.NewMockDnOSSGateWay(ctrl)

		ava := &avatar{
			avatarDB:   avatarDB,
			ossGateway: oss,
			logger:     common.NewLogger(),
		}

		testErr := errors.New("xxx")
		ctx := context.Background()
		visitor := &interfaces.Visitor{
			ID:   interfaces.SystemAuditAdmin,
			Type: interfaces.App,
		}
		Convey("evfs GetLocalOSSInfo error", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			_, err := ava.getAvaliableOSS(ctx, visitor)
			assert.Equal(t, err, testErr)
		})

		Convey("站点信息为空，报错", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
			_, err := ava.getAvaliableOSS(ctx, visitor)
			assert.Equal(t, err, rest.NewHTTPError("no available oss", rest.InternalServerError, nil))
		})

		info := interfaces.OSSInfo{ID: "xxx", BDefault: false}
		temp := []interfaces.OSSInfo{info}
		Convey("success", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(temp, nil)
			out, err := ava.getAvaliableOSS(ctx, visitor)
			assert.Equal(t, err, nil)
			assert.Equal(t, out, "xxx")
		})

		info1 := interfaces.OSSInfo{ID: "xxx1", BDefault: true}
		temp = []interfaces.OSSInfo{info, info1}
		Convey("success1", func() {
			oss.EXPECT().GetLocalEnabledOSSInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(temp, nil)
			out, err := ava.getAvaliableOSS(ctx, visitor)
			assert.Equal(t, err, nil)
			assert.Equal(t, out, "xxx1")
		})
	})
}
