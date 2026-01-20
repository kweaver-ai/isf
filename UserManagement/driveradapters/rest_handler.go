// Package driveradapters AnyShare 公共接口处理层
package driveradapters

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// RestHandler driveradapters 通用 RESTfual API Handler 接口
type RestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}
type restHandler struct {
	combine interfaces.LogicsCombine
	hydra   interfaces.Hydra
}

var (
	once    sync.Once
	handler RestHandler
)

// NewRESTHandler 通用 restful api handler 对象
func NewRESTHandler() RestHandler {
	once.Do(func() {
		handler = &restHandler{
			combine: logics.NewCombine(),
			hydra:   newHydra(),
		}
	})

	return handler
}

// RegisterPrivate 注册内部API
func (h *restHandler) RegisterPrivate(engine *gin.Engine) {
	// 注册探针
	engine.GET("/health/ready", h.getHealth)
	engine.GET("/health/alive", h.getAlive)

	engine.POST("/api/user-management/v1/names", h.convertIDToNameV1)
	engine.POST("/api/user-management/v2/names", h.convertIDToNameV2)
	engine.GET("/api/user-management/v1/emails", h.getEmails)
	engine.POST("/api/user-management/v1/search-org", h.getUserAndDepartmentInRange)
}

// RegisterPublic 注册外部API
func (h *restHandler) RegisterPublic(engine *gin.Engine) {
	// 注册探针
	engine.GET("/health/ready", h.getHealth)
	engine.GET("/health/alive", h.getAlive)

	engine.GET("/api/user-management/v1/search-in-org-tree", observable.MiddlewareTrace(common.SvcARTrace), h.searchInTree)
	engine.GET("/api/user-management/v1/console/search-in-org-tree", observable.MiddlewareTrace(common.SvcARTrace), h.consoleSearchInTree)
}

type jsonValueDesc = rest.JSONValueDesc

func (h *restHandler) getHealth(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "OK")
}

func (h *restHandler) getAlive(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.String(http.StatusOK, "OK")
}

func (h *restHandler) convertIDToNameV1(c *gin.Context) {
	h.convertIDToName(c, false)
}

func (h *restHandler) convertIDToNameV2(c *gin.Context) {
	h.convertIDToName(c, true)
}

