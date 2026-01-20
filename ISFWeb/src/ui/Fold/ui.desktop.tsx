import React from 'react'
import classnames from 'classnames'
import FoldBase from './ui.base'
import UIIcon from '../UIIcon/ui.desktop'
import Expand from '../Expand/ui.desktop'
import styles from './styles.desktop'

export default class Fold extends FoldBase {
    render() {
        const { label, children, className, iconClassName, labelProps: { className: labelClassName } } = this.props
        const { open } = this.state

        return (
            <div className={classnames(styles['container'], className)}>
                <div
                    className={classnames(styles['label'], labelClassName)}
                    onClick={this.toggle.bind(this)}
                >
                    {label}
                    <UIIcon className={classnames(styles['fold-icon'], iconClassName)} code={open ? '\uf04c' : '\uf04e'} size={16} color="#9e9e9e" />
                </div>
                <Expand open={open}>
                    {children}
                </Expand>
            </div>
        )
    }
}