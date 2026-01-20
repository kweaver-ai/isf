package drivenadapters

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/rest"
	jsoniter "github.com/json-iterator/go"
	"github.com/satori/uuid"

	"Authentication/common"
	"Authentication/interfaces"
)

type hydraPublic struct {
	log        common.Logger
	publicAddr string
	rawClient  *http.Client
}

type hydraAdmin struct {
	log        common.Logger
	adminAddr  string
	rawClient  *http.Client
	rawClient2 *http.Client
}

var (
	aOnce sync.Once
	pOnce sync.Once
	ha    *hydraAdmin
	hp    *hydraPublic
)

// NewHydraAdmin 创建hydra admin接口操作对象
func NewHydraAdmin() *hydraAdmin {
	config := common.SvcConfig
	aOnce.Do(func() {
		ha = &hydraAdmin{
			log:        common.NewLogger(),
			rawClient:  httpclient.NewRawHTTPClient(),
			rawClient2: common.NewRawHTTPClient2(),
			adminAddr:  fmt.Sprintf("http://%s:%d", config.OAuthAdminHost, config.OAuthAdminPort),
		}
	})

	return ha
}

// NewHydraPublic 创建hydra public接口操作对象
func NewHydraPublic() *hydraPublic {
	config := common.SvcConfig
	pOnce.Do(func() {
		hp = &hydraPublic{
			log:        common.NewLogger(),
			rawClient:  httpclient.NewRawHTTPClient(),
			publicAddr: fmt.Sprintf("http://%s:%d", config.OAuthPublicHost, config.OAuthPublicPort),
		}
	})

	return hp
}

func (h *hydraPublic) AuthorizeRequest(reqInfo *interfaces.AuthorizeInfo) (challenge string, context interface{}, err error) {
	var target string
	state := uuid.NewV4().String()
	switch reqInfo.ResponseType {
	case "code":
		target = fmt.Sprintf("%v/oauth2/auth?client_id=%v&redirect_uri=%v&response_type=%v&scope=%v&state=%v",
			h.publicAddr, reqInfo.ClientID, url.QueryEscape(reqInfo.RedirectURI), url.QueryEscape(reqInfo.ResponseType),
			url.QueryEscape(reqInfo.Scope), state)
	case "token":
		target = fmt.Sprintf("%v/oauth2/auth?client_id=%v&redirect_uri=%v&response_type=%v&scope=%v&state=%v",
			h.publicAddr, reqInfo.ClientID, url.QueryEscape(reqInfo.RedirectURI), url.QueryEscape(reqInfo.ResponseType),
			url.QueryEscape(reqInfo.Scope), state)
	case "token id_token":
		target = fmt.Sprintf("%v/oauth2/auth?client_id=%v&redirect_uri=%v&response_type=%v&scope=%v&state=%v&nonce=%v",
			h.publicAddr, reqInfo.ClientID, url.QueryEscape(reqInfo.RedirectURI), url.QueryEscape(reqInfo.ResponseType),
			url.QueryEscape(reqInfo.Scope), state, state)
	}

	location, cookies, err := h.doRequest("GET", target, nil, nil)
	if err != nil {
		h.log.Errorf("Authoriza request failed: %v, url: %v", err, target)
		return "", nil, err
	}

	u, err := url.Parse(location)
	if err != nil {
		return "", nil, err
	}

	// 正确返回
	// https://10.2.176.213:443/oauth2/signin?login_challenge=bbbbd20cd9614bb3bc4bb1ccede2d2fd

	// 报错示例
	// https://10.2.176.204:9010/callback?error=unsupported_response_type&error_description=The+authorization+server+does+not+support
	// +obtaining+a+token+using+this+method&error_hint=The+client+is+not+allowed+to+request+response_type+%22id_token+toke%22.&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4
	// https://10.2.176.204:9010/callback?error=invalid_scope&error_description=The+requested+scope+is+invalid%2C+unknown%2C+or+malformed
	// &error_hint=The+OAuth+2.0+Client+is+not+allowed+to+request+scope+%22offlin%22.&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4
	// https://10.2.176.203:3333/error?error=invalid_client&error_description=Client+authentication+failed+%28e.g.%2C+unknown+client
	// %2C+no+client+authentication+included%2C+or+unsupported+authentication+method%29&error_hint=The+requested+OAuth+2.0+Client+does+not+exist.
	// https://10.2.176.203:3333/error?error=invalid_request&error_description=The+request+is+missing+a+required+parameter%2C+includes+an
	// +invalid+parameter+value%2C+includes+a+parameter+more+than+once%2C+or+is+otherwise+malformed&error_hint=The+%22redirect_uri%22+parameter
	// +does+not+match+any+of+the+OAuth+2.0+Client%27s+pre-registered+redirect+urls.

	// 检查login_challenge存在
	loginChallenge, ok := u.Query()["login_challenge"]
	if !ok {
		return "", nil, rest.NewHTTPError(location, rest.BadRequest, nil)
	}

	return loginChallenge[0], cookies, nil
}

