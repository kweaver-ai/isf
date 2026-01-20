package httpinject

import (
	"sync"

	"AuditLog/common/helpers/hlartrace"
	"AuditLog/drivenadapters/httpaccess/efasthttp"
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

var (
	eFastOnce    sync.Once
	eFastImpl    ihttpaccess.EFastHttpAcc
	eFastArTrace api.Tracer
)

func NewEFastHttpAcc() ihttpaccess.EFastHttpAcc {
	eFastOnce.Do(func() {
		// 1. arTrace
		eFastArTrace = hlartrace.NewARTrace()

		// 2. um configuration
		efConf := config.GetEFastConfig()

		// 3. efast
		eFastImpl = efasthttp.NewEFastHttpAcc(
			eFastArTrace,
			logcmp.GetLogger(),
			efConf,
		)
	})

	eFastArTrace.SetInternalSpanName("efast")

	return eFastImpl
}
