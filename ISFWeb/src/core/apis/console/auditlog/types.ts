import { OpenAPI } from '../index'

/**
 * 获取报表数据源列表参数
 */
export interface GetReportCenterDataSourceListPrams {
    /**
     * 数据源组ID
     */
    datasource_group_id?: number;

    /**
     * 数据源名称(模糊搜索)
     */
    name?: string;

    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit: number;

    /**
     * 分页偏移量
     */
    offset: number;
}

/**
 * 获取报表数据源列表响应
 */
interface GetReportCenterDataSourceListResp {
    /**
     * 记录总数
     */
    total_count: number;

    /**
     * 数据记录
     */
    entries: ReadonlyArray<ReportCenterDataSourceObj>;
}

/**
 * 获取报表数据源列表、获取报表数据源详情、新建数据源、编辑数据源
 */
export interface ReportCenterDataSourceObj {
    /**
     * ID
     */
    id?: number;

    /**
     * 创建时间
     */
    created_at?: string;

    /**
     * 创建人
     */
    created_by?: string;

    /**
     * 修改时间
     */
    updated_at?: string;

    /**
     * 创建人
     */
    updated_by?: string;

    /**
     * 数据源名称
     */
    name: string;

    /**
     * 数据源组ID
     */
    datasource_group_id: number;

    /**
     * 数据源组名称
     */
    datasource_group_name?: string;

    /**
     * API前缀
     */
    api_prefix: string;

    /**
     * 数据源字段
     */
    rc_datasource_fields?: ReadonlyArray<DataSourceFields>;

    /**
     * 默认排序字段
     */
    default_sort_field?: string;

    /**
     * 默认排序方向:"asc" "desc"
     */
    default_sort_direction?: DefaultSortDirection;

    /**
     * 唯一且递增字段，此字段用于提高分页查询性能
     */
    unique_incremental_field?: string;

    /**
     * 报表数据记录ID，用于标识唯一一条报表数据记录
     */
    id_field: string;
}

/**
 * 数据源字段
 */
export interface DataSourceFields {
    /**
     * 数据源字段ID
     */
    id?: number;

    /**
     * 字段
     */
    field: string;

    /**
     * 字段标题
     */
    field_title: string;

    /**
     * 自定义字段标题(获取到的元数据没有，新建时前端自动加入)
     */
    field_title_custom?: string;

    /**
     * 是否为键值对字段: 0-否，1-是
     */
    is_kv_field: number;

    /**
     * 是否可排序: 0-否，1-是
     */
    is_can_sort: number;

    /**
     * 是否可搜索字段: 0-否，1-是
     */
    is_can_search?: number;

    /**
     * 支持的搜索字段类型数组
     * 索字段类型（1：文本框，2：下拉列表，3：文本范围，4：日期范围，5：组织结构）
     */
    search_field_support_type?: ReadonlyArray<number>;

    /**
     * 显示类型（1：文本，2：链接，3：图片）
     */
    show_type?: number;

    /**
     * 是否权限控制字段: 0-否，1-是
     */
    is_pms_ctrl_field?: number;

    /**
     * 是否组织架构字段: 0-否，1-是
     */
    is_org_structure_field?: number;

    /**
     * 组织架构字段配置
     */
    org_structure_field_config?: OrgStructureFieldConfig;

    /**
     * 搜索字段配置
     */
    search_field_config: SearchFieldConfig;
}

/**
 * 组织架构字段配置
 */
export interface OrgStructureFieldConfig {
    /**
     * 选择类型
     * 使用四位二进制表示（后面可以根据需要进行扩展，比如扩展成5位二进制），从低位到高位分别代表：1：人员，2：部门，3：用户组，4：联系人。转成10进制，分别对应：1（0001）：可选人员； 2（0010）：可选部门；3（0011）可选人员和部门；4（0100）：可选用户组；5（0101）：可选人员和用户组；6（0110）：可选部门和用户组；7（0111）：可选人员、部门和用户组；8（1000）：可选联系人；9（1001）：可选人员和联系人；10（1010）：可选部门和联系人；11（1011）：可选人员、部门和联系人；12（1100）：可选用户组和联系人；13（1101）：可选人员、用户组和联系人；14（1110）：可选部门、用户组和联系人；15（1111）：可选人员、部门、用户组和联系人
     * 传递对应的10进制数值即可
     */
    select_type: number;

