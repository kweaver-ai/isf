package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"policy_mgnt/decision"
	"policy_mgnt/general"
	"policy_mgnt/network"
	"policy_mgnt/test"
	"policy_mgnt/test/mock_descision"
	"policy_mgnt/test/mock_general"
	"policy_mgnt/test/mock_network"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func mockNetworkMgnt(t *testing.T) (*gomock.Controller, *gomock.Controller, *gomock.Controller, *mock_network.MockManagement, *mock_general.MockManagement, *mock_descision.MockPolicyDecision) {
	ctrl := gomock.NewController(t)
	gctrl := gomock.NewController(t)
	pctrl := gomock.NewController(t)
	mgnt := mock_network.NewMockManagement(ctrl)
	gmgnt := mock_general.NewMockManagement(gctrl)
	pmgnt := mock_descision.NewMockPolicyDecision(pctrl)
	return ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt
}

func (n *networkHandler) MockAddRouters(r *gin.RouterGroup) {
	r.GET("/user-login/network-restriction/network", n.searchNetwork)
	r.POST("/user-login/network-restriction/network", n.addNetwork)
	r.GET("/user-login/network-restriction/network/:id", n.getNetworkByID)
	r.PUT("/user-login/network-restriction/network/:id", n.editNetwork)
	r.DELETE("/user-login/network-restriction/network/:id", n.deleteNetwork)
	r.GET("/user-login/network-restriction/network/:id/accessor", n.searchAccessors)
	r.GET("/user-login/network-restriction/accessor/:id/network", n.getNetworksByAccessorID)
	r.POST("/user-login/network-restriction/network/:id/accessor", n.addAccessors)
	r.DELETE("/user-login/network-restriction/network/:id/accessor/:accessor_id", n.deleteAccessors)
	// 策略引擎暂时没有接入oauth，省略token检查
	//r.GET("/policy-data/bundle.tar.gz", n.getBundle)
}

func mockAddNetworkRoute(mgnt network.Management, gmgnt general.Management, pmgnt decision.PolicyDecision) *gin.Engine {
	// 关闭oauth
	viper.Set("oauth_on", false)
	router := gin.Default()
	group := router.Group("/api/policy-mgnt/v1")
	h := newNetworkHandlerWithMgnt(mgnt, gmgnt, pmgnt)
	h.MockAddRouters(group)
	return router
}

// 搜索网段
func TestSearchNetwork(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	mockData := []models.WebNetworkRestriction{models.WebNetworkRestriction{}}
	mockCount := 1

	var tests = []struct {
		input     []interface{}
		realinput []interface{}
		wantcode  int
		params    []string
	}{
		// ok
		{[]interface{}{"hehe", 1, 2}, []interface{}{"hehe", 1, 2}, 200, []string{}},
		// limit > 1000
		{[]interface{}{"hehe", 1, 2000}, []interface{}{"hehe", 1, 200}, 400, []string{"limit"}},
		// offset < 0
		{[]interface{}{"hehe", -200, 20}, []interface{}{"hehe", 0, 20}, 400, []string{"offset"}},
		// limit < 0
		{[]interface{}{"hehe", 0, -1}, []interface{}{"hehe", 0, -1}, 400, []string{"limit"}},
		// offset type error
		{[]interface{}{"hehe", "a", -1}, []interface{}{1, 0, -1}, 400, []string{"offset"}},
		// limit type error
		{[]interface{}{"hehe", 1, "a"}, []interface{}{1, 0, -1}, 400, []string{"limit"}},
	}

	for _, test := range tests {
		var inputKey, inputStart, inputLimit interface{}
		if len(test.input) == 3 {
			inputKey = test.input[0]
			inputStart = test.input[1]
			inputLimit = test.input[2]
		}
		realinputKey := test.realinput[0]
		realinputStart := test.realinput[1]
		realinputLimit := test.realinput[2]
		// 状态码不为200，不检查调用层
		if test.wantcode == 200 {
			mgnt.EXPECT().SearchNetwork(realinputKey, realinputStart, realinputLimit).Return(mockData, mockCount, nil)
		}
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network?key_word=%v&offset=%v&limit=%v",
			inputKey, inputStart, inputLimit), nil)
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}

	// case: search without query params
	mgnt.EXPECT().SearchNetwork("", 0, 20).Return(mockData, mockCount, nil)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/policy-mgnt/v1/user-login/network-restriction/network", nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

