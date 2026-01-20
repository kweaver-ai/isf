import React from 'react';
import { isFunction, noop, debounce, trim } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import TextInput from '../TextInput';

interface SearchInputProps {
    /**
     * 搜索框提示
     */
    placeholder?: string;

    /**
     * 输入框的值
     */
    value: string;

    /**
     * width，包含盒模型的padding和border
     */
    width?: number | string;

    /**
     * 触发搜索的延迟时间
     * TODO 默认值
     */
    delay?: number;

    /**
     * 支持清除
     */
    allowClear?: boolean;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 搜索输入字符最大长度
     */
    maxLength?: number;

    /**
     * 自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 搜索函数
     */
    loader: (key: string) => any;

    /**
     * 获取数据时触发
     */
    onFetch?: (key: string, process: Promise<any>) => any;

    /**
     * 数据加载完成时触发
     */
    onLoad?: (results: any) => any;

    onError?: (errorEvent: SweetUIEvent) => any;

    /**
     * 输入值发生变化时触发，传递value值
     */
    onValueChange?: (event: SweetUIEvent<string>) => void;

    /**
     * 输入限制函数
     */
    validator?: (value: string) => boolean;

    /**
     * 键盘输入时触发
     */
    onKeyDown?: (event: KeyboardEvent) => any;

    /**
     * 点击时触发
     */
    onClick?: (event: MouseEvent) => any;

    /**
     * 聚焦时的回调
     * @param event 文本框对象
     */
    onFocus?(event: FocusEvent): any;

    /**
     * 失焦时的回调
     * @param event 文本框对象
     */
    onBlur?(event: FocusEvent): any;

    /**
     * 回车触发
     */
    onPressEnter?: (event: KeyboardEvent) => any;

    onMounted: (ref: HTMLInputElement) => void;
}

interface SearchInputState {
    /**
     * 搜索关键字
     */
    value: string;
}

export default class SearchInput extends React.Component<SearchInputProps, SearchInputState> {

    static defaultProps = {
        disabled: false,

        validator: (_value: string) => true,

        autoFocus: false,

        delay: 300,

        loader: noop,

        onFetch: noop,

        onLoad: noop,

        onClick: noop,

        onFocus: noop,

        onBlur: noop,

        onChange: noop,

        onKeyDown: noop,

        onPressEnter: noop,
    }

    state = {
        value: '',
    }

    input: HTMLInputElement;

    // 延迟触发搜索的定时器
    timeout: number | null = null;

    // 正在执行的搜索
    process: Promise<any> | null = null;

    static getDerivedStateFromProps({ value }: { value: string }, prevState: SearchInputState) {
        if (value !== prevState.value) {
            return {
                value,
            }
        }

        return null
    }

    // 延迟 执行搜索
    private debounceLoad = debounce(this.load, this.props.delay, { leading: false, trailing: true })

    /**
     * 文本框变化触发搜索
     * @param key 关键字
     */
    private handleValueChange = (event: SweetUIEvent): void => {
        const { detail } = event
        if (this.process) {
            try {
                // 如果实现了abort方法则尝试调用
                this.process.abort();
            } catch (ex) {

            }

            this.process = null;
        }
        this.debounceLoad(detail);
        this.fireChangeEvent(detail)
    }

    /**
     * 触发搜索
     * @param key 输入值
     */
    public async load(key: string): Promise<any> {
        if (this.props.loader && (!key || trim(key))) {
            const value = trim(key)
            const process = this.props.loader(value)
            this.process = process;
            this.fireFetchEvent(value, process);
            try {
                const result = await this.promisify(process);
                this.fireLoadEvent(result);
            } catch (ex) {
                this.fireErrorEvent(ex)
            }
        }
    }

    /**
     * 触发点击搜索
     * @param ref 文本框对象
     */
    private handleClick = (event: MouseEvent): void => {
        if (!this.props.disabled) {
            this.triggerLoad(event.target.value);
            this.dispatchClickEvent(event);
        }
    }

