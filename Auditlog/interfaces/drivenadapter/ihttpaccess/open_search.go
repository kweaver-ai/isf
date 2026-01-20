package ihttpaccess

import "context"

//go:generate mockgen -source=./open_search.go -destination ./httpaccmock/open_search_mock.go -package httpaccmock
type OpsHttpAcc interface {
	CreateInterfaceNoID(ctx context.Context, index string, doc interface{}) (err error)

	BatchCreate(ctx context.Context,
		index string, data []map[string]interface{}, isWithID bool) (err error)

	BatchCreateInterface(ctx context.Context, index string, docs interface{}, isWithID bool) (err error)

	CreateIndex(ctx context.Context, index string, mapping, setting string) (err error)

	DeleteIndex(ctx context.Context, index string) (err error)

	DeleteDocByField(ctx context.Context, index string, field string, value interface{}) (err error)

	DeleteDocsByFieldRange(ctx context.Context, index string, field string, from, to interface{}) (err error)
}