// 添加网段
func TestAddNetwork(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	// case: 允许添加
	policy := models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

	var params models.NetworkRestriction
	mockID := "network_01"
	mgnt.EXPECT().AddNetwork(&params).Return(mockID, nil)

	var mockNet models.NetworkRestriction
	mockBody, _ := json.Marshal(mockNet)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/policy-mgnt/v1/user-login/network-restriction/network", bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// case: 不允许添加
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":false}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/policy-mgnt/v1/user-login/network-restriction/network", bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)

	// case: 传入body字段错误
	var tests = []struct {
		input    []byte
		wantcode int
		params   []string
	}{
		{[]byte(`{"name": 1, "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"name"}},
		{[]byte(`{"name": "", "start_ip": 1, "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"start_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": 1, "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"end_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": 1, "netmask": "", "net_type": ""}`), 400, []string{"ip_address"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": 1, "net_type": ""}`), 400, []string{"netmask"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": 1}`), 400, []string{"net_type"}},
	}

	for _, test := range tests {
		policy = models.Policy[[]byte]{
			Name:    "network_restriction",
			Default: []byte(`{"is_enabled":false}`),
			Value:   []byte(`{"is_enabled":true}`),
			Locked:  false,
		}
		// 允许添加
		gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

		resp = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/api/policy-mgnt/v1/user-login/network-restriction/network", bytes.NewReader(test.input))
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}
}

// 根据id获取网段
func TestGetNetwork(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	mockID := "network1"
	mockRes := models.WebNetworkRestriction{}
	mgnt.EXPECT().GetNetworkByID(mockID).Return(mockRes, nil)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

// 根据访问者id获取网段
func TestGetNetworksByAccessorID(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	var tests = []struct {
		input     []interface{}
		realinput []interface{}
		wantcode  int
		params    []string
	}{
		// ok
		{[]interface{}{"accessorid1", 1, 2}, []interface{}{"accessorid1", 1, 2}, 200, []string{}},
		// limit > 1000
		{[]interface{}{"accessorid1", 1, 2000}, []interface{}{"accessorid1", 1, 200}, 400, []string{"limit"}},
		// offset < 0
		{[]interface{}{"accessorid1", -200, 20}, []interface{}{"accessorid1", 0, 20}, 400, []string{"offset"}},
		//  limit < 0
		{[]interface{}{"accessorid1", 0, -1}, []interface{}{"accessorid1", 0, -1}, 400, []string{"limit"}},
		// offset type error
		{[]interface{}{"accessorid1", "a", -1}, []interface{}{1, 0, -1}, 400, []string{"offset"}},
		// limit type error
		{[]interface{}{"accessorid1", 0, "a"}, []interface{}{1, 0, -1}, 400, []string{"limit"}},
	}

	// mock data
	mockData := []models.WebNetworkRestriction{models.WebNetworkRestriction{}}
	mockCount := 2
	for _, test := range tests {
		var inputID, inputStart, inputLimit interface{}
		if len(test.input) == 3 {
			inputID = test.input[0]
			inputStart = test.input[1]
			inputLimit = test.input[2]
		}
		realinputID := test.realinput[0]
		realinputStart := test.realinput[1]
		realinputLimit := test.realinput[2]
		// 状态码不为200，不检查调用层
		if test.wantcode == 200 {
			mgnt.EXPECT().GetNetworksByAccessorID(realinputID, realinputStart, realinputLimit).Return(mockData, mockCount, nil)
		}
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET",
			fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/accessor/%v/network?offset=%v&limit=%v",
				inputID, inputStart, inputLimit), nil)
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}
}

