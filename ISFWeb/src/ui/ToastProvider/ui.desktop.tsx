import React from 'react'
import ToastProviderBase from './ui.base'
import PopOver from '../PopOver/ui.desktop'
import Toast from '../Toast/ui.desktop'
import styles from './styles.desktop'

class ShouldRerenderChildren extends React.PureComponent<any, any> {
    render() {
        return this.props.children
    }
}

export default class ToastProvider extends ToastProviderBase {

    static childContextTypes = ToastProviderBase.childContextTypes

    render() {
        const { toasts } = this.state
        return (
            <div
                role={this.props.role}
                className={this.props.className}
            >
                <PopOver
                    open={!!toasts.length}
                    anchorOrigin={['center', 100]}
                    targetOrigin={['center', 'top']}
                    freezable={false}
                    autoFix={false}
                    watch={!!toasts.length}
                >
                    {
                        toasts.map(([text, options], i) => (
                            <div
                                key={i}
                                className={styles['container']}
                                onMouseEnter={this.stop.bind(this)}
                                onMouseLeave={this.start.bind(this)}
                            >
                                <Toast {...options} onClose={this.handleClose.bind(this, i)}>
                                    {text}
                                </Toast>
                            </div>
                        ))
                    }
                </PopOver>
                <ShouldRerenderChildren>
                    {this.props.children}
                </ShouldRerenderChildren>
            </div>
        )
    }
}