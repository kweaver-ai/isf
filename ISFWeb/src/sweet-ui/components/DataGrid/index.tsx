import React from 'react';
import { ResizeObserver } from '@juggle/resize-observer';
import { isFunction, omit, isEqual, pick } from 'lodash';
import { clamp } from '@/util/formatters'
import { isDom } from '@/util/validators'
import { bindEvent, unbindEvent } from '@/util/browser';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import {
    DataTableColumn,
    Selection,
    EnableSelect,
    EnableMultiSelect,
    RowClassName,
    RowHoverClassName,
    RowKeyExtractor,
    CellKeyExtractor,
    DataTableEmptyComponent,
    DataTableRefreshingComponent,
    DataTableHeaderComponent,
    SelectionChangeHandler,
    RowDoubleClickedHandler,
    DataTableToolbarComponent,
    DataTableFooterComponent,
    RowExtraComponent,
    RowExtraMatcher,
    ExpandChangeHandler,
} from '../DataTable';
import DataTable from '../DataTable';
import DataGridHeader, { SortProps, RequestSortHandler } from './DataGridHeader';
import DataGridToolbar from './DataGridToolbar';
import DataGridFooter, { PageChangeHandler } from './DataGridFooter';
import { PageChangeEvent } from '../Pager';

type ToolbarComponent = React.FunctionComponent<{ data: ReadonlyArray<any>; selection: Selection }>;

type ParamsChangeHandler = (params: Params) => void;

type SelectAllChangeHandler = (event: SweetUIEvent<boolean>) => void;

type LazyLoadHandler = (event: SweetUIEvent<{ start: number; limit: number }>) => void;

type FilterHandler = (event: SweetUIEvent<Array<{ [key: string]: Array<any> }>>) => void;

type DragReorderHandler = (event: SweetUIEvent<Array<any>>) => void;

type onDragItemHandler = (event: SweetUIEvent<any>) => void;

export interface DataGridProps extends React.ClassAttributes<any> {
    /**
     * 角色
     */
    role?: string;

    /**
     * 列配置
     */
    columns: ReadonlyArray<DataTableColumn>;

    /**
     * 表格数据
     */
    data: ReadonlyArray<any>;

    /**
     * 懒加载数据的起始索引
     */
    start?: number;

    /**
     * 懒加载每次加载的最大数据条数
     */
    limit?: number;

    /**
     * 是否显示表头
     */
    headless?: boolean;

    /**
     * 是否显示表格的边框
     */
    showBorder?: boolean;

    /**
     * 是否正在刷新
     */
    refreshing?: boolean;

    /**
     * 表格整体高度，含表头和分页
     */
    height?: number | string;

    /**
     * 选中项
     */
    selection?: Selection;

    /**
     * 行高
     */
    rowHeight?: number;

    /**
     * 行key扩展
     */
    rowKeyExtractor?: RowKeyExtractor;

    /**
     * 行额外信息key扩展
     */
    rowExtraKeyExtractor?: RowKeyExtractor;

    /**
     * 单元格key扩展
     */
    cellKeyExtractor?: CellKeyExtractor;

    /**
     * 列表行的类名
     */
    rowClassName?: RowClassName;

    /**
     * 行悬浮样式
     */
    rowHoverClassName?: RowHoverClassName;

    /**
     * 是否允许选中
     */
    enableSelect?: EnableSelect;

    /**
     * 是否允许多选，只有 enableSelect: true 时有效
     */
    enableMultiSelect?: EnableMultiSelect;

    /**
     * 显示某一行的额外信息
     */
    showRowExtraOf?: RowExtraMatcher;

    /**
     * 渲染行额外内容
     */
    RowExtraComponent?: RowExtraComponent;

    /**
     * 工具栏内容
     */
    ToolbarComponent?: ToolbarComponent;

    /**
     * 正在刷新时的状态组件
     */
    RefreshingComponent?: DataTableRefreshingComponent;

