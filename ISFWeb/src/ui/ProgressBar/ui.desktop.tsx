import React from 'react';
import classnames from 'classnames';
import styles from './styles.desktop';

const ProgressBar: React.FunctionComponent<UI.ProgressBar.Props> = function ProgressBar({
    role,
    width = '100%',
    height,
    border,
    containerBackground,
    progressBackground,
    value,
    textAlign = 'right',
    renderValue = (value) => `${(value * 100).toFixed(0)}%`,
}) {
    return (
        <div
            role={role}
            className={styles['progressbar']}
            style={{
                width,
                height,
                border,
                backgroundColor: containerBackground,
            }}
        >
            <div className={classnames(
                styles['text'],
                {
                    [styles['left']]: textAlign === 'left',
                },
                {
                    [styles['right']]: textAlign === 'right',
                },
            )}>
                {
                    renderValue(value)
                }
            </div>
            <div style={{
                width: `${value * 100}%`,
                backgroundColor: progressBackground,
            }}
            className={styles['percentage']}
            >
            </div>
        </div>
    )
}

export default ProgressBar