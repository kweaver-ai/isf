import React from 'react'
import classnames from 'classnames'
import { ClassName } from '../helper'
import styles from './styles.desktop'

const AppBar: React.FunctionComponent<any> = function AppBar({ className, children, ...props } = {}) {
    return (
        <div className={classnames(styles['app-bar'], ClassName.BorderTopColor, className)} {...props}>{children}</div>
    )
}

export default AppBar