    /**
     * 是否多选: 0-否，1-是
     */
    is_multiple: number;
}

/**
 * 获取报表数据源参数
 */
export interface GetDataSourceGroupParams {
    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit: number;

    /**
     * 分页偏移量
     */
    offset: number;
}

/**
 * 获取报表数据源分组列表响应
 */
interface GetDataSourceGroupResp {
    /**
     * 记录总数
     */
    total_count: number;

    /**
     * 数据记录
     */
    entries: ReadonlyArray<DataSourceGroupObj>;
}

/**
 * 获取报表数据源分组信息
 */
export interface DataSourceGroupObj {
    /**
     * ID
     */
    id: number;

    /**
     * 创建时间
     */
    created_at?: string;

    /**
     * 创建人
     */
    created_by?: string;

    /**
     * 修改时间
     */
    updated_at?: string;

    /**
     * 创建人
     */
    updated_by?: string;

    /**
     * 数据源分类名称
     */
    name: string;

    /**
     * 是否默认分组: 0-否，1-是
     */
    is_default?: number;
}

/**
 * 获取报表数据源字段列表响应
 */
interface GetReportCenterDataSourceFieldsListByIdResp {
    /**
     * 记录总数
     */
    total_count: number;

    /**
     * 数据记录
     */
    entries: ReadonlyArray<DataSourceFieldsItem>;
}

/**
 * 数据源字段列表项信息
 */
export interface DataSourceFieldsItem {
    field_id: number;
    field: string;
    field_title_custom: string;
}

/**
 * 获取数据源元数据
 */
interface GetReportCenterDataSourceMetadataResp {
    /**
     * 默认排序方向:"asc" "desc"
     */
    default_sort_direction: DefaultSortDirection;

    /**
     * 默认排序字段
     */
    default_sort_field: string;

    /**
     * 唯一且递增字段，此字段用于提高分页查询性能
     */
    unique_incremental_field: string;

    /**
     * 字段列表
     */
    fields: ReadonlyArray<DataSourceFields>;
}

/**
 * 获取业务组列表参数
 */
interface GetBizGroupListParams {
    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit: number;

    /**
     * 分页偏移量
     */
    offset: number;
}

/**
 * 获取业务组列表响应
 */
interface GetBizGroupListResp {
    total_count: number;

    entries: ReadonlyArray<BizGroupListItem>;
}

/**
 * 业务组列表项信息
 */
export interface BizGroupListItem {
    /**
     * ID
     */
    id: number;

    /**
     * 创建时间
     */
    created_at: string;

    /**
     * 创建人
     */
    created_by: string;

    /**
     * 修改时间
     */
    updated_at: string;

    /**
     * 修改人
     */
    updated_by: string;

    /**
     * 名称
     */
    name: string;
}

/**
 * 创建业务组参数
 */
interface CreateBizGroupParams {
    /**
     * 名称
     */
    name: string;
}

/**
 * 创建业务组响应
 */
interface CreateBizGroupResp {
    /**
     * ID
     */
    id: number;
}

/**
 * 获取业务组详情参数
 */
interface GetBizGroupByIdParams {
    /**
     * ID
     */
    id: number;
}

/**
 * 删除业务组参数
 */
interface DeleteBizGroupParams {
    /**
     * ID
     */
    id: number;
}

/**
 * 更新业务组参数
 */
interface UpdateBizGroupParams {
    /**
     * ID
     */
    id: number;

    /**
     * 名称
     */
    name: string;
}

/**
 * 获取报表列表参数
 */
interface GetDataReportListParams {
    /**
     * 业务组ID
     */
    biz_group_id: number;

    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit: number;

    /**
     * 分页偏移量
     */
    offset: number;
}

/**
 * 获取报表列表响应
 */
interface GetDataReportListResp {
    total_count: number;

    entries: ReadonlyArray<DataReportListItem>;
}

