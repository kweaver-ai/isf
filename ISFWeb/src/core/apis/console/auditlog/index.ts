import { consolehttp } from '../../../openapiconsole'
import {
    GetDataReportDataList,
    GetDataReportFieldValuesList
} from './types'
/**
 * 获取报表数据列表
 */
export const getDataReportDataList: GetDataReportDataList = ({ id, ...data }, isHistory = false) => {
    return consolehttp('post', ['audit-log', 'v1', 'report-center', isHistory ? 'history' : 'active', id, 'list'], data, {})
}

/**
 * 获取报表字段值列表
 */
export const getDataReportFieldValuesList: GetDataReportFieldValuesList = ({ id, field, ...data }) => {
    return consolehttp('post', ['audit-log', 'v1', 'report-center', 'active', id, 'field', field, "values"], data, {})
}

/**
* 获取历史日志下载进度
*/
export const getGenCompressFileStatus = ({ taskid }) => {
    return consolehttp('get', ['audit-log', 'v1','history-log', 'download', taskid, 'progress'], {}, {})
}

/**
* 获取历史日志下载结果
*/
export const getDownloadResult = ({ taskid }) => {
    return consolehttp('get', ['audit-log', 'v1','history-log', 'download', taskid, 'result'], {}, {})
}

/**
 * 历史日志下载
 */
export const exportHistoryLog = (requestParams) => {
    return consolehttp('post', ['audit-log', 'v1','history-log', 'download', 'task'], requestParams, {})
}

// 转存周期单位
enum CycleUnit {
    Day = 'day',
    Week = 'week',
    Month = 'month',
    Year = 'year',
}

// 转存格式
enum DumpFormat {
    CSV = 'csv',
    XML = 'xml',
}

interface Configs {
    retention_period: number;
    retention_period_unit: CycleUnit;
    dump_time: string;
    dump_format: DumpFormat;
}

/**
 * 获取日志转存策略配置
 */
export const getLogStrategy = (): Promise<Configs> => {
    return consolehttp('get', ['audit-log', 'v1', 'log-strategy', 'dump'], {}, {})
}

/**
 * 更新日志转存策略配置
 */
export const updateLogStrategy = ({ field, retention_period, retention_period_unit, dump_time, dump_format }): Promise<void> => {
    return consolehttp('put', ['audit-log', 'v1', 'log-strategy', 'dump'], { retention_period, retention_period_unit, dump_time, dump_format }, { field })
}

/**
 * 获取历史日志下载是否需要加密
 */
export const getPasswordStatus = (): Promise<{status: boolean}> => {
    return consolehttp('get', ['audit-log', 'v1', 'history-log', 'download', 'pwdstatus'], {}, {})
}

/**
 * 更新历史日志下载是否需要加密
 */
export const updatePasswordStatus = ({ status }): Promise<void> => {
    return consolehttp('put', ['audit-log', 'v1', 'history-log', 'download', 'pwdstatus'], { status }, {})
}

/**
 * 获取日志查看范围策略配置
 */
export const getLogStrategyScope = ({ offset, limit, category }): Promise<any> => {
    return consolehttp('get', ['audit-log', 'v1', 'log-strategy', 'scope'], {}, { offset, limit, category })
}

/**
 * 更新日志查看范围策略配置
 */
export const updateLogStrategyScope = ({ id, type, category, role, scope }): Promise<any> => {
    return consolehttp('put', ['audit-log', 'v1', 'log-strategy', 'scope', id], { type, category, role, scope }, {})
}

/**
 * 新建日志查看范围策略配置
 */
export const addLogStrategyScope = ({ type, category, role, scope }): Promise<any> => {
    return consolehttp('post', ['audit-log', 'v1', 'log-strategy', 'scope'], { type, category, role, scope }, {})
}

/**
 * 删除日志查看范围策略配置
 */
export const deleteLogStrategyScope = (id) => {
    return consolehttp('delete', ['audit-log', 'v1', 'log-strategy', 'scope', id], {}, {})
}
