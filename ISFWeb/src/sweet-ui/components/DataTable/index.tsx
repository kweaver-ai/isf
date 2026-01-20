import React from 'react';
import classnames from 'classnames';
import { without, isFunction, includes, throttle } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import BaseTable, { BaseTableColumn, type RowKeyExtractor, type CellKeyExtractor, type RowExtraComponent, type RowExtraMatcher } from '../BaseTable';
import ContentView from '../ContentView';
import CheckBox from '../CheckBox';
import View from '../View';
import SweetIcon from '../SweetIcon'
import styles from './styles';

export { RowKeyExtractor, CellKeyExtractor };

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

/**
 * 行样式
 */
export type RowClassName = string;

/**
 * 行悬浮样式
 */
export type RowHoverClassName = string;

/**
 * 过滤菜单项
 */
type Filter = {
    text: string;
    value: string;
    checked: boolean;
}

/**
 * 列配置项
 */
export interface DataTableColumn extends BaseTableColumn {
    /**
     * 列标题
     */
    title?: string;

    /**
     * 列是否允许排序
     */
    sortable?: boolean;

    /**
     * 表头的过滤菜单项
     */
    filters?: ReadonlyArray<Filter>;

    /**
     * 表头的附加提示说明
     */
    tips?: string | React.ReactNode;
}

export interface ToolbarComponentProps {
    /**
     * 列表中的所有数据
     */
    data: ReadonlyArray<any>;

    /**
     * 当前选中项
     */
    selection: Selection;

    /**
     * 调用会选中所有项
     */
    selectAll: (selected: boolean) => void;
}

/**
 * 选中事件处理函数
 */
export type SelectionChangeHandler = (event: SweetUIEvent<Selection>) => void;

/**
 * 选中所有行事件处理函数
 */
export type ChangeSelectAllHandler = (event: SweetUIEvent<boolean>) => void;

/**
 * 双击事件处理
 */
export type RowDoubleClickedHandler = (event: SweetUIEvent<{ record: any; index: number }>) => void;

/**
 * 工具栏组件
 */
export type DataTableToolbarComponent = React.FunctionComponent<ToolbarComponentProps> | React.ReactElement<any> | null;

/**
 * 表头组件
 */
export type DataTableHeaderComponent =
    | React.FunctionComponent<{
        /**
         * 列表中的所有数据
         */
        data: ReadonlyArray<any>;

        /**
         * 当前选中项
         */
        selection: Selection;

        /**
         * 调用会选中所有项
         */
        selectAll: (selected: boolean) => void;

        /**
         * 列配置
         */
        columns: ReadonlyArray<DataTableColumn>;
    }>
    | React.ReactElement<any> |
    null;

/**
 * 页脚组件
 */
export type DataTableFooterComponent =
    | React.FunctionComponent<{
        /**
         * 列表中的所有数据
         */
        data: ReadonlyArray<any>;
    }>
    | React.ReactElement<any>
    | null;

/**
 * 内容为空提示组件
 */
export type DataTableEmptyComponent = React.FunctionComponent<{}> | React.ReactElement<any>;

/**
 * 正在刷新提示组件
 */
export type DataTableRefreshingComponent = React.ReactNode | React.SFC<{}> | React.Component<any, any>;

/**
 * 是否允许选中
 */
export type EnableSelect = boolean;

/**
 * 是否允许多选
 */
export type EnableMultiSelect = boolean;

/**
 * 展开/收起行
 */
export type ExpandChangeHandler = (event: SweetUIEvent<{ expandedKeys: ReadonlyArray<string>; record: any; expanded: boolean }>) => void;

export { RowExtraComponent }

export { RowExtraMatcher }

export interface DataTableProps extends React.ClassAttributes<void> {
    /**
     * 角色
     */
    role?: string;

    /**
     * 列配置数组
     */
    columns: ReadonlyArray<DataTableColumn>;

    /**
     * 列表中的所有数据
     */
    data: ReadonlyArray<any>;

    /**
     * 是否显示表格的边框
     */
    showBorder?: boolean;

    /**
     * 总体高度
     */
    height?: number | string;

    /**
     * 是否正在刷新
     */
    refreshing?: boolean;

    /**
     * 行高
     */
    rowHeight?: number;

    /**
     * 列表行的类名
     */
    rowClassName?: RowClassName;

