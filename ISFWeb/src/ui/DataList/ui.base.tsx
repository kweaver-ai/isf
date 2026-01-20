import React from 'react'
import { includes } from 'lodash'

export default class DataList extends React.PureComponent<UI.DataList.Props, UI.DataList.State> {
    static defaultProps = {
        selections: [],
        multiple: true,
        selecting: false,
    }

    state = {
        selections: this.props.selections,
    }

    /** 上一次点击项 */
    lastClickIndex = 0

    static getDerivedStateFromProps({ selections }, prevState) {
        if (selections !== prevState.selections) {
            return {
                selections,
            }
        }
        return null
    }

    /**
     * 切换单个选择状态
     * @param e 事件
     * @param data 数据
     */
    toggleSelect(e, item) {
        e.stopPropagation()
        const { multiple, onSelectionChange } = this.props
        const { selections } = this.state
        const { data, selectable } = item.props
        let nextSelections = []

        if (selectable) {
            nextSelections = includes(selections, data) ?
                selections.filter((selection) => selection !== data) :
                (multiple ? [...selections, data] : [data])
        }

        if (typeof onSelectionChange === 'function') {
            onSelectionChange(nextSelections, multiple)
        } else {
            this.setState({
                selections: nextSelections,
            })
        }
    }

    /**
     * 单击
     * @param e 事件
     * @param data 数据
     * @param index 下标
     */
    handleClick(e, item, index) {
        if (e.defaultPrevented) {
            return
        }
        const { onSelectionChange, onClick, children } = this.props
        const { selections } = this.state
        const { data, selectable } = item.props
        let nextSelections = [],
            { multiple } = this.props

        if (typeof onClick === 'function') {
            onClick(e, data, index)
        }

        if (selectable) {

            if (e.ctrlKey && !e.shiftKey) {
                nextSelections = includes(selections, data) ?
                    selections.filter((selection) => selection !== data) :
                    (multiple ? [...selections, data] : [data])
            } else if (!e.ctrlKey && e.shiftKey) {
                nextSelections = multiple ?
                    (
                        index > this.lastClickIndex ?
                            React.Children.toArray(children).slice(this.lastClickIndex, index + 1) :
                            React.Children.toArray(children).slice(index, this.lastClickIndex + 1)
                    ).filter(({ props }) => props.selectable).map(({ props }) => props.data)
                    : [data]
            } else {
                if (includes(selections, data) && selections.length === 1) {
                    nextSelections = []
                    this.lastClickIndex = 0
                } else {
                    nextSelections = [data]
                    this.lastClickIndex = index
                }
                multiple = false
            }
        }

        if (typeof onSelectionChange === 'function') {
            onSelectionChange(nextSelections, multiple)
        } else {
            this.setState({
                selections: nextSelections,
            })
        }
    }

    /**
     * 右击
     * @param e 右击事件
     * @param data 数据
     * @param index 下标
     */
    handleContextMenu(e, item, index) {

        const { onSelectionChange, onContextMenu } = this.props
        const { selections } = this.state
        const { data, selectable } = item.props
        let nextSelections = [],
            { multiple } = this.props

        if (typeof onContextMenu === 'function') {
            onContextMenu(e, data, index)
        }

        if (selectable) {
            if (multiple && includes(selections, data)) {
                nextSelections = selections
            } else {
                nextSelections = [data]
                multiple = false
            }
        }

        if (typeof onSelectionChange === 'function') {
            onSelectionChange(nextSelections, multiple)
        } else {
            this.setState({
                selections: nextSelections,
            })
        }
    }

    /**
     * 双击
     * @param e
     * @param data
     */
    handleDoubleClick(e, item, index) {
        const { onDoubleClick } = this.props
        if (typeof onDoubleClick === 'function') {
            onDoubleClick(e, item.props.data, index)
        }
    }
}