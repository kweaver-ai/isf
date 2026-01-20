import * as React from 'react';
import { noop, includes } from 'lodash';
import session from '@/util/session';
import { NodeType, FormatedNodeInfo } from '@/core/organization';
import { searchInOrgTree } from '@/core/apis/console/organization'
import { OrgType } from '@/core/apis/console/organization/types'
import { SysUserRoles, getRoleType, UserRole } from '@/core/role/role'
import { getNodeType } from '@/core/organization/organization';
import WebComponent from '../webcomponent';
import __ from './locale'

interface Props {
    /**
     * 关键字
     */
    value?: string;

    /**
     * 宽度
     */
    width?: string | number;

    /**
     * 选择搜索范围
     */
    selectType?: ReadonlyArray<NodeType>;

    /**
     * 是否显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

    /**
     * 是否显示未分配的用户
     */
    isShowUndistributed?: boolean;

    /**
     * 选择完自动填充搜索框
     */
    completeSearchKey?: boolean;

    /**
     * 自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * 检测 输入框的值是否发生改变
     */
    onValueChange?: (value: string) => boolean;

    /**
     * 选择事件
     */
    onSelectDep: (sharer: Record<string, any>) => any;

    /**
     * 角色id
     */
    roleId?: string;

    /**
     * 是否搜索框能够输入
     */
    canInput?: boolean;
}

interface State {
    // 搜索结果
    results: ReadonlyArray<FormatedNodeInfo>;
    // 搜索关键字
    searchKey: string;
    // 是否能够输入
    canInputValue: boolean;
}

export default class SearchDepBase extends WebComponent<Props, State> {
    static defaultProps = {
        isShowDisabledUsers: true,
        isShowUndistributed: false,
        onSelectDep: noop,
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
    }

    state: State = {
        results: [],

        searchKey: '',

        canInputValue: false,
    }

    autocomplete: HTMLInputElement;

    componentDidMount() {
        this.setState({
            searchKey: this.props.value || '',
        })
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.canInput !== prevState.canInputValue) {
            return {
                canInputValue: nextProps.canInput,
            }
        }
        return null
    }

    /**
     * 根据key获取部门
     * @return 部门数组
     */
    async getDepsByKey(key: string, start: number = 0): Promise<FormatedNodeInfo[]> {
        if (key) {
            const { user: { roles } } = session.get('isf.userInfo')
            const roleIds = roles.map((role) => role.id)
            const { roleId, selectType, isShowDisabledUsers, isShowUndistributed } = this.props
            const maxRole = getRoleType()

            const { users, departments } = await searchInOrgTree({
                keyword: key,
                role: [UserRole.Super, UserRole.Admin, UserRole.Security, UserRole.Audit].includes(maxRole) ?
                    maxRole
                    : (roleId && roleIds.includes(roleId)) ?
                        SysUserRoles[roleId]
                        : maxRole,
                type: [
                    ...(includes(selectType, NodeType.DEPARTMENT) || includes(selectType, NodeType.ORGANIZATION)) ? [OrgType.Dep] : [],
                    ...(includes(selectType, NodeType.USER) ? [OrgType.User] : []),
                ],
                ...(!isShowDisabledUsers ? { user_enabled: true } : {}),
                ...(!isShowUndistributed ? { user_assigned: true } : {}),
                offset: start,
                limit: 10,
            })

            return [
                ...(
                    departments ?
                        departments.entries.map((dep) => ({
                            id: dep.id,
                            name: dep.name,
                            type: getNodeType(dep),
                            parent_path: dep.parent_dep_path,
                            original: dep,
                        }))
                        : []
                ),
                ...(
                    users ?
                        users.entries.map((user) => ({
                            id: user.id,
                            name: user.name,
                            account: user.account,
                            type: getNodeType(user),
                            parent_path: Array.isArray(user.parent_dep_paths) ?
                                user.parent_dep_paths.length ?
                                    user.parent_dep_paths[0]
                                    : __('未分配组')
                                : '',
                            original: user,
                        }))
                        : []
                ),
            ]

        } else {
            return []
        }

    }

    /**
     * 获取搜索到的结果
     * @param results 部门数组
     */
    getSearchData(results: ReadonlyArray<FormatedNodeInfo>) {
        this.setState({
            results,
        })
    }

    /**
     * 选择搜索到的单个部门
     * @param dep 部门
     */
    selectItem(item: Record<string, any>): void {
        this.props.onSelectDep(item);
        if (this.props.completeSearchKey && item.name) {
            this.setState({
                searchKey: item.name,
            })
        }
        this.autocomplete.toggleActive(false);
    }

    handelChange(searchKey) {
        this.setState({
            searchKey,
        })
        this.props.onValueChange ? this.props.onValueChange(searchKey) : null
    }
    /**
     * 按下enter
     */
    handleEnter(e, selectIndex: number) {
        if (selectIndex >= 0) {
            this.selectItem(this.state.results[selectIndex])
        }
    }

    /**
     * 懒加载
     */
    protected lazyLoade = async (page: number, limit: number): Promise<void> => {
        this.setState({
            results: [
                ...this.state.results,
                ...(await this.getDepsByKey(this.state.searchKey, (page - 1) * limit)),
            ],
        })
    }
}
