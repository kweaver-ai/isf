package recdriveri

import (
	"context"
)

//go:generate mockgen -package mock -source rec.go -destination ../mock/rec.go
type IRecSvc interface {
	CreateByMapping(ctx context.Context) (err error)

	RemoveOldLog(ctx context.Context, index string, saveDays int) (err error)

	RemoveNotUseOpensearchIndexOnce(ctx context.Context) (err error)
}
