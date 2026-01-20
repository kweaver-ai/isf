import React from 'react';
import classnames from 'classnames';
import { isFunction, includes, isEqual, find, filter } from 'lodash';
import { createEventDispatcher } from '../../utils/event';
import FlatList, { ItemLayoutCalculator, ItemRenderer } from '../FlatList';
import View from '../View';
import CheckBox from '../CheckBox';
import styles from './styles';

/**
 * 单选选中项
 */
export type SingleSelection = string | number | boolean | object | null;

/**
 * 多选选中项
 */
export type MultiSelection = ReadonlyArray<any>;

/**
 * 可能为多选或单选
 */
export type Selection = SingleSelection | MultiSelection;

interface DataListProps {
    data: ReadonlyArray<any>;

    getItemLayout: ItemLayoutCalculator;

    /**
     * 是否允许选中
     */
    enableSelect?: boolean;

    /**
     * 是否允许多选，只有 enableSelect: true 时有效
     */
    enableMultiSelect?: boolean;

    /**
     * 定义如何渲染每一行数据
     */
    renderItem: ItemRenderer;

    /**
     * 列表高度
     */
    height?: number | string;

    /**
     * 列表最大高度
     */
    maxHeight?: number;

    /**
     * 当前选中项
     */
    selection: Selection;

    /**
     * 高亮项下标
     */
    highlightIndex?: number;

    /**
     * 列表是否正在刷新
     */
    refreshing?: boolean;

    /**
     * 列表项样式
     */
    listItemClassName?: string;

    /**
     * key值计算
     */
    keyExtractor?: (record: any, index: number) => string;

    /**
     * 选中项改变时触发，重复选中同一项不会触发
     */
    onSelectionChange?: (item: any) => any;

    /**
     * 鼠标悬浮项改变时触发
     */
    onItemHoverChange: ({ prevIndex, currentIndex }: { prevIndex: number; currentIndex: number }) => void;

    /**
     * 当本次将要加载完成时触发 TODO 改一个合适的方法名
     */
    onLoadingComplete?: () => void;

    /**
     * 列表正在刷新的指示
     */
    RefreshingIndicatorComponent?: React.FunctionComponent<{ key: string; start: number; limit: number }>;

    /**
     * 列表为空的指示
     */
    ListEmptyComponent?: React.FunctionComponent<{ key: string; start: number; limit: number }>;
}

interface DataListState {
    /**
     * 高亮项下标
     */
    highlightIndex: number;
}

export default class DataList extends React.Component<DataListProps, DataListState> {

    static defaultProps = {
        enableSelect: true,
    }

    state = {
        highlightIndex: -1,
    }

    flatListRef = null

    // shift多选时，会根据此下标位置计算选中范围
    lastClickItemIndex: number = 0;

    // 记录上一次高亮下标
    lasthighlightIndex: number = -1;

    static getDerivedStateFromProps({ highlightIndex }: DataListProps, prevState: DataListState) {
        if (highlightIndex !== prevState.highlightIndex) {
            return {
                highlightIndex,
            }
        }

        return null
    }

    /**
     * 保存FlatList ref
     */
    saveFlatListRef = (node) => {
        this.flatListRef = node
    }

    public scrollToIndex({ index, viewPosition }: { index: number; viewPosition: number }) {
        this.flatListRef && this.flatListRef.scrollToIndex({ index, viewPosition })
    }

    /**
     * 处理CheckBox选中状态改变
     */
    handleCheckBoxValueChange = (item: any, event: React.ChangeEvent<HTMLInputElement>) => {
        const { selection = [] } = this.props;

        if (Array.isArray(selection)) {
            const nextSelection = event.target.checked === true ? [...selection, item] : filter(selection, (s) => !isEqual(s, item))

            this.dispatchSelectionChangeEvent(nextSelection)
        }
    }

    /**
     * 处理CheckBox点击事件
     */
    handleCheckBoxClicked = (event: React.MouseEvent<HTMLInputElement>) => {
        event.stopPropagation();
    }

