package apiserver

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	redisC "policy_mgnt/infra/redis"
)

var (
	registerOnce      sync.Once
	registerSingleton *Register
)

type RedisSubscriber struct {
	Channel string
	Handler func([]byte)
}

type Register struct {
	errGroup *errgroup.Group
	ctx      context.Context
	redis    redisC.RedisConn
}

type RedisSub interface {
	RegisterRedis() []RedisSubscriber
}

func NewRegister() *Register {
	registerOnce.Do(func() {
		g, ctx := errgroup.WithContext(context.Background())
		registerSingleton = &Register{
			errGroup: g,
			ctx:      ctx,
			redis:    redisC.NewRedisConn(),
		}
	})
	return registerSingleton
}

// SubscribeRedis 订阅Redis
func (r *Register) SubscribeRedis(redisSubs []RedisSub) {
	for _, redisSub := range redisSubs {
		for _, subscriber := range redisSub.RegisterRedis() {
			func(sub RedisSubscriber) {
				r.errGroup.Go(func() error {
					redisErrChan := make(chan error, 1)
					// defer close(redisErrChan)
					go func() {
						redisErrChan <- r.redis.Subscribe(sub.Channel, sub.Handler)
					}()
					select {
					case <-r.ctx.Done():
						return r.ctx.Err()
					case err := <-redisErrChan:
						return err
					}
				})
			}(subscriber)
		}
	}
}
