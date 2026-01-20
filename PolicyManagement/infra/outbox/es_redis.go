package outbox

import (
	"errors"
	"fmt"

	redisC "policy_mgnt/infra/redis"

	outboxer "policy_mgnt/utils/gocommon/v2/outbox"
)

const (
	Channel = "channel"
)

type RedisES struct {
	client redisC.RedisConn
}

type options struct {
	channel string
}

func NewRedisES() *RedisES {
	return &RedisES{
		client: redisC.NewRedisConn(),
	}
}

func (p *RedisES) Send(evt *outboxer.OutboxMessage) error {
	opts, err := p.parseOptions(evt.Options)
	if err != nil {
		return fmt.Errorf("parseOptions: %w", err)
	}

	err = p.client.Publish(opts.channel, evt.Payload)
	if err != nil {
		return fmt.Errorf("p.client.Publish: %w", err)
	}

	return nil
}

func (r *RedisES) parseOptions(opts map[string]interface{}) (*options, error) {
	opt := options{}

	if data, ok := opts[Channel]; ok {
		opt.channel, ok = data.(string)
		if !ok {
			return nil, fmt.Errorf("invalid channel: %v", opts[Channel])
		}
	}
	if opt.channel == "" {
		return nil, errors.New("channel is an empty string")
	}

	return &opt, nil
}
