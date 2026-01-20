package apiserver

import (
	"encoding/json"
	"strings"

	"policy_mgnt/decision"
	"policy_mgnt/general"
	"policy_mgnt/utils"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/kweaver-ai/GoUtils/utilities"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -package mock_general -source ./general.go -destination test/mock_general/general_mock.go

type generalHandler struct {
	mgnt  general.Management
	pmgnt decision.PolicyDecision
}

func newGeneralHandler() (*generalHandler, error) {
	mgnt, err := general.NewManagement()
	if err != nil {
		return nil, err
	}
	pmgnt := decision.NewPolicyDecision()
	return newGeneralHandlerWithMgnt(mgnt, pmgnt), nil
}

func newGeneralHandlerWithMgnt(mgnt general.Management, pmgnt decision.PolicyDecision) *generalHandler {
	return &generalHandler{
		mgnt:  mgnt,
		pmgnt: pmgnt,
	}
}

func (h *generalHandler) AddRouters(r *gin.RouterGroup) {
	tokenCheck := oauth2Middleware(api.NewOAuth2())
	r.GET("/general", tokenCheck, h.getPolicyList)
	r.PUT("/general/:name/value", tokenCheck, h.setPolicyByName)
	r.PUT("/general/:name/state", tokenCheck, h.setPolicyState)
}

func policyToMap(mode string, policy models.Policy[[]byte]) (result map[string]interface{}) {
	var value json.RawMessage
	switch mode {
	case "default":
		value = policy.Default
	case "":
		fallthrough
	case "current":
		value = policy.Value
	}

	result = map[string]interface{}{
		"name":   policy.Name,
		"value":  value,
		"locked": policy.Locked,
	}
	return
}

func loopPolicyToMap(mode string, policyList []models.Policy[[]byte]) (result []interface{}) {
	result = make([]interface{}, len(policyList))
	for idx, policy := range policyList {
		result[idx] = policyToMap(mode, policy)
	}
	return
}

func getMode(c *gin.Context) (mode string) {
	mode = strings.TrimSpace(strings.ToLower(c.Query("mode")))
	switch mode {
	case "default":
	case "current":
	case "":
	default:
		params := []string{mode}
		err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		errorResponse(c, err)
	}
	return
}

func getQueryNames(c *gin.Context) (names []string) {
	nameString := strings.TrimSpace(strings.ToLower(c.Query("name")))
	if nameString == "" {
		names = []string{}
	} else {
		names = strings.Split(nameString, ",")
	}
	return
}

func (h *generalHandler) getPolicyList(c *gin.Context) {
	start, limit := getPageParams(c)
	if c.IsAborted() {
		return
	}

	mode := getMode(c)
	if c.IsAborted() {
		return
	}

	names := getQueryNames(c)
	if c.IsAborted() {
		return
	}

	policyList, count, err := h.mgnt.ListPolicy(start, limit, names)
	if err != nil {
		errorResponse(c, err)
		return
	}

	result := loopPolicyToMap(mode, policyList)

	listResponse(c, count, result)
}

func getNames(c *gin.Context) (names []string) {
	for _, name := range utilities.TrimDupStr(paramArray(c, "name")) {
		name = strings.TrimSpace(name)
		if name == "" {
			err := errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{}}})
			errorResponse(c, err)
			return
		}
		names = append(names, name)
	}
	return
}

// 检查url中的name是否和body中的一致
// 只存在于url中报错
// 只存在于body中忽略
func policyMatch(names []string, params []policyParam) (result []policyParam, err error) {
	var unknownPolicy []string
	for _, name := range names {
		found := false
		for _, param := range params {
			if param.Name == name {
				found = true
				result = append(result, param)
				break
			}
		}
		if !found {
			unknownPolicy = append(unknownPolicy, name)
		}
	}
	if len(unknownPolicy) > 0 {
		err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": unknownPolicy}})
	}
	return
}

func getForce(c *gin.Context) bool {
	forceStr := strings.TrimSpace(c.Query("force"))
	if forceStr == "" {
		forceStr = "false"
	}
	switch forceStr {
	case "true":
		return true
	case "false":
		return false
	default:
		params := []string{"force"}
		// ss
		errorResponse(c, errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}}))
		return false
	}
}

func (h *generalHandler) setPolicyByName(c *gin.Context) {
	force := getForce(c)
	if c.IsAborted() {
		return
	}
	names := getNames(c)
	if c.IsAborted() {
		return
	}

	err := validJsonData(c, utils.PolicesSchema)
	if err != nil {
		errorResponse(c, err)
		return
	}

	var params []policyParam
	err = c.ShouldBindJSON(&params)
	if err != nil {
		params := []string{}
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		errorResponse(c, err)
		return
	}

	// 检查传入的body是否为空
	if len(params) == 0 {
		params := []string{}
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		errorResponse(c, err)
		return
	}

	params, err = policyMatch(names, params)
	if err != nil {
		errorResponse(c, err)
		return
	}

	policies := make(map[string][]byte)
	for _, param := range params {
		policies[param.Name] = param.Value
	}

	err = h.mgnt.SetPolicyValue(policies, force)
	if err != nil {
		errorResponse(c, err)
		return
	}

	// 如果修改了访问者网段开关，通知更新
	if utilities.InStrSlice((&general.NetworkResitriction{}).Name(), names) {
		if err := h.pmgnt.PublishInit(); err != nil {
			errorResponse(c, err)
			return
		}
	}

	// 如果修改了未配置访问者网段访问开关，通知更新
	if utilities.InStrSlice((&general.NoNetworkPolicyAccessor{}).Name(), names) {
		if err := h.pmgnt.PublishInit(); err != nil {
			errorResponse(c, err)
			return
		}
	}

	// 如果修改了客户端登录选项，通知更新
	if utilities.InStrSlice((&general.ClientRestriction{}).Name(), names) {
		if err := h.pmgnt.PublishInit(); err != nil {
			errorResponse(c, err)
			return
		}
	}
}

func (h *generalHandler) setPolicyState(c *gin.Context) {
	names := getNames(c)
	var data map[string]interface{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		params := []string{}
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
		errorResponse(c, err)
		return
	}

	states := make(map[general.State]interface{})
	var unknownState []string
	for name, value := range data {
		switch name {
		case string(general.StateLocked):
			if _, ok := value.(bool); !ok {
				params := []string{name}
				err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
				errorResponse(c, err)
				return
			}
			states[general.StateLocked] = value
		default:
			unknownState = append(unknownState, name)
		}
	}

	if len(states) == 0 {
		err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": data}})
		errorResponse(c, err)
		return
	}

	if len(unknownState) > 0 {
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": unknownState}})
		errorResponse(c, err)
		return
	}

	err = h.mgnt.SetPolicyState(names, states)
	if err != nil {
		errorResponse(c, err)
		return
	}
}
