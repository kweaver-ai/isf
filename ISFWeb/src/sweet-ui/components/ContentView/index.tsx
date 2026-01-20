import React from 'react';
import { isFunction, isUndefined, omit } from 'lodash';
import classnames from 'classnames';
import { createEventDispatcher } from '../../utils/event';
import View from '../View';
import styles from './styles';

interface ContentViewProps extends React.ClassAttributes<void> {
    className?: string;

    style?: object;

    onOverflow?: () => void;

    onUnderflow?: () => void;

    onScroll?: () => void;
}

export default class ContentView extends React.Component<ContentViewProps, any> {
    /**
     * 容器元素
     */
    container?: HTMLElement | null;

    /**
     * 记住内容是否溢出，用于在DOM更新时判断是否需要触发事件
     */
    isContentOverflow?: boolean;

    render() {
        const { className, style, ...otherProps } = this.props

        return (
            <View
                onMounted={(container) => this.container = container}
                className={(classnames(styles['custom-scrollbar'], className))}
                style={style}
                onScroll={this.dispatchScrollEvent}
                {...omit(otherProps, 'onOverflow', 'onUnderflow')}
            >
                {
                    this.props.children
                }
            </View>
        )
    }

    componentDidMount() {
        this.detectOverflow()
    }

    componentDidUpdate() {
        this.detectOverflow()
    }

    private dispatchScrollEvent = createEventDispatcher(this.props.onScroll)

    /**
     * 检测内容是否溢出
     */
    private detectOverflow() {
        const { onOverflow, onUnderflow } = this.props

        if (this.container) {
            if ((isUndefined(this.isContentOverflow) || !this.isContentOverflow) && (this.container.scrollHeight > this.container.clientHeight)) {
                this.isContentOverflow = true

                if (isFunction(onOverflow)) {
                    onOverflow()
                }
            }
            else if ((isUndefined(this.isContentOverflow) || this.isContentOverflow) && (this.container.scrollHeight <= this.container.clientHeight)) {
                this.isContentOverflow = false

                if (isFunction(onUnderflow)) {
                    onUnderflow()
                }
            }
        }
    }
}
