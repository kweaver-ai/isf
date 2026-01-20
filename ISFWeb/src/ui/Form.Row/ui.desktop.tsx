import React from 'react';
import classNames from 'classnames';
import styles from './styles.desktop';

const FormRow: UI.FormRow.Component = function FormRow({ className = '', children, role }) {
    return (
        <div
            className={classNames(styles['row'], className)}
            role={role}
        >
            {
                children
            }
        </div>
    )
}
export default FormRow