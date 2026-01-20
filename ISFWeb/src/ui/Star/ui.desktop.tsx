import React from 'react'
import StarBase from './ui.base'
import FontIcon from '../FontIcon/ui.desktop'
import styles from './styles.desktop'
import __ from './locale'
import { range } from 'lodash'

export default class Star extends StarBase {
    render() {
        const { score, onStar, size, color, dashed, solid, dashedFallback, solidFallback } = this.props
        return (
            onStar && typeof onStar === 'function'
                ? <div className={styles['container']}>
                    {
                        range(1, 6).map((num) => (
                            <FontIcon
                                font="AnyShare"
                                onClick={() => onStar(num)}
                                key={num}
                                size={size}
                                color={num <= this.state.score ? color : '#aaa'}
                                code={num <= this.state.score ? solid : dashed}
                                fallback={num <= this.state.score ? solidFallback : dashedFallback}
                                onMouseOver={this.handleMouseEnter.bind(this, num)}
                                onMouseLeave={this.handleMouseLeave.bind(this)}
                                title={__('${score}分', { score: num })}
                            />
                        ))
                    }
                </div>
                : <div className={styles['container']} title={__('${score}分', { score })}>
                    {
                        range(1, Math.ceil(score) + 1).map((num) => (
                            num <= score
                                ? <span className={styles['star']} key={num} style={{ fontSize: size }}>
                                    <FontIcon font="AnyShare" code={solid} fallback={solidFallback} size={size} key={num} color={color} />
                                </span>
                                : <span className={styles['star']} key={num} style={{ fontSize: size }}>
                                    <FontIcon font="AnyShare" code={dashed} fallback={dashedFallback} size={size} color={color} />
                                    <span className={styles['star-highlight']} style={{ width: `${score % 1}em` }}>
                                        <FontIcon font="AnyShare" code={solid} fallback={solidFallback} size={size} color={color} />
                                    </span>
                                </span>
                        ))
                    }
                </div>
        )
    }
}