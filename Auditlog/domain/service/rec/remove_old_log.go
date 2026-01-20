package recsvc

import (
	"context"
	"time"

	recconsts "AuditLog/common/constants/recenums"
	"AuditLog/common/helpers"
)

func (l *recSvc) RemoveOldLog(ctx context.Context, index string, saveDays int) (err error) {
	// 删除日志
	// conf.SaveDays之前的时间戳
	oldTimeStr := time.Now().UTC().
		AddDate(0, 0, -saveDays).
		Format(time.RFC3339)

	err = l.opsHttpAcc.DeleteDocsByFieldRange(ctx, index, recconsts.CreatedFieldName, nil, oldTimeStr)
	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "recSvc.RemoveOldLog")
		return
	}

	return
}
