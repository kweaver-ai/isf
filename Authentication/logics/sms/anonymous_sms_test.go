package sms

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/go-playground/assert"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

func newaSMS(db interfaces.DBAnonymousSMS, sharemgnt interfaces.DnShareMgnt) *anonymousSMS {
	// pem 解码
	blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
	// X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	pattern := `^1[3456789]\d{9}$`
	re := regexp.MustCompile(pattern)

	return &anonymousSMS{
		dbAnonymousSMS: db,
		sharemgnt:      sharemgnt,
		numerics:       [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		smsWidth:       6,
		smsExpiration:  time.Minute * 2,
		re:             re,
		privateKey:     privateKey,
		logger:         common.NewLogger(),
		i18n: common.NewI18n(common.I18nMap{
			i18nSMSVCodeExpired: {
				interfaces.SimplifiedChinese:  "验证码已过期",
				interfaces.TraditionalChinese: "驗證碼已過期",
				interfaces.AmericanEnglish:    "The verification code has expired",
			},
			i18nSMSVerificationFailed: {
				interfaces.SimplifiedChinese:  "验证码校验失败",
				interfaces.TraditionalChinese: "驗證碼校驗失敗",
				interfaces.AmericanEnglish:    "The CAPTCHA verification failed. Please try again",
			},
			i18nSMSInvalidPhoneNumber: {
				interfaces.SimplifiedChinese:  "手机号不合法",
				interfaces.TraditionalChinese: "手機號不合法",
				interfaces.AmericanEnglish:    "Invalid tel number",
			},
		}),
	}
}

//nolint:lll
func TestCreateAndSendVCode(t *testing.T) {
	Convey("CreateAndSendVCode", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)
		db := mock.NewMockDBAnonymousSMS(ctrl)
		sharemgnt := mock.NewMockDnShareMgnt(ctrl)
		aSMS := newaSMS(db, sharemgnt)
		aSMS.trace = trace

		visitor := interfaces.Visitor{
			Language: interfaces.AmericanEnglish,
		}
		ctx := context.Background()
		phoneNumber := "dxvDBvRM19QkzJzxwiBP+70XXYvSr/I0eEMh4ZqvWqPvkmMCNXS+5qDxVG5wEBtNsJbL2nUVpekzjYUKYFeEby4MWakfCNkEgy8QPbGv4/fpcZNk2J1EHgcc5i2gaf2ZXSYhLcjkprScSYpu9D/Qr5akysCsxdg1oz2s4v5Vu31l+R+jZnnFXQxA7dcK+RiuGaa/E6gFbfOAJKorZxUkvwNe0tO4VG21XQxB/3BhvMuyz95QvmzJmp5HXR6UtKD8HM4xWTWm1014id8Ryc8IABiM4B688sSgI43ORhHSLs/sAYNqMjp8c4cTwvaHJB/56TkqInA9tTpkJ7KriKernA=="
		Convey("invalid phone_number", func() {
			phoneNumber = "IVGrqLBCKmlbClaGfMSjtTriORSOkdCmRi/p4skCpqo1+KgCx3K/qI+7ZkSBHEFxAfOQLnwQXbWdOKs1+WLuiWnKg7zcdZzJRu+rRzAxIU1jVNqiFRYxTReShvr6m+SU58iAEUeBKpNmXr2irIEsgNG3E9RukQ5eQqo0wMR+sgowVd6JOo8xWSY37ZItY7DFTcOlGy8SyvojpxovQEkq+mYZwtzNkMt20aa9VWDb91K2XiVQwE08+J6c3scYrlwinEjVVx0idVbmmzchJdQY7RQlqyrMPwvdMW4vw1Y+9HE2vWJlTIz3m4bDrBt9z61pbHISuV6g2fIGH/c9yniEZg=="
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()

			_, err := aSMS.CreateAndSendVCode(ctx, &visitor, phoneNumber, "xx")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "Invalid tel number"))
		})
		Convey("DeleteRecordWithinValidityPeriod failed", func() {
			tmpErr := fmt.Errorf("db error")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().DeleteRecordWithinValidityPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tmpErr)

			_, err := aSMS.CreateAndSendVCode(ctx, &visitor, phoneNumber, "xx")

			assert.Equal(t, err, tmpErr)
		})
		Convey("Create failed", func() {
			tmpErr := fmt.Errorf("db error")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().DeleteRecordWithinValidityPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tmpErr)

			_, err := aSMS.CreateAndSendVCode(ctx, &visitor, phoneNumber, "xx")

			assert.Equal(t, err, tmpErr)
		})
		Convey("failed, send vcode failed", func() {
			tmpErr := fmt.Errorf("send error")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().DeleteRecordWithinValidityPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			sharemgnt.EXPECT().UsrmSendAnonymousSMSVCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(tmpErr)

			_, err := aSMS.CreateAndSendVCode(ctx, &visitor, phoneNumber, "xx")

			assert.Equal(t, err, tmpErr)
		})
		Convey("success", func() {
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().DeleteRecordWithinValidityPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			sharemgnt.EXPECT().UsrmSendAnonymousSMSVCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			vcodeID, err := aSMS.CreateAndSendVCode(ctx, &visitor, phoneNumber, "xx")

			assert.Equal(t, err, nil)
			assert.NotEqual(t, vcodeID, "")
		})
	})
}

