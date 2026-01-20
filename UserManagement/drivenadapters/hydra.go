package drivenadapters

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/kweaver-ai/go-lib/httpclient"
	jsoniter "github.com/json-iterator/go"

	"UserManagement/common"
)

type hydra struct {
	adminAddress  string
	publicAddress string
	log           common.Logger
	client        *http.Client
}

type registerInfo struct {
	Name           string   `json:"client_name"`
	Secret         string   `json:"client_secret,omitempty"`
	GrantTypes     []string `json:"grant_types"`
	ResponseTypes  []string `json:"response_types"`
	Scope          string   `json:"scope"`
	ClientLifespan string   `json:"client_credentials_grant_access_token_lifespan,omitempty"`
}

var (
	hOnce sync.Once
	h     *hydra
)

// NewHydra 创建授权服务
func NewHydra() *hydra {
	hOnce.Do(func() {
		config := common.SvcConfig
		h = &hydra{
			adminAddress:  fmt.Sprintf("http://%s:%d", config.OAuthAdminHost, config.OAuthAdminPort),
			publicAddress: fmt.Sprintf("http://%s:%d", config.OAuthPublicHost, config.OAuthPublicPort),
			log:           common.NewLogger(),
			client:        httpclient.NewRawHTTPClient(),
		}
	})

	return h
}

// DeleteConsentAndLogin 删除认证与授权会话
func (h *hydra) DeleteConsentAndLogin(clientID, userID string) (err error) {
	if userID == "" {
		return nil
	}
	// 删除授权
	var target string
	if clientID == "" {
		target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/consent?subject=%s&all=true", h.adminAddress, userID)
	} else {
		target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/consent?subject=%s&client=%s", h.adminAddress, userID, clientID)
	}

	code, header, respBody, err := h.doRequest("DELETE", target, nil)
	if err != nil {
		return err
	}

	// 如果不是204或者404，则报错
	if code != http.StatusNoContent && code != http.StatusNotFound {
		err = fmt.Errorf("code:%v,header:%v,body:%v", code, header, string(respBody))
		return
	}

	// 删除认证
	target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/login?subject=%s", h.adminAddress, userID)

	code, header, respBody, err = h.doRequest("DELETE", target, nil)
	if err != nil {
		return err
	}

	// 如果不是204或者404，则报错
	if code != http.StatusNoContent && code != http.StatusNotFound {
		err = fmt.Errorf("code:%v,header:%v,body:%v", code, header, string(respBody))
		return
	}

	return
}

// Register 客户端注册， lifespan 单位为小时
func (h *hydra) Register(name, password string, lifespan int) (id string, err error) {
	target := fmt.Sprintf("%v/admin/clients", h.adminAddress)

	info := registerInfo{
		Name:          name,
		Secret:        password,
		GrantTypes:    []string{"client_credentials"},
		ResponseTypes: []string{"token"},
		Scope:         "all",
	}

	if lifespan > 0 {
		info.ClientLifespan = fmt.Sprintf("%dh", lifespan)
	}

	code, header, respBody, err := h.doRequest("POST", target, info)
	if err != nil {
		return
	}

	// 如果不是201，则报错
	if code != http.StatusCreated {
		return "", fmt.Errorf("code:%v,header:%v,body:%v", code, header, string(respBody))
	}

	respParam := make(map[string]interface{})
	err = jsoniter.Unmarshal(respBody, &respParam)
	if err != nil {
		return "", err
	}

	return respParam["client_id"].(string), nil
}

// Delete 客户端删除
func (h *hydra) Delete(id string) (err error) {
	target := fmt.Sprintf("%v/admin/clients/%v", h.adminAddress, id)

	code, header, respBody, err := h.doRequest("DELETE", target, nil)
	if err != nil {
		return err
	}

	// 如果不是204，则报错
	if code != http.StatusNoContent {
		err = fmt.Errorf("code:%v,header:%v,body:%v", code, header, string(respBody))
		return
	}

	return nil
}

// Update 客户端更新
func (h *hydra) Update(id, name, password string) (err error) {
	target := fmt.Sprintf("%v/admin/clients/%v", h.adminAddress, id)

	infos := make([]map[string]interface{}, 0)
	if name != "" {
		info1 := map[string]interface{}{
			"op":    "replace",
			"path":  "/client_name",
			"value": name,
		}
		infos = append(infos, info1)
	}

	if password != "" {
		info2 := map[string]interface{}{
			"op":    "replace",
			"path":  "/client_secret",
			"value": password,
		}
		infos = append(infos, info2)
	}

	code, header, respBody, err := h.doRequest("PATCH", target, infos)
	if err != nil {
		return err
	}

	// 如果不是200，则报错
	if code != http.StatusOK {
		err = fmt.Errorf("code:%v,header:%v,body:%v", code, header, string(respBody))
		return
	}

	return nil
}

func (h *hydra) GenerateToken(clientID, clientSecret string) (token string, err error) {
	target := fmt.Sprintf("%v/oauth2/token", h.publicAddress)

	// 生成Basic认证头
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req, err := http.NewRequest("POST", target, strings.NewReader("grant_type=client_credentials&scope=all"))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := h.client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			common.NewLogger().Errorln(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("url:%v,code:%v,header:%v,body:%v", target, res.StatusCode, res.Header, string(resBody))
	}

	tokenResp := make(map[string]interface{})
	err = jsoniter.Unmarshal(resBody, &tokenResp)
	return tokenResp["access_token"].(string), err
}

func (h *hydra) DeleteClientToken(clientID string) (err error) {
	target := fmt.Sprintf("%v/admin/oauth2/tokens?client_id=%s", h.adminAddress, clientID)

	code, header, respBody, err := h.doRequest("DELETE", target, nil)
	if err != nil {
		h.log.Errorln("delete client token failed, err:", err)
		return err
	}

	if code != http.StatusNoContent {
		return fmt.Errorf("url:%v,code:%v,header:%v,body:%v", target, code, header, string(respBody))
	}

	return nil
}

func (h *hydra) doRequest(method, target string, reqBody interface{}) (code int, header http.Header, resBody []byte, err error) {
	var buffer io.Reader
	var payload []byte
	if reqBody == nil {
		buffer = nil
	} else {
		payload, err = json.Marshal(reqBody)
		if err != nil {
			return
		}
		buffer = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, target, buffer)
	if err != nil {
		return
	}

	res, err := h.client.Do(req)
	if err != nil {
		return
	}

	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			common.NewLogger().Errorln(closeErr)
		}
	}()

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return code, header, nil, err
	}

	code = res.StatusCode
	header = res.Header
	return
}
