import React from 'react'
import SliderBase from './ui.base'
import classnames from 'classnames'
import styles from './styles.desktop'

export default class Slider extends SliderBase {

    render() {
        let { min, max, step, value, axis, size, length, thickness, backgroundColor, foregroundColor, className } = this.props;
        return (
            <div ref={(container) => this.container = container} className={classnames(styles['container'], className)} style={{
                backgroundColor,
                width: axis === 'x' ? length : size,
                height: axis === 'y' ? length : (thickness ? thickness : size),
                top: -size / 2,
            }} onMouseDown={this.handleMouseDown.bind(this)}>
                <div className={styles['coverd']} style={{
                    backgroundColor: foregroundColor,
                    width: axis === 'x' ? `${value / max * 100}%` : size,
                    height: axis === 'x' ? (thickness ? thickness : size) : `${value / max * 100}%`,
                }}>
                    <span className={styles['handler']} style={{
                        backgroundColor: foregroundColor,
                        width: 2 * size,
                        height: 2 * size,
                        top: axis === 'x' ? (thickness ? -(size - thickness / 2) : -size / 2) : -size,
                        right: axis === 'x' ? -size : -size / 2,
                        borderTopRightRadius: size,
                        borderBottomRightRadius: size,
                        borderBottomLeftRadius: size,
                        borderTopLeftRadius: size,
                    }}></span>
                </div>
            </div>
        )
    }
}