    /**
     * 列表为空时的显示内容
     */
    EmptyComponent?: DataTableEmptyComponent;

    /**
     * 选中项发生变化时触发
     */
    onSelectionChange?: SelectionChangeHandler;

    /**
     * 双击行触发
     */
    onRowDoubleClicked?: RowDoubleClickedHandler;

    /**
     * 分页配置，如果为空则不显示分页
     */
    DataGridPager?: DataGridPagerProps;

    /**
     * 表头配置
     */
    DataGridHeader?: DataGridHeaderProps;

    /**
     * 工具栏配置
     */
    DataGridToolbar?: DataGridToolbarProps;

    /**
     * 排序参数
     */
    sort?: SortProps;

    /**
     * 排序处理
     */
    onRequestSort?: RequestSortHandler;

    /**
     * 参数改变(同时传递分页和排序参数)
     */
    onParamsChange?: ParamsChangeHandler;

    /**
     * 全选状态改变
     */
    onSelectAllChange?: SelectAllChangeHandler;

    /**
     * 响应懒加载事件
     */
    onRequestLazyLoad?: LazyLoadHandler;

    /**
     * 过滤处理
     */
    onRequestFilter?: FilterHandler;

    /**
     * 获取过滤菜单项key值
     */
    getFilterKey?: (filter: any) => string;

    /**
     * 拖拽项
     */
    onDragItem?: onDragItemHandler;

    /**
     * 拖拽排序结束后的回调
     */
    onDragReorder?: DragReorderHandler;

    /**
     * 是否可展开
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
     * 展开/收起子行
     */
    onExpand?: ExpandChangeHandler;

    /**
     * 是否支持伸缩列
     */
    isResizable: boolean;

    /**
     * contentViewClassName
     */
    contentViewClassName?: string;
}

interface DataGridHeaderProps {
    /**
     * 是否显示【全选】工具
     */
    enableSelectAll?: boolean;
}

interface DataGridToolbarProps {
    /**
     * 是否在工具栏显示【全选】工具
     */
    enableSelectAll?: boolean;
}

interface DataGridPagerProps {
    /**
     * 页码
     */
    page: number;

    /**
     * 分页大小
     */
    size: number;

    /**
     * 页码发生变化时触发
     */
    onPageChange: PageChangeHandler;

    /**
     * 数据总条目数
     */
    total: number;
}

export interface DataGridState {
    /**
     * 选中项
     */
    selection?: Selection;

    /**
     * 页码
     */
    page?: number;

    /**
     * 排序参数
     */
    sort?: SortProps;

    /**
     * 列配置
     */
    columns: ReadonlyArray<DataTableColumn>;

    /**
     * 列宽是否固定（不由table自动分配）
     */
    isColsFixed: boolean;

    /**
     * 列宽是否移动中
     */
    isMoving: boolean;
}

/**
 * 参数传递(传递分页和排序参数)
 */
export interface Params {
    /**
     * 页码
     */
    page?: number;

    /**
     * 排序参数
     */
    sort?: SortProps;
}

/**
 * 默认开始页
 */
const DEFAULT_PAGE = 1;

/**
 * 默认列最小宽度
 */
const ColumnMinWidth = 50;

/**
 * 默认列最大宽度
 */
const ColumnMaxWidth = 10000;

/**
 * 数组重新排序
 * @param arr 数据
 * @param fromIndex 原始索引
 * @param toIndex 排序后的索引
 */
function reorder(arr: Array<any>, fromIndex: number, toIndex: number) {
    // eslint-disable-next-line no-restricted-properties
    let item = arr.splice(fromIndex, 1)[0];
    // eslint-disable-next-line no-restricted-properties
    arr.splice(toIndex, 0, item);
    return arr;
}

export default class DataGrid extends React.PureComponent<DataGridProps, DataGridState> {
    static defaultProps = {
        isResizable: true,
    }