/**
 * 报表列表项信息
 */
export interface DataReportListItem {
    /**
     * ID
     */
    id: number;

    /**
     * 报表名称
     */
    name: string;
}

/**
 * 创建报表参数
 */
interface CreateDataReportParams {
    /**
     * 报表名称
     */
    name: string;

    /**
     * 报表数据源ID
     */
    rc_datasource_id: number;

    /**
     * 显示字段ID数据
     */
    show_field_ids: ReadonlyArray<number>;

    /**
     * 搜索字段
     */
    search_fields: ReadonlyArray<SearchFieldsItem>;

    /**
     * 业务组ID
     */
    biz_group_id: number;
}

/**
 * 搜索字段项信息
 */
interface SearchFieldsItem {
    /**
     * 字段ID
     */
    rc_field_id: number;

    /**
     * 字段类型
     */
    field_type: FieldType;

    /**
     * 是否必填
     */
    is_required: IsRequired;
}

/**
 * 字段类型
 */
export enum FieldType {
    None = 0,
    Text = 1,
    Select = 2,
    TextRange = 3,
    DateRange = 4,
    Org = 5,
}

/**
 * 是否必填
 */
export enum IsRequired {
    No = 0,
    Yes = 1,
}

/**
 * 创建报表响应
 */
interface CreateDataReportResp {
    /**
     * ID
     */
    id: number;
}

/**
 * 获取报表详情参数
 */
interface GetDataReportByIdParams {
    /**
     * ID
     */
    id: number;
}

/**
 * 获取报表详情响应
 */
export interface GetDataReportByIdResp {
    /**
     * ID
     */
    id: number;

    /**
     * 报表名称
     */
    name: string;

    /**
     * 报表数据源ID
     */
    rc_datasource_id: number;

    /**
     * 报表数据源名称
     */
    rc_datasource_name: string;

    /**
     * 创建时间
     */
    created_at: string;

    /**
     * 创建人
     */
    created_by: string;

    /**
     * 修改时间
     */
    updated_at: string;

    /**
     * 修改人
     */
    updated_by: string;

    /**
     * 显示字段ID数据
     */
    show_field_ids: ReadonlyArray<number>;

    /**
     * 搜索字段
     */
    search_fields: ReadonlyArray<SearchFieldsItem>;
}

/**
 * 删除报表参数
 */
interface DeleteDataReportParams {
    /**
     * ID
     */
    id: number;
}

/**
 * 更新报表参数
 */
interface UpdateDataReportParams {
    /**
     * ID
     */
    id: number;

    /**
     * 报表名称
     */
    name: string;

    /**
     * 显示字段ID数据
     */
    show_field_ids: ReadonlyArray<number>;

    /**
     * 搜索字段
     */
    search_fields: ReadonlyArray<SearchFieldsItem>;
}

/**
 * 获取报表配置信息参数
 */
interface GetDataReportConfigParams {
    /**
     * ID
     */
    id: number;
}

/**
 * 获取报表配置信息响应
 */
interface GetDataReportConfigResp {
    /**
     * 报表数据源ID
     */
    rc_datasource_id: number;

    /**
     * 显示字段数组
     */
    show_fields: ReadonlyArray<ShowFieldsItem>;

    /**
     * 搜索字段
     */
    search_fields: ReadonlyArray<SearchFieldsItem2>;

    /**
     * 数据源配置
     */
    datasource_config: DatasourceConfig;
}

export interface SearchFieldsItem2 {
    id: number;

    field: string;

    field_title_custom: string;

    field_type: FieldType;

    search_is_required: IsRequired;

    search_field_config: SearchFieldConfig;

    org_structure_field_config: OrgStructureFieldConfig;
}

export interface SearchFieldConfig {
    is_can_search_by_api: boolean;

    dependent_fields: ReadonlyArray<DependentField>;

    search_label: string;

    support_types?: ReadonlyArray<number>;
}

export interface DependentField {
    field: string;

    is_must: boolean;
}

export interface DatasourceConfig {
    /**
     * 表示id的字段名
     */
    id_field: string;

