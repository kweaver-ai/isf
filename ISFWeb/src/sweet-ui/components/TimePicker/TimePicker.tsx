import React from 'react';
import classnames from 'classnames';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import TextInput from '../TextInput';
import Trigger from '../Trigger';
import View from '../View';
import SweetIcon from '../SweetIcon';
import Panel from './Panel';
import { generateConfig, getStartOfDay, generateShowHourMinuteSecond } from './util';
import { DateType, SharedProps } from './interface';
import styles from './styles';

interface TimePickerProps extends SharedProps {
    /**
     * 当前时间
     */
    value?: DateType;

    /**
     * 宽度ing
     */
    width?: number | string;

    /**
     * 默认时间
     */
    defaultValue?: DateType;

    /**
     * 打开面板时默认选中时间，当value/defaultValue都不存在时使用
     */
    defaultOpenValue?: DateType;

    /**
     * 显示的时间格式
     */
    format?: string;

    /**
     * 禁用全部操作
     */
    disabled?: boolean;

    /**
     * 是否显示清除按钮
     */
    allowClear?: boolean;

    /**
     * 提示文字
     */
    placeholder?: string;

    /**
     * 选择时间发生变化的回调
     */
    onValueChange?: (event: SweetUIEvent<{ date: DateType; timeString: string }>) => void;
}

interface TimePickerState {
    /**
     * 当前时间
     */
    value: DateType;

    /**
     * 输入框输入值
     */
    timeInput: string;

    /**
     * 时间面板的打开状态
     */
    open: boolean;
}

export default class TimePicker extends React.Component<TimePickerProps, TimePickerState> {
    static defaultProps = {
        allowClear: true,
        defaultOpenValue: getStartOfDay(),
        format: 'HH:mm:ss',
        hourStep: 1,
        minuteStep: 1,
        secondStep: 1,
    };

    constructor(props: TimePickerProps, ...args: any) {
        super(props);
        const { defaultValue, defaultOpenValue, value, format } = props;

        this.state = {
            value: value || defaultValue || defaultOpenValue,
            timeInput: value || defaultValue ? generateConfig.format(value || defaultValue, format) : '',
            open: false,
        };
    }

    timeInput: TextInput | null = null;

    static getDerivedStateFromProps({ value, format }: TimePickerProps, prevState: TimePickerState) {
        if (value && value !== prevState.value) {
            return {
                value,
                timeInput: generateConfig.format(value, format),
            };
        }
        return null;
    }

    /**
     * 输入内容改变时触发
     */
    handleInputValueChange = (input: string) => {
        this.setState({
            timeInput: input,
        });
        const date =
            input.trim() === ''
                ? this.props.defaultValue || this.props.defaultOpenValue
                : this.getValueFromFormatString(input, this.props.format);
        // 输入格式合法，更新value
        if (date && this.isValidTime(date)) {
            this.setState({
                value: date,
            });

            this.dispatchValueChangeEvent({ date, timeString: input });
        }
    };

    isValidTime = (date: DateType) => {
        const {
            disabledHours,
            disabledMinutes,
            disabledSeconds,
            hourStep = 1,
            minuteStep = 1,
            secondStep = 1,
        } = this.props;

        if (date) {
            const hour = generateConfig.getHour(date);
            const minute = generateConfig.getMinute(date);
            const second = generateConfig.getSecond(date);

            const disabledHourOptions = typeof disabledHours === 'function' && disabledHours();
            const disabledMinuteOptions = typeof disabledMinutes === 'function' && disabledMinutes(hour);
            const disabledSecondOptions = typeof disabledSeconds === 'function' && disabledSeconds(hour, minute);

            if (
                (disabledHourOptions && disabledHourOptions.indexOf(hour) >= 0) ||
                (disabledMinuteOptions && disabledMinuteOptions.indexOf(minute) >= 0) ||
                (disabledSecondOptions && disabledSecondOptions.indexOf(second) >= 0) ||
                hour % hourStep !== 0 ||
                minute % minuteStep !== 0 ||
                second % secondStep !== 0
            ) {
                return false;
            }
            return true;
        }
        return false;
    };

