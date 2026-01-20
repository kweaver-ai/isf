// Package drivenadapters 适配器eacp
package drivenadapters

import (
	"context"
	"fmt"
	"sync"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/observable"

	"Authentication/common"
	"Authentication/interfaces"
)

var (
	eOnce sync.Once
	e     *eacp
)

type eacp struct {
	httpClient  httpclient.HTTPClient
	log         common.Logger
	privateAddr string
	trace       observable.Tracer
}

// NewEacp 创建eacp接口操作对象
func NewEacp() *eacp {
	config := common.SvcConfig
	eOnce.Do(func() {
		e = &eacp{
			httpClient:  httpclient.NewHTTPClient(common.SvcARTrace),
			log:         common.NewLogger(),
			privateAddr: fmt.Sprintf("http://%s:%d", config.EacpPrivateHost, config.EacpPrivatePort),
			trace:       common.SvcARTrace,
		}
	})

	return e
}

func (s *eacp) ThirdPartyAuthentication(ctx context.Context, visitor *interfaces.Visitor, req *interfaces.ThirdPartyAuthInfo) (*interfaces.LoginInfo, error) {
	var err error
	s.trace.SetClientSpanName("适配器-第三方认证")
	newCtx, span := s.trace.AddClientTrace(ctx)
	defer func() { s.trace.TelemetrySpanEnd(span, err) }()

	permInfo := map[string]interface{}{
		"thirdpartyid": req.Credential.ID,
		"params":       req.Credential.Params,
		"device": map[string]interface{}{
			"name":        req.Device.Name,
			"client_type": req.Device.ClientType,
			"description": req.Device.Description,
			"udids":       req.Udids,
		},
		"ip": req.IP,
	}

	target := fmt.Sprintf("%v/api/eacp/v1/auth1/getbythirdparty", s.privateAddr)

	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}

	_, resParam, err := s.httpClient.Post(newCtx, target, headers, permInfo)
	if err != nil {
		s.log.Errorf("ThirdParty authentication failed: %v, url: %v", err, target)
		return nil, err
	}

	info := &interfaces.LoginInfo{
		Subject: resParam.(map[string]interface{})["user_id"].(string),
		Context: resParam.(map[string]interface{})["context"],
	}
	return info, nil
}
