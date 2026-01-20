package paramutils

import (
	"context"
	"fmt"

	"AuditLog/common"
	"AuditLog/common/constants/lsconsts"
	"AuditLog/errors"
	"AuditLog/models/rcvo"
)

type AvailableParam struct {
	SearchFields  []string
	OrderFields   []string
	KeyWordFields []string
	DataFields    []string
}

// CategoryCheck 检查 category 是否合法
func CategoryCheck(ctx context.Context, category string) error {
	if !common.InArray(category, common.AllLogType) {
		return errors.NewCtx(ctx, errors.BadRequestErr, "invalid category", "")
	}

	return nil
}

// LimitCheck 检查 limit 是否合法
func LimitCheck(ctx context.Context, limit int, defaultLimit int) error {
	if limit <= 0 || limit > defaultLimit {
		return errors.NewCtx(ctx, errors.BadRequestErr, "invalid limit", "")
	}

	return nil
}

// ParamsCheck 检查 entry 是否在 candidates 中
func ParamsCheck(ctx context.Context, entry, candidates []string, tag string) error {
	for _, v := range entry {
		if v == "" {
			continue
		}

		if !common.InArray(v, candidates) {
			cause := fmt.Sprintf("invalid %s field", tag)
			return errors.NewCtx(ctx, errors.BadRequestErr, cause, nil)
		}
	}

	return nil
}

// GetAvaliableParams 获取可用的参数
func GetAvaliableParams(getParams func() (*rcvo.ReportMetadataRes, error)) (avaliableParams *AvailableParam) {
	avaliableParams = &AvailableParam{
		DataFields:    make([]string, 0),
		SearchFields:  make([]string, 0),
		KeyWordFields: make([]string, 0),
		OrderFields:   make([]string, 0),
	}

	metaData, err := getParams()
	if err != nil {
		return
	}

	for _, v := range metaData.Fields {
		avaliableParams.DataFields = append(avaliableParams.DataFields, v.Field)
		if v.IsCanSearch == 1 {
			avaliableParams.SearchFields = append(avaliableParams.SearchFields, v.Field)
			if v.SearchFieldConfig.IsCanSearchByApi {
				avaliableParams.KeyWordFields = append(avaliableParams.KeyWordFields, v.Field)
			}
		}

		if v.IsCanSort == 1 {
			avaliableParams.OrderFields = append(avaliableParams.OrderFields, v.Field)
		}
	}

	return
}

// DumpFieldsCheck 检查转存策略字段是否合法
func DumpFieldsCheck(ctx context.Context, fields []string) error {
	for _, v := range fields {
		if !common.InArray(v, lsconsts.AllDumpFields) {
			return errors.NewCtx(ctx, errors.BadRequestErr, "invalid dump field: "+v, "")
		}
	}

	return nil
}
