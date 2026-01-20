package dbaulid

import (
	"context"
	"fmt"

	uniqidenums "AuditLog/common/enums/uniqueid"
	"AuditLog/common/helpers/dbhelper"
	"AuditLog/common/utils"
	"AuditLog/models/persistence"
)

// GenUniqID 生成一个唯一ID
func (repo *ulidRepo) GenUniqID(ctx context.Context, flag uniqidenums.UniqueIDFlag) (id string, err error) {
	maxRetry := 5
	for i := 0; i < maxRetry; i++ {
		id, err = repo.genUniqID(ctx, flag)
		if err != nil {
			continue
		}

		if id != "" {
			break
		}
	}

	if id == "" {
		err = fmt.Errorf("[%s]: failed to generate unique id, err: %w", "GenUniqID", err)
	}

	return
}

//nolint:unparam
func (repo *ulidRepo) genUniqID(ctx context.Context, flag uniqidenums.UniqueIDFlag) (id string, err error) {
	_po := &persistence.UniqueID{}
	sr := dbhelper.NewSQLRunner(repo.db, repo.logger)

	id = utils.UlidMake()
	_po.ID = id
	_po.Flag = flag

	_, err = sr.FromPo(_po).InsertStruct(_po)
	if err != nil {
		id = ""
	}

	return
}

func (repo *ulidRepo) DelUniqID(ctx context.Context, flag uniqidenums.UniqueIDFlag, id string) (err error) {
	_po := &persistence.UniqueID{}
	sr := dbhelper.NewSQLRunner(repo.db, repo.logger)

	_, err = sr.FromPo(_po).
		WhereEqual("f_id", id).
		WhereEqual("f_flag", flag).
		Delete()

	return
}
