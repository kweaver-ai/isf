import __ from './locale';
import * as loadImg from './assets/loading.gif';
import * as Empty from './assets/empty.png';
import * as NoSearch from './assets/NoSearch.png';
import * as FailedImg from './assets/loadFailed.png';
import * as NoSyncPlan from './assets/noSyncPlan.png';
import * as NoExamine from './assets/noExamine.png';
import * as FileEmpty from './assets/file.png';
import * as ClientFileEmpty from './assets/file_client.png';
import * as NoDocFlow from './assets/noDocFlow.png';
import * as ClientNoSearch from './assets/noSearch_client.png';
import * as ClientLoadFaild from './assets/loadFaild_client.png';

/**
 * 列表提示
 */
export const enum ListTipStatus {
    /**
     * 不显示提示
     */
    None,

    /**
     * 提示加载中
     */
    Loading,

    /**
     * 列表为空
     */
    Empty,

    /**
     * 提示加载失败
     */
    LoadFailed,

    /**
     * 无匹配的搜索内容
     */
    NoSearchResults,

    /**
     * 无同步任务
     */
    NoSyncPlan,

    /**
     * 无审核任务
     */
    NoExamine,

    /**
     * 组织结构为空
     */
    OrgEmpty,

    /**
     * 组织结果/列表 为空-客户端
     */
    ClientOrgEmpty,

    /**
     * 白色背景的加载中
     */
    LightLoading,

    /**
     * 客户端搜索为空
     */
    ClientNoSearch,

    /**
     * 客户端加载失败
     */
    ClientLoadFaild,

    /**
     * 没有文档流转
     */
    NoDocFlow,
}

/**
 * 列表提示语
 */
export const ListTipMessage = {
    [ListTipStatus.Empty]: __('列表为空'),
    [ListTipStatus.NoSearchResults]: __('抱歉，没有找到符合条件的结果'),
    [ListTipStatus.LoadFailed]: __('加载失败'),
    [ListTipStatus.OrgEmpty]: __('暂无可选的用户或部门'),
    [ListTipStatus.ClientOrgEmpty]: __('暂无可选的用户或部门'),
    [ListTipStatus.ClientNoSearch]: __('抱歉，没有找到相关内容'),
    [ListTipStatus.ClientLoadFaild]: __('抱歉，无法完成加载'),
}

/**
 * 兼容客户端
 */
export const getListTipMsg = () => {
    return {
        [ListTipStatus.Empty]: __('列表为空'),
        [ListTipStatus.NoSearchResults]: __('抱歉，没有找到符合条件的结果'),
        [ListTipStatus.LoadFailed]: __('加载失败'),
        [ListTipStatus.OrgEmpty]: __('暂无可选的用户或部门'),
        [ListTipStatus.ClientOrgEmpty]: __('暂无可选的用户或部门'),
        [ListTipStatus.ClientNoSearch]: __('抱歉，没有找到相关内容'),
        [ListTipStatus.ClientLoadFaild]: __('抱歉，无法完成加载'),
    }
}

/**
 * 根据list和searchKey决定显示什么提示
 */
export function getTipStatus<T>(list: ReadonlyArray<T>, searchKey: string): ListTipStatus {
    if (!list.length) {
        // 如果list长度为0
        return searchKey.trim() ? ListTipStatus.NoSearchResults : ListTipStatus.Empty;
    }

    return ListTipStatus.None;
}

export const ListTipStatusMapImg = {
    [ListTipStatus.Empty]: Empty,
    [ListTipStatus.Loading]: loadImg,
    [ListTipStatus.NoSearchResults]: NoSearch,
    [ListTipStatus.LoadFailed]: FailedImg,
    [ListTipStatus.NoSyncPlan]: NoSyncPlan,
    [ListTipStatus.NoExamine]: NoExamine,
    [ListTipStatus.OrgEmpty]: FileEmpty,
    [ListTipStatus.ClientOrgEmpty]: ClientFileEmpty,
    [ListTipStatus.ClientNoSearch]: ClientNoSearch,
    [ListTipStatus.ClientLoadFaild]: ClientLoadFaild,
    [ListTipStatus.NoDocFlow]: NoDocFlow,
}