import React from 'react';
import classnames from 'classnames';
import { isFunction, noop } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import BaseInput, { ValueChangeEvent } from '../BaseInput';
import styles from './styles';

interface TextInputProps extends React.ClassAttributes<void> {
    /**
     * 文本框类型
     */
    type?: 'text' | 'password' | 'email' | 'number' | 'search' | 'tel' | 'url' | 'time' | 'month';

    /**
     * 输入区域宽度，不传则宽度由输入的内容决定
     */
    width?: string;

    /**
     * 输入区域最小宽度，不传width时有效
     */
    minWidth?: number;

    /**
     * 输入区域最大宽度，不传width时有效
     */
    maxWidth?: number;

    /**
     * 默认输入内容
     */
    defaultValue?: string;

    /**
     * 输入值
     */
    value: string;

    /**
     * className
     */
    className?: string;

    /**
     * 自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 聚焦时选中
     */
    selectOnFocus?: [number] | [number, number] | boolean;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 是否只读
     */
    readOnly?: boolean;

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * 文本框输入字符最大长度
     */
    maxLength?: number;

    /**
     *  是否必填
     */
    required?: boolean;

    /**
     * 输入值发生变化时触发，传递value值
     */
    onValueChange?: (event: ValueChangeEvent) => void;

    /**
     * 点击事件
     */
    onClick?: (event: React.MouseEvent<HTMLInputElement>) => void;

    /**
     * 输入限制函数
     */
    validator?: (value: string) => boolean;

    /**
     * 聚焦事件
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 失去焦点事件
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 键盘输入时触发
     */
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 按下回车键触发
     */
    onPressEnter?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 渲染完成后触发
     */
    onMounted?: (ref: HTMLElement) => void;

    /**
     * 粘贴事件
     */
    onPaste?: (event: SweetUIEvent<ClipboardEvent>) => void;

    /**
     * css样式
     */
    style?: React.CSSProperties;
}

interface TextInputState {
    value: string;
    width: number | string;
}

export default class TextInput extends React.Component<TextInputProps, TextInputState> {
    static defaultProps = {
        minWidth: 80,
        maxWidth: 500,
        selectOnFocus: false,
        disabled: false,
        onClick: noop,
        type: 'text',
        validator: () => true,
    };

    constructor(props: TextInputProps, ...args: any[]) {
        super(props);
        this.state = {
            value: (typeof props.value === 'undefined' ? props.defaultValue : props.value) || '',
            width: 80,
        };
    }

    input: HTMLInputElement | null = null;

    static getDerivedStateFromProps(nextProps: TextInputProps, prevState: TextInputState) {
        if ('value' in nextProps && nextProps.value !== prevState.value) {
            return {
                value: nextProps.value,
            };
        }
        return null;
    }

    componentDidMount() {
        this.props.autoFocus && this.handleAutoFocus();
        if (!this.props.width) {
            this.resizeTextInput(this.props.placeholder || this.state.value);
        }
    }

    componentDidUpdate(prevProps: TextInputProps, prevState: TextInputState) {
        if (this.props.placeholder && this.props.placeholder !== prevProps.placeholder) {
            this.resizeTextInput(this.props.placeholder);
        }
    }

    /**
     * 更新值
     */
    private updateValue(value: string, callback?: () => void) {
        this.setState({ value }, callback);
    }

    /**
     * 处理输入框聚焦
     */
    private handleAutoFocus = () => {
        if (this.input) {
            this.input.focus();
        }
    };

    /**
     * 处理输入值发生变化
     */
    private handleInputValueChange = (event: ValueChangeEvent) => {
        const { detail } = event;
        const { value } = this.state;

        if (
            detail !== value &&
            ((!this.props.required && detail === '') ||
                (isFunction(this.props.validator) && this.props.validator(detail)))
        ) {
            !this.props.width && this.resizeTextInput(detail || this.props.placeholder || '');
            this.dispatchValueChangeEvent(detail);
        } else {
            event.preventDefault();
        }
    };

    /**
     * 改变输入框的宽度
     */
    private resizeTextInput = (text: string) => {
        const { minWidth = 80, maxWidth = 500 } = this.props;

        const span: HTMLSpanElement = document.createElement('span');
        document.body.appendChild(span);
        span.style.visibility = 'hidden';

        span.innerText = text;
        const width: number = span.offsetWidth + 10;
        span.parentNode && span.parentNode.removeChild(span);

        this.setState({
            width: width < maxWidth ? Math.max(width, minWidth) : maxWidth,
        });
    };

    private saveInput = (node: HTMLInputElement) => {
        this.input = node;
    };

    /**
     * 处理粘贴 todo
     */
    handlePaste = () => {
        setTimeout(() => {
            this.dispatchPasteEvent(this.state.value);
        });
    };

    private dispatchPasteEvent = createEventDispatcher(this.props.onPaste);

    /**
     * 处理鼠标点击
     */
    handleClick = (event: React.MouseEvent<HTMLInputElement>) => {
        isFunction(this.props.onClick) && this.props.onClick(event);
        event.stopPropagation();
    };

    /**
     * 处理键盘输入
     */
    handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.keyCode === 13) {
            event.preventDefault();
            isFunction(this.props.onPressEnter) && this.props.onPressEnter(event);
            return;
        }
        isFunction(this.props.onKeyDown) && this.props.onKeyDown(event);
    };

    /**
     * 处理聚焦
     */
    private handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        const { selectOnFocus } = this.props;

        if (selectOnFocus && this.input) {
            this.input.select();
            if (Array.isArray(selectOnFocus)) {
                this.input.selectionStart = selectOnFocus[0];
                this.input.selectionEnd = selectOnFocus[1] || this.input.value.length;
            }
        }
        isFunction(this.props.onFocus) && this.props.onFocus(event);
    };

    /**
     * 处理失焦
     */
    private handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
        isFunction(this.props.onBlur) && this.props.onBlur(event);
    };

    /**
     * 输入框聚焦
     */
    public focus() {
        this.input && this.input.focus();
    }

    /**
     * 输入框失焦
     */
    public blur() {
        this.input && this.input.blur();
    }

    /**
     * 输入框选中输入内容
     */
    public select() {
        this.input && this.input.select();
    }

    /**
     * 清空输入内容
     */
    clear() {
        this.updateValue('', () => this.resizeTextInput(''));
    }

    /**
     * 触发onValueChange
     */
    private dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange, ({ detail }) => {
        this.setState({ value: detail });
    });

    render() {
        const { className, disabled, placeholder, readOnly, maxLength, style, type } = this.props;
        const { value, width } = this.state;

        return (
            <BaseInput
                type={type}
                style={{ ...style, width: this.props.width || width }}
                className={classnames(className, styles['text-input'], { [styles['disabled']]: disabled })}
                {...{ value, disabled, placeholder, readOnly, maxLength }}
                onMounted={this.saveInput}
                onValueChange={this.handleInputValueChange.bind(this)}
                onClick={this.handleClick}
                onKeyDown={this.handleKeyDown}
                onFocus={this.handleFocus}
                onBlur={this.handleBlur}
                onPaste={this.handlePaste}
            />
        );
    }
}
