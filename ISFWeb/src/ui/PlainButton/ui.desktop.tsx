import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash'
import UIIcon from '../UIIcon/ui.desktop';
import styles from './styles.desktop';

const PlainButton: React.FunctionComponent<UI.PlainButton.Props> = function PlainButton({
    type = 'button',
    disabled = false,
    icon,
    minWidth = 80,
    width,
    className,
    children,
    fallback,
    onClick = noop,
    size = 13,
    ...props
}) {
    return (
        <button
            className={
                classnames(
                    className,
                    styles['plainbutton'],
                    styles['box-sizing-border-box'],
                    {
                        [styles['disabled']]: disabled,
                    },
                )
            }
            type={type}
            disabled={disabled}
            style={{ minWidth, width }}
            onClick={onClick}
            {...props}
        >
            {
                icon ?
                    <span className={styles['icon']}>
                        <UIIcon
                            size={size}
                            code={icon}
                            fallback={fallback}
                        />
                    </span> :
                    null
            }
            {
                children
            }
        </button>
    )
}

export default PlainButton