    /**
     * 当前表格是否是全选状态
     */
    isSelectAll: boolean = false;

    state = {
        selection: this.props.selection,
        page: this.props.DataGridPager ? this.props.DataGridPager.page : DEFAULT_PAGE,
        sort: this.props.sort,
        columns: this.props.columns,
        isColsFixed: false,
        isMoving: false,
    };

    /**
     * 当前选择的列过滤项
     */
    filterObject: { [key: string]: Array<any> } = {};

    /**
     * 正在伸缩的列
     */
    resizingColumn: {
        /**
         * 列配置项
         */
        column: DataTableColumn;

        /**
         * 开始伸缩前鼠标的位置
         */
        mouseStartX: number;

        /**
         * 列的index值
         */
        index: number;
    } = {}

    /**
     * DataTable的ref
     */
    tableRoot: any = null

    /**
     * DataGridHeader的ref
     */
    gridHeader: any = null

    /**
     * header容器监视器
     */
    resizeObserver: any = null

    componentDidMount() {

        if (this.gridHeader && this.props.isResizable) {
            this.setState({
                columns: this.gridHeader.getCols(),
                isColsFixed: true,
            })

            let headerWidth = 0

            this.resizeObserver = new ResizeObserver((entries, observer) => {
                const inlineSize = entries[0].borderBoxSize[0].inlineSize

                if (inlineSize !== headerWidth) {
                    const tableWidth = (entries[0].target.firstChild as HTMLElement).offsetWidth

                    if (tableWidth !== inlineSize) {
                        if (tableWidth < inlineSize || !headerWidth) {
                            this.setState({
                                columns: this.state.columns.map((col) => {
                                    const { width, minWidth = ColumnMinWidth, maxWidth = ColumnMaxWidth } = col

                                    return {
                                        ...col,
                                        width: clamp((Number(width) / tableWidth) * inlineSize, minWidth, maxWidth),
                                    }
                                }),
                                isColsFixed: true,
                            })
                        } else {
                            const race = inlineSize / headerWidth

                            this.setState({
                                columns: this.state.columns.map((col) => {
                                    const { width, minWidth = ColumnMinWidth, maxWidth = ColumnMaxWidth } = col

                                    return {
                                        ...col,
                                        width: clamp((Number(width) * race), minWidth, maxWidth),
                                    }
                                }),
                                isColsFixed: true,
                            })
                        }
                    }
                }

                headerWidth = inlineSize
            })

            this.tableRoot.header && isDom(this.tableRoot.header) && this.resizeObserver.observe(this.tableRoot.header)
        }
    }

    componentWillUnmount() {
        this.resizeObserver && this.resizeObserver.disconnect()
    }

