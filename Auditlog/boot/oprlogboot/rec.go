package oprlogboot

import (
	"context"
	"log"

	recinject "AuditLog/domain/service/inject/rec"
	"AuditLog/driveradapters/task/rectask"
)

func Init() {
	//if helpers.IsLocalDev() {
	//	return
	//}
	// 1. 当open_search索引不存在时创建索引
	ctx := context.Background()

	svc := recinject.NewRecSvc()

	err := svc.CreateByMapping(ctx)
	if err != nil {
		log.Fatalf("[recboot][Init]recsvc.CreateByMapping failed: %v", err)
	}

	// 2. 删除不用的opensearch索引
	err = svc.RemoveNotUseOpensearchIndexOnce(ctx)
	if err != nil {
		log.Fatalf("[recboot][Init]recsvc.RemoveNotUseOpensearchIndexOnce failed: %v", err)
	}

	// 3. 开启任务（删除超过保存时间的日志）
	go rectask.NewRemoveOldLogTask().Run()
}
