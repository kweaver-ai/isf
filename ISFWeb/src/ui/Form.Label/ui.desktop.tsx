import React from 'react';
import classnames from 'classnames';
import styles from './styles.desktop';

const FormLabel: UI.FormLabel.Component = function FormLabel({
    align,
    className,
    colon = false,
    children,
    ...props }) {
    return (
        <label
            className={classnames(
                styles['label'],
                {
                    [styles['align-top']]: align === 'top',
                    [styles['colon']]: colon,
                },
                className,
            )}
            {...props}
        >
            {
                children
            }
        </label>
    )
}

export default FormLabel