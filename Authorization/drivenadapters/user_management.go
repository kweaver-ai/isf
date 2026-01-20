// Package drivenadapters 当前微服务依赖的其他服务
package drivenadapters

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/mitchellh/mapstructure"

	"Authorization/common"
	"Authorization/interfaces"
)

type orgNameIDInfo struct {
	UserIDs      map[string]string
	DepartIDs    map[string]string
	ContactorIDs map[string]string
	GroupIDs     map[string]string
	AppIDs       map[string]string
}

type orgIDInfo struct {
	UserIDs      []string
	DepartIDs    []string
	ContactorIDs []string
	GroupIDs     []string
	AppIDs       []string
}

var (
	usermgntOnce sync.Once
	usermgnt     *usermgntSvc
)

var accessorIDsWarningLen = 5000

type usermgntSvc struct {
	baseURL         string
	log             common.Logger
	traceHTTPClient httpclient.HTTPClient
	roleTypeMap     map[string]interfaces.SystemRoleType
}

// NewUserMgnt 创建UserMgnt服务处理对象
func NewUserMgnt() *usermgntSvc {
	usermgntOnce.Do(func() {
		config := common.SvcConfig
		roleTypeMap := map[string]interfaces.SystemRoleType{
			"super_admin": interfaces.SuperAdmin,
			"sys_admin":   interfaces.SystemAdmin,
			"audit_admin": interfaces.AuditAdmin,
			"sec_admin":   interfaces.SecurityAdmin,
			"org_manager": interfaces.OrganizationAdmin,
			"org_audit":   interfaces.OrganizationAudit,
			"normal_user": interfaces.NormalUser,
		}
		usermgnt = &usermgntSvc{
			baseURL:         fmt.Sprintf("http://%s:%d/api/user-management", config.UserMgntPrivateHost, config.UserMgntPrivatePort),
			log:             common.NewLogger(),
			traceHTTPClient: httpclient.NewHTTPClient(common.SvcARTrace),
			roleTypeMap:     roleTypeMap,
		}
	})

	return usermgnt
}

// GetAccessorIDsByUserID 获取指定用户的访问令牌
func (u *usermgntSvc) GetAccessorIDsByUserID(ctx context.Context, userID string) (accessorIDs []string, err error) {
	headers := map[string]string{
		"x-error-code": "string",
	}
	target := fmt.Sprintf("%s/v1/users/%s/accessor_ids", u.baseURL, userID)
	respParam, err := u.traceHTTPClient.Get(ctx, target, headers)
	if err != nil {
		u.log.Errorf("GetAccessorIdsByUserID failed:%v, url:%v", err, target)
		return
	}

	accessorArr := respParam.([]any)

	length := len(accessorArr)
	// 访问令牌数量过多时，会使后面的逻辑处理变慢，如数据库的查询
	// 增加告警日志，方便排查问题, 如审核测试场景，内部组过多但未清理
	if length > accessorIDsWarningLen {
		u.log.Warnf("user_id:%v, accessor_ids len:%v", userID, length)
	}
	accessorIDs = make([]string, 0, length)
	for _, v := range accessorArr {
		accessorIDs = append(accessorIDs, v.(string))
	}

	return
}

// GetUserRolesByUserID 通过用户id获取角色
func (u *usermgntSvc) GetUserRolesByUserID(ctx context.Context, userID string) (roleTypes []interfaces.SystemRoleType, err error) {
	headers := map[string]string{
		"x-error-code": "string",
	}
	fields := "roles"
	target := fmt.Sprintf("%s/v1/users/%s/%s", u.baseURL, userID, fields)
	respParam, err := u.traceHTTPClient.Get(ctx, target, headers)
	if err != nil {
		u.log.Errorf("GetUserRolesByUserID failed:%v, url:%v", err, target)
		return
	}
	info := respParam.([]any)[0]
	rolesParam := info.(map[string]any)["roles"].([]any)
	for _, val := range rolesParam {
		roleType, ok := u.roleTypeMap[val.(string)]
		if !ok {
			err = errors.New("role type conversion error")
			return
		}
		roleTypes = append(roleTypes, roleType)
	}
	return
}

