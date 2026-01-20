// Package tclient thrift客户端创建函数
package tclient

import (
	"reflect"
	"strconv"

	"github.com/kweaver-ai/go-lib/thrift"
)

// NewTClient 创建thrift客户端
func NewTClient(newFunc, clientPtrPtr interface{}, ip string, port int) (transport thrift.TTransport, err error) {
	socket, err := thrift.NewTSocket(ip + ":" + strconv.Itoa(port))
	if err != nil {
		return
	}
	transport = thrift.NewTBufferedTransport(socket, 8192)
	err = transport.Open()
	if err != nil {
		return
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	argsV := make([]reflect.Value, 2)
	argsV[0] = reflect.ValueOf(transport)
	argsV[1] = reflect.ValueOf(protocolFactory)
	newFuncV := reflect.ValueOf(newFunc)
	clientV := newFuncV.Call(argsV)
	reflect.ValueOf(clientPtrPtr).Elem().Set(clientV[0])
	return
}
