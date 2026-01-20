package efastcmp

import (
	"context"
	"errors"
	"fmt"

	"AuditLog/common/utils"
	"AuditLog/infra/cmp/efastcmp/eferr"
	"AuditLog/infra/cmp/efastcmp/eftypes"
	"AuditLog/infra/cmp/httpclientcmp"
)

func (e *EFast) GetInfoByPath(ctx context.Context, path, token string) (isNotExists bool, ret *eftypes.Path2GnsResponse, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	url := fmt.Sprintf("%s/v1/file/getinfobypath", e.getPublicUrlPrefix())

	// 1、构建参数
	req := eftypes.Path2GnsReq{
		Namepath: path,
	}

	// 2、调用接口
	opt := httpclientcmp.WithToken(token)
	c := httpclientcmp.NewHTTPClient(e.arTrace, opt)

	resp, err := c.PostJSONExpect2xx(ctx, url, req)
	respErr := &httpclientcmp.CommonRespError{}

	if errors.As(err, &respErr) {
		if respErr.Code == eferr.FileOrDirNotFound {
			isNotExists = true
			err = nil
		}

		return
	}

	if err != nil {
		return
	}

	// 3、处理返回结果
	ret = &eftypes.Path2GnsResponse{}

	err = utils.JSON().Unmarshal([]byte(resp), &ret)
	if err != nil {
		return
	}

	return
}

func (e *EFast) Path2Gns(ctx context.Context, path, token string) (isNotExists bool, gns string, err error) {
	isNotExists, ret, err := e.GetInfoByPath(ctx, path, token)
	if err != nil {
		return
	}

	if isNotExists {
		return
	}

	gns = ret.DocId

	return
}