    /**
     * 默认排序字段
     */
    default_sort_field: string;

    /**
     * 默认排序方向
     */
    default_sort_direction: DefaultSortDirection;

    /**
     * 唯一且递增字段，此字段用于提高分页查询性能
        此字段需要满足下面的条件
        1、唯一且递增
        2、建立了索引（比如主键）
        当提供此字段时
        报表中心将在调用“/xx/list”接口时，在order_by中自动加上此字段的last_field_value。此时不再提供offset参数
        此时业务方可使用unique_incremental_field>unique_incremental_field_last_value limit xxx的方式进行分页
        如果排序字段不是此字段，业务方可将此字段作为第二个排序字段
        示例：select xx from table where sort_field>=sort_field_last_value and id>last_id order by sort_field xx, id xx limit xxx
        这里面id为unique_incremental_field，sort_field为default_sort_field或用户指定的排序字段
        sort_field>=sort_field_last_value 或者 sort_field>sort_field_last_value 由业务方自己决定，主要考虑sort_field的值是否可能重复
        当sort_field也为唯一且递增字段时，sql可以为：select xx from table where sort_field>sort_field_last_value order by sort_field xx limit xxx（业务方自行判断）
        当没有提供此字段时，报表中心将在调用“/xx/list”接口时提供offset字段。业务方可以使用offset和limit进行分页（这样分页当offset比较大时，性能可能会比较差）
     */
    unique_incremental_field: string;
}

export enum DefaultSortDirection {
    Asc = 'asc',

    Desc = 'desc',
}

/**
 * 显示字段数组项信息
 */
export interface ShowFieldsItem {
    /**
     * 数据源字段ID
     */
    id: number;

    /**
     * 字段标识
     */
    field: string;

    /**
     * 自定义字段标题
     */
    field_title_custom: string;

    /**
     * 是否可排序
     */
    is_can_sort: IsCanSort;

    /**
     * 显示类型
     */
    show_type: ShowType;
}

/**
 * 是否可排序
 */
enum IsCanSort {
    No = 0,
    Yes = 1,
}

/**
 * 显示类型
 */
export enum ShowType {
    Text = 1,
    Link = 2,
    Img = 3,
    Date = 4,
    Minute = 5,
    Time = 6,
    Download = 7,
    DcpiPmsFields = 101,
}

/**
 * 获取报表数据列表参数
 */
interface GetDataReportDataListParams {
    /**
     * 数据源ID
     */
    id: number;

    /**
     * 每页显示数量
     */
    limit: number;

    /**
     * 分页偏移量（当order_by中数组第一个元素last_field_value字段存在时，此字段不存在）
     */
    offset: number;

    /**
     * 搜索条件（key：搜索字段，val：搜索字段值）
     * 关于搜索条件的说明
        当搜索字段类型为文本框时，此字段为文本框输入的值
        类型可为string或number
        当搜索字段类型为下拉列表时，此字段为下拉列表选中的值
        类型可为string或number
        当搜索字段类型为文本范围时，此字段为文本范围输入的值
        类型为数组，数组第一个元素为开始值，第二个元素为结束值
        当搜索字段类型为日期范围时，此字段为日期组件选中的值
        类型为数组，数组第一个元素为开始值，第二个元素为结束值
        当搜索字段类型为组织结构时，此字段为组织结构组件选中的值
        当单选时，类型为string或number
        当多选时，类型为数组，数组元素类型为string或number
     */
    condition?: Record<string, string | number | [string | number, string | number] | (string | number)[]>;

    /**
     * 排序字段
        当用户没有主动排序时，为服务API的默认排序字段和默认排序方向（业务服务metadata接口中提供）
        数组的第一个元素为用户指定的排序字段或默认排序字段
        如果数组的第一个元素对应的字段不是unique_incremental_field，那么此时数组的第二个元素将为unique_incremental_field
     */
    order_by?: {
        /**
         * 字段标识
         */
        field: string;

        /**
         * 排序方向
         */
        direction: DefaultSortDirection;

        /**
         * 上一页最后一条记录的当前字段对应的字段值
         */
        last_field_value?: string;
    }[];
}