    componentDidUpdate(prevProps: DataGridProps, prevState: DataGridState) {
        const { selection, data, DataGridPager, columns, isResizable, headless } = this.props;
        const isMultiSelect = prevProps.enableSelect && prevProps.enableMultiSelect;
        const nextSelection = selection || prevState.selection;

        // 当数据发生变化
        if (data !== undefined && data !== prevProps.data) {
            if (Array.isArray(data) && data.length > 0) {
                // 多选情况下
                // 如果数据发生变化，且当前是全选状态，则新插入的数据也应当是选中状态
                if (isMultiSelect) {
                    if (this.isSelectAll) {
                        this.dispatchSelectionChangeEvent(data);
                    } else {
                        if (nextSelection !== prevState.selection) {
                            this.dispatchSelectionChangeEvent(nextSelection);
                        }
                    }
                } else {
                    // 单选情况下
                    if (nextSelection && !data.includes(nextSelection)) {
                        this.dispatchSelectionChangeEvent(null);
                    } else {
                        if (nextSelection !== prevState.selection) {
                            this.dispatchSelectionChangeEvent(nextSelection);
                        }
                    }
                }
            } else {
                // 当数据被清空
                this.dispatchSelectionChangeEvent(isMultiSelect ? [] : null);
            }
        } else {
            // 当选中项发生变化
            if (selection !== undefined && nextSelection !== prevState.selection) {
                this.dispatchSelectionChangeEvent(nextSelection);
            }
        }

        if (columns !== undefined && (
            columns.length !== prevProps.columns.length ||
            columns !== prevProps.columns
        )) {
            if (this.gridHeader && isResizable) {
                const headerCols = this.gridHeader.getCols()

                if (
                    columns.length !== prevProps.columns.length ||
                    columns.some((col) => !isEqual(
                        omit(col, 'renderCell', 'filters', 'sort', 'title'),
                        omit(prevProps.columns.find((h) => h.key === col.key), 'renderCell', 'filters', 'sort', 'title'),
                    )) ||
                    this.state.columns.some((col) => !isEqual(
                        pick(col, 'minWidth'),
                        pick(headerCols.find((h) => h.key === col.key), 'minWidth'),
                    ))
                ) {
                    this.setState({
                        isColsFixed: false,
                        columns,
                    }, () => {
                        const hCols = this.gridHeader.getCols()
                        this.setState({
                            columns: columns.map((col) => {
                                const { width, minWidth } = (hCols.find((sCol) => sCol.key === col.key)) as DataTableColumn

                                return (
                                    {
                                        minWidth,
                                        ...col,
                                        width,
                                    }
                                )
                            }),
                            isColsFixed: true,
                        })
                    })
                } else {
                    this.setState({
                        columns: columns.map((col) => {
                            const { width, minWidth } = (this.state.columns.find((sCol) => sCol.key === col.key)) as DataTableColumn

                            return (
                                {
                                    minWidth,
                                    ...col,
                                    width,
                                }
                            )
                        }),
                        isColsFixed: true,
                    })
                }
            }
            else {
                this.setState({
                    columns,
                })
            }
        }

        if (DataGridPager && DataGridPager.page && DataGridPager.page !== prevState.page) {
            this.setState({
                page: DataGridPager.page,
            })
        }
    }

    /**
     * 触发选项改变事件
     */
    private dispatchSelectionChangeEvent = createEventDispatcher(
        this.props.onSelectionChange,
        ({ detail: selection }) => {
            this.setState({ selection }, () => {
                const willSelectAll = this.willSelectAll(selection);

                if (willSelectAll !== this.isSelectAll) {
                    this.dispatchSelectAllChangeEvent(willSelectAll);
                }
            });
        },
    );

    /**
     * 触发双击行事件
     */
    private dispatchRowDoubleClickedEvent = createEventDispatcher(this.props.onRowDoubleClicked);

    /**
     * 鼠标移入行时触发事件
     */
    private dispatchonRowEnterEvent = createEventDispatcher(this.props.onRowEnter);

    /**
     * 鼠标移出行时触发事件
     */
    private dispatchonRowLeaveEvent = createEventDispatcher(this.props.onRowLeave);

    /**
     * 触发全选状态变化事件
     */
    private dispatchSelectAllChangeEvent = createEventDispatcher(
        this.props.onSelectAllChange,
        ({ detail: isSelectAll }) => {
            this.isSelectAll = isSelectAll;
        },
    );

    /**
     * 触发懒加载事件
     */
    private dispatchLazyLoadEvent = createEventDispatcher(this.props.onRequestLazyLoad);

    /**
     * 触发展开事件
     */
    private dispatchExpandIconClickedEvent = createEventDispatcher(this.props.onExpand);

    /**
     * 判断当前是否选中所有数据
     */
    private willSelectAll = (nextSelection = []): boolean => {
        const { enableSelect, enableMultiSelect } = this.props;

        if (enableSelect) {
            if (enableMultiSelect) {
                return (
                    Array.isArray(this.props.data) &&
                    this.props.data.length > 0 &&
                    nextSelection.length === this.props.data.length
                );
            } else {
                return false;
            }
        } else {
            return false;
        }
    };

