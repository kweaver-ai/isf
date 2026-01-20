package efastcmp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/common/utils"
	"AuditLog/infra/cmp/efastcmp/dto/efastarg"
	"AuditLog/infra/cmp/efastcmp/dto/efastret"
	"AuditLog/infra/cmp/efastcmp/eferr"
	"AuditLog/infra/cmp/httpclientcmp"
)

// GetFsMetadata 获取文件系统（文件/目录/文档库）元数据-批量 【内部接口】
// fs: 文件系统 file system
// 【注意】会自动去掉不存在的文件id，返回的结果中不会包含不存在的文件id
func (e *EFast) GetFsMetadata(ctx context.Context, args *efastarg.GetFsMetadataArgDto) (ret efastret.GetFsMetadataRetDto, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	var (
		loopCount int
		maxLoop   = 3
	)

	ret = make(efastret.GetFsMetadataRetDto, 0)

	if len(args.IDs) == 0 && len(args.ObjIDs) == 0 {
		err = errors.New("[GetFsMetadata]file_ids不能为空")
		return
	}

	c := httpclientcmp.NewHTTPClient(e.arTrace)

Loop:
	argDto := efastarg.NewGetFsMetadataEFArgDto(args)

	apiUrl := fmt.Sprintf("%s/v1/items-batch/%v", e.getUrlPrefix(),
		args.Fields.ToPathString(),
	)

	e.logger.Infof("GetFsMetadata apiUrl: %s", apiUrl)

	if helpers.IsLocalDev() {
		// mock data
		ret = efastret.GetFsMetadataRetDto{
			{
				ID:         "gns://D42F2729C56E489A948985D4E75C5813/4e8bfbda-d99c-11eb-35b9-24e8e050xxx5",
				Name:       "test.txt",
				DocLibType: enums.DocLibTypeStrCustom,
				Path:       "a/test.txt",
			},
			{
				ID:         "gns://D42F2729C56E489A948985D4E75C5813",
				Name:       "a",
				DocLibType: enums.DocLibTypeStrCustom,
				Path:       "a",
			},
		}

		return
	}

	resp, err := c.PostJSONExpect2xx(ctx, apiUrl, argDto)

	respErr := &httpclientcmp.CommonRespError{}
	if errors.As(err, &respErr) {
		loopCount++
		if loopCount > maxLoop {
			return nil, errors.Wrap(err, "获取文件信息失败")
		}

		if respErr.Code == eferr.FileOrDirNotFound && respErr.Detail != nil {
			var notExistsIDs []string
			if _notExistsIDs, ok := respErr.Detail["ids"]; ok {
				notExistsIDs = utils.MustStrSlice2(_notExistsIDs)
			}

			// 去掉不存在的文件id
			args.IDs = utils.Difference(args.IDs, notExistsIDs)

			if len(args.IDs) == 0 {
				// 当去掉不存在的后为空时，返回
				err = nil
				return
			}

			goto Loop
		}
	}

	if err != nil {
		return
	}

	err = utils.JSON().Unmarshal([]byte(resp), &ret)
	if err != nil {
		return
	}

	return
}

func (e *EFast) GetFsMetadataMap(ctx context.Context, docIds []string, fields []efastarg.IbField) (m map[string]*efastret.FsMetadata, err error) {
	m = make(map[string]*efastret.FsMetadata)

	dto := &efastarg.GetFsMetadataArgDto{
		IDs:    docIds,
		Fields: fields,
	}

	var ret efastret.GetFsMetadataRetDto

	ret, err = e.GetFsMetadata(ctx, dto)
	if err != nil {
		return
	}

	for _, v := range ret {
		m[v.ID] = v
	}

	return
}

func (e *EFast) GetOneFsName(ctx context.Context, docId string) (name string, err error) {
	fields := []efastarg.IbField{
		efastarg.IbFieldName,
	}

	docLibInfoMap, err := e.GetFsMetadataMap(ctx, []string{docId}, fields)
	if err != nil {
		return
	}

	docLibInfo, ok := docLibInfoMap[docId]
	if !ok {
		return
	}

	name = docLibInfo.Name

	return
}

func (e *EFast) GetFsMetadataNameMap(ctx context.Context, docIds []string) (m map[string]string, err error) {
	m = make(map[string]string)

	fields := []efastarg.IbField{
		efastarg.IbFieldName,
	}

	docLibInfoMap, err := e.GetFsMetadataMap(ctx, docIds, fields)
	if err != nil {
		return
	}

	for k, v := range docLibInfoMap {
		m[k] = v.Name
	}

	return
}
