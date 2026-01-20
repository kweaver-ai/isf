package complete_info

import (
	"context"

	"AuditLog/common/enums/oprlogenums"
	"AuditLog/common/helpers/panichelper"
	"AuditLog/common/utils"
	oprlogeo "AuditLog/domain/entity/oprlogeo"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

type CompleteInfo struct {
	documentHttpAcc ihttpaccess.DocumentHttpAcc
	umHttpAcc       ihttpaccess.UmHttpAcc
	eFastHttpAcc    ihttpaccess.EFastHttpAcc
}

func NewCompleteInfo(documentHttpAcc ihttpaccess.DocumentHttpAcc, umHttpAcc ihttpaccess.UmHttpAcc, eFastHttpAcc ihttpaccess.EFastHttpAcc) *CompleteInfo {
	return &CompleteInfo{
		documentHttpAcc: documentHttpAcc,
		umHttpAcc:       umHttpAcc,
		eFastHttpAcc:    eFastHttpAcc,
	}
}

func (l *CompleteInfo) DoCompleteInfo(ctx context.Context, maps []map[string]interface{}, eos []*oprlogeo.LogEntry, bizType oprlogenums.BizType) (err error) {
	defer panichelper.RecoveryAndSetErr(logcmp.GetLogger(), &err)

	if len(eos) == 0 {
		return
	}

	// 1. 补全用户部门信息
	err = l.completeDepInfo(ctx, maps, eos)
	if err != nil {
		return
	}

	// 2. 补全文档库信息
	err = l.completeDocLibInfo(ctx, maps, eos, bizType)
	if err != nil {
		return
	}

	return
}

// mergeEosToMap 将entryEos中的信息合并到maps中
func (l *CompleteInfo) mergeEosToMap(maps []map[string]interface{}, eos []*oprlogeo.LogEntry) (err error) {
	for i := range maps {
		err = utils.MergeMapInterface(maps[i], eos[i])
		if err != nil {
			return
		}
	}

	return
}
