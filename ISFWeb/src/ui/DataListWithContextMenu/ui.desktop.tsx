import React from 'react'
import DataListWithContextMenuBase from './ui.base'
import DataList from '../DataList/ui.desktop'

export default class DataListWithContextMenu extends DataListWithContextMenuBase {

    static Item = DataList.Item

    render() {
        const { contextMenu: ContextMenu, ...otherProps } = this.props
        const { open, contextMenuPosition, selections } = this.state
        return (
            <div onContextMenu={this.handleContextMenu.bind(this)}>
                <DataList {...otherProps} onSelectionChange={this.handleSelectionChange.bind(this)} />
                {
                    ContextMenu ?
                        <ContextMenu
                            open={open}
                            selections={selections}
                            position={contextMenuPosition}
                            onClose={this.closeContextMenu.bind(this)}
                        /> :
                        null
                }
            </div>
        )
    }
}