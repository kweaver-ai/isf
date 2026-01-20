import React from 'react';
import { noop, without, union, get, includes } from 'lodash';

export default class DataGridBase extends React.Component<UI.DataGrid.Props, void> {

    static defaultProps = {
        fields: [],

        height: 'auto',

        onClickRow: noop,

        onDblClickRow: noop,

        onSelectionChange: noop,

        onPageChange: noop,

        select: false,

        locator: (data) => 0,

        paginator: false,

        lazyLoad: false,
    }

    // 是否是多选
    isMultiSelect: boolean = get(this.props.select, 'multi') === true

    // 使用shift复选时，需要根据lastSelectIndex与当前点击的行的index，确定选中范围
    lastSelectedIndex: number = 0;

    selection: Array<any> & any = [];

    state = {
    }

    protected onClickRow(row, index, event) {
        switch (true) {
            case event.shiftKey: {
                const range = this.lastSelectedIndex < index ? [this.lastSelectedIndex, index + 1] : [index, this.lastSelectedIndex + 1]
                const nextSelection = this.props.children.slice.apply(this.props.children, range);
                this.rangeSelect(nextSelection);
                break;
            }

            case event.ctrlKey:
                if (this.isSelected(row)) {
                    this.multiUnselect(row);
                } else {
                    this.multiSelect(row);
                }
                break;

            default:
                this.lastSelectedIndex = index;
                this.toggleSelected(row);
                break;
        }

        this.props.onClickRow(row);

        event.stopPropagation();
    }

    protected toggleSelected(row, { multi = false } = {}) {
        // 单选模式
        if (this.props.select === true || (get(this.props.select, 'multi') !== undefined && !this.isMultiSelect)) {
            if (this.selection === row && !get(this.props.select, 'required')) {
                this.singleUnselect(row);
            } else {
                this.singleSelect(row);
            }
        }
        else if (this.isMultiSelect) {
            if (multi) {
                if (this.isSelected(row)) {
                    this.multiUnselect(row)
                } else {
                    this.multiSelect(row)
                }
            } else {
                if (this.selection.length === 1 && this.selection[0] === row) {
                    this.singleUnselect(row);
                } else {
                    this.singleSelect(row)
                }
            }
        }
    }

    /**
     * 选择范围内的文档
     * @param records 选择范围
     */
    private rangeSelect(records) {
        this.selection = records;
        this.forceUpdate();
        this.fireSelectionChangeEvent(this.selection);
    }

    /**
     * 单项选择
     */
    private singleSelect(record: Record<string, any>) {
        if (record === this.selection && get(this.props.select, 'required') && !this.isMultiSelect) {
            return;
        }
        this.selection = this.isMultiSelect ? [record] : record;
        this.forceUpdate();
        this.fireSelectionChangeEvent(this.selection);
    }

    /**
     * 单项取消选中
     */
    private singleUnselect(record: Record<string, any>) {
        this.selection = this.isMultiSelect ? this.getUnselectSelection(record) : null;
        this.forceUpdate();
        this.fireSelectionChangeEvent(this.selection)
    }

    private multiUnselect(records) {
        this.selection = this.getUnselectSelection(records);
        this.forceUpdate();
        this.fireSelectionChangeEvent(this.selection);
    }

    private multiSelect(records) {
        this.selection = this.getMultiSelection(records);
        this.forceUpdate();
        this.fireSelectionChangeEvent(this.selection);
    }

    /**
     * 选中数据对象
     * @param record 数据对象
     */
    private getMultiSelection(records: Record<string, any> | ReadonlyArray<Record<string, any>>) {
        return union(this.selection, Array.isArray(records) ? records : [records])
    }

    /**
     * 取消选中数据对象
     * @param record 数据对象
     */
    private getUnselectSelection(records: Record<string, any> | ReadonlyArray<Record<string, any>>): ReadonlyArray<Record<string, any>> {
        return Array.isArray(records) ? without(this.selection, ...records) : without(this.selection, records);
    }

    private isSelected(record) {
        return this.isMultiSelect ? includes(this.selection, record) : this.selection === record
    }

    /**
     * 选中项改变时触发
     */
    private fireSelectionChangeEvent(selection) {
        this.props.onSelectionChange(selection);
    }

    protected toggleSelectAll() {
        this.rangeSelect(this.selection.length === this.props.children.length ? [] : this.props.children.concat());
    }

    protected handleScroll(scroll) {
        this.setState({ scroll })
    }

    /**
     * 翻页
     * @param page 页码
     * @param limit 每页条数
     */
    protected changePage(page: number, limit: number) {
        this.firePageChangeEvent(page, limit);
    }

    /**
     * 翻页时触发
     * @param page 页码
     * @param limit 每页条数
     */
    private firePageChangeEvent(page: number, limit: number) {
        this.props.onPageChange(page, limit);
    }
}