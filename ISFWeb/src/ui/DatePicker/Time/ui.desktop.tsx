import React from 'react';
import classnames from 'classnames';
import { range } from 'lodash';
import LinkChip from '../../LinkChip/ui.desktop';
import UIIcon from '../../UIIcon/ui.desktop';
import styles from './styles.desktop';

export interface TimePanelProps {
    /**
     * 当前选择的date对象，包含指定的时间点
     */
    selectDate: Date;

    /**
     * 起始日期
     */
    start?: Date;

    /**
     * 结束日期
     */
    end?: Date;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 选中时间点时触发
     */
    onSelectTime: (date: Date) => void;

    /**
     * 点击时间面板标题触发
     */
    onClickTime: () => void;
}

/**
 * 指定网格列数
 */
const Col = 4;

/**
 * 时间段总数
 */
const TotalNumber = 25;

const TimePanel: React.FunctionComponent<TimePanelProps> = function TimePanel({
    selectDate,
    start,
    end,
    disabled,
    onSelectTime,
    onClickTime,
}) {
    return (
        <div className={styles['panel']}>
            <div className={styles['time-header']}>
                <UIIcon className={styles['time-icon']} size="16" code={'\uf022'} onClick={onClickTime} />
                <span className={styles['title']} onClick={onClickTime}>
                    {`${selectDate.getHours() < 10 ? `0${selectDate.getHours()}` : selectDate.getHours()}:${selectDate.getMinutes() < 10
                        ? `0${selectDate.getMinutes()}`
                        : selectDate.getMinutes()}`}
                </span>
            </div>
            <div className={styles['time-body']}>{renderPanel(disabled, selectDate, start, end, onSelectTime)}</div>
        </div>
    );
};

const renderPanel = (disabled, selectDate, start, end, onSelectTime) => {
    return range(Math.ceil(TotalNumber / Col)).map((row, index) => (
        <div key={index}>
            {range(Col).map((col, colindex) => {
                return (
                    <div
                        key={colindex}
                        className={classnames(
                            { [styles['cell']]: index * Col + colindex <= 24 },
                            {
                                [styles['selected']]:
                                    Math.abs(
                                        selectDate.getTime() - caculDate(selectDate, index * Col + colindex).getTime(),
                                    ) /
                                    1000 <
                                    60,
                            },
                            {
                                [styles['disabled']]:
                                    disabled ||
                                    (start &&
                                        caculDate(selectDate, index * Col + colindex).getTime() < start.getTime()) ||
                                    (end && caculDate(selectDate, index * Col + colindex).getTime() > end.getTime()),
                            },
                        )}
                    >
                        <LinkChip
                            key={index * Col + colindex}
                            className={styles['time']}
                            disabled={
                                disabled ||
                                (start && caculDate(selectDate, index * Col + colindex).getTime() < start.getTime()) ||
                                (end && caculDate(selectDate, index * Col + colindex).getTime() > end.getTime())
                            }
                            onClick={() => onSelectTime(caculDate(selectDate, index * Col + colindex))}
                        >
                            {index * Col + colindex > 24 ? null : index * Col + colindex === 24 ? (
                                '23:59'
                            ) : (
                                (index * Col + colindex) < 10 ? `0${index * Col + colindex}:00` : `${index * Col + colindex}:00`
                            )}
                        </LinkChip>
                    </div>
                );
            })}
        </div>
    ));
};

const caculDate = (date: Date, time: number) => {
    return new Date(
        date.getFullYear(),
        date.getMonth(),
        date.getDate(),
        time === 24 ? 23 : time,
        time === 24 ? 59 : 0,
        0,
    );
};

TimePanel.defaultProps = {
    disabled: false,
};

export default TimePanel;
