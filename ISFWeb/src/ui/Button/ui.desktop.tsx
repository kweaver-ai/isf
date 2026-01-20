import React from 'react';
import { noop } from 'lodash';
import classnames from 'classnames';
import { ClassName } from '../helper';
import UIIcon from '../UIIcon/ui.desktop';
import styles from './styles.desktop';

const Button: React.FunctionComponent<UI.Button.Props> = function Button({
    role,
    type = 'button',
    disabled = false,
    theme = 'regular',
    icon,
    minWidth,
    width,
    textColor,
    className,
    onClick = noop,
    onMouseDown = noop,
    children,
    fallback,
    size = 13,
    ...otherProps
}) {
    return (
        <button
            role={role}
            type={type}
            style={{ minWidth, width, color: textColor }}
            disabled={disabled}
            className={
                classnames(
                    styles['button'],
                    styles[theme],
                    className,
                    { [ClassName.BackgroundColor]: theme === 'oem' },
                    {
                        [styles['disabled']]: disabled,
                    },
                )
            }
            onClick={(event) => !disabled && onClick(event)}
            onMouseDown={onMouseDown}
            {...otherProps}
        >
            {
                icon ?
                    <span className={styles['icon']} >
                        <UIIcon
                            size={size}
                            code={icon}
                            fallback={fallback}
                            color={theme === 'dark' ? '#fff' : '#757575'}
                        />
                    </span > :
                    null
            }
            <span>
                {
                    children
                }
            </span>
        </button >
    )
}

export default Button