func (h *hydraPublic) convertOuterURLToInner(redirectURL string) string {
	return h.publicAddr + "/oauth2/auth" + strings.Split(redirectURL, "/oauth2/auth")[1]
}

func (h *hydraPublic) VerifierLogin(redirURL string, context interface{}) (challenge string, newContext interface{}, err error) {
	redirURL = h.convertOuterURLToInner(redirURL)
	location, newCookies, err := h.doRequest("GET", redirURL, context.([]*http.Cookie), nil)
	if err != nil {
		h.log.Errorf("Verifier login request failed: %v, url: %v", err, redirURL)
		return "", nil, err
	}

	u, err := url.Parse(location)
	if err != nil {
		return "", nil, err
	}

	// 正确返回
	// https://10.2.176.213:443/oauth2/consent?consent_challenge=bbbbd20cd9614bb3bc4bb1ccede2d2fd

	// 报错示例
	// https://10.2.176.204:9010/callback#error=request_forbidden&error_description=The+request+is+not+allowed&
	// error_hint=You+are+not+allowed+to+perform+this+action.&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4

	// 检查consent_challenge存在
	consentChallenge, ok := u.Query()["consent_challenge"]
	if !ok {
		return "", nil, rest.NewHTTPError(location, rest.BadRequest, nil)
	}

	return consentChallenge[0], newCookies, nil
}

