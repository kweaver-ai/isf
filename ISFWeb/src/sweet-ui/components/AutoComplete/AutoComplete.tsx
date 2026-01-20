import React from 'react';
import classnames from 'classnames'
import { isFunction, noop, isArray, isEqual, trim } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import { ItemLayoutCalculator, ItemRenderer } from '../FlatList';
import SearchBox from '../SearchBox';
import SweetIcon from '../SweetIcon'
import DataList, { Selection } from '../DataList';
import Portal from '../Portal';
import View from '../View';
import Locator from '../PopOver/Locator';
import Button from '../Button';
import Icon from '../Icon'
import CheckBox from '../CheckBox';
import styles from './styles';
import __ from './locale';
import loading from './assets/loading.gif'
import AppConfigContext from '@/core/context/AppConfigContext';

export interface AutoCompleteProps {
    /**
     * 搜索关键字
     */
    value: string;

    /**
     * 搜索框提示
     */
    placeholder?: string;

    /**
     * width，包含盒模型的padding和border
     */
    width?: number | string;

    /**
     * 搜索框后的图标
     */
    iconOnAfter?: string | React.ReactElement<any>;

    /**
     * 搜索框前的图标
     */
    iconOnBefore?: string | React.ReactElement<any>;

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
     * 自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 搜索输入字符最大长度
     */
    maxLength?: number;

    /**
     * 列表数据
     */
    data: ReadonlyArray<any>;

    /**
     * 加载数据的起始索引
     * 默认为0
     */
    start?: number;

    /**
     * 每次加载的最大数据条数
     * 默认为20
     */
    limit?: number;

    /**
     * 列表容器显示的最大高度
     */
    maxHeight?: number;

    /**
     * 输入限制函数
     */
    validator?: (value: string) => boolean;

    /**
     * 校验状态
     */
    validateStatus?: 'error' | 'normal';

    /**
     * 是否多选
     */
    enableMultiSelect?: boolean;

    /**
     * 是否可以全选
     */
    enableSelectAll?: boolean;

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
     * 点击搜索框
     */
    onClickInput?(event: MouseEvent): any;

    /**
     * 搜索函数
     */
    loader: ({ key, start, limit }: { key: string; start: number; limit: number }) => Promise<ReadonlyArray<any>>;

    /**
     * 获取数据时触发
     */
    onFetch?: (key: string, process: Promise<any>) => any;

    /**
     * 数据加载完成时触发
     */
    onLoad?: (results: any) => any;

    /**
     * 数据加载错误
     */
    onError?: (errorEvent: SweetUIEvent<any>) => any;

    /**
     * 数据尺寸
     */
    getItemLayout: ItemLayoutCalculator;

    /**
     * 定义如何渲染每一行数据
     */
    renderItem: ItemRenderer;

    /**
     * 额外的项
     */
    dropdownRender?: () => React.ReactDOM;

    /**
     * key值计算
     */
    keyExtractor?: (record: any, index: number) => string;

    /**
     * 高亮选项发生变化时触发
     */
    onHighlightChange: ({ prevIndex, currentIndex }: { prevIndex: number; currentIndex: number }) => void;

    // TODO 上下键切换高亮项时，需要界面区分循环/不循环两种模式
    /**
     * 上下键切换高亮项时，控制列表是否循环滚动
     */
    circle: boolean;

    /**
     * 被选中时调用，参数为选中项
     */
    onSelect?: (item: any) => void;

    /**
     * 当前选中项
     */
    selection: Selection;

    /**
     * 按回车键时触发，如果当前有高亮项则作为参数传递
     */
    onPressEnter: (item: any | null) => void;

    /**
     * 输入值发生变化时触发，传递value值
     */
    onValueChange?: (event: SweetUIEvent<string>) => void;

    /**
     * 鼠标划入
     */
    onMouseEnter?: (event: MouseEvent) => void;

    /**
     * 鼠标划出
     */
    onMouseLeave?: (event: MouseEvent) => void;

    /**
     * 弹出层容器
     */
    element?: HTMLElement;

    /**
     * 列表为空时显示
     */
    ListEmptyComponent?: React.FunctionComponent<{ key: string; start: number; limit: number }>;

