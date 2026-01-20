import React from 'react'
import { CheckBox } from '@/sweet-ui'
import UIIcon from '../UIIcon/ui.desktop';
import { Icon, Centered, LinkChip } from '../ui.desktop';
import TreeBase, { NodeStatus, SelectStatus, SelectType, CascadeDirection } from './ui.base'
import classnames from 'classnames'
import loading from './assets/loading.gif';
import styles from './styles.desktop'
import __ from './locale';

export default class Tree extends TreeBase {

    static defaultProps = {
        data: [],
        indent: 10,
        selectType: SelectType.UNRESTRICTED,
        cascadeDirection: CascadeDirection.DOUBLE_SIDED,

        ExpandedIcon: <span className={styles['switch-icon']}><UIIcon code={'\uf04c'} size={18} /></span>,
        ExpandingIcon: <span className={styles['switch-icon']}><UIIcon code={'\uf04c'} size={18} /></span>,
        UnexpandedIcon: <span className={styles['switch-icon']}><UIIcon code={'\uf04e'} size={18} /></span>,
        LeafIcon: <span className={styles['switch-icon']}>{' '}</span>,

        getNodeChildren: () => null,
        renderNode: () => null,
    }

    /**
    * 渲染树
    * @param nodeId
    */
    private renderTreeNode(nodeId = '') {
        const dataGroup = this.treeData[nodeId]
        if (dataGroup && dataGroup.length) {
            const {
                checkbox,
                renderNode,
                indent,
                isLeaf,
                ExpandedIcon,
                ExpandingIcon,
                UnexpandedIcon,
                LeafIcon,
                disabled,
            } = this.props
            const { nodeStatus, selectStatus } = this.state
            const ids = nodeId.split('.')
            const depth = ids.length - 1

            return (
                <ul className={classnames(styles['tree'], { [styles['expand']]: nodeStatus[nodeId] === NodeStatus.EXPANDED })}>
                    {
                        dataGroup.map((data, index, group) => {
                            const showCheckBox = typeof checkbox === 'function' ? checkbox(data, index, group) : checkbox
                            const currentNodeId = `${nodeId}.${index}`
                            const status = nodeStatus[currentNodeId]
                            const nodeDisabledStatus = this.getDisabledStatus(data)
                            let SwitchIcon = null

                            if (isLeaf(data, index, group)) {
                                SwitchIcon = LeafIcon
                            } else {
                                switch (status) {
                                    case NodeStatus.EXPANDED:
                                        SwitchIcon = ExpandedIcon
                                        break
                                    case NodeStatus.EXPANDING:
                                        SwitchIcon = ExpandingIcon
                                        break
                                    default:
                                        SwitchIcon = UnexpandedIcon
                                        break
                                }
                            }

                            return (
                                data.isLoadMore ?
                                    <div
                                        key={'load-more'}
                                        title={__('加载更多')}
                                        className={styles['load-more']}
                                        style={{ paddingLeft: indent * depth }}
                                        onClick={() => this.handleLoadMoreUsers(currentNodeId, data)}>
                                        <span> {__('加载更多')}</span>
                                        <UIIcon
                                            className={styles['load-more-icon']}
                                            code={'\uf10d'}
                                            size={16}
                                        />

                                    </div>
                                    :
                                    (
                                        data.isLoading ?
                                            <Centered key={'loading-icon'}>
                                                <Icon
                                                    url={loading}
                                                    size={20}
                                                />
                                            </Centered>
                                            :
                                            <li
                                                className={styles['node']} key={currentNodeId}
                                                onDoubleClick={(e) => { this.handleDoubleClick(e, data, currentNodeId) }}
                                            >

                                                <div
                                                    className={classnames(
                                                        styles['wrapper'],
                                                        { [styles['selected']]: (selectStatus[currentNodeId] === SelectStatus.TRUE || selectStatus[currentNodeId] === SelectStatus.HALF) && !showCheckBox },
                                                    )}
                                                    style={{ paddingLeft: indent * depth }}
                                                    onClick={
                                                        !disabled && !nodeDisabledStatus ?
                                                            (e) => {
                                                                e.stopPropagation();
                                                                this.toggleSelect(currentNodeId, data);
                                                            } : undefined
                                                    }
                                                >
                                                    <div
                                                        className={styles['switch']}
                                                        onClick={
                                                            !disabled && !isLeaf(data, index, group)
                                                                ? (e) => {
                                                                    e.stopPropagation();
                                                                    this.toggleExpand(currentNodeId)
                                                                } : undefined}
                                                    >
                                                        {SwitchIcon}
                                                    </div>
                                                    {
                                                        showCheckBox ?
                                                            (
                                                                <span className={styles['check-box']}>
                                                                    <CheckBox
                                                                        disabled={nodeDisabledStatus}
                                                                        checked={selectStatus[currentNodeId] === SelectStatus.TRUE || selectStatus[currentNodeId] === SelectStatus.HALF}
                                                                        onChange={() => this.toggleSelect(currentNodeId, data)}
                                                                    />
                                                                </span>
                                                            )
                                                            : null
                                                    }
                                                    <LinkChip
                                                        className={classnames(styles['name'], { [styles['disabled']]: (disabled || nodeDisabledStatus) })}
                                                    >
                                                        {
                                                            renderNode(data, index, group)
                                                        }
                                                    </LinkChip>
                                                </div>
                                                {
                                                    nodeStatus[currentNodeId] === NodeStatus.EXPANDING ?
                                                        <Centered key={'loading-icon'}>
                                                            <Icon
                                                                url={loading}
                                                                size={20}
                                                            />
                                                        </Centered>
                                                        :
                                                        this.renderTreeNode(currentNodeId)
                                                }
                                            </li>
                                    )
                            )
                        })
                    }
                </ul>
            )
        }
        return null
    }

    render() {
        return this.renderTreeNode()
    }
}