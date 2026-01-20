import * as React from 'react';
import { noop, includes, omit } from 'lodash';
import session from '@/util/session';
import { ShareMgnt } from '@/core/thrift';
import { NodeType, getNodeType, formatNodes, getSubUsers } from '@/core/organization';
import { getOrganizations } from '@/core/department/department';
import { getSubDepartments, getDepartmentOfUsersCount } from '@/core/thrift/sharemgnt'
import { ListTipStatus } from '../ListTipComponent/helper'
import WebComponent from '../webcomponent';
import { UsersLimit, DefaultPage } from './helper';
import __ from './locale';

interface Props extends React.Props<void> {
    /**
     * 用户id
     */
    userid: string;

    /**
     * 可选范围
     */
    selectType: Array<NodeType>;

    /**
     * 角色id
     */
    roleId?: string;

    /**
     * 选中节点时触发
     */
    onSelectionChange?: (node) => any;

    /**
     * 禁用节点时触发
     */
    getNodeStatus?: (node) => any;

    /**
     * 是否显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

    /**
     * 是否展示未分配组
     */
    isShowUndistributed?: boolean;

    /**
     * 禁用树
     */
    disabled: boolean;

    /**
     * 是否根据父节点状态禁用子节点
     */
    isDisableChildrenByParent: boolean;

     /**
      * 是否请求普通用户的组织架构
      */
     isRequestNormal?: boolean;
}

interface State {
    /**
     * 节点
     */
    nodes: ReadonlyArray<any>;

    /**
     * 提示状态
     */
    listTipStatus: ListTipStatus;
}

export default class OrganizationTreeBase extends WebComponent<Props, State> {
    static defaultProps = {
        userid: session.get('isf.userid'),
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        isShowDisabledUsers: true,
        isShowUndistributed: false,
        onSelectionChange: noop,
        getNodeStatus: () => ({ disabled: false }),
        disabled: false,
        isDisableChildrenByParent: false,
    }

    props: Props;

    state: State = {
        nodes: [],
        listTipStatus: ListTipStatus.Loading,
    }

    async componentDidMount() {
        await this.getOrganization()
    }

    /**
     * 获取组织结构信息
     */
    public getOrganization = async (): Promise<void> => {
        const {
            roleId,
            userid,
            isShowUndistributed,
            isRequestNormal,
        } = this.props

        let roots: any[] = []

        try {
            if (roleId) {
                roots = await ShareMgnt('UsrRolem_GetSupervisoryRootOrg', [userid || session.get('isf.userid'), roleId])
            } else {
                roots = await getOrganizations(isRequestNormal)

                if (isShowUndistributed) {
                    const count = await getDepartmentOfUsersCount(['-1'])

                    const undistributedNode = {
                        id: '-1',
                        name: __('未分配组'),
                        isUndistributedNode: true,
                        subUserCount: count,
                    }

                    roots = [undistributedNode, ...roots]
                }
            }
            this.setState({
                nodes: roots,
                listTipStatus: roots.length ? ListTipStatus.None : ListTipStatus.OrgEmpty,
            })
        } catch (ex) {
            this.setState({
                nodes: [],
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    /**
     * 加载子节点
     * @param node 节点
     */
    protected loadSubs = (node: any): Promise<void> => {
        node.children = [{ isLoading: true }]

        this.forceUpdate()

        return Promise.all([
            getSubDepartments([node.id]),
            includes(this.props.selectType, NodeType.USER) ?
                getSubUsers({
                    node,
                    page: DefaultPage,
                    limit: UsersLimit,
                    isShowDisabledUsers: this.props.isShowDisabledUsers as boolean,
                })
                : { users: [], page: DefaultPage },
        ]).then(([deps, { users, page }]) => {
            if (users.length > 0 && users.length < node.subUserCount) {
                node.children = [
                    ...users,
                    {
                        parentNode: node,
                        isLoading: false,
                        isLoadMore: true,
                        currentPage: page + 1,
                    },
                    ...deps,
                ]
            } else {
                node.children = [...users, ...deps]
            }
            this.addNodeStatusByChildren(node)

            this.forceUpdate()
        }).catch(() => {
            node.children = []

            this.addNodeStatusByChildren(node)

            this.forceUpdate()
        })
    }

    /**
     * 给节点添加状态
     */
    private addNodeStatusByChildren(node) {
        const nodeStatus = node.parent && node.parent.nodeStatus && node.parent.nodeStatus.disabled && this.props.isDisableChildrenByParent ?
            { disabled: true }
            : typeof this.props.getNodeStatus === 'function' && this.props.getNodeStatus(node)

        if (node.children) {
            node.children = node.children.map((n) => {
                this.addNodeStatusByChildren(n)

                if (n.hasOwnProperty('isLoading')) {
                    return n
                } else {
                    return {
                        ...n,
                        parent: { ...omit(node, 'children'), nodeStatus },
                    }
                }
            })
        }
    }

    /**
     * 分页加载更多用户
     * @param node 当前节点
     */
    protected loadMoreUsers = async (node) => {
        const { currentPage, parentNode } = node

        const users = parentNode.children.filter((item) => item.hasOwnProperty('user'))

        const departments = parentNode.children.filter((item) => !item.hasOwnProperty('user') && !item.hasOwnProperty('isLoadMore'))

        parentNode.children = [...users, { isLoading: true }, ...departments]

        this.forceUpdate()

        const { users: moreUsers, page } = await getSubUsers({
            node: parentNode,
            page: currentPage,
            limit: UsersLimit,
            isShowDisabledUsers: this.props.isShowDisabledUsers as boolean,
        })

        if ((moreUsers.length + (currentPage - 1) * UsersLimit) < parentNode.subUserCount) {
            parentNode.children = [
                ...users,
                ...moreUsers,
                { parentNode, isLoadMore: true, isLoading: false, currentPage: page + 1 },
                ...departments,
            ]
        } else {
            parentNode.children = [...users, ...moreUsers, ...departments]
        }

        this.addNodeStatusByChildren(parentNode)

        this.forceUpdate()
    }

    /**
     * 触发选中事件
     * @param selection 选中的节点数据
     */
    protected fireSelectionChangeEvent = async (node: any) => {
        const type = getNodeType(node)
        if (includes(this.props.selectType, type)) {
            const formatedNode = await formatNodes([node])
            this.props.onSelectionChange && this.props.onSelectionChange(Array.isArray(formatedNode) ? formatedNode[0] : formatedNode)
        }
    }
}