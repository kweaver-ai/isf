package infra

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"AuditLog/common"
	"AuditLog/infra/config"
)

const (
	MasterSlaveType string = "master-slave"
	StandaloneType  string = "standalone"
	SentinelType    string = "sentinel"
	ClusterType     string = "cluster"
)

var (
	redisOnce   sync.Once
	redisClient redis.UniversalClient
)

func ConnectRedis(conf *config.RedisConfig) redis.UniversalClient {
	redisOnce.Do(func() {
		ctx := context.Background()

		switch conf.ConnectType {
		case MasterSlaveType:
			for {
				redisClient = masterSlave(conf)
				if err := redisClient.Ping(ctx).Err(); err != nil {
					time.Sleep(time.Duration(3) * time.Second)
				} else {
					break
				}
			}
		case StandaloneType:
			for {
				redisClient = standalone(conf)
				if err := redisClient.Ping(ctx).Err(); err != nil {
					common.SvcConfig.Logger.Errorf("redis connect failed, type: %v, err: %v", conf.ConnectType, err)
					time.Sleep(time.Duration(3) * time.Second)
				} else {
					break
				}
			}
		case SentinelType:
			for {
				redisClient = sentinel(conf)
				if err := redisClient.Ping(ctx).Err(); err != nil {
					common.SvcConfig.Logger.Errorf("redis connect failed, type: %v, err: %v", conf.ConnectType, err)
					time.Sleep(time.Duration(3) * time.Second)
				} else {
					break
				}
			}
		case ClusterType:
			for {
				redisClient = cluster(conf)
				if err := redisClient.Ping(ctx).Err(); err != nil {
					common.SvcConfig.Logger.Errorf("redis connect failed, type: %v, err: %v", conf.ConnectType, err)
					time.Sleep(time.Duration(3) * time.Second)
				} else {
					break
				}
			}
		}

		common.SvcConfig.Logger.Infof("redis connect success, type: %v", conf.ConnectType)
	})

	return redisClient
}

// masterSlave 主从模式
func masterSlave(conf *config.RedisConfig) *redis.Client {
	if conf.MasterHost == "" {
		conf.MasterHost = "proton-redis-proton-redis.resource.svc.cluster.local"
	}

	if conf.MasterPort == "" {
		conf.Port = "6379"
	}

	opt := &redis.Options{
		Addr:               conf.MasterHost + ":" + conf.MasterPort,
		Password:           conf.Password,
		DB:                 conf.DB,
		MaxRetries:         conf.MaxRetries,
		PoolSize:           conf.PoolSize,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(conf.IdleCheckFrequency) * time.Second,
		MaxConnAge:         time.Duration(conf.MaxConnAge) * time.Second,
		PoolTimeout:        time.Duration(conf.PoolTimeout) * time.Second,
	}

	return redis.NewClient(opt)
}

// standalone 标准模式客户端
func standalone(conf *config.RedisConfig) *redis.Client {
	if conf.Host == "" {
		conf.Host = "proton-redis-proton-redis.resource.svc.cluster.local"
	}

	if conf.Port == "" {
		conf.Port = "6379"
	}

	opt := &redis.Options{
		Addr:               conf.Host + ":" + conf.Port,
		Username:           conf.UserName,
		Password:           conf.Password,
		DB:                 conf.DB,
		MaxRetries:         conf.MaxRetries,
		PoolSize:           conf.PoolSize,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(conf.IdleCheckFrequency) * time.Second,
		MaxConnAge:         time.Duration(conf.MaxConnAge) * time.Second,
		PoolTimeout:        time.Duration(conf.PoolTimeout) * time.Second,
	}

	return redis.NewClient(opt)
}

// sentinel 哨兵模式客户端
func sentinel(conf *config.RedisConfig) *redis.Client {
	if conf.MasterGroupName == "" {
		conf.MasterGroupName = "mymaster"
	}

	if conf.SentinelPwd == "" {
		panic("conf.SentinelPwd is empty")
	}

	if conf.SentinelHost == "" {
		conf.SentinelHost = "proton-redis-proton-redis-sentinel.resource.svc.cluster.local"
	}

	if conf.SentinelPort == "" {
		conf.SentinelPort = "26379"
	}

	opt := redis.FailoverOptions{
		MasterName:         conf.MasterGroupName,
		SentinelAddrs:      []string{fmt.Sprintf("%v:%v", conf.SentinelHost, conf.SentinelPort)},
		SentinelPassword:   conf.SentinelPwd,
		Username:           conf.UserName,
		Password:           conf.Password,
		DB:                 conf.DB,
		MaxRetries:         conf.MaxRetries,
		PoolSize:           conf.PoolSize,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(conf.IdleCheckFrequency) * time.Second,
		MaxConnAge:         time.Duration(conf.MaxConnAge) * time.Second,
		PoolTimeout:        time.Duration(conf.PoolTimeout) * time.Second,
	}

	return redis.NewFailoverClient(&opt)
}

// cluster 集群模式客户端
func cluster(conf *config.RedisConfig) *redis.ClusterClient {
	hosts := strings.Split(conf.Host, ",")
	addrs := make([]string, 0, len(hosts))
	for _, host := range hosts {
		if strings.Contains(host, ":") {
			addrs = append(addrs, host)
		} else {
			addrs = append(addrs, fmt.Sprintf("%v:%v", host, conf.Port))
		}
	}
	opt := redis.ClusterOptions{
		Addrs:              addrs,
		Username:           conf.UserName,
		Password:           conf.Password,
		MaxRetries:         conf.MaxRetries,
		PoolSize:           conf.PoolSize,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Second,
		IdleCheckFrequency: time.Duration(conf.IdleCheckFrequency) * time.Second,
		MaxConnAge:         time.Duration(conf.MaxConnAge) * time.Second,
		PoolTimeout:        time.Duration(conf.PoolTimeout) * time.Second,
	}

	return redis.NewClusterClient(&opt)
}
