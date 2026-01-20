package common

import (
	"fmt"

	sdk "github.com/kweaver-ai/proton-mq-sdk-go"

	"Authorization/interfaces"
)

// NewMQClient 获取消息队列客户端
func NewMQClient() (interfaces.MQClient, error) {
	c, err := sdk.NewProtonMQClientFromFile(SvcConfig.MQConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("sdk.NewProtonMQClientFromFile: %w", err)
	}
	return c, nil
}
