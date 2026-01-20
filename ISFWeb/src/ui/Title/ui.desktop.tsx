import React from 'react'
import classnames from 'classnames'
import TitleBase from './ui.base'
import PopOver from '../PopOver/ui.desktop'
import styles from './styles.desktop'

export default class Title extends TitleBase {

    /**
     * title限制宽度
     * @param ref
     */
    setTitleSize(ref) {
        if (ref) {
            const maxWidth = Math.min(window.innerWidth, 612)
            if (ref.offsetWidth > maxWidth) {
                this.setState({
                    whiteSpace: 'pre-wrap',
                    wordWrap: 'break-word',
                    overflowWrap: 'break-word',
                    width: maxWidth - 12,
                })
            }
            const [x, y] = this.state.position
            if (y + ref.offsetHeight + 1 > window.innerHeight) {
                this.setState({
                    position: [x, y - ref.offsetHeight - 20],
                })
            }
        }
    }

    render() {
        return (
            <div
                role={this.props.role}
                className={classnames(styles['container'], { [styles['inline']]: this.props.inline })}
                onMouseEnter={this.handleMouseEnter}
                onMouseLeave={this.handleMouseLeave}
            >
                {
                    this.props.children
                }
                <PopOver
                    open={this.state.open}
                    anchorOrigin={this.state.position}
                    freezable={false}
                    watch={true}
                    style={{ zIndex: 10000 }}
                    element={this.props.element}
                >
                    {
                        typeof this.props.content === 'string' && this.props.content !== '' ?
                            <div
                                className={classnames(styles['title'], this.props.className)}
                                style={{
                                    width: this.state.width,
                                    whiteSpace: this.state.whiteSpace,
                                    wordWrap: this.state.wordWrap,
                                    overflowWrap: this.state.overflowWrap,
                                }}
                                ref={this.setTitleSize.bind(this)}
                            >
                                {this.props.content}
                            </div> :
                            this.props.content
                    }
                </PopOver>
            </div>
        )
    }
}