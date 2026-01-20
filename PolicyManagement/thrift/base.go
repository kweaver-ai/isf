package thrift

import (
	"fmt"
	"os"

	"policy_mgnt/tapi/sharemgnt"

	"github.com/kweaver-ai/go-lib/thrift"
)

type AShareMgntClient struct {
	client    *sharemgnt.NcTShareMgntClient
	transport *thrift.TBufferedTransport
}

func newShareMgntClient() (aShareMgntClient *AShareMgntClient, err error) {
	var transport thrift.TTransport
	hostPort := fmt.Sprintf("%s:%d", os.Getenv("SHAREMGNT_THRIFT_HOST"), sharemgnt.NCT_SHAREMGNT_PORT)
	transport, err = thrift.NewTSocket(hostPort)
	if err != nil {
		return
	}
	transport = thrift.NewTBufferedTransport(transport, 8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := sharemgnt.NewNcTShareMgntClientFactory(transport, protocolFactory)
	return &AShareMgntClient{client: client, transport: transport.(*thrift.TBufferedTransport)}, nil
}