    /**
     * 列表正在刷新的指示
     */
    RefreshingIndicatorComponent?: React.FunctionComponent<{ key: string; start: number; limit: number }>;

    /**
     * 搜索结果出错
     */
    ErrorComponent?: React.FunctionComponent<{ key: string; start: number; limit: number }>;

    /**
     * 当列表滚动接近底部时加载更多的指示
     */
    EndReachedIndicatorComponent?: React.FunctionComponent<void>;

}

interface AutoCompleteState {
    /**
     * 是否正在加载数据
     */
    refreshing: boolean;

    /**
     * 是否正在加载更多数据
     */
    loadingMore: boolean;

    /**
     * 当前高亮项下标
     */
    highlightIndex: number;

    /**
     * 下拉列表激活状态
     */
    active: boolean;

    /**
     * 搜索框的值
     */
    value: string;

    /**
     * 加载状态
     */
    status: SearchStatus;

    /**
     * 是否是全选状态
     */
    isSelectedAll: boolean;

    /**
     * 当前选中项
     */
    selection: Selection;
}

/**
 * 搜索状态
 */
export enum SearchStatus {
    /**
     * 无操作
     */
    Pending,

    /**
     * 正在搜索
     */
    Fetching,

    /**
     * 搜索出错
     */
    SearchError,

    /**
     * 完成搜索
     */
    Ok,
}

export enum KeyCode {
    Tab = 9,
    Enter = 13,
    UpArrow = 38,
    DownArrow = 40,
}

export default class AutoComplete extends React.Component<AutoCompleteProps, AutoCompleteState> {

    static contextType = AppConfigContext

    static defaultProps = {
        data: [],
        start: 0,
        limit: 20,
        maxHeight: 300,
        validator: (value: string) => true,
        onFocus: noop,
        onBlur: noop,
        loader: noop,
        onHighlightChange: noop,
        circle: false,
        onSelect: noop,
        onPressEnter: noop,
    }

    state = {
        value: this.props.value,
        refreshing: false,
        loadingMore: false,
        highlightIndex: -1,
        active: false,
        status: SearchStatus.Pending,
        isSelectedAll: false,
        selection: this.props.enableMultiSelect ? [] : null,
    }

    // 用作下拉列表的定位锚点
    container: HTMLDivElement

    // 搜索框 ref
    searchBox: React.ReactNode

    dataListRef: React.ReactNode

    /**
     * 弹出层元素
     */
    popupInstance: Element | null = null;

    start = 0;

    keepActive = false

    static getDerivedStateFromProps({ value }: { value: string }, prevState: AutoCompleteState) {
        if (value !== prevState.value) {
            return {
                value,
                highlightIndex: -1,
                isSelectedAll: false,
            }
        }

        return null
    }

    componentDidUpdate(prevProps, prevState) {
        const { selection, data, enableSelectAll, value } = this.props

        if (selection !== this.state.selection) {
            this.setState({
                selection,
                isSelectedAll: isArray(selection) && data.every((item) => selection.includes(item)),
            })
        }

        if (data !== prevProps.data && enableSelectAll && isArray(selection)) {
            this.setState({
                isSelectedAll: selection.length >= data.length && data.every((item) => selection.includes(item)),
            })
        }
    }

    private saveContainerRef = (node: HTMLDivElement) => {
        this.container = node
    }

    /**
     * 保存搜索框 ref
     */
    private saveSearchBoxRef = (node: React.ReactNode) => {
        this.searchBox = node
    }

    private saveDataListRef = (node: React.ReactNode) => {
        this.dataListRef = node
    }

    private savePopupRef = (node: Element) => {
        this.popupInstance = node
    }

    private getContainer = () => {
        const element = isFunction(this.props.element) ? this.props.element() : this.props.element || this.context?.element
        
        const popupContainer = document.createElement('div');
        (((element && isArray(element) ? element[0] : element) || window.document.querySelector('body')) as HTMLBodyElement).appendChild(popupContainer as HTMLDivElement);

        return popupContainer
    }

    public focus = () => {
        this.searchBox.focus()
    }

    public blur = () => {
        this.searchBox.blur()
    }

    private preventHideResults = () => {
        this.keepActive = true
    }

