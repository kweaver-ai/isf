import React from 'react';
import classnames from 'classnames';
import { generateWeekDays, endOfDay, startOfDay } from '@/util/date';
import LinkChip from '../LinkChip/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop'
import CalendarBase from './ui.base';
import { getLocaleWeeks } from './helper';
import styles from './styles.desktop';
import __ from './locale';

export default class Calendar extends CalendarBase {
    render() {
        const [start, end] = this.props.selectRange;
        const [allowFrom, allowTo] = [start ? startOfDay(start, { type: 'GMT' }) : start, end ? endOfDay(end, { type: 'GMT' }) : end];
        const localeWeeks = getLocaleWeeks()

        return (
            <div className={styles.calendar}>
                <table>
                    <thead>
                        {
                            <tr>
                                {
                                    generateWeekDays(this.props.firstOfDay).map((day) => {
                                        return (
                                            <th className={styles.day} key={day}>
                                                {
                                                    localeWeeks[day]
                                                }
                                            </th>
                                        )
                                    })
                                }
                            </tr>
                        }
                    </thead>
                    <tbody>
                        {
                            this.state.weeks.map((week, index) => {
                                return (
                                    <tr key={index}>
                                        {
                                            week.map((date, index) => {
                                                return (
                                                    <td
                                                        key={index}
                                                        className={classnames(styles.cell, { [styles.selected]: this.matchSelected(date), [styles['invalid']]: !date })}>
                                                        {
                                                            date ?
                                                                <LinkChip
                                                                    key={date}
                                                                    className={styles.date}
                                                                    disabled={this.props.disabled || (allowFrom && date < allowFrom) || (allowTo && date > allowTo)}
                                                                    onClick={() => this.clickHandler(date)}
                                                                >
                                                                    {
                                                                        date.getDate()
                                                                    }
                                                                </LinkChip>
                                                                : null
                                                        }
                                                    </td>
                                                )
                                            })
                                        }
                                    </tr>
                                )
                            })
                        }
                        {
                            this.props.time ?
                                <tr className={styles['time-layout']}>
                                    <td colSpan={7} className={styles['time-padding']}>
                                        <UIIcon
                                            disabled={this.props.disabled}
                                            size="16"
                                            code={'\uf022'}
                                            color={'#757575'}
                                            className={styles['time-icon']}
                                            onClick={this.props.onClickTime}
                                        />
                                        <span
                                            className={classnames(styles['time-text'], { [styles['disabled']]: this.props.disabled })}
                                            onClick={this.props.onClickTime}
                                        >
                                            {__('选择时间点：')}
                                            <span className={styles['time']}>
                                                {this.props.time}
                                            </span>
                                        </span>
                                    </td>
                                </tr>
                                : null
                        }
                    </tbody>
                </table>
            </div>
        )
    }
}