func (h *hydraPublic) VerifierConsent(redirURL, responseType string, context interface{}) (*interfaces.TokenInfo, error) {
	redirURL = h.convertOuterURLToInner(redirURL)
	location, _, err := h.doRequest("GET", redirURL, context.([]*http.Cookie), nil)
	if err != nil {
		h.log.Errorf("Verifier consent failed: %v, url: %v", err, redirURL)
		return nil, err
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	// responseType == "code" 正确返回
	// https://10.2.176.204:9010/callback?code=zeEWYPu6KEhBYLUjUKFiohouhfBDw1qclX1AkB41PuM.RZfgl_hDgJ37pBAwXcWDgoUYhrb7XwwRmJF71Zq1ViI&scope=offline%20openid&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4
	// responseType == "token"  正确返回
	// https://10.2.176.204:9010/callback#access_token=Mj_sxiJg81hYcARjtmQj_4A_W2tUUki0J6TuADfhl4w.Nqvq4r1IqcesRlHRVIsiQoYdy6-0D1k4_-Ea83zxaPQ
	// &expires_in=3599&scope=offline%20openid&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4&token_type=bearer
	// responseType == "token id_token"  正确返回
	// https://10.2.176.204:9010/callback#access_token=e6pF0P2V6fqX_nNMzPhKN_71taoOlxOdVhIsjuuWJQo.mkHJZsPOqZjO4N6l2gGQLE7JthnL3DsmX30RTlaL8vo&expires_in=3600
	// &id_token=eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo0ZWI4NjBlYy01YzQ5LTQ4ZWItOTliYS00MTM3MTg0ODJmYjAiLCJ0eXAiOiJKV1QifQ.
	// eyJhdF9oYXNoIjoicVR0YWZnTkNsYjdkQ2FNdkxCQmt0USIsImF1ZCI6WyJrZWh1Il0sImF1dGhfdGltZSI6MTU5MzUwMjk3MCwiZXhwIjoxNTkzNTA2NTgwLCJpYXQiOj
	// E1OTM1MDI5ODAsImlzcyI6Imh0dHA6Ly8xMC4yLjE4MC4xODE6MzAwMDEvIiwianRpIjoiNWIyY2EzMmQtZGI2Yy00MjJhLWFlNTEtODRkZjdhMjEzNjNhIiwibm9uY2Ui
	// OiJjOWI3YmI3MS1mZTI1LTQ4NjMtYjJhZC0zYjM5Yzc5YTgyYTQiLCJyYXQiOjE1OTM1MDI5NjQsInNpZCI6IjViMWI4NjU2LTExMzMtNDUzYS1iNzcwLTA1ZDdkNjRhNT
	// hkZiIsInN1YiI6IjNjNzFmMjVjLTdkMjEtMTFlYS04ZmRhLTAwNTA2NSJ9.PuF-ZtVZMF8XzIBJ6fUShGIYPHPU0K3K0VhwHXaWsvN-f1IAwGhzGCuMIPD6rDChelaKj2S
	// kHeRVwMrv9VOLg39JwRxsALlITXRsjsj7vWM8oEOAI56AXbnkHCCXNi7rDZg9TlRnl7ZibirNqTT2eZcfj7PdADhNNj_g3meZrSU5Y_MSsdmabfwXOGHVGvptS9HoC7kcE
	// HmEAboi7bb7gH0NnTjXSPVBxvJlP85fjLuSH-GJgrZDaUfzrMygYTgRDzL-MSrMmIjh2z3BCztBpccgtu2QZLq9Sp4vDz3UKi4um9w3MYhJvckEjoVKb0s__6P87uLxrof
	// Kbq_GS-e1NW76mJkDzRff_FIbV9Ztr7XIGCHojJa3AiD1y3P06iVSo24mgEJpAa476_36aS3voxQLCerwijjQr_zUmns1xfpmiy7KDoXWzxJA3FD6x5KLE4fU004FotfZJ
	// 0MwPwVkaFPv7fC0GEjPeB6CJK4BOOiNpgWu6cL2bfqIEuTwVnd_0HMyihlYQ4Jilez9wt4TfZP5kQgbDV3uLtTw7YT2keijaz__8GcC9eX6Zrc5_UEOls6Ws1zrzQofAF7
	// re6vy3GSi42MjiY4IWn9JXqRx8Guqf3vsPPPGDd4jtHswb7LJymFgdfmt8sbp4N4ekbDRGVMdNW92hqpSBCJ6sNLPqZM
	// &scope=offline%20openid&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4&token_type=bearer

	// 报错示例
	// https://10.2.176.204:9010/callback?error=unsupported_response_type&error_description=The+authorization+server+does+not+support+obtaining
	// +a+token+using+this+method&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4
	// https://10.2.176.204:9010/callback#error=invalid_request&error_description=The+request+is+missing+a+required+parameter%252C+includes+an
	// +invalid+parameter+value%252C+includes+a+parameter+more+than+once%252C+or+is+otherwise+malformed&error_hint=Parameter+%2522nonce%2522+must
	// +be+set+when+using+the+OpenID+Connect+Hybrid+Flow.&state=c9b7bb71-fe25-4863-b2ad-3b39c79a82a4

	// 获取返回参数
	tInfo := &interfaces.TokenInfo{}
	switch responseType {
	case "code":
		code := u.Query()["code"][0]
		scope := u.Query()["scope"][0]
		if (code != "") && (scope != "") {
			tInfo.Code = code
			tInfo.Scope = scope
			tInfo.ResponseType = "code"
		} else {
			return nil, rest.NewHTTPError(location, rest.BadRequest, nil)
		}
	case "token":
		f, err := url.ParseQuery(u.Fragment)
		if err != nil {
			return nil, err
		}
		AccessToken := f.Get("access_token")
		ExpirsesIn := f.Get("expires_in")
		Scope := f.Get("scope")
		TokenType := f.Get("token_type")

		if (AccessToken == "") || (ExpirsesIn == "") || (Scope == "") || (TokenType == "") {
			return nil, rest.NewHTTPError(location, rest.BadRequest, nil)
		}

		tInfo.AccessToken = AccessToken
		tInfo.ExpirsesIn, err = strconv.ParseInt(ExpirsesIn, 10, 64)
		if err != nil {
			return nil, err
		}
		tInfo.Scope = Scope
		tInfo.TokenType = TokenType
	case "token id_token":
		f, err := url.ParseQuery(u.Fragment)
		if err != nil {
			return nil, err
		}

		AccessToken := f.Get("access_token")
		ExpirsesIn := f.Get("expires_in")
		IDToken := f.Get("id_token")
		Scope := f.Get("scope")
		TokenType := f.Get("token_type")

		if (AccessToken == "") || (ExpirsesIn == "") || (IDToken == "") || (Scope == "") || (TokenType == "") {
			return nil, rest.NewHTTPError(location, rest.BadRequest, nil)
		}

		tInfo.AccessToken = AccessToken
		tInfo.ExpirsesIn, err = strconv.ParseInt(ExpirsesIn, 10, 64)
		if err != nil {
			return nil, err
		}
		tInfo.IDToken = IDToken
		tInfo.Scope = Scope
		tInfo.TokenType = TokenType
		tInfo.ResponseType = "token id_token"
	}

	return tInfo, nil
}

func (h *hydraPublic) GetTokenEndpoint() (tokenEndpoint string, err error) {
	target := fmt.Sprintf("%s/.well-known/openid-configuration", h.publicAddr)
	res, err := h.rawClient.Get(target)
	if err != nil {
		h.log.Errorln(target, err)
		return
	}

	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			h.log.Errorln(closeErr)
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("code:%v, body:%v", res.StatusCode, string(resBody))
		return
	}

	resp := make(map[string]interface{})
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return
	}

	tokenEndpoint = resp["token_endpoint"].(string)
	return
}

