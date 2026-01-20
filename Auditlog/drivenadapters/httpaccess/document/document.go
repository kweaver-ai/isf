package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"AuditLog/common"
	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/drivenadapters"
	"AuditLog/drivenadapters/httpaccess/document/docaccret"
	"AuditLog/gocommon/api"
	"AuditLog/interfaces/drivenadapter/ihttpaccess"
)

var (
	dOnce sync.Once
	d     *document
)

type document struct {
	urlPrefix  string
	httpClient api.Client
	logger     api.Logger
}

func NewDocument() ihttpaccess.DocumentHttpAcc {
	dOnce.Do(func() {
		d = &document{
			urlPrefix:  common.SvcConfig.DocumentPrivateProtocol + "://" + common.SvcConfig.DocumentPrivateHost + ":" + common.SvcConfig.DocumentPrivatePort + "/api/document/v1",
			httpClient: drivenadapters.HTTPClient,
			logger:     drivenadapters.Logger,
		}
	})

	return d
}

// 批量获取文档库信息
func (d *document) GetBatchDocLibInfos(docIDs []string) (infos []*docaccret.DocLibItem, err error) {
	var resp *http.Response

	ctx := context.Background()
	fields := "doc_lib_type,name"
	addr := d.urlPrefix + fmt.Sprintf("/batch-doc-libs/%v", fields)
	reqBody := map[string]interface{}{
		"method": "GET",
		"ids":    docIDs,
	}

	reqBodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	// mock
	if helpers.IsLocalDev() {
		infos = []*docaccret.DocLibItem{
			{
				ID:   "gns://D42F2729C56E489A948985D4E75C5813",
				Name: "文档库1",
				Type: enums.DocLibTypeStrPersonal,
			},
			{
				ID:   "gns://D42F2729C56E489A948985D4E75C4813",
				Name: "文档库2",
				Type: enums.DocLibTypeStrDepartment,
			},
		}

		return
	}

	resp, err = d.httpClient.Post(ctx, addr, enums.HTTPHctJSON, bytes.NewReader(reqBodyByte))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBodyByte, err := io.ReadAll(resp.Body)
	if err = utils.JSON().Unmarshal(respBodyByte, &infos); err != nil {
		return
	}

	return
}
