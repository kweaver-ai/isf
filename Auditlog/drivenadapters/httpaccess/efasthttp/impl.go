package efasthttp

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"AuditLog/common/enums"
	"AuditLog/common/helpers"
	"AuditLog/infra/cmp/efastcmp/dto/efastarg"
	"AuditLog/infra/cmp/efastcmp/dto/efastret"
)

// GetDocLibType 获取文档库类型
// 【注意】：当gns对应的文件不存在时，会返回空字符串，不会报错（todo:考虑下这样是否合理）
func (e *eFastHttpAcc) GetDocLibType(ctx context.Context, gns string) (docLibType enums.DocLibType, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	if helpers.IsLocalDev() {
		// 模拟数据
		docLibType = enums.DocLibTypeStrCustom
		return
	}

	dto := &efastarg.GetFsMetadataArgDto{
		IDs: []string{gns},
		Fields: []efastarg.IbField{
			efastarg.IbFieldDocLibTypes,
		},
	}

	var ret efastret.GetFsMetadataRetDto

	ret, err = e.eFast.GetFsMetadata(ctx, dto)
	if err != nil {
		helpers.RecordErrLogWithPos(e.logger, err, "eFastHttpAcc.GetDocLibType")
		return "", errors.Wrap(err, "获取文档库类型失败")
	}

	if len(ret) == 0 {
		docLibType = ""
		return
	}

	docLibType = ret[0].DocLibType

	return
}

// CheckObjExists 检查文件或目录是否存在
// ids：文件或目录gns
// 返回不存在的文件或目录gns
func (e *eFastHttpAcc) CheckObjExists(ctx context.Context, ids []string) (notExistsIDs []string, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	return e.eFast.CheckObjExists(ctx, ids)
}

// CheckOneObjExists 检查某文件或目录是否存在
// id：文件或目录gns
// exist：是否存在
func (e *eFastHttpAcc) CheckOneObjExists(ctx context.Context, id string) (exist bool, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	ids := []string{id}

	notExistsIDs, err := e.eFast.CheckObjExists(ctx, ids)
	if err != nil {
		return
	}

	if len(notExistsIDs) == 0 {
		exist = true
	}

	return
}

func (e *eFastHttpAcc) Path2Gns(ctx context.Context, path, asToken string) (isNotExists bool, gns string, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	if helpers.IsLocalDev() {
		// 模拟数据
		isNotExists = false
		//nolint:lll
		gns = "gns://73D7C5DA8F694BFC9FD72E1D49568C8C/BBC492D00CA945FEAAFFE81F3FB60FA8/4056603F3BA54AE5959E797853440EC5/D7D254FB33CF43F799D03673D7582459/7D819026A9CE4420ADD8EBDB7E5CE2B6"

		return
	}

	isNotExists, gns, err = e.eFast.Path2Gns(ctx, path, asToken)

	return
}

func (e *eFastHttpAcc) Gns2Path(ctx context.Context, gns []string) (pathMap map[string]string, err error) {
	dto := &efastarg.GetFsMetadataArgDto{
		Fields: []efastarg.IbField{
			efastarg.IbFieldPaths,
		},
	}
	if len(gns) > 0 {
		if strings.HasPrefix(gns[0], "gns://") {
			dto.IDs = gns
		} else {
			dto.ObjIDs = gns
		}
	}

	var ret efastret.GetFsMetadataRetDto

	ret, err = e.eFast.GetFsMetadata(ctx, dto)
	if err != nil {
		helpers.RecordErrLogWithPos(e.logger, err, "eFastHttpAcc.GetPaths")
		return map[string]string{}, errors.Wrap(err, "获取文档路径名失败")
	}

	pathMap = make(map[string]string, len(ret))
	for _, item := range ret {
		pathMap[item.ID] = item.Path
	}

	return
}

func (e *eFastHttpAcc) CreateMultiLevelDir(ctx context.Context, parentDocID, path, asToken string) (dirDocID string, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	if helpers.IsLocalDev() {
		// 模拟数据
		//nolint:lll
		dirDocID = "gns://73D7C5DA8F694BFC9FD72E1D49568C8C/BBC492D00CA945FEAAFFE81F3FB60FA8/4056603F3BA54AE5959E797853440EC5/D7D254FB33CF43F799D03673D7582459/7D819026A9CE4420ADD8EBDB7E5CE2B6"
		return
	}

	req := &efastarg.CreateMultiLevelDirReq{
		DocId: parentDocID,
		Path:  path,
	}

	ret, err := e.eFast.CreateMultiLevelDir(ctx, req, asToken)
	if err != nil {
		return
	}

	dirDocID = ret.DocId

	return
}

func (e *eFastHttpAcc) GetOneFsName(ctx context.Context, docID string) (name string, err error) {
	ctx, span := e.arTrace.AddInternalTrace(ctx)
	defer func() { e.arTrace.TelemetrySpanEnd(span, err) }()

	if helpers.IsLocalDev() {
		// 模拟数据
		name = "xx库"
		return
	}

	name, err = e.eFast.GetOneFsName(ctx, docID)
	if err != nil {
		helpers.RecordErrLogWithPos(e.logger, err, "eFastHttpAcc.GetOneFsName")
		return "", errors.Wrap(err, "获取名称失败")
	}

	return
}
