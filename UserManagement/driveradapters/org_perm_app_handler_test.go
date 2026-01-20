// Package driveradapters 应用账户组织管理权限设置测试
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
	"go.uber.org/mock/gomock"

	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

var (
	strUserID = "user_id"
)

func newOrgPermAppRESTHandler(hydra interfaces.Hydra, perm interfaces.LogicsOrgPermApp) *orgPermAppHandler {
	return &orgPermAppHandler{
		hydra:      hydra,
		orgPermApp: perm,
		appPermOrgStrType: map[string]interfaces.OrgType{
			"user":       interfaces.User,
			"department": interfaces.Department,
			"group":      interfaces.Group,
		},
		appPermOrgTypeStr: map[interfaces.OrgType]string{
			interfaces.User:       "user",
			interfaces.Department: "department",
			interfaces.Group:      "group",
		},
		appOrgPermTypeStr: map[interfaces.AppOrgPermValue]string{
			interfaces.Modify: "modify",
			interfaces.Read:   "read",
		},
		appOrgPermStrType: map[string]interfaces.AppOrgPermValue{
			"modify": interfaces.Modify,
			"read":   interfaces.Read,
		},
	}
}

func TestGetAppOrgPerm(t *testing.T) {
	Convey("获取应用账户权限信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		o := mock.NewMockLogicsOrgPermApp(ctrl)
		app := newOrgPermAppRESTHandler(h, o)

		app.RegisterPublic(r)

		const target = "/api/user-management/v1/app-perms"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    false,
			VisitorID: "user_id",
		}

		Convey("token验证失败", func() {
			tempTarget := target + "/xxx/user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo.Active = true
		Convey("objects存在错误的对象", func() {
			tempTarget := target + "/xxx/user1,department"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("objects对象为空", func() {
			tempTarget := target + "/xxx/,department"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取权限信息失败", func() {
			tempTarget := target + "/xxx/user"
			testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().GetAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, testErr.Cause)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("数据获取成功", func() {
			tempTarget := target + "/xxx/user,department"
			infos := make([]interfaces.AppOrgPerm, 0)
			info1 := interfaces.AppOrgPerm{}
			info1.Subject = strUserID
			info1.Object = interfaces.User
			info1.Value = 1
			info2 := interfaces.AppOrgPerm{}
			info2.Subject = strUserID
			info2.Object = interfaces.Department
			info2.Value = 2
			infos = append(infos, info1, info2)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().GetAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(infos, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var outData []interface{}
			err := jsoniter.Unmarshal(respBody, &outData)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(outData), 2)

			tempData1 := outData[0].(map[string]interface{})
			assert.Equal(t, tempData1["subject"].(string), info1.Subject)
			assert.Equal(t, tempData1["object"].(string), app.appPermOrgTypeStr[info1.Object])
			tempPerm1, ok := tempData1["perms"].([]interface{})
			assert.Equal(t, ok, true)
			assert.Equal(t, len(tempPerm1), 1)
			assert.Equal(t, tempPerm1[0].(string), app.appOrgPermTypeStr[interfaces.Modify])

			tempData2 := outData[1].(map[string]interface{})
			assert.Equal(t, tempData2["subject"].(string), info2.Subject)
			assert.Equal(t, tempData2["object"].(string), app.appPermOrgTypeStr[info2.Object])
			tempPerm2, ok := tempData2["perms"].([]interface{})
			assert.Equal(t, ok, true)
			assert.Equal(t, len(tempPerm2), 1)
			assert.Equal(t, tempPerm2[0].(string), app.appOrgPermTypeStr[interfaces.Read])

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:dupl,funlen,gocyclo
func TestUpdateAppOrgPerm(t *testing.T) {
	Convey("更新应用账户权限信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		o := mock.NewMockLogicsOrgPermApp(ctrl)
		app := newOrgPermAppRESTHandler(h, o)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/app-perms"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    false,
			VisitorID: "user_id",
		}

		Convey("token验证失败", func() {
			tempTarget := target + "/xxx/user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo.Active = true
		Convey("url objects存在错误的对象", func() {
			tempTarget := target + "/xxx/user1,department"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url objects为空，设置失败", func() {
			tempTarget := target + "/xxx/,user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url objects存在重复对象", func() {
			tempTarget := target + "/xxx/user,user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "url objects are not uniqued")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内不存在subject", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"object": "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "body[0].subject is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内不存在object", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "body[0].object is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内不存在perms", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "xxx",
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "body[0].perms is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内subject为1，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": 1,
				"object":  "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].subject should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内object为1，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  1,
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].object should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内perms为1，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "xxx",
				"perms":   1,
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].perms should be slice")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内subject为nil，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": nil,
				"object":  "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].subject should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内object为nil，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  nil,
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].object should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内perms为nil，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "xxx",
				"perms":   nil,
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].perms should be slice")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内subject为空，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "",
				"object":  "user",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request subject error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内object为空，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body object error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内perms为空，类型错误", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms":   []string{},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body perms are empty")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("subjects id 不一致，报错", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx1",
				"object":  "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request subject error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request object非枚举，报错", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "xxx",
				"perms": []string{
					"xxxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body object error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request perms非枚举，报错", func() {
			tempTarget := target + "/xxx/user"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"xxx",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body perms error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url object和body objet数量不同，报错", func() {
			tempTarget := target + "/xxx/user,department"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("url object和body objet不一致，报错", func() {
			tempTarget := target + "/xxx/user,department"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			data2 := gin.H{
				"subject": "xxx",
				"object":  "group",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1, data2})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body objects are not same with url objects")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("body重复设置权限，报错", func() {
			tempTarget := target + "/xxx/user,department"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			data2 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			data3 := gin.H{
				"subject": "xxx",
				"object":  "department",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1, data2, data3})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "request body objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SetAppOrgPerm报错，报错", func() {
			tempTarget := target + "/xxx/user,department"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			data2 := gin.H{
				"subject": "xxx",
				"object":  "department",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1, data2})
			testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().SetAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "objects error what err")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			tempTarget := target + "/xxx/user,department"
			data1 := gin.H{
				"subject": "xxx",
				"object":  "user",
				"perms": []string{
					"read",
				},
			}
			data2 := gin.H{
				"subject": "xxx",
				"object":  "department",
				"perms": []string{
					"read",
				},
			}
			jsonData, _ := jsoniter.Marshal([]map[string]interface{}{data1, data2})

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().SetAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteAppOrgPerm(t *testing.T) {
	Convey("删除应用账户权限信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		o := mock.NewMockLogicsOrgPermApp(ctrl)
		app := newOrgPermAppRESTHandler(h, o)

		app.RegisterPublic(r)

		target := "/api/user-management/v1/app-perms"
		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    false,
			VisitorID: "user_id",
		}

		Convey("token验证失败", func() {
			tempTarget := target + "/xxx/user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo.Active = true
		Convey("objects存在错误的对象", func() {
			tempTarget := target + "/xxx/user1"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("objects存在空", func() {
			tempTarget := target + "/xxx/,user"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "objects error")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("删除权限信息失败", func() {
			tempTarget := target + "/xxx/user"
			testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().DeleteAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, testErr.Cause)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("数据删除成功", func() {
			tempTarget := target + "/xxx/user,department"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			o.EXPECT().DeleteAppOrgPerm(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("DELETE", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusNoContent)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
