import React from 'react';
import classnames from 'classnames';
import { get, isFunction, isArray } from 'lodash';
import View from '../View';
import TableRow from './TableRow';
import styles from './styles';

export type Record = any;

export type RenderCell = (value: any, record: any, index: number, indent: number, expanded: boolean) => any;

export interface BaseTableColumn {
    /**
     * 宽度
     */
    width?: string | number;

    /**
     * 最小宽度
     */
    minWidth?: number;

    /**
     * 最大宽度
     */
    maxWidth?: number;

    /**
     * 键
     */
    key?: string;

    /**
     * 单列渲染
     */
    renderCell?: RenderCell;
}

export type RowClassNameGenerator = (record: any, index: number) => string;

export type CellClassNameGenerator = (record: any, index: number, key?: string) => string;

export type RowKeyExtractor = (record: any) => string;

export type CellKeyExtractor = (record: any, key?: string) => string;

/**
 * 要显示额外内容的行
 */
export type RowExtraMatcher = number | Record;

/**
 * 渲染行额外内容
 */
export type RowExtraComponent =
    | React.SFC<{ record: Record; index: number }>
    | React.Component<React.ClassAttributes<{ record: Record; index: number }>, any>;

/**
 * 点击行事件处理函数
 */
export type RowClickedHandler = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => void;

/**
 * 双击行事件处理函数
 */
export type RowDoubleClickedHandler = (
    record: any,
    index: number,
    event: React.MouseEvent<HTMLTableRowElement>,
) => void;

/**
 * 点击展开图标
 */
export type ExpandIconClickedHandler = (
    record: any,
    expanded: boolean,
    event: React.MouseEvent<HTMLTableRowElement>,
) => void;

export interface BaseTableProps extends React.ClassAttributes<void> {
    /**
     * 数据列配置
     */
    columns: ReadonlyArray<BaseTableColumn>;

    /**
     * 根元素className
     */
    className?: string;

    /**
     * 表格的数据源
     */
    data?: ReadonlyArray<Record>;

    /**
     * 每一行的className
     */
    rowClassName?: string | RowClassNameGenerator;

    /**
     * 每个单元格的className
     */
    cellClassName?: string | CellClassNameGenerator;

    /**
     * 悬浮在一行上时的className
     */
    rowHoverClassName?: string | RowClassNameGenerator;

    /**
     * 显示某一行的额外信息
     */
    showRowExtraOf?: RowExtraMatcher;

    /**
     * 显示行的额外信息
     */
    RowExtraComponent?: RowExtraComponent;

    /**
     * 行额外信息的外部自定义样式
     */
    rowExtraClassName?: string;

    /**
     * 生成每行对应的唯一key
     */
    rowKeyExtractor?: RowKeyExtractor;

    /**
     * 生成行额外信息对应的唯一key
     */
    rowExtraKeyExtractor?: RowKeyExtractor;

    /**
     * 生成每单元格对应的唯一key
     */
    cellKeyExtractor?: CellKeyExtractor;

    /**
     * 鼠标悬浮在某一行时触发
     */
    onRowEnter?: (record: any, index: number) => void;

    /**
     * 鼠标离开某一行时触发
     */
    onRowLeave?: (record: any, index: number) => void;

    /**
     * 鼠标单击某一行时触发
     */
    onRowClicked?: RowClickedHandler;

    /**
     * 鼠标双击某一行时触发
     */
    onRowDoubleClicked?: RowDoubleClickedHandler;

    /**
     * 拖拽结束事件回调
     */
    onDragEnd?: (dragIndex: number, hoverIndex: number) => void;

    /**
     * table row索引
     */
    index?: number;

    /**
     * table row对象
     */
    record?: any;

    /**
     * 是否可以展开
     */
    expandable?: boolean;

    /**
     * 每一行唯一的key名，一般为data数组中的唯一标识，如'id'
     */
    rowKeyName?: string;

    /**
     * 展开的children的名字，默认为'children'
     */
    childrenColumnName?: string;

    /**
     * 展开行的key
     */
    expandedKeys: ReadonlyArray<string>;

    /**
     * 点击展开图标
     */
    onExpandIconClicked: ExpandIconClickedHandler;

    /**
     * 是否支持伸缩列
     */
    isResizable?: boolean;

    /**
     * 列宽是否固定（不由table自动分配）
     */
    isColsFixed: boolean;

    /**
     * 界面显示的数据
     */
    viewData?: ReadonlyArray<Record>;

    /**
     * table Y轴移动的距离
     */
    translateY: number;
}

/**
 * 默认的渲染单元格方法
 * @param value 单元格的值
 * @param record 整行记录
 * @returns any
 */
const renderCellDefault: RenderCell = (value, record) => value;

