// Package assertion 逻辑层
package assertion

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"

	"github.com/kweaver-ai/go-lib/rest"

	"Authentication/common"
	"Authentication/interfaces"
	"Authentication/interfaces/mock"
)

func setGinMode() func() {
	old := gin.Mode()
	gin.SetMode(gin.TestMode)
	return func() {
		gin.SetMode(old)
	}
}

const (
	appID  = "b550af01-06d0-446d-be5b-b44cfcd97906"
	userID = "ebb09008-fb85-11ed-9b9f-761dc86fda9f"
	ass    = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjQwMWU2OTYyLWFiNTAtNGM5MC1iZjA3LWY5ZDIzZDRkZWYzYiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly8xMC40LjY4LjY5OjQ0My9vYXV0aDIvdG9rZW4iXSwi" +
		"ZXhwIjoxNjkyOTUxMjE0LCJleHQiOnsiYWNjb3VudF90eXBlIjoib3RoZXIiLCJjbGllbnRfaWQiOiJiNTUwYWYwMS0wNmQwLTQ0NmQtYmU1Yi1iNDRjZmNkOTc5MDYiLCJjbGllbnRfdHlwZSI6ImFwcCIsImxvZ2luX2l" +
		"wIjoiIiwidWRpZCI6IiIsInZpc2l0b3JfdHlwZSI6InJlYWxuYW1lIn0sImlhdCI6MTY5Mjk0NzYxNCwiaXNzIjoiYXV0aGVudGljYXRpb24iLCJqdGkiOiJkYmIzMWIxMi1kODA2LTRjNzItOWY0Mi00ZjFiMGQyZDM1ZWQi" +
		"LCJuYmYiOjE2OTI5NDc2MTQsInN1YiI6ImViYjA5MDA4LWZiODUtMTFlZC05YjlmLTc2MWRjODZmZGE5ZiJ9.IbTsVZbjOtqSpC1odS6r6Nq3Qzes6Z0lHXU-z6AWYOTKJoU8WyNdSoowzwhN3F2X35af3uijEFIf60I33A6Q" +
		"2SAgD25Nn4RWhWqWzA-BhDhJDnkqhvwtW0EieZFCXNiFZhGPSH6d5DRFQuhnF4A9sYb1HlMt7iVfPrvUPPoRWp8YaiN2ABvi3FmAxfwx0YUbTun2TYEuvuS39ge6loToYhUQXluZl6R36LHT8It7jJteHfNIiZtCxpFA4J4Vr" +
		"P1WjjEFiSVXNYaspv1iNtw2fp8aGDWE56or8G2KvSTKPQyhQvcVO5dmScVA9tN8JDEDuwUnclwYNR1SiiAnUWZ-v8F-p9k327cQOi7jOp7fEPhjnlITw5bTfEyuqVv8bcJEr4_n5kC39BPvRSKYOH8t_vICyNqHTBcIqlJOrzT" +
		"qo9cZBWZXD8lO3DgztZasjc1MIUAewPEQA3mqN2PKsFMuJ3_CafCRTpAraZLyw87hRsQ3LTcy25FxAs1oS1G6xKiOIVXHP7hFWQzIJvr-U9aZ_3AC1rdbb-PKnP2W4d6IO1ZnWoSE6iVuvDfnIIVqoPX_fmbVUDE7LL3wR3Ou4Y" +
		"5U3W0djsB5w3NZ2K3Z4WQIbvUM5p1J6cIvOXCjc6bhO3i-ztC9qS2E-4cFDCB7TAS8FffcMzh3dtfMS27PudZ7twA"
)

func newAssertion(userMgnt interfaces.DnUserManagement, hydraAdmin interfaces.DnHydraAdmin, hydraPublic interfaces.DnHydraPublic, access interfaces.AccessTokenPerm) *assertion {
	blockPvt, _ := pem.Decode([]byte(pvtKeyStr))
	privateKey, _ := x509.ParsePKCS1PrivateKey(blockPvt.Bytes)
	return &assertion{
		userMgnt:        userMgnt,
		hydraAdmin:      hydraAdmin,
		hydraPublic:     hydraPublic,
		accessTokenPerm: access,
		log:             common.NewLogger(),
		issuer:          "authentication",
		kid:             "401e6962-ab50-4c90-bf07-f9d23d4def3b",
		privateKey:      privateKey,
	}
}

