import React from 'react'
import classnames from 'classnames'
import PopOver from '../PopOver/ui.desktop'
import Item from '../PopMenu.Item/ui.desktop'
import styles from './styles.desktop'

const PopMenu: React.FunctionComponent<UI.PopMenu.Props> = function PopMenu({ children, className, ...otherProps }) {
    return (
        <PopOver {...otherProps}>
            <ul className={classnames(styles['list'], className)}>
                {children}
            </ul>
        </PopOver>
    )
}

PopMenu.Item = Item

export default PopMenu