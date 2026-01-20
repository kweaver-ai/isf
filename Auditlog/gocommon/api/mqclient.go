package api

import (
	"context"
	"fmt"
	"sync"

	msqclient "github.com/kweaver-ai/proton-mq-sdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	mqOnce     sync.Once
	mqClient   MQClient
	configPath string = "/sysvol/conf/mq_config.yaml"
)

type MQClient interface {
	Publish(string, []byte) (err error)
	Subscribe(topic, channel string, cmd func([]byte) error)
}

// mqclient 生产者
type mqclient struct {
	log                      Logger
	prontonMQClient          msqclient.ProtonMQClient
	trace                    Tracer
	pollIntervalMilliseconds int64
	maxInFlight              int
}

// NewMQClient 创建消息队列
func NewMQClient() MQClient {
	mqOnce.Do(func() {
		mqSDK, err := msqclient.NewProtonMQClientFromFile(configPath)
		if err != nil {
			panic(fmt.Sprintf("ERROR: new mq client failed: %v\n", err))
		}
		mqClient = &mqclient{
			log:                      NewTelemetryLogger(),
			prontonMQClient:          mqSDK,
			trace:                    NewARTrace(),
			pollIntervalMilliseconds: int64(100),
			maxInFlight:              16,
		}
	})

	return mqClient
}

// Publish mq生产者
func (m *mqclient) Publish(topic string, msg []byte) (err error) {
	err = m.prontonMQClient.Pub(topic, msg)
	if err != nil {
		m.log.Errorln(err)
	}
	return
}

// Subscribe mq消费者
func (m *mqclient) Subscribe(topic, channel string, cmd func([]byte) error) {
	ctx := context.Background()
	go func(ctx context.Context, topic, channel string, cmd func([]byte) error) {
		var span trace.Span
		var err error
		var data []byte
		ctx, span = m.trace.AddConsumerTrace(ctx, topic)
		defer func() {
			span.SetAttributes(attribute.String("msg", string(data)))
			m.trace.TelemetrySpanEnd(span, err)
		}()
		err = m.prontonMQClient.Sub(topic, channel, func(data []byte) error {
			return cmd(data)
		}, m.pollIntervalMilliseconds, m.maxInFlight)
		if err != nil {
			m.log.Errorln(err)
		}
	}(ctx, topic, channel, cmd)
}
