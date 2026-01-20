package drivenadapters

import (
	"context"
	"sync"

	"github.com/kweaver-ai/go-lib/observable"
	"github.com/kweaver-ai/go-lib/rest"
	"github.com/kweaver-ai/go-lib/tclient"

	"Authentication/common"
	"Authentication/tapi/ethriftexception"
	"Authentication/tapi/sharemgnt"
)

var (
	sOnce sync.Once
	s     *shareMgnt
)

type shareMgnt struct {
	logger common.Logger
	trace  observable.Tracer
}

// NewShareMgnt 创建shareMgnt接口操作对象
func NewShareMgnt() *shareMgnt {
	sOnce.Do(func() {
		s = &shareMgnt{
			logger: common.NewLogger(),
			trace:  common.SvcARTrace,
		}
	})

	return s
}

// UsrmSMSValidate 校验短信验证码
func (s *shareMgnt) UsrmSMSValidate(userID, vcode string) (err error) {
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	err = shareMgntClient.Usrm_SMSValidate(context.Background(), userID, vcode)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.UsrmSMSValidate error: %v", err)
		return
	}

	return
}

// UsrmOTPValidate 校验动态密码
func (s *shareMgnt) UsrmOTPValidate(userID, otp string) (err error) {
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	err = shareMgntClient.Usrm_OTPValidate(context.Background(), userID, otp)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.UsrmOTPValidate error: %v", err)
		return
	}

	return
}

// UsrmIMAGECodeValidate 校验图形验证码
func (s *shareMgnt) UsrmIMAGECodeValidate(uuid, vcode string) (err error) {
	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	err = shareMgntClient.Usrm_IMAGECodeValidate(context.Background(), uuid, vcode)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.UsrmIMAGECodeValidate error: %v", err)
		return
	}

	return
}

// UsrmDomainAuth 域认证
func (s *shareMgnt) UsrmDomainAuth(ctx context.Context, loginName, domainPath, password string, ldapType int32) (result bool, err error) {
	s.trace.SetClientSpanName("出栈适配器层-域认证")
	newCtx, span := s.trace.AddClientTrace(ctx)
	defer func() { s.trace.TelemetrySpanEnd(span, err) }()

	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	result, err = shareMgntClient.Usrm_DomainAuth(newCtx, loginName, ldapType, domainPath, password)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.UsrmDomainAuth error: %v", err)
		return
	}

	return
}

// UsrmThirdAuth 第三方认证
func (s *shareMgnt) UsrmThirdAuth(ctx context.Context, loginName, password string) (result bool, err error) {
	s.trace.SetClientSpanName("出栈适配器层-第三方认证")
	newCtx, span := s.trace.AddClientTrace(ctx)
	defer func() { s.trace.TelemetrySpanEnd(span, err) }()

	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	result, err = shareMgntClient.Usrm_ThirdAuth(newCtx, loginName, password)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.UsrmThirdAuth error: %v", err)
		return
	}

	return
}

// UsrmSendAnonymousSMSVCode 发送匿名账户短信验证码
func (s *shareMgnt) UsrmSendAnonymousSMSVCode(ctx context.Context, phoneNumber, vcode string) (err error) {
	s.trace.SetClientSpanName("出栈适配器层-发送匿名账户短信验证码")
	newCtx, span := s.trace.AddClientTrace(ctx)
	defer func() { s.trace.TelemetrySpanEnd(span, err) }()

	var shareMgntClient *sharemgnt.NcTShareMgntClient
	transport, err := tclient.NewTClient(sharemgnt.NewNcTShareMgntClientFactory, &shareMgntClient, common.SvcConfig.ShareMgntHost, common.SvcConfig.ShareMgntPort)
	if err != nil {
		s.logger.Errorf("Create shareMgntClient error: %v", err)
		return err
	}
	defer func() {
		if closeErr := transport.Close(); closeErr != nil {
			s.logger.Errorln(closeErr)
		}
	}()

	err = shareMgntClient.Usrm_SendSMSVCode(newCtx, phoneNumber, vcode)
	if err != nil {
		s.logger.Errorf("Call shareMgntClient.Usrm_SendSMSVCode error: %v", err)
		return convertThriftErrToRest(err)
	}

	return nil
}

func convertThriftErrToRest(err error) (restError error) {
	switch v := err.(type) {
	case *ethriftexception.NcTException:
		restError = rest.NewHTTPErrorV2(rest.InternalServerError, v.GetExpMsg())
	default:
		restError = rest.NewHTTPErrorV2(rest.InternalServerError, err.Error())
	}

	return
}
