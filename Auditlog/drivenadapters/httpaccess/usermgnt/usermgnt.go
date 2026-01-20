package usermgnt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"AuditLog/common"
	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/models"
)

var (
	uOnce sync.Once
	uM    *userMgnt
)

type userMgnt struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewUserMgnt() interfaces.UserMgntRepo {
	uOnce.Do(func() {
		uM = &userMgnt{
			urlPrefix:  common.SvcConfig.UserMgntPrivateProtocol + "://" + common.SvcConfig.UserMgntPrivateHost + ":" + common.SvcConfig.UserMgntPrivatePort + "/api/user-management/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})
	return uM
}

// 获取用户信息
func (u *userMgnt) GetUserInfoByID(userIDs []string) (userinfos []models.User, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	userinfos = make([]models.User, 0)
	fields := "name,account,parent_deps,roles,csf_level,enabled,frozen,email,telephone,groups"

	if len(userIDs) == 0 {
		return
	}
	// 切片去重
	userIDs = common.ArrayRemoveDuplicate(userIDs)
	userIDsStr := strings.Join(userIDs, ",")
	addr := u.urlPrefix + fmt.Sprintf("/users/%v/%v", userIDsStr, fields)
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetUserInfoByID:%v\n", err)
		return
	}

	statusCode = resp.StatusCode
	if statusCode == http.StatusNotFound {
		return
	}
	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetUserInfoByID url: %v, statusCode: %v", addr, statusCode)
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&userinfos)
	if err != nil {
		u.logger.Infof("ERROR: GetUserInfoByID:%v\n", err)
		return
	}

	return
}

// 获取应用账户信息
func (u *userMgnt) GetAppInfoByID(id string) (appInfos models.App, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := u.urlPrefix + fmt.Sprintf("/apps/%v", id)
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetAppInfoByID:%v\n", err)
		return
	}
	statusCode = resp.StatusCode
	if statusCode == http.StatusNotFound {
		return
	}
	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetAppInfoByID url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&appInfos)
	if err != nil {
		u.logger.Infof("ERROR: GetAppInfoByID:%v\n", err)
		return
	}

	return
}

// 获取部门下所有用户ID
func (u *userMgnt) GetDeptAllUserIDs(deptID string) (userIDs map[string][]string, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := u.urlPrefix + fmt.Sprintf("/departments/%v/all_user_ids", deptID)
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetDeptAllUserIDs:%v\n", err)
		return
	}
	statusCode = resp.StatusCode

	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetDeptAllUserIDs url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&userIDs)
	if err != nil {
		u.logger.Infof("ERROR: GetDeptAllUserIDs:%v\n", err)
		return
	}

	return
}

// 获取角色用户ID
func (u *userMgnt) GetUserIDsByRoleNames(roleNames []string) (roleMemberInfos []*models.RoleMemberInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := u.urlPrefix + fmt.Sprintf("/role-members/%v", strings.Join(roleNames, ","))
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetUserIDsByRoleNames:%v\n", err)
		return
	}
	statusCode = resp.StatusCode

	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetUserIDsByRoleNames url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&roleMemberInfos)
	if err != nil {
		u.logger.Infof("ERROR: GetUserIDsByRoleNames:%v\n", err)
		return
	}

	return
}

// 获取指定部门的部门信息
func (u *userMgnt) GetDeptInfoByIDs(deptIDs []string) (deptInfo []*models.DeptInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := u.urlPrefix + fmt.Sprintf("/departments/%v/name,parent_deps,managers", strings.Join(deptIDs, ","))
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetDeptInfoByIDs:%v\n", err)
		return
	}
	statusCode = resp.StatusCode

	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetDeptInfoByIDs url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&deptInfo)
	if err != nil {
		u.logger.Infof("ERROR: GetDeptInfoByID:%v\n", err)
		return
	}

	return
}

// 获取指定级别的部门
func (u *userMgnt) GetDepsByLevel(level int) (deptInfos []*models.DepInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := u.urlPrefix + fmt.Sprintf("/departments?level=%v", level)
	resp, err = u.httpClient.Get(ctx, addr)
	if err != nil {
		u.logger.Errorf("GetDepsByLevel:%v\n", err)
		return
	}
	statusCode = resp.StatusCode

	if statusCode != http.StatusOK {
		u.logger.Warnf("ERROR: GetDepsByLevel url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&deptInfos)
	if err != nil {
		u.logger.Infof("ERROR: GetDepsByLevel:%v\n", err)
		return
	}

	return
}
