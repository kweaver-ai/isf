package redisaccess

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"AuditLog/drivenadapters"
	"AuditLog/interfaces"
)

var (
	dlmOnce sync.Once
	dlmLock *dlm
)

type dlm struct {
	redisClient redis.Cmdable
	expiration  time.Duration
	contexts    map[string]context.CancelFunc
	uniqueID    string
}

// NewDLM 获取DLM实例
func NewDLM() interfaces.DLM {
	dlmOnce.Do(func() {
		dlmLock = &dlm{
			redisClient: drivenadapters.RedisClient,
			expiration:  time.Second * 10,
			contexts:    make(map[string]context.CancelFunc),
			uniqueID:    uuid.New().String(),
		}
	})
	return dlmLock
}

// TryLock 获取锁
func (d *dlm) TryLock(key string) (success bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	success, err = d.redisClient.SetNX(ctx, key, d.uniqueID, d.expiration).Result()
	if err != nil {
		cancel()
		return false, fmt.Errorf("[redisaccess.DLM.TryLock]->[RedisClient.SetNX]->%w", err)
	}

	if !success {
		cancel()
		return false, nil
	}

	d.contexts[key] = cancel
	go d.keepLock(ctx, key)
	return true, nil
}

// UnLock 解锁，使用Lua脚本保证原子性
func (d *dlm) UnLock(key string) error {
	// Lua脚本，保证原子性
	const unlockScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        end
        return 0
    `

	ctx := context.TODO()
	result, err := d.redisClient.Eval(ctx, unlockScript, []string{key}, d.uniqueID).Result()
	if err != nil {
		return fmt.Errorf("[redisaccess.DLM.UnLock]->[RedisClient.Eval]->%w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("锁不存在或已被其他客户端持有")
	}

	if cancel, ok := d.contexts[key]; ok {
		cancel()
		delete(d.contexts, key)
	}
	return nil
}

// keepLock 优化锁续期逻辑
func (d *dlm) keepLock(ctx context.Context, key string) {
	ticker := time.NewTicker(d.expiration / 2)
	defer ticker.Stop()

	// 刷新锁的过期时间
	const refreshScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("expire", KEYS[1], ARGV[2])
        end
        return 0
    `

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result, err := d.redisClient.Eval(
				ctx,
				refreshScript,
				[]string{key},
				d.uniqueID,
				int(d.expiration.Seconds()),
			).Result()
			if err != nil {
				fmt.Printf("[redisaccess.DLM.keepLock]->[RedisClient.Eval]->%v\n", err)
				continue
			}

			if result.(int64) == 0 {
				// 锁已经不属于我们了，退出续期
				if cancel, ok := d.contexts[key]; ok {
					cancel()
					delete(d.contexts, key)
				}
				return
			}
		}
	}
}
