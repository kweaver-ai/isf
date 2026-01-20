import React from 'react';
import { isFunction } from 'lodash';
import classnames from 'classnames';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import SweetIcon from '../SweetIcon';
import styles from './styles';

/**
 * 最大可以为15为整数
 */
const MAX_SAFE_INTEGER = 999999999999999;

/**
 * 最小可以为-15为整数
 */
const MIN_SAFE_INTEGER = -999999999999999;

/**
 * 最小浮点数之间的差值
 */
const EPSILON = Math.pow(2, -52);

enum IncrementDirection {
    /**
     * 向上
     */
    Up = 1,

    /**
     * 向下
     */
    Down = -1,
}

interface NumberBoxProps {
    /**
     * 初始值
     */
    defaultNumber?: number;

    /**
     * 当前值，null表示空值
     */
    value?: number | null;

    /**
     * 输入框宽度
     */
    width?: number | string;

    /**
     * className
     */
    className?: string;

    /**
     * 禁用
     */
    disabled?: boolean;

    /**
     * 最小值
     */
    min?: number;

    /**
     * 最大值
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
     * 文本框状态
     */
    status?: 'normal' | 'error';

    /**
     * 占位符
     */
    placeholder?: string;

    /**
     * role
     */
    role?: string;

    /**
     * 文本框限制字符长度
     */
    maxLength?: number;

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
     * 鼠标进入文本域时触发（外层View）
     */
    onMouseEnter?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 鼠标移除文本域时候触发（外层View）
     */
    onMouseLeave?: (event: React.MouseEvent<HTMLElement>) => void;

    /**
     * 点击事件
     */
    onClick?: (event: React.MouseEvent<HTMLElement>) => void;
}

interface NumberBoxState {
    focused: boolean;
    inputValue: number | string;
    value: number | string;
}

export default class NumberBox extends React.Component<NumberBoxProps, NumberBoxState> {
    static defaultProps = {
        disabled: false,
        autoFocus: false,
        width: 80,
        selectOnFocus: false,
        max: MAX_SAFE_INTEGER,
        min: MIN_SAFE_INTEGER,
        status: 'normal',
    };

    constructor(props: NumberBoxProps, ...args: any[]) {
        super(props);

        const initialVal = 'value' in props ? props.value : props.defaultNumber;
        const validValue = this.getValidValue(initialVal);

        this.state = {
            inputValue: validValue,
            value: validValue,
            focused: props.autoFocus!,
        };
    }

    componentDidUpdate(prevProps: NumberBoxProps) {
        if ('defaultNumber' in this.props && !('value' in this.props) && this.props.defaultNumber !== prevProps.defaultNumber) {
            const validValue = this.getValidValue(this.props.defaultNumber)

            this.setState({
                inputValue: validValue,
                value: validValue,
            })
        }

        if ('value' in this.props && this.props.value !== prevProps.value) {
            if (!this.inputting) {
                const validValue = this.getValidValue(this.props.value);
                this.setState({
                    inputValue: validValue,
                    value: validValue,
                })
            }
        }

        if ('value' in this.props && this.toNumber(this.props.value) !== this.toNumber(this.state.inputValue)) {
            const validValue = this.getValidValue(this.props.value);
            this.setState({
                inputValue: validValue,
                value: validValue,
            })
        }

        // 当最大值/最小值改变时，触发onValueChange事件
        const nextValue = 'value' in this.props ? this.props.value : this.state.inputValue;
        const { onValueChange, max, min } = prevProps;

        if (
            typeof this.props.max === 'number' &&
            this.props.max !== max &&
            typeof nextValue === 'number' &&
            nextValue > this.props.max &&
            isFunction(onValueChange)
        ) {
            this.setState({
                value: this.props.max,
                inputValue: this.props.max,
            });
            this.dispatchValueChangeEvent(this.props.max);
        }
        if (
            typeof this.props.min === 'number' &&
            this.props.min !== min &&
            typeof nextValue === 'number' &&
            nextValue < this.props.min &&
            isFunction(onValueChange)
        ) {
            this.setState({
                value: this.props.min,
                inputValue: this.props.min,
            });
            this.dispatchValueChangeEvent(this.props.min);
        }
    }

    /**
     * 输入框ref
     */
    numberInputRef = React.createRef<HTMLInputElement>();

    /**
     * 表示是否正在输入
     */
    inputting: boolean = false;

    /**
     * 将传入参数转换为合法的格式   ---  string
     */
    getValidValue = (value: any) => {
        // 如果传递过来的是 '' | null | 或者形如 -  -. aa 等非法数字都被处理为 ''
        return this.isNotCompleteNumber(value) ? '' : this.getPrecisionValue(this.getValidValueByRange(value))
    };

