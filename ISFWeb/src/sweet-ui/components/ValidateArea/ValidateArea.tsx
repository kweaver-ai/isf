import React from 'react';
import { isFunction } from 'lodash';
import { SweetUIEvent } from '../../utils/event';
import ValidateTip from '../ValidateTip';
import TextArea from '../TextArea';

interface ValidateAreaProps {
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
     * 文本域是否禁用状态
     */
    disabled?: boolean;

    /**
     * 是否控制光标位置
     */
    cursorControl?: boolean;

    /**
     * 文本域输入内容的最大长度(正整数)
     */
    maxLength?: number;

    /**
     * 文本域是否支持编辑
     */
    readOnly?: boolean;

    /**
     * 文本域占位提示
     */
    placeholder?: string;

    /**
     * 文本域是否必填
     */
    required?: boolean;

    /**
     * 文本域默认内容
     */
    defaultValue?: string;

    /**
     * 文本域宽度
     */
    width?: number;

    /**
     * 文本域高度
     */
    height?: number;

    /**
     * 文本框padding
     */
    paddingRight?: number;

    /**
     * 文本域输入验证
     * @param value 文本值
     */
    validator?: (input: string) => boolean;

    /**
     * 文本域内容变化时的回调
     */
    onValueChange?: (e: SweetUIEvent<string>) => void;

    /**
     * 文本域按下回车的回调
     */
    onPressEnter?: (e: React.KeyboardEvent) => void;

    /**
     * 文本域聚焦事件回调
     */
    onFocus?: (e: React.FocusEvent<HTMLElement>) => void;

    /**
     * 文本域失焦事件回调
     */
    onBlur?: (e: React.FocusEvent<HTMLElement>) => void;

    /**
     * 鼠标进入文本域时触发
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface ValidateAreaState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateArea extends React.Component<ValidateAreaProps, ValidateAreaState> {
    state = {
        active: false,
        hover: false,
    }

    private handleFocus = (event: React.FocusEvent<HTMLElement>) => {
        this.setState({ active: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    private handleBlur = (event: React.FocusEvent<HTMLElement>) => {
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
        const { value, popoverDistance, target, getContainer, className, validateState, validateMessages, ...otherProps } = this.props
        const { active, hover } = this.state

        return (
            <ValidateTip
                placement={'right'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer, className }}
            >
                <TextArea
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