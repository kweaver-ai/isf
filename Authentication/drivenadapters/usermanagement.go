package drivenadapters

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"

	"github.com/kweaver-ai/go-lib/httpclient"

	"Authentication/common"
	"Authentication/interfaces"
)

var (
	uOnce sync.Once
	u     *userManagement
)

type userManagement struct {
	httpClient            httpclient.HTTPClient
	httpClient2           httpclient.HTTPClient
	log                   common.Logger
	trace                 observable.Tracer
	privateAddr           string
	roleTypeMap           map[string]interfaces.RoleType
	authTypeMap           map[string]interfaces.AuthType
	ladpServerTypeEnumMap map[string]interfaces.LDAPServerType
}

// NewUserManagement 创建UserManagement接口操作对象
func NewUserManagement() *userManagement {
	uOnce.Do(func() {
		config := common.SvcConfig
		u = &userManagement{
			httpClient:  httpclient.NewHTTPClient(common.SvcARTrace),
			httpClient2: httpclient.NewHTTPClient(common.SvcARTrace),
			log:         common.NewLogger(),
			trace:       common.SvcARTrace,
			roleTypeMap: map[string]interfaces.RoleType{
				"super_admin": interfaces.SuperAdmin,
				"sys_admin":   interfaces.SystemAdmin,
				"audit_admin": interfaces.AuditAdmin,
				"sec_admin":   interfaces.SecurityAdmin,
				"org_manager": interfaces.OrganizationAdmin,
				"org_audit":   interfaces.OrganizationAudit,
				"normal_user": interfaces.NormalUser,
			},
			authTypeMap: map[string]interfaces.AuthType{
				"local":  interfaces.Local,
				"domain": interfaces.Domain,
				"third":  interfaces.Third,
			},
			ladpServerTypeEnumMap: map[string]interfaces.LDAPServerType{
				"windows_ad": interfaces.WindowAD,
				"other_ldap": interfaces.OtherLDAP,
			},
			privateAddr: fmt.Sprintf("http://%s:%d", config.UserManagementPrivateHost, config.UserManagementPrivatePort),
		}
	})

	return u
}

