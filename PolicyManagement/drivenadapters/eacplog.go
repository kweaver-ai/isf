// Package drivenadapters 日志处理接口
package drivenadapters

import (
	"fmt"
	"policy_mgnt/common/config"
	"policy_mgnt/interfaces"
	"strings"
	"sync"
	"time"

	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/satori/uuid"

	"policy_mgnt/common"
)

var (
	logOnce sync.Once
	l       *eacplogSvc
)

type msgInfo struct {
	operation   string
	description string

	objectType string
	objectID   string
	objectName string

	level string
	typ   string
}

const (
	motUpdate = "update"

	obTUser = "user"

	lTInfo = "INFO"

	typManage = "management"
)

type eacplogSvc struct {
	log           common.Logger
	client        msqclient.ProtonMQClient
	userTypeMap   map[interfaces.VisitorType]string
	clientTypeMap map[interfaces.ClientType]string
	auditLogTopic string
}

// NewEacpLog 创建日志处理对象
func NewEacpLog() *eacplogSvc {
	logOnce.Do(func() {
		client, err := msqclient.NewProtonMQClientFromFile("/sysvol/conf/mq_config.yaml")
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		l = &eacplogSvc{
			log:    common.NewLogger(),
			client: client,
			userTypeMap: map[interfaces.VisitorType]string{
				interfaces.RealName:  "authenticated_user",
				interfaces.Anonymous: "anonymous_user",
				interfaces.App:       "app",
			},
			clientTypeMap: map[interfaces.ClientType]string{
				interfaces.Unknown:      "unknown",
				interfaces.IOS:          "ios",
				interfaces.Android:      "android",
				interfaces.WindowsPhone: "windows_phone", // 审计日志暂时不支持，也应该没有
				interfaces.Windows:      "windows",
				interfaces.MacOS:        "mac_os",
				interfaces.Web:          "web",
				interfaces.MobileWeb:    "mobile_web",
				interfaces.Nas:          "nas", // 审计日志暂时不支持，也应该没有
				interfaces.ConsoleWeb:   "console_web",
				interfaces.DeployWeb:    "deploy_web",
				interfaces.Linux:        "linux",
				interfaces.APP:          "app",
			},
			auditLogTopic: "isf.audit_log.log",
		}
	})

	return l
}

// publish 消息发送
func (e *eacplogSvc) publish(topic string, msg interface{}) error {
	// 发送消息
	message, err := jsoniter.Marshal(msg)
	if err != nil {
		e.log.Errorln(err)
		return err
	}

	err = e.client.Pub(topic, message)
	if err != nil {
		e.log.Errorln(err)
		return err
	}
	return nil
}

func (e *eacplogSvc) writeLog(visitor *interfaces.Visitor, info *msgInfo) (err error) {
	//
	strUUID := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		e.log.Errorln("eacplog write log uuid err :%v, topic:%v", err)
		return
	}

	// 消息整合
	body := make(map[string]interface{})
	body["operation"] = info.operation
	body["description"] = info.description
	body["op_time"] = time.Now().UnixNano() // 纳秒时间戳

	tempOperator := make(map[string]interface{})
	tempOperator["id"] = visitor.ID
	tempOperator["type"] = e.userTypeMap[visitor.Type]
	tempAgent := make(map[string]interface{})
	tempAgent["type"] = e.clientTypeMap[visitor.ClientType]
	tempAgent["ip"] = visitor.IP
	tempAgent["mac"] = visitor.Mac
	tempOperator["agent"] = tempAgent
	body["operator"] = tempOperator

	tempObject := make(map[string]interface{})
	tempObject["type"] = info.objectType
	if info.objectType != "" {
		tempObject["id"] = info.objectID
	}
	if info.objectName != "" {
		tempObject["name"] = info.objectName
	}
	body["object"] = tempObject

	tempLogFrom := make(map[string]interface{})
	tempLogFrom["package"] = "information-security-fabric"
	tempservice := make(map[string]interface{})
	tempservice["name"] = "policy-management"
	tempLogFrom["service"] = tempservice
	body["log_from"] = tempLogFrom

	body["level"] = info.level
	body["out_biz_id"] = strUUID
	body["type"] = info.typ

	err = e.publish(e.auditLogTopic, body)
	if err != nil {
		e.log.Errorln("eacplog write log err :%v, topic:%v", err, e.auditLogTopic)
	}
	return
}

// OpAddAuthorizedProducts 新增产品授权
func (e *eacplogSvc) OpAddAuthorizedProducts(visitor *interfaces.Visitor, name string, products []string) (err error) {
	msg := loadString("IDS_ADD_AUTHORIZED_PRODUCTS")
	if config.Config.Language == "en_US" {
		msg = fmt.Sprintf(msg, strings.Join(products, ","), name)
	} else {
		msg = fmt.Sprintf(msg, name, strings.Join(products, ","))
	}
	info := &msgInfo{
		operation:   motUpdate,
		description: msg,
		objectType:  obTUser,
		objectName:  name,
		level:       lTInfo,
		typ:         typManage,
	}

	return e.writeLog(visitor, info)
}

// OpDeleteAuthorizedProducts 删除产品授权
func (e *eacplogSvc) OpDeleteAuthorizedProducts(visitor *interfaces.Visitor, name string, products []string) (err error) {
	msg := loadString("IDS_DELETE_AUTHORIZED_PRODUCTS")
	if config.Config.Language == "en_US" {
		msg = fmt.Sprintf(msg, strings.Join(products, ","), name)
	} else {
		msg = fmt.Sprintf(msg, name, strings.Join(products, ","))
	}
	info := &msgInfo{
		operation:   motUpdate,
		description: msg,
		objectType:  obTUser,
		objectName:  name,
		level:       lTInfo,
		typ:         typManage,
	}

	return e.writeLog(visitor, info)
}

// OpUpdateAuthorizedProducts 更新产品授权
func (e *eacplogSvc) OpUpdateAuthorizedProducts(visitor *interfaces.Visitor, name string, currentProducts []string, futureProducts []string) (err error) {
	msg := loadString("IDS_UPDATE_AUTHORIZED_PRODUCTS")
	msg = fmt.Sprintf(msg, name, strings.Join(currentProducts, ","), strings.Join(futureProducts, ","))
	info := &msgInfo{
		operation:   motUpdate,
		description: msg,
		objectType:  obTUser,
		objectName:  name,
		level:       lTInfo,
		typ:         typManage,
	}

	return e.writeLog(visitor, info)
}