    /**
     * 改变组件内部状态，并通过回调方法将改变后的状态抛出
     * @param param0 接收page和sort
     * @param callback 回调函数
     */
    public changeParams({ page, sort }: Params = {}, callback?: Function) {
        this.setState(
            {
                page: page || this.state.page,
                sort: sort || this.state.sort,
            },
            () => {
                if (isFunction(this.props.onParamsChange)) {
                    this.props.onParamsChange({ page, sort });
                }
                isFunction(callback) && callback({ page: this.state.page, sort: this.state.sort });
            },
        );
    }

    /**
     * 排序事件处理
     * @param sort 排序参数
     */
    private requestSort(sort: SortProps) {
        this.changeParams(
            {
                sort,
            },
            ({ sort }: Params) => {
                if (isFunction(this.props.onRequestSort)) {
                    sort && this.props.onRequestSort(sort);
                }
            },
        );
    }

    /**
     * 过滤事件处理
     * @param filterItem
     */
    private requestFilter(filterItem: { key: string; filters: Array<any> }) {
        const { key, filters } = filterItem;
        const newOptions = omit(this.filterObject, key);

        if (filters.length > 0) {
            this.filterObject = { ...newOptions, [key]: filters };
        } else {
            this.filterObject = newOptions;
        }

        this.dispatchRequestFilterEvent(this.filterObject);
    }

    dispatchRequestFilterEvent = createEventDispatcher(this.props.onRequestFilter);

    /**
     * 分页事件处理
     * @param event 分页事件参数
     */
    private pageChange(event: PageChangeEvent) {
        const page = event.detail.page;
        this.changeParams(
            {
                page,
            },
            ({ page }: Params) => {
                if (this.props.DataGridPager && isFunction(this.props.DataGridPager.onPageChange)) {
                    this.props.DataGridPager.onPageChange({ ...event, detail: { ...event.detail, page } });
                }
            },
        );
    }

    handleLazyLoad = () => {
        const { limit } = this.props;
        if (limit) {
            // 由上层判断是否执行下一次加载
            this.dispatchLazyLoadEvent({ start: this.props.data.length, limit });
        }
    };

    /**
     * 拖拽结束事件的回调
     */
    handleDragEnd = (dragIndex: number, hoverIndex: number) => {
    };

    dispatchonDragItemEvent = createEventDispatcher(this.props.onDragItem)

    dispatchDragReorderEvent = createEventDispatcher(this.props.onDragReorder)

    /**
     * 开始伸缩列
     */
    startResize = (event: React.MouseEvent<HTMLElement>, column: DataTableColumn, index: number) => {
        this.resizingColumn = {
            column,
            index,
            mouseStartX: event.clientX,
        }

        this.setState({
            isMoving: true,
        })

        bindEvent(document, 'mousemove', this.moveResize)
        bindEvent(document, 'mouseup', this.endResize)
    }

    /**
    * 随鼠标移动伸缩列
    * @param event 鼠标移动事件对象
    */
    moveResize = (event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault();
        const { mouseStartX, column: { width, minWidth = ColumnMinWidth, maxWidth = ColumnMaxWidth }, index } = this.resizingColumn

        const currwidth = event.clientX - mouseStartX + Number(width)

        this.setState(({ columns }) => {
            const nextColumns = [...columns]

            nextColumns[index] = {
                ...nextColumns[index],
                width: clamp(currwidth, minWidth, maxWidth),
            }

            return { columns: nextColumns, isColsFixed: true }
        })
    }

    /**
     * 伸缩结束
     */
    endResize = () => {
        this.setState({
            isMoving: false,
        })

        unbindEvent(document, 'mousemove', this.moveResize);
        unbindEvent(document, 'mouseup', this.endResize);
    }

