package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"policy_mgnt/decision"
	"policy_mgnt/general"
	"policy_mgnt/test"
	"policy_mgnt/test/mock_descision"
	"policy_mgnt/test/mock_general"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/gocommon/api"
	"policy_mgnt/utils/models"
)

func mockGeneralMgnt(t *testing.T) (*gomock.Controller, *gomock.Controller, *mock_general.MockManagement, *mock_descision.MockPolicyDecision) {
	ctrl := gomock.NewController(t)
	pctrl := gomock.NewController(t)
	mgnt := mock_general.NewMockManagement(ctrl)
	pmgnt := mock_descision.NewMockPolicyDecision(pctrl)

	return ctrl, pctrl, mgnt, pmgnt
}

func (h *generalHandler) MockAddRouters(r *gin.RouterGroup) {
	r.GET("/general", h.getPolicyList)
	r.PUT("/general/:name/value", h.setPolicyByName)
	r.PUT("/general/:name/state", h.setPolicyState)
}

func mockAddGeneralRoute(mgnt general.Management, pmgnt decision.PolicyDecision) *gin.Engine {
	// 关闭oauth
	viper.Set("oauth_on", false)
	router := gin.Default()
	group := router.Group("/api/policy-management/v1")
	h := newGeneralHandlerWithMgnt(mgnt, pmgnt)
	h.MockAddRouters(group)
	return router
}

func assertPolicy(t *testing.T, data interface{}, assertPolicy models.Policy[[]byte], assertValue []byte) {
	policy := data.(map[string]interface{})
	assert.Equal(t, policy["name"], assertPolicy.Name)
	assert.Equal(t, policy["locked"], assertPolicy.Locked)
	value, _ := json.Marshal(policy["value"])
	assert.JSONEq(t, string(value), string(assertValue))
}

// struct 转 map
func TestPolicyToMap(t *testing.T) {
	policyList := []models.Policy[[]byte]{
		models.Policy[[]byte]{
			Name:    "test_policy1",
			Default: []byte(`{"enable":false}`),
			Value:   []byte(`{"enable":true}`),
			Locked:  false,
		},
		models.Policy[[]byte]{
			Name:    "test_policy2",
			Default: []byte(`{"enable":false}`),
			Value:   []byte(`{"enable":true}`),
			Locked:  false,
		},
	}
	resultCurrent := []interface{}{
		map[string]interface{}{
			"name":   "test_policy1",
			"value":  json.RawMessage(`{"enable":true}`),
			"locked": false,
		},
		map[string]interface{}{
			"name":   "test_policy2",
			"value":  json.RawMessage(`{"enable":true}`),
			"locked": false,
		},
	}
	resultDefault := []interface{}{
		map[string]interface{}{
			"name":   "test_policy1",
			"value":  json.RawMessage(`{"enable":false}`),
			"locked": false,
		},
		map[string]interface{}{
			"name":   "test_policy2",
			"value":  json.RawMessage(`{"enable":false}`),
			"locked": false,
		},
	}

	var result []interface{}
	// mode no set
	result = loopPolicyToMap("", policyList)
	assert.Equal(t, result, resultCurrent)

	// current
	result = loopPolicyToMap("current", policyList)
	assert.Equal(t, result, resultCurrent)

	// default
	result = loopPolicyToMap("default", policyList)
	assert.Equal(t, result, resultDefault)
}

// 检查策略是否匹配，不在body报错
func TestPolicyMatchError(t *testing.T) {
	names := []string{"test_policy1", "test_policy2"}
	policies := []policyParam{
		policyParam{
			Name:  names[0],
			Value: []byte(`{"enable":false}`),
		},
		policyParam{
			Name:  "unknown",
			Value: []byte(`{"enable":false}`),
		},
	}
	_, err := policyMatch(names, policies)
	assert.Equal(t, err, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"test_policy2"}}}))
}

// 检查策略是否匹配，不在url忽略
func TestPolicyMatchIgnore(t *testing.T) {
	names := []string{"test_policy1"}
	incorrentReq := []policyParam{
		policyParam{
			Name:  names[0],
			Value: []byte(`{"enable":false}`),
		},
		policyParam{
			Name:  "nomatch",
			Value: []byte(`{"enable":false}`),
		},
	}
	correntReq := []policyParam{
		policyParam{
			Name:  names[0],
			Value: []byte(`{"enable":false}`),
		},
	}
	policies, err := policyMatch(names, incorrentReq)
	// v1, _ := json.Marshal(correntReq)
	// v2, _ := json.Marshal(policies)
	// assert.Equal(t, v1, v2)
	assert.Equal(t, correntReq, policies)
	assert.Nil(t, err)
}

func setUpGetNamesRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/:name", func(c *gin.Context) {
		names := getNames(c)
		if c.IsAborted() {
			return
		}
		c.JSON(200, names)
	})
	return router
}

// 获取所有策略内容
func TestListPolicy(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	mockCount := 1
	policy := models.Policy[[]byte]{
		Name:    "test_policy",
		Default: []byte(`{"enable":false}`),
		Value:   []byte(`{"enable":true}`),
		Locked:  false,
	}

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	mockNames := []string{}
	mgnt.EXPECT().ListPolicy(0, 20, mockNames).Return([]models.Policy[[]byte]{policy}, mockCount, nil)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/policy-management/v1/general", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	data := test.AssertListResponse(t, resp.Body.Bytes(), mockCount)
	assertPolicy(t, data[0], policy, policy.Value)
}

// 获取指定策略内容
func TestListPolicyByName(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	mockCount := 1
	policy := models.Policy[[]byte]{
		Name:    "test_policy",
		Default: []byte(`{"enable":false}`),
		Value:   []byte(`{"enable":true}`),
		Locked:  false,
	}

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	mockNames := []string{"test_policy"}
	mgnt.EXPECT().ListPolicy(0, 20, mockNames).Return([]models.Policy[[]byte]{policy}, mockCount, nil)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/policy-management/v1/general?name=test_policy", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	data := test.AssertListResponse(t, resp.Body.Bytes(), mockCount)
	assertPolicy(t, data[0], policy, policy.Value)
}

// 获取所有策略内容，分页参数错误
func TestListPolicyInvalidPage(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()

	router := mockAddGeneralRoute(mgnt, pmgnt)

	req := httptest.NewRequest("GET", "/api/policy-management/v1/general?offset=a", nil)
	params := []string{"offset"}
	err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
	test.CoverError(t, req, router, err)
}

// 获取所有策略内容，未知模式
func TestListPolicyUnknownMode(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()

	router := mockAddGeneralRoute(mgnt, pmgnt)

	req := httptest.NewRequest("GET", "/api/policy-management/v1/general?mode=a", nil)
	err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"a"}}})
	test.CoverError(t, req, router, err)
}

// 获取所有策略内容，未知错误
func TestListPolicyUnknownError(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	err := fmt.Errorf("unknown")
	mockNames := []string{}
	mgnt.EXPECT().ListPolicy(0, 20, mockNames).Return(nil, 0, err)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	req := httptest.NewRequest("GET", "/api/policy-management/v1/general", nil)
	assertErr := errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: err.Error()})
	test.CoverError(t, req, router, assertErr)
}

// 设置策略成功
func TestPutPolicy(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	names := []string{"network_restriction", "password_strength_meter", "client_restriction"}
	clientResValue := map[string]bool{
		"pc_web":     false,
		"mobile_web": false,
		"windows":    false,
		"mac":        false,
		"android":    false,
		"ios":        false,
		"linux":      false,
	}
	clientResValueByte, _ := json.Marshal(clientResValue)
	params := []map[string]interface{}{
		map[string]interface{}{
			"name": names[0],
			"value": map[string]interface{}{
				"is_enabled": false,
			},
		},
		map[string]interface{}{
			"name": names[1],
			"value": map[string]interface{}{
				"enable": false,
				"length": 11,
			},
		},
		map[string]interface{}{
			"name":  names[2],
			"value": clientResValue,
		},
	}
	policies := map[string][]byte{
		names[0]: []byte(`{"is_enabled":false}`),
		names[1]: []byte(`{"enable":false,"length":11}`),
		names[2]: clientResValueByte,
	}

	mgnt.EXPECT().SetPolicyValue(policies, false).Return(nil)
	pmgnt.EXPECT().PublishInit().Return(nil)
	pmgnt.EXPECT().PublishInit().Return(nil)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	resp := httptest.NewRecorder()
	reqBody, _ := json.Marshal(params)
	req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ",")+"/value", bytes.NewReader(reqBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}

