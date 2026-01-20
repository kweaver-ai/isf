import React from 'react'
import PropTypes from 'prop-types';

export default class ToastProviderBase extends React.Component<any, any> {
    static childContextTypes = {
        toast: PropTypes.func,
    }

    state = {
        toasts: [],
    }

    timer

    getChildContext() {
        return {
            toast: this.toast.bind(this),
        }
    }

    componentDidMount() {
        this.start()
    }

    componentWillUnmount() {
        this.stop()
    }

    /**
     * 开始清理消息队列
     */
    start() {
        clearTimeout(this.timer)
        const { toasts } = this.state
        if (toasts.length) {
            this.timer = setTimeout(() => {
                this.setState({
                    toasts: this.state.toasts.slice(1),
                })
                this.start()
            }, 3000 / toasts.length)
        } else {
            this.timer = undefined
        }
    }

    /**
     * 停止清除消息
     */
    stop() {
        clearTimeout(this.timer)
    }

    /**
     * toast消息入队
     * @param toast
     */
    toast(text, options = {}) {
        this.setState({
            toasts: [...this.state.toasts, [text, options]],
        }, () => {
            if (!this.timer) {
                this.start()
            }
        })
    }

    /**
     * 关闭单条消息
     * @param i 消息下标
     */
    handleClose(i) {
        this.stop()
        const { toasts } = this.state
        this.setState({
            toasts: [...toasts.slice(0, i), ...toasts.slice(i + 1)],
        })
    }
}