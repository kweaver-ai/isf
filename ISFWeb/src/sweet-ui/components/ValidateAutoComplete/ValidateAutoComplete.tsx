import React from 'react';
import { isFunction } from 'lodash';
import ValidateTip from '../ValidateTip';
import AutoComplete, { type AutoCompleteProps } from '../AutoComplete';

interface ValidateAutoCompleteProps extends AutoCompleteProps {
    /**
     * 验证状态
     */
    validateState: any;

    /**
     * 验证信息
     */
    validateMessages: {
        [key: string]: string;
    };

    /**
     * 气泡距文本框的距离，不包括箭头，单位 px
     */
    popoverDistance?: number;

    /**
     * 气泡样式
     */
    className?: string;

    /**
     * 气泡获取浮层渲染父节点，默认渲染到body上
     */
    getContainer: () => HTMLElement;

    /**
     * 气泡事件监听的目标元素，默认window
     */
    target?: HTMLElement;

    /**
     * 文本域聚焦事件回调
     */
    onFocus?: (e: FocusEvent) => void;

    /**
     * 文本域失焦事件回调
     */
    onBlur?: (e: FocusEvent) => void;

    /**
     * 鼠标进入文本域时触发
     */
    onMouseEnter?: (event: MouseEvent) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: MouseEvent) => void;
}

interface ValidateAutoCompleteState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateAutoComplete extends React.Component<ValidateAutoCompleteProps, ValidateAutoCompleteState> {
    state = {
        active: false,
        hover: false,
    }

    private handleFocus = (event: FocusEvent) => {
        this.setState({ active: true })

        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    private handleBlur = (event: FocusEvent) => {
        this.setState({ active: false })

        isFunction(this.props.onBlur) && this.props.onBlur(event)
    }

    private handleMouseEnter = (event: MouseEvent) => {
        this.setState({ hover: true })

        isFunction(this.props.onMouseEnter) && this.props.onMouseEnter(event)
    }

    private handleMouseLeave = (event: MouseEvent) => {
        this.setState({ hover: false })

        isFunction(this.props.onMouseLeave) && this.props.onMouseLeave(event)
    }

    render() {
        const {
            popoverDistance,
            target,
            getContainer,
            className,
            validateState,
            validateMessages,
            ...otherProps
        } = this.props

        const {
            active,
            hover,
        } = this.state

        return (
            <ValidateTip
                placement={'rightTop'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer, className }}
            >
                <AutoComplete
                    validateStatus={validateState in validateMessages ? 'error' : 'normal'}
                    {...otherProps}
                    element={isFunction(getContainer) ? getContainer() : undefined}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseLeave={this.handleMouseLeave}
                />
            </ValidateTip>
        );
    }
}