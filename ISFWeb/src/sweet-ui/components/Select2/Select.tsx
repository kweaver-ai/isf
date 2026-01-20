import React from 'react';
import classnames from 'classnames';
import { isFunction, isEqual } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import Trigger from '../Trigger';
import Menu from '../Menu';
import SelectOption from './Option';
import SingleSelection from './Selector/SingleSelector'
import styles from './styles';
import AppConfigContext from '@/core/context/AppConfigContext';

interface SelectProps {
    /**
     * 当前选中的条目
     */
    value?: any;

    /**
     * placeholder
     */
    placeholder?: string;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 文本框状态
     */
    status?: 'normal' | 'error';

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
     *  下拉框层级
     */
    popupZIndex?: number;

    /**
     * 弹出层展开时是否冻结滚动条
     */
    freeze?: boolean;

    /**
     * 选中时触发
     */
    onChange?: (event: SweetUIEvent<string>) => void;

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
}

interface SelectState {
    /**
     * 选择项的标识
     */
    value: any;

    /**
     * 选择项的内容
     */
    label: string | React.ReactNode;

    /**
     * 是否聚焦状态
     */
    active: boolean;
}

export default class Select2 extends React.Component<SelectProps, SelectState> {
    static contextType = AppConfigContext
    static defaultProps = {
        disabled: false,
        selectorWidth: 200,
        menuWidth: 200,
        menuMaxHeight: 190,
        status: 'normal',
        freeze: false,
    };

    static Option = SelectOption;

    state = {
        value: this.props.value,
        label: '',
        active: false,
    };

    /**
     * 输入框ref
     */
    inputRef: null | HTMLElement = null;

    componentDidMount() {
        this.updateSelected(this.props.children)
    }

    componentDidUpdate(prevProps: SelectProps, prevState: SelectState) {
        if (this.props.children !== prevProps.children || this.props.value !== prevState.value) {
            this.updateSelected(this.props.children);
        }
    }

    private updateSelected(options: any) {
        const optionsList = Array.isArray(options) ? options : [options];
        const selectedItem = optionsList.filter((option) => isEqual(option.props.value, this.props.value))

        this.setState({
            value: this.props.value,
            label: selectedItem.length === 1 ? selectedItem[0].props.children : '',
        });
    }

    /**
     * 打开或者关闭弹窗触发
     */
    handlePopupVisibleChange = (event: SweetUIEvent<boolean>) => {
        this.dispatchPopupVisibleChangeEvent(event.detail);
    };

    dispatchPopupVisibleChangeEvent = createEventDispatcher(this.props.onPopupVisibleChange);

    /**
     * 点击文本框
     */
    handleClick(setPopupVisibleOnClick, e) {
        const { disabled } = this.props

        e.preventDefault();
        e.stopPropagation();
        if (disabled) {
            return
        }
        this.setState(() => ({ active: true }))

        setPopupVisibleOnClick()

        this.inputRef && this.inputRef.focus()
    }

    /**
     * 文本框聚焦
     */
    handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        const { disabled, onFocus } = this.props
        if (disabled) {
            return
        }
        this.setState(() => ({ active: true }))

        onFocus && isFunction(onFocus) && onFocus(event)
    }

    /**
     * 文本框失焦
     */
    handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
        const { onBlur } = this.props
        this.setState(() => ({ active: false }))

        onBlur && isFunction(onBlur) && onBlur(event)
    }

    /**
     * 选中项发生改变
     */
    handleSelect = (e, item, close) => {
        const { disabled, onClick, value, children } = item.props

        if (disabled) {
            return
        }

        onClick && isFunction(onClick) && onClick(e)

        if (!e.defaultPrevented) {
            this.setState(() => ({
                value: value,
                label: children,
                active: true,
            }))

            this.inputRef && this.inputRef.focus()
            this.dispatchSelectEvent(value)
            close()
        }
    };

    dispatchSelectEvent = createEventDispatcher(this.props.onChange);

    render() {
        const { children, status, disabled, placeholder, selectorWidth, selectorStyle, selectorClassName, menuWidth, menuMaxHeight, menuClassName, onMouseEnter, onMouseLeave, role, element } = this.props;
        const { value, label, active } = this.state;
        const rootElement = element || this.context && this.context.element

        return (
            <Trigger
                role={role}
                element={rootElement}
                renderer={
                    ({ setPopupVisibleOnClick, role }) => (
                        <SingleSelection
                            key={'singleSelection'}
                            label={label}
                            width={selectorWidth}
                            className={classnames(styles['wrapper'], selectorClassName)}
                            style={selectorStyle}
                            placeholder={placeholder}
                            active={active}
                            status={status}
                            disabled={disabled}
                            onClick={(e) => this.handleClick(setPopupVisibleOnClick, e)}
                            onMounted={(node) => this.inputRef = node}
                            onFocus={this.handleFocus}
                            onBlur={this.handleBlur}
                            onMouseEnter={onMouseEnter}
                            onMouseLeave={onMouseLeave}
                            role={role}
                        />
                    )}
                freeze={this.props.freeze}
                onPopupVisibleChange={this.handlePopupVisibleChange}
                popupZIndex={this.props.popupZIndex}
            >
                {
                    ({ close }) => (
                        <Menu
                            width={menuWidth}
                            maxHeight={menuMaxHeight}
                            className={classnames(styles['select-menu'], menuClassName)}
                        >
                            {
                                React.Children.map(children, (child, index) =>
                                    React.cloneElement(child, {
                                        selected: isEqual(child.props.value, value),
                                        onClick: (e) => this.handleSelect(e, child, close),
                                    }),
                                )
                            }
                        </Menu>
                    )
                }
            </Trigger>
        );
    }
}
