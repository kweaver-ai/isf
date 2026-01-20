import React from 'react'
import classnames from 'classnames'
import { ResizeObserver } from '@juggle/resize-observer'
import Button from '../Button'
import __ from './locale'
import styles from './styles'

interface Props {
    /**
     * 展开收起按钮是否固定
     */
    isFixed: boolean;

    /**
     * 最少展示行数
     */
    rows: number;

    /**
     * 行高
     */
    lineHeight: number;

    /**
     * 是否默认展开
     */
    defaultUnFold: boolean;
}

interface State {
    /**
     * 是否展开
     */
    isUnFold: boolean;

    /**
     * 是否超出最小展示行数
     */
    isMoreThanMin: boolean;
}

export default class TextCollapse extends React.Component<Props, State> {
    static defaultProps = {
        isFixed: false,
        rows: 2,
        lineHeight: 30,
        defaultUnFold: true,
    }

    state = {
        isUnFold: this.props.defaultUnFold,
        isMoreThanMin: false,
    }

    resizeObserver: any = null

    descBox: HTMLDivElement | null = null

    componentDidMount() {
        if (this.descBox) {
            this.resizeObserver = new ResizeObserver((entries, observer) => {
                const height = this.descBox.scrollHeight

                if (height / this.props.lineHeight > this.props.rows) {
                    this.setState({
                        isMoreThanMin: true,
                    })
                } else {
                    this.setState({
                        isMoreThanMin: false,
                    })
                }
            })

            this.resizeObserver.observe(this.descBox)
        }
    }

    componentWillUnmount() {
        this.resizeObserver && this.resizeObserver.disconnect()
    }

    /**
     * 展开/收起
     */
    private fold = () => {
        this.setState({
            isUnFold: !this.state.isUnFold,
        })
    }

    render() {
        const { isFixed, rows, lineHeight } = this.props
        const { isUnFold, isMoreThanMin } = this.state
        return (
            <div className={styles['wrapper']}>
                <div
                    className={styles['activity-desc']}
                    style={{
                        WebkitLineClamp: isUnFold ? 'unset' : rows,
                        maxHeight: isUnFold ? '100%' : rows * lineHeight,
                        lineHeight: `${lineHeight}px`,
                        paddingRight: isFixed ? 40 : 0,
                    }}
                    ref={(descBox) => this.descBox = descBox}
                >
                    {
                        isFixed ?
                            <>
                                {
                                    isMoreThanMin && !isUnFold ?
                                        <div className={styles['btn']}>
                                            {[
                                                '...  ',
                                                <Button
                                                    key={'unfold'}
                                                    theme={'text'}
                                                    icon={'arrowDown'}
                                                    iconSize={16}
                                                    onClick={this.fold}
                                                />,
                                            ]}
                                        </div>
                                        : null
                                }
                                {
                                    isMoreThanMin && isUnFold ?
                                        <Button
                                            className={styles['btn']}
                                            theme={'text'}
                                            icon={'arrowUp'}
                                            iconSize={16}
                                            onClick={this.fold}
                                        />
                                        : null
                                }
                                {this.props.children}
                            </>
                            :
                            <>
                                {
                                    isMoreThanMin && !isUnFold ?
                                        <div className={classnames(styles['btn'], styles['text-btn'])}>
                                            {[
                                                '...  ',
                                                <Button
                                                    key={'unfold'}
                                                    theme={'text'}
                                                    onClick={this.fold}
                                                >
                                                    {__('展开')}
                                                </Button>,
                                            ]}
                                        </div>
                                        : null
                                }
                                {this.props.children}
                                {
                                    isMoreThanMin && isUnFold ?
                                        <Button
                                            theme={'text'}
                                            className={styles['btn-no-absolute']}
                                            onClick={this.fold}
                                        >
                                            {__('收起')}
                                        </Button>
                                        : null
                                }
                            </>
                    }
                </div>
            </div>
        )
    }
}