// 编辑网段
func TestEditNetwork(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	// case: 允许编辑
	policy := models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

	var params models.NetworkRestriction
	mockID := "network_01"
	mgnt.EXPECT().EditNetwork(mockID, &params).Return(nil)
	pmgnt.EXPECT().PublishInit().Return(nil)

	var mockNet models.NetworkRestriction
	mockBody, _ := json.Marshal(mockNet)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// case: 不允许编辑
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":false}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)

	// case: 传入body字段错误
	var tests = []struct {
		input    []byte
		wantcode int
		params   []string
	}{
		{[]byte(`123123`), 400, []string{""}},
		{[]byte(`{"name": 1, "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"name"}},
		{[]byte(`{"name": "", "start_ip": 1, "end_ip": "", "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"start_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": 1, "ip_address": "", "netmask": "", "net_type": ""}`), 400, []string{"end_ip"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": 1, "netmask": "", "net_type": ""}`), 400, []string{"ip_address"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": 1, "net_type": ""}`), 400, []string{"netmask"}},
		{[]byte(`{"name": "", "start_ip": "", "end_ip": "", "ip_address": "", "netmask": "", "net_type": 1}`), 400, []string{"net_type"}},
	}

	for _, test := range tests {
		policy = models.Policy[[]byte]{
			Name:    "network_restriction",
			Default: []byte(`{"is_enabled":false}`),
			Value:   []byte(`{"is_enabled":true}`),
			Locked:  false,
		}
		// 允许
		gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

		resp = httptest.NewRecorder()
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), bytes.NewReader(test.input))
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}
}

// 删除网段
func TestDeleteNetwork(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	// case: 允许
	policy := models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

	mockID := "network_01"
	mgnt.EXPECT().DeleteNetwork(mockID).Return(nil)
	pmgnt.EXPECT().PublishInit().Return(nil)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// case: 不允许
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":false}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s", mockID), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}

// 查询访问者
func TestSearchAccessors(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	var mockData []models.AccessorInfo
	mockID := "network_01"
	mockCount := 1

	var tests = []struct {
		input     []interface{}
		realinput []interface{}
		wantcode  int
		params    []string
	}{
		// ok
		{[]interface{}{"hehe", 1, 2}, []interface{}{"hehe", 1, 2}, 200, []string{}},
		// limit > 1000
		{[]interface{}{"hehe", 1, 2000}, []interface{}{"hehe", 1, 200}, 400, []string{"limit"}},
		// offset < 0
		{[]interface{}{"hehe", -200, 20}, []interface{}{"hehe", 0, 20}, 400, []string{"offset"}},
		// limit < 0
		{[]interface{}{"hehe", 0, -1}, []interface{}{"hehe", 0, -1}, 400, []string{"limit"}},
		// offset type error
		{[]interface{}{"hehe", "a", -1}, []interface{}{1, 0, -1}, 400, []string{"offset"}},
		// limit type error
		{[]interface{}{"hehe", 0, "a"}, []interface{}{1, 0, -1}, 400, []string{"limit"}},
	}

	for _, test := range tests {
		var inputKey, inputStart, inputLimit interface{}
		if len(test.input) == 3 {
			inputKey = test.input[0]
			inputStart = test.input[1]
			inputLimit = test.input[2]
		}
		realinputKey := test.realinput[0]
		realinputStart := test.realinput[1]
		realinputLimit := test.realinput[2]
		// 状态码不为200，不检查调用层
		if test.wantcode == 200 {
			mgnt.EXPECT().SearchAccessors(mockID, realinputKey, realinputStart, realinputLimit).Return(mockData, mockCount, nil)
		}
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET",
			fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%v/accessor?key_word=%v&offset=%v&limit=%v",
				mockID, inputKey, inputStart, inputLimit), nil)
		router.ServeHTTP(resp, req)
		assert.Equal(t, test.wantcode, resp.Code)
		// 状态码不为200，判断detail
		if resp.Code != 200 {
			if test.wantcode == 400 {
				expectErr := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}})
				value, _ := json.Marshal(expectErr)
				assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
			}
		}
	}

	// case: search without query params
	mgnt.EXPECT().SearchAccessors(mockID, "", 0, 20).Return(mockData, mockCount, nil)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET",
		fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%v/accessor", mockID), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