/**
 * 获取报表数据列表响应
 */
interface GetDataReportDataListResp {
    /**
     * 数据记录
     */
    entries: ReadonlyArray<Record<string, string | number>>;
}

/**
 * 获取报表字段值列表参数
 */
interface GetDataReportFieldValuesListParams {
    /**
     * 数据源ID
     */
    id: number;

    /**
     * 字段标识
     */
    field: string;

    /**
     * 每页显示数量
     */
    limit?: number;

    /**
     * 分页偏移量
     */
    offset: number;

    /**
     * 查询条件（key：条件字段，val：条件字段值）
     * 该字段是否生效取决于getDataReportConfig接口返回结果中search_fields数组中对应字段的search_field_config字段中的dependent_fields字段是否为空
     */
    condition?: Record<string, string | number | (string | number)[]>;

    /**
     * 搜索关键字
     */
    keyword?: string;
}

/**
 * 获取报表字段值列表响应
 */
interface GetDataReportFieldValuesListResp {
    /**
     * 字段值列表
     */
    entries: ReadonlyArray<{
        value_code: string | number;
        value_name: string;
    }>;

    total_count: number;
}

/**
 * 获取下载任务列表参数
 */
interface GetDataReportDownloadTaskListParams {
    /**
     * 报表ID
     */
    report_id?: number;

    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit: number;

    /**
     * 分页偏移量
     */
    offset: number;

    /**
     * 任务名称
     */
    name?: string;
}

/**
 * 获取下载任务列表响应
 */
interface GetDataReportDownloadTaskListResp {
    /**
     * 记录总数
     */
    total_count: number;

    /**
     * 数据记录
     */
    entries: ReadonlyArray<DataReportDownloadTaskListItem>;
}

/**
 * 下载任务列表项信息
 */
export interface DataReportDownloadTaskListItem {
    /**
     * 任务ID
     */
    id: number;

    /**
     * 任务名称
     */
    name: string;

    /**
     * 报表ID
     */
    report_id: number;

    /**
     * 报表名称
     */
    report_name: string;

    /**
     * 报表业务组名称
     */
    biz_group_name: string;

    /**
     * 任务状态
     */
    status: Status;

    /**
     * 当status 为2时计算耗时（单位：秒）
     */
    last_exec_start_at?: number;
}

/**
 * 任务状态
 */
export enum Status {
    /**
     * 待执行
     */
    Pending = 1,

    /**
     * 执行中
     */
    Running = 2,

    /**
     * 已完成
     */
    Success = 3,

    /**
     * 执行失败
     */
    Failed = 4,

    /**
     * 下载文件已失效
     */
    DownloadFileInvalid = 5,

    /**
     * 空（前端初始化使用）
     */
    None = 6,
}

/**
 * 获取下载任务详情参数
 */
interface GetDataReportDownloadTaskDetailByIdParams {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 获取下载任务详情响应
 */
export interface GetDataReportDownloadTaskDetailByIdResp {
    /**
     * 任务名称
     */
    name: string;

    /**
     * 业务组名称
     */
    biz_group_name: string;

    /**
     * 报表名称
     */
    report_name: string;

    /**
     * 过滤条件
     */
    dl_payload_human_json: string; // JSON字符串，前端展示用，需要具有很好的可读性

    /**
     * 任务状态
     */
    status: Status;

    /**
     * 失败原因
     */
    fail_msg: string;

    /**
     * 最后一次执行开始时间
     */
    last_exec_start_at: number; // <date-time>

    /**
     * 最后一次执行结束时间
     */
    last_exec_end_at?: number; // <date-time>

    /**
     * 导出类型：1.全部、2.自定义、3.用户选择数据
     * 当值为1时，不传下面的 ids 和 max_number
     * 当值为2时，传 max_number
     * 当值为3时，传 ids
     */
    export_type: number;

    /**
     * 导出的数据id列表
     */
    ids?: ReadonlyArray<number>;

    /**
     * 导出的数据的长度（与ids对应）
     */
    select_number?: number;

