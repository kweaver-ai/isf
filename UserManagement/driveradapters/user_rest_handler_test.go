// Package driveradapters user AnyShare  用户逻辑接口处理层
package driveradapters

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	gerrors "github.com/kweaver-ai/go-lib/error"
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

const (
	strName111  = "Name111"
	strXXX      = "xxx"
	strEmail1   = "xxx@qq.com"
	strThirdID1 = "third_idxxx"
	strRemark1  = "remark1"
)

type UserBaseDataMock struct {
	Roles          []string                      `json:"roles" binding:"required"`
	Priority       int                           `json:"priority" binding:"required"`
	Enabled        bool                          `json:"enabled" binding:"required"`
	CsfLevel       int                           `json:"csf_level" binding:"required"`
	Name           string                        `json:"name" binding:"required"`
	ParentDeps     [][]interfaces.ObjectBaseInfo `json:"parent_deps" binding:"required"`
	ParentDepPaths []string                      `json:"parent_dep_paths" binding:"required"`
	Account        string                        `json:"account" binding:"required"`
	Frozen         bool                          `json:"frozen" binding:"required"`
	Authenticated  bool                          `json:"authenticated" binding:"required"`
	ID             string                        `json:"id" binding:"required"`
	Email          string                        `json:"email" binding:"required"`
	Telephone      string                        `json:"telephone" binding:"required"`
	ThirdAttr      string                        `json:"third_attr" binding:"required"`
	ThirdID        string                        `json:"third_id" binding:"required"`
	AuthType       string                        `json:"auth_type" binding:"required"`
	Manager        map[string]interface{}        `json:"manager"`
	CreatedAt      string                        `json:"created_at" binding:"required"`
}

type UserOwnDataMock struct {
	Avatar string `json:"avatar_url" binding:"required"`
}

func newUserRESTHandler(user interfaces.LogicsUser, h interfaces.Hydra, av interfaces.LogicsAvatar, role interfaces.LogicsRole) UserRestHandler {
	roleNameMap := make(map[interfaces.Role]string)
	roleNameMap[interfaces.SystemRoleSuperAdmin] = EnumSuperAdmin
	roleNameMap[interfaces.SystemRoleSysAdmin] = EnumSysAdmin
	roleNameMap[interfaces.SystemRoleAuditAdmin] = EnumAuditAdmin
	roleNameMap[interfaces.SystemRoleSecAdmin] = EnumSecAdmin
	roleNameMap[interfaces.SystemRoleOrgManager] = EnumOrgManager
	roleNameMap[interfaces.SystemRoleOrgAudit] = EnumOrgAudit
	roleNameMap[interfaces.SystemRoleNormalUser] = EnumNormaluser

	authTypeMap := make(map[interfaces.AuthType]string)
	authTypeMap[interfaces.Local] = EnumLocal
	authTypeMap[interfaces.Domain] = EnumDomain
	authTypeMap[interfaces.Third] = EnumThird

	PwdRetrievalMap := make(map[interfaces.PwdRetrievalStatus]string)
	PwdRetrievalMap[interfaces.PRSAvaliable] = EnumAvaliable
	PwdRetrievalMap[interfaces.PRSInvalidAccount] = EnumInvalidUser
	PwdRetrievalMap[interfaces.PRSDisableUser] = EnumDisableUser
	PwdRetrievalMap[interfaces.PRSUnablePWDRetrieval] = EnumUnablePWDRetrieval
	PwdRetrievalMap[interfaces.PRSNonLocalUser] = EnumNonLocalUser
	PwdRetrievalMap[interfaces.PRSEnablePwdControl] = EnumEnablePWDControl

	schema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(pwdErrInfoSchemaStr))
	incrementUserInfoSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(incrementUserInfoSchemaStr))
	getUserInfoSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getUserInfoSchemaStr))
	markerSchema, _ := gojsonschema.NewSchema(gojsonschema.NewStringLoader(markerSchemaStr))
	return &userRestHandler{
		user:                    user,
		roleNameMap:             roleNameMap,
		role:                    role,
		hydra:                   h,
		avatar:                  av,
		authTypeMap:             authTypeMap,
		pwdRetrievalMap:         PwdRetrievalMap,
		pwdErrInfoSchema:        schema,
		incrementUserInfoSchema: incrementUserInfoSchema,
		getUserInfoSchema:       getUserInfoSchema,
		markerSchema:            markerSchema,
	}
}

