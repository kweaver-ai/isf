import React from 'react';
import classnames from 'classnames';
import { isFunction } from 'lodash';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import styles from './styles';

interface CheckBoxProps {
    className?: string;

    id?: string;

    /**
     * 是否禁用
     */
    disabled?: boolean;

    /**
     * 是否选中
     */
    checked?: boolean;

    /**
     * 是否半选状态
     */
    indeterminate?: boolean;

    /**
     * 点击时触发的处理函数
     */
    onClick?: (event: React.MouseEvent<HTMLLabelElement>) => void;

    /**
     * 选中状态发生变化时触发
     */
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;

    /**
     * 选中状态发生变化时触发，传递选中状态作为参数
     */
    onCheckedChange?: (event: SweetUIEvent<boolean>) => void;
}

const CheckBox: React.FunctionComponent<CheckBoxProps> = function CheckBox({
    role,
    disabled,
    checked,
    className,
    id,
    onClick,
    onChange,
    onCheckedChange,
    indeterminate,
    children,
}) {
    const handleChangeEvent = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (!disabled) {
            createEventDispatcher(onCheckedChange)(event.target.checked);

            if (isFunction(onChange)) {
                onChange(event);
            }
        }
    };

    const handleClickEvent = (event: React.MouseEvent<HTMLLabelElement>) => {
        if (!disabled) {
            if (isFunction(onClick)) {
                onClick(event);
            } else {
                event.stopPropagation();
            }
        }
    };

    return (
        <label
            role={role}
            className={classnames(styles['checkbox-wrapper'], className, {
                [styles['disabled']]: disabled,
            })}
            onClick={handleClickEvent}
        >
            <View
                inline={true}
                className={classnames(styles['checkbox'], {
                    [styles['checkbox-checked']]: checked,
                    [styles['disabled']]: disabled,
                    [styles['checkbox-indeterminate']]: indeterminate,
                })}
            >
                <input
                    id={id}
                    type="checkbox"
                    className={classnames(styles['checkbox-input'], { [styles['disabled']]: disabled })}
                    onChange={handleChangeEvent}
                    {...{ disabled, checked }}
                />
            </View>
            {children !== undefined ? (
                <View inline={true} className={classnames(styles['checkbox-text'], { [styles['disabled']]: disabled })}>
                    {children}
                </View>
            ) : null}
        </label>
    );
};

export default CheckBox;
