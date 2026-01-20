import React from 'react';
import { isFunction } from 'lodash';
import { SweetUIEvent } from '../../utils/event';
import ValidateTip from '../ValidateTip';
import Select2 from '../Select2';
import SelectOption from '../Select2/Option';
import AppConfigContext from '@/core/context/AppConfigContext';

interface ValidateSelectProps {
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
     * 当前选中的条目
     */
    value: any;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * placeholder
     */
    placeholder: string;

    /**
     * 选择器的宽度，默认200
     */
    selectorWidth?: number;

    /**
     * 选择器的className
     */
    selectorClassName?: string;

    /**
     * 选择器样式
     */
    selectorStyle?: React.CSSProperties;

    /**
     * 下拉菜单宽度，默认200，与选择器同宽
     */
    menuWidth?: number;

    /**
     * 下拉菜单最大高度，包括条目高度30*n + 上下padding 4*2 + 上下 border 1*2
     */
    menuMaxHeight?: number;

    /**
     * 下拉菜单的className
     */
    menuClassName?: string;

    /**
     * 选中时触发
     */
    onChange: (event: SweetUIEvent<string>) => void;

    /**
     * 下拉选项展开状态变化时触发
     */
    onPopupVisibleChange?: (event: SweetUIEvent<boolean>) => void;

    /**
    * 输入框聚焦时触发
    */
    onFocus?: (event: React.FocusEvent<HTMLInputElement>) => void;

    /**
     * 输入框失焦时触发
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

    role?: string;
}

interface ValidateSelectState {
    /**
     * 是否处于active
     */
    active: boolean;

    /**
     * 是否处于hover
     */
    hover: boolean;
}

export default class ValidateSelect extends React.Component<ValidateSelectProps, ValidateSelectState> {
    static contextType = AppConfigContext;
    static Option = SelectOption;

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
        const { value, popoverDistance, target, getContainer, className, validateState, validateMessages, onChange, role, ...otherProps } = this.props
        const { active, hover } = this.state

        const rootContainer = getContainer || (() => this.context?.element);

        return (
            <ValidateTip
                role={role}
                placement={'rightTop'}
                content={validateMessages[validateState]}
                visible={(validateState in validateMessages) && (active || hover)}
                tipStatus={'error'}
                {...{ popoverDistance, target, getContainer: rootContainer, className }}
            >
                <Select2
                    value={value}
                    status={validateState in validateMessages ? 'error' : 'normal'}
                    {...otherProps}
                    onChange={onChange}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onMouseEnter={this.handleMouseEnter}
                    onMouseLeave={this.handleMouseLeave}
                >
                    {this.props.children}
                </Select2>
            </ValidateTip>
        );
    }
}