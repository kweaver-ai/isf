// Package driveradapters user AnyShare  用户逻辑接口处理层
package driveradapters

import (
	_ "embed" // 标准用法
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	gerrors "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// UserRestHandler user RESTfual API Handler 接口
type UserRestHandler interface {
	// RegisterPrivate 注册内部API
	RegisterPrivate(engine *gin.Engine)

	// RegisterPublic 注册外部API
	RegisterPublic(engine *gin.Engine)
}
type userRestHandler struct {
	hydra                   interfaces.Hydra
	user                    interfaces.LogicsUser
	role                    interfaces.LogicsRole
	avatar                  interfaces.LogicsAvatar
	roleNameMap             map[interfaces.Role]string
	authTypeMap             map[interfaces.AuthType]string
	pwdRetrievalMap         map[interfaces.PwdRetrievalStatus]string
	pwdErrInfoSchema        *gojsonschema.Schema
	incrementUserInfoSchema *gojsonschema.Schema
	getUserInfoSchema       *gojsonschema.Schema
	markerSchema            *gojsonschema.Schema
}

type userListQueryInfo struct {
	direction    interfaces.Direction
	limit        int
	hasMarker    bool
	createdStamp int64
	userID       string
}

// 角色Enum信息
const (
	EnumSuperAdmin string = "super_admin"
	EnumSysAdmin   string = "sys_admin"
	EnumAuditAdmin string = "audit_admin"
	EnumSecAdmin   string = "sec_admin"
	EnumOrgManager string = "org_manager"
	EnumOrgAudit   string = "org_audit"
	EnumNormaluser string = "normal_user"
)

// 认证类型Enum信息
const (
	EnumLocal  string = "local"
	EnumDomain string = "domain"
	EnumThird  string = "third"
)

// 密码找回状态Enum信息
const (
	EnumAvaliable          string = "available"
	EnumInvalidUser        string = "invalid_account"
	EnumUnablePWDRetrieval string = "unable_pwd_retrieval"
	EnumDisableUser        string = "disable_user"
	EnumNonLocalUser       string = "non_local_user"
	EnumEnablePWDControl   string = "enable_pwd_control"
)

var roleEnumIDMap = map[string]interfaces.Role{
	EnumSuperAdmin: interfaces.SystemRoleSuperAdmin,
	EnumSysAdmin:   interfaces.SystemRoleSysAdmin,
	EnumAuditAdmin: interfaces.SystemRoleAuditAdmin,
	EnumSecAdmin:   interfaces.SystemRoleSecAdmin,
	EnumOrgManager: interfaces.SystemRoleOrgManager,
	EnumOrgAudit:   interfaces.SystemRoleOrgAudit,
	EnumNormaluser: interfaces.SystemRoleNormalUser,
}

// LDAPServerType类型 枚举
var ladpServerTypeEnumMap = map[interfaces.LDAPServerType]string{
	interfaces.WindowAD:  "windows_ad",
	interfaces.OtherLDAP: "other_ldap",
}

// 用户账号密码验证失败Enum详细信息
var authFailedReasonEnumMap = map[interfaces.AuthFailedReason]string{
	interfaces.InvalidPassword:            "invalid_password",
	interfaces.InitialPassword:            "initial_password",
	interfaces.PasswordNotSafe:            "password_not_safe",
	interfaces.UnderControlPasswordExpire: "under_control_password_expire",
	interfaces.PasswordExpire:             "password_expire",
}
var (
	//go:embed jsonschema/user_auth/update_pwd_err_info.json
	pwdErrInfoSchemaStr string

	//go:embed jsonschema/user_management/increment_update_user_info.json
	incrementUserInfoSchemaStr string

	//go:embed jsonschema/user_management/get_user_info.json
	getUserInfoSchemaStr string

	//go:embed jsonschema/user_management/marker.json
	markerSchemaStr string
)

var (
	uOnce    sync.Once
	uHandler UserRestHandler

	strAccount   = "account"
	strThirdAttr = "third_attr"
	strName1     = "name"
	strEmail     = "email"
	strThirdID   = "third_id"

	directionStringMap = map[interfaces.Direction]string{
		interfaces.Desc: strDesc,
		interfaces.Asc:  strAsc,
	}
)

// NewUserRESTHandler user handler 对象
func NewUserRESTHandler() UserRestHandler {
	uOnce.Do(func() {
		pwdErrInfoSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(pwdErrInfoSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		incrementUserInfoSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(incrementUserInfoSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		getUserInfoSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(getUserInfoSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		markerSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(markerSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

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

		uHandler = &userRestHandler{
			user:                    logics.NewUser(),
			roleNameMap:             roleNameMap,
			hydra:                   newHydra(),
			role:                    logics.NewRole(),
			avatar:                  logics.NewAvatar(),
			authTypeMap:             authTypeMap,
			pwdRetrievalMap:         PwdRetrievalMap,
			pwdErrInfoSchema:        pwdErrInfoSchema,
			incrementUserInfoSchema: incrementUserInfoSchema,
			getUserInfoSchema:       getUserInfoSchema,
			markerSchema:            markerSchema,
		}
	})

	return uHandler
}

// RegisterPrivate 注册内部API
func (h *userRestHandler) RegisterPrivate(engine *gin.Engine) {
	engine.GET("/api/user-management/v1/users/:user_id/:fields", h.getUserInfo)
	engine.POST("/api/user-management/v1/batch-get-user-info", observable.MiddlewareTrace(common.SvcARTrace), h.getUserBaseInfoByPost)
	engine.GET("/api/user-management/v1/pwd-retrieval-method", h.getPWDRetrievalMethod)
	engine.GET("/api/user-management/v1/account-match", h.getUserInfoByAccount)
	engine.GET("/api/user-management/v1/user-auth", h.userAuth)
	engine.PUT("/api/user-management/v1/users/:user_id/pwd_err_info", h.updatePwdErrInfo)
	engine.GET("/api/user-management/v1/org_managers/:org_manager_ids/:fields", h.getOrgManagersInfo)
	engine.PATCH("/api/user-management/v1/users/:user_id", h.incrementUpdateUserInfoInternal)
	engine.GET("/api/user-management/v1/user-name-existed", observable.MiddlewareTrace(common.SvcARTrace), h.userNameExistdCheck)

	engine.GET("/api/user-management/v2/users", observable.MiddlewareTrace(common.SvcARTrace), h.getUserList)
}

// RegisterPublic 注册外部API
func (h *userRestHandler) RegisterPublic(engine *gin.Engine) {
	engine.GET("/api/user-management/v1/users/:user_ids/:fields", h.getUserBaseInfoInScope)
	engine.PUT("/api/user-management/v1/management/users/:user_id/:fields", h.updateUserInfo)
	engine.PATCH("/api/user-management/v1/management/users/:user_id", observable.MiddlewareTrace(common.SvcARTrace), h.incrementUpdateUserInfoExternal)
	engine.GET("/api/user-management/v1/profile/:fields", observable.MiddlewareTrace(common.SvcARTrace), h.getNorlmalUserInfo)
	engine.POST("/api/user-management/v1/profile/avatar", observable.MiddlewareTrace(common.SvcARTrace), h.updateAvatar)

	engine.GET("/api/user-management/v1/avatars/:user_ids", observable.MiddlewareTrace(common.SvcARTrace), h.getAvatarByIDs)
	engine.GET("/api/user-management/v1/query-user-by-telephone", h.getInfoByTel)

	// 管理控制台  主界面 部门下搜索
	engine.GET("/api/user-management/v1/console/search-users/:fields", observable.MiddlewareTrace(common.SvcARTrace), h.searchUsers)
}

// getUserInfo 获取用户信息
func (h *userRestHandler) getUserInfo(c *gin.Context) {
	fields := c.Param("fields")
	switch fields {
	case "department_ids":
		h.getAllBelongDepartmentID(c)
	case "accessor_ids":
		h.getAccessorIDsOfUser(c)
	default:
		h.getUserBaseInfo(c)
	}
}

// userNameCheck 检查用户名
func (h *userRestHandler) userNameExistdCheck(c *gin.Context) {
	name, ok := c.GetQuery("name")
	if !ok {
		err := rest.NewHTTPErrorV2(rest.BadRequest, "invalid name")
		rest.ReplyError(c, err)
		return
	}

	result, err := h.user.CheckUserNameExistd(c, name)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	outInfo := make(map[string]interface{})
	outInfo["result"] = result

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// searchUsers 用户成员搜索
func (h *userRestHandler) searchUsers(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 参数获取和检查
	keyScope := interfaces.UserSearchInDepartKeyScope{}
	key := interfaces.UserSearchInDepartKey{}

	key.DepartmentID, keyScope.BDepartmentID = c.GetQuery("department_id")
	key.Code, keyScope.BCode = c.GetQuery("code")
	key.Name, keyScope.BName = c.GetQuery("name")
	key.Account, keyScope.BAccount = c.GetQuery("account")
	key.ManagerName, keyScope.BManagerName = c.GetQuery("manager_name")
	key.DirectDepartCode, keyScope.BDirectDepartCode = c.GetQuery("direct_department_code")
	key.Position, keyScope.BPosition = c.GetQuery("position")

	var err error
	key.Offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || key.Offset < 0 {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid offset type")
		rest.ReplyError(c, err)
		return
	}

	key.Limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || (key.Limit < 1 || key.Limit > 1000) {
		err = rest.NewHTTPErrorV2(rest.BadRequest, "invalid limit type")
		rest.ReplyError(c, err)
		return
	}

	roleID, dRet := c.GetQuery("role")
	role := roleEnumIDMap[roleID]
	if !dRet || role == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid role", rest.SetDetail(map[string]interface{}{"params": "role"})))
		return
	}

	// 获取请求参数
	fields := strings.Split(c.Param("fields"), ",")
	userInfoRange, err := h.handlerInUserDBInfoRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取用户搜索信息
	userInfos, num, err := h.user.SearchUsers(c, &visitor, &keyScope, &key, userInfoRange, role)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 整理数据
	outInfo := make(map[string]interface{})
	entriesParm := make([]interface{}, 0)

	for k := range userInfos {
		userTemp := h.handleUserInfo(fields, &userInfos[k], true)
		userTemp["id"] = userInfos[k].ID
		userTemp["type"] = strUser

		entriesParm = append(entriesParm, userTemp)
	}

	outInfo["entries"] = entriesParm
	outInfo["total_count"] = num

	rest.ReplyOK(c, http.StatusOK, outInfo)
}

// getUserInfoByAccount 通过账户名匹配账户信息
func (h *userRestHandler) getUserInfoByAccount(c *gin.Context) {
	account, ok := c.GetQuery("account")
	if !ok || account == "" {
		rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "account"}))
		return
	}
	var enableIDCardLogin, enablePrefixMatch bool
	idCardLogin, ok := c.GetQuery("id_card_login")
	if ok {
		var err error
		enableIDCardLogin, err = strconv.ParseBool(idCardLogin)
		if err != nil {
			rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "id_card_login"}))
			return
		}
	}
	prefixMatch, ok := c.GetQuery("prefix_match")
	if ok {
		var err error
		enablePrefixMatch, err = strconv.ParseBool(prefixMatch)
		if err != nil {
			rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "prefix_match"}))
			return
		}
	}

	result, userBaseInfo, err := h.user.GetUserInfoByAccount(account, enableIDCardLogin, enablePrefixMatch)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	outData := make(map[string]interface{})
	outData["result"] = result
	if result {
		user := make(map[string]interface{})
		user["id"] = userBaseInfo.ID
		user["account"] = userBaseInfo.Account
		user["auth_type"] = h.authTypeMap[userBaseInfo.AuthType]
		user["pwd_err_cnt"] = userBaseInfo.PwdErrCnt
		user["pwd_err_last_time"] = userBaseInfo.PwdErrLastTime
		user["disable_status"] = !userBaseInfo.Enabled
		user["ldap_server_type"] = ladpServerTypeEnumMap[userBaseInfo.LDAPType]
		user["domain_path"] = userBaseInfo.DomainPath

		outData["user"] = user
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// userAuth 本地认证
func (h *userRestHandler) userAuth(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok || id == "" {
		rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "id"}))
		return
	}
	plainPassword, ok := c.GetQuery("password")
	if !ok || plainPassword == "" {
		rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest, map[string]interface{}{"params": "password"}))
		return
	}

	result, reason, err := h.user.UserAuth(id, plainPassword)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	outData := make(map[string]interface{})
	outData["result"] = result
	if !result {
		outData["reason"] = authFailedReasonEnumMap[reason]
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// updatePwdErrInfo 更新用户密码错误信息
func (h *userRestHandler) updatePwdErrInfo(c *gin.Context) {
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, h.pwdErrInfoSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}
	userID := c.Param("user_id")
	pwdErrCnt := int(jsonReq["pwd_err_cnt"].(float64))
	pwdErrLastTime := int64(jsonReq["pwd_err_last_time"].(float64))

	err := h.user.UpdatePwdErrInfo(userID, pwdErrCnt, pwdErrLastTime)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 响应204
	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// getAvatarByIDs 获取
func (h *userRestHandler) getAvatarByIDs(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取所有用户ID
	userIDs := strings.Split(c.Param("user_ids"), ",")
	oGetRange := interfaces.UserBaseInfoRange{
		ShowAvatar: true,
	}
	infos, err := h.user.GetUsersBaseInfo(c, &visitor, userIDs, oGetRange, true)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 返回值处理
	outData := make([]interface{}, len(infos))
	for k := range infos {
		tempData := make(map[string]string)
		tempData["id"] = infos[k].ID
		tempData["avatar_url"] = infos[k].Avatar
		outData[k] = tempData
	}
	rest.ReplyOK(c, http.StatusOK, outData)
}

// getNorlmalUserInfo 普通用户获取自己信息
func (h *userRestHandler) getNorlmalUserInfo(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取请求参数
	fields := strings.Split(c.Param("fields"), ",")
	var oGetRange interfaces.UserBaseInfoRange
	for _, v := range fields {
		switch v {
		case "avatar_url":
			oGetRange.ShowAvatar = true
		default:
			err := rest.NewHTTPError("invalid params", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取用户信息
	info, err := h.user.GetNorlmalUserInfo(c, &visitor, oGetRange)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 处理返回
	outData := h.handleUserOwnInfo(fields, &info)
	rest.ReplyOK(c, http.StatusOK, outData)
}

// updateUserInfo 更新头像信息
func (h *userRestHandler) updateAvatar(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 检查是否为表单上传
	formData, err := c.MultipartForm()
	if err != nil {
		err = rest.NewHTTPError(fmt.Sprintf("invalid request, multipart/form-data format error: %v", err), rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 检查参数
	files, ok := formData.File["avatar"]
	if !ok {
		err = rest.NewHTTPError("invalid params, avatar is required", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	// 只允许上传一个文件
	if len(files) != 1 {
		err = rest.NewHTTPError("invalid params, only support one file", rest.BadRequest, nil)
		rest.ReplyError(c, err)
		return
	}

	file := files[0]
	typ := file.Header.Get("Content-Type")

	// 获取文件缓存
	tmpFile, err := file.Open()
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	data := make([]byte, file.Size)
	_, err = tmpFile.Read(data)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}
	err = tmpFile.Close()
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 更新头像
	err = h.avatar.Update(c, &visitor, typ, data)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// updateUserInfo 更新用户信息
func (h *userRestHandler) updateUserInfo(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 获取用户信息,获取修改信息
	userID := c.Param("user_id")
	fields := strings.Split(c.Param("fields"), ",")
	var oUpdateRange interfaces.UserUpdateRange
	for _, v := range fields {
		switch v {
		case "password":
			oUpdateRange.UpdatePWD = true
		default:
			err := rest.NewHTTPError("invalid params", rest.BadRequest, nil)
			rest.ReplyError(c, err)
			return
		}
	}

	// 获取请求参数
	var jsonV interface{}
	if rErr := rest.GetJSONValue(c, &jsonV); rErr != nil {
		rest.ReplyError(c, rErr)
		return
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: oUpdateRange.UpdatePWD}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		rest.ReplyError(c, cErr)
		return
	}

	// 获取参数
	jsonObj := jsonV.(map[string]interface{})
	var userInfo interfaces.UserBaseInfo
	userInfo.ID = userID
	if oUpdateRange.UpdatePWD {
		userInfo.Password = jsonObj["password"].(string)
	}

	// 修改用户信息
	err := h.user.ModifyUserInfo(&visitor, oUpdateRange, &userInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// 增量更新用户信息内部接口
func (h *userRestHandler) incrementUpdateUserInfoInternal(c *gin.Context) {
	h.incrementUpdateUserInfo(c, nil)
}

// incrementUpdateUserInfo 增量更新用户信息
func (h *userRestHandler) incrementUpdateUserInfoExternal(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	h.incrementUpdateUserInfo(c, &visitor)
}

// 增量更新用户信息通用接口
func (h *userRestHandler) incrementUpdateUserInfo(c *gin.Context, visitor *interfaces.Visitor) {
	// 获取用户信息,获取修改信息
	userID := c.Param("user_id")
	var oUpdateRange interfaces.UserUpdateRange

	// jsonschema校验
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, h.incrementUserInfoSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	var userInfo interfaces.UserBaseInfo
	userInfo.ID = userID

	if jsonReq["custom_attr"] != nil {
		oUpdateRange.CustomAttr = true
		userInfo.CustomAttr = jsonReq["custom_attr"].(map[string]interface{})
	}

	// 增量修改用户信息
	err := h.user.IncrementModifyUserInfo(c, visitor, oUpdateRange, &userInfo)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusNoContent, nil)
}

// getAllBelongDepartmentID 根据用户id 获取用户所属部门id
func (h *userRestHandler) getAllBelongDepartmentID(c *gin.Context) {
	userID := c.Param("user_id")

	// 数据库获取deptid
	deptIDs, err := h.user.GetAllBelongDepartmentIDs(userID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	c.Writer.Header().Set("Location", "/api/user-management/v1/users/"+userID+"/department_ids")
	rest.ReplyOK(c, http.StatusOK, deptIDs)
}

// getAccessorIDsOfUser 获取指定用户的访问令牌
func (h *userRestHandler) getAccessorIDsOfUser(c *gin.Context) {
	userID := c.Param("user_id")

	// 获取指定用户的访问令牌
	accessorIDs, err := h.user.GetAccessorIDsOfUser(userID)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	rest.ReplyOK(c, http.StatusOK, accessorIDs)
}

// getUserList 获取用户列表
func (h *userRestHandler) getUserList(c *gin.Context) {
	// 检查fields范围
	tempFields := c.QueryArray("fields")
	fields := map[string]bool{
		"name":       true,
		"account":    true,
		"email":      true,
		"telephone":  true,
		"enabled":    true,
		"created_at": true,
		"frozen":     true,
	}
	for _, v := range tempFields {
		if _, ok := fields[v]; !ok {
			err := gerrors.NewError(gerrors.PublicBadRequest, "invalid field")
			rest.ReplyErrorV2(c, err)
			return
		}
	}

	// 获取用户信息类型检测
	userInfoRange, err := h.handlerInUserDBInfoRange(tempFields)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取其他参数，处理marker
	info, err := h.getUseristQueryParam(c)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取用户信息，必须要获取时间用于marker
	userInfoRange.ShowCreatedAt = true
	infos, num, hasNext, err := h.user.GetUserList(c, userInfoRange, info.direction, info.hasMarker, info.createdStamp, info.userID, info.limit)
	if err != nil {
		rest.ReplyErrorV2(c, err)
		return
	}

	// 获取marker数据
	nextMarkerStr := ""
	if hasNext {
		lastInfo := infos[len(infos)-1]
		nextMarkerMap := make(map[string]interface{})
		nextMarkerMap["created_at"] = lastInfo.CreatedAt
		nextMarkerMap["user_id"] = lastInfo.ID
		nextMarkerMap["direction"] = directionStringMap[info.direction]

		nextMarkerJSON, err := json.Marshal(nextMarkerMap)
		if err != nil {
			rest.ReplyErrorV2(c, err)
		}
		nextMarkerStr = base64.StdEncoding.EncodeToString(nextMarkerJSON)
	}

	// 数据整理
	out := make(map[string]interface{})
	outMsgs := make([]map[string]interface{}, 0)
	for i := range infos {
		tempData := h.handleUserInfo(tempFields, &infos[i], false)
		tempData["id"] = infos[i].ID
		outMsgs = append(outMsgs, tempData)
	}

	out["entries"] = outMsgs
	out["total_count"] = num
	out["next_marker"] = nextMarkerStr
	rest.ReplyOK(c, http.StatusOK, out)
}

// getUserInfo 获取用户基本信息
func (h *userRestHandler) getUserBaseInfo(c *gin.Context) {
	// 需要获取的用户信息类别
	userIDs := strings.Split(c.Param("user_id"), ",")
	fields := strings.Split(c.Param("fields"), ",")

	// 获取用户信息类型检测
	userInfoRange, err := h.handlerInUserDBInfoRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取相关的用户信息
	visitor := interfaces.Visitor{
		ErrorCodeType: GetErrorCodeType(c),
	}
	infos, err := h.user.GetUsersBaseInfo(c, &visitor, userIDs, userInfoRange, true)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 用户角色枚举处理
	outData := []map[string]interface{}{}
	for k := range infos {
		tempData := h.handleUserInfo(fields, &infos[k], false)
		tempData["id"] = infos[k].ID
		outData = append(outData, tempData)
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// getUserInfo 获取用户基本信息
func (h *userRestHandler) getUserBaseInfoByPost(c *gin.Context) {
	// jsonschema校验
	var jsonReq map[string]interface{}
	if err := validateAndBindGin(c, h.getUserInfoSchema, &jsonReq); err != nil {
		rest.ReplyError(c, err)
		return
	}

	userIDs := make([]string, 0)
	for _, v := range jsonReq["user_ids"].([]interface{}) {
		userIDs = append(userIDs, v.(string))
	}

	if jsonReq["method"].(string) != strHTTPGET {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid method"))
		return
	}

	bStrict := true
	if date, ok := jsonReq["strict"]; ok {
		bStrict = date.(bool)
	}

	fields := make([]string, 0)
	tempFields := jsonReq["fields"].([]interface{})
	for _, v := range tempFields {
		fields = append(fields, v.(string))
	}

	// 获取用户信息类型检测
	userInfoRange, err := h.handlerInUserDBInfoRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取相关的用户信息
	visitor := interfaces.Visitor{
		ErrorCodeType: GetErrorCodeType(c),
	}
	infos, err := h.user.GetUsersBaseInfo(c, &visitor, userIDs, userInfoRange, bStrict)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 用户角色枚举处理
	outData := []map[string]interface{}{}
	for k := range infos {
		tempData := h.handleUserInfo(fields, &infos[k], false)
		tempData["id"] = infos[k].ID
		outData = append(outData, tempData)
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// getInfoByTel 根据手机号获取用户基本信息
func (h *userRestHandler) getInfoByTel(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	telephone, ok := c.GetQuery("telephone")
	if !ok || telephone == "" {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid telephone"))
		return
	}

	fields, ok := c.GetQueryArray("field")
	if !ok {
		rest.ReplyError(c, rest.NewHTTPErrorV2(rest.BadRequest, "invalid field"))
		return
	}

	// 获取用户信息类型检测
	userInfoRange, err := h.handlerQueryUserRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取相关的用户信息
	result, info, err := h.user.GetUserBaseInfoByTelephone(c, &visitor, telephone, userInfoRange)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 用户角色枚举处理
	outData := map[string]interface{}{}
	outData["result"] = result
	if result {
		fields = append(fields, "telephone")
		tempData := h.handleUserInfo(fields, &info, false)
		tempData["id"] = info.ID

		outData["user"] = tempData
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// getUserBaseInfoInScope 获取调用者管辖范围内的用户信息
func (h *userRestHandler) getUserBaseInfoInScope(c *gin.Context) {
	// token验证
	visitor, vErr := verify(c, h.hydra)
	if vErr != nil {
		rest.ReplyError(c, vErr)
		return
	}

	// 需要获取的用户信息类别
	userIDs := strings.Split(c.Param("user_ids"), ",")
	fields := strings.Split(c.Param("fields"), ",")

	roleID, dRet := c.GetQuery("role")
	role := roleEnumIDMap[roleID]
	if !dRet || role == "" {
		rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest,
			map[string]interface{}{"params": "role"}))
		return
	}

	// 获取用户信息类型检测
	userInfoRange, err := h.handlerOutUserDBInfoRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取相关的用户信息
	info, err := h.user.GetUserBaseInfoInScope(&visitor, role, userIDs, userInfoRange)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 用户角色枚举处理vv
	outData := make([]interface{}, 0)
	var index int
	for index < len(info) {
		temp := h.handleUserInfo(fields, &info[index], false)
		temp["id"] = info[index].ID
		temp["type"] = "user"

		outData = append(outData, temp)
		index++
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// getPWDRetrievalMethod 根据账户名获取用户找回密码方式
func (h *userRestHandler) getPWDRetrievalMethod(c *gin.Context) {
	account, dRet := c.GetQuery("account")
	if !dRet || account == "" {
		rest.ReplyError(c, rest.NewHTTPError("invalid type", rest.BadRequest,
			map[string]interface{}{"params": "account"}))
		return
	}

	info, err := h.user.GetPWDRetrievalMethodByAccount(account)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	outData := make(map[string]interface{})
	outData["status"] = h.pwdRetrievalMap[info.Status]
	if info.Status == interfaces.PRSAvaliable {
		if info.BTelephone {
			outData["telephone"] = info.Telephone
		}
		if info.BEmail {
			outData[strEmail] = info.Email
		}
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// getOrgManagersInfo 获取组织管理员信息
func (h *userRestHandler) getOrgManagersInfo(c *gin.Context) {
	// 获取请求参数
	orgManagerIDs := strings.Split(c.Param("org_manager_ids"), ",")
	fields := strings.Split(c.Param("fields"), ",")

	// 获取用户信息类型检测
	orgManagerInfoRange, err := h.handlerOutOrgManagerInfoRange(fields)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 获取相关的组织管理员信息
	info, err := h.role.GetOrgManagersInfo(orgManagerIDs, orgManagerInfoRange)
	if err != nil {
		rest.ReplyError(c, err)
		return
	}

	// 用户角色枚举处理
	outData := make([]interface{}, 0)
	for k := range info {
		tempData := h.handleOrgManagerInfo(fields, &info[k])
		tempData["id"] = &info[k].ID
		outData = append(outData, tempData)
	}

	rest.ReplyOK(c, http.StatusOK, outData)
}

// 处理用户基本信息范围
func (h *userRestHandler) handlerOutUserDBInfoRange(fields []string) (userInfoRange interfaces.UserBaseInfoRange, err error) {
	for _, v := range fields {
		switch v {
		case strName1:
			userInfoRange.ShowName = true
		case "parent_dep_paths":
			userInfoRange.ShowParentDepPaths = true
		case strAccount:
			userInfoRange.ShowAccount = true
		default:
			err = rest.NewHTTPError("invalid type", rest.BadRequest,
				map[string]interface{}{"params": []string{0: v}})
			return
		}
	}
	return
}

// 处理组织管理员基本信息范围
func (h *userRestHandler) handlerOutOrgManagerInfoRange(fields []string) (orgManagerInfoRange interfaces.OrgManagerInfoRange, err error) {
	for _, v := range fields {
		switch v {
		case "sub_user_ids":
			orgManagerInfoRange.ShowSubUserIDs = true
		default:
			err = rest.NewHTTPError("invalid type", rest.BadRequest,
				map[string]interface{}{"params": []string{0: v}})
			return
		}
	}
	return
}

// 处理用户基本信息范围
//
//nolint:gocyclo
func (h *userRestHandler) handlerInUserDBInfoRange(fields []string) (userInfoRange interfaces.UserBaseInfoRange, err error) {
	for _, v := range fields {
		switch v {
		case "roles":
			userInfoRange.ShowRoles = true
		case "enabled":
			userInfoRange.ShowEnable = true
		case "priority":
			userInfoRange.ShowPriority = true
		case "csf_level":
			userInfoRange.ShowCSFLevel = true
		case strName1:
			userInfoRange.ShowName = true
		case "parent_deps":
			userInfoRange.ShowParentDeps = true
		case strAccount:
			userInfoRange.ShowAccount = true
		case "frozen":
			userInfoRange.ShowFrozen = true
		case "authenticated":
			userInfoRange.ShowAuthenticated = true
		case strEmail:
			userInfoRange.ShowEmail = true
		case "telephone":
			userInfoRange.ShowTelNumber = true
		case strThirdAttr:
			userInfoRange.ShowThirdAttr = true
		case strThirdID:
			userInfoRange.ShowThirdID = true
		case "auth_type":
			userInfoRange.ShowAuthType = true
		case "groups":
			userInfoRange.ShowGroups = true
		case "oss_id":
			userInfoRange.ShowOssID = true
		case "custom_attr":
			userInfoRange.ShowCustomAttr = true
		case strManager:
			userInfoRange.ShowManager = true
		case "remark":
			userInfoRange.ShowRemark = true
		case "code":
			userInfoRange.ShowCode = true
		case "position":
			userInfoRange.ShowPosition = true
		case "created_at":
			userInfoRange.ShowCreatedAt = true
		case "csf_level2":
			userInfoRange.ShowCSFLevel2 = true
		default:
			err = rest.NewHTTPError("invalid type", rest.BadRequest,
				map[string]interface{}{"params": []string{0: v}})
			return
		}
	}
	return
}

// 处理输出用户角色
func (h *userRestHandler) rolesToRoleNames(roles []interfaces.Role) (roleNames []string) {
	// 用户角色枚举处理
	roleNames = make([]string, 0)
	for _, v := range roles {
		name := h.roleNameMap[v]
		if name != "" {
			roleNames = append(roleNames, name)
		}
	}
	return
}

// 处理输出用户基本信息
//
//nolint:gocyclo
func (h *userRestHandler) handleUserInfo(fields []string, info *interfaces.UserBaseInfo, addParentDepCode bool) (outData map[string]interface{}) {
	// 处理输出用户角色
	roleNames := h.rolesToRoleNames(info.VecRoles)

	// 返回信息整理
	outData = make(map[string]interface{})
	groupInfos := make([]interface{}, 0)
	for _, v := range fields {
		switch v {
		case "roles":
			outData["roles"] = roleNames
		case "enabled":
			outData["enabled"] = info.Enabled
		case "priority":
			outData["priority"] = info.Priority
		case "csf_level":
			outData["csf_level"] = info.CSFLevel
		case strName1:
			outData[strName1] = info.Name
		case "parent_deps":
			tempAll := make([]interface{}, 0, len(info.ParentDeps))
			for k := range info.ParentDeps {
				temp := make([]interface{}, 0, len(info.ParentDeps[k]))
				for k1 := range info.ParentDeps[k] {
					tempDepInfo := make(map[string]interface{})
					tempDepInfo["id"] = info.ParentDeps[k][k1].ID
					tempDepInfo["name"] = info.ParentDeps[k][k1].Name
					tempDepInfo["type"] = strDepartment

					if addParentDepCode {
						tempDepInfo["code"] = info.ParentDeps[k][k1].Code
					}

					temp = append(temp, tempDepInfo)
				}
				tempAll = append(tempAll, temp)
			}
			outData["parent_deps"] = tempAll
		case strAccount:
			outData[strAccount] = info.Account
		case "parent_dep_paths":
			outData["parent_dep_paths"] = info.ParentDepPaths
		case "frozen":
			outData["frozen"] = info.Frozen
		case "authenticated":
			outData["authenticated"] = info.Authenticated
		case strEmail:
			outData[strEmail] = info.Email
		case "telephone":
			outData["telephone"] = info.TelNumber
		case strThirdAttr:
			outData[strThirdAttr] = info.ThirdAttr
		case strThirdID:
			outData[strThirdID] = info.ThirdID
		case "auth_type":
			outData["auth_type"] = h.authTypeMap[info.AuthType]
		case "groups":
			for _, value := range info.Groups {
				groupInfo := make(map[string]interface{})
				groupInfo["id"] = value.ID
				groupInfo[strName1] = value.Name
				groupInfo["type"] = "group"
				groupInfos = append(groupInfos, groupInfo)
			}
			outData["groups"] = groupInfos
		case "oss_id":
			outData["oss_id"] = info.OssID
		case "custom_attr":
			outData["custom_attr"] = info.CustomAttr
		case strManager:
			if info.Manager.ID != "" {
				managerInfo := make(map[string]interface{})
				managerInfo["id"] = info.Manager.ID
				managerInfo["name"] = info.Manager.Name
				managerInfo["type"] = strUser
				outData[strManager] = managerInfo
			}
		case "remark":
			outData["remark"] = info.Remark
		case "created_at":
			outData["created_at"] = time.Unix(info.CreatedAt, 0).Format(time.RFC3339)
		case "position":
			outData["position"] = info.Position
		case "code":
			outData["code"] = info.Code
		case "csf_level2":
			outData["csf_level2"] = info.CSFLevel2
		}
	}
	return
}

// 处理输出组织管理员信息
//
//nolint:gocritic
func (h *userRestHandler) handleOrgManagerInfo(fields []string, info *interfaces.OrgManagerInfo) (outData map[string]interface{}) {
	// 返回信息整理
	outData = make(map[string]interface{})
	for _, v := range fields {
		switch v {
		case "sub_user_ids":
			if len(info.SubUserIDs) > 0 {
				outData["sub_user_ids"] = info.SubUserIDs
			} else {
				outData["sub_user_ids"] = []string{}
			}
		}
	}
	return
}

// 处理输出用户自己的信息
func (h *userRestHandler) handleUserOwnInfo(fields []string, info *interfaces.UserBaseInfo) (outData map[string]interface{}) {
	// 返回信息整理
	outData = make(map[string]interface{})
	for _, v := range fields {
		if v == "avatar_url" {
			outData["avatar_url"] = info.Avatar
		}
	}
	return
}

// 处理用户基本信息范围
func (h *userRestHandler) handlerQueryUserRange(fields []string) (userInfoRange interfaces.UserBaseInfoRange, err error) {
	userInfoRange.ShowTelNumber = true
	for _, v := range fields {
		switch v {
		case strName1:
			userInfoRange.ShowName = true
		case strAccount:
			userInfoRange.ShowAccount = true
		case strEmail:
			userInfoRange.ShowEmail = true
		case strThirdID:
			userInfoRange.ShowThirdID = true
		default:
			err = rest.NewHTTPError("invalid type", rest.BadRequest,
				map[string]interface{}{"params": []string{0: v}})
			return
		}
	}
	return
}

func (h *userRestHandler) getUseristQueryParam(c *gin.Context) (info userListQueryInfo, err error) {
	// 获取marker
	markerStr := ""
	markerStr, info.hasMarker = c.GetQuery("marker")

	// 获取分页数据
	info.limit, err = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		err = gerrors.NewError(gerrors.PublicBadRequest, "limit is illeagal")
		return
	}
	if info.limit < 1 || info.limit > 1000 {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid limit([ 1 .. 1000 ])")
		return
	}

	// 如果存在marker，则不获取direction等信息
	if info.hasMarker {
		// base64 解码
		var markerJSONData []byte
		markerJSONData, err = base64.StdEncoding.DecodeString(markerStr)
		if err != nil {
			err = gerrors.NewError(gerrors.PublicBadRequest, "invalid marker base64")
			return
		}

		// 获取参数
		var jsonReq interface{}
		if err = validateAndBindNewError(markerJSONData, h.markerSchema, &jsonReq); err != nil {
			return
		}

		// 获取参数
		info.createdStamp = int64(jsonReq.(map[string]interface{})["created_at"].(float64))
		info.userID = jsonReq.(map[string]interface{})["user_id"].(string)
		strDirection := jsonReq.(map[string]interface{})["direction"].(string)
		info.direction = directionMap[strDirection]
		return
	}

	// 获取排序方向
	directionStr, dRet := c.GetQuery(("direction"))
	if dRet && directionStr != strAsc && directionStr != strDesc {
		err = gerrors.NewError(gerrors.PublicBadRequest, "invalid direction type")
		return
	}

	if !dRet {
		directionStr = strDesc
	}
	info.direction = directionMap[directionStr]

	return
}
