import * as React from 'react';
import Tree from '@/ui/Tree/ui.desktop';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { Icon, Centered } from '@/ui/ui.desktop';
import { NodeType, isLeaf } from '@/core/organization';
import ListTipComponent from '../ListTipComponent/component.view';
import { ListTipStatus, ListTipMessage } from '../ListTipComponent/helper';
import OrganizationTreeBase from './component.base';
import * as loading from './assets/loading.gif';
import { getNodeName, getIcon } from './helper';
import styles from './styles.view';
import __ from './locale';

export default class OrganizationTree extends OrganizationTreeBase {
    /**
     * 格式化节点数据
     * @param node 节点数据
     */
    private formatter = (node: any): JSX.Element => {
        return (
            <span className={styles['node']} title={getNodeName(node)}>
                <span className={styles['icon']}>
                    <UIIcon {...getIcon(node)} size={16} />
                </span>
                <span className={styles['name']}>
                    {
                        getNodeName(node)
                    }
                </span>
            </span>
        );
    }

    /**
     * 根据节点数组递归生成JSX树节点
     * @param nodes 节点数组
     */
    private generateNodes = (nodes: ReadonlyArray<any> = []): ReadonlyArray<JSX.Element> => {
        return nodes.map((node) => {
            const nodeStatus = node.parent && node.parent.nodeStatus && node.parent.nodeStatus.disabled && this.props.isDisableChildrenByParent ?
                { disabled: true }
                : this.props.getNodeStatus && this.props.getNodeStatus(node)

            return (
                node.isLoading ?
                    <Centered
                        role={'ui-centered'}
                        key={'loading-icon'}
                    >
                        <Icon
                            role={'ui-icon'}
                            url={loading}
                            size={20}
                        />
                    </Centered>
                    : (
                        node.isLoadMore ?
                            <div
                                key={'load-more'}
                                title={__('加载更多')}
                                className={styles['load-more']}
                                onClick={() => this.loadMoreUsers(node)}>
                                <span> {__('加载更多')}</span>
                                <UIIcon
                                    role={'ui-uiicon'}
                                    className={styles['load-more-icon']}
                                    code={'\uf10d'}
                                    size={16}
                                />
                            </div>
                            :
                            <Tree.Node
                                role={'ui-tree.node'}
                                isLeaf={isLeaf(node, this.props.selectType)}
                                data={node}
                                key={node.id}
                                disabled={this.props.disabled}
                                formatter={this.formatter}
                                loader={this.loadSubs}
                                getStatus={() => (nodeStatus)}
                            >
                                {
                                    node.children ? this.generateNodes(node.children) : null
                                }
                            </Tree.Node>
                    )

            )
        })
    }

    /**
     * 获取提示
     */
    private getListTipMessage = (selectType?: ReadonlyArray<NodeType>): { [status: number]: string } => {
        return {
            ...ListTipMessage,
            [ListTipStatus.OrgEmpty]: selectType && selectType.includes(NodeType.USER) ?
                __('暂无可选的用户或部门') : __('暂无可选的部门'),
        }
    }

    render() {
        const { nodes, listTipStatus } = this.state
        const { selectType } = this.props
        if (listTipStatus === ListTipStatus.None) {
            return (
                <div className={styles['tree-wrp']}>
                    <Tree
                        role={'ui-tree'}
                        disabled={this.props.disabled}
                        onSelectionChange={this.fireSelectionChangeEvent.bind(this)}
                    >
                        {
                            this.generateNodes(nodes)
                        }
                    </Tree>
                </div>
            )
        }

        return (
            <div className={styles['tree-wrp']}>
                <div className={styles['list-tip']}>
                    <ListTipComponent
                        listTipStatus={listTipStatus}
                        listTipMessage={this.getListTipMessage(selectType)}
                        isInDialog={true}
                    />
                </div>
            </div>
        )

    }
}