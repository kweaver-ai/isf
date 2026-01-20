import React from 'react'
import { isFunction } from 'lodash';

/**
 * @todo
 * 1. 补充类型
 * 2. 添加 onDataLoad，异步加载数据时触发
 */

/**
 * 节点展开状态
 */
export enum NodeStatus {
    UNEXPANDED,
    EXPANDING,
    EXPANDED,
}

/**
 * 选择类型
 */
export enum SelectStatus {
    TRUE = 1,
    HALF = 0.5,
    FALSE = 0,
}

/**
 * 选择类型
 */
export enum SelectType {
    /**
     * 禁止选中
     */
    NONE = 0,

    /**
     * 同级单选
     */
    SINGLE = 1,

    /**
     * 同级多选
     */
    MULTIPLE = 2,

    /**
     * 级联单选
     */
    CASCADE_SINGLE = 3,

    /**
     * 级联多选
     */
    CASCADE_MULTIPLE = 4,

    /**
     * 无限制
     */
    UNRESTRICTED = 5,
}

export enum CascadeDirection {
    /**
     * 仅支持向下级联，不支持向上级联
     */
    DOWN,

    /**
     * 双向级联
     */
    DOUBLE_SIDED,
}

// 每次加载的用户数量
const UsersLimit = 150;

export default class Tree extends React.Component<any, any> {
    static defaultProps = {
        disabled: false,

        /**
         * 是否是组织树
         */
        isOrgTree: true,
    }

    state = {
        nodeStatus: {},
        selectStatus: {},
        loaded: false,
    }

    treeData: { [id: string]: Array<any> } = {}
    loadStatus = {}
    selectedIds: Array<string> = []

    constructor(props, context) {
        super(props, context)
        this.getSelections = this.getSelections.bind(this)
    }

