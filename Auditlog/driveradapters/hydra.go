package driveradapters

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/httpclient"
	jsoniter "github.com/json-iterator/go"

	"AuditLog/common"
	"AuditLog/interfaces"
)

type hydra struct {
	adminAddress   string
	log            common.Logger
	client         *http.Client
	visitorTypeMap map[string]interfaces.VisitorType
	accountTypeMap map[string]interfaces.AccountType
	clientTypeMap  map[string]interfaces.ClientType
}

var (
	hOnce sync.Once
	h     *hydra
)

// newHydra 创建授权服务
func newHydra() *hydra {
	hOnce.Do(func() {
		config := common.SvcConfig
		visitorTypeMap := map[string]interfaces.VisitorType{
			"realname":  interfaces.RealName,
			"anonymous": interfaces.Anonymous,
			"business":  interfaces.App,
		}
		accountTypeMap := map[string]interfaces.AccountType{
			"other":   interfaces.Other,
			"id_card": interfaces.IDCard,
		}
		clientTypeMap := map[string]interfaces.ClientType{
			"unknown":       interfaces.Unknown,
			"ios":           interfaces.IOS,
			"android":       interfaces.Android,
			"windows_phone": interfaces.WindowsPhone,
			"windows":       interfaces.Windows,
			"mac_os":        interfaces.MacOS,
			"web":           interfaces.Web,
			"mobile_web":    interfaces.MobileWeb,
			"nas":           interfaces.Nas,
			"console_web":   interfaces.ConsoleWeb,
			"deploy_web":    interfaces.DeployWeb,
			"linux":         interfaces.Linux,
			"app":           interfaces.APP,
		}
		h = &hydra{
			adminAddress:   fmt.Sprintf("http://%s:%s", config.OAuthAdminHost, config.OAuthAdminPort),
			log:            common.NewLogger(),
			client:         httpclient.NewRawHTTPClient(),
			visitorTypeMap: visitorTypeMap,
			accountTypeMap: accountTypeMap,
			clientTypeMap:  clientTypeMap,
		}
	})

	return h
}

// Introspect token内省
func (h *hydra) Introspect(token string) (info interfaces.TokenIntrospectInfo, err error) {
	target := fmt.Sprintf("%v/admin/oauth2/introspect", h.adminAddress)
	resp, err := h.client.Post(target, "application/x-www-form-urlencoded",
		bytes.NewReader([]byte(fmt.Sprintf("token=%v", token))))
	if err != nil {
		h.log.Errorln(err)
		return
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			common.NewLogger().Errorln(closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		err = errors.New(string(body))
		return
	}

	respParam := make(map[string]interface{})
	err = jsoniter.Unmarshal(body, &respParam)
	if err != nil {
		return
	}

	// 令牌状态
	info.Active = respParam["active"].(bool)
	if !info.Active {
		return
	}

	// 访问者ID
	info.VisitorID = respParam["sub"].(string)
	// Scope 权限范围
	info.Scope = respParam["scope"].(string)
	// 客户端ID
	info.ClientID = respParam["client_id"].(string)
	// 客户端凭据模式
	if info.VisitorID == info.ClientID {
		info.VisitorTyp = interfaces.App
		return
	}
	// 以下字段 只在非客户端凭据模式时才存在
	// 访问者类型
	info.VisitorTyp = h.visitorTypeMap[respParam["ext"].(map[string]interface{})["visitor_type"].(string)]

	// 匿名用户
	if info.VisitorTyp == interfaces.Anonymous {
		return
	}

	// 实名用户
	if info.VisitorTyp == interfaces.RealName {
		// 登陆IP
		info.LoginIP = respParam["ext"].(map[string]interface{})["login_ip"].(string)
		// 设备ID
		info.Udid = respParam["ext"].(map[string]interface{})["udid"].(string)
		// 登录账号类型
		info.AccountTyp = h.accountTypeMap[respParam["ext"].(map[string]interface{})["account_type"].(string)]
		// 设备类型
		info.ClientTyp = h.clientTypeMap[respParam["ext"].(map[string]interface{})["client_type"].(string)]
		return
	}

	return
}
