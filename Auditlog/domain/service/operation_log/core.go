package oprsvc

import (
	"context"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_log"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/gin-gonic/gin"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
)

func (l *oprLogSvc) HandleMsg(ctx context.Context, msg []byte, bizType oprlogenums.BizType) (err error) {
	// 往msg中添加biz_type
	if bizType.IsServerBizType() {
		msg, err = utils.AddKeyToJSONArrayBys(msg, "biz_type", bizType)
		if err != nil {
			helpers.RecordErrLogWithPos(l.logger, err, "oprsvc.HandleMsg", "AddKeyToJSONArrayBys")
			err = nil

			return
		}
	}

	// 1. 解析请求体
	logMaps, entryEos, err := l.parseMsg(ctx, bizType, msg)
	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "oprsvc.HandleMsg", "parseMsg")
		return
	}

	// 2. 发送到AR
	l.sendToAR(logMaps, bizType)

	// 3. 推荐相关处理 （异步处理）
	go l.recHandle(entryEos, logMaps, bizType)

	// 本地开发api测试用，将logMaps通过http返回
	if helpers.IsLocalDev() {
		if c, ok := ctx.(*gin.Context); ok {
			var out interface{} = logMaps
			if len(logMaps) == 1 {
				out = logMaps[0]
			}

			c.JSON(201, out)
			c.Abort()
		}
	}

	return
}

func (l *oprLogSvc) sendToAR(logMaps []map[string]interface{}, bizType oprlogenums.BizType) {
	// 每个logMap发送一条日志
	for _, logMap := range logMaps {
		ar_log.Logger.InfoField(field.MallocJsonField(logMap), string(bizType))
	}

	l.logger.Infof("[上报运营日志到AR][sendToAR]: 运营日志已调用AR SDK上报, bizType: %s", bizType)
}

// WriteOperationLogToMQ 记录运营日志 to mq
func (l *oprLogSvc) WriteOperationLogToMQ(ctx context.Context, topic oprlogenums.OperationLogTopic, msgByte []byte) (err error) {
	if helpers.IsLocalDev() {
		err = l.HandleMsg(ctx, msgByte, topic.GetBizType())
	} else {
		err = l.mqClient.Publish(string(topic), msgByte)
	}

	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "oprsvc.WriteOperationLogToMQ", "mqClient.Publish")
	} else {
		l.logger.Infof("[API上报运营日志][WriteOperationLogToMQ]: 运营日志发生到mq成功, topic: %s", topic)
	}

	return
}
