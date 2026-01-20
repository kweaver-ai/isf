import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash';
import UIIcon from '../UIIcon/ui.desktop';
import Title from '../Title/ui.desktop'
import styles from './styles.desktop';

const InlineButton: React.FunctionComponent<UI.InlineButton.Props> = function InlineButton({ role, size, title, code, fallback, iconSize, disabled, onClick, className, type = 'button', color = '#757575', element, ...otherProps }) {
    return (
        title ?
            <Title
                role={role}
                timeout={0}
                content={title}
                element={element}
            >
                <button
                    className={classnames(styles['inline-button'], { [styles['disabled']]: disabled }, className)}
                    style={{ width: size, height: size, lineHeight: `${size - 2}px` }}
                    type={type}
                    onClick={disabled ? noop : onClick}
                    {...otherProps}
                >
                    <UIIcon
                        className={styles['icon']}
                        code={code}
                        fallback={fallback}
                        color={color}
                        size={iconSize}
                    />
                </button>
            </Title>
            :
            <button
                className={classnames(styles['inline-button'], { [styles['disabled']]: disabled }, className)}
                style={{ width: size, height: size, lineHeight: `${size - 2}px` }}
                disabled={disabled}
                type={type}
                onClick={onClick}
                {...otherProps}
            >
                <UIIcon
                    className={styles['icon']}
                    code={code}
                    fallback={fallback}
                    color={color}
                    size={iconSize}
                />
            </button>
    )
}

InlineButton.defaultProps = {
    size: 24,
    title: '',
    iconSize: 16,
    disabled: false,
    onClick: noop,
}

export default InlineButton;