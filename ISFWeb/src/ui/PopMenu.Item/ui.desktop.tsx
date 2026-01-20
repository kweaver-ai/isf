import React from 'react'
import classnames from 'classnames'
import UIIcon from '../UIIcon/ui.desktop'
import styles from './styles.desktop'
import { ClassName } from '../helper';

const PopMenuItem: React.FunctionComponent<any> = function PopMenuItem({ labelClassName, icon, onDOMNodeMount, label, className, children, size = 13, ...otherProps }) {
    return (
        <li
            className={classnames(
                styles['item'],
                className,
                { [styles['padding']]: typeof icon === 'undefined' },
                ClassName.Color__Hover,
            )}
            {...otherProps}
            ref={(ref) => typeof onDOMNodeMount === 'function' ? onDOMNodeMount(ref) : null}
        >
            {
                typeof icon !== 'undefined' ?
                    typeof icon === 'string' ?
                        <UIIcon code={icon} size={16} className={styles['icon']} /> :
                        icon :
                    null
            }
            {
                label ?
                    <span
                        className={classnames(styles['label'], labelClassName)}
                        style={{ fontSize: size }}>
                        {label}
                    </span> : null
            }
            <span className={styles['children']}>{children}</span>
        </li>
    )
}

export default PopMenuItem