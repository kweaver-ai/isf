import React from 'react';
import styles from './styles.desktop';

const PanelMain: UI.PanelMain.Component = function PanelMain({ children, role }) {
    return (
        <div
            className={styles['container']}
            role={role}
        >
            <div className={styles['padding']}>
                {
                    children
                }
            </div>
        </div>
    )
}

export default PanelMain;