package rcvo

// 依赖字段
type DependentField struct {
	Field  string `json:"field"`
	IsMust bool   `json:"is_must"`
}

// 搜索字段配置
type ReportSearchFieldConfig struct {
	DependentFields  []DependentField `json:"dependent_fields"`
	SupportTypes     []int            `json:"support_types" validate:"required,oneof=1 2 3 4 5"`
	IsCanSearchByApi bool             `json:"is_can_search_by_api" validate:"required"`

	// description: |
	//  - 搜索项标签（搜索项的显示名称）
	//  - not required（如果为空，则使用字段的field_title_custom作为搜索项标签）
	SearchLabel string `json:"search_label"`
}

// 组织架构字段配置
type ReportOrgStructureFieldConfig struct {
	SelectType int `json:"select_type" validate:"oneof=1 2 3 4 5 6 7 8 9 10 11 12 13 14 15"`
	IsMultiple int `json:"is_multiple" validate:"oneof=0 1"`
}

// 报表字段
type ReportField struct {
	Field                   string                         `json:"field" validate:"required"`                            // 字段
	FieldTitle              string                         `json:"field_title" validate:"required"`                      // 字段标题
	IsKvField               int                            `json:"is_kv_field" validate:"required,oneof=0 1"`            // 是否为键值对字段
	IsCanSort               int                            `json:"is_can_sort" validate:"required,oneof=0 1"`            // 是否可排序
	IsCanSearch             int                            `json:"is_can_search" validate:"required,oneof=0 1"`          // 是否可搜索
	ShowType                int                            `json:"show_type" validate:"required,oneof=1 2 3"`            // 显示类型
	IsPmsCtrlField          int                            `json:"is_pms_ctrl_field" validate:"required,oneof=0 1"`      // 是否权限控制字段
	IsOrgStructureField     int                            `json:"is_org_structure_field" validate:"required,oneof=0 1"` // 是否为组织架构字段
	SearchFieldConfig       ReportSearchFieldConfig        `json:"search_field_config" validate:"required"`              // 搜索字段配置
	OrgStructureFieldConfig *ReportOrgStructureFieldConfig `json:"org_structure_field_config"`                           // 组织架构字段配置
}

// 获取报表数据源元数据响应
type ReportMetadataRes struct {
	Fields                 []ReportField `json:"fields"`
	DefaultSortField       string        `json:"default_sort_field"`
	DefaultSortDirection   string        `json:"default_sort_direction" validate:"oneof=asc desc"`
	UniqueIncrementalField string        `json:"unique_incremental_field"`
	IdField                string        `json:"id_field" validate:"required"`
}

// 排序字段
type OrderField struct {
	Field          string `json:"field" validate:"required"`
	Direction      string `json:"direction" validate:"required,oneof=asc desc"`
	LastFieldValue string `json:"last_field_value"`
}

// 排序字段列表
type OrderFields []OrderField

// 获取报表数据列表请求
type ReportGetDataListReq struct {
	Limit     int            `json:"limit" validate:"omitempty,min=1,max=5000"`
	Offset    int            `json:"offset" validate:"omitempty,min=0"`
	Condition map[string]any `json:"condition"`
	IDs       []string       `json:"ids"`
	OrderBy   OrderFields    `json:"order_by"`
}

// 获取活跃日志数据列表请求
type ReportGetActiveDataListReq struct {
	ReportGetDataListReq
	IDs []uint64 `json:"ids"`
}

// 活跃日志报表数据
type ActiveLogReport struct {
	ID          string `json:"log_id"`
	UserName    string `json:"user_name"`
	CreatedTime int64  `json:"date"`
	IP          string `json:"ip"`
	Mac         string `json:"mac"`
	Msg         string `json:"msg"`
	ExMsg       string `json:"exmsg"`
	OpType      string `json:"op_type"`
	UserPaths   string `json:"user_paths"`
	Level       string `json:"level"`
	ObjName     string `json:"obj_name"`
	ObjType     string `json:"obj_type"`
}

// 活跃日志报表数据列表
type ActiveLogReports []ActiveLogReport

// 活跃日志报表数据列表响应
type ActiveReportListRes struct {
	Entries    ActiveLogReports `json:"entries" validate:"required"`
	TotalCount int              `json:"total_count" validate:"required"`
}

// 历史日志报表数据
type HistoryLogReport struct {
	ID       string `json:"id"`
	FileName string `json:"name"`
	Size     string `json:"size"`
	DumpDate int64  `json:"dump_date"`
}

// 历史日志报表数据列表
type HistoryLogReports []HistoryLogReport

// 历史日志报表数据列表响应
type HistoryReportListRes struct {
	Entries    HistoryLogReports `json:"entries" validate:"required"`
	TotalCount int               `json:"total_count" validate:"required"`
}

// 获取报表字段值请求
type ReportGetFieldValuesReqBody struct {
	Limit     int            `json:"limit" validate:"omitempty,min=1,max=5000"`
	Offset    int            `json:"offset" validate:"omitempty,min=0"`
	Condition map[string]any `json:"condition"`
	KeyWord   string         `json:"keyword"`
}