// BatchGetUserInfoByID 批量获取用户的基础信息
func (u *usermgntSvc) BatchGetUserInfoByID(ctx context.Context, userIDs []string) (userInfoMap map[string]interfaces.UserInfo, err error) {
	userInfoMap = make(map[string]interfaces.UserInfo)
	var userIDsStr string
	if len(userIDs) == 0 {
		return
	}
	for i, userID := range userIDs {
		userIDsStr += userID
		if i != len(userIDs)-1 {
			userIDsStr += ","
		}
	}
	fields := "account,name,csf_level,frozen,roles,email,telephone,third_attr,third_id,parent_deps"
	target := fmt.Sprintf("%s/v1/users/%s/%s", u.baseURL, userIDsStr, fields)
	headers := map[string]string{
		"x-error-code": "string",
	}
	respParam, err := u.traceHTTPClient.Get(ctx, target, headers)
	if err != nil {
		u.log.Errorf("BatchGetUserInfoByID failed:%v, url:%v", err, target)
		return
	}
	infos := respParam.([]any)
	for i := range infos {
		info := infos[i].(map[string]any)
		userInfo, errTmp := u.convertUserInfo(info)
		if errTmp != nil {
			return userInfoMap, errTmp
		}
		userInfo.ID = info["id"].(string)
		userInfo.UserType = interfaces.AccessorUser
		userInfoMap[userInfo.ID] = userInfo
	}
	return
}

func (u *usermgntSvc) convertUserInfo(info map[string]any) (userInfo interfaces.UserInfo, err error) {
	userInfo = interfaces.UserInfo{
		Account:    info["account"].(string),
		VisionName: info["name"].(string),
		CsfLevel:   int(info["csf_level"].(float64)),
		Frozen:     info["frozen"].(bool),
		Roles:      make(map[interfaces.SystemRoleType]bool),
		Email:      info["email"].(string),
		Telephone:  info["telephone"].(string),
		ThirdAttr:  info["third_attr"].(string),
		ThirdID:    info["third_id"].(string),
	}
	roles := info["roles"].([]any)
	for _, val := range roles {
		roleType, ok := u.roleTypeMap[val.(string)]
		if !ok {
			err = errors.New("role type conversion error")
			return
		}
		userInfo.Roles[roleType] = true
	}

	err = mapstructure.Decode(info["parent_deps"], &userInfo.ParentDeps)
	if err != nil {
		return interfaces.UserInfo{}, err
	}

	return
}

// GetParentDepartmentsByDepartmentID 根据部门ID获取父部门信息
func (u *usermgntSvc) GetParentDepartmentsByDepartmentID(ctx context.Context, departmentID string) (parentDeps []interfaces.Department, err error) {
	target := fmt.Sprintf("%s/v1/departments/%s/parent_deps", u.baseURL, departmentID)
	headers := map[string]string{
		"x-error-code": "string",
	}
	respParam, err := u.traceHTTPClient.Get(ctx, target, headers)
	if err != nil {
		u.log.Errorf("GetParentDepartmentsByDepartmentID failed:%v, url:%v", err, target)
		return
	}

	info := respParam.([]any)[0]
	err = mapstructure.Decode(info.(map[string]any)["parent_deps"], &parentDeps)
	if err != nil {
		return nil, err
	}
	return
}

