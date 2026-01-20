package httpinject

import (
	"AuditLog/common/helpers/hlartrace"
	"AuditLog/drivenadapters/httpaccess/opshttp"
	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/logcmp"
	"AuditLog/infra/config"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

var (
	opsImpl    ihttpaccess.OpsHttpAcc
	opsArTrace api.Tracer
)

func NewOpsHttpAcc() (ihttpaccess.OpsHttpAcc, error) {
	defer func() {
		if opsArTrace != nil {
			opsArTrace.SetInternalSpanName("open-search")
		}
	}()

	if opsImpl != nil {
		return opsImpl, nil
	}

	// ops configuration
	conf := config.GetConfig().OpenSearch

	// ar trace
	opsArTrace = hlartrace.NewARTrace()

	// ops
	_opsImpl, err := opshttp.NewOpsHttpAcc(
		opsArTrace,
		logcmp.GetLogger(),
		conf,
	)
	if err != nil {
		return nil, err
	}

	opsImpl = _opsImpl

	return opsImpl, nil
}
