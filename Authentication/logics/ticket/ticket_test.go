package ticket

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

func newLoTicket(dbTicket interfaces.DBTicket, hydraAdmin interfaces.DnHydraAdmin, userMgnt interfaces.DnUserManagement) *ticket {
	// pem 解码
	blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
	// X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
	if err != nil {
		common.NewLogger().Fatalln(err)
	}
	return &ticket{
		dbTicket:         dbTicket,
		userMgnt:         userMgnt,
		hydraAdmin:       hydraAdmin,
		privateKey:       privateKey,
		ticketExpiration: time.Minute * 5,
		logger:           common.NewLogger(),
	}
}

func TestNewTicket(t *testing.T) {
	Convey("Validate", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tic := NewTicket()
		assert.NotEqual(t, tic, nil)
	})
}

//nolint:lll
func TestCreateTicket(t *testing.T) {
	Convey("CreateTicket", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbTicket := mock.NewMockDBTicket(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		loTicket := newLoTicket(dbTicket, hydraAdmin, userMgnt)
		trace := mock.NewMockTraceClient(ctrl)
		loTicket.trace = trace

		reqInfo := &interfaces.TicketReq{
			ClientID:     "0c8839f4-894c-452c-ae96-911df5e04c64",
			RefreshToken: "TgtQys63IBFOvU7kbxVAldtWfg7mKlgkWHV7PYqev3iVKM348oCRLbSUTQOhejECOYPKXQx5ZZLrJWqPcz/r9wziPRmPeKYGS6HdjmUx2s1obj6fQ3X+xccgSpacMFP+jBIxWkz8iSY5+LfJbuPwFJ8z1YvxvFYuzPENYlgm0MgpY1Dyxf74DWlDd09QmRPdvH8WPzXjuMe5Zp07zWE+vjW0Er958QM9B8efjkuEHRdUWYZg/VR2U/2X9rphIG/IyazlstTo5TQk9Nj2yEpUiUHxk9Kl8I14Sye0WFi4Yp3TcFS1G/QVWPKHBLL+j5XI2Sb8qH2ejFlhLbaqRWGuVg==",
		}
		refreshTokenIntrospectInfo := &interfaces.RefreshTokenIntrospectInfo{
			ClientID: "7c00d61e-bba0-4bd5-b18c-1853d83e01c3",
			Sub:      "dfc9b098-dac4-11ee-b50a-028586548cf7",
		}

		ctx := context.Background()
		visitor := interfaces.Visitor{
			ErrorCodeType: interfaces.Number,
		}
		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("DecodeAndDecrypt error", func() {
			reqInfo.RefreshToken = "MTExCg=="

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "crypto/rsa: decryption error"))
		})

		Convey("DecodeAndDecrypt error2", func() {
			reqInfo.RefreshToken = "111"

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "illegal base64 data at input byte 0"))
		})

		Convey("IntrospectRefreshToken, token is not active", func() {
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(nil, rest.NewHTTPErrorV2(rest.BadRequest, "token is not active"))

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "token is not active"))
		})

		Convey("IntrospectRefreshToken, invalid token use", func() {
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(nil, rest.NewHTTPErrorV2(rest.BadRequest, "Invalid token use"))

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "Invalid token use"))
		})

		Convey("CreateTicket, invalid clientID", func() {
			reqInfo.ClientID = "7c00d61e-bba0-4bd5-b18c-1853d83e01c3"
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(refreshTokenIntrospectInfo, nil)

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "The client cannot apply for single sign-on credentials for itself"))
		})

		Convey("CreateTicket, GetUserInfo failed", func() {
			refreshTokenIntrospectInfo.Sub = "xxx"
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(refreshTokenIntrospectInfo, nil)
			userMgnt.EXPECT().GetUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, rest.NewHTTPErrorV2(rest.URINotExist, "not found"))

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.URINotExist, "not found"))
		})

		Convey("CreateTicket, GetClientInfo failed", func() {
			reqInfo.ClientID = "xxx"
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(refreshTokenIntrospectInfo, nil)
			userMgnt.EXPECT().GetUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			hydraAdmin.EXPECT().GetClientInfo(gomock.Any()).Return(nil, rest.NewHTTPErrorV2(rest.URINotExist, "not found"))

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.URINotExist, "not found"))
		})

		Convey("CreateTicket, db unavailable", func() {
			tmpErr := errors.New("unknown error")
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(refreshTokenIntrospectInfo, nil)
			userMgnt.EXPECT().GetUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			hydraAdmin.EXPECT().GetClientInfo(gomock.Any()).Return(nil, nil)
			dbTicket.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tmpErr)

			_, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.Equal(t, err, tmpErr)
		})

		Convey("CreateTicket, success", func() {
			hydraAdmin.EXPECT().IntrospectRefreshToken(gomock.Any()).Return(refreshTokenIntrospectInfo, nil)
			userMgnt.EXPECT().GetUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			hydraAdmin.EXPECT().GetClientInfo(gomock.Any()).Return(nil, nil)
			dbTicket.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

			ticketID, err := loTicket.CreateTicket(ctx, &visitor, reqInfo)

			assert.NotEqual(t, ticketID, "")
			assert.Equal(t, err, nil)
		})
	})
}