// 添加访问者
func TestAddAccessors(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	// case: 允许添加
	policy := models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

	params := []models.AccessorInfo{
		models.AccessorInfo{
			AccessorId:   "user_id",
			AccessorType: "user",
		},
	}
	mockID := "network_01"
	var mockRes []*api.MultiStatus
	mgnt.EXPECT().AddAccessors(mockID, params).Return(mockRes)

	mockBody := []byte(`[{"accessor_id": "user_id", "accessor_type": "user"}]`)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%v/accessor", mockID),
		bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMultiStatus, resp.Code)

	// case: 不允许添加
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":false}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%v/accessor", mockID),
		bytes.NewReader(mockBody))
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMultiStatus, resp.Code)
	wantbody, _ := json.Marshal([]*api.MultiStatus{api.MultiStatusObject("", nil,
		errors.ErrNoPermission(&api.ErrorInfo{Cause: "Network restriction function is not enabled."}))})
	assert.JSONEq(t, string(wantbody), string(resp.Body.Bytes()))

	// case: 传入body字段错误
	var tests = []struct {
		input    []byte
		wantcode int
		params   []string
		cause    string
	}{
		// body为null
		{[]byte(``), 400, []string{"request body"}, "Request body is needed."},
		// body不为json
		{[]byte(`123123`), 400, []string{"(root)"}, "(root): Invalid type. Expected: array, given: integer"},
		// id类型错误
		{[]byte(`[{"accessor_id": 1, "accessor_type": "user"}]`), 400, []string{"0.accessor_id"},
			"0.accessor_id: Invalid type. Expected: string, given: integer"},
		// id为null
		{[]byte(`[{"accessor_id": null, "accessor_type": "user"}]`), 400, []string{"0.accessor_id"},
			"0.accessor_id: Invalid type. Expected: string, given: null"},
		// type类型错误
		{[]byte(`[{"accessor_id": "user1", "accessor_type": 111}]`), 400, []string{"0.accessor_type"},
			"0.accessor_type: Invalid type. Expected: string, given: integer"},
		// type为null
		{[]byte(`[{"accessor_id": "user1", "accessor_type": null}]`), 400, []string{"0.accessor_type"},
			"0.accessor_type: Invalid type. Expected: string, given: null"},
		// type不为user或department
		{[]byte(`[{"accessor_id": "user1", "accessor_type": "xxx"}]`), 400, []string{"0.accessor_type"},
			"0.accessor_type: 0.accessor_type must be one of the following: \"user\", \"department\""},
		// 第二项错误
		{[]byte(`[{"accessor_id": "", "accessor_type": "user"}, {"accessor_id": null, "accessor_type": "user"}]`), 400,
			[]string{"1.accessor_id"}, "1.accessor_id: Invalid type. Expected: string, given: null"},
		// 网段不存在
		{[]byte(`[{"accessor_id": "user1", "accessor_type": "user"}]`), 404, []string{"id"}, ""},
	}

	for _, test := range tests {
		policy = models.Policy[[]byte]{
			Name:    "network_restriction",
			Default: []byte(`{"is_enabled":false}`),
			Value:   []byte(`{"is_enabled":true}`),
			Locked:  false,
		}
		// 允许添加
		gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
		// 根据期望错误码来mock返回值
		if test.wantcode == 404 {
			var params []models.AccessorInfo
			json.Unmarshal(test.input, &params)
			mst := api.MultiStatusObject("", nil, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}}))
			mockRes := []*api.MultiStatus{mst}
			mgnt.EXPECT().AddAccessors(mockID, params).Return(mockRes)
		}

		resp = httptest.NewRecorder()
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%v/accessor", mockID),
			bytes.NewReader(test.input))
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusMultiStatus, resp.Code)
		// 判断body
		var wantErr *api.Error
		switch test.wantcode {
		case 400:
			wantErr = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": test.params}, Cause: test.cause})
		case 404:
			wantErr = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": test.params}, Cause: test.cause})
		}
		wantRes := api.MultiStatusObject("", nil, wantErr)
		value, _ := json.Marshal([]*api.MultiStatus{wantRes})
		assert.JSONEq(t, string(value), string(resp.Body.Bytes()))
	}
}