    render() {
        const {
            data,
            headless,
            showBorder,
            enableSelect,
            enableMultiSelect,
            DataGridPager,
            ToolbarComponent,
            EmptyComponent,
            RefreshingComponent,
            getFilterKey,
            expandable,
            rowKeyName,
            childrenColumnName,
            expandedKeys,
            isResizable,
            contentViewClassName,
            filterElement,
        } = this.props;

        const { page, sort, columns, isColsFixed } = this.state;

        const Header: DataTableHeaderComponent = headless
            ? null
            : ({ columns, data, selectAll }) => {
                const { enableSelectAll = false } = this.props.DataGridHeader || {};

                return (
                    <DataGridHeader
                        ref={(gridHeader) => this.gridHeader = gridHeader}
                        enableSelectAll={enableSelect && enableMultiSelect && enableSelectAll}
                        isSelectAllChecked={this.willSelectAll(this.state.selection)}
                        onRequestSort={this.requestSort.bind(this)}
                        onRequestFilter={this.requestFilter.bind(this)}
                        onRequesStartResize={this.startResize}
                        isColsFixed={isColsFixed}
                        isMoving={this.state.isMoving}
                        {...{ columns, data, selectAll, sort, getFilterKey, isResizable, filterElement }}
                    />
                )
            }

        const Tools: DataTableToolbarComponent | null = ToolbarComponent
            ? ({ data, selection, selectAll }) => {
                const { enableSelectAll = false } = this.props.DataGridToolbar || {};

                return (
                    <DataGridToolbar
                        ToolbarComponent={ToolbarComponent}
                        isSelectAllChecked={this.willSelectAll(this.state.selection)}
                        {...{ data, selection, selectAll, enableSelectAll }}
                    />
                );
            }
            : null;

        const Footer: DataTableFooterComponent | null = DataGridPager ? (
            <DataGridFooter
                DataGridPager={{
                    ...DataGridPager,
                    onPageChange: this.pageChange.bind(this),
                    page,
                    total: DataGridPager.total || data.length,
                }}
            />
        ) : null;

        const Table = (role?) => (
            <DataTable
                role={role}
                ref={(tableRoot) => this.tableRoot = tableRoot}
                columns={columns}
                isResizable={isResizable}
                isColsFixed={isColsFixed}
                showBorder={showBorder}
                data={this.props.data}
                rowHeight={this.props.rowHeight}
                refreshing={this.props.refreshing}
                height={this.props.height}
                selection={this.state.selection}
                enableSelect={this.props.enableSelect}
                enableMultiSelect={this.props.enableMultiSelect}
                rowKeyExtractor={this.props.rowKeyExtractor}
                rowClassName={this.props.rowClassName}
                rowExtraKeyExtractor={this.props.rowExtraKeyExtractor}
                rowHoverClassName={this.props.rowHoverClassName}
                cellKeyExtractor={this.props.cellKeyExtractor}
                onSelectionChange={this.dispatchSelectionChangeEvent}
                onRowDoubleClicked={this.dispatchRowDoubleClickedEvent}
                onRowEnter={this.dispatchonRowEnterEvent}
                onRowLeave={this.dispatchonRowLeaveEvent}
                contentViewClassName={contentViewClassName}
                showRowExtraOf={this.props.showRowExtraOf}
                RowExtraComponent={this.props.RowExtraComponent}
                EmptyComponent={EmptyComponent}
                RefreshingComponent={RefreshingComponent}
                HeaderComponent={Header}
                ToolbarComponent={Tools}
                FooterComponent={Footer}
                onEndReached={this.handleLazyLoad}
                onDragEnd={this.handleDragEnd}
                expandable={expandable}
                expandedKeys={expandedKeys}
                rowKeyName={rowKeyName}
                childrenColumnName={childrenColumnName}
                onExpand={this.dispatchExpandIconClickedEvent}
            />
        );

        return Table(this.props.role);
    }
}
