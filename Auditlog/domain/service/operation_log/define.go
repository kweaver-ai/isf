package oprsvc

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"

	"AuditLog/gocommon/api"
	"AuditLog/interfaces"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
	oprlogdriveri "AuditLog/interfaces/driveradapter/operation_log"
)

type oprLogSvc struct {
	logger       api.Logger
	userMgntRepo interfaces.UserMgntRepo
	mqClient     api.MQClient
	dbPool       *sqlx.DB

	documentHttpAcc ihttpaccess.DocumentHttpAcc
	umHttpAcc       ihttpaccess.UmHttpAcc
	efastHttpAcc    ihttpaccess.EFastHttpAcc
	opsHttpAcc      ihttpaccess.OpsHttpAcc
}

func NewOprLogSvc(
	logger api.Logger,
	userMgntRepo interfaces.UserMgntRepo,
	documentHttpAcc ihttpaccess.DocumentHttpAcc,
	mqClient api.MQClient,
	dbPool *sqlx.DB,
	userHttpAcc ihttpaccess.UmHttpAcc,
	efastHttpAcc ihttpaccess.EFastHttpAcc,
	oprHttpAcc ihttpaccess.OpsHttpAcc,
) oprlogdriveri.IOprLogSvc {
	svc := &oprLogSvc{
		logger:          logger,
		userMgntRepo:    userMgntRepo,
		documentHttpAcc: documentHttpAcc,
		mqClient:        mqClient,
		dbPool:          dbPool,
		umHttpAcc:       userHttpAcc,
		efastHttpAcc:    efastHttpAcc,
		opsHttpAcc:      oprHttpAcc,
	}

	return svc
}
