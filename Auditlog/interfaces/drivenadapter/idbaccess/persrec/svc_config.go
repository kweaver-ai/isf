package persrecrepoi

import (
	"context"
)

//go:generate mockgen -package persrecdbmock -source svc_config.go -destination ./persrecdbmock/svc_config.go
type IPersSvcConfigRepo interface {
	Set(ctx context.Context, key, val string) (err error)
	Get(ctx context.Context, key string) (val string, err error)
}
