// Package main
package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"

	"Authentication/common"
	"Authentication/dbaccess"
	"Authentication/drivenadapters"
	"Authentication/driveradapters/conf"
	"Authentication/interfaces"
	amock "Authentication/interfaces/mock"
	"Authentication/logics"
)

var (
	mockOnce      sync.Once
	mock          sqlmock.Sqlmock
	chandler      conf.RESTHandler
	enginePrivate *gin.Engine
	enginePublic  *gin.Engine

	hydraHandler    interfaces.DepHTTPSvc
	userMgntHandler interfaces.DepHTTPSvc
)

func TestNewAuthentication(t *testing.T) {
	mockOnce.Do(func() {
		// mock 出站适配器
		TestInitUserMgntService(t)
		TestInitHydraService(t)

		// 创建 sqlmock 数据库连接和mock
		var db *sqlx.DB
		db, mock, _ = sqlx.New()

		// dbPool注入
		dbaccess.SetDBPool(db)

		// dbaccess 依赖注入
		logics.SetDBConf(dbaccess.NewConf())
		logics.SetDBPool(db)

		common.InitARTrace("test")
		// drivenadapters 依赖注入
		logics.SetDnUserManagement(drivenadapters.NewUserManagement())

		chandler = conf.NewRESTHandler()

		gin.SetMode(gin.TestMode)

		// 注册内部API
		enginePrivate = gin.New()
		enginePrivate.Use(gin.Recovery())
		chandler.RegisterPrivate(enginePrivate)
		// 注册外部API
		enginePublic = gin.New()
		enginePublic.Use(gin.Recovery())
		chandler.RegisterPublic(enginePublic)
	})
}

func handle(t *testing.T, w http.ResponseWriter, r *http.Request, handler interfaces.DepHTTPSvc) {
	target := r.URL.EscapedPath()
	method := r.Method
	body, err := io.ReadAll(r.Body)
	assert.Equal(t, err, nil)
	var reqBody map[string]interface{}
	if method == "POST" {
		reqBody = make(map[string]interface{})
		err = jsoniter.Unmarshal(body, &reqBody)
		assert.Equal(t, err, nil)
	}
	resCode, resBody := handler.HandleRequest(method, target, reqBody)
	w.WriteHeader(resCode)
	_, err = w.Write(resBody)
	assert.Equal(t, err, nil)
}

func TestInitUserMgntService(t *testing.T) {
	// http mock,模拟UserMgnt服务
	tsUserMgnt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(t, w, r, userMgntHandler)
	}))

	// 解析服务地址
	u, err := url.Parse(tsUserMgnt.URL)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// 配置注入
	port, _ := strconv.Atoi(u.Port())
	common.SvcConfig.UserManagementPrivateHost = u.Hostname()
	common.SvcConfig.UserManagementPrivatePort = port
}

func TestInitHydraService(t *testing.T) {
	// http mock,模拟Hydra服务
	tsHydra := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.EscapedPath()
		method := r.Method
		body, _ := io.ReadAll(r.Body)
		resCode, resBody := hydraHandler.HandleRequest(method, url, string(body))
		w.WriteHeader(resCode)
		_, err := w.Write(resBody)
		assert.Equal(t, err, nil)
	}))

	// 解析服务地址
	h, err := url.Parse(tsHydra.URL)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// 配置注入
	port, _ := strconv.Atoi(h.Port())
	common.SvcConfig.OAuthAdminHost = h.Hostname()
	common.SvcConfig.OAuthAdminPort = port
}

func controllerInjection(ctrl *gomock.Controller) (hydraMockObj, userMgntMockObj *amock.MockDepHTTPSvc) {
	hydraMockObj = amock.NewMockDepHTTPSvc(ctrl)
	hydraHandler = hydraMockObj

	userMgntMockObj = amock.NewMockDepHTTPSvc(ctrl)
	userMgntHandler = userMgntMockObj

	return
}

func getHTTPSvcResInfo() (httpSvcResMap map[string]interface{}) {
	httpSvcResMap = make(map[string]interface{})
	tokenInfo, _ := jsoniter.Marshal(gin.H{
		"active":    true,
		"client_id": "xxx",
		"sub":       "a5d47ec5-231f-35f5-1111-9194b66134a5",
		"scope":     "xxx",
		"ext": gin.H{
			"visitor_type": "realname",
			"login_ip":     "xx.xx.xx.xx",
			"udid":         "xxx",
			"account_type": "other",
			"client_type":  "web",
		},
	})
	httpSvcResMap["tokenInfo"] = tokenInfo

	adminRoleType, _ := jsoniter.Marshal([]map[string]interface{}{
		gin.H{
			"roles": []interface{}{"super_admin"},
		},
	})
	httpSvcResMap["adminRoleType"] = adminRoleType

	return
}

func TestSetConf(t *testing.T) {
	TestNewAuthentication(t)

	resInfo := getHTTPSvcResInfo()
	Convey("set conf", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hydraMock, _ := controllerInjection(ctrl)

		// 允许以任意顺序匹配期望
		mock.MatchExpectationsInOrder(false)

		Convey("remember_for 为int64最大值+1，抛错", func() {
			reqParam := map[string]interface{}{
				"remember_for": uint64(9223372036854775807 + 1),
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)
			hydraMock.EXPECT().HandleRequest("POST", "/admin/oauth2/introspect", "token=this-is-test-token").Return(http.StatusOK, resInfo["tokenInfo"])
			req := httptest.NewRequest("PUT", "/api/authentication/v1/config/remember_for", bytes.NewReader(reqParamByte))
			req.Header.Add("Authorization", "Bearer this-is-test-token")
			w := httptest.NewRecorder()
			enginePublic.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
