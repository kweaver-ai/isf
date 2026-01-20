import React from 'react';
import { isFunction } from 'lodash';
import { SweetUIEvent } from '../../utils/event';
import ValidateTip from '../ValidateTip';
import NumberBox from '../NumberBox';

interface ValidateNumberProps {
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
     * 数字框初始值
     */
    defaultNumber?: number;

    /**
     * 数字框当前值，null表示空值
     */
    value?: number | null;

    /**
     * 数字框宽度
     */
    width?: number | string;

    /**
     * 数字框禁用
     */
    disabled?: boolean;

    /**
     * 数字框最小值
     */
    min?: number;

    /**
     * 数字框最大值
     */
    max?: number;

    /**
     * 浮点数值精度，指定保留小数位数，要求是非负整数
     */
    precision?: number;

    /**
     * 按下鼠标向上/向下键时的步进，可以是小数
     */
    step?: number;

    /**
     * 当渲染数字框时，焦点是否自动落在输入框元素上
     */
    autoFocus?: boolean;

    /**
     * 设置数字输入框是否只读
     */
    readOnly?: boolean;

    /**
     * 数字框聚焦时自动选中内容
     */
    selectOnFocus?: [number] | [number, number] | boolean;

    /**
     * 文本框限制字符长度
     */
    maxLength?: number;

    /**
     * 数字框占位符
     */
    placeholder?: string;

    /**
     * 数字框数值发生变化时触发
     */
    onValueChange?: (event: SweetUIEvent<number | null>) => void;

    /**
     * 数字框聚焦事件回调
     */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 数字框失焦事件回调
     */
    onBlur?: (event: React.FocusEvent<HTMLInputElement>, value: string | number) => void;

    /**
     * 鼠标进入文本域时触发
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface ValidateNumberState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateNumber extends React.Component<ValidateNumberProps, ValidateNumberState> {
    state = {
        active: false,
        hover: false,
    }

    private handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        this.setState({ active: true });
        isFunction(this.props.onFocus) && this.props.onFocus(event)
    }

    private handleBlur = (event: React.FocusEvent<HTMLInputElement>, val: number | string) => {
        this.setState({ active: false });
        isFunction(this.props.onBlur) && this.props.onBlur(event, val)
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
        const { role, value, popoverDistance, target, getContainer, className, validateState, validateMessages, ...otherProps } = this.props
        const { active, hover } = this.state

        return (
            <ValidateTip
                role={role}
                placement={'rightTop'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer, className }}
            >
                <NumberBox
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