// ReportGetFieldValuesReq 获取报表字段值列表请求
type ReportGetFieldValuesReq struct {
	Field string `json:"field" validate:"required"`
	ReportGetFieldValuesReqBody
}

// ReportFieldValue 报表字段值
type ReportFieldValue struct {
	ValueCode string `json:"value_code"`
	ValueName any    `json:"value_name"`
}

// ReportFieldValues 报表字段值列表
type ReportFieldValues []ReportFieldValue

// ReportFieldValuesRes 获取报表字段值列表响应
type ReportFieldValuesRes struct {
	TotalCount int               `json:"total_count" validate:"required"`
	Entries    ReportFieldValues `json:"entries" validate:"required"`
}

// DCResponse 报表中心响应
type DCResponse struct {
	ID int `json:"id"`
}

// DCNewDataSourceGroupBody 新建数据源组请求
type DCNewDataSourceGroupBody struct {
	Name           string `json:"name" validate:"required"`              // 数据源组名称
	IsSystemConfig int    `json:"is_system_config" validate:"oneof=0 1"` // 是否为系统配置
}

// DCDataSourceField 数据源字段
type DCDataSourceField struct {
	ReportField
	FieldTitleCustom string `json:"field_title_custom" validate:"required"` // 自定义字段标题
}

// DCNewDataSourceBody 新建数据源请求
type DCNewDataSourceBody struct {
	Name                   string              `json:"name" validate:"required"`                         // 数据源名称
	GroupID                int                 `json:"datasource_group_id" validate:"required"`          // 数据源组ID
	ApiPrefix              string              `json:"api_prefix" validate:"required"`                   // 数据源API前缀
	Fields                 []DCDataSourceField `json:"rc_datasource_fields" validate:"required"`         // 数据源字段
	DefaultSortField       string              `json:"default_sort_field"`                               // 默认排序字段
	DefaultSortDirection   string              `json:"default_sort_direction" validate:"oneof=asc desc"` // 默认排序方向
	UniqueIncrementalField string              `json:"unique_incremental_field"`                         // 唯一增量字段
	IdField                string              `json:"id_field" validate:"required"`                     // ID字段
	IsSystemConfig         int                 `json:"is_system_config" validate:"oneof=0 1"`            // 是否为系统配置
	InternalLabel          string              `json:"internal_label"`                                   // 内部标签
}

// DCNewBizGroupBody 新建报表业务组请求
type DCNewBizGroupBody struct {
	Name           string `json:"name" validate:"required"`              // 业务组名称
	IsSystemConfig int    `json:"is_system_config" validate:"oneof=0 1"` // 是否为系统配置
}

// DCSearchField 搜索字段
type DCSearchField struct {
	Field      string `json:"field"`                                     // 字段
	RcFieldID  int    `json:"rc_field_id" validate:"required"`           // 字段ID
	FieldType  int    `json:"field_type" validate:"required"`            // 字段类型
	IsRequired int    `json:"is_required" validate:"required,oneof=0 1"` // 是否必填
}

// DCNewReportBody 新建报表请求
type DCNewReportBody struct {
	Name           string          `json:"name" validate:"required"`              // 报表名称
	BizGroupID     int             `json:"biz_group_id" validate:"required"`      // 业务组ID
	DataSourceID   int             `json:"rc_datasource_id" validate:"required"`  // 数据源ID
	ShowFieldIDs   []int           `json:"show_field_ids" validate:"required"`    // 显示字段ID列表
	SearchFields   []DCSearchField `json:"search_fields" validate:"required"`     // 搜索字段列表
	RcLabel        string          `json:"rc_label"`                              // 内部标签
	IsSystemConfig int             `json:"is_system_config" validate:"oneof=0 1"` // 是否为系统配置
	IsExportable   int             `json:"is_exportable" validate:"oneof=0 1"`    // 是否可导出
}

// DCDataSourceFieldItem 数据源字段
type DCDataSourceFieldItem struct {
	ID               int    `json:"field_id"`
	Field            string `json:"field"`
	FieldTitleCustom string `json:"field_title_custom"`
}

// DCDataSourceFieldsRes 数据源字段列表响应
type DCDataSourceFieldsRes struct {
	TotalCount int                     `json:"total_count" validate:"required"`
	Entries    []DCDataSourceFieldItem `json:"entries" validate:"required"`
}

// DCLogReportInfo 审计日志报表信息
type DCLogReportInfo struct {
	Name           string          `json:"name"`             // 报表名称
	ApiPrefix      string          `json:"api_prefix"`       // 数据源api前缀
	Label          string          `json:"label"`            // 报表标签
	ShowFieldIDs   []int           `json:"show_field_ids"`   // 显示字段ID列表
	IsSystemConfig int             `json:"is_system_config"` // 是否为系统配置
	IsExportable   int             `json:"is_exportable"`    // 是否可导出
	SearchFields   []DCSearchField `json:"search_fields"`    // 搜索字段列表
	DataSourceID   int             `json:"rc_datasource_id"` // 数据源ID
}
