// Package drivenadapters 消息队列
package drivenadapters

import (
	"sync"

	jsoniter "github.com/json-iterator/go"

	"Authentication/common"
	"Authentication/interfaces"
)

type messageBroker struct {
	log         common.Logger
	redisClient interfaces.RedisConn
}

var (
	msgOnce sync.Once
	m       *messageBroker
)

// NewMessageBroker 创建消息发送对象
func NewMessageBroker() *messageBroker {
	msgOnce.Do(func() {
		m = &messageBroker{
			log:         common.NewLogger(),
			redisClient: common.NewRedisConn(),
		}
	})
	return m
}

// AnonymousSmsExpUpdated 更新匿名登录短信验证码过期时间
func (m *messageBroker) AnonymousSmsExpUpdated(smsExpiration int) error {
	payload := map[string]interface{}{
		"anonymous_sms_expiration": smsExpiration,
	}
	payloadBytes, err := jsoniter.Marshal(payload)
	if err != nil {
		m.log.Errorln("failed to marshal anonymous sms expiration config, err:", err)
		return err
	}
	return m.redisClient.Publish("authentication.config.anonymous_sms_expiration.updated", payloadBytes)
}