    /**
     * 鼠标悬浮在列表行上的样式
     */
    rowHoverClassName?: RowHoverClassName;

    /**
     * 计算行的key属性
     */
    rowKeyExtractor?: RowKeyExtractor;

    /**
     * 计算行额外信息的key属性
     */
    rowExtraKeyExtractor?: RowKeyExtractor;

    /**
     * 计算单元格的key属性
     */
    cellKeyExtractor?: CellKeyExtractor;

    /**
     * 是否允许选中
     */
    enableSelect?: EnableSelect;

    /**
     * 是否允许多选，仅当 enableSelect={true} 时有效
     */
    enableMultiSelect?: EnableMultiSelect;

    /**
     * 显示某一行的额外信息
     */
    showRowExtraOf: RowExtraMatcher;

    /**
     * 选中项，selection属性必须是data的子集或者其中一项
     *
     * @example
     * ```jsx
     * // 不会选中，因为selection和data中的第一项不指向同一对象
     * const data = [{name: 1}]
     * <DataTable data={data} selection={{name: 1}} />
     *
     * // 会选中，因为selection === data
     * const data = [{name: 1}]
     * <DataTable data={data} selection={data[0]} />
     * ```
     */
    selection?: Selection;

    /**
     * 鼠标悬浮在某一行时触发
     */
    onRowEnter?: (record: any, index: number) => void;

    /**
     * 鼠标离开某一行时触发
     */
    onRowLeave?: (record: any, index: number) => void;

    /**
     * 渲染行额外内容
     */
    RowExtraComponent?: RowExtraComponent;

    /**
     * 工具栏组件
     */
    ToolbarComponent?: DataTableToolbarComponent;

    /**
     * 列表头组件
     */
    HeaderComponent?: DataTableHeaderComponent;

    /**
     * 列表页脚组件
     */
    FooterComponent?: DataTableFooterComponent;

    /**
     * 数据为空时显示的组件
     */
    EmptyComponent?: DataTableEmptyComponent;

    /**
     * 数据为空时显示的组件
     */
    RefreshingComponent?: DataTableRefreshingComponent;

    /**
     * 选中项发生变化时触发
     */
    onSelectionChange?: SelectionChangeHandler;

    /**
     * 双击行触发
     */
    onRowDoubleClicked?: RowDoubleClickedHandler;

    /**
     * 当列表滚动接近底部时触发
     */
    onEndReached?: () => void;

    /**
     * 当拖拽结束后触发
     */
    onDragEnd?: (dragIndex: number, hoverIndex: number) => void;

    /**
     * 是否允许展开(树形展开) 与RowExtraComponent不同
     */
    expandable?: boolean;

    /**
     * 每一行唯一的key名，一般为data数组中的唯一标识，默认为'id'
     */
    rowKeyName?: string;

    /**
     * 展开的children的名字，默认为'children'
     */
    childrenColumnName?: string;

    /**
     * 展开行的key，rowKeyName对应的value
     */
    expandedKeys?: ReadonlyArray<string>;

    /**
     * 展开/收起行
     */
    onExpand?: ExpandChangeHandler;

    /**
     * 是否支持伸缩列
     */
    isResizable?: boolean;

    /**
     * 列宽是否固定（不由table自动分配）
     */
    isColsFixed: boolean;

    /**
     * contentViewClassName
     */
    contentViewClassName?: string;
}

export interface DataTableState {
    selection: Selection;

    layout: {
        header: HTMLElement | null;
        toolbar: HTMLElement | null;
        footer: HTMLElement | null;
    };

    overflow?: boolean;

    expandedKeys: ReadonlyArray<string>;

    /**
     * 可视区域数据
     */
    viewData: ReadonlyArray<any>;

    /**
     * table Y轴移动的距离
     */
    translateY: number;
}

/**
 * 默认的渲染单元格方法
 * @param value 值
 * @return 返回输入值
 */
const renderCellDefault = (value: any) => value;

/**
 * 行高数值
 */
const RowHeight = 50;

export default class DataTable extends React.PureComponent<DataTableProps, DataTableState> {
    static defaultProps = {
        rowKeyName: 'id',
        childrenColumnName: 'children',
        isResizable: false,
        isColsFixed: false,
        rowHeight: RowHeight,
    }

