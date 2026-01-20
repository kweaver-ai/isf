import * as React from "react";
import classNames from "classnames";
import omit from "omit.js";
import Group from "antd/lib/input/Group";
import Search from "antd/lib/input/Search";
import TextArea from "antd/lib/input/TextArea";
import { Password } from "./Password";
import { Omit, LiteralUnion } from "antd/lib/_util/type";
import ClearableLabeledInput, { hasPrefixSuffix } from "./ClearableLabeledInput";
import { ConfigConsumer, ConfigConsumerProps } from "antd/lib/config-provider";
import SizeContext, { SizeType } from "antd/lib/config-provider/SizeContext";
import devWarning from "antd/lib/_util/devWarning";

export interface InputProps
    extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "size" | "prefix" | "type" | "value"> {
    value: string;
    prefixCls?: string;
    size?: SizeType;
    // ref: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#%3Cinput%3E_types
    type?: LiteralUnion<
        | "button"
        | "checkbox"
        | "color"
        | "date"
        | "datetime-local"
        | "email"
        | "file"
        | "hidden"
        | "image"
        | "month"
        | "number"
        | "password"
        | "radio"
        | "range"
        | "reset"
        | "search"
        | "submit"
        | "tel"
        | "text"
        | "time"
        | "url"
        | "week",
        string
    >;
    onPressEnter?: React.KeyboardEventHandler<HTMLInputElement>;
    addonBefore?: React.ReactNode;
    addonAfter?: React.ReactNode;
    prefix?: React.ReactNode;
    suffix?: React.ReactNode;
    allowClear?: boolean;
}

export function fixControlledValue<T>(value: T) {
    if (typeof value === "undefined" || value === null) {
        return "";
    }
    return value;
}

export function resolveOnChange(
    target: HTMLInputElement | HTMLTextAreaElement,
    e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement> | React.MouseEvent<HTMLElement, MouseEvent>,
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void
) {
    if (onChange) {
        let event = e;
        if (e.type === "click") {
            // click clear icon
            event = Object.create(e);
            event.target = target;
            event.currentTarget = target;
            // const originalInputValue = target.value;
            // change target ref value cause e.target.value should be '' when clear input
            target.value = "";
            onChange(event as React.ChangeEvent<HTMLInputElement>);
            // reset target ref value
            // target.value = originalInputValue;
            return;
        }
        onChange(event as React.ChangeEvent<HTMLInputElement>);
    }
}

export function getInputClassName(prefixCls: string, size?: SizeType, disabled?: boolean, direction?: any) {
    return classNames(prefixCls, {
        [`${prefixCls}-sm`]: size === "small",
        [`${prefixCls}-lg`]: size === "large",
        [`${prefixCls}-disabled`]: disabled,
        [`${prefixCls}-rtl`]: direction === "rtl",
    });
}

export interface InputState {
    focused: boolean;
}

class Input extends React.Component<InputProps, InputState> {
    static Group: typeof Group;

    static Search: typeof Search;

    static TextArea: typeof TextArea;

    static Password: typeof Password;

    static defaultProps = {
        type: "text",
    };

    input!: HTMLInputElement;

    clearableInput!: ClearableLabeledInput;

    direction: any = "ltr";

    constructor(props: InputProps) {
        super(props);
        this.state = {
            focused: false,
        };
    }

    // Since polyfill `getSnapshotBeforeUpdate` need work with `componentDidUpdate`.
    // We keep an empty function here.
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    componentDidUpdate() {}

    getSnapshotBeforeUpdate(prevProps: InputProps) {
        if (hasPrefixSuffix(prevProps) !== hasPrefixSuffix(this.props)) {
            devWarning(
                this.input !== document.activeElement,
                "Input",
                `When Input is focused, dynamic add or remove prefix / suffix will make it lose focus caused by dom structure change. Read more: https://ant.design/components/input/#FAQ`
            );
        }
        return null;
    }

    focus = () => {
        this.input.focus();
    };

    blur() {
        this.input.blur();
    }

    select() {
        this.input.select();
    }

    saveClearableInput = (input: ClearableLabeledInput) => {
        this.clearableInput = input;
    };

    saveInput = (input: HTMLInputElement) => {
        this.input = input;
    };

    onFocus: React.FocusEventHandler<HTMLInputElement> = (e) => {
        const { onFocus } = this.props;
        this.setState({ focused: true });
        if (onFocus) {
            onFocus(e);
        }
    };

    onBlur: React.FocusEventHandler<HTMLInputElement> = (e) => {
        const { onBlur } = this.props;
        this.setState({ focused: false });
        if (onBlur) {
            onBlur(e);
        }
    };

    handleReset = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
        // 修改input value属性值重置实现
        resolveOnChange(this.input, e, this.props.onChange);
    };

    renderInput = (prefixCls: string, size: SizeType | undefined, input: ConfigConsumerProps["input"] = {}) => {
        const { className, addonBefore, addonAfter, size: customizeSize, disabled } = this.props;
        // Fix https://fb.me/react-unknown-prop
        const otherProps = omit(this.props, [
            "prefixCls",
            "onPressEnter",
            "addonBefore",
            "addonAfter",
            "prefix",
            "suffix",
            "allowClear",
            // Input elements must be either controlled or uncontrolled,
            // specify either the value prop, or the defaultValue prop, but not both.
            "defaultValue",
            "size",
            "inputType",
            "value", // 删除input的value属性
        ] as any);
        return (
            <input
                autoComplete={input.autoComplete}
                {...otherProps}
                onChange={this.handleChange}
                onFocus={this.onFocus}
                onBlur={this.onBlur}
                onKeyDown={this.handleKeyDown}
                className={classNames(getInputClassName(prefixCls, customizeSize || size, disabled, this.direction), {
                    [className!]: className && !addonBefore && !addonAfter,
                })}
                ref={this.saveInput}
            />
        );
    };

    handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        // 修改input value属性值change事件实现
        resolveOnChange(this.input, e, this.props.onChange);
    };

    handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        const { onPressEnter, onKeyDown } = this.props;
        if (e.keyCode === 13 && onPressEnter) {
            onPressEnter(e);
        }
        if (onKeyDown) {
            onKeyDown(e);
        }
    };

    renderComponent = ({ getPrefixCls, direction, input }: ConfigConsumerProps) => {
        const { focused } = this.state;
        const { prefixCls: customizePrefixCls } = this.props;
        const { value, ...otherProps } = this.props; // 修改value属性读取方式
        const prefixCls = getPrefixCls("input", customizePrefixCls);
        this.direction = direction;

        return (
            <SizeContext.Consumer>
                {(size) => (
                    <ClearableLabeledInput
                        size={size}
                        {...otherProps}
                        prefixCls={prefixCls}
                        inputType="input"
                        value={fixControlledValue(value)}
                        element={this.renderInput(prefixCls, size, input)}
                        handleReset={this.handleReset}
                        ref={this.saveClearableInput}
                        direction={direction}
                        focused={focused}
                        triggerFocus={this.focus}
                    />
                )}
            </SizeContext.Consumer>
        );
    };

    render() {
        return <ConfigConsumer>{this.renderComponent}</ConfigConsumer>;
    }
}

export default Input;
