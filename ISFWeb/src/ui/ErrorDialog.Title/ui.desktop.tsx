import React from 'react';
import styles from './styles.desktop'

const ErrorDialogTitle: React.FunctionComponent<any> = function ErrorDialog({
    children,
}) {
    return (
        <div className={styles['title']}>
            {
                children
            }
        </div>
    )
}

export default ErrorDialogTitle