func TestValidate(t *testing.T) {
	Convey("Validate VCode", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		trace := mock.NewMockTraceClient(ctrl)
		db := mock.NewMockDBAnonymousSMS(ctrl)
		sharemgnt := mock.NewMockDnShareMgnt(ctrl)
		aSMS := newaSMS(db, sharemgnt)
		aSMS.trace = trace

		visitor := &interfaces.Visitor{Language: interfaces.SimplifiedChinese}
		ctx := context.Background()
		Convey("failed, GetInfoByID error", func() {
			tmpErr := fmt.Errorf("error")
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().GetInfoByID(gomock.Any(), gomock.Any()).Return(nil, tmpErr)

			_, err := aSMS.Validate(ctx, visitor, "", "", "")

			assert.Equal(t, err, tmpErr)
		})
		Convey("failed, vcode expired", func() {
			info := &interfaces.AnonymousSMSInfo{}
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().GetInfoByID(gomock.Any(), gomock.Any()).Return(info, nil)

			_, err := aSMS.Validate(ctx, visitor, "", "", "")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSVCodeExpired, interfaces.SimplifiedChinese)))
		})
		Convey("failed, invalid vcode", func() {
			info := &interfaces.AnonymousSMSInfo{Content: "123456", CreateTime: time.Now().Unix()}
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().GetInfoByID(gomock.Any(), gomock.Any()).Return(info, nil)

			_, err := aSMS.Validate(ctx, visitor, "", "", "")

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, aSMS.i18n.Load(i18nSMSVerificationFailed, interfaces.SimplifiedChinese)))
		})
		Convey("success", func() {
			info := &interfaces.AnonymousSMSInfo{CreateTime: time.Now().Unix()}
			trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
			trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
			trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
			db.EXPECT().GetInfoByID(gomock.Any(), gomock.Any()).Return(info, nil)

			_, err := aSMS.Validate(ctx, visitor, "", "", "")

			assert.Equal(t, err, nil)
		})
	})
}

func TestUpdateSMSExpiration(t *testing.T) {
	Convey("更新匿名登录短信验证码过期时间", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymousSMS(ctrl)
		sharemgnt := mock.NewMockDnShareMgnt(ctrl)
		aSMS := newaSMS(db, sharemgnt)
		exp := time.Duration(2) * time.Minute
		aSMS.UpdateSMSExpiration(2)
		assert.Equal(t, exp, aSMS.smsExpiration)
	})
}

func TestInitSMSExpiration(t *testing.T) {
	Convey("初始化匿名登录验证码过期时间", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock.NewMockDBAnonymousSMS(ctrl)
		conf := mock.NewMockConf(ctrl)
		sharemgnt := mock.NewMockDnShareMgnt(ctrl)
		aSMS := newaSMS(db, sharemgnt)
		aSMS.conf = conf

		cfg := interfaces.Config{}

		Convey("初始化失败", func() {
			conf.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg, errors.New("test"))
			err := aSMS.initSMSExpiration()
			assert.Equal(t, err, errors.New("test"))
		})

		Convey("初始化成功", func() {
			conf.EXPECT().GetConfig(gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg, nil)
			err := aSMS.initSMSExpiration()
			assert.Equal(t, err, nil)
		})
	})
}
