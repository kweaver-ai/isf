package oprsvc

import (
	"context"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/domain/service/operation_log/complete_info"
	"AuditLog/infra/json_schema/jsc_opr_log"
)

func (l *oprLogSvc) parseMsg(ctx context.Context, bizType oprlogenums.BizType, msg []byte) (logMaps []map[string]interface{}, entryEos []*oprlogeo.LogEntry, err error) {
	// 1. 解析请求体 to map
	logMaps, err = l.MsgToMap(msg)
	if err != nil {
		return
	}

	// 1.1 打印debug日志
	l.logger.Debugf("Write operation log, biz_tye: [%s], Topic:[%s], logMaps: %+v", bizType, bizType.ToTopic(), logMaps)

	// 2. 解析请求体 to struct
	// var recModel recagg.LogEntryForManyObj
	err = utils.JSON().Unmarshal(msg, &entryEos)
	if err != nil {
		return
	}

	// 3. 补全缺少的信息
	if bizType.IsClientBizType() || jsc_opr_log.IsSpecialClientType(bizType, entryEos[0].Operation) {
		c := complete_info.NewCompleteInfo(l.documentHttpAcc, l.umHttpAcc, l.efastHttpAcc)

		err = c.DoCompleteInfo(ctx, logMaps, entryEos, bizType)
		if err != nil {
			helpers.RecordErrLogWithPos(l.logger, err, "oprsvc.parseMsg", "CompleteInfo")
			err = nil

			return
		}
	}

	return
}

func (l *oprLogSvc) MsgToMap(msg []byte) (logMaps []map[string]interface{}, err error) {
	logMaps = make([]map[string]interface{}, 0)

	err = utils.JSON().Unmarshal(msg, &logMaps)
	if err != nil {
		return
	}

	return
}
