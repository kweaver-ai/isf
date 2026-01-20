package drivenadapters

import (
	"context"
	"fmt"
	"net/http"
	"policy_mgnt/common/config"
	"policy_mgnt/interfaces"
	"strconv"
	"sync"

	"policy_mgnt/common"

	"github.com/kweaver-ai/go-lib/httpclient"
	"github.com/kweaver-ai/go-lib/observable"
)

var (
	licenseOnce sync.Once
	lices       *license
)

type license struct {
	log        common.Logger
	rawClient  *http.Client
	client     httpclient.HTTPClient
	licenseURL string
	trace      observable.Tracer
}

// NewLicense 创建许可证驱动
func NewLicense() *license {
	licenseOnce.Do(func() {
		config := config.Config.ProtonApplicationConfig
		lices = &license{
			log:        common.NewLogger(),
			rawClient:  httpclient.NewRawHTTPClient(),
			client:     httpclient.NewHTTPClient(common.SvcARTrace),
			licenseURL: fmt.Sprintf("http://%s:%s/api/proton-application/v1/license?combine=true&status=actived&dimension=user", config.Host, config.Port),
			trace:      common.SvcARTrace,
		}
	})

	return lices
}

// GetLicenses 获取许可证
func (l *license) GetLicenses(ctx context.Context) (infos map[string]interfaces.License, err error) {
	l.trace.SetClientSpanName("获取许可证")
	newCtx, span := l.trace.AddClientTrace(ctx)
	defer func() { l.trace.TelemetrySpanEnd(span, err) }()

	infos = make(map[string]interfaces.License)
	respParam, err := l.client.Get(newCtx, l.licenseURL, nil)
	if err != nil {
		l.log.Errorf("drivenadapters license GetLicenses err: %v", err)
		return
	}

	// 解析数据
	temp := respParam.(map[string]interface{})
	for k, v := range temp {
		tempV := v.([]interface{})

		if len(tempV) == 0 {
			continue
		}

		tempLicense := tempV[0].(map[string]interface{})
		// 获取已授权人数
		var totalUserQuota int
		authorized, ok := tempLicense["authorized"].(string)
		if ok {
			totalUserQuota, err = strconv.Atoi(authorized)
			if err != nil {
				l.log.Errorln("drivenadapters license GetLicenses strconv.Atoi err: %v", err)
				continue
			}
			l.log.Infof("drivenadapters license GetLicenses authorized: %d", totalUserQuota)
		} else {
			continue
		}

		infos[k] = interfaces.License{
			Product:        k,
			TotalUserQuota: totalUserQuota,
		}

	}
	return infos, nil
}
