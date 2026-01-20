import React from 'react';
import styles from './styles.desktop';

export default function HeadBar({ children }) {
    return (
        <div>
            <div className={styles['title']}>
                {
                    children
                }
            </div>
            <div className={styles['line']}></div>
        </div>
    )
}