    protected handleFetch = (event: SweetUIEvent<any>) => {
        this.toggleActive(true)
        const { detail: { key, process } } = event

        this.setState({
            value: key,
            status: this.start === 0 ? SearchStatus.Fetching : this.state.status,
            refreshing: this.start === 0 ? true : false,
        })

        this.dispatchFetchEvent({ key, process })
    }

    private handleLoad = (event: SweetUIEvent<any>) => {
        const { detail } = event

        this.toggleActive(true)

        this.setState({
            status: SearchStatus.Ok,
            refreshing: false,
            loadingMore: false,
        })

        this.dispatchLoadEvent(detail)
    }

    private handleError = (event: SweetUIEvent<any>) => {
        const { detail } = event

        if (detail !== 'CANCEL') {
            this.setState({ status: SearchStatus.SearchError })

            this.dispatchErrorEvent(detail)
        }
    }

    private handleFocus = (event: SweetUIEvent<FocusEvent>) => {
        this.focus()
        this.dispatchFocusEvent(event)
    }

    private handleBlur = (event: SweetUIEvent<FocusEvent>) => {
        if (this.keepActive) {
            this.searchBox.focus()
        } else {
            this.blur()

            this.dispatchBlurEvent(event)
        }
        // event.preventDefault() // 问题：不选中结果项直接失焦搜索框时，list部分仍显示 TODO处理收起列表

        this.keepActive = false
    }

    public toggleActive = (active: boolean) => {
        this.setState({ active })

        if (!active) {
            this.setState({
                status: SearchStatus.Pending,
                highlightIndex: -1,
            })

            this.start = 0
        }
    }

    protected handleValueChange = (event: SweetUIEvent<any>) => {
        const { detail } = event

        if (trim(detail)) {
            // SearchBox 有输入且值不为空的时候重置start和limit，并触发loader
            this.start = 0

            this.setState({
                status: SearchStatus.Fetching,
            })
        } else {
            this.toggleActive(false)
        }

        this.dispatchValueChangeEvent(detail)
    }

    /**
     * 键盘事件处理，仅改变高亮项下标
     */
    private handleKeyDown = (event: SweetUIEvent<KeyboardEvent>) => {
        if (this.dataListRef) {
            switch (event.detail.keyCode) {
                case KeyCode.DownArrow:
                    event.preventDefault ? event.preventDefault() : (event.returnValue = false);
                    this.setState({
                        highlightIndex: this.state.highlightIndex + 1 >= this.props.data.length ?
                            this.props.circle ? 0 : this.state.highlightIndex
                            : this.state.highlightIndex + 1,
                    }, () => {
                        this.dataListRef.scrollToIndex({
                            index: this.state.highlightIndex,
                            viewPosition: 'bottom',
                        })
                    })

                    break

                case KeyCode.UpArrow:
                    event.preventDefault ? event.preventDefault() : (event.returnValue = false);
                    this.setState({
                        highlightIndex: this.state.highlightIndex - 1 < 0 ?
                            this.props.circle ? this.props.data.length - 1 : 0
                            : this.state.highlightIndex - 1,
                    }, () => {
                        this.dataListRef.scrollToIndex({ index: this.state.highlightIndex, viewPosition: 'top' })
                    })
                    break

                default:
                    break
            }
        }
    }

    /**
     * 输入框enter事件
     */
    private handlePressEnter = (event: SweetUIEvent<KeyboardEvent>) => {
        const { enableMultiSelect, data } = this.props

        this.dispatchPressEnterEvent(
            enableMultiSelect ?
                {
                    selection: [data[this.state.highlightIndex]],
                    isSelectedAll: this.state.isSelectedAll,
                }
                : data[this.state.highlightIndex],
        )

        if (!this.props.enableMultiSelect) {
            this.toggleActive(false)
        }
    }

    /**
     * 选中项发生变化
     */
    private handleSelectionChange = (event: SweetUIEvent<any>) => {
        const { detail } = event

        if (this.props.enableMultiSelect && isArray(detail)) {
            const newSelections = this.state.value ?
                Array.from(new Set([...this.state.selection.filter((item) => !this.props.data.includes(item)), ...detail]))
                : detail

            this.setState({
                selection: newSelections,
                isSelectedAll: isEqual(newSelections, this.props.data),
            }, () => {
                this.dispatchSelectEvent({ selection: newSelections, isSelectedAll: this.state.isSelectedAll })
            })
        } else {

            this.setState({
                selection: detail,
            })

            this.dispatchSelectEvent(detail)

            this.toggleActive(false)

        }
    }

