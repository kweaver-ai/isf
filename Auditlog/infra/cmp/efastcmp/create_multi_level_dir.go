package efastcmp

import (
	"context"
	"fmt"

	"AuditLog/common/utils"
	"AuditLog/infra/cmp/efastcmp/dto/efastarg"
	"AuditLog/infra/cmp/efastcmp/dto/efastret"
	"AuditLog/infra/cmp/httpclientcmp"
)

func (e *EFast) CreateMultiLevelDir(ctx context.Context, req *efastarg.CreateMultiLevelDirReq, token string) (ret *efastret.CreateMultiLevelDirRsp, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	url := fmt.Sprintf("%s/v1/dir/createmultileveldir", e.getPublicUrlPrefix())

	// 2、调用接口
	opt := httpclientcmp.WithToken(token)
	c := httpclientcmp.NewHTTPClient(e.arTrace, opt)

	resp, err := c.PostJSONExpect2xx(ctx, url, req)
	if err != nil {
		return
	}

	// 3、处理返回结果
	ret = &efastret.CreateMultiLevelDirRsp{}

	err = utils.JSON().Unmarshal([]byte(resp), &ret)
	if err != nil {
		return
	}

	return
}
