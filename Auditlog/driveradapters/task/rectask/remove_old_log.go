package rectask

import (
	"context"
	"log"
	"time"

	"AuditLog/common/helpers"
	recinject "AuditLog/domain/service/inject/rec"
	"AuditLog/infra/config"
	"AuditLog/infra/config/mapping"
	recdriveri "AuditLog/interfaces/driveradapter/rec"
)

type RemoveOldLogTask struct {
	svc     recdriveri.IRecSvc
	recConf *config.RecConf
}

func NewRemoveOldLogTask() *RemoveOldLogTask {
	conf := config.GetConfig().Rec

	svc := recinject.NewRecSvc()

	return &RemoveOldLogTask{
		recConf: conf,
		svc:     svc,
	}
}

func (t *RemoveOldLogTask) Run() {
	if helpers.IsAaronLocalDev() {
		return
	}
	logPrefix := "[rec_task][RemoveOldLogTask]"

	defer func() {
		if _err := recover(); _err != nil {
			log.Printf("%s: run panic recover, err:%v", logPrefix, _err)
		}
	}()

	ctx := context.Background()

	// 每一小时执行一次（多pod部署时，暂时不考虑多个pod多次执行的问题）
	duration := time.Second * time.Duration(t.recConf.RemoveOldLogTaskIntervalSecond)

	ticker := time.NewTicker(duration)

	for range ticker.C {
		log.Printf("%s: run start\n", logPrefix)

		t.do(ctx)

		log.Printf("%s: run end\n", logPrefix)
	}
}

func (t *RemoveOldLogTask) do(ctx context.Context) {
	for index := range mapping.RecMappingMap {
		err := t.svc.RemoveOldLog(ctx, string(index), t.recConf.SaveDays)
		if err != nil {
			log.Printf("[rec_task][RemoveOldLogTask] remove old log failed, index: %s, err: %v", index, err)
			continue
		}

		// 间隔2s
		time.Sleep(time.Second * 2)
	}
}
