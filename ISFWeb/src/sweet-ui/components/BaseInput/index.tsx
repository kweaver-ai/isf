import React from 'react';
import classnames from 'classnames';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import styles from './styles';

/**
 * 文本值变化事件对象
 */
export type ValueChangeEvent = SweetUIEvent<string>;

interface BaseInputProps {
    /**
     * 文本值
     */
    value?: string | number;

    /**
     * className
     */
    className?: string;

    /**
     * 文本框类型
     */
    type?: 'text' | 'password' | 'email' | 'number' | 'search' | 'tel' | 'url' | 'time' | 'month';

    /**
     * css样式
     */
    style?: React.CSSProperties;

    /**
     * 禁用
     */
    disabled?: boolean;

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * 只读
     */
    readOnly?: boolean;

    /**
     * 输入值发生变化时触发，传递value值
     */
    onValueChange?: (event: ValueChangeEvent) => void;

    /**
     * 渲染完成后触发
     */
    onMounted?: (ref: HTMLInputElement) => void;

    /**
     * 点击输入框时触发
     */
    onClick?: (event: React.MouseEvent<HTMLInputElement>) => void;

    /**
     * 键盘输入时触发
     */
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 输入框聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 输入框失焦时触发
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 粘贴事件
     */
    onPaste?: () => void;
}

const BaseInput: React.SFC<BaseInputProps> = function BaseInput({
    value = '',
    className,
    type,
    style,
    disabled = false,
    placeholder,
    readOnly,
    onMounted,
    onValueChange,
    onClick,
    onFocus,
    onBlur,
    onPaste,
    onKeyDown,
}) {
    const dispatchValueChangeEvent = createEventDispatcher(onValueChange);

    return (
        <input
            type={type}
            ref={onMounted}
            value={value}
            style={style}
            className={classnames(styles['base-input'], { [styles['disabled']]: disabled }, className)}
            onChange={(event) => dispatchValueChangeEvent((event.target as HTMLInputElement).value)}
            {...{ onClick, onFocus, onBlur, onPaste, disabled, placeholder, readOnly, onKeyDown }}
        />
    );
};

BaseInput.defaultProps = {
    type: 'text',
};

export default BaseInput;
