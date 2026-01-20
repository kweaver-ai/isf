package drivenadapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"policy_mgnt/common/config"
	"policy_mgnt/interfaces"
	"sync"

	"policy_mgnt/common"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/observable"
)

var (
	userManagementOnce      sync.Once
	userManagementSingleton *userManagement
)

type userManagement struct {
	log               common.Logger
	rawClient         *http.Client
	client            httpclient.HTTPClient
	userManagementURL string
	trace             observable.Tracer
	mapStringRole     map[string]interfaces.Role
}

// NewUserManagement 创建用户管理驱动
func NewUserManagement() *userManagement {
	userManagementOnce.Do(func() {
		config := config.Config.UserMgmtPvt
		userManagementSingleton = &userManagement{
			log:               common.NewLogger(),
			rawClient:         httpclient.NewRawHTTPClient(),
			client:            httpclient.NewHTTPClient(common.SvcARTrace),
			userManagementURL: fmt.Sprintf("http://%s:%s/api/user-management/v1", config.Host, config.Port),
			mapStringRole: map[string]interfaces.Role{
				"super_admin": interfaces.SystemRoleSuperAdmin,
				"sys_admin":   interfaces.SystemRoleSysAdmin,
				"sec_admin":   interfaces.SystemRoleSecAdmin,
				"audit_admin": interfaces.SystemRoleAuditAdmin,
				"org_manager": interfaces.SystemRoleOrgManager,
				"org_audit":   interfaces.SystemRoleOrgAudit,
				"normal_user": interfaces.SystemRoleNormalUser,
			},
			trace: common.SvcARTrace,
		}
	})

	return userManagementSingleton
}

// GetUserInfos 获取许可证
func (l *userManagement) GetUserInfos(ctx context.Context, ids []string) (infos map[string]interfaces.UserInfo, err error) {
	l.trace.SetClientSpanName("获取用户信息")
	newCtx, span := l.trace.AddClientTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	// 整合请求参数
	infos = make(map[string]interfaces.UserInfo)
	body, err := json.Marshal(map[string]interface{}{
		"method":   "GET",
		"user_ids": ids,
		"fields":   []string{"name", "roles"},
	})
	if err != nil {
		l.log.Errorln("user_management GetUserInfos request body marshal err: %v", err)
		return
	}

	tempUrl := l.userManagementURL + "/batch-get-user-info"
	_, respParam, err := l.client.Post(newCtx, tempUrl, nil, body)
	if err != nil {
		l.log.Errorln("user_management GetUserInfos err: %v", err)
		return
	}

	// 解析数据
	temp := respParam.([]interface{})
	for k := range temp {
		tempUser := temp[k].(map[string]interface{})
		tempRoles := make(map[interfaces.Role]bool, 0)
		for _, role := range tempUser["roles"].([]interface{}) {
			tempRoles[l.mapStringRole[role.(string)]] = true
		}
		tempID := tempUser["id"].(string)
		infos[tempID] = interfaces.UserInfo{
			ID:    tempID,
			Name:  tempUser["name"].(string),
			Roles: tempRoles,
		}
	}
	return infos, nil
}
