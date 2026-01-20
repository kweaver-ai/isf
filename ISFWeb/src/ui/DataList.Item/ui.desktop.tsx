import React from 'react'
import { noop } from 'lodash'
import classnames from 'classnames'
import Expand from '../Expand/ui.desktop'
import { CheckBox } from '@/sweet-ui'
import styles from './styles.desktop'

const DataListItem: React.FunctionComponent<UI.DataListItem.Props> = function DataListItem({
    className,
    checkBoxClassName,
    selectable = true,
    checkbox = true,
    selected = false,
    expandable = false,
    expandContent = null,
    expanded = false,
    children,
    onClick = noop,
    onContextMenu = noop,
    onDoubleClick = noop,
    onToggleSelect = noop,
    ...otherProps
}) {
    return (
        <li
            className={classnames(styles['item'], { [styles['selected']]: selected })}
        >
            <div
                className={classnames(styles['wrapper'], { [styles['checkbox-padding']]: checkbox }, className)}
                onClick={onClick}
                onContextMenu={onContextMenu}
                onDoubleClick={onDoubleClick}
                {...otherProps}
            >
                {
                    checkbox ?
                        <div className={classnames(styles['checkbox'], checkBoxClassName)}>
                            {
                                selectable ?
                                    <CheckBox
                                        checked={selected}
                                        onClick={(event) => event.stopPropagation()}
                                        onChange={onToggleSelect}
                                    /> : null
                            }
                        </div>
                        : null
                }
                {children}
            </div>
            {
                expandable ?
                    <Expand
                        open={expanded}
                        onContextMenu={(e) => e.preventDefault()}
                    >
                        {expandContent}
                    </Expand> :
                    null
            }
        </li>
    )
}

DataListItem.defaultProps = {
    selectable: true,
    checkbox: true,
    expandable: false,
    expandContent: null,
    expanded: false,
}

export default DataListItem