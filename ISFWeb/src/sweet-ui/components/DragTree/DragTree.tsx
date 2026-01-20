import React from 'react'
import RcTree from 'rc-tree'
import SweetIcon from '../SweetIcon'
import Icon from '../Icon'
import { TreeNodeProps, DragTreeProps } from './helper'
import styles from './styles.view.css'
import loading from './assets/loading.gif';

export default function DragTree(props: DragTreeProps) {
    const {
        selectable = true,
        draggable = true,
        role,
        ...restProps
    } = props

    const switcherIcon = (treeNode: TreeNodeProps) => {
        if (treeNode.isLeaf) {
            return ''
        }

        if (treeNode.loading) {
            return <Icon src={loading} />
        } else {
            return (
                <SweetIcon
                    name={treeNode.expanded ? 'arrowDown' : 'arrowRight'}
                />
            )
        }
    }

    return (
        <div className={styles['tree']} role={role}>
            <RcTree
                prefixCls={'drag-tree'}
                selectable={selectable}
                draggable={draggable}
                switcherIcon={switcherIcon}
                {...restProps}
            />
        </div>
    )
}