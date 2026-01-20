import React from 'react'
import { noop } from 'lodash'
import classnames from 'classnames'
import styles from './styles.mobile.css'

const TopDrawer: React.FunctionComponent<Components.TopDrawer.Props> = function TopDrawer({ open, mask, position, children, className, onClickMask, ...otherProps }) {
    return (
        <div className={styles['container']}>
            {
                mask ?
                    <div
                        className={classnames(styles['mask'], { [styles['show']]: open }, { [styles['mask-top']]: position === 'top' })}
                        onClick={onClickMask}
                    ></div> :
                    null
            }
            <div className={classnames(styles['drawer'], styles[position], { [styles['open']]: open }, { [styles['drawer-top']]: position === 'top' }, className)} {...otherProps}>{children}</div>
        </div>
    )
}

TopDrawer.defaultProps = {
    open: false,
    mask: true,
    position: 'bottom',
    onClickMask: noop,
}

export default TopDrawer