import * as React from 'react';
import { noop } from 'lodash';
import { formatTime } from '@/util/formatters';
import { today } from '@/util/date';
import WebComponent from '../webcomponent';
import __ from './locale';

interface Props {
    /**
     * 控件宽度
     */
    width?: number;

    /**
     * 初始日期微秒
     */
    value?: number;

    /**
     * 对齐方式
     * @example
     * ```
     * <ValidityBox2 dropAlign="left bottom" />
     * ```
     */
    dropAlign?: string;

    /**
     * 是否允许永久有效
     */
    allowPermanent?: boolean;

    /**
     * 日期选择范围
     * @example
     * ```
     *  <ValidityBox2 selectRange=[start, end] />
     *  <ValidityBox2 selectRange=[, end] />
     *  <ValidityBox2 selectRange=[start, ] />
     * ```
     */
    selectRange?: ReadonlyArray<Date>;

    /**
     * 默认选中从当前算起后的日期
     */
    defaultSelect?: number;

    /**
     * 选项改变时触发
     * @param date 日期对象
     */
    onChange?: (date: Date) => any;

    /**
     * 自定义样式
     */
    className?: string;
}

interface State {
    /**
     * 当前选中的日期微秒
     */
    value?: number;

    /**
     * 当前是否显示日历控件
     */
    active?: boolean;
}

export default class ValidityBox2Base extends WebComponent<Props, State> {
    static defaultProps = {
        onChange: noop,

        dropAlign: 'bottom left',

        allowPermanent: false,

        selectRange: [new Date()],

        defaultSelect: 30,

        disabled: false,
    }

    state = {
        value: this.props.value || -1,
    }

    select = null

    deactivePrevented: boolean = false  // 是否选中了Menu区域

    availableSelected: boolean = false  // 有效选中

    popOverClose = noop  // Popover的关闭函数

    componentDidMount() {
        this.setState({
            value: this.props.value || -1,
        })
    }

    static getDerivedStateFromProps({ value }, prevState){
        if(value !== prevState.value){
            return { value }
        }else{
            return null
        }
    }

    /**
     * 更新数据值
     * value 有效期长度，微秒单位
     */
    private updateValue(value: number, { active = false } = {}) {
        this.setState({ value, active });
        this.fireChangeEvent(value);
    }

    /**
     * 设定选中的日期
     * @param value 日期对象
     * @param close 弹出层关闭方法
     */
    protected setValidity(value: Date, close?) {
        this.updateValue(value.getTime() * 1000);
        // 选中了一个可选范围的日期
        close()
    }

    /**
     * 切换勾选永久有效
     * @param close 弹出层关闭方法
     * @param permanent 是否勾选永久有效
     * @param permVal 用来表示永久有效的value值
     */
    protected switchPermanent(close, permanent: boolean, permVal: number): void {
        if (permanent) {
            this.updateValue(permVal);
            // 勾选永久有效，需收起日期的popover
            close()
        } else {
            let defaultDate = today().setDate(today().getDate() + this.props.defaultSelect);
            this.updateValue(defaultDate * 1000, { active: true });
        }
    }

    /**
     * 格式化选中值
     * @param date 选中的日期
     */
    protected validityFormatter(date: number): string {
        if (date === -1) {
            return __('永久有效');
        } else {
            return formatTime(date, 'yyyy-MM-dd')
        }
    }

    /**
     * 触发选中事件
     * @param value 日期微秒
     */
    private fireChangeEvent(value: number): void {
        this.props.onChange(value);
    }

    /**
     * 切换选中状态
     */
    protected toggleActive() {
        this.setState({
            active: !this.state.active,
        })
    }

    protected onSelectBlur(e) {
        if (this.deactivePrevented) {
            // 当日期下拉出现时，select需要focuse
            this.select.focus();
        } else {
            this.setState({
                active: false,
            })
        }

        this.deactivePrevented = false;
    }
}