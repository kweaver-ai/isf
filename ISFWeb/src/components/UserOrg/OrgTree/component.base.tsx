import * as React from 'react'
import { includes, noop, omit } from 'lodash';
import session from '@/util/session/index';
import { sortDepartment, getSubDepartments } from '@/core/thrift/sharemgnt/sharemgnt'
import { SystemRoleType } from '@/core/role/role'
import { getOrganizations } from '@/core/department/department'
import { ErrorCode } from '@/core/thrift/sharemgnt/errcode'
import { Message2, Toast } from '@/sweet-ui'
import { UIIcon, Title } from '@/ui/ui.desktop'
import { DataNode, Key } from '@/sweet-ui/components/DragTree/helper'
import { ListTipStatus } from '../../ListTipComponent/helper'
import { getIcon } from '../../OrganizationTree/helper'
import WebComponent from '../../webcomponent'
import { SpecialDep } from '../helper'
import __ from './locale'
import styles from './styles.view';

interface OrgTreeProps extends React.Props<void> {
    /**
     * 用户id
     */
    userid: string;

    /**
     * 是否禁用空间站
     */
    isKjzDisabled: boolean;

    /**
     * 选中项改变
     */
    onRequestSelectDep: (selectedDep: DataNode) => any;
}

interface OrgTreeState {
    /**
     * 加载状态
     */
    listTipStatus: ListTipStatus;

    /**
     * 树节点
     */
    nodes: ReadonlyArray<DataNode>;

    /**
     * 选中的节点key
     */
    selectedKey: string;

    /**
     * 展开的节点key
     */
    expandedKeys: ReadonlyArray<string>;
}

export default class OrgTreeBase extends WebComponent<OrgTreeProps, OrgTreeState> {
    static defaultProps = {
        onRequestSelectDep: noop,
    }

    state = {
        listTipStatus: ListTipStatus.Loading,
        nodes: [],
        selectedKey: '',
        expandedKeys: [],
    }

    /**
     * 两个特殊的分组
     */
    specialGroup = []

    treeWrapper

    async componentDidMount() {
        await this.initOrgTree()
    }

    /**
     * 在组件销毁后设置state，防止内存泄漏
     */
    componentWillUnmount() {
        this.setState = (state, callback) => {
            return
        }
    }

    public initOrgTree = async (): Promise<void> => {
        try {
            let datas = await getOrganizations()

            const isAdminOrSupper = includes([SystemRoleType.Supper, SystemRoleType.Admin], session.get('isf.userInfo').user.roles[0].id)

            datas = [
                { id: '-1', name: __('未分配组'), depart_existed: false, user_existed: false },
                { id: '-2', name: __('所有用户'), depart_existed: false, user_existed: false },
                ...datas,
            ]

            this.specialGroup = datas[0] && datas[0].id === SpecialDep.Unassigned && isAdminOrSupper ? this.formatNode(datas.slice(0, 2)) : []

            const treeNodes = this.formatNode(!!this.specialGroup.length || !isAdminOrSupper ? datas.slice(2) : datas)

            const selectedKey = treeNodes.length ? treeNodes[0].key : !!this.specialGroup.length && isAdminOrSupper ? this.specialGroup[1].key : ''

            this.setState({
                expandedKeys: [],
                nodes: [],
                selectedKey: '',
                listTipStatus: ListTipStatus.None,
            }, () => {
                this.setState({
                    nodes: treeNodes,
                    selectedKey,
                })

                if (selectedKey) {
                    this.props.onRequestSelectDep(treeNodes.length ? treeNodes[0] : this.specialGroup[1])
                }
            })
        } catch (ex) {
            this.setState({
                listTipStatus: ListTipStatus.LoadFailed,
            })
        }
    }

    /**
     * 选中节点
     */
    protected seletectNode = (selectedKeys: ReadonlyArray<string>, { node }: { node: DataNode }): void => {
        if (selectedKeys.length) {
            this.setState({
                selectedKey: selectedKeys[0],
            })

            this.props.onRequestSelectDep(node)
        }
    }

    /**
     * 展开节点
     */
    protected expandNode = (expandedKeys: ReadonlyArray<string>, { node }: { node: DataNode }): void => {
        this.setState({
            expandedKeys,
        }, () => {
            if (node && node.pos && node.pos.split('-').length > 5) {
                setTimeout(() => {
                    this.treeWrapper.scrollLeft = (node.pos.split('-').length - 5) * 24
                }, 100)
            }
        })
    }

