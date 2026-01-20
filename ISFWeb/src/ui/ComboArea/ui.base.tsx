import React from 'react';
import { last, uniqBy, filter, includes, isFunction, noop } from 'lodash';
import { mapKeyCode } from '@/util/browser';

interface Props extends React.ClassAttributes<void> {
    /**
     * className
     */
    className?: string;

    /**
     * 宽度，含padding和border
     */
    width?: number | string;

    /**
     * 高度，含padding和border
     */
    height?: number | string;

    /**
     * 最小高度
     */
    minHeight?: number;

    /**
     * 最大允许高度
     */
    maxHeight?: number;

    /**
     * 临时解决方案，使用FlexTextBox2
     */
    useNewFlextInput?: boolean;

    /**
     * 最大允许宽度
     */
    maxWidth?: number;

    /**
     * 只读
     */
    readOnly?: boolean;

    /**
     * 是否不可编辑
     * 不可编辑状态下仍然可以删除已有项
     */
    uneditable?: boolean;

    /**
     * 禁用
     */
    disabled?: boolean;

    /**
     * 占位文本
     */
    placeholder?: string;

    /**
     * 自动创建Chip的分割字符
     */
    spliter?: Array<string>;

    /**
     * 文本框中的输入值
     */
    inputValue?: string;

    /**
     * 初始数据
     */
    value?: Array<any>;

    /**
     * 校验状态
     */
    validateState?: number | string;

    /**
     * 校验提示语
     */
    validateMessages?: ReadonlyArray<string> | object;

    /**
     * 格式化函数
     */
    formatter?: (o: any) => string;

    /**
     * 验证是否可以自动创建Chip
     * @params input 输入的值
     * @params data 已生成的数据
     * @return 返回验证是否通过
     */
    validator?: (input: string, data: Array<any>) => boolean;

    /**
     * 数据改变时触发
     */
    onChange?: (data: ReadonlyArray<any>) => any;

    /**
     * 输入框值发生改变时触发
     * @param value 输入值
     */
    onInputValueChange?: (value: string) => void;

    /**
     * 通过keydown增加chip时触发
     */
    onRequestAddChipBykeyDown?: ({ value, isValid }: { value: string; isValid: boolean }) => void;
}

interface State {
    /**
     * 数据
     */
    value: Array<any>;

    /**
     * 显示输入框为空的提示语
     */
    isHoverWarning: boolean;
}

export default class ComboAreaBase extends React.PureComponent<Props, State> {
    static defaultProps = {
        validateMessages: [],

        value: [],

        readOnly: false,

        uneditable: false,

        disabled: false,

        minHeight: 50,

        maxHeight: 100,

        onChange: noop,

        placeholder: '',

        spliter: [],

        formatter: (val) => val,

        validator: (_val) => true,
    }

    state: State = {
        value: this.props.value || [],
        isHoverWarning: false,
    }

    /**
     * 输入区域ref
     */
    input = null

    componentDidUpdate(prevProps, prevState) {
        if (this.props.value && this.props.value !== prevProps.value && this.props.value !== prevState.value) {
            this.setState({
                value: this.props.value,
            })
        }
    }

    private triggerChange = () => {
        isFunction(this.props.onChange) && this.props.onChange(this.state.value);
    }

    private addChip(chips) {
        this.setState({
            value: uniqBy(this.state.value.concat(chips)),
        }, this.triggerChange)
    }

    protected removeChip(chip) {
        this.setState({
            value: filter(this.state.value, (o) => o !== chip),
        }, this.triggerChange);
    }

    protected focusInput() {
        if (!this.props.uneditable) {
            if (this.input && typeof (this.input.focus) === 'function') {
                this.input.focus();
            }
        }
    }

    protected pasteHandler(value: string) {
        setTimeout(() => {
            this.batchAddChips(value);
        })
    }

    /**
     * 批量添加Chips
     */
    private batchAddChips(value) {
        const chips = value
            .split(new RegExp(this.props.spliter.join('|')))
            .filter((chip) => this.props.validator(chip, this.state.value));

        this.addChip(chips);
        this.clearInput();
    }

    /**
     * 验证输入值并添加
     * @param value 输入值
     */
    protected validateInput(value: string, { isKeyDown = false } = {}): void {
        let isValid = false

        if (value && this.props.validator(value, this.state.value)) {
            this.addChip(value);
            this.clearInput();

            isValid = true
        }

        if (isKeyDown) {
            value && this.props.onRequestAddChipBykeyDown && this.props.onRequestAddChipBykeyDown({ value, isValid })
        }
    }

    protected blurHandler(e: React.FocusEvent<HTMLAnchorElement>) {
        this.validateInput(this.state.inputValue)
        isFunction(this.props.onBlur) && this.props.onBlur();
    }

    protected keyDownHandler(e) {
        const input = this.state.inputValue

        // 增加Chip
        if (includes(this.props.spliter, mapKeyCode(e.keyCode)) || e.keyCode === 13) {

            this.validateInput(input, { isKeyDown: true })

            e.preventDefault ? e.preventDefault() : (e.returnValue = false);

        }
        // 删除输入或Chip
        else if (e.keyCode === 8) {
            if (!input) {
                if (this.state.value.length) {
                    this.removeChip(last(this.state.value));
                }
            }
        }
    }

    /**
     * 处理文本框输入值变化
     * @param value
     */
    protected handleValueChange = (value: string) => {
        this.setState({ inputValue: value }, () => {
            if (isFunction(this.props.onInputValueChange)) {
                this.props.onInputValueChange(value)
            }
        })
    }

    private clearInput() {
        this.setState({
            inputValue: '',
        })
        if (this.input && typeof (this.input.clear) === 'function') {
            this.input.clear();
        }
    }

    protected blur() {
        if (this.input && typeof (this.input.blur) === 'function') {
            this.input.blur();
        }
    }

    protected saveFlexInput = (input) => {
        this.input = input
    }

    /**
     * 鼠标悬浮
     */
    protected mouseOver = () => {
        this.setState({
            isHoverWarning: true,
        })
    }

    /**
     * 鼠标移除
     */
    protected mouseOut = () => {
        this.setState({
            isHoverWarning: false,
        })
    }
}