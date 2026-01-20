import React from 'react';
import classnames from 'classnames';
import { CheckBox } from '@/sweet-ui'
import LazyLoader from '../LazyLoader/ui.desktop';
import Paginator from '../Paginator/ui.desktop';
import DataGridRow from '../DataGrid.Row/ui.desktop';
import DataGridCell from '../DataGrid.Cell/ui.desktop'
import DataGridBase from './ui.base';
import styles from './style.desktop';

export default class DataGrid extends DataGridBase {
    static Row = DataGridRow;

    static Cell = DataGridCell;

    render() {
        return (
            <div className={classnames(styles['datagrid'], this.props.className)} style={{ height: this.props.height }}>
                {
                    !this.props.headless ? (
                        <div className={styles['header']}>
                            <table className={styles['table']}>
                                <colgroup>
                                    {
                                        // 多选复选框
                                        this.isMultiSelect ? <col width={12} /> : null
                                    }
                                    {
                                        this.props.fields.map(({ width }, index) => (
                                            <col
                                                key={index}
                                                width={width}
                                            />
                                        ))
                                    }
                                </colgroup>
                                {
                                    this.props.headless !== true ? (
                                        <thead className={styles['heads']}>
                                            <tr>
                                                {
                                                    this.isMultiSelect ?
                                                        <th className={styles['checkbox-cell']}>
                                                            <CheckBox
                                                                checked={this.props.children.length && this.selection.length === this.props.children.length}
                                                                onChange={this.toggleSelectAll.bind(this)}
                                                                onClick={(event) => event.stopPropagation()}
                                                            />
                                                        </th> : null
                                                }
                                                {
                                                    this.props.fields.map(({ label }, index) => (
                                                        <th
                                                            key={index}
                                                            className={styles['cell']}
                                                        >
                                                            {
                                                                label
                                                            }
                                                        </th>
                                                    ))
                                                }
                                            </tr>
                                        </thead>
                                    ) : null
                                }
                            </table>
                        </div>
                    ) : null
                }
                <div className={classnames(styles['body-wrap'], { [styles['headless-body']]: this.props.headless, [styles['static-body']]: !this.props.height || this.props.height === 'auto' })} style={{ bottom: this.props.paginator ? 34 : 0 }}>
                    <LazyLoader
                        scroll={this.state.scroll}
                        onScroll={this.handleScroll.bind(this)}
                        limit={this.lazyLoad && this.lazyLoad.limit}
                        onChange={this.lazyLoad && this.changePage.bind(this)}
                    >
                        <table className={styles['table']}>
                            <colgroup>
                                {
                                    this.isMultiSelect ? <col width={12} /> : null
                                }
                                {
                                    this.props.fields.map(({ width }, index) => (
                                        <col
                                            key={index}
                                            width={width}
                                        />
                                    ))
                                }
                            </colgroup>
                            <tbody>
                                {
                                    React.Children.map(this.props.children, (row, index) => React.cloneElement(row, {
                                        checkbox: this.isMultiSelect,
                                        selected: this.isMultiSelect ? this.selection.includes(row) : this.selection === row,
                                        onClick: this.onClickRow.bind(this, row, index),
                                        onCheckChange: this.toggleSelected.bind(this, row),
                                    }))
                                }
                            </tbody>
                        </table>
                    </LazyLoader>
                </div>
                {
                    this.props.paginator ?
                        <div className={styles['footer']}>
                            <div className={styles['footer-padding']}>
                                <Paginator {...this.props.paginator} onChange={this.changePage.bind(this)} />
                            </div>
                        </div>
                        : null
                }
            </div>
        )
    }
}