func (h *hydraPublic) AssertionForToken(clientID, clientSecret, assertion string) (*interfaces.TokenInfo, error) {
	reqParam := fmt.Sprintf("scope=%s&grant_type=%v&assertion=%v", "all", "urn:ietf:params:oauth:grant-type:jwt-bearer", assertion)

	return h.getAccessToken(clientID, clientSecret, reqParam)
}

func (h *hydraPublic) getAccessToken(clientID, clientSecret, reqParam string) (*interfaces.TokenInfo, error) {
	auth := clientID + ":" + clientSecret
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	reqHeader := make(map[string]string)
	reqHeader["Authorization"] = basicAuth
	reqHeader["Content-Type"] = "application/x-www-form-urlencoded"
	target := fmt.Sprintf("%s/oauth2/token", h.publicAddr)

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader([]byte(reqParam)))
	if err != nil {
		return nil, err
	}
	h.addHeaders(req, reqHeader)

	res, err := h.rawClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			h.log.Errorln(closeErr)
		}
	}()
	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, rest.NewHTTPError(string(resBodyByte), rest.Unauthorized, nil)
		}
		err = fmt.Errorf("code:%v,header:%v,body:%v", res.StatusCode, res.Header, string(resBodyByte))
		return nil, err
	}
	var jsonBody interface{}
	err = json.Unmarshal(resBodyByte, &jsonBody)
	if err != nil {
		return nil, err
	}

	token := &interfaces.TokenInfo{}
	token.AccessToken = jsonBody.(map[string]interface{})["access_token"].(string)
	token.ExpirsesIn = int64(jsonBody.(map[string]interface{})["expires_in"].(float64))
	token.Scope = jsonBody.(map[string]interface{})["scope"].(string)
	token.TokenType = jsonBody.(map[string]interface{})["token_type"].(string)

	return token, nil
}

