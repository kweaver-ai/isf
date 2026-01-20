// Package mq message quene 消息队列处理
package mq

import (
	"net/http"
	"reflect"
	"strings"
	"sync"

	errorv2 "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"Authentication/common"
	sessionSchema "Authentication/driveradapters/jsonschema/sessions_schema"
	userManagementSchema "Authentication/driveradapters/jsonschema/user_management_schema"
	"Authentication/driveradapters/util"
	"Authentication/interfaces"
	accesstokenperm "Authentication/logics/access_token_perm"
	"Authentication/logics/session"
)

// MQHandler 消息队列处理接口
type MQHandler interface {
	// Subscribe 订阅mq消息
	Subscribe()
}

var (
	mOnce    sync.Once
	mq       MQHandler
	clientID string
)

type mqHandler struct {
	mqClient                 MQClient
	accessTokenPerm          interfaces.AccessTokenPerm
	hydraSession             interfaces.HydraSession
	log                      common.Logger
	sessionsSchema           *gojsonschema.Schema
	userDeleteSchema         *gojsonschema.Schema
	userPasswordModifySchema *gojsonschema.Schema
	userStatusChanageSchema  *gojsonschema.Schema
}

type jsonValueDesc = rest.JSONValueDesc

// NewMQHandler 创建MQHandler
func NewMQHandler() MQHandler {
	mOnce.Do(func() {
		sessionsSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(sessionSchema.HydraSessionsSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userDeleteSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserDeleteSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userPasswordModifySchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserPasswordModifySchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userStatusChanageSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userManagementSchema.UserStatusChangeSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		handler := mqHandler{
			mqClient:                 NewMQClient(),
			hydraSession:             session.NewHydraSession(),
			accessTokenPerm:          accesstokenperm.NewAccessTokenPerm(),
			log:                      common.NewLogger(),
			sessionsSchema:           sessionsSchema,
			userDeleteSchema:         userDeleteSchema,
			userPasswordModifySchema: userPasswordModifySchema,
			userStatusChanageSchema:  userStatusChanageSchema,
		}
		mq = &handler
	})

	return mq
}

// Subscribe 订阅mq消息
func (m *mqHandler) Subscribe() {
	channel := "Authentication"
	topicFuncMap := make(map[string]func([]byte) error)
	// 发布-订阅消息
	topicFuncMap["core.app.deleted"] = m.appDeleted
	topicFuncMap["authentication.hydra.sessions.delete"] = m.sessionDelete
	topicFuncMap["core.user.delete"] = m.userDelete
	topicFuncMap["user_management.user.password.modified"] = m.userPasswordModify
	topicFuncMap["user_management.user.status.changed"] = m.userStatusChange

	for topic, f := range topicFuncMap {
		m.mqClient.Subscribe(topic, channel, f)
	}
}

// sessionDelete 删除consent Session
func (m *mqHandler) sessionDelete(message []byte) error {
	var jsonV interface{}
	if err := util.ValidateAndBind(message, m.sessionsSchema, &jsonV); err != nil {
		m.log.Errorln("delete hydra session event, invalid messge params: ", err)
		return nil
	}

	userID := jsonV.(map[string]interface{})["user_id"].(string)
	if userID == "" {
		m.log.Errorln("delete hydra session event, invalid messge params: The user_id is empty")
		return nil
	}

	_, ok := jsonV.(map[string]interface{})["client_id"]
	if ok {
		clientID = jsonV.(map[string]interface{})["client_id"].(string)
	}

	if err := m.hydraSession.Delete(userID, clientID); err != nil {
		m.log.Errorf("delete hydra session event failed, err:%v", err)
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// userDelete 用户删除
func (m *mqHandler) userDelete(message []byte) error {
	var jsonV interface{}
	if err := util.ValidateAndBind(message, m.userDeleteSchema, &jsonV); err != nil {
		m.log.Errorln("delete hydra session event, invalid messge params: ", err)
		return nil
	}

	userID := jsonV.(map[string]interface{})["id"].(string)
	if userID == "" {
		m.log.Errorln("delete hydra session event, invalid messge params: The user_id is empty")
		return nil
	}

	if err := m.hydraSession.Delete(userID, clientID); err != nil {
		m.log.Errorf("delete hydra session event failed, err:%v", err)
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// userPasswordModify 用户密码修改
func (m *mqHandler) userPasswordModify(message []byte) error {
	var jsonV interface{}
	if err := util.ValidateAndBind(message, m.userPasswordModifySchema, &jsonV); err != nil {
		m.log.Errorln("delete hydra session event, invalid messge params: ", err)
		return nil
	}

	userID := jsonV.(map[string]interface{})["user_id"].(string)
	if userID == "" {
		m.log.Errorln("delete hydra session event, invalid messge params: The user_id is empty")
		return nil
	}

	if err := m.hydraSession.Delete(userID, clientID); err != nil {
		m.log.Errorf("delete hydra session event failed, err:%v", err)
		if m.needReTry(err) {
			return err
		}
	}

	return nil
}

// userStatusChange 用户状态变更
func (m *mqHandler) userStatusChange(message []byte) error {
	var jsonV interface{}
	if err := util.ValidateAndBind(message, m.userStatusChanageSchema, &jsonV); err != nil {
		m.log.Errorln("delete hydra session event, invalid messge params: ", err)
		return nil
	}
	status := jsonV.(map[string]interface{})["status"].(bool)

	if !status {
		userID := jsonV.(map[string]interface{})["user_id"].(string)
		if userID == "" {
			m.log.Errorln("delete hydra session event, invalid messge params: The user_id is empty")
			return nil
		}

		if err := m.hydraSession.Delete(userID, clientID); err != nil {
			m.log.Errorf("delete hydra session event failed, err:%v", err)
			if m.needReTry(err) {
				return err
			}
		}
	}
	return nil
}

// appDeleted 应用账户被删除事件
func (m *mqHandler) appDeleted(message []byte) error {
	appID, err := m.getIDInNsqMsg(message)
	if err != nil {
		m.log.Errorf("appDeleted getIDInNsqMsg error:%v", err)
		return nil
	}

	err = m.accessTokenPerm.AppDeleted(appID)
	if err != nil {
		m.log.Errorf("appDeleted error:%v", err)
		if m.needReTry(err) {
			return err
		}
	} else {
		m.log.Infof("Delete app perm successfully: appId: %s", appID)
	}

	return nil
}

func (m *mqHandler) getIDInNsqMsg(message []byte) (id string, err error) {
	var jsonV interface{}
	err = jsoniter.Unmarshal(message, &jsonV)
	if err != nil {
		return
	}

	objDesc := make(map[string]*jsonValueDesc)
	objDesc["id"] = &jsonValueDesc{Kind: reflect.String, Required: true}
	reqParamsDesc := &jsonValueDesc{Kind: reflect.Map, Required: true, ValueDesc: objDesc}
	err = rest.CheckJSONValue("body", jsonV, reqParamsDesc)
	if err != nil {
		return
	}

	id = jsonV.(map[string]interface{})["id"].(string)
	if id == "" {
		err = rest.NewHTTPError("The id is empty.", rest.BadRequest, nil)
	}

	return
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
	case *rest.HTTPError:
		// 本服务 有标识的错误
		statusCode := e.Code
		if statusCode/100000000 == http.StatusInternalServerError/100 {
			// 内部错误
			return true
		}
		// 外部错误
		return false
	case *errorv2.Error:
		// 其他自研服务 抛出的错误
		code := strings.Split(e.Code, ".")[1]
		if code == errorv2.InternalServerError || code == errorv2.ServiceUnavailable {
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
