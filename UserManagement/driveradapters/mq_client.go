// Package driveradapters MQ客户端
package driveradapters

import (
	"sync"

	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"

	"UserManagement/common"
)

// MQClient MQ客户端接口
type MQClient interface {
	// Subscribe 订阅mq消息
	Subscribe(topic, channel string, cmd func([]byte) error)
}

var (
	mqOnce   sync.Once
	mqClient MQClient
)

type msgQueue struct {
	log                      common.Logger
	client                   msqclient.ProtonMQClient
	pollIntervalMilliseconds int64
	maxInFlight              int
}

// NewMQClient 创建消息队列
func NewMQClient() MQClient {
	mqOnce.Do(func() {
		client, err := msqclient.NewProtonMQClientFromFile("/sysvol/conf/mq_config.yaml")
		if err != nil {
			common.NewLogger().Fatalln(err)
		}
		mqClient = &msgQueue{
			log:                      common.NewLogger(),
			client:                   client,
			pollIntervalMilliseconds: 100,
			maxInFlight:              200,
		}
	})

	return mqClient
}

// Subscribe 服务订阅
func (m *msgQueue) Subscribe(topic, channel string, cmd func([]byte) error) {
	go func() {
		err := m.client.Sub(topic, channel, cmd, m.pollIntervalMilliseconds, m.maxInFlight)
		m.log.Errorln(err)
	}()
}