func (h *hydraPublic) addHeaders(req *http.Request, reqHeader map[string]string) {
	for k, v := range reqHeader {
		if v != "" {
			req.Header.Set(k, v)
		}
	}
}

func (h *hydraPublic) doRequest(method, target string, reqCookies []*http.Cookie, reqBody interface{}) (location string, cookies []*http.Cookie, err error) {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, target, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	for _, cookie := range reqCookies {
		req.AddCookie(cookie)
	}

	res, err := h.rawClient.Do(req)
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
		return
	}

	switch {
	case res.StatusCode == http.StatusFound:
	case res.StatusCode == http.StatusSeeOther:
	case res.StatusCode >= http.StatusBadRequest && res.StatusCode < http.StatusInternalServerError:
		cause := fmt.Sprintf("code:%v,header:%v,body:%v", res.StatusCode, res.Header, string(resBody))
		err = rest.NewHTTPError(cause, rest.BadRequest, nil)
		return
	default:
		err = fmt.Errorf("code:%v,header:%v,body:%v", res.StatusCode, res.Header, string(resBody))
		return
	}

	cookies = res.Cookies()
	location = res.Header.Get("Location")
	return
}

func (h *hydraAdmin) GetLoginRequestInformation(loginChallenge string) (*interfaces.DeviceInfo, error) {
	target := fmt.Sprintf("%v/admin/oauth2/auth/requests/login?login_challenge=%v", h.adminAddr, loginChallenge)

	body, err := h.doRequest(http.MethodGet, target, nil, []int{http.StatusOK})
	if err != nil {
		h.log.Errorf("Get login request information failed: %v, url: %v", err, target)
		return nil, err
	}

	description := ""
	name := ""
	device := body.(map[string]interface{})["client"].(map[string]interface{})["metadata"].(map[string]interface{})["device"].(map[string]interface{})
	if device["description"] != nil {
		description = device["description"].(string)
	}
	if device["name"] != nil {
		name = device["name"].(string)
	}
	return &interfaces.DeviceInfo{
		ClientType:  device["client_type"].(string),
		Description: description,
		Name:        name,
	}, nil
}

func (h *hydraAdmin) AcceptLoginRequest(subject, loginChallenge string) (string, error) {
	permInfo := map[string]interface{}{
		"subject": subject,
	}

	target := fmt.Sprintf("%v/admin/oauth2/auth/requests/login/accept?login_challenge=%v", h.adminAddr, loginChallenge)

	body, err := h.doRequest(http.MethodPut, target, permInfo, []int{http.StatusOK})
	if err != nil {
		h.log.Errorf("Accept login request failed: %v, url: %v", err, target)
		return "", err
	}

	return body.(map[string]interface{})["redirect_to"].(string), nil
}

func (h *hydraAdmin) AcceptConsentRequest(scope, consentChallenge string, context interface{}) (string, error) {
	scopeArr := strings.Split(scope, " ")
	permInfo := map[string]interface{}{
		"grant_scope": scopeArr,
		"session": map[string]interface{}{
			"access_token": context,
		},
	}

	target := fmt.Sprintf("%v/admin/oauth2/auth/requests/consent/accept?consent_challenge=%v", h.adminAddr, consentChallenge)

	body, err := h.doRequest(http.MethodPut, target, permInfo, []int{http.StatusOK})
	if err != nil {
		h.log.Errorf("Accept consent request failed: %v, url: %v", err, target)
		return "", err
	}

	return body.(map[string]interface{})["redirect_to"].(string), nil
}

