import React from 'react';
import { isFunction, isNumber } from 'lodash';
import classnames from 'classnames';
import UIIcon from '@/ui/UIIcon/ui.desktop';
import { DataTableColumn } from '../../DataTable';
import CheckBox from '../../CheckBox';
import View from '../../View';
import SweetIcon from '../../SweetIcon';
import DataGridFilter from '../DataGridFilter';
import styles from './styles';

/**
 * 排序方式
 */
export enum SortType {
    /**
     * 降序
     */
    DESC,

    /**
     * 升序
     */
    ASC
}

/**
 * 排序处理
 */
export type RequestSortHandler = (sort: SortProps) => void;

/**
 * 排序参数
 */
export interface SortProps {
    /**
     * 排序字段
     */
    key: string;

    /**
     * 排序方式 : 降序 升序
     */
    type: SortType;
}

interface DataGridHeaderProps {
    columns: ReadonlyArray<DataTableColumn>;

    data: ReadonlyArray<any>;

    isSelectAllChecked?: boolean;

    enableSelectAll?: boolean;

    selectAll: (checked: boolean) => void;

    onRequestSort?: RequestSortHandler;

    sort?: SortProps;

    onRequestFilter?: (filter: { key: string; filters: Array<string> }) => void;

    getFilterKey?: (filter: any) => string;

    /**
     * 是否支持伸缩列
     */
    isResizable?: boolean;

    /**
     * 列宽是否固定（不由table自动分配）
     */
    isColsFixed?: boolean;

    /**
     * 开始拖拽改变列宽
     */
    onRequesStartResize?: (event: React.MouseEvent<HTMLElement>, col: DataTableColumn, index: number) => void;

    /**
     * 是否正在拖拽移动
     */
    isMoving?: boolean;
}

export default class DataGridHeader extends React.PureComponent<DataGridHeaderProps, any> {
    static defaultProps = {
        data: [],
        isSelectAllChecked: false,
        enableSelectAll: false,
        isResizable: false,
        isColsFixed: false,
    }

    colsRefs: ReadonlyArray<any> = []

    titleNodes: ReadonlyArray<any> = []

    public getCols = () => {
        const { enableSelectAll, sort } = this.props

        return (
            (this.colsRefs.length && this.titleNodes.length) ?
                this.props.columns.map((col, index) => {
                    const { minWidth: colMinWidth } = col
                    const minWidth = this.titleNodes[index].getBoundingClientRect().width + 20 + ((index === 0 && enableSelectAll) ? 25 : 0)
                    const width = this.colsRefs[index].getBoundingClientRect().width

                    return (
                        {
                            ...col,
                            minWidth: (isNumber(colMinWidth) && minWidth < colMinWidth) ? colMinWidth : minWidth,
                            width: (width < minWidth ? minWidth : width),
                        }
                    )
                })
                : this.props.columns
        )
    }

    render() {
        const {
            data,
            isSelectAllChecked,
            enableSelectAll,
            columns,
            sort,
            onRequestSort,
            onRequestFilter,
            getFilterKey,
            selectAll,
            onRequesStartResize,
            isResizable,
            isColsFixed,
            isMoving,
            filterElement,
        } = this.props

        return (
            <table
                className={styles['header']}
                style={{ width: (isResizable && isColsFixed) ? 'auto' : '100%' }}
            >
                <colgroup>
                    {
                        columns.map((col) => (
                            <col
                                key={col.key}
                                width={col.width}
                            />
                        ))
                    }
                </colgroup>
                <thead>
                    <tr>
                        {
                            columns.map((col, index) => (
                                <th
                                    key={col.key}
                                    className={styles['head']}
                                    ref={(col) => this.colsRefs = [...this.colsRefs, col]}
                                >
                                    <View
                                        className={styles['col-layout']}
                                        style={{ width: (isResizable && isColsFixed) ? col.width : 'auto' }}
                                    >
                                        {
                                            enableSelectAll && index === 0 ? (
                                                <CheckBox
                                                    disabled={!data || !data.length}
                                                    className={styles['select-all']}
                                                    checked={isSelectAllChecked}
                                                    onChange={(event) => selectAll(event.target.checked)}
                                                />
                                            ) : null
                                        }
                                        <div
                                            ref={(titleNode) => this.titleNodes = [...this.titleNodes, titleNode]}
                                            className={classnames(styles['head-title'])}
                                        >
                                            <span
                                                className={classnames(styles['title'], { [styles['pointer']]: col.sortable })}
                                                onClick={() => { col.sortable && isFunction(onRequestSort) && onRequestSort({ key: col.key, type: sort.type === SortType.ASC ? SortType.DESC : SortType.ASC }) }}
                                            >
                                                {col.title}
                                            </span>

                                            {
                                                col.sortable ?
                                                    <span className={styles['inline']}>
                                                        <View className={styles['sort']}>
                                                            <View
                                                                className={styles['sort-up']}
                                                                onClick={() => { isFunction(onRequestSort) && onRequestSort({ key: col.key, type: SortType.ASC }) }}
                                                            >
                                                                <SweetIcon
                                                                    className={styles['icon']}
                                                                    size={12}
                                                                    color={col.key === sort.key && sort.type === SortType.ASC ? '#779EEA' : '#A5A8B4'}
                                                                    name={'sortUp'}
                                                                />
                                                            </View>
                                                            <View
                                                                className={styles['sort-down']}
                                                                onClick={() => { isFunction(onRequestSort) && onRequestSort({ key: col.key, type: SortType.DESC }) }}
                                                            >
                                                                <SweetIcon
                                                                    className={styles['icon']}
                                                                    size={12}
                                                                    color={col.key === sort.key && sort.type === SortType.DESC ? '#779EEA' : '#A5A8B4'}
                                                                    name={'sortDown'}
                                                                />
                                                            </View>
                                                        </View>
                                                    </span> : null
                                            }
                                            {
                                                col.filters && col.filters.length ?
                                                    <DataGridFilter
                                                        filterKey={col.key}
                                                        filters={col.filters}
                                                        element={filterElement}
                                                        onFilterChange={onRequestFilter}
                                                        getFilterKey={getFilterKey}
                                                    />
                                                    : null
                                            }
                                            {
                                                col.tips && col.tips !== '' ?
                                                    <UIIcon
                                                        className={styles['ui-icon']}
                                                        code={'\uf055'}
                                                        size={'13px'}
                                                        title={col.tips}
                                                        color={'#999'}
                                                    />
                                                    : null
                                            }
                                        </div>
                                    </View>
                                    {
                                        isResizable && data && data.length ?
                                            (
                                                <span
                                                    className={'sweetui-grid-resizer'}
                                                    style={{ opacity: isMoving ? 1 : 0 }}
                                                    onMouseDown={(event) => isFunction(onRequesStartResize) && onRequesStartResize(event, col, index)}
                                                ></span>
                                            )
                                            : null
                                    }
                                </th>
                            ))
                        }
                    </tr>
                </thead>
            </table>
        )
    }
}
