package redishelper

import (
	"context"
	"errors"
	"time"

	redis "github.com/go-redis/redis/v8"

	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/infra"
	"AuditLog/infra/config"

	"golang.org/x/sync/singleflight"
)

var ErrNotSupportInLocalEnv = errors.New("redishelper: not support in local env")

func SetStruct(ctx context.Context, rdb redis.Cmdable, key string, value interface{}, ttl time.Duration) error {
	if helpers.IsLocalDev() {
		return nil
	}

	jsonStr, err := utils.JSON().MarshalToString(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, jsonStr, ttl).Err()
}

var sfgGetStruct singleflight.Group

func GetStruct(ctx context.Context, rdb redis.Cmdable, key string, value interface{}) error {
	// 可根据需要打开或关闭
	if helpers.IsLocalDev() {
		return ErrNotSupportInLocalEnv
	}

	jsonInter, err, _ := sfgGetStruct.Do(key, func() (interface{}, error) {
		return rdb.Get(ctx, key).Result()
	})

	if err != nil {
		return err
	}

	//nolint:forcetypeassert
	return utils.JSON().UnmarshalFromString(jsonInter.(string), value)
}

func GetRedisClient() (redisClient redis.UniversalClient) {
	redisClient = infra.ConnectRedis(config.GetRedisConfig())
	return
}