    state: DataTableState = {
        selection: this.props.selection || (this.props.enableMultiSelect ? [] : null),
        layout: {
            header: null,
            toolbar: null,
            footer: null,
        },
        expandedKeys: this.props.expandedKeys || [],
        viewData: this.props.data,
        translateY: 0,
    };

    toolbar: HTMLElement | null = null;

    header: HTMLElement | null = null;

    footer: HTMLElement | null = null;

    /**
     * shift多选时，会根据此下标位置计算选中范围
     */
    lastRowClickIndex = 0;

    /**
     * DataTable内容区域
     */
    contentView: HTMLDivElement | null = null;

    /**
     * 本次渲染，当前列表是否触发过滚动到最后
     */
    hasEndReached: boolean = false;

    /**
     * 可视区域数据开始的index值
     */
    startIndex: number = 0;

    /**
     * 可视区域可容纳的data个数
     */
    viewDataLength: number = 0;

    /**
     * contentView变化监视器
     */
    contentViewObserver: any = null;

    render() {
        const {
            role,
            data = [],
            height = 'auto',
            showBorder,
            refreshing,
            rowHeight = RowHeight,
            rowHoverClassName,
            enableSelect,
            enableMultiSelect,
            ToolbarComponent,
            HeaderComponent,
            FooterComponent,
            EmptyComponent,
            RefreshingComponent,
            onDragEnd,
            expandable,
            childrenColumnName,
            rowKeyName,
            isResizable,
            isColsFixed,
            contentViewClassName,
        } = this.props;

        const { viewData, selection, translateY, overflow, layout: { toolbar, header, footer } } = this.state;
        const [firstColumn, ...restColumns] = this.props.columns;
        const handleCheckBoxClicked = this.handleCheckBoxClicked;
        const handleCheckBoxValueChange = this.handleCheckBoxValueChange;
        const hasNestChildren = this.hasNestChildren
        const handleExpandIconClicked = this.handleExpandIconClicked
        const columns = enableSelect && enableMultiSelect && Array.isArray(selection) ?
            // 多选模式下，在第一个单元格内插入复选框
            [
                {
                    ...firstColumn,
                    renderCell: (value: any, record: any, index: number, indent: number, expanded: boolean) => (
                        <View className={styles['multi-select-cell']}>
                            <View
                                className={styles['checkbox-area']}
                                onDoubleClick={this.handleCheckBoxDoubleClick}
                            >
                                <CheckBox
                                    className={styles['checkbox']}
                                    checked={includes(selection, record)}
                                    onClick={(event) => handleCheckBoxClicked(event)}
                                    onChange={(event) => handleCheckBoxValueChange(record, event)}
                                />
                            </View>
                            <View className={styles['cell-content']}>
                                {(firstColumn.renderCell || renderCellDefault)(value, record, index, indent, expanded)}
                            </View>
                        </View>
                    ),
                },
                ...restColumns,
            ] :
            expandable ?
                (
                    [
                        {
                            ...firstColumn,
                            renderCell: (value: any, record: any, index: number, indent: number, expanded: boolean) => (
                                <View className={styles['expandable-select-cell']}>
                                    <View>
                                        <span
                                            style={{ paddingLeft: `${15 * indent}px` }}
                                            className={styles['indent']}
                                        />
                                        <span className={styles['expand-icon']}>
                                            {
                                                hasNestChildren(record) ?
                                                    <SweetIcon
                                                        name={expanded ? 'arrowDown' : 'arrowRight'}
                                                        size={16}
                                                        onClick={(event) => handleExpandIconClicked(record, !expanded, event)}
                                                    />
                                                    : null
                                            }
                                        </span>
                                    </View>
                                    <View className={styles['cell-content']}>
                                        {(firstColumn.renderCell || renderCellDefault)(value, record, index, indent, expanded)}
                                    </View>
                                </View>
                            ),
                        },
                        ...restColumns,
                    ]
                ) :
                this.props.columns;

        // 固定高度布局时，使用绝对定位，使头/脚位置固定
        const fixedLayout = height !== 'auto';

        return (
            <div
                role={role}
                className={classnames(styles['root'], { [styles['show-border']]: showBorder })}
                style={{ height }}
            >
                {
                    ToolbarComponent ? (
                        <div
                            ref={(node) => (this.toolbar = node)}
                            className={classnames({ [styles['fixed-layout']]: fixedLayout })}
                            style={{
                                top: 0,
                            }}
                        >
                            {
                                isFunction(ToolbarComponent) ?
                                    (
                                        <ToolbarComponent
                                            data={data}
                                            selection={selection}
                                            selectAll={enableSelect && enableMultiSelect ? this.selectAll : undefined}
                                        />
                                    ) :
                                    ToolbarComponent
                            }
                        </div>
                    ) : null
                }
                {
                    HeaderComponent ? (
                        <div
                            className={styles['header-warraper']}
                            style={{
                                top: fixedLayout && toolbar !== null ? toolbar.clientHeight : 0,
                            }}
                        >
                            <div
                                ref={(node) => (this.header = node)}
                                className={classnames(
                                    {
                                        [styles['content-overflow']]: overflow,
                                    },
                                    styles['header'],
                                )}
                            >
                                {
                                    isFunction(HeaderComponent) ?
                                        (
                                            <HeaderComponent
                                                columns={this.props.columns}
                                                data={data}
                                                selection={selection}
                                                isResizable={isResizable}
                                                selectAll={enableSelect && enableMultiSelect ? this.selectAll : undefined}
                                            />
                                        ) :
                                        (
                                            HeaderComponent
                                        )
                                }
                            </div>
                        </div>
                    ) : null
                }
                <ContentView
                    ref={this.saveContentView}
                    onScroll={this.handleScrollLazily.bind(this)}
                    onOverflow={this.setOverflow.bind(this, true)}
                    onUnderflow={this.setOverflow.bind(this, false)}
                    className={classnames(
                        {
                            [styles['fixed-layout']]: fixedLayout,
                        },
                        contentViewClassName,
                    )}
                    style={{
                        top: fixedLayout ? 0 : void 0,
                        bottom: fixedLayout ? 0 : void 0,
                        overflow: 'auto',
                        background: '#fff',
                        marginTop: fixedLayout ?
                            (toolbar !== null ? toolbar.clientHeight : 0) + (header !== null ? header.clientHeight : 0)
                            : 0,
                        marginBottom: fixedLayout ? (footer !== null ? footer.clientHeight : 0) : 0,
                    }}
                >
                    {
                        Array.isArray(data) && data.length ? (
                            <BaseTable
                                className={styles['datatable']}
                                data={data}
                                columns={columns}
                                isResizable={isResizable}
                                isColsFixed={isColsFixed}
                                viewData={viewData}
                                translateY={translateY}
                                showRowExtraOf={this.props.showRowExtraOf}
                                RowExtraComponent={this.props.RowExtraComponent}
                                rowExtraClassName={styles['row-extra']}
                                rowExtraKeyExtractor={this.props.rowExtraKeyExtractor}
                                cellKeyExtractor={this.props.cellKeyExtractor}
                                rowKeyExtractor={this.props.rowKeyExtractor}
                                rowClassName={(record, index) =>
                                    (enableMultiSelect && Array.isArray(selection) ? includes(selection, record) : selection === record) ?
                                        classnames(styles['row-selected'], styles['row'], this.props.rowClassName) :
                                        classnames(styles['row'], this.props.rowClassName)
                                }
                                rowHoverClassName={classnames(
                                    styles['row-hovered'],
                                    rowHoverClassName,
                                )}
                                cellClassName={styles['cell']}
                                onRowClicked={this.handleRowClicked}
                                onRowDoubleClicked={this.handleRowDoubleClicked}
                                onRowEnter={this.onRowEnter}
                                onRowLeave={this.onRowLeave}
                                onDragEnd={onDragEnd}
                                expandable={expandable}
                                expandedKeys={this.state.expandedKeys}
                                rowKeyName={rowKeyName}
                                childrenColumnName={childrenColumnName}
                                onExpandIconClicked={this.handleExpandIconClicked}
                            />
                        ) : null
                    }
                </ContentView>
                {
                    !refreshing && EmptyComponent && (!data || !data.length) ? (
                        <View
                            className={classnames({ [styles['fixed-layout']]: fixedLayout })}
                            style={{
                                top: fixedLayout ? 0 : void 0,
                                bottom: fixedLayout ? 0 : void 0,
                                overflow: 'auto',
                                marginTop: fixedLayout ?
                                    (toolbar !== null ? toolbar.clientHeight : 0) + (header !== null ? header.clientHeight : 0) :
                                    0,
                                marginBottom: fixedLayout ? (footer !== null ? footer.clientHeight : 0) : 0,
                            }}
                        >
                            {
                                isFunction(EmptyComponent) ?
                                    <EmptyComponent /> :
                                    EmptyComponent
                            }
                        </View>
                    ) : null
                }
                {
                    refreshing && RefreshingComponent ? (
                        <View
                            className={classnames({ [styles['fixed-layout']]: fixedLayout })}
                            style={{
                                top: fixedLayout ? 0 : void 0,
                                bottom: fixedLayout ? 0 : void 0,
                                overflow: 'auto',
                                marginTop: fixedLayout
                                    ? (toolbar !== null ? toolbar.clientHeight : 0) +
                                    (header !== null ? header.clientHeight : 0)
                                    : 0,
                                marginBottom: fixedLayout ? (footer !== null ? footer.clientHeight : 0) : 0,
                            }}
                        >
                            {
                                isFunction(RefreshingComponent) ?
                                    <RefreshingComponent /> :
                                    RefreshingComponent
                            }
                        </View>
                    ) : null
                }
                {
                    FooterComponent ? (
                        <div
                            ref={(node) => (this.footer = node)}
                            className={classnames({ [styles['fixed-layout']]: fixedLayout })}
                            style={{
                                bottom: fixedLayout ? 0 : 'auto',
                            }}
                        >
                            {
                                isFunction(FooterComponent) ?
                                    <FooterComponent data={data} /> :
                                    FooterComponent
                            }
                        </div>
                    ) : null
                }
            </div>
        )
    }

