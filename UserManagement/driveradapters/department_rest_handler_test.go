package driveradapters

import (
	"errors"
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

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/interfaces/mock"
)

func newDepartHandler(dp interfaces.LogicsDepartment, h interfaces.Hydra) *departRestHandler {
	return &departRestHandler{
		depart: dp,
		hydra:  h,
	}
}
func TestGetAccessorIDsOfDepart(t *testing.T) {
	Convey("getAccessorIDsOfDepart", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id/accessor_ids"

		Convey("different", func() {
			assert.Equal(t, 1, 1)
		})

		Convey("getAccessorIDsOfUser error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			dp.EXPECT().GetAccessorIDsOfDepartment(gomock.Any()).Return(nil, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			groupIDs := []string{"dept_id", "xxxx"}
			dp.EXPECT().GetAccessorIDsOfDepartment(gomock.Any()).Return(groupIDs, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []string{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas, groupIDs)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartMemberIDs(t *testing.T) {
	Convey("getDepartMemberIDs", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id/member_ids"

		Convey("getDepartMemberIDs error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			var info interfaces.DepartMemberID
			dp.EXPECT().GetDepartMemberIDs(gomock.Any()).Return(info, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("different", func() {
			assert.Equal(t, 2, 2)
		})

		Convey("Success", func() {
			var info interfaces.DepartMemberID
			info.DepartIDs = append(info.DepartIDs, "xxxx")
			info.UserIDs = append(info.UserIDs, "yyyy")
			dp.EXPECT().GetDepartMemberIDs(gomock.Any()).Return(info, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			var strParmas interfaces.DepartMemberID
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas, info)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartManagers(t *testing.T) {
	Convey("getDepartManagers", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id/managers,manager,code,enabled"

		Convey("Success", func() {
			tempData := []interfaces.DepartInfo{
				{
					ID: "aa",
					Managers: []interfaces.NameInfo{
						{
							ID:   "xx1",
							Name: "name1",
						},
					},
					Manager: interfaces.NameInfo{
						ID:   "xx1",
						Name: "name1",
					},
					Enabled: true,
					Code:    "code1",
				},
				{
					ID: "bb",
					Managers: []interfaces.NameInfo{
						{
							ID:   "xx2",
							Name: "name2",
						},
					},
					Manager: interfaces.NameInfo{
						ID:   "aa",
						Name: "",
					},
					Enabled: true,
					Code:    "code2",
				},
			}
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(tempData, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, len(strParmas), 2)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success2", func() {
			tempData := []interfaces.DepartInfo{
				{
					ID: "bb",
					Managers: []interfaces.NameInfo{
						{
							ID:   "xx2",
							Name: "name2",
						},
					},
					Manager: interfaces.NameInfo{
						ID:   "",
						Name: "",
					},
				},
			}
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(tempData, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, len(strParmas), 1)

			temp := strParmas[0].(map[string]interface{})
			_, ok := temp[strManager]
			assert.Equal(t, ok, false)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success1", func() {
			tempData := []interfaces.DepartInfo{
				{
					ID: "aa",
					Managers: []interfaces.NameInfo{
						{
							ID:   "xx1",
							Name: "name1",
						},
					},
					Manager: interfaces.NameInfo{
						ID:   "xx1",
						Name: "name1",
					},
				},
			}
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(tempData, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, len(strParmas), 1)

			temp := strParmas[0].(map[string]interface{})
			assert.Equal(t, temp[strManager].(map[string]interface{})["id"], "xx1")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getDepartManagers not found", func() {
			tempData := make([]interfaces.DepartInfo, 0)
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(tempData, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make([]interface{}, 0)
			_ = jsoniter.Unmarshal(respBody, &strParmas)

			assert.Equal(t, len(strParmas), 0)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getDepartManagers error", func() {
			tmpErr := errors.New("test")
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, tmpErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartAllUserInfo(t *testing.T) {
	Convey("getDepartInfo", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id/all_user_ids"

		Convey("fields error", func() {
			target += ",xxxx"

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields user_enabled err error", func() {
			target += "?user_enabled=xxx"

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetAllDepartUserIDs error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			info := make([]string, 0)
			dp.EXPECT().GetAllDepartUserIDs(gomock.Any(), gomock.Any()).Return(info, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			info := make([]string, 0)
			info = append(info, "xxxx")
			dp.EXPECT().GetAllDepartUserIDs(gomock.Any(), gomock.Any()).Return(info, nil)

			req := httptest.NewRequest("GET", target+"?user_enabled=true", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)

			tmp := strParmas["all_user_ids"].([]interface{})
			assert.Equal(t, err, nil)
			assert.Equal(t, tmp[0].(string), info[0])

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartMemberInfo(t *testing.T) {
	Convey("获取部门子部门和子用户信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/department-members/:department_id/"

		Convey("token失效-报错", func() {
			tempTarget := target + "users"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
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
			VisitorID: "department_id",
		}
		Convey("fields参数错误-报错", func() {
			tempTarget := target + "users,xxxx?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset为-1-报错 ", func() {
			tempTarget := target + "users?role=sys_admin&offset=-1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit为0-报错 ", func() {
			tempTarget := target + "users?role=sys_admin&limit=0"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit为1001-报错 ", func() {
			tempTarget := target + "users?role=sys_admin&limit=1001"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role参数错误-报错 ", func() {
			tempTarget := target + "users?role=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role参数未传-报错", func() {
			tempTarget := target + "users"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetDepartMemberInfo error", func() {
			tempTarget := target + "users?role=sys_admin"
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().GetDepartMemberInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, nil, 0, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		temp1 := interfaces.DepartInfo{
			ID:            "ID1",
			Name:          "Name1",
			IsRoot:        true,
			BDepartExistd: true,
			BUserExistd:   true,
		}

		temp2 := interfaces.ObjectBaseInfo{
			ID:   "ID2",
			Name: "Name2",
			Type: "user",
		}

		depInfo := []interfaces.DepartInfo{temp1}
		userInfo := []interfaces.ObjectBaseInfo{temp2}

		Convey("success", func() {
			tempTarget := target + "users,departments?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().GetDepartMemberInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(depInfo, 2, userInfo, 4, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			respParam := make(map[string]ListInfo)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["users"].TotalCount, 4)
			outTemp := respParam["users"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp["id"], temp2.ID)
			assert.Equal(t, respParam["departments"].TotalCount, 2)
			outTemp1 := respParam["departments"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp1["id"], temp1.ID)
			assert.Equal(t, outTemp1["name"], temp1.Name)
			assert.Equal(t, outTemp1["type"], "department")
			assert.Equal(t, outTemp1["is_root"].(bool), true)
			assert.Equal(t, outTemp1["user_existed"].(bool), true)
			assert.Equal(t, outTemp1["depart_existed"].(bool), true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestGetManagementDepartMemberInfo(t *testing.T) {
	Convey("管理控制台获取部门子部门和子用户信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/department-members/:department_id/"

		Convey("token失效-报错", func() {
			tempTarget := target + "users"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
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
			VisitorID: "department_id",
		}
		Convey("fields参数错误-报错", func() {
			tempTarget := target + "xxxx?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset为-1-报错 ", func() {
			tempTarget := target + "departments?role=sys_admin&offset=-1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit为0-报错 ", func() {
			tempTarget := target + "departments?role=sys_admin&limit=0"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit为1001-报错 ", func() {
			tempTarget := target + "departments?role=sys_admin&limit=1001"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role参数错误-报错 ", func() {
			tempTarget := target + "departments?role=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role参数未传-报错", func() {
			tempTarget := target + "departments"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("GetDepartMemberInfo error", func() {
			tempTarget := target + "departments?role=sys_admin"
			testErr := rest.NewHTTPError("errorx", 403019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().GetDepartMemberInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, nil, 0, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		dep1 := interfaces.ObjectBaseInfo{
			ID:   "parent1",
			Name: "parent1",
			Type: "department",
			Code: "code1",
		}

		dep2 := interfaces.ObjectBaseInfo{
			ID:   "parent2",
			Name: "parent2",
			Type: "department",
			Code: "code2",
		}

		temp1 := interfaces.DepartInfo{
			ID:            "ID1",
			Name:          "Name1",
			IsRoot:        true,
			BDepartExistd: true,
			BUserExistd:   true,
			Manager: interfaces.NameInfo{
				ID:   "user1",
				Name: "name1",
			},
			Code:    "code1",
			Enabled: true,
			Remark:  "remark1",
			Email:   "email1",
			ParentDeps: []interfaces.ObjectBaseInfo{
				dep1,
				dep2,
			},
		}

		depInfo := []interfaces.DepartInfo{temp1}

		Convey("success", func() {
			tempTarget := target + "departments?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().GetDepartMemberInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(depInfo, 2, nil, 0, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)

			respParam := make(map[string]ListInfo)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam["departments"].TotalCount, 2)
			outTemp1 := respParam["departments"].Entries[0].(map[string]interface{})
			assert.Equal(t, outTemp1["id"], temp1.ID)
			assert.Equal(t, outTemp1["name"], temp1.Name)
			assert.Equal(t, outTemp1["type"], "department")
			assert.Equal(t, outTemp1["is_root"].(bool), true)
			assert.Equal(t, outTemp1["depart_existed"].(bool), true)
			assert.Equal(t, outTemp1["manager"].(map[string]interface{})["id"], temp1.Manager.ID)
			assert.Equal(t, outTemp1["manager"].(map[string]interface{})["name"], temp1.Manager.Name)
			assert.Equal(t, outTemp1["manager"].(map[string]interface{})["type"], "user")
			assert.Equal(t, outTemp1["code"], temp1.Code)
			assert.Equal(t, outTemp1["enabled"].(bool), temp1.Enabled)
			assert.Equal(t, outTemp1[strRemark], temp1.Remark)
			assert.Equal(t, outTemp1["email"], temp1.Email)
			assert.Equal(t, len(outTemp1["parent_deps"].([]interface{})), 2)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[0].(map[string]interface{})["id"], dep1.ID)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[0].(map[string]interface{})["name"], dep1.Name)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[0].(map[string]interface{})["type"], dep1.Type)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[0].(map[string]interface{})["code"], dep1.Code)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[1].(map[string]interface{})["id"], dep2.ID)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[1].(map[string]interface{})["name"], dep2.Name)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[1].(map[string]interface{})["type"], dep2.Type)
			assert.Equal(t, outTemp1["parent_deps"].([]interface{})[1].(map[string]interface{})["code"], dep2.Code)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartInfo(t *testing.T) {
	Convey("获取部门信息-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/ID1/"

		Convey("url参数类型错误", func() {
			tempTarget := target + "xxxxx"

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, respParam.Cause, "invalid params")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorx", 403019001, nil)
		Convey("GetDepartsInfo报错", func() {
			tempTarget := target + "managers"
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		departInfo := interfaces.DepartInfo{
			ID:   "ID1",
			Name: "Name1",
			ParentDeps: []interfaces.ObjectBaseInfo{
				{
					ID:   "sub1",
					Name: "name1",
					Type: "department",
				}, {
					ID:   "sub2",
					Name: "name2",
					Type: "department",
				},
			},
			Managers: []interfaces.NameInfo{
				{
					ID:   "user1",
					Name: "name1",
				}, {
					ID:   "user2",
					Name: "name2",
				},
			},
		}
		Convey("成功", func() {
			tempTarget := target + "name,parent_deps,managers"
			dp.EXPECT().GetDepartsInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.DepartInfo{departInfo}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(respParam), 1)

			temp1 := respParam[0].(map[string]interface{})
			assert.Equal(t, temp1["department_id"].(string), departInfo.ID)
			assert.Equal(t, temp1["name"].(string), departInfo.Name)

			tempParentDeps := temp1["parent_deps"].([]interface{})
			assert.Equal(t, len(tempParentDeps), 2)
			parentDep1 := tempParentDeps[0].(map[string]interface{})
			assert.Equal(t, parentDep1["id"].(string), "sub1")
			assert.Equal(t, parentDep1["name"].(string), "name1")
			assert.Equal(t, parentDep1["type"].(string), "department")
			parentDep2 := tempParentDeps[1].(map[string]interface{})
			assert.Equal(t, parentDep2["id"].(string), "sub2")
			assert.Equal(t, parentDep2["name"].(string), "name2")
			assert.Equal(t, parentDep2["type"].(string), "department")

			tempManagers := temp1["managers"].([]interface{})
			assert.Equal(t, len(tempManagers), 2)
			manager1 := tempManagers[0].(map[string]interface{})
			assert.Equal(t, manager1["id"].(string), "user1")
			assert.Equal(t, manager1["name"].(string), "name1")
			manager2 := tempManagers[1].(map[string]interface{})
			assert.Equal(t, manager2["id"].(string), "user2")
			assert.Equal(t, manager2["name"].(string), "name2")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartInfoByLevel(t *testing.T) {
	Convey("部门列举-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments"

		Convey("level参数不传，报错", func() {
			tempTarget := target

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, respParam.Cause, "invalid level")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("level为空字符串，获取部门列举，抛错400", func() {
			tempTarget := target + "?level="

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, rest.BadRequest)
			assert.Equal(t, respParam.Cause, "invalid level")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("errorx", 403019001, nil)
		infos := make([]interfaces.ObjectBaseInfo, 0)
		Convey("获取部门层级信息错误，报错", func() {
			tempTarget := target + "?level=0"
			dp.EXPECT().GetDepartsInfoByLevel(gomock.Any()).AnyTimes().Return(infos, testErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam, testErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		depInfo1 := interfaces.ObjectBaseInfo{ID: "ID1", Name: "Name1", Type: "Typ1", ThirdID: "third_id1"}
		depInfo2 := interfaces.ObjectBaseInfo{ID: "ID2", Name: "Name2", Type: "Typ2", ThirdID: "third_id2"}
		infos = append(infos, depInfo1, depInfo2)
		Convey("获取部门层级信息成功", func() {
			tempTarget := target + "?level=0"
			dp.EXPECT().GetDepartsInfoByLevel(gomock.Any()).AnyTimes().Return(infos, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(respParam), 2)

			temp1 := respParam[0].(map[string]interface{})
			assert.Equal(t, temp1["id"], depInfo1.ID)
			assert.Equal(t, temp1["name"], depInfo1.Name)
			assert.Equal(t, temp1["type"], depInfo1.Type)
			assert.Equal(t, temp1["third_id"], depInfo1.ThirdID)
			temp2 := respParam[1].(map[string]interface{})
			assert.Equal(t, temp2["id"], depInfo2.ID)
			assert.Equal(t, temp2["name"], depInfo2.Name)
			assert.Equal(t, temp2["type"], depInfo2.Type)
			assert.Equal(t, temp2["third_id"], depInfo2.ThirdID)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetDepartAllUserInfos(t *testing.T) {
	Convey("getDepartAllUserInfos", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id/all_users"

		user1 := interfaces.UserBaseInfo{
			ID:        "ID1",
			Name:      "name1",
			Account:   "account1",
			Email:     "email1",
			TelNumber: "telephone1",
			ThirdAttr: "thridattr1",
			ThirdID:   "thirdid1",
		}
		user3 := interfaces.UserBaseInfo{
			ID:        "ID3",
			Name:      "name3",
			Account:   "account3",
			Email:     "email3",
			TelNumber: "telephone3",
			ThirdAttr: "thridattr3",
			ThirdID:   "thirdid3",
		}
		info := []interfaces.UserBaseInfo{user1, user3}
		Convey("GetAllDepartUserInfos error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			dp.EXPECT().GetAllDepartUserInfos(gomock.Any()).Return(info, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("Success", func() {
			dp.EXPECT().GetAllDepartUserInfos(gomock.Any()).Return(info, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(strParmas), 2)

			temp1 := strParmas[0].(map[string]interface{})
			assert.Equal(t, temp1["id"].(string), user1.ID)
			assert.Equal(t, temp1["name"].(string), user1.Name)
			assert.Equal(t, temp1["account"].(string), user1.Account)
			assert.Equal(t, temp1["email"].(string), user1.Email)
			assert.Equal(t, temp1["telephone"].(string), user1.TelNumber)
			assert.Equal(t, temp1[strThirdAttr].(string), user1.ThirdAttr)
			assert.Equal(t, temp1["third_id"].(string), user1.ThirdID)

			temp2 := strParmas[1].(map[string]interface{})
			assert.Equal(t, temp2["id"].(string), user3.ID)
			assert.Equal(t, temp2["name"].(string), user3.Name)
			assert.Equal(t, temp2["account"].(string), user3.Account)
			assert.Equal(t, temp2["email"].(string), user3.Email)
			assert.Equal(t, temp2["telephone"].(string), user3.TelNumber)
			assert.Equal(t, temp2[strThirdAttr].(string), user3.ThirdAttr)
			assert.Equal(t, temp2["third_id"].(string), user3.ThirdID)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestDeleteManageDepart(t *testing.T) {
	Convey("deleteManageDepart", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPublic(r)
		target := "/api/user-management/v1/management/departments/:department_id"

		Convey("token失效-报错", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "department_id",
		}
		tmpErr := errors.New("test")
		Convey("DeleteDepart-报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().DeleteDepart(gomock.Any(), gomock.Any()).AnyTimes().Return(tmpErr)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			dp.EXPECT().DeleteDepart(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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

func TestDeleteDepart(t *testing.T) {
	Convey("deleteDepart", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/departments/:department_id"

		tmpErr := errors.New("test")
		Convey("DeleteDepart-报错", func() {
			dp.EXPECT().DeleteDepart(gomock.Any(), gomock.Any()).AnyTimes().Return(tmpErr)

			req := httptest.NewRequest("DELETE", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			dp.EXPECT().DeleteDepart(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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

//nolint:funlen
func TestSearchManageDepart(t *testing.T) {
	Convey("searchManageDepart", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		dp := mock.NewMockLogicsDepartment(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newDepartHandler(dp, h)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/console/search-departments/"
		Convey("token 校验失败", func() {
			tempTarget := target + "name?role=1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: false}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("包含不存在的field", func() {
			tempTarget := target + "xxxx?role=1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid params fields")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role参数错误-报错", func() {
			tempTarget := target + "name?role=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid params role")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("param里面enabled不是true或者false之类", func() {
			tempTarget := target + "name?role=super_admin&enabled=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid enabled parameter")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset < 0", func() {
			tempTarget := target + "name?role=super_admin&offset=-1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid start type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit > 1000", func() {
			tempTarget := target + "name?role=super_admin&limit=1001"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit < 1", func() {
			tempTarget := target + "name?role=super_admin&limit=0"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("SearchDeparts error", func() {
			tempTarget := target + "name?role=super_admin&enabled=true"
			tempErr := rest.NewHTTPErrorV2(rest.Forbidden, "invalid params")
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			dp.EXPECT().SearchDeparts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, tempErr)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam, tempErr)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		manager1 := interfaces.NameInfo{
			ID: "1", Name: "1",
		}
		manager2 := interfaces.NameInfo{
			ID: "2", Name: "2",
		}

		parentDep1 := interfaces.ObjectBaseInfo{
			ID: "1xx", Name: "1xx", Type: strDepartment, Code: "code1",
		}
		parentDep2 := interfaces.ObjectBaseInfo{
			ID: "2xx", Name: "2xxx", Type: strDepartment, Code: "code2",
		}

		depart1 := interfaces.DepartInfo{
			ID: "1", Name: "test1", Code: "test1", Manager: manager1, Remark: "test1", Email: "test1", ParentDeps: []interfaces.ObjectBaseInfo{parentDep1, parentDep2}, Enabled: true,
		}
		depart2 := interfaces.DepartInfo{
			ID: "2", Name: "test2", Code: "test2", Manager: manager2, Remark: "test2", Email: "test2", ParentDeps: []interfaces.ObjectBaseInfo{}, Enabled: true,
		}
		testData1 := []interfaces.DepartInfo{
			depart1,
			depart2,
		}
		Convey("success", func() {
			tempTarget := target + "name,code,manager,remark,email,parent_deps,enabled?role=super_admin&enabled=true"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(interfaces.TokenIntrospectInfo{Active: true}, nil)
			dp.EXPECT().SearchDeparts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testData1, 24, nil)

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			resParam := map[string]interface{}{}
			err := jsoniter.Unmarshal(respBody, &resParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, int(resParam["total_count"].(float64)), 24)

			datas := resParam["entries"].([]interface{})
			assert.Equal(t, len(datas), 2)
			assert.Equal(t, datas[0].(map[string]interface{})["id"], "1")
			assert.Equal(t, datas[0].(map[string]interface{})["name"], "test1")
			assert.Equal(t, datas[0].(map[string]interface{})["code"], "test1")
			assert.Equal(t, datas[0].(map[string]interface{})["manager"].(map[string]interface{})["id"], manager1.ID)
			assert.Equal(t, datas[0].(map[string]interface{})["manager"].(map[string]interface{})["name"], manager1.Name)
			assert.Equal(t, datas[0].(map[string]interface{})["manager"].(map[string]interface{})["type"], strUser)
			assert.Equal(t, datas[0].(map[string]interface{})[strRemark], "test1")
			assert.Equal(t, datas[0].(map[string]interface{})["email"], "test1")
			assert.Equal(t, len(datas[0].(map[string]interface{})["parent_deps"].([]interface{})), 2)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[0].(map[string]interface{})["id"], parentDep1.ID)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[0].(map[string]interface{})["name"], parentDep1.Name)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[0].(map[string]interface{})["type"], strDepartment)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[0].(map[string]interface{})["code"], parentDep1.Code)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[1].(map[string]interface{})["id"], parentDep2.ID)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[1].(map[string]interface{})["name"], parentDep2.Name)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[1].(map[string]interface{})["type"], parentDep2.Type)
			assert.Equal(t, datas[0].(map[string]interface{})["parent_deps"].([]interface{})[1].(map[string]interface{})["code"], parentDep2.Code)
			assert.Equal(t, datas[0].(map[string]interface{})["enabled"], true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