    /**
     * 获得有效值
     * 非数字直接返回，数字则检查是否在有效值（min,max）范围内，如没有则用有效值还原
     */
    getValidValueByRange = (value: any, min = this.props.min, max = this.props.max) => {
        let val = Number(value);

        if (this.isNotCompleteNumber(value)) {
            return value;
        }
        if (typeof min === 'number' && val < min) {
            val = min;
        }
        if (typeof max === 'number' && val > max) {
            val = max;
        }

        return val;
    }

    /**
     * 获取精度
     */
    getPrecision: () => number | undefined = () => {
        const { precision } = this.props

        return (typeof precision === 'number' && precision >= 0) ? precision : undefined // undefined 表示没有精度限制，即无小数位数限制，0表示整数
    };

    /**
     * 根据精度展示为相应的值   ---  string
     */
    getPrecisionValue = (value: any) => {
        const precision = this.getPrecision();

        return this.isNotCompleteNumber(value) ? value : precision !== undefined ? Number(value).toFixed(precision) : Number(value).toString()
    }

    /**
     * 能够转换为Number类型的转为Number类型
     * @memberof NumberBox
     */
    toNumber = (num: any) => {
        // 如果是不完整的输入，则返回字符串类型，否则返回数字类型
        return this.isNotCompleteNumber(num) ? num : Number(num)
    }

    /**
     * 数字输入框发生输入事件
     */
    handleValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const val = event.target.value.trim()
        if (this.state.focused) {
            this.inputting = true;
        }

