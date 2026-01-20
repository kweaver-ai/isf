package opensearchcmp

import (
	"github.com/opensearch-project/opensearch-go"

	"AuditLog/gocommon/api"
	"AuditLog/infra/cmp/icmp"
)

type OpsCmp struct {
	address  string
	username string
	password string

	arTrace api.Tracer
	logger  api.Logger

	client *opensearch.Client
}

type OpsCmpConf struct {
	Address  string
	Username string
	Password string

	ArTrace api.Tracer
	Logger  api.Logger
}

var _ icmp.IOpsCmp = &OpsCmp{}

func NewOpsCmp(conf *OpsCmpConf) (cmp icmp.IOpsCmp, err error) {
	o := &OpsCmp{
		address:  conf.Address,
		username: conf.Username,
		password: conf.Password,

		arTrace: conf.ArTrace,
		logger:  conf.Logger,
	}

	err = o.newClient()
	if err != nil {
		return
	}

	cmp = o

	return
}
