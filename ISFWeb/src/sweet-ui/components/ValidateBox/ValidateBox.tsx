import React from 'react';
import { isFunction } from 'lodash';
import { SweetUIEvent } from '../../utils/event';
import ValidateTip from '../ValidateTip';
import TextBox from '../TextBox';
import AppConfigContext from '@/core/context/AppConfigContext';

interface ValidateBoxProps {
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
     * 文本域内容
     */
    value: string;

    /**
     * 文本域禁用状态
     */
    disabled?: boolean;

    /**
     * 文本域width，包含盒模型的padding和border
     */
    width?: string;

    /**
     * 文本框默认内容(建议非受控场景下使用)
     */
    defaultValue?: string;

    /**
     * 文本域自动聚焦
     */
    autoFocus?: boolean;

    /**
     * 文本域聚焦时选中
     */
    selectOnFocus?: [number] | [number, number] | boolean;

    /**
     * 文本域css样式
     */
    style?: React.CSSProperties;

    /**
     * 文本域输入内容发生变化时触发
     */
    onValueChange?: (event: SweetUIEvent<string>) => void;

    /**
     * 文本域按下enter键事触发
     */
    onPressEnter?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 文本域键盘按下时触发
     */
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void;

    /**
     * 文本域聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 文本域失焦时触发
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
         * 鼠标进入文本域时触发
         */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 文本域输入限制函数
     */
    validator?: (value: string) => boolean;
}

interface ValidateBoxState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateBox extends React.Component<ValidateBoxProps, ValidateBoxState> {
    static contextType = AppConfigContext;
    state = {
        active: false,
        hover: false,
    }

    private handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        this.setState({ active: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    private handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
        this.setState({ active: false });
        isFunction(this.props.onBlur) && this.props.onBlur(event)
    }

    private handleMouseEnter = (event: React.MouseEvent<HTMLElement>) => {
        this.setState({ hover: true });
        isFunction(this.props.onMouseEnter) && this.props.onMouseEnter(event)
    }

    private handleMouseLeave = (event: React.MouseEvent<HTMLElement>) => {
        this.setState({ hover: false });
        isFunction(this.props.onMouseLeave) && this.props.onMouseLeave(event)
    }

    render() {
        const { value, popoverDistance, target, getContainer, className, validateState, validateMessages, role, ...otherProps } = this.props
        const { active, hover } = this.state
        const rootContainer = getContainer || (() => this.context?.element)

        return (
            <ValidateTip
                placement={'rightTop'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer: rootContainer, className }}
                role={role}
            >
                <TextBox
                    value={value}
                    status={validateState in validateMessages ? 'error' : 'normal'}
                    {...otherProps}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseLeave={this.handleMouseLeave}
                />
            </ValidateTip>
        );
    }
}