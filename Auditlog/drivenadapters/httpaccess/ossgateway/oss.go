package ossgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
	"AuditLog/interfaces"
	"AuditLog/models"
)

var (
	oOnce sync.Once
	oM    *ossGateway
)

type ossGateway struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewOssGateway() interfaces.OssGatewayRepo {
	oOnce.Do(func() {
		conf := config.GetOssGatewayConf()
		oM = &ossGateway{
			urlPrefix:  conf.Private.Protocol + "://" + conf.Private.Host + ":" + strconv.Itoa(conf.Private.Port) + "/api/ossgateway/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})
	return oM
}

// GetLocalOSSInfo 获取本地对象存储信息
func (o *ossGateway) GetLocalOSSInfo() (res []*models.OSSInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := o.urlPrefix + "/local-storages?enabled=true"
	resp, err = o.httpClient.Get(ctx, addr)
	if err != nil {
		o.logger.Errorf("[GetLocalOSSInfo]: Get url: %v, failed:%v\n", addr, err)
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[GetLocalOSSInfo]: Get url: %v, statusCode: %v", addr, statusCode)
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		o.logger.Infof("[GetLocalOSSInfo]: Decode failed:%v\n", err)
		return
	}

	return
}

// GetUploadInfo 获取上传信息
func (o *ossGateway) GetUploadInfo(ossID, objName string) (res *models.OSSUploadInfo, statusCode int, err error) {
	fileSize := 1024 * 1024 * 1024 // 1G
	var resp *http.Response
	ctx := context.Background()
	// 此URL为自适应分片/普通上传，在文件大小为1G时一定为分片上传，如果有问题，与oss网关联系
	addr := fmt.Sprintf(
		"%s/upload-info/%s/%s?internal_request=true&file_size=%d&request_method=PUT",
		o.urlPrefix, ossID, objName, fileSize,
	)
	resp, err = o.httpClient.Get(ctx, addr)
	if err != nil {
		o.logger.Errorf("[UploadInfo]: Get url: %v, failed:%v\n", addr, err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[UploadInfo]: Get url: %v, statusCode: %v", addr, statusCode)
		return
	}

	resTmp := map[string]interface{}{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&resTmp)
	if err != nil {
		o.logger.Infof("[UploadInfo]: Decode failed:%v\n", err)
		return
	}

	partSize, _ := strconv.Atoi(resTmp["partsize"].(string))
	maxNum, _ := strconv.Atoi(resTmp["max_num"].(string))
	res = &models.OSSUploadInfo{
		UploadID:   resTmp["upload_id"].(string),
		UploadType: resTmp["upload_type"].(string),
		PartSize:   partSize,
		MaxNum:     maxNum,
	}

	return
}

// GetUploadPartRequestInfo 获取上传分片的请求信息
func (o *ossGateway) GetUploadPartRequestInfo(ossID, objName, uploadID string, partNumber int) (res *models.OSSRequestInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := fmt.Sprintf(
		"%s/uploadpart/%s/%s?internal_request=true&part_id=%d&upload_id=%s",
		o.urlPrefix, ossID, objName, partNumber, uploadID,
	)
	resp, err = o.httpClient.Get(ctx, addr)
	if err != nil {
		o.logger.Errorf("[UploadPart]: Get url: %v, failed:%v\n", addr, err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[UploadPart]: Get url: %v, statusCode: %v", addr, statusCode)
		return
	}

	var jsonResponse map[string]models.OSSRequestInfo
	// 解析响应体
	if err = json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		o.logger.Warnf("[UploadPart] failed to decode response body: %v", err)
		return
	}

	partInfo, ok := jsonResponse[strconv.Itoa(partNumber)]
	if !ok {
		err = fmt.Errorf("part number %d not found in response", partNumber)
		o.logger.Warnf("[UploadPart] missing part info: %v", err)
		return
	}

	res = &models.OSSRequestInfo{
		Method:      partInfo.Method,
		URL:         partInfo.URL,
		Headers:     partInfo.Headers,
		RequestBody: partInfo.RequestBody,
	}

	return
}

// UploadPartByURL 上传分片
func (o *ossGateway) UploadPartByURL(url string, method string, body string, headers map[string]string) (res *models.OSSUploadPartInfo, statusCode int, err error) {
	ctx := context.Background()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		o.logger.Errorf("[UploadPartByURL]: NewRequest failed:%v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := o.httpClient.Do(ctx, req)
	if err != nil {
		o.logger.Errorf("[UploadPartByURL]: Do url: %v, failed:%v\n", url, err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[UploadPartByURL]: Do url: %v, statusCode: %v", url, statusCode)
		return
	}

	// 获取resp中的etag
	etag := resp.Header.Get("ETag")
	res = &models.OSSUploadPartInfo{
		Etag: etag,
		Size: len(body),
	}

	return
}

// GetCompleteUploadRequestInfo 获取完成上传的请求信息
func (o *ossGateway) GetCompleteUploadRequestInfo(ossID, objName, uploadID string, multiPartInfo map[int]models.OSSUploadPartInfo) (res *models.OSSRequestInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := fmt.Sprintf(
		"%s/completeupload/%s/%s?internal_request=true&upload_id=%s",
		o.urlPrefix, ossID, objName, uploadID,
	)

	reqBody := make(map[string]string)
	for partNum, info := range multiPartInfo {
		// 移除 ETag 中的引号
		etag := strings.Trim(info.Etag, `"`)
		reqBody[fmt.Sprintf("%d", partNum)] = etag
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		o.logger.Warnf("[CompleteUpload]: Marshal request body failed:%v\n", err)
		return
	}

	resp, err = o.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		o.logger.Errorf("[CompleteUpload]: Post url: %v, failed:%v\n", addr, err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		rep := map[string]string{}
		if err = json.NewDecoder(resp.Body).Decode(&rep); err != nil {
			o.logger.Warnf("[CompleteUpload]: Post url: %v, statusCode: %v, body: %v", addr, statusCode, err)
		} else {
			o.logger.Warnf("[CompleteUpload]: Post url: %v, statusCode: %v, body: %v", addr, statusCode, rep)
		}
		return
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		o.logger.Infof("[CompleteUpload]: Decode failed:%v\n", err)
		return
	}

	return
}

// CompleteUpload 完成上传
func (o *ossGateway) CompleteUploadByURL(url string, method string, body string, headers map[string]string) (resp *http.Response, statusCode int, err error) {
	ctx := context.Background()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		o.logger.Errorf("[CompleteUpload]: NewRequest failed:%v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	respTmp, err := o.httpClient.Do(ctx, req)
	if err != nil {
		o.logger.Errorf("[CompleteUpload]: Do failed:%v\n", err)
		return
	}
	defer respTmp.Body.Close()

	statusCode = respTmp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[CompleteUpload]: Do url: %v, statusCode: %v", url, statusCode)
		return
	}

	resp = respTmp

	return
}

func (o *ossGateway) GetDownLoadInfo(ossID, objName, fileName string, isInternal bool) (res *models.OSSRequestInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := fmt.Sprintf(
		"%s/download/%s/%s?internal_request=%t&type=query_string",
		o.urlPrefix, ossID, objName, isInternal,
	)
	if fileName != "" {
		addr += "&save_name=" + fileName
	}
	resp, err = o.httpClient.Get(ctx, addr)
	if err != nil {
		o.logger.Errorf("[GetDownLoadInfo]: Get url: %v, failed:%v\n", addr, err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[GetDownLoadInfo]: Get url: %v, statusCode: %v", addr, statusCode)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		o.logger.Infof("[GetDownLoadInfo]: Decode failed:%v\n", err)
		return
	}

	return
}

func (o *ossGateway) DownloadBlockByURL(url string, method string, body string, headers map[string]string, start, end int64) (resp *http.Response, statusCode int, err error) {
	ctx := context.Background()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		o.logger.Errorf("[DownloadBlockByURL]: NewRequest failed:%v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err = o.httpClient.Do(ctx, req)
	if err != nil {
		o.logger.Errorf("[DownloadBlockByURL]: Do failed:%v\n", err)
		return
	}

	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[DownloadBlockByURL]: Do url: %v, statusCode: %v", url, statusCode)
		return
	}

	return
}

// GetAvailableOSSID 获取可用对象存储ID
func (o *ossGateway) GetAvailableOSSID() (ossID string, err error) {
	ossInfos, _, err := o.GetLocalOSSInfo()
	if err != nil {
		return
	}

	for i, ossInfo := range ossInfos {
		if i == 0 {
			ossID = ossInfo.ID
		}

		if ossInfo.Default {
			ossID = ossInfo.ID
			break
		}
	}

	return
}

// GetDeleteRequestInfo 获取删除对象的请求协议
func (o *ossGateway) GetDeleteRequestInfo(ossID, objName string) (res *models.OSSRequestInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := fmt.Sprintf(
		"%s/delete/%s/%s?internal_request=true",
		o.urlPrefix, ossID, objName,
	)
	resp, err = o.httpClient.Get(ctx, addr)
	if err != nil {
		o.logger.Errorf("[GetDeleteRequestInfo]: Get url: %v, failed:%v\n", addr, err)
		return
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		o.logger.Warnf("[GetDeleteRequestInfo]: Get url: %v, statusCode: %v", addr, statusCode)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		o.logger.Infof("[GetDeleteRequestInfo]: Decode failed:%v\n", err)
		return
	}

	return
}

// DeleteObjectByURL 删除对象
func (o *ossGateway) DeleteObjectByURL(url string, method string, body string, headers map[string]string) (resp *http.Response, statusCode int, err error) {
	ctx := context.Background()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		o.logger.Errorf("[DeleteObjectByURL]: NewRequest failed:%v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err = o.httpClient.Do(ctx, req)
	if err != nil {
		o.logger.Errorf("[DeleteObjectByURL]: Do failed:%v\n", err)
		return
	}

	return
}
