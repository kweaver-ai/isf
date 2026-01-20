package ihttpaccess

import (
	"AuditLog/drivenadapters/httpaccess/document/docaccret"
)

//go:generate mockgen -source=./document.go -destination ./httpaccmock/document_mock.go -package httpaccmock
type DocumentHttpAcc interface {
	GetBatchDocLibInfos(docIDs []string) (infos []*docaccret.DocLibItem, err error)
}
