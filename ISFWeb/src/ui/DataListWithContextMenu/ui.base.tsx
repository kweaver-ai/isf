import React from 'react'
import { noop } from 'lodash'

export default class DataListWithContextMenuBase extends React.Component<UI.DataListWithContextMenu.Props, UI.DataListWithContextMenu.State> {

    static defaultProps = {
        onSelectionChange: noop,
        contentMenu: null,
    }

    state = {
        open: false,
        selections: this.props.selections || [],
        contextMenuPosition: [0, 0],
    }

    /**
     * 选中项变化
     * @param selections
     * @param multiple
     */
    handleSelectionChange(selections, multiple) {
        this.setState({
            selections,
        }, () => {
            this.props.onSelectionChange(this.state.selections, multiple)
        })
    }

    /**
     * 右键
     * @param e
     */
    handleContextMenu(e) {
        if (!e.defaultPrevented) {
            e.preventDefault()
            this.setState({
                open: true,
                contextMenuPosition: [e.clientX, e.clientY],
            })
        }
        this.setState({
            selections: [],
        }, () => {
            this.props.onSelectionChange(this.state.selections, this.props.multiple)
        })
    }

    /**
     * 关闭右键菜单
     */
    closeContextMenu() {
        this.setState({
            open: false,
        })
    }
}