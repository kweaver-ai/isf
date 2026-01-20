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

type CreateGroupNoNameParams struct {
	Notes string `json:"notes" binding:"required"`
}

type CreateGroupNoNotesParams struct {
	Name string `json:"name" binding:"required"`
}

type CreateGroupParams struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes" binding:"required"`
}

type MembersParams struct {
	ID         string `json:"id" binding:"required"`
	MemberType string `json:"type" binding:"required"`
}

type AddOrDeleteMembersParams struct {
	Method  string          `json:"method" binding:"required"`
	Members []MembersParams `json:"members" binding:"required"`
}

type GetMembersIDParam struct {
	Method   string   `json:"method" binding:"required"`
	GroupIDs []string `json:"group_ids" binding:"required"`
}

func newGroupRESTHandler(groupLogics interfaces.LogicsGroup, userLogics interfaces.LogicsUser, com interfaces.LogicsCombine, hydraClient interfaces.Hydra) GroupRestHandler {
	getGroupMembersSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getGroupMembersSchemaStr))
	if err != nil {
		panic(err)
	}

	createGroupSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(createGroupSchemaStr))
	if err != nil {
		panic(err)
	}

	return &groupRestHandler{
		hydra:   hydraClient,
		group:   groupLogics,
		user:    userLogics,
		combine: com,
		memberStringTypes: map[string]int{
			"user":       1,
			"department": 2,
		},
		memberIntTypes: map[int]string{
			1: "user",
			2: "department",
		},
		getGroupMembersSchema: getGroupMembersSchema,
		createGroupSchema:     createGroupSchema,
	}
}

