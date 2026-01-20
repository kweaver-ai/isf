// Package driveradapters department AnyShare  部门逻辑接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// DepartRestHandler depart RESTfual API Handler 接口
type DepartRestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册公共API
	RegisterPublic(engine *gin.Engine)
}
type departRestHandler struct {
	depart              interfaces.LogicsDepartment
	hydra               interfaces.Hydra
	getDepartInfoSchema *gojsonschema.Schema
}

const (
	strParentDeps = "parent_deps"
	strDepartment = "department"
	strManager    = "manager"
	strCode       = "code"
	strEnabled    = "enabled"
	strRemark     = "remark"
)

var (
	dOnce    sync.Once
	dHandler DepartRestHandler

	//go:embed jsonschema/depart/get_depart_info.json
	getDepartInfoSchemaStr string
)

var roleManageEnumIDMap = map[string]interfaces.Role{
	EnumSuperAdmin: interfaces.SystemRoleSuperAdmin,
	EnumSysAdmin:   interfaces.SystemRoleSysAdmin,
	EnumSecAdmin:   interfaces.SystemRoleSecAdmin,
	EnumOrgManager: interfaces.SystemRoleOrgManager,
}

// NewDepartRESTHandler depart handler 对象
func NewDepartRESTHandler() DepartRestHandler {
	dOnce.Do(func() {
		getDepartInfoSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getDepartInfoSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		dHandler = &departRestHandler{
			depart:              logics.NewDepartment(),
			hydra:               newHydra(),
			getDepartInfoSchema: getDepartInfoSchema,
		}
	})

	return dHandler
}

// RegisterPrivate 注册内部API
func (h *departRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/user-management/v1/departments/:department_id/:fields", h.getDepartFieldInfo)
	engine.GET("/api/user-management/v1/departments", h.getDepartInfoByLevel)
	engine.DELETE("/api/user-management/v1/departments/:department_id", h.deleteDepart)

	engine.POST("/api/user-management/v1/batch-get-department-info", observable.MiddlewareTrace(common.SvcARTrace), h.getDepartFieldInfoByPost)
}

// RegisterPublic 注册公共API
func (h *departRestHandler) RegisterPublic(engine *gin.Engine) {
	// 弹窗相关
	engine.GET("/api/user-management/v1/department-members/:department_id/:fields", h.getDepartMemberInfo)

	// 管理控制台获取根部门信息
	engine.GET("/api/user-management/v1/management/department-members/:department_id/:fields", h.getManageDepartMemberInfo)

	// 管理控制台删除部门
	engine.DELETE("/api/user-management/v1/management/departments/:department_id", h.deleteManageDepart)

	// 管理控制台搜索部门
	engine.GET("/api/user-management/v1/console/search-departments/:fields", observable.MiddlewareTrace(common.SvcARTrace), h.searchManageDepart)
}

// getDepartFieldInfo 按照参数获取部门信息
func (h *departRestHandler) getDepartFieldInfo(c *gin.Context) {
	fields := c.Param("fields")

	switch fields {
	case "accessor_ids":
		h.getAccessorIDsOfDepart(c)
	case "member_ids":
		h.getDepartMemberIDs(c)
	case "all_user_ids":
		h.getDepartAllUserIDs(c)
	case "all_users":
		h.getDepartAllUserInfos(c)
	default:
		h.getDepartInfo(c)
	}
}

