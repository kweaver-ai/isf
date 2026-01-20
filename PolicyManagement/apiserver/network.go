package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"policy_mgnt/decision"
	"policy_mgnt/general"
	"policy_mgnt/network"
	"policy_mgnt/utils"
	"policy_mgnt/utils/errors"
	"policy_mgnt/utils/models"

	"policy_mgnt/utils/gocommon/api"

	"github.com/kweaver-ai/GoUtils/utilities"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

type networkHandler struct {
	mgnt  network.Management
	gmgnt general.Management
	pmgnt decision.PolicyDecision
}

func newNetworkHandler() (*networkHandler, error) {
	mgnt, err := network.NewManagement()
	if err != nil {
		return nil, err
	}
	gmgnt, err := general.NewManagement()
	if err != nil {
		return nil, err
	}
	pmgnt := decision.NewPolicyDecision()

	return newNetworkHandlerWithMgnt(mgnt, gmgnt, pmgnt), nil
}

func newNetworkHandlerWithMgnt(mgnt network.Management, gmgnt general.Management, pmgnt decision.PolicyDecision) *networkHandler {
	return &networkHandler{
		mgnt:  mgnt,
		gmgnt: gmgnt,
		pmgnt: pmgnt,
	}
}

func (n *networkHandler) AddRouters(r *gin.RouterGroup) {
	tokenCheck := oauth2Middleware(api.NewOAuth2())
	r.GET("/user-login/network-restriction/network", tokenCheck, n.searchNetwork)
	r.POST("/user-login/network-restriction/network", tokenCheck, n.addNetwork)
	r.GET("/user-login/network-restriction/network/:id", tokenCheck, n.getNetworkByID)
	r.PUT("/user-login/network-restriction/network/:id", tokenCheck, n.editNetwork)
	r.DELETE("/user-login/network-restriction/network/:id", tokenCheck, n.deleteNetwork)
	r.GET("/user-login/network-restriction/network/:id/accessor", tokenCheck, n.searchAccessors)
	r.GET("/user-login/network-restriction/accessor/:id/network", tokenCheck, n.getNetworksByAccessorID)
	r.POST("/user-login/network-restriction/network/:id/accessor", tokenCheck, n.addAccessors)
	r.DELETE("/user-login/network-restriction/network/:id/accessor/:accessor_id", tokenCheck, n.deleteAccessors)
	// 策略引擎暂时没有接入oauth，省略token检查
	//r.GET("/policy-data/bundle.tar.gz", n.getBundle)
}

func (n *networkHandler) RegisterRedis() []RedisSubscriber {
	return []RedisSubscriber{
		{Channel: decision.PolicyInitDataChannel, Handler: n.policyInitHandler},
		{Channel: decision.PolicyUpdateDataChannel, Handler: n.policyChangeHandler},
	}
}

func (n *networkHandler) policyChangeHandler(data []byte) {
	l := api.NewLogger()
	var upDateData models.UserDepartRelation
	if err := jsoniter.Unmarshal(data, &upDateData); err != nil {
		l.Errorln("jsoniter.Unmarshal: ", err)
		return
	}

	// TODO 什么情况下会失败，失败了是否触发增量更新
	if err := n.pmgnt.IncrementalUpdatePolicy(context.Background(), upDateData.UserId, upDateData.DepartmentPaths); err != nil {
		l.Errorln("d.ld.IncrementalUpdatePolicy: ", err)
	}
}

func (n *networkHandler) policyInitHandler(_ []byte) {
	// TODO 什么情况下会失败，失败了是否触发全量更新
	l := api.NewLogger()
	l.Infof("Update policy start.")
	if err := n.pmgnt.InitOpaData(); err != nil { //更新包时需调用proton推送全量接口
		l.Errorf("Update data failed:%s.", err)
	} else {
		l.Infof("Update data complete.")
	}
}

// 获取访问者网段开关状态
func (n *networkHandler) getNetworkRestrictionState() (result bool, err error) {
	res, _, err := n.gmgnt.ListPolicy(0, 1, []string{"network_restriction"})
	if err != nil {
		return
	}

	value := res[0].Value
	var networkRestriction general.NetworkResitriction
	err = json.Unmarshal(value, &networkRestriction)
	if err != nil {
		return
	}
	result = networkRestriction.IsEnabled
	return
}

// 开关关闭时，只能查看信息
func (n *networkHandler) checkAvaiable() (err error) {
	state, err := n.getNetworkRestrictionState()

	if err != nil {
		return
	}
	if !state {
		err = errors.ErrNoPermission(&api.ErrorInfo{Cause: "Network restriction function is not enabled."})
		return
	}
	return
}

func (n *networkHandler) searchNetwork(c *gin.Context) {
	start, limit := getPageParams(c)
	if c.IsAborted() {
		return
	}
	keyWord := strings.TrimSpace(c.DefaultQuery("key_word", ""))

	res, count, err := n.mgnt.SearchNetwork(keyWord, start, limit)
	if err != nil {
		errorResponse(c, err)
		return
	}

	listResponse(c, count, res)
}

