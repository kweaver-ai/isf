// Package driveradapters message quene 消息队列处理
package driveradapters

import (
	_ "embed" // 标准用法
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"

	errorv2 "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"UserManagement/common"
	"UserManagement/interfaces"
	"UserManagement/logics"
)

// MQHandler 消息队列处理接口
type MQHandler interface {
	// Subscribe 订阅mq消息
	Subscribe()

	// OnUserDeleted 用户被删除 topic core.user.delete
	OnUserDeleted(message []byte) (err error)
}

var (
	pOnce sync.Once
	mq    MQHandler
)

type mqHandler struct {
	mqClient            MQClient
	anonymous           interfaces.LogicsAnonymous
	app                 interfaces.LogicsApp
	orgPermApp          interfaces.LogicsOrgPermApp
	logger              common.Logger
	combine             interfaces.LogicsCombine
	event               interfaces.LogicsEvent
	depart              interfaces.LogicsDepartment
	idsSchema           *gojsonschema.Schema
	orgNameModifySchema *gojsonschema.Schema
}

var (
	//go:embed jsonschema/common/ids.json
	idsSchemaStr string

	//go:embed jsonschema/mq/org_name_modify.json
	orgNameModifySchemaStr string
)

// NewMQHandler 创建MQHandler
func NewMQHandler() MQHandler {
	pOnce.Do(func() {
		idsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(idsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		orgNameModifySchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(orgNameModifySchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		handler := mqHandler{
			mqClient:            NewMQClient(),
			anonymous:           logics.NewAnonymous(),
			app:                 logics.NewApp(),
			logger:              common.NewLogger(),
			orgPermApp:          logics.NewOrgPermApp(),
			combine:             logics.NewCombine(),
			event:               logics.NewEvent(),
			depart:              logics.NewDepartment(),
			idsSchema:           idsSchema,
			orgNameModifySchema: orgNameModifySchema,
		}
		mq = &handler
	})

	return mq
}

// Subscribe 订阅mq消息
func (m *mqHandler) Subscribe() {
	channel := "UserManagement"
	topicFuncMap := make(map[string]func([]byte) error)
	topicFuncMap["core.anonymity.set"] = m.onCreateAnonymous
	topicFuncMap["core.anonymity.delete"] = m.onDeleteAnonymous
	topicFuncMap["core.user_management.app.delete"] = m.onDeleteApp
	topicFuncMap["core.app.name.modified"] = m.UpdateAppName
	topicFuncMap["core.user.delete"] = m.OnUserDeleted
	topicFuncMap["core.dept.delete"] = m.OnDepartDeleted
	topicFuncMap["user_management.org_manager.changed"] = m.OrgManagerChanged
	topicFuncMap["core.org.name.modify"] = m.OnOrgNameChanged
	for topic, f := range topicFuncMap {
		m.mqClient.Subscribe(topic, channel, f)
	}
}

func (m *mqHandler) getIDInMsg(message []byte) (id string, err error) {
	var jsonV interface{}
	if err := jsoniter.Unmarshal(message, &jsonV); err != nil {
		return "", err
	}
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if err := rest.CheckJSONValue("body", jsonV, reqParamsDesc); err != nil {
		return "", err
	}

	id = jsonV.(map[string]interface{})["id"].(string)
	return
}

// OrgManagerChanged 部门管理员变更事件
func (m *mqHandler) OrgManagerChanged(message []byte) (err error) {
	// 消息解析 参数检查
	log := common.NewLogger()

	var jsonV interface{}
	if err = validateAndBind(message, m.idsSchema, &jsonV); err != nil {
		log.Errorf("OrgManagerChanged getIDsInMsg error:%v", err)
		// 防止 NSQ 重发消息 返回nil
		return nil
	}

	// 更新配额信息
	deptIDs := make([]string, 0)
	for _, v := range jsonV.(map[string]interface{})["ids"].([]interface{}) {
		deptIDs = append(deptIDs, v.(string))
	}

	err = m.event.OrgManagerChanged(deptIDs)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// OnDepartDeleted 部门被删除
func (m *mqHandler) OnDepartDeleted(message []byte) (err error) {
	// 消息解析 参数检查
	log := common.NewLogger()
	deptID, err := m.getIDInMsg(message)
	if err != nil {
		log.Errorf("OnDepartDeleted getIDsInMsg error:%v", err)
		// 防止 NSQ 重发消息 返回nil
		return nil
	}

	// 部门被删除事件
	err = m.event.DeptDeleted(deptID)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// OnUserDeleted 用户被删除
func (m *mqHandler) OnUserDeleted(message []byte) (err error) {
	// 消息解析 参数检查
	log := common.NewLogger()
	userID, err := m.getIDInMsg(message)
	if err != nil {
		log.Errorf("OnUserDeleted getIDsInMsg error:%v", err)
		// 防止 NSQ 重发消息 返回nil
		return nil
	}

	// 部门被删除事件
	err = m.event.UserDeleted(userID)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// OnOrgNameChanged 组织架构名称变更
func (m *mqHandler) OnOrgNameChanged(message []byte) (err error) {
	// 消息解析 参数检查
	log := common.NewLogger()

	var jsonV interface{}
	if err = validateAndBind(message, m.orgNameModifySchema, &jsonV); err != nil {
		log.Errorf("OnOrgNameChanged validateAndBind error:%v", err)
		// 防止 NSQ 重发消息 返回nil
		return nil
	}

	// 获取信息
	id := jsonV.(map[string]interface{})["id"].(string)
	newName := jsonV.(map[string]interface{})["new_name"].(string)
	strType := jsonV.(map[string]interface{})["type"].(string)

	if strType == strUser {
		err = m.event.UserNameChanged(id, newName)
		if err != nil {
			if m.needReTry(err) {
				return err
			}
		}
	}

	return nil
}

// UpdateAppName 更新应用账户名称
func (m *mqHandler) UpdateAppName(message []byte) (err error) {
	var jsonV interface{}
	if err = jsoniter.Unmarshal(message, &jsonV); err != nil {
		return nil
	}
	objDesc := make(map[string]*jsonValueDesc)
	objDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	objDesc["new_name"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	if err = rest.CheckJSONValue("body", jsonV, reqParamsDesc); err != nil {
		m.logger.Errorln("update app name event ,invalid messge params:", err)
		return nil
	}
	var info interfaces.AppInfo
	info.ID = jsonV.(map[string]interface{})["id"].(string)
	info.Name = jsonV.(map[string]interface{})["new_name"].(string)

	// 更新应用账户名称
	err = m.orgPermApp.UpdateAppName(&info)
	if err != nil {
		m.logger.Errorln("update app name event ,UpdateAppName error")
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

func (m *mqHandler) onDeleteApp(message []byte) (err error) {
	var jsonV interface{}
	if err = jsoniter.Unmarshal(message, &jsonV); err != nil {
		m.logger.Errorln("delete app event ,messge body json is error")
		return nil
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		m.logger.Errorln("delete app event ,invalid messge params:", cErr)
		return nil
	}

	// 获取具体请求参数
	jsonObj := jsonV.(map[string]interface{})
	id := jsonObj["id"].(string)

	// 通过id 获取用户和部门名称信息
	err = m.app.DeleteApp(nil, id)
	if err != nil {
		m.logger.Errorln("delete app event ,Delete is error")
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// createAnonymous 创建匿名账户
func (m *mqHandler) onCreateAnonymous(message []byte) (err error) {
	var jsonV interface{}
	if err = jsoniter.Unmarshal(message, &jsonV); err != nil {
		m.logger.Errorln("create anonymous event ,messge body json is error")
		return nil
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	paramDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["password"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["limited_times"] = &jsonValueDesc{Kind: reflect.Float64, Required: true}
	paramDesc["expires_at"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["type"] = &jsonValueDesc{Kind: reflect.String, Required: false}
	paramDesc["verify_mobile"] = &jsonValueDesc{Kind: reflect.Bool, Required: false}

	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		m.logger.Errorln("create anonymous event ,invalid messge params:", cErr)
		return nil
	}

	// 获取具体请求参数
	var info interfaces.AnonymousInfo
	jsonObj := jsonV.(map[string]interface{})
	info.Password = jsonObj["password"].(string)
	info.LimitedTimes = int32(jsonObj["limited_times"].(float64))
	info.ID = jsonObj["id"].(string)
	if paramDesc["type"].Exist {
		info.Type = jsonObj["type"].(string)
	}
	// 匿名文档共享
	if strings.HasPrefix(info.ID, "AA") {
		info.VerifyMobile = jsonObj["verify_mobile"].(bool)
	}

	if cErr := m.checkPassword(info.Password); cErr != nil {
		m.logger.Errorln("create anonymous event ,messge password is error")
		return nil
	}

	if cErr := m.stringToTimeStamp(jsonObj["expires_at"].(string), &info.ExpiresAtStamp); cErr != nil {
		m.logger.Errorln("create anonymous event ,messge expires_at is error")
		return nil
	}

	if cErr := m.checkLimitedTimes(info.LimitedTimes); cErr != nil {
		m.logger.Errorln("create anonymous event ,messge limited_times is error")
		return nil
	}

	// 通过id 获取用户和部门名称信息
	err = m.anonymous.Create(info)
	if err != nil {
		m.logger.Errorln("create anonymous event ,Create is error")
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// deleteAnonymous 删除匿名账户
func (m *mqHandler) onDeleteAnonymous(message []byte) (err error) {
	var jsonV interface{}
	if err = jsoniter.Unmarshal(message, &jsonV); err != nil {
		m.logger.Errorln("delete anonymous event ,messge body json is error")
		return nil
	}

	// 检查请求参数与文档是否匹配
	paramDesc := make(map[string]*jsonValueDesc)
	strDesc := make(map[string]*jsonValueDesc)
	strDesc["element"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	paramDesc["ids"] = &jsonValueDesc{Kind: reflect.Slice, Required: true, ValueDesc: strDesc}

	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: paramDesc}
	if cErr := rest.CheckJSONValue("body", jsonV, reqParamsDesc); cErr != nil {
		m.logger.Errorln("delete anonymous event ,invalid messge params:", cErr)
		return nil
	}

	// 获取具体请求参数
	jsonObj := jsonV.(map[string]interface{})
	anonymousIDs := jsonObj["ids"].([]interface{})

	// 通过ID删除匿名账户
	errIDs := make([]string, 0)
	for _, v := range anonymousIDs {
		id := v.(string)
		err = m.anonymous.DeleteByID(id)
		if err != nil {
			errIDs = append(errIDs, id)
		}
	}

	if len(errIDs) != 0 {
		out := strings.Join(errIDs, ", ")
		m.logger.Errorln("delete anonymous event ,delete by ID is error, error anonymity id :" + out)
		if m.needReTry(err) {
			return err
		}
	}
	return nil
}

func (m *mqHandler) stringToTimeStamp(timeStr string, stamp *int64) error {
	timeStamp, err := rest.StringToTimeStamp(timeStr)
	if err != nil {
		return rest.NewHTTPError(err.Error(), rest.BadRequest,
			map[string]interface{}{"params": []string{0: "expires_at"}})
	}
	if (timeStamp != 0) && (timeStamp < common.Now().UnixNano()) {
		return rest.NewHTTPError("invalid expires at", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "expires_at"}})
	}
	*stamp = timeStamp

	return nil
}

func (m *mqHandler) checkPassword(password string) error {
	str := "^([a-zA-Z0-9~!%#$@\\-_.]{4,100})$"
	reg := regexp.MustCompile(str)
	if password != "" && !reg.MatchString(password) {
		return rest.NewHTTPError("invalid password", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "password"}})
	}

	return nil
}

func (m *mqHandler) checkLimitedTimes(limitedTimes int32) error {
	if limitedTimes != -1 && limitedTimes < 1 {
		return rest.NewHTTPError("invalid limited times", rest.BadRequest,
			map[string]interface{}{"params": []string{0: "limited_times"}})
	}
	return nil
}

func (m *mqHandler) needReTry(err error) bool {
	// 错误模型说明：
	// 错误类型分为外部错误和内部错误。每种错误可用错误码对其进行标识，错误码分为通用错误码和指定错误码。
	// 外部错误：由于接口调用方传参或调用方式不对引起的错误。标识错误码使用400系列，即以"4"开头。
	// 内部错误：由于服务自身代码问题引起的错误。标识错误码使用500系列，即以"5"开头。
	// 通用错误码：所有服务公用的错误码，一般用于标识调用方无需特殊处理的错误。当前有如下通用错误码：
	//          - 400000000：客户端请求错误
	//          - 401000000：未授权或授权已过期
	//          - 500000000：服务端内部错误
	// 指定错误码：服务自定义的错误码，以某唯一的错误码对该错误进行标识（在errors模块中定义），一般用于标识调用方需特殊处理的错误，如UI界面提示错误信息。
	// rest.HTTPError：用于存储本服务的错误信息，必有标识错误码
	// errorv2.Error: 存储错误码是string类型的错误信息，必有标识错误码
	// errors.New("xxx")：用于表示本服务产生的内部错误，当作为RESTFul API错误返回时，会在返回处使用"500000000"通用错误码进行标识。如需用指定错误码标识，则使用rest.HTTPError来表示。
	if err == nil {
		return false
	}
	// 外部错误 不需要重试 (访问者不存在、文件不存在) 可能返回值404 考虑一些非标准的错误码
	// 内部错误 服务网络不通 、服务内部崩溃、数据库崩溃 500 503 等 出错 保留审核信息 需要重试
	switch e := err.(type) {
	case *errorv2.Error:
		code := strings.Split(e.Code, ".")[1]
		if code == errorv2.InternalServerError || code == errorv2.ServiceUnavailable {
			// 内部错误
			return true
		}
		// 外部错误
		return false
	case *rest.HTTPError:
		// 本服务 有标识的错误
		statusCode := e.Code
		if statusCode/100000000 == http.StatusInternalServerError/100 {
			// 内部错误
			return true
		}
		// 外部错误
		return false
	default:
		// 本服务 未标识的的内部错误。若作为RESTFul API错误返回时，会在返回处使用"500000000"通用错误码进行标识
		return true
	}
}
