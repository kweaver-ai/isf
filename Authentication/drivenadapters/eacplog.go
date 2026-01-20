package drivenadapters

import (
	"sync"

	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/satori/uuid"

	"Authentication/common"
	"Authentication/interfaces"
)

var (
	logOnce sync.Once
	l       *eacplogSvc
)

type msgInfo struct {
	msg      string
	exMsg    string
	logType  logType
	logLevel logLevel
	opType   int32
	objID    string
}

type eacplogSvc struct {
	log         common.Logger
	client      msqclient.ProtonMQClient
	userTypeMap map[interfaces.VisitorType]string
	logTypeMap  map[logType]string
}

// logType 日志类型
type logType int

const (
	// ltLogin 登录日志
	ltLogin logType = 10

	// ltManage 管理日志
	ltManage logType = 11

	// ltOperation 操作日志
	ltOperation logType = 12
)

// logLevel 消息级别
type logLevel int

const (
	// llInfo 信息
	llInfo logLevel = 1

	// llWarn 警告
	llWarn logLevel = 2
)

// manageOpType 管理日志操作类型
type manageOpType int

const (
	// mtCreate 创建
	mtCreate manageOpType = 1

	// mtDelete 删除
	mtDelete manageOpType = 4
)

// NewEacpLog 创建日志处理对象
func NewEacpLog() *eacplogSvc {
	logOnce.Do(func() {
		client, err := msqclient.NewProtonMQClientFromFile("/sysvol/conf/service_conf/mq_config.yaml")
		if err != nil {
			common.NewLogger().Fatalln(err)
		}

		l = &eacplogSvc{
			log:    common.NewLogger(),
			client: client,
			userTypeMap: map[interfaces.VisitorType]string{
				interfaces.RealName:  "authenticated_user",
				interfaces.Anonymous: "anonymous_user",
				interfaces.Business:  "app",
			},
			logTypeMap: map[logType]string{
				ltLogin:     "as.audit_log.log_login",
				ltManage:    "as.audit_log.log_management",
				ltOperation: "as.audit_log.log_operation",
			},
		}
	})

	return l
}

// publish 消息发送
func (e *eacplogSvc) Publish(topic string, msg interface{}) error {
	return e.publish(topic, msg)
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

// 记录日志
func (e *eacplogSvc) writeLog(visitor *interfaces.Visitor, info *msgInfo) (err error) {
	//
	strUUID := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		e.log.Errorln("eacplog write log uuid err :%v, topic:%v", err)
		return
	}

	// 消息整合
	body := make(map[string]interface{})
	body["user_id"] = visitor.ID
	body["user_type"] = e.userTypeMap[visitor.Type]
	body["level"] = info.logLevel
	body["date"] = common.Now().UnixNano() / 1e3
	body["ip"] = visitor.IP
	body["mac"] = visitor.Mac
	body["msg"] = info.msg
	body["ex_msg"] = info.exMsg
	body["user_agent"] = visitor.UserAgent
	body["op_type"] = info.opType
	body["out_biz_id"] = strUUID
	body["obj_id"] = info.objID

	err = e.publish(e.logTypeMap[info.logType], body)
	if err != nil {
		e.log.Errorln("eacplog write log err :%v, topic:%v", err, e.logTypeMap[info.logType])
	}
	return
}

func (e *eacplogSvc) OpSetAppAccessTokenPerm(visitor *interfaces.Visitor, appName string) error {
	msg := loadString("IDS_SET_ACCESSTOKENPERM", appName)
	info := &msgInfo{
		msg:      msg,
		exMsg:    loadString("IDS_ACCESSTOKENPERM_EXMSG"),
		logType:  ltManage,
		logLevel: llInfo,
		opType:   int32(mtCreate),
	}

	return e.writeLog(visitor, info)
}

func (e *eacplogSvc) OpDeleteAppAccessTokenPerm(visitor *interfaces.Visitor, appName string) error {
	msg := loadString("IDS_DELETE_ACCESSTOKENPERM", appName)
	info := &msgInfo{
		msg:      msg,
		exMsg:    loadString("IDS_ACCESSTOKENPERM_EXMSG"),
		logType:  ltManage,
		logLevel: llWarn,
		opType:   int32(mtDelete),
	}

	return e.writeLog(visitor, info)
}
