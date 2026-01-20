// Package driveradapters group AnyShare  用户组逻辑接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/language"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"

	gerrors "github.com/kweaver-ai/go-lib/error"
)

// ListInfo 列举结构
type ListInfo struct {
	Entries    []interface{} `json:"entries" binding:"required"`
	TotalCount int           `json:"total_count" binding:"required"`
}

// GroupMemberInfo 用户组成员输出信息结构
type GroupMemberInfo struct {
	ID              string   `json:"id" binding:"required"`
	MemberType      string   `json:"type" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	DepartmentNames []string `json:"department_names" binding:"required"`
}

// MemberInfo 用户组成员输出信息结构（不包含直属部门信息）
type MemberInfo struct {
	ID         string `json:"id" binding:"required"`
	MemberType string `json:"type" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

// GMInfo 用户组成员输出信息结构（包含所属组信息）
type GMInfo struct {
	ID         string   `json:"id" binding:"required"`
	MemberType string   `json:"type" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	GroupNames []string `json:"group_names" binding:"required"`
}

// GroupRestHandler driveradapters 用户组 RESTfual API Handler 接口
type GroupRestHandler interface {
	// RegisterPublic 注册公共API
	RegisterPublic(engine *gin.Engine)

	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)
}
type groupRestHandler struct {
	hydra                 interfaces.Hydra
	group                 interfaces.LogicsGroup
	user                  interfaces.LogicsUser
	combine               interfaces.LogicsCombine
	memberIntTypes        map[int]string
	memberStringTypes     map[string]int
	getGroupMembersSchema *gojsonschema.Schema
	createGroupSchema     *gojsonschema.Schema
}

type listQueryInfo struct {
	direction  interfaces.Direction
	sort       interfaces.SortFiled
	keyword    string
	offset     int
	limit      int
	hasKeyword bool
}

const (
	strName        = "name"
	strDateCreated = "date_created"
	strDesc        = "desc"
	strAsc         = "asc"
	strHTTPGET     = "GET"
	strUser        = "user"
)

var (
	gOnce    sync.Once
	ghandler GroupRestHandler
	sortMap  = map[string]interfaces.SortFiled{
		strDateCreated: interfaces.DateCreated,
		strName:        interfaces.Name,
	}
	directionMap = map[string]interfaces.Direction{
		strDesc: interfaces.Desc,
		strAsc:  interfaces.Asc,
	}
	langMatcher = language.NewMatcher([]language.Tag{
		language.SimplifiedChinese,
		language.TraditionalChinese,
		language.AmericanEnglish,
	})
	langMap = map[language.Tag]interfaces.Language{
		language.SimplifiedChinese:  interfaces.SimplifiedChinese,
		language.TraditionalChinese: interfaces.TraditionalChinese,
		language.AmericanEnglish:    interfaces.AmericanEnglish,
	}

	errCodeTypeMap = map[string]interfaces.ErrorCodeType{
		"string": interfaces.Str,
		"number": interfaces.Number,
	}
)

var (
	//go:embed jsonschema/group/get_group_members.json
	getGroupMembersSchemaStr string

	//go:embed jsonschema/group/create_group.json
	createGroupSchemaStr string
)

// NewGroupRESTHandler 用户组 restful api handler 对象
func NewGroupRESTHandler() GroupRestHandler {
	gOnce.Do(func() {
		getGroupMembersSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getGroupMembersSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		createGroupSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(createGroupSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		ghandler = &groupRestHandler{
			hydra:   newHydra(),
			group:   logics.NewGroup(),
			user:    logics.NewUser(),
			combine: logics.NewCombine(),
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
	})

	return ghandler
}

// RegisterPrivate 注册内部API
func (g *groupRestHandler) RegisterPrivate(engine *gin.Engine) {
	// 组成员管理
	engine.POST("/api/user-management/v1/group-members", g.getMembersID)
}

// RegisterPublic 注册外部API
func (g *groupRestHandler) RegisterPublic(engine *gin.Engine) {
	// 组管理

	engine.POST("/api/user-management/v1/management/groups", observable.MiddlewareTrace(common.SvcARTrace), g.createGroup)
	engine.DELETE("/api/user-management/v1/management/groups/:id", g.deleteGroup)
	engine.PUT("/api/user-management/v1/management/groups/:id/:fields", g.modifyGroup)
	engine.GET("/api/user-management/v1/management/groups", g.getGroup)
	engine.GET("/api/user-management/v1/management/groups/:id", g.getGroupByID)
	engine.GET("/api/user-management/v1/search-in-group", g.searchGroupByKey)
	engine.GET("/api/user-management/v1/groups", g.getGroupOnClient)

	// 组成员管理
	engine.GET("/api/user-management/v1/management/group-members/:id", observable.MiddlewareTrace(common.SvcARTrace), g.getGroupMemebers)
	engine.POST("/api/user-management/v1/management/group-members/:id", g.addOrDeleteGroupMembers)
	engine.GET("/api/user-management/v1/group-members/:id", g.getMembersOnClient)
	engine.GET("/api/user-management/v1/console/group-members/:id/user-match", observable.MiddlewareTrace(common.SvcARTrace), g.userMatch)
	engine.GET("/api/user-management/v1/console/search-users-in-group", observable.MiddlewareTrace(common.SvcARTrace), g.searchInAllGroupOrg)
}

// getMembersID 根据用户组id数组批量获取用户成员和部门成员id
func (g *groupRestHandler) getMembersID(c *gin.Context) {
	// 获取请求参数
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, g.getGroupMembersSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取具体请求参数
	method := jsonReq["method"].(string)
	groupIDsObj := jsonReq["group_ids"].([]interface{})
	var groupIDs []string
	for _, id := range groupIDsObj {
		groupIDs = append(groupIDs, id.(string))
	}

	bOnlyShowEnableUser := false
	if tempUserEnabled, ok := jsonReq["user_enabled"]; ok {
		bOnlyShowEnableUser = tempUserEnabled.(bool)
		if !bOnlyShowEnableUser {
			cErr := rest.NewHTTPError("invalid param 'user_enabled'", rest.BadRequest, nil)
			rest.ReplyError(c, cErr)
			return
		}
	}

	// 检测参数是否合法
	if method != strHTTPGET {
		cErr := rest.NewHTTPError("invalid param 'method'", rest.BadRequest, nil)
		rest.ReplyError(c, cErr)
		return
	}

	if len(groupIDs) == 0 {
		cErr := rest.NewHTTPError("param 'group_ids' is empty", rest.BadRequest, nil)
		rest.ReplyError(c, cErr)
		return
	}

	// 通过id 获取用户和部门名称信息
	visitor := interfaces.Visitor{
		Language:      getXLang(c),
		ErrorCodeType: GetErrorCodeType(c),
	}
	userIds, departmentIds, err := g.group.GetGroupMembersID(&visitor, groupIDs, !bOnlyShowEnableUser)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应200
	resp := map[string]interface{}{
		"user_ids":       userIds,
		"department_ids": departmentIds,
	}

	rest.ReplyOK(c, http.StatusOK, resp)
}

// createGroup 创建组
func (g *groupRestHandler) createGroup(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, g.createGroupSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取具体请求参数
	name := jsonReq[strName].(string)
	notes := ""
	if strNotes, ok := jsonReq["notes"]; ok {
		notes = strNotes.(string)
	}

	groupIds := []string{}
	if strGroupIDs, ok := jsonReq["group_ids_of_members"]; ok {
		for _, v := range strGroupIDs.([]interface{}) {
			groupIds = append(groupIds, v.(string))
		}
	}

	strID, err := g.group.AddGroup(c, &visitor, name, notes, groupIds)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/user-management/v1/management/groups/%s", strID))
	rest.ReplyOK(c, http.StatusCreated, gin.H{"id": strID})
}

// deleteGroup 删除组
func (g *groupRestHandler) deleteGroup(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取删除用户组ID
	groupID := c.Param("id")

	// 删除用户组
	err := g.group.DeleteGroup(&visitor, groupID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// modifyGroup 修改组
func (g *groupRestHandler) modifyGroup(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取修改用户组 存在哪些参数
	groupID := c.Param("id")
	fields := strings.Split(c.Param("fields"), ",")
	var nameExist, notesExist bool
	for _, v := range fields {
		if v == strName {
			nameExist = true
		} else if v == "notes" {
			notesExist = true
		}
	}

	// 获取请求参数bool
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	if nameExist {
		paramDesc[strName] = &jsonValueDesc{Kind: reflect.String, Required: true}
	}

	if notesExist {
		paramDesc["notes"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	jsonObj := jsonV.(map[string]interface{})
	var name, notes string
	if nameExist {
		name = jsonObj[strName].(string)
		err := checkName(name)
		if err != nil {
			rest.ReplyError(c, err)
			return
		}
	}
	if notesExist {
		maxLen := 300
		notes = jsonObj["notes"].(string)
		if utf8.RuneCountInString(notes) > maxLen {
			err := rest.NewHTTPError("param notes is illegal", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 修改用户组
	err := g.group.ModifyGroup(&visitor, groupID, name, nameExist, notes, notesExist)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// getGroup 列举组
func (g *groupRestHandler) getGroup(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取列举用户组信息
	queryInfo, err := getListQueryParam(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	searchInfo := interfaces.SearchInfo{
		Direction:  queryInfo.direction,
		Keyword:    queryInfo.keyword,
		Sort:       queryInfo.sort,
		Offset:     queryInfo.offset,
		Limit:      queryInfo.limit,
		HasKeyWord: queryInfo.hasKeyword,
	}

	// 列举用户组
	num, groupInfos, err := g.group.GetGroup(&visitor, searchInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	var outInfo ListInfo
	outInfo.TotalCount = num
	outInfo.Entries = make([]interface{}, 0)
	for _, v := range groupInfos {
		data := make(map[string]interface{})
		data["id"] = v.ID
		data["name"] = v.Name
		data["notes"] = v.Notes
		outInfo.Entries = append(outInfo.Entries, data)
	}

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// getGroupByID 根据组ID获取组信息
func (g *groupRestHandler) getGroupByID(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取用户组ID
	groupID := c.Param("id")

	// 获取用户组信息
	groupInfo, err := g.group.GetGroupByID(&visitor, groupID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	out := make(map[string]interface{})
	out["id"] = groupInfo.ID
	out["name"] = groupInfo.Name
	out["notes"] = groupInfo.Notes
	rest.ReplyOK(c, http.StatusOK, out)
}

// searchInAllGroupOrg 在组成员完整组织架构下根据显示名模糊搜索用户
func (g *groupRestHandler) searchInAllGroupOrg(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取groupid
	groupID, bExist := c.GetQuery("group_id")
	if !bExist {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid group_id")
		rest.ReplyError(c, err)
		return
	}

	// 获取搜索关键字
	userName, bExist := c.GetQuery("key")
	if !bExist {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid key")
		rest.ReplyError(c, err)
		return
	}

	// 参数获取和检查
	nStart, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || nStart < 0 {
		err = rest.NewHTTPError("invalid offset type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	var nLimit int
	nLimit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (nLimit < 1 || nLimit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	//  用户匹配
	num, userIDs, uInfos, mInfos, err := g.group.SearchInAllGroupOrg(c, &visitor, groupID, userName, nStart, nLimit)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	outInfo := make(map[string]interface{})
	outInfo["total_count"] = num
	entriesData := make([]interface{}, 0)
	for _, v := range userIDs {
		tempInfo := make(map[string]interface{})
		tempInfo["id"] = uInfos[v].ID
		tempInfo["type"] = strUser
		tempInfo["name"] = uInfos[v].Name
		tempInfo["parent_deps"] = g.handleNameInfos(uInfos[v].ParentDeps)

		tempData := make([]interface{}, 0, len(mInfos[uInfos[v].ID]))
		for i := range mInfos[uInfos[v].ID] {
			tempData = append(tempData, g.handleGroupMemberInfo2(&mInfos[uInfos[v].ID][i]))
		}
		tempInfo["group_members"] = tempData

		entriesData = append(entriesData, tempInfo)
	}
	outInfo["entries"] = entriesData

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// userMatch 组成员下成员匹配
func (g *groupRestHandler) userMatch(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	groupID := c.Param("id")
	// 获取搜索关键字
	userName, bExist := c.GetQuery("name")
	if !bExist {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid name")
		rest.ReplyError(c, err)
		return
	}

	//  用户匹配
	bExist, uInfo, mInfos, err := g.group.UserMatch(c, &visitor, groupID, userName)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	outInfo := make(map[string]interface{})
	outInfo["result"] = bExist
	if bExist {
		outInfo["id"] = uInfo.ID
		outInfo["parent_deps"] = g.handleNameInfos(uInfo.ParentDeps)

		entriesData := make([]interface{}, 0, len(mInfos))
		for i := range mInfos {
			entriesData = append(entriesData, g.handleGroupMemberInfo2(&mInfos[i]))
		}
		outInfo["group_members"] = entriesData
	}

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// getGroupMemebers 列举组成员
func (g *groupRestHandler) getGroupMemebers(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	groupID := c.Param("id")
	queryInfo, err := getListQueryParam(c)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	searchInfo := interfaces.SearchInfo{
		Direction:  queryInfo.direction,
		Keyword:    queryInfo.keyword,
		Sort:       queryInfo.sort,
		Offset:     queryInfo.offset,
		Limit:      queryInfo.limit,
		HasKeyWord: queryInfo.hasKeyword,
	}

	// 列举组成员
	num, infos, err := g.group.GetGroupMembers(c, &visitor, groupID, searchInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	outInfo := make(map[string]interface{})
	outInfo["total_count"] = num

	entriesData := make([]interface{}, 0, len(infos))
	for i := range infos {
		entriesData = append(entriesData, g.handleGroupMemberInfo(&infos[i]))
	}
	outInfo["entries"] = entriesData

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// addGroupMemebers 添加和删除组成员
func (g *groupRestHandler) addOrDeleteGroupMembers(c *gin.Context) {
	// token验证和userID获取
	visitor, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	groupID := c.Param("id")

	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["method"] = &jsonValueDesc{Kind: reflect.String, Required: true}

	memberDesc := make(map[string]*jsonValueDesc)
	memberDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	memberDesc["type"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	membersDesc := make(map[string]*jsonValueDesc)
	membersDesc["element"] = &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: memberDesc}
	paramDesc["members"] = &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: membersDesc}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 检查参数是否合法
	jsonObj := jsonV.(map[string]interface{})
	method := jsonObj["method"].(string)
	members := jsonObj["members"].([]interface{})
	memberInfos := make(map[string]interfaces.GroupMemberInfo)
	for _, v := range members {
		member := interfaces.GroupMemberInfo{}
		temp := v.(map[string]interface{})
		member.ID = temp["id"].(string)
		strType := temp["type"].(string)
		var ok bool
		member.MemberType, ok = g.memberStringTypes[strType]
		if !ok {
			err := rest.NewHTTPError("param member type is illegal", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
		if member.ID == "" {
			err := rest.NewHTTPError("param member id is illegal", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}

		memberInfos[member.ID] = member
	}
	if method != "POST" && method != "DELETE" {
		err := rest.NewHTTPError("param method is illegal", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 批量添加或者删除用户组成员
	cErr := g.group.AddOrDeleteGroupMemebers(&visitor, method, groupID, memberInfos)
	if cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// 在组织下搜索组信息
func (g *groupRestHandler) searchGroupByKey(c *gin.Context) {
	// token验证
	_, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数检查
	var searchInfo interfaces.SearchClientInfo
	err := g.checkGMSearchParams(c, &searchInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 列举用户组
	tempInfo, err := g.combine.SearchGroupAndMemberInfoByKey(&searchInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 参数整理
	var groupInfo ListInfo
	groupInfo.Entries = make([]interface{}, 0)
	var memberInfo ListInfo
	memberInfo.Entries = make([]interface{}, 0)
	outData := make(map[string]interface{})
	var tmp GMInfo
	if searchInfo.BShowMember {
		for _, v := range tempInfo.MemberInfos {
			tmp.ID = v.ID
			tmp.Name = v.Name
			tmp.MemberType = g.memberIntTypes[v.NType]
			tmp.GroupNames = v.GroupNames
			memberInfo.Entries = append(memberInfo.Entries, tmp)
		}
		memberInfo.TotalCount = tempInfo.MemberNum

		outData["members"] = memberInfo
	}
	if searchInfo.BShowGroup {
		for _, v := range tempInfo.GroupInfos {
			groupInfo.Entries = append(groupInfo.Entries, v)
		}
		groupInfo.TotalCount = tempInfo.GroupNum

		outData["groups"] = groupInfo
	}

	// 设置返回参数
	rest.ReplyOK(c, http.StatusOK, outData)
}

// getGroupOnClient 客户端组列举接口
func (g *groupRestHandler) getGroupOnClient(c *gin.Context) {
	// token验证
	_, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数获取和检查
	nStart, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || nStart < 0 {
		err = rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	var nLimit int
	nLimit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (nLimit < 1 || nLimit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 列举组
	info, num, err := g.group.GetGroupOnClient(nStart, nLimit)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	outData := map[string]interface{}{
		"entries":     info,
		"total_count": num,
	}
	// 设置返回参数
	rest.ReplyOK(c, http.StatusOK, outData)
}

// getMembersOnClient
func (g *groupRestHandler) getMembersOnClient(c *gin.Context) {
	// token验证和userID获取
	_, vErr := verify(c, g.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	groupID := c.Param("id")

	// 参数获取和检查
	nStart, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || nStart < 0 {
		err = rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	var nLimit int
	nLimit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (nLimit < 1 || nLimit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 列举组成员
	info, num, err := g.group.GetMemberOnClient(groupID, nStart, nLimit)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 设置返回参数
	var outInfo ListInfo
	outInfo.TotalCount = num
	outInfo.Entries = make([]interface{}, 0)

	var tmp MemberInfo
	for _, v := range info {
		tmp.ID = v.ID
		tmp.Name = v.Name
		tmp.MemberType = g.memberIntTypes[v.NType]

		outInfo.Entries = append(outInfo.Entries, tmp)
	}

	// 设置返回参数
	rest.ReplyOK(c, http.StatusOK, outInfo)
}

func (g *groupRestHandler) checkGMSearchParams(c *gin.Context, searchInfo *interfaces.SearchClientInfo) (err error) {
	keyword, dRet := c.GetQuery("keyword")
	if !dRet || keyword == "" {
		err = rest.NewHTTPError("invalid keyword type", rest.BadRequest, nil)
		return
	}
	searchInfo.Keyword = keyword

	nStart, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || nStart < 0 {
		err = rest.NewHTTPError("invalid start type", rest.BadRequest, nil)
		return
	}
	searchInfo.Offset = nStart

	var nLimit int
	nLimit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (nLimit < 1 || nLimit > 1000) {
		err = rest.NewHTTPError("invalid limit type", rest.BadRequest, nil)
		return
	}
	searchInfo.Limit = nLimit

	var strTypes []string
	strTypes, dRet = c.GetQueryArray("type")
	if !dRet {
		err = rest.NewHTTPError("invalid type", rest.BadRequest, nil)
		return
	}
	for _, v := range strTypes {
		switch v {
		case "member":
			searchInfo.BShowMember = true
		case "group":
			searchInfo.BShowGroup = true
		default:
			err = rest.NewHTTPError("invalid type", rest.BadRequest, nil)
			return
		}
	}
	return nil
}

// handleGroupMemberInfos 处理用户组成员信息
func (g *groupRestHandler) handleGroupMemberInfo(input *interfaces.GroupMemberInfo) interface{} {
	tempMemberInfo := make(map[string]interface{})
	tempMemberInfo["id"] = input.ID
	tempMemberInfo["type"] = g.memberIntTypes[input.MemberType]
	tempMemberInfo["name"] = input.Name
	tempMemberInfo["department_names"] = input.DepartmentNames
	tempMemberInfo["parent_deps"] = g.handleNameInfos(input.ParentDeps)
	return tempMemberInfo
}

// handleGroupMemberInfos2 处理用户组成员信息,只处理部分信息
func (g *groupRestHandler) handleGroupMemberInfo2(input *interfaces.GroupMemberInfo) interface{} {
	tempMemberInfo := make(map[string]interface{})
	tempMemberInfo["id"] = input.ID
	tempMemberInfo["type"] = g.memberIntTypes[input.MemberType]
	tempMemberInfo["name"] = input.Name
	tempMemberInfo["parent_deps"] = g.handleNameInfos(input.ParentDeps)
	return tempMemberInfo
}

// handleNameInfos 处理名称信息
func (g *groupRestHandler) handleNameInfos(input [][]interfaces.NameInfo) []interface{} {
	tempParentDeps := make([]interface{}, 0, len(input))
	for _, v1 := range input {
		if len(v1) == 0 {
			continue
		}
		temp1 := make([]interface{}, 0, len(v1))
		for _, v2 := range v1 {
			temp2 := make(map[string]interface{})
			temp2["id"] = v2.ID
			temp2["name"] = v2.Name
			temp2["type"] = strDepartment

			temp1 = append(temp1, temp2)
		}
		tempParentDeps = append(tempParentDeps, temp1)
	}
	return tempParentDeps
}

func checkName(name string) (err error) {
	ret := !strings.ContainsAny(name, " ") &&
		!strings.ContainsAny(name, "|") &&
		!strings.ContainsAny(name, "\\") &&
		!strings.ContainsAny(name, "/") &&
		!strings.ContainsAny(name, ":") &&
		!strings.ContainsAny(name, "*") &&
		!strings.ContainsAny(name, "?") &&
		!strings.ContainsAny(name, "\"") &&
		!strings.ContainsAny(name, ">") &&
		!strings.ContainsAny(name, "<") &&
		utf8.RuneCountInString(name) <= 128 &&
		utf8.RuneCountInString(name) >= 1
	if !ret {
		return rest.NewHTTPError("param name is illegal", rest.BadRequest, nil)
	}
	return nil
}

func getListQueryParam(c *gin.Context) (info listQueryInfo, err error) {
	// 获取搜索关键字
	info.keyword, info.hasKeyword = c.GetQuery(("keyword"))

	// 获取排序方向
	directionStr, dRet := c.GetQuery(("direction"))
	if dRet && directionStr != strAsc && directionStr != strDesc {
		err = rest.NewHTTPError("invalid direction type", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "direction"}})
		return
	}

	if !dRet {
		directionStr = strDesc
	}
	info.direction = directionMap[directionStr]

	// 获取排序类型
	sortStr, sRet := c.GetQuery(("sort"))
	if sRet && sortStr != strDateCreated && sortStr != strName {
		err = rest.NewHTTPError("invalid sort type", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "sort"}})
		return
	}

	if !sRet {
		sortStr = strDateCreated
	}
	info.sort = sortMap[sortStr]

	// 获取数据下标
	info.offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		err = rest.NewHTTPError("offset is illeagal", rest.BadRequest, nil)
		return
	}
	if info.offset < 0 {
		err = rest.NewHTTPError("invalid offset(>=0)", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "offset"}})
		return
	}

	// 获取分页数据
	info.limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		err = rest.NewHTTPError("limit is illeagal", rest.BadRequest, nil)
		return
	}
	if info.limit < 1 || info.limit > 1000 {
		err = rest.NewHTTPError("invalid limit([ 1 .. 1000 ])", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "limit"}})
		return
	}

	return
}

func GetErrorCodeType(c *gin.Context) interfaces.ErrorCodeType {
	errCodeType := strings.ToLower(c.GetHeader("x-error-code"))
	if val, isExist := errCodeTypeMap[errCodeType]; isExist {
		return val
	}
	return interfaces.Number
}

// verify token有效性检查
func verify(c *gin.Context, hydra interfaces.Hydra) (visitor interfaces.Visitor, err error) {
	tokenID := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenID, "Bearer ")
	info, err := hydra.Introspect(token)
	if err != nil {
		return
	}

	if !info.Active {
		err = rest.NewHTTPError("token expired", rest.Unauthorized, nil)
		return
	}

	visitor = interfaces.Visitor{
		ID:            info.VisitorID,
		TokenID:       tokenID,
		IP:            c.ClientIP(),
		Mac:           c.GetHeader("X-Request-MAC"),
		UserAgent:     c.GetHeader("User-Agent"),
		Type:          info.VisitorTyp,
		Language:      getXLang(c),
		ErrorCodeType: GetErrorCodeType(c),
	}

	return
}

// verify token有效性检查
func verifyNewError(c *gin.Context, hydra interfaces.Hydra) (visitor interfaces.Visitor, err error) {
	tokenID := c.GetHeader("Authorization")
	token := strings.TrimPrefix(tokenID, "Bearer ")
	info, err := hydra.Introspect(token)
	if err != nil {
		return
	}

	if !info.Active {
		err = gerrors.NewError(gerrors.PublicUnauthorized, "token expired")
		return
	}

	visitor = interfaces.Visitor{
		ID:            info.VisitorID,
		TokenID:       tokenID,
		IP:            c.ClientIP(),
		Mac:           c.GetHeader("X-Request-MAC"),
		UserAgent:     c.GetHeader("User-Agent"),
		Type:          info.VisitorTyp,
		Language:      getXLang(c),
		ErrorCodeType: GetErrorCodeType(c),
	}

	return
}

// getXLang 解析获取 Header x-language
func getXLang(c *gin.Context) interfaces.Language {
	tag, _ := language.MatchStrings(langMatcher, getBCP47(c.GetHeader("x-language")))
	return langMap[tag]
}

// getBCP47 将约定的语言标签转换为符合BCP47标准的语言标签
// 默认值为 zh-Hans, 中国大陆简体中文
// https://www.rfc-editor.org/info/bcp47
func getBCP47(s string) string {
	switch strings.ToLower(s) {
	case "zh_cn", "zh-cn":
		return "zh-Hans"
	case "zh_tw", "zh-tw":
		return "zh-Hant"
	case "en_us", "en-us":
		return "en-US"
	default:
		return "zh-Hans"
	}
}
