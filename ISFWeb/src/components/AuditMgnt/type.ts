import { ListTipStatus } from '../ListTipComponent/helper'
import { GetDataReportByIdResp, DataReportListItem } from '@/core/apis/console/auditlog/types';
import __ from './locale';

export interface BizGroupItemDetail {
    id: string;
    name: string;
}

export type DataReportItemDetail = { dataReportItem: GetDataReportByIdResp | null } & { loadStatus: ListTipStatus };

export interface BizGroupItem {
    id: string;
    name: string;
    children: ReadonlyArray<DataReportListItem>;
}

export enum ActionType {
    None,

    CreateBizGroup,

    RenameBizGroup,

    DeleteBizGroup,

    CreateDataReport,

    EditDataReport,

    DeleteDataReport,

    ExportRecord,
}

export enum ValidateStatus {
    Normal,

    NameEmpty,

    BizGroupNameInvalid,

    BizGroupNameDuplicate,

    BizGroupNotExist,

    DataReportNameInvalid,

    DataReportNameDuplicate,

    ExportDataReportNameInvalid,

    ExportDataReportDuplicate,

    DataSourceNotExist,
}

export const ValidateMessages = {
    [ValidateStatus.NameEmpty]: __('此项不允许为空。'),
    [ValidateStatus.BizGroupNameInvalid]: __('名称不能包含 \\ / : * ? " < > | 特殊字符, 长度不能超过128个字符。'),
    [ValidateStatus.DataReportNameInvalid]: __('名称不能包含 \\ / : * ? " < > | 特殊字符, 长度不能超过128个字符。'),
    [ValidateStatus.ExportDataReportNameInvalid]: __('名称不能包含 \\ / : * ? " < > | 特殊字符, 长度不能超过256个字符。'),
    [ValidateStatus.BizGroupNameDuplicate]: __('业务组名称已存在。'),
    [ValidateStatus.DataReportNameDuplicate]: __('报表名称已存在。'),
    [ValidateStatus.BizGroupNotExist]: __('业务组已不存在。'),
    [ValidateStatus.DataSourceNotExist]: __('数据源已不存在。'),
}

export const DefaultSelectedDataReportInfo = {
    parentBizGroupId: "audit-log",
    parentBizGroupName: __('审计日志'),
    selectedDataReportId: "operation",
    selectedDataReportName: __('操作日志'),
}

/**
 * 组件tab类型
 */
export enum TabType {
    /**
     * 组织
     */
    Org = 'org',

    /**
     * 用户组
     */
    Group = 'group',

    /**
     * 匿名用户
     */
    Anonymous = 'anonymous',

    /**
     * 应用账户
     */
    App = 'app',
}