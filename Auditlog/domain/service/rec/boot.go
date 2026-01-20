package recsvc

import (
	"context"

	recconsts "AuditLog/common/constants/recenums"
	"AuditLog/common/helpers"
	"AuditLog/infra/config/mapping"
)

func (l *recSvc) CreateByMapping(ctx context.Context) (err error) {
	if helpers.IsAaronLocalDev() {
		return
	}

	for index, mapping := range mapping.RecMappingMap {
		err = l.opsHttpAcc.CreateIndex(ctx, string(index), mapping, recconsts.IndexSetting)
		if err != nil {
			return
		}
	}

	return
}

func (l *recSvc) RemoveNotUseOpensearchIndexOnce(ctx context.Context) (err error) {
	// 加分布式锁，后续的步骤在锁内执行
	mu := l.dmlCmp.NewMutex("RemoveNotUseOpensearchIndexOnce")
	err = mu.Lock(ctx)
	if err != nil {
		return
	}

	defer func() {
		_err := mu.Unlock()
		if _err != nil {
			l.logger.Errorln("[recSvc][RemoveNotUseOpensearchIndexOnce]: dlm unlock failed:", _err)
		}
	}()

	flag := "operation_log.RemoveNotUseOpensearchIndexOnce.is_deleted"
	deletedValFalg := "true"

	// 1. 判断是否已经删除
	val, err := l.svcConfigRepo.Get(ctx, flag)
	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "recSvc.RemoveNotUseOpensearchIndexOnce.svcConfigRepo.Get")
		return
	}

	if val == deletedValFalg {
		return
	}

	// 2. 删除索引
	indexes := []string{
		"index_llm_rec_opr_log_menu_button_click",
		"index_llm_rec_opr_log_dir_visit_cd",
		"index_llm_rec_opr_log_doc_operation",
	}

	for _, index := range indexes {
		err = l.opsHttpAcc.DeleteIndex(ctx, index)
		if err != nil {
			helpers.RecordErrLogWithPos(l.logger, err, "recSvc.RemoveNotUseOpensearchIndexOnce.DeleteIndex")
			return
		}
	}

	// 3. 标记已删除
	err = l.svcConfigRepo.Set(ctx, flag, deletedValFalg)
	if err != nil {
		helpers.RecordErrLogWithPos(l.logger, err, "recSvc.RemoveNotUseOpensearchIndexOnce.svcConfigRepo.Set")
		return
	}

	return
}
