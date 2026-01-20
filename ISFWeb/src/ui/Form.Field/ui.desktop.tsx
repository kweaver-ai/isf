import React from 'react';
import classnames from 'classnames'
import styles from './styles.desktop';

const FormField: UI.FormField.Component = function FormField({ children, role, className = '', isRequired = false }) {
    return (
        <div
            className={classnames(
                styles['field'],
                {
                    [styles['required']]: isRequired,
                },
                className,
            )}
            role={role}
        >
            {
                children
            }
        </div>
    )
}

export default FormField