func TestGetAllBelongDepartmentID(t *testing.T) {
	Convey("getAllBelongDepartmentID", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/users/:user_id/department_ids"

		Convey("getAllBelongDepartmentID error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 400019001, nil)
			userLogics.EXPECT().GetAllBelongDepartmentIDs(gomock.Any()).Return(nil, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
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

		Convey("Success", func() {
			deptIDs := []string{"dept_id"}
			userLogics.EXPECT().GetAllBelongDepartmentIDs(gomock.Any()).Return(deptIDs, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []string{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas, deptIDs)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetAccessorIDsOfUser(t *testing.T) {
	Convey("getAccessorIDsOfUser", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/users/:user_id/accessor_ids"

		Convey("getAccessorIDsOfUser error", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			userLogics.EXPECT().GetAccessorIDsOfUser(gomock.Any()).Return(nil, tempErr)

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
			userLogics.EXPECT().GetAccessorIDsOfUser(gomock.Any()).Return(groupIDs, nil)

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

func TestGetUserBaseInfo(t *testing.T) {
	Convey("获取用户信息-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/users/xasdasd/xxxx"
		Convey("fileds包含规定外的字段-报错", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		target = "/api/user-management/v1/users/:user_id/roles"
		Convey("用户不存在-报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			var basInfo []interfaces.UserBaseInfo
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(basInfo, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 503)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		target = "/api/user-management/v1/users/:user_id/roles,priority,csf_level,enabled,parent_deps,name,email,telephone,third_attr,third_id,auth_type"
		Convey("获取所有信息-成功", func() {
			basInfo := interfaces.UserBaseInfo{}
			basInfo.VecRoles = make([]interfaces.Role, 0)
			basInfo.VecRoles = append(basInfo.VecRoles, interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin, interfaces.SystemRoleAuditAdmin,
				interfaces.SystemRoleSecAdmin, interfaces.SystemRoleOrgManager, interfaces.SystemRoleOrgAudit, interfaces.SystemRoleNormalUser)
			basInfo.CSFLevel = 22
			basInfo.Enabled = true
			basInfo.Priority = 324
			basInfo.Name = strName111
			basInfo.Email = strEmail1
			basInfo.TelNumber = "123456"
			basInfo.ThirdAttr = strThirdAttr
			basInfo.ThirdID = strThirdID1
			basInfo.AuthType = interfaces.Local

			temp := interfaces.ObjectBaseInfo{
				ID:   "ID1",
				Name: "Name1",
				Type: "department",
			}

			temp1 := []interfaces.ObjectBaseInfo{temp}
			basInfo.ParentDeps = [][]interfaces.ObjectBaseInfo{temp1}

			usersInfo := []interfaces.UserBaseInfo{basInfo}
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(usersInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []UserBaseDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas[0].Roles, []string{EnumSuperAdmin, EnumSysAdmin, EnumAuditAdmin, EnumSecAdmin, EnumOrgManager, EnumOrgAudit, EnumNormaluser})
			assert.Equal(t, strParmas[0].Enabled, basInfo.Enabled)
			assert.Equal(t, strParmas[0].Priority, basInfo.Priority)
			assert.Equal(t, strParmas[0].CsfLevel, basInfo.CSFLevel)
			assert.Equal(t, strParmas[0].Name, basInfo.Name)
			assert.Equal(t, strParmas[0].ParentDeps, basInfo.ParentDeps)
			assert.Equal(t, strParmas[0].Email, basInfo.Email)
			assert.Equal(t, strParmas[0].Telephone, basInfo.TelNumber)
			assert.Equal(t, strParmas[0].ThirdAttr, basInfo.ThirdAttr)
			assert.Equal(t, strParmas[0].ThirdID, basInfo.ThirdID)
			assert.Equal(t, strParmas[0].AuthType, EnumLocal)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestGetUserBaseInfoByPost(t *testing.T) {
	Convey("获取用户信息POST-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/batch-get-user-info"
		Convey("fileds包含规定外的字段-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["user_ids"] = []string{strXXX}
			tmpBody["fields"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("body不包含fields-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["user_ids"] = []string{strXXX}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "(root): fields is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("ids不存在-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = strXXX
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "(root): user_ids is required")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("ids不为数组-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = strXXX
			tmpBody["user_ids"] = strXXX
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "user_ids: Invalid type. Expected: array, given: string")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("ids不为字符串数组-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = strXXX
			tmpBody["user_ids"] = []int{1, 2}
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "user_ids.0: Invalid type. Expected: string, given: integer; user_ids.1: Invalid type. Expected: string, given: integer")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("methods不存在-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["user_ids"] = []string{"zzzz"}
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "(root): method is required")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("method不为“GET”-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = strXXX
			tmpBody["user_ids"] = []string{"xxxx"}
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "invalid method")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
		Convey("用户不存在-报错", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["user_ids"] = []string{"xxxx"}
			tmpBody["fields"] = []string{strThirdAttr}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			usersInfo := []interfaces.UserBaseInfo{}
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(usersInfo, tempErr)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 503)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取所有信息-成功", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["user_ids"] = []string{"xxxx"}
			tmpBody["fields"] = []string{strThirdAttr, "roles", "priority", "enabled", "csf_level", "name", "parent_deps", strAccount, "frozen",
				"authenticated", "email", "telephone", "third_id", "auth_type", "groups", "oss_id", "custom_attr", strManager, "created_at"}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			basInfo := interfaces.UserBaseInfo{}
			basInfo.VecRoles = make([]interfaces.Role, 0)
			basInfo.VecRoles = append(basInfo.VecRoles, interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin, interfaces.SystemRoleAuditAdmin,
				interfaces.SystemRoleSecAdmin, interfaces.SystemRoleOrgManager, interfaces.SystemRoleOrgAudit, interfaces.SystemRoleNormalUser)
			basInfo.CSFLevel = 22
			basInfo.Enabled = true
			basInfo.Priority = 324
			basInfo.Name = strName111
			basInfo.Email = strEmail1
			basInfo.TelNumber = "123456"
			basInfo.ThirdAttr = strThirdAttr
			basInfo.ThirdID = strThirdID1
			basInfo.AuthType = interfaces.Local
			basInfo.Manager.ID = "manager_id"
			time1 := time.Now()
			basInfo.CreatedAt = time1.Unix()
			temp := interfaces.ObjectBaseInfo{
				ID:   "ID1",
				Name: "Name1",
				Type: "department",
			}

			temp1 := []interfaces.ObjectBaseInfo{temp}
			basInfo.ParentDeps = [][]interfaces.ObjectBaseInfo{temp1}

			usersInfo := []interfaces.UserBaseInfo{basInfo}
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(usersInfo, nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []UserBaseDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas[0].Roles, []string{EnumSuperAdmin, EnumSysAdmin, EnumAuditAdmin, EnumSecAdmin, EnumOrgManager, EnumOrgAudit, EnumNormaluser})
			assert.Equal(t, strParmas[0].Enabled, basInfo.Enabled)
			assert.Equal(t, strParmas[0].Priority, basInfo.Priority)
			assert.Equal(t, strParmas[0].CsfLevel, basInfo.CSFLevel)
			assert.Equal(t, strParmas[0].Name, basInfo.Name)
			assert.Equal(t, strParmas[0].ParentDeps, basInfo.ParentDeps)
			assert.Equal(t, strParmas[0].Email, basInfo.Email)
			assert.Equal(t, strParmas[0].Telephone, basInfo.TelNumber)
			assert.Equal(t, strParmas[0].ThirdAttr, basInfo.ThirdAttr)
			assert.Equal(t, strParmas[0].ThirdID, basInfo.ThirdID)
			assert.Equal(t, strParmas[0].AuthType, EnumLocal)
			assert.Equal(t, strParmas[0].Manager["id"], basInfo.Manager.ID)
			assert.Equal(t, strParmas[0].Manager["type"], strUser)
			assert.Equal(t, strParmas[0].CreatedAt, time1.Format(time.RFC3339))

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取所有信息11-成功,但是manager为null", func() {
			tmpBody := make(map[string]interface{})
			tmpBody["method"] = mstrGET
			tmpBody["user_ids"] = []string{"xxxx"}
			tmpBody["fields"] = []string{strThirdAttr, "roles", "priority", "enabled", "csf_level", "name", "parent_deps", strAccount, "frozen",
				"authenticated", "email", "telephone", "third_id", "auth_type", "groups", "oss_id", "custom_attr", strManager}
			reqBody, _ := jsoniter.Marshal(tmpBody)

			basInfo := interfaces.UserBaseInfo{}
			basInfo.VecRoles = make([]interfaces.Role, 0)
			basInfo.VecRoles = append(basInfo.VecRoles, interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin, interfaces.SystemRoleAuditAdmin,
				interfaces.SystemRoleSecAdmin, interfaces.SystemRoleOrgManager, interfaces.SystemRoleOrgAudit, interfaces.SystemRoleNormalUser)
			basInfo.CSFLevel = 22
			basInfo.Enabled = true
			basInfo.Priority = 324
			basInfo.Name = strName111
			basInfo.Email = strEmail1
			basInfo.TelNumber = "123456"
			basInfo.ThirdAttr = strThirdAttr
			basInfo.ThirdID = strThirdID1
			basInfo.AuthType = interfaces.Local
			basInfo.Manager.ID = ""

			temp := interfaces.ObjectBaseInfo{
				ID:   "ID1",
				Name: "Name1",
				Type: "department",
			}

			temp1 := []interfaces.ObjectBaseInfo{temp}
			basInfo.ParentDeps = [][]interfaces.ObjectBaseInfo{temp1}

			usersInfo := []interfaces.UserBaseInfo{basInfo}
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(usersInfo, nil)

			req := httptest.NewRequest("POST", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []UserBaseDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas[0].Roles, []string{EnumSuperAdmin, EnumSysAdmin, EnumAuditAdmin, EnumSecAdmin, EnumOrgManager, EnumOrgAudit, EnumNormaluser})
			assert.Equal(t, strParmas[0].Enabled, basInfo.Enabled)
			assert.Equal(t, strParmas[0].Priority, basInfo.Priority)
			assert.Equal(t, strParmas[0].CsfLevel, basInfo.CSFLevel)
			assert.Equal(t, strParmas[0].Name, basInfo.Name)
			assert.Equal(t, strParmas[0].ParentDeps, basInfo.ParentDeps)
			assert.Equal(t, strParmas[0].Email, basInfo.Email)
			assert.Equal(t, strParmas[0].Telephone, basInfo.TelNumber)
			assert.Equal(t, strParmas[0].ThirdAttr, basInfo.ThirdAttr)
			assert.Equal(t, strParmas[0].ThirdID, basInfo.ThirdID)
			assert.Equal(t, strParmas[0].AuthType, EnumLocal)
			assert.Equal(t, strParmas[0].Manager, nil)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetUserInfoByAccount(t *testing.T) {
	Convey("通过账户名匹配账户信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, nil, nil, nil)
		testURestHandler.RegisterPrivate(r)

		var user interfaces.UserBaseInfo
		target := "/api/user-management/v1/account-match"
		Convey("账户匹配失败--缺少account请求参数", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("账户匹配失败--服务内部错误", func() {
			tmpErr := fmt.Errorf("xx err")
			userLogics.EXPECT().GetUserInfoByAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, user, tmpErr)

			req := httptest.NewRequest("GET", target+"?account=1", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("账户匹配失败--账户不存在", func() {
			userLogics.EXPECT().GetUserInfoByAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, user, nil)

			req := httptest.NewRequest("GET", target+"?account=1&id_card_login=true&prefix_match=true", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, res.(map[string]interface{})["result"].(bool), false)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("账户匹配成功", func() {
			user.ID = "003a75b4-1f0d-4a7e-929c-497a41b9037c"
			user.AuthType = interfaces.Local
			user.PwdErrCnt = 0
			user.PwdErrLastTime = time.Now().Unix()
			user.Enabled = true
			user.LDAPType = interfaces.LDAPServerType(0)
			user.DomainPath = "xx"
			userLogics.EXPECT().GetUserInfoByAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, user, nil)

			req := httptest.NewRequest("GET", target+"?account=1&id_card_login=true&prefix_match=true", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, res.(map[string]interface{})["result"].(bool), true)
			user := res.(map[string]interface{})["user"]
			assert.Equal(t, user.(map[string]interface{})["id"].(string), "003a75b4-1f0d-4a7e-929c-497a41b9037c")
			assert.Equal(t, int(user.(map[string]interface{})["pwd_err_cnt"].(float64)), 0)
			assert.Equal(t, int64(user.(map[string]interface{})["pwd_err_last_time"].(float64)), time.Now().Unix())
			assert.Equal(t, user.(map[string]interface{})["disable_status"].(bool), false)
			assert.Equal(t, user.(map[string]interface{})["ldap_server_type"].(string), "")
			assert.Equal(t, user.(map[string]interface{})["domain_path"].(string), "xx")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUserAuth(t *testing.T) {
	Convey("本地认证-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, nil, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/user-auth"
		Convey("请求参数缺失--id", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("请求参数缺失--password", func() {
			req := httptest.NewRequest("GET", target+"?id=xx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("服务内部错误", func() {
			tmpErr := fmt.Errorf("xx err")
			userLogics.EXPECT().UserAuth(gomock.Any(), gomock.Any()).Return(false, interfaces.AuthFailedReason(0), tmpErr)

			req := httptest.NewRequest("GET", target+"?id=xx&password=xx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("认证失败--密码错误", func() {
			userLogics.EXPECT().UserAuth(gomock.Any(), gomock.Any()).Return(false, interfaces.InvalidPassword, nil)

			req := httptest.NewRequest("GET", target+"?id=1&password=xx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, res.(map[string]interface{})["result"].(bool), false)
			assert.Equal(t, res.(map[string]interface{})["reason"].(string), authFailedReasonEnumMap[interfaces.InvalidPassword])

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("认证成功", func() {
			userLogics.EXPECT().UserAuth(gomock.Any(), gomock.Any()).Return(true, interfaces.AuthFailedReason(0), nil)

			req := httptest.NewRequest("GET", target+"?id=1&password=xx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, result.StatusCode, http.StatusOK)
			assert.Equal(t, res.(map[string]interface{})["result"].(bool), true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUpdatePwdErrInfo(t *testing.T) {
	Convey("更新密码错误信息-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, nil, nil, nil)
		testURestHandler.RegisterPrivate(r)

		tmpBody := make(map[string]interface{})
		target := "/api/user-management/v1/users/user_id_xxx/pwd_err_info"
		Convey("请求参数缺失--pwd_err_cnt", func() {
			tmpBody["pwd_err_last_time"] = time.Now().Unix()
			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("请求参数缺失--pwd_err_last_time", func() {
			tmpBody["pwd_err_cnt"] = 1

			reqBody, _ := jsoniter.Marshal(tmpBody)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("更新失败--内部错误", func() {
			tmpBody["pwd_err_cnt"] = 1
			tmpBody["pwd_err_last_time"] = time.Now().Unix()

			reqBody, _ := jsoniter.Marshal(tmpBody)

			userLogics.EXPECT().UpdatePwdErrInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("xx"))

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		Convey("更新成功", func() {
			tmpBody["pwd_err_cnt"] = 1
			tmpBody["pwd_err_last_time"] = time.Now().Unix()

			reqBody, _ := jsoniter.Marshal(tmpBody)

			userLogics.EXPECT().UpdatePwdErrInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
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

func TestBatchGetUserBaseInfo(t *testing.T) {
	Convey("批量获取用户信息-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/users/xasdasd,zzzz/xxxx"
		Convey("fileds包含规定外的字段-报错", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		target = "/api/user-management/v1/users/:user_id,xxx/roles"
		Convey("用户不存在-报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			var basInfo interfaces.UserBaseInfo
			testInfo := []interfaces.UserBaseInfo{basInfo}
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testInfo, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 503)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		target = "/api/user-management/v1/users/:user_id,xxx/roles,priority,csf_level,enabled,name,account,frozen,authenticated,auth_type"
		Convey("获取所有信息-成功", func() {
			basInfo := interfaces.UserBaseInfo{}
			basInfo.VecRoles = make([]interfaces.Role, 0)
			basInfo.VecRoles = append(basInfo.VecRoles, interfaces.SystemRoleSuperAdmin, interfaces.SystemRoleSysAdmin, interfaces.SystemRoleAuditAdmin,
				interfaces.SystemRoleSecAdmin, interfaces.SystemRoleOrgManager, interfaces.SystemRoleOrgAudit, interfaces.SystemRoleNormalUser)
			basInfo.CSFLevel = 22
			basInfo.Enabled = true
			basInfo.Priority = 324
			basInfo.Name = strName111
			basInfo.Account = "kkkk"
			basInfo.Frozen = true
			basInfo.Authenticated = true
			basInfo.AuthType = interfaces.Domain

			temp := interfaces.ObjectBaseInfo{
				ID:   "ID1",
				Name: "Name1",
				Type: "department",
			}

			temp1 := []interfaces.ObjectBaseInfo{temp}
			basInfo.ParentDeps = [][]interfaces.ObjectBaseInfo{temp1}

			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]interfaces.UserBaseInfo{basInfo}, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []UserBaseDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas[0].Roles, []string{EnumSuperAdmin, EnumSysAdmin, EnumAuditAdmin, EnumSecAdmin, EnumOrgManager, EnumOrgAudit, EnumNormaluser})
			assert.Equal(t, strParmas[0].Enabled, basInfo.Enabled)
			assert.Equal(t, strParmas[0].Priority, basInfo.Priority)
			assert.Equal(t, strParmas[0].CsfLevel, basInfo.CSFLevel)
			assert.Equal(t, strParmas[0].Name, basInfo.Name)
			assert.Equal(t, strParmas[0].Account, basInfo.Account)
			assert.Equal(t, strParmas[0].Frozen, basInfo.Frozen)
			assert.Equal(t, strParmas[0].Authenticated, basInfo.Authenticated)
			assert.Equal(t, strParmas[0].ID, basInfo.ID)
			assert.Equal(t, strParmas[0].AuthType, EnumDomain)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetUserBaseInfoInScope(t *testing.T) {
	Convey("获取用户信息-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/users/xasdasd,xxx"

		Convey("token失效-报错", func() {
			tempTarget := target + "/name"
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
			target += "/name1?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("角色参数错误-报错", func() {
			target += "/name?role=sys_admin1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("角色不传-报错", func() {
			target += "/name"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("内部错误-报错", func() {
			target += "/name?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
			var basInfo []interfaces.UserBaseInfo
			userLogics.EXPECT().GetUserBaseInfoInScope(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(basInfo, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 503)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取所有参数信息-成功", func() {
			target += "/name,account,parent_dep_paths?role=sys_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			temp := interfaces.UserBaseInfo{
				ID:             strXXX,
				Name:           "Name1",
				Account:        "Account1",
				ParentDepPaths: []string{"xxx/xxxx"},
			}

			basInfo := []interfaces.UserBaseInfo{temp}
			userLogics.EXPECT().GetUserBaseInfoInScope(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(basInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := []UserBaseDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			for _, v := range strParmas {
				assert.Equal(t, v.ParentDepPaths, temp.ParentDepPaths)
				assert.Equal(t, v.Account, temp.Account)
				assert.Equal(t, v.Name, temp.Name)
			}

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUpdateUserInfo(t *testing.T) {
	Convey("修改用户信息-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/management/users/xasdasd/password"

		Convey("token失效-报错", func() {
			tempTarget := target + ",password1"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
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
			VisitorID: "department_id",
		}

		Convey("fields参数错误-报错", func() {
			target += ",password1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("请求报文中不包含必填信息-报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("PUT", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tmpBody := map[string]interface{}{
			"password": "some-secret",
		}
		reqBody, _ := jsoniter.Marshal(tmpBody)

		Convey("修改信息失败-报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 404019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().ModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNotFound)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("修改信息成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().ModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("PUT", target, bytes.NewReader(reqBody))
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

func TestGetNorlmalUserInfo(t *testing.T) {
	Convey("普通用户获取自己信息-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/profile"

		Convey("token失效-报错", func() {
			tempTarget := target + "/avatar_url"
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
			tmptarget := target + "/password1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("GET", tmptarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
		userInfo := interfaces.UserBaseInfo{
			Avatar: "ur;l:xxxx",
		}
		Convey("GetUsersBaseInfo-报错", func() {
			tmptarget := target + "/avatar_url"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetNorlmalUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)

			req := httptest.NewRequest("GET", tmptarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			tmptarget := target + "/avatar_url"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetNorlmalUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

			req := httptest.NewRequest("GET", tmptarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := UserOwnDataMock{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas.Avatar, userInfo.Avatar)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:dupl
func TestUpdateAvatar(t *testing.T) {
	Convey("更新个人头像-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		av := mock.NewMockLogicsAvatar(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, av, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/profile/avatar"

		Convey("token失效-报错", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
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

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    true,
			VisitorID: "department_id",
		}

		tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
		userInfo := make([]interfaces.UserBaseInfo, 1)
		Convey("非表单上传 报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)

			req := httptest.NewRequest("POST", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("表单内不存在file键值对", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)
			// 创建表单发送
			contentType, rd, err := createMultiPartRequest(true, false)
			assert.Equal(t, err, nil)

			req := httptest.NewRequest("POST", target, rd)
			req.Header.Set("Content-Type", contentType)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("表单内file文件超过1", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)
			// 创建表单发送
			contentType, rd, err := createMultiPartRequest(false, true)
			assert.Equal(t, err, nil)

			req := httptest.NewRequest("POST", target, rd)
			req.Header.Set("Content-Type", contentType)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)
			av.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			// 创建表单发送
			contentType, rd, err := createMultiPartRequest(false, false)
			assert.Equal(t, err, nil)

			req := httptest.NewRequest("POST", target, rd)
			req.Header.Set("Content-Type", contentType)
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

func createMultiPartRequest(bErrFile, bSecondFIle bool) (contentType string, w io.Reader, err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 在表单中创建一个文件字段
	h := make(textproto.MIMEHeader)
	if bErrFile {
		h.Set("Content-Disposition",
			`form-data; name="avatar1"; filename="xx.png"`)
	} else {
		h.Set("Content-Disposition",
			`form-data; name="avatar"; filename="xx.png"`)
	}

	formFile, err := writer.CreatePart(h)
	if err != nil {
		return "", nil, err
	}

	// 文件拷贝
	buffer := make([]byte, 48*1024-1)
	_, err = io.Copy(formFile, bytes.NewReader(buffer))
	if err != nil {
		return "", nil, err
	}

	if bSecondFIle {
		formFile, err := writer.CreatePart(h)
		if err != nil {
			return "", nil, err
		}

		// 文件拷贝
		buffer := make([]byte, 100)
		_, err = io.Copy(formFile, bytes.NewReader(buffer))
		if err != nil {
			return "", nil, err
		}
	}

	// 关闭
	if err := writer.Close(); err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body, nil
}

func TestGetAvatarsByID(t *testing.T) {
	Convey("根据用户ID获取头像-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		av := mock.NewMockLogicsAvatar(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, av, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/avatars/xxx"

		introspectInfo := interfaces.TokenIntrospectInfo{
			Active:    false,
			VisitorID: "department_id",
		}
		Convey("token失效-报错", func() {
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

		introspectInfo.Active = true
		tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
		userInfo := make([]interfaces.UserBaseInfo, 1)
		Convey("GetUsersBaseInfo报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, tempErr)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		userInfo1 := interfaces.UserBaseInfo{
			ID:     "xxx1",
			Avatar: "url1",
		}
		userInfo[0] = userInfo1
		Convey("成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUsersBaseInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfo, nil)

			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			var strParmas []interface{}
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, len(strParmas), 1)

			outData, ok := strParmas[0].(map[string]interface{})
			assert.Equal(t, ok, true)
			assert.Equal(t, outData["id"], userInfo1.ID)
			assert.Equal(t, outData["avatar_url"], userInfo1.Avatar)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetPWDRetrievalMethod(t *testing.T) {
	Convey("根据账户名获取用户找回密码方式", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/pwd-retrieval-method"

		Convey("query内不包含account", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("xxxx", rest.Forbidden, nil)
		var info interfaces.PwdRetrievalInfo
		Convey("GetPWDRetrievalMethodByAccount报错", func() {
			userLogics.EXPECT().GetPWDRetrievalMethodByAccount(gomock.Any()).AnyTimes().Return(info, testErr)

			testURL := target + "?account=xxx"
			req := httptest.NewRequest("GET", testURL, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		info.BEmail = true
		info.BTelephone = true
		info.Email = "email1"
		info.Telephone = "telephone1"
		info.Status = interfaces.PRSAvaliable
		Convey("成功", func() {
			userLogics.EXPECT().GetPWDRetrievalMethodByAccount(gomock.Any()).AnyTimes().Return(info, nil)

			testURL := target + "?account=xxx"
			req := httptest.NewRequest("GET", testURL, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam["status"], "available")
			assert.Equal(t, respParam["email"], info.Email)
			assert.Equal(t, respParam["telephone"], info.Telephone)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetOrgManagersInfo(t *testing.T) {
	Convey("获取组织管理员信息", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		roleLogics := mock.NewMockLogicsRole(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, roleLogics)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/org_managers/xxxx/"

		Convey("url内不包含sub_user_ids", func() {
			testURL := target + "roles"
			req := httptest.NewRequest("GET", testURL, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, nil, err)
			assert.Equal(t, respParam.Cause, "invalid type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("xxxx", rest.Forbidden, nil)
		Convey("GetOrgManagersInfo报错", func() {
			roleLogics.EXPECT().GetOrgManagersInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)

			testURL := target + "sub_user_ids"
			req := httptest.NewRequest("GET", testURL, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
		info := []interfaces.OrgManagerInfo{
			{
				ID:         "test",
				SubUserIDs: []string{"test1", "test2"},
			},
			{
				ID:         "test1",
				SubUserIDs: []string{},
			},
		}

		Convey("成功", func() {
			roleLogics.EXPECT().GetOrgManagersInfo(gomock.Any(), gomock.Any()).AnyTimes().Return(info, nil)

			testURL := target + "sub_user_ids"
			req := httptest.NewRequest("GET", testURL, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			respParam := make([]map[string]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestIncrementUpdateUserInfo(t *testing.T) {
	Convey("增量修改用户信息-外部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		common.InitARTrace("test")
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/management/users/xasdasd"

		Convey("token失效-报错", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("PATCH", target, http.NoBody)
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

		Convey("请求报文中不包含必填信息-报错", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest("PATCH", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tmpBody := map[string]interface{}{
			"custom_attr": map[string]interface{}{"test": "tset"},
		}
		reqBody, _ := jsoniter.Marshal(tmpBody)

		Convey("修改信息失败-报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 404019001, nil)
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().IncrementModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)

			req := httptest.NewRequest("PATCH", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNotFound)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("修改信息成功", func() {
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().IncrementModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("PATCH", target, bytes.NewReader(reqBody))
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

func TestIncrementUpdateUserInfoInternal(t *testing.T) {
	Convey("增量修改用户信息-内部接口-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		common.InitARTrace("test")
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/users/xasdasd"

		Convey("请求报文中不包含必填信息-报错", func() {
			req := httptest.NewRequest("PATCH", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tmpBody := map[string]interface{}{
			"custom_attr": map[string]interface{}{"test": "tset"},
		}
		reqBody, _ := jsoniter.Marshal(tmpBody)

		Convey("修改信息失败-报错", func() {
			tempErr := rest.NewHTTPError("用户不存在", 404019001, nil)
			userLogics.EXPECT().IncrementModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tempErr)

			req := httptest.NewRequest("PATCH", target, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusNotFound)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("修改信息成功", func() {
			userLogics.EXPECT().IncrementModifyUserInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			req := httptest.NewRequest("PATCH", target, bytes.NewReader(reqBody))
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

func TestGetUserBaseInfoByTelephones(t *testing.T) {
	Convey("根据手机号取用户信息POST-内部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/query-user-by-telephone"
		Convey("token失效-报错", func() {
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, target, http.NoBody)
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

		Convey("fileds包含规定外的字段-报错", func() {
			tempTarget := target + "?telephone=xxx&field=xx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("query不包含fields-报错", func() {
			tempTarget := target + "?telephone=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "invalid field")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("telephones不存在-报错", func() {
			tempTarget := target + "?telephone="
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 400)

			respBody, _ := io.ReadAll(result.Body)
			var res interface{}
			_ = jsoniter.Unmarshal(respBody, &res)

			assert.Equal(t, res.(map[string]interface{})["cause"], "invalid telephone")
			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		tempErr := rest.NewHTTPError("用户不存在", 503000000, nil)
		Convey("用户不存在-报错", func() {
			tempTarget := target + "?telephone=xxx&field=email"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			usersInfo := interfaces.UserBaseInfo{}
			userLogics.EXPECT().GetUserBaseInfoByTelephone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, usersInfo, tempErr)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, 503)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取所有信息-成功", func() {
			tempTarget := target + "?telephone=xxx&field=third_id&field=account&field=name&field=email"

			basInfo := interfaces.UserBaseInfo{}
			basInfo.Name = strName111
			basInfo.Email = strEmail1
			basInfo.TelNumber = "123456"
			basInfo.ThirdID = strThirdID1
			basInfo.ID = strUserID
			basInfo.Account = strAccount

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUserBaseInfoByTelephone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, basInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["result"].(bool), true)
			data, ok := strParmas["user"].(map[string]interface{})
			assert.Equal(t, ok, true)
			assert.Equal(t, data["id"].(string), basInfo.ID)
			assert.Equal(t, data["telephone"].(string), basInfo.TelNumber)
			assert.Equal(t, data["account"].(string), basInfo.Account)
			assert.Equal(t, data["name"].(string), basInfo.Name)
			assert.Equal(t, data["email"].(string), basInfo.Email)
			assert.Equal(t, data["third_id"].(string), basInfo.ThirdID)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("获取所有信息-成功1", func() {
			tempTarget := target + "?telephone=xxx&field=third_id&field=account&field=name&field=email"

			basInfo := interfaces.UserBaseInfo{}
			basInfo.Name = strName111
			basInfo.Email = strEmail1
			basInfo.TelNumber = "123456"
			basInfo.ThirdID = strThirdID1
			basInfo.ID = strUserID
			basInfo.Account = strAccount

			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().GetUserBaseInfoByTelephone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, basInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["result"].(bool), false)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:funlen
func TestSearchUsers(t *testing.T) {
	Convey("搜索用户接口-外部-接口层", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("xxtest")

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPublic(r)

		target := "/api/user-management/v1/console/search-users/"
		Convey("token失效-报错", func() {
			tempTarget := target + "xxx"
			introspectInfo := interfaces.TokenIntrospectInfo{
				Active:    false,
				VisitorID: "department_id",
			}
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
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
		Convey("没有role", func() {
			tempTarget := target + "xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid role")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset 错误", func() {
			tempTarget := target + "xxx?offset=aa"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid offset type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("offset 错误 < 0", func() {
			tempTarget := target + "xxx?offset=-1"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid offset type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit 错误", func() {
			tempTarget := target + "xxx?limit=aa"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit 错误 > 1000", func() {
			tempTarget := target + "xxx?limit=1001"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid limit type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("role 错误 不是枚举", func() {
			tempTarget := target + "xxx?limit=101&role=xxx"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid role")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("fields 错误 不是枚举", func() {
			tempTarget := target + "xxx?limit=101&role=super_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), "invalid type")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPErrorV2(rest.BadRequest, "testerr")
		Convey("SearchUsers 错误 不是枚举", func() {
			tempTarget := target + "name?limit=101&role=super_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().SearchUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, testErr)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, strParmas["cause"].(string), testErr.Cause)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		num := 20
		tempxxx1 := interfaces.ObjectBaseInfo{
			ID:   "sub1",
			Name: "name1",
			Type: "user",
			Code: "code1",
		}
		temp1 := []interfaces.ObjectBaseInfo{
			tempxxx1,
		}
		tempxx2 := interfaces.ObjectBaseInfo{
			ID:   "sub2",
			Name: "name2",
			Type: "user",
			Code: "code2",
		}
		temp2 := []interfaces.ObjectBaseInfo{
			tempxx2,
		}
		testUser1 := interfaces.UserBaseInfo{
			ID:      strThirdID,
			Name:    strName1,
			Account: strAccount,
			Remark:  strRemark1,
			ParentDeps: [][]interfaces.ObjectBaseInfo{
				temp1,
				temp2,
			},
			Manager: interfaces.NameInfo{
				ID:   "user2",
				Name: "name2",
			},
			VecRoles:  []interfaces.Role{interfaces.SystemRoleOrgAudit, interfaces.SystemRoleNormalUser},
			CSFLevel:  12,
			AuthType:  interfaces.Local,
			Priority:  14,
			CreatedAt: 111111,
			Enabled:   true,
			Code:      "code1",
			Position:  "Position1",
			Frozen:    true,
			CSFLevel2: 13,
		}

		Convey("success", func() {
			tempTarget := target + "account,name,remark,parent_deps,roles,csf_level,auth_type,priority,created_at,enabled,code,manager,position,frozen,csf_level2?limit=101&role=super_admin"
			h.EXPECT().Introspect(gomock.Any()).AnyTimes().Return(introspectInfo, nil)
			userLogics.EXPECT().SearchUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]interfaces.UserBaseInfo{testUser1}, num, nil)

			req := httptest.NewRequest(mstrGET, tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, result.StatusCode, http.StatusOK)

			respBody, _ := io.ReadAll(result.Body)
			strParmas := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &strParmas)
			assert.Equal(t, err, nil)
			assert.Equal(t, int(strParmas["total_count"].(float64)), num)

			data1 := strParmas["entries"].([]interface{})
			assert.Equal(t, len(data1), 1)

			data2 := data1[0].(map[string]interface{})
			assert.Equal(t, data2["id"].(string), testUser1.ID)
			assert.Equal(t, data2["name"].(string), testUser1.Name)
			assert.Equal(t, data2["account"].(string), testUser1.Account)
			assert.Equal(t, data2["type"].(string), strUser)
			assert.Equal(t, int(data2["csf_level"].(float64)), testUser1.CSFLevel)
			assert.Equal(t, data2["auth_type"].(string), "local")
			assert.Equal(t, int(data2["priority"].(float64)), testUser1.Priority)
			assert.Equal(t, data2["created_at"], "1970-01-02T14:51:51+08:00")
			assert.Equal(t, data2["enabled"].(bool), testUser1.Enabled)
			assert.Equal(t, data2["code"], testUser1.Code)
			assert.Equal(t, data2["position"], testUser1.Position)
			assert.Equal(t, data2["frozen"].(bool), testUser1.Frozen)
			assert.Equal(t, int(data2["csf_level2"].(float64)), testUser1.CSFLevel2)

			testManager := data2["manager"].(map[string]interface{})
			testManager["id"] = "user2"
			testManager["name"] = "name2"
			testManager["type"] = "user"

			testRoles := data2["roles"].([]interface{})
			assert.Equal(t, len(testRoles), 2)
			assert.Equal(t, testRoles[0].(string), "org_audit")
			assert.Equal(t, testRoles[1].(string), "normal_user")

			testParentDeps := data2["parent_deps"].([]interface{})
			assert.Equal(t, len(testParentDeps), 2)
			testParent1s := testParentDeps[0].([]interface{})
			assert.Equal(t, len(testParent1s), 1)
			testParent1 := testParent1s[0].(map[string]interface{})
			assert.Equal(t, testParent1["id"], "sub1")
			assert.Equal(t, testParent1["name"], "name1")
			assert.Equal(t, testParent1["type"], strDepartment)
			assert.Equal(t, testParent1["code"], "code1")

			testParent2s := testParentDeps[1].([]interface{})
			assert.Equal(t, len(testParent2s), 1)
			testParent2 := testParent2s[0].(map[string]interface{})
			assert.Equal(t, testParent2["id"], "sub2")
			assert.Equal(t, testParent2["name"], "name2")
			assert.Equal(t, testParent2["type"], strDepartment)
			assert.Equal(t, testParent2["code"], "code2")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestUserNameExistdCheck(t *testing.T) {
	Convey("userNameExistdCheck", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)
		target := "/api/user-management/v1/user-name-existed"

		Convey("query里面没有name", func() {
			nTarget := target + "?xxxx=1"

			req := httptest.NewRequest("GET", nTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid name")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPError("xxxx", rest.Forbidden, nil)
		Convey("CheckUserNameExistd报错", func() {
			userLogics.EXPECT().CheckUserNameExistd(gomock.Any(), gomock.Any()).AnyTimes().Return(false, testErr)

			nTarget := target + "?name=1"

			req := httptest.NewRequest("GET", nTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			userLogics.EXPECT().CheckUserNameExistd(gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)

			nTarget := target + "?name=1"

			req := httptest.NewRequest("GET", nTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam["result"].(bool), true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestCheckUserNameExistd(t *testing.T) {
	Convey("检查用户名是否存在", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v1/user-name-existed"

		Convey("未传name-报错", func() {
			req := httptest.NewRequest("GET", target, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPError("error", 503000000, nil)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid name")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("CheckUserNameExistd报错", func() {
			userLogics.EXPECT().CheckUserNameExistd(gomock.Any(), gomock.Any()).AnyTimes().Return(false, errors.New("test"))

			nTarget := target + "?name=1"

			req := httptest.NewRequest("GET", nTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()

			assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("success", func() {
			userLogics.EXPECT().CheckUserNameExistd(gomock.Any(), gomock.Any()).AnyTimes().Return(true, nil)

			nTarget := target + "?name=1"

			req := httptest.NewRequest("GET", nTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam["result"].(bool), true)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

func TestGetUserList(t *testing.T) {
	Convey("用户列举", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v2/users"

		Convey("参数传递报错，fields传入无效值", func() {
			req := httptest.NewRequest("GET", target+"?fields=xxxx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "xxx")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Description, "invalid field")
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("getUseristQueryParam 报错", func() {
			testErr := gerrors.NewError(gerrors.PublicConflict, "xxx")

			userLogics.EXPECT().GetUserList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, 0, false, testErr)
			req := httptest.NewRequest("GET", target+"?fields=name", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusConflict)

			respParam := gerrors.NewError(gerrors.PublicInternalServerError, "xxx")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicConflict)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		userInfos := []interfaces.UserBaseInfo{
			{
				ID:        strID,
				Account:   strENUS,
				Name:      strAsc,
				Enabled:   true,
				Email:     strEmail,
				TelNumber: strID,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
				Frozen:    true,
			},
		}

		Convey("success", func() {
			userLogics.EXPECT().GetUserList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(userInfos, 10, true, nil)
			req := httptest.NewRequest("GET", target+"?fields=name&fields=account&fields=enabled&fields=email&fields=telephone&fields=created_at&fields=frozen&direction=asc&limit=10", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make(map[string]interface{})
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, int(respParam["total_count"].(float64)), 10)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["id"].(string), strID)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["account"].(string), strENUS)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["name"].(string), strAsc)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["enabled"].(bool), true)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["email"].(string), strEmail)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["telephone"].(string), strID)
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["created_at"].(string), "2025-01-01T00:00:00+08:00")
			assert.Equal(t, respParam["entries"].([]interface{})[0].(map[string]interface{})["frozen"].(bool), true)

			// marker
			marker := respParam["next_marker"].(string)
			markerBytes, _ := base64.StdEncoding.DecodeString(marker)
			markerJSON := make(map[string]interface{})
			err = jsoniter.Unmarshal(markerBytes, &markerJSON)
			assert.Equal(t, err, nil)
			assert.Equal(t, markerJSON["direction"].(string), "asc")
			assert.Equal(t, int64(markerJSON["created_at"].(float64)), time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local).Unix())
			assert.Equal(t, markerJSON["user_id"].(string), strID)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}

//nolint:dupl,funlen
func TestGetUseristQueryParam(t *testing.T) {
	Convey("GetUseristQueryParam", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		common.InitARTrace("test")

		userLogics := mock.NewMockLogicsUser(ctrl)
		h := mock.NewMockHydra(ctrl)
		testURestHandler := newUserRESTHandler(userLogics, h, nil, nil)
		testURestHandler.RegisterPrivate(r)

		target := "/api/user-management/v2/users"
		Convey("limit =0", func() {
			req := httptest.NewRequest("GET", target+"?limit=0", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "invalid limit([ 1 .. 1000 ])")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit =1001", func() {
			req := httptest.NewRequest("GET", target+"?limit=1001", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "invalid limit([ 1 .. 1000 ])")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("limit不为数字", func() {
			req := httptest.NewRequest("GET", target+"?limit=xxxx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "limit is illeagal")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("direction = xxxx", func() {
			req := httptest.NewRequest("GET", target+"?direction=xxxx", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "invalid direction type")
		})

		Convey("marker = xxxx", func() {
			req := httptest.NewRequest("GET", target+"?marker=aaaaaa", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "invalid marker base64")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中没有created_at", func() {
			temp := map[string]interface{}{
				"user_id":   "123",
				"direction": "asc",
			}
			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "(root): created_at is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中created_at 不是数字", func() {
			temp := map[string]interface{}{
				"created_at": "123",
				"direction":  "asc",
				"user_id":    "123",
			}
			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "created_at: Invalid type. Expected: integer, given: string")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中没有user_id", func() {
			temp := map[string]interface{}{
				"created_at": 123,
				"direction":  "asc",
			}
			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "(root): user_id is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中user_id 不是字符串", func() {
			temp := map[string]interface{}{
				"created_at": 123,
				"direction":  "asc",
				"user_id":    123,
			}

			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "user_id: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中没有direction", func() {
			temp := map[string]interface{}{
				"created_at": 123,
				"user_id":    "123",
			}
			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "(root): direction is required")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中direction 不是asc或desc", func() {
			temp := map[string]interface{}{
				"created_at": 123,
				"direction":  "xxxx",
				"user_id":    "123",
			}

			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "direction: direction must be one of the following: \"asc\", \"desc\"")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		Convey("marker 中 direction不为string", func() {
			temp := map[string]interface{}{
				"created_at": 123,
				"direction":  123,
				"user_id":    "123",
			}

			tempBytes, _ := jsoniter.Marshal(temp)
			marker := base64.StdEncoding.EncodeToString(tempBytes)
			req := httptest.NewRequest("GET", target+"?marker="+marker, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := gerrors.NewError(gerrors.PublicConflict, "")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Code, gerrors.PublicBadRequest)
			assert.Equal(t, respParam.Description, "direction: Invalid type. Expected: string, given: integer")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