func (u *userManagement) AnonymousAuthentication(ctx context.Context, visitor *interfaces.Visitor, account, password, referrer string) (bool, error) {
	var err error
	u.trace.SetClientSpanName("适配器-匿名认证")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	permInfo := map[string]interface{}{
		"account":  account,
		"password": password,
	}

	target := fmt.Sprintf("%v/api/user-management/v1/anonymity-auth", u.privateAddr)
	headers := map[string]string{
		"x-referrer":   referrer,
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	_, resParam, err := u.httpClient.Post(newCtx, target, headers, permInfo)
	if err != nil {
		u.log.Errorf("Anonymous authentication failed: %v, url: %v", err, target)
		return false, err
	}

	return resParam.(map[string]interface{})["result"].(bool), nil
}

// GetUserNameByUserID 通过用户id获取用户名
func (u *userManagement) GetUserRolesByUserID(ctx context.Context, visitor *interfaces.Visitor, userID string) (roleTypes []interfaces.RoleType, err error) {
	u.trace.SetClientSpanName("适配器-获取用户角色")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	fields := "roles"
	target := fmt.Sprintf("%s/api/user-management/v1/users/%s/%s", u.privateAddr, userID, fields)
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	respParam, err := u.httpClient.Get(newCtx, target, headers)
	if err != nil {
		u.log.Errorf("GetUserRolesByUserID failed:%v, url:%v", err, target)
		return
	}
	rolesParam := respParam.([]interface{})[0].(map[string]interface{})["roles"].([]interface{})
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

// GetAppInfo 获取应用账户信息
func (u *userManagement) GetAppInfo(ctx context.Context, visitor *interfaces.Visitor, appID string) (info interfaces.AppInfo, err error) {
	u.trace.SetClientSpanName("适配器-获取应用账户信息")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	target := fmt.Sprintf("%s/api/user-management/v1/apps/%s", u.privateAddr, appID)
	respParam, err := u.httpClient.Get(newCtx, target, map[string]string{"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType]})
	if err != nil {
		u.log.Errorf("GetAppInfo failed:%v, url:%v", err, target)
		return
	}
	info.ID = appID
	info.Name = respParam.(map[string]interface{})["name"].(string)
	return
}

// AccountMatch 账户匹配
func (u *userManagement) AccountMatch(ctx context.Context, visitor *interfaces.Visitor, account string, idCardLogin, prefixMatch bool) (result bool, info interfaces.UserBaseInfo, err error) {
	u.trace.SetClientSpanName("适配器-账户匹配")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	target := fmt.Sprintf("%s/api/user-management/v1/account-match?account=%s&id_card_login=%v&prefix_match=%v", u.privateAddr, url.QueryEscape(account), idCardLogin, prefixMatch)
	resParam, err := u.httpClient.Get(newCtx, target, map[string]string{"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType]})
	if err != nil {
		u.log.Errorf("AccountMatch failed:%v, url:%v", err, target)
		return
	}

	result = resParam.(map[string]interface{})["result"].(bool)
	if result {
		userInfo := resParam.(map[string]interface{})["user"]
		info.ID = userInfo.(map[string]interface{})["id"].(string)
		info.Account = userInfo.(map[string]interface{})["account"].(string)
		info.AuthType = u.authTypeMap[userInfo.(map[string]interface{})["auth_type"].(string)]
		info.PwdErrCnt = int(userInfo.(map[string]interface{})["pwd_err_cnt"].(float64))
		info.PwdErrLastTime = int64(userInfo.(map[string]interface{})["pwd_err_last_time"].(float64))
		info.DisableStatus = userInfo.(map[string]interface{})["disable_status"].(bool)
		info.LDAPType = u.ladpServerTypeEnumMap[userInfo.(map[string]interface{})["ldap_server_type"].(string)]
		info.DomainPath = userInfo.(map[string]interface{})["domain_path"].(string)
	}
	return
}

// UserAuth 账户认证
func (u *userManagement) UserAuth(ctx context.Context, visitor *interfaces.Visitor, userID, password string) (result bool, reason string, err error) {
	u.trace.SetClientSpanName("UserMgnt-用户身份认证")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	target := fmt.Sprintf("%s/api/user-management/v1/user-auth?id=%s&password=%s", u.privateAddr, userID, url.QueryEscape(password))
	resParam, err := u.httpClient2.Get(newCtx, target, map[string]string{"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType]})
	if err != nil {
		u.log.Errorf("UserAuth failed:%v, url:%v", err, fmt.Sprintf("%s/api/user-management/v1/user-auth?id=%s", u.privateAddr, userID))
		return
	}
	result = resParam.(map[string]interface{})["result"].(bool)
	if !result {
		reason = resParam.(map[string]interface{})["reason"].(string)
	}

	return
}

// UpdatePWDErrInfo 更新账户密码错误信息
func (u *userManagement) UpdatePWDErrInfo(ctx context.Context, visitor *interfaces.Visitor, userID string, pwdErrCnt int) (err error) {
	u.trace.SetClientSpanName("UserMgnt-更新账户密码错误信息")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	target := fmt.Sprintf("%s/api/user-management/v1/users/%s/pwd_err_info", u.privateAddr, userID)
	reqBody := map[string]interface{}{
		"pwd_err_cnt":       pwdErrCnt,
		"pwd_err_last_time": time.Now().Unix(),
	}
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	statusCode, _, err := u.httpClient2.Put(newCtx, target, headers, reqBody)
	if err != nil {
		u.log.Errorf("UpdatePWDErrInfo failed:%v, url:%v", err, target)
		return
	}

	if statusCode != http.StatusNoContent {
		err = rest.NewHTTPError("UpdatePWDErrInfo failed", statusCode, nil)
	}

	return
}

func (u *userManagement) GetAnonymityInfoByID(ctx context.Context, visitor *interfaces.Visitor, anonymityID string) (verifyMobile bool, err error) {
	u.trace.SetClientSpanName("UserMgnt-根据ID获取匿名账户信息")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	target := fmt.Sprintf("%s/api/user-management/v1/anonymity/%s", u.privateAddr, anonymityID)

	resParam, err := u.httpClient2.Get(newCtx, target, map[string]string{"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType]})
	if err != nil {
		return false, err
	}

	return resParam.(map[string]interface{})["verify_mobile"].(bool), nil
}

// GetUserInfo 获取实名账户信息
func (u *userManagement) GetUserInfo(ctx context.Context, visitor *interfaces.Visitor, userID string) (info *interfaces.UserBaseInfo, err error) {
	u.trace.SetClientSpanName("UserMgnt-获取实名账户信息")
	newCtx, span := u.trace.AddClientTrace(ctx)
	defer func() { u.trace.TelemetrySpanEnd(span, err) }()

	fields := "account"
	target := fmt.Sprintf("%s/api/user-management/v1/users/%s/%s", u.privateAddr, userID, fields)
	resParam, err := u.httpClient2.Get(newCtx, target, map[string]string{"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType]})
	if err != nil {
		return nil, err
	}

	info = &interfaces.UserBaseInfo{
		Account: resParam.([]interface{})[0].(map[string]interface{})["account"].(string),
	}
	return info, nil
}
