import React from 'react';
import classnames from 'classnames';
import { noop } from 'lodash';
import { ClassName } from '@/ui/helper';
import SweetIcon from '../SweetIcon';
import View from '../View';
import styles from './styles';

/**
 * 按钮主题
 */
export type Theme = 'regular' | 'oem' | 'gray' | 'dark' | 'text' | 'noborder';

/**
 * 按钮类型
 */
export type Type = 'button' | 'submit' | 'reset';

/**
 * 按钮尺寸
 */
export type Size = 'normal' | 'auto';

/**
 * 鼠标点击时调用的处理函数
 * @param event 鼠标事件
 */
export type ClickEventHandler = (event: React.MouseEvent<HTMLButtonElement>) => void;

export interface ButtonProps {
    /**
     * 内联样式
     */
    style?: React.CSSProperties;

    /**
     * 样式class
     */
    className?: string;

    /**
     * 如果为true，则按钮所有的交互都将被禁用
     */
    disabled?: boolean;

    /**
     * 按钮主题
     */
    theme?: Theme;

    /**
     * 设置button原生的type值
     */
    type?: Type;

    /**
     * 设置按钮尺寸
     */
    size?: Size;

    /**
     * 设置按钮的图标类型, 类型为图标名称或图标元素
     */
    icon?: string | React.ReactNode;

    /**
     * 按钮内的图标大小
    */
    iconSize?: number;

    /**
     * 按钮宽度
     */
    width?: number | string;

    /**
     * 按钮高度
     */
    height?: number | string;

    /**
     * 点击按钮时的回调
     * @param event 鼠标事件
     */
    onClick?: ClickEventHandler;

    /**
     * 鼠标按下时的回调
     */
    onMouseDown?: ClickEventHandler;
}

const Button: React.FunctionComponent<ButtonProps> = function Button({
    disabled = false,
    type = 'button',
    theme = 'regular',
    icon,
    width,
    height,
    size,
    className,
    onClick = noop,
    onMouseDown = noop,
    children,
    iconSize = 16,
    style,
    ...otherProps
}) {
    return (
        <button
            type={type}
            style={{ ...style, width, height }}
            disabled={disabled}
            className={classnames(
                styles['button'],
                styles[theme],
                className,
                { [styles['disabled']]: disabled },
                { [ClassName.BackgroundColor]: theme === 'oem' },
                { [styles['normal']]: (size !== 'auto' && theme !== 'text' && theme !== 'noborder') },
            )}
            onClick={(event) => !disabled && onClick(event)}
            onMouseDown={(event) => !disabled && onMouseDown(event)}
            {...otherProps}
        >
            {icon ? typeof icon === 'string' ? (
                <SweetIcon
                    size={iconSize}
                    name={icon}
                    color={theme === 'dark' || theme === 'oem' ? '#fff' : '#505050'}
                    className={classnames(styles['icon'], { [styles['disabled']]: disabled })}
                />
            ) : (
                icon
            ) : null}
            {
                children !== 'undefined' ?
                    <View inline={true}>
                        {children}
                    </View>
                    : null
            }
        </button>
    );
};

export default Button;