//nolint:staticcheck,gocyclo
func (u *usermgntSvc) GetNameByAccessorIDs(ctx context.Context, accessorIDs map[string]interfaces.AccessorType) (accessorNames map[string]string, err error) {
	var orgInfo orgIDInfo
	orgInfo.UserIDs = make([]string, 0)
	orgInfo.DepartIDs = make([]string, 0)
	orgInfo.ContactorIDs = make([]string, 0)
	orgInfo.GroupIDs = make([]string, 0)
	orgInfo.AppIDs = make([]string, 0)
	for accessorID, accessorType := range accessorIDs {
		if accessorType == interfaces.AccessorUser {
			orgInfo.UserIDs = append(orgInfo.UserIDs, accessorID)
		} else if accessorType == interfaces.AccessorDepartment {
			orgInfo.DepartIDs = append(orgInfo.DepartIDs, accessorID)
		} else if accessorType == interfaces.AccessorContactor {
			orgInfo.ContactorIDs = append(orgInfo.ContactorIDs, accessorID)
		} else if accessorType == interfaces.AccessorGroup {
			orgInfo.GroupIDs = append(orgInfo.GroupIDs, accessorID)
		} else if accessorType == interfaces.AccessorApp {
			orgInfo.AppIDs = append(orgInfo.AppIDs, accessorID)
		}
	}

	orgNameInfo, err := u.getOrgNameIDInfo(ctx, &orgInfo)
	if err != nil {
		u.log.Errorf("GetNameByAccessorID err:%v", err)
	}
	accessorNames = make(map[string]string)
	for accessorID, accessorType := range accessorIDs {
		if accessorType == interfaces.AccessorUser {
			if value, ok := orgNameInfo.UserIDs[accessorID]; ok {
				accessorNames[accessorID] = value
			}
		} else if accessorType == interfaces.AccessorDepartment {
			if value, ok := orgNameInfo.DepartIDs[accessorID]; ok {
				accessorNames[accessorID] = value
			}
		} else if accessorType == interfaces.AccessorContactor {
			if value, ok := orgNameInfo.ContactorIDs[accessorID]; ok {
				accessorNames[accessorID] = value
			}
		} else if accessorType == interfaces.AccessorGroup {
			if value, ok := orgNameInfo.GroupIDs[accessorID]; ok {
				accessorNames[accessorID] = value
			}
		} else if accessorType == interfaces.AccessorApp {
			if value, ok := orgNameInfo.AppIDs[accessorID]; ok {
				accessorNames[accessorID] = value
			}
		}
	}
	return
}

func (u *usermgntSvc) getOrgNameIDInfo(ctx context.Context, orgInfo *orgIDInfo) (orgNameInfo orgNameIDInfo, err error) {
	tmpInfo := map[string]any{
		"method":         "GET",
		"user_ids":       orgInfo.UserIDs,
		"department_ids": orgInfo.DepartIDs,
		"contactor_ids":  orgInfo.ContactorIDs,
		"group_ids":      orgInfo.GroupIDs,
		"app_ids":        orgInfo.AppIDs,
	}
	target := fmt.Sprintf("%v/v2/names", u.baseURL)
	headers := map[string]string{
		"x-error-code": "string",
	}
	_, respParam, err := u.traceHTTPClient.Post(ctx, target, headers, tmpInfo)
	if err != nil {
		u.log.Errorf("getOrgNameIDInfo failed: %v, url: %v", err, target)
		return
	}

	userNameInfos := respParam.(map[string]any)["user_names"].([]any)
	orgNameInfo.UserIDs = make(map[string]string)
	for _, x := range userNameInfos {
		id := x.(map[string]any)["id"].(string)
		name := x.(map[string]any)["name"].(string)
		orgNameInfo.UserIDs[id] = name
	}
	orgNameInfo.DepartIDs = make(map[string]string)
	departNameInfos := respParam.(map[string]any)["department_names"].([]any)
	for _, x := range departNameInfos {
		id := x.(map[string]any)["id"].(string)
		name := x.(map[string]any)["name"].(string)
		orgNameInfo.DepartIDs[id] = name
	}
	orgNameInfo.ContactorIDs = make(map[string]string)
	conatctorNameInfos := respParam.(map[string]any)["contactor_names"].([]any)
	for _, x := range conatctorNameInfos {
		id := x.(map[string]any)["id"].(string)
		name := x.(map[string]any)["name"].(string)
		orgNameInfo.ContactorIDs[id] = name
	}
	orgNameInfo.GroupIDs = make(map[string]string)
	groupNameInfos := respParam.(map[string]any)["group_names"].([]any)
	for _, x := range groupNameInfos {
		id := x.(map[string]any)["id"].(string)
		name := x.(map[string]any)["name"].(string)
		orgNameInfo.GroupIDs[id] = name
	}
	orgNameInfo.AppIDs = make(map[string]string)
	appNameInfos := respParam.(map[string]any)["app_names"].([]any)
	for _, x := range appNameInfos {
		id := x.(map[string]any)["id"].(string)
		name := x.(map[string]any)["name"].(string)
		orgNameInfo.AppIDs[id] = name
	}
	return
}
