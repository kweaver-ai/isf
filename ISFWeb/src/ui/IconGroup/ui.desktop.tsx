import React from 'react'
import classnames from 'classnames'
import Item from '../IconGroup.Item/ui.desktop'
import styles from './styles.desktop'

const IconGroup: React.FunctionComponent<UI.IconGroup.Props> = function IconGroup({ className, children, ...otherProps }) {
    return (
        <div className={classnames(styles['container'], className)} {...otherProps}>
            {children}
        </div>
    )
}

IconGroup.Item = Item

export default IconGroup