package ihttpaccess

import (
	"context"

	"AuditLog/common/enums"
)

//go:generate mockgen -source=./efast.go -destination ./httpaccmock/efast_mock.go -package httpaccmock
type EFastHttpAcc interface {
	GetDocLibType(ctx context.Context, gns string) (docLibType enums.DocLibType, err error)

	CheckObjExists(ctx context.Context, ids []string) (notExistsIDs []string, err error)

	CheckOneObjExists(ctx context.Context, id string) (exist bool, err error)

	Path2Gns(ctx context.Context, path, asToken string) (isNotExists bool, gns string, err error)

	CreateMultiLevelDir(ctx context.Context, parentDocID, path, token string) (dirDocID string, err error)

	GetOneFsName(ctx context.Context, docID string) (name string, err error)

	Gns2Path(ctx context.Context, gns []string) (pathMap map[string]string, err error)
}
