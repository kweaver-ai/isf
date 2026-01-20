package oprinject

import (
	"fmt"
	"sync"

	"AuditLog/common"
	oprsvc "AuditLog/domain/service/operation_log"
	"AuditLog/drivenadapters/httpaccess/document"
	"AuditLog/drivenadapters/httpaccess/httpinject"
	"AuditLog/drivenadapters/httpaccess/usermgnt"
	"AuditLog/gocommon/api"
	"AuditLog/infra"
	oprlogdriveri "AuditLog/interfaces/driveradapter/operation_log"
)

var (
	oprLogSvcOnce sync.Once
	oprLogSvcImpl oprlogdriveri.IOprLogSvc
)

func NewOprLogSvc() oprlogdriveri.IOprLogSvc {
	opsHttpAcc, err := httpinject.NewOpsHttpAcc()
	if err != nil {
		err = fmt.Errorf("[NewOprLogSvc]:NewOpsHttpAcc failed: %w", err)
		panic(err)
	}

	oprLogSvcOnce.Do(func() {
		oprLogSvcImpl = oprsvc.NewOprLogSvc(
			common.SvcConfig.Logger,
			usermgnt.NewUserMgnt(),
			document.NewDocument(),
			api.NewMQClient(),
			infra.NewDBPool(),
			httpinject.NewUmHttpAcc(),
			httpinject.NewEFastHttpAcc(),
			opsHttpAcc,
		)
	})

	return oprLogSvcImpl
}
