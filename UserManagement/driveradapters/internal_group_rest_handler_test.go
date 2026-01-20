// Package driveradapters group AnyShare  内部组逻辑接口处理层
//
//nolint:funlen
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

const (
	strID = "adasasddsad"
)

func newInternalGroupRestHandler(g interfaces.LogicsInternalGroup) *internalGroupRestHandler {
	return &internalGroupRestHandler{
		group: g,
		orgTypeToString: map[interfaces.OrgType]string{
			interfaces.User: "user",
		},
		stringToOrgType: map[string]interfaces.OrgType{
			"user": interfaces.User,
		},
	}
}

func TestInternalGroupCreateGroup(t *testing.T) {
	Convey("创建内部组", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := mock.NewMockLogicsInternalGroup(ctrl)
		group := newInternalGroupRestHandler(g)

		group.RegisterPrivate(r)

		const target = "/api/user-management/v1/internal-groups"
		testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)
		Convey("AddGroup报错，失败", func() {
			g.EXPECT().AddGroup().AnyTimes().Return(strID, testErr)

			req := httptest.NewRequest("POST", target, http.NoBody)
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

		Convey("成功", func() {
			g.EXPECT().AddGroup().AnyTimes().Return(strID, nil)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusCreated)
			assert.Equal(t, result.Header["Location"][0], "/api/user-management/v1/internal-groups/"+strID)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestInternalGroupDeleteGroup(t *testing.T) {
	Convey("删除内部组", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := mock.NewMockLogicsInternalGroup(ctrl)
		group := newInternalGroupRestHandler(g)

		group.RegisterPrivate(r)

		const target = "/api/user-management/v1/internal-groups"
		testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)
		Convey("DeleteGroup", func() {
			tempTarget := target + "/xxxx,xxxx"
			g.EXPECT().DeleteGroup(gomock.Any()).AnyTimes().Return(testErr)

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

		Convey("成功", func() {
			tempTarget := target + "/xxxx,xxxx"
			g.EXPECT().DeleteGroup(gomock.Any()).AnyTimes().Return(nil)

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

func TestInternalGetMembersByID(t *testing.T) {
	Convey("获取内部组成员信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := mock.NewMockLogicsInternalGroup(ctrl)
		group := newInternalGroupRestHandler(g)

		group.RegisterPrivate(r)

		const target = "/api/user-management/v1/internal-group-members"
		testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)
		Convey("DeleteGroup", func() {
			tempTarget := target + "/xxxx"
			g.EXPECT().GetGroupMemberByID(gomock.Any()).AnyTimes().Return(nil, testErr)

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

		out := make([]interfaces.InternalGroupMember, 0)
		data1 := interfaces.InternalGroupMember{
			ID:   strID,
			Type: interfaces.User,
		}
		data2 := interfaces.InternalGroupMember{
			ID:   strID + "1",
			Type: interfaces.User,
		}
		out = append(out, data1, data2)
		Convey("成功", func() {
			tempTarget := target + "/xxxx"
			g.EXPECT().GetGroupMemberByID(gomock.Any()).AnyTimes().Return(out, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(respParam), 2)

			temp1 := respParam[0].(map[string]interface{})
			assert.Equal(t, temp1["id"], data1.ID)
			assert.Equal(t, temp1["type"], "user")

			temp2 := respParam[1].(map[string]interface{})
			assert.Equal(t, temp2["id"], data2.ID)
			assert.Equal(t, temp2["type"], "user")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestInternalUpdateMembers(t *testing.T) {
	Convey("更新内部组成员", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := mock.NewMockLogicsInternalGroup(ctrl)
		group := newInternalGroupRestHandler(g)

		group.RegisterPrivate(r)

		const target = "/api/user-management/v1/internal-group-members"
		testErr := rest.NewHTTPError("objects error what err", rest.Forbidden, nil)
		Convey("request body 非json", func() {
			tempTarget := target + "/xxxx"
			jsonData, _ := jsoniter.Marshal("xxxx")

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member没有id", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"type": "user",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "body[0].id is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member内id为int", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"type": "user",
				"id":   1,
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].id should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member没有type", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id": "user",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "body[0].type is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member内type为int", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id":   "user",
				"type": 1,
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "type of body[0].type should be string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member内id为空", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id":   "",
				"type": "user",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "param member id is illegal")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("request body内member内type非法", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id":   strID,
				"type": "xxxxx",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "param member type is illegal")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("UpdateMembers报错", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id":   strID,
				"type": "user",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			g.EXPECT().UpdateMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)
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

		Convey("success", func() {
			tempTarget := target + "/xxxx"
			data1 := gin.H{
				"id":   strID,
				"type": "user",
			}
			jsonData, _ := jsoniter.Marshal([]interface{}{data1})

			g.EXPECT().UpdateMembers(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
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