func TestCreateGroupSuccess(t *testing.T) {
	Convey("CreateGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("user-management")

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups"

		nTempNotes := "xxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxy" +
			"xxxxxxxxxyxxxx😃xxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxy" +
			"xxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxy" +
			"xxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxy" +
			"xxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxxy" +
			"xxxxxxxxxyxxxxxxxxxyxxxxxxxxxyxxxxxxxxx哈xxxxxxxxxy"

		strTempName := "xxxxx😃xxx哈xxxxxxxxx哈xxxxxxxxx哈xxxxxxxxx哈" +
			"xxxxxxxxx哈xxxxxxxxx哈xxxxxxxxx哈xxxxxxxxx哈" +
			"xxxxxxxxx哈xxxxxxxx哈xxxxxxxxx哈xxxxxxxxx哈" +
			"xxxxxxxx"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
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

		Convey("AddGroup is err", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			testErr := rest.NewHTTPError("errorx", 409019000, nil)
			g.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("xxx", testErr)

			reqParam := CreateGroupParams{
				Name:  strTempName,
				Notes: nTempNotes,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().AddGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return("xxx", nil)

			reqParam := CreateGroupParams{
				Name:  strTempName,
				Notes: nTempNotes,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.Header["Location"][0], "/api/user-management/v1/management/groups/xxx")
			assert.Equal(t, result.StatusCode, http.StatusCreated)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteGroup(t *testing.T) {
	Convey("DeleteGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups/xxxxx"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("DeleteGroup is err", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			testErr := rest.NewHTTPError("errorx", 404019000, nil)
			g.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNotFound)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
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

func TestGetGroupByID(t *testing.T) {
	Convey("GetGroupByID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups/xxxxx"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupByID is err", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			info := interfaces.GroupInfo{
				ID:   "xxx",
				Name: "xxxyy",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userInfo := make(map[interfaces.Role]bool)
			userInfo[interfaces.SystemRoleSuperAdmin] = true
			testErr := rest.NewHTTPError("errorx", 404019000, nil)
			g.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).AnyTimes().Return(info, testErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNotFound)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			info := interfaces.GroupInfo{
				ID:   "xxx",
				Name: "xxxyy",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userInfo := make(map[interfaces.Role]bool)
			userInfo[interfaces.SystemRoleSuperAdmin] = true
			g.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).AnyTimes().Return(info, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := interfaces.GroupInfo{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, info)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestModifyGroupParamsFail(t *testing.T) {
	Convey("modifyGroupParamsFail", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups/yyyyy/"

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("Name is int", func() {
			tempTarget := target + "notes,name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type CreateGroupParamsNameInt struct {
				Name  int    `json:"name" binding:"required"`
				Notes string `json:"notes" binding:"required"`
			}

			reqParam := CreateGroupParamsNameInt{
				Name: 1,
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"x",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("notes is int", func() {
			tempTarget := target + "notes,name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type CreateGroupParamsNotesInt struct {
				Name  string `json:"name" binding:"required"`
				Notes int    `json:"notes" binding:"required"`
			}

			reqParam := CreateGroupParamsNotesInt{
				Name:  "xxxxxxx",
				Notes: 1,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestModifyGroupFail(t *testing.T) {
	Convey("modifyGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups/xxxx/"

		Convey("token expired", func() {
			tempTarget := target + "xxxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("no name param", func() {
			tempTarget := target + "name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := CreateGroupNoNameParams{
				Notes: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("no notes param", func() {
			tempTarget := target + "notes"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := CreateGroupNoNotesParams{
				Name: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("no notes or name param", func() {
			tempTarget := target + "notes,name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := CreateGroupNoNotesParams{
				Name: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("notes is illegal", func() {
			tempTarget := target + "notes,name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type CreateGroupParamsNameInt struct {
				Name  int    `json:"name" binding:"required"`
				Notes string `json:"notes" binding:"required"`
			}

			reqParam := CreateGroupParamsNameInt{
				Name: 1,
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"x",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestModifyGroup(t *testing.T) {
	Convey("modifyGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups/xxxx/"

		Convey("ModifyGroup err ", func() {
			tempTarget := target + "notes,name"
			testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().ModifyGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			reqParam := CreateGroupParams{
				Name: "xxxx",
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 409019000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			tempTarget := target + "notes,name"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().ModifyGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			reqParam := CreateGroupParams{
				Name: "xxxx",
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
					"xxxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("PUT", tempTarget, bytes.NewReader(reqParamByte))
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

func TestGetGroupParamsFail(t *testing.T) {
	Convey("GetGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups?"

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("GetGroup Direction err ", func() {
			tmpTarget := target + "direction=xxx"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroup sort err ", func() {
			tmpTarget := target + "sort=xxx"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroup offset err 2", func() {
			tmpTarget := target + "offset=-1"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroup limit err 1", func() {
			tmpTarget := target + "limit=0"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroup limit err 2", func() {
			tmpTarget := target + "limit=1001"

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			_ = jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetGroup(t *testing.T) {
	Convey("GetGroup", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/groups?"

		Convey("token expired", func() {
			tmpTarget := target + "sort=desc"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("GetGroup err ", func() {
			var outInfo []interfaces.GroupInfo
			tmpTarget := target + "sort=date_created"
			testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroup(gomock.Any(), gomock.Any()).AnyTimes().Return(0, outInfo, testErr)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 409019000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroup Success ", func() {
			var outInfo []interfaces.GroupInfo
			outListInfo := ListInfo{
				Entries: make([]interface{}, 0),
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroup(gomock.Any(), gomock.Any()).AnyTimes().Return(0, outInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var respParam ListInfo
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, outListInfo)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetGroupMembersParams(t *testing.T) {
	Convey("GetGroupMembers", t, func() {
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)
		common.InitARTrace("user-management")

		testGRestHandler.RegisterPublic(r)
		const xxx string = "/api/user-management/v1/management/group-members/"
		target := xxx

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("GetGroupMembers sort err ", func() {
			tmpTarget := target + "xxxx?sort=xxxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers direction err ", func() {
			tmpTarget := target + "xxxx?direction=xxxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)
			assert.Equal(t, 3, 3)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers limit err 1", func() {
			tmpTarget := target + "xxxx?limit=xxxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, 2, 2)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers limit err 2", func() {
			tmpTarget := target + "xxxx?limit=0"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam.Code, 400000000)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers limit err 3", func() {
			tmpTarget := target + "xxxx?limit=1001"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers offset err ", func() {
			tmpTarget := target + "xxxx?offset=xxxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, respParam.Code, 400000000)
			assert.Equal(t, nil, err)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembers offset err 1", func() {
			tmpTarget := target + "xxxx?offset=-1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 400000000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetGroupMembers(t *testing.T) {
	Convey("GetGroupMembers", t, func() {
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)
		common.InitARTrace("user-management")

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/group-members/"

		Convey("token expired", func() {
			tmpTarget := target + "xxxxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("GetGroupMembers err ", func() {
			var outInfo []interfaces.GroupMemberInfo
			tmpTarget := target + "xxxx?sort=name"
			testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(0, outInfo, testErr)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := rest.HTTPError{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 409019000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetGroupMembersSuccess(t *testing.T) {
	Convey("GetGroupMembers", t, func() {
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)
		common.InitARTrace("user-management")

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/group-members/"

		Convey("token expired", func() {
			tmpTarget := target + "xxxxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		tempNam1 := interfaces.NameInfo{
			ID:   strUserID,
			Name: strAccount,
		}
		temp1Namel1s := []interfaces.NameInfo{tempNam1}
		data1 := interfaces.GroupMemberInfo{
			ID:              strUserID,
			Name:            strAccount,
			MemberType:      1,
			DepartmentNames: []string{strUserID, strAccount},
			ParentDeps:      [][]interfaces.NameInfo{temp1Namel1s},
		}

		Convey("GetGroupMembers Success ", func() {
			tmpTarget := target + "xxxx?sort=name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroupMembers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(100, []interfaces.GroupMemberInfo{data1}, nil)

			req := httptest.NewRequest("GET", tmpTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var respParam map[string]interface{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam["total_count"], float64(100))
			outdata := respParam["entries"].([]interface{})
			assert.Equal(t, len(outdata), 1)
			outdata1 := outdata[0].(map[string]interface{})
			assert.Equal(t, outdata1["id"], data1.ID)
			assert.Equal(t, outdata1["type"], "user")
			assert.Equal(t, outdata1["name"], data1.Name)
			outtemp1 := outdata1["department_names"].([]interface{})
			assert.Equal(t, outtemp1[0].(string), data1.DepartmentNames[0])
			assert.Equal(t, outtemp1[1].(string), data1.DepartmentNames[1])
			outdata2 := outdata1["parent_deps"].([]interface{})
			assert.Equal(t, len(outdata2), 1)
			otdata3 := outdata2[0].([]interface{})
			assert.Equal(t, len(otdata3), 1)
			otdata4 := otdata3[0].(map[string]interface{})
			otdata4["id"] = tempNam1.ID
			otdata4["name"] = tempNam1.Name
			otdata4["type"] = "department"

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAddAndDeleteGroupMembersParamsFail(t *testing.T) {
	Convey("AddAndDeleteGroupMembersParamsFail", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		const xxx string = "/api/user-management/v1/management/group-members/xxxx"
		target := xxx

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("method param ileagal", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type AddOrDeleteMembersParamsInt struct {
				Method  int             `json:"method" binding:"required"`
				Members []MembersParams `json:"members" binding:"required"`
			}

			reqParam := AddOrDeleteMembersParamsInt{
				Method: 1,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type param ileagal1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type MembersParamsIDInt struct {
				ID         int    `json:"id" binding:"required"`
				MemberType string `json:"type" binding:"required"`
			}

			type AddOrDeleteMembersParams struct {
				Method  string               `json:"method" binding:"required"`
				Members []MembersParamsIDInt `json:"members" binding:"required"`
			}

			member := MembersParamsIDInt{
				ID:         1,
				MemberType: "xxxxx",
			}
			reqParam := AddOrDeleteMembersParams{
				Method: "POST",
				Members: []MembersParamsIDInt{
					0: member,
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type param ileagal2", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			type MembersParamsTypeInt struct {
				ID         string `json:"id" binding:"required"`
				MemberType int    `json:"type" binding:"required"`
			}

			type AddOrDeleteMembersParams struct {
				Method  string                 `json:"method" binding:"required"`
				Members []MembersParamsTypeInt `json:"members" binding:"required"`
			}

			member := MembersParamsTypeInt{
				ID:         "xxxx",
				MemberType: 1,
			}
			reqParam := AddOrDeleteMembersParams{
				Method: "POST",
				Members: []MembersParamsTypeInt{
					0: member,
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAddAndDeleteGroupMembersFail(t *testing.T) {
	Convey("AddAndDeleteGroupMembersFail", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/group-members/xxxx"

		Convey("token expired", func() {
			tempTarget := target + "xxxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("POST", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("method param ileagal", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			reqParam := AddOrDeleteMembersParams{
				Method: "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("type param ileagal", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			member := MembersParams{
				ID:         "xxxx",
				MemberType: "xxxxx",
			}
			reqParam := AddOrDeleteMembersParams{
				Method: "POST",
				Members: []MembersParams{
					0: member,
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestAddAndDeleteGroupMembersParamsSuccess(t *testing.T) {
	Convey("AddAndDeleteGroupMembersParamsSuccess", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/group-members/xxxx"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}

			tempTarget := target + "xxxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("POST", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("AddOrDeleteGroupMemebers err", func() {
			testErr := rest.NewHTTPError("errorxxx", 409019000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().AddOrDeleteGroupMemebers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testErr)

			member := MembersParams{
				ID:         "xxxx",
				MemberType: "user",
			}
			reqParam := AddOrDeleteMembersParams{
				Method: "DELETE",
				Members: []MembersParams{
					0: member,
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, 409019000)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("AddOrDeleteGroupMemebers Success", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().AddOrDeleteGroupMemebers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			member := MembersParams{
				ID:         "xxxx",
				MemberType: "department",
			}
			reqParam := AddOrDeleteMembersParams{
				Method: "DELETE",
				Members: []MembersParams{
					0: member,
				},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
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

func TestGetMembersIDFail(t *testing.T) {
	Convey("GetMembersID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/group-members"

		Convey("no GroupIDs param", func() {
			reqParam := GetMembersIDParam{
				Method: "GET",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("unsupported method", func() {
			reqParam := GetMembersIDParam{
				Method:   "xxxxxx",
				GroupIDs: []string{"xxxxxx"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("unsupported show_unable_user", func() {
			reqParam := map[string]interface{}{
				"method":           "GET",
				"group_ids":        []string{"xxxxxx"},
				"user_enabled":     "xxxx",
				"show_unable_user": "xxxx",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("unsupported show_unable_user false", func() {
			reqParam := map[string]interface{}{
				"method":       "GET",
				"group_ids":    []string{"xxxxxx"},
				"user_enabled": "false",
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GroupIDs param is empty", func() {
			reqParam := GetMembersIDParam{
				Method:   "GET",
				GroupIDs: []string{},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupMembersID error", func() {
			tempErr := rest.NewHTTPError("用户组不存在", 400019003, nil)
			g.EXPECT().GetGroupMembersID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, tempErr)

			reqParam := GetMembersIDParam{
				Method:   "GET",
				GroupIDs: []string{"xxxxxx"},
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, tempErr, respParam)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetMembersIDSuccess(t *testing.T) {
	Convey("GetMembersID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/group-members"

		Convey("success", func() {
			testUserIds := []string{"user_id"}
			testDepartmentIds := []string{"department_id"}
			g.EXPECT().GetGroupMembersID(gomock.Any(), gomock.Any(), gomock.Any()).Return(testUserIds, testDepartmentIds, nil)

			reqParam := map[string]interface{}{
				"method":       "GET",
				"group_ids":    []string{"xxxxxx"},
				"user_enabled": true,
			}
			reqParamByte, _ := jsoniter.Marshal(reqParam)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqParamByte))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var respParam interface{}
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSearchGroupByKeyParams(t *testing.T) {
	Convey("SearchGroupByKeyParams", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/search-in-group"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}

			tempTarget := target
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("param keyword err", func() {
			testErr := rest.NewHTTPError("invalid keyword type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyxxx=1&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param start err", func() {
			testErr := rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&offset=-1&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param start err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&offset=nil&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param offset err2", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&offset=nil&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&limit=nil&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err2", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&limit=nil&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err3", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&limit=1001&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err4", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&limit=-1&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param type err", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxxxx&type=zzz"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestSearchGroupByKey(t *testing.T) {
	Convey("SearchGroupByKeyword", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/search-in-group"

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("param keyword empty", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=&type=zzz"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param type err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?keyword=xxx&type=null"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SearchGroupByKeyword err", func() {
			var tmp interfaces.GMSearchOutInfo
			testErr := rest.NewHTTPError("xxxx", 503000000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			c.EXPECT().SearchGroupAndMemberInfoByKey(gomock.Any()).AnyTimes().Return(tmp, testErr)
			target += "?keyword=xxx&type=member"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SearchGroupByKeyword Success", func() {
			tmp := interfaces.NameInfo{
				ID:   "xxx",
				Name: "yyy",
			}
			tmpUser := interfaces.MemberInfo{
				ID:         "yyy",
				Name:       "kkkk",
				NType:      1,
				GroupNames: []string{"tttttt"},
			}
			tmpXX := interfaces.GMSearchOutInfo{
				GroupInfos:  []interfaces.NameInfo{tmp},
				GroupNum:    1,
				MemberInfos: []interfaces.MemberInfo{tmpUser},
				MemberNum:   1,
			}

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			c.EXPECT().SearchGroupAndMemberInfoByKey(gomock.Any()).AnyTimes().Return(tmpXX, nil)
			target += "?keyword=xxx&type=member&type=group"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			var respParam map[string]interface{}
			err := jsoniter.Unmarshal(respBody, &respParam)

			assert.Equal(t, nil, err)

			tmp1Interface := respParam["members"]
			tmp2Interface := respParam["groups"]
			testMembers := tmp1Interface.(map[string]interface{})
			testGroups := tmp2Interface.(map[string]interface{})

			t1 := MemberInfo{
				ID:         "yyy",
				Name:       "kkkk",
				MemberType: "user",
			}

			sliceMembers := (testMembers["entries"]).([]interface{})
			sliceGroups := (testGroups["entries"]).([]interface{})

			tMember := (sliceMembers[0]).(map[string]interface{})
			tGroup := (sliceGroups[0]).(map[string]interface{})

			assert.Equal(t, int(testMembers["total_count"].(float64)), 1)
			assert.Equal(t, tMember["id"].(string), t1.ID)
			assert.Equal(t, tMember["name"].(string), t1.Name)
			assert.Equal(t, tMember["type"].(string), t1.MemberType)
			assert.Equal(t, int(testGroups["total_count"].(float64)), 1)
			assert.Equal(t, tGroup["id"].(string), tmp.ID)
			assert.Equal(t, tGroup["name"].(string), tmp.Name)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetGroupOnClient(t *testing.T) {
	Convey("GetGroupOnClient", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/groups"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}

			tempTarget := target
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("param start err", func() {
			testErr := rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?offset=-1"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param start err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?offset=nil"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err", func() {
			testErr := rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?limit=-2"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?limit=nil"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetGroupOnClient err", func() {
			testErr := rest.NewHTTPError("xxxx", 503000000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroupOnClient(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, testErr)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getGroupOnClient Success", func() {
			tmp := interfaces.NameInfo{
				ID:   "xxx",
				Name: "zzz",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetGroupOnClient(gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.NameInfo{tmp}, 2, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			var respParam ListInfo
			err := jsoniter.Unmarshal(respBody, &respParam)

			assert.Equal(t, nil, err)
			assert.Equal(t, len(respParam.Entries), 1)
			assert.Equal(t, respParam.TotalCount, 2)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetMembersOnClient(t *testing.T) {
	Convey("GetMembersOnClient", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/group-members/xxxxxx"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}

			tempTarget := target
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "user_id",
		}

		Convey("param start err", func() {
			testErr := rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?offset=-1"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param start err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?offset=nil"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err", func() {
			testErr := rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?limit=-2"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param limit err1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			target += "?limit=nil"
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetMemberOnClient err", func() {
			testErr := rest.NewHTTPError("xxxx", 503000000, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetMemberOnClient(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, testErr)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			respBody, _ := io.ReadAll(result.Body)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getMembersOnClient Success", func() {
			tmp := interfaces.MemberSimpleInfo{
				ID:    "xxx",
				Name:  "zzz",
				NType: 1,
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().GetMemberOnClient(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.MemberSimpleInfo{tmp}, 2, nil)
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			var respParam ListInfo
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, len(respParam.Entries), 1)
			assert.Equal(t, respParam.TotalCount, 2)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestUserMatch(t *testing.T) {
	Convey("UserMatch", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("user-management")

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/console/group-members/xxxx/user-match"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorx", 409019000, nil)
		Convey("name is not exist", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid name")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		var uInfo interfaces.GroupMemberInfo
		var mInfos []interfaces.GroupMemberInfo
		Convey("UserMatch error", func() {
			tempTarget := target + "?name=1"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().UserMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, uInfo, mInfos, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		D1NameInfo := interfaces.NameInfo{ID: "d1d1", Name: "d1"}
		D2NameInfo := interfaces.NameInfo{ID: "d2d1", Name: "d2"}
		D3NameInfo := interfaces.NameInfo{ID: "d3d1", Name: "d3"}
		D4NameInfo := interfaces.NameInfo{ID: "d4d1", Name: "d4"}
		D5NameInfo := interfaces.NameInfo{ID: "d5d1", Name: "d5"}
		D6NameInfo := interfaces.NameInfo{ID: "d6d1", Name: "d6"}

		path1 := []interfaces.NameInfo{D1NameInfo, D2NameInfo}
		path2 := []interfaces.NameInfo{D2NameInfo, D3NameInfo}
		path3 := []interfaces.NameInfo{D4NameInfo, D5NameInfo}
		path4 := []interfaces.NameInfo{D1NameInfo, D6NameInfo}

		uInfo.ID = strAccount
		uInfo.ParentDeps = [][]interfaces.NameInfo{path1, path3}

		member1 := interfaces.GroupMemberInfo{
			ID:         strUserID,
			Name:       strAccount,
			MemberType: 1,
			ParentDeps: [][]interfaces.NameInfo{path2, path4},
		}
		member2 := interfaces.GroupMemberInfo{
			ID:         strUserID,
			Name:       strAccount,
			MemberType: 2,
			ParentDeps: [][]interfaces.NameInfo{path1, path4},
		}

		mInfos = []interfaces.GroupMemberInfo{member1, member2}
		Convey("success", func() {
			tempTarget := target + "?name=1"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().UserMatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, uInfo, mInfos, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam["result"].(bool), true)
			assert.Equal(t, respParam["id"].(string), uInfo.ID)
			testParentDeps := respParam["parent_deps"].([]interface{})
			assert.Equal(t, len(testParentDeps), 2)
			testParentDep1 := testParentDeps[0].([]interface{})
			assert.Equal(t, len(testParentDep1), 2)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["id"], D1NameInfo.ID)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["name"], D1NameInfo.Name)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["type"], "department")
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["id"], D2NameInfo.ID)
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["name"], D2NameInfo.Name)
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["type"], "department")
			testParentDep2 := testParentDeps[1].([]interface{})
			assert.Equal(t, len(testParentDep2), 2)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["id"], D4NameInfo.ID)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["name"], D4NameInfo.Name)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["type"], "department")
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["id"], D5NameInfo.ID)
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["name"], D5NameInfo.Name)
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["type"], "department")

			testGroupMembers := respParam["group_members"].([]interface{})
			assert.Equal(t, len(testGroupMembers), 2)

			testGM1 := testGroupMembers[0].(map[string]interface{})
			assert.Equal(t, testGM1["name"].(string), member1.Name)
			assert.Equal(t, testGM1["id"].(string), member1.ID)
			assert.Equal(t, testGM1["type"].(string), "user")

			testGM1ParentDeps := testGM1["parent_deps"].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps), 2)
			testGM1ParentDeps11 := testGM1ParentDeps[0].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps11), 2)
			testGM1P1 := testGM1ParentDeps11[0].(map[string]interface{})
			assert.Equal(t, testGM1P1["name"].(string), D2NameInfo.Name)
			assert.Equal(t, testGM1P1["id"].(string), D2NameInfo.ID)
			assert.Equal(t, testGM1P1["type"].(string), "department")
			testGM1P2 := testGM1ParentDeps11[1].(map[string]interface{})
			assert.Equal(t, testGM1P2["name"].(string), D3NameInfo.Name)
			assert.Equal(t, testGM1P2["id"].(string), D3NameInfo.ID)
			assert.Equal(t, testGM1P2["type"].(string), "department")
			testGM1ParentDeps12 := testGM1ParentDeps[1].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps12), 2)
			testGM1P111 := testGM1ParentDeps12[0].(map[string]interface{})
			assert.Equal(t, testGM1P111["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM1P111["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM1P111["type"].(string), "department")
			testGM1P112 := testGM1ParentDeps12[1].(map[string]interface{})
			assert.Equal(t, testGM1P112["name"].(string), D6NameInfo.Name)
			assert.Equal(t, testGM1P112["id"].(string), D6NameInfo.ID)
			assert.Equal(t, testGM1P112["type"].(string), "department")

			testGM2 := testGroupMembers[1].(map[string]interface{})
			assert.Equal(t, testGM2["name"].(string), member2.Name)
			assert.Equal(t, testGM2["id"].(string), member2.ID)
			assert.Equal(t, testGM2["type"].(string), "department")

			testGM2ParentDeps := testGM2["parent_deps"].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps), 2)
			testGM2ParentDeps21 := testGM2ParentDeps[0].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps21), 2)
			testGM2P1 := testGM2ParentDeps21[0].(map[string]interface{})
			assert.Equal(t, testGM2P1["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM2P1["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM2P1["type"].(string), "department")
			testGM2P2 := testGM2ParentDeps21[1].(map[string]interface{})
			assert.Equal(t, testGM2P2["name"].(string), D2NameInfo.Name)
			assert.Equal(t, testGM2P2["id"].(string), D2NameInfo.ID)
			assert.Equal(t, testGM2P2["type"].(string), "department")
			testGM2ParentDeps22 := testGM2ParentDeps[1].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps22), 2)
			testGM2P21 := testGM2ParentDeps22[0].(map[string]interface{})
			assert.Equal(t, testGM2P21["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM2P21["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM2P21["type"].(string), "department")
			testGM2P22 := testGM2ParentDeps22[1].(map[string]interface{})
			assert.Equal(t, testGM2P22["name"].(string), D6NameInfo.Name)
			assert.Equal(t, testGM2P22["id"].(string), D6NameInfo.ID)
			assert.Equal(t, testGM2P22["type"].(string), "department")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen,dupl
func TestSearchInAllGroupOrg(t *testing.T) {
	Convey("searchInAllGroupOrg", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("user-management")

		h := mock.NewMockHydra(ctrl)
		g := mock.NewMockLogicsGroup(ctrl)
		u := mock.NewMockLogicsUser(ctrl)
		c := mock.NewMockLogicsCombine(ctrl)
		testGRestHandler := newGroupRESTHandler(g, u, c, h)

		testGRestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/console/search-users-in-group"

		Convey("token expired", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("group_id is not exist", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid group_id")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("key is not exist", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid key")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset is not valid", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx&&key=xxx&&offset=xx"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid offset type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset i < 0", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx&&key=xxx&&offset=-1"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid offset type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit i < 0", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx&&key=xxx&&limit=-1"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit i > 1000", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx&&key=xxx&&limit=1001"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit is not valid", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			tempTarget := target + "?group_id=xxxx&&key=xxx&&limit=xxxx"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorx", 409019000, nil)
		uInfo := make(map[string]interfaces.GroupMemberInfo)
		mInfos := make(map[string][]interfaces.GroupMemberInfo)
		Convey("SearchInAllGroupOrg error", func() {
			tempTarget := target + "?key=xxx&&group_id=xxxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().SearchInAllGroupOrg(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(12, []string{}, uInfo, mInfos, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		D1NameInfo := interfaces.NameInfo{ID: "d1d1", Name: "d1"}
		D2NameInfo := interfaces.NameInfo{ID: "d2d1", Name: "d2"}
		D3NameInfo := interfaces.NameInfo{ID: "d3d1", Name: "d3"}
		D4NameInfo := interfaces.NameInfo{ID: "d4d1", Name: "d4"}
		D5NameInfo := interfaces.NameInfo{ID: "d5d1", Name: "d5"}
		D6NameInfo := interfaces.NameInfo{ID: "d6d1", Name: "d6"}

		path1 := []interfaces.NameInfo{D1NameInfo, D2NameInfo}
		path2 := []interfaces.NameInfo{D2NameInfo, D3NameInfo}
		path3 := []interfaces.NameInfo{D4NameInfo, D5NameInfo}
		path4 := []interfaces.NameInfo{D1NameInfo, D6NameInfo}

		tempUserInfo := interfaces.GroupMemberInfo{
			ID:         strAccount,
			Name:       strName111,
			ParentDeps: [][]interfaces.NameInfo{path1, path3},
		}
		uInfo[strAccount] = tempUserInfo

		member1 := interfaces.GroupMemberInfo{
			ID:         strUserID,
			Name:       strAccount,
			MemberType: 1,
			ParentDeps: [][]interfaces.NameInfo{path2, path4},
		}
		member2 := interfaces.GroupMemberInfo{
			ID:         strUserID,
			Name:       strAccount,
			MemberType: 2,
			ParentDeps: [][]interfaces.NameInfo{path1, path4},
		}

		mInfos[strAccount] = []interfaces.GroupMemberInfo{member1, member2}
		Convey("success", func() {
			tempTarget := target + "?key=a&&group_id=xxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    true,
				VisitorID: "user_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			g.EXPECT().SearchInAllGroupOrg(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(12, []string{strAccount}, uInfo, mInfos, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)

			assert.Equal(t, respParam["total_count"].(float64), float64(12))

			tempParas := respParam["entries"].([]interface{})
			assert.Equal(t, len(tempParas), 1)
			tempParas1 := tempParas[0].(map[string]interface{})
			assert.Equal(t, tempParas1["id"].(string), uInfo[strAccount].ID)
			assert.Equal(t, tempParas1["type"].(string), strUser)
			assert.Equal(t, tempParas1["name"].(string), strName111)
			testParentDeps := tempParas1["parent_deps"].([]interface{})
			assert.Equal(t, len(testParentDeps), 2)
			testParentDep1 := testParentDeps[0].([]interface{})
			assert.Equal(t, len(testParentDep1), 2)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["id"], D1NameInfo.ID)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["name"], D1NameInfo.Name)
			assert.Equal(t, testParentDep1[0].(map[string]interface{})["type"], "department")
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["id"], D2NameInfo.ID)
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["name"], D2NameInfo.Name)
			assert.Equal(t, testParentDep1[1].(map[string]interface{})["type"], "department")
			testParentDep2 := testParentDeps[1].([]interface{})
			assert.Equal(t, len(testParentDep2), 2)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["id"], D4NameInfo.ID)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["name"], D4NameInfo.Name)
			assert.Equal(t, testParentDep2[0].(map[string]interface{})["type"], "department")
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["id"], D5NameInfo.ID)
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["name"], D5NameInfo.Name)
			assert.Equal(t, testParentDep2[1].(map[string]interface{})["type"], "department")

			testGroupMembers := tempParas1["group_members"].([]interface{})
			assert.Equal(t, len(testGroupMembers), 2)

			testGM1 := testGroupMembers[0].(map[string]interface{})
			assert.Equal(t, testGM1["name"].(string), member1.Name)
			assert.Equal(t, testGM1["id"].(string), member1.ID)
			assert.Equal(t, testGM1["type"].(string), "user")

			testGM1ParentDeps := testGM1["parent_deps"].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps), 2)
			testGM1ParentDeps11 := testGM1ParentDeps[0].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps11), 2)
			testGM1P1 := testGM1ParentDeps11[0].(map[string]interface{})
			assert.Equal(t, testGM1P1["name"].(string), D2NameInfo.Name)
			assert.Equal(t, testGM1P1["id"].(string), D2NameInfo.ID)
			assert.Equal(t, testGM1P1["type"].(string), "department")
			testGM1P2 := testGM1ParentDeps11[1].(map[string]interface{})
			assert.Equal(t, testGM1P2["name"].(string), D3NameInfo.Name)
			assert.Equal(t, testGM1P2["id"].(string), D3NameInfo.ID)
			assert.Equal(t, testGM1P2["type"].(string), "department")
			testGM1ParentDeps12 := testGM1ParentDeps[1].([]interface{})
			assert.Equal(t, len(testGM1ParentDeps12), 2)
			testGM1P111 := testGM1ParentDeps12[0].(map[string]interface{})
			assert.Equal(t, testGM1P111["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM1P111["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM1P111["type"].(string), "department")
			testGM1P112 := testGM1ParentDeps12[1].(map[string]interface{})
			assert.Equal(t, testGM1P112["name"].(string), D6NameInfo.Name)
			assert.Equal(t, testGM1P112["id"].(string), D6NameInfo.ID)
			assert.Equal(t, testGM1P112["type"].(string), "department")

			testGM2 := testGroupMembers[1].(map[string]interface{})
			assert.Equal(t, testGM2["name"].(string), member2.Name)
			assert.Equal(t, testGM2["id"].(string), member2.ID)
			assert.Equal(t, testGM2["type"].(string), "department")

			testGM2ParentDeps := testGM2["parent_deps"].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps), 2)
			testGM2ParentDeps21 := testGM2ParentDeps[0].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps21), 2)
			testGM2P1 := testGM2ParentDeps21[0].(map[string]interface{})
			assert.Equal(t, testGM2P1["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM2P1["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM2P1["type"].(string), "department")
			testGM2P2 := testGM2ParentDeps21[1].(map[string]interface{})
			assert.Equal(t, testGM2P2["name"].(string), D2NameInfo.Name)
			assert.Equal(t, testGM2P2["id"].(string), D2NameInfo.ID)
			assert.Equal(t, testGM2P2["type"].(string), "department")
			testGM2ParentDeps22 := testGM2ParentDeps[1].([]interface{})
			assert.Equal(t, len(testGM2ParentDeps22), 2)
			testGM2P21 := testGM2ParentDeps22[0].(map[string]interface{})
			assert.Equal(t, testGM2P21["name"].(string), D1NameInfo.Name)
			assert.Equal(t, testGM2P21["id"].(string), D1NameInfo.ID)
			assert.Equal(t, testGM2P21["type"].(string), "department")
			testGM2P22 := testGM2ParentDeps22[1].(map[string]interface{})
			assert.Equal(t, testGM2P22["name"].(string), D6NameInfo.Name)
			assert.Equal(t, testGM2P22["id"].(string), D6NameInfo.ID)
			assert.Equal(t, testGM2P22["type"].(string), "department")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