    static getDerivedStateFromProps({ enableSelect, selection, data }: DataTableProps, prevState: DataTableState) {
        if (enableSelect) {
            if (selection !== undefined && selection !== prevState.selection) {
                return {
                    selection,
                }
            }
        }

        return null
    }

    componentDidMount() {
        this.setState({
            layout: {
                toolbar: this.toolbar,
                header: this.header,
                footer: this.footer,
            },
        });
    }

    componentWillUnmount() {
        this.contentViewObserver && this.contentViewObserver.disconnect()
    }

    componentDidUpdate(prevProps: DataTableProps, prevState: DataTableState) {
        const { onSelectionChange, data, expandedKeys } = this.props

        if (this.state.selection !== prevState.selection) {
            createEventDispatcher(onSelectionChange)(this.state.selection)
        }

        if (data !== prevProps.data) {
            this.setState({
                viewData: data,
            })

            this.hasEndReached = false;
        }

        if (expandedKeys && expandedKeys !== prevProps.expandedKeys) {
            this.setState({
                expandedKeys,
            })
        }
    }

    saveContentView = (contentView: any) => {
        this.contentView = contentView && contentView.container
    }

    /**
     * 处理容器滚动事件
     * @param event UIEvent事件对象
     */
    private handleScroll = (event: SweetUIEvent<MouseWheelEvent>) => {
        if (this.contentView) {
            const { scrollTop, clientHeight, scrollHeight, scrollLeft } = this.contentView

            if (this.header && this.header.scrollLeft !== scrollLeft) {
                this.header.scrollLeft = scrollLeft
            } else {
                const { rowHeight = RowHeight } = this.props

                if (scrollHeight - (scrollTop + clientHeight) <= rowHeight) {
                    if (!this.hasEndReached) {
                        const { onEndReached } = this.props

                        if (isFunction(onEndReached)) {
                            this.hasEndReached = true;
                            this.dispatchEndReachedEvent()
                        }
                    }
                }
            }
        }
    }

