import { chunk, first, last } from 'lodash';

/**
 * 根据年／月构建当月所有日期的日期对象
 * @param year 年
 * @param month 月
 */
// eslint-disable-next-line
export function generateDaysOfMonth(year: number, month: number, { startsFromZero = false } = {} as { startsFromZero: boolean }): ReadonlyArray<Date> {
    let days: ReadonlyArray<Date> = [];
    let date = 1;
    let dateOfMonth = startsFromZero ? new Date(year, month - 1, date, 0, 0, 0) : new Date(year, month - 1, date, 23, 59, 59);

    do {
        days = [...days, dateOfMonth];
        dateOfMonth = startsFromZero ? new Date(year, month - 1, ++date, 0, 0, 0) : new Date(year, month - 1, ++date, 23, 59, 59);
    } while (dateOfMonth.getMonth() + 1 === month);

    return days;
}

/**
 * 构造星期数组
 * @params firstOfDay 一周的开始
 * @return 返回代表days的数组
 */
export function generateWeekDays(firstOfDay = 0): ReadonlyArray<number> {
    let day = firstOfDay;
    let ret: ReadonlyArray<number> = [];

    while (ret.length < 7) {

        ret = [...ret, day]

        if (day === 6) {
            day = 0;
        } else {
            day++;
        }
    }

    return ret;
}

/**
 * 构建日期Matrix
 * @param year 年
 * @param month 月
 * @param [firstOfDay] 从周几开始，0代表周日，依次递增
 */
// eslint-disable-next-line
export function generateWeeksOfMonth(year: number, month: number, firstOfDay: number = 0, { startsFromZero = false } = {} as { startsFromZero: boolean }): Array<Array<Date>> {
    const week = generateWeekDays(firstOfDay);
    const dates = generateDaysOfMonth(year, month, { startsFromZero });
    const dateStartIndex = week.indexOf(first(dates).getDay());
    const dateEndIndex = week.indexOf(last(dates).getDay());
    const datesOfPrevMonth = new Array(dateStartIndex);
    const datesOfNextMonth = new Array(7 - (dateEndIndex + 1));

    return chunk(datesOfPrevMonth.concat(dates, datesOfNextMonth), 7);
}

/**
 * 从日期字符串转换毫秒数
 * @params dateString 日期字符串
 */
export function getUTCTime(dateString: string) {
    let y, M, d, h, m, s;

    // (ISO 8601标准) 例：2017-12-14 , 2017-12-11T14:50:55+08:00
    if (dateString.match(/^\d{4}(-?\d{2}){2}([\sT]\d{2}:\d{2}:\d{2}([\+\-\s]\d{2}:\d{2})?)?/)) {
        const { date = '1970-01-01', time = '00:00:00', zone = '+00:00' } = dateString.match(/^(\d{4}(-?\d{2}){2})|([\sT]\d{2}:\d{2}:\d{2})|([\+\-\s]\d{2}:\d{2})/g).reduce((prev, currentValue, index) => {
            return {
                ...prev,
                date: currentValue.match(/\d{4}(-?\d{2}){2}/) ? currentValue : prev['date'],
                time: currentValue.match(/[\sT]\d{2}:\d{2}:\d{2}/) ? currentValue : prev['time'],
                zone: currentValue.match(/[\+\-\s]\d{2}:\d{2}/) ? currentValue : prev['zone'],
            }
        }, {})
        // zone指定时区，可以是：Z (UTC)、+hh:mm、-hh:mm
        const [hh, mm] = zone.split(':');
        const [, t = '00:00:00'] = time.split(/[\sT]/);

        [y = 0, M = 0, d = 0] = date.split('-').map(Number);
        [h = 0, m = 0, s = 0] = t.split(':').map(Number);
        h = h - Number(hh);
        m = m - Number(mm);
    } else {
        let [fullDate, time] = dateString.split(/\s+/);

        [h = 0, m = 0, s = 0] = time ? time.split(':').map(Number) : [];
        if (fullDate.match(/^\d{1,2}\/\d{1,2}\/\d{4}$/)) {
            [M = 0, d = 0, y = 0] = fullDate.split('/').map(Number);
        } else if (fullDate.match(/\d{4}(-\d{1,2}){2}/)) {
            [y = 0, M = 0, d = 0] = fullDate.split('-').map(Number);
        } else if (fullDate.match(/\d{4}(\.\d{1,2}){2}/)) {
            [y = 0, M = 0, d = 0] = fullDate.split('.').map(Number);
        }
    }

    return Date.UTC(y, M - 1, d, h, m, s);
}

/**
 * 获取某一天的00:00:00
 * @param date
 * @param type UTC | GMT
 */
// eslint-disable-next-line
export function startOfDay(date: Date = new Date(), { type = 'UTC' } = {} as { type: string }) {
    switch (type) {
        case 'GMT':
            return new Date(date.getFullYear(), date.getMonth(), date.getDate(), 0, 0, 0, 0).getTime()
        case 'UTC':
        default:
            return Date.UTC(date.getFullYear(), date.getMonth(), date.getDate(), 0, 0, 0, 0);
    }
}

/**
 * 获取某一天的23:59:59
 * @param date
 * @param type UTC | GMT
 */
// eslint-disable-next-line
export function endOfDay(date: Date = new Date(), { type = 'UTC' } = {} as { type: string }) {
    switch (type) {
        case 'GMT':
            return new Date(date.getFullYear(), date.getMonth(), date.getDate(), 23, 59, 59, 999).getTime();
        case 'UTC':
        default:
            return Date.UTC(date.getFullYear(), date.getMonth(), date.getDate(), 23, 59, 59, 999);
    }

}

/**
 * 计算时钟
 * @param totalSeconds 总秒数
 */
export function clock(totalSeconds: number) {
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds - hours * 3600) / 60);
    const seconds = totalSeconds - hours * 3600 - minutes * 60;

    return { hours, minutes, seconds }
}

/**
 * 构件今天日期对象，设定为当天23:59:59
 */
export function today(): Date {
    const now = new Date();
    return new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59);
}
