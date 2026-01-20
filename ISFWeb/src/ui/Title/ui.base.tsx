import React from 'react'

export default class Title extends React.Component<UI.Title.Props, any> {

    static defaultProps = {
        timeout: 300,
        content: undefined,
    }

    state = {
        open: false,
        position: [0, 0],
        width: 'auto',
        whiteSpace: 'pre',
        wordBreak: 'normal',
    }

    timer

    constructor(props, context) {
        super(props, context)
        this.handleMouseEnter = this.handleMouseEnter.bind(this)
        this.handleMouseLeave = this.handleMouseLeave.bind(this)
    }

    handleMouseEnter(e) {
        const position = [e.clientX, e.clientY + 20]
        if (this.props.content !== undefined) {
            this.timer = setTimeout(() => {
                this.setState({
                    open: true,
                    position,
                })
            }, this.props.timeout)
        }
    }

    handleMouseLeave() {
        clearTimeout(this.timer)

        /**
         * setTimeout 执行 setState 修复 edge 下mouseleave 事件导致鼠标移出时 hover 样式没实时生效
         */
        setTimeout(() => {
            this.setState({
                open: false,
                width: 'auto',
                whiteSpace: 'pre',
                wordBreak: 'normal',
            })
        })
    }

    componentWillUnmount() {
        if (this.timer) {
            clearTimeout(this.timer)
            this.timer = null;
        }
    }
}