// 设置策略的参数错误
func TestPutInvalidPolicy(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)
	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	router := mockAddGeneralRoute(mgnt, pmgnt)

	var tests = []struct {
		names    []string
		body     []byte
		wantcode int
		params   []string
		cause    string
	}{
		// body为null
		{[]string{"password_strength_meter"}, []byte(``), 400, []string{"request body"}, "Request body is needed."},
		// body不为json
		{[]string{"password_strength_meter"}, []byte(`"is_enabled"`), 400, []string{"(root)"},
			"(root): Invalid type. Expected: array, given: string"},
		// name不为指定字符串
		{[]string{"password_strength_meter"}, []byte(`[{"name": "123", "value": ""}]`), 400, []string{"0", "0.name", "0.value"},
			"0: Must validate at least one schema (anyOf); 0.name: 0.name must be one of the following: \"password_strength_meter\"; 0.value: Invalid type. Expected: object, given: string"},
		// name正确，body字段不符
		{[]string{"password_strength_meter"}, []byte(`[{"name": "network_restriction", "value": {"is_enabled1": true}}]`), 400, []string{"0", "0.value"},
			"0: Must validate at least one schema (anyOf); 0.value: is_enabled is required"},
		// name正确，body字段的值类型不符
		{[]string{"password_strength_meter"}, []byte(`[{"name": "network_restriction", "value": {"is_enabled": 111}}]`), 400, []string{"0", "0.value.is_enabled"},
			"0: Must validate at least one schema (anyOf); 0.value.is_enabled: Invalid type. Expected: boolean, given: integer"},
		// name正确，body中不包含name中的字段
		{[]string{"user_document"}, []byte(`[{"name": "network_restriction", "value": {"is_enabled": true}}]`), 404, []string{"user_document"}, ""},
	}

	for _, test := range tests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(test.names, ",")+"/value", bytes.NewReader(test.body))
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}, Cause: test.cause})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
			if test.wantcode == 404 {
				expectErr := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": test.params}, Cause: test.cause})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}
}

// 设置策略，策略已锁定
func TestPutPolicyLocked(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	names := []string{"network_restriction"}
	params := []map[string]interface{}{
		map[string]interface{}{
			"name": names[0],
			"value": map[string]interface{}{
				"is_enabled": false,
			},
		},
	}
	policies := map[string][]byte{
		names[0]: []byte(`{"is_enabled":false}`),
	}
	err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"policys": names}})
	mgnt.EXPECT().SetPolicyValue(policies, false).Return(err)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	reqBody, _ := json.Marshal(params)
	req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+names[0]+"/value", bytes.NewReader(reqBody))
	test.CoverError(t, req, router, err)
}

// 锁定策略
func TestLockPolicy(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	names := []string{"test_policy1", "test_policy2"}
	states := map[general.State]interface{}{
		general.StateLocked: true,
	}
	mgnt.EXPECT().SetPolicyState(names, states).Return(nil)

	router := mockAddGeneralRoute(mgnt, pmgnt)

	resp := httptest.NewRecorder()
	reqBody, _ := json.Marshal(map[string]bool{"locked": true})
	req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ",")+"/state", bytes.NewReader(reqBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}

// 锁定策略，非法请求体
func TestLockPolicyIncorrectBody(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	names := []string{"test_policy1", "test_policy2"}

	router := mockAddGeneralRoute(mgnt, pmgnt)

	req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ",")+"/state", bytes.NewReader([]byte(`["name"]`)))
	err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{}}})
	test.CoverError(t, req, router, err)

	reqBody, _ := json.Marshal(map[string]string{string(general.StateLocked): "true"})
	req = httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ",")+"/state", bytes.NewReader(reqBody))
	err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"locked"}}})
	test.CoverError(t, req, router, err)
}

// 设置策略状态，未知状态
func TestSetPolicyStateNotFound(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, pctrl, mgnt, pmgnt := mockGeneralMgnt(t)
	defer ctrl.Finish()
	defer pctrl.Finish()
	names := []string{"test_policy1", "test_policy2"}
	params := map[string]interface{}{
		string(general.StateLocked): true,
		"unknown1":                  true,
	}

	router := mockAddGeneralRoute(mgnt, pmgnt)

	reqBody, _ := json.Marshal(params)
	req := httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ";")+"/state", bytes.NewReader(reqBody))
	err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"unknown1"}}})
	test.CoverError(t, req, router, err)

	// 传入空数据
	params1 := []byte("{}")
	req = httptest.NewRequest("PUT", "/api/policy-management/v1/general/"+strings.Join(names, ";")+"/state", bytes.NewReader(params1))
	err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": map[string]interface{}{}}})
	test.CoverError(t, req, router, err)
}
