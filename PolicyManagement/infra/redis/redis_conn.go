package redis

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"

	"policy_mgnt/common/config"

	clog "policy_mgnt/utils/gocommon/v2/log"
)

//go:generate mockgen -source=./redis_conn.go -package=mock_logics -destination=../../test/mock_logics/redis_conn_mock.go
const (
	// RedisTypeSentinel redis哨兵模式
	RedisTypeSentinel = "sentinel"
	// RedisTypeStandalone redis单机模式
	RedisTypeStandalone = "standalone"
	// RedisTypeMasterSlave redis主从模式
	RedisTypeMasterSlave = "master-slave"
	// RedisTypeCluster redis集群模式
	RedisTypeCluster = "cluster"
)

var (
	rdConnOnce sync.Once
	rdConn     *redisConn
)

type redisConn struct {
	writeCli   *redis.Client // 负责写操作
	readCli    *redis.Client // 负责读操作
	logger     clog.Logger
	clusterCli *redis.ClusterClient // 集群客户端
	isCluster  bool                 // 是否为集群模式
}

// RedisConn redis客户端接口
type RedisConn interface {
	// Subscribe 消息订阅
	Subscribe(channel string, cmd func([]byte)) error

	// Publish 消息发布
	Publish(channel string, message []byte) error
}

// NewRedisConn 创建Redis连接对象
func NewRedisConn() RedisConn {
	rdConnOnce.Do(func() {
		redisConfig := config.Config.Redis
		rdConn = &redisConn{
			logger: clog.NewLogger(),
		}

		switch redisConfig.ConnectType {
		case RedisTypeSentinel:
			rdConn.writeCli = redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:       redisConfig.ConnectInfo.MasterGroupName,
				SentinelAddrs:    []string{fmt.Sprintf("%s:%d", redisConfig.ConnectInfo.SentinelHost, redisConfig.ConnectInfo.SentinelPort)},
				Username:         redisConfig.ConnectInfo.Username,
				Password:         redisConfig.ConnectInfo.Password,
				SentinelUsername: redisConfig.ConnectInfo.SentinelUsername,
				SentinelPassword: redisConfig.ConnectInfo.SentinelPassword,
			})
			rdConn.readCli = rdConn.writeCli
		case RedisTypeStandalone:
			rdConn.writeCli = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisConfig.ConnectInfo.Host, redisConfig.ConnectInfo.Port),
				Username: redisConfig.ConnectInfo.Username,
				Password: redisConfig.ConnectInfo.Password,
			})
			rdConn.readCli = rdConn.writeCli
		case RedisTypeMasterSlave:
			// 华为云环境下，尝试区分主从节点
			rdConn.writeCli = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisConfig.ConnectInfo.MasterHost, redisConfig.ConnectInfo.MasterPort),
				Username: redisConfig.ConnectInfo.Username,
				Password: redisConfig.ConnectInfo.Password,
			})
			if redisConfig.ConnectInfo.SlaveHost != "" {
				rdConn.readCli = redis.NewClient(&redis.Options{
					Addr:     fmt.Sprintf("%s:%d", redisConfig.ConnectInfo.SlaveHost, redisConfig.ConnectInfo.SlavePort),
					Username: redisConfig.ConnectInfo.Username,
					Password: redisConfig.ConnectInfo.Password,
				})
			} else {
				rdConn.readCli = rdConn.writeCli
			}
		case RedisTypeCluster:
			clusterAddrs := rdConn.parseClusterAddrs(redisConfig.ConnectInfo.Host, redisConfig.ConnectInfo.Port)

			clusterOptions := &redis.ClusterOptions{Addrs: clusterAddrs}
			if redisConfig.ConnectInfo.Password != "" {
				clusterOptions.Password = redisConfig.ConnectInfo.Password
				if redisConfig.ConnectInfo.Username != "" {
					clusterOptions.Username = redisConfig.ConnectInfo.Username
				}
			}

			rdConn.clusterCli = redis.NewClusterClient(clusterOptions)
			rdConn.isCluster = true
		default:
			rdConn.logger.Fatalf("redis connect type should be one of %s, %s, %s, %s", RedisTypeSentinel, RedisTypeMasterSlave, RedisTypeStandalone, RedisTypeCluster)
		}

		// 连接测试
		pong := "PONG"
		if rdConn.isCluster {
			s, err := rdConn.clusterCli.Ping(context.Background()).Result()
			if err != nil || s != pong {
				rdConn.logger.Fatalf("create redis cluster client failed:%v", err)
			}
		} else {
			s, err := rdConn.writeCli.Ping(context.Background()).Result()
			if err != nil || s != pong {
				rdConn.logger.Fatalf("create redis write client failed:%v", err)
			}

			s, err = rdConn.readCli.Ping(context.Background()).Result()
			if err != nil || s != pong {
				rdConn.logger.Fatalf("create redis read client failed:%v", err)
			}
		}
	})

	return rdConn
}

// Subscribe 消息订阅
func (r *redisConn) Subscribe(channel string, cmd func([]byte)) error {
	var sub *redis.PubSub

	if r.isCluster {
		sub = r.clusterCli.Subscribe(context.Background(), channel)
	} else {
		sub = r.readCli.Subscribe(context.Background(), channel)
	}

	_, err := sub.Receive(context.Background())
	if err != nil {
		return fmt.Errorf("Receive: %w", err)
	}
	ch := sub.Channel()

	// TODO 需要errGroup吗
	go func() {
		// 获取消息推送请求信息
		for msg := range ch {
			// 解析推送信息
			cmd([]byte(msg.Payload))
		}
	}()

	return nil
}

// Publish 消息发布
func (r *redisConn) Publish(channel string, message []byte) error {
	var err error
	if r.isCluster {
		_, err = r.clusterCli.Publish(context.Background(), channel, string(message)).Result()
	} else {
		_, err = r.writeCli.Publish(context.Background(), channel, string(message)).Result()
	}

	if err != nil {
		r.logger.Errorf("redis Publish %s err:%v", channel, err)
	}
	return err
}

// parseClusterAddrs 解析集群地址
func (r *redisConn) parseClusterAddrs(host string, defaultPort int) []string {
	if host == "" {
		r.logger.Fatalf("cluster host cannot be empty")
		return nil
	}

	hostList := strings.Split(host, ",")
	var clusterAddrs []string

	for _, addr := range hostList {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}

		if strings.Contains(addr, ":") {
			clusterAddrs = append(clusterAddrs, addr)
		} else {
			clusterAddrs = append(clusterAddrs, fmt.Sprintf("%s:%d", addr, defaultPort))
		}
	}

	if len(clusterAddrs) == 0 {
		r.logger.Fatalf("no valid cluster addresses found in host: %s", host)
	}

	return clusterAddrs
}
