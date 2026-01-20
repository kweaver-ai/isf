import React from 'react'

export default class ExpandBase extends React.Component<UI.Expand.Props, UI.Expand.State> {
    static defaultProps = {
        open: false,
    }

    state = {
        marginTop: 0,
        loaded: false,
        animation: false,
    }

    animationTimer = null

    content: HTMLDivElement | null = null

    componentDidMount() {
        this.setState({
            marginTop: this.content ? -this.content.offsetHeight : 0,
            loaded: true,
        }, () => {
            this.clearAniTimer()
            this.animationTimer = setTimeout(() => {
                this.setState({
                    animation: true,
                })
            })
        })
    }

    componentDidUpdate() {
        const marginTop = this.content ? -this.content.offsetHeight : 0
        // 范围取值，因为offsetHeight取值为四舍五入，避免造成死循环
        if (marginTop > this.state.marginTop + 1 || marginTop < this.state.marginTop - 1) {
            this.setState({
                marginTop,
            })
        }
    }

    clearAniTimer = () => {
        if (this.animationTimer) {
            clearTimeout(this.animationTimer);
            this.animationTimer = null;
        }
    }

    componentWillUnmount() {
        this.clearAniTimer()
    }
}