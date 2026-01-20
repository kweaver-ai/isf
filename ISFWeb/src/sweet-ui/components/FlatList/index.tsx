import React from 'react';
import classnames from 'classnames';
import { throttle, isFunction, isEqual } from 'lodash';
import View from '../View';
import styles from './styles';

/**
 * 数据的渲染模型
 * @property length 数据渲染的高度
 * @property offset 该条数据渲染距离列表顶端的偏移量
 * @property index 该条数据的下标
 */
interface ItemLayout {
    length: number;
    offset: number;
    index: number;
}

/**
 * 数据渲染尺寸计算
 */
export type ItemLayoutCalculator = (record: any, index: number) => ItemLayout;

/**
 * 渲染Item
 */
export type ItemRenderer = (record: any, index: number) => React.ReactNode;

/**
 * 当列表滚动接近底部时触发
 */
type EndReachedHandler = () => void;

interface FlatListProps extends React.ClassAttributes<void> {
    /**
     * 数据
     */
    data: ReadonlyArray<any>;

    /**
     * 如何渲染每一行数据
     */
    renderItem: ItemRenderer;

    /**
     * 获取每行数据的渲染高度，用于计算最终要渲染的数据切片
     * @param record 该行数据
     * @param index 该行数据下标
     */
    getItemLayout: ItemLayoutCalculator;

    /**
     * 根节点上的className
     */
    className?: string;

    /**
     * 高度
     */
    height?: number | string;

    /**
     * 最大高度
     */
    maxHeight?: number;

    /**
     * key值计算
     */
    keyExtractor?: (record: any, index: number) => string;

    /**
     * 滚动容器时触发
     */
    onScroll?: (event: React.UIEvent<HTMLDivElement>) => void;

    /**
     * 当列表滚动接近底部时触发
     */
    onEndReached?: EndReachedHandler;
}

interface FlatListState {
    itemsLayout: ReadonlyArray<ItemLayout>;

    /**
     * 切片范围，[start, end)
     */
    sliceRange: ReadonlyArray<number>;

    paddingTop?: number;

    paddingBottom?: number;
}

export default class FlatList extends React.Component<FlatListProps, FlatListState> {
    container: HTMLDivElement | null = null;

    state: FlatListState = {
        itemsLayout: [],
        sliceRange: [],
    };

    // 当前列表是否触发过滚动到最后
    hasEndReached: boolean = false;

    componentDidMount() {
        this.updateItemsLayout(this.props.data, this.setItemsLayout);
        this.updateListHeight()
    }

    componentDidUpdate(prevProps: FlatListProps, prevState: FlatListState) {
        this.updateListHeight()
        if (!isEqual(this.props.data, prevProps.data)) {
            this.hasEndReached = false;
            this.updateItemsLayout(this.props.data, this.setItemsLayout);
        }
    }

    /**
     * 更新列表高度
     */
    updateListHeight() {// TODO datalist
        if (this.props.maxHeight) {
            const { itemsLayout } = this.state;
            const [firstItem, ...last] = itemsLayout;

            if (firstItem && firstItem.length * this.props.data.length < this.props.maxHeight) {
                this.container.style.height = `${firstItem.length * this.props.data.length}px`
            } else {
                this.container.style.height = `${this.props.maxHeight}px`
            }
        }
    }

    /**
     * 计算并更新所有item布局
     * @param data
     * @param callback
     */
    private updateItemsLayout(data: ReadonlyArray<any> = [], callback?: () => void) {
        const { getItemLayout } = this.props;

        this.setState(
            {
                itemsLayout: data.map((item, index) => ({ index, ...getItemLayout(item, index) })),
            },
            callback,
        );
    }