    async componentDidMount() {
        const roots = await this.props.data
        this.treeData = { '': roots }
        this.loadStatus = { '': true }
        this.toggleExpand('')
        if (typeof this.props.getDefaultSelectStatus === 'function') {
            const selectStatus = {}
            roots.forEach((data, i) => {
                const id = `.${i}`
                selectStatus[id] = this.props.getDefaultSelectStatus(data, id, roots)
            })
            this.setState({
                selectStatus,
            })
        }

        if (typeof this.props.isDefaultExpanded === 'function') {
            for (const [data, index] of Array.from(roots, (data, index) => [data, index])) {
                if (!this.props.isLeaf(data, index, roots) && this.props.isDefaultExpanded(data, `.${index}`, roots)) {
                    await this.toggleExpand(`.${index}`)
                }
            }
        }
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.props.getNodeChildren !== prevProps.getNodeChildren) {
            /**
             * 树数据改变，重置树数据、选择状态
             */
            this.treeData = { '': this.props.data }
            this.loadStatus = { '': true }
            this.selectedIds = []
            this.setState({
                nodeStatus: {},
                selectStatus: {},
            })
            this.toggleExpand('')
        } else if (this.props.selectType !== prevProps.selectType) {
            /**
             * 选择类型改变，重置树选择状态
             */
            this.selectedIds = []
            this.setState({
                selectStatus: {},
            })
        } else if (this.props.data !== prevProps.data) {
            this.treeData = { '': this.props.data }
            this.loadStatus = { '': true }
            this.setState({
                nodeStatus: {},
            }, () => {
                this.toggleExpand('')
            })
        }
    }

    /**
     * 获取节点状态
     */
    protected getDisabledStatus(node: any): boolean {
        return isFunction(this.props.onRequestGetStatus) ? this.props.onRequestGetStatus(node).disabled : false
    }

    /**
     * 双击事件回调
     * @param e 事件对象
     */
    handleDoubleClick = (e: Event, data: object, currentNodeId: string): void => {
        const {
            disabled,
            isLeaf,
        } = this.props

        if (!disabled && !isLeaf(data)) {
            this.toggleExpand(currentNodeId);
            e.stopPropagation();
        } else {
            e.stopPropagation()
        }
    }

    /**
     * 展开节点
     * @param id
     */
    public async toggleExpand(id: string) {
        const { getNodeChildren, selectType } = this.props
        const { nodeStatus, selectStatus } = this.state
        switch (nodeStatus[id]) {
            case NodeStatus.EXPANDED:
                this.setState({
                    nodeStatus: { ...nodeStatus, [id]: NodeStatus.UNEXPANDED },
                })
                return
            case NodeStatus.EXPANDING:
                return
            default:
                if (!this.loadStatus[id]) {
                    this.setState({
                        nodeStatus: { ...nodeStatus, [id]: NodeStatus.EXPANDING },
                    })
                    const lastDotIndex = id.lastIndexOf('.')
                    this.treeData[id] = await getNodeChildren(this.treeData[id.slice(0, lastDotIndex)][id.slice(lastDotIndex + 1)])

                    this.loadStatus[id] = true
                    this.setState(({ nodeStatus, selectStatus: { ...nextSelectStatus } }) => {

                        if (typeof this.props.getDefaultSelectStatus === 'function') {
                            nextSelectStatus = { ...selectStatus }
                            this.treeData[id].forEach((data, index) => {
                                nextSelectStatus[`${id}.${index}`] = this.props.getDefaultSelectStatus(data, `${id}.${index}`, this.treeData[id])
                            })
                        } else if (selectType === SelectType.CASCADE_MULTIPLE) {
                            nextSelectStatus = { ...selectStatus }
                            this.treeData[id].forEach((data, index) => {
                                nextSelectStatus[`${id}.${index}`] = selectStatus[id] || SelectStatus.FALSE
                            })
                        }
                        if (typeof this.props.isDefaultExpanded === 'function') {
                            this.treeData[id].forEach((data, index) => {
                                if (!this.props.isLeaf(data, index, this.treeData[id]) && this.props.isDefaultExpanded(data, `${id}.${index}`, this.treeData[id])) {
                                    this.toggleExpand(`${id}.${index}`)
                                }
                            })
                        }

                        return {
                            nodeStatus: { ...nodeStatus, [id]: NodeStatus.EXPANDED },
                            selectStatus: nextSelectStatus,
                        }
                    })
                } else {
                    this.setState(({ nodeStatus }) => {
                        return { nodeStatus: { ...nodeStatus, [id]: NodeStatus.EXPANDED } }
                    })
                }
                return
        }
    }

    /**
     * 级联选择子节点
     * @param nodeId
     * @param status
     * @param selectStatus
     */
    private cascadeSelectChildren(nodeId, status, selectStatus) {
        selectStatus[nodeId] = status
        if (this.treeData[nodeId]) {
            this.treeData[nodeId].forEach((data, i) => {
                this.cascadeSelectChildren(`${nodeId}.${i}`, status, selectStatus)
            })
        }
    }

    /**
     * 分页加载更多用户
     * @param node 当前节点
     */
    protected async handleLoadMoreUsers(id, node) {
        const { selectType, isOrgTree = true } = this.props;
        const { selectStatus } = this.state;
        let nextSelectStatus = { ...selectStatus };
        const { parentNode, currentPage } = node;
        const lastDotIndex = id.lastIndexOf('.');
        const nodeId = id.slice(0, lastDotIndex)

        if (isOrgTree) {
            const users = this.treeData[nodeId].filter((item) => item.hasOwnProperty('user'));
            const departments = this.treeData[nodeId].filter((item) => !item.hasOwnProperty('user') && !item.hasOwnProperty('isLoadMore'));

            this.treeData[nodeId] = [...users, { isLoading: true }, ...departments];
            this.forceUpdate();
            const result = await this.props.loadMoreUsers(node);

            const isLoadMore = (result.length + (currentPage - 1) * UsersLimit) < parentNode.subUserCount

            if (isLoadMore) {
                this.treeData[nodeId] = [...users, ...result, { parentNode, isLoadMore: true, currentPage: node.currentPage + 1 }, ...departments];
                for (let i = 0; i < this.treeData[nodeId].length; i++) {
                    if (i < users.length) {
                        nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i}`] || SelectStatus.FALSE
                    } else if (i < (users.length + result.length)) {
                        if (selectType === SelectType.CASCADE_MULTIPLE) {
                            const status = selectStatus[nodeId] === SelectStatus.TRUE ? SelectStatus.TRUE : SelectStatus.FALSE;
                            nextSelectStatus[`${nodeId}.${i}`] = status;
                        } else {
                            nextSelectStatus[`${nodeId}.${i}`] = SelectStatus.FALSE;
                        }
                    } else if (i === (users.length + result.length)) {
                        nextSelectStatus[`${nodeId}.${i}`] = SelectStatus.FALSE;
                    } else {
                        nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i - result.length}`]
                    }
                }
            } else {
                this.treeData[nodeId] = [...users, ...result, ...departments];
                for (let i = 0; i < this.treeData[nodeId].length; i++) {
                    if (i < users.length) {
                        nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i}`] || SelectStatus.FALSE
                    } else if (i < (users.length + result.length)) {
                        if (selectType === SelectType.CASCADE_MULTIPLE) {
                            const status = selectStatus[nodeId] === SelectStatus.TRUE ? SelectStatus.TRUE : SelectStatus.FALSE;
                            nextSelectStatus[`${nodeId}.${i}`] = status;
                        } else {
                            nextSelectStatus[`${nodeId}.${i}`] = SelectStatus.FALSE;
                        }
                    } else {
                        nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i - result.length + 1}`]
                    }
                }
            }
        } else {
            const children = this.treeData[nodeId].slice(0, -1)

            this.treeData[nodeId] = [...children, { isLoading: true }]

            this.forceUpdate()

            const result = await this.props.loadMoreUsers(node)

            this.treeData[nodeId] = [...children, ...result]

            const isLoadMore = ((result.length + (currentPage - 1) * UsersLimit) < parentNode.subNodesCount) || false

            if (isLoadMore) {
                this.treeData[nodeId] = [...this.treeData[nodeId], { parentNode, isLoadMore: true, currentPage: node.currentPage + 1 }]
            }

            for (let i = 0; i < this.treeData[nodeId].length; i++) {
                if (i < children.length) {
                    nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i}`] || SelectStatus.FALSE
                } else if (i < (children.length + result.length)) {
                    if (selectType === SelectType.CASCADE_MULTIPLE) {
                        const status = selectStatus[nodeId] === SelectStatus.TRUE ? SelectStatus.TRUE : SelectStatus.FALSE;
                        nextSelectStatus[`${nodeId}.${i}`] = status;
                    } else {
                        nextSelectStatus[`${nodeId}.${i}`] = SelectStatus.FALSE;
                    }
                } else {
                    nextSelectStatus[`${nodeId}.${i}`] = selectStatus[`${nodeId}.${i - result.length + (isLoadMore ? 0 : 1)}`]
                }
            }
        }

        this.setState({
            selectStatus: nextSelectStatus,
        }, () => {
            if (typeof this.props.onSelectionChange === 'function') {
                this.props.onSelectionChange(this.getSelections())
            }
        })

    }

    /**
     * 级联选择父节点
     * @param nodeId
     * @param status
     * @param selectStatus
     */
    private cascadeSelectParent(nodeId, status, selectStatus) {
        // 仅支持向下级联
        if (this.props.cascadeDirection === CascadeDirection.DOWN) {
            selectStatus[nodeId] = status
            if (status === SelectStatus.FALSE) { // 取消勾选子节点
                nodeId.split('.').forEach((id, i, arr) => {
                    const upperId = arr.slice(0, arr.length - 1 - i).join('.')
                    selectStatus[upperId] = status
                })
            }
        } else {
            nodeId.split('.').forEach((id, i, arr) => {
                const upperId = arr.slice(0, arr.length - i).join('.')
                if (this.treeData[upperId] && this.treeData[upperId].find((data, i) => selectStatus[`${upperId}.${i}`] !== status)) {
                    selectStatus[upperId] = SelectStatus.HALF
                } else {
                    selectStatus[upperId] = status
                }
            })
        }
    }

    /**
     * 选择节点
     * @param id
     */
    public toggleSelect(id: string, select?: any) {
        const { selectType } = this.props
        const { selectStatus } = this.state

        let nextSelectStatus = {}

        switch (selectType) {

            case SelectType.SINGLE:
                nextSelectStatus = { [id]: SelectStatus.TRUE }
                this.selectedIds = nextSelectStatus[id] === SelectStatus.TRUE ? [id] : []
                break

            case SelectType.MULTIPLE: {
                this.selectedIds = []
                const parentNodeId = id.slice(0, id.lastIndexOf('.'))
                this.treeData[parentNodeId].forEach((data, i) => {
                    const currentNodeId = `${parentNodeId}.${i}`
                    if (currentNodeId === id) {
                        nextSelectStatus[currentNodeId] = selectStatus[id] === SelectStatus.TRUE ? SelectStatus.FALSE : SelectStatus.TRUE
                    } else {
                        nextSelectStatus[currentNodeId] = selectStatus[currentNodeId]
                    }

                    if (nextSelectStatus[currentNodeId] === SelectStatus.TRUE) {
                        this.selectedIds = [...this.selectedIds, currentNodeId]
                    }
                })
                break
            }

            case SelectType.CASCADE_SINGLE: {
                this.selectedIds = []
                const hasSelectedChild = Object.keys(selectStatus).find((key) => selectStatus[key] === SelectStatus.TRUE && key.startsWith(id) && key !== id)
                let nextStatus = !hasSelectedChild && selectStatus[id] === SelectStatus.TRUE ? SelectStatus.FALSE : SelectStatus.TRUE
                id.split('.').forEach((id, index, arr) => {
                    const currentNodeId = arr.slice(0, index + 1).join('.')
                    nextSelectStatus[currentNodeId] = nextStatus
                    if (index > 0 && nextStatus === SelectStatus.TRUE) {
                        this.selectedIds = [...this.selectedIds, currentNodeId]
                    }
                })
                break
            }

            case SelectType.CASCADE_MULTIPLE: {
                const status = selectStatus[id] === SelectStatus.TRUE ? SelectStatus.FALSE : SelectStatus.TRUE
                nextSelectStatus = { ...selectStatus }
                this.cascadeSelectChildren(id, status, nextSelectStatus)
                this.cascadeSelectParent(id, status, nextSelectStatus)
                break
            }

            case SelectType.UNRESTRICTED:
                nextSelectStatus = {
                    ...selectStatus,
                    [id]: selectStatus[id] === SelectStatus.TRUE ? SelectStatus.FALSE : SelectStatus.TRUE,
                }
                break

            default: break
        }
        this.setState({ selectStatus: nextSelectStatus }, () => {
            if (typeof this.props.onSelectionChange === 'function') {
                this.props.onSelectionChange(this.getSelections(), select, id)
            }
        })
    }

    /**
     * 获取级联多选选择项
     * @param nodeId
     */
    private getCascadeSelections(nodeId = '') {
        if (nodeId && this.state.selectStatus[nodeId] === SelectStatus.TRUE) {
            this.selectedIds = [...this.selectedIds, nodeId]
        } else if (this.treeData[nodeId] && this.treeData[nodeId].length) {
            this.treeData[nodeId].forEach((data, index) => {
                if (!data.hasOwnProperty('isLoadMore')) {
                    this.getCascadeSelections(`${nodeId}.${index}`)
                }
            })
        }
    }

    /**
     * 获取无限制选择项
     * @param nodeId
     */
    private getUnRestrictSelections(nodeId = '') {
        if (nodeId && this.state.selectStatus[nodeId] === SelectStatus.TRUE) {
            this.selectedIds = [...this.selectedIds, nodeId]
        }
        if (this.treeData[nodeId] && this.treeData[nodeId].length) {
            this.treeData[nodeId].forEach((data, index) => this.getUnRestrictSelections(`${nodeId}.${index}`))
        }
    }

    /**
     * 获取选中节点
     */
    public getSelections() {
        switch (this.props.selectType) {
            case SelectType.CASCADE_MULTIPLE:
                this.selectedIds = []
                this.getCascadeSelections()
                break
            case SelectType.UNRESTRICTED:
                this.selectedIds = []
                this.getUnRestrictSelections()
            default:
                break
        }
        return this.selectedIds.map((selectedId) => {
            const lastDotIndex = selectedId.lastIndexOf('.')
            return this.treeData[selectedId.slice(0, lastDotIndex)][selectedId.slice(lastDotIndex + 1)]
        })
    }

    /**
     * 取消所有选择
     */
    public cancelSelections() {
        this.selectedIds = []
        this.setState({
            selectStatus: {},
        })
    }
}