    /**
     * 处理悬浮项发生改变
     */
    private handleItemOverChange = (event) => {
        const { detail: { prevIndex, currentIndex } } = event

        if (currentIndex !== -1) {
            this.setState({
                highlightIndex: currentIndex,
            })

            this.dispatchHighlightChange({ prevIndex, currentIndex })
        }
    }

    /**
     * 加载更多
     */
    private loadMore = () => {
        if (this.searchBox && this.props.data.length % this.props.limit === 0) {

            // 显示加载更多的提示
            this.setState({
                loadingMore: true,
            }, () => {
                this.start = this.start + this.props.limit
                this.searchBox.load(this.state.value)
            })
        }
    }

    /**
     * 重试
     */
    private reload = () => {
        this.searchBox.load(this.state.value)
    }

    /**
     * 全选
     */
    private selectAll = (event: MouseEvent) => {
        this.searchBox.focus()

        const { selection, value } = this.state

        const nextSelection = event.target.checked ?
            [...selection, ...this.props.data.filter((item) => !selection.includes(item))]
            : value ? selection.filter((item) => !this.props.data.includes(item)) : []

        this.setState({
            isSelectedAll: !this.state.isSelectedAll,
            selection: nextSelection,
        }, () => {
            this.dispatchSelectEvent({ selection: nextSelection, isSelectedAll: this.state.isSelectedAll })
        })
    }

    private handleClickAway = (e: MouseEvent) => {
        if (this.popupInstance && !this.popupInstance.contains(e.target)) {
            this.toggleActive(false)
        } else {
            this.searchBox.focus()
        }
    }

    /**
     * 清空输入框的值
     */
    public clearInput = () => {
        this.searchBox.clearInput()
    }

    /**
     * 懒加载
     */
    private lazyLoader = (key: string) => this.props.loader({ key, start: this.start, limit: this.props.limit });