func (h *hydraAdmin) PublicRegister(client *interfaces.RegisterInfo) (*interfaces.ClientInfo, error) {
	target := fmt.Sprintf("%v/admin/clients", h.adminAddr)

	payload, err := json.Marshal(*client)
	if err != nil {
		return nil, err
	}

	resp, err := h.rawClient.Post(target, "application/json", bytes.NewReader(payload))
	if err != nil {
		h.log.Errorln(err)
		return nil, err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			common.NewLogger().Errorln(closeErr)
		}
	}()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, interfaces.HydraError{
			Body:   resBody,
			Status: resp.StatusCode,
		}
	}

	respParam := make(map[string]interface{})
	err = jsoniter.Unmarshal(resBody, &respParam)
	if err != nil {
		return nil, err
	}

	clientInfo := &interfaces.ClientInfo{
		ClientID:     respParam["client_id"].(string),
		ClientSecret: respParam["client_secret"].(string),
	}

	return clientInfo, nil
}

func (h *hydraAdmin) IntrospectRefreshToken(refreshToken string) (info *interfaces.RefreshTokenIntrospectInfo, err error) {
	target := fmt.Sprintf("%v/admin/oauth2/introspect", h.adminAddr)

	resp, err := h.rawClient.Post(target, "application/x-www-form-urlencoded", strings.NewReader("token="+refreshToken))
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			h.log.Errorln(closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		return nil, errors.New(string(body))
	}

	respParam := make(map[string]interface{})
	err = jsoniter.Unmarshal(body, &respParam)
	if err != nil {
		return nil, err
	}

	// 令牌状态
	if !respParam["active"].(bool) {
		return nil, rest.NewHTTPErrorV2(rest.Unauthorized, "token is not active")
	}

	// 必须是刷新令牌
	if respParam["token_use"].(string) != "refresh_token" {
		return nil, rest.NewHTTPErrorV2(rest.BadRequest, "Invalid token use")
	}

	info = &interfaces.RefreshTokenIntrospectInfo{}
	info.ClientID = respParam["client_id"].(string)
	info.Sub = respParam["sub"].(string)

	// 访问者ID
	return info, nil
}

func (h *hydraAdmin) GetClientInfo(clientID string) (info *interfaces.ClientInfo, err error) {
	target := fmt.Sprintf("%v/admin/clients/%v", h.adminAddr, clientID)

	resBody, err := h.doRequest(http.MethodGet, target, nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}

	info = &interfaces.ClientInfo{}
	info.ClientID = resBody.(map[string]interface{})["client_id"].(string)

	return info, nil
}

func (h *hydraAdmin) doRequest(method, target string, reqBody interface{}, expectCodes []int) (resBody interface{}, err error) {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, target, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	res, err := h.rawClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		closeErr := res.Body.Close()
		if closeErr != nil {
			common.NewLogger().Errorln(closeErr)
		}
	}()

	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	flag := false
	for _, expectCode := range expectCodes {
		if expectCode == res.StatusCode {
			flag = true
			break
		}
	}
	if !flag {
		err = fmt.Errorf("code:%v,header:%v,body:%v", res.StatusCode, res.Header, string(resBodyByte))
		return nil, err
	}

	err = jsoniter.Unmarshal(resBodyByte, &resBody)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (h *hydraAdmin) CreateTrustedPair(publicKey *rsa.PublicKey, issuer, kid string) (err error) {
	nBytes := publicKey.N.Bytes()
	const eBytesLen = 4
	// 根据当前的Hydra版本测试，JWT Grant过期后依旧可以使用。
	// JWT Grant的过期时间，底层数据类型为 timestamp，会有2038问题。
	// 暂时将JWT Grant有效期设为10年，待后续优化。
	const years = 10
	eBytes := make([]byte, eBytesLen)
	binary.BigEndian.PutUint32(eBytes, uint32(publicKey.E))
	nBase64 := strings.TrimRight(base64.URLEncoding.EncodeToString(nBytes), "=")
	eBase64 := strings.TrimRight(base64.URLEncoding.EncodeToString(eBytes), "=")

	reqParam := make(map[string]interface{})
	reqParam["allow_any_subject"] = true
	reqParam["expires_at"] = time.Now().AddDate(years, 0, 0)
	reqParam["issuer"] = issuer
	reqParam["jwk"] = map[string]interface{}{
		"use": "sig",
		"kty": "RSA",
		"kid": kid,
		"alg": "RS256",
		"n":   nBase64,
		"e":   eBase64,
	}
	reqParam["scope"] = []string{"all"}

	target := fmt.Sprintf("%s/admin/trust/grants/jwt-bearer/issuers", h.adminAddr)
	// 创建受信任关系，返回状态码解析：
	// http.StatusCreated：受信任关系创建成功。
	// http.StatusConflict：受信任关系之前就已经创建好了。
	_, err = h.doRequest(http.MethodPost, target, reqParam, []int{http.StatusCreated, http.StatusConflict})
	if err != nil {
		h.log.Errorln(target, err)
		return
	}

	return
}

