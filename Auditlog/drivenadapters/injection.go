/* 出栈注入层 */
package drivenadapters

import (
	"github.com/kweaver-ai/proton-rds-sdk-go/sqlx"
	"github.com/go-redis/redis/v8"

	"AuditLog/gocommon/api"
)

var (
	DBPool      *sqlx.DB
	Logger      api.Logger
	HTTPClient  api.Client
	RedisClient redis.Cmdable
)

func SetDBPool(i *sqlx.DB) {
	DBPool = i
}

func SetLogger(i api.Logger) {
	Logger = i
}

func SetHTTPClient(i api.Client) {
	HTTPClient = i
}

func SetRedisClient(i redis.Cmdable) {
	RedisClient = i
}
