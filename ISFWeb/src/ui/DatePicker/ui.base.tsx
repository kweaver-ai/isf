import React from 'react';
import { noop } from 'lodash';

interface DatePickerProps {
    /**
     * 选中日期时触发
     */
    onChange: (date: Date) => void;

    /**
     * 日历翻页时触发
     */
    onChangeCalendar: Function;

    /**
     * 选择的起始日期范围
     */
    selectRange: Array<Date>;

    /**
     * 选中的日期对象从当天00:00:00开始
     */
    startsFromZero: boolean;

    /**
     * 是否支持选择时间
     */
    showTime: boolean;
}

interface DatePickerState {
    /**
     * 日期面板的状态
     */
    mode: 'time' | 'date';

    /**
     * 当前选中的日期
     */
    value: Date;

    /**
     * 日期面板显示的年份
     */
    year: number;

    /**
     * 日期面板显示的月份
     */
    month: number;
}

export default class DatePickerBase extends React.PureComponent<DatePickerProps, DatePickerState> {
    static defaultProps = {
        onChange: noop,

        onChangeCalendar: noop,

        selectRange: [],

        startsFromZero: false,
    }

    state = {
        value: this.props.value,

        year: (this.props.value || new Date()).getFullYear(),

        month: (this.props.value || new Date()).getMonth() + 1,

        mode: 'date',
    }

    static getDerivedStateFromProps({ disabled, value }, prevState) {
        if (!disabled && value && (!prevState.value || (prevState.value && value.getTime() !== prevState.value.getTime()))) {
            return {
                value,
                year: value.getFullYear(),
                month: value.getMonth() + 1,
            }
        }
        return null
    }

    flipYear(yearChange) {
        if (this.state.year + yearChange < 1970) {
            return;
        }
        this.setState({
            year: this.state.year + yearChange,
        })

        this.fireChangeCalendarEvent();
    }

    flipMonth(monthChange) {

        const nextMonth = this.state.month + monthChange;

        if (nextMonth > 12) {
            this.flipYear(1);

            this.setState({
                month: 1,
            })
        }
        else if (nextMonth < 1) {

            if (this.state.year === 1970) {
                return;
            }

            this.flipYear(-1);

            this.setState({
                month: 12,
            })
        } else {
            this.setState({
                month: this.state.month + monthChange,
            })
        }

        this.fireChangeCalendarEvent();
    }

    fireChangeCalendarEvent() {
        this.props.onChangeCalendar();
    }

    protected handleClickTime = () => {
        if(!this.props.disabled) {
            this.setState({
                mode: 'time',
            })
        }
    }

    protected handleReturnDatePanel = () => {
        this.setState({
            mode: 'date',
        })
    }

    /**
     * 选中时间点时触发
     */
    protected handleSelectTime = (selectedDate) => {
        this.setState({
            value: selectedDate,
        })
        this.props.onChange(selectedDate)
        this.handleReturnDatePanel()
    }

    /**
     * 选中日期时触发
     */
    protected handleSelectDate = (selectedDate) => {
        if(this.props.showTime) {
            const { selectRange: [start, end] } = this.props;
            const { value } = this.state;
            let hours = value.getHours()

            if(start && selectedDate.toDateString() === start.toDateString()) {
                // 选中日期是开始日期，忽略时间点
                hours = Math.max(start.getMinutes() > 0 ? start.getHours() + 1 : start.getHours(), hours)
            }
            if(end && selectedDate.toDateString() === end.toDateString()) {
                // 选中日期是结束日期，忽略时间点
                hours = Math.min(end.getHours(), hours)
            }

            const date = new Date(
                selectedDate.getFullYear(),
                selectedDate.getMonth(),
                selectedDate.getDate(),
                hours > 23 ? 23 : hours,
                hours < 23 ? 0 : value.getMinutes(),
                0,
            );

            this.setState({
                value: date,
            })
            this.props.onChange(date)
        } else {
            this.props.onChange(selectedDate)
        }
    }
}