    /**
     * 鼠标单击一行时触发
     */
    handleClickItem = (event, item, index) => {
        const { enableSelect } = this.props

        if (enableSelect) {
            const { enableMultiSelect, data = [], selection = [] } = this.props;
            let nextSelection

            if (enableMultiSelect && Array.isArray(selection)) {
                if (event.ctrlKey) { // ctrl键被按下
                    nextSelection = includes(selection, item) ? filter(selection, (s) => !isEqual(s, item)) : [...selection, item]

                    this.lastClickItemIndex = index
                } else if (event.shiftKey) { // shift键被按下
                    if (this.lastClickItemIndex < index) {
                        nextSelection = data.slice(this.lastClickItemIndex, index + 1)
                    } else if (this.lastClickItemIndex > index) {
                        nextSelection = data.slice(index, this.lastClickItemIndex + 1)
                    } else {
                        nextSelection = [data[this.lastClickItemIndex]]
                    }
                } else {
                    nextSelection = includes(selection, item) ? (selection.length > 1 ? [item] : []) : [item]

                    this.lastClickItemIndex = index
                }
            } else {
                nextSelection = item === selection ? null : item
            }
            this.dispatchSelectionChangeEvent(nextSelection)
        }
    }

    private dispatchSelectionChangeEvent = createEventDispatcher(this.props.onSelectionChange)

    /**
     * 鼠标双击时触发
     */
    private handleDoubleClickItem = (event, item, index) => {
        this.dispatchItemDoubleClickedEvent({ item, index })
    }

    private dispatchItemDoubleClickedEvent = createEventDispatcher(this.props.onItemDoubleClicked)

    /**
     * 处理鼠标移入列表
     */
    private handleMouseMove(index: number) {
        if (this.lasthighlightIndex !== index) {
            this.dispatchItemHoverChangeEvent({ prevIndex: this.lasthighlightIndex, currentIndex: index })
            this.lasthighlightIndex = index;
        }
    }

    /**
     * 处理鼠标移出列表
     */
    private handleMouseLeave = () => {
        this.dispatchItemHoverChangeEvent({ prevIndex: this.lasthighlightIndex, currentIndex: -1 })
        this.lasthighlightIndex = -1
    }

    private dispatchItemHoverChangeEvent = createEventDispatcher(this.props.onItemHoverChange)

    render() {
        const {
            height,
            selection,
            getItemLayout,
            maxHeight,
            enableMultiSelect,
            data,
            renderItem,
            onLoadingComplete,
            refreshing,
            RefreshingIndicatorComponent,
            ListEmptyComponent,
            keyExtractor,
            highlightIndex,
            listItemClassName,
        } = this.props;

        const renderFlatListItem = (dataItem: any, index: number) => (
            <View
                className={classnames(
                    styles['datalist-item'],
                    listItemClassName,
                    {
                        [styles['checkbox-padding']]: enableMultiSelect,
                        [styles['highlight']]: highlightIndex === index,
                        [styles['selected']]: enableMultiSelect ? includes(selection, dataItem) : isEqual(selection, dataItem),
                    },
                )}
                onMouseMove={() => this.handleMouseMove(index)}
                onMouseLeave={this.handleMouseLeave}
                onClick={(event) => this.handleClickItem(event, dataItem, index)}
                onDoubleClick={(event) => this.handleDoubleClickItem(event, dataItem, index)}
            >
                {
                    this.props.enableMultiSelect ?
                        <CheckBox
                            className={styles['checkbox']}
                            checked={includes(selection, dataItem) || !!find(selection, dataItem)}
                            onClick={(event) => this.handleCheckBoxClicked(event)}
                            onChange={(event) => this.handleCheckBoxValueChange(dataItem, event)}

                        />
                        : null
                }
                {renderItem(dataItem, index)}
            </View>
        );

        return (
            <View style={{ height }}>
                {
                    refreshing && RefreshingIndicatorComponent ? (
                        <View className={styles['refreshing-indicator-wrapper']}>
                            {isFunction(RefreshingIndicatorComponent) ? (
                                <RefreshingIndicatorComponent />
                            ) : null}
                        </View>
                    ) : null
                }
                {
                    data && data.length ?
                        <FlatList
                            ref={this.saveFlatListRef}
                            renderItem={renderFlatListItem}
                            {...{ keyExtractor, data, getItemLayout, height, maxHeight }}
                            onEndReached={onLoadingComplete}
                        />
                        : ListEmptyComponent ?
                            <View className={styles['empty-component-wrapper']}>
                                {isFunction(ListEmptyComponent) ? <ListEmptyComponent /> : ListEmptyComponent}
                            </View>
                            : null
                }
            </View>
        )
    }
}