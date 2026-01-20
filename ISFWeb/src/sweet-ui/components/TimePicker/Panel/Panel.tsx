import React from 'react';
import { DateType, SharedProps } from '../interface';
import { generateConfig, leftPad } from '../util';
import TimeUnitColumn, { Unit } from './TimeUnitColumn';
import styles from './styles';

function generateTimeUnits(start: number, end: number, step: number, disabledUnits: Array<number> | undefined) {
    const units: Unit[] = [];
    for (let i = start; i <= end; i += step) {
        // eslint-disable-next-line
        units.push({
            label: leftPad(i, 2),
            value: i,
            disabled: (disabledUnits || []).includes(i),
        });
    }
    return units;
}

interface PanelProps extends SharedProps {
    open?: boolean;
    format?: string;
    showHour?: boolean;
    showMinute?: boolean;
    showSecond?: boolean;
    hideDisabledOptions?: boolean;
    defaultValue?: DateType;
    value?: DateType | null;
    onSelect: (value: DateType) => void;
}

const Panel: React.FunctionComponent<PanelProps> = function Panel({
    open,
    value,
    showHour,
    showMinute,
    showSecond,
    hourStep = 1,
    minuteStep = 1,
    secondStep = 1,
    disabledHours,
    disabledMinutes,
    disabledSeconds,
    hideDisabledOptions,
    onSelect,
}) {
    const columns: {
        node: React.ReactElement;
    }[] = [];

    const hour = value ? generateConfig.getHour(value) : -1;
    const minute = value ? generateConfig.getMinute(value) : -1;
    const second = value ? generateConfig.getSecond(value) : -1;

    const hours = generateTimeUnits(0, 23, hourStep, disabledHours && disabledHours());

    const minutes = generateTimeUnits(0, 59, minuteStep, disabledMinutes && disabledMinutes(hour));

    const seconds = generateTimeUnits(0, 59, secondStep, disabledSeconds && disabledSeconds(hour, minute));

    const setTime = (newHour: number, newMinute: number, newSecond: number) => {
        let newDate = value || generateConfig.getNow();

        const mergedHour = Math.max(0, newHour);
        const mergedMinute = Math.max(0, newMinute);
        const mergedSecond = Math.max(0, newSecond);

        newDate = generateConfig.setSecond(newDate, mergedSecond);
        newDate = generateConfig.setMinute(newDate, mergedMinute);
        newDate = generateConfig.setHour(newDate, mergedHour);

        return newDate;
    };

    function addColumnNode(
        condition: boolean | undefined,
        node: React.ReactElement,
        columnValue: number,
        units: Unit[],
        onColumnSelect: (diff: number) => void,
    ) {
        if (condition !== false) {
            // eslint-disable-next-line
            columns.push({
                node: React.cloneElement(node, {
                    value: columnValue,
                    onSelect: onColumnSelect,
                    units,
                    hideDisabledOptions,
                    open,
                }),
            });
        }
    }

    // Hour
    addColumnNode(showHour, <TimeUnitColumn key="hour" />, hour, hours, (num) => {
        onSelect(setTime(num, minute, second));
    });

    // Minute
    addColumnNode(showMinute, <TimeUnitColumn key="minute" />, minute, minutes, (num) => {
        onSelect(setTime(hour, num, second));
    });

    // Second
    addColumnNode(showSecond, <TimeUnitColumn key="second" />, second, seconds, (num) => {
        onSelect(setTime(hour, minute, num));
    });

    return (
        <div className={styles['time-panel-content']}>
            {columns.map(({ node }) => node)}
        </div>
    );
};

export default Panel;
