import React from 'react';
import classnames from 'classnames';
import { createEventDispatcher, SweetUIEvent } from '../../utils/event';
import View from '../View';
import styles from './styles';

interface RadioProps {
    /**
     * 单选框值
     */
    value: any;

    /**
     * HTML name
     */
    name?: string;

    /**
     * HTML id
     */
    id?: string;

    /**
     * 禁用状态
     */
    disabled?: boolean;

    /**
     * 选中状态
     */
    checked?: boolean;

    /**
     * 默认选中状态
     */
    defaultChecked?: boolean;

    /**
     * 选中状态发生变化时触发
     */
    onChange?: (event: SweetUIEvent<{ checked: boolean; value: any }>) => void;
}

const Radio: React.FunctionComponent<RadioProps> = ({
    role,
    children,
    id,
    name,
    value,
    checked = false,
    defaultChecked = false,
    disabled = false,
    onChange,
}) => {
    const handleChangeEvent = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (!disabled) {
            createEventDispatcher(onChange)({ event, value });
        }
    };

    return (
        <label role={role} className={classnames(styles['radio-wrapper'], { [styles['disabled']]: disabled })}>
            <View
                inline={true}
                className={classnames(styles['radio'], {
                    [styles['radio-checked']]: checked || defaultChecked,
                    [styles['disabled']]: disabled,
                })}
            >
                <input
                    type="radio"
                    className={classnames(styles['radio-input'], { [styles['disabled']]: disabled })}
                    {...{ id, name, value, disabled }}
                    checked={checked || defaultChecked}
                    onChange={handleChangeEvent}
                />
            </View>
            {children !== undefined ? (
                <View inline={true} className={classnames(styles['radio-text'], { [styles['disabled']]: disabled })}>
                    {children}
                </View>
            ) : null}
        </label>
    );
};

export default Radio;