    /**
     * 导出前 x 条数据
     */
    max_number?: number;
}

/**
 * 创建下载任务参数
 */
interface CreateDataReportDownloadTaskParams {
    /**
     * 任务名称
     */
    name?: string;

    /**
     * 报表ID
     */
    report_id: number;

    /**
     * 过滤条件
     */
    dl_payload_json: string; // JSON字符串

    /**
     * 过滤条件
     */
    dl_payload_json_human: string; // JSON字符串，前端展示用，需要具有很好的可读性

    /**
     * 报表记录id数据，当使用选择数据进行导出，需要传递该字段
     */
    ids?: (string | number)[];

    /**
     * 导出前 x 条数据，当使用自定义数据量导出时，需要传递该字段
     */
    max_number?: number;
}

/**
 * 创建下载任务响应
 */
interface CreateDataReportDownloadTaskResp {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 取消下载任务参数
 */
interface CancelDataReportDownloadTaskParams {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 删除下载任务参数
 */
interface DeleteDataReportDownloadTaskParams {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 获取下载文件地址参数
 */
interface GetDataReportDownloadFileUrlParams {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 获取下载文件地址响应
 */
interface GetDataReportDownloadFileUrlResp {
    /**
     * 下载文件ID
     */
    file_id: string;

    /**
     * oss id
     */
    oss_id: string;
}

/**
 * 重新执行下载任务参数
 */
interface ReExecuteDataReportDownloadTaskParams {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 重新执行下载任务响应
 */
interface ReExecuteDataReportDownloadTaskResp {
    /**
     * 任务ID
     */
    id: number;
}

/**
 * 获取数据源元数据
 */
interface GetReportCenterDataSourceMetadataResp {
    /**
     * 默认排序方向:"asc" "desc"
     */
    default_sort_direction: DefaultSortDirection;

    /**
     * 默认排序字段
     */
    default_sort_field: string;

    /**
     * 唯一且递增字段，此字段用于提高分页查询性能
     */
    unique_incremental_field: string;

    /**
     * 字段列表
     */
    fields: ReadonlyArray<DataSourceFields>;

