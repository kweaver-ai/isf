package driveradapters

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/mock/gomock"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

type RegisterParam struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type,omitempty"`
}

type UpdateParam struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

const (
	appTarget = "/api/user-management/v1/apps"
)

func newAppRESTHandler(hydra interfaces.Hydra, app interfaces.LogicsApp, user interfaces.LogicsUser) AppRestHandler {
	appTokenGenerateSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(appTokenGenerateSchemaStr))

	return &appRestHandler{
		hydra: hydra,
		app:   app,
		user:  user,
		appType: map[string]interfaces.AppType{
			"general":   interfaces.General,
			"specified": interfaces.Specified,
			"internal":  interfaces.Internal,
		},
		credentialType: map[interfaces.CredentialType]string{
			interfaces.CredentialTypePassword: "password",
			interfaces.CredentialTypeToken:    "token",
		},
		appTokenGenerateSchema: appTokenGenerateSchema,
	}
}

func mockRequest(needBody bool, method, target string, body io.Reader, r http.Handler) (result *http.Response, resBody map[string]interface{}) {
	req := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	result = w.Result()

	if needBody {
		message, _ := io.ReadAll(result.Body)
		_ = jsoniter.Unmarshal(message, &resBody)
	}

	return
}

func TestGeneralAppRegisterParam1(t *testing.T) {
	Convey("通用账号注册", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPublic(r)

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("不传password，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"name": "test",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("不传name，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"password": "some-secret",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("注册通用应用账户成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := RegisterParam{
				Name:     "(σﾟ∀ﾟ)σ..☆哎哟不错哦❤haha666",
				Password: "some-secret",
			}

			a.EXPECT().RegisterApp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("this-is-id", nil)
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, resBody := mockRequest(true, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, resBody["id"].(string), "this-is-id")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("注册通用应用账户成功1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"name":     "(σﾟ∀ﾟ)σ..☆哎哟不错哦❤haha666",
				"password": "",
			}

			a.EXPECT().RegisterApp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("this-is-id", nil)
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, resBody := mockRequest(true, "POST", appTarget, bytes.NewReader(reqParamByte), r)
			assert.Equal(t, result.StatusCode, http.StatusCreated)

			assert.Equal(t, resBody["id"].(string), "this-is-id")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSpecifiedAppRegisterParam3(t *testing.T) {
	Convey("专用账户注册", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPrivate(r)

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("不传name，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"password": "aaaaaaaaaaaaaaaa",
				"type":     "internal",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("不传password，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"name": "aaaaaaaaaaaaaaaa",
				"type": "internal",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type传internal和specified字段以外的字段，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := RegisterParam{
				Name:     "aaaaaaaaaaaaaaaa",
				Password: "some-secret",
				Type:     "specifiedddddddd",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type传空字符串/None，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := RegisterParam{
				Name:     "aaaaaaaaaaaaaaaa",
				Password: "some-secret",
				Type:     "",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type不传，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := RegisterParam{
				Name:     "aaaaaaaaaaaaaaaa",
				Password: "some-secret",
				Type:     "",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			result, _ := mockRequest(false, "POST", appTarget, bytes.NewReader(reqParamByte), r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteApp(t *testing.T) {
	Convey("删除应用账户", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/apps/test-id"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("删除应用账户，删除成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			a.EXPECT().DeleteApp(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			result, _ := mockRequest(true, "DELETE", target, nil, r)
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUpdateApp(t *testing.T) {
	Convey("更新应用账户", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/apps/test-id/"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("fileds传name和password以外的字段，抛错 ", func() {
			reqParam := UpdateParam{
				Name: "test",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			result, _ := mockRequest(false, "PUT", target+"test", bytes.NewReader(reqParamByte), r)
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		reqParam := UpdateParam{
			Name:     "test",
			Password: "testtest",
		}
		reqParamByte, _ := jsoniter.Marshal(reqParam)

		Convey("更新者为admin管理员/超级管理员，更新成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			a.EXPECT().UpdateApp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			result, _ := mockRequest(false, "PUT", target+"name,password", bytes.NewReader(reqParamByte), r)
			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAppList(t *testing.T) {
	Convey("获取应用账户列表", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/apps"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("获取成功", func() {
			info := &[]interfaces.AppInfo{
				{
					ID:             "test1",
					Name:           "test1",
					CredentialType: interfaces.CredentialTypePassword,
				},
				{
					ID:             "test2",
					Name:           "test2",
					CredentialType: interfaces.CredentialTypeToken,
				},
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			a.EXPECT().AppList(gomock.Any(), gomock.Any()).AnyTimes().Return(info, 2, nil)
			result, _ := mockRequest(false, "GET", target, nil, r)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("sort传字符串以外的类型，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?sort=1", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("direction传asc、desc以外的字段，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?direction=ascasdf", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset传整形以外的类型，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?offset=asdfasd", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset小于0，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?offset=-5", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit传整形以外的类型，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?limit=asdfasdf", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit小于1，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?limit=-5", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit大于1000，抛错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			result, _ := mockRequest(false, "GET", target+"?limit=9999", nil, r)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetApp(t *testing.T) {
	Convey("test get app", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPrivate(r)

		target := "/api/user-management/v1/apps/aaaaaa"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("success", func() {
			info := &interfaces.AppInfo{
				ID:   "test1",
				Name: "test1",
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			a.EXPECT().GetApp(gomock.Any()).AnyTimes().Return(info, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGenerateAppToken(t *testing.T) {
	Convey("test generate app token", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		common.InitARTrace("test")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		a := mock.NewMockLogicsApp(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		app := newAppRESTHandler(h, a, u)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/console/app-tokens"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    false,
			VisitorID: "user_id",
		}

		Convey("hydra token introspect failed", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo.Active = true
		Convey("request body 不包含id", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("requst body id 不为string", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := map[string]interface{}{
				"id": 123,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GenerateAppToken error", func() {
			reqParam := map[string]interface{}{
				"id": "test-id",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			testErr := rest.NewHTTPError("error", 503000000, nil)
			a.EXPECT().GenerateAppToken(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("", testErr)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GenerateAppToken success", func() {
			reqParam := map[string]interface{}{
				"id": "test-id",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			a.EXPECT().GenerateAppToken(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("token1", nil)
			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			body, _ := io.ReadAll(result.Body)
			msg := make(map[string]interface{})
			_ = jsoniter.Unmarshal(body, &msg)
			assert.Equal(t, msg["token"].(string), "token1")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
