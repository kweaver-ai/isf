import { ListTipStatus } from '../../ListTipComponent/helper';
import { UserGroup } from '../helper';
import __ from './locale';

export interface MemberGridProps extends React.Props<void> {
    /**
     * 选中的用户组
     */
    selectedGroup: UserGroup.GroupInfo;

    /**
     * 用户组列表状态
     */
    groupStatus: ListTipStatus;

    /**
     * 更新用户组列表
     */
    onRequestUpdateGroup: () => void;
}

export interface MemberGridState {
    /**
     * 成员列表信息
     */
    data: {
        /**
         * 成员
         */
        members: ReadonlyArray<UserGroup.MemberInfo>;

        /**
         * 总数
         */
        total: number;

        /**
         * 当前页码
         */
        page: number;
    };

    /**
     * 搜索关键字
     */
    searchKey: string;

    /**
     * 列表状态
     */
    listTipStatus: ListTipStatus;

    /**
     * 选中项
     */
    selections: ReadonlyArray<UserGroup.MemberInfo>;

    /**
     * 是否显示添加成员列表
     */
    isShowAdd: boolean;

    /**
     * 搜索框是否聚焦
     */
    searchBoxIsOnFocus: boolean;

    /**
     * 搜索类型：用户显示名、用户组成员名
     */
    searchType: SearchType;
}

export enum SearchType {
    /**
     * 用户组成员
     */
    UserGroupMemberName = 1,

    /**
     * 用户显示名
     */
    UserDisplayNmae = 2,
}

export const SearchTypeText = {
    [SearchType.UserGroupMemberName]: __('搜索用户组成员'),
    [SearchType.UserDisplayNmae]: __('搜索用户'),
}