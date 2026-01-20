import React from 'react';
import { isFunction } from 'lodash';
import ValidateTip from '../ValidateTip';
import ComboArea from '../ComboArea';

interface ValidateComboAreaProps {
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
     * 标签，支持数组或者对象数组，配合formatter函数显示标签的内容
     */
    value: ReadonlyArray<any>;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 文本框宽度
     */
    width?: number;

    /**
     * 文本框高度
     */
    height?: number;

    /**
     * placeholder
     */
    placeholder: string;

    /**
     * 文本域css样式
     */
    style?: React.CSSProperties;

    /**
     * tag标签样式
     */
    tagClassName?: string;

    /**
     * 标签发生改变时触发，添加或者删除标签
     */
    onChange?: (value: ReadonlyArray<any>) => void;

    /**
     * 文本域聚焦时触发
     */
    onFocus?: (event: React.FocusEvent<Element>) => void;

    /**
     * 文本域失焦时触发
     */
    onBlur?: (event: React.FocusEvent<Element>) => void;

    /**
         * 鼠标进入文本域时触发
         */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface ValidateComboAreaState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateComboArea extends React.Component<ValidateComboAreaProps, ValidateComboAreaState> {
    state = {
        active: false,
        hover: false,
    }

    private handleFocus = (event: React.FocusEvent<Element>) => {
        this.setState({ active: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    private handleBlur = (event: React.FocusEvent<Element>) => {
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
        const { value, popoverDistance, target, getContainer, className, validateState, validateMessages, onChange, ...otherProps } = this.props
        const { active, hover } = this.state

        return (
            <ValidateTip
                placement={'rightTop'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer, className }}
            >
                <ComboArea
                    value={value}
                    status={validateState in validateMessages ? 'error' : 'normal'}
                    {...otherProps}
                    onChange={onChange}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseLeave={this.handleMouseLeave}
                />
            </ValidateTip>
        );
    }
}