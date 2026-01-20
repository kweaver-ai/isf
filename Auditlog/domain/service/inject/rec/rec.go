package recdriveri

import (
	"fmt"
	"sync"

	"AuditLog/common"
	"AuditLog/common/constants"
	"AuditLog/common/helpers/dlmhelper"
	persrecadmininject "AuditLog/domain/service/inject/persrec/admin"
	recsvc "AuditLog/domain/service/rec"
	persrec_db "AuditLog/drivenadapters/db/persrec"
	"AuditLog/drivenadapters/httpaccess/httpinject"
	"AuditLog/infra/cmp/redisdlmcmp"
	recdriveri "AuditLog/interfaces/driveradapter/rec"
)

var (
	oprLogSvcOnce sync.Once
	oprLogSvcImpl recdriveri.IRecSvc
)

func NewRecSvc() recdriveri.IRecSvc {
	opsHttpAcc, err := httpinject.NewOpsHttpAcc()
	if err != nil {
		err = fmt.Errorf("[NewRecSvc]:NewOpsHttpAcc failed: %w", err)
		panic(err)
	}

	repoBase := persrecadmininject.GetPersRepoBase()

	oprLogSvcOnce.Do(func() {
		oprLogSvcImpl = recsvc.NewRecSvc(
			common.SvcConfig.Logger,
			opsHttpAcc,
			persrec_db.NewSvcConfigRepo(repoBase),
			redisdlmcmp.NewRedisDlmCmp(getDlmConf()),
		)
	})

	return oprLogSvcImpl
}

func getDlmConf() (dlmConf *redisdlmcmp.RedisDlmCmpConf) {
	redisKeyPrefix := constants.RedisKeyPrefix + ":oprlog"

	dlmConf = dlmhelper.GetDefaultDlmConf(redisKeyPrefix)

	return
}
