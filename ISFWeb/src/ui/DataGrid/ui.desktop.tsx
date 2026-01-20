import React from 'react';
import classnames from 'classnames';
import { map, get, noop } from 'lodash';
import DataGridBase from './ui.base';
import DataGridField from '../DataGrid.Field/ui.desktop';
import LazyLoader from '../LazyLoader/ui.desktop';
import Paginator from '../Paginator/ui.desktop';
import { CheckBox } from '@/sweet-ui'
import styles from './style.desktop';

export default class DataGrid extends DataGridBase {

    static Field = DataGridField;

    render() {
        const { EmtpyComponent, data } = this.props;

        return (
            <div className={classnames(styles['datagrid'], this.props.className)} style={{ height: this.props.height }}>
                <div className={styles['header-wrap']}>
                    <div className={classnames({ [styles['header-overflow']]: this.state.isOverflow })}>
                        <table className={styles['table']}>
                            <colgroup>
                                {
                                    // 多选复选框
                                    this.isMultiSelect ? <col width={5} /> : null
                                }
                                {
                                    this.props.children
                                }
                            </colgroup>
                            {
                                this.props.headless !== true ?
                                    <thead className={styles['heads']}>
                                        <tr className={styles['heads-tr']}>
                                            {
                                                this.isMultiSelect ?
                                                    <th className={styles['checkbox-cell']}>
                                                        <CheckBox
                                                            className={styles['checkbox']}
                                                            disabled={this.props.data.length === 0 ? true : false}
                                                            checked={this.props.data.length && this.state.selection.length === this.props.data.length ? true : false}
                                                            onChange={this.toggleSelectAll.bind(this)}
                                                            onClick={(event) => event.stopPropagation()}
                                                        />
                                                    </th> :
                                                    null
                                            }
                                            {
                                                React.Children.map(this.props.children, (child) => {
                                                    if (!child) {return;}

                                                    const { props: { key, label } } = child;

                                                    return (
                                                        <th key={key} className={styles['cell']} style={{ height: this.props.headHeight }}>
                                                            {
                                                                label
                                                            }
                                                        </th>
                                                    )
                                                })
                                            }
                                        </tr>
                                    </thead> :
                                    null
                            }
                        </table>
                    </div>
                </div>

                <div className={classnames(styles['body-wrap'], { [styles['headless-body']]: this.props.headless, [styles['static-body']]: !this.props.height || this.props.height === 'auto' })} style={{ bottom: this.props.paginator ? 39 : 0 }}>
                    {
                        EmtpyComponent && data.length === 0 ?
                            EmtpyComponent
                            :
                            <LazyLoader
                                ref={(ref) => this.lazyLoadList = ref}
                                scroll={this.state.scroll}
                                // onScroll={this.handleScroll.bind(this)}
                                limit={this.lazyLoad && this.lazyLoad.limit}
                                onChange={this.lazyLoad ? this.changePage.bind(this) : noop}
                            >
                                <table className={styles['table']}>
                                    <colgroup>
                                        {
                                            this.isMultiSelect ? <col width={5} /> : null
                                        }
                                        {
                                            this.props.children
                                        }
                                    </colgroup>
                                    <tbody className={styles['content']}>
                                        {
                                            map(this.props.data, (record, index) => {
                                                const { checkbox, disabled } = this.props.getRecordStatus(record, index);

                                                return (
                                                    <tr
                                                        key={this.props.getKey(record)}
                                                        onDoubleClick={(e) => this.onDblClickRow(e, record)}
                                                        onClick={(e) => this.onClickRow(e, record, index)}
                                                        className={classnames(
                                                            styles['record'],
                                                            this.props.rowHoverClass,
                                                            { [styles['selected']]: this.isSelected(record), [styles['strap']]: this.props.strap && (index % 2 === 0) },
                                                        )}
                                                    >
                                                        {
                                                            // 多选复选框
                                                            this.isMultiSelect ?
                                                                <td className={styles['checkbox-cell']}>
                                                                    <CheckBox
                                                                        className={styles['checkbox']}
                                                                        disabled={disabled === true || (checkbox && get(checkbox, 'disabled') === true)}
                                                                        checked={this.isSelected(record)}
                                                                        onClick={(event) => event.stopPropagation()}
                                                                        onChange={() => { this.toggleSelected(record, { multi: true }); }} />
                                                                </td>
                                                                : null
                                                        }
                                                        {
                                                            React.Children.map(this.props.children, (child, index) => {
                                                                if (!child) {return;}

                                                                const { props: { field, formatter, className, align } } = child;

                                                                return (
                                                                    <td
                                                                        className={classnames(
                                                                            styles['cell'], className,
                                                                            {
                                                                                [styles['text-align-center']]: align === 'center',
                                                                                [styles['text-align-right']]: align === 'right',
                                                                            },
                                                                        )}
                                                                        key={index}
                                                                    >
                                                                        {
                                                                            formatter(record[field], record, { checkbox, disabled })
                                                                        }
                                                                    </td>
                                                                )
                                                            })
                                                        }
                                                    </tr>
                                                )
                                            })
                                        }
                                    </tbody>
                                </table>
                            </LazyLoader>
                    }
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
            </div >
        )
    }
}