func TestGetAssertionByUserID(t *testing.T) {
	Convey("TestGetAssertionByUserID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		access := mock.NewMockAccessTokenPerm(ctrl)
		trace := mock.NewMockTraceClient(ctrl)
		as := newAssertion(userMgnt, hydraAdmin, hydraPublic, access)
		as.trace = trace

		visitor := &interfaces.Visitor{ID: appID, Type: interfaces.Business}
		ctx := context.Background()
		testErr := errors.New("some error")

		trace.EXPECT().SetInternalSpanName(gomock.Any()).AnyTimes()
		trace.EXPECT().AddInternalTrace(gomock.Any()).AnyTimes()
		trace.EXPECT().TelemetrySpanEnd(gomock.Any(), gomock.Any()).AnyTimes()
		Convey("访问者不是app类型", func() {
			visitor.Type = interfaces.RealName
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.BadRequest)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Visitor type must be app.")
		})

		Convey("GetUserRolesByUserID失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, testErr)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		Convey("检查访问令牌权限失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, testErr)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		Convey("没有访问令牌权限", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(false, nil)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Forbidden)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "This app doesn't have permission to access token.")
		})

		Convey("获取受信任关系失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			hydraAdmin.EXPECT().GetKidTrustedPairByIssuer(gomock.Any()).Return(map[string]bool{}, testErr)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		Convey("未创建过受信任关系，创建受信任关系失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			hydraAdmin.EXPECT().GetKidTrustedPairByIssuer(gomock.Any()).Return(map[string]bool{}, nil)
			hydraAdmin.EXPECT().CreateTrustedPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		Convey("生成断言，获取tokenEndpoint失败", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			hydraAdmin.EXPECT().GetKidTrustedPairByIssuer(gomock.Any()).Return(map[string]bool{}, nil)
			hydraAdmin.EXPECT().CreateTrustedPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			hydraPublic.EXPECT().GetTokenEndpoint().Return("", testErr)
			_, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, testErr)
		})

		Convey("生成断言", func() {
			userMgnt.EXPECT().GetUserRolesByUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.RoleType{}, nil)
			access.EXPECT().CheckAppAccessTokenPerm(gomock.Any()).Return(true, nil)
			hydraAdmin.EXPECT().GetKidTrustedPairByIssuer(gomock.Any()).Return(map[string]bool{}, nil)
			hydraAdmin.EXPECT().CreateTrustedPair(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			hydraPublic.EXPECT().GetTokenEndpoint().Return("", nil)
			assertion, err := as.GetAssertionByUserID(ctx, visitor, userID)
			assert.Equal(t, err, nil)
			assert.NotEqual(t, assertion, nil)
		})
	})
}

func TestTokenHook(t *testing.T) {
	Convey("TestTokenHook", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userMgnt := mock.NewMockDnUserManagement(ctrl)
		hydraAdmin := mock.NewMockDnHydraAdmin(ctrl)
		hydraPublic := mock.NewMockDnHydraPublic(ctrl)
		access := mock.NewMockAccessTokenPerm(ctrl)
		as := newAssertion(userMgnt, hydraAdmin, hydraPublic, access)

		Convey("生成断言与使用断言的app不一致", func() {
			_, err := as.TokenHook(ass, userID)
			assert.Equal(t, err.(*rest.HTTPError).Code, rest.Forbidden)
			assert.Equal(t, err.(*rest.HTTPError).Cause, "Invalid client_id.")
		})

		Convey("解析成功", func() {
			res, err := as.TokenHook(ass, appID)
			assert.Equal(t, err, nil)
			assert.Equal(t, res["visitor_type"].(string), "realname")
			assert.Equal(t, res["login_ip"].(string), "")
			assert.Equal(t, res["account_type"].(string), "other")
			assert.Equal(t, res["udid"].(string), "")
			assert.Equal(t, res["client_type"].(string), "app")
		})
	})
}
