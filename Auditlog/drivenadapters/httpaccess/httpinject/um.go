package httpinject

import (
	"sync"

	"AuditLog/common/helpers/hlartrace"
	"AuditLog/drivenadapters/httpaccess/umhttpaccess"
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

var (
	umOnce    sync.Once
	umImpl    ihttpaccess.UmHttpAcc
	umArTrace api.Tracer
)

func NewUmHttpAcc() ihttpaccess.UmHttpAcc {
	umOnce.Do(func() {
		// 1. arTrace
		umArTrace = hlartrace.NewARTrace()

		// 2. um configuration
		umConf := config.GetUserManagementConfig()

		// 3. um
		umImpl = umhttpaccess.NewUmHttpAcc(
			umArTrace,
			logcmp.GetLogger(),
			umConf,
		)
	})

	umArTrace.SetInternalSpanName("user-management")

	return umImpl
}