export default class BaseTable extends React.PureComponent<BaseTableProps, any> {
    static defaultProps = {
        columns: [],
        data: [],
        isResizable: false,
    };

    /**
     * 点击行时触发
     */
    handleRowClicked = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onRowClicked) && this.props.onRowClicked(record, index, event);
    };

    /**
     * 双击行触发
     */
    handleRowDoubleClicked = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onRowDoubleClicked) && this.props.onRowDoubleClicked(record, index, event);
    };

    /**
     * 鼠标移入行时触发
     */
    handleRowEnter = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onRowEnter) && this.props.onRowEnter(record, index);
    };

    /**
     * 鼠标移出行时触发
     */
    handleRowLeave = (record: any, index: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onRowLeave) && this.props.onRowLeave(record, index);
    };

    /**
     * 点击展开图标
     */
    handleTriggerExpand = (record: any, expanded: boolean, event: React.MouseEvent<HTMLTableRowElement>) => {
        isFunction(this.props.onExpandIconClicked) && this.props.onExpandIconClicked(record, expanded, event);
    }

    render() {
        const {
            className,
            data = [],
            rowKeyExtractor,
            rowExtraKeyExtractor,
            cellKeyExtractor,
            rowHoverClassName,
            rowClassName,
            cellClassName,
            columns,
            showRowExtraOf,
            RowExtraComponent,
            rowExtraClassName,
            onDragEnd,
            rowKeyName = 'id',
            childrenColumnName = 'children',
            expandable,
            expandedKeys,
            isResizable,
            isColsFixed,
            viewData = [],
            translateY = 0,
        } = this.props;

        const BaseTableRow = TableRow;

        const getBaseTableRows = (record: any, index: number, indent: number, expanded: boolean) => {
            return (
                [
                    <BaseTableRow
                        key={isFunction(rowKeyExtractor) ? rowKeyExtractor(record) : `${indent}-${index}`}
                        className={isFunction(rowClassName) ? rowClassName(record, index) : rowClassName}
                        hoverClassName={
                            isFunction(rowHoverClassName) ? rowHoverClassName(record, index) : rowHoverClassName
                        }
                        onClick={(event) => this.handleRowClicked(record, index, event)}
                        onDoubleClick={(event) => this.handleRowDoubleClicked(record, index, event)}
                        onEnter={(event) => this.handleRowEnter(record, index, event)}
                        onLeave={(event) => this.handleRowLeave(record, index, event)}
                        // 以下是拖拽特性用到的方法
                        onDragEnd={onDragEnd}
                        dropOverClassName={styles['drag-over']}
                        index={index}
                        record={record}
                    >
                        {
                            columns.map((column, columIndex) => (
                                <td
                                    key={isFunction(cellKeyExtractor) ? cellKeyExtractor(record, column.key) : `${indent}-${columIndex}`}
                                    className={
                                        isFunction(cellClassName) ? (
                                            cellClassName(record, columIndex, column.key)
                                        ) : (cellClassName)
                                    }
                                >
                                    <View
                                        className={styles['cell-content']}
                                        style={{ width: (isResizable && isColsFixed) ? column.width : 'auto' }}
                                    >
                                        {(column.renderCell || renderCellDefault)(
                                            column.key ? get(record, column.key) : null,
                                            record,
                                            index,
                                            indent,
                                            expanded,
                                        )}
                                    </View>
                                </td>
                            ))
                        }
                    </BaseTableRow>,
                    RowExtraComponent && (showRowExtraOf === index || showRowExtraOf === record) ?
                        (
                            <tr
                                className={rowExtraClassName}
                                key={isFunction(rowExtraKeyExtractor) ? rowExtraKeyExtractor(record) : `extrarow-${index}`}
                            >
                                <td colSpan={columns.length}>
                                    {
                                        <RowExtraComponent
                                            {...{ record, index }}
                                        />
                                    }
                                </td>
                            </tr>
                        ) : null,
                ]
            )
        }

        const getRows = (datas: ReadonlyArray<any>, indent = 0) => {
            return datas.map((record, index) => {

                const expanded = expandedKeys && expandedKeys.includes(record[rowKeyName])

                let rows: ReadonlyArray<any> = [getBaseTableRows(record, index, indent, expanded)]

                if (expandable && expanded) {
                    rows = [...rows, getRows(isArray(record[childrenColumnName]) ? record[childrenColumnName] : [record[childrenColumnName]], indent + 1)]
                }

                return rows
            })
        }

        return (
            <table
                className={classnames(styles['table'], className)}
                style={{ width: (isResizable && isColsFixed) ? 'auto' : '100%', transform: `translate3d(0px, ${translateY}px, 0)` }}
            >
                <colgroup>{columns.map((col, index) => <col key={index} width={col.width} />)}</colgroup>
                <tbody>
                    {getRows(viewData)}
                </tbody>
            </table>
        );
    }
}
