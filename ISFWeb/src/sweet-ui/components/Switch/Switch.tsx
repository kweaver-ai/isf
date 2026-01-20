import React from 'react';
import { isFunction } from 'lodash';
import classnames from 'classnames';
import styles from './styles';

export interface SwitchProps {
    /**
     *  是否禁用
     */
    disabled?: boolean;

    /**
     * 指定当前是否选中
     */
    checked?: boolean;

    /**
     * 样式class
     */
    className?: string;

    /**
     * role
     */
    role?: string;

    /**
     * 开关状态变化时的回调
     */
    onChange: (event: { event: React.MouseEvent; detail: boolean }) => void;
}

const Switch: React.FunctionComponent<SwitchProps> = function Switch({
    disabled = false,
    checked = false,
    onChange,
    className,
    role,
}) {
    const handleClickEvent = (event, checked: boolean) => {
        if (disabled) {
            return
        }

        isFunction(onChange) && onChange({ event, detail: !checked })
    };

    return (
        <button
            type="button"
            role={role}
            className={classnames(
                styles['switch'],
                {
                    [styles['switch-checked']]: checked,
                    [styles['switch-unchecked']]: !checked,
                    [styles['disabled']]: disabled,
                },
                className,
            )}
            onClick={(event) => handleClickEvent(event, checked)}
        />
    );
};

export default Switch;
