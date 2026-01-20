import moment from 'moment';
import { DateType } from './interface';

export const generateConfig = {
    // get
    getNow: () => moment(),
    getHour: (date: DateType) => date.hour(),
    getMinute: (date: DateType) => date.minute(),
    getSecond: (date: DateType) => date.second(),

    // set
    setHour: (date: DateType, hour: number) => {
        const clone = date.clone();
        clone.hour(hour);
        return clone;
    },
    setMinute: (date: DateType, minute: number) => {
        const clone = date.clone();
        clone.minute(minute);
        return clone;
    },
    setSecond: (date: DateType, second: number) => {
        const clone = date.clone();
        clone.second(second);
        return clone;
    },

    format: (date: DateType, format: string) => {
        const clone = date.clone();
        return clone.format(format);
    },

    parse: (text: string, formats: string[]) => {
        const fallbackFormatList: string[] = [];

        for (let i = 0; i < formats.length; i += 1) {
            let format = formats[i];
            let formatText = text;

            const date = moment(formatText, format, true);
            if (date.isValid()) {
                return date;
            }
        }

        for (let i = 0; i < fallbackFormatList.length; i += 1) {
            const date = moment(text, fallbackFormatList[i], false);

            if (date.isValid()) {
                return date;
            }
        }

        return null;
    },
};

export function leftPad(str: string | number, length: number, fill: string = '0') {
    let current = String(str);
    if (current.length < length) {
        current = `${fill}${str}`;
    }
    return current;
}

const scrollIds = new Map<HTMLElement, number>();

export function scrollTo(element: HTMLElement, to: number, duration: number) {
    if (scrollIds.get(element)) {
        cancelAnimationFrame(scrollIds.get(element)!);
    }

    if (duration <= 0) {
        scrollIds.set(
            element,
            requestAnimationFrame(() => {
                element.scrollTop = to;
            }),
        );

        return;
    }
    const difference = to - element.scrollTop;
    const perTick = difference / duration * 10;

    scrollIds.set(
        element,
        requestAnimationFrame(() => {
            element.scrollTop += perTick;
            if (element.scrollTop !== to) {
                scrollTo(element, to, duration - 10);
            }
        }),
    );
}

/**
 * 获取今天的0时0分0秒
 */
export function getStartOfDay() {
    return moment().startOf('day');
}

/**
 * 根据时间格式返回是否显示时、分、秒
 * @param format 时间格式
 */
export function generateShowHourMinuteSecond(format: string) {
    return {
        showHour: format.indexOf('H') > -1 || format.indexOf('h') > -1,
        showMinute: format.indexOf('m') > -1,
        showSecond: format.indexOf('s') > -1,
    };
}
