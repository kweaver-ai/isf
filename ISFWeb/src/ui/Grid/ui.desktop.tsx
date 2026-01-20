import React from 'react'
import GridBase from './ui.base'
import { chunk } from 'lodash'

export default class Grid extends GridBase {
    render() {
        const { children, cols, ...otherProps } = this.props
        return (
            <div {...otherProps}>
                {
                    chunk(children, cols).map((rows, index) => (
                        <div key={index}>
                            {
                                rows.map((col, colIndex) => (
                                    <div
                                        key={colIndex}
                                        style={{ width: `${100 / cols}%`, display: 'inline-block' }}
                                    >{col}</div>
                                ))
                            }
                        </div>
                    ))
                }
            </div>
        )
    }
}