// getDepartAllUserInfos 获取指定部门的所有用户信息
func (h *departRestHandler) getDepartAllUserInfos(c *gin.Context) {
	departID := c.Param("department_id")

	// 获取指定部门的成员ID
	infos, err := h.depart.GetAllDepartUserInfos(departID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make([]interface{}, 0)
	for k := range infos {
		temp := make(map[string]interface{})
		temp["id"] = infos[k].ID
		temp["name"] = infos[k].Name
		temp["account"] = infos[k].Account
		temp["email"] = infos[k].Email
		temp["telephone"] = infos[k].TelNumber
		temp[strThirdAttr] = infos[k].ThirdAttr
		temp["third_id"] = infos[k].ThirdID

		out = append(out, temp)
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// getAccessorIDsOfDepart 获取指定部门的访问令牌
func (h *departRestHandler) getAccessorIDsOfDepart(c *gin.Context) {
	departID := c.Param("department_id")

	// 获取指定部门的访问令牌
	accessorIDs, err := h.depart.GetAccessorIDsOfDepartment(departID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusOK, accessorIDs)
}

// getDepartMemberIDs 获取指定部门的成员ID
func (h *departRestHandler) getDepartMemberIDs(c *gin.Context) {
	departID := c.Param("department_id")

	// 获取指定部门的成员ID
	info, err := h.depart.GetDepartMemberIDs(departID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusOK, info)
}

// getDepartAllUserIDs 获取指定部门的所有用户ID信息
func (h *departRestHandler) getDepartAllUserIDs(c *gin.Context) {
	departID := c.Param("department_id")

	// 参数获取，是否只显示启用用户
	bOnlyShowEnableUser := false
	bOnlyShowEnableUserStr, ok := c.GetQuery("user_enabled")
	if ok {
		var err error
		bOnlyShowEnableUser, err = strconv.ParseBool(bOnlyShowEnableUserStr)
		if err != nil || !bOnlyShowEnableUser {
			err := rest.NewHTTPError("invalid user_enabled", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取指定部门的成员ID
	info, err := h.depart.GetAllDepartUserIDs(departID, !bOnlyShowEnableUser)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make(map[string]interface{})
	out["all_user_ids"] = info

	rest.ReplyOK(c, http.StatusOK, out)
}

// getManageDepartMemberInfo
func (h *departRestHandler) getManageDepartMemberInfo(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数获取
	departID, info, err := h.getConsoleDepartMembersParam(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	info.BShowSubDepart = true
	info.BShowDepartManager = true
	info.BShowDepartRemark = true
	info.BShowDepartCode = true
	info.BShowDepartEmail = true
	info.BShowDepartEnabled = true
	info.BShowDepartParentDeps = true

	// 获取指定部门的成员ID
	depInfo, depNum, _, _, err := h.depart.GetDepartMemberInfo(&visitor, departID, info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数info
	outInfo := make(map[string]interface{})
	if info.BShowDeparts {
		departParams := make(map[string]interface{})
		entriesParm := make([]interface{}, 0)

		for k := range depInfo {
			depInfoParams := make(map[string]interface{})
			depInfoParams["id"] = depInfo[k].ID
			depInfoParams["name"] = depInfo[k].Name
			depInfoParams["type"] = strDepartment
			depInfoParams["is_root"] = depInfo[k].IsRoot
			depInfoParams["depart_existed"] = depInfo[k].BDepartExistd
			depInfoParams["email"] = depInfo[k].Email
			depInfoParams[strRemark] = depInfo[k].Remark
			depInfoParams["code"] = depInfo[k].Code
			depInfoParams["enabled"] = depInfo[k].Enabled

			if depInfo[k].Manager.ID != "" {
				managerInfo := make(map[string]interface{})
				managerInfo["id"] = depInfo[k].Manager.ID
				managerInfo["name"] = depInfo[k].Manager.Name
				managerInfo["type"] = strUser
				depInfoParams["manager"] = managerInfo
			}

			parentDeps := make([]interface{}, 0)
			for k1 := range depInfo[k].ParentDeps {
				parentDep := make(map[string]interface{})
				parentDep["id"] = depInfo[k].ParentDeps[k1].ID
				parentDep["name"] = depInfo[k].ParentDeps[k1].Name
				parentDep["type"] = strDepartment
				parentDep["code"] = depInfo[k].ParentDeps[k1].Code
				parentDeps = append(parentDeps, parentDep)
			}
			depInfoParams["parent_deps"] = parentDeps

			entriesParm = append(entriesParm, depInfoParams)
		}

		departParams["total_count"] = depNum
		departParams["entries"] = entriesParm
		outInfo["departments"] = departParams
	}

	// 数据整理
	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// deleteDepart 删除部门
func (h *departRestHandler) deleteDepart(c *gin.Context) {
	// 参数获取
	departID := c.Param("department_id")
	err := h.depart.DeleteDepart(nil, departID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

func (h *departRestHandler) searchManageDepart(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	fields := strings.Split(c.Param("fields"), ",")
	scope, err := h.handleSearchDepartsFields(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	strRole, dRet := c.GetQuery("role")
	role, ok := roleEnumIDMap[strRole]
	if !dRet || !ok {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid params role")
		rest.ReplyError(c, err)
		return
	}

	var strEnabledData string
	ks := interfaces.DepartSearchKeyScope{}
	k := interfaces.DepartSearchKey{}
	k.Name, ks.BName = c.GetQuery("name")
	k.Code, ks.BCode = c.GetQuery("code")
	k.Remark, ks.BRemark = c.GetQuery(strRemark)
	k.DirectDepartCode, ks.BDirectDepartCode = c.GetQuery("direct_department_code")
	strEnabledData, ks.BEnabled = c.GetQuery("enabled")
	if ks.BEnabled {
		k.Enabled, err = strconv.ParseBool(strEnabledData)
		if err != nil {
			err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid enabled parameter")
			rest.ReplyError(c, err)
			return
		}
	}
	k.ManagerName, ks.BManagerName = c.GetQuery("manager_name")

	k.Offset, k.Limit, err = h.getLimitInfo(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	infos, num, err := h.depart.SearchDeparts(c, &visitor, &scope, &ks, &k, role)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make(map[string]interface{})
	out["total_count"] = num
	tempEntries := make([]interface{}, 0)
	for k := range infos {
		temp := make(map[string]interface{})

		temp["id"] = infos[k].ID
		temp["type"] = strDepartment
		handleSearchDepartsData(&infos[k], fields, temp)
		tempEntries = append(tempEntries, temp)
	}
	out["entries"] = tempEntries

	rest.ReplyOK(c, http.StatusOK, out)
}

func handleSearchDepartsData(info *interfaces.DepartInfo, fields []string, data map[string]interface{}) {
	for _, v := range fields {
		switch v {
		case strName:
			data["name"] = info.Name
		case strRemark:
			data[strRemark] = info.Remark
		case strEmail:
			data["email"] = info.Email
		case strParentDeps:
			parentDeps := make([]interface{}, 0)
			for _, v := range info.ParentDeps {
				parentDep := make(map[string]interface{})
				parentDep["id"] = v.ID
				parentDep["name"] = v.Name
				parentDep["code"] = v.Code
				parentDep["type"] = v.Type
				parentDeps = append(parentDeps, parentDep)
			}
			data["parent_deps"] = parentDeps
		case strManager:
			if info.Manager.ID != "" {
				manager := make(map[string]interface{})
				manager["id"] = info.Manager.ID
				manager["name"] = info.Manager.Name
				manager["type"] = strUser
				data["manager"] = manager
			}
		case strCode:
			data["code"] = info.Code
		case strEnabled:
			data["enabled"] = info.Enabled
		}
	}
}

func (h *departRestHandler) handleSearchDepartsFields(fields []string) (scope interfaces.DepartInfoScope, err error) {
	for _, v := range fields {
		switch v {
		case strName:
			scope.BShowName = true
		case strRemark:
			scope.BRemark = true
		case strEmail:
			scope.BEmail = true
		case strParentDeps:
			scope.BParentDeps = true
		case strManager:
			scope.BManager = true
		case strCode:
			scope.BCode = true
		case strEnabled:
			scope.BEnabled = true
		default:
			err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid params fields")
			return scope, err
		}
	}

	return scope, nil
}

// deleteManageDepart 管理控制台删除部门
func (h *departRestHandler) deleteManageDepart(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数获取
	departID := c.Param("department_id")
	err := h.depart.DeleteDepart(&visitor, departID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// getDepartMemberInfo 获取指定部门的成员信息
func (h *departRestHandler) getDepartMemberInfo(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数获取
	departID, info, err := h.getDepMembersParmas(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	info.BShowSubDepart = true
	info.BShowSubUser = true

	// 获取指定部门的成员ID
	depInfo, depNum, userInfo, userNum, err := h.depart.GetDepartMemberInfo(&visitor, departID, info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数info
	outInfo := make(map[string]interface{})
	if info.BShowDeparts {
		departParams := make(map[string]interface{})
		entriesParm := make([]interface{}, 0)

		for k := range depInfo {
			depInfoParams := make(map[string]interface{})
			depInfoParams["id"] = depInfo[k].ID
			depInfoParams["name"] = depInfo[k].Name
			depInfoParams["type"] = strDepartment
			depInfoParams["is_root"] = depInfo[k].IsRoot
			depInfoParams["depart_existed"] = depInfo[k].BDepartExistd
			depInfoParams["user_existed"] = depInfo[k].BUserExistd
			entriesParm = append(entriesParm, depInfoParams)
		}

		departParams["total_count"] = depNum
		departParams["entries"] = entriesParm
		outInfo["departments"] = departParams
	}

	if info.BShowUsers {
		var userOutInfo ListInfo
		userOutInfo.TotalCount = userNum
		userOutInfo.Entries = make([]interface{}, 0)
		for _, v := range userInfo {
			userOutInfo.Entries = append(userOutInfo.Entries, v)
		}
		outInfo["users"] = userOutInfo
	}

	// 数据整理
	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// getDepartInfo 获取部门信息
func (h *departRestHandler) getDepartInfo(c *gin.Context) {
	// 获取参数信息
	depIDs := strings.Split(c.Param("department_id"), ",")
	fields := strings.Split(c.Param("fields"), ",")

	var scope interfaces.DepartInfoScope
	for _, v := range fields {
		switch v {
		case strName:
			scope.BShowName = true
		case strParentDeps:
			scope.BParentDeps = true
		case "managers":
			scope.BManagers = true
		case strManager:
			scope.BManager = true
		case strCode:
			scope.BCode = true
		case strEnabled:
			scope.BEnabled = true
		case strThirdID:
			scope.BThirdID = true
		default:
			err := rest.NewHTTPError("invalid params", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取部门信息
	infos, err := h.depart.GetDepartsInfo(depIDs, scope, true)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make([]interface{}, 0, len(infos))
	for k := range infos {
		tempDepInfo := h.handleDepartInfos(&infos[k], scope)
		out = append(out, tempDepInfo)
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// getDepartInfo 获取部门信息
func (h *departRestHandler) getDepartFieldInfoByPost(c *gin.Context) {
	// jsonschema校验
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, h.getDepartInfoSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	depIDs := make([]string, 0)
	for _, v := range jsonReq["department_ids"].([]interface{}) {
		depIDs = append(depIDs, v.(string))
	}

	fields := make([]string, 0)
	for _, v := range jsonReq["fields"].([]interface{}) {
		fields = append(fields, v.(string))
	}

	var scope interfaces.DepartInfoScope
	for _, v := range fields {
		switch v {
		case strName:
			scope.BShowName = true
		case strParentDeps:
			scope.BParentDeps = true
		case "managers":
			scope.BManagers = true
		case strManager:
			scope.BManager = true
		case strCode:
			scope.BCode = true
		case strEnabled:
			scope.BEnabled = true
		case strThirdID:
			scope.BThirdID = true
		default:
			err := rest.NewHTTPError("invalid params", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取部门信息
	infos, err := h.depart.GetDepartsInfo(depIDs, scope, false)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make([]interface{}, 0, len(infos))
	for k := range infos {
		tempDepInfo := h.handleDepartInfos(&infos[k], scope)
		out = append(out, tempDepInfo)
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// 处理返回值
func (h *departRestHandler) handleDepartInfos(info *interfaces.DepartInfo, scope interfaces.DepartInfoScope) (tempDepInfo map[string]interface{}) {
	tempDepInfo = make(map[string]interface{})
	tempDepInfo["department_id"] = info.ID
	if scope.BShowName {
		tempDepInfo["name"] = info.Name
	}
	if scope.BParentDeps {
		parentDeps := make([]interface{}, 0)
		for _, v1 := range info.ParentDeps {
			tempParentDep := make(map[string]interface{})
			tempParentDep["id"] = v1.ID
			tempParentDep["name"] = v1.Name
			tempParentDep["type"] = v1.Type
			parentDeps = append(parentDeps, tempParentDep)
		}
		tempDepInfo["parent_deps"] = parentDeps
	}
	if scope.BManagers {
		managers := make([]interface{}, 0)
		for _, v1 := range info.Managers {
			temp := make(map[string]interface{})
			temp["id"] = v1.ID
			temp["name"] = v1.Name
			managers = append(managers, temp)
		}
		tempDepInfo["managers"] = managers
	}

	if scope.BManager {
		if info.Manager.ID != "" {
			manager := make(map[string]interface{})
			manager["id"] = info.Manager.ID
			manager["type"] = strUser
			tempDepInfo[strManager] = manager
		}
	}

	if scope.BCode {
		tempDepInfo[strCode] = info.Code
	}

	if scope.BEnabled {
		tempDepInfo[strEnabled] = info.Enabled
	}

	if scope.BThirdID {
		tempDepInfo[strThirdID] = info.ThirdID
	}

	return
}

// getDepartInfoByLevel 根据部门层级获取部门信息
func (h *departRestHandler) getDepartInfoByLevel(c *gin.Context) {
	// 检查参数
	level, ok := c.GetQuery("level")
	nLevel, err := strconv.Atoi(level)
	if !ok || err != nil {
		err := rest.NewHTTPError("invalid level", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 获取部门信息
	infos, err := h.depart.GetDepartsInfoByLevel(nLevel)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make([]interface{}, 0, len(infos))
	for _, v := range infos {
		tempDepInfo := make(map[string]interface{})
		tempDepInfo["id"] = v.ID
		tempDepInfo["name"] = v.Name
		tempDepInfo["type"] = v.Type
		tempDepInfo["third_id"] = v.ThirdID

		out = append(out, tempDepInfo)
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// getDepMembersParmas 获取分页信息
func (h *departRestHandler) getDepMembersParmas(c *gin.Context) (depID string, out interfaces.OrgShowPageInfo, err error) {
	depID = c.Param("department_id")
	fields := strings.Split(c.Param("fields"), ",")
	for _, v := range fields {
		if v == "users" {
			out.BShowUsers = true
		} else if v == "departments" {
			out.BShowDeparts = true
		} else {
			err = rest.NewHTTPError("invalid fileds type", rest.BadRequest, nil)
			return
		}
	}

	// 参数获取和检查
	out.Offset, out.Limit, err = h.getLimitInfo(c)

	role, dRet := c.GetQuery("role")
	out.Role = roleEnumIDMap[role]
	if !dRet || out.Role == "" {
		err = rest.NewHTTPError("invalid type", rest.BadRequest,
			map[string]interface{}{"params": "role"})
		return
	}
	return
}

func (h *departRestHandler) getConsoleDepartMembersParam(c *gin.Context) (depID string, out interfaces.OrgShowPageInfo, err error) {
	depID = c.Param("department_id")
	fields := strings.Split(c.Param("fields"), ",")
	for _, v := range fields {
		if v == "departments" {
			out.BShowDeparts = true
		} else {
			err = rest.NewHTTPError("invalid fileds type", rest.BadRequest, nil)
			return
		}
	}

	// 参数获取和检查
	out.Offset, out.Limit, err = h.getLimitInfo(c)

	// 不支持普通用户/组织审计员/审计管理员
	var rOk bool
	role, dRet := c.GetQuery("role")
	out.Role, rOk = roleManageEnumIDMap[role]

	if !dRet || !rOk {
		err = rest.NewHTTPError("invalid type", rest.BadRequest,
			map[string]interface{}{"params": "role"})
		return
	}

	return
}

func (h *departRestHandler) getLimitInfo(c *gin.Context) (offset, limit int, err error) {
	// 参数获取和检查
	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		err = rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (limit < 1 || limit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		return
	}

	return offset, limit, nil
}