    /**
     * 拖拽
     * @param node 目标节点
     * @param dragNode 拖拽的节点
     * @param dropPosition 拖动的位置
     * @param dropToGap 是否拖动至缝隙
     */
    protected drop = async ({ node, dragNode, dropPosition, dropToGap }): Promise<void> => {
        if (node.pos.split('-').slice(0, -1).join('') !== dragNode.pos.split('-').slice(0, -1).join('')) {
            Toast.open(__('不允许跨部门/组织调整'))

            return
        }

        if (dropToGap) {
            try {
                const dropPos = node.pos.split('-')
                dropPosition = dropPosition - Number(dropPos[dropPos.length - 1])

                const data = [...this.state.nodes]

                // 找到dragNodeObject并在当前位置移除
                let dragObj

                this.loop(data, dragNode.key, (item, index, arr) => {
                    arr.splice(index, 1) // eslint-disable-line

                    dragObj = item
                })

                // 找到NodeObject及其index并获取当前层级的array
                let levelArr
                let nodeIndex
                let targetNode = node

                this.loop(data, node.key, (item, index, arr) => {
                    levelArr = arr
                    nodeIndex = index

                    if (node.dragOverGapBottom) {
                        targetNode = levelArr[index + 1] || null
                    }
                })

                await sortDepartment([this.props.userid, dragNode.key, targetNode ? targetNode.key : ''])

                // 放置在目标节点的上方
                if (dropPosition === -1) {
                    levelArr.splice(nodeIndex, 0, dragObj) // eslint-disable-line
                }
                // 放置在目标节点的下方
                else {
                    levelArr.splice(nodeIndex + 1, 0, dragObj) // eslint-disable-line
                }

                this.setState({
                    nodes: data,
                    selectedKey: dragNode.key,
                })

                this.props.onRequestSelectDep(dragNode)

            } catch (ex) {
                if (ex && ex.error && ex.error.errID) {
                    switch (ex.error.errID) {
                        case ErrorCode.DepNotExist:
                            if (await Message2.info({ message: __('部门“${srcDep}”已不存在。', { srcDep: dragNode.data.name }) })) {
                                this.deleteNode(dragNode.data)
                            }
                            break

                        case ErrorCode.SortTargetDepNotExist:
                            if (await Message2.info({ message: __('目标部门“${targetDep}”已不存在。', { targetDep: node.data.name }) })) {
                                this.deleteNode(node.data)
                                this.setState({
                                    selectedKey: node.data.id,
                                })
                            }
                            break

                        default:
                            Message2.info({ message: ex.errMsg || '' })
                            break
                    }
                }
            }
        }
    }

    /**
     * 新增节点
     */
    public addNode = async (nodeInfo: Core.ShareMgnt.ncTDepartmentInfo, parentInfo: Core.ShareMgnt.ncTDepartmentInfo | null = null): Promise<void> => {
        const { nodes, expandedKeys } = this.state

        let newNode = this.formatNode([nodeInfo])[0]

        let datas = [...nodes]

        let newExpandedKeys = [...expandedKeys]

        if (parentInfo && parentInfo.id) {
            await this.loop(datas, parentInfo.id, async (item, index, levelArr) => {
                newNode = { ...newNode, parent: omit(item, 'children') }

                const children = item.children ?
                    [...item.children, newNode]
                    : item.isLeaf ?
                        [newNode]
                        : await this.getSubNodes(item)

                levelArr.splice(index, 1, { ...item, isLeaf: false, children }) // eslint-disable-line

                if (!expandedKeys.some((key) => key === item.key)) {
                    newExpandedKeys = [...newExpandedKeys, item.key]
                }
            })
        } else {
            datas = [...datas, newNode]
        }

        this.setState({
            nodes: datas,
            selectedKey: nodeInfo.id,
            expandedKeys: newExpandedKeys,
        })

        this.props.onRequestSelectDep(newNode)
    }

    /**
     * 编辑节点
     */
    public updateNode = async (nodeInfo: Core.ShareMgnt.ncTDepartmentInfo, params: any, isUpdateChildren: boolean): Promise<void> => {
        let datas = [...this.state.nodes]

        let updatedNode

        await this.loop(datas, nodeInfo.id, (item, index, levelArr) => {
            updatedNode = { ...this.formatNode([{ ...nodeInfo, ...params }], item.parent || null)[0], children: item.children || undefined }

            levelArr.splice(index, 1, updatedNode) // eslint-disable-line

            if (isUpdateChildren && updatedNode.children) {
                updatedNode.children = this.updateNodesParams(updatedNode.children, params)
            }
        })

        this.setState({
            nodes: datas,
            selectedKey: nodeInfo.id,
        })

        this.props.onRequestSelectDep(updatedNode)
    }

    /**
     * 删除节点
     */
    public deleteNode = async (nodeInfo: Core.ShareMgnt.ncTDepartmentInfo): Promise<DataNode> => {
        const { expandedKeys, nodes } = this.state

        const datas = [...nodes]

        let nextDelectedNode

        let exKeys = expandedKeys.filter((key) => key !== nodeInfo.id)

        let deletedNode

        await this.loop(datas, nodeInfo.id, async (item, index, levelArr) => {
            deletedNode = item

            levelArr.splice(index, 1) // eslint-disable-line

            const arrLength = levelArr.length

            if (nodeInfo.is_root) {
                if (arrLength) {
                    nextDelectedNode = index <= arrLength - 1 ? levelArr[index] : levelArr[arrLength - 1]
                } else {
                    nextDelectedNode = this.specialGroup[1]
                }
            } else {
                nextDelectedNode = item.parent

                if (!arrLength) {
                    await this.loop(datas, nextDelectedNode.key, (parentNode, index, parentLevelArr) => {
                        parentNode.isLeaf = true

                        nextDelectedNode = { ...parentNode }

                        exKeys = exKeys.filter((key) => key !== nextDelectedNode.key)
                    })
                }
            }
        })

        if (nextDelectedNode) {
            this.setState({
                nodes: datas,
                selectedKey: nextDelectedNode.key,
                expandedKeys: exKeys,
            })

            this.props.onRequestSelectDep(nextDelectedNode)
        }

        return deletedNode
    }

