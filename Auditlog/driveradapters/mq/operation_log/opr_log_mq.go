package oprlogmq

import (
	"context"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/helpers"
	"AuditLog/common/helpers/panichelper"
	"AuditLog/common/utils"
	"AuditLog/infra/json_schema/jsc_opr_log"
)

func (m *oprLogMqHandler) Subscribe() {
	if helpers.IsAaronLocalDev() {
		m.localTest()
		return
	}

	channel := "audit_log"
	for _, topic := range oprlogenums.GetAllOLT() {

		handler := m.getHandlerByTopic(topic)
		m.client.Subscribe(string(topic), channel, handler)
	}
}

func (m *oprLogMqHandler) getHandlerByTopic(topic oprlogenums.OperationLogTopic) func([]byte) error {
	return func(msg []byte) (err error) {
		defer func() {
			if _err := recover(); _err != nil {
				panicLogMsg := panichelper.PanicTraceErrLog(_err)
				m.logger.Errorln(panicLogMsg)

				err = fmt.Errorf("[getHandlerByTopic][panic recovery]: %v", _err)
			}
		}()

		err = m.commonHandle(topic.GetBizType(), msg)
		if err != nil {
			return
		}

		return
	}
}

func (m *oprLogMqHandler) commonHandle(bizType oprlogenums.BizType, msg []byte) (err error) {
	ctx := context.Background()

	// debug模式下打印日志
	if helpers.IsOprLogShowLogForDebug() {
		var formattedJSON string

		formattedJSON, err = utils.FormatJSONString(string(msg))
		if err != nil {
			m.logger.Warnf("[运营日志上报][oprLogMqHandler] biz_tye: [%s], utils.FormatJSONString(string(msg)) err %v", bizType, err)
			err = nil
			return
		}

		fmt.Printf("[运营日志mq接收到日志]:\n%v\n", formattedJSON)
	}

	// 1. 检查业务类型
	if !bizType.Check() {
		m.logger.Warnf("[运营日志上报][oprLogMqHandler] biz_tye: [%s], invalid bizType", bizType)
		err = nil
		return
	}

	result := gjson.GetBytes(msg, "@this")
	if !result.IsArray() {
		msg = utils.JSONObjectToArray(msg)
	}

	// 2. 校验日志格式
	operation, invalidFields, err := jsc_opr_log.ValidateOprLogJSONSchema(msg, bizType)
	if err != nil {
		m.logger.Errorf("[运营日志上报][oprLogMqHandler] biz_tye: [%s], jsc_opr_log.ValidateOprLogJSONSchema(msg, bizType) err %v", bizType, err)
		err = nil
		return
	}

	// 有错误字段，返回
	if len(invalidFields) != 0 {

		invalidFieldsPretty := make([]string, len(invalidFields))
		for i, field := range invalidFields {
			invalidFieldsPretty[i] = fmt.Sprintf("%d: [%s]", i+1, field)
		}

		m.logger.Errorf("[运营日志上报][oprLogMqHandler] biz_tye: [%s], operation: [%s], invalidFields: ( %s ), log content: ( %s )", bizType, operation, strings.Join(invalidFieldsPretty, " , "), string(msg))
		return
	}

	// 3. 领域服务处理
	err = m.oprLogSvc.HandleMsg(ctx, msg, bizType)

	return
}

func (m *oprLogMqHandler) localTest() {
	// 1. DocumentDomainSync
	topic := oprlogenums.DocumentDomainSync.ToTopic()
	handler := m.getHandlerByTopic(topic)

	err := handler(getMockMsgSingleDocumentDomainSync())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgSingle_DocumentDomainSync()) err %v", err))
	}

	//	2. DirVisit
	topic = oprlogenums.DirVisit.ToTopic()
	handler = m.getHandlerByTopic(topic)
	err = handler(getMockMsgSingleDirVisit())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgSingleDirVisit()) err %v", err))
	}

	//	3. MenuButtonClick
	topic = oprlogenums.MenuButtonClick.ToTopic()
	handler = m.getHandlerByTopic(topic)
	// 3.1 single
	err = handler(getMockMsgSingleInvalidButtonClick())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgSingleInvalidButtonClick()) err %v", err))
	}
	//	3.2 many
	err = handler(getMockMsgManyMenuButtonClick())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgManyMenuButtonClick()) err %v", err))
	}

	// 4. docOperation
	topic = oprlogenums.DocOperation.ToTopic()
	handler = m.getHandlerByTopic(topic)
	err = handler(getMockMsgDocOperation())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgDocOperation()) err %v", err))
	}

	// 5. kcOperation
	topic = oprlogenums.KcOperation.ToTopic()
	handler = m.getHandlerByTopic(topic)
	err = handler(getMockMsgKcOperation())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgKcOperation()) err %v", err))
	}

	// 6. ClientOperation
	topic = oprlogenums.ClientOperation.ToTopic()
	handler = m.getHandlerByTopic(topic)
	err = handler(getMockMsgClientOperation())
	if err != nil {
		panic(fmt.Sprintf("handler(getMockMsgClientOperation()) err %v", err))
	}

	// // kcOperation
	// topic = oprlogenums.KcOperation.ToTopic()
	// handler = m.getHandlerByTopic(topic)
	// err = handler(getMockMsgKcOperation())
	// if err != nil {
	// 	panic(fmt.Sprintf("handler(getMockMsgKcOperation()) err %v", err))
	// }
}
