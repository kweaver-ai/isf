// Package driveradapters 角色管理测试
package driveradapters

import (
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

func newRoleRESTHandler(role interfaces.LogicsRole) *roleRestHandler {
	return &roleRestHandler{
		roleNameIDMap: map[string]interfaces.Role{
			EnumSuperAdmin: interfaces.SystemRoleSuperAdmin,
			EnumSysAdmin:   interfaces.SystemRoleSysAdmin,
			EnumAuditAdmin: interfaces.SystemRoleAuditAdmin,
			EnumSecAdmin:   interfaces.SystemRoleSecAdmin,
			EnumOrgManager: interfaces.SystemRoleOrgManager,
			EnumOrgAudit:   interfaces.SystemRoleOrgAudit,
		},
		roleIDNameMap: map[interfaces.Role]string{
			interfaces.SystemRoleSuperAdmin: EnumSuperAdmin,
			interfaces.SystemRoleSysAdmin:   EnumSysAdmin,
			interfaces.SystemRoleAuditAdmin: EnumAuditAdmin,
			interfaces.SystemRoleSecAdmin:   EnumSecAdmin,
			interfaces.SystemRoleOrgManager: EnumOrgManager,
			interfaces.SystemRoleOrgAudit:   EnumOrgAudit,
		},
		role: role,
	}
}

func TestGetRoleMembersByRoleIDs(t *testing.T) {
	Convey("获取角色成员", t, func() {
		test := setGinMode()
		defer test()
		r := gin.New()
		r.Use(gin.Recovery())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := mock.NewMockLogicsRole(ctrl)
		roleLogics := newRoleRESTHandler(o)

		common.InitARTrace("xxtest")

		roleLogics.RegisterPrivate(r)

		const target = "/api/user-management/v1/role-members/"
		Convey("role错误", func() {
			tempTarget := target + "xxxx"

			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusBadRequest)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "invalid role")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		testErr := rest.NewHTTPErrorV2(rest.Forbidden, "xxxx")
		Convey("GetUserIDsByRoleIDs error", func() {
			o.EXPECT().GetUserIDsByRoleIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, testErr)
			tempTarget := target + "super_admin,sys_admin,audit_admin,sec_admin,org_manager,org_audit"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusForbidden)

			respParam := rest.NewHTTPErrorV2(403000000, "error")
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)
			assert.Equal(t, respParam.Cause, "xxxx")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		// 测试成功
		data := make(map[interfaces.Role][]string)
		data[interfaces.SystemRoleSuperAdmin] = []string{"1"}
		data[interfaces.SystemRoleSysAdmin] = []string{"4", "5"}
		data[interfaces.SystemRoleAuditAdmin] = []string{"7"}
		data[interfaces.SystemRoleOrgManager] = nil
		data[interfaces.SystemRoleOrgAudit] = []string{}
		Convey("success", func() {
			o.EXPECT().GetUserIDsByRoleIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(data, nil)
			tempTarget := target + "super_admin,sys_admin,audit_admin,sec_admin,org_manager,org_audit"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			assert.Equal(t, len(respParam), 6)
			tempData := make(map[string]interface{})
			for _, v := range respParam {
				tempData[v.(map[string]interface{})["role"].(string)] = v.(map[string]interface{})["members"]
			}

			assert.Equal(t, len(tempData[EnumSuperAdmin].([]interface{})), 1)
			assert.Equal(t, len(tempData[EnumSysAdmin].([]interface{})), 2)
			assert.Equal(t, len(tempData[EnumAuditAdmin].([]interface{})), 1)
			assert.Equal(t, len(tempData[EnumOrgManager].([]interface{})), 0)
			assert.Equal(t, len(tempData[EnumOrgAudit].([]interface{})), 0)
			assert.Equal(t, len(tempData[EnumSecAdmin].([]interface{})), 0)

			assert.Equal(t, tempData[EnumSuperAdmin].([]interface{})[0].(map[string]interface{})["id"], "1")
			assert.Equal(t, tempData[EnumSysAdmin].([]interface{})[0].(map[string]interface{})["id"], "4")
			assert.Equal(t, tempData[EnumSysAdmin].([]interface{})[1].(map[string]interface{})["id"], "5")
			assert.Equal(t, tempData[EnumAuditAdmin].([]interface{})[0].(map[string]interface{})["id"], "7")

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})

		// 测试成功, 去重检查
		data = make(map[interfaces.Role][]string)
		data[interfaces.SystemRoleSuperAdmin] = []string{"1"}
		Convey("success1", func() {
			o.EXPECT().GetUserIDsByRoleIDs(gomock.Any(), gomock.Any()).AnyTimes().Return(data, nil)
			tempTarget := target + "super_admin,super_admin,super_admin,super_admin,super_admin,super_admin"
			req := httptest.NewRequest("GET", tempTarget, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			result := w.Result()
			respBody, _ := io.ReadAll(result.Body)

			assert.Equal(t, result.StatusCode, http.StatusOK)

			respParam := make([]interface{}, 0)
			err := jsoniter.Unmarshal(respBody, &respParam)
			assert.Equal(t, err, nil)

			assert.Equal(t, len(respParam), 1)

			if err := result.Body.Close(); err != nil {
				assert.Equal(t, err, nil)
			}
		})
	})
}
