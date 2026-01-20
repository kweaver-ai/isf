package kc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
	"AuditLog/interfaces"
	"AuditLog/models"
)

var (
	kcOnce sync.Once
	kc     *Kc
)

type Kc struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewKc() interfaces.KcRepo {
	kcOnce.Do(func() {
		conf := config.GetKcConf()
		kc = &Kc{
			urlPrefix:  conf.Private.Protocol + "://" + conf.Private.Host + ":" + strconv.Itoa(conf.Private.Port) + "/api/pri-kc-mc/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})
	return kc
}

func (k *Kc) GetUserInfoByIDS(userIDs []string) (res []*models.KcUserInfo, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := k.urlPrefix + "/search-user?filter_as_id"

	body := map[string]interface{}{
		"filter_as_ids": userIDs,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err = k.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}
	statusCode := resp.StatusCode

	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		k.logger.Warnf("ERROR: GetDeptInfoByIDs url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	kcRes := &models.KcUserInfoRes{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&kcRes)
	if err != nil {
		k.logger.Infof("ERROR: GetDeptInfoByID:%v\n", err)
		return
	}

	// 获取成功
	if kcRes.Code == 200033800 {
		res = make([]*models.KcUserInfo, len(kcRes.Data))
		for i := range kcRes.Data {
			res[i] = &kcRes.Data[i]
		}
	}

	return
}
