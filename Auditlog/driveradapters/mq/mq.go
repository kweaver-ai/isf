package driveradapters

import (
	"database/sql"
	_ "embed"
	"errors"
	"strings"
	"sync"

	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/xeipuuv/gojsonschema"

	"AuditLog/common"
	"AuditLog/common/conf"
	"AuditLog/drivenadapters/httpaccess/usermgnt"
	"AuditLog/drivenadapters/redisaccess"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/logics"
	"AuditLog/models"
)

type mqHandler struct {
	logger         api.Logger
	client         api.MQClient
	logMgnt        interfaces.LogMgnt
	outbox         interfaces.Outbox
	dbPool         *sqlx.DB
	userMgntRepo   interfaces.UserMgntRepo
	cache          sync.Map
	dlm            interfaces.DLM
	lockPrefix     string
	auditLogSchema *gojsonschema.Schema
}

var (
	mqOnce sync.Once
	m      interfaces.MQHandler

	//go:embed jsonschema/audit_log_schema.json
	auditLogSchemaStr string
)

func NewMQHandler() interfaces.MQHandler {
	mqOnce.Do(func() {
		dbConfig := &sqlx.DBConfig{}
		_ = common.Configure(dbConfig, "dbrw.yaml")
		dbPool, _ := sqlx.NewDB(dbConfig)
		o := logics.NewOutbox("client_log")
		l := logics.NewLogMgnt()
		u := usermgnt.NewUserMgnt()

		auditLogSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(auditLogSchemaStr))
		if err != nil {
			panic(err)
		}

		m = &mqHandler{
			logger:         common.SvcConfig.Logger,
			client:         api.NewMQClient(),
			logMgnt:        l,
			outbox:         o,
			dbPool:         dbPool,
			userMgntRepo:   u,
			cache:          sync.Map{},
			dlm:            redisaccess.NewDLM(),
			lockPrefix:     common.LogPrefix,
			auditLogSchema: auditLogSchema,
		}
	})

	return m
}

func (m *mqHandler) Subscribe() {
	// 审计日志
	channel := "audit_log"
	topicFuncMap := map[string]func([]byte) error{
		"as.audit_log.log_login":      m.loginLog,
		"as.audit_log.log_management": m.managementLog,
		"as.audit_log.log_operation":  m.operationLog,
		"isf.audit_log.log":           m.auditLog,
	}

	for t, f := range topicFuncMap {
		m.client.Subscribe(t, channel, f)
	}
}

func (m *mqHandler) loginLog(msg []byte) (err error) {
	err = m.commonLog("login", msg)
	if err != nil {
		return
	}

	return
}

func (m *mqHandler) managementLog(msg []byte) (err error) {
	err = m.commonLog("management", msg)
	if err != nil {
		return
	}

	return
}

func (m *mqHandler) operationLog(msg []byte) (err error) {
	err = m.commonLog("operation", msg)
	if err != nil {
		return
	}

	return
}

