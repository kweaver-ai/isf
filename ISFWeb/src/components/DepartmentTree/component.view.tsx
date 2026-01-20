import * as React from 'react'
import DepartmentTreeBase from './component.base'
import Tree from '@/ui/Tree2/ui.desktop'
import UIIcon from '@/ui/UIIcon/ui.desktop'
import Icon from '@/ui/Icon/ui.desktop'
import { isLeaf, getNodeIcon, NodeType } from '@/core/organization'
import { userDocLibImg, classifiedUserDocLibImg } from '@/core/doclibs/assets';
import ListTipComponent from '../ListTipComponent/component.view'
import { ListTipStatus, ListTipMessage } from '../ListTipComponent/helper'
import __ from './locale'
import styles from './styles.desktop'

export default class DepartmentTree extends DepartmentTreeBase {

    renderNode = (node: any): JSX.Element => {
        return (
            <div
                role={'ui-title'}
                title={node.name || node.user && node.user.displayName || ''}
                key={node.id || ''}
            >
                {
                    node.isDocLib ?
                        <Icon
                            url={this.doclibIconDisabled ? classifiedUserDocLibImg : userDocLibImg}
                            size={20}
                        />
                        : <UIIcon
                            role={'ui-uiicon'}
                            {...getNodeIcon(node)}
                            size={16}
                        />
                }
                <span className={styles['name']}>{node.name || (node.user && node.user.displayName)}</span>
            </div>
        )
    }

    /**
     * 获取提示
     */
    getListTipMessage = (selectType?: ReadonlyArray<NodeType>): { [str: number]: string } => {
        return {
            ...ListTipMessage,
            [ListTipStatus.OrgEmpty]: selectType && selectType.includes(NodeType.USER) ?
                __('暂无可选的用户或部门') : __('暂无可选的部门'),
        }
    }

    render() {
        const {
            listTipStatus,
            root,
        } = this.state

        const {
            disabled,
            nodeSelectType,
            cascadeDirection,
            selectType,
            onSelectionChange,
        } = this.props

        if (listTipStatus === ListTipStatus.None) {
            return (
                <div className={styles['tree-wrp']}>
                    <Tree
                        role={'ui-tree2'}
                        disabled={disabled}
                        selectType={nodeSelectType}
                        cascadeDirection={cascadeDirection}
                        checkbox={true}
                        data={root}
                        isLeaf={(node) => isLeaf(node, selectType)}
                        renderNode={this.renderNode}
                        getNodeChildren={this.getNodeChildren}
                        loadMoreUsers={this.loadMoreUsers}
                        ref={this.treeRef}
                        onSelectionChange={onSelectionChange}
                    />
                </div >
            )
        }

        return (
            <div className={styles['list']}>
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