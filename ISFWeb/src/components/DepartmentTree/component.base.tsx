import * as React from 'react'
import { includes, noop } from 'lodash';
import { getSubDepartments } from '@/core/thrift/sharemgnt/sharemgnt';
import { CascadeDirection, SelectType as NodeSelectType } from '@/ui/Tree2/ui.base';
import { getOrganizations } from '@/core/department/department';
import { getConfidentialConfig } from '@/core/apis/eachttp/config/config';
import { NodeType, isLeaf, ExtraRoot, getSubUsers, formatNodes } from '@/core/organization';
import { ListTipStatus } from '../ListTipComponent/helper';
import { UsersLimit, DefaultPage } from './helper';

interface DepartmentTreeProps {
    /**
     * 节点选择类型
     */
    nodeSelectType: NodeSelectType;

    /**
     * 组织树选择的类型（组织，部门，用户）
     */
    selectType: ReadonlyArray<NodeType>;

    /**
     * 树是否禁用
     */
    disabled: boolean;

    /**
     * 级联方向
     */
    cascadeDirection: CascadeDirection;

    /**
     * 额外根节点
     */
    extraRoots: ReadonlyArray<ExtraRoot>;

    /**
     * 是否显示被禁用的用户
     */
    isShowDisabledUsers?: boolean;

     /**
      * 是否请求普通用户的组织架构
      */
     isRequestNormal?: boolean;

    /**
     * 选中项改变
     */
    onSelectionChange: (selecctions: ReadonlyArray<any>, select: object, id: string) => void;
}

interface DepartmentTreeState {
    /**
     * 根节点
     */
    root: ReadonlyArray<any>;

    /**
     * 列表提示状态
     */
    listTipStatus: ListTipStatus;
}

export default class DepartmentTree extends React.PureComponent<DepartmentTreeProps, DepartmentTreeState> {
    static defaultProps = {
        selectType: [NodeType.ORGANIZATION, NodeType.DEPARTMENT],
        disabled: false,
        cascadeDirection: CascadeDirection.DOUBLE_SIDED,
        extraRoots: [],
        isShowDisabledUsers: true,
        onSelectionChange: noop,
    }

    state: DepartmentTreeState = {
        root: [],
        listTipStatus: ListTipStatus.Loading,
    }

    tree = null

    /**
     * 获取选中项
     */
    getSelections = async () => {
        if (this.tree) {
            const nodes = (this.tree as any).getSelections()
            return await formatNodes(nodes)
        }

        return []
    }

    treeRef = (tree) => {
        this.tree = tree
    }

    /**
     * 是否使用磁盘样式图标（涉密模式）
     */
    doclibIconDisabled: boolean = false;

    async componentDidMount() {
        try {
            const { extraRoots } = this.props

            const data = await getOrganizations(this.props.isRequestNormal)
            this.doclibIconDisabled = await getConfidentialConfig('doclib_icon_disabled')

            this.setState({
                root: [...extraRoots, ...data],
                listTipStatus: !extraRoots.length && !data.length ? ListTipStatus.OrgEmpty : ListTipStatus.None,
            })
        } catch (ex) {
            this.setState({
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    /**
     * 点击获取子节点
     */
    getNodeChildren = (node: any): Promise<ReadonlyArray<any>> | undefined => {
        if (!isLeaf(node, this.props.selectType)) {
            return this.getChildNode(node)
        }
    }

    /**
     * 获取子节点
     */
    private async getChildNode(node: any): Promise<ReadonlyArray<any>> {
        try {
            const [deps, { users, page }] = await Promise.all([
                getSubDepartments([node.id]),
                includes(this.props.selectType, NodeType.USER) ?
                    getSubUsers({
                        node,
                        page: DefaultPage,
                        limit: UsersLimit,
                        isShowDisabledUsers: this.props.isShowDisabledUsers as boolean,
                    })
                    : { users: [], page: DefaultPage },
            ])
            if (users.length > 0 && users.length < node.subUserCount) {
                return [...users, { parentNode: node, isLoadMore: true, currentPage: page + 1 }, ...deps];
            } else {
                return [...users, ...deps];
            }
        } catch (ex) {
            return []
        }
    }

    /**
     * 分页加载更多用户
     * @param node 当前节点
     */
    protected loadMoreUsers = async (node: any) => {
        const { currentPage, parentNode } = node

        try {
            const { users, page } = await getSubUsers({
                node: parentNode,
                page: currentPage,
                limit: UsersLimit,
                isShowDisabledUsers: this.props.isShowDisabledUsers as boolean,
            })

            return users
        } catch (ex) {
            return []
        }
    }

    /**
     * 取消所有选择节点
     */
    cancelSelections = () => this.tree ? (this.tree as any).cancelSelections() : null
}