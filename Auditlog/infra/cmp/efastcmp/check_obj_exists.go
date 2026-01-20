package efastcmp

import (
	"context"
	"fmt"

	"AuditLog/common/utils"
	"AuditLog/infra/cmp/efastcmp/dto/efastarg"
	"AuditLog/infra/cmp/efastcmp/dto/efastret"
	"AuditLog/infra/cmp/httpclientcmp"
)

// CheckObjExists 检查文件或目录是否存在
func (e *EFast) CheckObjExists(ctx context.Context, ids []string) (notExistsIDs []string, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	if len(ids) == 0 {
		panic("[CheckObjExists]file_ids不能为空")
	}

	// 1、构建参数
	args := &efastarg.GetFsMetadataArgDto{
		IDs: ids,
		Fields: []efastarg.IbField{
			efastarg.IbFieldName,
		},
	}

	// 2、调用接口
	c := httpclientcmp.NewHTTPClient(e.arTrace)

	argDto := efastarg.NewGetFsMetadataEFArgDto(args)
	apiUrl := fmt.Sprintf("%s/v1/items-batch/%v", e.getUrlPrefix(),
		args.Fields.ToPathString(),
	)

	ret := efastret.GetFsMetadataRetDto{}

	resp, err := c.PostJSONExpect2xx(ctx, apiUrl, argDto)
	if err != nil {
		return
	}

	// 3、处理返回结果
	err = utils.JSON().Unmarshal([]byte(resp), &ret)
	if err != nil {
		return
	}

	existsIDs := make([]string, 0, len(ids))
	for _, metadata := range ret {
		existsIDs = append(existsIDs, metadata.ID)
	}

	notExistsIDs = utils.Difference(ids, existsIDs)

	return
}
