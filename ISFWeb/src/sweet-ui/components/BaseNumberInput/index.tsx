import React from 'react';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import BaseInput, { ValueChangeEvent } from '../BaseInput';

interface BaseNumberInput {
    /**
     * 文本值
     */
    value: number;

    /**
     * 最小值
     */
    min?: number;

    /**
     * 最大值
     */
    max?: number;

    /**
     * className
     */
    className?: string;

    /**
     * 禁用
     */
    disabled?: boolean;

    /**
     * 输入值发生变化时触发
     */
    onValueChange?: (event: SweetUIEvent<number>) => void;

    /**
     * 输入值为空时触发
     */
    onEmptyValue?: (event: SweetUIEvent<string>) => void;
}

const BaseNumberInput: React.SFC<BaseNumberInput> = function BaseNumberInput({
    value,
    min = -Infinity,
    max = Infinity,
    className,
    disabled = false,
    onValueChange,
    onEmptyValue,
    ...restProps
}) {
    // 触发onValueChange
    const dispatchValueChangeEvent = createEventDispatcher(onValueChange);

    // 触发onEmptyValue
    const dispatchEmptyValueEvent = createEventDispatcher(onEmptyValue);

    const handleValueChange = (event: ValueChangeEvent) => {
        const { detail: value } = event;

        if (value === '') {
            dispatchEmptyValueEvent(value);
        } else {
            const valueNumber = Number(value);
            if (!isNaN(valueNumber) && valueNumber >= min && valueNumber <= max) {
                dispatchValueChangeEvent(valueNumber);
            }
        }
    };

    return <BaseInput type="text" onValueChange={handleValueChange} {...{ value, min, max, className, disabled }} {...restProps} />;
};

BaseNumberInput.defaultProps = {};

export default BaseNumberInput;
