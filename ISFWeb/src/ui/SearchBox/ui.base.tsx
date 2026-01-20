import React from 'react';
import { noop, isFunction } from 'lodash';

export default class SearchBoxBase extends React.PureComponent<UI.SearchBox.Props, UI.SearchBox.State> {
    static defaultProps = {
        disabled: false,

        value: '',

        validator: (_value) => true,

        autoFocus: false,

        icon: '\uf01e',

        delay: 300,

        loader: noop,

        onFetch: noop,

        onLoad: noop,

        onLoadFailed: noop,

        onFocus: noop,

        onBlur: noop,

        onChange: noop,

        onEnter: noop,

        onClick: noop,

        onKeyDown: noop,
    }

    state: UI.SearchBox.State = {
        value: this.props.value,
        focus: false,
    }

    /*
     * 延迟触发搜索的定时器
     */
    timeout: number | null = null;

    searchInput: HTMLInputElement;

    /*
     * 正在执行的搜索
     */
    process: Promise<any> | null = null;

    constructor(props, ...args: any[]) {
        super(props);
        this.state = {
            focus: props.autoFocus,
        };
    }

    static getDerivedStateFromProps({ value }, prevState) {
        if (value !== prevState.value) {
            return {
                value,
            }
        }
        return null
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.state.value !== prevProps.value && this.state.value !== prevState.value) {
            this.fireChangeEvent(this.state.value)
        }
    }

    /**
     * 触发搜索
     * @param input 值
     */
    public load(input: string): void {
        this.searchInput.load(input)
    }

    /**
     * 输入发生变化时触发
     * @param value 文本值
     */
    protected handleChange(value: string): void {
        this.updateValue(value);
    }

    /**
     * 设置聚焦状态
     */
    protected handleFocus(event): void {
        this.setState({ focus: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event);
    }

    /**
     * 设置失焦状态
     */
    protected handleBlur(event): void {
        this.setState({ focus: false });
        isFunction(this.props.onBlur) && this.props.onBlur(event);
    }

    /**
     * 清空值
     */
    protected clearInput(): void {
        this.updateValue('');
        this.load('');
    }

    /**
     * 更新值并触发onChange事件
     * @param value 值
     */
    protected updateValue(value: string): void {
        this.setState({
            value,
        }, () => this.fireChangeEvent(value));
    }

    /**
     * 触发文本框变化事件
     * @param key 文本框输入值
     */
    private fireChangeEvent(key: string): void {
        isFunction(this.props.onChange) && this.props.onChange(key);
    }

    /**
     * 清空输入框的值
     */
    public clearInput() { // eslint-disable-line
        this.searchInput.clearInput()
        this.searchInput.focus();
    }

    /**
     * 输入框聚焦
     */
    public focus() {
        this.searchInput.focus();
    }

    /**
     * 输入框失焦
     */
    public blur() {
        this.searchInput.blur();
    }
}