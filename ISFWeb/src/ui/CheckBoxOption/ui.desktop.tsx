import React from 'react';
import classnames from 'classnames';
import CheckBox from '../CheckBox/ui.desktop';
import styles from './styles.desktop';

export default function CheckBoxOption({ role, className, children, disabled, ...props }: UI.CheckBoxOption.Props) {
    return (
        <label role={role} className={classnames(styles['container'], className)}>
            <CheckBox
                disabled={disabled}
                {...props}
            />
            <span className={classnames(styles['text'], { [styles['disabled']]: disabled })}>
                {
                    children
                }
            </span>
        </label>
    )
}