    /**
     * 报表数据记录ID，用于标识唯一一条报表数据记录
     */
    id_field: string;
}

/**
 * 获取权限控制字段详情
 */
interface GetDcpiPmsFieldsParams {
    /**
     * 数据控制权限项编号
     */
    code: string;
}

/**
 * 获取权限控制字段详情
 */
interface GetDcpiPmsFieldsResp {
    /**
     * 权限字段配置JSON字符串
     */
    pms_fields: string;
}

/** ================================================== 函数声明 ======================================================== */

/**
 * 获取数据源分组列表
 */
export type GetDataSourceGroupList = OpenAPI<GetDataSourceGroupParams, GetDataSourceGroupResp>;

/**
 * 新建数据源分组列表
 */
export type CreateDataSourceGroup = OpenAPI<{ name: string }, { id: number }>;

/**
 * 修改数据源分组
 */
export type EditDataSourceGroupById = OpenAPI<{ id: number; name: string }, void>;

/**
 * 删除数据源分组
 */
export type DeleteDataSourceGroupById = OpenAPI<{ id: number }, void>;

/**
 * 获取数据源分组详情
 */
export type GetDataSourceGroupDetailById = OpenAPI<{ id: number }, DataSourceGroupObj>;

/**
 * 获取报表数据源列表
 */
export type GetReportCenterDataSourceList = OpenAPI<GetReportCenterDataSourceListPrams, GetReportCenterDataSourceListResp>;

/**
 * 获取报表数据源详情
 */
export type GetReportCenterDataSourceDetailById = OpenAPI<{ id: number }, ReportCenterDataSourceObj>;

/**
 * 获取报表数据源字段列表
 */
export type GetReportCenterDataSourceFieldsListById = OpenAPI<{
    /**
     * 数据源ID
     */
    id: number;

    /**
     * 每页显示数量（默认值是10），最大值是1000
     */
    limit?: number;

    /**
     * 分页偏移量
     */
    offset: number;

    /**
     * 自定义字段标题（支持模糊搜索）
     */
    field_title_custom?: string;
}, GetReportCenterDataSourceFieldsListByIdResp>;

/**
 * 新建报表数据源
 */
export type CreateReportCenterDataSource = OpenAPI<ReportCenterDataSourceObj, { id: number }>;

/**
 * 获取报表数据源元数据
 */
export type GetReportCenterDataSourceMetadata = OpenAPI<{ api_prefix: string }, GetReportCenterDataSourceMetadataResp>;

/**
 * 修改报表数据源
 */
export type EditReportCenterDataSourceById = OpenAPI<ReportCenterDataSourceObj, void>;

/**
 * 删除报表数据源
 */
export type DeleteReportCenterDataSourceById = OpenAPI<{ id: number }, void>;

/**
 * 获取业务组列表
 */
export type GetBizGroupList = OpenAPI<GetBizGroupListParams, GetBizGroupListResp>;

/**
 * 创建业务组
 */
export type CreateBizGroup = OpenAPI<CreateBizGroupParams, CreateBizGroupResp>;

/**
 * 获取业务组详情
 */
export type GetBizGroupById = OpenAPI<GetBizGroupByIdParams, BizGroupListItem>;

/**
 * 删除业务组
 */
export type DeleteBizGroup = OpenAPI<DeleteBizGroupParams, void>;

/**
 * 更新业务组
 */
export type UpdateBizGroup = OpenAPI<UpdateBizGroupParams, void>;

/**
 * 获取报表列表
 */
export type GetDataReportList = OpenAPI<GetDataReportListParams, GetDataReportListResp>;

/**
 * 创建报表
 */
export type CreateDataReport = OpenAPI<CreateDataReportParams, CreateDataReportResp>;

/**
 * 获取报表详情
 */
export type GetDataReportById = OpenAPI<GetDataReportByIdParams, GetDataReportByIdResp>;

/**
 * 删除报表
 */
export type DeleteDataReport = OpenAPI<DeleteDataReportParams, void>;

/**
 * 更新报表
 */
export type UpdateDataReport = OpenAPI<UpdateDataReportParams, void>;

/**
 * 获取报表配置信息
 */
export type GetDataReportConfig = OpenAPI<GetDataReportConfigParams, GetDataReportConfigResp>;

/**
 * 获取报表数据列表
 */
export type GetDataReportDataList = OpenAPI<GetDataReportDataListParams, GetDataReportDataListResp>;

/**
 * 获取报表字段值列表
 */
export type GetDataReportFieldValuesList = OpenAPI<GetDataReportFieldValuesListParams, GetDataReportFieldValuesListResp>;

/**
 * 获取下载任务列表
 */
export type GetDataReportDownloadTaskList = OpenAPI<GetDataReportDownloadTaskListParams, GetDataReportDownloadTaskListResp>;

/**
 * 获取下载任务详情
 */
export type GetDataReportDownloadTaskDetailById = OpenAPI<GetDataReportDownloadTaskDetailByIdParams, GetDataReportDownloadTaskDetailByIdResp>;

/**
 * 创建下载任务
 */
export type CreateDataReportDownloadTask = OpenAPI<CreateDataReportDownloadTaskParams, CreateDataReportDownloadTaskResp>;

/**
 * 取消下载任务
 */
export type CancelDataReportDownloadTask = OpenAPI<CancelDataReportDownloadTaskParams, void>;

/**
 * 删除下载任务
 */
export type DeleteDataReportDownloadTask = OpenAPI<DeleteDataReportDownloadTaskParams, void>;

/**
 * 获取下载文件地址
 */
export type GetDataReportDownloadFileUrl = OpenAPI<GetDataReportDownloadFileUrlParams, GetDataReportDownloadFileUrlResp>;

/**
 * 重新执行下载任务（复制一条新的下载任务，状态为待执行）
 */
export type ReExecuteDataReportDownloadTask = OpenAPI<ReExecuteDataReportDownloadTaskParams, ReExecuteDataReportDownloadTaskResp>;

/**
 * 获取权限控制字段详情
 */
export type GetDcpiPmsFields = OpenAPI<GetDcpiPmsFieldsParams, GetDcpiPmsFieldsResp>;