func (n *networkHandler) getNetworkByID(c *gin.Context) {
	id := c.Param("id")
	res, err := n.mgnt.GetNetworkByID(id)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (n *networkHandler) editNetwork(c *gin.Context) {
	if err := n.checkAvaiable(); err != nil {
		errorResponse(c, err)
		return
	}
	id := c.Param("id")
	var params models.NetworkRestriction
	if err := parseBody(c, &params); err != nil {
		errorResponse(c, err)
		return
	}
	utilities.TrimStruct(&params)
	if err := n.mgnt.EditNetwork(id, &params); err != nil {
		errorResponse(c, err)
		return
	}

	if err := n.pmgnt.PublishInit(); err != nil {
		errorResponse(c, err)
		return
	}
}

func (n *networkHandler) deleteNetwork(c *gin.Context) {
	if err := n.checkAvaiable(); err != nil {
		errorResponse(c, err)
		return
	}
	id := c.Param("id")
	if err := n.mgnt.DeleteNetwork(id); err != nil {
		errorResponse(c, err)
		return
	}

	if err := n.pmgnt.PublishInit(); err != nil {
		errorResponse(c, err)
		return
	}
}

func (n *networkHandler) addNetwork(c *gin.Context) {
	if err := n.checkAvaiable(); err != nil {
		errorResponse(c, err)
		return
	}
	var params models.NetworkRestriction
	if err := parseBody(c, &params); err != nil {
		errorResponse(c, err)
		return
	}
	utilities.TrimStruct(&params)

	networkID, err := n.mgnt.AddNetwork(&params)
	if err != nil {
		errorResponse(c, err)
		return
	}
	c.Redirect(http.StatusCreated, "network/"+networkID)
}

// 检查multistatus中是否包含指定的status
// 只要有一个，返回true，否则返回false
func checkStatusFromMutiRes(res []*api.MultiStatus, status int) bool {
	for _, r := range res {
		if r.Status == status {
			return true
		}
	}
	return false
}

func (n *networkHandler) addAccessors(c *gin.Context) {
	if err := n.checkAvaiable(); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		c.JSON(http.StatusMultiStatus, []*api.MultiStatus{mst})
		return
	}
	id := c.Param("id")

	if err := validJsonData(c, utils.AccessorsSchema); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		c.JSON(http.StatusMultiStatus, []*api.MultiStatus{mst})
		return
	}

	var accessorInfos []models.AccessorInfo
	if err := parseBody(c, &accessorInfos); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		c.JSON(http.StatusMultiStatus, []*api.MultiStatus{mst})
		return
	}

	res := n.mgnt.AddAccessors(id, accessorInfos)
	c.JSON(http.StatusMultiStatus, res)

	if checkStatusFromMutiRes(res, http.StatusCreated) {
		if err := n.pmgnt.PublishInit(); err != nil {
			errorResponse(c, err)
			return
		}
	}
}

func (n *networkHandler) searchAccessors(c *gin.Context) {
	id := c.Param("id")
	start, limit := getPageParams(c)
	if c.IsAborted() {
		return
	}
	keyWord := strings.TrimSpace(c.DefaultQuery("key_word", ""))

	res, count, err := n.mgnt.SearchAccessors(id, keyWord, start, limit)
	if err != nil {
		errorResponse(c, err)
		return
	}

	listResponse(c, count, res)
}

func (n *networkHandler) getNetworksByAccessorID(c *gin.Context) {
	id := c.Param("id")
	start, limit := getPageParams(c)
	if c.IsAborted() {
		return
	}

	res, count, err := n.mgnt.GetNetworksByAccessorID(id, start, limit)
	if err != nil {
		errorResponse(c, err)
		return
	}

	listResponse(c, count, res)
}

func (n *networkHandler) deleteAccessors(c *gin.Context) {
	if err := n.checkAvaiable(); err != nil {
		mst := api.MultiStatusObject("", nil, err)
		c.JSON(http.StatusMultiStatus, []*api.MultiStatus{mst})
		return
	}
	networkID := c.Param("id")
	accessorIDs, err := paramSetTrim(c, "accessor_id")
	if err != nil {
		mst := api.MultiStatusObject("", nil, err)
		c.JSON(http.StatusMultiStatus, []*api.MultiStatus{mst})
		return
	}

	res := n.mgnt.DeleteAccessors(networkID, accessorIDs)
	c.JSON(http.StatusMultiStatus, res)

	if checkStatusFromMutiRes(res, http.StatusOK) {
		if err := n.pmgnt.PublishInit(); err != nil {
			errorResponse(c, err)
			return
		}
	}
}

func (n *networkHandler) flushAccessors(userID string) error {
	accessorID := strings.TrimSpace(userID)
	if err := n.mgnt.FlushAccessors([]string{accessorID}); err != nil {
		return err
	}

	if err := n.pmgnt.PublishInit(); err != nil {
		return err
	}
	return nil
}

func (n *networkHandler) updateUserDepartmentRelation(userId string, departmentIds []string) error {
	userID := strings.TrimSpace(userId)
	if err := n.pmgnt.PublishChange(userID, departmentIds); err != nil {
		return err
	}
	return nil
}
