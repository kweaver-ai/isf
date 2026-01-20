import React from 'react';
import classnames from 'classnames'
import styles from './styles.desktop';

/**
 * @param selected 是否选中
 * @param onMouseOver 鼠标移动要上面去
 */
export default function AutoCompleteListItem({ role, children, onMouseOver, selected, onMount }) {
    return (
        <li
            role={role}
            className={classnames(styles['autocomplte-list-item'], { [styles['selected']]: selected })}
            onMouseOver={onMouseOver}
            ref={onMount}
        >
            {
                children
            }
        </li>
    )
}