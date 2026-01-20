import React from 'react';
import styles from './styles.desktop'

const ErrorDialogDetail: React.FunctionComponent<any> = function ErrorDialog({
    children,
}) {
    return (
        <div className={styles['detail']}>
            {
                children
            }
        </div>
    )
}

export default ErrorDialogDetail