// 删除访问者
func TestDeleteAccessors(t *testing.T) {
	teardown := test.SetUpGin(t)
	defer teardown(t)

	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
	defer ctrl.Finish()
	defer gctrl.Finish()
	defer pctrl.Finish()
	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

	// case: 允许
	policy := models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)

	mockID := "network_01"
	mockAccIDs := []string{"user1", "user2"}
	var mockRes []*api.MultiStatus
	mgnt.EXPECT().DeleteAccessors(mockID, mockAccIDs).Return(mockRes)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s/accessor/%s",
		mockID, strings.Join(mockAccIDs, ",")), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMultiStatus, resp.Code)

	// case: 网段不存在
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":true}`),
		Locked:  false,
	}
	mst := api.MultiStatusObject("", nil, errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}}))
	mockRes = []*api.MultiStatus{mst}
	mgnt.EXPECT().DeleteAccessors(mockID, mockAccIDs).Return(mockRes)

	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s/accessor/%s",
		mockID, strings.Join(mockAccIDs, ",")), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMultiStatus, resp.Code)
	wantbody, _ := json.Marshal([]*api.MultiStatus{api.MultiStatusObject("", nil,
		errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{"id"}}}))})
	assert.JSONEq(t, string(wantbody), string(resp.Body.Bytes()))

	// case: 不允许
	policy = models.Policy[[]byte]{
		Name:    "network_restriction",
		Default: []byte(`{"is_enabled":false}`),
		Value:   []byte(`{"is_enabled":false}`),
		Locked:  false,
	}
	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy[[]byte]{policy}, 1, nil)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/policy-mgnt/v1/user-login/network-restriction/network/%s/accessor/%s",
		mockID, strings.Join(mockAccIDs, ",")), nil)
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMultiStatus, resp.Code)
	wantbody, _ = json.Marshal([]*api.MultiStatus{api.MultiStatusObject("", nil,
		errors.ErrNoPermission(&api.ErrorInfo{Cause: "Network restriction function is not enabled."}))})
	assert.JSONEq(t, string(wantbody), string(resp.Body.Bytes()))
}

// 获取传给OPA的 访问者网段具体数据
// func TestGetBundle(t *testing.T) {
// 	teardown := test.SetUpGin(t)
// 	defer teardown(t)
// 	defer os.Remove("allow.rego")
// 	defer os.Remove("boundle.tar.gz")
// 	defer os.Remove("data.json")
// 	defer os.Remove("get_policy.rego")

// 	ctrl, gctrl, pctrl, mgnt, gmgnt, pmgnt := mockNetworkMgnt(t)
// 	defer ctrl.Finish()
// 	defer gctrl.Finish()
// 	defer pctrl.Finish()
// 	router := mockAddNetworkRoute(mgnt, gmgnt, pmgnt)

// 	if _, err := os.Stat(bundlePath); err == nil {
// 		os.Remove(bundlePath)
// 	}

// 	policy := models.Policy{
// 		Name:    "network_restriction",
// 		Default: []byte(`{"is_enabled":false}`),
// 		Value:   []byte(`{"is_enabled":true}`),
// 		Locked:  false,
// 	}
// 	gmgnt.EXPECT().ListPolicy(0, 1, []string{"network_restriction"}).Return([]models.Policy{policy}, 1, nil)

// 	mockRes := make(map[string]interface{})
// 	mgnt.EXPECT().GetNetworkData(mockRes).Return(nil)

// 	mockResult := make(map[string]interface{})
// 	pmgnt.EXPECT().GetWatermarkPolicyData(mockResult).Return(nil)

// 	resp := httptest.NewRecorder()
// 	req := httptest.NewRequest("GET", "/api/policy-mgnt/v1/policy-data/bundle.tar.gz", nil)
// 	router.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusOK, resp.Code)
// }
