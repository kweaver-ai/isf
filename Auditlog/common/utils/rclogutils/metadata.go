package rclogutils

import (
	"context"
	"encoding/json"

	"AuditLog/common/constants/rclogconsts"
	"AuditLog/locale"
	"AuditLog/models/rcvo"
)

var (
	acReportMetadata      *rcvo.ReportMetadataRes
	historyReportMetadata *rcvo.ReportMetadataRes
)

// 获取活跃日志报表元数据
func GetActiveMetadata() (meta *rcvo.ReportMetadataRes, err error) {
	if acReportMetadata != nil {
		return acReportMetadata, nil
	}

	ctx := context.TODO()

	idField := rcvo.ReportField{
		Field:                   rclogconsts.LogID,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogID),
		IsKvField:               0,
		IsCanSort:               0,
		IsCanSearch:             0,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig:       rcvo.ReportSearchFieldConfig{},
	}

	createTimeField := rcvo.ReportField{
		Field:                   rclogconsts.Date,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogDate),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                6,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{4},
			IsCanSearchByApi: false,
		},
	}

	levelField := rcvo.ReportField{
		Field:                   rclogconsts.LogLevel,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogLevel),
		IsKvField:               1,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{2},
			IsCanSearchByApi: true,
		},
	}

	userNameField := rcvo.ReportField{
		Field:                   rclogconsts.UserName,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogUser),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	macField := rcvo.ReportField{
		Field:                   rclogconsts.Mac,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogMac),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	ipField := rcvo.ReportField{
		Field:                   rclogconsts.IP,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogIP),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	userPathsField := rcvo.ReportField{
		Field:                   rclogconsts.UserPaths,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogUserPaths),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	opTypeField := rcvo.ReportField{
		Field:                   rclogconsts.OpType,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogOpType),
		IsKvField:               1,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{2},
			IsCanSearchByApi: true,
		},
	}

	msgField := rcvo.ReportField{
		Field:                   rclogconsts.Msg,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogMsg),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	exmsgField := rcvo.ReportField{
		Field:                   rclogconsts.ExMsg,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogExMsg),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	acReportMetadata = &rcvo.ReportMetadataRes{
		Fields: []rcvo.ReportField{
			idField,
			levelField,
			createTimeField,
			macField,
			ipField,
			userNameField,
			userPathsField,
			opTypeField,
			msgField,
			exmsgField,
		},
		DefaultSortField:     "date",
		DefaultSortDirection: "desc",
		IdField:              rclogconsts.LogID,
	}

	return acReportMetadata, nil
}

// GetHistoryMetadata 获取历史审计日志元数据
func GetHistoryMetadata(ctx context.Context) (meta *rcvo.ReportMetadataRes, err error) {
	if historyReportMetadata != nil {
		return historyReportMetadata, nil
	}

	idField := rcvo.ReportField{
		Field:                   rclogconsts.ID,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogID),
		IsKvField:               0,
		IsCanSort:               0,
		IsCanSearch:             0,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig:       rcvo.ReportSearchFieldConfig{},
	}

	dumpDateField := rcvo.ReportField{
		Field:                   rclogconsts.DumpDate,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogDumpDate),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             0,
		ShowType:                6,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig:       rcvo.ReportSearchFieldConfig{},
	}

	fileNameField := rcvo.ReportField{
		Field:                   rclogconsts.Name,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogFileName),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             1,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig: rcvo.ReportSearchFieldConfig{
			SupportTypes:     []int{1},
			IsCanSearchByApi: false,
		},
	}

	sizeField := rcvo.ReportField{
		Field:                   rclogconsts.Size,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogSize),
		IsKvField:               0,
		IsCanSort:               1,
		IsCanSearch:             0,
		ShowType:                1,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig:       rcvo.ReportSearchFieldConfig{},
	}

	operationField := rcvo.ReportField{
		Field:                   rclogconsts.Operation,
		FieldTitle:              locale.GetI18nCtx(ctx, locale.RCLogOperation),
		IsKvField:               0,
		IsCanSort:               0,
		IsCanSearch:             0,
		ShowType:                7,
		IsPmsCtrlField:          0,
		IsOrgStructureField:     0,
		OrgStructureFieldConfig: nil,
		SearchFieldConfig:       rcvo.ReportSearchFieldConfig{},
	}

	historyReportMetadata = &rcvo.ReportMetadataRes{
		Fields: []rcvo.ReportField{
			idField,
			fileNameField,
			dumpDateField,
			sizeField,
			operationField,
		},
		DefaultSortField:     "dump_date",
		DefaultSortDirection: "desc",
		IdField:              rclogconsts.ID,
	}

	return historyReportMetadata, nil
}

// GetIDFromErrorResponse 从错误响应中提取ID
func GetIDFromErrorResponse(errResp map[string]interface{}) (id int, ok bool) {
	if detail, exists := errResp["detail"].(map[string]interface{}); exists {
		if idNum, exists := detail["id"].(json.Number); exists {
			if idInt64, err := idNum.Int64(); err == nil {
				return int(idInt64), true
			}
		}
	}

	return 0, false
}
