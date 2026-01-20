package doccenter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"AuditLog/drivenadapters"
	"AuditLog/gocommon/api"
	"AuditLog/infra/config"
	"AuditLog/interfaces"
	"AuditLog/models/rcvo"
)

var (
	dOnce sync.Once
	dM    *docCenter
)

type docCenter struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewDocCenter() interfaces.DocCenterRepo {
	dOnce.Do(func() {
		conf := config.GetDocCenterConf()
		dM = &docCenter{
			urlPrefix:  conf.Private.Protocol + "://" + conf.Private.Host + ":" + strconv.Itoa(conf.Private.Port) + "/api/doc-center/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})
	return dM
}

// NewDataSourceGroup 新建数据源组
func (d *docCenter) NewDataSourceGroup(body *rcvo.DCNewDataSourceGroupBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := d.urlPrefix + "/report-center/datasourcegroup"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err = d.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		decoder := json.NewDecoder(resp.Body)
		decoder.UseNumber()
		err = decoder.Decode(&errResp)
		if err != nil {
			return
		}
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}

	return res, statusCode, errResp, nil
}

// NewDataSource 新建数据源
func (d *docCenter) NewDataSource(body *rcvo.DCNewDataSourceBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := d.urlPrefix + "/report-center/datasource"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err = d.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		decoder := json.NewDecoder(resp.Body)
		decoder.UseNumber()
		err = decoder.Decode(&errResp)
		if err != nil {
			return
		}
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}

	return res, statusCode, errResp, nil
}

// NewReport 新建报表
func (d *docCenter) NewReport(body *rcvo.DCNewReportBody) (res *rcvo.DCResponse, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := d.urlPrefix + "/report-center/report"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err = d.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}

	return res, statusCode, nil
}

// NewReportGroup 新建报表组
func (d *docCenter) NewReportGroup(body *rcvo.DCNewBizGroupBody) (res *rcvo.DCResponse, statusCode int, errResp map[string]interface{}, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := d.urlPrefix + "/report-center/bizgroup"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	resp, err = d.httpClient.Post(ctx, addr, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		decoder := json.NewDecoder(resp.Body)
		decoder.UseNumber()
		err = decoder.Decode(&errResp)
		if err != nil {
			return
		}
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}

	return res, statusCode, errResp, nil
}

// GetDataSourceFields 获取数据源字段
func (d *docCenter) GetDataSourceFields(dataSourceID int) (res *rcvo.DCDataSourceFieldsRes, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := d.urlPrefix + fmt.Sprintf("/report-center/datasource/%v/fields", dataSourceID)
	resp, err = d.httpClient.Get(ctx, addr)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}

	return res, statusCode, nil
}
