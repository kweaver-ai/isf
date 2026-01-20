import React from 'react'
import classnames from 'classnames'
import { includes } from 'lodash'
import DataListBase from './ui.base'
import Item from '../DataList.Item/ui.desktop'
import styles from './styles.desktop'

export default class DataList extends DataListBase {

    static Item = Item

    render() {
        const { className, children } = this.props
        const { selections } = this.state
        return (
            <ul className={classnames(styles['list'], className)}>
                {
                    React.Children.toArray(children).map((item: React.ReactElement<UI.DataListItem.Props>, index) => {
                        const { data, ...otherProps } = item.props
                        return React.cloneElement(item, {
                            selected: includes(selections, data),
                            onToggleSelect: (e) => this.toggleSelect(e, item, index),
                            onClick: (e) => this.handleClick(e, item, index),
                            onDoubleClick: (e) => this.handleDoubleClick(e, item, index),
                            onContextMenu: (e) => this.handleContextMenu(e, item, index),
                            ...otherProps,
                        })
                    })
                }
            </ul>
        )
    }
}