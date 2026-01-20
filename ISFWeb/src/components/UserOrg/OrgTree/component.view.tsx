import * as React from 'react'
import classnames from 'classnames'
import { DragTree } from '@/sweet-ui'
import { UIIcon, Title } from '@/ui/ui.desktop'
import { NodeType, getNodeIconByType } from '@/core/organization';
import { ListTipStatus } from '../../ListTipComponent/helper';
import ListTipComponent from '../../ListTipComponent/component.view';
import OrgTreeBase from './component.base'
import __ from './locale'
import styles from './styles.view';

export default class OrgTree extends OrgTreeBase {
    render() {
        const {
            listTipStatus,
            nodes,
            selectedKey,
            expandedKeys,
        } = this.state

        return (
            <div className={styles['tree']}>
                <div
                    ref={(tree) => (this.treeWrapper = tree)}
                    className={styles['tree-wrapper']}
                >
                    {
                        listTipStatus !== ListTipStatus.None ?
                            <ListTipComponent
                                listTipStatus={listTipStatus}
                            />
                            :
                            <div className={styles['drag-tree']}>
                                {
                                    this.specialGroup.map((group) => (
                                        <div
                                            className={styles['special-node']}
                                            key={group.key}
                                        >
                                            <span className={styles['switcher']}></span>
                                            <span
                                                className={classnames(styles['node-wrapper'], { [styles['selected']]: selectedKey === group.data.id })}
                                                onClick={() => this.seletectNode([group.key], { node: group })}
                                            >
                                                <UIIcon
                                                    role={'ui-uiicon'}
                                                    className={styles['icon']}
                                                    size={16}
                                                    {...getNodeIconByType(NodeType.DEPARTMENT)}
                                                />
                                                <Title
                                                    role={'ui-title'}
                                                    inline={true}
                                                    content={group.data.name}
                                                >
                                                    {group.data.name}
                                                </Title>
                                            </span>
                                        </div>
                                    ))
                                }
                                {
                                    nodes && nodes.length ?
                                        <DragTree
                                            role={'sweetui-dragtree'}
                                            draggable={this.props.isKjzDisabled}
                                            treeData={nodes}
                                            selectedKeys={[selectedKey]}
                                            expandedKeys={expandedKeys}
                                            onSelect={this.seletectNode}
                                            onExpand={this.expandNode}
                                            onDrop={this.drop}
                                            loadData={this.loadData}
                                        />
                                        : null
                                }
                            </div>
                    }
                </div>
                {
                    this.props.isKjzDisabled ?
                        <div className={styles['org-tree-remark']}>
                            {__('注：可直接拖动部门/组织调整顺序，不允许跨部门/组织调整')}
                        </div>
                        : null
                }
            </div>
        )
    }
}