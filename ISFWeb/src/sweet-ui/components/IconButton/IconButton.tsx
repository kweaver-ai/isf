import React from 'react';
import classnames from 'classnames';
import BaseButton, { BaseButtonProps } from '../BaseButton';
import styles from './styles';

/**
 * 传递到BaseButton上的props
 */
type BaseButtonReceivedProps = 'disabled' | 'style' | 'className' | 'onClick';

interface IconButtonProps extends Pick<BaseButtonProps, BaseButtonReceivedProps> {
    /**
     * 图标
     */
    icon: React.ReactNode;

    /**
     * 按钮尺寸
     */
    size?: number;

    /**
     * 样式class
     */
    className?: string;
}

const IconButton: React.FunctionComponent<IconButtonProps> = function IconButton({
    role,
    icon,
    size = 24,
    disabled,
    className,
    ...restProps
}) {
    return (
        <BaseButton
            role={role}
            className={classnames(styles['icon-button'], { [styles['disabled']]: disabled }, className)}
            style={{ width: size, height: size }}
            {...{ ...restProps, disabled }}
        >
            {icon}
        </BaseButton>
    );
};

export default IconButton;