func (h *hydraAdmin) GetKidTrustedPairByIssuer(issuer string) (trustedPair map[string]bool, err error) {
	trustedPair = make(map[string]bool)
	target := fmt.Sprintf("%s/admin/trust/grants/jwt-bearer/issuers?issuer=%s", h.adminAddr, issuer)
	body, err := h.doRequest(http.MethodGet, target, nil, []int{http.StatusOK})
	if err != nil {
		h.log.Errorln(target, err)
		return
	}

	resp := body.([]interface{})
	for i := range resp {
		pair := resp[i].(map[string]interface{})
		publicKey, pok := pair["public_key"].(map[string]interface{})
		if pok {
			kid, ok := publicKey["kid"].(string)
			if ok {
				trustedPair[kid] = true
			}
		}
	}
	return
}

func (h *hydraAdmin) SetAppAsUserAgent(clientID string) (err error) {
	target := fmt.Sprintf("%s/admin/clients/%s", h.adminAddr, clientID)
	patchBody := make([]map[string]interface{}, 0)
	patchBody = append(patchBody,
		map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata",
			"value": map[string]interface{}{"device": map[string]interface{}{"client_type": "app"}},
		},
		map[string]interface{}{
			"op":    "replace",
			"path":  "/grant_types",
			"value": []string{"client_credentials", "urn:ietf:params:oauth:grant-type:jwt-bearer"},
		},
	)

	_, err = h.doRequest(http.MethodPatch, target, patchBody, []int{http.StatusOK})
	if err != nil {
		h.log.Errorln(target, err)
	}
	return
}

func (h *hydraAdmin) DeleteSession(userID, clientID string) error {
	if userID == "" {
		return nil
	}

	// 删除授权会话
	var target string
	if clientID == "" {
		target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/consent?subject=%s&all=true", h.adminAddr, userID)
	} else {
		target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/consent?subject=%s&client=%s", h.adminAddr, userID, clientID)
	}
	if err := h.doRequestWithoutTimeout(http.MethodDelete, target, nil); err != nil {
		return err
	}

	// 删除认证会话
	target = fmt.Sprintf("%v/admin/oauth2/auth/sessions/login?subject=%s", h.adminAddr, userID)
	if err := h.doRequestWithoutTimeout(http.MethodDelete, target, nil); err != nil {
		return err
	}

	return nil
}

func (h *hydraAdmin) doRequestWithoutTimeout(method, target string, reqBody interface{}) error {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, target, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	res, err := h.rawClient2.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if (res.StatusCode < http.StatusOK) || (res.StatusCode >= http.StatusMultipleChoices) {
		return fmt.Errorf("code:%v,header:%v,target:%v", res.StatusCode, res.Header, target)
	}

	return nil
}
