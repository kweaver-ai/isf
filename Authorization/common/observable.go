package common

import (
	"github.com/kweaver-ai/go-lib/observable"
)

// SvcARTrace ARTrace实例
var SvcARTrace observable.Tracer

// InitARTrace 初始化ARTrace实例
func InitARTrace(serviceName string) {
	SvcARTrace = observable.NewARTrace(serviceName)
}