func TestValidate(t *testing.T) {
	Convey("Validate", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbTicket := mock.NewMockDBTicket(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		userMgnt := mock.NewMockDnUserManagement(ctrl)
		loTicket := newLoTicket(dbTicket, hydraAdmin, userMgnt)
		trace := mock.NewMockTraceClient(ctrl)
		loTicket.trace = trace

		ticketID := "01HVQWAXXCSA98YT2BY6Q0KKRS"
		clientID := "0c8839f4-894c-452c-ae96-911df5e04c64"
		tickInfo := &interfaces.TicketInfo{
			ID:         "01HVQWAXXCSA98YT2BY6Q0KKRS",
			UserID:     "dfc9b098-dac4-11ee-b50a-028586548cf7",
			ClientID:   "0c8839f4-894c-452c-ae96-911df5e04c64",
			CreateTime: time.Now().Add(-time.Minute * 2).Unix(),
		}

		ctx := context.Background()
		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("GetTicketByID error", func() {
			dbTicket.EXPECT().GetTicketByID(gomock.Any(), gomock.Any()).Return(nil, rest.NewHTTPErrorV2(rest.Unauthorized, "invalid ticket"))

			_, err := loTicket.Validate(ctx, ticketID, clientID)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.Unauthorized, "invalid ticket"))
		})

		Convey("failed, Invalid clientID", func() {
			tickInfo.ClientID = "other client id"
			dbTicket.EXPECT().GetTicketByID(gomock.Any(), gomock.Any()).Return(tickInfo, nil)

			_, err := loTicket.Validate(ctx, ticketID, clientID)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "Invalid clientID"))
		})

		Convey("failed, ticket expired", func() {
			tickInfo.CreateTime = time.Now().Add(-time.Minute * 10).Unix()
			dbTicket.EXPECT().GetTicketByID(gomock.Any(), gomock.Any()).Return(tickInfo, nil)

			_, err := loTicket.Validate(ctx, ticketID, clientID)

			assert.Equal(t, err, rest.NewHTTPErrorV2(rest.BadRequest, "Ticket expired"))
		})

		Convey("failed, db unavailable", func() {
			tmpErr := errors.New("unknow error")
			dbTicket.EXPECT().GetTicketByID(gomock.Any(), gomock.Any()).Return(tickInfo, nil)
			dbTicket.EXPECT().DeleteByIDs(gomock.Any(), gomock.Any()).Return(tmpErr)

			_, err := loTicket.Validate(ctx, ticketID, clientID)

			assert.Equal(t, err, tmpErr)
		})

		Convey("success", func() {
			dbTicket.EXPECT().GetTicketByID(gomock.Any(), gomock.Any()).Return(tickInfo, nil)
			dbTicket.EXPECT().DeleteByIDs(gomock.Any(), gomock.Any()).Return(nil)

			userID, err := loTicket.Validate(ctx, ticketID, clientID)

			assert.Equal(t, userID, tickInfo.UserID)
			assert.Equal(t, err, nil)
		})
	})
}