    dispatchEndReachedEvent = createEventDispatcher(this.props.onEndReached)

    /**
     * 容器滚动事件懒触发
     */
    private handleScrollLazily = throttle(this.handleScroll, 0, { trailing: true, leading: false });

    /**
     * 触发selectionChange事件
     */
    private dispatchSelectionChangeEvent = createEventDispatcher(this.props.onSelectionChange, ({ detail: selection }) => {
        this.setState({ selection });
    });

    /**
     * 触发双击事件
     */
    private dispatchRowDoubleClickedEvent = createEventDispatcher(this.props.onRowDoubleClicked)

    /**
     * 点击CheckBox时，阻止触发上层<tr>的选中事件，避免复选操作被影响
     */
    private handleCheckBoxClicked = (event: React.MouseEvent<HTMLInputElement>) => {
        event.stopPropagation();
    }

    /**
     * 复选框变化时执行复选操作
     */
    private handleCheckBoxValueChange = (record: any, event: React.ChangeEvent<HTMLInputElement>) => {
        const { selection = [] } = this.state;

        if (Array.isArray(selection)) {
            const nextSelection = event.target.checked === true ? [...selection, record] : without(selection, record);

            this.dispatchSelectionChangeEvent(nextSelection);
        }
    }

    /**
     * preventDefault 复选框双击事件
     * @param e 事件对象
     */
    private handleCheckBoxDoubleClick(e: React.SyntheticEvent<HTMLDivElement>) {
        e.preventDefault()
    }

