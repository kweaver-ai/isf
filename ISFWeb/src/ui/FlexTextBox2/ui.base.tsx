import React from 'react';
import { isFunction } from 'lodash';

interface Props extends React.Props<any> {
    /**
     * 输入区域宽度，不传则宽度由输入的内容决定
     */
    width?: number;

    /**
     * 输入区域最小宽度，不传width时有效
     */
    minWidth?: number;

    /**
     * 输入区域最大宽度，不传width时有效
     */
    maxWidth?: number;

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
     * 点击事件
     */
    onClick?: (event: MouseEvent) => any;

    /**
     * 输入限制函数
     */
    validator?: (value: string) => boolean;

    /**
     * 聚焦事件
     */
    onFocus?: (event: FocusEvent) => any;

    /**
     * 失去焦点事件
     */
    onBlur?: (event: FocusEvent) => any;

    /**
     * 键盘输入时触发
     */
    onKeyDown?: (event: KeyboardEvent) => any;

    /**
     * 按下回车键触发
     */
    onPressEnter?: (event: KeyboardEvent) => any;

    /**
     * 粘贴事件
     */
    onPaste?: (event: React.ClipboardEvent<HTMLElement>) => void;

    /**
     * 输入值变化时触发
     */
    onValueChange: (value: string) => void;

}

interface State {
    value: string;

    placeholder?: string;
}

export default class FlexTextBoxBase extends React.Component<Props, State> {
    static defaultProps = {
        minWidth: 80,
        validator: () => true,
    }

    state = {
        value: '',
        width: 80,
    }

    input: HTMLInputElement;

    componentDidMount() {
        this.props.autoFocus && this.handleAutoFocus();
        if (!this.props.width) {
            this.resizeTextInput(this.props.placeholder || this.state.value)
        }
    }

    static getDerivedStateFromProps({ value, placeholder }, prevState) {
        if(value !== prevState.value) {
            return {
                value,
            }
        }
        return null
    }

    componentDidUpdate(prevProps, prevState) {
        if(!this.state.value && this.props.placeholder !== prevProps.placeholder) {
            this.resizeTextInput(this.props.placeholder)
        }
    }

    /**
     * 更新值
     */
    protected updateValue(value, callback) {
        this.setState({ value }, callback);
    }

    /**
     * 处理输入框聚焦
     */
    protected handleAutoFocus = () => {
        const { selectOnFocus } = this.props;

        this.input.focus();
        if (selectOnFocus) {
            this.input.select();
            if (Array.isArray(selectOnFocus)) {
                this.input.selectionStart = selectOnFocus[0]
                this.input.selectionEnd = selectOnFocus[1] || this.input.value.length
            }
        }
    }

    /**
     * 处理输入值发生变化
     */
    protected handleValueChange = (event) => {

        const input = event.target.value;

        const { value } = this.state;

        if (input !== value && ((!this.props.required && input === '') || this.props.validator(input))) {
            !this.props.width && this.resizeTextInput(input || this.props.placeholder)
            this.setState({
                value: input,
            })
            isFunction(this.props.onValueChange) && this.props.onValueChange(input)
        } else {
            event.preventDefault()
        }
    }

    /**
     * 改变输入框的宽度
     */
    public resizeTextInput = (text) => {
        const { minWidth, maxWidth } = this.props;

        const span = document.createElement('span');
        span.style = {
            position: 'absolute',
            top: '100%',
            left: '100%',
            visibility: 'hidden',
        }
        document.body.appendChild(span);

        span.innerText = text;
        const width = span.offsetWidth + 10;
        span.parentNode.removeChild(span);

        this.setState({
            width: width < maxWidth ? Math.max(width, minWidth) :maxWidth,
        })
    }

    protected saveInput = (node: HTMLInputElement) => {
        this.input = node;
    }

    /**
     * 处理鼠标点击
     */
    handleClick = (event: MouseEvent) => {
        isFunction(this.props.onClick) && this.props.onClick(event)
        event.stopPropagation();
    }

    /**
     * 处理键盘输入
     */
    handleKeyDown = (event: KeyboardEvent) => {
        if (event.keyCode === 13) {
            event.preventDefault ?
                event.preventDefault()
                : (event.returnValue = false);
            this.props.onPressEnter(event)
        } else {
            this.props.onKeyDown(event)
        }
    }

    /**
     * 处理聚焦
     */
    protected handleFocus = (event: FocusEvent) => {
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    /**
     * 处理失焦
     */
    protected handleBlur = (event: FocusEvent) => {
        this.blur()
        this.props.onBlur(event)
    }

    /**
     * 处理粘贴
     */
    handlePaste = (event) => {
        setTimeout(() => {
            isFunction(this.props.onPaste) && this.props.onPaste(this.state.value)
        })

    }

    /**
     * 输入框聚焦
     */
    public focus() {
        this.input.focus();
    }

    /**
     * 输入框失焦
     */
    public blur() {
        this.input.blur();
    }

    /**
     * 输入框选中输入内容
     */
    select() {
        this.input.select()
    }

    clear() {
        this.updateValue('', () => this.resizeTextInput(''))
    }

}