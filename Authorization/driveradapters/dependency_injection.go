package driveradapters

import (
	"Authorization/interfaces"
)

var mqClient interfaces.MQClient

// SetMQClient 设置实例
func SetMQClient(i interfaces.MQClient) {
	mqClient = i
}
