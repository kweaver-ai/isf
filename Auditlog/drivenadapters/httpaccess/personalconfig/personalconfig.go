package personalconfig

import (
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
	personalConfigOnce sync.Once
	personalConfig     *PersonalConfig
)

type PersonalConfig struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewPersonalConfig() interfaces.PersonalConfigRepo {
	personalConfigOnce.Do(func() {
		conf := config.GetPersonalConfigConf()
		personalConfig = &PersonalConfig{
			urlPrefix:  conf.Private.Protocol + "://" + conf.Private.Host + ":" + strconv.Itoa(conf.Private.Port) + "/api/personal-config/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})
	return personalConfig
}

func (p *PersonalConfig) GetModuleInfoByName(moduleName string) (res *models.ServiceModuleInfo, statusCode int, err error) {
	var resp *http.Response
	ctx := context.Background()
	addr := p.urlPrefix + "/deployment/module/" + moduleName
	resp, err = p.httpClient.Get(ctx, addr)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode

	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		p.logger.Warnf("ERROR: GetModuleInfoByName url: %v, statusCode: %v", addr, statusCode)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		p.logger.Warnf("ERROR: GetModuleInfoByName url: %v, err: %v", addr, err)
		return
	}

	return
}
