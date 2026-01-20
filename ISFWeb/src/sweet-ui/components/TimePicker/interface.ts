import moment from 'moment';

/**
 * 时间值类型
 */
export type DateType = moment.Moment;

export interface SharedProps {
    /**
     * 小时选项间隔
     */
    hourStep?: number;

    /**
     * 分钟选项间隔
     */
    minuteStep?: number;

    /**
     * 秒选项间隔
     */
    secondStep?: number;

    /**
     * 禁止选择部分小时
     */
    disabledHours?: () => Array<number>;

    /**
     * 禁止选择部分分钟
     */
    disabledMinutes?: (selectedHour: number) => Array<number>;

    /**
     * 禁止选择部分秒
     */
    disabledSeconds?: (selectedHour: number, selectedMinute: number) => Array<number>;
}