package oprsvc

import (
	"context"
	"fmt"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/helpers"
	"AuditLog/common/types/rectypes"
	"AuditLog/common/vars/recvars"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/domain/entity/receo"
)

func (l *oprLogSvc) recHandle(eos []*oprlogeo.LogEntry, logMaps []map[string]interface{}, bizType oprlogenums.BizType) {
	defer func() {
		if e := recover(); e != nil {
			_err := fmt.Errorf("panic: %v", e)
			helpers.RecordErrLogWithPos(l.logger, _err, "oprLogSvc.recHandle", "recover")
		}
	}()

	if !recvars.IsRecBizType(bizType) {
		return
	}

	l.logger.Infof("biz_tye: [%s], recHandle eos: %v", bizType, eos)

	var (
		err error
		ctx = context.Background()
		// operation = eos[0].Operation
		operation = ""
	)

	// 1. Get the index
	index := rectypes.GetIndexByOprBizType(bizType, operation)

	// 2. Batch create logs to open_search
	// recEos := receo.GetEosByOprLogEos(eos)
	// err = l.opsHttpAcc.BatchCreateInterface(ctx, string(index), recEos, false)

	// if err != nil {
	// 	helpers.RecordErrLogWithPos(l.logger, err, "oprLogSvc.recHandle")
	// 	return
	// }

	// 3. AR数据同步到AS
	docs := receo.GetDocs(logMaps)
	err = l.opsHttpAcc.BatchCreateInterface(ctx, string(index), docs, false)

	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "oprLogSvc.recHandle")
		return
	}
}
