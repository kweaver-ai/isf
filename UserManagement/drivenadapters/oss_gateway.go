// Package drivenadapters 消息队列
package drivenadapters

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/observable"

	"UserManagement/common"
	"UserManagement/interfaces"
)

type ossGateWay struct {
	log          common.Logger
	rawClient    *http.Client
	client       httpclient.HTTPClient
	adminAddress string
	trace        observable.Tracer
}

var (
	ossOnce sync.Once
	oss     *ossGateWay
)

// ossRequest OSS请求
type ossRequest struct {
	method  string
	url     string
	headers map[string]string
}

// NewOSSGateWay 创建OSS网关服务
func NewOSSGateWay() *ossGateWay {
	ossOnce.Do(func() {
		config := common.SvcConfig
		oss = &ossGateWay{
			log:          common.NewLogger(),
			rawClient:    httpclient.NewRawHTTPClient(),
			client:       httpclient.NewHTTPClient(common.SvcARTrace),
			adminAddress: fmt.Sprintf("http://%s:%d/api/ossgateway", config.OSSGateWayPrivateHost, config.OSSGateWayPrivatePort),
			trace:        common.SvcARTrace,
		}
	})

	return oss
}

// GetDownloadURL 获取下载文件URL
func (o *ossGateWay) GetDownloadURL(ctx context.Context, visitor *interfaces.Visitor, ossID, key string) (strURL string, err error) {
	o.trace.SetClientSpanName("适配器-获取下载文件URL")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	url := fmt.Sprintf("%s/v1/download/%s/%s?internal_request=falset&type=query_string", o.adminAddress, ossID, key)
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	respParam, err := o.client.Get(newCtx, url, headers)
	if err != nil {
		o.log.Errorf("oss_gateway GetDownloadURL err: %v", err)
		return
	}

	// 解析数据
	info := o.convertOSSResponse(respParam)
	return info.url, nil
}

// UploadFile 上传文件
func (o *ossGateWay) UploadFile(ctx context.Context, visitor *interfaces.Visitor, ossID, key string, data []byte) (err error) {
	o.trace.SetClientSpanName("适配器-上传文件")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	// 获取上传URL
	url := fmt.Sprintf("%s/v1/upload/%s/%s?request_method=PUT&internal_request=true", o.adminAddress, ossID, key)
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	respParam, err := o.client.Get(newCtx, url, headers)
	if err != nil {
		o.log.Errorf("oss_gateway UploadFile get upload url error: %v", err)
		return
	}

	// 解析数据
	requestInfo := o.convertOSSResponse(respParam)

	// 开始上传
	req, err := http.NewRequest(requestInfo.method, requestInfo.url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	for k, v := range requestInfo.headers {
		if v != "" {
			req.Header.Add(k, v)
		}
	}

	resp, err := o.rawClient.Do(req)
	if err != nil {
		o.log.Errorf("UploadFile error: %s, requestInfo:%+v", err.Error(), requestInfo)
		return
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			o.log.Errorln(closeErr)
		}
	}()

	// 如果接口存在问题，直接返回错误
	body, err := io.ReadAll(resp.Body)
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		err = fmt.Errorf("code:%v,header:%v,body:%v", resp.StatusCode, resp.Header, string(body))
		return
	}
	return
}

// DeleteFile 删除文件
func (o *ossGateWay) DeleteFile(ctx context.Context, visitor *interfaces.Visitor, ossID, key string) (err error) {
	o.trace.SetClientSpanName("适配器-删除文件")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	// 获取删除文件URL
	url := fmt.Sprintf("%s/v1/delete/%s/%s?internal_request=true", o.adminAddress, ossID, key)
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	respParam, err := o.client.Get(newCtx, url, headers)
	if err != nil {
		o.log.Errorf("oss_gateway DeleteFile get download url error: %v", err)
		return
	}

	// 解析数据
	info := o.convertOSSResponse(respParam)

	// 删除文件
	req, err := http.NewRequest(info.method, info.url, http.NoBody)
	if err != nil {
		return
	}

	for k, v := range info.headers {
		if v != "" {
			req.Header.Add(k, v)
		}
	}
	resp, err := o.rawClient.Do(req)
	if err != nil {
		o.log.Errorf("DeleteFile error: %s, requestInfo:%+v", err.Error(), info)
		return
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			o.log.Errorln(closeErr)
		}
	}()

	// 如果接口存在问题，直接返回错误
	body, err := io.ReadAll(resp.Body)
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		err = fmt.Errorf("code:%v,header:%v,body:%v", resp.StatusCode, resp.Header, string(body))
		return
	}
	return
}

// convertOSSResponse oss response解析
func (o *ossGateWay) convertOSSResponse(body interface{}) (info ossRequest) {
	temp := body.(map[string]interface{})

	info.method = temp["method"].(string)
	info.url = temp["url"].(string)

	tempMap := temp["headers"].(map[string]interface{})
	info.headers = make(map[string]string)
	for k, v := range tempMap {
		info.headers[k] = v.(string)
	}

	return
}

// GetLocalEnabledOSSInfo 获取本地站点下可用的存储信息
func (o *ossGateWay) GetLocalEnabledOSSInfo(ctx context.Context, visitor *interfaces.Visitor) (out []interfaces.OSSInfo, err error) {
	o.trace.SetClientSpanName("适配器-获取本地站点下可用的存储信息")
	newCtx, span := o.trace.AddClientTrace(ctx)
	defer func() { o.trace.TelemetrySpanEnd(span, err) }()

	url := fmt.Sprintf("%s/v1/local-storages?enabled=true", o.adminAddress)
	headers := map[string]string{
		"x-error-code": ErrCodeTypeToStr[visitor.ErrorCodeType],
	}
	respParam, err := o.client.Get(newCtx, url, headers)
	if err != nil {
		o.log.Errorf("oss_gateway GetLocalEnabledOSSInfo err: %v", err)
		return
	}

	// 解析数据
	tempOSSInfos := respParam.([]interface{})
	out = make([]interfaces.OSSInfo, len(tempOSSInfos))
	for k, v := range tempOSSInfos {
		tmp := v.(map[string]interface{})
		data := interfaces.OSSInfo{}

		data.ID = tmp["id"].(string)
		data.BDefault = tmp["default"].(bool)

		out[k] = data
	}

	return
}