    /**
     * 触发双击
     */
    private handleRowDoubleClicked = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        this.dispatchRowDoubleClickedEvent({ record, index })
    }

    /**
     * 触发移入行
     */
    private onRowEnter = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        this.dispatchonRowEnterEvent({ record, index })
    }

    /**
     * 触发移出行
     */
    private onRowLeave = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        this.dispatchonRowLeaveEvent({ record, index })
    }

    /**
     * 触发移出行事件
     */
    private dispatchonRowEnterEvent = createEventDispatcher(this.props.onRowEnter)

    /**
     * 触发移出行事件
     */
    private dispatchonRowLeaveEvent = createEventDispatcher(this.props.onRowLeave)

    /**
     * 点击展开按钮触发
     */
    private dispatchonExpandIconClickEvent = createEventDispatcher(this.props.onExpand)

    /**
     * 点中一行时触发
     */
    private handleRowClicked = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        const { enableSelect } = this.props

        if (enableSelect) {
            const { enableMultiSelect, data = [] } = this.props
            const { selection = [] } = this.state
            let nextSelection

            if (enableMultiSelect && Array.isArray(selection)) {
                if (event.ctrlKey) {
                    nextSelection = includes(selection, record) ? without(selection, record) : [...selection, record]

                    this.lastRowClickIndex = index
                } else if (event.shiftKey) {
                    if (this.lastRowClickIndex < index) {
                        nextSelection = data.slice(this.lastRowClickIndex, index + 1)
                    } else if (this.lastRowClickIndex > index) {
                        nextSelection = data.slice(index, this.lastRowClickIndex + 1)
                    } else {
                        nextSelection = [data[this.lastRowClickIndex]]
                    }
                } else {
                    nextSelection = includes(selection, record) ? (selection.length > 1 ? [record] : []) : [record]

                    this.lastRowClickIndex = index
                }
            } else {
                nextSelection = record === selection ? null : record
            }

            this.dispatchSelectionChangeEvent(nextSelection)
        }
    }

    /**
     * 选中所有行
     */
    private selectAll = (selected: boolean) => {
        const nextSelection = selected ? this.props.data || [] : [];

        this.dispatchSelectionChangeEvent(nextSelection);
    };

    /**
     * 设置内容溢出状态
     */
    private setOverflow = (overflow: boolean) => {
        this.setState({ overflow })
    }

    /**
     * 点击展开图标
     */
    private handleExpandIconClicked = (record: any, expanded: boolean, event: React.MouseEvent<HTMLTableRowElement>) => {
        event.stopPropagation()

        const expandedKey = this.props.rowKeyName && record[this.props.rowKeyName]

        if (expandedKey) {
            let newExpandedKeys = [...this.state.expandedKeys]

            if (expanded) {
                newExpandedKeys = [...newExpandedKeys, expandedKey]
            } else {
                newExpandedKeys = newExpandedKeys.filter((key) => key !== expandedKey)
            }

            this.setState({
                expandedKeys: newExpandedKeys,
            })

            this.dispatchonExpandIconClickEvent({ expandedKeys: newExpandedKeys, record, expanded })
        }
    }

    /**
     * 是否有子节点
     */
    private hasNestChildren = (record: any): boolean => {
        const { childrenColumnName } = this.props

        if (childrenColumnName && record && record[childrenColumnName]) {
            if (Array.isArray(record[childrenColumnName])) {
                return !!record[childrenColumnName].length
            }

            return true
        }

        return false
    }
}