    /**
     * 根据getItemLayout()计算渲染切片和内容上下padding值
     */
    private setItemsLayout() {
        if (this.container) {
            const { data, onEndReached } = this.props;

            // 获取视口的尺寸并计算所有数据的高度，将落在视口内的数据的下标范围作为要显示的数据的索引返回
            const { clientHeight, scrollTop } = this.container;
            const { itemsLayout } = this.state;

            const { paddingTop, sliceRange, paddingBottom } = itemsLayout.reduce(
                (prev, { length, offset, index }, i) => {
                    const { paddingTop, sliceRange, paddingBottom, contentOffset } = prev;
                    const [start, end] = sliceRange;
                    const nextContentOffset = contentOffset + length;

                    // 当内容尚未进入到视口
                    if (nextContentOffset < scrollTop) {
                        return {
                            ...prev,
                            paddingTop: nextContentOffset,
                            contentOffset: nextContentOffset,
                        };
                    } else {
                        // 当内容进入视口
                        // 内容刚开始进入视口
                        if (start === -1) {
                            return {
                                ...prev,
                                sliceRange: [i, end],
                                contentOffset: nextContentOffset,
                            };
                        }
                        // 内容开始溢出视口
                        if (nextContentOffset >= scrollTop + clientHeight) {
                            // 刚好溢出视口
                            if (end === -1) {
                                return {
                                    ...prev,
                                    sliceRange: [start, i + 1],
                                    contentOffset: nextContentOffset,
                                };
                            } else {
                                // 视口外的内容
                                return {
                                    ...prev,
                                    paddingBottom: paddingBottom + length,
                                    contentOffset: nextContentOffset,
                                };
                            }
                        }

                        // 落在视口内的内容，正常渲染
                        return {
                            ...prev,
                            contentOffset: nextContentOffset,
                        };
                    }
                },
                { paddingTop: 0, sliceRange: [-1, -1], paddingBottom: 0, contentOffset: 0 },
            );

            this.setState({ paddingTop, sliceRange, paddingBottom }, () => {
                const [, end] = sliceRange;

                // 滚动到底部
                if (this.state.paddingTop && end === data.length && !this.hasEndReached) {
                    if (isFunction(onEndReached)) {
                        this.hasEndReached = true;
                        onEndReached();
                    }
                }
            });
        }
    }

    /**
     * 处理容器滚动事件
     * @param event UIEvent事件对象
     */
    private handleScroll(event: React.UIEvent<HTMLDivElement>) {
        const { onScroll } = this.props;

        this.setItemsLayout();

        if (isFunction(onScroll)) {
            onScroll(event);
        }
    }

    /**
     * 容器滚动事件懒触发
     */
    private handleScrollLazily = throttle(this.handleScroll, 1000 / 24, { trailing: true, leading: false });

    saveListRef = (node) => {
        this.container = node;
    }

    /**
     * 控制滚动条滚动到指定下标的列表项
     * @param index
     * @param viewPosition 值为'top'时，index指定的列表项位于顶部，'bottom'位于底部
     */
    public scrollToIndex({ index: paramIndex, viewPosition }) {// 尽可能接近指示的位置
        const { itemsLayout, sliceRange } = this.state;
        const [start, end] = sliceRange;
        const { clientHeight, scrollTop } = this.container;
        // eslint-disable-next-line no-redeclare
        const { index, length, offset } = itemsLayout.find((item) => item.index === paramIndex)
        const itemOffset = length * (index + 1) // 这里获得的length只是预估的高度，不一定是真实的高度 TODO待修改

        if (viewPosition === 'top') { // 顶对齐
            // 判断相对视口的位置
            if ((this.props.data.length - index) * length >= clientHeight) {
                this.container.scrollTop = length * index
            } else {
                this.container.scrollTop = this.props.data.length * length - clientHeight
            }
        } else if (viewPosition === 'bottom') { // 底对齐
            if (index <= start) { // 全部或部分在视口以上
                if (itemOffset >= clientHeight) {
                    this.container.scrollTop = Math.abs(itemOffset - scrollTop)
                } else {
                    this.container.scrollTop = 0
                }
            } else if (index >= end - 1) { // 全部或部分在视口以下
                this.container.scrollTop = itemOffset - clientHeight
            }
        }
    }

    render() {
        const {
            data,
            renderItem,
            height,
            className,
            keyExtractor,
        } = this.props;
        const { paddingTop, sliceRange, paddingBottom } = this.state;
        const [start, end] = sliceRange;

        return (
            <View
                onMounted={this.saveListRef}
                onScroll={this.handleScrollLazily.bind(this)}
                className={classnames(styles['root'], className)}
                style={{ height }}
            >
                {
                    // 只有计算好了视口中要显示的数据切片才能进行渲染
                    start !== -1 || end !== -1 ? (
                        <ol className={styles['list']} style={{ paddingTop, paddingBottom }}>
                            {data
                                .slice(start, end === -1 ? undefined : end)
                                .map((record, index) => (
                                    <li
                                        key={isFunction(keyExtractor) ? keyExtractor(record, start + index) : index}
                                    >
                                        {renderItem(record, start + index)}
                                    </li>
                                ))}
                        </ol>
                    ) : null
                }
            </View>
        );
    }
}
