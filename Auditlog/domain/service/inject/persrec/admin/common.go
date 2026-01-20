package persrecadmininject

import (
	"AuditLog/common/helpers/hlartrace"
	persrec_db "AuditLog/drivenadapters/db/persrec"
	"AuditLog/infra"
	"AuditLog/infra/cmp/logcmp"
)

func GetPersRepoBase() *persrec_db.RepoBase {
	logger := logcmp.GetLogger()
	arTracer := hlartrace.NewARTrace()

	return persrec_db.NewRepoBase(logger, infra.NewDBPool(), arTracer)
}
