import React from 'react';
import FlexBox from '../FlexBox/ui.desktop';
import Calendar from '../Calendar/ui.desktop';
import UIIcon from '../UIIcon/ui.desktop';
import TimePanel from './Time/ui.desktop';
import DatePickerBase from './ui.base';
import styles from './styles.desktop';
import prevmonth from './assets/prevmonth.png';
import nextmonth from './assets/nextmonth.png';
import prevyear from './assets/prevyear.png';
import nextyear from './assets/nextyear.png';

const DefaultTime = '23:59';

export default class DatePicker extends DatePickerBase {
    render() {
        const { mode, value } = this.state;
        const [start, end] = this.props.selectRange;

        if (mode === 'time') {
            return (
                <TimePanel
                    selectDate={this.state.value}
                    {...{ start, end }}
                    onClickTime={this.handleReturnDatePanel}
                    onSelectTime={this.handleSelectTime}
                />
            );
        }

        return (
            <div
                role={this.props.role}
                className={styles.container}
                onMouseDown={this.props.onDatePickerClick}
            >
                <div className={styles.monthPanel}>
                    <FlexBox>
                        <FlexBox.Item align="left middle">
                            <span className={styles.navWrapper}>
                                <UIIcon
                                    disabled={this.props.disabled}
                                    size="12"
                                    code={'\uf010'}
                                    fallback={prevyear}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        this.flipYear(-1);
                                    }}
                                />
                            </span>
                            <span className={styles.navWrapper}>
                                <UIIcon
                                    disabled={this.props.disabled}
                                    size="12"
                                    code={'\uf012'}
                                    fallback={prevmonth}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        this.flipMonth(-1);
                                    }}
                                />
                            </span>
                        </FlexBox.Item>
                        <FlexBox.Item align="center middle">{`${this.state.year}-${this.state.month}`}</FlexBox.Item>
                        <FlexBox.Item align="right middle">
                            <span className={styles.navWrapper}>
                                <UIIcon
                                    disabled={this.props.disabled}
                                    size="12"
                                    code={'\uf011'}
                                    fallback={nextmonth}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        this.flipMonth(1);
                                    }}
                                />
                            </span>
                            <span className={styles.navWrapper}>
                                <UIIcon
                                    disabled={this.props.disabled}
                                    size="12"
                                    code={'\uf00f'}
                                    fallback={nextyear}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        this.flipYear(1);
                                    }}
                                />
                            </span>
                        </FlexBox.Item>
                    </FlexBox>
                </div>
                <div className={styles.calendar}>
                    <Calendar
                        disabled={this.props.disabled}
                        selectRange={this.props.selectRange}
                        year={this.state.year}
                        month={this.state.month}
                        select={this.props.value}
                        startsFromZero={this.props.startsFromZero}
                        onSelect={this.handleSelectDate}
                        time={this.props.showTime ?
                            (value ? `${value.getHours() < 10 ? `0${value.getHours()}` : value.getHours()}
                            : ${value.getMinutes() < 10 ? `0${value.getMinutes()}` : value.getMinutes()}`
                                : DefaultTime)
                            : null}
                        onClickTime={this.handleClickTime}
                    />
                </div>
            </div>
        );
    }
}