        if (val === '') {
            this.setState({
                inputValue: val,
            })
            this.dispatchValueChangeEvent(val)
        } else {
            const { min } = this.props
            const precision = this.getPrecision()

            if (typeof min === 'number' && min >= 0 && val.includes('-') || // 最小值大于等于0时，不允许输入 -
                (precision === 0 && val.includes('.')) || // 精度为0时，不允许输入小数点
                (typeof precision === 'number' && precision > 0 && val.indexOf('.') !== -1 && val.split('.')[1].length > precision) || // 小数点后面的位数超过精度限定的长度，无法输入
                val.replace(/[^-]/g, '').length > 1 || // 不允许输入两个以上的-
                val.replace(/[^.]/g, '').length > 1 || // 不允许输入两个以上的.
                val.indexOf('-') > 0 || // 不允许在非第一的位置输入负号
                val.indexOf('.') === 0 || // 不允许在第一个的位置输入 .
                val.indexOf('-.') === 0 || // 不允许直接输入-.
                !(/^[0-9\-\.]+$/.test(val)) || // 不允许输入数字|-|.以外的内容
                (
                    typeof precision === 'number' ?
                        (val.replace(/\-/g, '').match(/[0-9]*/)[0].length > (15 - precision))
                        : val.replace(/\.|\-/g, '').length > 15
                ) // 除小数点和负号，整数部分+小数部分之和不允许超过15位
            ) {
                // 不合法的输入不响应
                return
            } else {
                this.setState({
                    inputValue: val,
                })
                this.dispatchValueChangeEvent(this.toNumber(val))
            }
        }
    };

    dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange);

    /**
     * 不完整的输入
     * '-' '' null => are not complete numbers
     */
    isNotCompleteNumber = (num: any) => {
        return isNaN(num) || num === '' || num === null;
    }

    /**
     * 输入框失焦时的处理函数
     */
    private handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
        this.inputting = false;
        this.setState({
            focused: false,
        });
        const { inputValue } = this.state;
        const val = this.isNotCompleteNumber(inputValue) ? inputValue : this.getPrecisionValue(this.getValidValueByRange(inputValue))

        this.setState({
            inputValue: val,
        })
        this.dispatchValueChangeEvent(this.toNumber(val));

        isFunction(this.props.onBlur) && this.props.onBlur(event, val);
    };

    /**
     * 输入框聚焦时的处理函数
     */
    private handleFocus = (event: React.FocusEvent<HTMLInputElement>) => {
        const { selectOnFocus } = this.props;

        this.setState({
            focused: true,
        });

        if (selectOnFocus && this.numberInputRef.current) {
            if (Array.isArray(selectOnFocus)) {
                // 指定格式化显示，必须指定聚焦选中起始位置
                this.numberInputRef.current.selectionStart = selectOnFocus[0];
                this.numberInputRef.current.selectionEnd = selectOnFocus[1] || this.numberInputRef.current.value.length;
            }
        }

        isFunction(this.props.onFocus) && this.props.onFocus(event);
    };

    /**
     * keydown事件触发时的处理函数
     */
    handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.keyCode === 38) {
            this.updateValueByIncrement(IncrementDirection.Up);
            e.preventDefault();
        } else if (e.keyCode === 40) {
            this.updateValueByIncrement(IncrementDirection.Down);
            e.preventDefault();
        }
    };

    /**
     * 点击向上的箭头
     */
    handleUp = () => {
        this.numberInputRef && this.numberInputRef.current && this.numberInputRef.current.focus()
        this.updateValueByIncrement(IncrementDirection.Up)
    }

    /**
     * 点击向下的箭头
     */
    handleDown = () => {
        this.numberInputRef && this.numberInputRef.current && this.numberInputRef.current.focus()
        this.updateValueByIncrement(IncrementDirection.Down);
    }

    /**
     * 点击上下箭头时，不触发文本框的失焦事件
     */
    iconMouseDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        e.preventDefault();
    }

    /**
     * 按增量更新值
     */
    updateValueByIncrement = (direction: IncrementDirection) => {
        if (this.props.disabled || this.props.readOnly) {
            return false;
        }

        const precision = this.getPrecision();
        const step = this.getStep(precision, this.props.step!);
        if (step === 0) {
            return false;
        }
        const { min, max } = this.props
        const { inputValue } = this.state;
        let val: string
        // 如果不是完整的输入，则不支持进行上下箭头的加减
        if (this.isNotCompleteNumber(inputValue)) {
            if (direction === IncrementDirection.Up && typeof min === 'number' && min !== MIN_SAFE_INTEGER) {
                val = min.toString()
            } else if (direction === IncrementDirection.Down && typeof max === 'number' && max !== MAX_SAFE_INTEGER) {
                val = max.toString()
            } else {
                return false
            }
        } else {
            val = this.getValidValueByRange(Number(inputValue) + step * direction);
        }

        // 当没有设定精度，按步长进行增量，会出现0.1+0.2=0.300000000004的场景，因此根据输入的值和步长中较大的精度作为toFixed的精度，解决该问题
        let finalPrecision = precision
        if (precision === undefined) {
            const inputValueString = inputValue.toString()
            const stepString = step.toString()
            const inputValuePrecision = inputValueString.indexOf('.') === -1 ? 0 : (inputValueString.slice(inputValueString.indexOf('.') + 1)).length // 输入值的精度
            const stepPrecision = stepString.indexOf('.') === -1 ? 0 : (stepString.slice(stepString.indexOf('.') + 1)).length  // 设置的步长的精度
            finalPrecision = Math.max(inputValuePrecision, stepPrecision)
        }

        const increValue = Number(val).toFixed(finalPrecision)

        this.setState(
            {
                inputValue: String(increValue).replace(/\.|\-/g, '').length > 15 ? inputValue : increValue,
            },
            () => {
                this.dispatchValueChangeEvent(this.toNumber(this.state.inputValue));
            },
        );
    };

    /**
     * 获取步进
     */
    private getStep = (precision: number | undefined, step: number) => {
        if (typeof precision === 'number') {
            // 如果 步进step 大于等于 精度precision，则将步进设置为 step
            return typeof step === 'number' && step >= Math.pow(10, Math.round(precision) * -1) ? step : 0// Math.round(precision) 确保参与计算的精度为整数
        } else {
            return typeof step === 'number' && step > EPSILON ? step : 0;
        }
    };

    render() {
        const { autoFocus, disabled, readOnly, placeholder, maxLength, width, step, status, onMouseEnter, onMouseLeave, onClick, className, role } = this.props;
        const { focused, inputValue } = this.state;
        const showArrow = typeof step === 'number' && step !== 0

        return (
            <View
                role={role}
                className={classnames(styles['number-box'], {
                    [styles[`${status}-focused`]]: focused,
                    [styles['disabled']]: disabled,
                    [styles[`${status}`]]: status === 'error',
                    [styles['box-padding']]: showArrow,
                    [styles['read-only']]: readOnly,
                })}
                style={{ width }}
                inline={true}
                onClick={(event) => isFunction(onClick) && onClick(event)}
                {...{ onMouseEnter, onMouseLeave }}
            >
                <input
                    type="text"
                    autoComplete="off"
                    className={classnames(styles['input'], { [styles['disabled']]: disabled }, className)}
                    ref={this.numberInputRef}
                    {...{ autoFocus, disabled, readOnly, placeholder, maxLength }}
                    value={inputValue}
                    onChange={this.handleValueChange}
                    onFocus={this.handleFocus}
                    onBlur={this.handleBlur}
                    onKeyDown={this.handleKeyDown}
                />
                {
                    status === 'error' ?
                        <SweetIcon
                            name={'caution'}
                            size={16}
                            color={'#e60012'}
                            className={classnames(styles['caution'], { [styles['arrow-visiable']]: showArrow })}
                        /> : null
                }
                {
                    // 步进step是不为0的数字时，显示上下箭头
                    showArrow && !disabled ?
                        <View className={styles['icon-wrapper']}>
                            <View onClick={this.handleUp} onMouseDown={this.iconMouseDown} className={styles['icon-item']}>
                                <SweetIcon name={'arrowUp'} size={14} />
                            </View>
                            <View onClick={this.handleDown} onMouseDown={this.iconMouseDown} className={styles['icon-item']}>
                                <SweetIcon name={'arrowDown'} size={14} />
                            </View>
                        </View>
                        : null
                }
            </View >
        );
    }
}