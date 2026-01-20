import React from 'react';
import classnames from 'classnames';
import styles from './styles';

/**
 * 鼠标点击时调用的处理函数
 * @param event 鼠标事件
 */
export type ClickEventHandler = (event: React.MouseEvent<HTMLButtonElement>) => any;

export interface BaseButtonProps {
    /**
     * 如果为true，则按钮所有的交互都将被禁用
     */
    disabled?: boolean;

    /**
     * 内联样式
     */
    style?: React.CSSProperties;

    /**
     * className
     */
    className?: string;

    /**
     * 鼠标点击时触发
     */
    onClick?: ClickEventHandler;
}

const BaseButton: React.FunctionComponent<BaseButtonProps> = function Button({
    disabled,
    style,
    onClick,
    className,
    children,
    ...otherProps
}) {
    return (
        <button className={classnames(styles['base-button'], className)} {...{ disabled, style, onClick, ...otherProps }}>
            {children}
        </button>
    );
};

BaseButton.defaultProps = {};

export default BaseButton;
