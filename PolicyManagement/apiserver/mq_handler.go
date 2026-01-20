// Package apiserver message quene 消息队列处理
package apiserver

import (
	_ "embed" // 标准用法
	"net/http"
	"strings"
	"sync"

	errorv2 "github.com/kweaver-ai/go-lib/error"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/ory/gojsonschema"

	"policy_mgnt/common"
	"policy_mgnt/interfaces"
	"policy_mgnt/logics"
)

// MQHandler 消息队列处理接口
type MQHandler interface {
	// Subscribe 订阅mq消息
	Subscribe()
}

var (
	pOnce sync.Once
	mq    MQHandler

	//go:embed jsonschema/mq/user_created.json
	userCreatedSchemaStr string

	//go:embed jsonschema/mq/user_status_changed.json
	userStatusChangedSchemaStr string

	//go:embed jsonschema/mq/id.json
	userDeletedSchemaStr string
)

type mqHandler struct {
	mqClient                MQClient
	log                     common.Logger
	userCreatedSchema       *gojsonschema.Schema
	userStatusChangedSchema *gojsonschema.Schema
	userDeletedSchema       *gojsonschema.Schema
	event                   interfaces.LogicsEvent
}

// NewMQHandler 创建MQHandler
func NewMQHandler() MQHandler {
	pOnce.Do(func() {
		userCreatedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userCreatedSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userStatusChangedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userStatusChangedSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		userDeletedSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(userDeletedSchemaStr))
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		handler := mqHandler{
			mqClient:                NewMQClient(),
			log:                     common.NewLogger(),
			userCreatedSchema:       userCreatedSchema,
			userStatusChangedSchema: userStatusChangedSchema,
			userDeletedSchema:       userDeletedSchema,
			event:                   logics.NewEvent(),
		}
		mq = &handler
	})

	return mq
}

// Subscribe 订阅mq消息
func (m *mqHandler) Subscribe() {
	channel := "PolicyManagement"
	topicFuncMap := make(map[string]func([]byte) error)
	topicFuncMap["core.user_management.user.created"] = m.onUserCreated
	topicFuncMap["user_management.user.status.changed"] = m.onUserStatusChanged
	topicFuncMap["core.user.delete"] = m.onUserDeleted
	for topic, f := range topicFuncMap {
		m.mqClient.Subscribe(topic, channel, f)
	}
}

func (m *mqHandler) onUserDeleted(message []byte) (err error) {
	// 参数检查
	var jsonV interface{}
	if err = validateAndBindNewError(message, m.userDeletedSchema, &jsonV); err != nil {
		m.log.Errorf("onUserDeleted validateAndBindGin error:%v, message:%s", err, string(message))
		return nil
	}

	userID := jsonV.(map[string]interface{})["id"].(string)
	err = m.event.UserDeleted(userID)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
	}

	return
}

func (m *mqHandler) onUserCreated(message []byte) (err error) {
	// 参数检查
	var jsonV interface{}
	if err = validateAndBindNewError(message, m.userCreatedSchema, &jsonV); err != nil {
		m.log.Errorf("onUserCreated validateAndBindGin error:%v, message:%s", err, string(message))
		// 防止 NSQ 重发消息 返回nil
		return nil
	}

	id := jsonV.(map[string]interface{})["id"].(string)
	err = m.event.UserCreated(id)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
	}

	return
}

func (m *mqHandler) onUserStatusChanged(message []byte) (err error) {
	// 参数检查
	var jsonV interface{}
	if err = validateAndBindNewError(message, m.userStatusChangedSchema, &jsonV); err != nil {
		m.log.Errorf("onUserStatusChanged validateAndBindGin error:%v, message:%s", err, string(message))
		return nil
	}

	userID := jsonV.(map[string]interface{})["user_id"].(string)
	status := jsonV.(map[string]interface{})["status"].(bool)
	err = m.event.UserStatusChanged(userID, status)
	if err != nil {
		if m.needReTry(err) {
			return err
		}
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