    private dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange)

    private dispatchFocusEvent = createEventDispatcher(this.props.onFocus)

    private dispatchBlurEvent = createEventDispatcher(this.props.onBlur)

    private dispatchFetchEvent = createEventDispatcher(this.props.onFetch)

    private dispatchLoadEvent = createEventDispatcher(this.props.onLoad)

    private dispatchErrorEvent = createEventDispatcher(this.props.onError)

    private dispatchHighlightChange = createEventDispatcher(this.props.onHighlightChange)

    private dispatchSelectEvent = createEventDispatcher(this.props.onSelect)

    private dispatchPressEnterEvent = createEventDispatcher(this.props.onPressEnter);

    render() {
        const {
            data,
            width,
            disabled,
            validateStatus,
            delay,
            allowClear,
            autoFocus,
            placeholder,
            validator,
            maxLength,
            iconOnAfter,
            iconOnBefore,
            circle,
            maxHeight,
            keyExtractor,
            renderItem,
            getItemLayout,
            ErrorComponent,
            enableMultiSelect,
            enableSelectAll,
            RefreshingIndicatorComponent,
            ListEmptyComponent,
            EndReachedIndicatorComponent,
            dropdownRender,
            onMouseEnter = noop,
            onMouseLeave = noop,
            onClickInput = noop,
        } = this.props;

        const {
            value,
            highlightIndex,
            refreshing,
            loadingMore,
            status,
            isSelectedAll,
            selection,
        } = this.state;

        return (
            <View
                onMounted={this.saveContainerRef}
                inline={true}
                style={{ width }}
                className={styles['autocomplete-box']}
                {...{ onMouseEnter, onMouseLeave }}
            >
                <SearchBox
                    {...{ value, width, disabled, maxLength, iconOnBefore, iconOnAfter, delay, allowClear, autoFocus, placeholder, validator }}
                    ref={this.saveSearchBoxRef}
                    className={classnames(
                        this.props.className,
                        { [styles['error']]: validateStatus === 'error' },
                    )}
                    loader={this.lazyLoader}
                    onValueChange={this.handleValueChange}
                    onFetch={this.handleFetch}
                    onLoad={this.handleLoad}
                    onError={this.handleError}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onClick={onClickInput}
                    onPressEnter={this.handlePressEnter.bind(this)}
                    onKeyDown={this.handleKeyDown}
                />
                {
                    validateStatus === 'error' ?
                        <SweetIcon
                            name={'caution'}
                            size={16}
                            color={'#e60012'}
                            className={styles['caution']}
                        /> : null
                }
                {
                    this.state.active && (
                        data.length ||
                        ListEmptyComponent ||
                        (loadingMore && EndReachedIndicatorComponent) ||
                        (status === SearchStatus.SearchError && ErrorComponent)
                    ) ?
                        <Portal getContainer={this.getContainer}>
                            <Locator
                                anchor={this.container}
                                anchorOrigin={['left', 'bottom']}
                                alignOrigin={['left', 'top']}
                                className={styles['locator']}
                                onMouseDown={this.handleClickAway}
                                element={this.props.element}
                            >
                                <View
                                    className={styles['list-panel']}
                                    style={{ width }}
                                    onMounted={this.savePopupRef}
                                    onMouseDown={this.preventHideResults}
                                >
                                    {
                                        status === SearchStatus.Ok ?
                                            [
                                                (
                                                    data.length && enableSelectAll ?
                                                        <CheckBox
                                                            className={styles['checkbox']}
                                                            checked={isSelectedAll}
                                                            indeterminate={!isSelectedAll && (!value && selection.length || data.some((item) => selection.includes(item)))}
                                                            onClick={(event) => event.stopPropagation()}
                                                            onChange={this.selectAll}
                                                        >
                                                            {__('全选')}
                                                        </CheckBox>
                                                        : null
                                                ),
                                                (
                                                    <DataList
                                                        key={'list'}
                                                        data={data}
                                                        ref={this.saveDataListRef}
                                                        {...{ circle, renderItem, keyExtractor, refreshing, highlightIndex, maxHeight, selection, enableMultiSelect, enableSelectAll }}
                                                        getItemLayout={isFunction(getItemLayout) ? getItemLayout : () => ({ length: 34 })}
                                                        onLoadingComplete={this.loadMore}
                                                        onSelectionChange={this.handleSelectionChange}
                                                        ListEmptyComponent={isFunction(ListEmptyComponent) ? () => ListEmptyComponent({ key: value, start: this.start, limit: this.props.limit }) : ListEmptyComponent}
                                                        RefreshingIndicatorComponent={() => isFunction(RefreshingIndicatorComponent) && RefreshingIndicatorComponent({ key: value, start: this.start, limit: this.props.limit })}
                                                        onItemHoverChange={this.handleItemOverChange}
                                                        listItemClassName={styles['list-item']}
                                                    />
                                                ),
                                            ]
                                            : null
                                    }
                                    {
                                        status === SearchStatus.Fetching ? (
                                            <div>
                                                {
                                                    isFunction(RefreshingIndicatorComponent) ?
                                                        <RefreshingIndicatorComponent />
                                                        : (
                                                            <div className={styles['tip']}>
                                                                <Icon src={loading} />
                                                                <span className={styles['loading']}>{__('加载中...')}</span>
                                                            </div>
                                                        )
                                                }
                                            </div>
                                        ) : null
                                    }
                                    {
                                        loadingMore && EndReachedIndicatorComponent ? (
                                            <div className={styles['endreached-component-wrapper']}>
                                                {isFunction(EndReachedIndicatorComponent) ? <EndReachedIndicatorComponent /> : null}
                                            </div>
                                        ) : null
                                    }
                                    {
                                        status === SearchStatus.SearchError ? (
                                            <div>
                                                {
                                                    isFunction(ErrorComponent) ?
                                                        <ErrorComponent />
                                                        : (
                                                            <div className={styles['tip']}>
                                                                {__('加载失败，')}
                                                                <Button
                                                                    theme={'text'}
                                                                    onClick={this.reload}
                                                                >
                                                                    {__('重试')}
                                                                </Button>
                                                            </div>
                                                        )
                                                }
                                            </div>
                                        ) : null
                                    }
                                    {
                                        isFunction(dropdownRender) ? dropdownRender() : null
                                    }
                                </View>
                            </Locator>
                        </Portal>
                        : null
                }
            </View>
        )
    }
}
