package apiserver

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"policy_mgnt/common"
	"policy_mgnt/dbaccess"
	"policy_mgnt/logics"
	"policy_mgnt/utils"
	"policy_mgnt/utils/errors"

	"policy_mgnt/utils/gocommon/api"

	"github.com/kweaver-ai/GoUtils/utilities"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/gin-gonic/gin"
)

// StartAPIServer start api server
func StartAPIServer() {
	// 初始化对象
	// dbTracePool依赖注入
	dbTracePool := common.NewDBTracePool()
	dbaccess.SetDBTracePool(dbTracePool)
	logics.SetDBTracePool(dbTracePool)
	logics.SetDBLicense(dbaccess.NewDBLicense())
	logics.SetDBConfig(dbaccess.NewDBConfig())
	logics.SetDBOutbox(dbaccess.NewOutbox())

	server := newServer()
	port := os.Getenv("API_SERVER_PORT")
	if _, err := strconv.Atoi(port); err != nil {
		port = "9603"
	}
	err := server.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func newServer() *gin.Engine {
	register := NewRegister()
	nh, err := newNetworkHandler()
	if err != nil {
		panic(err)
	}
	register.SubscribeRedis([]RedisSub{nh})
	server := gin.New()
	server.Use(api.Ginrus(), gin.Recovery())

	group := server.Group("/api/policy-management/v1")
	// 健康检查
	hh, err := newHealthHandler()
	if err != nil {
		panic(err)
	}
	hh.AddRouters(group)

	// 通用策略管理
	gh, err := newGeneralHandler()
	if err != nil {
		panic(err)
	}
	gh.AddRouters(group)

	// 访问者网段管理
	nh.AddRouters(group)

	// 访问者网段管理
	dh, err := newDecisionHandler()
	if err != nil {
		panic(err)
	}
	dh.AddRouters(group)

	// 许可证管理
	licenseGroup := server.Group("/api/license/v1")
	lg, err := newLicenseHandler()
	if err != nil {
		panic(err)
	}
	lg.AddRouters(licenseGroup)

	userDeleteChan := make(chan string, 50000)
	deptDeleteChan := make(chan string, 10000)
	departmentAddUserChan := make(chan utils.MsgDepartAddUser, 10000)
	userMovedChan := make(chan utils.MsgUserMoved, 10000)
	departmentRemoveUserChan := make(chan utils.MsgDepartRemoveUser, 10000)
	utils.InitSharemgntConsumer(userDeleteChan, deptDeleteChan, departmentAddUserChan, userMovedChan, departmentRemoveUserChan)
	// 协程：引擎的用户、部门变动，策略管理服务更新表信息
	go handleUser(userDeleteChan, nh)
	go handleDepartment(deptDeleteChan, nh)
	go handleDepartAddUser(departmentAddUserChan, nh)
	go handleUserMOved(userMovedChan, nh)
	go handleDepartRemoveUser(departmentRemoveUserChan, nh)

	// 文档库变动
	// go handleDocLibDelete(docLibDeleteChan, wh)
	// go handleDocLibUpdate(docLibUpdateChan, wh)
	return server
}

// 用户处理
func handleUser(userDeleteChan chan string, nh *networkHandler) {
	l := api.NewLogger()
	for {
		if userID, ok := <-userDeleteChan; ok {
			if err := nh.flushAccessors(userID); err != nil {
				l.Errorf("Delete accessor %v failed, err:%v.", userID, err)
			} else {
				l.Infof("Delete accessor %s complete.", userID)
			}
		}
	}
}

// 部门处理
func handleDepartment(deptDeleteChan chan string, nh *networkHandler) {
	l := api.NewLogger()
	for {
		if departmentID, ok := <-deptDeleteChan; ok {
			if err := nh.flushAccessors(departmentID); err != nil {
				l.Errorf("Delete accessor %v failed, err:%v.n", departmentID, err)
			} else {
				l.Infof("Delete accessor %s complete.n", departmentID)
			}
		}
	}
}

// 部门添加用户处理
func handleDepartAddUser(departmentAddUserChan chan utils.MsgDepartAddUser, nh *networkHandler) {
	l := api.NewLogger()
	for {
		if departAddUserInfo, ok := <-departmentAddUserChan; ok {
			if err := nh.updateUserDepartmentRelation(departAddUserInfo.ID, departAddUserInfo.DepartPaths); err != nil {
				l.Errorf("Departmenrt add user %v failed, err:%v.", departAddUserInfo.ID, err)
			} else {
				l.Infof("Departmenrt add user %v complete.", departAddUserInfo.ID)
			}
		}
	}
}

// 用户移动处理
func handleUserMOved(userMovedChan chan utils.MsgUserMoved, nh *networkHandler) {
	l := api.NewLogger()
	for {
		if userMovedInfo, ok := <-userMovedChan; ok {
			departPaths := []string{userMovedInfo.OldDepartPath, userMovedInfo.NewDepartPath}
			if err := nh.updateUserDepartmentRelation(userMovedInfo.ID, departPaths); err != nil {
				l.Errorf("Move user %v failed, err:%v.", userMovedInfo.ID, err)
			} else {
				l.Infof("Move user %s complete.", userMovedInfo.ID)
			}
		}
	}
}

// 部门移除用户处理
func handleDepartRemoveUser(departmentRemoveUserChan chan utils.MsgDepartRemoveUser, nh *networkHandler) {
	l := api.NewLogger()
	for {
		if departRemoveUserInfo, ok := <-departmentRemoveUserChan; ok {
			if err := nh.updateUserDepartmentRelation(departRemoveUserInfo.ID, departRemoveUserInfo.DepartPaths); err != nil {
				l.Errorf("Departmenrt remove user %v failed, err:%v.", departRemoveUserInfo.ID, err)
			} else {
				l.Infof("Departmenrt remove user %v complete.", departRemoveUserInfo.ID)
			}
		}
	}
}

// 通知opa更新策略内容
// func notifyUpdatePolicy(n api.Produce) error {
// 	if topic == "" {
// 		topic = os.Getenv("POLICY_TOPIC")
// 		if topic == "" {
// 			topic = "policy-data-topic"
// 		}
// 	}

// 	if len(message) == 0 {
// 		mes := os.Getenv("POLICY_MESSAGE")
// 		if mes == "" {
// 			mes = `{"policy_source_service": "policy-management"}`
// 		}
// 		message = []byte(mes)
// 	}

// 	if err := n.Publish(topic, message); err != nil {
// 		return err
// 	}
// 	return nil
// }

// 通道信号处理
type Throttle struct {
	SigChan        chan struct{} // 信号写入通道
	Period         time.Duration // 缓冲间隔
	FTimer         *time.Timer   // 定时器
	ProcessingChan chan struct{} // 结束通道
	AgainChan      chan struct{} // 是否再执行一次
}

func NewThrottle(sigChan, processingChan, againChan chan struct{}, period time.Duration, f func()) *Throttle {
	return &Throttle{
		SigChan:        sigChan,
		Period:         period,
		FTimer:         time.AfterFunc(period, f),
		ProcessingChan: processingChan,
		AgainChan:      againChan,
	}
}

// 处理规则：
// 1、指定缓冲时间内只要有新的调用请求，前面的请求作废，只用最后一个请求;
// 2、处理过程中:同时只能处理一个更新协程;
// 3、处理过程中:有新的处理请求，则这次处理结束后，继续处理一次;
func (t *Throttle) Start() {
	for {
		if _, ok := <-t.SigChan; ok {
			// 如果任务状态是已经结束，重置计时器
			if n := len(t.ProcessingChan); n == 0 {
				t.FTimer.Stop()
				t.FTimer.Reset(t.Period)
			} else { // 如果任务状态是正在进行，新的请求写入AgainChan
				// AgainChan中没有数据时，写入
				if len(t.AgainChan) == 0 {
					t.AgainChan <- struct{}{}
				}
			}
		}
	}
}

func getPageParams(c *gin.Context) (start int, limit int) {
	var err error
	startStr := c.Query("offset")
	limitStr := c.Query("limit")

	if startStr == "" {
		start = 0
	} else {
		start, err = strconv.Atoi(startStr)
		if err != nil {
			params := []string{"offset"}
			err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
			errorResponse(c, err)
			return
		}
		if start < 0 {
			err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"offset"}}})
			errorResponse(c, err)
			return
		}
	}

	if limitStr == "" {
		limit = 20
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			params := []string{"limit"}
			err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": params}})
			errorResponse(c, err)
			return
		}
		if limit < 1 || limit > 1000 {
			err := errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"limit"}}})
			errorResponse(c, err)
			return
		}
	}

	return
}