func (m *mqHandler) commonLog(logType string, msg []byte) (err error) {
	// 解析请求体
	var jsonV models.AuditLog
	err = jsoniter.Unmarshal(msg, &jsonV)
	if err != nil {
		m.logger.Warnf("Write audit %v log, json unmarshal err %v", logType, err)
		return nil
	}
	m.logger.Infof("Write audit %v log, Topic:[%s], msg: %+v", logType, "as.audit_log.log_"+logType, jsonV)

	invalideParams, cause := api.ValidJson(common.PostAuditLog, string(msg))
	// 无错误字段，返回
	if len(invalideParams) != 0 {
		m.logger.Warnf("Invalid parmas are %v, invalid cause is %v", invalideParams, cause)
		return
	}
	info := &models.ReceiveLogVo{
		Language:   common.SvcConfig.Languaue,
		LogType:    logType,
		LogContent: &jsonV,
	}

	err = m.logMgnt.ReceiveLog(info)
	if err != nil {
		var tx *sql.Tx
		tx, err = m.dbPool.Begin()
		if err != nil {
			return
		}

		// 异常时Rollback
		defer func() {
			switch err {
			case nil:
				// 提交事务
				err = tx.Commit()
				if err != nil {
					m.logger.Warnf("recordLog Transaction Commit Error:%v", err)
					return
				}
				// 触发outbox消息推送线程
				m.outbox.NotifyPushOutboxThread()
			default:
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					m.logger.Warnf("recordLog Transaction Rollback Error:%v", rollbackErr)
					return
				}
			}
		}()

		if logType == "login" {
			err = m.outbox.AddOutboxInfo(common.AuditLoginLogTopic, jsonV, tx)
		} else if logType == "management" {
			err = m.outbox.AddOutboxInfo(common.AuditManagementTopic, jsonV, tx)
		} else if logType == "operation" {
			err = m.outbox.AddOutboxInfo(common.AuditOperationTopic, jsonV, tx)
		} else {
			return
		}

		if err != nil {
			m.logger.Warnf("AddOutboxInfo err, log content: %v", jsonV)
		}
	}

	return nil
}

func (m *mqHandler) auditLog(msg []byte) (err error) {
	// 解析请求体
	jsonReq, err := m.checkAuditLogMsg(msg)
	if err != nil {
		m.logger.Warnf("audit_log auditLog invalid parmas, invalid cause is %v", err)
		return nil
	}

	// 获取日志信息
	logType := jsonReq["type"].(string)
	logInfo, err := m.getLogInfo(jsonReq)
	if err != nil {
		m.logger.Warnf("audit_log auditLog get log info err: %v", err)
		return
	}

	info := &models.ReceiveLogVo{
		Language:   common.SvcConfig.Languaue,
		LogType:    logType,
		LogContent: &logInfo,
	}

	err = m.logMgnt.ReceiveAuditLog(info)
	if err != nil {
		m.logger.Warnf("audit_log auditLog receive log info err: %v", err)

		var tx *sql.Tx
		tx, err = m.dbPool.Begin()
		if err != nil {
			return
		}

		// 异常时Rollback
		defer func() {
			switch err {
			case nil:
				// 提交事务
				err = tx.Commit()
				if err != nil {
					m.logger.Warnf("recordLog Transaction Commit Error:%v", err)
					return
				}
				// 触发outbox消息推送线程
				m.outbox.NotifyPushOutboxThread()
			default:
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					m.logger.Warnf("recordLog Transaction Rollback Error:%v", rollbackErr)
					return
				}
			}
		}()

		if logType == "login" {
			err = m.outbox.AddOutboxInfo(common.AuditLoginLogTopic, logInfo, tx)
		} else if logType == "management" {
			err = m.outbox.AddOutboxInfo(common.AuditManagementTopic, logInfo, tx)
		} else if logType == "operation" {
			err = m.outbox.AddOutboxInfo(common.AuditOperationTopic, logInfo, tx)
		} else {
			return
		}

		if err != nil {
			m.logger.Warnf("AddOutboxInfo err, log content: %v", logInfo)
		}
	}
	return nil
}

// validateAndBind 校验json数据
func (m *mqHandler) validateAndBind(body []byte, schema *gojsonschema.Schema, bind interface{}) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return err
	}
	if !result.Valid() {
		msgList := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			msgList = append(msgList, err.String())
		}
		return errors.New(strings.Join(msgList, "; "))
	}

	if err := jsoniter.Unmarshal(body, bind); err != nil {
		return err
	}

	return nil
}