    /**
     * 移动节点
     */
    public moveNode = async (nodeInfo: Core.ShareMgnt.ncTDepartmentInfo, targetNodeInfo: Core.ShareMgnt.ncTDepartmentInfo, params: any = null): Promise<void> => {
        let moveNode = await this.deleteNode(nodeInfo)

        if (params) {
            moveNode = this.updateNodesParams([moveNode], params)[0]
        }

        let datas = [...this.state.nodes]

        await this.loop(datas, targetNodeInfo.id, async (item, index, levelArr) => {
            item.isLeaf = false

            moveNode = { ...moveNode, parent: omit(item, 'children') }

            if (item.children) {
                item.children.push(moveNode) // eslint-disable-line
            } else {
                if (item.is_root) {
                    item.children = [moveNode]
                } else {
                    item.children = (await this.getSubNodes(item)).map((child) => child.key === moveNode.key ? moveNode : child)
                }
            }
        })

        this.setState({
            nodes: datas,
            expandedKeys: this.state.expandedKeys.filter((key) => key !== nodeInfo.id),
        })
    }

    /**
     * 异步加载数据
     */
    protected loadData = async (node: DataNode): Promise<void> => {
        if (node.children) {
            return Promise.resolve()
        } else {
            const subDeps = await this.getSubNodes(node)

            this.setState({
                nodes: this.updateTreeNodes(this.state.nodes, node.key, subDeps),
            })
        }
    }

    /**
     * 更新节点及子节点的属性
     */
    private updateNodesParams(nodes: ReadonlyArray<DataNode>, params: any): ReadonlyArray<DataNode> {
        let datas = [...nodes]

        datas.forEach((node) => {
            node.data = { ...node.data, ...params }

            if (node.children) {
                this.updateNodesParams(node.children, params);
            }
        })

        return datas
    }

    /**
     * 获取子节点
     */
    private getSubNodes = async (node: DataNode): Promise<ReadonlyArray<DataNode>> => {
        try {
            const subDeps = await getSubDepartments([node.data.id])

            return this.formatNode(subDeps, node)

        } catch (ex) {
            if (ex && ex.error && ex.error.errID) {
                switch (ex.error.errID) {
                    case ErrorCode.DepOrOrgNotExist:
                        if (await Message2.info({ message: __('部门“${srcDep}”已不存在。', { srcDep: node.data.name }) })) {
                            this.deleteNode(node.data)
                        }
                        break

                    default:
                        Message2.info({ message: ex.error.errMsg })
                        break
                }
            }

            return []
        }
    }

    /**
     * 格式化节点
     */
    private formatNode(datas: Core.ShareMgnt.ncTDepartmentInfo, parent: Core.ShareMgnt.ncTDepartmentInfo | null = null) {
        return datas.map((data) => ({
            key: data.id,
            title: (
                <Title
                    inline={true}
                    content={data.name}
                >
                    <span className={styles['title']}>{data.name}</span>
                </Title>
            ),
            isLeaf: data.id === SpecialDep.Unassigned || data.id === SpecialDep.AllUsers || (data.hasOwnProperty('depart_existed') && !data.depart_existed) || data.subDepartmentCount === 0,
            icon: (
                <UIIcon
                    className={styles['icon']}
                    size={16}
                    {...getIcon(data)}
                />
            ),
            data,
            parent: parent ? omit(parent, 'children') : null,
        }))
    }

    /**
     * 遍历节点
     */
    private loop = async (datas: ReadonlyArray<DataNode>, key: string, callback: Function): Promise<void> => {
        for (let i = 0; i < datas.length; i++) {
            if (datas[i].key === key) {
                return await callback(datas[i], i, datas)
            }

            if (datas[i].children) {
                await this.loop(datas[i].children, key, callback)
            }
        }
    }

    /**
     * 异步展开节点更新nodes
     */
    private updateTreeNodes = (list: ReadonlyArray<DataNode>, key: Key, children: ReadonlyArray<DataNode>): ReadonlyArray<DataNode> => {
        return list.map((node) => {
            if (node.key === key) {
                return {
                    ...node,
                    children,
                }
            } else if (node.children) {
                return {
                    ...node,
                    children: this.updateTreeNodes(node.children, key, children),
                }
            }

            return node
        })
    }
}