    /**
     * 触发搜索
     * @param key 检索关键字
     */
    private triggerLoad(value) {
        this.load(value);
    }

    /**
     * 处理聚焦
     */
    private handleFocus = (event: SweetUIEvent<FocusEvent>): void => {
        if (!this.props.disabled) {
            this.dispatchFocusEvent(event);
        }
    }

    /**
     * 处理失焦
     */
    private handleBlur = (event: SweetUIEvent<FocusEvent>): void => {
        if (!this.props.disabled) {
            if (this.timeout) {
                clearTimeout(this.timeout);
            }
            this.dispatchBlurEvent(event);
        }
    }

    /**
     * 触发文本框变化事件
     * @param key 文本框输入值
     */
    private fireChangeEvent(key: string) {
        this.dispatchValueChangeEvent(key);
    }

    /**
     * 触发搜索进程
     * @param process 搜索进程
     */
    private fireFetchEvent(key: string, process: Promise<any>) {
        this.dispatchFetchEvent({ key, process });
    }

    /**
     * 触发load事件
     * @param result 搜索结果
     */
    private fireLoadEvent(result: any): void {
        this.dispatchLoadEvent(result);
    }

    /**
     * 触发搜索出错事件
     */
    private fireErrorEvent(ex: any) {
        this.dispatchErrorEvent(ex)
    }

    /**
     * 处理键盘输入
     */
    private handleKeyDown = (event: KeyboardEvent) => {
        this.dispatchKeyDownEvent(event)
    }

    /**
     * 处理按下enter键
     */
    private handlePressEnter = (event: KeyboardEvent) => {
        this.dispatchPressEnterEvent(event)
    }

    /**
     * 清空输入框的值
     */
    public clearInput = () => {
        this.handleValueChange({ detail: '' })
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
     * 触发onClick
     */
    private dispatchClickEvent = createEventDispatcher(this.props.onClick);

    /**
     * 触发onFetch
     */
    private dispatchFetchEvent = createEventDispatcher(this.props.onFetch);

    /**
     * 触发onLoad
     */
    private dispatchLoadEvent = createEventDispatcher(this.props.onLoad);

    /**
     * 触发onError
     */
    private dispatchErrorEvent = createEventDispatcher(this.props.onError)

    /**
     * 触发onValueChange
     */
    private dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange);

    /**
     * 触发onFocus
     */
    private dispatchFocusEvent = createEventDispatcher(this.props.onFocus);

    /**
     * 触发onBlur
     */
    private dispatchBlurEvent = createEventDispatcher(this.props.onBlur, () => this.setState({ focus: false }));

    /**
     * 触发onKeyDown
     */
    private dispatchKeyDownEvent = createEventDispatcher(this.props.onKeyDown);

    /**
     * 触发onPressEnter
     */
    private dispatchPressEnterEvent = createEventDispatcher(this.props.onPressEnter);

    /**
     * 将任何输入Promise化
     * @param input 输入值
    */
    private promisify(input: any): Promise<any> {
        return isFunction(input && input.then) ? input : Promise.resolve(input)
    }

    /**
     * 渲染完成后获取input框ref
     */
    private saveInput = (node) => {
        this.input = node
    }

    render() {

        const { disabled, width, placeholder, maxLength, autoFocus, validator, ...otherProps } = this.props;
        const { value } = this.state;

        return (
            <TextInput
                ref={this.saveInput}
                {...{ width, value, maxLength, disabled, placeholder, autoFocus, validator }}
                style={{ background: 'none' }}
                onValueChange={this.handleValueChange}
                onClick={this.handleClick}
                onFocus={this.handleFocus}
                onBlur={this.handleBlur}
                onPressEnter={this.handlePressEnter}
                onKeyDown={this.handleKeyDown}
            />
        )
    }
}
