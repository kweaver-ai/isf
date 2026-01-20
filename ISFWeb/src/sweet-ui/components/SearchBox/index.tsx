import React from 'react';
import classnames from 'classnames';
import { isFunction, noop } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import Control from '../Control';
import SearchInput from '../SearchInput';
import SweetIcon from '../SweetIcon';
import View from '../View';
import styles from './styles';

interface SearchBoxProps {
    /**
     * 图标是否在前
     */
    iconOnBefore?: string | React.ReactElement<any>;

    /**
     * 图标是否在后
     */
    iconOnAfter?: string | React.ReactElement<any>;

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
     * 样式
     */
    className?: string;

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
}

interface SearchBoxState {
    /**
     * 搜索关键字
     */
    value: string;

    /**
     * 搜索框聚焦状态
     */
    focus: boolean;
}

export default class SearchBox extends React.Component<SearchBoxProps, SearchBoxState> {

    static defaultProps = {
        disabled: false,

        validator: (_value: string) => true,

        autoFocus: false,

        delay: 300,

        allowClear: true,

        iconOnBefore: 'search',

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
        value: this.props.value,
        focus: false,
    }

    searchInput: HTMLInputElement;

    // 延迟触发搜索的定时器
    timeout: number | null = null;

    // 正在执行的搜索
    process: Promise<any> | null = null;

    componentDidUpdate(prevProps: SearchBoxProps, prevState: SearchBoxState) {
        if (this.props.value !== prevState.value) {
            this.setState({
                value: this.props.value,
            })
        }
    }

    /**
     * 渲染完成后获取input框ref
     */
    private saveSearchInput = (node) => {
        this.searchInput = node
    }

    /**
     * 更新值并触发onChange事件
     * @param value 值
     */
    private updateValue(value: string): void {
        this.setState({
            value,
        }, () => this.fireChangeEvent(value));
    }

    /**
     * 文本框变化触发搜索
     * @param key 关键字
     */
    private handleValueChange = (event: SweetUIEvent): void => {
        const { detail } = event;

        this.updateValue(detail)
    }

    /**
     * 触发搜索
     * @param key 输入值
     */
    public async load(key: string): Promise<any> {
        this.searchInput.load(key)
    }

    /**
     * 触发点击搜索
     * @param ref 文本框对象
     */
    private handleClick = (event: MouseEvent): void => {
        this.dispatchClickEvent(event);
    }

    /**
     * 处理聚焦
     */
    private handleFocus = (event: SweetUIEvent<FocusEvent>): void => {
        if (!this.props.disabled) {
            this.setState({ focus: true })
            this.dispatchFocusEvent(event);
        }
    }

    /**
     * 处理失焦
     */
    private handleBlur = (event: SweetUIEvent<FocusEvent>): void => {
        if (!this.props.disabled) {
            this.setState({ focus: false })
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
     * 触发onValueChange
     */
    private dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange);

    /**
     * 清空输入框的值
     */
    public clearInput = () => {
        this.searchInput.clearInput();
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
     * 触发onFocus
     */
    private dispatchFocusEvent = createEventDispatcher(this.props.onFocus);

    /**
     * 触发onBlur
     */
    private dispatchBlurEvent = createEventDispatcher(this.props.onBlur);

    /**
     * 触发onKeyDown
     */
    private dispatchKeyDownEvent = createEventDispatcher(this.props.onKeyDown);

    /**
     * 触发onPressEnter
     */
    private dispatchPressEnterEvent = createEventDispatcher(this.props.onPressEnter);

    render() {
        // 是否显示清空按钮
        const showClear = ({ allowClear, disabled }, { value }): boolean => allowClear && !disabled && !!value;
        const {
            disabled,
            width,
            className,
            iconOnBefore,
            iconOnAfter,
            placeholder,
            maxLength,
            autoFocus,
            validator,
            loader,
            delay,
            onClick,
            onPressEnter,
            onKeyDown,
            onError,
        } = this.props;

        const { value, focus } = this.state;

        return (
            <Control {...{ width, disabled, focus, className }}>
                <div className={styles['box']}>
                    {
                        iconOnBefore && (
                            <View className={styles['icon-before']} inline={true}>
                                {
                                    isFunction(iconOnBefore) ?
                                        iconOnBefore()
                                        :
                                        (
                                            <SweetIcon
                                                name={iconOnBefore || 'search'}
                                                color={'#cfcfcf'}
                                                size={16}
                                            />
                                        )
                                }
                            </View>
                        )
                    }
                    <View
                        inline={true}
                        className={classnames(
                            styles['input-content'],
                            { [styles['clear-indent']]: !iconOnAfter && showClear(this.props, this.state) },
                        )}
                    >
                        <SearchInput
                            ref={this.saveSearchInput}
                            width={'100%'}
                            {...{
                                value, maxLength, disabled, placeholder, autoFocus, validator,
                                loader, delay, onClick, onPressEnter, onKeyDown, onError,
                            }}
                            onValueChange={this.handleValueChange}
                            onFetch={this.dispatchFetchEvent}
                            onLoad={this.dispatchLoadEvent}
                            onFocus={this.handleFocus}
                            onBlur={this.handleBlur}
                            onClick={this.handleClick}
                        />
                    </View>
                    {
                        showClear(this.props, this.state) ?
                            <View className={styles['icon-after']} inline={true}>
                                <SweetIcon
                                    size={15}
                                    name={'clear'}
                                    className={styles['chip-x-icon']}
                                    onClick={this.clearInput}
                                />
                            </View>
                            : null
                    }
                    {
                        iconOnAfter && (
                            <View className={styles['icon-after']} inline={true}>
                                {
                                    isFunction(iconOnAfter) ?
                                        iconOnAfter()
                                        :
                                        (
                                            <SweetIcon
                                                name={iconOnAfter || 'search'}
                                                color={'#cfcfcf'}
                                                size={16}
                                                onClick={() => { this.focus(); this.searchInput.load(value) }}
                                            />
                                        )
                                }
                            </View>
                        )
                    }
                </div>
            </Control>
        )
    }
}
