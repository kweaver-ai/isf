import React from 'react'
import classnames from 'classnames'
import { isFunction } from 'lodash'
import { ClassName } from '@/ui/helper';
import styles from './styles.view.css'
import SweetIcon from '../SweetIcon/index'

export enum Status {
    Error = 'error',
    Process = 'process',
    Finish = 'finish',
    Wait = 'wait',
}

interface StepInfo {
    className?: string;
    disabled?: boolean;
    stepIndex?: number;
    stepNumber?: number;
    status?: Status;
    title?: React.ReactNode;
    subTitle?: React.ReactNode;
    description?: React.ReactNode;
    tailContent?: React.ReactNode;
    icon?: React.ReactNode;
    extraNode?: React.ReactNode;
}

interface StepsProps {
    items: StepInfo[];
    current?: number;
    size?: 'default' | 'small';
    direction?: 'horizontal' | 'vertical';
    onChange?: (current: number) => void;
}

const Steps: React.FC<StepsProps> = ({
    items = [],
    direction = 'horizontal',
    onChange,
    current = 0,
    size = 'default',
}) => {

    const renderStep = (item: StepInfo, index: number) => {
        const mergedItem = { ...item }
        const stepNumber = index

        if (!mergedItem.status) {
            if (stepNumber === current) {
                mergedItem.status = Status.Process
            }
            else if (stepNumber < current) {
                mergedItem.status = Status.Finish
            } else {
                mergedItem.status = Status.Wait
            }
        }

        if (!mergedItem.icon) {
            if (mergedItem.status === Status.Finish) {
                mergedItem.icon = (
                    <SweetIcon
                        name={'selected'}
                        size={size === 'small' ? 14 : 18}
                    />
                )
            } else if (mergedItem.status === Status.Error) {
                mergedItem.icon = (
                    <SweetIcon
                        name={'x'}
                        size={size === 'small' ? 14 : 18}
                    />
                )
            } else {
                mergedItem.icon = (
                    <span>{stepNumber + 1}</span>
                )
            }
        }

        const mergedStatus = mergedItem.status || Status.Wait

        const classString = classnames(
            styles['item'],
            styles[`item-${mergedStatus}`],
            {
                [styles['item-active']]: stepNumber === current,
                [styles['item-disabled']]: mergedItem.disabled === true,
                [styles['item-small']]: size === 'small',
            },
            mergedItem.className,
        )

        return (
            <div
                key={index}
                className={classString}
                {...(!!onChange && !mergedItem.disabled) ? { role: 'button' } : {}}
                onClick={() => !mergedItem.disabled && isFunction(onChange) && onChange(index)}
            >
                <div
                    className={styles['item-container']}
                >
                    <div className={styles['item-tail']}>{mergedItem.tailContent}</div>
                    <div
                        className={classnames(
                            styles['item-icon'],
                            styles[`item-icon-${mergedStatus}`],
                            {
                                [ClassName.BackgroundColor]: mergedStatus === 'process',
                            },
                        )}
                    >
                        {mergedItem.icon}
                    </div>
                    <div className={styles['item-content']}>
                        <div className={styles['item-title']}>
                            {mergedItem.title}
                            {mergedItem.subTitle && (
                                <div
                                    className={styles['item-subtitle']}
                                >
                                    {mergedItem.subTitle}
                                </div>
                            )}
                        </div>
                        {mergedItem.description && <div className={styles['item-description']}>{mergedItem.description}</div>}
                    </div>
                </div>
                {mergedItem.extraNode}
            </div>
        )
    }

    return (
        <div className={classnames(styles['steps'], styles[`steps-${direction}`])}>
            {items.filter((item) => item).map(renderStep)}
        </div>
    )
}

export default Steps