// checkAuditLogMsg 校验审计日志消息
func (m *mqHandler) checkAuditLogMsg(msg []byte) (data map[string]interface{}, err error) {
	// 检查json格式
	var jsonReq map[string]interface{}
	err = m.validateAndBind(msg, m.auditLogSchema, &jsonReq)
	if err != nil {
		return nil, err
	}

	// 检查operation
	logType := common.MapLogType[jsonReq["type"].(string)]
	bExistOpType := false
	switch logType {
	case interfaces.LogType_Login:
		if _, ok := conf.MapLoginOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			bExistOpType = true
		}
	case interfaces.LogType_Management:
		if _, ok := conf.MapManageOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			bExistOpType = true
		}
	case interfaces.LogType_Operation:
		if _, ok := conf.MapOperOperTypeStrToint[jsonReq["operation"].(string)]; ok {
			bExistOpType = true
		}
	}

	if !bExistOpType {
		return nil, errors.New("operation is not enums")
	}

	// 当用户类型为internal_service时， 一定要传user name
	userData := jsonReq["operator"].(map[string]interface{})
	if userData["type"].(string) == "internal_service" {
		if _, ok := userData["name"]; !ok {
			return nil, errors.New("user name is required")
		}
	}

	// 当用户类型为authenticated_user或者anonymous_user时， 一定要传agent
	if userData["type"].(string) == "authenticated_user" || userData["type"].(string) == "anonymous_user" {
		if _, ok := userData["agent"]; !ok {
			return nil, errors.New("agent is required")
		}
	}

	// 检查object
	if object, ok := jsonReq["object"]; ok {
		objData := object.(map[string]interface{})
		if _, ok := conf.MapObjectTypeStrToint[objData["type"].(string)]; !ok {
			return nil, errors.New("object type is not enums")
		}
	}

	return jsonReq, nil
}

// getLogInfo 获取日志信息
func (m *mqHandler) getLogInfo(jsonReq map[string]interface{}) (logInfo models.AuditLog, err error) {
	// 获取日志信息
	logType := common.MapLogType[jsonReq["type"].(string)]

	switch logType {
	case interfaces.LogType_Login:
		logInfo.OpType = conf.MapLoginOperTypeStrToint[jsonReq["operation"].(string)]
	case interfaces.LogType_Management:
		logInfo.OpType = conf.MapManageOperTypeStrToint[jsonReq["operation"].(string)]
	case interfaces.LogType_Operation:
		logInfo.OpType = conf.MapOperOperTypeStrToint[jsonReq["operation"].(string)]
	}

	logInfo.Msg = jsonReq["description"].(string)
	logInfo.Date = int64(jsonReq["op_time"].(float64) / 1000)

	// 获取操作员信息
	operatorData := jsonReq["operator"].(map[string]interface{})
	logInfo.UserType = operatorData["type"].(string)
	logInfo.UserID = operatorData["id"].(string)
	if _, ok := operatorData["name"]; ok {
		logInfo.UserName = operatorData["name"].(string)
	}

	if v, ok := operatorData["agent"]; ok {
		agent := v.(map[string]interface{})
		logInfo.IP = agent["ip"].(string)
		logInfo.Mac = agent["mac"].(string)
	}

	// 获取object
	if object, ok := jsonReq["object"]; ok {
		objData := object.(map[string]interface{})
		logInfo.ObjType = conf.MapObjectTypeStrToint[objData["type"].(string)]
		if _, ok := objData["id"]; ok {
			logInfo.ObjID = objData["id"].(string)
		}
		if _, ok := objData["name"]; ok {
			logInfo.ObjName = objData["name"].(string)
		}
	}

	// 获取日志其他信息
	if exMsg, ok := jsonReq["ex_msg"]; ok {
		logInfo.Exmsg = exMsg.(string)
	}
	logInfo.Level = common.MapLevelType[jsonReq["level"].(string)]
	logInfo.OutBizID = jsonReq["out_biz_id"].(string)

	if detail, ok := jsonReq["detail"]; ok {
		detailBytes, err := jsoniter.Marshal(detail)
		if err != nil {
			return logInfo, err
		}
		logInfo.AdditionalInfo = string(detailBytes)
	}

	return logInfo, nil
}
