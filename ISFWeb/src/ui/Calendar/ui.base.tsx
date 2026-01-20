// eslint-disable-next-line spaced-comment
/// <reference path="./ui.base.d.ts" />

import React from 'react';
import { noop } from 'lodash';
import { generateWeeksOfMonth, startOfDay } from '@/util/date';

export default class CalendarBase extends React.PureComponent<UI.Calendar.Props, any> implements UI.Calendar.Base {
    static defaultProps = {
        year: new Date().getFullYear(),

        month: new Date().getMonth() + 1,

        date: new Date(),

        onSelect: noop,

        selectRange: [], // 允许选择的日期范围，数组项1为起始日起，数组项2为结束日起

        startsFromZero: false, // 选中的日期对象从当天00:00:00开始
    }

    state = {
        year: this.props.year,

        month: this.props.month,

        weeks: generateWeeksOfMonth(this.props.year, this.props.month, this.props.firstOfDay, { startsFromZero: this.props.startsFromZero }) || [],
    }

    today = startOfDay(new Date())

    static getDerivedStateFromProps({ year, month }, prevState) {
        if (year !== prevState.year || month !== prevState.month) {
            return {
                year,
                month,
            }
        }
        return null
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState.year !== this.state.year || prevState.month !== this.state.month) {
            this.updateCalendar(this.state.year, this.state.month)
        }
    }

    /**
     * 判断两日期是否是同一天
     * @param a 待比较的日期
     * @param b 待比较的日期
     * @return 返回是否是同一天
     */
    matchSelected(date: Date): boolean {
        return date
            && this.props.select
            && date.getFullYear() === this.props.select.getFullYear()
            && date.getMonth() + 1 === this.props.select.getMonth() + 1
            && date.getDate() === this.props.select.getDate()
    }

    /**
     * 更新日历月份
     * @param year 年份
     * @param month 月份
     */
    updateCalendar(year: string & number, month: string & number) {
        this.setState({
            weeks: generateWeeksOfMonth(year, month, this.props.firstOfDay, { startsFromZero: this.props.startsFromZero }),
        })
    }

    clickHandler(date) {
        this.props.onSelect(date);
    }
}