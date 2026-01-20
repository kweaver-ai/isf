import React from 'react';
import { noop } from 'lodash';
import __ from './locale';

export enum Validator {
    // 正常
    Ok,

    // 文本框为空
    NoValue,

    // 最大错误次数
    MaxValue,

    // 最小错误次数
    MinValue,
}

export default class NumberBoxBase extends React.PureComponent<UI.NumberBox.Props, UI.NumberBox.State> {
    static defaultProps = {
        disabled: false,

        validator: (_value) => true,

        autoFocus: false,

        value: '0',

        step: 1,

        onFocus: noop,

        onBlur: noop,

        onChange: noop,

        onEnter: noop,

        onClick: noop,

        onKeyDown: noop,
    }

    timer: number | undefined;

    interval: number | undefined;

    state: UI.NumberBox.State = {
        value: String(this.props.value),
        validateState: Validator.Ok,
        validateMessages: {
            [Validator.NoValue]: __('此输入项不能为空。'),
            [Validator.MaxValue]: this.props.ValidatorMessage && this.props.ValidatorMessage.max ? this.props.ValidatorMessage.max : '',
            [Validator.MinValue]: this.props.ValidatorMessage && this.props.ValidatorMessage.min ? this.props.ValidatorMessage.min : '',
        },
    }

    componentDidUpdate(prevProps, prevState) {
        if (String(this.props.value) !== String(prevProps.value)) {
            this.setState({
                value: String(this.props.value),
            })
        }
    }

    /**
     * 增加 按下事件
     */
    protected addValue() {
        const NumberValue = Number(this.state.value)
        if (this.props.max !== undefined && NumberValue < this.props.max && !this.props.disabled) {
            this.handleChange(String(NumberValue + this.props.step))
        }
    }

    /**
     * 减少 鼠标按下事件
     */
    protected subValue() {
        const NumberValue = Number(this.state.value)
        if (this.props.min !== undefined && 0 < NumberValue && !this.props.disabled) {
            this.handleChange(String(NumberValue + (this.props.step * -1)))
        }
    }

    /**
     * 循环增加或减少
     * @param step 步进
     */
    // private setTimer(step: number | undefined) {
    //     this.timer = setTimeout(() => {
    //         this.interval = setInterval(() => {
    //             this.timer = undefined;
    //             if (step > 0 && this.props.max !== undefined && Number(this.state.value) < this.props.max) {
    //                 this.changeValue(step);
    //             } else if (step < 0 && this.props.min !== undefined && this.props.min < Number(this.state.value)) {
    //                 this.changeValue(step);
    //             } else {
    //                 this.clearTimer()
    //             }
    //         }, 100)
    //     }, 500)
    // }

    /**
     * 结束增加或减少
     */
    // protected clearTimer() {
    //     if (this.timer) {
    //         clearTimeout(this.timer);
    //         this.timer = undefined;
    //         this.props.onChange(Number(this.state.value));
    //     }
    //     if (this.interval) {
    //         clearInterval(this.interval)
    //         this.interval = undefined;
    //         this.props.onChange(Number(this.state.value));
    //     }
    // }

    /**
     * 文本框输入
     */
    protected handleChange(value: string) {
        if (value === '') {
            this.setState({
                value: value,
                validateState: Validator.NoValue,
            })
            this.props.onChange(value);
        } else if (this.props.max && Number(value) > this.props.max) {
            // 超过最大限制自动修改为最大限制的值
            this.setState({
                value: String(this.props.max),
            })
            this.props.onChange(this.props.max);
        } else if (this.props.min && Number(value) < this.props.min) {
            // 小于最小限制，依旧允许输入，外层传入错误信息时，显示错误信息
            if (this.props.ValidatorMessage && this.props.ValidatorMessage.min) {
                this.setState({
                    validateState: Validator.MinValue,
                })
            }
            this.setState({
                value: value,
            })
            this.props.onChange(value);
        } else {
            this.setState({
                value: value,
                validateState: Validator.Ok,
            })

            this.props.onChange(value);
        }
    }
}