    /**
     * 清除输入框内容
     */
    handleClearInput = () => {
        const { defaultValue, defaultOpenValue } = this.props;

        this.setState({
            value: defaultValue || defaultOpenValue,
            timeInput: '',
        });

        this.dispatchValueChangeEvent({
            date: defaultValue || defaultOpenValue,
            timeString: this.getFormatStringFromValue(defaultValue || defaultOpenValue, this.props.format),
        });
    };

    /**
     * 面板打开/关闭时触发
     */
    handlePopupVisibleChange = (open: boolean) => {
        // 关闭时验证输入合法性
        if (!open) {
            const { defaultValue, format } = this.props;

            if (this.state.timeInput.trim() === '') {
                this.setState({
                    timeInput: defaultValue ? this.getFormatStringFromValue(defaultValue, format) : '',
                });
            } else {
                const date = this.getValueFromFormatString(this.state.timeInput, format);

                if (!date || !this.isValidTime(date)) {
                    // 输入格式不合法，还原更改前的输入值
                    this.setState({
                        timeInput: this.getFormatStringFromValue(this.state.value, format),
                    });
                }
            }

            // 手动失焦输入框
            this.timeInput && this.timeInput.blur();
        }

        this.setState({
            open,
        });
    };

    /**
     * 点击时间面板里的选项触发
     */
    handleSelectTime = (time: DateType) => {
        const timeString = this.getFormatStringFromValue(time, this.props.format);
        this.setState({
            value: time,
            timeInput: timeString,
        });
        this.dispatchValueChangeEvent({ date: time, timeString });
    };

    /**
     * time value转化为格式化字符串
     */
    getFormatStringFromValue = (date: DateType, format: string) => {
        return generateConfig.format(date, format);
    };

    /**
     * 格式化字符串转化为time value
     */
    getValueFromFormatString = (timeString: string, format: string) => {
        return generateConfig.parse(timeString, [format]);
    };

    saveTimeInput = (timeInput: TextInput) => {
        this.timeInput = timeInput;
    };

    dispatchValueChangeEvent = createEventDispatcher(this.props.onValueChange);

    renderClearButton() {
        const { value } = this.state;
        const { disabled } = this.props;
        if (!value || disabled) {
            return null;
        }

        return (
            <SweetIcon name={'clear'} className={classnames(styles['clear-button'])} onClick={this.handleClearInput} />
        );
    }

    render() {
        const {
            disabled,
            width = 200,
            placeholder,
            format,
            allowClear,
            hourStep,
            minuteStep,
            secondStep,
            disabledHours,
            disabledMinutes,
            disabledSeconds,
        } = this.props;
        const { value, timeInput, open } = this.state;

        const { showHour, showMinute, showSecond } = generateShowHourMinuteSecond(format);

        const inputElement = ({
            setPopupVisibleOnFocus,
        }: Partial<{
            setPopupVisibleOnClick: () => void;
            setPopupVisibleOnFocus: () => void;
            setPopupVisibleOnBlur: () => void;
        }>) => (
            <View
                key={'timePickerTrigger'}
                className={classnames(
                    styles['time-picker'],
                    { [styles['disabled']]: disabled },
                    { [styles['active']]: open },
                )}
                inline={true}
                style={{ width }}
            >
                <TextInput
                    disabled={disabled}
                    width={'100%'}
                    className={classnames(styles['time-input'])}
                    ref={this.saveTimeInput}
                    value={timeInput}
                    selectOnFocus={true}
                    placeholder={placeholder}
                    onFocus={setPopupVisibleOnFocus}
                    onValueChange={({ detail }) => this.handleInputValueChange(detail)}
                />
                {allowClear ? this.renderClearButton() : null}
                <SweetIcon
                    name={'time'}
                    className={classnames(styles['time-icon'], { [styles['disabled']]: disabled })}
                />
            </View>
        );

        return (
            <Trigger
                renderer={inputElement}
                freeze={false}
                onPopupVisibleChange={({ detail }) => this.handlePopupVisibleChange(detail)}
                triggerEvent={'focus'}
            >
                {() => (
                    <View className={styles['time-panel']}>
                        <Panel
                            {...{
                                value,
                                open,
                                showHour,
                                showMinute,
                                showSecond,
                                hourStep,
                                minuteStep,
                                secondStep,
                                disabledHours,
                                disabledMinutes,
                                disabledSeconds,
                            }}
                            onSelect={this.handleSelectTime}
                        />
                    </View>
                )}
            </Trigger>
        );
    }
}