// convertIDToName 根据用户/部门id 获取用户/部门显示名
func (h *restHandler) convertIDToName(c *gin.Context, bv2 bool) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["method"] = &jsonValueDesc{Kind: reflect.String, Required: true}

	strDesc := make(map[string]*jsonValueDesc)
	strDesc["element"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["user_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["department_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["contactor_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["group_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["app_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["strict"] = &jsonValueDesc{Kind: reflect.Bool, Required: false}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	jsonObj := jsonV.(map[string]interface{})
	strMethod := jsonObj["method"].(string)
	strUserIDs := getArrayParamFromJSON(paramDesc, jsonObj, "user_ids")
	strDeptIDs := getArrayParamFromJSON(paramDesc, jsonObj, "department_ids")
	strContactorIDs := getArrayParamFromJSON(paramDesc, jsonObj, "contactor_ids")
	strGroupIDs := getArrayParamFromJSON(paramDesc, jsonObj, "group_ids")
	strAppIDs := getArrayParamFromJSON(paramDesc, jsonObj, "app_ids")

	bStrict := true
	if date, ok := jsonObj["strict"]; ok {
		bStrict = date.(bool)
	}

	if cErr := checkMethod(strMethod); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查ID是否重复
	if bStrict {
		if cErr := checkStringUnique("user_ids", strUserIDs); cErr != nil {
			rest.ReplyError(c, cErr)
			return
		}
		if cErr := checkStringUnique("department_ids", strDeptIDs); cErr != nil {
			rest.ReplyError(c, cErr)
			return
		}
		if cErr := checkStringUnique("contactor_ids", strContactorIDs); cErr != nil {
			rest.ReplyError(c, cErr)
			return
		}
		if cErr := checkStringUnique("group_ids", strGroupIDs); cErr != nil {
			rest.ReplyError(c, cErr)
			return
		}
		if cErr := checkStringUnique("app_ids", strAppIDs); cErr != nil {
			rest.ReplyError(c, cErr)
			return
		}
	}

	// 通过id 获取用户和部门名称信息
	info := interfaces.OrgIDInfo{
		UserIDs:      strUserIDs,
		DepartIDs:    strDeptIDs,
		ContactorIDs: strContactorIDs,
		GroupIDs:     strGroupIDs,
		AppIDs:       strAppIDs,
	}

	visitor := interfaces.Visitor{
		Language:      getXLang(c),
		ErrorCodeType: GetErrorCodeType(c),
	}
	outInfo, err := h.combine.ConvertIDToName(&visitor, &info, bv2, bStrict)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	resp := map[string]interface{}{
		"user_names":       outInfo.UserNames,
		"department_names": outInfo.DepartNames,
		"contactor_names":  outInfo.ContactorNames,
		"group_names":      outInfo.GroupNames,
		"app_names":        outInfo.AppNames,
	}
	rest.ReplyOK(c, http.StatusOK, resp)
}

// getUserAndDepartmentInRange 获取范围内的用户和部门
func (h *restHandler) getUserAndDepartmentInRange(c *gin.Context) {
	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["method"] = &jsonValueDesc{Kind: reflect.String, Required: true}

	strDesc := make(map[string]*jsonValueDesc)
	strDesc["element"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["user_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["department_ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: false, ValueDesc: strDesc}
	paramDesc["scope"] = &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: strDesc}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	jsonObj := jsonV.(map[string]interface{})
	strMethod := jsonObj["method"].(string)
	var strUserIDs []string
	if paramDesc["user_ids"].Exist {
		userObj := jsonObj["user_ids"].([]interface{})
		for _, id := range userObj {
			strUserIDs = append(strUserIDs, id.(string))
		}
	}
	var strDeptIDs []string
	if paramDesc["department_ids"].Exist {
		deptObj := jsonObj["department_ids"].([]interface{})
		for _, id := range deptObj {
			strDeptIDs = append(strDeptIDs, id.(string))
		}
	}
	rangeObj := jsonObj["scope"].([]interface{})
	var strRangeDeptIDs []string
	for _, id := range rangeObj {
		strRangeDeptIDs = append(strRangeDeptIDs, id.(string))
	}
	if cErr := checkMethod(strMethod); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 通过id 获取用户和部门名称信息
	userIDs, deptIDs, err := h.combine.GetUserAndDepartmentInScope(strUserIDs, strDeptIDs, strRangeDeptIDs)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	resp := map[string]interface{}{
		"user_ids":       userIDs,
		"department_ids": deptIDs,
	}
	rest.ReplyOK(c, http.StatusOK, resp)
}

// getEmails 根据用户/部门id 获取用户/部门邮箱地址
func (h *restHandler) getEmails(c *gin.Context) {
	// 获取请求参数
	strUserIDs, userRet := c.GetQueryArray("user_id")
	strDeptIDs, depRet := c.GetQueryArray("department_id")
	if !userRet && !depRet {
		err := rest.NewHTTPError("invalid param user_id and department_id", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 通过id 获取用户和部门邮箱信息
	info := interfaces.OrgIDInfo{
		UserIDs:   strUserIDs,
		DepartIDs: strDeptIDs,
	}

	visitor := interfaces.Visitor{
		Language:      getXLang(c),
		ErrorCodeType: GetErrorCodeType(c),
	}
	outInfo, err := h.combine.GetEmails(&visitor, &info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	resp := map[string]interface{}{
		"user_emails":       outInfo.UserEmails,
		"department_emails": outInfo.DepartEmails,
	}
	rest.ReplyOK(c, http.StatusOK, resp)
}

// consoleSearchInTree 管理控制台组织架构树展开搜索用户和部门信息
func (h *restHandler) consoleSearchInTree(c *gin.Context) {
	// token验证
	visitor, err := verify(c, h.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// query参数获取
	info, err := h.getSearchInTreeParams(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	strOnlyShowEnabledUser, ok := c.GetQuery("user_enabled")
	if ok {
		info.BOnlyShowEnabledUser, err = strconv.ParseBool(strOnlyShowEnabledUser)
		if err != nil || !info.BOnlyShowEnabledUser {
			err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid user_enabled")
			rest.ReplyError(c, err)
			return
		}
	}

	strOnlyShowAssignedUser, ok := c.GetQuery("user_assigned")
	if ok {
		info.BOnlyShowAssignedUser, err = strconv.ParseBool(strOnlyShowAssignedUser)
		if err != nil || !info.BOnlyShowAssignedUser {
			err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid user_assigned")
			rest.ReplyError(c, err)
			return
		}
	}

	// 搜索info
	var userInfos []interfaces.SearchUserInfo
	var departInfos []interfaces.SearchDepartInfo
	var userNum, departNum int
	userInfos, departInfos, userNum, departNum, err = h.combine.SearchInOrgTree(c, &visitor, info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make(map[string]interface{})
	if info.BShowUsers {
		var outUsers ListInfo
		outUsers.Entries = make([]interface{}, 0)
		for _, v := range userInfos {
			outUsers.Entries = append(outUsers.Entries, v)
		}
		outUsers.TotalCount = userNum
		out["users"] = outUsers
	}

	if info.BShowDeparts {
		var outDeparts ListInfo
		outDeparts.Entries = make([]interface{}, 0)
		for _, v := range departInfos {
			outDeparts.Entries = append(outDeparts.Entries, v)
		}
		outDeparts.TotalCount = departNum
		out["departments"] = outDeparts
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// searchInTree 组织架构树展开搜索用户和部门信息
func (h *restHandler) searchInTree(c *gin.Context) {
	// token验证
	visitor, err := verify(c, h.hydra)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	info, err := h.getSearchInTreeParams(c)
	info.BOnlyShowEnabledUser = true
	info.BOnlyShowAssignedUser = true
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 搜索info
	var userInfos []interfaces.SearchUserInfo
	var departInfos []interfaces.SearchDepartInfo
	var userNum, departNum int
	userInfos, departInfos, userNum, departNum, err = h.combine.SearchInOrgTree(c, &visitor, info)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 数据整理
	out := make(map[string]interface{})
	if info.BShowUsers {
		var outUsers ListInfo
		outUsers.Entries = make([]interface{}, 0)
		for _, v := range userInfos {
			outUsers.Entries = append(outUsers.Entries, v)
		}
		outUsers.TotalCount = userNum
		out["users"] = outUsers
	}

	if info.BShowDeparts {
		var outDeparts ListInfo
		outDeparts.Entries = make([]interface{}, 0)
		for _, v := range departInfos {
			outDeparts.Entries = append(outDeparts.Entries, v)
		}
		outDeparts.TotalCount = departNum
		out["departments"] = outDeparts
	}

	rest.ReplyOK(c, http.StatusOK, out)
}

// getArrayParamFromJSON 从json中获取string数组
func getArrayParamFromJSON(json map[string]*jsonValueDesc, jsonObj map[string]interface{}, key string) (info []string) {
	if json[key].Exist {
		obj := jsonObj[key].([]interface{})
		for _, id := range obj {
			info = append(info, id.(string))
		}
	}
	return
}

// checkMethod 检查方法是否正确
func checkMethod(method string) error {
	const strMethod string = "GET"
	if method != strMethod {
		return rest.NewHTTPError("invalid method", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "method"}})
	}
	return nil
}

// checkStringUnique 检查数组是否重复
func checkStringUnique(param string, arry []string) error {
	sort.Strings(arry)

	bIsUnique := true
	var strLast string
	for _, v := range arry {
		if v == strLast {
			bIsUnique = false
			break
		}
		strLast = v
	}

	if !bIsUnique {
		return rest.NewHTTPError(fmt.Sprintf("param %s is not unique", param), rest.BadRequest, nil)
	}
	return nil
}

// getSearchInTreeParams 搜索参数获取
func (h *restHandler) getSearchInTreeParams(c *gin.Context) (info interfaces.OrgShowPageInfo, err error) {
	// 检查请求参数与文档是否匹配
	strTypes, ret := c.GetQueryArray("type")
	if !ret {
		err = rest.NewHTTPError("invalid type type", rest.BadRequest, nil)
		return
	}

	for _, v := range strTypes {
		switch v {
		case strUser:
			info.BShowUsers = true
		case strDepartment:
			info.BShowDeparts = true
		default:
			err = rest.NewHTTPError("invalid type type", rest.BadRequest, nil)
			return
		}
	}

	info.Offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || info.Offset < 0 {
		err = rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
		return
	}

	info.Limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (info.Limit < 1 || info.Limit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		return
	}

	info.Keyword, ret = c.GetQuery("keyword")
	if !ret || info.Keyword == "" {
		err = rest.NewHTTPError("invalid keyword", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "keyword"}})
		return
	}

	role, ret := c.GetQuery("role")
	info.Role = roleEnumIDMap[role]
	if !ret || info.Role == "" {
		err = rest.NewHTTPError("invalid type", rest.BadRequest,
			map[string]interface{}{"params": "role"})
		return
	}
	return
}