func paramArray(c *gin.Context, key string) []string {
	// 根据API规范，使用逗号“,”作为 Path Param Array 的分割符
	value := c.Param(key)
	return strings.Split(value, ",")
}

func listResponse(c *gin.Context, count int, data interface{}) {
	c.JSON(200, gin.H{
		"count": count,
		"data":  data,
	})
}

func errorResponse(c *gin.Context, err error) {
	apiErr, ok := err.(*api.Error)
	if !ok {
		apiErr = errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: err.Error()})
	}

	tmpErr := rest.NewHTTPErrorV2(apiErr.Code, apiErr.Cause)

	if apiErr.Detail != nil {
		if temp, ok := apiErr.Detail.(map[string]interface{}); ok {
			tmpErr.Detail = temp
		} else {
			tmpErr.Detail = map[string]interface{}{"detail": apiErr.Detail}
		}
	}
	tmpErr.Description = ""
	tmpErr.Message = apiErr.Message

	rest.ReplyError(c, tmpErr)
	c.Abort()
}

// 解析body中的json数据，如果数据类型不符，报错
// 错误信息中包含解析失败的field
func parseBody(c *gin.Context, params interface{}) (err error) {
	err = c.ShouldBindJSON(&params)
	if err != nil {
		var invalideParams []string
		var field string
		if err1, ok := err.(*json.UnmarshalTypeError); ok {
			field = err1.Field
		}
		invalideParams = append(invalideParams, field)
		err = errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": invalideParams}})
		return
	}
	return
}

// jsonshema 对body进行校验
// TODO: 改为中间件
func validJsonData(c *gin.Context, schema string) error {
	// 解决读取两次body失败
	var bodyBytes []byte
	if c.Request.Body == nil {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"request body"}}})
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"request body"}}})
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": []string{"request body"}}, Cause: "Request body is needed."})
	}
	invalideParams, cause := utilities.ValidJson(schema, string(bodyBytes))
	// 无错误字段，返回
	if len(invalideParams) == 0 {
		return nil
	}

	return errors.ErrBadRequestPublic(&api.ErrorInfo{Detail: map[string]interface{}{"invalid_params": invalideParams}, Cause: cause})
}

// 获取url中多个资源id或名称
// 包括去重、去空格，如果为空字符串，报错
func paramSetTrim(c *gin.Context, identifier string) (res []string, err error) {
	for _, each := range utilities.TrimDupStr(paramArray(c, identifier)) {
		each = strings.TrimSpace(each)
		if each == "" {
			err = errors.ErrNotFound(&api.ErrorInfo{Detail: map[string]interface{}{"notfound_params": []string{identifier}}})
			return
		}
		res = append(res, each)
	}
	return
}
